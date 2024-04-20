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

type GFpost struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty"`
	IDstr               string         `bson:"id_str"`
	Tstr                string         `bson:"t"`                // "post"
	DeletedBool         bool           `bson:"deleted_bool"`
	ClientTypeStr       string         `bson:"client_type_str"`  // "gchrome_ext" //type of the client that created the post
	TitleStr            string         `bson:"title_str"`
	DescriptionStr      string         `bson:"description_str"`
	CreationDatetimeStr string         `bson:"creation_datetime_str"`
	PosterUserNameStr   string         `bson:"poster_user_name_str"` // user-name of the user that posted this post

	//------------
	// GF_IMAGES
	ThumbnailURLstr string                     `bson:"thumbnail_url_str"` // SYMPHONY 0.3
	ImagesIDsLst    []gf_images_core.GFimageID `bson:"images_ids_lst"`
	
	//------------
	PostElementsLst []*GFpostElement `bson:"post_elements_lst"`
	
	//------------
	// TAGGING

	TagsLst  []string      `bson:"tags_lst"`
	NotesLst []*GFpostNote `bson:"notes_lst"` // SYMPHONY 0.3 - notes are chunks of text (for now) that can be attached to a post
	
	//------------
	// COLORS

	// every event can have multiple colors assigned to it
	ColorsLst []string `bson:"colors_lst"`

	//------------
}

type GFpostNote struct {
	UserIDstr           string `bson:"user_id_str"`
	BodyStr             string `bson:"body_str"`
	CreationDatetimeStr string `bson:"creation_datetime_str"`
}

//------------------------------------------------

func CreateNewPost(pPostInfoMap map[string]interface{}, pRuntimeSys *gf_core.RuntimeSys) (*GFpost, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_post.Create_new_post()")
	pRuntimeSys.LogFun("INFO",      "pPostInfoMap - "+fmt.Sprint(pPostInfoMap))

	// IMPORTANT!! - not all posts have "tags_lst" element, check if this is fine or if should be enforced
	// assert(pPostInfoMap.containsKey("tags_lst"));
	// assert(pPostInfoMap["tags_lst"] is List);

	post_title_str := pPostInfoMap["title_str"].(string)
	
	//--------------------
	// POST ELEMENTS

	post_elements_infos_lst := pPostInfoMap["post_elements_lst"].([]interface{})
	post_elements_lst       := create_post_elements(post_elements_infos_lst, post_title_str, pRuntimeSys)
	pRuntimeSys.LogFun("INFO","post_elements_lst - "+fmt.Sprint(post_elements_lst))

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
	if _, ok := pPostInfoMap["thumbnail_url_str"]; !ok {
		thumbnail_url_str = ""
	} else {
		thumbnail_url_str = pPostInfoMap["thumbnail_url_str"].(string)
	}

	// IMAGES_IDS
	var images_ids_lst []gf_images_core.Gf_image_id
	if _, ok := pPostInfoMap["images_ids_lst"]; !ok {

		// "images_ids_lst" key was not present
		images_ids_lst = []gf_images_core.Gf_image_id{}
	} else if ids_lst,ok := pPostInfoMap["images_ids_lst"].([]gf_images_core.Gf_image_id); !ok {
		// if pPostInfoMap is coming from mongodb, its of []interface{} type, so a "type conversion" is done)
		images_ids_lst = []gf_images_core.Gf_image_id(ids_lst)
	} else {
		images_ids_lst = ids_lst
	}

	//-------------------------
	// TAGS
	var tagsLst []string

	if _, ok := pPostInfoMap["tags_lst"]; !ok {
		// if "tags_lst" key is missing, assign an empty string
		tagsLst = []string{}
	} else if input_tags_lst,ok := pPostInfoMap["tags_lst"].([]string); ok {
		// "tags_lst" is present in pPostInfoMap and is of type []string
		tagsLst = input_tags_lst

		// //CAITION!! - if pPostInfoMap is coming from mongodb, its of []interface{} type, so a "type conversion" is done.
		// //            is this ever coming from mongodb? posts are deserialized by the mgo lib, not in this function ever. 
		// tags_lst = []string(tags_lst)
	} else {
		// "tags_lst" is not of type []string
		tagsLst = pPostInfoMap["tags_lst"].([]string)
	}

	//-------------------------
	// NOTES
	var notesLst []*GFpostNote
	if _, ok := pPostInfoMap["notes_lst"]; !ok {
		notesLst = []*GFpostNote{}
	} else {
		notesInfosLst := pPostInfoMap["notes_lst"].([]map[string]interface{})
		notesLst       = createPostNotes(notesInfosLst, pRuntimeSys)
	}

	//-------------------------
	// COLORS
	var colors_lst []string
	if _,ok := pPostInfoMap["colors_lst"]; !ok {
		colors_lst = []string{}
	} else if tagsLst, ok := pPostInfoMap["colors_lst"].([]string); !ok {
		//if pPostInfoMap is coming from mongodb, its of []interface{} type, so a "type conversion" is done)
		colors_lst = []string(tagsLst)
	} else {
		colors_lst = pPostInfoMap["colors_lst"].([]string)
	}

	//--------------------
	post := &GFpost{
		IDstr:               id_str,
		Tstr:                "post",
		DeletedBool:         false,
		ClientTypeStr:       pPostInfoMap["client_type_str"].(string),
		TitleStr:            post_title_str,
		DescriptionStr:      pPostInfoMap["description_str"].(string),
		CreationDatetimeStr: creation_datetime_str,
		PosterUserNameStr:   pPostInfoMap["poster_user_name_str"].(string),
		ThumbnailURLstr:     thumbnail_url_str,
		ImagesIDsLst:        images_ids_lst,
		PostElementsLst:     post_elements_lst,
		TagsLst:             tagsLst,
		NotesLst:            notesLst,
		ColorsLst:           colors_lst,
	}
	return post, nil
}

//------------------------------------------------	
// a post has to first be created, and only then can it be published

func publish(p_post_title_str string, pRuntimeSys *gf_core.RuntimeSys) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_post.publish()")




	


}

//------------------------------------------------
func createPostNotes(p_raw_notes_lst []map[string]interface{}, pRuntimeSys *gf_core.RuntimeSys) []*GFpostNote {

	notes_lst := []*GFpostNote{}
	for _, noteMap := range p_raw_notes_lst {
		
		snippet := &GFpostNote{
			UserIDstr: "anonymous",
			BodyStr:   noteMap["body_str"].(string),
		}
		notes_lst = append(notes_lst, snippet)
	}
	return notes_lst
}

//------------------------------------------------
func GetPostURL(pPostTitleStr string) string {



	return ""


}