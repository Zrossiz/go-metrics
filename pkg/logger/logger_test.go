package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel zapcore.Level
		expectedError bool
	}{
		{"debug", "debug", zap.DebugLevel, false},
		{"info", "info", zap.InfoLevel, false},
		{"warn", "warn", zap.WarnLevel, false},
		{"error", "error", zap.ErrorLevel, false},
		{"default", "invalid", zap.InfoLevel, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, err := New(tt.level)
			if tt.expectedError {
				assert.NoError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLevel, log.AtomicLevel.Level())
			}
		})
	}
}
