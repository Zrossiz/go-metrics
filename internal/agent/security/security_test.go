package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	publicKey := &privateKey.PublicKey

	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	tempFile, err := os.CreateTemp("", "public_key_*.pem")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(publicKeyPEM)
	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	retrievedPublicKey, err := GetPublicKey(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, publicKey.N, retrievedPublicKey.N)
	assert.Equal(t, publicKey.E, retrievedPublicKey.E)
}

func TestEncryptBody(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	publicKey := &privateKey.PublicKey

	originalBody := []byte("Test data for encryption")

	encryptedBody, err := EncryptBody(originalBody, publicKey)
	assert.NoError(t, err)

	decodedBody, err := base64.StdEncoding.DecodeString(encryptedBody)
	assert.NoError(t, err)

	decryptedBody, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedBody)
	assert.NoError(t, err)
	assert.Equal(t, originalBody, decryptedBody)
}

func TestGetPublicKey_InvalidPath(t *testing.T) {
	_, err := GetPublicKey("nonexistent_file.pem")
	assert.Error(t, err)
}

func TestGetPublicKey_InvalidPEM(t *testing.T) {
	tempFile, err := os.CreateTemp("", "invalid_key_*.pem")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("invalid pem data"))
	assert.NoError(t, err)

	err = tempFile.Close()
	assert.NoError(t, err)

	_, err = GetPublicKey(tempFile.Name())
	assert.Error(t, err)
}

func TestEncryptBody_InvalidKey(t *testing.T) {
	invalidKey := &rsa.PublicKey{}
	_, err := EncryptBody([]byte("test data"), invalidKey)
	assert.Error(t, err)
}
