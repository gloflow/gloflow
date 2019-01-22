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

package gf_images_lib

import (
	"fmt"
	"strings"
	"github.com/globalsign/mgo"
	"github.com/gloflow/gloflow/go/gf_core"
)
//--------------------------------------------------
func db_index__init(p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {

	//---------------------
	//FIELDS  - "t","flows_names_lst","origin_url_str"
	//QUERIES - flows_db__images_exist() issues these queries

	doc_type__index := mgo.Index{
		Key:        []string{"t","flows_names_lst","origin_url_str"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := p_runtime_sys.Mongodb_coll.EnsureIndex(doc_type__index)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "duplicate key error index") {
			//ignore, index already exists
		} else {
			gf_err := gf_core.Error__create(`failed to create db index on fields - {"t","flows_names_lst","origin_url_str"}`,
				"mongodb_ensure_index_error",nil,err,"gf_images_lib",p_runtime_sys)
			return gf_err
		}
	}

	//DEPRECATED!! - flow_name_str field is deprecated in favor of flows_names_lst, 
	//               but some of the old img records still use it and havent been migrated yet. 
	//               so we're creating this index until migration is complete, and then 
	//               it should be removed.
	//FIELDS  - "t","flow_name_str","origin_url_str"
	//QUERIES - flows_db__images_exist() issues these queries

	doc_type__index = mgo.Index{
		Key:        []string{"t","flow_name_str","origin_url_str"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = p_runtime_sys.Mongodb_coll.EnsureIndex(doc_type__index)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "duplicate key error index") {
			//ignore, index already exists
		} else {
			gf_err := gf_core.Error__create(`failed to create db index on fields - {"t","flow_name_str","origin_url_str"}`,
				"mongodb_ensure_index_error",nil,err,"gf_images_lib",p_runtime_sys)
			return gf_err
		}
	}
	//---------------------

	return nil

}