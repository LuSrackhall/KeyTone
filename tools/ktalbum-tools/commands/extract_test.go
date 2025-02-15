package commands

import (
	"os"
	"testing"
)

func TestExtract(t *testing.T) {
	// 创建一个临时测试文件
	testFile := "test.ktalbum"
	outputFile := "test.zip"
	
	// 清理测试文件
	defer func() {
		os.Remove(testFile)
		os.Remove(outputFile)
	}()

	// 测试文件不存在的情况
	err := Extract(testFile, outputFile, true)
	if err == nil {
		t.Error("应该返回错误：文件不存在")
	}
} 