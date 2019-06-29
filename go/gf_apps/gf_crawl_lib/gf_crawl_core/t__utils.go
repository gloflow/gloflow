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
	"testing"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func t__cleanup__test_page_imgs(p_test__crawler_name_str string, p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "t__utils.t__cleanup__test_page_imgs()")

	_,err := p_runtime_sys.Mongodb_db.C("gf_crawl").RemoveAll(bson.M{
			"t":                bson.M{"$in": []string{"crawler_page_img", "crawler_page_img_ref",},},
			"crawler_name_str": p_test__crawler_name_str,
		})
	if err != nil {
		panic(err)
	}
}

//-------------------------------------------------
func T__init(p_test *testing.T) (*gf_core.Runtime_sys, *Gf_crawler_runtime) {

	test__mongodb_host_str      := "127.0.0.1"
	test__mongodb_db_name_str   := "gf_tests"
	test__cluster_node_type_str := "master"
	
	log_fun      := gf_core.Init_log_fun()
	mongodb_db   := gf_core.Mongo__connect(test__mongodb_host_str, test__mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_crawl_tests",
		Log_fun:          log_fun,
		Mongodb_db:       mongodb_db,
		Mongodb_coll:     mongodb_coll,
	}
	//-------------
	//ELASTICSEARCH
	esearch_client, gf_err := gf_core.Elastic__get_client(runtime_sys)
	if gf_err != nil {
		p_test.Errorf("failed to get ElasticSearch client in test initialization")
		return nil, nil
	}
	//-------------
	//S3
	s3_test_info := gf_core.T__get_s3_info(runtime_sys)
	//-------------

	crawler_runtime := &Gf_crawler_runtime{
		Events_ctx:            nil,
		Esearch_client:        esearch_client,
		S3_info:               s3_test_info.Gf_s3_info,
		Cluster_node_type_str: test__cluster_node_type_str,
	}

	return runtime_sys, crawler_runtime
}