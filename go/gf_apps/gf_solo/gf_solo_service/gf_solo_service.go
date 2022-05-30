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
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
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
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Run(pConfig *GF_config,
	pRuntimeSys *gf_core.RuntimeSys) {

	//-------------
	// CONFIG
	portMetricsInt := 9110

	portInt, err := strconv.Atoi(pConfig.Port_str)
	if err != nil {
		panic(err)
	}

	portAdminInt, err := strconv.Atoi(pConfig.Port_admin_str)
	if err != nil {
		panic(err)
	}

	//-------------
	
	yellow := color.New(color.BgYellow).Add(color.FgBlack).SprintFunc()
	green  := color.New(color.BgGreen).Add(color.FgBlack).SprintFunc()

	fmt.Printf("%s%s\n", yellow("GF_SOLO"), green("==============="))

	//-------------
	user, err := user.Current()
	if err != nil {
        panic(err)
    }
	fmt.Printf("(%s), dir (%s)\n", user.Username, user.HomeDir)

	// VALIDATOR
	validator := gf_core.Validate__init()
	pRuntimeSys.Validator = validator

	// HTTP_MUX
	gfSoloHTTPmux := http.NewServeMux()

	//-------------
	// GF_LANDING_PAGE
	// landing_page goes first, its handlers, because it contains the root path handler ("/")
	// and that should match first.
	gf_landing_page_lib.Init_service(pConfig.Templates_paths_map,
		gfSoloHTTPmux,
		pRuntimeSys)

	//-------------
	// GF_IDENTITY

	gf_identity__service_info := &gf_identity_lib.GF_service_info{
		Name_str:                       "gf_identity",
		Domain_base_str:                pConfig.Domain_base_str,
		AuthLoginURLstr:                "/landing/main", // on email confirm redirect user to this
		AuthLoginSuccessRedirectURLstr: "/v1/home/main", // on login success redirecto to home
		Enable_events_app_bool:                  true,
		Enable_user_creds_in_secrets_store_bool: true,
		Enable_email_bool:                       true,
		Enable_email_require_confirm_for_login_bool: true,

		// ADD!! - for now regular users are not required to MFA confirm for login.
		//         there should be an option for users to be able to enable this
		//         individually if they so desire.
		Enable_mfa_require_confirm_for_login_bool: false,
	}
	gfErr := gf_identity_lib.InitService(gfSoloHTTPmux,
		gf_identity__service_info,
		pRuntimeSys)
	if gfErr != nil {
		return
	}

	//-------------
	// GF_ADMIN - its started in a separate goroutine and listening on a diff
	//            port than the main service.
	sentry_hub_clone := sentry.CurrentHub().Clone()
	go func(p_local_hub *sentry.Hub) {

		adminHTTPmux := http.NewServeMux()

		admin_service_info := &gf_admin_lib.GF_service_info{
			Name_str:                                "gf_admin",
			Admin_email_str:                         pConfig.Admin_email_str,
			Enable_events_app_bool:                  true,
			Enable_user_creds_in_secrets_store_bool: true,
			Enable_email_bool:                       true,
		}


		// IMPORTANT!! - since admin is listening on its own port, and likely its own domain
		//               we want further isolation from main app handlers by
		//               instantiating gf_identity handlers dedicated to admin.
		admin_identity__service_info := &gf_identity_lib.GF_service_info{
			Name_str:                        "gf_admin_identity",
			Domain_base_str:                 pConfig.Domain_admin_base_str,
			Admin_mfa_secret_key_base32_str: pConfig.Admin_mfa_secret_key_base32_str,
			AuthLoginURLstr:                 "/v1/admin/login_ui", // on email confirm redirect user to this

			// FEATURE_FLAGS
			Enable_events_app_bool:                  true,
			Enable_user_creds_in_secrets_store_bool: true,
			Enable_email_bool:                       true,
			Enable_email_require_confirm_for_login_bool: true,
			Enable_mfa_require_confirm_for_login_bool:   true, // admins have to MFA confirm to login
			
		}

		gfErr := gf_admin_lib.InitNewService(pConfig.Templates_paths_map,
			admin_service_info,
			admin_identity__service_info,
			adminHTTPmux,
			p_local_hub,
			pRuntimeSys)
		if gfErr != nil {
			return
		}

		// SERVER_INIT - blocking
		gf_rpc_lib.Server__init_with_mux(portAdminInt, adminHTTPmux)

	}(sentry_hub_clone)

	//-------------
	// GF_HOME

	homeServiceInfo := &gf_home_lib.GFserviceInfo{
		AuthLoginURLstr: "/landing/main", // if not logged in redirect users to this
	}

	gfErr = gf_home_lib.InitService(pConfig.Templates_paths_map,
		homeServiceInfo,
		gfSoloHTTPmux,
		pRuntimeSys)
	if gfErr != nil {
		return
	}

	//-------------
	// GF_IMAGES

	// CONFIG
	gf_images__config, gfErr := gf_images_core.Config__get(pConfig.Images__config_file_path_str,
		pRuntimeSys)
	if gfErr != nil {
		return
	}
	
	gf_images__service_info := &gf_images_core.GFserviceInfo{
		Mongodb_host_str:                           pConfig.Mongodb_host_str,
		Mongodb_db_name_str:                        pConfig.Mongodb_db_name_str,

		Images_store_local_dir_path_str:            gf_images__config.Store_local_dir_path_str,
		Images_thumbnails_store_local_dir_path_str: gf_images__config.Thumbnails_store_local_dir_path_str,
		Media_domain_str:                           gf_images__config.Media_domain_str,
		Images_main_s3_bucket_name_str:             gf_images__config.Main_s3_bucket_name_str,

		AWS_access_key_id_str:                      pConfig.AWS_access_key_id_str,
		AWS_secret_access_key_str:                  pConfig.AWS_secret_access_key_str,
		AWS_token_str:                              pConfig.AWS_token_str,

		Templates_paths_map: pConfig.Templates_paths_map,

		// on user trying to access authed endpoint while not logged in, redirect to this
		AuthLoginURLstr: "/landing/main",
	}

	jobs_mngr_ch := gf_images_lib.Init_service(gfSoloHTTPmux,
		gf_images__service_info,
		gf_images__config,
		pRuntimeSys)

	//-------------
	// GF_ANALYTICS
	
	gf_analytics__service_info := &gf_analytics_lib.GF_service_info{

		Crawl__config_file_path_str:      pConfig.Crawl__config_file_path_str,
		Crawl__cluster_node_type_str:     pConfig.Crawl__cluster_node_type_str,
		Crawl__images_local_dir_path_str: pConfig.Crawl__images_local_dir_path_str,

		Media_domain_str:       gf_images__config.Media_domain_str,
		Py_stats_dirs_lst:      pConfig.Analytics__py_stats_dirs_lst,
		Run_indexer_bool:       pConfig.Analytics__run_indexer_bool,
		Elasticsearch_host_str: pConfig.Elasticsearch_host_str,

		AWS_access_key_id_str:     pConfig.AWS_access_key_id_str,
		AWS_secret_access_key_str: pConfig.AWS_secret_access_key_str,
		AWS_token_str:             pConfig.AWS_token_str,

		Templates_paths_map: pConfig.Templates_paths_map,
	}
	gf_analytics_lib.Init_service(gf_analytics__service_info,
		gfSoloHTTPmux,
		pRuntimeSys)

	//-------------
	// GF_PUBLISHER

	// FIX!! - find a soloution where gf_solo gf_publisher functions can invoke
	//         gf_images functions in the same process if in non-distributed mode.
	//         specifying gf_images host
	//         is there because of the default distributed design that assumes
	//         gf_publisher and gf_images run as separate processes.
	gf_images_service_host_port_str := "127.0.0.1"
	gf_images_runtime_info := &gf_publisher_lib.GF_images_extern_runtime_info{
		Jobs_mngr:               nil, // indicates not to send in-process messages to jobs_mngr goroutine, instead use HTTP REST API of gf_images
		Service_host_port_str:   gf_images_service_host_port_str,
		Templates_dir_paths_map: pConfig.Templates_paths_map,
	}
	
	gf_publisher_lib.Init_service(gfSoloHTTPmux,
		gf_images_runtime_info,
		pRuntimeSys)

	//-------------
	// GF_TAGGER
	gf_tagger_lib.Init_service(pConfig.Templates_paths_map,
		jobs_mngr_ch,
		gfSoloHTTPmux,
		pRuntimeSys)

	//-------------
	// GF_ML
	gf_ml_lib.InitService(gfSoloHTTPmux, pRuntimeSys)

	//-------------
	// GF_WEB3_MONITOR
	
	web3Config := &gf_eth_core.GF_config{
		AlchemyAPIkeyStr: pConfig.AlchemyAPIkeyStr,
	}
	gf_web3_lib.InitService(gfSoloHTTPmux,
		web3Config,
		pRuntimeSys)

	//-------------
	// METRICS - start prometheus metrics endpoint, and get core_metrics
	coreMetrics := gf_core.MetricsInit("/metrics", portMetricsInt)
	pRuntimeSys.Metrics = coreMetrics
	
	//-------------
	// SERVER_INIT - blocking
	gf_rpc_lib.ServerInitWithMux(portInt, gfSoloHTTPmux)
}

