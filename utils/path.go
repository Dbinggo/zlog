package utils

import (
	"path/filepath"
	"runtime"
)

// 获取项目根路径
func GetRootPath(myPath string) string {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("Something wrong with getting root path")
	}
	absPath, err := filepath.Abs(fileName)
	rootPath := filepath.Dir(filepath.Dir(absPath))
	if err != nil {
		panic(any(err))
	}
	return filepath.Join(rootPath, myPath)
}
