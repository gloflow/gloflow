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

package gf_publisher_core

import (
	"fmt"
	"time"
	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//------------------------------------------------
type Gf_post struct {
	Id                    primitive.ObjectID `bson:"_id,omitempty"`
	Id_str                string         `bson:"id_str"`
	T_str                 string         `bson:"t"`                //"post"
	Deleted_bool          bool           `bson:"deleted_bool"`
	Client_type_str       string         `bson:"client_type_str"`  //"gchrome_ext" //type of the client that created the post
	Title_str             string         `bson:"title_str"`
	Description_str       string         `bson:"description_str"`
	Creation_datetime_str string         `bson:"creation_datetime_str"`
	Poster_user_name_str  string         `bson:"poster_user_name_str"` //user-name of the user that posted this post

	//------------
	// GF_IMAGES
	Thumbnail_url_str string                        `bson:"thumbnail_url_str"` //SYMPHONY 0.3
	Images_ids_lst    []gf_images_core.Gf_image_id `bson:"images_ids_lst"`
	
	//------------
	Post_elements_lst []*Gf_post_element `bson:"post_elements_lst"`
	
	//------------
	Tags_lst  []string                   `bson:"tags_lst"`
	Notes_lst []*Gf_post_note            `bson:"notes_lst"`      //SYMPHONY 0.3 - notes are chunks of text (for now) that can be attached to a post
 
	//every event can have multiple colors assigned to it
	Colors_lst []string                  `bson:"colors_lst"`
}

type Gf_post_note struct {
	User_id_str           string `bson:"user_id_str"`
	Body_str              string `bson:"body_str"`
	Creation_datetime_str string `bson:"creation_datetime_str"`
}

//------------------------------------------------
func Create_new_post(p_post_info_map map[string]interface{}, p_runtime_sys *gf_core.RuntimeSys) (*Gf_post, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_post.Create_new_post()")
	p_runtime_sys.LogFun("INFO",      "p_post_info_map - "+fmt.Sprint(p_post_info_map))

	// IMPORTANT!! - not all posts have "tags_lst" element, check if this is fine or if should be enforced
	// assert(p_post_info_map.containsKey("tags_lst"));
	// assert(p_post_info_map["tags_lst"] is List);

	post_title_str := p_post_info_map["title_str"].(string)
	
	//--------------------
	// POST ELEMENTS

	post_elements_infos_lst := p_post_info_map["post_elements_lst"].([]interface{})
	post_elements_lst       := create_post_elements(post_elements_infos_lst, post_title_str, p_runtime_sys)
	p_runtime_sys.LogFun("INFO","post_elements_lst - "+fmt.Sprint(post_elements_lst))

	//--------------------
	// CREATION DATETIME

	// strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0,'f',10,64)
	creation_datetime_str := time.Now().String()

	//--------------------
	// FIX!! - id_str should be a hash, just as events/images have their ID"s generated as hashes
	
	id_str := fmt.Sprintf("%s:%s", post_title_str, creation_datetime_str)

	//--------------------
	// OPTIONAL VALUES

	// THUMBNAIL_URL
	var thumbnail_url_str string
	if _, ok := p_post_info_map["thumbnail_url_str"]; !ok {
		thumbnail_url_str = ""
	} else {
		thumbnail_url_str = p_post_info_map["thumbnail_url_str"].(string)
	}

	// IMAGES_IDS
	var images_ids_lst []gf_images_core.Gf_image_id
	if _, ok := p_post_info_map["images_ids_lst"]; !ok {

		// "images_ids_lst" key was not present
		images_ids_lst = []gf_images_core.Gf_image_id{}
	} else if ids_lst,ok := p_post_info_map["images_ids_lst"].([]gf_images_core.Gf_image_id); !ok {
		// if p_post_info_map is coming from mongodb, its of []interface{} type, so a "type conversion" is done)
		images_ids_lst = []gf_images_core.Gf_image_id(ids_lst)
	} else {
		images_ids_lst = ids_lst
	}

	//-------------------------
	// TAGS
	var tags_lst []string

	if _,ok := p_post_info_map["tags_lst"]; !ok {
		// if "tags_lst" key is missing, assign an empty string
		tags_lst = []string{}
	} else if input_tags_lst,ok := p_post_info_map["tags_lst"].([]string); ok {
		// "tags_lst" is present in p_post_info_map and is of type []string
		tags_lst = input_tags_lst

		// //CAITION!! - if p_post_info_map is coming from mongodb, its of []interface{} type, so a "type conversion" is done.
		// //            is this ever coming from mongodb? posts are deserialized by the mgo lib, not in this function ever. 
		// tags_lst = []string(tags_lst)
	} else {
		// "tags_lst" is not of type []string
		tags_lst = p_post_info_map["tags_lst"].([]string)
	}

	//-------------------------
	// NOTES
	var notes_lst []*Gf_post_note
	if _,ok := p_post_info_map["notes_lst"]; !ok {
		notes_lst = []*Gf_post_note{}
	} else {
		notes_infos_lst := p_post_info_map["notes_lst"].([]map[string]interface{})
		notes_lst        = create_post_notes(notes_infos_lst, p_runtime_sys)
	}

	//-------------------------
	// COLORS
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
	post := &Gf_post{
		Id_str:                id_str,
		T_str:                 "post",
		Deleted_bool:          false,
		Client_type_str:       p_post_info_map["client_type_str"].(string),
		Title_str:             post_title_str,
		Description_str:       p_post_info_map["description_str"].(string),
		Creation_datetime_str: creation_datetime_str,
		Poster_user_name_str:  p_post_info_map["poster_user_name_str"].(string),
		Thumbnail_url_str:     thumbnail_url_str,
		Images_ids_lst:        images_ids_lst,
		Post_elements_lst:     post_elements_lst,
		Tags_lst:              tags_lst,
		Notes_lst:             notes_lst,
		Colors_lst:            colors_lst,
	}
	return post, nil
}

//------------------------------------------------	
// a post has to first be created, and only then can it be published

func publish(p_post_title_str string, p_runtime_sys *gf_core.RuntimeSys) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_post.publish()")
}

//------------------------------------------------
func create_post_notes(p_raw_notes_lst []map[string]interface{}, p_runtime_sys *gf_core.RuntimeSys) []*Gf_post_note {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_post.create_post_notes()")

	notes_lst := []*Gf_post_note{}
	for _, note_map := range p_raw_notes_lst {
		
		snippet := &Gf_post_note{
			User_id_str: "anonymous",
			Body_str:    note_map["body_str"].(string),
		}
		notes_lst = append(notes_lst, snippet)
	}
	return notes_lst
}