//-------------------------------------------------
func Runtime__get(pConfig_path_str string,
	p_external_plugins *gf_core.External_plugins,
	p_log_fun          func(string, string)) (*gf_core.Runtime_sys, *GF_config, error) {

	// CONFIG
	config_dir_path_str := path.Dir(pConfig_path_str)  // "./../config/"
	config_name_str     := path.Base(pConfig_path_str) // "gf_solo"
	
	config, err := Config__init(config_dir_path_str, config_name_str)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to load config")
		return nil, nil, err
	}

	//--------------------
	// SENTRY - ERROR_REPORTING
	if config.Sentry_endpoint_str != "" {

		sentry_endpoint_str := config.Sentry_endpoint_str
		sentry_samplerate_f := 1.0
		sentry_trace_handlers_map := map[string]bool{
			
		}
		err := gf_core.Error__init_sentry(sentry_endpoint_str,
			sentry_trace_handlers_map,
			sentry_samplerate_f)
		if err != nil {
			panic(err)
		}

		defer sentry.Flush(2 * time.Second)
	}

	//--------------------
	// RUNTIME_SYS
	runtimeSys := &gf_core.Runtime_sys{
		Service_name_str: "gf_solo",
		Log_fun:          p_log_fun,

		// SENTRY - enable it for error reporting
		Errors_send_to_sentry_bool: true,

		// EXTERNAL_PLUGINS
		External_plugins: p_external_plugins,
	}
	
	//--------------------
	// MONGODB
	mongodb_host_str := config.Mongodb_host_str
	mongodb_url_str  := fmt.Sprintf("mongodb://%s", mongodb_host_str)
	fmt.Printf("mongodb_host    - %s\n", mongodb_host_str)
	fmt.Printf("mongodb_db_name - %s\n", config.Mongodb_db_name_str)

	mongodb_db, _, gf_err := gf_core.Mongo__connect_new(mongodb_url_str,
		config.Mongodb_db_name_str,
		nil,
		runtimeSys)
	if gf_err != nil {
		return nil, nil, gf_err.Error
	}

	runtimeSys.Mongo_db   = mongodb_db
	runtimeSys.Mongo_coll = mongodb_db.Collection("data_symphony")
	fmt.Printf("mongodb connected...\n")

	//--------------------
	return runtimeSys, config, nil
}