/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

package gf_tagger_lib

import (
	"fmt"
	"time"
	"context"
	"text/template"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
)

//---------------------------------------------------
type GF_bookmark struct {
	V_str                string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`
	User_id_str          gf_core.GF_ID      `bson:"user_id_str"` // creator user of the bookmark

	Url_str         string   `bson:"url_str"`
	Description_str string   `bson:"description_str"`
	Tags_lst        []string `bson:"tags_lst"`

	// SCREENSHOT
	Screenshot_image_id_str            gf_images_core.GF_image_id `bson:"screenshot_image_id_str"`
	Screenshot_image_thumbnail_url_str string                      `bson:"screenshot_image_thumbnail_url_str"`
}

type GF_bookmark_extern struct {
	Id_str               gf_core.GF_ID `json:"id_str"`
	Creation_unix_time_f float64       `json:"creation_unix_time_f"`
	Url_str              string        `json:"url_str"`
	Description_str      string        `json:"description_str"`
	Tags_lst             []string      `json:"tags_lst"`
}

// INPUT
type GF_bookmark__input_create struct {
	User_id_str     gf_core.GF_ID
	Url_str         string   `mapstructure:"url_str"         validate:"required,min=5,max=400"`
	Description_str string   `mapstructure:"description_str" validate:"min=1,max=600"`
	Tags_lst        []string `mapstructure:"tags_lst"        validate:""`
}

// INPUT
type GF_bookmark__input_get struct {
	Response_format_str string
	User_id_str         gf_core.GF_ID
}

// OUTPUT
type GF_bookmark__output_get struct {
	Bookmarks_lst         []*GF_bookmark_extern
	Template_rendered_str string
}

//---------------------------------------------------
// GET
func bookmarks__pipeline__get(p_input *GF_bookmark__input_get,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_ctx                    context.Context,
	pRuntimeSys              *gf_core.RuntimeSys) (*GF_bookmark__output_get, *gf_core.GFerror) {



	// DB
	bookmarks_lst, gf_err := db__bookmark__get_all(p_input.User_id_str,
		p_ctx,
		pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}


	var output *GF_bookmark__output_get

	//------------------------
	// HTML
	if p_input.Response_format_str == "html" {
		
		// RENDER_TEMPLATE
		template_rendered_str, gf_err := render_bookmarks(bookmarks_lst,
			p_tmpl,
			p_subtemplates_names_lst,
			pRuntimeSys)
		if gf_err != nil {
			return nil, gf_err
		}


		output = &GF_bookmark__output_get{
			Template_rendered_str: template_rendered_str,
		}

	//------------------------
	// JSON
	} else if p_input.Response_format_str == "json" {
		bookmarks_small_lst := []*GF_bookmark_extern{}
		for _, b := range bookmarks_lst {

			bookmark_small := &GF_bookmark_extern{
				Id_str:               b.Id_str,
				Creation_unix_time_f: b.Creation_unix_time_f,
				Url_str:              b.Url_str,
				Description_str:      b.Description_str,
				Tags_lst:             b.Tags_lst,
			}
			bookmarks_small_lst = append(bookmarks_small_lst, bookmark_small)
		}

		output = &GF_bookmark__output_get{
			Bookmarks_lst: bookmarks_small_lst,
		}
	}

	//------------------------
	
	return output, nil
}

//---------------------------------------------------
// CREATE
func bookmarks__pipeline__create(p_input *GF_bookmark__input_create,
	pImagesJobsMngr gf_images_jobs_core.JobsMngr,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// VALIDATE

	gfErr := gf_core.ValidateStruct(p_input, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	userIDstr            := gf_core.GF_ID(p_input.User_id_str)
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	unique_vals_for_id_lst := []string{
		p_input.Url_str,
		string(userIDstr),
	}
	IDstr := gf_core.IDcreate(unique_vals_for_id_lst,
		creation_unix_time_f)
	
	bookmark := &GF_bookmark{
		V_str:                "0",
		Id_str:               IDstr,
		Deleted_bool:         false,
		Creation_unix_time_f: creation_unix_time_f,
		User_id_str:          userIDstr,

		Url_str:         p_input.Url_str,
		Description_str: p_input.Description_str,
		Tags_lst:        p_input.Tags_lst,
	}

	// DB
	gfErr = db__bookmark__create(bookmark, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------
	// SCREENSHOT

	// IMPORTANT!! - only run bookmark screenshoting if a images_jobs_mngr was
	//               supplied to run image processing.
	if pImagesJobsMngr != nil {

		go func() {

			ctx := context.Background()
			gfErr := bookmarks__pipeline__screenshot(p_input.Url_str,
			 IDstr,
				ctx,
				pImagesJobsMngr,
				pRuntimeSys)
			if gfErr != nil {
				return
			}
		}()
	}

	//------------------------
	return nil
}

//---------------------------------------------------
// SCREENSHOTS
//---------------------------------------------------
func bookmarks__pipeline__screenshot(pURLstr string,
	pBookmarkIDstr  gf_core.GF_ID,
	pCtx            context.Context,
	pImagesJobsMngr gf_images_jobs_core.JobsMngr,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	//-----------------
	// SCREENSHOT_CREATE
	bookmark_local_image_name_str := fmt.Sprintf("%s.png", pBookmarkIDstr)

	gfErr := bookmarks__screenshot_create(pURLstr,
		bookmark_local_image_name_str,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	//-----------------
	// GF_IMAGES_JOBS__RUN
	images_to_process_lst := []gf_images_jobs_core.GF_image_local_to_process{
		{
			Local_file_path_str: bookmark_local_image_name_str,
		},
	}
	client_type_str := "gf_tagger_bookmarks"
	flows_names_lst := []string{"bookmarks", }
	
	_, job_expected_outputs_lst, gfErr := gf_images_jobs_client.RunLocalImgs(client_type_str,
		images_to_process_lst,
		flows_names_lst,
		pImagesJobsMngr,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}


	screenshot_image_id_str                  := job_expected_outputs_lst[0].Image_id_str
	screenshot_image_thumbnail_small_url_str := job_expected_outputs_lst[0].Thumbnail_small_relative_url_str

	//-----------------
	// DB_UPDATE - updated bookmark with screenshot image information
	gfErr = db__bookmark__update_screenshot(pBookmarkIDstr,
		screenshot_image_id_str,
		screenshot_image_thumbnail_small_url_str,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	//-----------------
	return nil
}

//---------------------------------------------------
func bookmarks__screenshot_create(pURLstr string,
	pTargetLocalFilePathStr string,
	pRuntimeSys             *gf_core.RuntimeSys) *gf_core.GFerror {

	//------------------------
	// SCREENSHOT
	cmdLst := []string{
		"google-chrome", // "chromium",
		"--headless",
		
		//-----------------
		// FIX!! - figure out a way to eliminate usage of "no-sandbox"
		//         (possibly with Docker seccomps: "docker run --security-opt seccomp=chrome.json", while still
		//         keeping it simple for endusers to run their own GF container instances).
		//         
		// needed to run headless Chrome in containers, even when container doesnt run as root user.
		// otherwise error is reported:
		// "Failed to move to new namespace: PID namespaces supported, Network namespace supported, but failed: errno = Operation not permitted"
		// "--no-sandbox",
		
		//-----------------

		"--disable-gpu",

		//-----------------
		// RESOLUTION
		// ADD!! - screenshot mobile resolutions as well, along with the desktop resolution.
		"--window-size=1920,1080",

		//-----------------
		"--screenshot",
		"--hide-scrollbars", // Hide scrollbars from screenshots
		pURLstr,
	}
	_, _, gfErr := gf_core.CLIrunStandard(cmdLst, nil, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------
	// RENAME_SCREENSHOT_FILE
	_, _, gfErr = gf_core.CLIrunStandard([]string{"mv", "screenshot.png", pTargetLocalFilePathStr},
		nil, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	return nil
}