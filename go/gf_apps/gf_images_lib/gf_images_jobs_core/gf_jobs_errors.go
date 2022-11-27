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

type JobError struct {
	Type_str             string `bson:"type_str"`  //"fetcher_error"|"transformer_error"
	Error_str            string `bson:"error_str"` //serialization of the golang error
	Image_source_url_str string `bson:"image_source_url_str"`
}

//-------------------------------------------------

func jobErrorSend(pJobErrorTypeStr string,
	pGFerr             *gf_core.GFerror,
	pImageSourceURLstr string,
	pImageIDstr        gf_images_core.GFimageID,
	pJobIDstr          string,
	pJobUpdatesCh      chan JobUpdateMsg,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	pRuntimeSys.LogFun("ERROR", fmt.Sprintf("fetching image failed - %s - %s", pImageSourceURLstr, pGFerr.Error))

	error_str := fmt.Sprint(pGFerr.Error)
	
	update_msg := JobUpdateMsg{
		Name_str:             pGFerr.Type_str,
		Type_str:             JOB_UPDATE_TYPE__ERROR,
		ImageIDstr:           pImageIDstr,
		Image_source_url_str: pImageSourceURLstr,
		Err_str:              error_str,
	}

	pJobUpdatesCh <- update_msg

	//-------------------------------------------------
	go func() {
		_ = jobErrorPersist(pJobIDstr,
			pJobErrorTypeStr,
			error_str,
			pImageSourceURLstr,
			pRuntimeSys)
	}()
	
	//-------------------------------------------------
	return nil
}

//-------------------------------------------------

func jobErrorPersist(pJobIDstr string,
	pErrorTypeStr      string,
	pErrorStr          string,
	pImageSourceURLstr string,
	pRuntimeSys        *gf_core.RuntimeSys) *gf_core.GFerror {

	jobError := JobError{
		Type_str:             pErrorTypeStr,
		Error_str:            pErrorStr,
		Image_source_url_str: pImageSourceURLstr,
	}

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":      "img_running_job",
			"id_str": pJobIDstr,
		},
		bson.M{
			"$push": bson.M{"errors_lst": jobError,},
		})

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to update img_running_job type document in mongodb, to add a job error",
			"mongodb_update_error",
			map[string]interface{}{
				"job_id_str":           pJobIDstr,
				"error_type_str":       pErrorTypeStr,
				"error_str":            pErrorStr,
				"image_source_url_str": pImageSourceURLstr,
			},
			err, "gf_images_jobs", pRuntimeSys)
		return gfErr
	}

	return nil
}