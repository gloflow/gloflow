package gf_domains_lib

import (
	"fmt"
	"time"
	"net/http"
	"text/template"
	"gf_core"
	"gf_rpc_lib"
)
//-------------------------------------------------
func Init_handlers(p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains_handlers.Init_handlers()")

	//---------------------
	//TEMPLATES
	template_path_str         := "./templates/gf_domains_browser.html"
	domains_browser__tmpl,err := template.New("gf_domains_browser.html").ParseFiles(template_path_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to parse a template",
			"template_create_error",
			&map[string]interface{}{"template_path_str":template_path_str,},
			err,"gf_images_lib",p_runtime_sys)
		return gf_err
	}
	//---------------------

	//---------------------
	//POSTS_ELEMENTS
	http.HandleFunc("/a/domains/browser",func(p_resp http.ResponseWriter,
											p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO","INCOMING HTTP REQUEST - /a/domains/browser ----------")

		if p_req.Method == "GET" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			//--------------------
			//response_format_str - "json"|"html"

			qs_map := p_req.URL.Query()
			fmt.Println(qs_map)

			/*//response_format_str - "j"(for json)|"h"(for html)
			response_format_str := gf_rpc_lib.Get_response_format(qs_map,
															p_log_fun)*/
			//--------------------
			//GET DOMAINS FROM DB
			domains_lst,gf_err := db__get_domains(p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/domains/browser",
									"rpc_handler failed getting domains",
									gf_err,p_resp,p_runtime_sys)
				return
			}
			//--------------------
			//RENDER TEMPLATE
			gf_err = domains_browser__render_template(domains_lst,
											domains_browser__tmpl,
											p_resp,
											p_runtime_sys)
			if gf_err != nil {
				gf_rpc_lib.Error__in_handler("/a/domains/browser",
									"failed to render domains_browser page",
									gf_err,p_resp,p_runtime_sys)
				return
			}

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/a/domains/browser",
										start_time__unix_f,
										end_time__unix_f,
										p_runtime_sys)
			}()
		}
	})
	//---------------------

	return nil
}