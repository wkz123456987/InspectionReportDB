package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetLongTransactionAnd2PCInfo 用于获取长事务和2PC事务相关信息，并以表格形式展示，同时输出相关建议。
func GetLongTransactionAnd2PCInfo(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取长事务和2PC事务相关信息...")
	resultWriter.WriteResult("\n### 长事务和2PC相关信息:\n")

	// 获取长事务相关信息
	result1 := ConnectPostgreSQL("[QUERY_LONG_TRANSACTION_INFO]")
	resultWriter.WriteResult("\n##### 长事务相关信息:\n")
	if len(result1) > 0 {
		buffer1 := &bytes.Buffer{}
		writer1 := tablewriter.NewWriter(buffer1)
		writer1.SetAutoFormatHeaders(true)
		writer1.SetHeader([]string{"数据库名", "用户名", "查询语句", "事务开始时间", "事务持续时间", "查询开始时间", "查询持续时间", "状态"})

		for _, row := range result1 {
			writer1.Append(row)
		}

		writer1.Render()
		resultWriter.WriteResult(buffer1.String())
	} else {
		resultWriter.WriteResult("未查询到事务持续时长（长事务）超过30分钟的相关信息")
	}

	// 获取2PC相关信息
	result2 := ConnectPostgreSQL("[QUERY_2PC_INFO]")
	resultWriter.WriteResult("\n#### 2PC相关信息:\n")
	if len(result2) > 0 {
		buffer2 := &bytes.Buffer{}
		writer2 := tablewriter.NewWriter(buffer2)
		writer2.SetAutoFormatHeaders(true)
		writer2.SetHeader([]string{"2PC事务ID", "2PC事务GID", "开始时间", "所属用户", "数据库名", "2PC持续时间"})
		for _, row := range result2 {
			writer2.Append(row)
		}

		writer2.Render()
		resultWriter.WriteResult(buffer2.String())
	} else {
		resultWriter.WriteResult("未查询到2PC持续时长超过30分钟的相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 长事务过程中产生的垃圾，无法回收，建议不要在数据库中运行LONG SQL，或者错开DML高峰时间去运行LONG SQL。2PC事务一定要记得尽快结束掉，否则可能会导致数据库膨胀。
	`
	resultWriter.WriteResult(suggestion)
}
