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

package main

import (
	"fmt"
	"strconv"
	"gopkg.in/mgo.v2"
	"apps/gf_publisher_lib"
	"apps/gf_images_lib"
)
//------------------------------------------------
type Featured_post struct {
	Title_str            string
	Image_url_str        string
	Url_str              string
	Images_number_int    int	
}

type Featured_img struct {
	Title_str                      string
	Image_url_str                  string
	Image_thumbnail_medium_url_str string
	Image_origin_page_url_str      string
	Creation_unix_time_str         string
}
//------------------------------------------
//IMAGES
//------------------------------------------
func get_featured_imgs(p_max_random_cursor_position_int int, //500
	p_elements_num_to_get_int int, //5
	p_mongodb_coll            *mgo.Collection,
	p_log_fun                 func(string,string)) ([]*Featured_img,error) {
	p_log_fun("FUN_ENTER","gf_featured.get_featured_imgs()")

	imgs_lst,err := gf_images_lib.DB__get_random_imgs_range(p_elements_num_to_get_int,
		p_max_random_cursor_position_int,
		"general", //p_flow_name_str
		p_mongodb_coll,
		p_log_fun)

	if err != nil {
		return nil,err
	}

	featured_imgs_lst := []*Featured_img{}
	for _,img := range imgs_lst {

		featured := &Featured_img{
			Title_str:                     img.Title_str,
			Image_url_str:                 img.Thumbnail_small_url_str,
			Image_thumbnail_medium_url_str:img.Thumbnail_medium_url_str,
			Image_origin_page_url_str:     img.Origin_page_url_str,
			Creation_unix_time_str:        strconv.FormatFloat(img.Creation_unix_time_f,'f',6,64),
		}
		featured_imgs_lst = append(featured_imgs_lst,featured)
	}

	return featured_imgs_lst,nil
}
//------------------------------------------
//POSTS
//------------------------------------------
func get_featured_posts(p_max_random_cursor_position_int int, //500
	p_elements_num_to_get_int int, //5
	p_mongodb_coll            *mgo.Collection,
	p_log_fun                 func(string,string)) ([]*Featured_post,error) {
	p_log_fun("FUN_ENTER","gf_featured.get_featured_posts()")

	//gets posts starting in some random position (time wise), 
	//and as many as specified after that random point
	posts_lst,err := gf_publisher_lib.DB__get_random_posts_range(p_elements_num_to_get_int,
		p_max_random_cursor_position_int,
		p_mongodb_coll,
		p_log_fun)
	if err != nil {
		return nil,err
	}

	featured_posts_lst := posts_to_featured(posts_lst, p_log_fun)
	return featured_posts_lst,nil
}
//------------------------------------------
func posts_to_featured(p_posts_lst []*gf_publisher_lib.Post, p_log_fun func(string,string)) []*Featured_post {
	p_log_fun("FUN_ENTER","gf_featured.posts_to_featured()")

	featured_posts_lst := []*Featured_post{}
	for _,post := range p_posts_lst {
		featured          := post_to_featured(post, p_log_fun)
		featured_posts_lst = append(featured_posts_lst,featured)
	}

	//CAUTION!! - in some cases image_src is null or "error", in which case it should not 
	//            be included in the final output. This is due to issues in the gf_image and 
	//			  gf_publisher internal apps
	featured_elements_with_no_errors_lst := []*Featured_post{}
	for _,featured := range featured_posts_lst {
		p_log_fun("INFO","featured.Image_url_str - "+featured.Image_url_str)
		if featured.Image_url_str == "" || featured.Image_url_str == "error" {
			err_msg_str := fmt.Sprintf("post with title [%s] has a image_src that is [%s]", featured.Title_str, featured.Image_url_str)
			p_log_fun("ERROR",err_msg_str)
		} else {
			featured_elements_with_no_errors_lst = append(featured_elements_with_no_errors_lst,featured)
		}
	}

	return featured_elements_with_no_errors_lst
}
//------------------------------------------
func post_to_featured(p_post *gf_publisher_lib.Post, p_log_fun func(string,string)) *Featured_post {
	p_log_fun("FUN_ENTER","gf_featured.post_to_featured()")

	post_url_str := fmt.Sprintf("/posts/%s",p_post.Title_str)
	p_log_fun("INFO","p_post.Thumbnail_url_str - "+p_post.Thumbnail_url_str)

	featured := &Featured_post{
		Title_str:        p_post.Title_str,
		Image_url_str:    p_post.Thumbnail_url_str,
		Url_str:          post_url_str,
		Images_number_int:len(p_post.Images_ids_lst),
	}
	return featured
}