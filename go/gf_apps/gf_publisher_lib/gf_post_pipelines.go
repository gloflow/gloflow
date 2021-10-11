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
	"io"
	"encoding/json"
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
)

//------------------------------------------------
// CREATE_POST
func Pipeline__create_post(p_post_info_map map[string]interface{},
	p_gf_images_runtime_info *GF_images_extern_runtime_info,
	p_runtime_sys            *gf_core.Runtime_sys) (*gf_publisher_core.Gf_post, string, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_post_pipelines.Pipeline__create_post()")

	//----------------------
	// VERIFY INPUT
	max_title_chars_int       := 100
	max_description_chars_int := 1000
	post_element_tag_max_int  := 20

	p_runtime_sys.Log_fun("INFO", "p_post_info_map - "+fmt.Sprint(p_post_info_map))
	verified_post_info_map, gf_err := gf_publisher_core.Verify_external_post_info(p_post_info_map,
		max_title_chars_int,
		max_description_chars_int,
		post_element_tag_max_int,
		p_runtime_sys)
	if gf_err != nil {
		return nil, "", gf_err
	}

	//----------------------
	// CREATE POST
	post, gf_err := gf_publisher_core.Create_new_post(verified_post_info_map, p_runtime_sys)
	if gf_err != nil {
		return nil, "", gf_err
	}

	p_runtime_sys.Log_fun("INFO","post - "+fmt.Sprint(post))

	//----------------------
	// PERSIST POST
	gf_err = gf_publisher_core.DB__create_post(post, p_runtime_sys)
	if gf_err != nil {
		return nil, "", gf_err
	}

	//----------------------
	//IMAGES
	//IMPORTANT - long-lasting image operation
	images_job_id_str, img_gf_err := process_external_images(post, p_gf_images_runtime_info, p_runtime_sys)
	if img_gf_err != nil {
		return nil, "", img_gf_err
	}

	//----------------------

	return post, images_job_id_str, nil
}

//------------------------------------------------
func Pipeline__get_post(p_post_title_str string,
	p_response_format_str    string,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_resp                   io.Writer,
	p_runtime_sys            *gf_core.Runtime_sys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_pipelines.Pipeline__get_post()")

	post, gf_err := gf_publisher_core.DB__get_post(p_post_title_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//------------------
	// HACK!!
	// some of the post_adt.tags_lst have tags that are empty strings (")
	// which showup as artifacts in the HTML since each tag gets a <div></div>
	// so here a post_adt is modified in place. 
	// this will over time correct/remove empty string tags, but the source cause of this
	// (on post tagging/creation) is still there, so find that and fix it.
	whole_tags_lst := []string{}
	for _, tag_str := range post.Tags_lst {
		if tag_str != "" {
			whole_tags_lst = append(whole_tags_lst, tag_str)
		}
	}
	post.Tags_lst = whole_tags_lst
	
	//------------------

	switch p_response_format_str {
		//------------------
		// HTML RENDERING
		case "html":

			// SCALABILITY!!
			// ADD!! - cache this result in redis, and server it from there
			//         only re-generate the template every so often
			//         or figure out some quick way to check if something changed
			gf_err := post__render_template(post, p_tmpl, p_subtemplates_names_lst, p_resp, p_runtime_sys)
			if gf_err != nil {
				return gf_err
			}

		//------------------
		// JSON EXPORT
		case "json":

			post_byte_lst,err := json.Marshal(post)
			if err != nil {
				gf_err := gf_core.Error__create("failed to serialize a Post into JSON form",
					"json_marshal_error",
					map[string]interface{}{"post_title_str": p_post_title_str,},
					err, "gf_publisher_lib", p_runtime_sys)
				return gf_err
			}
			post_str := string(post_byte_lst)
			p_resp.Write([]byte(post_str))

		//------------------
	}
	return nil
}