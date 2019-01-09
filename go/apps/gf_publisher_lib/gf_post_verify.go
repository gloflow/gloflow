package gf_publisher_lib

import (
	"errors"
	"fmt"
	"strings"
	"github.com/globalsign/mgo"
	"gf_core"
)
//---------------------------------------------------
//external post_info is the one that comes from outside the system
//(it does not have an id assigned to it)

func verify_external_post_info(p_post_info_map map[string]interface{},
	p_max_title_chars_int       int, //100
	p_max_description_chars_int int, //1000
	p_post_element_tag_max_int  int, //20
	p_mongodb_coll              *mgo.Collection,
	p_log_fun                   func(string,string)) (map[string]interface{},error) {
	p_log_fun("FUN_ENTER","gf_post_verify.verify_external_post_info()")
		
	//default_main_image_url_str := "http://gloflow.com/gf_publisher/objects/images/gf_landingpage_3.png"

	//-------------------
	//TYPE
	if _,ok := p_post_info_map["client_type_str"]; !ok {
		return nil,errors.New("post client_type_str not supplied")
	}
	//-------------------
	//TITLE

	if _,ok := p_post_info_map["title_str"]; !ok {
		return nil,errors.New("post title_str not supplied")
	}
	title_str := p_post_info_map["title_str"].(string)

	if len(title_str) > p_max_title_chars_int {
		return nil,errors.New(fmt.Sprintf("title_str is longer (%s) then the max allowed number of chars (%s)", len(title_str), p_max_title_chars_int))
	}

	//ATTENTION!!
	//FB is removing/having problems with these symbols in url endings, and since the url to posts is composed of 
	//the post title, FB breaks these links
	//so striping them off right here avoids that

	clean_title_str   := title_str
	replace_chars_lst := []string{"[",",",":","#","%","&","!","]","$"}
	for _,c := range replace_chars_lst {
		strings.Replace(clean_title_str,c,"",-1)
	}
	//-------------------
	//DESCRIPTION
		
	if _,ok := p_post_info_map["description_str"]; !ok {
		return nil,errors.New("post description_str not supplied")
	}
	description_str := p_post_info_map["description_str"].(string)

	if len(description_str) > p_max_description_chars_int {
		return nil,errors.New(fmt.Sprintf("description_str is longer (%s) then the max allowed number of chars (%s)", len(description_str), p_max_description_chars_int))
	}
	//-------------------
	//POST ELEMENTS
	err := verify_post_elements(p_post_info_map, p_post_element_tag_max_int, p_log_fun)
	if err != nil {
		return nil,err
	}
	//-------------------	
	//TAGS
	tags_lst,err := verify_tags(p_post_info_map, p_log_fun)
	if err != nil {
		return nil,err
	}
	//-------------------
	if _,ok := p_post_info_map["poster_user_name_str"]; !ok {
		return nil,errors.New("post poster_user_name_str not supplied")
	}

	if _,ok := p_post_info_map["post_elements_lst"]; !ok {
		return nil,errors.New("post post_elements_lst not supplied")
	}

	//"id_str" - not included here since p_post_info_map comes from outside the system
	//           and the internal id"s are for now not passed outside (or coming in from outside)
	verified_post_info_map := map[string]interface{}{
		"client_type_str":     p_post_info_map["client_type_str"].(string),
		"title_str":           clean_title_str,
		"description_str":     description_str,
		"poster_user_name_str":p_post_info_map["poster_user_name_str"].(string),
		"post_elements_lst":   p_post_info_map["post_elements_lst"],
		"tags_lst":            tags_lst,
	}
	
	return verified_post_info_map,nil
}
//---------------------------------------------------
func verify_tags(p_post_info_map map[string]interface{}, p_log_fun func(string,string)) ([]string,error) { 
	p_log_fun("FUN_ENTER","gf_post_verify.verify_tags()")
		
	if _,ok := p_post_info_map["tags_str"]; !ok {
		return nil,errors.New("p_post_info_map doesnt contain the tags_str key")
	}

	input_tags_str := p_post_info_map["tags_str"].(string)
	tags_lst       := strings.Split(input_tags_str," ")

	p_log_fun("INFO","input_tags_str - "+fmt.Sprint(input_tags_str))
	p_log_fun("INFO","tags_lst       - "+fmt.Sprint(tags_lst))

	return tags_lst,nil
}
//---------------------------------------------------
func verify_post_elements(p_post_info_map map[string]interface{},
	p_post_element_tag_max_int int,
	p_log_fun                  func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_post_verify.verify_post_elements()")
	
	if _,ok := p_post_info_map["post_elements_lst"]; !ok {
		return errors.New("p_post_info_map doesnt contain the post_elements_lst key")
	}
	post_elements_lst := p_post_info_map["post_elements_lst"].([]interface{})

	//verify each individiaul post_element
	for _,post_element := range post_elements_lst {
		post_element_map := post_element.(map[string]interface{})
		err := verify_post_element(post_element_map, p_post_element_tag_max_int, p_log_fun)
		if err != nil {
			return err
		}

		//------------------------
		//SECURITY
		//ADD!! - have a external-url checking routines/whitelists/blacklists
		//        and other url sanitization routines,
		//        to prevent various XSS attacks
		//------------------------
	}

	return nil
}
//---------------------------------------------------
func verify_post_element(p_post_element_info_map map[string]interface{},
	p_post_element_tag_max_int int, //20
	p_log_fun                  func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_post_verify.verify_post_element()")
	p_log_fun("INFO"     ,"p_post_element_info_map - "+fmt.Sprint(p_post_element_info_map))

	post_element_type_str := p_post_element_info_map["type_str"].(string)
	
	if !(post_element_type_str == "link"  ||
		post_element_type_str == "image"  ||
		post_element_type_str == "video"  ||
		post_element_type_str == "iframe" ||
		post_element_type_str == "text") {

		return errors.New("post_element_type_str is not link|image|video|iframe|text - "+post_element_type_str)
	}
	
	//--------------
	if (post_element_type_str == "link"  ||
		post_element_type_str == "image" ||
		post_element_type_str == "video" ||
		post_element_type_str == "iframe") {	 

		//FIX!! - newe versions of post_element_info_dict format use extern_url_str
		//        instead of url_str. so when all post"s in the DB are updated to this format
		//        remove p_post_element_info_dict.containsKey("url_str") from this assert
		if !(gf_core.Map_has_key(p_post_element_info_map,"url_str") ||
			gf_core.Map_has_key(p_post_element_info_map,"extern_url_str")) {
			return errors.New("p_post_element_info_map doesnt contain url_str|extern_url_str")
		}
	}
	//--------------
	//TAGS       

	if pe_tags_lst,ok := p_post_element_info_map["tags_lst"]; ok {
		for _,tag_str := range pe_tags_lst.([]string) {

			if len(tag_str) <= p_post_element_tag_max_int {
				return errors.New(fmt.Sprintf("tag (%s) is longer then max chars per tag (%d)", tag_str, p_post_element_tag_max_int))
			}
		}
	}
	//--------------

	return nil
}