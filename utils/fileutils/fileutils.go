package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogWriter 结构体包含日志文件的写入句柄
type LogWriter struct {
	File *os.File
}

// ResultWriter 结构体包含结果文件的写入句柄
type ResultWriter struct {
	File *os.File
}

// NewLogWriter 创建并返回一个LogWriter实例
func NewLogWriter(logDir, logFileName string) (*LogWriter, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	// 创建或打开日志文件
	logFilePath := filepath.Join(logDir, logFileName)
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &LogWriter{File: file}, nil
}

// NewResultWriter 创建并返回一个ResultWriter实例
func NewResultWriter(resultDir, resultFileName string) (*ResultWriter, error) {
	// 确保结果目录存在
	if err := os.MkdirAll(resultDir, 0755); err != nil {
		return nil, err
	}
	// 创建或打开结果文件
	resultFilePath := filepath.Join(resultDir, resultFileName)
	file, err := os.OpenFile(resultFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &ResultWriter{File: file}, nil
}

// WriteLog 追加写入日志内容
func (lw *LogWriter) WriteLog(message string) error {
	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05.000 MST")
	// 创建带有时间戳的日志消息
	logMessage := fmt.Sprintf("%s %s", now, message)
	fmt.Println(logMessage) // 打印带时间戳的日志消息
	_, err := lw.File.WriteString(logMessage + "\n")
	if err != nil {
		return err
	}
	return lw.File.Sync() // 确保数据被写入磁盘
}

// WriteResult 追加写入结果内容
func (rw *ResultWriter) WriteResult(message string) error {
	_, err := rw.File.WriteString(message + "\n")
	if err != nil {
		return err
	}
	return rw.File.Sync() // 确保数据被写入磁盘
}

// Close 关闭文件
func (lw *LogWriter) Close() error {
	return lw.File.Close()
}

// Close 关闭文件
func (rw *ResultWriter) Close() error {
	return rw.File.Close()
}
