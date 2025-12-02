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
	"context"
	"net/http"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/go-playground/validator"
)

//-------------------------------------------------
// RUNTIME_SYS

type RuntimeSys struct {
	
	AppNameStr	   string
	ServiceNameStr string
	EnvStr         string
	Debug_bool     bool // if debug mode is enabled (some places will print extra info in debug mode)
	LogFun         func(string, string)
	LogNewFun      GFlogFun

	// DB
	SQLdb      *sql.DB
	SQLdsnStr  string
	Mongo_db   *mongo.Database
	Mongo_coll *mongo.Collection // main mongodb collection to use when none is specified
	
	// ERRORS
	Errors_send_to_mongodb_bool bool // if errors should be persisted to Mongodb
	ErrorsSendToSentryBool  bool // if errors should be sent to Sentry service

	NamesPrefixStr string

	Validator *validator.Validate

	// PLUGINS
	ExternalPlugins *ExternalPlugins

	Metrics *GFmetrics

	// HTTP_PROXY
	// if a http proxy should be use this value is set
	// "http://proxy:8888"
	HTTPproxyServerURIstr string

	// SENTRY - used to pass the DNS to py sentry clients.
	SentryDSNstr string

	// EVENTS
	EnableEventsAppBool bool
}

//-------------------------------------------------
// PLUGINS

type ExternalPlugins struct {

	//---------------------------
	// RPC_HANDLERS
	RPChandlersGetCallback func(*RuntimeSys) ([]HTTPhandlerInfo, *GFerror)
	RPCreqPreProcessCallback func(*http.Request, http.ResponseWriter, context.Context, *RuntimeSys) (bool, *GFerror)

	// CORS_DOMAINS - domains that are allowed to access the API, beyond the domain
	// 			  	  that the API is hosted on.
	CORSoriginDomainsLst []string

	//---------------------------
	// IDENTITY

	IdentitySessionValidateApiKeyCallback func(string, *http.Request, context.Context, *RuntimeSys) (bool, GF_ID, *GFerror)
	
	//---------------------------
	// IMAGES
	ImageFilterMetadataCallback func(map[string]interface{}) map[string]interface{}
	
	//---------------------------
	// EVENTS
	EventCallback func(string, map[string]interface{}, string, GF_ID, *RuntimeSys) *GFerror

	//---------------------------
	// SECRETS
	SecretCreateCallback func(string, map[string]interface{}, string, *RuntimeSys) *GFerror
	SecretGetCallback    func(string, *RuntimeSys) (map[string]interface{}, *GFerror)

	//---------------------------
	// EMAIL
	
	// called on every sending of email in the system
	EmailSendingCallback func(string, string, string, string, *RuntimeSys) *GFerror

	//---------------------------
	/*
	// NFT
	// get metadata on defined fetchers
	NFTgetFetchersMetaCallback func() map[string]map[string]interface{}

	// fetch NFTs for a particular owner account using a particular fetcher (with name)
	NFTfetchForOwnerAddressCallback func(string, string, *RuntimeSys) *GFerror
	*/
	//---------------------------
}