package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetArchiveAndAutoVacuumSettings 用于获取是否开启归档、自动垃圾回收相关设置信息，并以表格形式展示，同时输出相关建议。
func GetArchiveAndAutoVacuumSettings(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取是否开启归档、自动垃圾回收相关设置信息...")
	resultWriter.WriteResult("\n###  是否开启归档、自动垃圾回收相关设置信息:\n")
	// 获取是否开启归档、自动垃圾回收设置信息
	result := ConnectPostgreSQL("[QUERY_ARCHIVE_AND_AUTOVACUUM_SETTINGS]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"名称", "设置值"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到是否开启归档、自动垃圾回收相关设置信息")
		resultWriter.WriteResult("未查询到是否开启归档、自动垃圾回收相关设置信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 如果当前的wal文件和最后一个归档失败的wal文件之间相差很多个文件，建议尽快排查归档失败的原因，以便修复，否则pg_wal目录可能会撑爆。
	`
	resultWriter.WriteResult(suggestion)
}
