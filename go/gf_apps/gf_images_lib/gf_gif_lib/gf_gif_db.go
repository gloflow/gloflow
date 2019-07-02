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

package gf_gif_lib

import (
	"fmt"
	"time"
	"net/url"
	"github.com/globalsign/mgo/bson"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)

//--------------------------------------------------
func gif_db__create(p_image_source_url_str string,
	p_image_origin_page_url_str string,
	p_img_width_int             int,
	p_img_height_int            int,
	p_frames_num_int            int,
	p_frames_s3_urls_lst        []string,
	p_runtime_sys               *gf_core.Runtime_sys) (*Gf_gif,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_gif_db.gif_db__create()")

	img_title_str, gf_err := gf_images_utils.Get_image_title_from_url(p_image_source_url_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("%f:gif",creation_unix_time_f)
	gf_url_str           := fmt.Sprintf("/images/d/gifs/%s.gif",img_title_str)

	//--------------
	origin_page_url, err := url.Parse(p_image_origin_page_url_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse GIF's origin_page url when creating a DB record",
			"url_parse_error",
			map[string]interface{}{
				"image_source_url_str":      p_image_source_url_str,
				"image_origin_page_url_str": p_image_origin_page_url_str,
			},
			err, "gf_gif_lib", p_runtime_sys)
		return nil, gf_err
	}
	//--------------

	gif := &Gf_gif{
		Id_str:                     id_str,
		T_str:                      "gif",
		Creation_unix_time_f:       creation_unix_time_f,
		Deleted_bool:               false,
		Valid_bool:                 true,
		Title_str:                  img_title_str,
		Gf_url_str:                 gf_url_str,
		Origin_url_str:             p_image_source_url_str,
		Origin_page_url_str:        p_image_origin_page_url_str,
		Origin_page_domain_str:     origin_page_url.Host,
		Width_int:                  p_img_width_int,
		Height_int:                 p_img_height_int,
		Preview_frames_num_int:     p_frames_num_int,
		Preview_frames_s3_urls_lst: p_frames_s3_urls_lst,
		Tags_lst:                   []string{},
	}

	err = p_runtime_sys.Mongodb_coll.Insert(gif)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a GIF in mongodb",
			"mongodb_insert_error",
			map[string]interface{}{
				"image_source_url_str":      p_image_source_url_str,
				"image_origin_page_url_str": p_image_origin_page_url_str,
			},
			err,"gf_gif_lib",p_runtime_sys)
		return nil, gf_err
	}
	return gif, nil
}

//--------------------------------------------------
func gif_db__delete(p_id_str string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_gif_db.gif_db__delete()")

	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"t":      "gif",
			"id_str": p_id_str,
		},
		bson.M{
			"$set": bson.M{"deleted_bool": true,},
		})
	
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to mark a GIF as deleted in mongodb",
			"mongodb_update_error",
			map[string]interface{}{"gif_id_str": p_id_str,},
			err, "gf_gif_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}

//--------------------------------------------------
func gif_db__get_by_img_id(p_gf_img_id_str string,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_gif,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_gif_db.gif_db__get_by_img_id()")

	var gif Gf_gif
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{
			"t":                   "gif",
			"deleted_bool":        false,
			"gf_image_id_str":     p_gf_img_id_str,
			"title_str":           bson.M{"$exists":true,},
			"origin_page_url_str": bson.M{"$exists":true,},
			"tags_lst":            bson.M{"$exists":true,},
		}).One(&gif)

	if fmt.Sprint(err) == "not found" {
		gf_err := gf_core.Mongo__handle_error("GIF with gf_img_id_str not found",
			"mongodb_not_found_error",
			map[string]interface{}{"gf_img_id_str": p_gf_img_id_str,},
			err, "gf_gif_lib", p_runtime_sys)
		return nil, gf_err
	}

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("GIF with gf_img_id_str failed the DB find operation",
			"mongodb_find_error",
			map[string]interface{}{"gf_img_id_str": p_gf_img_id_str,},
			err, "gf_gif_lib", p_runtime_sys)
		return nil, gf_err
	}

	spew.Dump(gif)

	return &gif, nil
}

//--------------------------------------------------
func gif_db__get_by_origin_url(p_origin_url_str string,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_gif, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_gif_db.gif_db__get_by_origin_url()")

	var gif Gf_gif
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{
			"t":                   "gif",
			"deleted_bool":        false,
			"origin_url_str":      p_origin_url_str,
			"title_str":           bson.M{"$exists":true,},
			"origin_page_url_str": bson.M{"$exists":true,},
			"tags_lst":            bson.M{"$exists":true,},
		}).One(&gif)

	if fmt.Sprint(err) == "not found" {
		gf_err := gf_core.Mongo__handle_error("GIF with origin_url_str not found",
			"mongodb_not_found_error",
			map[string]interface{}{"origin_url_str": p_origin_url_str,},
			err, "gf_gif_lib", p_runtime_sys)
		return nil,gf_err
	}

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("GIF with origin_url_str failed the DB find operation",
			"mongodb_find_error",
			map[string]interface{}{"origin_url_str": p_origin_url_str,},
			err, "gf_gif_lib", p_runtime_sys)
		return nil,gf_err
	}

	return &gif,nil
}

//--------------------------------------------------
func gif_db__get_page(p_cursor_start_position_int int, //0
	p_elements_num_int int,                //50
	p_runtime_sys      *gf_core.Runtime_sys) ([]*Gf_gif, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_gif_db.gif_db__get_page()")

	gifs_lst := []*Gf_gif{}

	//descending - true - sort the latest items first
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{
			"t":                      "gif",
			"valid_bool":             true,
			"preview_frames_num_int": bson.M{"$gte":0},
			"title_str":              bson.M{"$exists":true,},
			"origin_page_url_str":    bson.M{"$exists":true,},
			"tags_lst":               bson.M{"$exists":true,},
		}).
		Sort("-creation_unix_time_f"). //descending:true
		Skip(p_cursor_start_position_int).
		Limit(p_elements_num_int).
		All(&gifs_lst)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("GIFs pages failed to be retreived",
			"mongodb_find_error",
			map[string]interface{}{
				"cursor_start_position_int": p_cursor_start_position_int,
				"elements_num_int":          p_elements_num_int,
			},
			err, "gf_gif_lib", p_runtime_sys)
		return nil,gf_err
	}

	return gifs_lst,nil
}

//--------------------------------------------------
func gif_db__update_image_id(p_gif_id_str string,
	p_image_id_str gf_images_utils.Gf_image_id,
	p_runtime_sys  *gf_core.Runtime_sys) *gf_core.Gf_error {

	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"t":      "gif",
			"id_str": p_gif_id_str,
		},
		bson.M{"$set": bson.M{"gf_image_id_str": p_image_id_str,},})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to mark a GIF's gf_image_id_str in mongodb",
			"mongodb_update_error",
			map[string]interface{}{
				"gif_id_str":   p_gif_id_str,
				"image_id_str": p_image_id_str,
			},
			err, "gf_gif_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}