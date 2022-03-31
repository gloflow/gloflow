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

package gf_image_storage

import (
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type GFstorage struct {
	TypeStr string // "local_fs" | "s3" | "ipfs"
	S3      GFstorageS3
}

type GFstorageS3 struct {
	BucketNameStr string
	Info          *gf_core.GF_s3_info
}

type GFstorageIPFS struct {

}

//---------------------------------------------------
// IPFS
func IPFSinit() {

}

//---------------------------------------------------
// FILE_PUT
func FilePut(pSourceFilePathStr string,
	pTargetFilePathStr string,
	pStorage           *GF_image_storage,
	pRuntimeSys        *gf_core.Runtime_sys) *gf_core.GF_error {

	if pStorage.TypeStr == "s3" {
		gfErr := gf_core.S3__upload_file(pSourceFilePathStr,
			pTargetFilePathStr,
			pStorage.S3.Bucket_name_str,
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
func FileGet(pStorage *GF_image_storage,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GF_error {




	if pStorage.TypeStr == "s3" {
		
	}


	return nil

}