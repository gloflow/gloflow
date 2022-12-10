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
	// "fmt"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

type GFjwtTokenVal     string
type GFjwtSecretKeyVal string
type GFjwtSecretKey struct {
	Vstr                 string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	IDstr                gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	CreationUNIXtimeF    float64            `bson:"creation_unix_time_f"`

	Val                 GFjwtSecretKeyVal   `bson:"val_str"`
	UserIdentifierStr   string              `bson:"user_identifier_str"`
}

type GFjwtClaims struct {
	UserIdentifierStr string `json:"user_identifier_str"`
	jwt.StandardClaims
}

//---------------------------------------------------
// PIPELINE__GENERATE

func JWTpipelineGenerate(pUserIdentifierStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (GFjwtTokenVal, *gf_core.GFerror) {

	//----------------------
	// JWT_SECRET_KEY_GENERATE
	jwtSecretKeyValStr, gfErr := JWTgenerateSecretSigningKey(pUserIdentifierStr,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	//----------------------
	// JWT_GENERATE
	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0

	jwtTokenVal, gfErr := jwtGenerate(pUserIdentifierStr,
		jwtSecretKeyValStr,
		creationUNIXtimeF,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	//----------------------

	return jwtTokenVal, nil
}

//---------------------------------------------------

// generate and store in the DB the secret key thats used
// to sign new JWT tokens
func JWTgenerateSecretSigningKey(pUserIdentifierStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (GFjwtSecretKeyVal, *gf_core.GFerror) {

	jwtSecretKeyValStr := GFjwtSecretKeyVal(gf_core.StrRandom())

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	jwtID := jwtGenerateID(pUserIdentifierStr, creationUNIXtimeF)

	jwtSecretKey := &GFjwtSecretKey{
		Vstr:              "0",
		IDstr:             jwtID,
		Deleted_bool:      false,
		CreationUNIXtimeF: creationUNIXtimeF,
		Val:               jwtSecretKeyValStr,
		UserIdentifierStr: pUserIdentifierStr,
	}

	// DB_CREATE__SECRET_KEY
	gfErr := dbJWTsecretKeyCreate(jwtSecretKey, pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return jwtSecretKeyValStr, nil
}

//---------------------------------------------------
// GENERATE

func jwtGenerate(pUserIdentifierStr string,
	pJWTsecretKeyVal   GFjwtSecretKeyVal,
	pCreationUNIXtimeF float64,
	pRuntimeSys        *gf_core.RuntimeSys) (GFjwtTokenVal, *gf_core.GFerror) {


	issuerStr := "gf"
	_, jwtTokenTTLsecInt  := GetSessionTTL()
	expirationUNIXtimeInt := int64(pCreationUNIXtimeF) + jwtTokenTTLsecInt

	// CLAIMS
	claims := GFjwtClaims{
		pUserIdentifierStr,
		jwt.StandardClaims{
			ExpiresAt: expirationUNIXtimeInt,
			Issuer:    issuerStr, 
		},
	}

	// NEW_TOKEN
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SIGNING - to be able to verify using the same secret_key that in the future
	//           a received token is valid and unchanged.
	jwtTokenValStr, err := jwtToken.SignedString([]byte(pJWTsecretKeyVal))
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to to update user info",
			"crypto_jwt_sign_token_error",
			map[string]interface{}{
				"user_identifier_str": pUserIdentifierStr,
			},
			err, "gf_session", pRuntimeSys)
		return GFjwtTokenVal(""), gfErr
	}

	return GFjwtTokenVal(jwtTokenValStr), nil
}

//---------------------------------------------------

func jwtGenerateID(pUserIdentifierStr string,
	pCreationUNIXtimeF float64) gf_core.GF_ID {
	
	fields_for_id_lst := []string{
		pUserIdentifierStr,
	}
	gf_id_str := gf_core.IDcreate(fields_for_id_lst,
		pCreationUNIXtimeF)
	return gf_id_str
}

//---------------------------------------------------

// validate a supplied JWT token, and return a user identifier
// stored in the JWT token
func JWTpipelineValidate(pJWTtokenVal GFjwtTokenVal,
	pCtx        context .Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// VALIDATE
	validBool, userIdentifierStr, gfErr := JWTvalidate(pJWTtokenVal,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	if !validBool {
		gfErr := gf_core.ErrorCreate("JWT token supplied for validation is invalid",
			"crypto_jwt_verify_token_invalid_error",
			map[string]interface{}{
				"jwt_token_val_str": pJWTtokenVal,
			},
			nil, "gf_session", pRuntimeSys)
		return "", gfErr
	}

	return userIdentifierStr, nil
}

//---------------------------------------------------
// VALIDATE

func JWTvalidate(pJWTtokenVal GFjwtTokenVal,
	pCtx         context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, string, *gf_core.GFerror) {

	claims := &GFjwtClaims{}
	jwtToken, err := jwt.ParseWithClaims(string(pJWTtokenVal),
		claims,
		func(pJWTtoken *jwt.Token) (interface{}, error) {

			userIdentifierStr := pJWTtoken.Claims.(*GFjwtClaims).UserIdentifierStr

			// DB_GET
			jwtSecretKey, gfErr := dbJWTsecretKeyGet(userIdentifierStr, pCtx, pRuntimeSys)
			if gfErr != nil {
				return nil, gfErr.Error
			}

			return []byte(jwtSecretKey.Val), nil
		})

	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to verify a JWT token",
			"crypto_jwt_verify_token_error",
			map[string]interface{}{
				"jwt_token_val_str": pJWTtokenVal,
			},
			err, "gf_session", pRuntimeSys)
		return false, "", gfErr
	}

	validBool         := jwtToken.Valid
	userIdentifierStr := jwtToken.Claims.(*GFjwtClaims).UserIdentifierStr
	
	return validBool, userIdentifierStr, nil
}

//---------------------------------------------------
// DB
//---------------------------------------------------

func dbJWTsecretKeyCreate(pJWTsecretKey *GFjwtSecretKey,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	collNameStr := "gf_auth_jwt"

	gfErr := gf_core.MongoInsert(pJWTsecretKey,
		collNameStr,
		map[string]interface{}{
			"id_str":              pJWTsecretKey.IDstr,
			"user_identifier_str": pJWTsecretKey.UserIdentifierStr,
			"caller_err_msg_str":  "failed to create jwt_secret_key for a user in a DB",
		},
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//---------------------------------------------------

func dbJWTsecretKeyGet(pUserIdentifierStr string,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (*GFjwtSecretKey, *gf_core.GFerror) {

	findOpts := options.Find()
	findOpts.SetSort(map[string]interface{}{"creation_unix_time_f": -1}) // descending - true - sort the latest items first
	
	dbCursor, gfErr := gf_core.MongoFind(bson.M{
			"user_identifier_str": string(pUserIdentifierStr),
			"deleted_bool":        false,
		},
		findOpts,
		map[string]interface{}{
			"user_identifier_str": pUserIdentifierStr,
			"caller_err_msg_str":  "failed to get jwt_secret_key for a user from DB",
		},
		pRuntimeSys.Mongo_db.Collection("gf_auth_jwt"),
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}



	var jwtSecretKeysLst []*GFjwtSecretKey
	err := dbCursor.All(pCtx, &jwtSecretKeysLst)
	if err != nil {
		gfErr := gf_core.MongoHandleError("failed to get DB results of query to get latest JWT key ",
			"mongodb_cursor_all",
			map[string]interface{}{
				"user_identifier_str": pUserIdentifierStr,
			},
			err, "gf_session", pRuntimeSys)
		return nil, gfErr
	}

	jwtSecretKey := jwtSecretKeysLst[0]

	return jwtSecretKey, nil
}