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

package gf_crawl_lib

import (
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
)

//--------------------------------------------------
type Gf_crawler_cluster_worker struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               string        `bson:"id_str"               json:"id_str"`
	T_str                string        `bson:"t"                    json:"t"`            // "crawler_cluster_worker"
	Creation_unix_time_f float64       `bson:"creation_unix_time_f" json:"creation_unix_time_f"`
	Ext_name_str         string        `bson:"name_str"             json:"ext_name_str"` // externally-supplied worker name
}

type Gf_json_msg__link__get_unresolved struct {
	Link_id_str           string  `json:"link_id_str"`
	Fetch_id_str          string  `json:"fetch_id_str"`
	Fetch_creation_time_f float64 `json:"fetch_creation_time_f"`
}

//--------------------------------------------------
func cluster__client(p_req_type_str string,
	p_runtime   *gf_crawl_core.GFcrawlerRuntime,
	pRuntimeSys *gf_core.Runtime_sys) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_crawl_cluster.cluster__client()")
	switch p_req_type_str {
		case "register_worker":

		case "create__page_img":

		case "link__get_unresolved":

		case "link__mark_as_resolved":
	}
}

//--------------------------------------------------
func cluster__register_worker(p_ext_worker_name_str string,
	p_runtime   *gf_crawl_core.GFcrawlerRuntime,
	pRuntimeSys *gf_core.Runtime_sys) (*Gf_crawler_cluster_worker, *gf_core.GFerror) {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_crawl_cluster.cluster__register_worker()")

	id_str               := "crawler_cluster_worker__"+fmt.Sprint()
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	worker := &Gf_crawler_cluster_worker{
		Id_str:               id_str,
		T_str:                "crawler_cluster_worker",
		Creation_unix_time_f: creation_unix_time_f,
		Ext_name_str:         p_ext_worker_name_str,
	}

	//------------
	// DB

	ctx           := context.Background()
	coll_name_str := "gf_crawl"
	gf_err        := gf_core.Mongo__insert(worker,
		coll_name_str,
		map[string]interface{}{
			"ext_worker_name_str": p_ext_worker_name_str,
			"caller_err_msg_str":  "failed to insert a Gf_crawler_cluster_worker into the DB in order to register it",
		},
		ctx,
		pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}

	/*err := pRuntimeSys.Mongodb_db.C("gf_crawl").Insert(worker)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a Gf_crawler_cluster_worker in mongodb in order to register it",
			"mongodb_insert_error",
			map[string]interface{}{
				"ext_worker_name_str": p_ext_worker_name_str,
			},
			err, "gf_crawl_lib", pRuntimeSys)
		return nil, gf_err
	}*/

	//------------

	return worker, nil
}

