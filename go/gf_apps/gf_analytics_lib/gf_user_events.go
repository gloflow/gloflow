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

package gf_analytics_lib

import (
	"fmt"
	"time"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"crypto/md5"
	"encoding/hex"
	"context"
	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type Gf_user_event_input struct {
	Type_str string                 `json:"type_str"`
	Data_map map[string]interface{} `json:"d"`
}

type Gf_user_event_req_ctx struct {
	User_ip_str          string `json:"user_ip_str"      bson:"user_ip_str"`
	User_agent_str       string `json:"user_agent_str"   bson:"user_agent_str"`
	Browser_name_str     string `json:"browser_name_str" bson:"browser_name_str"`
	Browser_ver_str      string `json:"browser_ver_str"  bson:"browser_ver_str"`
	Os_name_str          string `json:"os_name_str"      bson:"os_name_str"`
	Os_ver_str           string `json:"os_ver_str"       bson:"os_ver_str"`
	Cookies_str          string `json:"cookies_str"      bson:"cookies_str"`
}

type Gf_user_event struct {
	Id                   primitive.ObjectID     `json:"-"                    bson:"_id,omitempty"`
	Id_str               string                 `json:"id_str"               bson:"id_str"`
	T_str                string                 `json:"-"                    bson:"t"` //"usr_event"
	Creation_unix_time_f float64                `json:"creation_unix_time_f" bson:"creation_unix_time_f"`
	Event_data_map       map[string]interface{} `json:"event_data_map"       bson:"event_data_map"`
	Session_id_str       string                 `json:"session_id_str"       bson:"session_id_str"`
	Req_ctx              Gf_user_event_req_ctx  `json:"req_ctx"              bson:"req_ctx"`
	time__unix_f         float64                `json:"time__unix_f"         bson:"time__unix_f"`
}

//-------------------------------------------------
func user_event__parse_input(p_req *http.Request,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_user_event_input, string, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_user_events.user_event__parse_input()")

	//--------------------
	input             := Gf_user_event_input{}
	body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
	err               := json.Unmarshal(body_bytes_lst, &input)
	
	//--------------------
	session_id_str := session__get_id_cookie(p_req, p_resp, p_runtime_sys)

	//--------------------

	if err != nil {
		gf_err := gf_core.Error__create("failed to parse json http input for user_event",
			"json_unmarshal_error",
			nil, err, "gf_analytics", p_runtime_sys)
		return nil, "", gf_err
	}
	return &input, session_id_str, nil
}

//-------------------------------------------------
func user_event__create(p_input *Gf_user_event_input,
	p_session_id_str string,
	p_gf_req_ctx     *Gf_user_event_req_ctx,
	p_runtime_sys    *gf_core.Runtime_sys) *gf_core.Gf_error {
	
	creation_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
	//--------------------
	// ID
	h := md5.New()
	h.Write([]byte(fmt.Sprint(creation_time__unix_f)))
	h.Write([]byte("user_event"))
	sum        := h.Sum(nil)
	id_hex_str := hex.EncodeToString(sum)
	
	//--------------------
	
	gf_user_event := &Gf_user_event{
		Id_str:               id_hex_str,
		T_str:                "usr_event",
		Creation_unix_time_f: creation_time__unix_f,
		Event_data_map:       p_input.Data_map,
		Session_id_str:       p_session_id_str,
		Req_ctx:              *p_gf_req_ctx,
	}

	ctx           := context.Background()
	coll_name_str := p_runtime_sys.Mongo_coll.Name()
	gf_err        := gf_core.Mongo__insert(gf_user_event,
		coll_name_str,
		map[string]interface{}{
			"caller_err_msg_str": "failed to insert a user_event into the DB",
		},
		ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	/*err := p_runtime_sys.Mongodb_coll.Insert(gf_user_event)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a user_event in mongodb",
			"mongodb_insert_error",
			map[string]interface{}{},
			err, "gf_analytics", p_runtime_sys)
		return gf_err
	}*/
		
	return nil
}

//-------------------------------------------------
func session__get_id_cookie(p_req *http.Request,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.Runtime_sys) string {

	cookie, _ := p_req.Cookie("gf")  
	if cookie == nil {
		session_id_str := session__create_id_cookie(p_req, p_resp, p_runtime_sys)
		return session_id_str
	} else {
		session_id_str := cookie.Value
		return session_id_str
	}
}

//-------------------------------------------------
func session__create_id_cookie(p_req *http.Request,
	p_resp        http.ResponseWriter,
	p_runtime_sys *gf_core.Runtime_sys) string {

	current_time__unix_f := float64(time.Now().UnixNano())/1000000000.0
	ip_str               := p_req.RemoteAddr
	session_id_str       := fmt.Sprintf("%f_%s", current_time__unix_f, ip_str)

	p_runtime_sys.Log_fun("INFO", "session_id_str - "+session_id_str)

	new_cookie := http.Cookie{Name:"gf", Value:session_id_str}
	http.SetCookie(p_resp, &new_cookie)

	return session_id_str
}