package zlog

import (
	"context"
	"github.com/dbinggo/zlog/zapx"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"testing"
)

func TestZlog(t *testing.T) {
	SetLogger(zapLogger, false)
	ctx := context.Background()
	ctx = AddFiled(ctx, zap.String("traceId", "123456"))
	Debugf("redis中未找到数据")
	DebugfCtx(ctx, "redis中未找到数据")

}
func TestLogx(t *testing.T) {
	zapLogger := zapx.Develop()

	SetLogger(zapLogger, false)
	ctx := context.Background()
	ctx = AddFiled(ctx, zap.String("traceId", "123456"))

	zlogger := NewLogger()
	zeroLogger := NewZeroLogger(zlogger, false)
	logx.SetWriter(zeroLogger)

	ctx = AddFiled(ctx, zap.String("traceId", "123456"))
	//logx.Info("redis中未找到数据")
	//logx.Info("redis中未找到数据")
	logx.WithContext(ctx).Debug("redis中未找到数据")
}
