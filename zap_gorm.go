package zlog

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type gormLogger struct {
	logger *Zlogger
	level  logger.LogLevel
}

var _ logger.Interface = (*gormLogger)(nil) // 接口实现检查

//// Interface logger interface
//type Interface interface {
//	LogMode(LogLevel) Interface
//	Info(context.Context, string, ...interface{})
//	Warn(context.Context, string, ...interface{})
//	Error(context.Context, string, ...interface{})
//	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
//}

func NewGormLogger(logger *Zlogger) *gormLogger {
	return &gormLogger{
		logger: logger,
	}
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

func (l *gormLogger) Info(ctx context.Context, format string, value ...interface{}) {
	if l.level < logger.Info {
		return
	}

	l.logger.WithContext(ctx).WithCallerSkip(6).infof(format, value)
}

func (l *gormLogger) Warn(ctx context.Context, format string, value ...interface{}) {
	if l.level < logger.Warn {
		return
	}

	l.logger.WithContext(ctx).WithCallerSkip(6).warnf(format, value)
}

func (l *gormLogger) Error(ctx context.Context, format string, value ...interface{}) {
	if l.level < logger.Error {
		return
	}

	l.logger.WithContext(ctx).WithCallerSkip(6).errorf(format, value)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 请求和返回条数
	sql, rows := fc()

	// 通用字段
	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("time", elapsed),
		zap.Int64("rows", rows),
	}
	// Gorm 错误
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if l.level >= logger.Warn {
				if l.logger.formatJson() {
					l.logger.WithContext(ctx).WithCallerSkip(6).warnField("Database ErrRecordNotFound", logFields...)
				} else {
					l.logger.WithContext(ctx).WithCallerSkip(6).warnf("Database ErrRecordNotFound sql: %s ,time: %v ,rows: %d", sql, elapsed, rows)
				}

			}
		} else {
			// 其他错误使用 error 等级
			if l.level >= logger.Error {
				if l.logger.formatJson() {
					l.logger.WithContext(ctx).WithCallerSkip(6).errorField("Database Error", logFields...)
				} else {
					l.logger.WithContext(ctx).WithCallerSkip(6).errorf("Database Error sql: %s ,time: %v ,rows: %d", sql, elapsed, rows)
				}
			}
		}
	}

	// 慢查询日志
	if elapsed > (200 * time.Millisecond) {
		if l.logger.formatJson() {
			l.logger.WithContext(ctx).WithCallerSkip(6).warnField("Database Slow Log", logFields...)
		} else {
			l.logger.WithContext(ctx).WithCallerSkip(6).warnf("Database Slow Log sql: %s ,time: %v ,rows: %d", sql, elapsed, rows)
		}
	}

	// 记录所有 SQL 请求
	if l.logger.formatJson() {
		l.logger.WithContext(ctx).WithCallerSkip(6).infoField("Database Query", logFields...)
	} else {
		l.logger.WithContext(ctx).WithCallerSkip(6).infof("Database Query sql: %s ,time: %v ,rows: %d", sql, elapsed, rows)
	}
}
