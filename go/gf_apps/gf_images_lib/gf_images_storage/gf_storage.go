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

package gf_images_storage

import (
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type GFimageStorageOpDef struct {
	SourceFilePathStr string
	TargetFilePathStr string
	S3bucketNameStr   string
}

type GFimageStorage struct {
	TypeStr string // "local_fs" | "s3" | "ipfs"
	S3      *GFstorageS3
	IPFS    *GFstorageIPFS
}

type GFstorageS3 struct {
	Info *gf_core.GFs3Info
}

type GFstorageIPFS struct {

}

//---------------------------------------------------
// IPFS
func Init(pTypeStr string,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimageStorage, *gf_core.GFerror) {


	var S3storage *GFstorageS3
	if pTypeStr == "s3" {

		// get new S3 client, and get AWS creds from environment
		S3info, gfErr := gf_core.S3init("", "", "", pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
		S3storage = &GFstorageS3{
			Info: S3info,
		}
	}



	storage := &GFimageStorage{
		TypeStr: pTypeStr,
		S3:      S3storage,
	}

	return storage, nil


}

//---------------------------------------------------
// FILE_PUT
func FilePut(pOpDef *GFimageStorageOpDef,
	pStorage    *GFimageStorage,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	if pStorage.TypeStr == "s3" {
		_, gfErr := gf_core.S3uploadFile(pOpDef.SourceFilePathStr,
			pOpDef.TargetFilePathStr,
			pOpDef.S3bucketNameStr,
			pStorage.S3.Info,
			pRuntimeSys)
			
		if gfErr != nil {
			return gfErr
		}
	}

	return nil

}

//---------------------------------------------------
// FILE_GET
func FileGet(pStorage *GFimageStorage,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GF_error {




	if pStorage.TypeStr == "s3" {
		
	}


	return nil

}