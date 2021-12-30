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

package gf_crypto

import (
	// "fmt"
    "strings"
	"crypto/ecdsa"
	"golang.org/x/crypto/sha3"
	"github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"
)

//-------------------------------------------------
func Eth_generate_keys() (string, string, string, error) {

    //---------------------------
	// PRIVATE_KEY_GENERATE
	private_key, err := crypto.GenerateKey()
    if err != nil {
        return "", "", "", err
    }

    private_key_bytes_lst := crypto.FromECDSA(private_key)

    // generate private_key hex, and remove the "0x" prefix
    private_key_hex_str := hexutil.Encode(private_key_bytes_lst) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

    //---------------------------
    // PUBLIC_KEY_GENERATE
    public_key := private_key.Public()
    public_key_ECDSA, ok := public_key.(*ecdsa.PublicKey)
    if !ok {
        return "", "", "", err
    }

    public_key_bytes_lst := crypto.FromECDSAPub(public_key_ECDSA)
    public_key_hex_str   := hexutil.Encode(public_key_bytes_lst)

    // [4:] - removing the first 4 bytes:
    //  - first 2  - "0x"
    //  - second 2 - "04" - this is a prefix to indicate an uncompressed public key ("03" is for a compressed one)
    // fmt.Println(hexutil.Encode(public_key_bytes_lst)[4:])

    //---------------------------
    // ADDRESS_GENERATE - from public_key

    // PubkeyToAddress() - helper method, gives the same result as the bellow
    //                     derivation using Keccak256 hash approach.
    // address_hex_str := crypto.PubkeyToAddress(*public_key_ECDSA).Hex()

    hash := sha3.NewLegacyKeccak256()
    hash.Write(public_key_bytes_lst[1:])
    // take the last 40 chars (20 bytes) of the pub_key hash, and get a hex of that
    address_hex_str := hexutil.Encode(hash.Sum(nil)[12:]) // 0x96216849c49358b10257cb55b28ea603c874b05e

    //---------------------------
	return private_key_hex_str, public_key_hex_str, address_hex_str, nil
}

//-------------------------------------------------
func Eth_sign_data(p_data_to_sign_str string,
	p_private_key_hex_str string) (string, error) {

    // clearn private_key hex
    var private_key_hex_clean_str string;
    if (strings.HasPrefix(p_private_key_hex_str, "0x")) {
        private_key_hex_clean_str = strings.TrimPrefix(p_private_key_hex_str, "0x")
    } else {
        private_key_hex_clean_str = p_private_key_hex_str
    }

    // parse private_key hex
	private_key, err := crypto.HexToECDSA(private_key_hex_clean_str)
    if err != nil {
        return "", err
    }

	data_to_sign_bytes_lst := []byte(p_data_to_sign_str)
	data_to_sign_hash := crypto.Keccak256Hash(data_to_sign_bytes_lst)

    //---------------------------
    // SIGN
	signature, err := crypto.Sign(data_to_sign_hash.Bytes(), private_key)
	if err != nil {
		return "", err
	}

	signature_hex_str := hexutil.Encode(signature)

    //---------------------------
	return signature_hex_str, err
}