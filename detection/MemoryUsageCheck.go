package detection

import (
	"GoBasic/utils/fileutils"
	"fmt"
	"strconv"
	"strings"
)

func MemoryUsageCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始巡检远程内存使用率...")
	RemoteMemoryUsageCheck(GetSSHConfig(logWriter), logWriter, resultWriter)
}

// RemoteMemoryUsageCheck 获取远程内存使用率并展示
func RemoteMemoryUsageCheck(sshConf SSHConfig, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	result, err := ExecuteRemoteCommand(sshConf, "free")
	if err != nil {
		logWriter.WriteLog("执行远程命令失败: " + err.Error())
		return
	}

	processMemoryUsageResult(result, resultWriter)
}

func processMemoryUsageResult(result string, resultWriter *fileutils.ResultWriter) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	// 写入标题
	header := "### 1.2、内存使用率:\n"
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	memoryUsageHeader := "| 内存使用率 |\n|------------|"
	resultWriter.WriteResult(memoryUsageHeader)

	// 重新解析结果提取内存使用率数据并添加到表格
	for _, line := range lines {
		if strings.Contains(line, "Mem") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				used := fields[2]
				total := fields[1]
				usageRate := fmt.Sprintf("%.0f%%", float64(atoi(used))/float64(atoi(total))*100)
				// Markdown 表格的单元格
				resultWriter.WriteResult(fmt.Sprintf("| %s |\n", usageRate))
				hasData = true
			}
		}
	}

	if !hasData {
		resultWriter.WriteResult("未查询到远程内存使用率相关信息")
		return
	}

	// 写入建议
	suggestion := "**建议:** \n   > 注意检查业务中内存占用高的原因. "
	resultWriter.WriteResult(suggestion)
}

func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
