package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetStandbyInfo 用于获取备库信息，并以表格形式展示。
func GetStandbyInfo(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取备库信息...")
	resultWriter.WriteResult("\n###  备库信息:\n")
	// 获取备库信息
	result := ConnectPostgreSQL("[QUERY_STANDBY_INFO]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"用户名", "应用程序名称", "客户端地址"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到备库信息相关信息")
		resultWriter.WriteResult("未查询到备库信息相关信息")
	}
}
