/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_eth_monitor_lib

import (
	"fmt"
	// "time"
	"context"
	"net/http"
	log "github.com/sirupsen/logrus"
	// "github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	// "github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_blocks"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_tx"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_indexer"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_eth_worker"
)

//-------------------------------------------------
func InitService() {






	
}

//-------------------------------------------------
func Run_service(p_runtime *gf_eth_core.GF_runtime) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_eth_monitor_service.Run_service()")

	//-------------
	/*// SENTRY
	sentry_endpoint_str := p_runtime.Config.Sentry_endpoint_str
	sentry_samplerate_f := 1.0
	sentry_trace_handlers_map := map[string]bool{
		"GET /gfethm/v1/block/index":   true,
		"GET /gfethm/v1/tx/trace/plot": true,
		"GET /gfethm/v1/block":         true,
		"GET /gfethm/v1/miner": true,
		// "/gfethm/v1/block": true,
		"GET /gfethm/v1/peers": true,
		// "http__master__get_block": true,
		// "http__master__get_block": true,
	}
	err := gf_core.Error__init_sentry(sentry_endpoint_str,
		sentry_trace_handlers_map,
		sentry_samplerate_f)
	if err != nil {
		panic(err)
	}

	defer sentry.Flush(2 * time.Second)*/

	sentry_endpoint_uri_str := p_runtime.Config.Sentry_endpoint_str
	gf_eth_core.Sentry__init(sentry_endpoint_uri_str)

	//-------------
	// METRICS
	port_metrics_int := 9110

	metrics, gf_err := gf_eth_core.Metrics__init(port_metrics_int)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	
	gf_eth_blocks.Init_continuous_metrics(metrics, p_runtime)
	gf_eth_tx.Init_continuous_metrics(metrics, p_runtime)

	// causing errors
	// gf_eth_core.Eth_tx_trace__init_continuous_metrics(metrics, p_runtime)
	
	gf_eth_core.Eth_peers__init_continuous_metrics(metrics, p_runtime)

	//-------------
	// QUEUE
	if p_runtime.Config.Events_consume_bool {
		queue_name_str  := p_runtime.Config.AWS_SQS_queue_str
		queue_info, err := Event__init_queue(queue_name_str, metrics)
		if err != nil {
			fmt.Println("failed to initialize event queue")
			panic(err)
		}

		// QUEUE_START_CONSUMING

		ctx := context.Background()
		event__start_sqs_consumer(queue_info, ctx, metrics, p_runtime)
	}

	//-------------
	// WORKER_DISCOVERY
	get_hosts_fn, _ := gf_eth_worker.Discovery__init(p_runtime)
	
	//-------------
	// INDEXER

	indexer_cmds_ch, indexer_job_updates_new_consumer_ch, gf_err := gf_eth_indexer.Init(get_hosts_fn, metrics, p_runtime)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
	// HANDLERS
	gf_err = InitHandlers(get_hosts_fn,
		indexer_cmds_ch,
		indexer_job_updates_new_consumer_ch,
		metrics,
		p_runtime)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
	// METRICS_SERVER
	gf_eth_core.Metrics__init_server(port_metrics_int)

	//-------------
	port_str := p_runtime.Config.Port_str

	// p_runtime.Runtime_sys.Log_fun("INFO", fmt.Sprintf("STARTING HTTP SERVER - PORT - %s", port_str))
	log.WithFields(log.Fields{"port": port_str,}).Info("STARTING HTTP SERVER")

	sentry_handler := sentryhttp.New(sentryhttp.Options{}).Handle(http.DefaultServeMux)
	http_err       := http.ListenAndServe(fmt.Sprintf(":%s", port_str), sentry_handler)

	if http_err != nil {
		msg_str := fmt.Sprintf("cant start listening on port - ", port_str)
		p_runtime.Runtime_sys.Log_fun("ERROR", msg_str)
		p_runtime.Runtime_sys.Log_fun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}