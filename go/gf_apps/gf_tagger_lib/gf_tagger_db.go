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

package gf_tagger_lib

import (
	"fmt"
	"context"
	// "gopkg.in/mgo.v2/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
// BOOKMARKS
//---------------------------------------------------
func db__bookmark__create(p_bookmark *GF_bookmark,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

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
func db__bookmark__get_all(p_user_id_str gf_core.GF_ID,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) ([]*GF_bookmark, *gf_core.GF_error) {



	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"creation_unix_time_f": 1})
	
	db_cursor, gf_err := gf_core.Mongo__find(bson.M{
			"user_id_str":  p_user_id_str,
			"deleted_bool": false,
		},
		find_opts,
		map[string]interface{}{
			"user_id_str":        p_user_id_str,
			"caller_err_msg_str": "failed to get bookmarks for a user",
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
		gf_err := gf_core.Mongo__handle_error("failed to get mongodb results of query to get all Bookmarks",
			"mongodb_cursor_all",
			map[string]interface{}{
				"user_id_str": p_user_id_str,
			},
			err, "gf_tagger_lib", p_runtime_sys)
		return nil, gf_err
	}


	return bookmarks_lst, nil
}

//---------------------------------------------------
// VAR
//---------------------------------------------------
func db__get_objects_with_tag_count(p_tag_str string,
	p_object_type_str string,
	p_runtime_sys     *gf_core.Runtime_sys) (int64, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_db.db__get_objects_with_tag_count()")

	switch p_object_type_str {
		case "post":

			ctx := context.Background()
			count_int, err := p_runtime_sys.Mongo_coll.CountDocuments(ctx, bson.M{
				"t":        "post",
				"tags_lst": bson.M{"$in": []string{p_tag_str,}},
			})
			
			/*count_int, err := p_runtime_sys.Mongodb_coll.Find(bson.M{
					"t":        "post",
					"tags_lst": bson.M{"$in": []string{p_tag_str,}},
				}).Count()*/

			if err != nil {
				gf_err := gf_core.Mongo__handle_error(fmt.Sprintf("failed to count of posts with tag - %s", p_tag_str),
					"mongodb_find_error",
					map[string]interface{}{
						"tag_str":         p_tag_str,
						"object_type_str": p_object_type_str,
					},
					err, "gf_tagger", p_runtime_sys)
				return 0, gf_err
			}
			return count_int, nil
	}
	return 0, nil
}

//---------------------------------------------------
// POSTS
//---------------------------------------------------
func db__get_post_notes(p_post_title_str string,
	p_runtime_sys *gf_core.Runtime_sys) ([]*GF_note, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_db.db__get_post_notes()")

	post, gf_err := gf_publisher_core.DB__get_post(p_post_title_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	post_notes_lst := post.Notes_lst
	notes_lst      := []*GF_note{}
	for _,s := range post_notes_lst {

		note := &GF_note{
			User_id_str:           s.User_id_str,
			Body_str:              s.Body_str,
			Target_obj_id_str:     post.Title_str,
			Target_obj_type_str:   "post",
			Creation_datetime_str: s.Creation_datetime_str,
		}
		notes_lst = append(notes_lst,note)
	}
	p_runtime_sys.Log_fun("INFO", "got # notes - "+fmt.Sprint(len(notes_lst)))
	return notes_lst, nil
}

//---------------------------------------------------
func db__add_post_note(p_note *GF_note,
	p_post_title_str string,
	p_runtime_sys    *gf_core.Runtime_sys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_db.db__add_post_note()")

	//--------------------
	post_note := &gf_publisher_core.Gf_post_note{
		User_id_str:           p_note.User_id_str,
		Body_str:              p_note.Body_str,
		Creation_datetime_str: p_note.Creation_datetime_str,
	}

	//--------------------
	
	ctx := context.Background()
	_, err := p_runtime_sys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":         "post",
			"title_str": p_post_title_str,
		}, 
		bson.M{"$push": bson.M{"notes_lst": post_note},
	})
	
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a gf_post in a mongodb with a new note",
			"mongodb_update_error",
			map[string]interface{}{
				"post_title_str":p_post_title_str,
				"note":          p_note,
			},
			err, "gf_tagger", p_runtime_sys)
		return gf_err
	}
	return nil
}

