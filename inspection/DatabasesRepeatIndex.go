package inspection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
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
	header := "\n###  重复创建的索引:\n"
	resultWriter.WriteResult(header)
	for _, db := range dbList {
		if db == "" {
			continue
		}
		printRepeatIndexTable(db, logWriter, resultWriter)
	}
	resultWriter.WriteResult("    建议:")
	resultWriter.WriteResult("        > 当创建重复索引后，不会对数据库的性能产生优化作用，反而会产生一些维护上的成本，请删除重复索引")
}

// printRepeatIndexTable 打印指定数据库的重复索引表格
func printRepeatIndexTable(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "表名", "索引名"})

	// 获取指定数据库中重复索引信息
	repeatIndexInfoResult := ConnectPostgreSQL("[QUERY_REPEAT_INDEX_INFO]", db)
	if len(repeatIndexInfoResult) > 0 {
		for _, row := range repeatIndexInfoResult {
			writer.Append(row)
		}
		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未检测到重复创建的索引信息", db))
		resultWriter.WriteResult(fmt.Sprintf("在数据库 %s 中未检测到重复创建的索引信息", db))
	}
}
