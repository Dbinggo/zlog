package zapx

import (
	"github.com/dbinggo/zlog"
	"github.com/zeromicro/go-zero/core/logx"
)

// 从zero的配置文件中获得配置 减少代码侵入性
func InitLoggerFromZero(conf logx.LogConf, basePath string) *zlog.Zlogger {
	// zap 配置
	var zapConfig = ZapConfig{
		// 是否为 json格式
		Format: conf.Encoding,
		// bug 等级
		Level: conf.Level,
		// 是否开启彩色（Info 颜色）
		Colour: conf.Encoding == "plain" && conf.Mode == "console",
		// 日志存储路径 会在路径下生成 info.log error.log
		FilePath: conf.Path,
		// 是否存储日志
		File: conf.Mode == "file",
		// 是否在控制台输出
		Terminal: conf.Mode == "console",
	}

	logger := GetLogger(zapConfig)
	zlog.SetLogger(logger, conf.Encoding == "json", basePath)
	zlogger := zlog.NewLogger()
	return zlogger
}

// Develop 开发模式
func Develop(basePath string) *zlog.Zlogger {
	return InitLoggerFromZero(logx.LogConf{
		Mode:     "console",
		Encoding: "plain",
	}, basePath)
}
