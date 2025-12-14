package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env string) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if strings.EqualFold(env, "dev") || strings.EqualFold(env, "development") {
		level = zapcore.DebugLevel
	}
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encCfg)
	output := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, output, level)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)), nil
}
