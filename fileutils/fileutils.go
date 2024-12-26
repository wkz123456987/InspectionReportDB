package fileutils

import (
	"os"
)

// CreateFile 创建文件，如果文件已存在则打开文件准备追加内容
func CreateFile(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// WriteLog 写入日志文件
func WriteLog(file *os.File, message string) error {
	_, err := file.WriteString(message + "\n")
	return err
}

// WriteResult 写入结果文件
func WriteResult(file *os.File, message string) error {
	_, err := file.WriteString(message + "\n")
	return err
}
