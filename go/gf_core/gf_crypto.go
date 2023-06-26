
/*
MIT License

Copyright (c) 2022 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_core

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
)

//-------------------------------------------------

// generate new RSA private/public key-pair (4096 bit) and encoded it into
// the PEM format and return as two separate strings.
func CryptoGenerateKeysAsPEM() (string, string) {

	pubKey, privKey := CryptoGenerateKeys()
	pubKeyPEMstr, privKeyPEMstr := CryptoConvertKeysToPEM(privKey, pubKey)

	return pubKeyPEMstr, privKeyPEMstr 
}

//-------------------------------------------------

// generate RSA private/public keys (4096 bit)
func CryptoGenerateKeys() (*rsa.PublicKey, *rsa.PrivateKey) {
    privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
    return &privkey.PublicKey, privkey 
}

//-------------------------------------------------
// PEM
//-------------------------------------------------

// parse a private key from a PEM string
func CryptoParseKeysFromPEM(pPublicKeyPEMstr string,
	pPrivateKeyPEMstr string,
	pRuntimeSys       *RuntimeSys) (*rsa.PublicKey, *rsa.PrivateKey, *GFerror) {
	
	//------------------------
	// find the next PEM formatted block (certificate, private key etc) in the input
    pubBlock, _ := pem.Decode([]byte(pPublicKeyPEMstr))
    if pubBlock == nil {
		gfErr := ErrorCreate("failed to parse a PEM block from a string, containing a public key",
			"crypto_pem_decode",
			map[string]interface{}{},
			nil, "gf_core", pRuntimeSys)
		return nil, nil, gfErr
    }

    publicKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
    if err != nil {
        gfErr := ErrorCreate("failed to parse a x509 formated public-key from a PEM block",
			"crypto_x509_parse",
			map[string]interface{}{},
			err, "gf_core", pRuntimeSys)
		return nil, nil, gfErr
    }

	//------------------------
	// find the next PEM formatted block (certificate, private key etc) in the input
    privBlock, _ := pem.Decode([]byte(pPrivateKeyPEMstr))
    if privBlock == nil {
		gfErr := ErrorCreate("failed to parse a PEM block from a string, containing a private key",
			"crypto_pem_decode",
			map[string]interface{}{},
			nil, "gf_core", pRuntimeSys)
		return nil, nil, gfErr
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
    if err != nil {
        gfErr := ErrorCreate("failed to parse a x509 formated private-key from a PEM block",
			"crypto_x509_parse",
			map[string]interface{}{},
			err, "gf_core", pRuntimeSys)
		return nil, nil, gfErr
    }

	//------------------------

    return publicKey, privateKey, nil
}

//-------------------------------------------------

func CryptoConvertKeysToPEM(pPrivateKey *rsa.PrivateKey,
	pPublicKey *rsa.PublicKey) (string, string) {

	//------------------------
	// PUBLIC_KEY
	
	pubKeyPEMstr := CryptoConvertPubKeyToPEM(pPublicKey)
	//------------------------
	// PRIVATE_KEY
	privKeyBytesLst := x509.MarshalPKCS1PrivateKey(pPrivateKey)
	privKeyPEM := pem.EncodeToMemory(
			&pem.Block{
					Type:  "RSA PRIVATE KEY",
					Bytes: privKeyBytesLst,
			},
	)

	privKeyPEMstr := string(privKeyPEM)

	//------------------------
	return pubKeyPEMstr, privKeyPEMstr
}

//-------------------------------------------------

func CryptoConvertPubKeyToPEM(pPublicKey *rsa.PublicKey) string {
	pubKeyBytesLst := x509.MarshalPKCS1PublicKey(pPublicKey)
	pubKeyPEM := pem.EncodeToMemory(
			&pem.Block{
					Type:  "RSA PUBLIC KEY",
					Bytes: pubKeyBytesLst,
			},
	)

	pubKeyPEMstr := string(pubKeyPEM)
	return pubKeyPEMstr
}