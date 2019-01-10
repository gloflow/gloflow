package main

import (
	"errors"
	"fmt"
	"strings"
	"gopkg.in/mgo.v2"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib"
	"github.com/gloflow/gloflow/go/apps/gf_publisher_lib"
)
//---------------------------------------------------
//p_tags_str      - :String - "," separated list of strings
//p_object_id_str - :String - this is an external identifier for an object, not necessarily its internal. 
//                            for posts - their p_object_extern_id_str is their Title, but internally they have
//                                        another ID.

func add_tags_to_object(p_tags_str string,
	p_object_type_str      string,
	p_object_extern_id_str *string,
	p_mongodb_coll         *mgo.Collection,
	p_log_fun              func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_tagger.add_tags_to_object()")

	if p_object_type_str != "post" &&
		p_object_type_str != "image" &&
		p_object_type_str != "event" {

		return errors.New("p_object_type_str ("+p_object_type_str+") is not of supported type (post|image|event)")
	}
	
	tags_lst,err := parse_tags(p_tags_str,
		500, //p_max_tags_bulk_size_int        int, //500
		20,  //p_max_tag_characters_number_int int, //20	
		p_log_fun)
	if err != nil {
		return err
	}

	p_log_fun("INFO","tags_lst - "+fmt.Sprint(tags_lst))
	//---------------
	//POST
	
	switch p_object_type_str {
		//---------------
		//POST
		case "post":
			post_title_str := p_object_extern_id_str
			exists_bool,_  := gf_publisher_lib.DB__check_post_exists(post_title_str, p_mongodb_coll, p_log_fun)
			if exists_bool {
				p_log_fun("INFO","POST EXISTS")
				err := db__add_tags_to_post(post_title_str, tags_lst, p_mongodb_coll, p_log_fun)
				return err
			} else {
				return errors.New(fmt.Sprintf("post with title (%s) doesnt exist", post_title_str))
			}

		//---------------
		//IMAGE
		case "image":
			image_id_str := p_object_extern_id_str
			
			exists_bool,_ := gf_images_lib.DB__image_exists(image_id_str, p_mongodb_coll, p_log_fun)
			if exists_bool {
				err := db__add_tags_to_image(image_id_str, tags_lst, p_mongodb_coll, p_log_fun)
				if err != nil {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("image with id (%s) doesnt exist",image_id_str))
			}
		//---------------
	}

	return nil
}
//---------------------------------------------------
func get_objects_with_tags(p_tags_lst []string,
	p_object_type_str string,
	p_page_index_int  int,
	p_page_size_int   int,
	p_mongodb_coll    *mgo.Collection,
	p_log_fun         func(string,string)) (map[string][]map[string]interface{},error) {
	p_log_fun("FUN_ENTER","gf_tagger.get_objects_with_tags()")
		
	objects_with_tags_map := map[string][]map[string]interface{}{}
	for _,tag_str := range p_tags_lst {

		objects_with_tag_lst,err := get_objects_with_tag(tag_str,
			p_object_type_str,
			p_page_index_int,
			p_page_size_int,
			p_mongodb_coll,
			p_log_fun)

		if err != nil {
			return nil,err
		}
		objects_with_tags_map[tag_str] = objects_with_tag_lst
	}
	return objects_with_tags_map,nil
}
//---------------------------------------------------
func get_objects_with_tag(p_tag_str string,
	p_object_type_str string,
	p_page_index_int  int,
	p_page_size_int   int,
	p_mongodb_coll    *mgo.Collection,
	p_log_fun         func(string,string)) ([]map[string]interface{},error) {
	p_log_fun("FUN_ENTER","gf_tagger.get_objects_with_tag()")
	p_log_fun("INFO"     ,"p_object_type_str - "+p_object_type_str)

	if p_object_type_str != "post" {
		return nil,errors.New("p_object_type_str is not 'post'")
	}

	posts_with_tag_lst,err := db__get_posts_with_tag(p_tag_str,
		p_page_index_int,
		p_page_size_int,
		p_mongodb_coll,
		p_log_fun)
	if err != nil {
		return nil,err
	}

	//package up info of each post that was found with tag 
	min_posts_infos_lst := []map[string]interface{}{}
	for _,post := range posts_with_tag_lst {
		
		post_info_map := map[string]interface{}{
			"title_str":              post.Title_str,
			"tags_lst":               post.Tags_lst,
			"url_str":                "/posts/"+post.Title_str,
			"object_type_str":        p_object_type_str,
			"thumbnail_small_url_str":post.Thumbnail_url_str,
		}

		min_posts_infos_lst = append(min_posts_infos_lst,post_info_map)
	}

	objects_infos_lst := min_posts_infos_lst
	return objects_infos_lst,nil
}
//---------------------------------------------------
func parse_tags(p_tags_str string,
	p_max_tags_bulk_size_int        int, //500
	p_max_tag_characters_number_int int, //20
	p_log_fun                       func(string,string)) ([]string,error) {
	p_log_fun("FUN_ENTER","gf_tagger.parse_tags()")
	
	tags_lst := strings.Split(p_tags_str," ")
	//---------------------
	if len(tags_lst) > p_max_tags_bulk_size_int {
		return nil,errors.New("too many tags supplied - max is "+fmt.Sprint(p_max_tags_bulk_size_int))
	}
	//---------------------
	for _,tag_str := range tags_lst {
		if len(tag_str) > p_max_tag_characters_number_int {
			return nil,errors.New(fmt.Sprintf("tag (%s) is too long - max is (%s)", tag_str, fmt.Sprint(p_max_tags_bulk_size_int)))
		}
	}
	//---------------------
	return tags_lst,nil
}