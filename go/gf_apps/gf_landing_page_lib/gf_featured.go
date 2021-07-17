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

package gf_landing_page_lib

import (
	"fmt"
	"strconv"
	"net/url"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
)

//------------------------------------------------
type Gf_featured_post struct {
	Title_str            string
	Image_url_str        string
	Url_str              string
	Images_number_int    int	
}

type Gf_featured_img struct {
	Title_str                      string
	Image_url_str                  string
	Image_thumbnail_medium_url_str string
	Image_origin_page_url_str      string // for each featured image this is the URL used in links
	Image_origin_page_url_host_str string // this is displayed in the user UI for each featured image
	Creation_unix_time_str         string
}

//------------------------------------------
// IMAGES
//------------------------------------------
func get_featured_imgs(p_max_random_cursor_position_int int, // 500
	p_elements_num_to_get_int int, // 5
	p_runtime_sys             *gf_core.Runtime_sys) ([]*Gf_featured_img, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_featured.get_featured_imgs()")

	imgs_lst, err := gf_images_utils.DB__get_random_imgs_range(p_elements_num_to_get_int,
		p_max_random_cursor_position_int,
		"general", // p_flow_name_str
		p_runtime_sys)

	if err != nil {
		return nil, err
	}

	featured_imgs_lst := []*Gf_featured_img{}
	for _, img := range imgs_lst {

		// FIX!! - create a proper gf_er
		origin_page_url, err := url.Parse(img.Origin_page_url_str)
		if err != nil {
			continue
		}

		featured := &Gf_featured_img{
			Title_str:                      img.Title_str,
			Image_url_str:                  img.Thumbnail_small_url_str,
			Image_thumbnail_medium_url_str: img.Thumbnail_medium_url_str,
			Image_origin_page_url_str:      img.Origin_page_url_str,
			Image_origin_page_url_host_str: origin_page_url.Host,
			Creation_unix_time_str:         strconv.FormatFloat(img.Creation_unix_time_f, 'f', 6, 64),
		}
		featured_imgs_lst = append(featured_imgs_lst, featured)
	}
	return featured_imgs_lst, nil
}

//------------------------------------------
// POSTS
//------------------------------------------
func get_featured_posts(p_max_random_cursor_position_int int, // 500
	p_elements_num_to_get_int int, // 5
	p_runtime_sys             *gf_core.Runtime_sys) ([]*Gf_featured_post, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_featured.get_featured_posts()")

	//gets posts starting in some random position (time wise), 
	//and as many as specified after that random point
	posts_lst, gf_err := gf_publisher_core.DB__get_random_posts_range(p_elements_num_to_get_int,
		p_max_random_cursor_position_int,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	featured_posts_lst := posts_to_featured(posts_lst, p_runtime_sys)
	return featured_posts_lst, nil
}

//------------------------------------------
func posts_to_featured(p_posts_lst []*gf_publisher_core.Gf_post, p_runtime_sys *gf_core.Runtime_sys) []*Gf_featured_post {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_featured.posts_to_featured()")

	featured_posts_lst := []*Gf_featured_post{}
	for _, post := range p_posts_lst {
		featured          := post_to_featured(post, p_runtime_sys)
		featured_posts_lst = append(featured_posts_lst, featured)
	}

	// CAUTION!! - in some cases image_src is null or "error", in which case it should not 
	//             be included in the final output. This is due to past issues/bugs in the gf_image and 
	//             gf_publisher.
	featured_elements_with_no_errors_lst := []*Gf_featured_post{}
	for _, featured := range featured_posts_lst {
		p_runtime_sys.Log_fun("INFO", "featured.Image_url_str - "+featured.Image_url_str)

		//
		if featured.Image_url_str == "" || featured.Image_url_str == "error" {
			err_msg_str := fmt.Sprintf("post with title [%s] has a image_src that is [%s]", featured.Title_str, featured.Image_url_str)
			p_runtime_sys.Log_fun("ERROR", err_msg_str)
		} else {
			featured_elements_with_no_errors_lst = append(featured_elements_with_no_errors_lst, featured)
		}
	}

	return featured_elements_with_no_errors_lst
}

//------------------------------------------
func post_to_featured(p_post *gf_publisher_core.Gf_post, p_runtime_sys *gf_core.Runtime_sys) *Gf_featured_post {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_featured.post_to_featured()")

	post_url_str := fmt.Sprintf("/posts/%s", p_post.Title_str)
	p_runtime_sys.Log_fun("INFO", "p_post.Thumbnail_url_str - "+p_post.Thumbnail_url_str)

	featured := &Gf_featured_post{
		Title_str:         p_post.Title_str,
		Image_url_str:     p_post.Thumbnail_url_str,
		Url_str:           post_url_str,
		Images_number_int: len(p_post.Images_ids_lst),
	}
	return featured
}