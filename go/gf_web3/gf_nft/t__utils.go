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

package gf_nft

import (
	"fmt"
	"testing"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	// "github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func TindexAddress(pAddressStr string,
	pChainStr    string,
	pHTTPagent   *gorequest.SuperAgent,
	pTestPortInt int,
	pTest        *testing.T) {

	fmt.Println("====================================")
	fmt.Println("test NFT INDEX ADDRESS")

	urlStr  := fmt.Sprintf("http://localhost:%d/v1/web3/nft/index_address", pTestPortInt)
	fmt.Println("URL", urlStr)
	
	dataMap := map[string]string{
		"address_str": pAddressStr,
		"chain_str":   pChainStr,
	}
	dataBytesLst, _ := json.Marshal(dataMap)

	_, bodyStr, errs := pHTTPagent.Post(urlStr).
		Send(string(dataBytesLst)).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		pTest.FailNow()
	}
	
	bodyMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
		fmt.Println(err)
		pTest.FailNow()
	}

	spew.Dump(bodyMap)

	assert.True(pTest, bodyMap["status"].(string) != "ERROR", "nft address indexing http request failed")



}