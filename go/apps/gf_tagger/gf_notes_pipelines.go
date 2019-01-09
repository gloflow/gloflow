package main

import (
	"fmt"
	"time"
	"strconv"
	"errors"
	"strings"
	"net/http"
	"gopkg.in/mgo.v2"
)

//---------------------------------------------------
type Note struct {
	User_id_str           string `json:"user_id_str"         bson:"user_id_str"`         //user_id of the user that attached this note
	Body_str              string `json:"body_str"            bson:"body_str"`
	Target_obj_id_str     string `json:"target_obj_id_str"   bson:"target_obj_id_str"`   //object_id to which this note is attached
	Target_obj_type_str   string `json:"target_obj_type_str" bson:"target_obj_type_str"` //"post"|"image"|"video"
	Creation_datetime_str string `json:"creation_datetime_str"`
}
//---------------------------------------------------
func pipeline__add_note(p_input_data_map map[string]interface{},
	p_mongo_coll *mgo.Collection,
	p_log_fun    func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_notes_pipelines.pipeline__add_note()")

	//----------------
	//INPUT
	if _,ok := p_input_data_map["otype"]; !ok {
		return errors.New("note 'otype' not supplied")
	}

	if _,ok := p_input_data_map["o_id"]; !ok {
		return errors.New("note 'o_id' not supplied")
	}

	if _,ok := p_input_data_map["body"]; !ok {
		return errors.New("note 'body' not supplied")
	}

	object_type_str      := strings.TrimSpace(p_input_data_map["otype"].(string))
	object_extern_id_str := strings.TrimSpace(p_input_data_map["o_id"].(string))
	body_str             := strings.TrimSpace(p_input_data_map["body"].(string))
	//----------------

	if  object_type_str == "post" {

		post_title_str        := object_extern_id_str
		creation_datetime_str := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0,'f',10,64)

		note := &Note{
			User_id_str:          "anonymous",
			Body_str:             body_str,
			Target_obj_id_str:    post_title_str,
			Target_obj_type_str:  object_type_str,
			Creation_datetime_str:creation_datetime_str,
		}

		fmt.Println(">>>>>>>>>>>>>>> =============")
		fmt.Println(note)

		err := db__add_post_note(note,
			&post_title_str,
			p_mongo_coll,
			p_log_fun)
		if err != nil {
			return err
		}

	} else {
		return errors.New("object_type_str ["+object_type_str+"] not implemented yet")
	}
	return nil
}
//---------------------------------------------------
func pipeline__get_notes(p_req *http.Request,
	p_mongo_coll *mgo.Collection,
	p_log_fun    func(string,string)) ([]*Note,error) {
	p_log_fun("FUN_ENTER","gf_notes_pipelines.pipeline__get_notes()")

	//-----------------
	//INPUT

	qs_map := p_req.URL.Query()

	if _,ok := qs_map["otype"]; !ok {
		return nil,errors.New("'otype' not supplied")
	}

	if _,ok := qs_map["o_id"]; !ok {
		return nil,errors.New("'o_id' not supplied")
	}

	object_type_str      := strings.TrimSpace(qs_map["otype"][0])
	object_extern_id_str := strings.TrimSpace(qs_map["o_id"][0])
	//-----------------

	tagger_notes_lst := []*Note{}
	if object_type_str == "post" {

		var err error
		post_title_str   := object_extern_id_str
		notes_lst,err := db__get_post_notes(&post_title_str, p_mongo_coll, p_log_fun)

		for _,s := range notes_lst {
			note := &Note{
				User_id_str:        s.User_id_str,
				Body_str:           s.Body_str,
				Target_obj_id_str:  post_title_str,
				Target_obj_type_str:object_type_str,
			}
			tagger_notes_lst = append(tagger_notes_lst,note)

		}

		if err != nil {
			return nil,err
		}
	} else {
		return nil,errors.New("object_type_str ["+object_type_str+"] not implemented yet")
	}

	return tagger_notes_lst,nil
}