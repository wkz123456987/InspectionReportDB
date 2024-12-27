package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// CheckDBConnections 函数用于获取数据库连接相关信息，并以表格形式展示，同时输出相关建议。
func CheckDBConnections(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库连接相关信息...")
	resultWriter.WriteResult("\n### 数据库连接信息:\n")
	result := ConnectPostgreSQL("[QUERY_DB_CONNECTIONS]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"总连接", "已使用连接", "剩余给超级用户连接", "剩余给普通用户连接"})

		writer.Append(result[0])
		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到有效数据")
		resultWriter.WriteResult("未查询到有效数据")
	}

	// 输出建议内容
	suggestion := `
    建议:
        > 给超级用户和普通用户设置足够的连接, 以免不能登录数据库.
	`
	resultWriter.WriteResult(suggestion)
}
