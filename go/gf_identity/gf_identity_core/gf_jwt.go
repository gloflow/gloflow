/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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
	"time"
	"context"
	"crypto/rsa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/golang-jwt/jwt"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFjwtTokenVal     string
type GFjwtSecretKeyVal string
type GFjwtSecretKeyPEMval string

type GFjwtSecret struct {
	Vstr                 string                `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID    `bson:"_id,omitempty"`
	IDstr                gf_core.GF_ID         `bson:"id_str"`
	DeletedBool          bool                  `bson:"deleted_bool"`
	CreationUNIXtimeF    float64               `bson:"creation_unix_time_f"`
	Val                  GFjwtSecretKeyPEMval  `bson:"val_str"`
}

type GFjwtClaims struct {
	UserIdentifierStr string `json:"user_identifier_str"`
	jwt.StandardClaims
}

//---------------------------------------------------
// GENERATE
//---------------------------------------------------

func JWTpipelineGenerate(pUserIdentifierStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFjwtTokenVal, *gf_core.GFerror) {

	privKey, gfErr := jwtGetSigningPrivKey(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	fmt.Println(privKey)
	signKey := privKey

	//----------------------
	// JWT_GENERATE
	tokenValStr, gfErr := jwtGenerate(pUserIdentifierStr,
		signKey,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//----------------------

	return tokenValStr, nil
}

//---------------------------------------------------

func jwtGenerate(pUserIdentifierStr string,
	pSignKey    *rsa.PrivateKey,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFjwtTokenVal, *gf_core.GFerror) {
	
	issuerStr := "gf"
	_, jwtTokenTTLsecInt  := GetSessionTTL()
	creationUNIXtimeF     := float64(time.Now().UnixNano())/1000000000.0
	expirationUNIXtimeInt := int64(creationUNIXtimeF) + jwtTokenTTLsecInt


	// CLAIMS
	claims := GFjwtClaims{
		pUserIdentifierStr,
		jwt.StandardClaims{
			ExpiresAt: expirationUNIXtimeInt,
			Issuer:    issuerStr, 
		},
	}

	// claims := token.Claims.(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(10 * time.Minute)
	// claims["authorized"] = true
	// claims["user"] = "username"


	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	jwtTokenSignedStr, err := jwtToken.SignedString(pSignKey)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to sign JWT token for user",
			"crypto_jwt_sign_token_error",
			map[string]interface{}{
				"user_identifier_str": pUserIdentifierStr,
			},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}


	tokenValStr := GFjwtTokenVal(jwtTokenSignedStr)
	return &tokenValStr, nil
}

//---------------------------------------------------

// used only by users that self-host and dont use a dedicated secret store.
// instead they store all data in the DB for max simplicity of hosting.
func JWTgenerateSigningSecretIfAbsent(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	// if a secrets "get" callback is specified we're making an assumption
	// that the user has setup some sort of secrets store, and that there is a
	// JWT signing secret in there that can be used...
	// therefore there's no need to create it in the DB from scratch.
	// ADD!! - have a more robust was (flag) for checking if there is a 
	//         secret store setup for JWT secret fetching.
	if pRuntimeSys.ExternalPlugins.SecretGetCallback != nil {

	} else {
	
		gfErr := JWTgenerateSigningSecret(pCtx, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}
	}

	return nil
}

//---------------------------------------------------

// generate and store in the DB the secret key thats used
// to sign new JWT tokens. this is only done if the user is self-hosting
// and doesnt have want to use a secrets store where they place the secret 
// separatelly from GF (and GF only fetches it from the secret store).
// this is also done only once on startup when that secret is detected
// not to exist.
func JWTgenerateSigningSecret(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {


	privKeyPEMstr := gf_core.CryptoGeneratePrivKeyAsPEM()


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
		Val:               GFjwtSecretKeyPEMval(privKeyPEMstr),
	}

	// DB_CREATE__SECRET_KEY
	gfErr := dbJWTcreateSigningSecret(jwtSecret, pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func jwtGetSigningPrivKey(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*rsa.PrivateKey, *gf_core.GFerror) {

	var jwtSigningSecretPEMvalStr string

	// SECRETS_STORE
	if pRuntimeSys.ExternalPlugins.SecretGetCallback != nil {

		secretNameStr := fmt.Sprintf("gf_jwt_signing_secret_%s", pRuntimeSys.EnvStr)

		// SECRET_GET
		secretMap, gfErr := pRuntimeSys.ExternalPlugins.SecretGetCallback(secretNameStr,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		jwtSigningSecretPEMvalFromAWSstr := secretMap["val_str"].(string)
		jwtSigningSecretPEMvalStr = jwtSigningSecretPEMvalFromAWSstr

	} else {

		// DB
		jwtSigningSecretPEMvalFromDBstr, gfErr := dbJWTgetSigningSecret(pCtx, pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}

		jwtSigningSecretPEMvalStr = string(jwtSigningSecretPEMvalFromDBstr.Val)
	}

	// parse PEM
	privKey, gfErr := gf_core.CryptoParsePrivKeyFromPEM(jwtSigningSecretPEMvalStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return privKey, nil
}

//---------------------------------------------------
// VALIDATE
//---------------------------------------------------

func JWTpipelineValidate(pJWTtokenVal GFjwtTokenVal,
	pCtx        context .Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	return "", nil
}

func JWTvalidate(pJWTtokenVal GFjwtTokenVal,
	pPubKey     *rsa.PublicKey,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {

	jwtToken, err := jwt.Parse(string(pJWTtokenVal), func(pToken *jwt.Token) (interface{}, error) {

		return pPubKey, nil
	})

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to verify a JWT token",
			"crypto_jwt_verify_token_error",
			map[string]interface{}{
				"jwt_token_val_str": pJWTtokenVal,
			},
			err, "gf_identity_core", pRuntimeSys)
		return false, "", gfErr
	}

	validBool         := jwtToken.Valid
	userIdentifierStr := jwtToken.Claims.(*GFjwtClaims).UserIdentifierStr

	return validBool, userIdentifierStr, nil
}

//---------------------------------------------------
// DB
//---------------------------------------------------

func dbJWTcreateSigningSecret(pJWTsecret *GFjwtSecret,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	collNameStr := "gf_auth_jwt_secret"

	gfErr := gf_core.MongoInsert(pJWTsecret,
		collNameStr,
		map[string]interface{}{
			"id_str":             pJWTsecret.IDstr,
			"caller_err_msg_str": "failed to create jwt_secret_key in the DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

// there should be only one valid (non-deleted) jwt signing secret in the DB,
// used for all users.
func dbJWTgetSigningSecret(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFjwtSecret, *gf_core.GFerror) {

	findOpts := options.FindOne()
	
	jwtSecret := GFjwtSecret{}
	collNameStr := "gf_auth_jwt_secret"
	err := pRuntimeSys.Mongo_db.Collection(collNameStr).FindOne(pCtx, bson.M{
			"deleted_bool": false,
		},
		findOpts).Decode(&jwtSecret)

	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get a JWT signing secret from the DB",
			"mongodb_find_error",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return nil, gfErr
	}

	return &jwtSecret, nil
}

//---------------------------------------------------

// check if a JWT signing secret exists in a DB, when a secrets-storage
// backend is not being used (by users that self-host and use the DB for everything).
func dbJWTexistsSigningSecret(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	collNameStr := "gf_auth_jwt_secret"
	countInt, gfErr := gf_core.MongoCount(bson.M{
			"deleted_bool":  false,
		},
		map[string]interface{}{
			"caller_err_msg": "failed to check if there is a JWT signing secret in the DB",
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