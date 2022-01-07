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

package gf_images_core

import (
	"fmt"
	"context"
	"time"
	"math/rand"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func DB__put_image(p_image *GF_image,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {
	
	// UPSERT
	query  := bson.M{"t": "img", "id_str": p_image.Id_str,}
	gf_err := gf_core.Mongo__upsert(query,
		p_image,
		map[string]interface{}{"image_id_str": p_image.Id_str,},
		p_runtime_sys.Mongo_coll,
		p_ctx, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}

//---------------------------------------------------
func DB__get_image(p_image_id_str GF_image_id,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_image, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_db.DB__get_image()")
	


	ctx := context.Background()
	var image GF_image

	q             := bson.M{"t": "img", "id_str": p_image_id_str}
	coll_name_str := p_runtime_sys.Mongo_coll.Name()
	err           := p_runtime_sys.Mongo_db.Collection(coll_name_str).FindOne(ctx, q).Decode(&image)
	if err != nil {

		// FIX!! - a record not being found in the DB is possible valid state. it should be considered
		//         if this should not return an error but instead just a "nil" value for the record.
		if err == mongo.ErrNoDocuments {
			gf_err := gf_core.Mongo__handle_error("image does not exist in mongodb",
				"mongodb_not_found_error",
				map[string]interface{}{"image_id_str": p_image_id_str,},
				err, "gf_images_core", p_runtime_sys)
			return nil, gf_err
		}
		
		gf_err := gf_core.Mongo__handle_error("failed to get image from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_core", p_runtime_sys)
		return nil, gf_err
	}


	/*var image Gf_image
	err := p_runtime_sys.Mongo_coll.Find(bson.M{"t": "img", "id_str": p_image_id_str}).One(&image)

	// FIX!! - a record not being found in the DB is possible valid state. it should be considered
	//         if this should not return an error but instead just a "nil" value for the record.
	if fmt.Sprint(err) == "not found" {
		gf_err := gf_core.Mongo__handle_error("image does not exist in mongodb",
			"mongodb_not_found_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_core", p_runtime_sys)
		return nil, gf_err
	}

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get image from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_core", p_runtime_sys)
		return nil, gf_err
	}*/
	
	return &image, nil
}

//---------------------------------------------------
func DB__image_exists(p_image_id_str string, p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_db.DB__image_exists()")

	ctx := context.Background()
	c, err := p_runtime_sys.Mongo_coll.CountDocuments(ctx, bson.M{"t": "img", "id_str": p_image_id_str})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to check if image exists in the DB",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_core", p_runtime_sys)
		return false, gf_err
	}

	if c > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

//---------------------------------------------------
func DB__get_random_imgs_range(p_imgs_num_to_get_int int, // 5
	p_max_random_cursor_position_int int, // 2000
	p_flow_name_str                  string,
	p_runtime_sys                    *gf_core.Runtime_sys) ([]*GF_image, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_db.DB__get_random_imgs_range()")

	// reseed the random number source
	rand.Seed(time.Now().UnixNano())
	
	random_cursor_position_int := rand.Intn(p_max_random_cursor_position_int) // new Random().nextInt(p_max_random_cursor_position_int)
	p_runtime_sys.Log_fun("INFO", "imgs_num_to_get_int        - "+fmt.Sprint(p_imgs_num_to_get_int))
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
	
	var imgs_lst []*Gf_image
	err := cursor.All(ctx, &imgs_lst)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get mongodb results of query to get Images",
			"mongodb_cursor_all",
			map[string]interface{}{
				"imgs_num_to_get_int":            p_imgs_num_to_get_int,
				"max_random_cursor_position_int": p_max_random_cursor_position_int,
				"flow_name_str":                  p_flow_name_str,
			},
			err, "gf_images_core", p_runtime_sys)
		return nil, gf_err
	}

	/*var imgs_lst []*gf_images_core.Gf_image
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
			err, "gf_images_core", p_runtime_sys)
		return nil, gf_err
	}*/

	return imgs_lst, nil
}