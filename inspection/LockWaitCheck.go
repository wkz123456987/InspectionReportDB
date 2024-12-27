package inspection

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// LockWaitCheck 函数用于检查锁等待情况，并以表格形式打印相关信息，同时输出相关建议。
func LockWaitCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始检查锁等待情况...")
	// 获取锁等待信息
	result := ConnectPostgreSQL("[QUERY_LOCK_WAIT_INFO]")
	if len(result) > 0 {
		resultWriter.WriteResult("\n### 锁等待:\n")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{
			"锁类型", "读锁模式", "读锁用户", "读锁数据库", "关联关系", "读锁进程ID", "读锁页面",
			"读锁元组", "读锁事务开始时间", "读锁查询开始时间", "读锁锁定时长", "读锁查询语句",
			"写锁模式", "写锁进程ID", "写锁页面", "写锁元组", "写锁事务开始时间", "写锁查询开始时间", "写锁锁定时长", "写锁查询语句",
		})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到锁等待相关信息")
		resultWriter.WriteResult("未查询到锁等待相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 锁等待状态, 反映业务逻辑的问题或者SQL性能有问题, 建议深入排查持锁的SQL.
	`
	resultWriter.WriteResult(suggestion)
}
