/*
GloFlow media management/publishing system
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

package gf_domains_lib

import (
	"fmt"
	"net/url"
	"github.com/globalsign/mgo/bson"
	"github.com/fatih/color"
	"github.com/gloflow/gloflow/go/gf_core"
)
//---------------------------------------------------
/*example post:
{
	"_id" : ObjectId("54bc39f02c42e565d20c432a"),
	"poster_user_name_str" : "Ivan",
	"creation_datetime" : "2012-09-24T04:15:02.315146",
	"colors_lst" : [ ],
	"post_elements_lst" : [
		{
			"height_str" : null,
			"colors_lst" : [ ],
			"description_str" : null,
			"meta_dict" : {
				
			},
			"post_index_3_tpl" : [
				0,
				0,
				0
			],
			"id_str" : "pub_pe:Apollo program - Lunar Module - jetison:(0, 0, 0)",
			"tags_lst" : [ ],
			"type_str" : "link",
			"extern_url_str" : "http://www.hq.nasa.gov/alsj/S66-10987.jpg",
			"width_str" : null
		},
		{
			"height_str" : null,
			"img_thumbnail_large_url_str" : "/images/d/1b8cab084a27a768cc87e662facc2ccb_thumb_large.jpeg",
			"colors_lst" : [ ],
			"img_thumbnail_small_url_str" : "/images/d/1b8cab084a27a768cc87e662facc2ccb_thumb_small.jpeg",
			"meta_dict" : {
				
			},
			"post_index_3_tpl" : [
				0,
				0,
				0
			],
			"description_str" : null,
			"img_thumbnail_medium_url_str" : "/images/d/1b8cab084a27a768cc87e662facc2ccb_thumb_medium.jpeg",
			"id_str" : "pub_pe:Apollo program - Lunar Module - jetison:(1, 0, 0)",
			"tags_lst" : [ ],
			"type_str" : "image",
			"extern_url_str" : "http://www.hq.nasa.gov/alsj/S66-10987.jpg",
			"width_str" : null
		}
	],
	"description_str" : "",
	"title_str" : "Apollo program - Lunar Module - jetison",
	"images_ids_lst" : [
		"1b8cab084a27a768cc87e662facc2ccb"
	],
	"id_str" : "Apollo program - Lunar Module - jetison:2012-09-24 04:15:02.315146",
	"tags_lst" : [
		"nasa",
		"nasa_apollo",
		"moon",
		"space"
	],
	"type_str" : "gchrome_ext",
	"main_image_url_str" : "",
	"main_image_thumbnail_medium_url_str" : null,
	"t" : "post"
}*/
//---------------------------------------------------
type Gf_domain_posts struct {
	Name_str  string `bson:"name_str"`
	Count_int int    `bson:"count_int"`
}
//---------------------------------------------------
func Get_domains_posts__mongo(p_runtime_sys *gf_core.Runtime_sys) ([]Gf_domain_posts, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains__posts.Get_domains_posts__mongo()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	p_runtime_sys.Log_fun("INFO",cyan("AGGREGATE POSTS DOMAINS ")+yellow(">>>>>>>>>>>>>>>"))

	pipe := p_runtime_sys.Mongodb_coll.Pipe([]bson.M{

		bson.M{"$match":bson.M{"t":"post",},},

		//{"\$limit":10}, //DEBUGGING
		//---------------------
		//IMPORTANT!! - for each document create a new one that only contains the "post_elements_lst"
		//              field, since this is where the url string is contained, in the post_element
		//              of type 'link'
		bson.M{"$project":bson.M{
				"_id":              false, //suppression of the "_id" field
				"post_elements_lst":"$post_elements_lst",
			},
		},
		//---------------------
		//ATTENTION!! - potentially creates a large number of docs. for large amounts of original docs
		//              with large post_elements_lst members, this may overflow memory.
		//              figure out how to horizontally partition this op. 
		// 
		//create as many new docs as there are post_elements in the current doc,
		//where the new docs contain all the fields of the original source doc
		bson.M{"$unwind":"$post_elements_lst",},
		//---------------------
		//"$type_str" - attribute of "post_elements_lst" elements
		bson.M{"$match":bson.M{"post_elements_lst.type_str":"link",},},
		bson.M{"$project":bson.M{
				"_id":                false, //suppression of the "_id" field
				"post_extern_url_str":"$post_elements_lst.extern_url_str",
			},
		},
	})
	
	type Posts_extern_urls_result struct {
		Post_extern_url_str string `bson:"post_extern_url_str"`
	}

	results_lst := []Posts_extern_urls_result{}
	err         := pipe.All(&results_lst)

	if err != nil {
		gf_err := gf_core.Error__create("failed to run an aggregation pipeline to get domains posts",
			"mongodb_aggregation_error",
			nil,err,"gf_domains_lib",p_runtime_sys)
		return nil,gf_err
	}
	//---------------
	//IMPORTANT!! - application level join. move this to Db with the apporpriate "domain_str" field

	parsed_domains_map := map[string]int{}
	for _,r := range results_lst {

		u,err := url.Parse(r.Post_extern_url_str)
		if err != nil {

			continue
		}
		domain_str := u.Host


		if _,ok := parsed_domains_map[domain_str]; ok {
			parsed_domains_map[domain_str] += 1
		} else {
			parsed_domains_map[domain_str] = 1
		}
	}
	//---------------
	domain_posts_lst := []Gf_domain_posts{}
	for domain_str,count_int := range parsed_domains_map {

		domain_posts := Gf_domain_posts{
			Name_str: domain_str,
			Count_int:count_int,
		}
		domain_posts_lst = append(domain_posts_lst, domain_posts)
	}

	p_runtime_sys.Log_fun("INFO",yellow(">>>>>>>> DOMAIN_POSTS FOUND - ")+cyan(fmt.Sprint(len(domain_posts_lst))))
	//---------------

	return domain_posts_lst, nil

	/*mongo console query - for testing
	db.posts.aggregate(
		{"$match":{"t":"post"}},
		{"$unwind":"$post_elements_lst"},
		{"$match":{"post_elements_lst.type_str":"link"}},
		{"$project":{"_id":0,"title_str":1,"url_domain_str":"$post_elements_lst.extern_url_domain_str"}},
		{"$group":{"_id":"$url_domain_str","titles_lst":{"$addToSet":"$title_str"},"count":{"$sum":1}}},
		{"$sort":{"count":-1}});*/
}