/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

package gf_solo_images_service

import (
	"fmt"
	"os/user"
	"strconv"
	"path"
	"time"
	"net/http"
	"github.com/getsentry/sentry-go"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// RUN
/*
runs the gf_images solo server. its a standalone server, independent of gf_solo.
it initializes:
- gf_identity (JSON auth handlers, key server, no login UI)
- gf_images (HTTP handlers, jobs manager, flows, storage)
- prometheus metrics endpoint
- HTTP handlers supplied by external hooks
and then blocks serving HTTP.
*/

func Run(pConfig *gf_core.GFconfig,
	pRuntimeSys *gf_core.RuntimeSys) {

	yellow := color.New(color.BgYellow).Add(color.FgBlack).SprintFunc()
	green  := color.New(color.BgGreen).Add(color.FgBlack).SprintFunc()

	pRuntimeSys.LogNewFun("INFO", fmt.Sprintf("%s%s", yellow("GF_SOLO_IMAGES"), green("===============")), nil)

	// flush buffered Sentry events on process exit
	defer sentry.Flush(2 * time.Second)

	// EVENTS
	enableEventsAppBool := true
	pRuntimeSys.EnableEventsAppBool = enableEventsAppBool

	//-------------
	// CONFIG
	portInt, err := strconv.Atoi(pConfig.PortStr)
	if err != nil {
		panic(err)
	}

	portMetricsInt, err := strconv.Atoi(pConfig.PortMetricsStr)
	if err != nil {
		pRuntimeSys.LogNewFun("WARN", "invalid or missing port_metrics, using default 9110", nil)
		portMetricsInt = 9110
	}

	pRuntimeSys.LogNewFun("DEBUG", "gf_solo_images service config", nil)
	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(pConfig)
	}

	// IMPORTANT!! - pure API server, there is no login page. empty URL means that auth failures
	//               are returned as errors instead of being redirected.
	authLoginURLstr := ""

	//-------------
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("(%s), dir (%s)\n", currentUser.Username, currentUser.HomeDir)

	// HTTP_MUX
	httpMux := http.NewServeMux()

	// METRICS
	// start prometheus metrics endpoint, and get core_metrics
	coreMetrics := gf_core.MetricsInit("/metrics", portMetricsInt)
	pRuntimeSys.Metrics = coreMetrics

	metricsRPCglobal := gf_rpc_lib.MetricsCreateGlobal()

	//-------------
	// GF_IDENTITY
	// auth is needed by the images upload/share/flows handlers.

	identityServiceInfo := &gf_identity_core.GFserviceInfo{
		NameStr: "gf_identity",

		// DOMAINS
		DomainBaseStr:           pConfig.DomainBaseStr,
		DomainForAuthCookiesStr: &pConfig.DomainForAuthCookiesStr,

		AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,
		AuthLoginURLstr:                authLoginURLstr,             // on email confirm redirect user to this
		AuthLoginSuccessRedirectURLstr: "",

		// EVENTS
		EnableEventsAppBool: enableEventsAppBool,

		// no HTML rendering
		DisableLoginUIbool: true,

		EnableUserCredsInSecretsStoreBool:     true,
		EnableEmailBool:                       true,
		EnableEmailRequireConfirmForLoginBool: true,
		EnableMFArequireConfirmForLoginBool:   false,
	}
	keyServer, gfErr := gf_identity.InitService(nil, // no templates
		httpMux,
		metricsRPCglobal,
		identityServiceInfo,
		pRuntimeSys)
	if gfErr != nil {
		pRuntimeSys.LogNewFun("ERROR", "failed to init gf_identity service", map[string]interface{}{"err": gfErr.Error.Error()})
		return
	}

	//-------------
	// GF_IMAGES

	// CONFIG
	imagesConfig, gfErr := gf_images_core.ConfigGet(pConfig.ImagesConfigFilePathStr,
		pConfig.ImagesUseNewStorageEngineBool,
		pConfig.IPFSnodeHostStr,
		pRuntimeSys)
	if gfErr != nil {
		pRuntimeSys.LogNewFun("ERROR", "failed to load gf_images config", map[string]interface{}{"err": gfErr.Error.Error()})
		return
	}

	imagesServiceInfo := &gf_images_core.GFserviceInfo{
		Port_str:      pConfig.PortStr,
		DomainBaseStr: pConfig.DomainBaseStr,

		Mongodb_host_str:    pConfig.MongoHostStr,
		Mongodb_db_name_str: pConfig.MongoDBnameStr,

		ImagesStoreLocalDirPathStr:           imagesConfig.ImagesStoreLocalDirPathStr,
		ImagesThumbnailsStoreLocalDirPathStr: imagesConfig.ImagesThumbnailsStoreLocalDirPathStr,
		VideoStoreLocalDirPathStr:            imagesConfig.VideoStoreLocalDirPathStr,
		MediaDomainStr:                       imagesConfig.MediaDomainStr,
		ImagesMainS3bucketNameStr:            imagesConfig.MainS3bucketNameStr,

		ConfigFilePathStr: pConfig.ImagesConfigFilePathStr,

		EmailSharingSenderAddressStr: imagesConfig.EmailSharingSenderAddressStr,

		//-------------------------
		// AUTH
		AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,

		// on user trying to access authed endpoint while not logged in, redirect to this
		AuthLoginURLstr: authLoginURLstr,
		KeyServer:       keyServer,

		//-------------------------

		// IMAGES_STORAGE
		UseNewStorageEngineBool: pConfig.ImagesUseNewStorageEngineBool,

		// IPFS
		IPFSnodeHostStr: imagesConfig.IPFSnodeHostStr,

		// EVENTS
		EnableEventsAppBool: enableEventsAppBool,
	}

	// all images HTTP handlers, jobs manager, flows, gifs, storage
	_ = gf_images_lib.InitSoloService(httpMux,
		imagesServiceInfo,
		imagesConfig,
		pRuntimeSys)

	//-------------
	// HOOKS - register external http handlers

	if pRuntimeSys.ExternalHooks != nil && pRuntimeSys.ExternalHooks.RPChandlersGetCallback != nil {

		//-------------
		// USER_RPC_HANDLERS
		handlersLst, handlersV2lst, gfErr := pRuntimeSys.ExternalHooks.RPChandlersGetCallback(httpMux,
			pRuntimeSys)
		if gfErr != nil {
			pRuntimeSys.LogNewFun("ERROR", "failed to get RPC handlers from external hooks", map[string]interface{}{"err": gfErr.Error.Error()})
			return
		}

		//-------------

		authSubsystemTypeStr := pConfig.AuthSubsystemTypeStr
		metricsGroupNameStr  := "gf_solo_images__hooks_rpc_handlers"

		gf_rpc_lib.CreateHandlersHTTP(handlersLst,
			httpMux,
			authSubsystemTypeStr,
			authLoginURLstr,
			keyServer,
			metricsGroupNameStr,
			pRuntimeSys)

		gf_rpc_lib.CreateHandlersV2http(handlersV2lst,
			httpMux,
			authSubsystemTypeStr,
			authLoginURLstr,
			keyServer,
			metricsGroupNameStr,
			pRuntimeSys)
	}

	//-------------
	// SERVER_INIT - blocking
	gf_rpc_lib.ServerInitWithMux("gf_solo_images", portInt, httpMux)

	//-------------
}

