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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_address"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
// VAR
//---------------------------------------------------
// GET_ALL_OBJECTS_TAGS

/*
very performance heavy function, fetches all tags in the system into memory.
really only meant to be used for small datasets of tags, and currently only for temporary
migrations. not practical for large numbers of objects with large numbers of tags.
*/
func dbMongoGetAllObjectsTags(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"t", bson.M{"$in": []string{"img", "post"}}},
			}},
		},

		// filter out objects that have tags_lst set to null or of length 0
		{
			{"$match", bson.D{
				{"tags_lst", bson.M{"$exists": true}},
				{"tags_lst", bson.M{"$ne": []string{}}},
			}},
		},
		{
			{"$project", bson.D{
				{"t",        true},
				{"id_str",   true},
				{"tags_lst", true},

				/*
				if the object doesnt container the user_id_str property or its set to null, set the
				user to be "anon". otherwise set it to the user_id_str value.
				*/
				{"user_id_str", bson.M{"$ifNull": []interface{}{"$user_id_str", "anon"}}},
			}},
		},
		{
			{"$unwind", "$tags_lst"},
		},
		{
			{"$project", bson.D{
				{"tag_str",  "$tags_lst"},  // rename tags_lst to tag_str
				{"t",        true},
				{"id_str",   true},
				{"user_id_str", true},
			}},
		},
	}
	cursor, err := pRuntimeSys.Mongo_coll.Aggregate(pCtx, pipeline)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to run DB aggregation to get all image tags",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_tagger_lib", pRuntimeSys)
		return nil, gfErr
	}
	defer cursor.Close(pCtx)
	
	allTagsLst := []map[string]interface{}{}
	err = cursor.All(pCtx, &allTagsLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get mongodb results of query to get all image tags",
			"mongodb_cursor_all",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}
	
	//--------------------
	// WEB3_ADDRESSES

	allAddressesTagsLst, gfErr := gf_address.DBmongoGetAllTags(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	allTagsLst = append(allTagsLst, allAddressesTagsLst...)

	//--------------------

	return allTagsLst, nil
}

//---------------------------------------------------

func dbMongoGetObjectsWithTagCount(pTagStr string,
	pObjectTypeStr string,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) (int64, *gf_core.GFerror) {

	var countInt int64

	switch pObjectTypeStr {
	
	// IMAGE
	case "image":

		imageCountInt, err := pRuntimeSys.Mongo_coll.CountDocuments(pCtx, bson.M{
			"t":        "img",
			"tags_lst": bson.M{"$in": []string{pTagStr,}},
		})

		if err != nil {
			gfErr := gf_core.MongoHandleError(fmt.Sprintf("failed to count of images with tag - %s in DB", pTagStr),
				"mongodb_find_error",
				map[string]interface{}{
					"tag_str":         pTagStr,
					"object_type_str": pObjectTypeStr,
				},
				err, "gf_tagger_lib", pRuntimeSys)
			return 0, gfErr
		}
		countInt = imageCountInt

	// POST
	case "post":
		countPostsInt, err := pRuntimeSys.Mongo_coll.CountDocuments(pCtx, bson.M{
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
		countInt = countPostsInt
	}
	return countInt, nil
}

//---------------------------------------------------

func dbMongoGetObjectsWithTag(pTagStr string,
	pTargetTypeStr string,
	pTargetType    interface{},
	pPageIndexInt  int,
	pPageSizeInt   int,
	pCtx           context.Context,
	pRuntimeSys    *gf_core.RuntimeSys) *gf_core.GFerror {

	findOpts := options.Find()
	findOpts.SetSort(map[string]interface{}{"creation_datetime": 1})
	findOpts.SetSkip(int64(pPageIndexInt))
	findOpts.SetLimit(int64(pPageSizeInt))

	cursor, gfErr := gf_core.MongoFind(bson.M{
			"t":        pTargetTypeStr,
			"tags_lst": bson.M{"$in": []string{pTagStr,}},
		},
		findOpts,
		map[string]interface{}{
			"tag_str":            pTagStr,
			"target_type_str":    pTargetTypeStr,
			"page_index_int":     pPageIndexInt,
			"page_size_int":      pPageSizeInt,
			"caller_err_msg_str": fmt.Sprintf("failed to get %s with specified tag in DB", pTargetTypeStr),
		},
		pRuntimeSys.Mongo_coll,
		pCtx,
		pRuntimeSys)

	if gfErr != nil {
		return gfErr
	}

	var err error
	switch target := pTargetType.(type) {
	
	// IMAGE
	case *[]*gf_images_core.GFimage:
		err = cursor.All(pCtx, target)

	// POST
	case *[]*gf_publisher_core.GFpost:
		err = cursor.All(pCtx, target)
	}

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get objects with specified tag in DB",
			"mongodb_cursor_decode",
			map[string]interface{}{
				"tag_str":        pTagStr,
				"page_index_int": pPageIndexInt,
				"page_size_int":  pPageSizeInt,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------
// IMAGES
//---------------------------------------------------

func dbMongoAddTagsToImage(pImageIDstr string,
	pTagsLst    []string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//--------------------
	// INITIALIZE_TAGS_ARRAY
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(
		pCtx,
		bson.M{
			"t":      "img",
			"id_str": pImageIDstr,
			"$or": []bson.M{
				{"tags_lst": nil},
				{"tags_lst": bson.M{"$not": bson.M{"$type": 4}}},
			},
		},
		bson.M{"$set": bson.M{"tags_lst": []string{}}})
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed initialize tags_lst property of image in the DB",
			"mongodb_update_error",
			map[string]interface{}{
				"image_id_str": pImageIDstr,
				"tags_lst":     pTagsLst,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gfErr
	}

	//--------------------
	// PUSH_TAGS
	updateQuery := bson.M{"$push": bson.M{
			"tags_lst": bson.M{

				// extend the tags_lst DB list with elements from pTagsLst
				"$each": pTagsLst,
			},
		},
	}

	_, err = pRuntimeSys.Mongo_coll.UpdateMany(pCtx, bson.M{
			"t":      "img",
			"id_str": pImageIDstr,
		},
		updateQuery)
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

	//--------------------
	return nil
}

//---------------------------------------------------
// POSTS
//---------------------------------------------------

func dbMongoGetPostNotes(pPostTitleStr string,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFnote, *gf_core.GFerror) {

	post, gfErr := gf_publisher_core.DBmongoGetPost(pPostTitleStr, pRuntimeSys)
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

func dbMongoAddPostNote(pNote *GFnote,
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
// ADD_TAGS_TO_POST
// FIX!! - add tag to post by its ID, not by its Title!

func dbMongoAddTagsToPost(pPostTitleStr string,
	pTagsLst    []string,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

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