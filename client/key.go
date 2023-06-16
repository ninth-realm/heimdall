package client

import "github.com/ninth-realm/heimdall/crypto"

const keyPrefixLength = 6

const keySuffixLength = 32

func generateAPIKey() (string, string, error) {
	prefix, err := crypto.GenerateRandHexString(keyPrefixLength)
	if err != nil {
		return "", "", err
	}

	suffix, err := crypto.GenerateRandBase64String(keySuffixLength)
	if err != nil {
		return "", "", err
	}

	return prefix, suffix, nil
}
