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

func DBmongoPutImage(pImage *GFimage,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "data_symphony"
	coll := pRuntimeSys.Mongo_db.Collection(collNameStr)

	// UPSERT
	query := bson.M{"t": "img", "id_str": pImage.IDstr,}
	gfErr := gf_core.MongoUpsert(query,
		pImage,
		map[string]interface{}{"image_id_str": pImage.IDstr,},
		coll,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func DBmongoGetImage(pImageIDstr GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFimage, *gf_core.GFerror) {

	collNameStr := "data_symphony"
	coll        := pRuntimeSys.Mongo_db.Collection(collNameStr)
	
	var image GFimage

	q   := bson.M{"t": "img", "id_str": pImageIDstr}
	err := coll.FindOne(pCtx, q).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			gfErr := gf_core.MongoHandleError("image does not exist in mongodb",
				"mongodb_not_found_error",
				map[string]interface{}{"image_id_str": pImageIDstr,},
				err, "gf_images_core", pRuntimeSys)
			return nil, gfErr
		}
		
		gfErr := gf_core.MongoHandleError("failed to get image from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": pImageIDstr,},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}
	
	return &image, nil
}

//---------------------------------------------------

func DBmongoImageExistsByID(pImageIDstr GFimageID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {
	
	collNameStr := "data_symphony"
	coll := pRuntimeSys.Mongo_db.Collection(collNameStr)

	c, err := coll.CountDocuments(pCtx, bson.M{"t": "img", "id_str": pImageIDstr})
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to check if image exists in the DB",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": pImageIDstr,},
			err, "gf_images_core", pRuntimeSys)
		return false, gfErr
	}

	if c > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

//---------------------------------------------------

