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
	"github.com/fatih/color"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//--------------------------------------------------
func Image__db_create(p_img *Gf_crawler_page_img,
	p_runtime     *Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.Gf_error) {
	//p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_db.Image__db_create()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//------------
	//MASTER
	if p_runtime.Cluster_node_type_str == "master" {

		c,err := p_runtime_sys.Mongodb_coll.Find(bson.M{
				"t":        "crawler_page_img",
				"hash_str": p_img.Hash_str,
			}).Count()
		if err != nil {
			gf_err := gf_core.Error__create("failed to count the number of crawler_page_img's in mongodb",
				"mongodb_find_error",
				&map[string]interface{}{
					"img_ref_url_str":             p_img.Url_str,
					"img_ref_origin_page_url_str": p_img.Origin_page_url_str,
				},
				err, "gf_crawl_core", p_runtime_sys)
			return false, gf_err
		}

		//crawler_page_img already exists, from previous crawls, so ignore it
		var exists_bool bool
		if c > 0 {
			p_runtime_sys.Log_fun("INFO",yellow(">>>>>>>> - DB PAGE_IMAGE ALREADY EXISTS >-- ")+cyan(p_img.Url_str))
			
			exists_bool = true
			return exists_bool, nil
		} else {
				
			//IMPORTANT!! - only insert the crawler_page_img if it doesnt exist in the DB already
			err = p_runtime_sys.Mongodb_coll.Insert(p_img)
			if err != nil {
				gf_err := gf_core.Error__create("failed to insert a crawler_page_img in mongodb",
					"mongodb_insert_error",
					&map[string]interface{}{
						"img_ref_url_str":             p_img.Url_str,
						"img_ref_origin_page_url_str": p_img.Origin_page_url_str,
					},
					err, "gf_crawl_core", p_runtime_sys)
				return false, gf_err
			}

			exists_bool = false
			return exists_bool, nil
		}
	}
	//------------
	//WORKER
	if p_runtime.Cluster_node_type_str == "worker" {

		//ADD!! - issue a HTTP request for this data to a remote 'master' node
	}
	//------------

	return false,nil
}
//--------------------------------------------------
func Image__db_create_ref(p_img_ref *Gf_crawler_page_img_ref,
	p_runtime     *Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	//p_log_fun("FUN_ENTER","gf_crawl_images_db.Image__db_create_ref()")

	if p_runtime.Cluster_node_type_str == "master" {

		c,err := p_runtime_sys.Mongodb_coll.Find(bson.M{
				"t":       "crawler_page_img_ref",
				"hash_str":p_img_ref.Hash_str,
			}).Count()
		if err != nil {
			gf_err := gf_core.Error__create("failed to count the number of crawler_page_img_ref's in mongodb",
				"mongodb_find_error",
				&map[string]interface{}{
					"img_ref_url_str":             p_img_ref.Url_str,
					"img_ref_origin_page_url_str": p_img_ref.Origin_page_url_str,
				},
				err, "gf_crawl_core", p_runtime_sys)
			return gf_err
		}

		//crawler_page_img already exists, from previous crawls, so ignore it
		if c > 0 {
			return nil
		} else {
				
			//IMPORTANT!! - only insert the crawler_page_img if it doesnt exist in the DB already
			err = p_runtime_sys.Mongodb_coll.Insert(p_img_ref)
			if err != nil {
				gf_err := gf_core.Error__create("failed to insert a crawler_page_img_ref in mongodb",
					"mongodb_insert_error",
					&map[string]interface{}{
						"img_ref_url_str":            p_img_ref.Url_str,
						"img_ref_origin_page_url_str":p_img_ref.Origin_page_url_str,
					},
					err, "gf_crawl_core", p_runtime_sys)
				return gf_err
			}
		}
	} else {

	}

	return nil
}
//--------------------------------------------------
func image__db_get(p_id_str string,
	p_runtime     *Gf_crawler_runtime,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_crawler_page_img,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_crawl_images_db.image__db_get()")

	var img Gf_crawler_page_img
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{
			"t":     "crawler_page_img",
			"id_str":p_id_str,
		}).One(&img)
	if err != nil {
		gf_err := gf_core.Error__create("failed to get crawler_page_img by ID from mongodb",
			"mongodb_find_error",
			&map[string]interface{}{"id_str":p_id_str,},
			err, "gf_crawl_core", p_runtime_sys)
		return nil, gf_err
	}

	return &img, nil
}