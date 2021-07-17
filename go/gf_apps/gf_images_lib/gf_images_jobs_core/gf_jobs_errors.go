// SPDX-License-Identifier: GPL-2.0
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

package gf_images_jobs_core

import (
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//-------------------------------------------------
type Job_Error struct {
	Type_str             string `bson:"type_str"`  //"fetcher_error"|"transformer_error"
	Error_str            string `bson:"error_str"` //serialization of the golang error
	Image_source_url_str string `bson:"image_source_url_str"`
}

//-------------------------------------------------
func job_error__send(p_job_error_type_str string,
	p_gf_err               *gf_core.Gf_error,
	p_image_source_url_str string,
	p_image_id_str         gf_images_core.Gf_image_id,
	p_job_id_str           string,
	p_job_updates_ch       chan Job_update_msg,
	p_runtime_sys          *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_errors.job_error__send()")

	p_runtime_sys.Log_fun("ERROR", fmt.Sprintf("fetching image failed - %s - %s", p_image_source_url_str, p_gf_err.Error))

	error_str := fmt.Sprint(p_gf_err.Error)
	
	update_msg := Job_update_msg{
		Name_str:             p_gf_err.Type_str,
		Type_str:             JOB_UPDATE_TYPE__ERROR,
		Image_id_str:         p_image_id_str,
		Image_source_url_str: p_image_source_url_str,
		Err_str:              error_str,
	}
	p_job_updates_ch <- update_msg
	//------------
	go func() {
		_ = job_error__persist(p_job_id_str,
			p_job_error_type_str,
			error_str,
			p_image_source_url_str,
			p_runtime_sys)
	}()
	
	//------------
	return nil
}

//-------------------------------------------------
func job_error__persist(p_job_id_str string,
	p_error_type_str       string,
	p_error_str            string,
	p_image_source_url_str string,
	p_runtime_sys          *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_jobs_errors.job_error__persist()")

	job_error := Job_Error{
		Type_str:             p_error_type_str,
		Error_str:            p_error_str,
		Image_source_url_str: p_image_source_url_str,
	}

	ctx := context.Background()
	_, err := p_runtime_sys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":      "img_running_job",
			"id_str": p_job_id_str,
		},
		bson.M{
			"$push": bson.M{"errors_lst": job_error,},
		})

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update img_running_job type document in mongodb, to add a job error",
			"mongodb_update_error",
			map[string]interface{}{
				"job_id_str":           p_job_id_str,
				"error_type_str":       p_error_type_str,
				"error_str":            p_error_str,
				"image_source_url_str": p_image_source_url_str,
			},
			err, "gf_images_jobs", p_runtime_sys)
		return gf_err
	}

	return nil
}