package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetMasterStandbyRole 用于获取主备库角色信息，并以表格形式展示。
func GetMasterStandbyRole(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取主备库角色信息...")
	resultWriter.WriteResult("\n### 3.17、数据库主备角色:\n")

	// 获取主备库角色信息
	result := ConnectPostgreSQL("[QUERY_MASTER_STANDBY_ROLE]")
	if len(result) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 数据库主备角色 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|-----------------|"
		resultWriter.WriteResult(separator)

		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s |", row[0]))
		}
	} else {
		logWriter.WriteLog("未查询到主备库角色相关信息")
		resultWriter.WriteResult("未查询到主备库角色相关信息")
	}
}
