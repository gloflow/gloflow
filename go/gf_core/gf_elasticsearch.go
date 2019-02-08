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

package gf_core

import (
	"fmt"
	"context"
	"github.com/olivere/elastic"
)
//-------------------------------------------------
func Elastic__get_client(p_runtime_sys *Runtime_sys) (*elastic.Client,*Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_elasticsearch.Elastic__get_client()")

	es_host_str := "127.0.0.1:9200"
	p_runtime_sys.Log_fun("INFO","es_host_str - " + es_host_str)

	elasticsearch_client,err := elastic.NewClient(elastic.SetURL("http://"+es_host_str))
	if err != nil {
		gf_err := Error__create("failed to insert a user_track_start into mongodb",
			"elasticsearch_get_client",
			&map[string]interface{}{"es_host_str":es_host_str,},
			err,"gf_core",p_runtime_sys)
		return nil,gf_err	
	}


	//ping elasticsearch server
	ctx                     := context.Background()
	ping_url_str            := fmt.Sprintf("http://%s",es_host_str)
	es_info, resp_code, err := elasticsearch_client.Ping(ping_url_str).Do(ctx)
	if err != nil {
		gf_err := Error__create("failed to insert a user_track_start into mongodb",
			"elasticsearch_ping",
			&map[string]interface{}{"ping_url_str":ping_url_str,},
			err,"gf_core",p_runtime_sys)
		return nil,gf_err	
	}

	fmt.Printf("elasticsearch - resp_code/version %d/%s\n", resp_code, es_info.Version.Number)

	return elasticsearch_client,nil
}