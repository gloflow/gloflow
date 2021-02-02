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

package gf_eth_monitor_lib

import (
	"fmt"
	
	// "context"
	// "strings"
	log "github.com/sirupsen/logrus"
	// "github.com/getsentry/sentry-go"
	"github.com/ethereum/go-ethereum/ethclient"
	// eth_types "github.com/ethereum/go-ethereum/core/types"
	// eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_core"
	// "github.com/davecgh/go-spew/spew"
)

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
		
			
		error_defs_map := gf_eth_monitor_core.Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to connect to Eth rpc-json API in gf_eth_monitor",
			"eth_rpc__dial",
			map[string]interface{}{"host": p_host_str,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
    }

	log.WithFields(log.Fields{"host": p_host_str, "port": p_geth_port_int}).Info("Connected to Ethereum node")
	

	return client, nil
}