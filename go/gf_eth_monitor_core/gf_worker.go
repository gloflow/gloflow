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

package gf_eth_monitor_core

import (
	"time"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_aws"
)

//-------------------------------------------------
// channel that receives channels that receive list of strings. 
// client sends a channel on which it is expecting to receive a response that is a list of strings.
type worker_inspector__get_hosts_ch chan chan []string

//-------------------------------------------------
func Worker__discovery__init(p_runtime_sys *gf_core.Runtime_sys) (func() []string, chan *gf_core.Gf_error) {

	// CONFIG
	update_period_sec     := 2*60 * time.Second // 2min
	target_instances_tags_lst := []map[string]string{
		{"Name": "gf_eth_monitor__worker__archive"},
		{"Name": "gf_eth_monitor__worker__fast"},
	}

	get_hosts_ch        := make(worker_inspector__get_hosts_ch, 100) // client requests for hosts are received on this channel
	new_instances_ch    := make(chan []*ec2.Instance, 10)            // new instances discovered are sent to this channel
	discovery_errors_ch := make(chan *gf_core.Gf_error, 100)         // errors durring instance discovery are sent to this channel

	//-------------------------------------------------
	// MAIN
	go func() {

		var hosts_lst []string
		for {
			select {
			//-----------------------------
			// MSG__NEW_INSTANCES
			case new_instances_lst := <- new_instances_ch:

				hosts_lst = []string{} // reset
				for _, inst := range new_instances_lst {



					inst__dns_name_str := *inst.PublicDnsName
					hosts_lst = append(hosts_lst, inst__dns_name_str)
				}
				
			//-----------------------------
			// MSG__GET_HOSTS
			case get_hosts__reply_ch := <- get_hosts_ch:

				get_hosts__reply_ch <- hosts_lst

			//-----------------------------
			}
		}
	}()

	//-------------------------------------------------
	// DISCOVER
	go func() {
		for {

			
			aws_instances_lst, gf_err := gf_aws.AWS_EC2__describe_instances__by_tags(target_instances_tags_lst, p_runtime_sys)
			if gf_err != nil {
				discovery_errors_ch <- gf_err

				// SLEEP - sleep after error as well
				time.Sleep(update_period_sec)
				continue
			}

			new_instances_ch <- aws_instances_lst
			
			// SLEEP
			time.Sleep(update_period_sec)
		}
	}()
	
	//-------------------------------------------------

	get_hosts_fn := func() []string {

		reply_ch := make(chan []string)
		get_hosts_ch <- reply_ch
		hosts_lst := <- reply_ch
		return hosts_lst
	}
	return get_hosts_fn, discovery_errors_ch
}