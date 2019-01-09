package gf_publisher_lib

import (
	"fmt"
	"time"
	"github.com/globalsign/mgo/bson"
)
//------------------------------------------------
type Post struct {
	Id                    bson.ObjectId `bson:"_id,omitempty"`
	Id_str                string        `bson:"id_str"`
	T_str                 string        `bson:"t"`                //"post"
	Client_type_str       string        `bson:"client_type_str"`  //"gchrome_ext" //type of the client that created the post
	Title_str             string        `bson:"title_str"`
	Description_str       string        `bson:"description_str"`
	Creation_datetime_str string        `bson:"creation_datetime_str"`
	Poster_user_name_str  string        `bson:"poster_user_name_str"` //user-name of the user that posted this post

	//------------
	//GF_IMAGES
	Thumbnail_url_str string            `bson:"thumbnail_url_str"` //SYMPHONY 0.3
	Images_ids_lst    []string          `bson:"images_ids_lst"`
	//------------
	Post_elements_lst []*Post_element   `bson:"post_elements_lst"`
	//------------
	Tags_lst  []string                  `bson:"tags_lst"`
	Notes_lst []*Post_note              `bson:"notes_lst"`      //SYMPHONY 0.3 - notes are chunks of text (for now) that can be attached to a post
 
	//every event can have multiple colors assigned to it
	Colors_lst []string                 `bson:"colors_lst"`
}

type Post_note struct {
	User_id_str           string `bson:"user_id_str"`
	Body_str              string `bson:"body_str"`
	Creation_datetime_str string `bson:"creation_datetime_str"`
}
//------------------------------------------------
func create_new_post(p_post_info_map map[string]interface{},
		p_log_fun func(string,string)) (*Post,error) {
	p_log_fun("FUN_ENTER","gf_post.create_new_post()")
	p_log_fun("INFO"     ,"p_post_info_map - "+fmt.Sprint(p_post_info_map))

	//IMPORTANT!! - not all posts have "tags_lst" element, check if this is fine or if should be enforced
	//assert(p_post_info_map.containsKey("tags_lst"));
	//assert(p_post_info_map["tags_lst"] is List);

	post_title_str := p_post_info_map["title_str"].(string)
	//--------------------
	//POST ELEMENTS

	post_elements_infos_lst := p_post_info_map["post_elements_lst"].([]interface{})
	post_elements_lst       := create_post_elements(post_elements_infos_lst,
											&post_title_str,
											p_log_fun)
	p_log_fun("INFO","post_elements_lst - "+fmt.Sprint(post_elements_lst))
	//--------------------
	//CREATION DATETIME

	//strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0,'f',10,64)
	creation_datetime_str := time.Now().String()
	//--------------------
	//FIX!! - id_str should be a hash, just as events/images have their ID"s generated as hashes
	
	id_str := fmt.Sprintf("%s:%s",post_title_str,creation_datetime_str)
	//--------------------
	//OPTIONAL VALUES

	//THUMBNAIL_URL
	var thumbnail_url_str string
	if _,ok := p_post_info_map["thumbnail_url_str"]; !ok {
		thumbnail_url_str = ""
	} else {
		thumbnail_url_str = p_post_info_map["thumbnail_url_str"].(string)
	}

	//IMAGES_IDS
	var images_ids_lst []string
	if _,ok := p_post_info_map["images_ids_lst"]; !ok {

		//"images_ids_lst" key was not present
		images_ids_lst = []string{}
	} else if ids_lst,ok := p_post_info_map["images_ids_lst"].([]string); !ok {
		//if p_post_info_map is coming from mongodb, its of []interface{} type, so a "type conversion" is done)
		images_ids_lst = []string(ids_lst)
	} else {
		images_ids_lst = ids_lst
	}
	//-------------------------
	//TAGS
	var tags_lst []string

	if _,ok := p_post_info_map["tags_lst"]; !ok {
		//if "tags_lst" key is missing, assign an empty string
		tags_lst = []string{}
	} else if input_tags_lst,ok := p_post_info_map["tags_lst"].([]string); ok {
		//"tags_lst" is present in p_post_info_map and is of type []string
		tags_lst = input_tags_lst

		////CAITION!! - if p_post_info_map is coming from mongodb, its of []interface{} type, so a "type conversion" is done.
		////            is this ever coming from mongodb? posts are deserialized by the mgo lib, not in this function ever. 
		//tags_lst = []string(tags_lst)
	} else {
		//"tags_lst" is not of type []string
		tags_lst = p_post_info_map["tags_lst"].([]string)
	}
	//-------------------------
	//NOTES
	var notes_lst []*Post_note
	if _,ok := p_post_info_map["notes_lst"]; !ok {
		notes_lst = []*Post_note{}
	} else {
		notes_infos_lst := p_post_info_map["notes_lst"].([]map[string]interface{})
		notes_lst        = create_post_notes(notes_infos_lst,
										p_log_fun)
	}
	//-------------------------
	//COLORS
	var colors_lst []string
	if _,ok := p_post_info_map["colors_lst"]; !ok {
		colors_lst = []string{}
	} else if tags_lst,ok := p_post_info_map["colors_lst"].([]string); !ok {
		//if p_post_info_map is coming from mongodb, its of []interface{} type, so a "type conversion" is done)
		colors_lst = []string(tags_lst)
	} else {
		colors_lst = p_post_info_map["colors_lst"].([]string)
	}
	//--------------------
	post := &Post{
		Id_str               :id_str,
		T_str                :"post",
		Client_type_str      :p_post_info_map["client_type_str"].(string),
		Title_str            :post_title_str,
		Description_str      :p_post_info_map["description_str"].(string),
		Creation_datetime_str:creation_datetime_str,
		Poster_user_name_str :p_post_info_map["poster_user_name_str"].(string),
		Thumbnail_url_str    :thumbnail_url_str,
		Images_ids_lst       :images_ids_lst,
		Post_elements_lst    :post_elements_lst,
		Tags_lst             :tags_lst,
		Notes_lst            :notes_lst,
		Colors_lst           :colors_lst,
	}
	
	return post,nil
}
//------------------------------------------------	
//a post has to first be created, and only then can it be published

func publish(p_post_title_str *string,
		p_log_fun func(string,string)) {
	p_log_fun("FUN_ENTER","gf_post.publish()")
}
//------------------------------------------------
func create_post_notes(p_raw_notes_lst []map[string]interface{},
				p_log_fun func(string,string)) []*Post_note {
	p_log_fun("FUN_ENTER","gf_post.create_post_notes()")

	notes_lst := []*Post_note{}
	for _,note_map := range p_raw_notes_lst {
		
		snippet := &Post_note{
			User_id_str:"anonymous",
			Body_str   :note_map["body_str"].(string),
		}
		notes_lst = append(notes_lst,snippet)
	}

	return notes_lst
}