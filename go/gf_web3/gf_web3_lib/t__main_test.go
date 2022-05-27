/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_web3_lib

import (
	"os"
	"time"
	"testing"
	"net/http"
	"context"
	// "github.com/gloflow/gloflow/go/gf_core"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_nft"
)

//---------------------------------------------------
func TestMain(m *testing.M) {

	
	runtime, _, err := gf_eth_core.TgetRuntime()
	if err != nil {
		panic(err)
	}

	// GF_WEB3_MONITOR_SERVICE
	testWeb3MonitorServicePortInt := 2000
	go func() {

		HTTPmux := http.NewServeMux()

		config := &gf_eth_core.GF_config{
			AlchemyAPIkeyStr: os.Getenv("GF_ALCHEMY_SERVICE_ACC__API_KEY"),
		}
		InitService(HTTPmux,
			config,
			runtime.RuntimeSys)
			
		gf_rpc_lib.ServerInitWithMux(testWeb3MonitorServicePortInt, HTTPmux)
	}()

	// GF_IDENTITY_SERVICE
	testIdentityServicePortInt := 2001
	go func() {

		gf_identity_lib.TestStartService(testIdentityServicePortInt,
			runtime.RuntimeSys)
	}()

	time.Sleep(2*time.Second) // let services startup

	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------
func TestAddressNFT(pTest *testing.T) {

	runtime, _, err := gf_eth_core.TgetRuntime()
	if err != nil {
		pTest.FailNow()
	}

	testWeb3MonitorServiceInt  := 2000
	testIdentityServicePortInt := 2001
	HTTPagent := gorequest.New()
	ctx       := context.Background()

	// CREATE_AND_LOGIN_NEW_USER
	gf_identity_lib.TestCreateAndLoginNewUser(pTest,
		HTTPagent,
		testIdentityServicePortInt,
		ctx,
		runtime.RuntimeSys)

	//--------------------
	// NFT_INDEX_ADDRESS
	testUserAddressEthStr := "0x4eDE0b31Fd116B8A00ADD6F449499Cd36b70AAE6"
	chainStr := "eth"
	gf_nft.TindexAddress(testUserAddressEthStr,
		chainStr,
		HTTPagent,
		testWeb3MonitorServiceInt,
		pTest)

	//--------------------
	
}