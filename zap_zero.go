package zlog

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

// 实现zero 的 logx的Writer接口
//
//	type Writer interface {
//		Alert(v any)
//		Close() error
//		Debug(v any, fields ...LogField)
//		Error(v any, fields ...LogField)
//		Info(v any, fields ...LogField)
//		Severe(v any)
//		Slow(v any, fields ...LogField)
//		Stack(v any)
//		Stat(v any, fields ...LogField)
//	}
type zeroLogger struct {
	logger     *zlogger
	formatJson bool
}

var _ logx.Writer = (*zeroLogger)(nil) // 接口实现检查
func NewZeroLogger(logger *zlogger, formatJson bool) logx.Writer {
	return &zeroLogger{logger: logger, formatJson: formatJson}
}

func (l *zeroLogger) Alert(v interface{}) {
	l.logger.WithCallerSkip(6).errorf(fmt.Sprint(v))
}

func (l *zeroLogger) Close() error {
	return l.logger.Sync()
}

func (l *zeroLogger) Debug(v interface{}, fields ...logx.LogField) {
	if l.formatJson {
		l.logger.WithCallerSkip(6).debugField(fmt.Sprint(v), toZapFields(fields...)...)
	} else {
		l.logger.WithCallerSkip(6).debugf(fmt.Sprint(v))
	}
}

func (l *zeroLogger) Error(v interface{}, fields ...logx.LogField) {
	if l.formatJson {
		l.logger.WithCallerSkip(6).errorField(fmt.Sprint(v), toZapFields(fields...)...)
	} else {
		l.logger.errorf(fmt.Sprint(v))
	}
}

func (l *zeroLogger) Info(v interface{}, fields ...logx.LogField) {
	if l.formatJson {
		l.logger.WithCallerSkip(6).infoField(fmt.Sprint(v), toZapFields(fields...)...)
	} else {
		l.logger.WithCallerSkip(6).infof(fmt.Sprint(v))
	}
}

func (l *zeroLogger) Severe(v interface{}) {
	l.logger.WithCallerSkip(6).errorf(fmt.Sprint(v))
}

func (l *zeroLogger) Slow(v interface{}, fields ...logx.LogField) {
	if l.formatJson {
		l.logger.WithCallerSkip(6).warnField(fmt.Sprint(v), toZapFields(fields...)...)
	} else {
		l.logger.WithCallerSkip(6).warnf(fmt.Sprint(v))
	}
}

func (l *zeroLogger) Stack(v interface{}) {
	if l.formatJson {
		l.logger.WithCallerSkip(6).errorf(fmt.Sprint(v), zap.Stack("stack"))
	} else {
		l.logger.WithCallerSkip(6).errorf(fmt.Sprint(v))
	}
}

func (l *zeroLogger) Stat(v interface{}, fields ...logx.LogField) {
	if l.formatJson {
		l.logger.WithCallerSkip(6).infoField(fmt.Sprint(v), toZapFields(fields...)...)
	} else {
		l.logger.WithCallerSkip(6).infof(fmt.Sprint(v))

	}
}

func toZapFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
