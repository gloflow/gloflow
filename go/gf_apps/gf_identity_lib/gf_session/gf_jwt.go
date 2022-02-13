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

package gf_session

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
type GF_jwt_token_val      string
type GF_jwt_secret_key_val string
type GF_jwt_secret_key struct {
	V_str                string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`

	Val                 GF_jwt_secret_key_val `bson:"val_str"`
	User_identifier_str string                `bson:"user_identifier_str"`
}

type GF_jwt_claims struct {
	User_identifier_str string `json:"user_identifier_str"`
	jwt.StandardClaims
}

//---------------------------------------------------
// PIPELINE__GENERATE
func JWT__pipeline__generate(p_user_identifier_str string, // p_user_address_eth GF_user_address_eth,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (GF_jwt_token_val, *gf_core.GF_error) {

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0

	// JWT_GENERATE
	jwt_secret_key_val_str := GF_jwt_secret_key_val(gf_core.Str_random())
	jwt_token_val, gf_err := jwt__generate(p_user_identifier_str,
		jwt_secret_key_val_str,
		creation_unix_time_f,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}

	jwt_id := jwt__generate_id(p_user_identifier_str, creation_unix_time_f)
	jwt_secret_key := &GF_jwt_secret_key{
		V_str:                "0",
		Id_str:               jwt_id,
		Deleted_bool:         false,
		Creation_unix_time_f: creation_unix_time_f,
		Val:                  jwt_secret_key_val_str,
		User_identifier_str:  p_user_identifier_str,
		// User_address_eth: p_user_address_eth,
	}

	// DB_CREATE__SECRET_KEY
	gf_err = db__jwt_secret_key__create(jwt_secret_key, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}

	return jwt_token_val, nil
}

//---------------------------------------------------
// GENERATE
func jwt__generate(p_user_identifier_str string, // p_user_address_eth GF_user_address_eth,
	p_jwt_secret_key_val   GF_jwt_secret_key_val,
	p_creation_unix_time_f float64,
	p_runtime_sys          *gf_core.Runtime_sys) (GF_jwt_token_val, *gf_core.GF_error) {


	issuer_str := "gf"
	jwt_token_ttl_sec_int    := int64(60*60*24*7) // 7 days
	expiration_unix_time_int := int64(p_creation_unix_time_f) + jwt_token_ttl_sec_int

	// CLAIMS
	claims := GF_jwt_claims{
		p_user_identifier_str,
		jwt.StandardClaims{
			ExpiresAt: expiration_unix_time_int,
			Issuer:    issuer_str, 
		},
	}

	// NEW_TOKEN
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SIGNING - to be able to verify using the same secret_key that in the future
	//           a received token is valid and unchanged.
	jwt_token_val_str, err := jwt_token.SignedString([]byte(p_jwt_secret_key_val))
	if err != nil {
		gf_err := gf_core.Error__create("failed to to update user info",
			"crypto_jwt_sign_token_error",
			map[string]interface{}{
				"user_identifier_str": p_user_identifier_str,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return GF_jwt_token_val(""), gf_err
	}

	return GF_jwt_token_val(jwt_token_val_str), nil
}

//---------------------------------------------------
func jwt__generate_id(p_user_identifier_str string, // p_user_address_eth GF_user_address_eth,
	p_creation_unix_time_f float64) gf_core.GF_ID {
	fields_for_id_lst := []string{
		p_user_identifier_str, // string(p_user_address_eth),
	}
	gf_id_str := gf_core.ID__create(fields_for_id_lst,
		p_creation_unix_time_f)
	return gf_id_str
}

//---------------------------------------------------
func jwt__pipeline__validate(p_jwt_token_val GF_jwt_token_val,
	p_ctx         context .Context,
	p_runtime_sys *gf_core.Runtime_sys) (string, *gf_core.GF_error) {

	// VALIDATE
	valid_bool, user_identifier_str, gf_err := jwt__validate(p_jwt_token_val,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}

	if !valid_bool {
		gf_err := gf_core.Error__create("JWT token supplied for validation is invalid",
			"crypto_jwt_verify_token_invalid_error",
			map[string]interface{}{
				"jwt_token_val_str": p_jwt_token_val,
			},
			nil, "gf_identity_lib", p_runtime_sys)
		return "", gf_err
	}

	return user_identifier_str, nil
}

//---------------------------------------------------
// VALIDATE
func jwt__validate(p_jwt_token_val GF_jwt_token_val,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (bool, string, *gf_core.GF_error) {

	claims := &GF_jwt_claims{}
	jwt_token, err := jwt.ParseWithClaims(string(p_jwt_token_val),
		claims,
		func(p_jwt_token *jwt.Token) (interface{}, error) {

			user_identifier_str := p_jwt_token.Claims.(*GF_jwt_claims).User_identifier_str

			// DB_GET
			jwt_secret_key, gf_err := db__jwt_secret_key__get(user_identifier_str, p_ctx, p_runtime_sys)
			if gf_err != nil {
				return nil, gf_err.Error
			}

			return []byte(jwt_secret_key.Val), nil
		})

	if err != nil {
		gf_err := gf_core.Error__create("failed to verify a JWT token",
			"crypto_jwt_verify_token_error",
			map[string]interface{}{
				"jwt_token_val_str": p_jwt_token_val,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return false, "", gf_err
	}

	valid_bool          := jwt_token.Valid
	user_identifier_str := jwt_token.Claims.(*GF_jwt_claims).User_identifier_str
	
	return valid_bool, user_identifier_str, nil
}

//---------------------------------------------------
// DB
//---------------------------------------------------
func db__jwt_secret_key__create(p_jwt_secret_key *GF_jwt_secret_key,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {

	coll_name_str := "gf_auth_jwt"

	gf_err := gf_core.Mongo__insert(p_jwt_secret_key,
		coll_name_str,
		map[string]interface{}{
			"id_str":              p_jwt_secret_key.Id_str,
			"user_identifier_str": p_jwt_secret_key.User_identifier_str,
			"caller_err_msg_str":  "failed to create jwt_secret_key for a user in a DB",
		},
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}

//---------------------------------------------------
func db__jwt_secret_key__get(p_user_identifier_str string,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_jwt_secret_key, *gf_core.GF_error) {

	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"creation_unix_time_f": -1}) // descending - true - sort the latest items first
	
	db_cursor, gf_err := gf_core.Mongo__find(bson.M{
			"user_identifier_str": string(p_user_identifier_str),
			"deleted_bool":        false,
		},
		find_opts,
		map[string]interface{}{
			"user_identifier_str": p_user_identifier_str,
			"caller_err_msg_str":  "failed to get jwt_secret_key for a user from DB",
		},
		p_runtime_sys.Mongo_db.Collection("gf_auth_jwt"),
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}



	var jwt_secret_keys_lst []*GF_jwt_secret_key
	err := db_cursor.All(p_ctx, &jwt_secret_keys_lst)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get DB results of query to get latest JWT key ",
			"mongodb_cursor_all",
			map[string]interface{}{
				"user_identifier_str": p_user_identifier_str,
			},
			err, "gf_tagger_lib", p_runtime_sys)
		return nil, gf_err
	}

	
	jwt_secret_key := jwt_secret_keys_lst[0]
	

	return jwt_secret_key, nil
}