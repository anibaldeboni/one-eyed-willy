package logger

import (
	"github.com/one-eyed-willy/pkg/envs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.SugaredLogger {
	var cfg zap.Config

	if envs.IsProd() {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	cfg.EncoderConfig.TimeKey = "@timestamp"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "log_level"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.StacktraceKey = "trace"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Encoding = "json"
	cfg.OutputPaths = []string{"stdout"}
	l, _ := cfg.Build()

	return l.Sugar()
}
