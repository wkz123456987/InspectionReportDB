package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetStandbyInfo 用于获取备库信息，并以表格形式展示。
func GetStandbyInfo(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取备库信息...")
	resultWriter.WriteResult("\n### 3.18、备库信息:\n")

	// 获取备库信息
	result := ConnectPostgreSQL("[QUERY_STANDBY_INFO]")
	if len(result) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 用户名 | 应用程序名称 | 客户端地址 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|--------|--------------|------------|"
		resultWriter.WriteResult(separator)

		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
				row[0], row[1], row[2]))
		}
	} else {
		logWriter.WriteLog("未查询到备库信息相关信息")
		resultWriter.WriteResult("未查询到备库信息相关信息")
	}
}
