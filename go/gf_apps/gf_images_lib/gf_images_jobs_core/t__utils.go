/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_images_jobs_core

import (
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
)

//---------------------------------------------------
func TgetJobsMngr(pImagesStoreLocalDirPathStr string,
	pImagesThumbnailsStoreLocalDirPathStr string,
	pVideoStoreLocalDirPathStr            string,
	pMediaDomainStr                       string,
	pRuntimeSys                           *gf_core.RuntimeSys) JobsMngr {


	imagesConfig := &gf_images_core.GFconfig{
		UseNewStorageEngineBool: true,
	}

	//---------------------
	// IMAGE_STORAGE
	storageConfig := &gf_images_storage.GFimageStorageConfig{
		TypesToProvisionLst: []string{"local"},
	}
	
	imageStorage, gfErr := gf_images_storage.Init(storageConfig, pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//---------------------
	// JOBS_MNGR
	jobsMngr := JobsMngrInit(pImagesStoreLocalDirPathStr,
		pImagesThumbnailsStoreLocalDirPathStr,
		pVideoStoreLocalDirPathStr,
		pMediaDomainStr,
		nil, // pLifecycleCallbacks
		imagesConfig,
		imageStorage,
		nil, // pS3info
		pRuntimeSys)

	//---------------------
	
	return jobsMngr
}