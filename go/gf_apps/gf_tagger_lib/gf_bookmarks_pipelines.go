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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/go-playground/validator"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//---------------------------------------------------
type GF_bookmark struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`
	User_id_str          gf_core.GF_ID      `bson:"user_id_str"` // creator user of the bookmark

	Url_str         string   `bson:"url_str"`
	Description_str string   `bson:"descr_str"`
	Tags_lst        []string `bson:"tags_lst"`
}

type GF_bookmark_small struct {
	Id_str               gf_core.GF_ID `json:"id_str"`
	Creation_unix_time_f float64       `json:"creation_unix_time_f"`
	Url_str              string        `json:"url_str"`
	Description_str      string        `json:"descr_str"`
	Tags_lst             []string      `json:"tags_lst"`
}

// INPUT
type GF_bookmark__input_create struct {
	User_id_str     gf_core.GF_ID
	Url_str         string   `mapstructure:"url"   validate:"required,min=5,max=300"`
	Description_str string   `mapstructure:"descr" validate:"min=1,max=600"`
	Tags_lst        []string `mapstructure:"tags"  validate:""`
}

// INPUT
type GF_bookmark__input_get_all struct {
	User_id_str gf_core.GF_ID
}

// OUTPUT
type GF_bookmark__output_get_all struct {
	Bookmarks_lst []*GF_bookmark_small `json:"bookmarks"`
}

//---------------------------------------------------
// GET_ALL
func bookmarks__pipeline__get_all(p_input *GF_bookmark__input_get_all,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_bookmark__output_get_all, *gf_core.Gf_error) {




	bookmarks_lst, gf_err := db__bookmark__get_all(p_input.User_id_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}


	bookmarks_small_lst := []*GF_bookmark_small{}
	for _, b := range bookmarks_lst {
		bookmark_small := &GF_bookmark_small{
			Id_str:               b.Id_str,
			Creation_unix_time_f: b.Creation_unix_time_f,
			Url_str:              b.Url_str,
			Description_str:      b.Description_str,
			Tags_lst:             b.Tags_lst,
		}
		bookmarks_small_lst = append(bookmarks_small_lst, bookmark_small)
	}

	output := &GF_bookmark__output_get_all{
		Bookmarks_lst: bookmarks_small_lst,
	}
	return output, nil
}

//---------------------------------------------------
// CREATE
func bookmarks__pipeline__create(p_input *GF_bookmark__input_create,
	p_images_jobs_mngr gf_images_jobs_core.Jobs_mngr,
	p_validator        *validator.Validate,
	p_ctx              context.Context,
	p_runtime_sys      *gf_core.Runtime_sys) *gf_core.Gf_error {



	//------------------------
	// VALIDATE
	gf_err := gf_core.Validate(p_input, p_validator, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//------------------------

	user_id_str          := gf_core.GF_ID(p_input.User_id_str)
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	fields_for_id_lst := []string{
		p_input.Url_str,
		string(user_id_str),
	}
	gf_id_str := gf_core.Image_ID__md5_create(fields_for_id_lst,
		creation_unix_time_f)
	
	bookmark := &GF_bookmark{
		Id_str:               gf_id_str,
		Deleted_bool:         false,
		Creation_unix_time_f: creation_unix_time_f,
		User_id_str:          user_id_str,

		Url_str:         p_input.Url_str,
		Description_str: p_input.Description_str,
		Tags_lst:        p_input.Tags_lst,
	}
	gf_err = db__bookmark__create(bookmark, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}




	go func() {

		ctx := context.Background()
		gf_err := bookmarks__pipeline__screenshot(p_input.Url_str,
			gf_id_str,
			ctx,
			p_images_jobs_mngr,
			p_runtime_sys)
		if gf_err != nil {
			return
		}


	}()

	return nil
}

//---------------------------------------------------
// SCREENSHOTS
//---------------------------------------------------
func bookmarks__pipeline__screenshot(p_url_str string,
	p_bookmark_id_str  gf_core.GF_ID,
	p_ctx              context.Context,
	p_images_jobs_mngr gf_images_jobs_core.Jobs_mngr,
	p_runtime_sys      *gf_core.Runtime_sys) *gf_core.Gf_error {


	bookmark_local_image_name_str := fmt.Sprintf("%s.png", p_bookmark_id_str)

	gf_err := bookmarks__screenshot_create(p_url_str,
		bookmark_local_image_name_str,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}





	

	return nil
}

//---------------------------------------------------
func bookmarks__screenshot_create(p_url_str string,
	p_target_local_file_path_str string,
	p_runtime_sys                *gf_core.Runtime_sys) *gf_core.Gf_error {

	//------------------------
	// SCREENSHOT
	cmd_lst := []string{
		"chromium",
		"--headless",
		"--disable-gpu",
		"--window-size=1920,1080",
		"--screenshot",
		"--hide-scrollbars", // Hide scrollbars from screenshots
		p_url_str,
	}
	_, _, gf_err := gf_core.CLI__run_standard(cmd_lst, nil, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//------------------------
	// RENAME_SCREENSHOT_FILE
	_, _, gf_err = gf_core.CLI__run_standard([]string{"mv", "screenshot.png", p_target_local_file_path_str},
		nil, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//------------------------

	return nil
}