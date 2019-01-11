package main

import (
	"fmt"
	"text/template"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)
//--------------------------------------------------
func render_objects_with_tag(p_tag_str string,
	p_tmpl           *template.Template,
	p_page_index_int int,
	p_page_size_int  int,
	p_resp           http.ResponseWriter,
	p_runtime_sys    *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_tagger_view.render_objects_with_tag()");

	//-----------------------------
	//FIX!! - SCALABILITY!! - get tag info on "image" and "post" types is a very long
	//                        operation, and should be done in some more efficient way,
	//                        as in with a mongodb aggregation pipeline

	objects_infos_lst, gf_err := get_objects_with_tag(p_tag_str,
		"post", //p_object_type_str
		p_page_index_int,
		p_page_size_int,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	posts_with_tag_lst := []map[string]interface{}{}
	for _,p_object_info_map := range objects_infos_lst {

		//----------------
		var post_thumbnail_url_str string
		thumb_small_str := p_object_info_map["thumbnail_small_url_str"].(string)

		if thumb_small_str == "" {
			
			//FIX!! - use some user-configurable value that is configured at startup
			//IMPORTANT!! - some "thumbnail_small_url_str" are blank strings (""),
			error_img_url_str     := "http://gloflow.com/images/d/gf_landing_page_logo.png"
			post_thumbnail_url_str = error_img_url_str
		} else {
			post_thumbnail_url_str = thumb_small_str
		}
		//----------------
		post_info_map := map[string]interface{}{
			"post_title_str":        p_object_info_map["title_str"].(string),
			"post_tags_lst":         p_object_info_map["tags_lst"].([]string),
			"post_url_str":          p_object_info_map["url_str"].(string),
			"post_thumbnail_url_str":post_thumbnail_url_str,
		}

		posts_with_tag_lst = append(posts_with_tag_lst,post_info_map)
	}
	//-----------------------------


	type tmpl_data struct {
		Tag_str                string
		Posts_with_tag_num_int int
		Images_with_tag_int    int
		Posts_with_tag_lst     []map[string]interface{}
	}

	object_type_str                  := "post"
	posts_with_tag_count_int, gf_err := db__get_objects_with_tag_count(p_tag_str, object_type_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}
	
	err = p_tmpl.Execute(p_resp,
		tmpl_data{
			Tag_str:               p_tag_str,
			Posts_with_tag_num_int:posts_with_tag_count_int,
			Images_with_tag_int:   0, //FIX!! - image tagging is now implemented, and so counting images with tag occurance should be done ASAP. 
			Posts_with_tag_lst:    posts_with_tag_lst,
		})

	if err != nil {
		gf_err := gf_core.Error__create("failed to render the objects_with_tag template",
			"template_render_error",
			&map[string]interface{}{
				"tag_str":p_tag_str,
			},
			err, "gf_tagger", p_runtime_sys)
		return gf_err
	}

	return nil
}	