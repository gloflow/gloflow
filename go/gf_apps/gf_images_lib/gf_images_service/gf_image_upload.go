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
	// "fmt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_identity/gf_policy"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_aws"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_storage"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_flows"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

// GFimageUploadInfo struct represents a single image upload sequence.
// It is both stored in the DB and returned to the initiating client in JSON form.
// It contains the ID of the future gf_image that will be created in the system to represent
// the image that the client is wanting to upload.
type GFimageUploadInfo struct {
	Id                primitive.ObjectID       `json:"-"                      bson:"_id,omitempty"`
	Tstr              string                   `json:"-"                      bson:"t"` // "img_upload_info"
	CreationUNIXtimeF float64                  `json:"creation_unix_time_f"   bson:"creation_unix_time_f"`
	NameStr           string                   `json:"name_str"               bson:"name_str"`
	ImageIDstr        gf_images_core.GFimageID `json:"upload_gf_image_id_str" bson:"upload_gf_image_id_str"`
	S3filePathStr     string                   `json:"-"                      bson:"s3_file_path_str"` // internal data, dont send to clients
	FlowsNamesLst     []string                 `json:"flows_names_lst"        bson:"flows_names_lst"`
	ClientTypeStr     string                   `json:"-"                      bson:"client_type_str"`  // internal data, dont send to clients
	PresignedURLstr   string                   `json:"presigned_url_str"      bson:"presigned_url_str"`
	UserID            gf_core.GF_ID            `json:"user_id_str"            bson:"user_id_str"`      // ID of the user starting this upload
}

type GFimageUploadMetrics struct {
	Id                 primitive.ObjectID       `bson:"_id,omitempty"`
	CreationUNIXtimeF  float64                  `bson:"creation_unix_time_f"`
	ImageIDstr         gf_images_core.GFimageID `bson:"upload_gf_image_id_str"`
	ClientTypeStr      string                   `bson:"client_type_str"`
	UploadClientDurationSecF         float64    `bson:"upload_client_duration_sec_f"`
	UploadClientTransferDurationSecF float64       `bson:"upload_client_transfer_duration_sec_f"`
	UserID                           gf_core.GF_ID `bson:"user_id_str"` // ID of the user starting this upload
}

//---------------------------------------------------

