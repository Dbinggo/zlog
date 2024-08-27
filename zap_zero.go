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
type zeroWriter struct {
	logger *Zlogger
}

const zeroCallerSkip = 0

var _ logx.Writer = (*zeroWriter)(nil) // 接口实现检查
func NewZeroWriter(logger *Zlogger) logx.Writer {
	return &zeroWriter{logger: logger}
}

func (l *zeroWriter) Alert(v interface{}) {
	l.logger.WithCallerSkip(zeroCallerSkip).errorf(fmt.Sprint(v))
}

func (l *zeroWriter) Close() error {
	return l.logger.Sync()
}

func (l *zeroWriter) Debug(v interface{}, fields ...logx.LogField) {

	l.logger.WithCallerSkip(zeroCallerSkip).debugField(fmt.Sprint(v), toZapFields(fields...)...)

}

func (l *zeroWriter) Error(v interface{}, fields ...logx.LogField) {

	l.logger.WithCallerSkip(zeroCallerSkip).errorField(fmt.Sprint(v), toZapFields(fields...)...)

}

func (l *zeroWriter) Info(v interface{}, fields ...logx.LogField) {

	l.logger.WithCallerSkip(zeroCallerSkip).infoField(fmt.Sprint(v), toZapFields(fields...)...)

}

func (l *zeroWriter) Severe(v interface{}) {
	l.logger.WithCallerSkip(zeroCallerSkip).errorf(fmt.Sprint(v))
}

func (l *zeroWriter) Slow(v interface{}, fields ...logx.LogField) {

	l.logger.WithCallerSkip(zeroCallerSkip).warnField(fmt.Sprint(v), toZapFields(fields...)...)

}

func (l *zeroWriter) Stack(v interface{}) {
	if l.logger.FormatJson() {
		l.logger.WithCallerSkip(zeroCallerSkip).errorf(fmt.Sprint(v), zap.Stack("stack"))
	} else {
		l.logger.WithCallerSkip(zeroCallerSkip).errorf(fmt.Sprint(v))
	}
}

func (l *zeroWriter) Stat(v interface{}, fields ...logx.LogField) {

	l.logger.WithCallerSkip(zeroCallerSkip).infoField(fmt.Sprint(v), toZapFields(fields...)...)

}

func toZapFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
