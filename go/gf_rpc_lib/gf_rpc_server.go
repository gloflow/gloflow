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

func ServerInitWithMux(pPortInt int,
	pMux *http.ServeMux) {

	log.WithFields(log.Fields{"port": pPortInt,}).Info("STARTING HTTP SERVER >>>>>>>>>>>")
	
	// IMPORTANT!! - wrap mux with Sentry.
	//               without this handler contexts are not initialized properly
	//               and creating spans and other sentry primitives fails with nul pointer exceptions.
	sentryHandler := sentryhttp.New(sentryhttp.Options{}).Handle(pMux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", pPortInt),
		Handler: sentryHandler,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.WithFields(log.Fields{
			"port": pPortInt,
			"err":  err,
		}).Fatal("server cant start listening")
		panic(-1)
	}
}

//-------------------------------------------------

func ServerInit(pPortInt int) {

	log.WithFields(log.Fields{"port": pPortInt,}).Info("STARTING HTTP SERVER >>>>>>>>>>>")

	sentryHandler := sentryhttp.New(sentryhttp.Options{}).Handle(http.DefaultServeMux)
	err           := http.ListenAndServe(fmt.Sprintf(":%d", pPortInt), sentryHandler)

	if err != nil {
		log.WithFields(log.Fields{
			"port": pPortInt,
			"err":  err,
		}).Fatal("server cant start listening")
		panic(-1)
	}
}