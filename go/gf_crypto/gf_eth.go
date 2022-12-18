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

func EthGenerateKeys() (string, string, string, error) {

    //---------------------------
	// PRIVATE_KEY_GENERATE
	privateKey, err := crypto.GenerateKey()
    if err != nil {
        return "", "", "", err
    }

    privateKeyBytesLst := crypto.FromECDSA(privateKey)

    // generate private_key hex, and remove the "0x" prefix
    privateKeyHexStr := hexutil.Encode(privateKeyBytesLst) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

    //---------------------------
    // PUBLIC_KEY_GENERATE
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        return "", "", "", err
    }

    publicKeyBytesLst := crypto.FromECDSAPub(publicKeyECDSA)
    publicKeyHexStr   := hexutil.Encode(publicKeyBytesLst)

    // [4:] - removing the first 4 bytes:
    //  - first 2  - "0x"
    //  - second 2 - "04" - this is a prefix to indicate an uncompressed public key ("03" is for a compressed one)
    // fmt.Println(hexutil.Encode(publicKeyBytesLst)[4:])

    //---------------------------
    // ADDRESS_GENERATE - from public_key

    // PubkeyToAddress() - helper method, gives the same result as the bellow
    //                     derivation using Keccak256 hash approach.
    // addressHexStr := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

    hash := sha3.NewLegacyKeccak256()
    hash.Write(publicKeyBytesLst[1:])
    // take the last 40 chars (20 bytes) of the pub_key hash, and get a hex of that
    addressHexStr := hexutil.Encode(hash.Sum(nil)[12:]) // 0x96216849c49358b10257cb55b28ea603c874b05e

    //---------------------------
	return privateKeyHexStr, publicKeyHexStr, addressHexStr, nil
}

//-------------------------------------------------

func EthSignData(pDataToSignStr string,
	pPrivateKeyHexStr string) (string, error) {

    // clean private_key hex
    var privateKeyHexCleanStr string;
    if (strings.HasPrefix(pPrivateKeyHexStr, "0x")) {
        privateKeyHexCleanStr = strings.TrimPrefix(pPrivateKeyHexStr, "0x")
    } else {
        privateKeyHexCleanStr = pPrivateKeyHexStr
    }

    // parse private_key hex
	privateKey, err := crypto.HexToECDSA(privateKeyHexCleanStr)
    if err != nil {
        return "", err
    }

	dataToSignBytesLst := []byte(pDataToSignStr)
	dataToSignHash := crypto.Keccak256Hash(dataToSignBytesLst)

    //---------------------------
    // SIGN
	signature, err := crypto.Sign(dataToSignHash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	signatureHexStr := hexutil.Encode(signature)

    //---------------------------
	return signatureHexStr, err
}