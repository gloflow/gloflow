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
	ImageSourceFilePathStr      string
	ImageTargetLocalFilePathStr string
	S3bucketNameStr             string
}

type GFputFromLocalOpDef struct {
	ImageSourceLocalFilePathStr string
	ImageTargetFilePathStr      string
	S3bucketNameStr             string
}

type GFcopyOpDef struct {
	ImageSourceFilePathStr    string
	SourceFileS3bucketNameStr string
	ImageTargetFilePathStr    string
	TargetFileS3bucketNameStr string
}

type GFgeneratePresignedURLopDef struct {
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

	switch pStorage.TypeStr {

	// LOCAL
	case "local":

		sourceFileLocalPathStr := pOpDef.ImageSourceLocalFilePathStr
		targetFileLocalPathStr := pOpDef.ImageTargetFilePathStr

		// local file get operation copies the image file from the local image storage dir
		// to some desired working dir path 
		gfErr := gf_core.FileCopy(sourceFileLocalPathStr,
			targetFileLocalPathStr,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

	// S3
	case "s3":
		_, gfErr := gf_core.S3uploadFile(pOpDef.ImageSourceLocalFilePathStr,
			pOpDef.ImageTargetFilePathStr,
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

	switch pStorage.TypeStr {

	// LOCAL
	case "local":
		
		sourceFileLocalPathStr := pOpDef.ImageSourceFilePathStr
		targetFileLocalPathStr := pOpDef.ImageTargetFilePathStr

		// local file get operation copies the image file from the local image storage dir
		// to some desired working dir path 
		gfErr := gf_core.FileCopy(sourceFileLocalPathStr,
			targetFileLocalPathStr,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

	// S3
	case "s3":
		gfErr := gf_core.S3copyFile(pOpDef.SourceFileS3bucketNameStr,
			pOpDef.ImageSourceFilePathStr,
			pOpDef.TargetFileS3bucketNameStr,
			pOpDef.ImageTargetFilePathStr,
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
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	switch pStorage.TypeStr {

	// LOCAL
	case "local":
		
		sourceFileLocalPathStr := pOpDef.ImageSourceFilePathStr
		targetFileLocalPathStr := pOpDef.ImageTargetLocalFilePathStr

		// local file get operation copies the image file from the local image storage dir
		// to some desired working dir path 
		gfErr := gf_core.FileCopy(sourceFileLocalPathStr,
			targetFileLocalPathStr,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	
	// S3
	case "s3":
		gfErr := gf_core.S3getFile(pOpDef.ImageSourceFilePathStr,
			pOpDef.ImageTargetLocalFilePathStr,
			pOpDef.S3bucketNameStr,
			pStorage.S3.Info,
			pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

	case "ipfs":
		fmt.Println("ipfs")
	}

	return nil
}

//---------------------------------------------------
func FileGeneratePresignedURL(pOpDef *GFgeneratePresignedURLopDef,
	pStorage    *GFimageStorage,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {


	var presignedURLstr string

	switch pStorage.TypeStr {

	// s3 is the only file storage for which generating presigned URL's makes sense
	case "s3":
		S3presignedURLstr, gfErr := gf_core.S3generatePresignedUploadURL( pOpDef.TargetFilePathStr,
			pOpDef.TargetFileS3bucketNameStr,
			pStorage.S3.Info,
			pRuntimeSys)
		if gfErr != nil {
			return "", gfErr
		}

		presignedURLstr = S3presignedURLstr
		
		
	}

	return presignedURLstr, nil
}