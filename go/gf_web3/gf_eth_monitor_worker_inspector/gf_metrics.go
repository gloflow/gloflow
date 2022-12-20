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
	"fmt"
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type GF_metrics struct {
	counter__http_req_num__get_blocks prometheus.Counter
}

//-------------------------------------------------

func metrics__init(p_port_int int) (*GF_metrics, *gf_core.GFerror) {


	//---------------------------
	// SQS_MSGS_NUM
	counter__http_req_num__get_blocks := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gf_eth_monitor_worker_inspector__http_req_num__get_blocks",
		Help: "number of HTTP requests received to get blocks",
	})
	
	//---------------------------

	prometheus.MustRegister(counter__http_req_num__get_blocks)


	metrics_router := mux.NewRouter()
	metrics_router.Handle("/metrics", promhttp.Handler())

	metrics_server := http.Server{
		Handler: metrics_router,
		Addr:    fmt.Sprintf(":%d", p_port_int),
	}

	go func() {

		log.WithFields(log.Fields{"port": p_port_int,}).Info("STARTING METRICS HTTP SERVER")

		// ADD!! - check for returned error here,
		//         and report this in some way to the user.
		metrics_server.ListenAndServe()
	}()






	metrics := &GF_metrics{
		counter__http_req_num__get_blocks: counter__http_req_num__get_blocks,
	}

	return metrics, nil
}