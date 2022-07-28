/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_images_service

import (
	"fmt"
	"time"
	"context"
	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
// GF_image_upload_info struct represents a single image upload sequence.
// It is both stored in the DB and returned to the initiating client in JSON form.
// It contains the ID of the future gf_image that will be created in the system to represent
// the image that the client is wanting to upload.
type GFimageUploadInfo struct {
	Id                primitive.ObjectID       `json:"-"                      bson:"_id,omitempty"`
	Tstr              string                   `json:"-"                      bson:"t"` // "img_upload_info"
	CreationUNIXtimeF float64                  `json:"creation_unix_time_f"   bson:"creation_unix_time_f"`
	NameStr           string                   `json:"name_str"               bson:"name_str"`
	UploadImageIDstr  gf_images_core.GFimageID `json:"upload_gf_image_id_str" bson:"upload_gf_image_id_str"`
	S3filePathStr     string                   `json:"-"                      bson:"s3_file_path_str"` // internal data, dont send to clients
	FlowsNamesLst     []string                 `json:"flows_names_lst"        bson:"flows_names_lst"`
	ClientTypeStr     string                   `json:"-"                      bson:"client_type_str"`  // internal data, dont send to clients
	PresignedURLstr   string                   `json:"presigned_url_str"      bson:"presigned_url_str"`
}

//---------------------------------------------------
// UploadInit initializes an file upload process.
// This will create a pre-signed S3 URL for the caller of this function to use
// for uploading of content to GF.
func UploadInit(pImageNameStr string,
	pImageFormatStr string,
	pFlowsNamesLst  []string,
	pClientTypeStr  string,
	pStorage        *gf_images_storage.GFimageStorage,
	pS3info         *gf_core.GFs3Info,
	pConfig         *gf_images_core.GFconfig,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFimageUploadInfo, *gf_core.GFerror) {
	
	//------------------
	// CHECK_IMAGE_FORMAT
	normalizedFormatStr, ok := gf_images_core.Image__check_image_format(pImageFormatStr, pRuntimeSys)
	if !ok {
		gfErr := gf_core.Error__create("image format is invalid that specified for image thats being prepared for uploading via upload__init",
			"verify__invalid_value_error",
			map[string]interface{}{"image_format_str": pImageFormatStr,},
			nil, "gf_images_lib", pRuntimeSys)
		return nil, gfErr
	}

	//------------------
	// GF_IMAGE_ID
	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	imagePathStr      := pImageNameStr
	uploadImageIDstr  := gf_images_core.Image_ID__create(imagePathStr, normalizedFormatStr, pRuntimeSys)

	

	
	


	//------------------
	// PRESIGN_URL

	var presignedURLstr   string
	var targetFilePathStr string
	var gfErr             *gf_core.GFerror

	// NEW_STORAGE
	if  pConfig.UseNewStorageEngineBool {

		imageFileNameStr := gf_images_core.ImageGetFilepathFromID(uploadImageIDstr, normalizedFormatStr)
		targetFilePathStr = imageFileNameStr

		op := &gf_images_storage.GFgeneratePresignedURLopDef{
			TargetFilePathStr: targetFilePathStr,
		}
		if pStorage.TypeStr == "s3" {
			op.TargetFileS3bucketNameStr = pStorage.S3.UploadsSourceS3bucketNameStr
		}
		presignedURLstr, gfErr = gf_images_storage.FileGeneratePresignedURL(op, pStorage, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

	} else {
		// LEGACY

		S3filePathStr := gf_images_core.S3getImageFilepath(uploadImageIDstr,
			normalizedFormatStr,
			pRuntimeSys)
		targetFilePathStr = S3filePathStr

		S3bucketNameStr := pConfig.Uploaded_images_s3_bucket_str // "gf--uploaded--img"

		pRuntimeSys.Log_fun("INFO", fmt.Sprintf("S3 generating presigned_url - bucket (%s) - file (%s)",
			S3bucketNameStr,
			S3filePathStr))

		presignedURLstr, gfErr = gf_core.S3generatePresignedUploadURL(S3filePathStr,
			S3bucketNameStr,
			pS3info,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}

	pRuntimeSys.Log_fun("INFO", fmt.Sprintf("S3 presigned URL - %s", presignedURLstr))
	
	//------------------
	
	uploadInfo := &GFimageUploadInfo{
		Tstr:              "img_upload_info",
		CreationUNIXtimeF: creationUNIXtimeF,
		NameStr:           pImageNameStr,
		UploadImageIDstr:  uploadImageIDstr,
		S3filePathStr:     targetFilePathStr,
		FlowsNamesLst:     pFlowsNamesLst,
		ClientTypeStr:     pClientTypeStr,
		PresignedURLstr:   presignedURLstr,
	}

	//------------------
	// DB
	gfErr = UploadDBputInfo(uploadInfo, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------
	return uploadInfo, nil
}

//---------------------------------------------------
// UploadComplete completes the image file upload sequence.
// It is run after the initialization stage, and after the client/caller conducts
// the upload operation.
func UploadComplete(pUploadImageIDstr gf_images_core.GF_image_id,
	pMetaMap    map[string]interface{},
	pJobsMngrCh chan gf_images_jobs_core.JobMsg,
	pRuntimeSys *gf_core.RuntimeSys) (*gf_images_jobs_core.GFjobRunning, *gf_core.GFerror) {
	
	// DB
	uploadInfo, gfErr := Upload_db__get_info(pUploadImageIDstr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	imageToProcessLst := []gf_images_jobs_core.GF_image_uploaded_to_process{
		{
			GF_image_id_str:  pUploadImageIDstr,
			S3_file_path_str: uploadInfo.S3filePathStr,
			Meta_map:         pMetaMap,
		},
	}

	// JOB
	runningJob, gfErr := gf_images_jobs_client.Run_uploaded_imgs(uploadInfo.ClientTypeStr,
		imageToProcessLst,
		uploadInfo.FlowsNamesLst,
		pJobsMngrCh,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}
	
	return runningJob, nil
}

//---------------------------------------------------
// DB
//---------------------------------------------------
func UploadDBputInfo(pUploadInfo *GFimageUploadInfo,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	ctx         := context.Background()
	collNameStr := "gf_images_upload_info"
	gfErr       := gf_core.Mongo__insert(pUploadInfo,
		collNameStr,
		map[string]interface{}{
			"upload_image_id_str": pUploadInfo.UploadImageIDstr,
			"caller_err_msg_str":  "failed to update/upsert image upload_info into the DB",
		},
		ctx,
		pRuntimeSys)
	
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------
func Upload_db__get_info(p_upload_gf_image_id_str gf_images_core.Gf_image_id,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimageUploadInfo, *gf_core.GF_error) {

	ctx := context.Background()

	var upload_info GFimageUploadInfo
	err := pRuntimeSys.Mongo_db.Collection("gf_images_upload_info").FindOne(ctx, bson.M{
			"t":                      "img_upload_info",
			"upload_gf_image_id_str": p_upload_gf_image_id_str,
		}).Decode(&upload_info)

	if err == mongo.ErrNoDocuments {
		gf_err := gf_core.Mongo__handle_error("image_upload_info does not exist in mongodb",
			"mongodb_not_found_error",
			map[string]interface{}{"upload_gf_image_id_str": p_upload_gf_image_id_str,},
			err, "gf_images_lib", pRuntimeSys)
		return nil, gf_err
	}

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get image_upload_info from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"upload_gf_image_id_str": p_upload_gf_image_id_str,},
			err, "gf_images_lib", pRuntimeSys)
		return nil, gf_err
	}

	return &upload_info, nil
}

//---------------------------------------------------
func Upload_db__put_image_upload_info(pImageUploadInfo *GFimageUploadInfo,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx         := context.Background()
	collNameStr := pRuntimeSys.Mongo_coll.Name()
	gfErr       := gf_core.Mongo__insert(pImageUploadInfo,
		collNameStr,
		map[string]interface{}{
			"upload_gf_image_id_str": pImageUploadInfo.UploadImageIDstr,
			"caller_err_msg_str":     "failed to update/upsert image upload_info into the DB",
		},
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}