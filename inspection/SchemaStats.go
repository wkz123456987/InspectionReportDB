package inspection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetSchemaStats 用于获取数据库中schema的统计情况并处理结果展示
func GetSchemaStats(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库中schema的统计情况...")
	resultWriter.WriteResult("\n### schema统计:\n")
	// 标记是否有数据库存在有效数据，初始化为false
	hasAnyData := false

	// 获取数据库列表
	resultDBList := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(resultDBList) > 0 {
		for _, db := range resultDBList {
			if db[0] == "" {
				continue
			}
			// 检查当前数据库是否有有效数据
			hasData := printSchemaStatsTable(db[0], logWriter, resultWriter)
			if hasData {
				hasAnyData = true
			}
		}
	}

	// 根据整体是否有数据决定输出内容
	if hasAnyData {
	} else {
		resultWriter.WriteResult("未查询到schema统计相关信息")
	}
	// 打印建议
	suggestion := `
		建议:
			> 主要关注pg_catalog的大小，若pg_catalog太大，需要排查是哪个系统表出现膨胀导致的.
		`
	resultWriter.WriteResult(suggestion)
}

// printSchemaStatsTable 打印指定数据库的schema统计情况表格
func printSchemaStatsTable(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) bool {
	// 标记当前数据库是否获取到有效数据，初始化为false
	hasData := false

	// 打印当前数据库的schema标题
	resultWriter.WriteResult(fmt.Sprintf("【%s】库的schema: \n", db))

	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"schemaName", "Byte", "MB", "GB"})

	// 获取指定数据库的schema统计信息
	result := ConnectPostgreSQL("[QUERY_SCHEMA_STATS]", db)
	if len(result) > 0 {
		for _, line := range result {
			writer.Append(line)
			hasData = true
		}
	}

	if hasData {
		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog(fmt.Sprintf("数据库 %s 中未查询到schema统计相关信息", db))
		resultWriter.WriteResult(fmt.Sprintf("数据库 %s 中未查询到schema统计相关信息", db))
	}
	return hasData
}
