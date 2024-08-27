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

package gf_images_lib

import (
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_service"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_flows"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_solo_service"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func InitSoloService(pHTTPmux *http.ServeMux,
	pServiceInfo *gf_images_core.GFserviceInfo,
	pConfig      *gf_images_core.GFconfig,
	pRuntimeSys  *gf_core.RuntimeSys) gf_images_jobs_core.JobsMngr {

	//-------------
	// METRICS
	metrics := gf_images_core.MetricsCreate("gf_images")

	//-------------
	// DB_INDEXES
	// IMPORTANT!! - make sure mongo has indexes build for relevant queries
	gf_images_service.DBmongoIndexInit(pRuntimeSys)

	//-------------
	// S3
	// REMOVE!! - usage of AWS creds here, they should be discovered
	//            by the AWS client from the environment.
	s3Info, gfErr := gf_aws.S3init(// pServiceInfo.AWS_access_key_id_str,
		// pServiceInfo.AWS_secret_access_key_str,
		// pServiceInfo.AWS_token_str,
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
		pServiceInfo.MediaDomainStr,
		pConfig,
		imageStorage,
		s3Info,
		metrics,
		pRuntimeSys)

	gf_images_jobs.InitHandlers(pHTTPmux, jobsMngrCh, pRuntimeSys)

	//-------------
	// IMAGE_FLOWS

	gfErr = gf_images_flows.Init(pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	// flows__templates_dir_path_str := pServiceInfo.Templates_dir_paths_map["flows_str"]
	gfErr = gf_images_flows.InitHandlers(pServiceInfo.AuthSubsystemTypeStr,
		pServiceInfo.AuthLoginURLstr,
		pServiceInfo.KeyServer,
		pHTTPmux,
		pServiceInfo.TemplatesPathsMap,
		jobsMngrCh,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//-------------
	// HANDLERS
	gfErr = gf_images_solo_service.InitHandlers(pServiceInfo.AuthSubsystemTypeStr,
		pServiceInfo.AuthLoginURLstr,
		pServiceInfo.KeyServer,
		pHTTPmux,
		jobsMngrCh,
		pConfig,
		pServiceInfo,
		pServiceInfo.MediaDomainStr,
		imageStorage,
		s3Info,
		metrics,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//------------------------
	
	// STATIC_FILE SERVING
	staticFilesURLbaseStr := "/images"
	localDirPathStr       := "./static"
	gf_core.HTTPinitStaticServingWithMux(staticFilesURLbaseStr,
		localDirPathStr,
		pHTTPmux,
		pRuntimeSys)

	//------------------------

	return jobsMngrCh
}