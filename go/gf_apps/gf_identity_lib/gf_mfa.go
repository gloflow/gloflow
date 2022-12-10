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

package gf_identity_lib

import (
	// "fmt"
	"time"
	"bytes"
	"strings"
	"strconv"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
)

//------------------------------------------------

type GFuserAuthMFAinputConfirm struct {
	UserNameStr        gf_identity_core.GFuserName `validate:"required,min=3,max=50"`
	ExternHtopValueStr string    `validate:"required,min=10,max=200"`
	SecretKeyBase32str string    `validate:"required,min=8,max=200"`
}

//------------------------------------------------

func mfaPipelineConfirm(pInput *GFuserAuthMFAinputConfirm,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	htopValueStr, gfErr := TOTPgenerateValue(pInput.SecretKeyBase32str,
		pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	if pInput.ExternHtopValueStr == htopValueStr {

		// get a preexisting login_attempt if one exists and hasnt expired for this user.
		// if it has then the user will have to restart the login flow
		// (which will create a new login_attempt).
		var loginAttempt *gf_identity_core.GFloginAttempt
		loginAttempt, gfErr = gf_identity_core.LoginAttemptGetIfValid(pInput.UserNameStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}

		// there is no login_attempt for this user thats active, or it has expired.
		// do nothing, forcing the user to restart the login process.
		if loginAttempt == nil {
			return false, nil
		}

		//------------------------
		// UPDATE_LOGIN_ATTEMPT
		// if password is valid then update the login_attempt 
		// to indicate that the password has been confirmed
		mfaConfirmBool := true
		updateOp := &gf_identity_core.GFloginAttemptUpdateOp{MFAconfirmedBool: &mfaConfirmBool}
		gfErr = gf_identity_core.DBloginAttemptUpdate(&loginAttempt.IDstr,
			updateOp,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}

		//------------------------
		return true, nil
		
	} else {
		return false, nil
	}
	return false, nil
}

//---------------------------------------------------
// TOTP - https://datatracker.ietf.org/doc/html/rfc6238

func TOTPgenerateValue(pSecretKeyStr string,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// index number of a 30s time period since start of Unix time.
	// TOTP is specified to use UTC time, so timezones dont have to be 
	// accounted for between multiple parties.
	intervalInt := time.Now().UTC().Unix() / 30

	pRuntimeSys.LogNewFun("DEBUG", "TOTP value generated", map[string]interface{}{"interval_int": intervalInt,})
	
	tokenStr, gfErr  := HOTPgenerateValue(pSecretKeyStr,
		intervalInt,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return tokenStr, nil
}

//---------------------------------------------------
// HOTP_GENERATE_TOKEN
// HOTP - https://datatracker.ietf.org/doc/html/rfc4226

// secret_key - symmetric key

func HOTPgenerateValue(pSecretKeyBase32str string,
	pTimeIntervalInt int64,
	pRuntimeSys      *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// expected length, can be 6-10 long, google-authenticator uses 6
	hotpTokenLengthInt := 6

	// convert secret_key to base32.
	// one way to represent Base32 numbers in a human-readable way is by using 
	// a standard 32-character set, such as the twenty-two upper-case letters A–V and the digits 0-9.
	// its done in order to allow for this limited alphabet.
	// first switch all characters to upper case just in case.
	keyUppercaseBase32str := strings.ToUpper(pSecretKeyBase32str)
	keyBytesLst, err      := base32.StdEncoding.DecodeString(keyUppercaseBase32str)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to base32 encode secret key for HOTP token generation",
			"generic_error",
			map[string]interface{}{
				"time_interval_int": pTimeIntervalInt,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return "", gfErr
	}

	// BIG_ENDIAN - byte order where the highest-order bit is stored first (lower memory address),
	//              and least last (higher address).
	// COUNTER - time_interval is the counter in the HOTP standard, defined to be 8 bytes.
	timeIntervalBytesLst := make([]byte, 8)
	binary.BigEndian.PutUint64(timeIntervalBytesLst, uint64(pTimeIntervalInt))

	// sign secret_key using HMAC-SHA1
	// SHA1 - 160-bit string (20 bytes), 40 digits in hex
	// HMAC - HMAC is a specific type of message authentication code involving a cryptographic
	//        hash function and a secret cryptographic key. As with any MAC, it may be used to
	//        simultaneously verify both the data integrity and authenticity of a message.
	//        trades off the need for a complex public key infrastructure by delegating the 
	//        key exchange to the communicating parties, who are responsible for establishing 
	//        and using a trusted channel to agree on the key prior to communication.
	// HMAC-SHA-1(K,C) - this is a hash that is a function of the key and counter (interval)
	hash := hmac.New(sha1.New, keyBytesLst)
	hash.Write(timeIntervalBytesLst)
	hash_bytes_lst := hash.Sum(nil)

	//-----------------------
	// TRUNCATE - using a subset of the hash
	// use last half-byte to choose the index in the hash to start from.
	// this number can be a maximum of 15 (0xf), whereas sha1 hash is 20 bytes,
	// which leaves exactly 4 bytes which is needed.
	offsetInt := (hash_bytes_lst[len(hash_bytes_lst)-1] & 0xf)

	var headerInt uint32
	
	// generate 4 byte (32 bit) chunk from hash, starting at offset
	chunk4bytesLst := hash_bytes_lst[offsetInt : offsetInt + 4]
	reader         := bytes.NewReader(chunk4bytesLst)
	err = binary.Read(reader, binary.BigEndian, &headerInt)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to base32 encode secret key for HOTP token generation",
			"io_reader_error",
			map[string]interface{}{
				"time_interval_int": pTimeIntervalInt,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return "", gfErr
	}

	//-----------------------

	/*code := int64(chunk4bytesLst[0])<<24 |
		int64(chunk4bytesLst[1])<<16 |
		int64(chunk4bytesLst[2])<<8 |
		int64(chunk4bytesLst[3])*/

	
	// 1_000_000 - is 10^d - "d" represents a "d" number of least significant digits,
	//             which are selected out here for usage.
	//             (ignoring the most significant bits)
	htopValue := (int(headerInt) & 0x7fffffff) % 1_000_000

	htopValueStr := strconv.Itoa(int(htopValue))

	htopValuePaddedStr := padWith0s(htopValueStr, hotpTokenLengthInt)
	return htopValuePaddedStr, nil
}

//---------------------------------------------------
// PAD_WITH_0S - append/pad string with 0's until its of target_length

func padWith0s(p_str string, p_target_length_int int) string {
	if len(p_str) >= p_target_length_int {
		return p_str
	}
	for i := (p_target_length_int - len(p_str)); i > 0; i-- {
		p_str = "0" + p_str
	}
	return p_str
}