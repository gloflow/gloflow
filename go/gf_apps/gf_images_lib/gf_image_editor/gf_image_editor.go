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

package gf_image_editor

import (
	"fmt"
	"os"
	"time"
	"context"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"image"
	"image/png"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------
type Gf_edited_image struct {
	Id                   primitive.ObjectID       `bson:"_id,omitempty"`
	Id_str               string                   `bson:"id_str"` 
	T_str                string                   `bson:"t"` // "img_edited"
	Creation_unix_time_f float64                  `bson:"creation_unix_time_f"`
	Source_image_id_str  gf_images_core.GFimageID `bson:"source_image_id_str"`
}

type Gf_edited_image__save__http_input struct {
	Title_str             string                   `json:"new_title_str"`         // title of the new edited_image
	Source_image_id_str   gf_images_core.GFimageID `json:"source_image_id_str"`   // id of the gf_image that has modification applied to it
	Source_flow_name_str  string                   `json:"source_flow_name_str"`  // which flow was the original image from
	Target_flow_name_str  string                   `json:"target_flow_name_str"`  // which flow the modified_image should be placed into
	Image_base64_data_str string                   `json:"image_base64_data_str"` // base64 encoded pixel data of the image
	Applied_filters_lst   []string                 `json:"applied_filters_lst"`   // list of filter names (in order) that were applied to the original image
	New_height_int        int                      `json:"new_height_int"`        // new dimensions in case of cropping/resizing
	New_width_int         int                      `json:"new_width_int"`         // new dimensions in case of cropping/resizing
}

type GFeditedImageProcessingInfo struct {
	png_image                 image.Image
	tmp_local_filepath_str    string
	image_origin_url_str      string
	image_origin_page_url_str string
	image_width_int           int
	image_height_int          int
}

//-------------------------------------------------
func save_edited_image__pipeline(p_handler_url_path_str string,
	p_req         *http.Request,
	p_resp        http.ResponseWriter, 
	pCtx          context.Context,
	pRuntimeSys   *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_image_editor.save_edited_image__pipeline()")

	//--------------------------
	// INPUT
	var input *Gf_edited_image__save__http_input
	body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
	err               := json.Unmarshal(body_bytes_lst, input)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to parse json edited_image_save http_input",
			"json_decode_error",
			map[string]interface{}{"handler_url_path_str": p_handler_url_path_str,},
			err, "gf_image_editor", pRuntimeSys)
		return gfErr
	}

	new_title_str       := input.Title_str
	source_image_id_str := input.Source_image_id_str

	//--------------------------
	// SAVE_BASE64_DATA_TO_FILE
	// IMPORTANT!! - save first, and then create a GFimages
	processingInfo, gfErr := saveEditedImage(source_image_id_str,
		input.Image_base64_data_str,
		pRuntimeSys)
	if err != nil {
		return gfErr
	}

	//--------------------------

	source_gf_image, gfErr := gf_images_core.DB__get_image(source_image_id_str, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	processingInfo.image_origin_url_str      = source_gf_image.Origin_url_str
	processingInfo.image_origin_page_url_str = source_gf_image.Origin_page_url_str

	gfErr = createImage(new_title_str,
		[]string{input.Target_flow_name_str,},
		processingInfo,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	//--------------------------

	return nil
}

//-------------------------------------------------
func saveEditedImage(p_source_image_id_str gf_images_core.GFimageID,
	p_image_base64_data_str string,
	pRuntimeSys             *gf_core.RuntimeSys) (*GFeditedImageProcessingInfo, *gf_core.GFerror) {
	
	//--------------------------
	// BASE64_DECODE

	image_byte_lst, err := base64.StdEncoding.DecodeString(p_image_base64_data_str)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to decode base64 string of image_data",
			"base64_decoding_error",
			map[string]interface{}{
				"source_image_id_str":   p_source_image_id_str,
				"image_base64_data_str": p_image_base64_data_str,
			},
			err, "gf_image_editor", pRuntimeSys)
		return nil, gfErr
	}

	//--------------------------
	// PNG

	image_reader   := bytes.NewReader(image_byte_lst)
	png_image, err := png.Decode(image_reader)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to encode png image_byte array while saving edited_image",
			"png_encoding_error",
			map[string]interface{}{
				"source_image_id_str":   p_source_image_id_str,
				"image_base64_data_str": p_image_base64_data_str,
			},
			err, "gf_image_editor", pRuntimeSys)
		return nil, gfErr
	}

	//--------------------------
	// FILE

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	tmp_local_filepath_str := fmt.Sprintf("/%f.png", creationUNIXtimeF)

	// FILE_CREATE
	file, err := os.Create(tmp_local_filepath_str)
	if err != nil {
		gfErr := gf_core.ErrorCreate("OS failed to create a file to save edited_image to FS",
			"file_create_error",
			map[string]interface{}{
				"source_image_id_str":    p_source_image_id_str,
				"tmp_local_filepath_str": tmp_local_filepath_str,
			},
			err, "gf_image_editor", pRuntimeSys)
		return nil, gfErr
	}
	defer file.Close()

	// FILE_WRITE_IMAGE
	err = png.Encode(file, png_image)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to encode png image_byte array while saving GIF frame to FS",
			"png_encoding_error",
			map[string]interface{}{"tmp_local_filepath_str": tmp_local_filepath_str,},
			err, "gf_image_editor", pRuntimeSys)
		return nil, gfErr
	}

	// FILE_SYNC
	if err := file.Sync(); err != nil {
		gfErr := gf_core.ErrorCreate("failed to decode jpen image_byte array while saving edited_image",
			"file_sync_error",
			map[string]interface{}{
				"source_image_id_str":    p_source_image_id_str,
				"tmp_local_filepath_str": tmp_local_filepath_str,
			},
			err, "gf_image_editor", pRuntimeSys)
		return nil, gfErr
	}

	//--------------------------
	// IMAGE_DIMENSIONS

	imageWidthInt, imageHeightInt := gf_images_core.GetImageDimensionsFromImage(png_image,
		pRuntimeSys)
	
	//--------------------------

	processingInfo := GFeditedImageProcessingInfo{
		png_image:              png_image,
		tmp_local_filepath_str: tmp_local_filepath_str,
		image_width_int:        imageWidthInt,
		image_height_int:       imageHeightInt,
	}

	return &processingInfo, nil
}

