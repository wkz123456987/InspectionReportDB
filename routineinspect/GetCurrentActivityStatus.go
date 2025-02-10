package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetCurrentActivityStatus 函数用于获取数据库当前活跃度状态信息，并以表格形式打印相关信息，同时输出相关建议。
func GetCurrentActivityStatus(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库当前活跃度状态信息...")
	resultWriter.WriteResult("\n### 3.1、当前活跃度:\n")

	// Markdown 表格的表头和分隔行
	tableHeader := "| 当前时间 | 状态 | count |"
	separator := "|----------|------|-------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	result := ConnectPostgreSQL("[QUERY_ACTIVITY_STATUS]")
	if len(result) > 0 {
		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
				row[0], row[1], row[2]))
		}
	} else {
		logWriter.WriteLog("未查询到有效数据")
		resultWriter.WriteResult("未查询到有效数据")
	}

	// 打印建议
	suggestion := "> 如果active状态很多, 说明数据库比较繁忙。如果idle in transaction很多, 说明业务逻辑设计可能有问题。如果idle很多, 可能使用了连接池, 并且可能没有自动回收连接到连接池的最小连接数。"
	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult(suggestion)
}
