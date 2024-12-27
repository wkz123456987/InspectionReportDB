package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetDatabaseUsage 用于获取数据库使用情况信息，并以表格形式展示，同时输出相关建议。
func GetDatabaseUsage(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库使用情况信息...")
	resultWriter.WriteResult("\n###  数据库使用情况:\n")
	// 获取数据库使用情况信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_USAGE]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库", "数据库大小"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到数据库使用情况相关信息")
		resultWriter.WriteResult("未查询到数据库使用情况相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 注意检查数据库的大小，是否需要清理历史数据。
	`
	resultWriter.WriteResult(suggestion)
}
