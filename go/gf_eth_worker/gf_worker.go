/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_eth_worker

import (
	"fmt"
	"time"
	"context"
	"strings"
	"github.com/getsentry/sentry-go"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_aws"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_core"
)

//-------------------------------------------------
// channel that receives channels that receive list of strings. 
// client sends a channel on which it is expecting to receive a response that is a list of strings.
type worker_inspector__get_hosts_ch chan chan []string
type Get_worker_hosts_fn            func(context.Context, *gf_eth_core.GF_runtime) []string

//-------------------------------------------------
func Discovery__init(p_runtime *gf_eth_core.GF_runtime) (func(context.Context, *gf_eth_core.GF_runtime) []string, chan *gf_core.GF_error) {

	// FIX!! - move this out of this function, into an ENV var or some kind of config
	worker_inspector_port_int := uint(9000)

	// CONFIG
	update_period_sec         := 2*60 * time.Second // 2min
	target_instances_tags_lst := [][]map[string]string{
		{{"Name": "gf_eth_monitor__worker__archive"}},
		// {{"Name": "gf_eth_monitor__worker__fast"}},
	}

	get_hosts_ch        := make(worker_inspector__get_hosts_ch, 100) // client requests for hosts are received on this channel
	new_instances_ch    := make(chan []*ec2.Instance, 10)            // new instances discovered are sent to this channel
	discovery_errors_ch := make(chan *gf_core.GF_error, 100)         // errors durring instance discovery are sent to this channel

	//-------------------------------------------------
	// MAIN
	go func() {

		var hosts_lst []string
		for {
			select {
			//-----------------------------
			// MSG__NEW_INSTANCES
			// when new instances are discovered and sent on the new_instances_ch
			// they're processed and added to the hosts_lst, for sending back to clients
			// requesting this info on get_hosts_ch

			case new_instances_lst := <- new_instances_ch:

				hosts_lst = []string{} // reset
				for _, ec2_inst := range new_instances_lst {



					inst__dns_name_str  := *ec2_inst.PublicDnsName
					inst__host_port_str := fmt.Sprintf("%s:%d", inst__dns_name_str, worker_inspector_port_int)

					// IMPORTANT!! - only register instances that have a public DNS name.
					//               they might not have that DNS name assigned, if they're 
					//               not reachable from outside their VPC, or if they're terminated...
					if inst__dns_name_str == "" {
						continue
					}

					hosts_lst = append(hosts_lst, inst__host_port_str)
				}
				
			//-----------------------------
			// MSG__GET_HOSTS - requests for hosts_lst made from clients of this goroutine are made
			//                  by sending the channel on which to send the response on.
			case get_hosts__reply_ch := <- get_hosts_ch:

				get_hosts__reply_ch <- hosts_lst

			//-----------------------------
			}
		}
	}()

	//-------------------------------------------------
	// DISCOVER - runs continuously
	go func() {
		for {

			aws_instances_all_lst := []*ec2.Instance{}
			for _, instance_tags_lst := range target_instances_tags_lst {
				aws_instances_lst, gf_err := gf_aws.AWS_EC2__describe_instances__by_tags(instance_tags_lst, p_runtime.Runtime_sys)
				if gf_err != nil {
					discovery_errors_ch <- gf_err

					// SLEEP - sleep after error as well
					// time.Sleep(update_period_sec)

					continue
				}

				aws_instances_all_lst = append(aws_instances_all_lst, aws_instances_lst...)
			} 


			

			new_instances_ch <- aws_instances_all_lst
			
			// SLEEP
			time.Sleep(update_period_sec)
		}
	}()
	
	//-------------------------------------------------
	// GET_WORKER_HOSTS__DYNAMIC_FN
	get_worker_hosts__dynamic_fn := func(p_ctx context.Context, p_runtime *gf_eth_core.GF_runtime) []string {


		fmt.Println("=============")
		fmt.Println(p_ctx)

		
		span__get_worker_hosts := sentry.StartSpan(p_ctx, "get_worker_hosts")
		span__get_worker_hosts.SetTag("workers_aws_discovery", fmt.Sprint(p_runtime.Config.Workers_aws_discovery_bool))

		var workers_inspectors_hosts_lst []string

		//---------------------
		// FROM_DISCOVERY
		if p_runtime.Config.Workers_aws_discovery_bool {

			reply_ch := make(chan []string)
			get_hosts_ch       <- reply_ch // send the channel (on which to receive the response) to the main hosts resolving goroutine 
			hosts_ports_lst := <- reply_ch // wait for a response from the resolving goroutine to the hosts request 

			workers_inspectors_hosts_lst = hosts_ports_lst

		//---------------------
		// FROM_CONFIG
		// REMOVE!? - is this used? 
		//            seems that get_worker_hosts__static_fn() gets used if p_runtime.Config.Workers_aws_discovery_bool == true
		//            so this "else" branch will never be used.
		} else {
			workers_inspectors_hosts_str := p_runtime.Config.Workers_hosts_str
			workers_inspectors_hosts_lst = strings.Split(workers_inspectors_hosts_str, ",")
		}

		//---------------------
		span__get_worker_hosts.Finish()

		return workers_inspectors_hosts_lst
	}

	//-------------------------------------------------
	// GET_WORKER_HOSTS__STATIC_FN
	get_worker_hosts__static_fn := func(p_ctx context.Context, p_runtime *gf_eth_core.GF_runtime) []string {
		worker_hosts_lst := strings.Split(p_runtime.Config.Workers_hosts_str, ",")
		return worker_hosts_lst
	}
	
	//-------------------------------------------------

	var get_worker_hosts_fn func(context.Context, *gf_eth_core.GF_runtime) []string
	if p_runtime.Config.Workers_aws_discovery_bool {
		get_worker_hosts_fn = get_worker_hosts__dynamic_fn
	} else {
		get_worker_hosts_fn = get_worker_hosts__static_fn
	}

	return get_worker_hosts_fn, discovery_errors_ch
}