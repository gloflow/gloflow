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

package gf_identity_core

import (
	"fmt"
	"context"
	"crypto/rsa"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_extern_services/gf_auth0"
)

//---------------------------------------------------

type GFjwtSecretKeyVal     string
type GFjwtPublicKeyPEMval  string
type GFjwtPrivateKeyPEMval string

type GFjwtSecret struct {
	Vstr                 string                `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID    `bson:"_id,omitempty"`
	IDstr                gf_core.GF_ID         `bson:"id_str"`
	DeletedBool          bool                  `bson:"deleted_bool"`
	CreationUNIXtimeF    float64               `bson:"creation_unix_time_f"`

	PublicKeyPEMstr      GFjwtPublicKeyPEMval  `bson:"public_key_pem_str"`
	PrivateKeyPEMstr     GFjwtPrivateKeyPEMval `bson:"private_key_pem_str"`
}

type GFjwtValidationKeyReq struct {
	authSubsystemTypeStr string
	responseCh           chan *rsa.PublicKey
}

type GFjwtSigningKeyReq struct {
	authSubsystemTypeStr string
	responseCh           chan *rsa.PrivateKey
}

type GFkeyServerInfo struct {
	GetJWTvalidationKeyCh chan GFjwtValidationKeyReq
	GetJWTsigningKeyCh    chan GFjwtSigningKeyReq
}

//---------------------------------------------------
// CLIENT
//---------------------------------------------------

func KSclientJWTgetValidationKey(pAuthSubsystemTypeStr string,
	pKeyServerInfo *GFkeyServerInfo,
	pRuntimeSys    *gf_core.RuntimeSys) (*rsa.PublicKey, *gf_core.GFerror) {

	responseCh := make(chan *rsa.PublicKey)
	req := GFjwtValidationKeyReq{
		authSubsystemTypeStr: pAuthSubsystemTypeStr,
		responseCh:           responseCh,
	}
	pKeyServerInfo.GetJWTvalidationKeyCh <- req

	publicKey := <- responseCh
	return publicKey, nil
}

//---------------------------------------------------

func ksClientJWTgetSigningKey(pAuthSubsystemTypeStr string,
	pKeyServerInfo *GFkeyServerInfo,
	pRuntimeSys    *gf_core.RuntimeSys) (*rsa.PrivateKey, *gf_core.GFerror) {

	responseCh := make(chan *rsa.PrivateKey)
	req := GFjwtSigningKeyReq{
		authSubsystemTypeStr: pAuthSubsystemTypeStr,
		responseCh:           responseCh,
	}
	pKeyServerInfo.GetJWTsigningKeyCh <- req

	privateKey := <- responseCh
	return privateKey, nil
}

//---------------------------------------------------
// SERVER
//---------------------------------------------------

// initialize a goroutine that servers requests from other goroutines
// for public/private keypairs
func KSinit(pAuth0initBool bool,
	pRuntimeSys *gf_core.RuntimeSys) (*GFkeyServerInfo, *gf_core.GFerror) {

	pRuntimeSys.LogNewFun("INFO", "initializing gf_identity keys server...", nil)

	//------------------------
	// JWT_SIGNING_KEY - generate it if the user is not using a secret store, where they
	//                   placed it independently.
	//                   this key is always fetched/generated, regardless if Auth0 is
	//                   activated or not, since even with Auth0 activated the Ethereum
	//                   auth method is available.     
	ctx := context.Background()
	publicKey, privateKey, gfErr := ksJWTgetKeysPipeline(ctx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	// AUTH0
	var auth0publicKey *rsa.PublicKey
	if pAuth0initBool {

		auth0config := gf_auth0.LoadConfig(pRuntimeSys)

		_, auth0publicKey, gfErr = gf_auth0.GetJWTpublicKeyForTenant(auth0config, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}

	//------------------------


	getJWTvalidationKeyCh := make(chan GFjwtValidationKeyReq, 100)
	getJWTsigningKeyCh := make(chan GFjwtSigningKeyReq, 10)

	go func() {
		for {
			select {
			
			// VALIDATION
			case req := <- getJWTvalidationKeyCh:

				pRuntimeSys.LogNewFun("DEBUG", "key_server - received request for JWT validation key",
					map[string]interface{}{"auth_subsystem_type_str": req.authSubsystemTypeStr,})

				switch req.authSubsystemTypeStr {
				case GF_AUTH_SUBSYSTEM_TYPE__USERPASS:
					req.responseCh <- publicKey

				case GF_AUTH_SUBSYSTEM_TYPE__ETH:
					req.responseCh <- publicKey

				case GF_AUTH_SUBSYSTEM_TYPE__AUTH0:
					req.responseCh <- auth0publicKey
				} 

			// SIGNING
			case req := <- getJWTsigningKeyCh:

				pRuntimeSys.LogNewFun("DEBUG", "key_server - received request for JWT signing key",
						map[string]interface{}{"auth_subsystem_type_str": req.authSubsystemTypeStr,})

				// signing using keys managed by key_server are only done for "userpass" and "eth" auth methods
				if req.authSubsystemTypeStr == GF_AUTH_SUBSYSTEM_TYPE__USERPASS ||
					req.authSubsystemTypeStr == GF_AUTH_SUBSYSTEM_TYPE__ETH {

					req.responseCh <- privateKey
				} else {
					// unsupported auth_subsystem_type
					req.responseCh <- nil
				}
			}
		}
	}()

	info := &GFkeyServerInfo{
		GetJWTvalidationKeyCh: getJWTvalidationKeyCh,
		GetJWTsigningKeyCh:    getJWTsigningKeyCh,
	}
	return info, nil
}

//---------------------------------------------------

// used only by users that self-host and dont use a dedicated secret store.
// instead they store all data in the DB for max simplicity of hosting.
func ksJWTgetKeysPipeline(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*rsa.PublicKey, *rsa.PrivateKey, *gf_core.GFerror) {

	// if a secrets "get" callback is specified we're making an assumption
	// that the user has setup some sort of secrets store, and that there is a
	// JWT signing secret in there that can be used...
	// therefore there's no need to create it in the DB from scratch.
	// ADD!! - have a more robust was (flag) for checking if there is a 
	//         secret store setup for JWT secret fetching.
	if pRuntimeSys.ExternalPlugins != nil && 
		pRuntimeSys.ExternalPlugins.SecretGetCallback != nil {

	} else {
		
		existsBool, gfErr := ksDBjwtExistsSecret(pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, nil, gfErr
		}

		// JWT keys dont exist in the DB, so generate them
		if !existsBool {
			pRuntimeSys.LogNewFun("DEBUG",
				"JWT keys (private/public RSA keys) to be generated since it doesnt exist, and no secrets_store to read from...", nil)

			publicKey, privateKey, gfErr := ksJWTgenerateKeys(pCtx, pRuntimeSys)
			if gfErr != nil {
				return nil, nil, gfErr
			}
			return publicKey, privateKey, nil

		} else {

			// JWT keys are found in the DB, so return that
			publicKey, privateKey, gfErr := ksJWTgetKeysFromStore(pCtx, pRuntimeSys)
			if gfErr != nil {
				return nil, nil, gfErr
			}
			return publicKey, privateKey, nil
		}
	}

	return nil, nil, nil
}

//---------------------------------------------------

// generate and store in the DB the secret key thats used
// to sign new JWT tokens. this is only done if the user is self-hosting
// and doesnt have want to use a secrets store where they place the secret 
// separatelly from GF (and GF only fetches it from the secret store).
// this is also done only once on startup when that secret is detected
// not to exist.
func ksJWTgenerateKeys(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*rsa.PublicKey, *rsa.PrivateKey, *gf_core.GFerror) {

	// GENERATE
	pubKeyPEMstr, privKeyPEMstr := gf_core.CryptoGenerateKeysAsPEM()

	//------------------------
	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	fieldsForIDlst := []string{
		"jwt_secret",
	}
	IDstr := gf_core.IDcreate(fieldsForIDlst,
		creationUNIXtimeF)

	jwtSecret := &GFjwtSecret{
		Vstr:              "0",
		IDstr:             IDstr,
		DeletedBool:       false,
		CreationUNIXtimeF: creationUNIXtimeF,
		
		PublicKeyPEMstr:   GFjwtPublicKeyPEMval(pubKeyPEMstr),
		PrivateKeyPEMstr:  GFjwtPrivateKeyPEMval(privKeyPEMstr),
	}

	// DB_CREATE
	gfErr := ksDBcreateJWTsecret(jwtSecret, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	//------------------------
	publicKey, privateKey, gfErr := gf_core.CryptoParseKeysFromPEM(pubKeyPEMstr, privKeyPEMstr, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}
	return publicKey, privateKey, nil
}

//---------------------------------------------------

func ksJWTgetKeysFromStore(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*rsa.PublicKey, *rsa.PrivateKey, *gf_core.GFerror) {
	

	pRuntimeSys.LogNewFun("DEBUG", "getting JWT keys (private/public RSA keys)...", nil)

	var jwtPrivateKeyPEMvalStr string // for JWT signing
	var jwtPublicKeyPEMvalStr string  // for JWT Validation

	// SECRETS_STORE
	if pRuntimeSys.ExternalPlugins != nil &&
		pRuntimeSys.ExternalPlugins.SecretGetCallback != nil {

		secretNameStr := fmt.Sprintf("gf_jwt_keypair_%s", pRuntimeSys.EnvStr)

		pRuntimeSys.LogNewFun("DEBUG", "getting keys from secrets store", map[string]interface{}{
			"secret_name_str": secretNameStr,
		})

		// SECRET_GET
		secretMap, gfErr := pRuntimeSys.ExternalPlugins.SecretGetCallback(secretNameStr,
			pRuntimeSys)
		if gfErr != nil {
			return nil, nil, gfErr
		}

		jwtPublicKeyPEMvalFromAWSstr := secretMap["public_key_pem_str"].(string)
		jwtPublicKeyPEMvalStr = jwtPublicKeyPEMvalFromAWSstr

		jwtPrivateKeyPEMvalFromAWSstr := secretMap["private_key_pem_str"].(string)
		jwtPrivateKeyPEMvalStr = jwtPrivateKeyPEMvalFromAWSstr

	} else {

		pRuntimeSys.LogNewFun("DEBUG", "getting keys from DB", nil)

		// DB
		jwtSecretFromDB, gfErr := ksDBgetJWTsecret(pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, nil, gfErr
		}

		jwtPublicKeyPEMvalStr  = string(jwtSecretFromDB.PublicKeyPEMstr)
		jwtPrivateKeyPEMvalStr = string(jwtSecretFromDB.PrivateKeyPEMstr)
	}

	// parse PEM
	publicKey, privateKey, gfErr := gf_core.CryptoParseKeysFromPEM(jwtPublicKeyPEMvalStr,
		jwtPrivateKeyPEMvalStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gfErr
	}

	return publicKey, privateKey, nil
}

//---------------------------------------------------
// DB
//---------------------------------------------------

func ksDBcreateJWTsecret(pJWTsecret *GFjwtSecret,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	collNameStr := "gf_auth_jwt_secret"

	gfErr := gf_core.MongoInsert(pJWTsecret,
		collNameStr,
		map[string]interface{}{
			"id_str":             pJWTsecret.IDstr,
			"caller_err_msg_str": "failed to create jwt_secret in the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

// there should be only one valid (non-deleted) jwt secret in the DB,
// used for all users.
func ksDBgetJWTsecret(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFjwtSecret, *gf_core.GFerror) {

	findOpts := options.FindOne()
	
	jwtSecret := GFjwtSecret{}
	collNameStr := "gf_auth_jwt_secret"
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx, bson.M{
			"deleted_bool": false,
		},
		findOpts).Decode(&jwtSecret)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get a jwt_secret from the DB",
			"mongodb_find_error",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	return &jwtSecret, nil
}

//---------------------------------------------------

// check if a JWT secret exists in a DB, when a secrets-storage
// backend is not being used (by users that self-host and use the DB for everything).
func ksDBjwtExistsSecret(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	collNameStr := "gf_auth_jwt_secret"
	countInt, gfErr := gf_core.MongoCount(bson.M{
			"deleted_bool": false,
		},
		map[string]interface{}{
			"caller_err_msg": "failed to check if there is a jwt_secret in the DB",
		},
		pRuntimeSys.Mongo_db.Collection(collNameStr),
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	if countInt > 0 {
		return true, nil
	}
	return false, nil
}