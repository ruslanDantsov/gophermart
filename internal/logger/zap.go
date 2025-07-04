package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var Log = zap.NewNop()

func Initialized(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	withCustomTimeLayout("2006-01-02 15:04:05")(&cfg)
	WithServiceName("GopherMart app")(&cfg)

	configuredLogger, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = configuredLogger
	return nil
}

type LoggerOption func(*zap.Config)

func withCustomTimeLayout(layout string) LoggerOption {
	return func(cfg *zap.Config) {
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(layout))
		}
	}
}

func WithServiceName(name string) LoggerOption {
	return func(cfg *zap.Config) {
		cfg.InitialFields = map[string]interface{}{
			"service": name,
		}
	}
}
