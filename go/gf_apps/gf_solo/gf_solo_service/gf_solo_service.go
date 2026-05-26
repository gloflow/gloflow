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
	"net/http"
	"github.com/getsentry/sentry-go"
	"github.com/fatih/color"
	gf_core "github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_admin_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_home_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_landing_page_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_analytics_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib"
	"github.com/gloflow/gloflow/go/gf_web3/gf_web3_lib"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang_server/gf_lang_service"
	"github.com/davecgh/go-spew/spew"
	// "github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib"
)

//-------------------------------------------------

func Run(pConfig *gf_core.GFconfig,
	pRuntimeSys *gf_core.RuntimeSys) {

	yellow := color.New(color.BgYellow).Add(color.FgBlack).SprintFunc()
	green  := color.New(color.BgGreen).Add(color.FgBlack).SprintFunc()

	pRuntimeSys.LogNewFun("INFO", fmt.Sprintf("%s%s", yellow("GF_SOLO"), green("===============")), nil)

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
	// GF_LANDING_PAGE
	// landing_page goes first, its handlers, because it contains the root path handler ("/")
	// and that should match first.
	gf_landing_page_lib.InitService(pConfig.TemplatesPathsMap,
		gfSoloHTTPmux,
		pRuntimeSys)

	//-------------
	// GF_IDENTITY

	gfIdentityServiceInfo := &gf_identity_core.GFserviceInfo{
		NameStr: "gf_identity",

		// DOMAINS
		DomainBaseStr:           pConfig.DomainBaseStr,
		DomainForAuthCookiesStr: &pConfig.DomainForAuthCookiesStr,

		AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,
		AuthLoginURLstr:                authLoginURLstr, // on email confirm redirect user to this
		AuthLoginSuccessRedirectURLstr: "/v1/home/main", // on login success redirecto to home

		// EVENTS
		EnableEventsAppBool: enableEventsAppBool,

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
		metricsRPCglobal,
		gfIdentityServiceInfo,
		pRuntimeSys)
	if gfErr != nil {
		return
	}

	//-------------
	/*
	GF_ADMIN - its started in a separate goroutine and listening on a diff
		port than the main service.
	*/
	sentryHubClone := sentry.CurrentHub().Clone()
	go func(pLocalHub *sentry.Hub) {

		adminHTTPmux := http.NewServeMux()

		adminServiceInfo := &gf_admin_lib.GFserviceInfo{
			NameStr:                           "gf_admin",
			AuthSubsystemTypeStr:              pConfig.AuthSubsystemTypeStr,
			AdminEmailStr:                     pConfig.AdminEmailStr,
			EnableEventsAppBool:               enableEventsAppBool,
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
			EnableEventsAppBool:                   enableEventsAppBool,
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
		AuthLoginURLstr:      authLoginURLstr, // if not logged in redirect users to this
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

	imagesJobsMngrCh := gf_images_lib.InitService(gfSoloHTTPmux,
		gfImagesServiceInfo,
		imagesConfig,
		pRuntimeSys)

	// Store JobsMngrCh in RuntimeSys for MCP tools and other subsystems
	// This enables flow integration for trader plots and other image operations
	pRuntimeSys.JobsMngrCh = imagesJobsMngrCh
	pRuntimeSys.LogNewFun("INFO", "JobsMngrCh stored in RuntimeSys", map[string]interface{}{
		"jobs_mngr_ch_nil": imagesJobsMngrCh == nil,
		"runtime_sys_ptr":  fmt.Sprintf("%p", pRuntimeSys),
	})

	//-------------
	// GF_ANALYTICS

	gfAnalyticsServiceInfo := &gf_analytics_lib.GFserviceInfo{

		Media_domain_str:  imagesConfig.MediaDomainStr,
		Py_stats_dirs_lst: pConfig.AnalyticsPyStatsDirsLst,
		Run_indexer_bool:  pConfig.AnalyticsRunIndexerBool,

		Templates_paths_map: pConfig.TemplatesPathsMap,

		// IMAGES_STORAGE
		ImagesUseNewStorageEngineBool: pConfig.ImagesUseNewStorageEngineBool,

		AuthSubsystemTypeStr: pConfig.AuthSubsystemTypeStr,
		AuthLoginURLstr:      authLoginURLstr, // on email confirm redirect user to this
		KeyServer:            keyServer,

		// Elasticsearch_host_str: pConfig.ElasticsearchHostStr,
		// Crawl__config_file_path_str:      pConfig.CrawlConfigFilePathStr,
		// Crawl__cluster_node_type_str:     pConfig.CrawlClusterNodeTypeStr,
		// Crawl__images_local_dir_path_str: pConfig.CrawlImagesLocalDirPathStr,
	}
	gf_analytics_lib.InitService(gfAnalyticsServiceInfo,
		metricsRPCglobal,
		gfSoloHTTPmux,
		pRuntimeSys)

	//-------------
	// GF_TAGGER
	gf_tagger_lib.InitService(pConfig.AuthSubsystemTypeStr,
		authLoginURLstr,
		keyServer,
		gfSoloHTTPmux,
		pConfig.TemplatesPathsMap,
		imagesJobsMngrCh,
		pRuntimeSys)

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
	// PLUGIN - register extern http handlers

	if pRuntimeSys.ExternalPlugins != nil && pRuntimeSys.ExternalPlugins.RPChandlersGetCallback != nil {

		//-------------
		// USER_RPC_HANDLERS
		pRuntimeSys.LogNewFun("INFO", "Calling plugins RPChandlersGetCallback", map[string]interface{}{
			"runtime_sys_ptr": fmt.Sprintf("%p", pRuntimeSys),
			"jobs_mngr_ch":    pRuntimeSys.JobsMngrCh,
		})
		handlersLst, handlersV2lst, gfErr := pRuntimeSys.ExternalPlugins.RPChandlersGetCallback(gfSoloHTTPmux, pRuntimeSys)
		if gfErr != nil {

			return
		}

		//-------------

		authSubsystemTypeStr := pConfig.AuthSubsystemTypeStr
		metricsGroupNameStr  := "gf_solo__plugin_rpc_handlers"

		gf_rpc_lib.CreateHandlersHTTP(handlersLst,
			gfSoloHTTPmux,
			authSubsystemTypeStr,
			authLoginURLstr,
			keyServer,
			metricsGroupNameStr,
			pRuntimeSys)

		gf_rpc_lib.CreateHandlersV2http(handlersV2lst,
			gfSoloHTTPmux,
			authSubsystemTypeStr,
			authLoginURLstr,
			keyServer,
			metricsGroupNameStr,
			pRuntimeSys)
	}

	//-------------
	// SERVER_INIT - blocking
	gf_rpc_lib.ServerInitWithMux("gf_solo", portInt, gfSoloHTTPmux)

	//-------------
}

//-------------------------------------------------
