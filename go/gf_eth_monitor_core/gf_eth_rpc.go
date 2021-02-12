/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_eth_monitor_core

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/getsentry/sentry-go"
	// eth_types "github.com/ethereum/go-ethereum/core/types"
	// eth_common "github.com/ethereum/go-ethereum/common"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Eth_rpc__call(p_input_json_str string,
	p_eth_node_host_str string,
	p_runtime_sys       *gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error) {
	

	eth_http_port_int := 8545

	//-----------------------
	// HTTP_POST

	timeout_sec := time.Second * 10
	client      := &http.Client{Timeout: timeout_sec,}

	url_str  := fmt.Sprintf("http://%s:%s", p_eth_node_host_str, eth_http_port_int)
	req, err := http.NewRequest("POST", url_str, strings.NewReader(p_input_json_str))
	if err != nil {
		gf_err := gf_core.Error__create("failed to construct HTTP POST request using JSON input",
			"http_client_req_error",
			map[string]interface{}{
				"url_str":        url_str,
				"input_json_str": p_input_json_str,
			},
			err, "gf_eth_monitor_core", p_runtime_sys)
		return nil, gf_err
	}
	req.Header.Set("Content-Type", "application/json")
	

	resp, err := client.Do(req)
	if err != nil {
		gf_err := gf_core.Error__create("failed to execute HTTP POST request to eth_rpc API",
			"http_client_req_error",
			map[string]interface{}{"url_str": url_str,},
			err, "gf_eth_monitor_core", p_runtime_sys)
		return nil, gf_err
	}

	//-----------------------
	// JSON_DECODE
	body_bytes_lst, _ := ioutil.ReadAll(resp.Body)

	var output_map map[string]interface{}
	err = json.Unmarshal(body_bytes_lst, &output_map)
	if err != nil {
		gf_err := gf_core.Error__create(fmt.Sprintf("failed to parse json response from gf_rpc_client"), 
			"json_decode_error",
			map[string]interface{}{"url_str": url_str,},
			err, "gf_eth_monitor_core", p_runtime_sys)
		return nil, gf_err
	}
	
	//-----------------------

	return output_map, nil
}

//-------------------------------------------------
// INIT
func Eth_rpc__init(p_host_str string,
	p_geth_port_int int,
	p_runtime_sys   *gf_core.Runtime_sys) (*ethclient.Client, *gf_core.Gf_error) {

	

	url_str := fmt.Sprintf("http://%s:%d", p_host_str, p_geth_port_int)

	client, err := ethclient.Dial(url_str)
    if err != nil {
		log.Fatal(err)

		log.WithFields(log.Fields{
			"url_str":   url_str,
			"geth_host": p_host_str,
			"port":      p_geth_port_int,
			"err":       err}).Fatal("failed to connect json-rpc connect to Eth node")
		
			
		error_defs_map := Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to connect to Eth rpc-json API in gf_eth_monitor",
			"eth_rpc__dial",
			map[string]interface{}{"host": p_host_str,},
			err, "gf_eth_monitor_core", error_defs_map, p_runtime_sys)
		return nil, gf_err
    }

	log.WithFields(log.Fields{"host": p_host_str, "port": p_geth_port_int}).Info("Connected to Ethereum node")
	

	return client, nil
}