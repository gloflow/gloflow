package gf_publisher_lib

import (
	"fmt"
	"net/http"
	"encoding/json"
	"text/template"
	"github.com/globalsign/mgo"
)
//------------------------------------------------
//CREATE_POST
func Pipeline__create_post(p_post_info_map map[string]interface{},
					p_gf_images_service_host_port_str *string,
					p_mongo_coll                      *mgo.Collection,
					p_log_fun                         func(string,string)) (*Post,*string,error) {
	p_log_fun("FUN_ENTER","gf_post_pipelines.Pipeline__create_post()")

	//----------------------
	//VERIFY INPUT
	max_title_chars_int       := 100
	max_description_chars_int := 1000
	post_element_tag_max_int  := 20

	p_log_fun("INFO","p_post_info_map - "+fmt.Sprint(p_post_info_map))
	verified_post_info_map,err := verify_external_post_info(p_post_info_map,
												max_title_chars_int,
												max_description_chars_int,
												post_element_tag_max_int,
												p_mongo_coll,
												p_log_fun)
	if err != nil {
		return nil,nil,err
	}
	//----------------------
	//CREATE POST
	post,err := create_new_post(verified_post_info_map,
							p_log_fun)
	if err != nil {
		return nil,nil,err
	}

	p_log_fun("INFO","post - "+fmt.Sprint(post))
	//----------------------
	//PERSIST POST
	err = DB__create_post(post,
					p_mongo_coll,
					p_log_fun)
	if err != nil {
		return nil,nil,err
	}
	//----------------------
	//IMAGES
	//IMPORTANT - long-lasting image operation
	images_job_id_str,img_err := process_external_images(post,
										p_gf_images_service_host_port_str,
										p_mongo_coll,
										p_log_fun)
	if img_err != nil {
		return nil,nil,img_err
	}
	//----------------------

	return post,images_job_id_str,nil
}
//------------------------------------------------
func Pipeline__get_post(p_post_title_str *string,
				p_response_format_str *string,
				p_tmpl                *template.Template,
				p_resp                http.ResponseWriter,
				p_mongo_coll          *mgo.Collection,
				p_log_fun             func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_post_pipelines.Pipeline__get_post()")

	post,err := DB__get_post(p_post_title_str,
						p_mongo_coll,
						p_log_fun)
	if err != nil {
		return err
	}

	//------------------
	//HACK!!
	//some of the post_adt.tags_lst have tags that are empty strings (")
	//which showup as artifacts in the HTML since each tag gets a <div></div>
	//so here a post_adt is modified in place. 
	//this will over time correct/remove empty string tags, but the source cause of this
	//(on post tagging/creation) is still there, so find that and fix it.
	whole_tags_lst := []string{}
	for _,tag_str := range post.Tags_lst {
		if tag_str != "" {
			whole_tags_lst = append(whole_tags_lst,tag_str)
		}
	}
	post.Tags_lst = whole_tags_lst
	//------------------

	switch *p_response_format_str {
		//------------------
		//HTML RENDERING
		case "html":

			//SCALABILITY!!
			//ADD!! - cache this result in redis, and server it from there
			//        only re-generate the template every so often
			//        or figure out some quick way to check if something changed
			err := post__render_template(post,
									p_tmpl,
									p_resp,
									p_log_fun)
			if err != nil {
				return err
			}
		//------------------
		//JSON EXPORT
		
		case "json":

			post_lst,err := json.Marshal(post)
			if err != nil {
				return err
			}
			post_str := string(post_lst)

			p_resp.Write([]byte(post_str))
		//------------------
	}

	return nil
}