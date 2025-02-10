package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// UnusedIndexesSinceLastCheck 函数用于获取各个数据库中未使用或使用较少的索引信息，并以表格形式展示，同时输出相关建议。
func UnusedIndexesSinceLastCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取非template数据库名称,未使用或使用较少的索引信息...")
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		logWriter.WriteLog("未查询到有效数据库名称")
		return
	}
	dbList := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbList[i] = row[0]
	}
	// 写入标题
	header := "\n###  2.4、上次巡检以来未使用或使用较少的索引:\n"
	resultWriter.WriteResult(header)
	for _, db := range dbList {
		if db == "" {
			continue
		}
		printUnusedIndexes(db, logWriter, resultWriter)
	}
	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult("> 建议和应用开发人员确认后, 删除不需要的索引.")
}

// printUnusedIndexes 打印指定数据库中未使用或使用较少的索引
func printUnusedIndexes(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	// 获取指定数据库中未使用或使用较少的索引信息
	unusedIndexesInfoResult := ConnectPostgreSQL("[QUERY_UNUSED_INDEXES_INFO]", db)
	if len(unusedIndexesInfoResult) == 0 {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未查询到上次巡检以来未使用或使用较少的索引信息", db))
		resultWriter.WriteResult(fmt.Sprintf("\n在数据库 %s 中未查询到上次巡检以来未使用或使用较少的索引信息\n", db))
		return
	}

	// 写入标题
	header := fmt.Sprintf("\n**数据库 %s 中未使用或使用较少的索引:**\n", db)
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	tableHeader := "| 当前数据库 | 模式名 | 表名 | 索引名 |"
	separator := "|------------|------|----|--------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range unusedIndexesInfoResult {
		// 假设row是一个包含所需字段的切片
		// 请确保row的长度至少为3，否则需要调整索引
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s |",
			row[0], row[1], row[2], row[3]))
	}
}
