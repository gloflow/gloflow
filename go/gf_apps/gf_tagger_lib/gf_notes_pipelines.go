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
	"time"
	"strconv"
	"strings"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFnote struct {
	UserIDstr           string `json:"user_id_str"         bson:"user_id_str"`         //user_id of the user that attached this note
	BodyStr             string `json:"body_str"            bson:"body_str"`
	TargetObjIDstr      string `json:"target_obj_id_str"   bson:"target_obj_id_str"`   //object_id to which this note is attached
	TargetObjTypeStr    string `json:"target_obj_type_str" bson:"target_obj_type_str"` //"post"|"image"|"video"
	CreationDatetimeStr string `json:"creation_datetime_str"`
}

//---------------------------------------------------

func notesPipelineAdd(p_input_data_map map[string]interface{},
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	//----------------
	// INPUT
	if _, ok := p_input_data_map["otype"]; !ok {
		gfErr := gf_core.ErrorCreate("note 'otype' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map": p_input_data_map,},
			nil, "gf_tagger", pRuntimeSys)
		return gfErr
	}

	if _, ok := p_input_data_map["o_id"]; !ok {
		gfErr := gf_core.ErrorCreate("note 'o_id' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map": p_input_data_map,},
			nil, "gf_tagger", pRuntimeSys)
		return gfErr
	}

	if _, ok := p_input_data_map["body"]; !ok {
		gfErr := gf_core.ErrorCreate("note 'body' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"input_data_map": p_input_data_map,},
			nil, "gf_tagger", pRuntimeSys)
		return gfErr
	}

	object_type_str      := strings.TrimSpace(p_input_data_map["otype"].(string))
	objectExternIDstr := strings.TrimSpace(p_input_data_map["o_id"].(string))
	body_str             := strings.TrimSpace(p_input_data_map["body"].(string))

	//----------------

	if object_type_str == "post" {

		post_title_str        := objectExternIDstr
		creation_datetime_str := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0, 'f', 10, 64)

		note := &GFnote{
			UserIDstr:           "anonymous",
			BodyStr:             body_str,
			TargetObjIDstr:      post_title_str,
			TargetObjTypeStr:    object_type_str,
			CreationDatetimeStr: creation_datetime_str,
		}

		gfErr := db__add_post_note(note, post_title_str, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}
	return nil
}

//---------------------------------------------------

func notesPipelineGet(pReq *http.Request,
	pRuntimeSys *gf_core.RuntimeSys) ([]*GFnote, *gf_core.GFerror) {

	//-----------------
	// INPUT
	qsMap := pReq.URL.Query()

	if _, ok := qsMap["otype"]; !ok {
		gfErr := gf_core.ErrorCreate("note 'otype' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"qs_map": qsMap,},
			nil, "gf_tagger", pRuntimeSys)
		return nil, gfErr
	}

	if _, ok := qsMap["o_id"]; !ok {
		gfErr := gf_core.ErrorCreate("note 'o_id' not supplied",
			"verify__missing_key_error",
			map[string]interface{}{"qs_map": qsMap,},
			nil, "gf_tagger", pRuntimeSys)
		return nil, gfErr
	}

	objectTypeStr     := strings.TrimSpace(qsMap["otype"][0])
	objectExternIDstr := strings.TrimSpace(qsMap["o_id"][0])

	//-----------------

	taggerNotesLst := []*GFnote{}
	if objectTypeStr == "post" {

		postTitleStr    := objectExternIDstr
		notesLst, gfErr := db__get_post_notes(postTitleStr, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
		
		for _, s := range notesLst {
			note := &GFnote{
				UserIDstr:        s.UserIDstr,
				BodyStr:          s.BodyStr,
				TargetObjIDstr:   postTitleStr,
				TargetObjTypeStr: objectTypeStr,
			}
			taggerNotesLst = append(taggerNotesLst, note)
		}
	}
	return taggerNotesLst, nil
}