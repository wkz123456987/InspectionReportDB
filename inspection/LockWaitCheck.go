package inspection

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// LockWaitCheck 函数用于检查锁等待情况，并以表格形式打印相关信息，同时输出相关建议。
func LockWaitCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始检查锁等待情况...")
	resultWriter.WriteResult("\n### 2.10、锁等待:\n")

	// 获取锁等待信息
	result := ConnectPostgreSQL("[QUERY_LOCK_WAIT_INFO]")
	if len(result) == 0 {
		logWriter.WriteLog("未查询到锁等待相关信息")
		resultWriter.WriteResult("未查询到锁等待相关信息")
		return
	}

	// Markdown 表格的表头和分隔行
	tableHeader := "| 锁类型 | 读锁模式 | 读锁用户 | 读锁数据库 | 关联关系 | 读锁进程ID | 读锁页面 | 读锁元组 | 读锁事务开始时间 | 读锁查询开始时间 | 读锁锁定时长 | 读锁查询语句 | 写锁模式 | 写锁进程ID | 写锁页面 | 写锁元组 | 写锁事务开始时间 | 写锁查询开始时间 | 写锁锁定时长 | 写锁查询语句 |"
	separator := "|--------|----------|----------|------------|------------|------------|----------|----------|--------------------|--------------------|----------------|----------------|----------|------------|----------|----------|--------------------|--------------------|----------------|----------------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	// 遍历结果并添加数据到Markdown表格中
	for _, row := range result {
		// 假设row是一个包含所需字段的切片
		resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13], row[14], row[15], row[16], row[17], row[18], row[19], row[20], row[21]))
	}

	// 打印建议
	suggestion := "> 锁等待状态, 反映业务逻辑的问题或者SQL性能有问题, 建议深入排查持锁的SQL."
	resultWriter.WriteResult("\n**建议:**\n")
	resultWriter.WriteResult(suggestion)
}
