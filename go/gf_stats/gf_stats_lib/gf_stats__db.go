/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

package gf_stats_lib

import (
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type Stat_db_coll struct {
	Docs_count_int             int64
	Size_bytes_int             int64
	Avg_obj_size_bytes_int     int64
	Storage_size_bytes_int     int64
	Index_total_size_bytes_int int64
}

//-------------------------------------------------

func Db_stats__coll(p_coll_name_str string,
	p_ctx         context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*Stat_db_coll, *gf_core.GFerror) {

	db   := pRuntimeSys.Mongo_db
	coll := db.Collection(p_coll_name_str)



	// DOCS_COUNT
	count_int, err := coll.CountDocuments(p_ctx, bson.M{})
	if err != nil {

		gf_err := gf_core.MongoHandleError(fmt.Sprintf("failed to DB count records in collection - %s", p_coll_name_str),
			"mongodb_count_error",
			map[string]interface{}{},
			err, "gf_stats_lib", pRuntimeSys)
		return nil, gf_err
	}






	r := db.RunCommand(p_ctx, bson.M{"collStats": p_coll_name_str})

	var r_map bson.M
	err = r.Decode(&r_map)

	if err !=nil {
		gf_err := gf_core.MongoHandleError(fmt.Sprintf("failed to DB decode collStats for collection - %s", p_coll_name_str),
			"mongodb_cursor_decode",
			map[string]interface{}{},
			err, "gf_stats_lib", pRuntimeSys)
		return nil, gf_err
	}



	coll_size__bytes_int             := r_map["size"]
	coll_avg_obj_size__bytes_int     := r_map["avgObjSize"]
	coll_storage_size__bytes_int     := r_map["storageSize"]
	coll_total_index_size__bytes_int := r_map["totalIndexSize"]





	db_coll_stats := &Stat_db_coll{
		Docs_count_int:             count_int,
		Size_bytes_int:             int64(coll_size__bytes_int.(int32)),
		Avg_obj_size_bytes_int:     int64(coll_avg_obj_size__bytes_int.(int32)),
		Storage_size_bytes_int:     int64(coll_storage_size__bytes_int.(int32)),
		Index_total_size_bytes_int: int64(coll_total_index_size__bytes_int.(int32)),
	}

	return db_coll_stats, nil
}