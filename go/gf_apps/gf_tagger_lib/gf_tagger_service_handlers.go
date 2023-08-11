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
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib/gf_tagger_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib/gf_bookmarks"
)

//-------------------------------------------------

func initHandlers(pTemplatesPathsMap map[string]string,
	pImagesJobsMngr gf_images_jobs_core.JobsMngr,
	pMux            *http.ServeMux,
	pRuntimeSys     *gf_core.RuntimeSys) *gf_core.GFerror {

	// TEMPLATES
	gfTemplates, gfErr := gf_tagger_core.TemplatesLoad(pTemplatesPathsMap, pRuntimeSys)
	if gfErr != nil {
		return gfErr
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
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//------------------
				// INPUT

				// FIX!! - this route should authenticated and this user_id
				//         should be fetched from the handler context.
				userID := gf_core.GF_ID("gf")

				inputMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var input gf_bookmarks.GFbookmarkInputCreate
				err := mapstructure.Decode(inputMap, &input)
				if err != nil {
					gfErr := gf_core.ErrorCreate("failed to load http input into GFbookmarkInputCreate struct",
						"mapstruct__decode",
						map[string]interface{}{},
						err, "gf_tagger_lib", pRuntimeSys)
					return nil, gfErr
				}

				input.User_id_str = "anonymous"

				//------------------

				gfErr = gf_bookmarks.PipelineCreate(&input,
					userID,
					pImagesJobsMngr,
					pCtx,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}
			}

			return nil, nil
		},
		pMux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	// CREATE
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/bookmarks/get",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			//------------------
			// INPUT
			
			qs_map := pReq.URL.Query()

			// response_format_str - "j"(for json) | "h"(for html)
			response_format_str := gf_rpc_lib.GetResponseFormat(qs_map, pRuntimeSys)

			
			input := &gf_bookmarks.GFbookmarkInputGet{
				Response_format_str: response_format_str,
				User_id_str:         "anonymous",
			}

			//------------------

			output, gfErr := gf_bookmarks.PipelineGet(input,
				gfTemplates.Bookmarks,
				gfTemplates.BookmarksSubtemplatesNamesLst,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}


			switch response_format_str { 
			case "json":
				data_map := map[string]interface{}{
					"bookmarks_lst": output.Bookmarks_lst,
				}
				return data_map, nil
		
			case "html":

				pResp.Write([]byte(output.TemplateRenderedStr))
				return nil, nil
			}


			return nil, nil

		},
		pMux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// NOTES
	//---------------------
	// CREATE
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/tags/notes/create",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {
			if pReq.Method == "POST" {

				//------------
				// INPUT
				iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//------------
	
				gfErr = notesPipelineAdd(iMap, pRuntimeSys)
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
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// GET_NOTES

	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/tags/notes/get",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				notesLst, gfErr := notesPipelineGet(pReq, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr 
				}

				dataMap := map[string]interface{}{"notes_lst": notesLst,}
				return dataMap, nil
			}

			return nil, nil
		},
		pMux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------
	// TAGS
	//---------------------
	// ADD_TAGS
	
	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/tags/create",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "POST" {

				//------------
				// INPUT
				iMap, gfErr := gf_core.HTTPgetInput(pReq, pRuntimeSys)
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
		true, // pStoreRunBool
		nil,
		pRuntimeSys)
	
	//---------------------
	// GET_OBJECTS_WITH_TAG

	gf_rpc_lib.CreateHandlerHTTPwithMux("/v1/tags/objects",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if pReq.Method == "GET" {

				objectsWithTagLst, gfErr := tagsPipelineGetObjects(pReq, pResp, 
					gfTemplates.TagObjects,
					gfTemplates.TagObjectsSubtemplatesNamesLst,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				// if the response_format was HTML then objects_with_tag_lst is nil,
				// in which case there is no json to send back
				if objectsWithTagLst != nil {

					dataMap := map[string]interface{}{"objects_with_tag_lst": objectsWithTagLst,}
					return dataMap, nil
				} else {
					return nil, nil
				}
			}

			return nil, nil
		},
		pMux,
		metrics,
		true, // pStoreRunBool
		nil,
		pRuntimeSys)

	//---------------------

	return nil
}