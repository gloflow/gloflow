package gf_publisher_lib

import (
	"errors"
	"github.com/globalsign/mgo"
	"apps/gf_images_lib"
)
//---------------------------------------------------
func process_external_images(p_post *Post,
	p_gf_images_service_host_port_str *string, //"http://127.0.0.1:2060"
	p_mongodb_coll                    *mgo.Collection,
	p_log_fun                         func(string,string)) (*string,error) {
	p_log_fun("FUN_ENTER","gf_post_images.process_external_images()")

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
	running_job_id_str,outputs_lst,err := gf_images_lib.Client__dispatch_process_extern_images(post_elements_images_urls_lst, //p_input_images_urls_lst
		post_elements_images_origin_pages_urls_str,       //p_input_images_origin_pages_urls_str
		&image_job_client_type_str,
		p_gf_images_service_host_port_str,
		p_log_fun)

	if err != nil {
		return nil,err
	}
	//-------------------
	image_ids_lst := []string{}
	for _,output := range outputs_lst {

		if _,ok := post_elements_map[output.Image_source_url_str]; !ok {
			return nil,errors.New("gf_images_lib client returned results for unknown image url - "+output.Image_source_url_str)
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

	p_log_fun("INFO","post_thumbnail_str - "+post_thumbnail_str)
	//----------------
	//persists the newly updated post (some of its post_elements have been updated
	//in the initiation of image post_elements)
	err = DB__update_post(p_post, p_mongodb_coll, p_log_fun)
	if err != nil {
		return nil,err
	}
	//----------------

	return running_job_id_str,nil
}