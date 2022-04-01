/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

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

package gf_images_lib

import (
	"fmt"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_gif_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_service"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_image_editor"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_flows"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Init_service(p_http_mux *http.ServeMux,
	pServiceInfo  *gf_images_core.GFserviceInfo,
	p_config      *gf_images_core.GF_config,
	p_runtime_sys *gf_core.Runtime_sys) gf_images_jobs_core.Jobs_mngr {

	//-------------
	// DB_INDEXES
	// IMPORTANT!! - make sure mongo has indexes build for relevant queries
	gf_images_service.DB_index__init(p_runtime_sys)

	//-------------
	// S3
	s3_info, gf_err := gf_core.S3__init(pServiceInfo.AWS_access_key_id_str,
		pServiceInfo.AWS_secret_access_key_str,
		pServiceInfo.AWS_token_str,
		p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
	// JOBS_MANAGER

	jobs_mngr_ch := gf_images_jobs.Init(pServiceInfo.Images_store_local_dir_path_str,
		pServiceInfo.Images_thumbnails_store_local_dir_path_str,
		pServiceInfo.Media_domain_str,
		p_config,
		s3_info,
		p_runtime_sys)

	/*jobs_mngr_ch := gf_images_jobs.Jobs_mngr__init(pServiceInfo.Images_store_local_dir_path_str,
		pServiceInfo.Images_thumbnails_store_local_dir_path_str,
		// pServiceInfo.Images_main_s3_bucket_name_str,
		pServiceInfo.Media_domain_str,
		p_config,
		s3_info,
		p_runtime_sys)*/

	//-------------
	// IMAGE_FLOWS

	// flows__templates_dir_path_str := pServiceInfo.Templates_dir_paths_map["flows_str"]
	gf_err = gf_images_flows.Init_handlers(pServiceInfo.AuthLoginURLstr,
		p_http_mux,
		pServiceInfo.Templates_paths_map,
		jobs_mngr_ch,
		p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
	// GIF
	gf_err = gf_gif_lib.Gif__init_handlers(p_http_mux, p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//-------------
	// IMAGE_EDITOR
	gf_image_editor.Init_handlers(p_http_mux, p_runtime_sys)
	
	//-------------
	/*gf_gif_lib.Init_img_to_gif_migration(*p_images_store_local_dir_path_str,
							*p_images_main_s3_bucket_name_str,
							s3_client,
							s3_uploader, //s3_client,
							mongodb_coll,
							p_log_fun)*/
	
	//-------------
	// JOBS_MANAGER
	gf_images_jobs.Jobs_mngr__init_handlers(p_http_mux, jobs_mngr_ch, p_runtime_sys)

	//-------------
	// HANDLERS
	gf_err = gf_images_service.InitHandlers(pServiceInfo.AuthLoginURLstr,
		p_http_mux,
		jobs_mngr_ch,
		p_config,
		pServiceInfo.Media_domain_str,
		s3_info,
		p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}

	//------------------------
	// DASHBOARD SERVING
	static_files__url_base_str := "/images"
	local_dir_path_str         := "./static"
	gf_core.HTTP__init_static_serving_with_mux(static_files__url_base_str,
		local_dir_path_str,
		p_http_mux,
		p_runtime_sys)

	//------------------------

	return jobs_mngr_ch
}

//-------------------------------------------------
// Run_service runs/starts the gf_images service in the same process as where its being called.
// An HTTP servr is started and listens on a supplied port.
// DB(MongoDB) connection is established as well.
// S3 client is initialized as a target file-system for image files.
func Run_service(p_mux *http.ServeMux,
	pServiceInfo   *gf_images_core.GFserviceInfo,
	p_init_done_ch chan bool,
	p_log_fun      func(string, string)) {
	p_log_fun("FUN_ENTER", "gf_images_service.Run_service()")

	p_log_fun("INFO", "")
	p_log_fun("INFO", " >>>>>>>>>>> STARTING GF_IMAGES SERVICE")
	p_log_fun("INFO", "")
	logo_str := `.           ..         .         
      .         .            .          .       .
            .         ..xxxxxxxxxx....               .       .             .
    .             MWMWMWWMWMWMWMWMWMWMWMWMW                       .
              IIIIMWMWMWMWMWMWMWMWMWMWMWMWMWMttii:        .           .
 .      IIYVVXMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWxx...         .           .
     IWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMx..
   IIWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWNMWMWMWMWMWMWMWMWMWMWMWMWMx..        .
    ""MWMWMWMWMWM"""""""".  .:..   ."""""MWMWMWMWMWMWMWMWMWMWMWMWMWti.
 .     ""   .   .: . :. : .  . :.  .  . . .  """"MWMWMWMWMWMWMWMWMWMWMWMWMti=
        . .   : . :   .  .'.' '....xxxxx...,'. '   ' ."""YWMWMWMWMWMWMWMWMWMW+
     ; .  .  . : . .' :  . ..XXXXXXXXXXXXXXXXXXXXx.         . YWMWMWMWMWMWMW
.    .  .  .    . .   .  ..XXXXXXXXWWWWWWWWWWWWWWWWXXXX.  .     .     """""""
        ' :  : . : .  ...XXXXXWWW"   W88N88@888888WWWWWXX.   .   .       . .
   . ' .    . :   ...XXXXXXWWW"    M88Ng8GGGG5G888^8M "WMBX.          .   ..  :
         :     ..XXXXXXXXWWW"     M88a8WWRWWWMW8oo88M   WWMX.     .    :    .
           "XXXXXXXXXXXXWW"       WN8s88WWWWW  W8@@@8M    BMBRX.         .  : :
  .       XXXXXXXX=MMWW":  .      W8N888WWWWWWWW88888W      XRBRXX.  .       .
     ....  ""XXXXXMM::::. .        W8@889WWWWWM8@8N8W      . . :RRXx.    .
         .....'''  MMM::.:.  .      W888N89999888@8W      . . ::::"RXV    .  :
 .       ..'''''      MMMm::.  .      WW888N88888WW     .  . mmMMMMMRXx
      ..' .            ""MMmm .  .       WWWWWWW   . :. :,miMM"""  : ""    .
   .                .       ""MMMMmm . .  .  .   ._,mMMMM"""  :  ' .  :
               .                  ""MMMMMMMMMMMMM""" .  : . '   .        .
          .              .     .    .                      .         .
.                                         .          .         .`
	p_log_fun("INFO", logo_str)

	//-------------
	// RUNTIME_SYS
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_images",
		Log_fun:          p_log_fun,
	}

	mongo_db, _, gf_err := gf_core.Mongo__connect_new(pServiceInfo.Mongodb_host_str,
		pServiceInfo.Mongodb_db_name_str,
		nil,
		runtime_sys)
	if gf_err != nil {
		return
	}
	runtime_sys.Mongo_db   = mongo_db
	runtime_sys.Mongo_coll = mongo_db.Collection("data_symphony")

	//-------------
	// CONFIG

	gf_config, gf_err := gf_images_core.Config__get(pServiceInfo.Config_file_path_str, runtime_sys)
	if gf_err != nil {
		return
	}

	//------------------------
	// INIT
	Init_service(p_mux, pServiceInfo, gf_config, runtime_sys)

	//------------------------
	// IMPORTANT!! - signal to user that server in this goroutine is ready to start listening 
	if p_init_done_ch != nil {
		p_init_done_ch <- true
	}

	//----------------------

	runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	runtime_sys.Log_fun("INFO", fmt.Sprintf("STARTING HTTP SERVER - PORT - %s", pServiceInfo.Port_str))
	runtime_sys.Log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	http_err := http.ListenAndServe(fmt.Sprintf(":%s", pServiceInfo.Port_str), nil)
	if http_err != nil {
		msg_str := fmt.Sprintf("cant start listening on port - %s", pServiceInfo.Port_str)
		runtime_sys.Log_fun("ERROR", msg_str)
		runtime_sys.Log_fun("ERROR", fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}