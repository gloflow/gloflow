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
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/globalsign/mgo/bson"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//--------------------------------------------------

func dbMongoCreate(p_image_source_url_str string,
	p_image_origin_page_url_str string,
	p_img_width_int             int,
	p_img_height_int            int,
	p_frames_num_int            int,
	p_frames_s3_urls_lst        []string,
	pRuntimeSys                 *gf_core.RuntimeSys) (*GFgif, *gf_core.GFerror) {

	imgTitleStr, gfErr := gf_images_core.GetImageTitleFromURL(p_image_source_url_str, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("%f:gif",creation_unix_time_f)
	gf_url_str           := fmt.Sprintf("/images/d/gifs/%s.gif", imgTitleStr)

	//--------------
	origin_page_url, err := url.Parse(p_image_origin_page_url_str)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to parse GIF's origin_page url when creating a DB record",
			"url_parse_error",
			map[string]interface{}{
				"image_source_url_str":      p_image_source_url_str,
				"image_origin_page_url_str": p_image_origin_page_url_str,
			},
			err, "gf_gif_lib", pRuntimeSys)
		return nil, gfErr
	}

	//--------------

	gif := &GFgif{
		Id_str:                     id_str,
		T_str:                      "gif",
		Creation_unix_time_f:       creation_unix_time_f,
		Deleted_bool:               false,
		Valid_bool:                 true,
		Title_str:                  imgTitleStr,
		Gf_url_str:                 gf_url_str,
		Origin_url_str:             p_image_source_url_str,
		Origin_page_url_str:        p_image_origin_page_url_str,
		Origin_page_domain_str:     origin_page_url.Host,
		Width_int:                  p_img_width_int,
		Height_int:                 p_img_height_int,
		Preview_frames_num_int:     p_frames_num_int,
		PreviewFramesS3urlsLst:     p_frames_s3_urls_lst,
		Tags_lst:                   []string{},
	}


	ctx           := context.Background()
	coll_name_str := pRuntimeSys.Mongo_coll.Name()
	gfErr         = gf_core.MongoInsert(gif,
		coll_name_str,
		map[string]interface{}{
			"image_source_url_str":      p_image_source_url_str,
			"image_origin_page_url_str": p_image_origin_page_url_str,
			"caller_err_msg_str":        "failed to insert a GIF into the DB",
		},
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return gif, nil
}

//--------------------------------------------------

func dbMongoDelete(p_id_str string,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":      "gif",
			"id_str": p_id_str,
		},
		bson.M{
			"$set": bson.M{"deleted_bool": true,},
		})
	
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to mark a GIF as deleted in mongodb",
			"mongodb_update_error",
			map[string]interface{}{"gif_id_str": p_id_str,},
			err, "gf_gif_lib", pRuntimeSys)
		return gfErr
	}
	return nil
}

//--------------------------------------------------

func dbMongoGetByImageID(p_gf_img_id_str string,
	pRuntimeSys *gf_core.RuntimeSys) (*GFgif, *gf_core.GFerror) {

	ctx := context.Background()

	var gif GFgif
	err := pRuntimeSys.Mongo_coll.FindOne(ctx, bson.M{
			"t":                   "gif",
			"deleted_bool":        false,
			"gf_image_id_str":     p_gf_img_id_str,
			"title_str":           bson.M{"$exists": true,},
			"origin_page_url_str": bson.M{"$exists": true,},
			"tags_lst":            bson.M{"$exists": true,},
		}).Decode(&gif)

	// FIX!! - a record not being found in the DB is possible valid state. it should be considered
	//         if this should not return an error but instead just a "nil" value for the record.
	if err == mongo.ErrNoDocuments {
		gfErr := gf_core.MongoHandleError("GIF with gf_img_id_str not found",
			"mongodb_not_found_error",
			map[string]interface{}{"gf_img_id_str": p_gf_img_id_str,},
			err, "gf_gif_lib", pRuntimeSys)
		return nil, gfErr
	}

	if err != nil {
		gfErr := gf_core.MongoHandleError("GIF with gf_img_id_str failed the DB find operation",
			"mongodb_find_error",
			map[string]interface{}{"gf_img_id_str": p_gf_img_id_str,},
			err, "gf_gif_lib", pRuntimeSys)
		return nil, gfErr
	}

	spew.Dump(gif)

	return &gif, nil
}

