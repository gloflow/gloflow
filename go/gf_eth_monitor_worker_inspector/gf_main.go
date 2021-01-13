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

package main

import (
	"os"
	"fmt"
	"net/http"
	log "github.com/sirupsen/logrus"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_lib"
)

//-------------------------------------------------
type GF_runtime struct {
	eth_rpc_client *ethclient.Client
	runtime_sys    *gf_core.Runtime_sys
}

//-------------------------------------------------
func main() {
	
	// log to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// log.SetLevel(log.WarnLevel)

	port_int         := 2000
	port_metrics_int := 9120
	worker_inspector__geth__host_str := "127.0.0.1"
	sentry_endpoint_str              := os.Getenv("GF_SENTRY_ENDPOINT")

	log_fun := gf_core.Init_log_fun()
	runtime, err := runtime__get(log_fun)
	if err != nil {
		panic(err)
	}

	
	//-------------
	// SENTRY
	
	sentry_samplerate_f := 1.0
	sentry_trace_handlers_map := map[string]bool{
		"/gfethm_worker_inspect/v1/blocks": true,
	}

	err = gf_core.Error__init_sentry(sentry_endpoint_str,
		sentry_trace_handlers_map,
		sentry_samplerate_f)
	if err != nil {
		panic(err)
	}

	//-------------
	// METRICS
	metrics, gf_err := metrics__init(port_metrics_int)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
	// ETH_CLIENT
	eth_client, gf_err := gf_eth_monitor_lib.Eth_rpc__init(worker_inspector__geth__host_str, runtime.runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	runtime.eth_rpc_client = eth_client

	//-------------
	// HANDLERS
	init_handlers(metrics, runtime)

	//-------------
	
	log.WithFields(log.Fields{"port": port_int,}).Info("STARTING HTTP SERVER")

	sentry_handler := sentryhttp.New(sentryhttp.Options{}).Handle(http.DefaultServeMux)
	http_err := http.ListenAndServe(fmt.Sprintf(":%d", port_int), sentry_handler)

	if http_err != nil {
		log.WithFields(log.Fields{"port": port_int, "err": http_err}).Fatal("cant start HTTP listening on port")
		panic(http_err)
	}
}

//-------------------------------------------------
func runtime__get(p_log_fun func(string, string)) (*GF_runtime, error) {

	//--------------------
	// RUNTIME_SYS
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_eth_monitor",
		Log_fun:          p_log_fun,
		
		// SENTRY - enable it for error reporting
		Errors_send_to_sentry_bool: true,
	}

	//--------------------
	// RUNTIME
	runtime := &GF_runtime{
		runtime_sys: runtime_sys,
	}

	//--------------------
	return runtime, nil
}