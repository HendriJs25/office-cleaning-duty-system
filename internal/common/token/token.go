package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func Generate() string {
	return rand.Text()
}

func Hash(rawToken string) string {
	digest := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(digest[:])
}
