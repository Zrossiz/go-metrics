package security

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/Zrossiz/go-metrics/internal/server/config"
)

type contextKey string

const (
	DecryptedMessageKey = contextKey("decryptedMessage")
)

func NewDecryptedContext(ctx context.Context, message []byte) context.Context {
	return context.WithValue(ctx, DecryptedMessageKey, message)
}

func DecryptedFromContext(ctx context.Context) []byte {
	if message, ok := ctx.Value(DecryptedMessageKey).([]byte); ok {
		return message
	}
	return nil
}

func GetPrivateKey(path string) (*rsa.PrivateKey, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(file)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func CheckCryptoBody(encryptedBody []byte) ([]byte, error) {
	encryptedMessage, err := base64.StdEncoding.DecodeString(string(encryptedBody))
	if err != nil {
		return make([]byte, 0), err
	}
	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, config.AppConfig.PrivateKey, encryptedMessage)
	if err != nil {
		return make([]byte, 0), err
	}

	return decryptedMessage, nil
}
