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

package gf_solo_service

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
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_ml_lib"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func Run(pConfig *GFconfig,
	pRuntimeSys *gf_core.RuntimeSys) {

	yellow := color.New(color.BgYellow).Add(color.FgBlack).SprintFunc()
	green  := color.New(color.BgGreen).Add(color.FgBlack).SprintFunc()

	pRuntimeSys.LogNewFun("INFO", fmt.Sprintf("%s%s", yellow("GF_SOLO_IMAGES"), green("===============")), nil)
	


	// EVENTS
	enableEventsAppBool := true
	pRuntimeSys.EnableEventsAppBool = enableEventsAppBool

	//-------------
	// CONFIG
	portMetricsInt := 9110

	portInt, err := strconv.Atoi(pConfig.PortStr)
	if err != nil {
		panic(err)
	}

	portAdminInt, err := strconv.Atoi(pConfig.PortAdminStr)
	if err != nil {
		panic(err)
	}
	
	pRuntimeSys.LogNewFun("DEBUG", "gf_solo service config", nil)
	if gf_core.LogsIsDebugEnabled() {
		spew.Dump(pConfig)
	}
	

	authLoginURLstr := "/v1/identity/login_ui"
	
	//-------------
	user, err := user.Current()
	if err != nil {
        panic(err)
    }
	fmt.Printf("(%s), dir (%s)\n", user.Username, user.HomeDir)

	// VALIDATOR
	validator := gf_core.ValidateInit()
	pRuntimeSys.Validator = validator

	// HTTP_MUX
	gfSoloHTTPmux := http.NewServeMux()

	// METRICS
	metricsRPCglobal := gf_rpc_lib.MetricsCreateGlobal()

	//-------------
	// GF_IMAGES

	// CONFIG
	imagesConfig, gfErr := gf_images_core.ConfigGet(pConfig.ImagesConfigFilePathStr,
		pConfig.ImagesUseNewStorageEngineBool,
		pConfig.IPFSnodeHostStr,
		pRuntimeSys)
	if gfErr != nil {
		return
	}
	
	gfImagesServiceInfo := &gf_images_core.GFserviceInfo{
		DomainBaseStr: pConfig.DomainBaseStr,
		Mongodb_host_str:                     pConfig.MongoHostStr,
		Mongodb_db_name_str:                  pConfig.MongoDBnameStr,

		ImagesStoreLocalDirPathStr:           imagesConfig.ImagesStoreLocalDirPathStr,
		ImagesThumbnailsStoreLocalDirPathStr: imagesConfig.ImagesThumbnailsStoreLocalDirPathStr,
		VideoStoreLocalDirPathStr:            imagesConfig.VideoStoreLocalDirPathStr,
		MediaDomainStr:                       imagesConfig.MediaDomainStr,
		ImagesMainS3bucketNameStr:            imagesConfig.MainS3bucketNameStr,

		TemplatesPathsMap: pConfig.TemplatesPathsMap,

		EmailSharingSenderAddressStr: imagesConfig.EmailSharingSenderAddressStr,

		//-------------------------
		// AUTH_SUBSYSTEM_TYPE
		AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,

		// AUTH
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

	imagesJobsMngrCh := gf_images_lib.InitSoloService(gfSoloHTTPmux,
		gfImagesServiceInfo,
		imagesConfig,
		pRuntimeSys)

	//-------------
	// GF_ML
	gf_ml_lib.InitService(gfSoloHTTPmux, pRuntimeSys)

	//-------------
	// METRICS - start prometheus metrics endpoint, and get core_metrics
	coreMetrics := gf_core.MetricsInit("/metrics", portMetricsInt)
	pRuntimeSys.Metrics = coreMetrics
	
	//-------------
	// REGISTER_EXTERN_HTTP_HANDLERS

	if pRuntimeSys.ExternalPlugins != nil && pRuntimeSys.ExternalPlugins.RPChandlersGetCallback != nil {

		//-------------
		// USER_RPC_HANDLERS
		handlersLst, gfErr := pRuntimeSys.ExternalPlugins.RPChandlersGetCallback(pRuntimeSys)
		if gfErr != nil {
			return
		}
		
		//-------------

		authSubsystemTypeStr := pConfig.AuthSubsystemTypeStr
		metricsGroupNameStr  := "gf_solo__plugin_rpc_handlers"
		
		gf_rpc_lib.CreateHandlersHTTP(metricsGroupNameStr,
			handlersLst,
			gfSoloHTTPmux,
			authSubsystemTypeStr,
			authLoginURLstr,
			keyServer,
			pRuntimeSys)
	}

	//-------------
	// SERVER_INIT - blocking
	gf_rpc_lib.ServerInitWithMux("gf_solo", portInt, gfSoloHTTPmux)

	//-------------
}

//-------------------------------------------------

func RuntimeGet(pConfigPathStr string,
	pExternalPlugins *gf_core.ExternalPlugins,
	pLogFun          func(string, string),
	pLogNewFun       gf_core.GFlogFun) (*gf_core.RuntimeSys, *GFconfig, error) {

	//--------------------
	// CONFIG
	configDirPathStr := path.Dir(pConfigPathStr)  // "./../config/"
	configNameStr    := path.Base(pConfigPathStr) // "gf_solo"
	
	config, err := ConfigInit(configDirPathStr, configNameStr)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to load config")
		return nil, nil, err
	}

	//--------------------
	// SENTRY - ERROR_REPORTING
	if config.SentryEndpointStr != "" {

		sentryEndpointStr := config.SentryEndpointStr
		sentrySampleRateDefaultF := 1.0
		sentryTracingRateForHandlersMap := map[string]float64{
			
		}
		err := gf_core.ErrorInitSentry(sentryEndpointStr,
			sentryTracingRateForHandlersMap,
			sentrySampleRateDefaultF)
		if err != nil {
			panic(err)
		}

		defer sentry.Flush(2 * time.Second)
	}

	//--------------------
	// RUNTIME_SYS
	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_solo",
		EnvStr:         config.EnvStr,
		LogFun:         pLogFun,
		LogNewFun:      pLogNewFun,

		// SENTRY - enable it for error reporting
		ErrorsSendToSentryBool: true,

		// EXTERNAL_PLUGINS
		ExternalPlugins: pExternalPlugins,

		// SENTRY_DSN
		SentryDSNstr: config.SentryEndpointStr,
	}
	
	//--------------------
	// SQL

	sqlDB, gfErr := gf_core.DBsqlConnect(config.SQLdbNameStr,
		config.SQLuserNameStr,
		config.SQLpassStr,
		config.SQLhostStr,
		runtimeSys)

	runtimeSys.SQLdb = sqlDB

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