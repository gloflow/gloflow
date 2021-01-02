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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_runtime struct {
	Eth_rpc_client *ethclient.Client
	Runtime_sys    *gf_core.Runtime_sys
}

//-------------------------------------------------
func main() {
	
	// log to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// log.SetLevel(log.WarnLevel)

	port_int         := 2000
	port_metrics_int := 9120
	worker_inspector__geth__host_str := "127.0.0.1"


	log_fun := gf_core.Init_log_fun()
	runtime, err := runtime__get(log_fun)
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
	eth_client := eth_rpc__init(worker_inspector__geth__host_str)
	runtime.Eth_rpc_client = eth_client

	//-------------
	// HANDLERS
	init_handlers(metrics, runtime)

	//-------------
	
	log.WithFields(log.Fields{"port": port_int,}).Info("STARTING HTTP SERVER")
	http_err := http.ListenAndServe(fmt.Sprintf(":%d", port_int), nil)

	if http_err != nil {
		log.WithFields(log.Fields{"port": port_int, "err": http_err}).Fatal("cant start HTTP listening on port")
		panic(http_err)
	}
}

//-------------------------------------------------
func runtime__get(p_log_fun func(string, string)) (*GF_runtime, error) {

	// RUNTIME_SYS
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_eth_monitor",
		Log_fun:          p_log_fun,	
	}

	//--------------------
	// RUNTIME
	runtime := &GF_runtime{
		Runtime_sys: runtime_sys,
	}

	return runtime, nil
}