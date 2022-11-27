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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
// VAR
//---------------------------------------------------

func db__get_objects_with_tag_count(pTagStr string,
	pObjectTypeStr string,
	pRuntimeSys    *gf_core.RuntimeSys) (int64, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_tagger_db.db__get_objects_with_tag_count()")

	switch pObjectTypeStr {
		case "post":

			ctx := context.Background()
			countInt, err := pRuntimeSys.Mongo_coll.CountDocuments(ctx, bson.M{
				"t":        "post",
				"tags_lst": bson.M{"$in": []string{pTagStr,}},
			})

			if err != nil {
				gfErr := gf_core.MongoHandleError(fmt.Sprintf("failed to count of posts with tag - %s in DB", pTagStr),
					"mongodb_find_error",
					map[string]interface{}{
						"tag_str":         pTagStr,
						"object_type_str": pObjectTypeStr,
					},
					err, "gf_tagger_lib", pRuntimeSys)
				return 0, gfErr
			}
			return countInt, nil
	}
	return 0, nil
}

//---------------------------------------------------
// POSTS
//---------------------------------------------------

func db__get_post_notes(pPostTitleStr string,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFnote, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_tagger_db.db__get_post_notes()")

	post, gfErr := gf_publisher_core.DBgetPost(pPostTitleStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	postNotesLst := post.NotesLst
	notesLst     := []*GFnote{}
	for _, s := range postNotesLst {

		note := &GFnote{
			UserIDstr:           s.UserIDstr,
			BodyStr:             s.BodyStr,
			TargetObjIDstr:      post.TitleStr,
			TargetObjTypeStr:    "post",
			CreationDatetimeStr: s.CreationDatetimeStr,
		}
		notesLst = append(notesLst, note)
	}
	pRuntimeSys.LogFun("INFO", "got # notes - "+fmt.Sprint(len(notesLst)))
	return notesLst, nil
}

//---------------------------------------------------

func db__add_post_note(pNote *GFnote,
	pPostTitleStr string,
	pRuntimeSys   *gf_core.RuntimeSys) *gf_core.GFerror {

	//--------------------
	postNote := &gf_publisher_core.GFpostNote{
		UserIDstr:           pNote.UserIDstr,
		BodyStr:             pNote.BodyStr,
		CreationDatetimeStr: pNote.CreationDatetimeStr,
	}

	//--------------------
	
	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":         "post",
			"title_str": pPostTitleStr,
		}, 
		bson.M{"$push": bson.M{"notes_lst": postNote},
	})
	
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to update a gf_post in a mongodb with a new note in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"post_title_str": pPostTitleStr,
				"note":           pNote,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}
	return nil
}

//---------------------------------------------------

func db__get_posts_with_tag(pTagStr string,
	p_page_index_int int,
	p_page_size_int  int,
	pRuntimeSys      *gf_core.RuntimeSys) ([]*gf_publisher_core.GFpost, *gf_core.GFerror) {

	// FIX!! - potentially DOESNT SCALE. if there is a huge number of posts
	//         with a tag, toList() will accumulate a large collection in memory. 
	//         instead use a Stream-oriented way where results are streamed lazily
	//         in some fashion

	ctx := context.Background()

	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"creation_datetime": -1})
	find_opts.SetSkip(int64(p_page_index_int))
    find_opts.SetLimit(int64(p_page_size_int))

	cursor, gfErr := gf_core.MongoFind(bson.M{
			"t":        "post",
			"tags_lst": bson.M{"$in": []string{pTagStr,}},
		},
		find_opts,
		map[string]interface{}{
			"tag_str":            pTagStr,
			"page_index_int":     p_page_index_int,
			"page_size_int":      p_page_size_int,
			"caller_err_msg_str": fmt.Sprintf("failed to get posts with specified tag in DB"),
		},
		pRuntimeSys.Mongo_coll,
		ctx,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}

	var posts_lst []*gf_publisher_core.GFpost
	err := cursor.All(ctx, &posts_lst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get posts with specified tag in DB",
			"mongodb_cursor_decode",
			map[string]interface{}{
				"tag_str":        pTagStr,
				"page_index_int": p_page_index_int,
				"page_size_int":  p_page_size_int,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}

	/*err := pRuntimeSys.Mongodb_coll.Find(bson.M{
			"t":        "post",
			"tags_lst": bson.M{"$in": []string{pTagStr,}},
		}).
		Sort("-creation_datetime"). // descending:true
		Skip(p_page_index_int).
		Limit(p_page_size_int).
		All(&posts_lst)

	if gfErr != nil {
		gfErr := gf_core.MongoHandleError("failed to get posts with specified tag",
			"mongodb_find_error",
			map[string]interface{}{
				"tag_str":        pTagStr,
				"page_index_int": p_page_index_int,
				"page_size_int":  p_page_size_int,
			},
			err, "gf_tagger", pRuntimeSys)
		return nil, gfErr
	}*/

	return posts_lst, nil
}

//---------------------------------------------------

func db__add_tags_to_post(pPostTitleStr string,
	pTagsLst    []string,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_tagger_db.db__add_tags_to_post()")

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":         "post",
			"title_str": pPostTitleStr,
		},
		bson.M{"$push": bson.M{"tags_lst": pTagsLst},
	})
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to update a gf_post with new tags in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"post_title_str": pPostTitleStr,
				"tags_lst":       pTagsLst,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}
	return nil
}

//---------------------------------------------------
// IMAGES
//---------------------------------------------------

func db__add_tags_to_image(pImageIDstr string,
	pTagsLst    []string,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":      "img",
			"id_str": pImageIDstr,
		},
		bson.M{"$push": bson.M{"tags_lst": pTagsLst},
	})
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to update a gf_image with new tags in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"image_id_str": pImageIDstr,
				"tags_lst":     pTagsLst,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}
	return nil
}