//-------------------------------------------------
func createImage(p_new_title_str string,
	p_images_flows_names_lst []string,
	pProcessingInfo          *GFeditedImageProcessingInfo,
	pCtx                     context.Context,
	pRuntimeSys              *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_image_editor.create_gf_image()")
	
	
	image_client_type_str                := "gf_image_editor" // IMPORTANT!! - since gf_image_editor is creating the image,
	image_format_str                     := "png"
	local_thumbnails_target_dir_path_str := "."
	small_thumb_max_size_px_int          := 200
	medium_thumb_max_size_px_int         := 400
	large_thumb_max_size_px_int          := 600
	
	//--------------------------
	// GF_IMAGE_ID
	imageURIstr := pProcessingInfo.tmp_local_filepath_str
	imageIDstr  := gf_images_core.CreateImageID(imageURIstr,
		pRuntimeSys)

	//--------------------------
	// THUMBNAILS
	imageThumbs, gfErr := gf_images_core.CreateThumbnails(pProcessingInfo.png_image,
		imageIDstr,
		image_format_str,
		local_thumbnails_target_dir_path_str,
		small_thumb_max_size_px_int,
		medium_thumb_max_size_px_int,
		large_thumb_max_size_px_int,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}

	//--------------------------

	gfImageInfo := &gf_images_core.GFimageNewInfo{
		IDstr:                          imageIDstr,
		Title_str:                      p_new_title_str,
		Flows_names_lst:                p_images_flows_names_lst,
		Image_client_type_str:          image_client_type_str,
		Origin_url_str:                 pProcessingInfo.image_origin_url_str,
		Origin_page_url_str:            pProcessingInfo.image_origin_page_url_str,
		Original_file_internal_uri_str: pProcessingInfo.tmp_local_filepath_str, // image_local_file_path_str,
		Thumbnail_small_url_str:        imageThumbs.Small_relative_url_str,
		Thumbnail_medium_url_str:       imageThumbs.Medium_relative_url_str,
		Thumbnail_large_url_str:        imageThumbs.Large_relative_url_str,
		Format_str:                     image_format_str,
		Width_int:                      pProcessingInfo.image_width_int,
		Height_int:                     pProcessingInfo.image_height_int,
	}

	// IMPORTANT!! - creates a GF_Image struct and stores it in the DB.
	//               every GIF in the system has its GF_Gif DB struct and GF_Image DB struct.
	//               these two structs are related by origin_url
	_, gfErr = gf_images_core.ImageCreateNew(gfImageInfo, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}