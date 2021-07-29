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

package gf_eth_indexer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_metrics struct {
	SQS__msgs_num__counter prometheus.Counter
}

//-------------------------------------------------
// INIT
func Metrics__init() (*GF_metrics, *gf_core.GF_error) {


	//---------------------------
	// SQS_MSGS_NUM
	counter__sqs_msgs_num := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gf_eth_monitor__sqs_msgs_num",
		Help: "number of AWS SQS messages received",
	})



	prometheus.MustRegister(counter__sqs_msgs_num)
	return nil, nil
}