//--------------------------------------------------

func dbMongoGetByOriginURL(p_origin_url_str string,
	pRuntimeSys *gf_core.RuntimeSys) (*GFgif, *gf_core.GFerror) {

	ctx := context.Background()

	var gif GFgif
	err := pRuntimeSys.Mongo_coll.FindOne(ctx, bson.M{
			"t":                   "gif",
			"deleted_bool":        false,
			"origin_url_str":      p_origin_url_str,
			"title_str":           bson.M{"$exists": true,},
			"origin_page_url_str": bson.M{"$exists": true,},
			"tags_lst":            bson.M{"$exists": true,},
		}).Decode(&gif)

	// FIX!! - a record not being found in the DB is possible valid state. it should be considered
	//         if this should not return an error but instead just a "nil" value for the record.
	if err == mongo.ErrNoDocuments {
		gfErr := gf_core.MongoHandleError("GIF with origin_url_str not found",
			"mongodb_not_found_error",
			map[string]interface{}{"origin_url_str": p_origin_url_str,},
			err, "gf_gif_lib", pRuntimeSys)
		return nil,gfErr
	}

	if err != nil {
		gfErr := gf_core.MongoHandleError("GIF with origin_url_str failed the DB find operation",
			"mongodb_find_error",
			map[string]interface{}{"origin_url_str": p_origin_url_str,},
			err, "gf_gif_lib", pRuntimeSys)
		return nil, gfErr
	}

	return &gif, nil
}

//--------------------------------------------------

func dbMongoGetPage(p_cursor_start_position_int int, // p_elements_num_int
	p_elements_num_int int,
	pRuntimeSys        *gf_core.RuntimeSys) ([]GFgif, *gf_core.GFerror) {

	ctx := context.Background()

	find_opts := options.Find()
    find_opts.SetSort(map[string]interface{}{"creation_unix_time_f": -1}) // descending - true - sort the latest items first
	find_opts.SetSkip(int64(p_cursor_start_position_int))
    find_opts.SetLimit(int64(p_elements_num_int))

	cursor, gfErr := gf_core.MongoFind(bson.M{
			"t":                      "gif",
			"valid_bool":             true,
			"preview_frames_num_int": bson.M{"$gte": 0},
			"title_str":              bson.M{"$exists": true,},
			"origin_page_url_str":    bson.M{"$exists": true,},
			"tags_lst":               bson.M{"$exists": true,},
		},
		find_opts,
		map[string]interface{}{
			"cursor_start_position_int": p_cursor_start_position_int,
			"elements_num_int":          p_elements_num_int,
			"caller_err_msg_str":        "GIFs pages failed to be retreived",
		},
		pRuntimeSys.Mongo_coll,
		ctx,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}
	
	gifs_lst := []GFgif{}
	for cursor.Next(ctx) {

		var gf_gif GFgif
		err := cursor.Decode(&gf_gif)
		if err != nil {
			gfErr := gf_core.MongoHandleError("failed to decode mongodb result of query to get GIFs",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_gif_lib", pRuntimeSys)

			return nil, gfErr
		}
		gifs_lst = append(gifs_lst, gf_gif)
	}

	return gifs_lst, nil
}

//--------------------------------------------------

func dbMongoUpdateImageID(p_gif_id_str string,
	p_image_id_str gf_images_core.GFimageID,
	pRuntimeSys  *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":      "gif",
			"id_str": p_gif_id_str,
		},
		bson.M{"$set": bson.M{"gf_image_id_str": p_image_id_str,},})

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to mark a GIF's gf_image_id_str in mongodb",
			"mongodb_update_error",
			map[string]interface{}{
				"gif_id_str":   p_gif_id_str,
				"image_id_str": p_image_id_str,
			},
			err, "gf_gif_lib", pRuntimeSys)
		return gfErr
	}
	return nil
}