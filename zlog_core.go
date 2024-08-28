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

const coreCallSkip = 4

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
		callerSkip: 2,
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
	l.WithCallerSkip(coreCallSkip).debugf(format, v...)
}
func Debugf(format string, v ...any) {
	l := &Zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(coreCallSkip).debugf(format, v...)

}
func InfofCtx(ctx context.Context, format string, v ...any) {
	l := &Zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(coreCallSkip).infof(format, v...)
}
func Infof(format string, v ...any) {
	l := &Zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(coreCallSkip).infof(format, v...)

}
func WarnfCtx(ctx context.Context, format string, v ...any) {
	l := &Zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(coreCallSkip).warnf(format, v...)
}
func Warnf(format string, v ...any) {
	l := &Zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(coreCallSkip).warnf(format, v...)

}
func ErrorfCtx(ctx context.Context, format string, v ...any) {
	l := &Zlogger{ctx: ctx, format: formatJson}
	l.WithCallerSkip(coreCallSkip).errorf(format, v...)
}
func Errorf(format string, v ...any) {
	l := &Zlogger{ctx: context.Background(), format: formatJson}
	l.WithCallerSkip(coreCallSkip).errorf(format, v...)

}

// ###########################################
// 通过对象来使用
func (l *Zlogger) Debugf(format string, v ...any) {
	l.WithCallerSkip(coreCallSkip).debugf(format, v)
}
func (l *Zlogger) Infof(format string, v ...any) {
	l.WithCallerSkip(coreCallSkip).infof(format, v)
}
func (l *Zlogger) Warnf(format string, v ...any) {
	l.WithCallerSkip(coreCallSkip).warnf(format, v)
}
func (l *Zlogger) Errorf(format string, v ...any) {
	l.WithCallerSkip(coreCallSkip).errorf(format, v)
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
	return l.callerSkip
}

func (l *Zlogger) FormatJson() bool {
	return l.format == "json"
}

func (l *Zlogger) addCaller(_logger *zap.Logger) (zap.Logger, string) {
	callerSkip := l.getCallerSkip()
	var file string
	var line int
	_v := make([]interface{}, 0)
	// 代表溯源到本项目路径
	if callerSkip == 0 {
		for i := 0; ; i++ {
			_, file, line, _ = runtime.Caller(i)
			if len(file) >= len(basePath) && file[:len(basePath)] == basePath {
				file = file[len(basePath)+1:]
				break
			}
		}
	} else if callerSkip < 0 {
		callerSkip = -callerSkip
		_, file, line, _ = runtime.Caller(callerSkip)

	} else {
		_, file, line, _ = runtime.Caller(callerSkip)
		file = file[len(basePath)+1:]
	}
	format := "%s:%d"
	_v = append(_v, file, line)
	if l.FormatJson() {
		_logger = _logger.With(zap.String(loggerCallerKey, fmt.Sprintf(format, file, line)))
		return *_logger, ""
	}
	return *_logger, fmt.Sprintf(format+" \t", _v...)
}

func (l *Zlogger) addTrace(ctx context.Context, _logger *zap.Logger) (zap.Logger, string) {
	traceId := trace.TraceIDFromContext(ctx)
	if traceId == "" {
		return *_logger, ""
	}
	if l.FormatJson() {
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
	if l.FormatJson() {
		_logger = _logger.With(zap.String(loggerSpanKey, spanId))
		return *_logger, ""
	}
	format := "%v\t"
	return *_logger, fmt.Sprintf(format, spanId)
}
func (l *Zlogger) addExField(ctx context.Context, _logger *zap.Logger, fieldMap map[string]string) (zap.Logger, string) {
	if exField := ctx.Value(loggerFieldKey); exField != nil {
		if l.FormatJson() {
			_logger = _logger.With(exField.([]zapcore.Field)...)
			return *_logger, ""
		} else {
			format := "%v \t"
			ret := ""
			for _, field := range exField.([]zapcore.Field) {
				// 如果传入的 fieldMap 含有这个key 那么就用fieldMap中的值
				if fieldString, exit := fieldMap[field.Key]; exit {
					ret += fmt.Sprintf(format, fieldString)
				} else {
					ret += fmt.Sprintf(format, field.String)
				}
				//if !slices.Contains([]string{loggerCallerKey, loggerTraceKey, loggerSpanKey}, field.Key) {
				//	ret += fmt.Sprintf(format, field.String)
				//}
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
	if caller, exist = fieldMap[loggerCallerKey]; !exist || l.FormatJson() {
		*logger, caller = l.addCaller(logger)
	}
	if traceId, exist = fieldMap[loggerTraceKey]; !exist || l.FormatJson() {
		*logger, traceId = l.addTrace(l.ctx, logger)
	}
	if spanId, exist = fieldMap[loggerSpanKey]; !exist || l.FormatJson() {
		*logger, spanId = l.addSpan(l.ctx, logger)
	}
	//*logger, caller = l.addCaller(logger)
	//*logger, traceId = l.addTrace(l.ctx, logger)
	//*logger, spanId = l.addSpan(l.ctx, logger)

	*logger, field = l.addExField(l.ctx, logger, fieldMap)
	if l.FormatJson() {
		newLine = ""
	} else {
		newLine = "\n"
	}
	return *logger, caller + traceId + spanId + field + newLine
}

// WithCallerSkip  携带跳过的层数
func (l *Zlogger) WithCallerSkip(skip int) *Zlogger {
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
func printStack() {
	// 获取调用栈的程序计数器
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	// 遍历调用栈帧
	for {
		frame, more := frames.Next()
		fmt.Printf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}
