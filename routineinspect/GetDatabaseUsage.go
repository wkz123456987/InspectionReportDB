package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetDatabaseUsage 用于获取数据库使用情况信息，并以表格形式展示，同时输出相关建议。
func GetDatabaseUsage(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库使用情况信息...")
	resultWriter.WriteResult("\n### 3.9、数据库使用情况:\n")

	// 获取数据库使用情况信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_USAGE]")
	if len(result) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 数据库 | 数据库大小 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|--------|------------|"
		resultWriter.WriteResult(separator)

		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s |",
				row[0], row[1]))
		}
	} else {
		logWriter.WriteLog("未查询到数据库使用情况相关信息")
		resultWriter.WriteResult("未查询到数据库使用情况相关信息")
	}

	// 打印建议
	suggestion := "> 注意检查数据库的大小，是否需要清理历史数据。"
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
