package zlog

import (
	"context"
	"fmt"
	"github.com/dbinggo/zlog/utils"
	"github.com/zeromicro/go-zero/core/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
)

const (
	logger_formate_json = "json"
	loggerKey           = "zapLogger"
	loggerFieldKey      = "field"
	loggerDebugKey      = "debug"
	loggerCallerKey     = "caller"
	loggerSpanKey       = "span"
	loggerTraceKey      = "trace"
	loggerCallerSkipKey = "callerSkip"

	//loggerPrefixKey = "prefix"
	// 分隔符

)

type zlogger struct {
	prefix     string
	format     string
	ctx        context.Context
	callerSkip int
}

// 全局logger 没有在上下问中找到logger会用这个
var globalZapLogger *zap.Logger = new(zap.Logger)
var formatJson = "json"

// SetLogger 注册logger
func SetLogger(zapLogger *zap.Logger, json bool) {
	globalZapLogger = zapLogger
	if !json {
		formatJson = "plain"
	}
}
func NewLogger() *zlogger {
	return &zlogger{
		format:     formatJson,
		ctx:        context.Background(),
		callerSkip: 0,
	}
}

// withContext 从上下文中拿到logger
func withContext(ctx context.Context) zap.Logger {
	if ctx == nil {
		return *globalZapLogger
	}
	logger := ctx.Value(loggerKey)
	if logger == nil {
		logger = *globalZapLogger
		context.WithValue(ctx, loggerKey, &logger)
	}
	return logger.(zap.Logger)
}

// ###########################################
// 可以直接通过 包名来使用
func DebugfCtx(ctx context.Context, format string, v ...any) {
	l := &zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(4).debugf(format, v...)
}
func Debugf(format string, v ...any) {
	l := &zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(4).debugf(format, v...)

}
func InfofCtx(ctx context.Context, format string, v ...any) {
	l := &zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(4).infof(format, v...)
}
func Infof(format string, v ...any) {
	l := &zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(4).infof(format, v...)

}
func WarnfCtx(ctx context.Context, format string, v ...any) {
	l := &zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(4).warnf(format, v...)
}
func Warnf(format string, v ...any) {
	l := &zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(4).warnf(format, v...)

}
func ErrorfCtx(ctx context.Context, format string, v ...any) {
	l := &zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(4).errorf(format, v...)
}
func Errorf(format string, v ...any) {
	l := &zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(4).errorf(format, v...)

}

// ###########################################
// 通过对象来使用
func (l *zlogger) Debugf(format string, v ...any) {
	l.WithCallerSkip(4).debugf(format, v)
}
func (l *zlogger) Infof(format string, v ...any) {
	l.WithCallerSkip(4).infof(format, v)
}
func (l *zlogger) Warnf(format string, v ...any) {
	l.WithCallerSkip(4).warnf(format, v)
}
func (l *zlogger) Errorf(format string, v ...any) {
	l.WithCallerSkip(4).errorf(format, v)
}

// ###########################################
// 内部方法
func (l *zlogger) debugf(format string, v ...any) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger)
	logger.Debug(fmt.Sprintf(exString+format, v...))
}
func (l *zlogger) infof(format string, v ...any) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger)
	logger.Info(fmt.Sprintf(exString+format, v...))
}
func (l *zlogger) warnf(format string, v ...any) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger)
	logger.Warn(fmt.Sprintf(exString+format, v...))
}
func (l *zlogger) errorf(format string, v ...any) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger)
	logger.Error(fmt.Sprintf(exString+format, v...))
}

// ###########################################

// ###########################################
// 通用功能方法
func (l *zlogger) getCallerSkip() int {
	if l.callerSkip == 0 {
		return 3
	}
	return l.callerSkip
}

func (l *zlogger) formatJson() bool {
	return l.format == "json"
}

func (l *zlogger) addCaller(_logger *zap.Logger) (zap.Logger, string) {
	format := "%s:%d"
	_, file, line, _ := runtime.Caller(l.getCallerSkip())
	_v := make([]interface{}, 0)
	file = file[len(utils.GetRootPath(""))+1:]
	_v = append(_v, file, line)
	if l.formatJson() {
		_logger = _logger.With(zap.String(loggerCallerKey, fmt.Sprintf(format, file, line)))
		return *_logger, ""
	}
	return *_logger, fmt.Sprintf(format+"\t", _v...)
}