// UploadInit initializes an file upload process.
// This will create a pre-signed S3 URL for the caller of this function to use
// for uploading of content to GF.
func UploadInit(pImageNameStr string,
	pImageFormatStr string,
	pFlowsNamesLst  []string,
	pClientTypeStr  string,
	pUserID         gf_core.GF_ID,
	pStorage        *gf_images_storage.GFimageStorage,
	pS3info         *gf_aws.GFs3Info,
	pConfig         *gf_images_core.GFconfig,
	pServiceInfo    *gf_images_core.GFserviceInfo,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFimageUploadInfo, *gf_core.GFerror) {
	
	//------------------
	/*
	CREATE_FLOWS - check if flows to which this image is being added exist,
		and create if its missing.
		if flow doesnt exist it is assigned to this user. if it does exist
		nothing happens, and subsequent policy verification will check if
		user is allowed to add images to the flow.
	*/

	gfErr := gf_images_flows.CreateIfMissingWithPolicy(pFlowsNamesLst,
		pUserID,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------
	// POLICY_VERIFY - verify user is allowed to upload an image into the specified flows.
	//                 raises error if policy rejects the op.
	
	opStr := gf_policy.GF_POLICY_OP__FLOW_ADD_IMG
	gfErr = gf_images_flows.VerifyPolicy(opStr,
		pFlowsNamesLst,
		pUserID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------
	// CHECK_IMAGE_FORMAT

	ok := gf_images_core.CheckImageFormat(pImageFormatStr, pRuntimeSys)
	if !ok {
		gfErr := gf_core.ErrorCreate("image format is invalid that specified for image thats being prepared for uploading via upload__init",
			"verify__invalid_value_error",
			map[string]interface{}{
				"image_format": pImageFormatStr,
				"user_id":      pUserID},
			nil, "gf_images_service", pRuntimeSys)
		return nil, gfErr
	}

	normalizedFormatStr := gf_images_core.NormalizeImageFormat(pImageFormatStr)

	//------------------
	// GF_IMAGE_ID
	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	imageURIstr       := pImageNameStr
	uploadImageIDstr  := gf_images_core.CreateImageID(imageURIstr, pRuntimeSys)

	//------------------
	// PRESIGN_URL

	var presignedURLstr   string
	var targetFilePathStr string

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

		pRuntimeSys.LogNewFun("DEBUG", "S3 generating presigned url", map[string]interface{}{
			"s3_bucket_name_str": S3bucketNameStr,
			"s3_file_path_str":   S3filePathStr,
		})

		presignedURLstr, gfErr = gf_aws.S3generatePresignedUploadURL(S3filePathStr,
			S3bucketNameStr,
			pS3info,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}

	pRuntimeSys.LogNewFun("DEBUG", "S3 presigned URL generated", map[string]interface{}{
		"presigned_url_str": presignedURLstr,})
	
	//------------------
	
	uploadInfo := &GFimageUploadInfo{
		Tstr:              "img_upload_info",
		CreationUNIXtimeF: creationUNIXtimeF,
		NameStr:           pImageNameStr,
		ImageIDstr:        uploadImageIDstr,
		S3filePathStr:     targetFilePathStr,
		FlowsNamesLst:     pFlowsNamesLst,
		ClientTypeStr:     pClientTypeStr,
		PresignedURLstr:   presignedURLstr,
		UserID:            pUserID,
	}

	//------------------
	// DB
	gfErr = dbPutUploadInfo(uploadInfo, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------
	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		eventMetaMap := map[string]interface{}{
			"user_id": pUserID,
		}
		gf_events.EmitApp(gf_images_core.GF_EVENT_APP__IMAGE_UPLOAD_INIT,
			eventMetaMap,
			pRuntimeSys.AppNameStr,
			pUserID,
			pCtx,
			pRuntimeSys)
	}

	//------------------

	return uploadInfo, nil
}

//---------------------------------------------------

// completes the image file upload sequence.
// It is run after the initialization stage, and after the client/caller conducts
// the upload operation.
func UploadComplete(pUploadImageIDstr gf_images_core.GFimageID,
	pFlowsNamesLst []string,
	pMetaMap       map[string]interface{},
	pUserID        gf_core.GF_ID,
	pJobsMngrCh    chan gf_images_jobs_core.JobMsg,
	pServiceInfo   *gf_images_core.GFserviceInfo,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (*gf_images_jobs_core.GFjobRunning, *gf_core.GFerror) {
	
	//------------------
	// POLICY_VERIFY - verify user is allowed to upload an image into the specified flows.
	//                 raises error if policy rejects the op.
	//
	opStr := gf_policy.GF_POLICY_OP__FLOW_ADD_IMG
	gfErr := gf_images_flows.VerifyPolicy(opStr,
		pFlowsNamesLst,
		pUserID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------

	// DB
	uploadInfo, gfErr := dbGetUploadInfo(pUploadImageIDstr, pUserID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	
	imageToProcessLst := []gf_images_jobs_core.GFimageUploadedToProcess{
		{
			GFimageIDstr:  pUploadImageIDstr,
			S3filePathStr: uploadInfo.S3filePathStr,
			MetaMap:       pMetaMap,
		},
	}

	// JOB
	runningJob, gfErr := gf_images_jobs_client.RunUploadedImages(uploadInfo.ClientTypeStr,
		imageToProcessLst,
		uploadInfo.FlowsNamesLst,
		pUserID,
		pJobsMngrCh,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}
	
	//------------------------
	// EVENT
	if pServiceInfo.EnableEventsAppBool {
		eventMetaMap := map[string]interface{}{
			"user_id": pUserID,
		}
		gf_events.EmitApp(gf_images_core.GF_EVENT_APP__IMAGE_UPLOAD_COMPLETE,
			eventMetaMap,
			pRuntimeSys.AppNameStr,
			pUserID,
			pCtx,
			pRuntimeSys)
	}

	//------------------------

	return runningJob, nil
}

//---------------------------------------------------

func UploadMetricsCreate(pUploadImageIDstr gf_images_core.GFimageID,
	pClientTypeStr  string,
	pMetricsDataMap map[string]interface{},
	pUserID         gf_core.GF_ID,
	pMetrics        *gf_images_core.GFmetrics,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {
	
	// VALIDATE
	if _, ok := pMetricsDataMap["upload_client_duration_sec_f"]; !ok {
		gfErr := gf_core.MongoHandleError("image upload metrics data is missing the 'upload_client_duration_sec_f' key",
			"verify__missing_key_error",
			map[string]interface{}{
				"upload_gf_image_id_str": pUploadImageIDstr,
				"user_id":                pUserID},
			nil, "gf_images_service", pRuntimeSys)
		return gfErr
	}

	if _, ok := pMetricsDataMap["upload_client_transfer_duration_sec_f"]; !ok {
		gfErr := gf_core.MongoHandleError("image upload metrics data is missing the 'upload_client_transfer_duration_sec_f' key",
			"verify__missing_key_error",
			map[string]interface{}{
				"upload_gf_image_id_str": pUploadImageIDstr,
				"user_id":                pUserID},
			nil, "gf_images_service", pRuntimeSys)
		return gfErr
	}

	uploadClientDurationSecF := pMetricsDataMap["upload_client_duration_sec_f"].(float64)
	uploadClientTransferDurationSecF := pMetricsDataMap["upload_client_transfer_duration_sec_f"].(float64)


	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	metrics := &GFimageUploadMetrics {
		CreationUNIXtimeF:  creationUNIXtimeF,
		ImageIDstr:         pUploadImageIDstr,
		ClientTypeStr:      pClientTypeStr,
		UploadClientDurationSecF:         uploadClientDurationSecF,
		UploadClientTransferDurationSecF: uploadClientTransferDurationSecF,
		UserID:                           pUserID,
	}

	// DB
	gfErr := dbUploadMetricsCreate(metrics,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	if pMetrics != nil {
		pMetrics.ImageUploadClientDurationGauge.Set(uploadClientDurationSecF)
		pMetrics.ImageUploadClientTransferDurationGauge.Set(uploadClientTransferDurationSecF)
	}

	return nil
}

//---------------------------------------------------
// DB
//---------------------------------------------------

func dbUploadMetricsCreate(pUploadMetrics *GFimageUploadMetrics,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {


	collNameStr := "gf_images_upload_metrics"
	gfErr       := gf_core.MongoInsert(pUploadMetrics,
		collNameStr,
		map[string]interface{}{
			"upload_image_id_str": pUploadMetrics.ImageIDstr,
			"caller_err_msg_str":  "failed to update/upsert image upload_info into the DB",
		},
		pCtx,
		pRuntimeSys)
	
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func dbPutUploadInfo(pUploadInfo *GFimageUploadInfo,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	collNameStr := "gf_images_upload_info"
	gfErr       := gf_core.MongoInsert(pUploadInfo,
		collNameStr,
		map[string]interface{}{
			"user_id":             pUploadInfo.UserID,
			"upload_image_id_str": pUploadInfo.ImageIDstr,
			"caller_err_msg_str":  "failed to update/upsert image upload_info into the DB",
		},
		pCtx,
		pRuntimeSys)
	
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func dbGetUploadInfo(pUploadImageIDstr gf_images_core.GFimageID,
	pUserID     gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimageUploadInfo, *gf_core.GFerror) {

	var uploadInfo GFimageUploadInfo
	err := pRuntimeSys.Mongo_db.Collection("gf_images_upload_info").FindOne(pCtx, bson.M{
			"t":                      "img_upload_info",
			"upload_gf_image_id_str": pUploadImageIDstr,
		}).Decode(&uploadInfo)

	if err == mongo.ErrNoDocuments {
		gfErr := gf_core.MongoHandleError("image_upload_info does not exist in mongodb",
			"mongodb_not_found_error",
			map[string]interface{}{
				"user_id":                pUserID,
				"upload_gf_image_id_str": pUploadImageIDstr,
			},
			err, "gf_images_service", pRuntimeSys)
		return nil, gfErr
	}

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get image_upload_info from mongodb",
			"mongodb_find_error",
			map[string]interface{}{
				"user_id":                pUserID,
				"upload_gf_image_id_str": pUploadImageIDstr,
			},
			err, "gf_images_service", pRuntimeSys)
		return nil, gfErr
	}

	if uploadInfo.UserID != pUserID {
		gfErr := gf_core.MongoHandleError("invalid user_id fetching upload_info from DB, user starting upload is not the same thats fetching it.",
			"user_incorrect",
			map[string]interface{}{
				"user_id":                pUserID,
				"upload_user_id":         uploadInfo.UserID,
				"upload_gf_image_id_str": pUploadImageIDstr,
			
			},
			err, "gf_images_service", pRuntimeSys)
		return nil, gfErr
	}

	return &uploadInfo, nil
}