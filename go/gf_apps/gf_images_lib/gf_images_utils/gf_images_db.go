/*
GloFlow media management/publishing system
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

package gf_images_utils

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//---------------------------------------------------
func DB__put_image(p_image *Gf_image,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_db.DB__put_image()")
	
	//spec          - a dict specifying elements which must be present for a document to be updated
	//upsert = True - insert doc if it doesnt exist, else just update
	_,err := p_runtime_sys.Mongodb_coll.Upsert(bson.M{"t":"img","id_str":p_image.Id_str,},p_image)
	if err != nil {
		gf_err := gf_core.Error__create("failed to update/upsert gf_image in a mongodb",
			"mongodb_update_error",
			&map[string]interface{}{"image_id_str":p_image.Id_str,},
			err,"gf_images_utils",p_runtime_sys)
		return gf_err
	}

	return nil
}

//---------------------------------------------------
func DB__get_image(p_image_id_str string,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_image,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_image_db.DB__get_image()")

	var image Gf_image
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{"t":"img","id_str":p_image_id_str}).One(&image)

	if fmt.Sprint(err) == "not found" {
		gf_err := gf_core.Error__create("image does not exist in mongodb",
			"mongodb_not_found_error",
			&map[string]interface{}{"image_id_str":p_image_id_str,},
			err,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}

	if err != nil {
		gf_err := gf_core.Error__create("failed to get image from mongodb",
			"mongodb_find_error",
			&map[string]interface{}{"image_id_str":p_image_id_str,},
			err,"gf_images_utils",p_runtime_sys)
		return nil,gf_err
	}
	
	return &image,nil
}