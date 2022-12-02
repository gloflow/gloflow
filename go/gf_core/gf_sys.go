/*
MIT License

Copyright (c) 2019 Ivan Trajkovic

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

package gf_core

import (
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/go-playground/validator"
)

//-------------------------------------------------
// RUNTIME_SYS

type RuntimeSys struct {
	Service_name_str string
	LogFun           func(string, string)
	LogNewFun        func(pMsgStr string, pGroupStr string, pLevelStr string, pMetaMap map[string]interface{})

	Mongo_db         *mongo.Database
	Mongo_coll       *mongo.Collection // main mongodb collection to use when none is specified
	Debug_bool       bool              // if debug mode is enabled (some places will print extra info in debug mode)

	// ERRORS
	Errors_send_to_mongodb_bool bool // if errors should be persisted to Mongodb
	ErrorsSendToSentryBool  bool // if errors should be sent to Sentry service

	Names_prefix_str string

	Validator *validator.Validate

	// PLUGINS
	ExternalPlugins *ExternalPlugins

	Metrics *GFmetrics

	// HTTP_PROXY
	// if a http proxy should be use this value is set
	// "http://proxy:8888"
	HTTPproxyServerURIstr string
}

//-------------------------------------------------
// PLUGINS

type ExternalPlugins struct {

	//---------------------------
	// EVENTS
	EventCallback func(string, map[string]interface{}, *RuntimeSys) *GFerror

	//---------------------------
	// SECRETS
	SecretCreateCallback func(string, map[string]interface{}, string, *RuntimeSys) *GFerror
	SecretGetCallback    func(string, *RuntimeSys) (map[string]interface{}, *GFerror)

	//---------------------------
	// NFT
	// get metadata on defined fetchers
	NFTgetFetchersMetaCallback func() map[string]map[string]interface{}

	// fetch NFTs for a particular owner account using a particular fetcher (with name)
	NFTfetchForOwnerAddressCallback func(string, string, *RuntimeSys) *GFerror

	//---------------------------
}