/*
GloFlow media management/publishing system
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

package gf_publisher_lib

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib"
)
//---------------------------------------------------
func process_external_images(p_post *Post,
	p_gf_images_service_host_port_str string, //"http://127.0.0.1:2060"
	p_runtime_sys                     *gf_core.Runtime_sys) (string, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_images.process_external_images()")

	//-------------------	
	//POST ELEMENTS IMAGES

	post_elements_images_urls_lst              := []string{}
	post_elements_images_origin_pages_urls_str := []string{}
	post_elements_map                          := map[string]*Post_element{}

	for _,post_element := range p_post.Post_elements_lst {
		if post_element.Type_str == "image" {
			image_url_str                             := post_element.Extern_url_str
			source_page_url_str                       := post_element.Source_page_url_str
			post_elements_images_urls_lst              = append(post_elements_images_urls_lst,             image_url_str)
			post_elements_images_origin_pages_urls_str = append(post_elements_images_origin_pages_urls_str,source_page_url_str)
			post_elements_map[image_url_str]           = post_element
		}
	}
	//-------------------
	image_job_client_type_str := "gf_publisher"
	running_job_id_str, outputs_lst, gf_err := gf_images_lib.Client__dispatch_process_extern_images(post_elements_images_urls_lst, //p_input_images_urls_lst
		post_elements_images_origin_pages_urls_str,       //p_input_images_origin_pages_urls_str
		image_job_client_type_str,
		p_gf_images_service_host_port_str,
		p_runtime_sys)

	if gf_err != nil {
		return "", gf_err
	}
	//-------------------
	image_ids_lst := []string{}
	for _,output := range outputs_lst {
		gf_images__output_img_source_url_str := output.Image_source_url_str

		if _,ok := post_elements_map[gf_images__output_img_source_url_str]; !ok {
			gf_err := gf_core.Error__create(fmt.Sprintf("gf_images_lib client returned results for unknown image url - "+gf_images__output_img_source_url_str),
				"verify__invalid_value_error",
				&map[string]interface{}{"gf_images__output_img_source_url_str":gf_images__output_img_source_url_str,},
				nil, "gf_publisher_lib", p_runtime_sys)
			return "", gf_err
		}

		//--------------------
		//IMPORTANT!! - this is IMAGE_JOB EXPECTED_OUTPUT - these image url's will be resolved
		//              at a later time when the job completes (job is a long-running process)
		post_element                             := post_elements_map[output.Image_source_url_str]
		post_element.Image_id_str                 = output.Image_id_str
		post_element.Img_thumbnail_small_url_str  = output.Thumbnail_small_relative_url_str
		post_element.Img_thumbnail_medium_url_str = output.Thumbnail_medium_relative_url_str
		post_element.Img_thumbnail_large_url_str  = output.Thumbnail_large_relative_url_str
		//--------------------

		image_ids_lst = append(image_ids_lst,output.Image_id_str)
	}

	//IMPORTANT!! - list of all images in this post
	p_post.Images_ids_lst = image_ids_lst
	//----------------
	//POST THUMBNAIL
	//IMPORTANT!! - first image in the list of images supplied for the post, is also used as the post thumbnail
	first_image_url_str     := outputs_lst[0].Image_source_url_str
	first_post_element      := post_elements_map[first_image_url_str]
	post_thumbnail_str      := first_post_element.Img_thumbnail_small_url_str
	p_post.Thumbnail_url_str = post_thumbnail_str

	p_runtime_sys.Log_fun("INFO","post_thumbnail_str - "+post_thumbnail_str)
	//----------------
	//persists the newly updated post (some of its post_elements have been updated
	//in the initiation of image post_elements)
	gf_err = DB__update_post(p_post, p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}
	//----------------

	return running_job_id_str, nil
}