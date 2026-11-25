package auth

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadPrivateKey(path string) (crypto.Signer, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Iterate over PEM blocks (handles files containing both EC PARAMETERS and EC PRIVATE KEY)
	for len(pemBytes) > 0 {
		var block *pem.Block
		block, pemBytes = pem.Decode(pemBytes)
		if block == nil {
			break
		}

		switch block.Type {
		case "EC PRIVATE KEY":
			return x509.ParseECPrivateKey(block.Bytes)
		case "RSA PRIVATE KEY":
			return x509.ParsePKCS1PrivateKey(block.Bytes)
		case "PRIVATE KEY": // PKCS#8 format (works for both RSA and ECDSA)
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			signer, ok := key.(crypto.Signer)
			if !ok {
				return nil, fmt.Errorf("parsed PKCS#8 key does not implement crypto.Signer")
			}
			return signer, nil
		}
	}
	return nil, fmt.Errorf("no valid private key block found in %s", path)
}
