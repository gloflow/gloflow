package main

import (
	"time"
	"fmt"
	"strconv"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
	"text/template"
	"github.com/globalsign/mgo"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/apps/gf_publisher_lib"
)
//-------------------------------------------------
func init_handlers(p_gf_images_service_host_port_str *string,
	p_mongodb_coll *mgo.Collection,
	p_log_fun    func(string,string)) error {
	p_log_fun("FUN_ENTER","gf_publisher_service_handlers.init_handlers()")

	//---------------------
	//TEMPLATES
	post__tmpl,err := template.New("gf_post.html").ParseFiles("./templates/gf_post.html")
	if err != nil {
		return err
	}

	posts_browser__tmpl,err := template.New("gf_posts_browser.html").ParseFiles("./templates/gf_posts_browser.html")
	if err != nil {
		return err
	}

	/*posts_stats__tmpl,err := template.New("gf_posts_stats.html").ParseFiles("./templates/gf_posts_stats.html")
	if err != nil {
		return err
	}*/
	//---------------------
	//HIDDEN DASHBOARD

	http.HandleFunc("/posts/dash/18956180__42115/",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/dash ----------")


	})
	//---------------------
	//POSTS
	
	http.HandleFunc("/posts/create",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/create ----------")

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//------------
			//INPUT
			post_info_map    := map[string]interface{}{}
			body_bytes_lst,_ := ioutil.ReadAll(p_req.Body)
		    err              := json.Unmarshal(body_bytes_lst,&post_info_map)

			if err != nil {
				p_log_fun("ERROR",fmt.Sprint(err))
				gf_rpc_lib.Error__in_handler("/posts/create", err, "create_post pipeline received bad post_info_map input", p_resp, p_mongodb_coll, p_log_fun)
				return
			}
			//------------

			_,images_job_id_str,err := gf_publisher_lib.Pipeline__create_post(post_info_map,
				p_gf_images_service_host_port_str,
				p_mongodb_coll,
				p_log_fun)

			if err != nil {
				p_log_fun("ERROR",fmt.Sprint(err))
				gf_rpc_lib.Error__in_handler("/posts/create", err, "create_post pipeline failed", p_resp, p_mongodb_coll, p_log_fun)
				return 
			}

			gf_rpc_lib.Http_Respond(map[string]interface{}{"images_job_id_str":*images_job_id_str}, "OK", p_resp, p_log_fun)

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/create", start_time__unix_f, end_time__unix_f, p_mongodb_coll, p_log_fun)
			}()
		}
	})
	
	//---------------------
	//POST_STATUS
	
	http.HandleFunc("/posts/status",func(p_resp http.ResponseWriter, p_req *http.Request) {
			p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/status ----------")
		})
	//---------------------
	/*http.HandleFunc("/posts/create_with_updates",func(p_resp http.ResponseWriter,
													p_req *http.Request) {
			p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/create_with_updates ----------")
		})*/

	http.HandleFunc("/posts/update",func(p_resp http.ResponseWriter, p_req *http.Request) {
			p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/update ----------")
		})

	http.HandleFunc("/posts/delete",func(p_resp http.ResponseWriter, p_req *http.Request) {
			p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/delete ----------")
		})
	//---------------------
	//BROWSER

	http.HandleFunc("/posts/browser",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/browser ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------
			//response_format_str - "json"|"html"

			qs_map := p_req.URL.Query()

			//response_format_str - "j"(for json)|"h"(for html)
			response_format_str := gf_rpc_lib.Get_response_format(qs_map, p_log_fun)
			//--------------------

			err := gf_publisher_lib.Render_initial_pages(response_format_str,
				6, //p_initial_pages_num_int int
				5, //p_page_size_int
				posts_browser__tmpl,
				p_resp,
				p_mongodb_coll,
				p_log_fun)

			if err != nil {
				p_log_fun("ERROR",fmt.Sprint(err))
				gf_rpc_lib.Error__in_handler("/posts/browser", err, "failed to render posts_browser initial page", p_resp, p_mongodb_coll, p_log_fun)
				return
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/browser", start_time__unix_f, end_time__unix_f, p_mongodb_coll, p_log_fun)
			}()
		}
	})
	//---------------------
	//GET_BROWSER_PAGE (slice of posts data series)
	http.HandleFunc("/posts/browser_page",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/browser_page ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
			
			//--------------------
			//INPUT

			qs_map := p_req.URL.Query()

			page_index_int := 0 //default - "h" - HTML
			if a_lst,ok := qs_map["pg_index"]; ok {
				page_index_int,_ = strconv.Atoi(a_lst[0]) //user supplied value
				if err != nil {
					p_log_fun("ERROR",fmt.Sprint(err))
					gf_rpc_lib.Error__in_handler("/posts/browser_page", err, "pg_index (page_index) is not an integer", p_resp, p_mongodb_coll, p_log_fun)
					return
				}
			}

			page_size_int := 10 //default - "h" - HTML
			if a_lst,ok := qs_map["pg_size"]; ok {
				page_size_int,err = strconv.Atoi(a_lst[0]) //user supplied value
				if err != nil {
					p_log_fun("ERROR",fmt.Sprint(err))
					gf_rpc_lib.Error__in_handler("/posts/browser_page", err, "pg_size (page_size) is not an integer", p_resp, p_mongodb_coll, p_log_fun)
					return
				}
			}
			//--------------------
			
			serialized_pages_lst,err := gf_publisher_lib.Get_posts_page(page_index_int, page_size_int, p_mongodb_coll, p_log_fun)
			if err != nil {
				p_log_fun("ERROR",fmt.Sprint(err))
				gf_rpc_lib.Error__in_handler("/posts/browser_page", err, "failed to get posts page", p_resp, p_mongodb_coll, p_log_fun)
				return
			}

			//------------
			//JSON RESPONSE

			r_lst,_ := json.Marshal(serialized_pages_lst)
			r_str   := string(r_lst)
			fmt.Fprintf(p_resp,r_str)
			//------------

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/browser_page", start_time__unix_f, end_time__unix_f, p_mongodb_coll, p_log_fun)
			}()
		}
	})
	//---------------------
	//GET POST
	http.HandleFunc("/posts/",func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_log_fun("INFO","INCOMING HTTP REQUEST - /posts/ ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------
			//response_format_str - "j"(for json)|"h"(for html)

			qs_map := p_req.URL.Query()

			//response_format_str - "j"(for json)|"h"(for html)
			response_format_str := gf_rpc_lib.Get_response_format(qs_map,
															p_log_fun)
			//--------------------
			//POST_TITLE

			url_str          := p_req.URL.Path
			url_elements_lst := strings.Split(url_str,"/")

			//IMPORTANT!! - "!=3" - because /a/b splits into {"","a","b",}
			if len(url_elements_lst) != 3 {
				p_log_fun("ERROR",fmt.Sprint(err))
				gf_rpc_lib.Error__in_handler("/posts/", err, "get_post url is not of proper format - "+url_str, p_resp, p_mongodb_coll, p_log_fun)
				return
			}

			raw_post_title_str := url_elements_lst[2]

			//IMPORTANT!! - replaceAll() - is used here because at the time of testing all titles were still
			//                             with their spaces (" ") encoded as "+". So for the title to be correct,
			//                             for lookups against the internal DB, this is decoded.
			//decodeComponent() - this decodes the percentage encoded symbols. it does not remove
			//                    "+" encoded spaces (" "), and the need for replaceAll()
			post_title_encoded_str := strings.Replace(raw_post_title_str,"+"," ",-1)
			post_title_str,err     := url.QueryUnescape(post_title_encoded_str)
			if err != nil {
				p_log_fun("ERROR",fmt.Sprint(err))
				gf_rpc_lib.Error__in_handler("/posts/", err, "post title cant be query_unescaped - "+post_title_encoded_str, p_resp, p_mongodb_coll, p_log_fun)
				return
			}
			p_log_fun("INFO","post_title_str - "+post_title_str)
			//--------------------

			err = gf_publisher_lib.Pipeline__get_post(&post_title_str,
				&response_format_str,
				post__tmpl,
				p_resp,
				p_mongodb_coll,
				p_log_fun)

			if err != nil {
				p_log_fun("ERROR",fmt.Sprint(err))
				gf_rpc_lib.Error__in_handler("/posts/", err, "get_post pipeline failed", p_resp, p_mongodb_coll, p_log_fun)
				return
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/posts/", start_time__unix_f, end_time__unix_f, p_mongodb_coll, p_log_fun)
			}()
		}
	})
	//---------------------
	//POSTS_ELEMENTS
	http.HandleFunc("/posts_elements/create",func(p_resp http.ResponseWriter, p_req *http.Request) {


		})
	//---------------------

	return nil
}