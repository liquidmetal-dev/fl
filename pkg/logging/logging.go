package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Configure(level string) error {
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("parsing log level %s: %w", level, err)
	}

	cfg := zap.NewProductionConfig()

	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	cfg.EncoderConfig.TimeKey = ""
	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level.SetLevel(parsedLevel)

	manager, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("building zap logger: %w", err)
	}
	zap.ReplaceGlobals(manager)

	return nil
}
