package main

import (
	"time"
	"strings"
	"net/http"
	"github.com/ianoshen/uaparser"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func init_handlers(p_runtime_sys *gf_core.Runtime_sys) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_analytics_handlers.init_handlers()")

	//--------------
	//USER_EVENT
	http.HandleFunc("/a/ue", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST --- /a/ue")

		if p_req.Method == "OPTIONS" {
			p_resp.Header().Set("Access-Control-Allow-Origin","*")
			p_resp.Header().Set("Access-Control-Allow-Origin","Origin, X-Requested-With, Content-Type, Accept")
		}

		if p_req.Method == "POST" {
			start_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

			ip_str       := p_req.RemoteAddr
			clean_ip_str := strings.Split(ip_str,":")[0]

			cookies_lst := p_req.Cookies()
			cookies_str := gf_core.HTTP__serialize_cookies(cookies_lst,p_runtime_sys)
			//-----------------
			//BROWSER INFORMATION
			user_agent_str := p_req.UserAgent()
			user_agent     := uaparser.Parse(user_agent_str)

			var browser_name_str string
			var browser_ver_str  string
			if user_agent.Browser != nil {
				browser_name_str = user_agent.Browser.Name
				browser_ver_str  = user_agent.Browser.Version
			}

			os_name_str    := user_agent.OS.Name
			os_version_str := user_agent.OS.Version
			//-----------------
			//INPUT
			input, session_id_str, gf_err := user_event__parse_input(p_req, p_resp, p_runtime_sys)
			if gf_err != nil {
				//IMPORTANT!! - this is a special case handler, we dont want it to return any standard JSON responses,
				//              this handler should be fire-and-forget from the users/clients perspective.
				return
			}
			//-----------------
						
			gf_req_ctx := &Gf_user_event_req_ctx {
				User_ip_str:      clean_ip_str,
				User_agent_str:   user_agent_str,
				Browser_name_str: browser_name_str,
				Browser_ver_str:  browser_ver_str,
				Os_name_str:      os_name_str,
				Os_ver_str:       os_version_str,
				Cookies_str:      cookies_str,
			}

			gf_err = user_event__create(input, session_id_str, gf_req_ctx, p_runtime_sys)
			if gf_err != nil {
				//IMPORTANT!! - this is a special case handler, we dont want it to return any standard JSON responses,
				//              this handler should be fire-and-forget from the users/clients perspective.
				return
			}
			//-----------------

			end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
		
			go func() {
				gf_rpc_lib.Store_rpc_handler_run("/a/ue", start_time__unix_f, end_time__unix_f, p_runtime_sys)
			}()
		}
	})
	//--------------
}