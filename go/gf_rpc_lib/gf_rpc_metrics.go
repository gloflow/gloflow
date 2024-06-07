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
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

//-------------------------------------------------

type GFmetrics struct {

	// HANDLERS
	HandlersCountersMap map[string]prometheus.Counter

	// AUTH_SESSION_INVALID - counter for when auth session validation fails when
	//                        request is received and validated
	HandlersAuthSessionInvalidCounter prometheus.Counter

	// AUTH_SESSION_CORS
	HandlersAuthSessionCORScounter prometheus.Counter
}

// FIX!! - make sure this passed to all services, currently only identity is using it.
type GFglobalMetrics struct {

	// AUTH
	HandlersAuthCounter prometheus.Counter
}

//-------------------------------------------------
// CREATE_GLOBAL

func MetricsCreateGlobal() *GFglobalMetrics {

	handlersAuthCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_rpc__handler_auth"), 
		Help: "number of auth requests received",
	})
	prometheus.MustRegister(handlersAuthCounter)

	metrics := &GFglobalMetrics{
		HandlersAuthCounter: handlersAuthCounter,
	}

	return metrics
}

//-------------------------------------------------
// CREATE_FOR_HANDLERS

func MetricsCreateForHandlers(pMetricsGroupNameStr string,
	pServiceNameStr       string,
	pHandlersEndpointsLst []string) *GFmetrics {
	

	handlersCountersMap := map[string]prometheus.Counter{}
	for _, handlerEndpointStr := range pHandlersEndpointsLst {

		handlerEndpointCleanStr := strings.ReplaceAll(handlerEndpointStr, "/", "_")
		nameStr := fmt.Sprintf("gf_rpc__handler_reqs_num__%s_%s", pServiceNameStr, handlerEndpointCleanStr)
		
		handlerReqsNumCounter := prometheus.NewCounter(prometheus.CounterOpts{
			Name: nameStr,
			Help: "handler number of requests",
		})
		prometheus.MustRegister(handlerReqsNumCounter)


		handlersCountersMap[handlerEndpointStr] = handlerReqsNumCounter
	}

	handlersAuthSessionInvalidCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_rpc__handler_auth_session_invalid_num__%s_%s", pServiceNameStr, pMetricsGroupNameStr), 
		Help: "number of invalid auth session requests received",
	})
	prometheus.MustRegister(handlersAuthSessionInvalidCounter)

	handlersAuthSessionCORScounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_rpc__handler_auth_session_cors_num__%s_%s", pServiceNameStr, pMetricsGroupNameStr), 
		Help: "number of CORS (cross-domain) auth session requests received",
	})
	prometheus.MustRegister(handlersAuthSessionCORScounter)

	metrics := &GFmetrics{
		HandlersCountersMap:               handlersCountersMap,
		HandlersAuthSessionInvalidCounter: handlersAuthSessionInvalidCounter,
		HandlersAuthSessionCORScounter:    handlersAuthSessionCORScounter,
	}

	return metrics
}