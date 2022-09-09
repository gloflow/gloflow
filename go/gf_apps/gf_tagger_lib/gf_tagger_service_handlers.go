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

package gf_tagger_lib

import (
	"context"
	"net/http"
	"github.com/mitchellh/mapstructure"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------
func initHandlers(pTemplatesPathsMap map[string]string,
	pImagesJobsMngr gf_images_jobs_core.JobsMngr,
	pMux            *http.ServeMux,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	// TEMPLATES
	gf_templates, gf_err := tmpl__load(pTemplatesPathsMap, pRuntimeSys)
	if gf_err != nil {
		return gf_err
	}

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/bookmarks/create",
		"/v1/bookmarks/get",
		"/v1/tags/notes/create",
		"/v1/tags/notes/get",
		"/v1/tags/create",
		"/v1/tags/objects",
	}
	metricsGroupNameStr := "main"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, "gf_tagger", handlers_endpoints_lst)
	
	//---------------------
	// BOOKMARKS
	//---------------------
	// CREATE
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/bookmarks/create",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			if p_req.Method == "POST" {

				//------------------
				// INPUT
				input_map, gf_err := gf_core.HTTPgetInput(p_req, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				var input GF_bookmark__input_create
				err := mapstructure.Decode(input_map, &input)
				if err != nil {
					gf_err := gf_core.ErrorCreate("failed to load http input into GF_bookmark__input_create struct",
						"mapstruct__decode",
						map[string]interface{}{},
						err, "gf_tagger_lib", pRuntimeSys)
					return nil, gf_err
				}

				input.User_id_str = "anonymous"

				//------------------

				gf_err = bookmarks__pipeline__create(&input,
					pImagesJobsMngr,
					pCtx,
					pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}
			}

			return nil, nil
		},
		pMux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)


	// CREATE
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/bookmarks/get",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			//------------------
			// INPUT
			
			qs_map := p_req.URL.Query()

			// response_format_str - "j"(for json) | "h"(for html)
			response_format_str := gf_rpc_lib.Get_response_format(qs_map, pRuntimeSys)

			
			input := &GF_bookmark__input_get{
				Response_format_str: response_format_str,
				User_id_str:         "anonymous",
			}

			//------------------

			output, gf_err := bookmarks__pipeline__get(input,
				gf_templates.bookmarks__tmpl,
				gf_templates.bookmarks__subtemplates_names_lst,
				pCtx,
				pRuntimeSys)
			if gf_err != nil {
				return nil, gf_err
			}


			switch response_format_str { 
			case "json":
				data_map := map[string]interface{}{
					"bookmarks_lst": output.Bookmarks_lst,
				}
				return data_map, nil
		
			case "html":

				p_resp.Write([]byte(output.Template_rendered_str))
				return nil, nil
			}


			return nil, nil

		},
		pMux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//---------------------
	// NOTES
	//---------------------
	// CREATE
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/tags/notes/create",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {
			if p_req.Method == "POST" {

				//------------
				// INPUT
				i_map, gf_err := gf_core.HTTPgetInput(p_req, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				//------------
	
				gf_err = notes__pipeline__add(i_map, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}
				
				data_map := map[string]interface{}{}
				return data_map, nil
			}

			return nil, nil
		},
		pMux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//---------------------
	// GET_NOTES

	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/tags/notes/get",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			if p_req.Method == "GET" {

				notes_lst, gf_err := notes__pipeline__get(p_req, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err 
				}

				data_map := map[string]interface{}{"notes_lst": notes_lst,}
				return data_map, nil
			}

			return nil, nil
		},
		pMux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//---------------------
	// TAGS
	//---------------------
	// ADD_TAGS
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/tags/create",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if p_req.Method == "POST" {

				//------------
				// INPUT
				iMap, gfErr := gf_core.HTTPgetInput(p_req, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------

				gfErr = pipelineAdd(iMap, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				dataMap := map[string]interface{}{}
				return dataMap, nil
			}
			
			return nil, nil
		},
		pMux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)
	
	//---------------------
	// GET_OBJECTS_WITH_TAG

	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/tags/objects",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.Gf_error) {

			if p_req.Method == "GET" {

				objects_with_tag_lst, gf_err := tags__pipeline__get_objects(p_req, p_resp, 
					gf_templates.tag_objects__tmpl,
					gf_templates.tag_objects__subtemplates_names_lst,
					pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				// if the response_format was HTML then objects_with_tag_lst is nil,
				// in which case there is no json to send back
				if objects_with_tag_lst != nil {

					data_map := map[string]interface{}{"objects_with_tag_lst": objects_with_tag_lst,}
					return data_map, nil
				} else {
					return nil, nil
				}
			}

			return nil, nil
		},
		pMux,
		metrics,
		true, // p_store_run_bool
		nil,
		pRuntimeSys)

	//---------------------

	return nil
}