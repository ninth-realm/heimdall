package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

func GenerateRandHexString(l int) (string, error) {
	buf := make([]byte, l)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf)[:6], nil
}

func GenerateRandBase64String(l int) (string, error) {
	buf := make([]byte, l)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}
