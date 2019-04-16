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

package gf_crawl_core

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)
//--------------------------------------------------
func images_s3__stage__store_images(p_crawler_name_str string,
	p_page_imgs__pipeline_infos_lst []*gf_page_img__pipeline_info,
	p_origin_page_url_str           string,
	p_s3_bucket_name_str            string,
	p_runtime                       *Gf_crawler_runtime,
	p_runtime_sys                   *gf_core.Runtime_sys) []*gf_page_img__pipeline_info {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_s3.images_s3__stage__store_images")

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")
	fmt.Println("IMAGES__GET_IN_PAGE    - STAGE - s3_store_images")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> -------------------------")

	for _, page_img__pinfo := range p_page_imgs__pipeline_infos_lst {

		//IMPORTANT!! - skip failed images
		if page_img__pinfo.gf_error != nil {
			continue
		}

		//IMPORTANT!! - skip images that have already been processed (and is in the DB)
		if page_img__pinfo.exists_bool {
			continue
		}

		//IMPORTANT!! - check image is not flagged as a NSFV image
		if page_img__pinfo.nsfv_bool {
			continue
		}

		//------------------
		//IMPORTANT!! - only store/persist if they are valid (of the right dimensions) or
		//              if they're a GIF (all GIF's are stored/persisted,
		//              even if they determined to be NSFV for some reason).

		if page_img__pinfo.page_img.Img_ext_str == "gif" || page_img__pinfo.page_img.Valid_for_usage_bool {

			gf_err := image_s3__upload(page_img__pinfo.page_img,
				page_img__pinfo.local_file_path_str,
				page_img__pinfo.thumbs,
				p_s3_bucket_name_str,
				p_runtime,
				p_runtime_sys)

			if gf_err != nil {
				t := "image_s3_upload__failed"
				m := "failed s3 uploading of image with img_url_str - "+page_img__pinfo.page_img.Url_str
				Create_error_and_event(t,m,map[string]interface{}{"origin_page_url_str":p_origin_page_url_str,}, page_img__pinfo.page_img.Url_str, p_crawler_name_str,
					gf_err, p_runtime, p_runtime_sys)
				page_img__pinfo.gf_error = gf_err
				continue //IMPORTANT!! - if an image processing fails, continue to the next image, dont abort
			}
		}
		//------------------
	}

	return p_page_imgs__pipeline_infos_lst
}

//--------------------------------------------------
func image_s3__upload(p_image *Gf_crawler_page_img,
	p_local_image_file_path_str string,
	p_image_thumbs              *gf_images_utils.Gf_image_thumbs,
	p_s3_bucket_name_str        string,
	p_runtime                   *Gf_crawler_runtime,
	p_runtime_sys               *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_s3.image_s3__upload()")


	gf_err := gf_images_utils.S3__store_gf_image(p_local_image_file_path_str,
		p_image_thumbs,
		p_s3_bucket_name_str,
		p_runtime.S3_info,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//------------------
	p_image.S3_stored_bool = true
	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"t":        "crawler_page_img",
			"hash_str": p_image.Hash_str,
		},
		bson.M{
			"$set": bson.M{"s3_stored_bool":true},
		})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update an crawler_page_img s3_stored flag by its hash",
			"mongodb_update_error",
			&map[string]interface{}{"image_hash_str": p_image.Hash_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return gf_err
	}
	//------------------

	return nil
}