func (l *zlogger) addTrace(ctx context.Context, _logger *zap.Logger) (zap.Logger, string) {
	traceId := trace.TraceIDFromContext(ctx)
	if traceId == "" {
		return *_logger, ""
	}
	if l.formatJson() {
		_logger = _logger.With(zap.String(loggerTraceKey, traceId))
		return *_logger, ""
	}
	format := "%v\t"
	return *_logger, fmt.Sprintf(format, traceId)
}
func (l *zlogger) addSpan(ctx context.Context, _logger *zap.Logger) (zap.Logger, string) {

	spanId := trace.SpanIDFromContext(ctx)
	if spanId == "" {
		return *_logger, ""
	}
	if l.formatJson() {
		_logger = _logger.With(zap.String(loggerSpanKey, spanId))
		return *_logger, ""
	}
	format := "%v\t"
	return *_logger, fmt.Sprintf(format, spanId)
}
func (l *zlogger) addExField(ctx context.Context, _logger *zap.Logger) (zap.Logger, string) {
	if exField := ctx.Value(loggerFieldKey); exField != nil {
		if l.formatJson() {
			_logger = _logger.With(exField.([]zapcore.Field)...)
			return *_logger, ""
		} else {
			format := "%v\t"
			ret := ""
			for _, field := range exField.([]zapcore.Field) {
				ret += fmt.Sprintf(format, field.String)
			}
			return *_logger, ret
		}

	}
	return *_logger, ""
}
func AddFiled(ctx context.Context, fields ...zapcore.Field) context.Context {
	if exField, ok := ctx.Value(loggerFieldKey).([]zapcore.Field); ok {
		exField = append(exField, fields...)
		return context.WithValue(ctx, loggerFieldKey, exField)
	}
	return context.WithValue(ctx, loggerFieldKey, fields)

}

//func addPrefix(ctx context.Context, _logger *zap.Logger) (zap.Logger, string) {
//	return *_logger, prefix
//}

// 构建 field
func (l *zlogger) buildField(logger *zap.Logger) (zap.Logger, string) {
	// 加入 caller
	var (
		caller  string
		traceId string
		spanId  string
		field   string
	)
	*logger, caller = l.addCaller(logger)
	*logger, traceId = l.addTrace(l.ctx, logger)
	*logger, spanId = l.addSpan(l.ctx, logger)
	*logger, field = l.addExField(l.ctx, logger)
	newLine := "\n"
	if l.formatJson() {
		newLine = ""

	}
	return *logger, caller + traceId + spanId + field + newLine
}

// WithCallerSkip  携带跳过的层数
func (l *zlogger) WithCallerSkip(skip int) *zlogger {
	if skip <= 0 {
		return l
	}
	return &zlogger{
		ctx:        l.ctx,
		callerSkip: skip,
		format:     l.format,
		prefix:     l.prefix,
	}
}
func (l *zlogger) Sync() error {
	logger := withContext(l.ctx)
	return logger.Sync()
}

// WithContext 携带上下文
func (l *zlogger) WithContext(ctx context.Context) *zlogger {
	return &zlogger{
		ctx:        ctx,
		callerSkip: l.callerSkip,
		format:     l.format,
		prefix:     l.prefix,
	}
}

// 供外部方法使用
func (l *zlogger) debugField(msg string, fields ...zap.Field) {
	logger := withContext(l.ctx)
	logger.Debug(msg, fields...)
}
func (l *zlogger) infoField(msg string, fields ...zap.Field) {
	logger := withContext(l.ctx)
	logger.Info(msg, fields...)
}
func (l *zlogger) warnField(msg string, fields ...zap.Field) {
	logger := withContext(l.ctx)
	logger.Warn(msg, fields...)
}
func (l *zlogger) errorField(msg string, fields ...zap.Field) {
	logger := withContext(l.ctx)
	logger.Error(msg, fields...)
}
