package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() error {
	config := zap.NewDevelopmentEncoderConfig()

	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	Log = zap.New(core, zap.AddCaller())

	return nil
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
