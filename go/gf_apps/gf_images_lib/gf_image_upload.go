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

package gf_images_lib

import (
	"fmt"
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
// Gf_image_upload_info struct represents a single image upload sequence.
// It is both stored in the DB and returned to the initiating client in JSON form.
// It contains the ID of the future gf_image that will be created in the system to represent
// the image that the client is wanting to upload.
type Gf_image_upload_info struct {
	Id                     bson.ObjectId               `json:"-"                      bson:"_id,omitempty"`
	T_str                  string                      `json:"-"                      bson:"t"` // "img_upload_info"
	Creation_unix_time_f   float64                     `json:"creation_unix_time_f"   bson:"creation_unix_time_f"`
	Name_str               string                      `json:"name_str"               bson:"name_str"`
	Upload_gf_image_id_str gf_images_utils.Gf_image_id `json:"upload_gf_image_id_str" bson:"upload_gf_image_id_str"`
	S3_file_path_str       string                      `json:"-"                      bson:"s3_file_path_str"` // internal data, dont send to clients
	Flows_names_lst        []string                    `json:"flows_names_lst"        bson:"flows_names_lst"`
	Client_type_str        string                      `json:"-"                      bson:"client_type_str"`  // internal data, dont send to clients
	Presigned_url_str      string                      `json:"presigned_url_str"      bson:"presigned_url_str"`
}

//---------------------------------------------------
// Upload__init initializes an file upload process.
// This will create a pre-signed S3 URL for the caller of this function to use
// for uploading of content to GF.
func Upload__init(p_image_name_str string,
	p_image_format_str string,
	p_flows_names_lst  []string,
	p_client_type_str  string,
	p_s3_info          *gf_core.Gf_s3_info,
	p_config           *gf_images_utils.Gf_config,
	p_runtime_sys      *gf_core.Runtime_sys) (*Gf_image_upload_info, *gf_core.Gf_error) {
	
	//------------------
	// CHECK_IMAGE_FORMAT
	normalized_format_str, ok := gf_images_utils.Image__check_image_format(p_image_format_str, p_runtime_sys)
	if !ok {
		gf_err := gf_core.Error__create("image format is invalid that specified for image thats being prepared for uploading via upload__init",
			"verify__invalid_value_error",
			map[string]interface{}{"image_format_str": p_image_format_str,},
			nil, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}

	//------------------
	// GF_IMAGE_ID
	creation_unix_time_f   := float64(time.Now().UnixNano())/1000000000.0
	image_path_str         := p_image_name_str
	upload_gf_image_id_str := gf_images_utils.Image_ID__create(image_path_str, normalized_format_str, p_runtime_sys)

	s3_file_path_str := gf_images_utils.S3__get_image_s3_filepath(upload_gf_image_id_str,
		normalized_format_str,
		p_runtime_sys)

	s3_bucket_name_str := p_config.Uploaded_images_s3_bucket_str // "gf--uploaded--img"

	//------------------
	// PRESIGN_URL
	p_runtime_sys.Log_fun("INFO", fmt.Sprintf("S3 generating presigned_url - bucket (%s) - file (%s)",
		s3_bucket_name_str,
		s3_file_path_str))

	presigned_url_str, gf_err := gf_core.S3__generate_presigned_url(s3_file_path_str,
		s3_bucket_name_str,
		p_s3_info,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	p_runtime_sys.Log_fun("INFO", fmt.Sprintf("S3 presigned URL - %s", presigned_url_str))

	//------------------
	
	upload_info := &Gf_image_upload_info{
		T_str:                  "img_upload_info",
		Creation_unix_time_f:   creation_unix_time_f,
		Name_str:               p_image_name_str,
		Upload_gf_image_id_str: upload_gf_image_id_str,
		S3_file_path_str:       s3_file_path_str,
		Flows_names_lst:        p_flows_names_lst,
		Client_type_str:        p_client_type_str,
		Presigned_url_str:      presigned_url_str,
	}


	// DB
	gf_err = Upload_db__put_info(upload_info, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}


	return upload_info, nil
}

//---------------------------------------------------
// Upload__complete completes the image file upload sequence.
// It is run after the initialization stage, and after the client/caller conducts
// the upload operation.
func Upload__complete(p_upload_gf_image_id_str gf_images_utils.Gf_image_id,
	p_jobs_mngr_ch    chan gf_images_jobs.Job_msg,
	p_s3_info         *gf_core.Gf_s3_info,
	p_runtime_sys     *gf_core.Runtime_sys) (*gf_images_jobs.Gf_running_job, *gf_core.Gf_error) {
	


	// DB
	upload_info, gf_err := Upload_db__get_info(p_upload_gf_image_id_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}


	image_to_process_lst := []gf_images_jobs.Gf_image_uploaded_to_process{
		gf_images_jobs.Gf_image_uploaded_to_process{
			Gf_image_id_str:  p_upload_gf_image_id_str,
			S3_file_path_str: upload_info.S3_file_path_str,
		},
	}

	running_job, gf_err := gf_images_jobs.Client__run_uploaded_imgs(upload_info.Client_type_str,
		image_to_process_lst,
		upload_info.Flows_names_lst,
		p_jobs_mngr_ch,
		p_runtime_sys)

	if gf_err != nil {
		return nil, gf_err
	}



	spew.Dump(running_job)
	


	return running_job, nil
}

//---------------------------------------------------
// DB
//---------------------------------------------------
func Upload_db__put_info(p_upload_info *Gf_image_upload_info,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {

	p_runtime_sys.Log_fun("INFO", "DB INSERT - img_upload_info")
	
	err := p_runtime_sys.Mongodb_db.C("gf_images_upload_info").Insert(p_upload_info)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update/upsert gf_image in a mongodb",
			"mongodb_insert_error",
			map[string]interface{}{"upload_gf_image_id_str": p_upload_info.Upload_gf_image_id_str,},
			err, "gf_images_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}

//---------------------------------------------------
func Upload_db__get_info(p_upload_gf_image_id_str gf_images_utils.Gf_image_id,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_image_upload_info, *gf_core.Gf_error) {

	var upload_info Gf_image_upload_info
	err := p_runtime_sys.Mongodb_db.C("gf_images_upload_info").Find(bson.M{
		"t":                      "img_upload_info",
		"upload_gf_image_id_str": p_upload_gf_image_id_str,
	}).One(&upload_info)

	if fmt.Sprint(err) == "not found" {
		gf_err := gf_core.Mongo__handle_error("image_upload_info does not exist in mongodb",
			"mongodb_not_found_error",
			map[string]interface{}{"upload_gf_image_id_str": p_upload_gf_image_id_str,},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get image_upload_info from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"upload_gf_image_id_str": p_upload_gf_image_id_str,},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}

	return &upload_info, nil
	
}
