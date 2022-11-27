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

package gf_publisher_core

import (
	"time"
	"fmt"
	"math/rand"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func DBgetPost(p_post_title_str string,
	pRuntimeSys *gf_core.RuntimeSys) (*GFpost, *gf_core.GFerror) {

	ctx := context.Background()

	var post GFpost
	err := pRuntimeSys.Mongo_coll.FindOne(ctx,
		bson.M{"t": "post", "title_str": p_post_title_str}).Decode(&post)

	// err := pRuntimeSys.Mongodb_coll.Find(bson.M{"t":"post", "title_str": p_post_title_str}).One(&post)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get a post from the DB",
			"mongodb_find_error",
			map[string]interface{}{"post_title_str": p_post_title_str,},
			err, "gf_publisher_core", pRuntimeSys)
		return nil, gfErr
	}

	return &post, nil
}

//---------------------------------------------------
// CREATE

func DBcreatePost(p_post *GFpost,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx           := context.Background()
	coll_name_str := pRuntimeSys.Mongo_coll.Name()
	gfErr        := gf_core.MongoInsert(p_post,
		coll_name_str,
		map[string]interface{}{
			"caller_err_msg_str": "failed to create a post into the DB",
		},
		ctx,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}

	/*err := pRuntimeSys.Mongodb_coll.Insert(p_post) // writeConcern: mongo.WriteConcern.ACKNOWLEDGED);
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to create a post in the DB",
			"mongodb_insert_error",
			map[string]interface{}{},
			err, "gf_publisher_core", pRuntimeSys)
		return gfErr
	}*/

	return nil
}

//---------------------------------------------------
// UPDATE

func DBupdatePost(pPost *GFpost, 
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":         "post",
			"title_str": pPost.TitleStr,
		},
		pPost)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to update a gf_post in a mongodb",
			"mongodb_update_error",
			map[string]interface{}{"post_title_str":pPost.TitleStr,},
			err, "gf_publisher_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// DELETE

func DBmarkAsDeletedPost(p_post_title_str string,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":         "post",
			"title_str": p_post_title_str,
		},
		bson.M{"$set": bson.M{"deleted_bool": true}})

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to mark as deleted a gf_post in a mongodb",
			"mongodb_update_error",
			map[string]interface{}{"post_title_str": p_post_title_str,},
			err, "gf_publisher_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// DELETE

func DBdeletePost(p_post_title_str string, pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.DeleteOne(ctx, bson.M{"title_str": p_post_title_str})
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to update a gf_post in a mongodb",
			"mongodb_delete_error",
			map[string]interface{}{"post_title_str": p_post_title_str,},
			err, "gf_publisher_core", pRuntimeSys)
		return gfErr
	}
	return nil
}

//---------------------------------------------------
// GET_POSTS_PAGE

func DBgetPostsPage(p_cursor_start_position_int int, // 0
	p_elements_num_int int, // 50
	pRuntimeSys        *gf_core.RuntimeSys) ([]*GFpost, *gf_core.GFerror) {

	ctx := context.Background()
	
	// descending - true - sort the latest items first
	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"creation_datetime_str": -1})
	find_opts.SetSkip(int64(p_cursor_start_position_int))
    find_opts.SetLimit(int64(p_elements_num_int))
	
	cursor, gfErr := gf_core.MongoFind(bson.M{"t": "post"},
		find_opts,
		map[string]interface{}{
			"cursor_start_position_int": p_cursor_start_position_int,
			"elements_num_int":          p_elements_num_int,
			"caller_err_msg_str":        "failed to get a posts page from the DB",
		},
		pRuntimeSys.Mongo_coll,
		ctx,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}

	posts_lst := []*GFpost{}
	err := cursor.All(ctx, &posts_lst)

	/*err := pRuntimeSys.Mongodb_coll.Find(bson.M{"t": "post"}).
		Sort("-creation_datetime_str"). // descending:true
		Skip(p_cursor_start_position_int).
		Limit(p_elements_num_int).
		All(&posts_lst)*/

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get a posts page from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"cursor_start_position_int": p_cursor_start_position_int,
				"elements_num_int":          p_elements_num_int,
			},
			err, "gf_publisher_core", pRuntimeSys)
		return nil, gfErr
	}

	return posts_lst, nil
}

//---------------------------------------------------
// REMOVE!! - is this a duplicate of DB__get_posts_page?

func DB__get_posts_from_offset(p_cursor_position_int int,
	p_posts_num_to_get_int int,
	pRuntimeSys            *gf_core.RuntimeSys) ([]*GFpost, *gf_core.GFerror) {

	ctx := context.Background()

	//----------------
	// IMPORTANT!! - because mongo"s skip() scans the collections number of docs
	//               to skip before it returns the results, p_posts_num_to_get_int should
	//               not be large, otherwise the performance will be very bad
	// assert(p_cursor_position_int < 500);

	//----------------

	find_opts := options.Find()
	find_opts.SetSkip(int64(p_cursor_position_int))
    find_opts.SetLimit(int64(p_posts_num_to_get_int))
	
	cursor, gfErr := gf_core.MongoFind(bson.M{"t": "post"},
		find_opts,
		map[string]interface{}{
			"cursor_start_position_int": p_cursor_position_int,
			"posts_num_to_get_int":      p_posts_num_to_get_int,
			"caller_err_msg_str":        "failed to get a posts page from the DB",
		},
		pRuntimeSys.Mongo_coll,
		ctx,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}

	posts_lst := []*GFpost{}
	err := cursor.All(ctx, &posts_lst)


	/*err := pRuntimeSys.Mongodb_coll.Find(bson.M{"t": "post"}).
		Skip(p_cursor_position_int).
		Limit(p_posts_num_to_get_int).
		All(&posts_lst)*/

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get a posts page from the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"cursor_start_position_int": p_cursor_position_int,
				"posts_num_to_get_int":      p_posts_num_to_get_int,
			},
			err, "gf_publisher_core", pRuntimeSys)
		return nil, gfErr
	}

	return posts_lst, nil
}

//---------------------------------------------------

func DB__get_random_posts_range(p_posts_num_to_get_int int, // 5
	p_max_random_cursor_position_int int, // 500
	pRuntimeSys                      *gf_core.RuntimeSys) ([]*GFpost, *gf_core.GFerror) {

	rand.Seed(time.Now().Unix())
	random_cursor_position_int := rand.Intn(p_max_random_cursor_position_int) //new Random().nextInt(p_max_random_cursor_position_int)
	pRuntimeSys.LogFun("INFO","random_cursor_position_int - "+fmt.Sprint(random_cursor_position_int))

	posts_in_random_range_lst, gfErr := DB__get_posts_from_offset(random_cursor_position_int,
		p_posts_num_to_get_int,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return posts_in_random_range_lst, nil
}

//---------------------------------------------------

func DB__check_post_exists(p_post_title_str string,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {
	
	ctx := context.Background()
	count_int, err := pRuntimeSys.Mongo_coll.CountDocuments(ctx, bson.M{
		"t":         "post",
		"title_str": p_post_title_str,
	})

	/*count_int,err := pRuntimeSys.Mongodb_coll.Find(bson.M{
			"t":         "post",
			"title_str": p_post_title_str,
		}).Count()*/

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to check if the post exists in DB",
			"mongodb_find_error",
			map[string]interface{}{"post_title_str": p_post_title_str,},
			err, "gf_publisher_core", pRuntimeSys)
		return false, gfErr
	}
	if count_int > 0 {
		return true, nil
	} else {
		return false, nil
	}
}