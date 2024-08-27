package test

import (
	"context"
	"github.com/dbinggo/zlog"
	"github.com/dbinggo/zlog/utils"
	"github.com/dbinggo/zlog/zapx"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	zlogger := zapx.Develop(utils.GetRootPath(""))
	zlog.NewGormLogger(zlogger)
	db, err := gorm.Open(mysql.Open("root:13131227873@tcp(123.207.73.185:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: zlog.NewGormLogger(zlogger),
	})
	// todo
	type yjddbTest struct {
		age int
	}
	err = db.AutoMigrate(yjddbTest{})
	if err != nil {
		zlog.WarnfCtx(context.Background(), "AutoMigrate %v", err)
	}
}
