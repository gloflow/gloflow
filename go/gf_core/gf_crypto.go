
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
	"fmt"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
)

//-------------------------------------------------

func CryptoGeneratePrivKeyAsPEM() string {

	privKey, pubKey := CryptoGenerateKeyPair()

	fmt.Println(pubKey)

	privKeyPEMstr := CryptoConvertPrivKeyToPEM(privKey)

	return privKeyPEMstr
}

func CryptoGenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
    privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
    return privkey, &privkey.PublicKey
}

//-------------------------------------------------
// PEM
//-------------------------------------------------

// parse a private key from a PEM string
func CryptoParsePrivKeyFromPEM(pPrivKeyPEMstr string,
	pRuntimeSys *RuntimeSys) (*rsa.PrivateKey, *GFerror) {

	// find the next PEM formatted block (certificate, private key etc) in the input
    block, _ := pem.Decode([]byte(pPrivKeyPEMstr))
    if block == nil {
		gfErr := ErrorCreate("failed to parse a PEM block from a string, containing a private key",
			"crypto_pem_decode",
			map[string]interface{}{},
			nil, "gf_core", pRuntimeSys)
		return nil, gfErr
    }

    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        gfErr := ErrorCreate("failed to parse a x509 formated private-key from a PEM block",
			"crypto_x509_parse",
			map[string]interface{}{},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
    }

    return priv, nil
}

//-------------------------------------------------

func CryptoConvertPrivKeyToPEM(pPrivKey *rsa.PrivateKey) string {

	privKeyBytesLst := x509.MarshalPKCS1PrivateKey(pPrivKey)
	privKeyPEM := pem.EncodeToMemory(
			&pem.Block{
					Type:  "RSA PRIVATE KEY",
					Bytes: privKeyBytesLst,
			},
	)

	privKeyPEMstr := string(privKeyPEM)
	return privKeyPEMstr
}

//-------------------------------------------------

