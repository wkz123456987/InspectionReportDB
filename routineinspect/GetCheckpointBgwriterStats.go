package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetCheckpointBgwriterStats 用于获取检查点、bgwriter统计信息，并以表格形式展示，同时输出相关建议。
func GetCheckpointBgwriterStats(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取检查点、bgwriter统计信息...")
	resultWriter.WriteResult("\n###  数据库检查点和bgwriter统计信息:\n")
	// 获取检查点、bgwriter统计信息
	result := ConnectPostgreSQL("[QUERY_CHECKPOINT_BGWRITER_STATS]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{
			"checkpoints_timed", "checkpoints_req", "checkpoint_write_time",
			"checkpoint_sync_time", "buffers_checkpoint", "buffers_clean",
			"maxwritten_clean", "buffers_backend", "buffers_backend_fsync",
			"buffers_alloc", "stats_reset",
		})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到检查点、bgwriter统计信息相关信息")
		resultWriter.WriteResult("未查询到检查点、bgwriter统计信息相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 如果检测结果显示checkpoint_write_time多，说明检查点持续时间长，检查点过程中产生了较多的脏页。
        > checkpoint_sync_time代表检查点开始时的shared buffer中的脏页被同步到磁盘的时间，如果时间过长，并且数据库在检查点时性能较差，考虑一下提升块设备的IOPS能力。
        > buffers_backend_fsync太多说明需要加大shared buffer 或者 减小bgwriter_delay参数。
	`
	resultWriter.WriteResult(suggestion)
}
