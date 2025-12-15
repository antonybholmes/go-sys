package sys

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

func MakeEd25519Key(key string) (ed25519.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}

	// Check the length of the decoded bytes
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("decoded public key has invalid length: expected %d, got %d", ed25519.PublicKeySize, len(publicKeyBytes))
	}

	return ed25519.PublicKey(publicKeyBytes), nil
}
