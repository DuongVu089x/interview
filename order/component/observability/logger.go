package observability

import (
	"github.com/DuongVu089x/interview/order/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger creates a new logger with rotation
func NewLogger(cfg *config.ObservabilityConfig) (*zap.Logger, error) {
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Logging.OutputPath,
		MaxSize:    cfg.Logging.MaxSize,
		MaxBackups: cfg.Logging.MaxBackups,
		MaxAge:     cfg.Logging.MaxAge,
		Compress:   cfg.Logging.Compress,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var level zapcore.Level
	err := level.UnmarshalText([]byte(cfg.Logging.Level))
	if err != nil {
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		level,
	)

	logger := zap.New(core)
	return logger, nil
}

// GetLogger returns a named logger
func GetLogger(name string) *zap.Logger {
	return zap.L().Named(name)
}