func DBmongoImagesExistByURLs(pImagesExternURLsLst []string,
	pFlowNameStr   string,
	pClientTypeStr string,
	pUserID        gf_core.GF_ID,
	pRuntimeSys    *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	//------------------------
	var queryMap bson.M
	if pFlowNameStr == "all" {

		// ALL_FLOWS
		queryMap = bson.M{
			"t": "img",

			// only check for images owned by the target user, or "anon" images not owned by anyone
			"$or": []bson.M{
				bson.M{"user_id_str": pUserID,},
				bson.M{"user_id_str": "anon",},
			},

			// IMPORTANT!! - return all images who's origin_url_str has a value
			//               thats in the list pImagesExternURLsLst
			"origin_url_str": bson.M{"$in": pImagesExternURLsLst,},
		}
	} else {

		// SPECIFIC_FLOWS
		queryMap = bson.M{
			"t": "img",
			//------------

			"$and": []bson.M{
				bson.M{
					"$or": []bson.M{

						/*
						DEPRECATED!! - if a img has a flow_name_str (most legacy image do,
							should be migrating to flows_names_lst),
							then match it with supplied flow_name_str.
						*/
						bson.M{"flow_name_str": pFlowNameStr,},

						bson.M{"flows_names_lst": bson.M{"$in": []string{pFlowNameStr,}}},
					},
				},

				// only check for images owned by the target user, or "anon" images not owned by anyone
				bson.M{
					"$or": []bson.M{
						bson.M{"user_id_str": pUserID,},
						bson.M{"user_id_str": "anon",},
					},
				},
			},

			/*
			"$or": []bson.M{

				// DEPRECATED!! - if a img has a flow_name_str (most due, but migrating to flows_names_lst),
				//                then match it with supplied flow_name_str.
				bson.M{"flow_name_str": pFlowNameStr,},

				// IMPORTANT!! - new approach, images can belong to multiple flows.
				//               check if the suplied flow_name_str in in the flows_names_lst list
				bson.M{"flows_names_lst": bson.M{"$in": []string{pFlowNameStr,}}},
			},
			*/

			//------------
			// IMPORTANT!! - return all images who's origin_url_str has a value
			//               thats in the list pImagesExternURLsLst
			"origin_url_str": bson.M{"$in": pImagesExternURLsLst,},
		}
	}

	//------------------------

	ctx := context.Background()

	projectionMap := bson.M{
		"creation_unix_time_f": 1,
		"id_str":               1,
		"origin_url_str":       1, // image url from a page
		"origin_page_url_str":  1, // page url from which the image url was extracted
		"flows_names_lst":      1, // flows in which this image is placed
		"tags_lst":             1, // tags attached to an image
	}

	findOpts := options.Find()
	findOpts.SetProjection(projectionMap)

	cursor, gfErr := gf_core.MongoFind(queryMap,
		findOpts,
		map[string]interface{}{
			"images_extern_urls_lst": pImagesExternURLsLst,
			"flow_name_str":          pFlowNameStr,
			"client_type_str":        pClientTypeStr,
			"caller_err_msg_str":     "failed to find images in flow when checking if images exist",
		},
		pRuntimeSys.Mongo_coll,
		ctx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	var existingImagesLst []map[string]interface{}
	err := cursor.All(ctx, &existingImagesLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to find images in flow when checking if images exist",
			"mongodb_cursor_all",
			map[string]interface{}{
				"images_extern_urls_lst": pImagesExternURLsLst,
				"flow_name_str":          pFlowNameStr,
				"client_type_str":        pClientTypeStr,
			},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}

	return existingImagesLst, nil
}

//---------------------------------------------------

func DBmongoGetRandomImagesRange(pImgsNumToGetInt int, // 5
	pMaxRandomCursorPositionInt int, // 2000
	pFlowNameStr                string,
	pUserID                     gf_core.GF_ID,
	pRuntimeSys                 *gf_core.RuntimeSys) ([]*GFimage, *gf_core.GFerror) {

	// reseed the random number source
	rand.Seed(time.Now().UnixNano())
	
	randomCursorPositionInt := rand.Intn(pMaxRandomCursorPositionInt) // new Random().nextInt(pMaxRandomCursorPositionInt)
	pRuntimeSys.LogNewFun("DEBUG", "imgs_num_to_get_int        - "+fmt.Sprint(pImgsNumToGetInt), nil)
	pRuntimeSys.LogNewFun("DEBUG", "random_cursor_position_int - "+fmt.Sprint(randomCursorPositionInt), nil)



	ctx := context.Background()

	find_opts := options.Find()
	find_opts.SetSkip(int64(randomCursorPositionInt))
    find_opts.SetLimit(int64(pImgsNumToGetInt))

	collNameStr := "data_symphony"
	coll := pRuntimeSys.Mongo_db.Collection(collNameStr)

	cursor, gfErr := gf_core.MongoFind(bson.M{
			"t":                    "img",
			"creation_unix_time_f": bson.M{"$exists": true,},
			"flows_names_lst":      bson.M{"$in": []string{pFlowNameStr},},
			//---------------------
			// IMPORTANT!! - this is the new member that indicates which page url (if not directly uploaded) the
			//               image came from. only use these images, since only they can be properly credited
			//               to the source site
			"origin_page_url_str": bson.M{"$exists": true,},
			
			//---------------------
		},
		find_opts,
		map[string]interface{}{
			"imgs_num_to_get_int":            pImgsNumToGetInt,
			"max_random_cursor_position_int": pMaxRandomCursorPositionInt,
			"flow_name_str":                  pFlowNameStr,
			"caller_err_msg_str":             "failed to get random img range from the DB",
		},
		coll,
		ctx,
		pRuntimeSys)

	if gfErr != nil {
		return nil, gfErr
	}
	
	var imgsLst []*GFimage
	err := cursor.All(ctx, &imgsLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get mongodb results of query to get Images",
			"mongodb_cursor_all",
			map[string]interface{}{
				"imgs_num_to_get_int":            pImgsNumToGetInt,
				"max_random_cursor_position_int": pMaxRandomCursorPositionInt,
				"flow_name_str":                  pFlowNameStr,
			},
			err, "gf_images_core", pRuntimeSys)
		return nil, gfErr
	}

	return imgsLst, nil
}

//---------------------------------------------------

func DBmongoAddTagsToImage(pImageID GFimageID,
	pTagsLst    []string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//--------------------
	// INITIALIZE_TAGS_ARRAY
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(
		pCtx,
		bson.M{
			"t":      "img",
			"id_str": pImageID,
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
				"image_id_str": pImageID,
				"tags_lst":     pTagsLst,
			},
			err, "gf_images_core", pRuntimeSys)
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
			"id_str": string(pImageID),
		},
		updateQuery)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to update a gf_image with new tags in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"image_id_str": string(pImageID),
				"tags_lst":     pTagsLst,
			},
			err, "gf_images_core", pRuntimeSys)
		return gfErr
	}

	//--------------------
	return nil
}