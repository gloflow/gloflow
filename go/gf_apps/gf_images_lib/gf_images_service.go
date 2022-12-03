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
	"strconv"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_gif_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_service"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_image_editor"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_flows"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func InitService(pHTTPmux *http.ServeMux,
	pServiceInfo *gf_images_core.GFserviceInfo,
	pConfig      *gf_images_core.GFconfig,
	pRuntimeSys  *gf_core.RuntimeSys) gf_images_jobs_core.JobsMngr {

	//-------------
	// METRICS
	metrics := gf_images_core.MetricsCreate("gf_images")

	//-------------
	// DB_INDEXES
	// IMPORTANT!! - make sure mongo has indexes build for relevant queries
	gf_images_service.DBindexInit(pRuntimeSys)

	//-------------
	// S3
	s3Info, gfErr := gf_aws.S3init(pServiceInfo.AWS_access_key_id_str,
		pServiceInfo.AWS_secret_access_key_str,
		pServiceInfo.AWS_token_str,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//---------------------
	// IMAGE_STORAGE

	storageConfig := &gf_images_storage.GFimageStorageConfig{
		TypesToProvisionLst:          []string{"s3", "ipfs"},
		IPFSnodeHostStr:              pConfig.IPFSnodeHostStr,
		ThumbsS3bucketNameStr:        pConfig.ImagesFlowToS3bucketMap["general"],
		UploadsSourceS3bucketNameStr: pConfig.Uploaded_images_s3_bucket_str,
		UploadsTargetS3bucketNameStr: pConfig.ImagesFlowToS3bucketMap["general"],
		ExternImagesS3bucketNameStr:  pConfig.ImagesFlowToS3bucketMap["general"],
	}

	imageStorage, gfErr := gf_images_storage.Init(storageConfig, pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//-------------
	// JOBS_MANAGER

	jobsMngrCh := gf_images_jobs.Init(pServiceInfo.ImagesStoreLocalDirPathStr,
		pServiceInfo.ImagesThumbnailsStoreLocalDirPathStr,
		pServiceInfo.VideoStoreLocalDirPathStr,
		pServiceInfo.Media_domain_str,
		pConfig,
		imageStorage,
		s3Info,
		pRuntimeSys)

	//-------------
	// IMAGE_FLOWS

	// flows__templates_dir_path_str := pServiceInfo.Templates_dir_paths_map["flows_str"]
	gfErr = gf_images_flows.InitHandlers(pServiceInfo.AuthLoginURLstr,
		pHTTPmux,
		pServiceInfo.Templates_paths_map,
		jobsMngrCh,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//-------------
	// GIF
	gfErr = gf_gif_lib.InitHandlers(pHTTPmux, pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	/*gf_gif_lib.Init_img_to_gif_migration(*p_images_store_local_dir_path_str,
		*p_images_main_s3_bucket_name_str,
		s3_client,
		s3_uploader, //s3_client,
		mongodb_coll,
		pLogFun)*/

	//-------------
	// IMAGE_EDITOR
	gf_image_editor.InitHandlers(pHTTPmux, pRuntimeSys)
	
	//-------------
	// JOBS_MANAGER
	gf_images_jobs.InitHandlers(pHTTPmux, jobsMngrCh, pRuntimeSys)

	//-------------
	// HANDLERS
	gfErr = gf_images_service.InitHandlers(pServiceInfo.AuthLoginURLstr,
		pHTTPmux,
		jobsMngrCh,
		pConfig,
		pServiceInfo.Media_domain_str,
		imageStorage,
		s3Info,
		metrics,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//------------------------
	// DASHBOARD SERVING
	static_files__url_base_str := "/images"
	local_dir_path_str         := "./static"
	gf_core.HTTPinitStaticServingWithMux(static_files__url_base_str,
		local_dir_path_str,
		pHTTPmux,
		pRuntimeSys)

	//------------------------

	return jobsMngrCh
}

//-------------------------------------------------

// Run_service runs/starts the gf_images service in the same process as where its being called.
// An HTTP servr is started and listens on a supplied port.
// DB(MongoDB) connection is established as well.
// S3 client is initialized as a target file-system for image files.
func RunService(pHTTPmux *http.ServeMux,
	pServiceInfo   *gf_images_core.GFserviceInfo,
	p_init_done_ch chan bool,
	pLogFun      func(string, string)) {
	pLogFun("FUN_ENTER", "gf_images_service.RunService()")

	pLogFun("INFO", "")
	pLogFun("INFO", " >>>>>>>>>>> STARTING GF_IMAGES SERVICE")
	pLogFun("INFO", "")
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
	pLogFun("INFO", logo_str)

	//-------------
	// RUNTIME_SYS
	
	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_images",
		LogFun:           pLogFun,
	}

	mongoDB, _, gfErr := gf_core.MongoConnectNew(pServiceInfo.Mongodb_host_str,
		pServiceInfo.Mongodb_db_name_str,
		nil,
		runtimeSys)
	if gfErr != nil {
		return
	}
	runtimeSys.Mongo_db   = mongoDB
	runtimeSys.Mongo_coll = mongoDB.Collection("data_symphony")

	//-------------
	// CONFIG
	config, gfErr := gf_images_core.ConfigGet(pServiceInfo.Config_file_path_str,
		pServiceInfo.UseNewStorageEngineBool,
		pServiceInfo.IPFSnodeHostStr,
		runtimeSys)
	if gfErr != nil {
		return
	}

	//------------------------
	// INIT
	InitService(pHTTPmux, pServiceInfo, config, runtimeSys)

	//------------------------
	// IMPORTANT!! - signal to user that server in this goroutine is ready to start listening 
	if p_init_done_ch != nil {
		p_init_done_ch <- true
	}

	//----------------------

	portInt, err := strconv.Atoi(pServiceInfo.Port_str)
	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	// SERVER_INIT - blocking
	gf_rpc_lib.ServerInitWithMux(portInt, pHTTPmux)
}