package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetArchiveAndAutoVacuumSettings 用于获取是否开启归档、自动垃圾回收相关设置信息，并以表格形式展示，同时输出相关建议。
func GetArchiveAndAutoVacuumSettings(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取是否开启归档、自动垃圾回收相关设置信息...")
	resultWriter.WriteResult("\n### 3.16、是否开启归档、自动垃圾回收相关设置信息:\n")

	// 获取是否开启归档、自动垃圾回收设置信息
	result := ConnectPostgreSQL("[QUERY_ARCHIVE_AND_AUTOVACUUM_SETTINGS]")
	if len(result) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 名称 | 设置值 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|------|--------|"
		resultWriter.WriteResult(separator)

		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s |",
				row[0], row[1]))
		}
	} else {
		logWriter.WriteLog("未查询到是否开启归档、自动垃圾回收相关设置信息")
		resultWriter.WriteResult("未查询到是否开启归档、自动垃圾回收相关设置信息")
	}

	// 打印建议
	suggestion := "> 如果当前的wal文件和最后一个归档失败的wal文件之间相差很多个文件，建议尽快排查归档失败的原因，以便修复，否则pg_wal目录可能会撑爆。"
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