//---------------------------------------------------
func db__get_posts_with_tag(p_tag_str string,
	p_page_index_int int,
	p_page_size_int  int,
	p_runtime_sys    *gf_core.Runtime_sys) ([]*gf_publisher_core.Gf_post, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_db.db__get_posts_with_tag()")
	p_runtime_sys.Log_fun("INFO",      fmt.Sprintf("p_tag_str - %s", p_tag_str))

	// FIX!! - potentially DOESNT SCALE. if there is a huge number of posts
	//         with a tag, toList() will accumulate a large collection in memory. 
	//         instead use a Stream-oriented way where results are streamed lazily
	//         in some fashion
		
	



	ctx := context.Background()

	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"creation_datetime": -1})
	find_opts.SetSkip(int64(p_page_index_int))
    find_opts.SetLimit(int64(p_page_size_int))

	cursor, gf_err := gf_core.Mongo__find(bson.M{
			"t":        "post",
			"tags_lst": bson.M{"$in": []string{p_tag_str,}},
		},
		find_opts,
		map[string]interface{}{
			"tag_str":            p_tag_str,
			"page_index_int":     p_page_index_int,
			"page_size_int":      p_page_size_int,
			"caller_err_msg_str": fmt.Sprintf("failed to get posts with specified tag"),
		},
		p_runtime_sys.Mongo_coll,
		ctx,
		p_runtime_sys)

	if gf_err != nil {
		return nil, gf_err
	}

	var posts_lst []*gf_publisher_core.Gf_post
	err := cursor.All(ctx, &posts_lst)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get posts with specified tag",
			"mongodb_cursor_decode",
			map[string]interface{}{
				"tag_str":        p_tag_str,
				"page_index_int": p_page_index_int,
				"page_size_int":  p_page_size_int,
			},
			err, "gf_tagger", p_runtime_sys)
		return nil, gf_err
	}

	/*err := p_runtime_sys.Mongodb_coll.Find(bson.M{
			"t":        "post",
			"tags_lst": bson.M{"$in": []string{p_tag_str,}},
		}).
		Sort("-creation_datetime"). // descending:true
		Skip(p_page_index_int).
		Limit(p_page_size_int).
		All(&posts_lst)

	if gf_err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get posts with specified tag",
			"mongodb_find_error",
			map[string]interface{}{
				"tag_str":        p_tag_str,
				"page_index_int": p_page_index_int,
				"page_size_int":  p_page_size_int,
			},
			err, "gf_tagger", p_runtime_sys)
		return nil, gf_err
	}*/

	return posts_lst, nil
}

//---------------------------------------------------
func db__add_tags_to_post(p_post_title_str string,
	p_tags_lst    []string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_db.db__add_tags_to_post()")

	ctx := context.Background()
	_, err := p_runtime_sys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":         "post",
			"title_str": p_post_title_str,
		},
		bson.M{"$push": bson.M{"tags_lst": p_tags_lst},
	})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a gf_post in mongodb with new tags",
			"mongodb_update_error",
			map[string]interface{}{
				"post_title_str":p_post_title_str,
				"tags_lst":      p_tags_lst,
			},
			err, "gf_tagger", p_runtime_sys)
		return gf_err
	}
	return nil
}

//---------------------------------------------------
// IMAGES
//---------------------------------------------------
func db__add_tags_to_image(p_image_id_str string,
	p_tags_lst    []string,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_db.db__add_tags_to_image()")

	ctx := context.Background()
	_, err := p_runtime_sys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":      "img",
			"id_str": p_image_id_str,
		},
		bson.M{"$push": bson.M{"tags_lst": p_tags_lst},
	})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update a gf_image in mongodb with new tags",
			"mongodb_update_error",
			map[string]interface{}{
				"image_id_str": p_image_id_str,
				"tags_lst":     p_tags_lst,
			},
			err, "gf_tagger", p_runtime_sys)
		return gf_err
	}
	return nil
}