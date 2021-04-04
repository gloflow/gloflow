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

package gf_images_lib

import (
	"fmt"
	"time"
	"math/rand"
	"context"
	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)

//---------------------------------------------------
func DB__put_image_upload_info(p_image_upload_info *Gf_image_upload_info,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {

	ctx           := context.Background()
	coll_name_str := p_runtime_sys.Mongo_coll.Name()
	gf_err        := gf_core.Mongo__insert(p_image_upload_info,
		coll_name_str,
		map[string]interface{}{
			"upload_gf_image_id_str": p_image_upload_info.Upload_gf_image_id_str,
			"caller_err_msg_str":     "failed to update/upsert gf_image in a mongodb",
		},
		ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	/*err := p_runtime_sys.Mongo_coll.Insert(p_image_upload_info)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update/upsert gf_image in a mongodb",
			"mongodb_insert_error",
			map[string]interface{}{"upload_gf_image_id_str": p_image_upload_info.Upload_gf_image_id_str,},
			err, "gf_images_lib", p_runtime_sys)
		return gf_err
	}*/

	return nil
}

//---------------------------------------------------
func DB__get_random_imgs_range(p_imgs_num_to_get_int int, // 5
	p_max_random_cursor_position_int int, // 2000
	p_flow_name_str                  string,
	p_runtime_sys                    *gf_core.Runtime_sys) ([]*gf_images_utils.Gf_image, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_image_db.DB__get_random_imgs_range()")

	rand.Seed(time.Now().Unix())
	random_cursor_position_int := rand.Intn(p_max_random_cursor_position_int) // new Random().nextInt(p_max_random_cursor_position_int)
	p_runtime_sys.Log_fun("INFO", "random_cursor_position_int - "+fmt.Sprint(random_cursor_position_int))



	ctx := context.Background()

	find_opts := options.Find()
	find_opts.SetSkip(int64(random_cursor_position_int))
    find_opts.SetLimit(int64(p_imgs_num_to_get_int))

	cursor, gf_err := gf_core.Mongo__find(bson.M{
			"t":                    "img",
			"creation_unix_time_f": bson.M{"$exists": true,},
			"flows_names_lst":      bson.M{"$in": []string{p_flow_name_str},},
			//---------------------
			// IMPORTANT!! - this is the new member that indicates which page url (if not directly uploaded) the
			//               image came from. only use these images, since only they can be properly credited
			//               to the source site
			"origin_page_url_str": bson.M{"$exists": true,},
			
			//---------------------
		},
		find_opts,
		map[string]interface{}{
			"imgs_num_to_get_int":            p_imgs_num_to_get_int,
			"max_random_cursor_position_int": p_max_random_cursor_position_int,
			"flow_name_str":                  p_flow_name_str,
			"caller_err_msg_str":             "failed to get random img range from the DB",
		},
		p_runtime_sys.Mongo_coll,
		ctx,
		p_runtime_sys)

	if gf_err != nil {
		return nil, gf_err
	}
	
	var imgs_lst []*gf_images_utils.Gf_image
	err := cursor.All(ctx, &imgs_lst)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get mongodb results of query to get Images",
			"mongodb_cursor_all",
			map[string]interface{}{
				"imgs_num_to_get_int":            p_imgs_num_to_get_int,
				"max_random_cursor_position_int": p_max_random_cursor_position_int,
				"flow_name_str":                  p_flow_name_str,
			},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}

	/*var imgs_lst []*gf_images_utils.Gf_image
	err := p_runtime_sys.Mongo_coll.Find(bson.M{
			"t":                    "img",
			"creation_unix_time_f": bson.M{"$exists": true,},
			"flows_names_lst":      bson.M{"$in": []string{p_flow_name_str},},
			//---------------------
			// IMPORTANT!! - this is the new member that indicates which page url (if not directly uploaded) the
			//               image came from. only use these images, since only they can be properly credited
			//               to the source site
			"origin_page_url_str": bson.M{"$exists": true,},
			
			//---------------------
		}).
		Skip(random_cursor_position_int).
		Limit(p_imgs_num_to_get_int).
		All(&imgs_lst)
		
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get random img range from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"imgs_num_to_get_int":            p_imgs_num_to_get_int,
				"max_random_cursor_position_int": p_max_random_cursor_position_int,
				"flow_name_str":                  p_flow_name_str,
			},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}*/

	return imgs_lst, nil
}

//---------------------------------------------------
func DB__image_exists(p_image_id_str string, p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_image_db.DB__image_exists()")


	ctx := context.Background()
	c, err := p_runtime_sys.Mongo_coll.CountDocuments(ctx, bson.M{"t": "img", "id_str": p_image_id_str})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to check if image exists in the DB",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_lib", p_runtime_sys)
		return false, gf_err
	}

	/*c,err := p_runtime_sys.Mongo_coll.Find(bson.M{"t": "img","id_str": p_image_id_str}).Count()
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to check if image exists in the DB",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_lib", p_runtime_sys)
		return false, gf_err
	}*/

	if c > 0 {
		return true, nil
	} else {
		return false, nil
	}
}