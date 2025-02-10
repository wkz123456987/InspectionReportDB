package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// IndexBloatCheck 函数用于检查数据库中索引膨胀情况，并以表格形式打印相关信息，同时输出相关建议。
func IndexBloatCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始检查数据库中索引膨胀情况...")
	resultWriter.WriteResult("\n###  2.6、检查数据库中索引膨胀情况:\n")
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 获取非template数据库名称
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		logWriter.WriteLog("未查询到有效数据库名称")
		resultWriter.WriteResult("未查询到有效数据库名称")
		return
	}
	dbList := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbList[i] = row[0]
	}

	for _, db := range dbList {
		if db == "" {
			continue
		}
		// 调用函数处理每个数据库的索引膨胀情况，更新hasData的值
		hasDataForDb := printIndexBloatTable(db, logWriter, resultWriter)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {

	} else {
		resultWriter.WriteResult("未查询到数据库中索引膨胀相关信息")
	}

	// 打印建议
	suggestion := "> 如果索引膨胀太大, 会影响性能, 建议重建索引, create index CONCURRENTLY...."
	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult(suggestion)
}

// printIndexBloatTable 打印指定数据库的索引膨胀情况表格
func printIndexBloatTable(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) bool {
	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 获取指定数据库中索引膨胀信息
	indexBloatInfoResult := ConnectPostgreSQL("[QUERY_INDEX_BLOAT_INFO]", db)
	if len(indexBloatInfoResult) == 0 {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未查询到索引膨胀相关信息", db))
		resultWriter.WriteResult(fmt.Sprintf("\n在数据库 %s 中未查询到索引膨胀相关信息\n", db))
		return currentHasData
	}

	// 写入标题
	header := fmt.Sprintf("\n**数据库 %s 的索引膨胀情况:**\n", db)
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	tableHeader := "| 数据库 | schema | 表名 | 表膨胀系数 | 索引名 | 索引膨胀系数 |"
	separator := "|--------|--------|------|------------|--------|------------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range indexBloatInfoResult {
		// 假设row是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s",
			row[0], row[1], row[2], row[3], row[4], row[5]))
		currentHasData = true
	}

	return currentHasData
}
