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

package gf_publisher_lib

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
)

//---------------------------------------------------
type Gf_images_client_result struct {
	image_ids_lst      []string
	running_job_id_str string
	post_thumbnail_str string
}

//---------------------------------------------------
func process_external_images(p_post *Gf_post,
	p_gf_images_runtime_info *Gf_images_extern_runtime_info,
	p_runtime_sys            *gf_core.Runtime_sys) (string, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_images.process_external_images()")

	//-------------------	
	//POST ELEMENTS IMAGES

	post_elements_images_urls_lst              := []string{}
	post_elements_images_origin_pages_urls_str := []string{}
	post_elements_map                          := map[string]*Gf_post_element{}

	for _,post_element := range p_post.Post_elements_lst {
		if post_element.Type_str == "image" {
			image_url_str                             := post_element.Extern_url_str
			origin_page_url_str                       := post_element.Origin_page_url_str
			post_elements_images_urls_lst              = append(post_elements_images_urls_lst,              image_url_str)
			post_elements_images_origin_pages_urls_str = append(post_elements_images_origin_pages_urls_str, origin_page_url_str)
			post_elements_map[image_url_str]           = post_element
		}
	}
	//-------------------
	image_job_client_type_str := "gf_publisher"

	var result *Gf_images_client_result
	var gf_err *gf_core.Gf_error

	//HTTP
	if p_gf_images_runtime_info.Jobs_mngr == nil {

		/*//FIX!! - when GF is compiled into a unified binary there shouldnt be an HTTP call to process images, this instead should be sending
		//        a message to the gf_images job manager that is running in the same process and gf_publisher here.
		running_job_id_str, outputs_lst, gf_err := gf_images_lib.Client__dispatch_process_extern_images(post_elements_images_urls_lst, //p_input_images_urls_lst
			post_elements_images_origin_pages_urls_str,       //p_input_images_origin_pages_urls_str
			image_job_client_type_str,
			p_gf_images_runtime_info.service_host_port_str,
			p_runtime_sys)

		if gf_err != nil {
			return "", gf_err
		}

		final_running_job_id_str = running_job_id_str
		final_outputs_lst        = outputs_lst*/



		result, gf_err = process_external_images__via_http(post_elements_map,
			post_elements_images_urls_lst,
			post_elements_images_origin_pages_urls_str,
			image_job_client_type_str,
			p_gf_images_runtime_info.Service_host_port_str,
			p_runtime_sys)
		if gf_err != nil {
			return "", nil	
		}

	//IN_PROCESS - for unified binary where both gf_publisher and gf_images are running in the same process
	} else {
		
		result, gf_err = process_external_images__in_process(post_elements_map,
			post_elements_images_urls_lst,
			post_elements_images_origin_pages_urls_str,
			image_job_client_type_str,
			p_gf_images_runtime_info.Jobs_mngr,
			p_runtime_sys)
		if gf_err != nil {
			return "", nil	
		}
	}
	//-------------------
	/*image_ids_lst := []string{}
	for _, output := range outputs_lst {
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

		image_ids_lst = append(image_ids_lst, output.Image_id_str)
	}*/

	//IMPORTANT!! - list of all images in this post
	p_post.Images_ids_lst = result.image_ids_lst
	//----------------
	/*//POST THUMBNAIL
	//IMPORTANT!! - first image in the list of images supplied for the post, is also used as the post thumbnail
	first_image_url_str     := outputs_lst[0].Image_source_url_str
	first_post_element      := post_elements_map[first_image_url_str]
	post_thumbnail_str      := first_post_element.Img_thumbnail_small_url_str*/

	p_post.Thumbnail_url_str = result.post_thumbnail_str
	p_runtime_sys.Log_fun("INFO", fmt.Sprintf("post_thumbnail_str - %s",result.post_thumbnail_str))
	//----------------
	//persists the newly updated post (some of its post_elements have been updated
	//in the initiation of image post_elements)
	gf_err = DB__update_post(p_post, p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}
	//----------------

	return result.running_job_id_str, nil
}

