/*
MIT License

Copyright (c) 2021 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_rpc_lib

import (
	"fmt"
	
	"github.com/prometheus/client_golang/prometheus"
	
)

//-------------------------------------------------
type GF_metrics struct {
	Handlers_counters_map map[string]prometheus.Counter
}

//-------------------------------------------------
// CREATE_FOR_HANDLER
func Metrics__create_for_handlers(p_handlers_endpoints_lst []string) *GF_metrics {
	

	handlers_counters_map := map[string]prometheus.Counter{}
	for _, handler_endpoint_str := range p_handlers_endpoints_lst {
		name_str                   := fmt.Sprintf("gf_rpc__handler_reqs_num__%s", handler_endpoint_str)
		handler__reqs_num__counter := prometheus.NewCounter(prometheus.CounterOpts{
			Name: name_str,
			Help: "handler number of requests",
		})

		prometheus.MustRegister(handler__reqs_num__counter)


		handlers_counters_map[handler_endpoint_str] = handler__reqs_num__counter
	}

	metrics := &GF_metrics{
		Handlers_counters_map: handlers_counters_map,
	}

	return metrics
}