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
	"time"
	"fmt"
	"math/rand"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//---------------------------------------------------
func DB__get_post(p_post_title_str string,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_post, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_db.DB__get_post()")

	var post Gf_post
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{"t":"post", "title_str":p_post_title_str}).One(&post)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get a post from the DB",
			"mongodb_find_error",
			&map[string]interface{}{"post_title_str":p_post_title_str,},
			err, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}

	return &post, nil
}
//---------------------------------------------------
//CREATE
func DB__create_post(p_post *Gf_post,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_db.DB__create_post()")

	err := p_runtime_sys.Mongodb_coll.Insert(p_post) //writeConcern: mongo.WriteConcern.ACKNOWLEDGED);
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to create a post in the DB",
			"mongodb_insert_error",
			&map[string]interface{}{},
			err, "gf_publisher_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}
//---------------------------------------------------
//UPDATE
func DB__update_post(p_post *Gf_post, 
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_db.DB__update_post()")

	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"t":        "post",
			"title_str":p_post.Title_str,
		},
		p_post)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a gf_post in a mongodb",
			"mongodb_update_error",
			&map[string]interface{}{"post_title_str":p_post.Title_str,},
			err, "gf_publisher_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}
//---------------------------------------------------
//DELETE
func DB__mark_as_deleted_post(p_post_title_str string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_db.DB__mark_as_deleted_post()")

	err := p_runtime_sys.Mongodb_coll.Update(bson.M{
			"t":        "post",
			"title_str":p_post_title_str,
		},
		bson.M{"$set":bson.M{"deleted_bool":true}})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to mark as deleted a gf_post in a mongodb",
			"mongodb_update_error",
			&map[string]interface{}{"post_title_str":p_post_title_str,},
			err, "gf_publisher_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}
//---------------------------------------------------
//DELETE
func DB___delete_post(p_post_title_str string, p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_post_db.DB__delete_post()")

	err := p_runtime_sys.Mongodb_coll.Remove(bson.M{"title_str": p_post_title_str})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a gf_post in a mongodb",
			"mongodb_delete_error",
			&map[string]interface{}{"post_title_str": p_post_title_str,},
			err, "gf_publisher_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}
//---------------------------------------------------
//GET_POSTS_PAGE
func DB__get_posts_page(p_cursor_start_position_int int, //0
	p_elements_num_int int, //50
	p_runtime_sys      *gf_core.Runtime_sys) ([]*Gf_post,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_db.DB__get_posts_page()")

	posts_lst := []*Gf_post{}
	//descending - true - sort the latest items first
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{"t":"post"}).
		Sort("-creation_datetime_str"). //descending:true
		Skip(p_cursor_start_position_int).
		Limit(p_elements_num_int).
		All(&posts_lst)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get a posts page from the DB",
			"mongodb_find_error",
			&map[string]interface{}{
				"cursor_start_position_int":p_cursor_start_position_int,
				"elements_num_int":         p_elements_num_int,
			},
			err, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}

	return posts_lst, nil
}
//---------------------------------------------------
//REMOVE!! - is this a duplicate of DB__get_posts_page?
func DB__get_posts_from_offset(p_cursor_position_int int,
	p_posts_num_to_get_int int,
	p_runtime_sys          *gf_core.Runtime_sys) ([]*Gf_post, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_db.DB__get_posts_from_offset()")

	//----------------
	//IMPORTANT!! - because mongo"s skip() scans the collections number of docs
	//				to skip before it returns the results, p_posts_num_to_get_int should
	//				not be large, otherwise the performance will be very bad
	//assert(p_cursor_position_int < 500);
	//----------------

	var posts_lst []*Gf_post
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{"t":"post"}).
		Skip(p_cursor_position_int).
		Limit(p_posts_num_to_get_int).
		All(&posts_lst)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get a posts page from the DB",
			"mongodb_find_error",
			&map[string]interface{}{
				"cursor_start_position_int":p_cursor_position_int,
				"posts_num_to_get_int":     p_posts_num_to_get_int,
			},
			err, "gf_publisher_lib", p_runtime_sys)
		return nil, gf_err
	}

	return posts_lst, nil
}
//---------------------------------------------------
func DB__get_random_posts_range(p_posts_num_to_get_int int, //5
	p_max_random_cursor_position_int int, //500
	p_runtime_sys                    *gf_core.Runtime_sys) ([]*Gf_post, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_db.DB__get_random_posts_range()")

	rand.Seed(time.Now().Unix())
	random_cursor_position_int := rand.Intn(p_max_random_cursor_position_int) //new Random().nextInt(p_max_random_cursor_position_int)
	p_runtime_sys.Log_fun("INFO","random_cursor_position_int - "+fmt.Sprint(random_cursor_position_int))

	posts_in_random_range_lst, gf_err := DB__get_posts_from_offset(random_cursor_position_int, p_posts_num_to_get_int, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	return posts_in_random_range_lst, nil
}
//---------------------------------------------------
func DB__check_post_exists(p_post_title_str string,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_post_db.DB__check_post_exists()")
	
	count_int,err := p_runtime_sys.Mongodb_coll.Find(bson.M{
			"t":        "post",
			"title_str":p_post_title_str,
		}).Count()

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to check if the post exists in DB",
			"mongodb_find_error",
			&map[string]interface{}{"post_title_str":p_post_title_str,},
			err, "gf_publisher_lib", p_runtime_sys)
		return false, gf_err
	}
	if count_int > 0 {
		return true, nil
	} else {
		return false, nil
	}
}