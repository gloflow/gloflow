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
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_crawl_lib/gf_crawl_core"
)

//-------------------------------------------------
func init_handlers(p_media_domain_str string,
	p_crawled_images_s3_bucket_name_str string,
	p_gf_images_s3_bucket_name_str string,
	p_templates_paths_map          map[string]string,
	p_http_mux                     *http.ServeMux,
	p_runtime                      *gf_crawl_core.GFcrawlerRuntime,
	p_runtime_sys                  *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//---------------------
	// TEMPLATES

	gf_templates, gf_err := tmpl__load(p_templates_paths_map, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//----------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/crawl/image/recent",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//------------------
				recent_images_lst, gf_err := gf_crawl_core.Images__db_get_recent(p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//------------------
				// OUTPUT
				data_map := map[string]interface{}{
					"recent_images_lst": recent_images_lst,
				}
				return data_map, nil

				//------------------
			}
			return nil, nil
		},
		p_http_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)

	//----------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/crawl/image/add_to_flow",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//--------------------------
				// INPUT
				i, gf_err := gf_core.HTTPgetInput(p_resp, p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				crawler_page_image_id_str := i["crawler_page_image_id_str"].(string)

				flows_names_lst := []string{}
				for _, s := range i["flows_names_lst"].([]interface{}) {
					flows_names_lst = append(flows_names_lst, s.(string))
				}

				//--------------------------
				gf_err = gf_crawl_core.FlowsAddExternImage(gf_crawl_core.Gf_crawler_page_image_id(crawler_page_image_id_str),
					flows_names_lst,

					p_media_domain_str,
					p_crawled_images_s3_bucket_name_str,
					p_gf_images_s3_bucket_name_str,
					p_runtime,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//------------------
				// OUTPUT
				output_map := map[string]interface{}{}
				return output_map, nil

				//------------------
			}
			return nil, nil
		},
		p_http_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)

	//----------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/crawl/search",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {
				
				query_term_str := p_req.URL.Query()["term"][0]
				p_runtime_sys.Log_fun("INFO", "query_term_str - "+query_term_str)

				// IMPORTANT!! - only query if the indexer is enabled
				if p_runtime.Esearch_client != nil {
					gf_err := gf_crawl_core.Index__query(query_term_str, p_runtime, p_runtime_sys)
					if gf_err != nil {
						return nil, gf_err
					}
				}
				//------------------
				// OUTPUT
				output_map := map[string]interface{}{}
				return output_map, nil

				//------------------
			}
			return nil, nil
		},
		p_http_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)

	//----------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/crawl/crawl_dashboard",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "GET" {

				//--------------------
				// RENDER TEMPLATE
				gf_err := dashboard__render_template(gf_templates.dashboard__tmpl,
					gf_templates.dashboard__subtemplates_names_lst,
					p_resp,
					p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
				return nil, nil
			}
			return nil, nil
		},
		p_http_mux,
		nil,
		true, // p_store_run_bool
		nil,
		p_runtime_sys)
	
	//--------------

	return nil
}