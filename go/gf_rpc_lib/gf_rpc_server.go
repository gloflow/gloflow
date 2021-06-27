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
	"net/http"
	log "github.com/sirupsen/logrus"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

//-------------------------------------------------
func Server__init(p_port_int int) {

	log.WithFields(log.Fields{"port": p_port_int,}).Info("STARTING HTTP SERVER >>>>>>>>>>>")

	sentry_handler := sentryhttp.New(sentryhttp.Options{}).Handle(http.DefaultServeMux)
	err            := http.ListenAndServe(fmt.Sprintf(":%d", p_port_int), sentry_handler)

	if err != nil {
		log.WithFields(log.Fields{
			"port": p_port_int,
			"err":  err,
		}).Fatal("server cant start listening")
		panic(-1)
	}
}