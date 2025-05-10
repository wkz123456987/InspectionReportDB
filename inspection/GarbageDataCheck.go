package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GarbageDataCheck 函数用于检查数据库中垃圾数据情况，并以表格形式打印相关信息，同时输出相关建议。
func GarbageDataCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始检查数据库中垃圾数据情况...")
	resultWriter.WriteResult("\n###   2.7、检查数据库中垃圾数据情况:\n")
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
		// 调用函数处理每个数据库的垃圾数据情况，更新hasData的值
		hasDataForDb := printGarbageDataTable(db, logWriter, resultWriter)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {

	} else {
		resultWriter.WriteResult("未查询到数据库中垃圾数据相关信息")
	}

	// 打印建议
	suggestion := " > 通常垃圾过多, 可能是因为无法回收垃圾, 或者回收垃圾的进程繁忙或没有及时唤醒, 或者没有开启autovacuum, 或在短时间内产生了大量的垃圾.可以等待autovacuum进行处理, 或者手工执行vacuum table."
	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult(suggestion)
}

// printGarbageDataTable 打印指定数据库的垃圾数据情况表格// printGarbageDataTable 打印指定数据库的垃圾数据情况表格
func printGarbageDataTable(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) bool {
	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 获取指定数据库中垃圾数据信息
	garbageDataInfoResult := ConnectPostgreSQL("[QUERY_GARBAGE_DATA_INFO]", db)
	if len(garbageDataInfoResult) == 0 {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未查询到垃圾数据相关信息", db))
		resultWriter.WriteResult(fmt.Sprintf("\n在数据库 %s 中未查询到垃圾数据相关信息\n", db))
		return false
	}

	// 写入标题
	header := fmt.Sprintf("\n**数据库 %s 的垃圾数据情况:**", db)
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	tableHeader := "| 数据库 | schema | 表名 | 死元组数量 |"
	separator := "|--------|--------|------|------------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range garbageDataInfoResult {
		// 假设row是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s |",
			row[0], row[1], row[2], row[3]))
		currentHasData = true
	}

	return currentHasData
}
