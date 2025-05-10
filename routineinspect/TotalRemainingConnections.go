package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// CheckDBConnections 函数用于获取数据库连接相关信息，并以表格形式展示，同时输出相关建议。
func CheckDBConnections(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库连接相关信息...")
	resultWriter.WriteResult("\n### 3.2、数据库连接信息:\n")

	// Markdown 表格的表头
	tableHeader := "| 总连接 | 已使用连接 | 剩余给超级用户连接 | 剩余给普通用户连接 |"
	resultWriter.WriteResult(tableHeader)

	// Markdown 表格的分隔行
	separator := "|--------|------------|---------------------|---------------------|"
	resultWriter.WriteResult(separator)

	result := ConnectPostgreSQL("[QUERY_DB_CONNECTIONS]")
	if len(result) > 0 {
		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s |",
				row[0], row[1], row[2], row[3]))
		}
	} else {
		logWriter.WriteLog("未查询到有效数据")
		resultWriter.WriteResult("未查询到有效数据")
	}

	// 打印建议
	suggestion := "> 给超级用户和普通用户设置足够的连接，以免不能登录数据库。"
	resultWriter.WriteResult("\n**建议:** \n" + suggestion)
}
