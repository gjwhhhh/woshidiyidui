package oss

import (
	"fmt"
	"os"
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
	default:
		fmt.Println("The current server system is not linux or windows, and the startup fails")
		os.Exit(-1)
	}
}
