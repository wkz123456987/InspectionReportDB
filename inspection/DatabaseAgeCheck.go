package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// DatabaseAgeCheck 函数用于检查数据库年龄情况，并以表格形式打印相关信息，同时输出相关建议。
func DatabaseAgeCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始检查数据库年龄情况...")
	// 获取数据库年龄信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_AGE_INFO]")
	if len(result) == 0 {
		logWriter.WriteLog("未查询到数据库年龄相关信息")
		resultWriter.WriteResult("未查询到数据库年龄相关信息")
		return
	}

	// 写入标题
	resultWriter.WriteResult("### 2.8、数据库年龄:\n")

	// Markdown 表格的表头和分隔行
	tableHeader := "| 数据库 | 年龄 | 年龄_剩余 |"
	separator := "|--------|------|----------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range result {
		// 假设row是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
			row[0], row[1], row[2]))
	}

	// 打印建议
	suggestion := "> 数据库的年龄正常情况下应该小于vacuum_freeze_table_age, 如果剩余年龄小于5亿, 建议人为干预, 将LONG SQL或事务杀掉后, 执行vacuum freeze."

	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult(suggestion)
}
