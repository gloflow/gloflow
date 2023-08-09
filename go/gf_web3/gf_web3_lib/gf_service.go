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

package gf_web3_lib

import (
	"fmt"
	"context"
	"net/http"
	log "github.com/sirupsen/logrus"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_blocks"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_tx"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_indexer"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_worker"
	"github.com/gloflow/gloflow/go/gf_web3/gf_address"
	"github.com/gloflow/gloflow/go/gf_web3/gf_nft"
)

//-------------------------------------------------

func InitService(pAuthSubsystemTypeStr string,
	pKeyServer        *gf_identity_core.GFkeyServerInfo,
	pHTTPmux          *http.ServeMux,
	pConfig           *gf_eth_core.GF_config,
	pImagesJobsMngrCh chan gf_images_jobs_core.JobMsg,
	pRuntimeSys       *gf_core.RuntimeSys) {

	//-------------
	// ADDRESS
	gfErr := gf_address.InitHandlers(pAuthSubsystemTypeStr,
		pKeyServer,
		pHTTPmux,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//-------------
	// NFT
	gfErr = gf_nft.InitHandlers(pKeyServer,
		pHTTPmux,
		pConfig,
		pImagesJobsMngrCh,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//-------------
}

//-------------------------------------------------

func RunService(pRuntime *gf_eth_core.GF_runtime) {

	//-------------
	/*// SENTRY
	sentry_endpoint_str := pRuntime.Config.Sentry_endpoint_str
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

	sentry_endpoint_uri_str := pRuntime.Config.Sentry_endpoint_str
	gf_eth_core.SentryInit(sentry_endpoint_uri_str)

	//-------------
	// METRICS
	port_metrics_int := 9110

	metrics, gfErr := gf_eth_core.Metrics__init(port_metrics_int)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	
	gf_eth_blocks.Init_continuous_metrics(metrics, pRuntime)
	gf_eth_tx.Init_continuous_metrics(metrics, pRuntime)

	// causing errors
	// gf_eth_core.Eth_tx_trace__init_continuous_metrics(metrics, pRuntime)
	
	gf_eth_core.Eth_peers__init_continuous_metrics(metrics, pRuntime)

	//-------------
	// QUEUE
	if pRuntime.Config.Events_consume_bool {
		queue_name_str  := pRuntime.Config.AWS_SQS_queue_str
		queue_info, err := Event__init_queue(queue_name_str, metrics)
		if err != nil {
			fmt.Println("failed to initialize event queue")
			panic(err)
		}

		// QUEUE_START_CONSUMING

		ctx := context.Background()
		eventStartSQSconsumer(queue_info, ctx, metrics, pRuntime)
	}

	//-------------
	// WORKER_DISCOVERY
	get_hosts_fn, _ := gf_eth_worker.Discovery__init(pRuntime)
	
	//-------------
	// INDEXER

	indexer_cmds_ch, indexer_job_updates_new_consumer_ch, gfErr := gf_eth_indexer.Init(get_hosts_fn, metrics, pRuntime)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//-------------
	// HANDLERS
	gfErr = InitHandlers(get_hosts_fn,
		indexer_cmds_ch,
		indexer_job_updates_new_consumer_ch,
		metrics,
		pRuntime)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//-------------
	// METRICS_SERVER
	gf_eth_core.Metrics__init_server(port_metrics_int)

	//-------------
	port_str := pRuntime.Config.Port_str

	// pRuntime.runtimeSys.LogFun("INFO", fmt.Sprintf("STARTING HTTP SERVER - PORT - %s", port_str))
	log.WithFields(log.Fields{"port": port_str,}).Info("STARTING HTTP SERVER")

	sentry_handler := sentryhttp.New(sentryhttp.Options{}).Handle(http.DefaultServeMux)
	http_err       := http.ListenAndServe(fmt.Sprintf(":%s", port_str), sentry_handler)

	if http_err != nil {
		msg_str := fmt.Sprintf("cant start listening on port - %s", port_str)
		pRuntime.RuntimeSys.LogFun("ERROR", msg_str)
		pRuntime.RuntimeSys.LogFun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}