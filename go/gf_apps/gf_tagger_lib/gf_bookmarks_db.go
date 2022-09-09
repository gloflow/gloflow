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
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------
func db__bookmark__create(p_bookmark *GF_bookmark,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.GF_error {

	coll_name_str := "gf_bookmarks"

	gf_err := gf_core.Mongo__insert(p_bookmark,
		coll_name_str,
		map[string]interface{}{
			"url_str":            p_bookmark.Url_str,
			"caller_err_msg_str": "failed to insert GF_bookmark into the DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	return nil
}

//---------------------------------------------------
func db__bookmark__update_screenshot(p_bookmark_id_str gf_core.GF_ID,
	p_screenshot_image_id_str            gf_images_core.GF_image_id,
	p_screenshot_image_thumbnail_url_str string,
	p_ctx                                context.Context,
	p_runtime_sys                        *gf_core.RuntimeSys) *gf_core.GF_error {

	

	coll := p_runtime_sys.Mongo_db.Collection("gf_bookmarks")
	_, err := coll.UpdateMany(p_ctx, bson.M{
		"id_str": p_bookmark_id_str,
	},
	bson.M{
		"$set": bson.M{
			"screenshot_image_id_str":            p_screenshot_image_id_str,
			"Screenshot_image_thumbnail_url_str": p_screenshot_image_thumbnail_url_str,
		},
	},)
	
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update DB ",
			"mongodb_update_error",
			map[string]interface{}{
				"bookmark_id_str":                    p_bookmark_id_str,
				"screenshot_image_id_str":            p_screenshot_image_id_str,
				"screenshot_image_thumbnail_url_str": p_screenshot_image_thumbnail_url_str,
			},
			err, "gf_tagger_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}

//---------------------------------------------------
func db__bookmark__get_all(p_user_id_str gf_core.GF_ID,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.RuntimeSys) ([]*GF_bookmark, *gf_core.GF_error) {



	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"creation_unix_time_f": -1}) // descending - true - sort the latest items first
	
	db_cursor, gf_err := gf_core.Mongo__find(bson.M{
			"user_id_str":  p_user_id_str,
			"deleted_bool": false,
		},
		find_opts,
		map[string]interface{}{
			"user_id_str":        p_user_id_str,
			"caller_err_msg_str": "failed to get bookmarks for a user from DB",
		},
		p_runtime_sys.Mongo_db.Collection("gf_bookmarks"),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}



	var bookmarks_lst []*GF_bookmark
	err := db_cursor.All(p_ctx, &bookmarks_lst)
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to get DB results of query to get all Bookmarks",
			"mongodb_cursor_all",
			map[string]interface{}{
				"user_id_str": p_user_id_str,
			},
			err, "gf_tagger_lib", p_runtime_sys)
		return nil, gf_err
	}


	return bookmarks_lst, nil
}