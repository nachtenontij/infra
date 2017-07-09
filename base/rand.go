package base

import (
	"crypto/rand"
	"encoding/hex"
)

// Generate a secret hexadecimal key with the given number of bytes of
// entropy (thus its length will be twice that).
func GenerateHexSecret(nBytes int) string {
	var rawKey []byte = make([]byte, nBytes)
	_, err := rand.Read(rawKey)
	if err != nil {
		panic("GenerateHexSecret: no random")
	}
	return hex.EncodeToString(rawKey[:])
}
