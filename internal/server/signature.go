package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func ValidateGitHubSignature(secret string, body []byte, sigHeader string) bool {
	if secret == "" {
		return false // 建议不允许空 secret
	}
	if sigHeader == "" {
		return false
	}
	const prefix = "sha256="
	if !strings.HasPrefix(sigHeader, prefix) {
		return false
	}
	givenSig := strings.TrimPrefix(sigHeader, prefix)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := mac.Sum(nil)
	expectedHex := hex.EncodeToString(expected)
	return hmac.Equal([]byte(givenSig), []byte(expectedHex))
}
