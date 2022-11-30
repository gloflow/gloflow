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
	"context"
	"strings"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func verifyAuthSignatureAllMethods(pSignatureStr GFauthSignature,
	pNonceStr       GFuserNonceVal,
	pUserAddressETH GFuserAddressETH,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {

	// first attempt - try to verify using the data_header
	validBool, gfErr := verifyAuthSignature(pSignatureStr,
		string(pNonceStr),
		pUserAddressETH,
		true,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}

	if !validBool {

		// second attempt - dont validate using the data header
		validBool, gfErr = verifyAuthSignature(pSignatureStr,
			string(pNonceStr),
			pUserAddressETH,
			false,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}
	}

	return validBool, nil
}

//---------------------------------------------------
/*p_signature_nonce_str - In cryptography, a nonce is an arbitrary number that can be used
	just once in a cryptographic communication. It is often a random or 
	pseudo-random number issued in an authentication protocol 
	to ensure that old communications cannot be reused in replay attacks.*/

// https://goethereumbook.org/signature-verify/

func verifyAuthSignature(pSignatureStr GFauthSignature,
	pDataStr                string,
	pUserAddressETH         GFuserAddressETH,
	pValidateDataHeaderBool bool,
	pCtx                    context.Context,
	pRuntimeSys             *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {
	
	var finalDataStr string
	if pValidateDataHeaderBool {
		finalDataStr = fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s",
			len(pDataStr),
			pDataStr)
	} else {
		finalDataStr = pDataStr
	}
	


	


	// decode signature from hex form
	sigDecodedBytesLst, err := hexutil.Decode(string(pSignatureStr))
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to hex-decode a signature supplied for validation",
			"crypto_hex_decode",
			map[string]interface{}{},
			err, "gf_identity_lib", pRuntimeSys)
		return false, gfErr
	}



	//------------------------
	// VALIDITY_CHECK_1

	/*https://github.com/ethereum/go-ethereum/issues/19751
	@Péter Szilágyi/karalabe - geth lead - Jun 24, 2019:
	Originally Ethereum used 27 / 28 (which internally is just 0 / 1, just some weird bitcoin legacy to add 27).
	Later when we needed to support chain IDs in the signatures, the V as changed to ID*2 + 35 / ID*2 + 35.
	However, both V's are still supported on mainnet (Homestead vs. EIP155).
	The code was messy to pass V's around from low level crypto primitives in 27/28 notation,
	and then later for EIP155 to subtract 27, then do the whole x2+35 magic.
	The current logic is that the low level crypto operations returns 0/1 (because that is the canonical V value),
	and the higher level signers (Frontier, Homestead, EIP155) convert that V to whatever Ethereum specs on top of secp256k1.
	Use the high level signers, don't use the secp256k1 library directly.
	If you use the low level crypto library directly, 
	you need to be aware of how generic ECC relates to Ethereum signatures.*/

	sigLastByte := sigDecodedBytesLst[64]
	if sigLastByte != 27 && sigLastByte != 28 && sigLastByte != 0 && sigLastByte != 1 {
		gfErr := gf_core.ErrorCreate("signature validation failed because the last byte (V value) of the signature is invalid",
			"crypto_signature_eth_last_byte_invalid_value",
			map[string]interface{}{
				"sig_last_byte": sigLastByte,
			},
			nil, "gf_identity_lib", pRuntimeSys)
		return false, gfErr
	}

	//------------------------


	// bring the last byte of the singature to 0|1.
	// some wallets use 0|1 as recovery values (27/28 is a legacy value), so deduct
	// 27 only if these legacy values are used.
	if sigDecodedBytesLst[64] == 27 || sigDecodedBytesLst[64] == 28 {
		sigDecodedBytesLst[64] -= 27
	}


	// removing the last/recovery_id byte from the singature. used for verification.
	sig_no_recovery_id_bytes_lst := []byte(sigDecodedBytesLst[:len(sigDecodedBytesLst)-1])



	dataHash := crypto.Keccak256Hash([]byte(finalDataStr))
	publicKeyECDSA, err := crypto.SigToPub(dataHash.Bytes(), sigDecodedBytesLst)
	if err != nil {
		gfErr := gf_core.ErrorCreate("signature validation failed because the last byte (V value) of the signature is invalid",
			"crypto_ec_recover_pubkey",
			map[string]interface{}{
				"sig_last_byte": sigLastByte,
			},
			err, "gf_identity_lib", pRuntimeSys)
		return false, gfErr
	}


	// CompressPubkey() - CompressPubkey encodes a public key to the 33-byte compressed format.
	// compressed pub-key - uncompressed public key as a concatenation of Eliptic Curve x and y2 (as oppose to y1).
	// 						this format was used earlier by wallets, but as the blockchain started to grow there was 
	//                      a need to compress a public key. That's why it was decided to create a compressed format
	//                      for a public key that would use two times less space in memory by removing the y coordinate.
	//                      (because it can be calculated from the x coordinate by passing it to the y^2 = x^3 + 3 equation).
	//                      now there is only x coordinate as a public key plus a prefix that defines whether
	//                      the y should be negative or positive.
	// they start with the:
	//		- prefix 03 (for negative y compressed public key)
	//		- 04 (for uncompressed public key)
	publicKeyBytesLst := crypto.CompressPubkey(publicKeyECDSA)

	// generate an Ethereum address that corresponds to a given Public Key.
	// secp256k1.CompressPubkey(pubkey.X, pubkey.Y)
	userAddressETHderivedStr := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	//------------------------
	// VALIDITY_CHECK_2
	// compare addresses, suplied and derived (from pubkey), and if not the same
	// return right away and declare signature as invalid.
	if strings.ToLower(userAddressETHderivedStr) != strings.ToLower(string(pUserAddressETH)) {
		fmt.Println("eth derived address and supplied eth address are not the same",
			userAddressETHderivedStr,
			string(pUserAddressETH))
		return false, nil
	}

	// VALIDITY_CHECK_3
	validBool := crypto.VerifySignature(publicKeyBytesLst,
		dataHash.Bytes(),
		sig_no_recovery_id_bytes_lst)

	//------------------------

	return validBool, nil
}