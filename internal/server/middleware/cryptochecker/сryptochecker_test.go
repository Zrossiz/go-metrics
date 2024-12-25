package cryptochecker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCryptoChecker для подмены реальной логики расшифровки
type MockCryptoChecker struct {
	mock.Mock
}

func (m *MockCryptoChecker) Decrypt(encryptedData []byte) ([]byte, error) {
	args := m.Called(encryptedData)
	return args.Get(0).([]byte), args.Error(1)
}

// Helper function to test successful decryption
func testSuccessfulDecryption(t *testing.T, mockChecker *MockCryptoChecker) {
	// Мокируем логику расшифровки
	mockChecker.On("Decrypt", []byte("encrypted data")).Return([]byte("decrypted data"), nil)

	// Создаем следующий обработчик, который будет проверять контекст
	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Проверяем, что контекст содержит расшифрованные данные
		decryptedMessage := r.Context().Value(security.DecryptedMessageKey) // Используем правильный ключ
		assert.NotNil(t, decryptedMessage)
		assert.Equal(t, "decrypted data", string(decryptedMessage.([]byte)))
		rw.WriteHeader(http.StatusOK)
	})

	// Создаем middleware с мокированным checkCryptoBody
	middleware := DecryptMiddleware(nextHandler, mockChecker.Decrypt)

	// Создаем запрос с зашифрованными данными в теле
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("encrypted data")))
	rr := httptest.NewRecorder()

	// Вызываем middleware
	middleware.ServeHTTP(rr, req)

	// Проверяем статус
	assert.Equal(t, http.StatusOK, rr.Code)
	mockChecker.AssertExpectations(t)
}

func TestDecryptMiddleware_SuccessfulDecryption(t *testing.T) {
	mockChecker := new(MockCryptoChecker)
	testSuccessfulDecryption(t, mockChecker)
}
