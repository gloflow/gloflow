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

package gf_crawl_core

import (
	"fmt"
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//--------------------------------------------------
type Gf_crawl_link_alloc struct {
	Id                   bson.ObjectId  `bson:"_id,omitempty"`
	Id_str               string         `bson:"id_str"`
	T_str                string         `bson:"t"`                         //"crawler_link_alloc"
	Creation_unix_time_f float64        `bson:"creation_unix_time_f"`
	Crawler_name_str     string         `bson:"crawler_name_str"`
	Block_size_int       int            `bson:"block_size_int"`
	Sleep_time_sec_int   int            `bson:"sleep_time_sec_int"`

	Last_run_unix_time_f      float64
	Current_link_block_id_str string    `bson:"current_link_block_id_str"`
}

type Gf_crawl_link_alloc_block struct {
	Id                       bson.ObjectId `bson:"_id,omitempty"`
	Id_str                   string        `bson:"id_str"`
	Creation_unix_time_f     float64       `bson:"creation_unix_time_f"`
	T_str                    string        `bson:"t"`                     //"crawler_link_alloc_block"
	Allocator_id_str         string        `bson:"allocator_id_str"`
	Unresolved_links_ids_lst []string      `bson:"unresolved_links_ids_lst"`
}

//--------------------------------------------------
func Link_alloc__init(p_crawler_name_str string, p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_links_allocator.Link_alloc__init()")

	allocator, gf_err := Link_alloc__create(p_crawler_name_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	go func() {
		for ;; {

			gf_err := Link_alloc__run(allocator, p_runtime_sys)
			if gf_err != nil {

			}

			//SLEEP
			sleep_length := time.Second * time.Duration(allocator.Sleep_time_sec_int)
			time.Sleep(sleep_length)
		}
	}()

	return nil
}

//--------------------------------------------------
func Link_alloc__create(p_crawler_name_str string, p_runtime_sys *gf_core.Runtime_sys) (*Gf_crawl_link_alloc, *gf_core.Gf_error) {

	block_size_int     := 100
	sleep_time_sec_int := 60*20 //20min

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("gf_crawl_link_alloc:%s:%f", p_crawler_name_str, creation_unix_time_f)

	//IMPORTANT!! - there can be multiple allocators operating in a single cluster. potentially they can have different allocations strategies,
	//              or may have different limitations on the range of values for various filters used by allocation function.
	allocator := &Gf_crawl_link_alloc{
		Id_str:                    id_str,
		T_str:                     "crawler_link_alloc", 
		Creation_unix_time_f:      creation_unix_time_f, 
		Crawler_name_str:          p_crawler_name_str,
		Block_size_int:            block_size_int,
		Sleep_time_sec_int:        sleep_time_sec_int,
		Current_link_block_id_str: "",
	}

	//DB
	err := p_runtime_sys.Mongodb_db.C("gf_crawl").Insert(allocator)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a crawl_link_alloc in mongodb",
			"mongodb_insert_error",
			map[string]interface{}{
				"crawler_name_str": p_crawler_name_str,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return nil, gf_err
	}
	
	return allocator, nil
}

//--------------------------------------------------
func Link_alloc__run(p_alloc *Gf_crawl_link_alloc, p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {

	alloc_run_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	//BLOCK
	new_block, gf_err := Link_alloc__create_links_block(p_alloc.Id_str, p_alloc.Crawler_name_str, p_alloc.Block_size_int, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
		
	p_alloc.Last_run_unix_time_f      = alloc_run_unix_time_f
	p_alloc.Current_link_block_id_str = new_block.Id_str

	return nil
}

//--------------------------------------------------
func Link_alloc__create_links_block(p_alloc_id_str string,
	p_crawler_name_str string,
	p_block_size_int int,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_crawl_link_alloc_block, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_crawl_links_allocator.Link_alloc__create_links_block()")


	query := p_runtime_sys.Mongodb_db.C("gf_crawl").Find(bson.M{
		"t":                    "crawler_page_outgoing_link",
		"crawler_name_str":     p_crawler_name_str, //get links that were discovered by this crawler
		"valid_for_crawl_bool": true,
		"fetched_bool":         false,

		//IMPORTANT!! - get all unresolved links that also dont have any errors associated
		//              with them. this way rep`eated processing of unresolved links that always cause 
		//              an error is avoided (wasted resources)
		"error_type_str": bson.M{"$exists":false,},
		"error_id_str":   bson.M{"$exists":false,},
	}).
	//IMPORTANT!! - sort by date of link creation/discovery, and get the links that were discovered first,
	//              ascending order of unix timestamps.
	Sort("$creation_unix_time_f: 1").
	Limit(p_block_size_int).
	Select(bson.M{"id_str": 1})


	var unresolved_links_ids_lst []string
	err := query.All(&unresolved_links_ids_lst)

	if err != nil {
		gf_err := gf_core.Error__create("failed to get a block of crawler_page_outgoing_links, to allocate for crawling",
			"mongodb_find_error",
			map[string]interface{}{
				"crawler_name_str": p_crawler_name_str,
				"block_size_int":   p_block_size_int,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return nil, gf_err
	}

	//-------------------
	//CREATE_LINK_ALLOCATOR_BLOCK

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("gf_crawl_link_alloc_block:%s:%f", p_crawler_name_str, creation_unix_time_f)

	block := &Gf_crawl_link_alloc_block{
		Id_str:                   id_str,
		T_str:                    "crawler_link_alloc_block",
		Creation_unix_time_f:     creation_unix_time_f,
		Allocator_id_str:         p_alloc_id_str,
		Unresolved_links_ids_lst: unresolved_links_ids_lst,
	}

	//DB
	err = p_runtime_sys.Mongodb_db.C("gf_crawl").Insert(block)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a crawl_link_alloc_block in mongodb",
			"mongodb_insert_error",
			map[string]interface{}{
				"allocator_id_str": p_alloc_id_str,
			},
			err, "gf_crawl_core", p_runtime_sys)
		return nil, gf_err
	}
	//-------------------

	return block, nil
}

