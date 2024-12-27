package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetDatabaseConnectionLimits 用于获取数据库连接限制相关信息，并以表格形式展示，同时输出相关建议。
func GetDatabaseConnectionLimits(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库连接限制相关信息...")
	resultWriter.WriteResult("\n###  数据库连接数限制:\n")
	// 获取数据库连接限制相关信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_CONNECTION_LIMITS]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库", "数据库连接限制", "数据库已使用连接"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到数据库连接限制相关信息")
		resultWriter.WriteResult("未查询到数据库连接限制相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 给数据库设置足够的连接数，使用alter database... CONNECTION LIMIT来设置。
	`
	resultWriter.WriteResult(suggestion)
}
