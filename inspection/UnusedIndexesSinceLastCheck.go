package inspection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
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
	header := "\n###  上次巡检以来未使用或使用较少的索引:\n"
	resultWriter.WriteResult(header)
	for _, db := range dbList {
		if db == "" {
			continue
		}
		printUnusedIndexes(db, logWriter, resultWriter)
	}
	resultWriter.WriteResult("    建议:")
	resultWriter.WriteResult("        > 建议和应用开发人员确认后, 删除不需要的索引.")
}

// printUnusedIndexes 打印指定数据库中未使用或使用较少的索引
func printUnusedIndexes(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "未使用数量"})

	// 获取指定数据库中未使用或使用较少的索引信息
	unusedIndexesInfoResult := ConnectPostgreSQL("[QUERY_UNUSED_INDEXES_INFO]", db)
	if len(unusedIndexesInfoResult) > 0 {
		for _, row := range unusedIndexesInfoResult {
			writer.Append(row)
		}
		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未查询到上次巡检以来未使用或使用较少的索引信息", db))
		resultWriter.WriteResult(fmt.Sprintf("在数据库 %s 中未查询到上次巡检以来未使用或使用较少的索引信息", db))
	}
}
