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
	"fmt"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// OP_DEFs
//---------------------------------------------------
type GFgetOpDef struct {
	ImageLocalFilePathStr string
	TargetFilePathStr     string
	S3bucketNameStr       string
}

type GFputFromLocalOpDef struct {
	SourceLocalFilePathStr string
	TargetFilePathStr      string
	S3bucketNameStr        string
}

type GFcopyOpDef struct {
	SourceFilePathStr         string
	SourceFileS3bucketNameStr string
	TargetFilePathStr         string
	TargetFileS3bucketNameStr string
}

//---------------------------------------------------
// VAR
//---------------------------------------------------
type GFimageStorage struct {
	TypeStr string // "local" | "s3" | "ipfs"
	Local   *GFstorageLocal
	S3      *GFstorageS3
	IPFS    *GFstorageIPFS
}

type GFstorageLocal struct {
	ThumbsDirPathStr        string
	UploadsSourceDirPathStr string
	UploadsTargetDirPathStr string
	ExternImagesDirPathStr  string
}

type GFstorageS3 struct {
	Info                         *gf_core.GFs3Info
	ThumbsS3bucketNameStr        string
	UploadsSourceS3bucketNameStr string
	UploadsTargetS3bucketNameStr string
	ExternImagesS3bucketNameStr  string
}

type GFstorageIPFS struct {
	Shell *ipfs.Shell
}

type GFimageStorageConfig struct {
	TypesToProvisionLst []string // list of storage types to initialize
	IPFSnodeHostStr     string

	// LOCAL_DIRS
	ThumbsDirPathStr        string
	UploadsSourceDirPathStr string
	UploadsTargetDirPathStr string
	ExternImagesDirPathStr  string

	// S3_BUCKETS
	ThumbsS3bucketNameStr        string
	UploadsSourceS3bucketNameStr string
	UploadsTargetS3bucketNameStr string
	ExternImagesS3bucketNameStr  string
}

//---------------------------------------------------
// INIT
func Init(pConfig *GFimageStorageConfig,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimageStorage, *gf_core.GFerror) {

	storage := &GFimageStorage{}
	for _, storageTypeStr := range pConfig.TypesToProvisionLst {
		
		switch storageTypeStr {

		//-------------
		// LOCAL
		case "local":

			localStorage := &GFstorageLocal{
				ThumbsDirPathStr:        pConfig.ThumbsDirPathStr,
				UploadsSourceDirPathStr: pConfig.UploadsSourceDirPathStr,
				UploadsTargetDirPathStr: pConfig.UploadsTargetDirPathStr,
				ExternImagesDirPathStr:  pConfig.ExternImagesDirPathStr,
			}
			storage.Local = localStorage

		//-------------
		// S3
		case "s3":

			// get new S3 client, and get AWS creds from environment
			S3info, gfErr := gf_core.S3init("", "", "", pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}
			S3storage := &GFstorageS3{
				Info:                         S3info,
				ThumbsS3bucketNameStr:        pConfig.ThumbsS3bucketNameStr,
				UploadsSourceS3bucketNameStr: pConfig.UploadsSourceS3bucketNameStr,
				UploadsTargetS3bucketNameStr: pConfig.UploadsTargetS3bucketNameStr,
			}
			storage.S3 = S3storage

		//-------------
		// IPFS
		case "ipfs":

			ipfsShell, gfErr := gf_core.IPFSinit(pConfig.IPFSnodeHostStr,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			IPFSstorage := &GFstorageIPFS{
				Shell: ipfsShell,
			}
			storage.IPFS = IPFSstorage

		//-------------
		}
	}

	return storage, nil
}

//---------------------------------------------------
// FILE_PUT
func FilePutFromLocal(pOpDef *GFputFromLocalOpDef,
	pStorage    *GFimageStorage,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	if pStorage.TypeStr == "s3" {
		_, gfErr := gf_core.S3uploadFile(pOpDef.SourceLocalFilePathStr,
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
// FILE_COPY
func FileCopy(pOpDef *GFcopyOpDef,
	pStorage    *GFimageStorage,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	if pStorage.TypeStr == "s3" {
		gfErr := gf_core.S3copyFile(pOpDef.SourceFileS3bucketNameStr,
			pOpDef.SourceFilePathStr,
			pOpDef.TargetFileS3bucketNameStr,
			pOpDef.TargetFilePathStr,
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
func FileGet(pOpDef *GFgetOpDef,
	pStorage    *GFimageStorage,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GFerror {




	switch pStorage.TypeStr {
	case "s3":
		fmt.Println("s3")

	case "ipfs":
		fmt.Println("ipfs")
	}

	



	return nil

}