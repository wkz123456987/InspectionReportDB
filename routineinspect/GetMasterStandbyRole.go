package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetMasterStandbyRole 用于获取主备库角色信息，并以表格形式展示。
func GetMasterStandbyRole(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取主备库角色信息...")
	resultWriter.WriteResult("\n###  数据库主备角色:\n")
	// 获取主备库角色信息
	result := ConnectPostgreSQL("[QUERY_MASTER_STANDBY_ROLE]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库主备角色"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到主备库角色相关信息")
		resultWriter.WriteResult("未查询到主备库角色相关信息")
	}
}
