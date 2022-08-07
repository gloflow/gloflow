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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------
func DBgetFlowsIDs(pFlowsNamesLst []string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]gf_core.GF_ID, *gf_core.GF_error) {

	flowsIDsLst := []gf_core.GF_ID{}
	for _, flowNameStr := range pFlowsNamesLst {
		flowIDstr, gfErr := DBgetID(flowNameStr, pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
		flowsIDsLst = append(flowsIDsLst, flowIDstr)
	}
	return flowsIDsLst, nil
}

//---------------------------------------------------
func DBgetID(pFlowNameStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (gf_core.GF_ID, *gf_core.GF_error) {

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
		gfErr := gf_core.Mongo__handle_error("failed to get user basic_info in the DB",
			"mongodb_find_error",
			map[string]interface{}{
				"flow_name_str": pFlowNameStr,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return "", gfErr
	}
	flowIDstr := gf_core.GF_ID(flowBasicInfoMap["id_str"].(gf_core.GF_ID))

	return flowIDstr, nil
}

//---------------------------------------------------
// GET_ALL
func DBgetAll(p_ctx context.Context,
	p_runtime_sys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GF_error) {

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
	cursor, err := p_runtime_sys.Mongo_coll.Aggregate(p_ctx, pipeline)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to run DB aggregation to get all flows names",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_images_flows", p_runtime_sys)
		return nil, gf_err
	}
	defer cursor.Close(p_ctx)
	
	all_flows_lst := []map[string]interface{}{}
	err = cursor.All(p_ctx, &all_flows_lst)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get mongodb results of query to get all flows names",
			"mongodb_cursor_all",
			map[string]interface{}{},
			err, "gf_images_flows", p_runtime_sys)
		return nil, gf_err
	}
	
	return all_flows_lst, nil
}

//---------------------------------------------------
func Flows_db__add_flow_name_to_image(p_flow_name_str string,
	p_image_gf_id_str gf_images_core.Gf_image_id,
	p_runtime_sys     *gf_core.RuntimeSys) *gf_core.GF_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_flows_db.Flows_db__add_flow_name_to_image()")
	
	ctx := context.Background()
	_, err := p_runtime_sys.Mongo_coll.UpdateMany(ctx, bson.M{
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
		gf_err := gf_core.Error__create("failed to add a flow to an existing image DB record",
			"mongodb_update_error",
			map[string]interface{}{
				"flow_name_str":   p_flow_name_str,
				"image_gf_id_str": p_image_gf_id_str,
			},
			err, "gf_images_lib", p_runtime_sys)
		return gf_err
	}

	return nil
}

//---------------------------------------------------
func flows_db__get_pages_total_num(p_flow_name_str string,
	p_page_size_int int,
	p_ctx           context.Context,
	p_runtime_sys   *gf_core.RuntimeSys) (int64, *gf_core.GF_error) {



	imgs_in_flow__count_int, gf_err := gf_core.Mongo__count(bson.M{
			"t":             "img",
			"flow_name_str": p_flow_name_str,
		},
		map[string]interface{}{
			"flow_name_str":      p_flow_name_str,
			"caller_err_msg_str": "failed to count the number of images in a flow",
		},
		p_runtime_sys.Mongo_coll,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return 0, gf_err
	}


	flow_pages_num_int := int64(math.Ceil(float64(imgs_in_flow__count_int) / float64(p_page_size_int)))
	return flow_pages_num_int, nil;
}

