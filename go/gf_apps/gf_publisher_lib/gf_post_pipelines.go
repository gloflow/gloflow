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

func PipelineCreatePost(pPostInfoMap map[string]interface{},
	pImagesRuntimeInfo *GF_images_extern_runtime_info,
	pRuntimeSys        *gf_core.RuntimeSys) (*gf_publisher_core.GFpost, string, *gf_core.GFerror) {

	//----------------------
	// VERIFY INPUT
	max_title_chars_int       := 100
	max_description_chars_int := 1000
	post_element_tag_max_int  := 20

	verified_post_info_map, gfErr := gf_publisher_core.Verify_external_post_info(pPostInfoMap,
		max_title_chars_int,
		max_description_chars_int,
		post_element_tag_max_int,
		pRuntimeSys)
	if gfErr != nil {
		return nil, "", gfErr
	}

	//----------------------
	// CREATE POST
	post, gfErr := gf_publisher_core.CreateNewPost(verified_post_info_map, pRuntimeSys)
	if gfErr != nil {
		return nil, "", gfErr
	}

	pRuntimeSys.LogFun("INFO","post - "+fmt.Sprint(post))

	//----------------------
	// PERSIST POST
	gfErr = gf_publisher_core.DBcreatePost(post, pRuntimeSys)
	if gfErr != nil {
		return nil, "", gfErr
	}

	//----------------------
	// IMAGES
	// IMPORTANT - long-lasting image operation
	imagesJobIDstr, gfErr := processExternalImages(post, pImagesRuntimeInfo, pRuntimeSys)
	if gfErr != nil {
		return nil, "", gfErr
	}

	//----------------------

	return post, imagesJobIDstr, nil
}

//------------------------------------------------

func PipelineGetPost(pPostTitleStr string,
	p_response_format_str string,
	p_tmpl                *template.Template,
	pSubtemplatesNamesLst []string,
	p_resp                io.Writer,
	pRuntimeSys           *gf_core.RuntimeSys) *gf_core.GFerror {

	post, gfErr := gf_publisher_core.DBgetPost(pPostTitleStr, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------
	// HACK!!
	// some of the post_adt.tags_lst have tags that are empty strings (")
	// which showup as artifacts in the HTML since each tag gets a <div></div>
	// so here a post_adt is modified in place. 
	// this will over time correct/remove empty string tags, but the source cause of this
	// (on post tagging/creation) is still there, so find that and fix it.
	wholeTagsLst := []string{}
	for _, tagStr := range post.TagsLst {
		if tagStr != "" {
			wholeTagsLst = append(wholeTagsLst, tagStr)
		}
	}
	post.TagsLst = wholeTagsLst
	
	//------------------

	switch p_response_format_str {
		//------------------
		// HTML RENDERING
		case "html":

			// SCALABILITY!!
			// ADD!! - cache this result in redis, and server it from there
			//         only re-generate the template every so often
			//         or figure out some quick way to check if something changed
			gfErr := post__render_template(post, p_tmpl, pSubtemplatesNamesLst, p_resp, pRuntimeSys)
			if gfErr != nil {
				return gfErr
			}

		//------------------
		// JSON EXPORT
		case "json":

			postByteLst, err := json.Marshal(post)
			if err != nil {
				gfErr := gf_core.ErrorCreate("failed to serialize a Post into JSON form",
					"json_marshal_error",
					map[string]interface{}{"post_title_str": pPostTitleStr,},
					err, "gf_publisher_lib", pRuntimeSys)
				return gfErr
			}
			postStr := string(postByteLst)
			p_resp.Write([]byte(postStr))

		//------------------
	}
	return nil
}