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

func initHandlers(pMediaDomainStr string,
	p_crawled_images_s3_bucket_name_str string,
	p_gf_images_s3_bucket_name_str      string,
	p_templates_paths_map               map[string]string,
	pHTTPmux                          *http.ServeMux,
	pRuntime                            *gf_crawl_core.GFcrawlerRuntime,
	pRuntimeSys                         *gf_core.RuntimeSys) *gf_core.GFerror {
	
	//---------------------
	// TEMPLATES

	gfTemplates, gfErr := tmpl__load(p_templates_paths_map, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	
	//----------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/crawl/image/recent",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//------------------
				recentImagesLst, gfErr := gf_crawl_core.ImagesDBgetRecent(pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				data_map := map[string]interface{}{
					"recent_images_lst": recentImagesLst,
				}
				return data_map, nil

				//------------------
			}
			return nil, nil
		},
		pHTTPmux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//----------------
	// FIX!! - move out of the crawler
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/crawl/image/add_to_flow",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//--------------------------
				// INPUT
				i, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				crawler_page_image_id_str := i["crawler_page_image_id_str"].(string)

				flows_names_lst := []string{}
				for _, s := range i["flows_names_lst"].([]interface{}) {
					flows_names_lst = append(flows_names_lst, s.(string))
				}

				//--------------------------
				gfErr = gf_crawl_core.FlowsAddExternImage(gf_crawl_core.GFcrawlerPageImageID(crawler_page_image_id_str),
					flows_names_lst,
					pMediaDomainStr,
					p_crawled_images_s3_bucket_name_str,
					p_gf_images_s3_bucket_name_str,
					pCtx,
					pRuntime,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------------
				// OUTPUT
				output_map := map[string]interface{}{}
				return output_map, nil

				//------------------
			}
			return nil, nil
		},
		pHTTPmux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//----------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/crawl/search",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {
				
				query_term_str := pReq.URL.Query()["term"][0]
				pRuntimeSys.LogFun("INFO", "query_term_str - "+query_term_str)

				// IMPORTANT!! - only query if the indexer is enabled
				if pRuntime.EsearchClient != nil {
					gfErr := gf_crawl_core.IndexQuery(query_term_str, pRuntime, pRuntimeSys)
					if gfErr != nil {
						return nil, gfErr
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
		pHTTPmux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//----------------
	gf_rpc_lib.CreateHandlerHTTPwithMux("/a/crawl/crawl_dashboard",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				//--------------------
				// RENDER TEMPLATE
				gfErr := dashboard__render_template(gfTemplates.dashboard__tmpl,
					gfTemplates.dashboard__subtemplates_names_lst,
					pResp,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
				return nil, nil
			}
			return nil, nil
		},
		pHTTPmux,
		nil,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)
	
	//--------------

	return nil
}