//---------------------------------------------------
func flows_db__get_page(p_flow_name_str string,
	p_cursor_start_position_int int, // 0
	p_elements_num_int          int, // 50
	p_ctx                       context.Context,
	p_runtime_sys               *gf_core.RuntimeSys) ([]*gf_images_core.GF_image, *gf_core.GF_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_flows_db.flows_db__get_page()")

	find_opts := options.Find()
    find_opts.SetSort(map[string]interface{}{"creation_unix_time_f": -1}) // descending - true - sort the latest items first
	find_opts.SetSkip(int64(p_cursor_start_position_int))
    find_opts.SetLimit(int64(p_elements_num_int))
	
	cursor, gf_err := gf_core.Mongo__find(bson.M{
			"t":   "img",
			"$or": []bson.M{

				// DEPRECATED!! - if a img has a flow_name_str (most due, but migrating to flows_names_lst),
				//                then match it with supplied flow_name_str.
				bson.M{"flow_name_str": p_flow_name_str,},

				// IMPORTANT!! - new approach, images can belong to multiple flows.
				//               check if the suplied flow_name_str in in the flows_names_lst list
				bson.M{"flows_names_lst": bson.M{"$in": []string{p_flow_name_str,}}},
			},
		},
		find_opts,
		map[string]interface{}{
			"flow_name_str":             p_flow_name_str,
			"cursor_start_position_int": p_cursor_start_position_int,
			"elements_num_int":          p_elements_num_int,
			"caller_err_msg_str":        "failed to get a page of images from a flow",
		},
		p_runtime_sys.Mongo_coll,
		p_ctx,
		p_runtime_sys)
	
	if gf_err != nil {
		return nil, gf_err
	}
	
	images_lst := []*gf_images_core.GF_image{}
	err        := cursor.All(p_ctx, &images_lst)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get a page of images from a flow",
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}

	return images_lst, nil
}

//-------------------------------------------------
func flows_db__images_exist(p_images_extern_urls_lst []string,
	p_flow_name_str   string,
	p_client_type_str string,
	p_runtime_sys     *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GF_error) {
	
	//------------------------
	var query_map bson.M
	if p_flow_name_str == "all" {

		// ALL_FLOWS
		query_map = bson.M{
			"t": "img",

			//IMPORTANT!! - return all images who's origin_url_str has a value
			//              thats in the list p_images_extern_urls_lst
			"origin_url_str": bson.M{"$in": p_images_extern_urls_lst,},
		}
	} else {

		// SPECIFIC_FLOWS
		query_map = bson.M{
			"t": "img",
			//------------
			"$or": []bson.M{

				// DEPRECATED!! - if a img has a flow_name_str (most due, but migrating to flows_names_lst),
				//                then match it with supplied flow_name_str.
				bson.M{"flow_name_str": p_flow_name_str,},

				// IMPORTANT!! - new approach, images can belong to multiple flows.
				//               check if the suplied flow_name_str in in the flows_names_lst list
				bson.M{"flows_names_lst": bson.M{"$in": []string{p_flow_name_str,}}},
			},

			//------------
			// IMPORTANT!! - return all images who's origin_url_str has a value
			//               thats in the list p_images_extern_urls_lst
			"origin_url_str": bson.M{"$in": p_images_extern_urls_lst,},
		}
	}

	//------------------------

	ctx := context.Background()

	projectionMap := bson.M{
		"creation_unix_time_f": 1,
		"id_str":               1,
		"origin_url_str":       1, // image url from a page
		"origin_page_url_str":  1, // page url from which the image url was extracted
	}

	findOpts := options.Find()
	findOpts.SetProjection(projectionMap)

	cursor, gfErr := gf_core.Mongo__find(query_map,
		findOpts,
		map[string]interface{}{
			"images_extern_urls_lst": p_images_extern_urls_lst,
			"flow_name_str":          p_flow_name_str,
			"client_type_str":        p_client_type_str,
			"caller_err_msg_str":     "failed to find images in flow when checking if images exist",
		},
		p_runtime_sys.Mongo_coll,
		ctx,
		p_runtime_sys)
	if gfErr != nil {
		return nil, gfErr
	}

	var existingImagesLst []map[string]interface{}
	err := cursor.All(ctx, &existingImagesLst)
	if err != nil {
		gfErr := gf_core.Mongo__handle_error("failed to find images in flow when checking if images exist",
			"mongodb_cursor_all",
			map[string]interface{}{
				"images_extern_urls_lst": p_images_extern_urls_lst,
				"flow_name_str":          p_flow_name_str,
				"client_type_str":        p_client_type_str,
			},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gfErr
	}

	/*err := p_runtime_sys.Mongodb_coll.Find(query_map).

				// select which fields to include in the results
				Select(bson.M{
					"creation_unix_time_f": 1,
					"id_str":               1,
					"origin_url_str":       1, // image url from a page
					"origin_page_url_str":  1, // page url from which the image url was extracted
				}).
				All(&existingImagesLst)

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to find images in flow when checking if images exist",
			"mongodb_find_error",
			map[string]interface{}{
				"images_extern_urls_lst": p_images_extern_urls_lst,
				"flow_name_str":          p_flow_name_str,
				"client_type_str":        p_client_type_str,
			},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}*/

	return existingImagesLst, nil
}