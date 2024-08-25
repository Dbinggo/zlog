package zlog

import (
	"context"
	"fmt"
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

// 这是一个 最底层的对象 其余接口实现的对象都调用此对象
/*
	-----------------   	-----------------
	｜  zeroWriter  ｜		｜  gormLogger  ｜
	-----------------		-----------------
			｜						｜
			——————————————————————————
						｜
		          --------------
				  ｜  Zlogger  ｜
*/
type Zlogger struct {
	prefix     string
	format     string
	ctx        context.Context
	callerSkip int
}

// 全局logger 没有在上下问中找到logger会用这个
var globalZapLogger *zap.Logger = new(zap.Logger)
var formatJson = "json"
var basePath = ""

// SetLogger 注册logger
func SetLogger(zapLogger *zap.Logger, json bool, path string) {
	globalZapLogger = zapLogger
	if !json {
		formatJson = "plain"
	}
	basePath = path
}
func NewLogger() *Zlogger {
	return &Zlogger{
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
	l := &Zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(4).debugf(format, v...)
}
func Debugf(format string, v ...any) {
	l := &Zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(4).debugf(format, v...)

}
func InfofCtx(ctx context.Context, format string, v ...any) {
	l := &Zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(4).infof(format, v...)
}
func Infof(format string, v ...any) {
	l := &Zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(4).infof(format, v...)

}
func WarnfCtx(ctx context.Context, format string, v ...any) {
	l := &Zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(4).warnf(format, v...)
}
func Warnf(format string, v ...any) {
	l := &Zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(4).warnf(format, v...)

}
func ErrorfCtx(ctx context.Context, format string, v ...any) {
	l := &Zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(4).errorf(format, v...)
}
func Errorf(format string, v ...any) {
	l := &Zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(4).errorf(format, v...)

}

// ###########################################
// 通过对象来使用
func (l *Zlogger) Debugf(format string, v ...any) {
	l.WithCallerSkip(4).debugf(format, v)
}
func (l *Zlogger) Infof(format string, v ...any) {
	l.WithCallerSkip(4).infof(format, v)
}
func (l *Zlogger) Warnf(format string, v ...any) {
	l.WithCallerSkip(4).warnf(format, v)
}
func (l *Zlogger) Errorf(format string, v ...any) {
	l.WithCallerSkip(4).errorf(format, v)
}

// ###########################################
// 内部方法
func (l *Zlogger) debugf(format string, v ...any) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger)
	logger.Debug(fmt.Sprintf(exString+format, v...))
}
func (l *Zlogger) infof(format string, v ...any) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger)
	logger.Info(fmt.Sprintf(exString+format, v...))
}
func (l *Zlogger) warnf(format string, v ...any) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger)
	logger.Warn(fmt.Sprintf(exString+format, v...))
}
func (l *Zlogger) errorf(format string, v ...any) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger)
	logger.Error(fmt.Sprintf(exString+format, v...))
}

// ###########################################

// ###########################################
// 通用功能方法
func (l *Zlogger) getCallerSkip() int {
	if l.callerSkip == 0 {
		return 3
	}
	return l.callerSkip
}

func (l *Zlogger) formatJson() bool {
	return l.format == "json"
}

func (l *Zlogger) addCaller(_logger *zap.Logger) (zap.Logger, string) {
	format := "%s:%d"
	_, file, line, _ := runtime.Caller(l.getCallerSkip())
	_v := make([]interface{}, 0)
	file = file[len(basePath)+1:]
	_v = append(_v, file, line)
	if l.formatJson() {
		_logger = _logger.With(zap.String(loggerCallerKey, fmt.Sprintf(format, file, line)))
		return *_logger, ""
	}
	return *_logger, fmt.Sprintf(format+"\t", _v...)
}

func (l *Zlogger) addTrace(ctx context.Context, _logger *zap.Logger) (zap.Logger, string) {
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
func (l *Zlogger) addSpan(ctx context.Context, _logger *zap.Logger) (zap.Logger, string) {

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
func (l *Zlogger) addExField(ctx context.Context, _logger *zap.Logger, fieldMap map[string]string) (zap.Logger, string) {
	if exField := ctx.Value(loggerFieldKey); exField != nil {
		if l.formatJson() {
			_logger = _logger.With(exField.([]zapcore.Field)...)
			return *_logger, ""
		} else {
			format := "%v\t"
			ret := ""
			for _, field := range exField.([]zapcore.Field) {
				// 如果传入的 fieldMap 含有这个key 那么就用fieldMap中的值
				if fieldString, exit := fieldMap[field.Key]; exit {
					ret += fmt.Sprintf(format, fieldString)
				} else {
					ret += fmt.Sprintf(format, field)
				}
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
func (l *Zlogger) buildField(logger *zap.Logger, fields ...zap.Field) (zap.Logger, string) {
	// 如果外部有传入 相关字段那就用外面传入的 如果没有用自己的
	fieldMap := make(map[string]string)
	for _, field := range fields {
		fieldMap[field.Key] = field.String
	}

	var (
		caller  string
		traceId string
		spanId  string
		field   string
		exist   bool
		newLine string
	)
	// 如果map 中不存在 caller 那么使用自己的caller
	if caller, exist = fieldMap["caller"]; !exist || l.formatJson() {
		*logger, caller = l.addCaller(logger)
	}
	if traceId, exist = fieldMap["trace"]; !exist || l.formatJson() {
		*logger, traceId = l.addTrace(l.ctx, logger)
	}
	if spanId, exist = fieldMap["span"]; !exist || l.formatJson() {
		*logger, spanId = l.addSpan(l.ctx, logger)
	}
	*logger, field = l.addExField(l.ctx, logger, fieldMap)
	if l.formatJson() {
		newLine = ""
	} else {
		newLine = "\n"
	}
	return *logger, caller + traceId + spanId + field + newLine
}

// WithCallerSkip  携带跳过的层数
func (l *Zlogger) WithCallerSkip(skip int) *Zlogger {
	if skip <= 0 {
		return l
	}
	return &Zlogger{
		ctx:        l.ctx,
		callerSkip: skip,
		format:     l.format,
		prefix:     l.prefix,
	}
}
func (l *Zlogger) Sync() error {
	logger := withContext(l.ctx)
	return logger.Sync()
}

// WithContext 携带上下文
func (l *Zlogger) WithContext(ctx context.Context) *Zlogger {
	return &Zlogger{
		ctx:        ctx,
		callerSkip: l.callerSkip,
		format:     l.format,
		prefix:     l.prefix,
	}
}

// 供外部方法使用
func (l *Zlogger) debugField(msg string, fields ...zap.Field) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger, fields...)
	logger.Debug(fmt.Sprintf(exString + msg))
}
func (l *Zlogger) infoField(msg string, fields ...zap.Field) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger, fields...)
	logger.Info(fmt.Sprintf(exString + msg))
}
func (l *Zlogger) warnField(msg string, fields ...zap.Field) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger, fields...)
	logger.Warn(fmt.Sprintf(exString + msg))
}
func (l *Zlogger) errorField(msg string, fields ...zap.Field) {
	logger := withContext(l.ctx)
	logger, exString := l.buildField(&logger, fields...)
	logger.Error(fmt.Sprintf(exString + msg))
}
