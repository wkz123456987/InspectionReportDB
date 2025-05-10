package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetLongTransactionAnd2PCInfo 用于获取长事务和2PC事务相关信息，并以表格形式展示，同时输出相关建议。
func GetLongTransactionAnd2PCInfo(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取长事务和2PC事务相关信息...")
	resultWriter.WriteResult("\n### 3.13、长事务和2PC相关信息:\n")

	// 获取长事务相关信息
	result1 := ConnectPostgreSQL("[QUERY_LONG_TRANSACTION_INFO]")
	resultWriter.WriteResult("\n#### 3.13.1、长事务相关信息:\n")
	if len(result1) > 0 {
		// Markdown 表格的表头
		tableHeader1 := "| 数据库名 | 用户名 | 查询语句 | 事务开始时间 | 事务持续时间 | 查询开始时间 | 查询持续时间 | 状态 |"
		resultWriter.WriteResult(tableHeader1)

		// Markdown 表格的分隔行
		separator1 := "|----------|--------|----------|--------------|--------------|--------------|--------------|------|"
		resultWriter.WriteResult(separator1)

		for _, row := range result1 {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |",
				row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7]))
		}
	} else {
		resultWriter.WriteResult("未查询到事务持续时长（长事务）超过30分钟的相关信息")
	}

	// 获取2PC相关信息
	result2 := ConnectPostgreSQL("[QUERY_2PC_INFO]")
	resultWriter.WriteResult("\n#### 3.13.2、2PC相关信息:\n")
	if len(result2) > 0 {
		// Markdown 表格的表头
		tableHeader2 := "| 2PC事务ID | 2PC事务GID | 开始时间 | 所属用户 | 数据库名 | 2PC持续时间 |"
		resultWriter.WriteResult(tableHeader2)

		// Markdown 表格的分隔行
		separator2 := "|------------|------------|----------|----------|----------|------------|"
		resultWriter.WriteResult(separator2)

		for _, row := range result2 {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |",
				row[0], row[1], row[2], row[3], row[4], row[5]))
		}
	} else {
		resultWriter.WriteResult("未查询到2PC持续时长超过30分钟的相关信息")
	}

	// 打印建议
	suggestion := "> 长事务过程中产生的垃圾，无法回收，建议不要在数据库中运行LONG SQL，或者错开DML高峰时间去运行LONG SQL。2PC事务一定要记得尽快结束掉，否则可能会导致数据库膨胀。"
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
