package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetUserConnectionLimits 用于获取用户连接数限制相关信息，并以表格形式展示，同时输出相关建议。
func GetUserConnectionLimits(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取用户连接数限制相关信息...")
	resultWriter.WriteResult("\n###  用户连接数限制:\n")
	// 获取用户连接数限制相关信息
	result := ConnectPostgreSQL("[QUERY_USER_CONNECTION_LIMITS]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"用户名", "用户连接限制", "当前用户已使用的连接数"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到用户连接数限制相关信息")
		resultWriter.WriteResult("未查询到用户连接数限制相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 给用户设置足够的连接数，使用alter role... CONNECTION LIMIT来设置。
	`
	resultWriter.WriteResult(suggestion)
}
