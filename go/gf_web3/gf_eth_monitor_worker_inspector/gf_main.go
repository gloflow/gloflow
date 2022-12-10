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
	"time"
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
)

//-------------------------------------------------
type GF_runtime struct {
	eth_rpc_host_str string
	eth_rpc_client   *ethclient.Client
	py_plugins       *gf_eth_core.GF_py_plugins
	runtimeSys       *gf_core.RuntimeSys
}

//-------------------------------------------------
func main() {
	
	// log to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// log.SetLevel(log.WarnLevel)

	port_int         := 9000
	port_metrics_int := 9120
	geth__port_int   := 8545
	geth__host_str               := os.Getenv("GF_GETH_HOST") // "127.0.0.1"
	sentry_endpoint_str          := os.Getenv("GF_SENTRY_ENDPOINT")
	py_plugins_base_dir_path_str := os.Getenv("GF_PY_PLUGINS_BASE_DIR_PATH")

	logFun, _ := gf_core.LogsInit()
	runtime, err := runtimeGet(geth__host_str, py_plugins_base_dir_path_str, logFun)
	if err != nil {
		panic(err)
	}

	//-------------
	// SENTRY
	
	sentry_samplerate_f := 1.0
	sentry_transaction_to_trace_map := map[string]bool{
		"GET /gfethm_worker_inspect/v1/account/info": true,
		"GET /gfethm_worker_inspect/v1/tx/trace":     true,
		"GET /gfethm_worker_inspect/v1/blocks":       true,
	}

	err = gf_core.Error__init_sentry(sentry_endpoint_str,
		sentry_transaction_to_trace_map,
		sentry_samplerate_f)
	if err != nil {
		panic(err)
	}

	defer sentry.Flush(2 * time.Second)

	//-------------
	// METRICS
	metrics, gf_err := metrics__init(port_metrics_int)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
	// ETH_CLIENT
	eth_client, gf_err := gf_eth_core.Eth_rpc__init(geth__host_str,
		geth__port_int,
		runtime.runtimeSys)
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
func runtimeGet(p_eth_rpc_host_str string,
	p_py_plugins_base_dir_path_str string,
	pLogFun                        func(string, string)) (*GF_runtime, error) {

	//--------------------
	// RUNTIME_SYS
	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_eth_monitor_worker_inspector",
		LogFun:           pLogFun,
		
		// SENTRY - enable it for error reporting
		ErrorsSendToSentryBool: true,
	}

	//--------------------
	// PY_PLUGINS
	py_plugins := &gf_eth_core.GF_py_plugins {
		Base_dir_path_str: p_py_plugins_base_dir_path_str,
	}

	//--------------------
	// RUNTIME
	runtime := &GF_runtime{
		eth_rpc_host_str: p_eth_rpc_host_str,
		runtimeSys:       runtimeSys,
		py_plugins:       py_plugins,
	}

	//--------------------
	return runtime, nil
}