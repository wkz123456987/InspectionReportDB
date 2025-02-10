package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetDatabaseConnectionLimits 用于获取数据库连接限制相关信息，并以表格形式展示，同时输出相关建议。
func GetDatabaseConnectionLimits(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库连接限制相关信息...")
	resultWriter.WriteResult("\n### 3.11、数据库连接数限制:\n")

	// 获取数据库连接限制相关信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_CONNECTION_LIMITS]")
	if len(result) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 数据库 | 数据库连接限制 | 数据库已使用连接 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|--------|----------------|-------------------|"
		resultWriter.WriteResult(separator)

		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
				row[0], row[1], row[2]))
		}
	} else {
		logWriter.WriteLog("未查询到数据库连接限制相关信息")
		resultWriter.WriteResult("未查询到数据库连接限制相关信息")
	}

	// 打印建议
	suggestion := "> 给数据库设置足够的连接数，使用alter database... CONNECTION LIMIT来设置。"
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
