package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetDBVersion 函数用于获取数据库版本信息，并进行展示（这里简单打印版本信息，可按需调整展示方式）
func GetDBVersion(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库版本信息...")
	result := ConnectPostgreSQL("[QUERY_DB_VERSION]")
	if len(result) > 0 {
		// 因为查询版本信息一般是单行数据，这里取第一行第一列作为版本内容
		version := result[0][0]
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库版本"})
		writer.Append([]string{version})
		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到有效数据")
		resultWriter.WriteResult("未查询到有效数据")
	}
}
