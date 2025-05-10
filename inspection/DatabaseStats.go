package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// DatabaseStats 函数用于获取数据库统计信息，包括回滚比例、命中比例、数据块读写时间、死锁、复制冲突，并以表格形式展示，同时输出相关建议。
func DatabaseStats(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库统计信息...")
	resultWriter.WriteResult("\n### 2.5、获取数据库统计信息,回滚比例, 命中比例, 数据块读写时间, 死锁, 复制冲突:\n")

	// 获取数据库统计信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_STATS]")
	if len(result) == 0 {
		logWriter.WriteLog("未查询到相关的数据库统计信息")
		resultWriter.WriteResult("未查询到相关的数据库统计信息")
		return
	}

	// Markdown 表格的表头和分隔行
	tableHeader := "| 数据库 | 回滚比例 | 命中比例 | 数据块读时间 | 数据块写时间 | 复制冲突 | 死锁 |"
	separator := "|--------|----------|----------|--------------|--------------|----------|------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range result {
		// 假设row是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |",
			row[0], row[1], row[2], row[3], row[4], row[5], row[6]))
	}

	// 打印建议
	suggestion := "> 回滚比例大说明业务逻辑可能有问题, 命中率小说明shared_buffer要加大, 数据块读写时间长说明块设备的IO性能要提升, 死锁次数多说明业务逻辑有问题, 复制冲突次数多说明备库可能在跑LONG SQL."
	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult(suggestion)
}
