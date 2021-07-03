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

package gf_images_lib

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"context"
	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
// IMPORTANT!! - image_flow's are ordered sequences of images, that the user creates and then
//               over time adds images to it... 

type Images_flow struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               string             `bson:"id_str"`
	T_str                string             `bson:"t"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`
	Name_str             string             `bson:"name_str"`
}

type Image_exists__check struct {
	Id                         primitive.ObjectID `bson:"_id,omitempty"`
	Id_str                     string             `bson:"id_str"`
	T_str                      string             `bson:"t"`
	Creation_unix_time_f       float64            `bson:"creation_unix_time_f"`
	Images_extern_urls_lst     []string           `bson:"images_extern_urls_lst"`
}

// //-------------------------------------------------
// // GET_MAPPING_TO_S3_BUCKETS
// func flows__get_mapping_to_s3_buckets() map[string]string {
// 	// FLOW_TO_S3_BUCKET_MAPPING - maps which image flows are going to use which S3 buckets
// 	//                             to store their images.
// 	flow_to_s3_bucket_map := map[string]string{
// 		"general":    "gf--img",
// 		"discovered": "gf--discovered--img", // "gf--img--discover"
// 	}
//
// 	return flow_to_s3_bucket_map
// }

//-------------------------------------------------
// GET_PAGE__PIPELINE
func flows__get_page__pipeline(p_req *http.Request,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.Runtime_sys) ([]*gf_images_utils.Gf_image, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_flows.flows__get_page__pipeline()")

	//--------------------
	// INPUT

	qs_map := p_req.URL.Query()

	flow_name_str := "general" // default
	if a_lst,ok := qs_map["fname"]; ok {
		flow_name_str = a_lst[0]
	}

	var err error
	page_index_int := 0 // default
	if a_lst, ok := qs_map["pg_index"]; ok {
		pg_index           := a_lst[0]
		page_index_int, err = strconv.Atoi(pg_index) // user supplied value
		
		if err != nil {
			gf_err := gf_core.Error__create("failed to parse integer pg_index query string arg",
				"int_parse_error",
				map[string]interface{}{"pg_index": pg_index,},
				err, "gf_images_lib", p_runtime_sys)
			return nil, gf_err
		}
	}

	page_size_int := 10 // default
	if a_lst, ok := qs_map["pg_size"]; ok {
		pg_size          := a_lst[0]
		page_size_int,err = strconv.Atoi(pg_size) // user supplied value
		if err != nil {
			gf_err := gf_core.Error__create("failed to parse integer pg_size query string arg",
				"int_parse_error",
				map[string]interface{}{"pg_size": pg_size,},
				err, "gf_images_lib", p_runtime_sys)
			return nil, gf_err
		}
	}

	p_runtime_sys.Log_fun("INFO",fmt.Sprintf("flow_name_str  - %s", flow_name_str))
	p_runtime_sys.Log_fun("INFO",fmt.Sprintf("page_index_int - %d", page_index_int))
	p_runtime_sys.Log_fun("INFO",fmt.Sprintf("page_size_int  - %d", page_size_int))

	//--------------------

	//--------------------
	// GET_PAGES
	cursor_start_position_int := page_index_int*page_size_int
	pages_lst, gf_err := flows_db__get_page(flow_name_str,  // "general", //p_flow_name_str
		cursor_start_position_int, // p_cursor_start_position_int
		page_size_int,             // p_elements_num_int
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//------------------
	return pages_lst, nil
}

//-------------------------------------------------
// IMAGES_EXIST_CHECK
func flows__images_exist_check(p_images_extern_urls_lst []string,
	p_flow_name_str   string,
	p_client_type_str string,
	p_runtime_sys     *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_flows.flows__images_exist_check()")

	existing_images_lst, gf_err := flows_db__images_exist(p_images_extern_urls_lst, p_flow_name_str, p_client_type_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	//-------------------------
	// PERSIST IMAGE_EXISTS_CHECK

	go func() {
		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		id_str               := fmt.Sprintf("img_exists_check:%f",creation_unix_time_f)
		
		check := Image_exists__check{
			Id_str:                 id_str,
			T_str:                  "img_exists_check",
			Creation_unix_time_f:   creation_unix_time_f,
			Images_extern_urls_lst: p_images_extern_urls_lst,
		}

		ctx           := context.Background()
		coll_name_str := p_runtime_sys.Mongo_coll.Name()
		_              = gf_core.Mongo__insert(check,
			coll_name_str,
			map[string]interface{}{
				"images_extern_urls_lst": p_images_extern_urls_lst,
				"flow_name_str":          p_flow_name_str,
				"client_type_str":        p_client_type_str,
				"caller_err_msg_str":     "failed to insert a img_exists_check into the DB",
			},
			ctx,
			p_runtime_sys)
			
		/*db_err := p_runtime_sys.Mongo_coll.Insert(check)
		if db_err != nil {
			_ = gf_core.Mongo__handle_error("failed to insert a img_exists_check in mongodb",
				"mongodb_insert_error",
				map[string]interface{}{
					"images_extern_urls_lst": p_images_extern_urls_lst,
					"flow_name_str":          p_flow_name_str,
					"client_type_str":        p_client_type_str,
				},
				db_err, "gf_images_lib", p_runtime_sys)
			return
		}*/
	}()

	//-------------------------

	return existing_images_lst, nil
}

//-------------------------------------------------
// ADD_EXTERN_IMAGE
func Flows__add_extern_image(p_image_extern_url_str string,
	p_image_origin_page_url_str string,
	p_flows_names_lst           []string,
	p_client_type_str           string,
	p_jobs_mngr_ch              chan gf_images_jobs_core.Job_msg,
	p_runtime_sys               *gf_core.Runtime_sys) (*string, *string, gf_images_utils.Gf_image_id, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_flows.Flows__add_extern_image()")
	p_runtime_sys.Log_fun("INFO",      fmt.Sprintf("p_flows_names_lst - %s",p_flows_names_lst))

	//------------------
	images_urls_to_process_lst := []gf_images_jobs_core.Gf_image_extern_to_process{
			gf_images_jobs_core.Gf_image_extern_to_process{
				Source_url_str:      p_image_extern_url_str,
				Origin_page_url_str: p_image_origin_page_url_str,
			},
		}
		
	running_job, job_expected_outputs_lst, gf_err := gf_images_jobs_core.Client__run_extern_imgs(p_client_type_str,
		images_urls_to_process_lst,
		p_flows_names_lst,
		p_jobs_mngr_ch,
		p_runtime_sys)

	if gf_err != nil {
		return nil, nil, gf_images_utils.Gf_image_id(""), gf_err
	}

	//------------------

	image_id_str                     := gf_images_utils.Gf_image_id(job_expected_outputs_lst[0].Image_id_str)
	thumbnail_small_relative_url_str := job_expected_outputs_lst[0].Thumbnail_small_relative_url_str

	return &running_job.Id_str, &thumbnail_small_relative_url_str, image_id_str, nil
}

//-------------------------------------------------
// CREATE
func flows__create(p_images_flow_name_str string,
	p_runtime_sys *gf_core.Runtime_sys) (*Images_flow, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_flows.flows__create()")

	id_str               := fmt.Sprintf("img_flow:%f", float64(time.Now().UnixNano())/1000000000.0)
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	flow := &Images_flow{
		Id_str:               id_str,
		T_str:                "img_flow",
		Name_str:             p_images_flow_name_str,
		Creation_unix_time_f: creation_unix_time_f,
	}

	ctx           := context.Background()
	coll_name_str := p_runtime_sys.Mongo_coll.Name()
	_ = gf_core.Mongo__insert(flow,
		coll_name_str,
		map[string]interface{}{
			"images_flow_name_str": p_images_flow_name_str,
			"caller_err_msg_str":   "failed to insert a image Flow into the DB",
		},
		ctx,
		p_runtime_sys)
	
	/*err := p_runtime_sys.Mongo_coll.Insert(*flow)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a image Flow in mongodb",
			"mongodb_insert_error",
			map[string]interface{}{
				"images_flow_name_str": p_images_flow_name_str,
			},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}*/

	return flow, nil
}