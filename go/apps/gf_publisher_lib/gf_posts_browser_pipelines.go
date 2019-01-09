package gf_publisher_lib

import (
	"fmt"
	"text/template"
	"net/http"
	"github.com/globalsign/mgo"
)
//------------------------------------------------
func Get_posts_page(p_page_index_int int,
	p_page_elements_num_int int,
	p_mongodb_coll          *mgo.Collection,
	p_log_fun               func(string,string)) ([]map[string]interface{},error) {
	p_log_fun("FUN_ENTER","gf_posts_browser_pipelines.Get_posts_page()")

	cursor_start_position_int := p_page_index_int*p_page_elements_num_int
	page_lst,err              := DB__get_posts_page(cursor_start_position_int, p_page_elements_num_int, p_mongodb_coll, p_log_fun)
	if err != nil {
		return nil,err
	}

	serialized_page_lst := []map[string]interface{}{}
	for _,post := range page_lst {
		post_map := map[string]interface{}{
			"title_str":            post.Title_str,
			"images_number_str":    len(post.Images_ids_lst),
			"creation_datetime_str":post.Creation_datetime_str,
			"thumbnail_url_str":    post.Thumbnail_url_str,
			"tags_lst":             post.Tags_lst,
		}
		serialized_page_lst = append(serialized_page_lst,post_map)
	}

	return serialized_page_lst,nil
}
//------------------------------------------------
//get initial pages - the pages that are rendered in the initial HTML template. 
//                    subsequent pages are loaded as AJAX requests, via HTTP API. 

func Render_initial_pages(p_response_format_str string,
	p_initial_pages_num_int int, //6
	p_page_size_int         int, //5
	p_tmpl                  *template.Template,
	p_resp                  http.ResponseWriter,
	p_mongodb_coll          *mgo.Collection,
	p_log_fun               func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_posts_browser_pipelines.Render_initial_pages()")
	
	posts_pages_lst := [][]*Post{}

	for i:=0;i<p_initial_pages_num_int;i++ {

		start_position_int := i*p_page_size_int
		//int end_position_int   = start_position_int+p_page_size_int;

		p_log_fun("INFO",fmt.Sprintf(">>>>>>> start_position_int - %d - %d", start_position_int, p_page_size_int))

		//initial page might be larger then subsequent pages, that are requested 
		//dynamically by the front-end
		page_lst,err := DB__get_posts_page(start_position_int, p_page_size_int, p_mongodb_coll, p_log_fun)
		if err != nil {
			return err
		}

		posts_pages_lst = append(posts_pages_lst,page_lst)
	}
	
	err := posts_browser__render_template(posts_pages_lst, p_tmpl, p_page_size_int, p_resp, p_log_fun)
	if err != nil {
		return err
	}

	return nil
}