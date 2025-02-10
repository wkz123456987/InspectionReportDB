package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetUsedDataTypeCounts 函数用于获取用户使用的数据类型统计信息，并以表格形式展示，同时输出相关建议。
func GetUsedDataTypeCounts(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取用户使用的数据类型统计信息...")
	resultWriter.WriteResult("\n### 3.5、用户使用了多少种数据类型:\n")
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		logWriter.WriteLog("未查询到有效数据库名称")
		resultWriter.WriteResult("未查询到有效数据库名称")
		return
	}
	dbNames := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbNames[i] = row[0]
	}

	// 用于存储所有数据类型统计信息的结果
	var allResult [][]string

	// 遍历每个数据库，获取数据类型及数量信息并合并结果
	for _, db := range dbNames {
		dataTypeCountsResult := ConnectPostgreSQL("[QUERY_USED_DATA_TYPE_COUNTS]", db)
		if len(dataTypeCountsResult) > 0 {
			allResult = append(allResult, dataTypeCountsResult...)
		}
	}

	// 根据是否有数据决定输出内容
	if len(allResult) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 当前数据库 | 数据类型名称 | 数量 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|------------|--------------|------|"
		resultWriter.WriteResult(separator)

		for _, row := range allResult {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
				row[0], row[1], row[2]))
		}
	} else {
		logWriter.WriteLog("未查询到用户使用的数据类型相关信息")
		resultWriter.WriteResult("未查询到用户使用的数据类型相关信息")
	}

	// 打印建议
	suggestion := "> 关注常用的数据类型，对于使用频率极低的数据类型可考虑是否合理，必要时进行优化调整。"
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
