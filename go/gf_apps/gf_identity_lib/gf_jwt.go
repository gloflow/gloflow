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

package gf_identity_lib

import (
	"time"
	"context"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type GF_jwt_val string
type GF_jwt_token struct {
	V_str                string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`

	Val_str          GF_jwt_val          `bson:"val_str"`
	User_address_eth GF_user_address_eth `bson:"user_address_eth"`
}

type GF_jwt_claims struct {
	User_address_eth GF_user_address_eth `json:"user_address_eth"`
	jwt.StandardClaims
}

//---------------------------------------------------
// PIPELINE__GENERATE
func jwt__pipeline__generate(p_user_address_eth GF_user_address_eth,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (GF_jwt_val, *gf_core.GF_error) {

	

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0



	// JWT_GENERATE

	signing_key_str := gf_core.Str_random()
	jwt_token_val, gf_err := jwt__generate(p_user_address_eth,
		signing_key_str,
		creation_unix_time_f,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}



	jwt_id := jwt__generate_id(p_user_address_eth, creation_unix_time_f)
	jwt_token := &GF_jwt_token{
		V_str:                "0",
		Id_str:               jwt_id,
		Creation_unix_time_f: creation_unix_time_f,
		Val_str:              jwt_token_val,
		User_address_eth:     p_user_address_eth,
	}



	// DB
	gf_err = db__jwt__create(jwt_token, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}




	return jwt_token_val, nil
}

//---------------------------------------------------
// GENERATE
func jwt__generate(p_user_address_eth GF_user_address_eth,
	p_signing_secret_key_str string,
	p_creation_unix_time_f   float64,
	p_runtime_sys            *gf_core.Runtime_sys) (GF_jwt_val, *gf_core.GF_error) {


	issuer_str := "gf"
	jwt_token_ttl_sec_int    := int64(60*60*24*7) // 7 days
	expiration_unix_time_int := int64(p_creation_unix_time_f) + jwt_token_ttl_sec_int

	// CLAIMS
	claims := GF_jwt_claims{
		p_user_address_eth,
		jwt.StandardClaims{
			ExpiresAt: expiration_unix_time_int,
			Issuer:    issuer_str, 
		},
	}

	// NEW_TOKEN
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SIGNING - to be able to verify using the same secret_key that in the future
	//           a received token is valid and unchanged.
	jwt_token_val_str, err := jwt_token.SignedString([]byte(p_signing_secret_key_str))
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to to update user info",
			"crypto_jwt_sign_token_error",
			map[string]interface{}{
				"user_address_eth": p_user_address_eth,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return GF_jwt_val(""), gf_err
	}

	return GF_jwt_val(jwt_token_val_str), nil
}

//---------------------------------------------------
func jwt__generate_id(p_user_address_eth GF_user_address_eth,
	p_creation_unix_time_f float64) gf_core.GF_ID {
	fields_for_id_lst := []string{
		string(p_user_address_eth),
	}
	gf_id_str := gf_core.Image_ID__md5_create(fields_for_id_lst,
		p_creation_unix_time_f)
	return gf_id_str
}

//---------------------------------------------------
func jwt__validate_from_req(p_user_eth_address GF_user_address_eth,
	p_req         *http.Request,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {




	



	return nil




}

//---------------------------------------------------
func jwt__validate(p_jwt_val GF_jwt_val,
	p_signing_secret_key_str string,
	p_runtime_sys            *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {

	claims := &jwt.MapClaims{}
	jwt_token, err := jwt.ParseWithClaims(string(p_jwt_val),
		claims,
		func(p_jwt_token *jwt.Token) (interface{}, error) {
			return []byte(p_signing_secret_key_str), nil
		})
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to verify a JWT token",
			"crypto_jwt_verify_token_error",
			map[string]interface{}{
				"jwt_val_str": p_jwt_val,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return false, gf_err
	}

	valid_bool := jwt_token.Valid

	return valid_bool, nil
}

//---------------------------------------------------
// DB
//---------------------------------------------------
func db__jwt__create(p_jwt_token *GF_jwt_token,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.GF_error {


	return nil
}