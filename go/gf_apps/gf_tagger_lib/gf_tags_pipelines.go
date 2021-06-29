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
	"text/template"
	"net/http"
	"io"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//---------------------------------------------------
// AUTHORIZED

func pipeline__add_tags(p_input_data_map map[string]interface{},
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_tagger_pipelines.pipeline__add_tags()")

	//----------------
	// INPUT
	if _,ok := p_input_data_map["otype"]; !ok {
		gf_err := gf_core.Error__create("input 'otype' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map":p_input_data_map,},
			nil, "gf_tagger", p_runtime_sys)
		return gf_err
	}

	if _,ok := p_input_data_map["o_id"]; !ok {
		gf_err := gf_core.Error__create("input 'o_id' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map":p_input_data_map,},
			nil, "gf_tagger", p_runtime_sys)
		return gf_err
	}

	if _,ok := p_input_data_map["tags"]; !ok {
		gf_err := gf_core.Error__create("input 'tags' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map":p_input_data_map,},
			nil, "gf_tagger", p_runtime_sys)
		return gf_err
	}

	object_type_str      := strings.TrimSpace(p_input_data_map["otype"].(string))
	object_extern_id_str := strings.TrimSpace(p_input_data_map["o_id"].(string))
	tags_str             := strings.TrimSpace(p_input_data_map["tags"].(string))

	//----------------
	gf_err := add_tags_to_object(tags_str,
		object_type_str,
		object_extern_id_str,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	//----------------
	return nil
}

//---------------------------------------------------
func pipeline__get_objects_with_tag(p_req *http.Request,
	p_resp                   io.Writer,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_runtime_sys            *gf_core.Runtime_sys) ([]map[string]interface{}, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_tagger_pipelines.pipeline__get_objects_with_tag()")

	//----------------
	// INPUT
	qs_map := p_req.URL.Query()
	var err error

	// response_format_str - "j"(for json)|"h"(for html)
	response_format_str := gf_rpc_lib.Get_response_format(qs_map, p_runtime_sys)

	if _,ok := qs_map["otype"]; !ok {
		gf_err := gf_core.Error__create("input 'otype' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"qs_map":qs_map,},
			nil, "gf_tagger", p_runtime_sys)
		return nil, gf_err
	}

	if _,ok := qs_map["tag"]; !ok {
		gf_err := gf_core.Error__create("input 'tag' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"qs_map":qs_map,},
			nil, "gf_tagger", p_runtime_sys)
		return nil, gf_err
	}

	// TrimSpace() - Returns the string without any leading and trailing whitespace.
	object_type_str := strings.TrimSpace(qs_map["otype"][0])
	tag_str         := strings.TrimSpace(qs_map["tag"][0])

	// PAGE_INDEX
	page_index_int := 0
	if a_lst,ok := qs_map["pg_index"]; ok {
		input_val          := a_lst[0]
		page_index_int, err = strconv.Atoi(input_val) //user supplied value
		if err != nil {
			gf_err := gf_core.Error__create("input pg_index not an integer",
				"verify__value_not_integer_error",
				map[string]interface{}{"input_val":input_val,},
				nil, "gf_tagger", p_runtime_sys)
			return nil, gf_err
		}
	}

	// PAGE_SIZE
	page_size_int := 10
	if a_lst,ok := qs_map["pg_size"]; ok {
		input_val         := a_lst[0]
		page_size_int, err = strconv.Atoi(input_val) //user supplied value
		if err != nil {
			gf_err := gf_core.Error__create("input pg_size not an integer",
				"verify__value_not_integer_error",
				map[string]interface{}{"input_val":input_val,},
				nil, "gf_tagger", p_runtime_sys)
			return nil, gf_err
		}
	}

	//----------------

	switch response_format_str {
		//------------------
		// HTML RENDERING
		case "html":
			p_runtime_sys.Log_fun("INFO","HTML RESPONSE >>")
			gf_err := render_objects_with_tag(tag_str,
				p_tmpl,
				p_subtemplates_names_lst,
				page_index_int,
				page_size_int,
				p_resp,
				p_runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

		//------------------
		// JSON EXPORT
		case "json":
			p_runtime_sys.Log_fun("INFO","JSON RESPONSE >>")
			objects_with_tag_lst, gf_err := get_objects_with_tag(tag_str,
				object_type_str,
				page_index_int,
				page_size_int,
				p_runtime_sys)
			if gf_err != nil {
				return nil, gf_err
			}

			// FIX!! - objects_with_tag_lst - have to be exported for external use, not just serialized
			//                                from their internal representation

			return objects_with_tag_lst, nil

		//------------------
	}
	return nil, nil
}