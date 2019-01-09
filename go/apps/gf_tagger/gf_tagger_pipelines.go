package main

import (
	"errors"
	"strings"
	"strconv"
	"text/template"
	"net/http"
	"gopkg.in/mgo.v2"
	
	"gf_rpc_lib"
)
//---------------------------------------------------
//AUTHORIZED

func pipeline__add_tags(p_input_data_map map[string]interface{},
	p_mongodb_coll *mgo.Collection,
	p_log_fun      func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_tagger_pipelines.pipeline__add_tags()")

	//----------------
	//INPUT
	if _,ok := p_input_data_map["otype"]; !ok {
		return errors.New("snippet 'otype' not supplied")
	}

	if _,ok := p_input_data_map["o_id"]; !ok {
		return errors.New("snippet 'o_id' not supplied")
	}

	if _,ok := p_input_data_map["tags"]; !ok {
		return errors.New("'tags' not supplied")
	}

	object_type_str      := strings.TrimSpace(p_input_data_map["otype"].(string))
	object_extern_id_str := strings.TrimSpace(p_input_data_map["o_id"].(string))
	tags_str             := strings.TrimSpace(p_input_data_map["tags"].(string))
	//----------------

	err := add_tags_to_object(tags_str,
		object_type_str,
		&object_extern_id_str,
		p_mongodb_coll,
		p_log_fun)
	if err != nil {
		return err
	}

	return nil
}
//---------------------------------------------------
func pipeline__get_objects_with_tag(p_req *http.Request,
	p_resp         http.ResponseWriter,
	p_tmpl         *template.Template,
	p_mongodb_coll *mgo.Collection,
	p_log_fun      func(string,string)) ([]map[string]interface{},error) {
	p_log_fun("FUN_ENTER","gf_tagger_pipelines.pipeline__get_objects_with_tag()")

	//----------------
	//INPUT

	qs_map := p_req.URL.Query()

	//response_format_str - "j"(for json)|"h"(for html)
	response_format_str := gf_rpc_lib.Get_response_format(qs_map, p_log_fun)
	if _,ok := qs_map["otype"]; !ok {
		return nil,errors.New("snippet 'otype' not supplied")
	}

	if _,ok := qs_map["tag"]; !ok {
		return nil,errors.New("'tag' not supplied")
	}

	//TrimSpace() - Returns the string without any leading and trailing whitespace.
	object_type_str := strings.TrimSpace(qs_map["otype"][0])
	tag_str         := strings.TrimSpace(qs_map["tag"][0]  )

	//PAGE_INDEX
	page_index_int := 0
	if a_lst,ok := qs_map["pg_index"]; ok {
		page_index_int,_ = strconv.Atoi(a_lst[0]) //user supplied value
		/*if err != nil {
			return nil,err
		}*/
	}

	//PAGE_SIZE
	page_size_int := 10
	if a_lst,ok := qs_map["pg_size"]; ok {
		page_size_int,_ = strconv.Atoi(a_lst[0]) //user supplied value
		/*if err != nil {
			return nil,err
		}*/
	}
	//----------------

	switch response_format_str {
		//------------------
		//HTML RENDERING
		case "html":
			p_log_fun("INFO","HTML RESPONSE >>")
			err := render_objects_with_tag(tag_str,
				p_tmpl,
				page_index_int,
				page_size_int,
				p_mongodb_coll,
				p_resp,
				p_log_fun)
			if err != nil {
				return nil,err
			}
		//------------------
		//JSON EXPORT
		
		case "json":
			p_log_fun("INFO","JSON RESPONSE >>")
			objects_with_tag_lst,err := get_objects_with_tag(tag_str,
				object_type_str,
				page_index_int,
				page_size_int,
				p_mongodb_coll,
				p_log_fun)
			if err != nil {
				return nil,err
			}

			//FIX!! - objects_with_tag_lst - have to be exported for external use, not just serialized
			//                                from their internal representation

			return objects_with_tag_lst,nil
		//------------------
	}

	return nil,nil
}