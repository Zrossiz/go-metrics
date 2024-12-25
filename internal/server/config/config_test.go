package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigFromEnv_Success(t *testing.T) {
	expectedDB := "db_dsn"
	expectedAddress := "address"
	expectedStoreInterval := "19s"
	expectedFileStoragePath := "file storage path"
	expectedLogLevel := "DEBUG"
	expectedCryptoKey := "crypto_key"
	expectedHashKey := "hash_key"

	os.Setenv("DATABASE_DSN", expectedDB)
	os.Setenv("ADDRESS", expectedAddress)
	os.Setenv("STORE_INTERVAL", expectedStoreInterval)
	os.Setenv("RESTORE", "true")
	os.Setenv("FILE_STORAGE_PATH", expectedFileStoragePath)
	os.Setenv("LOG_LEVEL", expectedLogLevel)
	os.Setenv("CRYPTO_KEY", expectedCryptoKey)
	os.Setenv("KEY", expectedHashKey)

	cfg, err := GetConfig()
	if err != nil {
		t.Errorf("err get config: %v", err)
	}

	assert.NotNil(t, cfg)
	assert.Equal(t, expectedDB, cfg.DBDSN)
	assert.Equal(t, expectedAddress, cfg.ServerAddress)
	assert.Equal(t, expectedStoreInterval, cfg.StoreInterval)
	assert.Equal(t, true, cfg.Restore)
	assert.Equal(t, expectedFileStoragePath, cfg.FileStoragePath)
	assert.Equal(t, expectedLogLevel, cfg.LogLevel)
	assert.Equal(t, expectedCryptoKey, cfg.PrivateKeyPath)
	assert.Equal(t, expectedHashKey, cfg.HashKey)
}
