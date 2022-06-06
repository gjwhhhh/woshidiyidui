package oss

import (
	"runtime"
)

var LocalFilePathPrefix string

// 初始化设置本地文件前缀
func init() {
	goos := runtime.GOOS
	// 根据系统决定文件缓存的存储路径
	switch goos {
	case "linux":
		LocalFilePathPrefix = "/usr/DouYinCache/"
	case "windows":
		LocalFilePathPrefix = "C:\\DouYinCache\\"
	case "darwin":
		LocalFilePathPrefix = "/Users/weixizi/douyin/"
	default:
		panic("The current server system is not linux or windows, and the startup fails")
	}
}