//--------------------------------------------------
func cluster__init_handlers(p_crawl_config_file_path_str string,
	p_runtime *gf_crawl_core.GFcrawlerRuntime,
	pRuntimeSys *gf_core.Runtime_sys) *gf_core.GFerror {
	pRuntimeSys.Log_fun("FUN_ENTER", "gf_crawl_cluster.cluster__init_handlers()")
	
	crawlers_map, gf_err := gf_crawl_core.Get_all_crawlers(p_crawl_config_file_path_str, pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}
	
	//----------------
	// REGISTER_WORKER
	http.HandleFunc("/a/crawl/cluster/register__worker", func(p_resp http.ResponseWriter, p_req *http.Request) {

		worker_name_str := p_req.URL.Query()["worker_name"][0]
		pRuntimeSys.Log_fun("INFO","worker_name_str - "+worker_name_str)

		worker, gf_err := cluster__register_worker(worker_name_str, p_runtime, pRuntimeSys)
		if gf_err != nil {
			return
		}

		//------------------
		// OUTPUT
		
		r_map := map[string]interface{}{
			"status_str":    "OK",
			"worker_id_str": worker.Id_str,
		}

		r_lst,_ := json.Marshal(r_map)
		r_str   := string(r_lst)
		fmt.Fprint(p_resp, r_str)

		//------------------
	})

	//----------------
	// IMAGES
	http.HandleFunc("/a/crawl/cluster/create__page_imgs", func(p_resp http.ResponseWriter, p_req *http.Request) {
		pRuntimeSys.Log_fun("INFO","INCOMING HTTP REQUEST -- /a/crawl/cluster/create__page_imgs ----------")
		if p_req.Method == "POST" {

			worker_id_str := p_req.URL.Query()["worker_id"][0]
			pRuntimeSys.Log_fun("INFO", "worker_id_str - "+worker_id_str)

			var imgs_lst []gf_crawl_core.Gf_crawler_page_image
			body_bytes_lst,_ := ioutil.ReadAll(p_req.Body)
			err              := json.Unmarshal(body_bytes_lst, &imgs_lst)
			if err != nil {
				panic(err)
				return
			}

			imgs_existed_lst := []bool{}
			for _,img := range imgs_lst {
				img_existed_bool, gf_err := gf_crawl_core.Image__db_create(&img, p_runtime, pRuntimeSys)
				if gf_err != nil {
					return
				}

				imgs_existed_lst = append(imgs_existed_lst, img_existed_bool)
			}

			//------------------
			// OUTPUT
			
			r_map := map[string]interface{}{
				"status_str":       "OK",
				"imgs_existed_lst": imgs_existed_lst,
			}

			r_lst,_ := json.Marshal(r_map)
			r_str   := string(r_lst)
			fmt.Fprint(p_resp,r_str)

			//------------------
		}
	})

	http.HandleFunc("/a/crawl/cluster/create__page_img_ref", func(p_resp http.ResponseWriter, p_req *http.Request) {
		pRuntimeSys.Log_fun("INFO", "INCOMING HTTP REQUEST -- /a/crawl/cluster/create__page_img_ref ----------")
		if p_req.Method == "POST" {

			var imgs_refs_lst []gf_crawl_core.Gf_crawler_page_image_ref
			body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
			err               := json.Unmarshal(body_bytes_lst,&imgs_refs_lst)
			if err != nil {
				panic(err)
				return
			}

			for _, img_ref := range imgs_refs_lst {
				gf_err := gf_crawl_core.Image__db_create_ref(&img_ref, p_runtime, pRuntimeSys)
				if gf_err != nil {
					return
				}
			}
		}
	})
	//-----------------
	// LINKS
	http.HandleFunc("/a/crawl/cluster/link__get_unresolved", func(p_resp http.ResponseWriter, p_req *http.Request) {
		pRuntimeSys.Log_fun("INFO", "INCOMING HTTP REQUEST -- /a/crawl/cluster/link__get_unresolved ----------")
		if p_req.Method == "GET" {

			crawler_name_str := p_req.URL.Query()["crawler_name_str"][0]
			pRuntimeSys.Log_fun("INFO","crawler_name_str - "+crawler_name_str)

			if crawler, ok := crawlers_map[crawler_name_str]; ok {

				// domains_lst := crawler.Domains_lst

				unresolved_link, gf_err := gf_crawl_core.Link__db_get_unresolved(crawler.Name_str, pRuntimeSys)
				if gf_err != nil {
					return
				}
				//------------------
				// OUTPUT
				
				r_map := map[string]interface{}{
					"status_str":      "OK",
					"unresolved_link": unresolved_link,
				}

				r_lst, _ := json.Marshal(r_map)
				r_str    := string(r_lst)
				fmt.Fprint(p_resp, r_str)

				//------------------
			} else {

			}
		}
	})

	http.HandleFunc("/a/crawl/cluster/link__mark_as_resolved", func(p_resp http.ResponseWriter, p_req *http.Request) {

		pRuntimeSys.Log_fun("INFO", "INCOMING HTTP REQUEST -- /a/crawl/cluster/link__mark_as_resolved ----------")
		if p_req.Method == "GET" {

			//---------------------
			// INPUT
			var input Gf_json_msg__link__get_unresolved
			body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
			err               := json.Unmarshal(body_bytes_lst, &input)
			if err != nil {
				panic(err)
				return 
			}

			//---------------------

			link, gf_err := gf_crawl_core.Link__db_get(input.Link_id_str, pRuntimeSys)
			if gf_err != nil {
				return
			}
			
			gf_err = gf_crawl_core.Link__db_mark_as_resolved(link,
				input.Fetch_id_str,
				input.Fetch_creation_time_f,
				pRuntimeSys)
			if gf_err != nil {
				return
			}

			//------------------
			// OUTPUT
			
			r_map := map[string]interface{}{
				"status_str": "OK",
			}

			r_lst,_ := json.Marshal(r_map)
			r_str   := string(r_lst)
			fmt.Fprint(p_resp, r_str)

			//------------------
		}
	})
	
	//-----------------
	return nil
}