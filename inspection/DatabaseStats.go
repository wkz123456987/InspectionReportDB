package inspection

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// DatabaseStats 函数用于获取数据库统计信息，包括回滚比例、命中比例、数据块读写时间、死锁、复制冲突，并以表格形式展示，同时输出相关建议。
func DatabaseStats(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取数据库统计信息...")
	resultWriter.WriteResult("\n###  获取数据库统计信息,回滚比例, 命中比例, 数据块读写时间, 死锁, 复制冲突:\n")
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "回滚比例", "命中比例", "数据块读时间", "数据块写时间", "复制冲突", "死锁"})

	// 获取数据库统计信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_STATS]")
	if len(result) > 0 {
		for _, row := range result {
			writer.Append(row)
		}
		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到相关的数据库统计信息")
		resultWriter.WriteResult("未查询到相关的数据库统计信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 回滚比例大说明业务逻辑可能有问题, 命中率小说明shared_buffer要加大, 数据块读写时间长说明块设备的IO性能要提升, 死锁次数多说明业务逻辑有问题, 复制冲突次数多说明备库可能在跑LONG SQL.
	`
	resultWriter.WriteResult(suggestion)
}
