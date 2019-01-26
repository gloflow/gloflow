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

package gf_crawl_lib

import (
	"fmt"
	"strings"
	"github.com/globalsign/mgo"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
)
//--------------------------------------------------
func db_index__init(p_runtime *gf_crawl_core.Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {

	//---------------------
	//FIELDS - "t"
	doc_type__index := mgo.Index{
	    Key:        []string{"t"}, //all stat queries first match on "t"
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
			gf_err := gf_core.Error__create(`failed to create db index on fields - {"t"}`,
				"mongodb_ensure_index_error", nil, err, "gf_crawl_lib", p_runtime_sys)
			return gf_err
		}
	}
	//---------------------
	//FIELDS - "t","hash_str"
	doc_type__index = mgo.Index{
	    Key:        []string{"t","hash_str"}, //all stat queries first match on "t"
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
			gf_err := gf_core.Error__create(`failed to create db index on fields - {"t","hash_str"}`,
				"mongodb_ensure_index_error", nil, err, "gf_crawl_lib", p_runtime_sys)
			return gf_err
		}
	}
	//---------------------
	return nil
}