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
	"fmt"
	"context"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func verify__auth_signature__all_methods(p_signature_str GF_auth_signature,
	p_nonce_str        GF_user_nonce_val,
	p_user_address_eth GF_user_address_eth,
	p_ctx              context.Context,
	p_runtime_sys      *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {

	
	// first attempt - try to verify using the data_header
	valid_bool, gf_err := verify__auth_signature(p_signature_str,
		string(p_nonce_str),
		p_user_address_eth,
		true,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return false, gf_err
	}



	if !valid_bool {

		// second attempt - dont validate using the data header
		valid_bool, gf_err = verify__auth_signature(p_signature_str,
			string(p_nonce_str),
			p_user_address_eth,
			false,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return false, gf_err
		}
	}

	return valid_bool, nil
}

//---------------------------------------------------
/*p_signature_nonce_str - In cryptography, a nonce is an arbitrary number that can be used
	just once in a cryptographic communication. It is often a random or 
	pseudo-random number issued in an authentication protocol 
	to ensure that old communications cannot be reused in replay attacks.*/

// https://goethereumbook.org/signature-verify/

func verify__auth_signature(p_signature_str GF_auth_signature,
	p_data_str                  string,
	p_user_address_eth          GF_user_address_eth,
	p_validate_data_header_bool bool,
	p_ctx                       context.Context,
	p_runtime_sys               *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {
	
	var final_data_str string
	if p_validate_data_header_bool {
		final_data_str = fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s",
			len(p_data_str),
			p_data_str)
	} else {
		final_data_str = p_data_str
	}
	


	


	// decode signature from hex form
	sig_decoded_bytes_lst, err := hexutil.Decode(string(p_signature_str))
	if err != nil {
		gf_err := gf_core.Error__create("failed to hex-decode a signature supplied for validation",
			"crypto_hex_decode",
			map[string]interface{}{},
			err, "gf_identity_lib", p_runtime_sys)
		return false, gf_err
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

	sig_last_byte := sig_decoded_bytes_lst[64]
	if sig_last_byte != 27 && sig_last_byte != 28 && sig_last_byte != 0 && sig_last_byte != 1 {
		gf_err := gf_core.Error__create("signature validation failed because the last byte (V value) of the signature is invalid",
			"crypto_signature_eth_last_byte_invalid_value",
			map[string]interface{}{
				"sig_last_byte": sig_last_byte,
			},
			nil, "gf_identity_lib", p_runtime_sys)
		return false, gf_err
	}

	//------------------------


	// bring the last byte of the singature to 0|1.
	// some wallets use 0|1 as recovery values (27/28 is a legacy value), so deduct
	// 27 only if these legacy values are used.
	if sig_decoded_bytes_lst[64] == 27 || sig_decoded_bytes_lst[64] == 28 {
		sig_decoded_bytes_lst[64] -= 27
	}


	// removing the last/recovery_id byte from the singature. used for verification.
	sig_no_recovery_id_bytes_lst := []byte(sig_decoded_bytes_lst[:len(sig_decoded_bytes_lst)-1])



	data_hash := crypto.Keccak256Hash([]byte(final_data_str))
	sig_public_key_ECDSA, err := crypto.SigToPub(data_hash.Bytes(), sig_decoded_bytes_lst)
	if err != nil {
		gf_err := gf_core.Error__create("signature validation failed because the last byte (V value) of the signature is invalid",
			"crypto_ec_recover_pubkey",
			map[string]interface{}{
				"sig_last_byte": sig_last_byte,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return false, gf_err
	}


	// CompressPubkey() - CompressPubkey encodes a public key to the 33-byte compressed format. 
	public_key_bytes_lst := crypto.CompressPubkey(sig_public_key_ECDSA)






	// generate an Ethereum address that corresponds to a given Public Key.
	user_address_eth_derived_str := crypto.PubkeyToAddress(*sig_public_key_ECDSA).Hex()


	//------------------------
	// VALIDITY_CHECK_2
	// compare addresses, suplied and derived (from pubkey), and if not the same
	// return right away and declare signature as invalid.
	if user_address_eth_derived_str != string(p_user_address_eth) {
		return false, nil
	}

	// VALIDITY_CHECK_3
	valid_bool := crypto.VerifySignature(public_key_bytes_lst,
		data_hash.Bytes(),
		sig_no_recovery_id_bytes_lst)

	//------------------------

	return valid_bool, nil
}