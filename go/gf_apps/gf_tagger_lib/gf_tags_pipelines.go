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
)

//---------------------------------------------------
// AUTHORIZED

func pipelineAdd(pInputDataMap map[string]interface{},
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

	object_type_str      := strings.TrimSpace(pInputDataMap["otype"].(string))
	object_extern_id_str := strings.TrimSpace(pInputDataMap["o_id"].(string))
	tags_str             := strings.TrimSpace(pInputDataMap["tags"].(string))


	var metaMap map[string]interface{}
	if _, ok := pInputDataMap["meta_map"]; ok {
		metaMap = pInputDataMap["meta_map"].(map[string]interface{})
	}


	if object_type_str == "address" {
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
	gfErr := addTagsToObject(tags_str,
		object_type_str,
		object_extern_id_str,
		metaMap,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//----------------
	return nil
}

//---------------------------------------------------

func tagsPipelineGetObjects(p_req *http.Request,
	p_resp                   io.Writer,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	pRuntimeSys            *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	//----------------
	// INPUT
	qs_map := p_req.URL.Query()
	var err error

	// responseFormatStr - "json" | "html"
	responseFormatStr := gf_rpc_lib.GetResponseFormat(qs_map, pRuntimeSys)

	if _, ok := qs_map["otype"]; !ok {
		gfErr := gf_core.ErrorCreate("input 'otype' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"qs_map":qs_map,},
			nil, "gf_tagger", pRuntimeSys)
		return nil, gfErr
	}

	if _, ok := qs_map["tag"]; !ok {
		gfErr := gf_core.ErrorCreate("input 'tag' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"qs_map":qs_map,},
			nil, "gf_tagger", pRuntimeSys)
		return nil, gfErr
	}

	// TrimSpace() - Returns the string without any leading and trailing whitespace.
	object_type_str := strings.TrimSpace(qs_map["otype"][0])
	tag_str         := strings.TrimSpace(qs_map["tag"][0])

	// PAGE_INDEX
	page_index_int := 0
	if a_lst, ok := qs_map["pg_index"]; ok {
		input_val          := a_lst[0]
		page_index_int, err = strconv.Atoi(input_val) //user supplied value
		if err != nil {
			gfErr := gf_core.ErrorCreate("input pg_index not an integer",
				"verify__value_not_integer_error",
				map[string]interface{}{"input_val":input_val,},
				nil, "gf_tagger", pRuntimeSys)
			return nil, gfErr
		}
	}

	// PAGE_SIZE
	page_size_int := 10
	if a_lst, ok := qs_map["pg_size"]; ok {
		input_val         := a_lst[0]
		page_size_int, err = strconv.Atoi(input_val) //user supplied value
		if err != nil {
			gfErr := gf_core.ErrorCreate("input pg_size not an integer",
				"verify__value_not_integer_error",
				map[string]interface{}{"input_val":input_val,},
				nil, "gf_tagger", pRuntimeSys)
			return nil, gfErr
		}
	}

	//----------------

	switch responseFormatStr {
		//------------------
		// HTML RENDERING
		case "html":
			pRuntimeSys.LogFun("INFO","HTML RESPONSE >>")
			gfErr := renderObjectsWithTag(tag_str,
				p_tmpl,
				p_subtemplates_names_lst,
				page_index_int,
				page_size_int,
				p_resp,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

		//------------------
		// JSON EXPORT
		case "json":
			pRuntimeSys.LogFun("INFO","JSON RESPONSE >>")
			objectsWithTagLst, gfErr := getObjectsWithTag(tag_str,
				object_type_str,
				page_index_int,
				page_size_int,
				pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr
			}

			// FIX!! - objectsWithTagLst - have to be exported for external use, not just serialized
			//                             from their internal representation

			return objectsWithTagLst, nil

		//------------------
	}
	return nil, nil
}