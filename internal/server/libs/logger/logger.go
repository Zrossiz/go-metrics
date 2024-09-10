package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	ZapLogger   *zap.Logger
	AtomicLevel zap.AtomicLevel
}

func New(level string) (*Logger, error) {
	var zapLevel zapcore.Level

	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	atomicLevel := zap.NewAtomicLevelAt(zapLevel)

	config := zap.Config{
		Level:            atomicLevel,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{ZapLogger: zapLogger, AtomicLevel: atomicLevel}, nil
}
