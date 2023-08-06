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
	// "fmt"
	"testing"
	"context"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_identity"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow/go/gf_web3/gf_address"
)

//---------------------------------------------------
func TestAddresses(pTest *testing.T) {


	runtime, _, err := gf_eth_core.TgetRuntime()
	if err != nil {
		pTest.FailNow()
	}

	testGFserviceInt           := 2000
	testIdentityServicePortInt := 2001
	HTTPagent := gorequest.New()
	ctx       := context.Background()

	// CREATE_AND_LOGIN_NEW_USER
	cookiesInRespLst := gf_identity.TestUserpassCreateAndLoginNewUser(pTest,
		HTTPagent,
		testIdentityServicePortInt,
		ctx,
		runtime.RuntimeSys)	


		
	
	testUserAddressEthStr := "0x4eDE0b31Fd116B8A00ADD6F449499Cd36b70AAE6"
	testAddressTypeStr := "observed"
	chainStr := "eth"
	
	// ADD_ADDRESS
	gf_address.TaddAddress(testUserAddressEthStr,
		testAddressTypeStr,
		chainStr,
		cookiesInRespLst,
		HTTPagent,
		testGFserviceInt,
		pTest)

	// GET_ALL_ADDRESSES
	gf_address.TgetAllAddresses(testAddressTypeStr,
		chainStr,
		cookiesInRespLst,
		HTTPagent,
		testGFserviceInt,
		pTest)

}