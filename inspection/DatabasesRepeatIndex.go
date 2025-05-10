package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// DatabasesRepeatIndex 函数用于检查数据库中重复创建的索引，并以表格形式打印相关信息，同时输出相关建议。
func DatabasesRepeatIndex(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取非template数据库名称, 重复索引信息...")
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
	header := "\n###  2.3、重复创建的索引:\n"
	resultWriter.WriteResult(header)
	for _, db := range dbList {
		if db == "" {
			continue
		}
		printRepeatIndexTable(db, logWriter, resultWriter)
	}
	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult("> 当创建重复索引后，不会对数据库的性能产生优化作用，反而会产生一些维护上的成本，请删除重复索引")
}

// printRepeatIndexTable 打印指定数据库的重复索引表格
func printRepeatIndexTable(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	// 获取指定数据库中重复索引信息
	repeatIndexInfoResult := ConnectPostgreSQL("[QUERY_REPEAT_INDEX_INFO]", db)
	if len(repeatIndexInfoResult) == 0 {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未检测到重复创建的索引信息", db))
		resultWriter.WriteResult(fmt.Sprintf("\n在数据库 %s 中未检测到重复创建的索引信息\n", db))
		return
	}

	// 写入标题
	header := fmt.Sprintf("\n**数据库 %s 的重复索引表格:**\n", db)
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	tableHeader := "| 数据库 | 表名 | 索引名 |"
	separator := "|--------|------|--------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range repeatIndexInfoResult {
		// 假设row是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
			row[0], row[1], row[2]))
	}
}
