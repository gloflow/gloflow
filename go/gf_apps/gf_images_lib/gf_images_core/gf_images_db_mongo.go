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

		// FIX!! - a record not being found in the DB is possible valid state. it should be considered
		//         if this should not return an error but instead just a "nil" value for the record.
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

func DBmongoImageExists(pImageIDstr GFimageID,
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

func DBmongoGetRandomImagesRange(pImgsNumToGetInt int, // 5
	pMaxRandomCursorPositionInt int, // 2000
	pFlowNameStr                string,
	pUserID                     gf_core.GF_ID,
	pRuntimeSys                 *gf_core.RuntimeSys) ([]*GFimage, *gf_core.GFerror) {

	// reseed the random number source
	rand.Seed(time.Now().UnixNano())
	
	random_cursor_position_int := rand.Intn(pMaxRandomCursorPositionInt) // new Random().nextInt(pMaxRandomCursorPositionInt)
	pRuntimeSys.LogNewFun("DEBUG", "imgs_num_to_get_int        - "+fmt.Sprint(pImgsNumToGetInt), nil)
	pRuntimeSys.LogNewFun("DEBUG", "random_cursor_position_int - "+fmt.Sprint(random_cursor_position_int), nil)



	ctx := context.Background()

	find_opts := options.Find()
	find_opts.SetSkip(int64(random_cursor_position_int))
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