//---------------------------------------------------
func process_external_images__via_http(p_post_elements_map map[string]*Gf_post_element,
	p_post_elements_images_urls_lst              []string,
	p_post_elements_images_origin_pages_urls_str []string,
	p_image_job_client_type_str                  string,
	p_gf_images_service_host_port_str            string,
	p_runtime_sys                                *gf_core.Runtime_sys) (*Gf_images_client_result, *gf_core.Gf_error) {

	//--------------------
	//HTTP
	running_job_id_str, outputs_lst, gf_err := gf_images_lib.Client__dispatch_process_extern_images(p_post_elements_images_urls_lst, //p_input_images_urls_lst
		p_post_elements_images_origin_pages_urls_str,       //p_input_images_origin_pages_urls_str
		p_image_job_client_type_str,
		p_gf_images_service_host_port_str,
		p_runtime_sys)

	if gf_err != nil {
		return nil, gf_err
	}
	//--------------------

	image_ids_lst := []string{}
	for _, output := range outputs_lst {
		gf_images__output_img_source_url_str := output.Image_source_url_str

		if _,ok := p_post_elements_map[gf_images__output_img_source_url_str]; !ok {
			gf_err := gf_core.Error__create(fmt.Sprintf("gf_images_lib client returned results for unknown image url - "+gf_images__output_img_source_url_str),
				"verify__invalid_value_error",
				map[string]interface{}{"gf_images__output_img_source_url_str":gf_images__output_img_source_url_str,},
				nil, "gf_publisher_lib", p_runtime_sys)
			return nil, gf_err
		}

		//--------------------
		//IMPORTANT!! - this is IMAGE_JOB EXPECTED_OUTPUT - these image url's will be resolved
		//              at a later time when the job completes (job is a long-running process)
		post_element                             := p_post_elements_map[output.Image_source_url_str]
		post_element.Image_id_str                 = output.Image_id_str
		post_element.Img_thumbnail_small_url_str  = output.Thumbnail_small_relative_url_str
		post_element.Img_thumbnail_medium_url_str = output.Thumbnail_medium_relative_url_str
		post_element.Img_thumbnail_large_url_str  = output.Thumbnail_large_relative_url_str
		image_ids_lst = append(image_ids_lst, output.Image_id_str)
		//--------------------
	}
	//--------------------
	//POST THUMBNAIL
	//IMPORTANT!! - first image in the list of images supplied for the post, is also used as the post thumbnail
	first_image_url_str := outputs_lst[0].Image_source_url_str
	first_post_element  := p_post_elements_map[first_image_url_str]
	post_thumbnail_str  := first_post_element.Img_thumbnail_small_url_str
	//--------------------

	result := &Gf_images_client_result{
		image_ids_lst:      image_ids_lst,
		running_job_id_str: running_job_id_str,
		post_thumbnail_str: post_thumbnail_str,
	}

	return result, nil
}

//---------------------------------------------------
func process_external_images__in_process(p_post_elements_map map[string]*Gf_post_element,
	p_post_elements_images_urls_lst              []string,
	p_post_elements_images_origin_pages_urls_str []string,
	p_image_job_client_type_str                  string,
	p_gf_images_jobs_mngr                        gf_images_jobs.Jobs_mngr,
	p_runtime_sys                                *gf_core.Runtime_sys) (*Gf_images_client_result, *gf_core.Gf_error) {

	//ADD!! - accept this flows_names argument from http arguments, not hardcoded as is here
	flows_names_lst := []string{"general",}

	images_to_process_lst := []gf_images_jobs.Image_to_process{}
	for i, image_url_str := range p_post_elements_images_urls_lst {
		
		origin_page_url_str := p_post_elements_images_origin_pages_urls_str[i]
		img_to_process      := gf_images_jobs.Image_to_process{
			Source_url_str:     image_url_str,
			Origin_page_url_str:origin_page_url_str,
		}
		images_to_process_lst = append(images_to_process_lst, img_to_process)
	}

	//--------------------
	//IN_PROCESS
	running_job, outputs_lst, gf_err := gf_images_jobs.Job__start(p_image_job_client_type_str,
		images_to_process_lst,
		flows_names_lst,
		p_gf_images_jobs_mngr,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	//--------------------
	
	image_ids_lst := []string{}
	for _, output := range outputs_lst {
		gf_images__output_img_source_url_str := output.Image_source_url_str

		if _,ok := p_post_elements_map[gf_images__output_img_source_url_str]; !ok {
			gf_err := gf_core.Error__create(fmt.Sprintf("gf_images_lib client returned results for unknown image url - "+gf_images__output_img_source_url_str),
				"verify__invalid_value_error",
				map[string]interface{}{"gf_images__output_img_source_url_str":gf_images__output_img_source_url_str,},
				nil, "gf_publisher_lib", p_runtime_sys)
			return nil, gf_err
		}

		//--------------------
		//IMPORTANT!! - this is IMAGE_JOB EXPECTED_OUTPUT - these image url's will be resolved
		//              at a later time when the job completes (job is a long-running process)
		post_element                             := p_post_elements_map[output.Image_source_url_str]
		post_element.Image_id_str                 = output.Image_id_str
		post_element.Img_thumbnail_small_url_str  = output.Thumbnail_small_relative_url_str
		post_element.Img_thumbnail_medium_url_str = output.Thumbnail_medium_relative_url_str
		post_element.Img_thumbnail_large_url_str  = output.Thumbnail_large_relative_url_str
		image_ids_lst = append(image_ids_lst, output.Image_id_str)
		//--------------------
	}
	//--------------------
	//POST THUMBNAIL
	//IMPORTANT!! - first image in the list of images supplied for the post, is also used as the post thumbnail
	first_image_url_str := outputs_lst[0].Image_source_url_str
	first_post_element  := p_post_elements_map[first_image_url_str]
	post_thumbnail_str  := first_post_element.Img_thumbnail_small_url_str
	//--------------------

	result := &Gf_images_client_result{
		image_ids_lst:      image_ids_lst,
		running_job_id_str: running_job.Id_str,
		post_thumbnail_str: post_thumbnail_str,
	}

	return result, nil
}