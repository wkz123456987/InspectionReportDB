package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetTablespaceUsage 函数用于获取表空间使用情况信息，并以表格形式展示，同时输出相关建议。
func GetTablespaceUsage(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取表空间使用情况信息...")
	resultWriter.WriteResult("\n### 3.8、表空间使用情况:\n")

	// 查询表空间使用情况信息
	tablespaceUsageResult := ConnectPostgreSQL("[QUERY_TABLESPACE_USAGE]")
	if len(tablespaceUsageResult) == 0 {
		logWriter.WriteLog("未查询到表空间使用情况信息")
		resultWriter.WriteResult("未查询到表空间使用情况信息")
		return
	}

	// Markdown 表格的表头
	tableHeader := "| 表空间名称 | 表空间路径 | 表空间大小 |"
	resultWriter.WriteResult(tableHeader)

	// Markdown 表格的分隔行
	separator := "|------------|------------|------------|"
	resultWriter.WriteResult(separator)

	// 遍历查询结果，输出表空间使用情况信息
	for _, row := range tablespaceUsageResult {
		// 假设 row 是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
			row[0], row[1], row[2]))
	}

	// 打印建议
	suggestion := "> 定期检查表空间的使用情况，确保表空间有足够的空间以避免存储问题。如果表空间不足，可以考虑扩展表空间或迁移数据到更大的存储设备。"
	resultWriter.WriteResult("\n**建议:** \n" + suggestion)
}
