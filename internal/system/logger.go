package system

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func getLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = Config.Application.Log.Outputs
	cfg.Level = Config.Application.Log.Level

	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}
