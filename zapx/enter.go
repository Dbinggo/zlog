package zapx

import (
	"github.com/dbinggo/zlog"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

func InitLogger(rest rest.RestConf) {
	// zap 配置
	var zapConfig = ZapConfig{
		// 是否为 json格式
		Format: rest.Log.Encoding,
		// bug 等级
		Level: rest.Log.Level,
		// 是否开启彩色（Info 颜色）
		Colour: rest.Log.Encoding == "plain" && rest.Log.Mode == "console",
		// 日志存储路径 会在路径下生成 info.log error.log
		FilePath: rest.Log.Path,
		// 是否存储日志
		File: rest.Log.Mode == "file",
		// 是否在控制台输出
		Terminal: rest.Log.Mode == "console",
	}

	logger := GetLogger(zapConfig)
	zlog.SetLogger(logger, rest.Log.Encoding == "json")
	zlogger := zlog.NewLogger()
	var zapWriter logx.Writer
	zapWriter = zlog.NewZeroLogger(zlogger, rest.Log.Encoding == "json")
	logx.SetWriter(zapWriter)
}
