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

package gf_images_service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type Browser__job_run_result struct {
	Id                           primitive.ObjectID `bson:"_id,omitempty"`
	T_str                        string        `bson:"t"`
	Img__id_str                  string        `bson:"img__id_str"`
	Img__dominant_color_str      string        `bson:"img__dominant_color_str"`
	Img__color_pallete_lst       []string      `bson:"img__color_pallete_lst"`
	Browser__unix_start_time_str float64       `bson:"browser__unix_start_time_str"`
	Browser__unix_end_time_str   float64       `bson:"browser__unix_end_time_str"`
	Browser__id_f                float64       `bson:"browser__id_f"`
}

//-------------------------------------------------
type Browser__ai_classify__job_run_result struct {
	Id                           primitive.ObjectID `bson:"_id,omitempty"`
	T_str                        string        `bson:"t"`
	Img__id_str                  string        `bson:"img__id_str"`
	Browser__unix_start_time_str float64       `bson:"browser__unix_start_time_str"`
	Browser__unix_end_time_str   float64       `bson:"browser__unix_end_time_str"`
	Browser__id_f                float64       `bson:"browser__id_f"`
}

//-------------------------------------------------
func Process__browser_image_calc_result(p_browser_jobs_runs_results_lst []map[string]interface{}, p_runtime_sys *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_images_calc.Process__browser_image_calc_result()")

	for _, m := range p_browser_jobs_runs_results_lst {

		color_pallete_lst := []string{}
		for _, c := range m["p"].([]interface{}) {
			color_pallete_lst = append(color_pallete_lst, c.(string))
		}

		image_id_str := m["i"].(string)


		browser_job_result := &Browser__job_run_result{
			T_str:                        "img__browser_run_job_result",
			Img__id_str:                  image_id_str,
			Img__dominant_color_str:      m["c"].(string),
			Img__color_pallete_lst:       color_pallete_lst,
			Browser__unix_start_time_str: m["st"].(float64),
			Browser__unix_end_time_str:   m["et"].(float64),
		}

		//IMPORTANT!! - unique browser ID is sometimes not supplied if that functionality is disabled by the user
		if m["f"] != nil {
			browser_job_result.Browser__id_f = m["f"].(float64)
		}

		ctx           := context.Background()
		coll_name_str := p_runtime_sys.Mongo_coll.Name()
		gf_err        := gf_core.MongoInsert(browser_job_result,
			coll_name_str,
			map[string]interface{}{"image_id_str": image_id_str,},
			ctx,
			p_runtime_sys)
		if gf_err != nil {
			return gf_err
		}
			
		/*err := p_runtime_sys.Mongo_coll.Insert(browser_job_result)
		if err != nil {

			gf_err := gf_core.MongoHandleError("failed to insert a Browser__job_run_result in mongodb",
				"mongodb_insert_error",
				map[string]interface{}{"image_id_str": image_id_str,},
				err, "gf_images_lib", p_runtime_sys)
			return gf_err
		}*/
	}
	return nil
}