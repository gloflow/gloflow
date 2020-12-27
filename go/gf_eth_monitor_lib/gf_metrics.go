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
	"net/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_metrics struct {
	
	counter__sqs_msgs_num                      prometheus.Counter
	counter__http_req_num__get_peers           prometheus.Counter
	counter__db_writes_num__new_peer_lifecycle prometheus.Counter
	counter__errs_num                          prometheus.Counter
}

//-------------------------------------------------
func metrics__init(p_port_int int) (*GF_metrics, *gf_core.Gf_error) {


	//---------------------------
	// SQS_MSGS_NUM
	counter__sqs_msgs_num := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gf_eth_monitor__sqs_msgs_num",
		Help: "number of session leave SDS HTTP requests received by the MMS",
	})

	// HTTP_REQ_NUM__GET_PEERS
	counter__http_req_num__get_peers := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gf_eth_monitor__http_req_num__get_peers",
		Help: "number of HTTP requests received to get peers",
	})

	// HTTP_REQ_NUM__GET_PEERS
	counter__db_writes_num__new_peer_lifecycle := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gf_eth_monitor__db_writes_num__new_peer_lifecycle",
		Help: "number of DB write operations for the new_peer_lifecycle",
	})

	// ERRS_NUM
	counter__errs_num := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gf_eth_monitor__errs_num",
		Help: "number of errors in ",
	})

	//---------------------------
	prometheus.MustRegister(counter__sqs_msgs_num)
	prometheus.MustRegister(counter__http_req_num__get_peers)
	prometheus.MustRegister(counter__db_writes_num__new_peer_lifecycle)
	prometheus.MustRegister(counter__errs_num)




	metrics_router := mux.NewRouter()
	metrics_router.Handle("/metrics", promhttp.Handler())

	metrics_server := http.Server{
		Handler: metrics_router,
		Addr:    fmt.Sprintf(":%d", p_port_int),
	}

	go func() {
		// ADD!! - check for returned error here,
		//         and report this in some way to the user.
		metrics_server.ListenAndServe()
	}()






	metrics := &GF_metrics{
		counter__sqs_msgs_num:                      counter__sqs_msgs_num,
		counter__http_req_num__get_peers:           counter__http_req_num__get_peers,
		counter__db_writes_num__new_peer_lifecycle: counter__db_writes_num__new_peer_lifecycle,
		counter__errs_num:                          counter__errs_num,
	}

	return metrics, nil
}