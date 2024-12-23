package security

import (
	"context"
	"crypto/rsa"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRSA struct {
	mock.Mock
}

const privateKeyPemExample = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA071jqS/oQG4vsDG2voa9I7doaYM3cdLTYQold4jBwmSl9eFf
/PCPSYK7Bo+ztRNR3y6LOnGGZHL/85nUricyUR92KXZorEI7DEr1j3V4IUuJ1M9f
uNKlDu4tH7mLWc7F1yk15j6t8/BSJGaXYjYdTfJJKpuhNAFkIwKUUORjFN+s1SCV
XjlXS2VkbYXquizep9aUcGhHWWB0zEzyH+KudY/+PuHqAQ1oCXpsi+ZVsJSF9yxn
L3p33BiwwnJMgt+MmE7Pzb1mKZlL/ugJUZ8KQVG2AAXYKvgnW3OhX82vA2lDvZ7k
gZ50Not7TY6NSPCz7nofrDlGFbUC1acWj49y9QIDAQABAoIBAA1PoqxILrmeM7FH
7s0be1h7dzzq2tl0/4TiNmjFWCj4WtaSTI6CMP/WEBvhfNKtXEBDlM0fxesar6yI
xy4LmMYwzCTfJMVvhHbJX5adM+gj745JWyMrWuUNZBjSeUg0D4vvM0w+NIFZYlCX
gnzSGhWEXcUn84Jpc/ofd2N+eBwftklYDzBSczVQ4bwRbknS4A1R83uB2Elj4fqw
4uYwyfa7dmZm3DlF7HABrO+SmxCG5hDUoc/NJn5CbbLekn9G5V36tZty5ia9iXdM
iVsxaW+d7O0pQ4SW68tkJQmhr9AOc1NeLXjqALaTCzYKDrXDUtDR8YDT92JvpI+A
tmUfvQECgYEA13F3k5mvgnRdudZp5J8FAB2LL7MaXFhTas2zpHqJrDrX5VCxpLrm
Eif8nb8LwA/AuuVizsVY3s/IFSRfgBgEl+/NqVSZ3+4BX3YFl9GESnyOMuMawqjE
5THCEd41KBntOMyK7O95uyuLABUxt5vFDi4zoE9Nkaqnmmo0T0jAxi0CgYEA+5ly
ypE8kEGmS9YK+flIEDyf2TdBwIz4cvf1agmHua8yjx8Sx4CtMbjY9vyaS91ZmOFn
AgTKN2g4ATDyJ4SeBcDbZbe9haJC6MEt9b4Rg46e8UZJ/Epf+PKFSwvdiEuEx7jQ
T5ydafhgTul+1q6JEF4Nr2iCzl2xAxEl3Enn5OkCgYEAoOLmRj5Vt9kIiRgamhU6
mbx2TZe1jtKS8MZOafzsRMbopSHel0LPPy3HU1HxB2t8JNXaNMlhNXr7UvaHrtPA
0mnNLq+z/WrycYRkZtyaqzlaw5ufR1DbQMEoyUkkbx71bR4qfQfU4zaAJf6t0wyr
WoycFxJBvg8v/HtlNvQAqb0CgYEAxLhkIrKQchKCngULrAwXJmrgaQxlYtJWaD4s
Ku6sqqirlXAsVMTtplTrf6JeWjcGGR0UV2W7XrskHvpQPEna7JCwesXBb71BJ4/0
CZLFSuG2sNvOeW8FvzaQte7fFfRGK4r7hWPlSLglRU4YGG97R8riVGYY8JYdE1LT
EXzPzhECgYAqwal3YSq0WdIauSNmiIBnaZzMfuByRrEW1opoklbVJexwu1QukRqy
2wl4Wdg1ZaB4Z22OAaIRV8lU+8JbDtTWNWHhxiUXRHtR4drIakxFCwHpt6SRP/8O
Y3gi47ujhkZ9RrADtHJmOJtoL3uuVgCmkulphl6ILbISFIY+fhkpmQ==
-----END RSA PRIVATE KEY-----`

func (m *MockRSA) DecryptPKCS1v15(rand io.Reader, priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	args := m.Called(rand, priv, ciphertext)
	return args.Get(0).([]byte), args.Error(1)
}

func TestNewDecryptedContext(t *testing.T) {
	message := []byte("test message")
	ctx := NewDecryptedContext(context.Background(), message)

	decryptedMessage := DecryptedFromContext(ctx)

	assert.Equal(t, message, decryptedMessage)
}

func TestDecryptedFromContext(t *testing.T) {
	message := []byte("test message")
	ctx := NewDecryptedContext(context.Background(), message)

	decryptedMessage := DecryptedFromContext(ctx)

	assert.Equal(t, message, decryptedMessage)
}

func TestGetPrivateKey_Success(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "private_key.pem")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(privateKeyPemExample))
	assert.Nil(t, err)

	privateKey, err := GetPrivateKey(tmpFile.Name())

	assert.Nil(t, err)
	assert.NotNil(t, privateKey)
}

func TestGetPrivateKey_Failure(t *testing.T) {
	_, err := GetPrivateKey("invalid_key_path")

	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

func TestCheckCryptoBody_Failure(t *testing.T) {
	encryptedBody := "invalid_encrypted_string"

	decryptedMessage, err := CheckCryptoBody([]byte(encryptedBody))

	assert.NotNil(t, err)
	assert.Empty(t, decryptedMessage)
}
