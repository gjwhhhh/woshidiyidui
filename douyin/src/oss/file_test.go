package oss

import (
	"bou.ke/monkey"
	"os"
	"testing"
)

// 测试文件上传
func TestFileUpLoad(t *testing.T) {
	// 获取当前测试目录路径
	dir, _ := os.Getwd()
	// 将本地文件路径转换为测试目录路径
	LocalFilePathPrefix = dir + "\\"
	// 取消删除本地文件这一操作
	monkey.Patch(os.Remove, func(name string) error {
		return nil
	})
	// 写在打桩
	defer monkey.UnpatchAll()
	// 控制台打印文件路径
	println(FileUpLoad("test.txt"))
}
