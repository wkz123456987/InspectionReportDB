package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// TableAgeCheck 函数用于检查表年龄情况，并以表格形式打印相关信息，同时输出相关建议。
func TableAgeCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始检查表年龄情况...")
	resultWriter.WriteResult("\n### 2.9、表年龄:\n")
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
		// 调用函数处理每个数据库的表年龄情况，更新hasData的值
		hasDataForDb := printTableAgeTable(db, logWriter, resultWriter)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {

	} else {
		resultWriter.WriteResult("未查询到表年龄相关信息")
	}

	// 打印建议
	suggestion := "> 表的年龄正常情况下应该小于vacuum_freeze_table_age, 如果剩余年龄小于5亿, 建议人为干预, 将LONG SQL或事务杀掉后, 执行vacuum freeze."

	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult(suggestion)
}

// printTableAgeTable 打印指定数据库的表年龄情况表格
func printTableAgeTable(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) bool {
	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 获取指定数据库中表年龄信息
	tableAgeInfoResult := ConnectPostgreSQL("[QUERY_TABLE_AGE_INFO]", db)
	if len(tableAgeInfoResult) == 0 {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未查询到表年龄相关信息", db))
		resultWriter.WriteResult(fmt.Sprintf("\n在数据库 %s 中未查询到表年龄相关信息\n", db))
		return false
	}

	// 写入标题
	header := fmt.Sprintf("\n**数据库 %s 的表年龄情况:**\n", db)
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	tableHeader := "| 数据库 | rolname | nspname | relkind | 表名 | 年龄 | 年龄_剩余 |"
	separator := "|--------|---------|---------|---------|------|------|----------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range tableAgeInfoResult {
		// 假设row是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |",
			row[0], row[1], row[2], row[3], row[4], row[5], row[6]))
		currentHasData = true
	}

	return currentHasData
}
