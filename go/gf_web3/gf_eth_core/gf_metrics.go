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

package gf_eth_core

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
	
	// SQS
	SQS__msgs_num__counter prometheus.Counter

	// PEERS
	Peers__http_req_num__get_peers__counter prometheus.Counter
	Peers__unique_names_num__gauge          prometheus.Gauge

	// BLOCK
	Block__db_count__gauge       prometheus.Gauge
	Block__indexed_num__counter prometheus.Counter

	// TX
	Tx__db_count__gauge      prometheus.Gauge
	Tx__indexed_num__counter prometheus.Counter

	// TX_TRACE
	Tx_trace__worker_inspector_durration__gauge prometheus.Gauge
	Tx_trace__py_plugin__plot_durration__gauge  prometheus.Gauge
	Tx_trace__db_count__gauge                   prometheus.Gauge

	// DB
	DB__writes_num__new_peer_lifecycle__counter  prometheus.Counter

	// ERRORS
	Errs_num__counter prometheus.Counter
}

//-------------------------------------------------
// INIT
func Metrics__init(p_port_int int) (*GF_metrics, *gf_core.GFerror) {


	//---------------------------
	// SQS_MSGS_NUM
	sqs_msgs_num__counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gf_eth_mon",
		Name: "sqs_msgs_num__count",
		Help: "number of AWS SQS messages received",
	})

	//---------------------------
	// PEERS

	// HTTP_REQ_NUM__GET_PEERS
	http_req_num__get_peers__counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gf_eth_mon",
		Name: "http_req_num__get_peers",
		Help: "number of HTTP requests received to get peers",
	})
	
	// PEERS_UNIQUE_NAMES_NUM
	peers_unique_names_num__gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gf_eth_mon",
			Name: "peers__unique_names_num__count",
			Help: "number of unique peer names",
		})
	
	//---------------------------
	// BLOCK
	block__db_count__gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gf_eth_mon",
			Name: "block__db_count__gauge",
			Help: "how many block records are in the DB",
		})

	
	// INDEXED_BLOCKS
	block__indexed_num__counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gf_eth_mon",
		Name: "block__indexed__num__count",
		Help: "number of blocks that were indexed",
	})

	//---------------------------
	// TX
	tx__db_count__gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gf_eth_mon",
			Name: "tx__db_count__gauge",
			Help: "how many tx records are in the DB",
		})
	
	// INDEXED_TXS
	tx__indexed_num__count := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gf_eth_mon",
		Name: "tx__indexed_num__count",
		Help: "number of tx's that were indexed",
	})

	//---------------------------
	// TX_TRACE
	
	tx_trace__worker_inspector_durration__gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gf_eth_mon",
			Name:      "tx_trace__worker_inspector_durration__gauge",
			Help:      "how long the tx_trace request takes to the worker_inspector",
		})

	tx_trace__py_plugin__plot_durration__gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gf_eth_mon",
			Name:      "tx_trace__py_plugin__plot_durration__gauge",
			Help:      "how long the tx_trace py_plugin plot execution takes",
		})

	tx_trace__db_count__gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "gf_eth_mon",
			Name:      "tx_trace__db_count__gauge",
			Help:      "how many tx_trace records are in the DB",
		})
	
	//---------------------------
	// DB_WRITES_NUM__NEW_PEER_LIFECYCLE
	db_writes_num__new_peer_lifecycle__counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gf_eth_mon",
		Name: "db_writes_num__new_peer_lifecycle__count",
		Help: "number of DB write operations for the new_peer_lifecycle",
	})

	//---------------------------
	// ERRS_NUM
	errs_num__counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gf_eth_mon",
		Name: "errs_num__count",
		Help: "number of errors",
	})

	//---------------------------
	prometheus.MustRegister(sqs_msgs_num__counter)
	prometheus.MustRegister(http_req_num__get_peers__counter)
	prometheus.MustRegister(peers_unique_names_num__gauge)
	prometheus.MustRegister(block__db_count__gauge)
	prometheus.MustRegister(block__indexed_num__counter)
	prometheus.MustRegister(tx__db_count__gauge)
	prometheus.MustRegister(tx__indexed_num__count)
	prometheus.MustRegister(tx_trace__worker_inspector_durration__gauge)
	prometheus.MustRegister(tx_trace__py_plugin__plot_durration__gauge)
	prometheus.MustRegister(tx_trace__db_count__gauge)
	prometheus.MustRegister(db_writes_num__new_peer_lifecycle__counter)
	prometheus.MustRegister(errs_num__counter)


	metrics := &GF_metrics{
		SQS__msgs_num__counter:                      sqs_msgs_num__counter,

		// PEERS
		Peers__http_req_num__get_peers__counter:     http_req_num__get_peers__counter,
		Peers__unique_names_num__gauge:              peers_unique_names_num__gauge,

		// BLOCK
		Block__db_count__gauge:                      block__db_count__gauge,
		Block__indexed_num__counter:                 block__indexed_num__counter,

		// TX
		Tx__db_count__gauge:                         tx__db_count__gauge,
		Tx__indexed_num__counter:                    tx__indexed_num__count,

		// TX_TRACE
		Tx_trace__worker_inspector_durration__gauge: tx_trace__worker_inspector_durration__gauge,
		Tx_trace__py_plugin__plot_durration__gauge:  tx_trace__py_plugin__plot_durration__gauge,
		Tx_trace__db_count__gauge:                   tx_trace__db_count__gauge,

		DB__writes_num__new_peer_lifecycle__counter: db_writes_num__new_peer_lifecycle__counter,
		Errs_num__counter:                           errs_num__counter,
	}

	return metrics, nil
}

//-------------------------------------------------
// INIT_SERVER
func Metrics__init_server(p_port_int int) {
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
}