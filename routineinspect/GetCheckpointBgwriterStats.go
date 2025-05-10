package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetCheckpointBgwriterStats 用于获取检查点、bgwriter统计信息，并以表格形式展示，同时输出相关建议。
func GetCheckpointBgwriterStats(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取检查点、bgwriter统计信息...")
	resultWriter.WriteResult("\n### 3.12、数据库检查点和bgwriter统计信息:\n")

	// 获取检查点、bgwriter统计信息
	result := ConnectPostgreSQL("[QUERY_CHECKPOINT_BGWRITER_STATS]")
	if len(result) > 0 {
		// Markdown 表格的表头
		tableHeader := "| checkpoints_timed | checkpoints_req | checkpoint_write_time | checkpoint_sync_time | buffers_checkpoint | buffers_clean | maxwritten_clean | buffers_backend | buffers_backend_fsync | buffers_alloc | stats_reset |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|-----------------------|-------------------|-----------------------|---------------------|----------------------|---------------|---------------------|-----------------|-----------------------|---------------|--------------|"
		resultWriter.WriteResult(separator)

		for _, row := range result {
			// 确保row切片中有11个元素（最后一个元素是日期时间类型，需要特殊处理）
			if len(row) == 11 {
				// 写入行数据
				resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |",
					row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], row[10]))
			}
		}
	} else {
		logWriter.WriteLog("未查询到检查点、bgwriter统计信息相关信息")
		resultWriter.WriteResult("未查询到检查点、bgwriter统计信息相关信息")
	}

	// 打印建议
	suggestion := "> - 如果检测结果显示checkpoint_write_time多，说明检查点持续时间长，检查点过程中产生了较多的脏页。\n" +
		" > - checkpoint_sync_time代表检查点开始时的shared buffer中的脏页被同步到磁盘的时间，如果时间过长，并且数据库在检查点时性能较差，考虑一下提升块设备的IOPS能力。\n" +
		" > - buffers_backend_fsync太多说明需要加大shared buffer 或者 减小bgwriter_delay参数。"
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
