package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// TablesWithTooManyIndexes 查找索引数超过4并且SIZE大于10MB的表
func TablesWithTooManyIndexes(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始查找索引数超过4并且SIZE大于10MB的表...")
	// 获取非template数据库名称
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
	header := "\n###  2.2、查找索引数超过4并且SIZE大于10MB的表\n"
	resultWriter.WriteResult(header)
	for _, db := range dbList {
		if db == "" {
			continue
		}
		printTablesWithTooManyIndexes(db, logWriter, resultWriter)
	}
	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult("> 索引数量太多, 影响表的增删改性能, 建议检查是否有不需要的索引.\n")
}

// printTablesWithTooManyIndexes 打印指定数据库中索引数超过4且SIZE大于10MB的表
func printTablesWithTooManyIndexes(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	// 获取指定数据库中索引数超过4且SIZE大于10MB的表信息
	tablesInfoResult := ConnectPostgreSQL("[QUERY_TABLES_WITH_TOO_MANY_INDEXES]", db)
	if len(tablesInfoResult) == 0 {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未查询到索引数超过4且SIZE大于10MB的表信息", db))
		resultWriter.WriteResult(fmt.Sprintf("\n在数据库 %s 中未查询到索引数超过4且SIZE大于10MB的表信息\n", db))
		return
	}

	// 写入标题
	header := fmt.Sprintf("\n**数据库 %s 中索引数超过4且SIZE大于10MB的表:**\n", db)
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	tableHeader := "| 数据库 | 模式 | 表名 | 表大小 | 索引数量 |"
	separator := "|--------|------|------|--------|----------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range tablesInfoResult {
		// 假设row是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s |",
			row[0], row[1], row[2], row[3], row[4]))
	}
}
