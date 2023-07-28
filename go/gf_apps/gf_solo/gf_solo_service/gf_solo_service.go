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
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_admin_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_home_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_landing_page_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_analytics_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_ml_lib"
	"github.com/gloflow/gloflow/go/gf_web3/gf_web3_lib"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang_server/gf_lang_service"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func Run(pConfig *GFconfig,
	pRuntimeSys *gf_core.RuntimeSys) {

	yellow := color.New(color.BgYellow).Add(color.FgBlack).SprintFunc()
	green  := color.New(color.BgGreen).Add(color.FgBlack).SprintFunc()

	pRuntimeSys.LogNewFun("INFO", fmt.Sprintf("%s%s\n", yellow("GF_SOLO"), green("===============")), nil)
	
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

	//-------------
	// GF_LANDING_PAGE
	// landing_page goes first, its handlers, because it contains the root path handler ("/")
	// and that should match first.
	gf_landing_page_lib.InitService(pConfig.TemplatesPathsMap,
		gfSoloHTTPmux,
		pRuntimeSys)

	//-------------
	// GF_IDENTITY

	gfIdentityServiceInfo := &gf_identity_core.GFserviceInfo{
		NameStr:       "gf_identity",
		DomainBaseStr: pConfig.DomainBaseStr,

		AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,

		AuthLoginURLstr:                       "/v1/identity/login_ui", // on email confirm redirect user to this
		AuthLoginSuccessRedirectURLstr:        "/v1/home/main", // on login success redirecto to home
		EnableEventsAppBool:                   true,
		EnableUserCredsInSecretsStoreBool:     true,
		EnableEmailBool:                       true,
		EnableEmailRequireConfirmForLoginBool: true,

		// ADD!! - for now regular users are not required to MFA confirm for login.
		//         there should be an option for users to be able to enable this
		//         individually if they so desire.
		EnableMFArequireConfirmForLoginBool: false,
	}
	keyServer, gfErr := gf_identity.InitService(pConfig.TemplatesPathsMap,
		gfSoloHTTPmux,
		gfIdentityServiceInfo,
		pRuntimeSys)
	if gfErr != nil {
		return
	}

	//-------------
	// GF_ADMIN - its started in a separate goroutine and listening on a diff
	//            port than the main service.
	sentryHubClone := sentry.CurrentHub().Clone()
	go func(pLocalHub *sentry.Hub) {

		adminHTTPmux := http.NewServeMux()

		adminServiceInfo := &gf_admin_lib.GFserviceInfo{
			NameStr:                           "gf_admin",
			AuthSubsystemTypeStr:              pConfig.AuthSubsystemTypeStr,
			AdminEmailStr:                     pConfig.AdminEmailStr,
			EnableEventsAppBool:               true,
			EnableUserCredsInSecretsStoreBool: true,
			EnableEmailBool:                   true,
		}


		// IMPORTANT!! - since admin is listening on its own port, and likely its own domain
		//               we want further isolation from main app handlers by
		//               instantiating gf_identity handlers dedicated to admin.
		adminIdentityServiceInfo := &gf_identity_core.GFserviceInfo{
			NameStr:       "gf_identity_admin",
			DomainBaseStr: pConfig.DomainAdminBaseStr,

			AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,

			AdminMFAsecretKeyBase32str: pConfig.AdminMFAsecretKeyBase32str,
			AuthLoginURLstr:            "/v1/admin/login_ui", // on email confirm redirect user to this

			// FEATURE_FLAGS
			EnableEventsAppBool:                   true,
			EnableUserCredsInSecretsStoreBool:     true,
			EnableEmailBool:                       true,
			EnableEmailRequireConfirmForLoginBool: true,
			EnableMFArequireConfirmForLoginBool:   true, // admins have to MFA confirm to login
			
		}

		gfErr := gf_admin_lib.InitNewService(pConfig.TemplatesPathsMap,
			adminServiceInfo,
			adminIdentityServiceInfo,
			adminHTTPmux,
			pLocalHub,
			pRuntimeSys)
		if gfErr != nil {
			return
		}

		// SERVER_INIT - blocking
		gf_rpc_lib.ServerInitWithMux("gf_solo_admin", portAdminInt, adminHTTPmux)

	}(sentryHubClone)

	//-------------
	// GF_HOME

	homeServiceInfo := &gf_home_lib.GFserviceInfo{
		AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,
		AuthLoginURLstr:      "/v1/identity/login_ui", // if not logged in redirect users to this
		KeyServer:            keyServer,
	}

	gfErr = gf_home_lib.InitService(pConfig.TemplatesPathsMap,
		homeServiceInfo,
		gfSoloHTTPmux,
		pRuntimeSys)
	if gfErr != nil {
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
		return
	}
	
	gfImagesServiceInfo := &gf_images_core.GFserviceInfo{
		Mongodb_host_str:                     pConfig.MongoHostStr,
		Mongodb_db_name_str:                  pConfig.MongoDBnameStr,

		ImagesStoreLocalDirPathStr:           imagesConfig.ImagesStoreLocalDirPathStr,
		ImagesThumbnailsStoreLocalDirPathStr: imagesConfig.ImagesThumbnailsStoreLocalDirPathStr,
		VideoStoreLocalDirPathStr:            imagesConfig.VideoStoreLocalDirPathStr,
		Media_domain_str:                     imagesConfig.Media_domain_str,
		Images_main_s3_bucket_name_str:       imagesConfig.Main_s3_bucket_name_str,

		Templates_paths_map: pConfig.TemplatesPathsMap,

		//-------------------------
		// AUTH_SUBSYSTEM_TYPE
		AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,

		// AUTH
		// on user trying to access authed endpoint while not logged in, redirect to this
		AuthLoginURLstr: "/v1/identity/login_ui",
		KeyServer:       keyServer,
		
		//-------------------------
		
		// IMAGES_STORAGE
		UseNewStorageEngineBool: pConfig.ImagesUseNewStorageEngineBool,

		// IPFS
		IPFSnodeHostStr: imagesConfig.IPFSnodeHostStr,
	}

	imagesJobsMngrCh := gf_images_lib.InitService(gfSoloHTTPmux,
		gfImagesServiceInfo,
		imagesConfig,
		pRuntimeSys)

	//-------------
	// GF_ANALYTICS
	
	gfAnalyticsServiceInfo := &gf_analytics_lib.GFserviceInfo{

		Crawl__config_file_path_str:      pConfig.CrawlConfigFilePathStr,
		Crawl__cluster_node_type_str:     pConfig.CrawlClusterNodeTypeStr,
		Crawl__images_local_dir_path_str: pConfig.CrawlImagesLocalDirPathStr,

		Media_domain_str:       imagesConfig.Media_domain_str,
		Py_stats_dirs_lst:      pConfig.AnalyticsPyStatsDirsLst,
		Run_indexer_bool:       pConfig.AnalyticsRunIndexerBool,
		Elasticsearch_host_str: pConfig.ElasticsearchHostStr,

		Templates_paths_map: pConfig.TemplatesPathsMap,

		// IMAGES_STORAGE
		ImagesUseNewStorageEngineBool: pConfig.ImagesUseNewStorageEngineBool,
	}
	gf_analytics_lib.InitService(gfAnalyticsServiceInfo,
		gfSoloHTTPmux,
		pRuntimeSys)

	//-------------
	// GF_PUBLISHER

	// FIX!! - find a soloution where gf_solo gf_publisher functions can invoke
	//         gf_images functions in the same process if in non-distributed mode.
	//         specifying gf_images host
	//         is there because of the default distributed design that assumes
	//         gf_publisher and gf_images run as separate processes.
	gfImagesServiceHostPortStr := "127.0.0.1"
	gfImagesRuntimeInfo := &gf_publisher_lib.GF_images_extern_runtime_info{
		Jobs_mngr:               nil, // indicates not to send in-process messages to jobs_mngr goroutine, instead use HTTP REST API of gf_images
		Service_host_port_str:   gfImagesServiceHostPortStr,
		Templates_dir_paths_map: pConfig.TemplatesPathsMap,
	}
	
	gf_publisher_lib.InitService(gfSoloHTTPmux,
		gfImagesRuntimeInfo,
		pRuntimeSys)

	//-------------
	// GF_TAGGER
	gf_tagger_lib.InitService(pConfig.TemplatesPathsMap,
		imagesJobsMngrCh,
		gfSoloHTTPmux,
		pRuntimeSys)

	//-------------
	// GF_ML
	gf_ml_lib.InitService(gfSoloHTTPmux, pRuntimeSys)

	//-------------
	// GF_WEB3
	
	web3Config := &gf_eth_core.GF_config{
		AlchemyAPIkeyStr: pConfig.AlchemyAPIkeyStr,
	}
	gf_web3_lib.InitService(pConfig.AuthSubsystemTypeStr,
		keyServer,
		gfSoloHTTPmux,
		web3Config,
		imagesJobsMngrCh,
		pRuntimeSys)
	
	//-------------
	// GF_LANG

	// ADD!! - expose this as a ENV var. the user should be able to 
	//         disable the language server from starting up and accepting
	//         programs for execution.
	enableLangBool := true

	if enableLangBool {
		langServiceInfo := &gf_lang_service.GFserviceInfo{
			NameStr: "gf_lang", 
		}
		langConfig := &gf_lang_service.GFconfig{}
		gfErr := gf_lang_service.InitService(gfSoloHTTPmux,
			langServiceInfo,
			langConfig,
			pRuntimeSys)
		if gfErr != nil {
			return
		}
	}

	//-------------
	// METRICS - start prometheus metrics endpoint, and get core_metrics
	coreMetrics := gf_core.MetricsInit("/metrics", portMetricsInt)
	pRuntimeSys.Metrics = coreMetrics
	
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