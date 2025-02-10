package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetUserConnectionLimits 用于获取用户连接数限制相关信息，并以表格形式展示，同时输出相关建议。
func GetUserConnectionLimits(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取用户连接数限制相关信息...")
	resultWriter.WriteResult("\n### 3.10、用户连接数限制:\n")

	// 获取用户连接数限制相关信息
	result := ConnectPostgreSQL("[QUERY_USER_CONNECTION_LIMITS]")
	if len(result) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 用户名 | 用户连接限制 | 当前用户已使用的连接数 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|--------|----------------|------------------------|"
		resultWriter.WriteResult(separator)

		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
				row[0], row[1], row[2]))
		}
	} else {
		logWriter.WriteLog("未查询到用户连接数限制相关信息")
		resultWriter.WriteResult("未查询到用户连接数限制相关信息")
	}

	// 打印建议
	suggestion := "> 给用户设置足够的连接数，使用alter role... CONNECTION LIMIT来设置。"
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
