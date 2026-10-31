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
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_service"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_flows"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_gif_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_solo_service"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

/*
InitSoloService initializes the complete gf_images service, independent of gf_solo.
it sets up:
- DB (SQL tables, Mongo indexes)
- S3 client and image storage
- images jobs manager and its HTTP handlers
- images flows (JSON API handlers)
- GIF HTTP handlers
- images HTTP handlers (get, share, upload, upload_metrics, browser calc results, health)

its a pure API server: no HTML/template rendering and no static files serving.
the caller is responsible for starting the HTTP server on the mux.
*/
func InitSoloService(pHTTPmux *http.ServeMux,
	pServiceInfo *gf_images_core.GFserviceInfo,
	pConfig      *gf_images_core.GFconfig,
	pRuntimeSys  *gf_core.RuntimeSys) gf_images_jobs_core.JobsMngr {

	//-------------
	// METRICS
	metricsGroupNameStr := "gf_images_solo__rpc_handlers"
	metrics := gf_images_core.MetricsCreate("gf_images")

	//-------------
	// DB
	ctx := context.Background()

	// SQL - images tables
	gfErr := gf_images_core.DBsqlCreateTables(ctx, pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	// DB_INDEXES
	// IMPORTANT!! - make sure mongo has indexes build for relevant queries
	gf_images_service.DBmongoIndexInit(pRuntimeSys)

	//-------------
	// S3
	// REMOVE!! - usage of AWS creds here, they should be discovered
	//            by the AWS client from the environment.
	s3Info, gfErr := gf_aws.S3init(pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//---------------------
	// IMAGE_STORAGE

	// S3 is the primary storage. IPFS is provisioned only if a node host is configured.
	storageTypesLst := []string{"s3"}
	if pConfig.IPFSnodeHostStr != "" {
		storageTypesLst = append(storageTypesLst, "ipfs")
	}

	storageConfig := &gf_images_storage.GFimageStorageConfig{
		TypesToProvisionLst:          storageTypesLst,
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

	// extern HTTP handlers for the images jobs manager
	gf_images_jobs.InitHandlers(pHTTPmux, jobsMngrCh, pRuntimeSys)

	//-------------
	// IMAGE_FLOWS

	gfErr = gf_images_flows.Init(pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	flowsHandlersLst, gfErr := gf_images_flows.InitHandlersAPI(pServiceInfo.AuthSubsystemTypeStr,
		pServiceInfo.AuthLoginURLstr,
		pServiceInfo.KeyServer,
		pHTTPmux,
		jobsMngrCh,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	// register flows handlers, which are defined in the new style (handler info structs).
	gf_rpc_lib.CreateHandlersHTTP(flowsHandlersLst,
		pHTTPmux,
		pServiceInfo.AuthSubsystemTypeStr,
		pServiceInfo.AuthLoginURLstr,
		pServiceInfo.KeyServer,
		metricsGroupNameStr,
		pRuntimeSys)

	//-------------
	// GIF
	gfErr = gf_gif_lib.InitHandlers(pHTTPmux, pRuntimeSys)
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

	return jobsMngrCh
}