package gf_identity_core

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func RefreshTokenGenerate(pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// 32 bytes = 256 bits of entropy
	// will result in a 43-character base64url encoded string
	tokenBytes := make([]byte, 32)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to generate cryptographically secure random bytes for refresh token",
			"crypto_random_generation_error",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return "", gfErr
	}

	// encode to base64url (URL-safe, no padding)
	tokenStr := base64.RawURLEncoding.EncodeToString(tokenBytes)

	return tokenStr, nil
}
