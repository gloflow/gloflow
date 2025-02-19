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

package gf_images_flows

import (
	// "fmt"
	"context"
	"math"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------

type GFflowMongo struct {
	Vstr              string             `bson:"v_str"` // schema_version
	Id                primitive.ObjectID `bson:"_id,omitempty"`
	IDstr             gf_core.GF_ID      `bson:"id_str"`
	CreationUNIXtimeF float64            `bson:"creation_unix_time_f"`
	NameStr           string             `bson:"name_str"`
}

//-------------------------------------------------

func DBmongoGetFlowsIDs(pFlowsNamesLst []string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]gf_core.GF_ID, *gf_core.GFerror) {

	flowsIDsLst := []gf_core.GF_ID{}
	for _, flowNameStr := range pFlowsNamesLst {
		flowIDstr, gfErr := DBmongoGetID(flowNameStr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
		flowsIDsLst = append(flowsIDsLst, flowIDstr)
	}
	return flowsIDsLst, nil
}

//---------------------------------------------------

func DBmongoGetID(pFlowNameStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GFerror) {

	collNameStr := "gf_flows"
	findOpts := options.FindOne()
	findOpts.Projection = map[string]interface{}{
		"id_str": 1,
	}
	
	flowBasicInfoMap := map[string]interface{}{}
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx,
		bson.M{
			"name_str": pFlowNameStr,	
		},
		findOpts).Decode(&flowBasicInfoMap)
		
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get user basic_info in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"flow_name_str": pFlowNameStr,
			},
			err, "gf_images_flows", pRuntimeSys)
		return "", gfErr
	}
	flowIDstr := gf_core.GF_ID(flowBasicInfoMap["id_str"].(gf_core.GF_ID))

	return flowIDstr, nil
}

//---------------------------------------------------
// GET_ALL

func DBmongoGetAll(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"t", "img"},
			}},
		},
		{
			{"$project", bson.D{
				// IMPORTANT!! - not inlucing the deprecated field "flow_name_str"
				//               since an increasingly small subset of old images is using it
				//               and new users and GF instances will never have it in their DB.
				{"flows_names_lst", true},
			}},
		},
		{
			{"$unwind", "$flows_names_lst"},
		},
		{
			{"$group", bson.D{
				{"_id", "$flows_names_lst"},
				{"count_int", bson.M{"$sum": 1}},
			}},
		},
		{
			{"$sort", bson.D{
				{"count_int", -1},
			}},
		},
	}
	cursor, err := pRuntimeSys.Mongo_coll.Aggregate(pCtx, pipeline)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to run DB aggregation to get all flows names",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}
	defer cursor.Close(pCtx)
	
	allFlowsLst := []map[string]interface{}{}
	err = cursor.All(pCtx, &allFlowsLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get mongodb results of query to get all flows names",
			"mongodb_cursor_all",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}
	
	return allFlowsLst, nil
}

//---------------------------------------------------

func DBmongoAddFlowNameToImage(p_flow_name_str string,
	p_image_gf_id_str gf_images_core.GFimageID,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {
	
	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":      "img",
			"id_str": p_image_gf_id_str,	
		},
		bson.M{

			// IMPORTANT!! - add an image to a flow
			"$addToSet": bson.M{
				"flows_names_lst": p_flow_name_str,
			},
		})
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to add a flow to an existing image DB record",
			"mongodb_update_error",
			map[string]interface{}{
				"flow_name_str":   p_flow_name_str,
				"image_gf_id_str": p_image_gf_id_str,
			},
			err, "gf_images_lib", pRuntimeSys)
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func dbMongoGetPagesTotalNum(pFlowNameStr string,
	pPageSizeInt int,
	pCtx         context.Context,
	pRuntimeSys  *gf_core.RuntimeSys) (int64, *gf_core.GFerror) {

	imagesCountInt, gfErr := gf_core.MongoCount(bson.M{
			"t":             "img",
			"flow_name_str": pFlowNameStr,
		},
		map[string]interface{}{
			"flow_name_str":      pFlowNameStr,
			"caller_err_msg_str": "failed to count the number of images in a flow",
		},
		pRuntimeSys.Mongo_coll,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return 0, gfErr
	}

	flowPagesNumInt := int64(math.Ceil(float64(imagesCountInt) / float64(pPageSizeInt)))
	return flowPagesNumInt, nil;
}

//---------------------------------------------------
// GET_PAGE

func dbMongoGetPage(pFlowNameStr string,
	pCursorStartPositionInt int, // 0
	pElementsNumInt         int, // 50
	pCtx                    context.Context,
	pRuntimeSys             *gf_core.RuntimeSys) ([]*gf_images_core.GFimage, *gf_core.GFerror) {

	findOpts := options.Find()
    findOpts.SetSort(map[string]interface{}{"creation_unix_time_f": -1}) // descending - true - sort the latest items first
	findOpts.SetSkip(int64(pCursorStartPositionInt))
    findOpts.SetLimit(int64(pElementsNumInt))
	
	cursor, gfErr := gf_core.MongoFind(bson.M{
			"t":   "img",
			"$or": []bson.M{

				// DEPRECATED!! - if a img has a flow_name_str (most due, but migrating to flows_names_lst),
				//                then match it with supplied flow_name_str.
				bson.M{"flow_name_str": pFlowNameStr,},

				// IMPORTANT!! - new approach, images can belong to multiple flows.
				//               check if the suplied flow_name_str in in the flows_names_lst list
				bson.M{"flows_names_lst": bson.M{"$in": []string{pFlowNameStr,}}},
			},
		},
		findOpts,
		map[string]interface{}{
			"flow_name_str":             pFlowNameStr,
			"cursor_start_position_int": pCursorStartPositionInt,
			"elements_num_int":          pElementsNumInt,
			"caller_err_msg_str":        "failed to get a page of images from a flow",
		},
		pRuntimeSys.Mongo_coll,
		pCtx,
		pRuntimeSys)
	
	if gfErr != nil {
		return nil, gfErr
	}
	
	imagesLst := []*gf_images_core.GFimage{}
	err       := cursor.All(pCtx, &imagesLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get a page of images from a flow",
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_images_flows", pRuntimeSys)
		return nil, gfErr
	}

	return imagesLst, nil
}