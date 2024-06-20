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
	"strings"
	"strconv"
	"context"
	"text/template"
	"net/http"
	"io"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_events"
	"github.com/gloflow/gloflow/go/gf_apps/gf_tagger_lib/gf_tagger_core"
)

//---------------------------------------------------

func pipelineGetAllTags(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {




	tagsLst, gfErr := dbMongoGetAllTags(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return tagsLst, nil

}

//---------------------------------------------------
// AUTHORIZED

func pipelineAdd(pInputDataMap map[string]interface{},
	pUserID     gf_core.GF_ID,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//----------------
	// INPUT
	if _, ok := pInputDataMap["otype"]; !ok {
		gfErr := gf_core.ErrorCreate("input 'otype' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map":pInputDataMap,},
			nil, "gf_tagger", pRuntimeSys)
		return gfErr
	}

	if _, ok := pInputDataMap["o_id"]; !ok {
		gfErr := gf_core.ErrorCreate("input 'o_id' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map": pInputDataMap,},
			nil, "gf_tagger", pRuntimeSys)
		return gfErr
	}

	if _, ok := pInputDataMap["tags"]; !ok {
		gfErr := gf_core.ErrorCreate("input 'tags' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map": pInputDataMap,},
			nil, "gf_tagger", pRuntimeSys)
		return gfErr
	}

	objectTypeStr     := strings.TrimSpace(pInputDataMap["otype"].(string))
	objectExternIDstr := strings.TrimSpace(pInputDataMap["o_id"].(string))
	tagsStr           := strings.TrimSpace(pInputDataMap["tags"].(string))


	var metaMap map[string]interface{}
	if _, ok := pInputDataMap["meta_map"]; ok {
		metaMap = pInputDataMap["meta_map"].(map[string]interface{})
	}


	if objectTypeStr == "address" {
		if metaMap == nil {
			gfErr := gf_core.ErrorCreate("tagging objects of type 'address' has to contain meta_map with",
				"verify__missing_key_error",
				map[string]interface{}{"input_data_map": pInputDataMap,},
				nil, "gf_tagger", pRuntimeSys)
			return gfErr
		}

		if _, ok := metaMap["chain_str"]; !ok {
			gfErr := gf_core.ErrorCreate("tagging objects of type 'address' has to contain meta_map with chain_str key",
				"verify__missing_key_error",
				map[string]interface{}{"input_data_map": pInputDataMap,},
				nil, "gf_tagger", pRuntimeSys)
			return gfErr	
		}
	}

	//----------------
	gfErr := addTagsToObject(tagsStr,
		objectTypeStr,
		objectExternIDstr,
		metaMap,
		pUserID,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//----------------

	//------------------
	// EVENT
	if pRuntimeSys.EnableEventsAppBool {
		eventMeta := map[string]interface{}{
			"tags":        tagsStr,
			"object_type": objectTypeStr,
			"object_id":   objectExternIDstr,
		}
		gf_events.EmitApp(gf_tagger_core.GF_EVENT_APP__TAG_ADD,
			eventMeta,
			pUserID,
			pCtx,
			pRuntimeSys)
	}

	//------------------
	return nil
}

//---------------------------------------------------
// PIPELINE_GET_OBJECTS

func tagsPipelineGetObjects(pReq *http.Request,
	pResp                 io.Writer,
	pTemplate             *template.Template,
	pSubtemplatesNamesLst []string,
	pCtx                  context.Context,
	pRuntimeSys           *gf_core.RuntimeSys) ([]map[string]interface{}, map[string]interface{}, *gf_core.GFerror) {

	//----------------
	// INPUT
	qsMap := pReq.URL.Query()
	var err error

	// responseFormatStr - "json" | "html"
	responseFormatStr := gf_rpc_lib.GetResponseFormat(qsMap, pRuntimeSys)

	if _, ok := qsMap["otype"]; !ok {
		/*
		gfErr := gf_core.ErrorCreate("input 'otype' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"qs_map": qsMap,},
			nil, "gf_tagger_lib", pRuntimeSys)
		*/

		abortReasonMap := map[string]interface{}{
			"msg_str":  "input 'otype' not supplied",
			"type_str": "verify__missing_key_error",
			"qs_map":   qsMap,
		}
		return nil, abortReasonMap, nil
	}

	if _, ok := qsMap["tag"]; !ok {
		/*
		gfErr := gf_core.ErrorCreate("input 'tag' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"qs_map": qsMap,},
			nil, "gf_tagger_lib", pRuntimeSys)
		*/
		abortReasonMap := map[string]interface{}{
			"msg_str":  "input 'tag' not supplied",
			"type_str": "verify__missing_key_error",
			"qs_map":   qsMap,
		}
		
		return nil, abortReasonMap, nil
	}

	// TrimSpace() - Returns the string without any leading and trailing whitespace.
	objectTypeStr := strings.TrimSpace(qsMap["otype"][0])
	tagStr        := strings.TrimSpace(qsMap["tag"][0])

	// PAGE_INDEX
	pageIndexInt := 0
	if aLst, ok := qsMap["pg_index"]; ok {
		inputVal := aLst[0]
		pageIndexInt, err = strconv.Atoi(inputVal) // user supplied value
		if err != nil {
			gfErr := gf_core.ErrorCreate("input pg_index not an integer",
				"verify__value_not_integer_error",
				map[string]interface{}{"input_val": inputVal,},
				nil, "gf_tagger_lib", pRuntimeSys)
			return nil, nil, gfErr
		}
	}

	// PAGE_SIZE
	pageSizeInt := 10
	if aLst, ok := qsMap["pg_size"]; ok {
		inputVal := aLst[0]
		pageSizeInt, err = strconv.Atoi(inputVal) // user supplied value
		if err != nil {
			gfErr := gf_core.ErrorCreate("input pg_size not an integer",
				"verify__value_not_integer_error",
				map[string]interface{}{"input_val": inputVal,},
				nil, "gf_tagger_lib", pRuntimeSys)
			return nil, nil, gfErr
		}
	}

	//----------------

	switch responseFormatStr {
		//------------------
		// HTML RENDERING
		case "html":
			pRuntimeSys.LogNewFun("DEBUG", "HTML RESPONSE >>", nil)
			
			templateRenderedStr, gfErr := renderObjectsWithTag(tagStr,
				pTemplate,
				pSubtemplatesNamesLst,
				pageIndexInt,
				pageSizeInt,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, nil, gfErr
			}

			pResp.Write([]byte(templateRenderedStr))

		//------------------
		// JSON EXPORT
		case "json":
			pRuntimeSys.LogNewFun("DEBUG", "JSON RESPONSE >>", nil)
			
			objectsWithTagLst, gfErr := exportObjectsWithTag(tagStr,
				objectTypeStr,
				pageIndexInt,
				pageSizeInt,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, nil, gfErr
			}

			/*
			FIX!! - objectsWithTagLst - have to be exported for external use, not just serialized
				from their internal representation
			*/
			return objectsWithTagLst, nil, nil

		//------------------
	}
	return nil, nil, nil
}