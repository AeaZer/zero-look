package zapx

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const callerSkipOffset = 3

type ZapWriter struct {
	logger *zap.Logger
}

type OptConf struct {
	ZapOpts  []zap.Option
	ZapCores []zapcore.Core
}

func getOption(opts []OptConf) OptConf {
	optConf := OptConf{
		ZapOpts:  make([]zap.Option, 0),
		ZapCores: make([]zapcore.Core, 0),
	}
	for _, opt := range opts {
		optConf = opt
	}
	return optConf
}

func NewZapWriter(optConfigs ...OptConf) (logx.Writer, error) {
	optConf := getOption(optConfigs)
	zapOpts := append(optConf.ZapOpts, zap.AddCallerSkip(callerSkipOffset))
	logger, err := zap.NewProduction(zapOpts...)
	if err != nil {
		return nil, err
	}
	if len(optConf.ZapCores) != 0 {
		logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			cors := append([]zapcore.Core{core}, optConf.ZapCores...)
			return zapcore.NewTee(cors...)
		}))
	}
	return &ZapWriter{
		logger: logger,
	}, nil
}

func (w *ZapWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *ZapWriter) Close() error {
	return w.logger.Sync()
}

func (w *ZapWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.Debug(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Severe(v interface{}) {
	w.logger.Fatal(fmt.Sprint(v))
}

func (w *ZapWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v), zap.Stack("stack"))
}

func (w *ZapWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func toZapFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
