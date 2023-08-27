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

package gf_eth_core

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
	"context"
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
	p_error_data_map    map[string]interface{},
	p_ctx               context.Context,
	pRuntimeSys       *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
	

	eth_http_port_int := 8545

	//-----------------------
	// HTTP_POST

	// IMPORTANT!! - some eth operations, like transactiont tracing, could take a significant amount of time.
	//               so timeout has to be sufficiently large.
	timeout_sec := time.Second * 60
	client      := &http.Client{Timeout: timeout_sec,}

	url_str  := fmt.Sprintf("http://%s:%d", p_eth_node_host_str, eth_http_port_int)
	req, err := http.NewRequestWithContext(p_ctx, "POST", url_str, strings.NewReader(p_input_json_str))
	if err != nil {

		error_data_map := map[string]interface{}{
			"url_str":        url_str,
			"input_json_str": p_input_json_str,
		}
		for k, v := range p_error_data_map {
			error_data_map[k] = v
		}
		gfErr := gf_core.ErrorCreate("failed to construct HTTP POST request using JSON input",
			"http_client_req_error",
			error_data_map,
			err, "gf_eth_core", pRuntimeSys)
		return nil, gfErr
	}
	req.Header.Set("Content-Type", "application/json")
	
	// EXECUTE
	resp, err := client.Do(req)
	if err != nil {
		error_data_map := map[string]interface{}{"url_str": url_str,}
		for k, v := range p_error_data_map {
			error_data_map[k] = v
		}
		gfErr := gf_core.ErrorCreate("failed to execute HTTP POST request to eth_rpc API",
			"http_client_req_error",
			error_data_map,
			err, "gf_eth_core", pRuntimeSys)
		return nil, gfErr
	}

	//-----------------------
	// JSON_DECODE
	body_bytes_lst, _ := ioutil.ReadAll(resp.Body)

	var output_map map[string]interface{}
	err = json.Unmarshal(body_bytes_lst, &output_map)
	if err != nil {
		error_data_map := map[string]interface{}{"url_str": url_str,}
		for k, v := range p_error_data_map {
			error_data_map[k] = v
		}
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from gf_rpc_client"), 
			"json_decode_error",
			error_data_map,
			err, "gf_eth_core", pRuntimeSys)
		return nil, gfErr
	}
	
	//-----------------------

	return output_map, nil
}

//-------------------------------------------------
// INIT

func Eth_rpc__init(p_host_str string,
	p_geth_port_int int,
	pRuntimeSys   *gf_core.RuntimeSys) (*ethclient.Client, *gf_core.GFerror) {

	

	url_str := fmt.Sprintf("http://%s:%d", p_host_str, p_geth_port_int)

	client, err := ethclient.Dial(url_str)
    if err != nil {
		log.Fatal(err)

		log.WithFields(log.Fields{
			"url_str":   url_str,
			"geth_host": p_host_str,
			"port":      p_geth_port_int,
			"err":       err}).Fatal("failed to connect json-rpc connect to Eth node")
		
			
		error_defs_map := ErrorGetDefs()
		gfErr := gf_core.ErrorCreateWithDefs("failed to connect to Eth rpc-json API in gf_eth_monitor",
			"eth_rpc__dial",
			map[string]interface{}{"host": p_host_str,},
			err, "gf_eth_core", error_defs_map, 1, pRuntimeSys)
		return nil, gfErr
    }

	log.WithFields(log.Fields{"host": p_host_str, "port": p_geth_port_int}).Info("Connected to Ethereum node")
	

	return client, nil
}