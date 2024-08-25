package test

import (
	"context"
	"github.com/dbinggo/zlog"
	"github.com/dbinggo/zlog/utils"
	"github.com/dbinggo/zlog/zapx"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"testing"
)

func TestZlog(t *testing.T) {
	zapx.Develop(utils.GetRootPath(""))
	ctx := context.Background()
	ctx = zlog.AddFiled(ctx, zap.String("traceId", "123456"))
	zlog.Debugf("redis中未找到数据")
	zlog.DebugfCtx(ctx, "redis中未找到数据")

}
func TestLogx(t *testing.T) {
	zlogger := zapx.Develop(utils.GetRootPath(""))
	writer := zlog.NewZeroWriter(zlogger)
	logx.SetWriter(writer)
	ctx := context.Background()
	ctx = zlog.AddFiled(ctx, zap.String("traceId", "123456"))

	ctx = zlog.AddFiled(ctx, zap.String("traceId", "123456"))
	//logx.Info("redis中未找到数据")
	//logx.Info("redis中未找到数据")
	logx.WithContext(ctx).Debug("redis中未找到数据")
}
func TestGormLogger(t *testing.T) {
	// todo
}
