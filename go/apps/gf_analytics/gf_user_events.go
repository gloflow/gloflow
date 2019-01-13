package main



type Gf_user_event struct {
	Id                   bson.ObjectId          `json:"-"                    bson:"_id,omitempty"`
	Id_str               string                 `json:"id_str"               bson:"id_str"` 
	T_str                string                 `json:"-"                    bson:"t"` //"usr_event"
	Creation_unix_time_f float64                `json:"creation_unix_time_f" bson:"creation_unix_time_f"`
	Event_data_map       map[string]interface{} `json:"event_data_map"       bson:"event_data_map"`
	Session_id_str       string                 `json:"session_id_str"       bson:"session_id_str"`
	User_ip_str          string                 `json:"user_ip_str"          bson:"user_ip_str"`
	User_agent_str       string                 `json:"user_agent_str"       bson:"user_agent_str"`
	Browser_name_str     string                 `json:"browser_name_str"     bson:"browser_name_str"`
	Browser_ver_str      string                 `json:"browser_ver_str"      bson:"browser_ver_str"`
	Os_name_str          string                 `json:"os_name_str"          bson:"os_name_str"`
	Os_ver_str           string                 `json:"os_ver_str"           bson:"os_ver_str"`
	Cookies_str          string                 `json:"cookies_str"          bson:"cookies_str"`
	time__unix_f         float64                `json:"time__unix_f"         bson:"time__unix_f"`
}
//-------------------------------------------------
func user_event__create(p_user_event *Gf_user_event,
	p_runtime_sys *gf_core.Runtime_sys) {


		









}