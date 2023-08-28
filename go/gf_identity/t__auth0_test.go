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
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func TestAuth0(pTest *testing.T) {

	fmt.Println(" TEST__IDENTITY_AUTH0 >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	serviceNameStr := "gf_identity_test"
	mongoHostStr := cliArgsMap["mongodb_host_str"].(string) // "127.0.0.1"
	sqlHostStr   := cliArgsMap["sql_host_str"].(string)
	runtimeSys := Tinit(serviceNameStr, mongoHostStr, sqlHostStr, logNewFun, logFun)
	ctx := context.Background()

	gfErr := gf_identity_core.DBsqlCreateTables(ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	auth0config := gf_auth0.LoadConfig(runtimeSys)

	auth0keyIDstr, auth0publicKey, gfErr := gf_auth0.JWTgetPublicKeyForTenant(auth0config, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	pubKeyPEMstr := gf_core.CryptoConvertPubKeyToPEM(auth0publicKey)

	runtimeSys.LogNewFun("INFO", "Auth0 keys...", map[string]interface{}{
		"auth0_key_id":  auth0keyIDstr,
		"auth0_pub_key": auth0publicKey,
		"auth0_pub_key_pem": pubKeyPEMstr,
	})




	
	//----------------------
	// LOGIN
	sessionID, gfErr := gf_identity_core.Auth0loginPipeline(ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	runtimeSys.LogNewFun("INFO", "Auth0 login pipeline complete...", map[string]interface{}{
		"session_id": sessionID,
	})

	//----------------------
	// GET_SESSION
	session, gfErr := gf_identity_core.DBsqlAuth0getSession(sessionID, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	spew.Dump(session)

	//----------------------

	userID, userNameStr := TestCreateUserInDB(pTest, ctx, runtimeSys)

	resolvedUserNameStr, gfErr := gf_identity_core.DBsqlGetUserNameByID(userID, ctx, runtimeSys)
	runtimeSys.LogNewFun("INFO", "user_name resolving from user_id succeeded...", map[string]interface{}{
		"user_name": resolvedUserNameStr,
	})

	//----------------------
	updateOp := &gf_identity_core.GFloginAttemptUpdateOp{
		UserID:      &userID,
		UserNameStr: &userNameStr,
	}
	gfErr = gf_identity_core.DBsqlLoginAttemptUpdateBySessionID(sessionID, updateOp, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	//----------------------
	// LOGOUT
	gfErr = gf_identity_core.Auth0logoutPipeline(sessionID, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	//----------------------
	// GET_SESSION
	session, gfErr = gf_identity_core.DBsqlAuth0getSession(sessionID, ctx, runtimeSys)
	if gfErr != nil {
		pTest.Fail()
	}

	spew.Dump(session)

	assert.True(pTest, session.DeletedBool,
		"auth0 session should be marked as deleted")

	//----------------------
}