//-------------------------------------------------

func RuntimeGet(pConfigPathStr string,
	pExternalHooks *gf_core.ExternalHooks,
	pLogFun        func(string, string),
	pLogNewFun     gf_core.GFlogFun) (*gf_core.RuntimeSys, *gf_core.GFconfig, error) {

	//--------------------
	// RUNTIME_SYS
	runtimeSys := &gf_core.RuntimeSys{
		AppNameStr:     "gf_solo_images",
		ServiceNameStr: "gf_solo_images",
		LogFun:         pLogFun,
		LogNewFun:      pLogNewFun,

		// VALIDATOR
		Validator: gf_core.ValidateInit(),

		// EXTERNAL_HOOKS
		ExternalHooks: pExternalHooks,
	}

	//--------------------
	// CONFIG
	configDirPathStr := path.Dir(pConfigPathStr)  // "./config/"
	configNameStr    := path.Base(pConfigPathStr) // "gf_solo_images"

	config, err := ConfigInit(configDirPathStr, configNameStr, runtimeSys)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to load config")
		return nil, nil, err
	}

	// IMPORTANT!! - gf_identity session validation reads runtimeSys.Config, it has to be set.
	runtimeSys.Config = config
	runtimeSys.EnvStr = config.EnvStr

	// GF_IDENTITY
	runtimeSys.IdentitySubsystemTypeStr = config.AuthSubsystemTypeStr

	//--------------------
	// SENTRY - ERROR_REPORTING
	if config.SentryDSNstr != "" {

		fmt.Println("Initializing Sentry error reporting...")

		runtimeSys.SentryDSNstr           = config.SentryDSNstr
		runtimeSys.ErrorsSendToDBbool     = true
		runtimeSys.ErrorsSendToSentryBool = true

		sentrySampleRateDefaultF        := 1.0
		sentryTracingRateForHandlersMap := map[string]float64{}
		err := gf_core.ErrorInitSentry(config.SentryDSNstr,
			sentryTracingRateForHandlersMap,
			sentrySampleRateDefaultF)
		if err != nil {
			return nil, nil, err
		}
	}

	//--------------------
	// SQL

	sqlDB, sqlDSNstr, gfErr := gf_core.DBsqlConnect(config.SQLdbNameStr,
		config.SQLuserNameStr,
		config.SQLpassStr,
		config.SQLhostStr,
		"require", // SSL mode - required for PostgreSQL 18
		runtimeSys)
	if gfErr != nil {
		return nil, nil, gfErr.Error
	}

	runtimeSys.SQLdb     = sqlDB
	runtimeSys.SQLdsnStr = sqlDSNstr

	//--------------------
	// MONGODB
	mongodbHostStr := config.MongoHostStr
	mongodbURLstr  := fmt.Sprintf("mongodb://%s", mongodbHostStr)
	fmt.Printf("mongodb_host    - %s\n", mongodbHostStr)
	fmt.Printf("mongodb_db_name - %s\n", config.MongoDBnameStr)

	mongodbDB, _, gfErr := gf_core.MongoConnectNew(mongodbURLstr,
		config.MongoDBnameStr,
		nil,
		runtimeSys)
	if gfErr != nil {
		return nil, nil, gfErr.Error
	}

	runtimeSys.Mongo_db   = mongodbDB
	runtimeSys.Mongo_coll = mongodbDB.Collection("data_symphony")
	fmt.Printf("mongodb connected...\n")

	//--------------------
	return runtimeSys, config, nil
}
