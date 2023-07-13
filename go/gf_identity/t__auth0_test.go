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

package gf_identity

import (
	"fmt"
	"testing"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

/*
func TestAuth0Login(pTest *testing.T) {
	fmt.Println(" TEST__IDENTITY_AUTH0_LOGIN >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	serviceNameStr := "gf_identity_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	runtimeSys := Tinit(serviceNameStr, mongoHostStr)



	authSubsystemTypeStr := "auth0"
	portInt := 2000
	TestStartService(authSubsystemTypeStr, portInt, runtimeSys)





	HTTPagent   := gorequest.New()





	urlStr := "/v1/identity/auth0/login"
	_, bodyStr, errs := HTTPagent.Get(urlStr).
		End()
	if (len(errs) > 0) {
		runtimeSys.LogNewFun("DEBUG", "eth preflight HTTP failed", map[string]interface{}{"errs": errs,})
		pTest.FailNow()
	}
}
*/

//-------------------------------------------------

func TestAuth0(pTest *testing.T) {

	fmt.Println(" TEST__IDENTITY_AUTH0 >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	serviceNameStr := "gf_identity_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	runtimeSys := Tinit(serviceNameStr, mongoHostStr)



	auth0config := gf_auth0.LoadConfig(runtimeSys)

	auth0keyIDstr, auth0publicKey, gfErr := gf_auth0.GetJWTpublicKeyForTenant(auth0config, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	pubKeyPEMstr := gf_core.CryptoConvertPubKeyToPEM(auth0publicKey)

	runtimeSys.LogNewFun("INFO", "Auth0 keys...", map[string]interface{}{
		"auth0_key_id":  auth0keyIDstr,
		"auth0_pub_key": auth0publicKey,
		"auth0_pub_key_pem": pubKeyPEMstr,
	})


}

