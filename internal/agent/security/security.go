package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func GetPublicKey(path string) (*rsa.PublicKey, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error read file: %v", err)
	}

	block, _ := pem.Decode(file)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("error decode PEM block")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parse RSA public key: %v", err)
	}

	return publicKey, nil
}

func EncryptBody(body []byte, pubKey *rsa.PublicKey) (string, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, body)
	if err != nil {
		return "", err
	}

	result := base64.StdEncoding.EncodeToString(encryptedData)

	return result, nil
}
