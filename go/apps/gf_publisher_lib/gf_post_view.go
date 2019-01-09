package gf_publisher_lib

import (
	"errors"
	"net/http"
	"text/template"
)
//--------------------------------------------------
func post__render_template(p_post *Post,
				p_tmpl    *template.Template,
				p_resp    http.ResponseWriter,
				p_log_fun func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_post_view.post__render_template()")
	
	template_post_elements_lst,err := package_post_elements_infos(p_post,
															p_log_fun)
	if err != nil {
		return err
	}

	image_post_elements_og_info_lst,err := get_image_post_elements_FBOpenGraph_info(p_post,
																				p_log_fun)
	if err != nil {
		return err
	}

	post_tags_lst := []string{}
	for _,tag_str := range p_post.Tags_lst {
		post_tags_lst = append(post_tags_lst,tag_str)
	}

	type tmpl_data struct {
		Post_title_str                  string
		Post_tags_lst                   []string
		Post_description_str            string
		Post_poster_user_name_str       string
		Post_thumbnail_url_str          string
		Post_elements_lst               []map[string]interface{}
		Image_post_elements_og_info_lst []map[string]string
	}
	
	/*template_info_map := map[string]interface{}{
		"post_title_str"                       :p_post.title_str,
		"post_tags_lst"                        :post_tags_lst,
		"post_description_str"                 :p_post.description_str,
		"post_poster_user_name_str"   template_str         :p_post.poster_user_name_str,
		"post_elements_lst"                    :template_post_elements_lst,
		"img_thumbnail_medium_absolute_url_str":post_thumbnail_url_str,
		"image_post_elements_og_info_lst"      :image_post_elements_og_info_lst,
	}

	final String template_str = p_template.renderString(template_info_map)
	return template_str;*/

	err = p_tmpl.Execute(p_resp,tmpl_data{
		Post_title_str                 :p_post.Title_str,
		Post_tags_lst                  :post_tags_lst,
		Post_description_str           :p_post.Description_str,
		Post_poster_user_name_str      :p_post.Poster_user_name_str,
		Post_thumbnail_url_str         :p_post.Thumbnail_url_str,
		Post_elements_lst              :template_post_elements_lst,
		Image_post_elements_og_info_lst:image_post_elements_og_info_lst,
	})

	if err != nil {
		return err
	}

	return nil
}
//--------------------------------------------------
func package_post_elements_infos(p_post *Post,
							p_log_fun func(string,string)) ([]map[string]interface{},error) {
	p_log_fun("FUN_ENTER","gf_post_view.package_post_elements_infos()")

	template_post_elements_lst := []map[string]interface{}{}

	for _,post_element := range p_post.Post_elements_lst {

		p_log_fun("INFO","post_element.Type_str - "+post_element.Type_str)

		if !(post_element.Type_str == "link" ||
			post_element.Type_str == "image" ||
			post_element.Type_str == "video" ||
			post_element.Type_str == "text") {
			return nil,errors.New("post_element type is not 'link'|'image'|'video'|'text' - "+post_element.Type_str)
		}

		post_element_tags_lst := []string{}
		for _,tag_str := range post_element.Tags_lst {
			post_element_tags_lst = append(post_element_tags_lst,tag_str)
		}

		switch post_element.Type_str {
			case "link":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool":false, //for mustache template conditionals
					"post_element_type__image_bool":false,
					"post_element_type__link_bool" :true,
					"post_element_description_str" :post_element.Description_str,
					"post_element_extern_url_str"  :post_element.Extern_url_str,
				}
				template_post_elements_lst = append(template_post_elements_lst,post_element_map)
				continue
			case "image":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool"            :false, //for mustache template conditionals
					"post_element_type__image_bool"            :true,
					"post_element_type__link_bool"             :false,
					"post_element_img_thumbnail_medium_url_str":post_element.Img_thumbnail_medium_url_str,
					"post_element_img_thumbnail_large_url_str" :post_element.Img_thumbnail_large_url_str,
					"tags_lst"                                 :post_element_tags_lst,
				}
				template_post_elements_lst = append(template_post_elements_lst,post_element_map)
				continue
			case "video":
				post_element_map := map[string]interface{}{
					"post_element_type__video_bool":true, //for mustache template conditionals
					"post_element_type__image_bool":false,
					"post_element_type__link_bool" :false,
					"post_element_extern_url_str"  :post_element.Extern_url_str,
					"tags_lst"                     :post_element_tags_lst,
				}
				template_post_elements_lst = append(template_post_elements_lst,post_element_map)
				continue
		}
	}

	return template_post_elements_lst,nil
}
//--------------------------------------------------
func get_image_post_elements_FBOpenGraph_info(p_post *Post,
								p_log_fun func(string,string)) ([]map[string]string,error) {
	p_log_fun("FUN_ENTER","gf_post_view.get_image_post_elements_FBOpenGraph_info()")

	image_post_elements_lst,err := get_post_elements_of_type(p_post,
													"image",
													p_log_fun)
	if err != nil {
		return nil,err
	}

	var top_image_post_elements_lst []*Post_element
	if len(image_post_elements_lst) > 5 {

		//getRange() - returns an Iterable<String>
		top_image_post_elements_lst = image_post_elements_lst[:5] //new List.from(image_post_elements_lst.getRange(0,5))
	} else { 
		top_image_post_elements_lst = image_post_elements_lst
	}

	//---------------------
	og_info_lst := []map[string]string{}
	for _,post_element := range top_image_post_elements_lst {
		d := map[string]string{
			"img_thumbnail_medium_absolute_url_str":post_element.Img_thumbnail_medium_url_str,
		}
		og_info_lst = append(og_info_lst,d)
	}
	//---------------------

	return og_info_lst,nil
}