package inspection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// DatabasesTop10 函数用于获取各个数据库中符合条件的表的相关信息（每个数据库取前10），并以表格形式展示，同时输出相关建议。
func DatabasesTop10(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取非template数据库名称,TOP 10 数据库size...")
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
	header := "\n### TOP 10 数据库SIZE信息:\n"
	resultWriter.WriteResult(header)
	for _, db := range dbList {
		if db == "" {
			continue
		}
		printTable(db, logWriter, resultWriter)
	}
	resultWriter.WriteResult("    建议:")
	resultWriter.WriteResult("        > 经验值: 单表超过8GB, 并且这个表需要频繁更新 或 删除+插入的话, 建议对表根据业务逻辑进行合理拆分后获得更好的性能, 以及便于对膨胀索引进行维护; 如果是只读的表, 建议适当结合SQL语句进行优化.\n")
}

// printTable 打印指定数据库的表格
func printTable(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "模式", "表名", "类型", "大小"})

	// 获取指定数据库中符合条件的表大小信息
	tableSizeInfoResult := ConnectPostgreSQL("[QUERY_TABLE_SIZE_INFO]", db)
	if len(tableSizeInfoResult) > 0 {
		for _, row := range tableSizeInfoResult {
			writer.Append(row)
		}
		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未查询到符合条件的表信息", db))
		resultWriter.WriteResult(fmt.Sprintf("在数据库 %s 中未查询到符合条件的表信息", db))
	}
}
