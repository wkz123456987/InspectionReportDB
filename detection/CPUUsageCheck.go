package detection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// CPUUsageCheck 读取配置文件并执行远程CPU使用率检查
func CPUUsageCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始巡检远程系统CPU使用率...")
	result, err := ExecuteRemoteCommand(GetSSHConfig(logWriter), "top -b -n 1")
	if err != nil {
		logWriter.WriteLog("执行远程命令失败: " + err.Error())
		return
	}
	processCPUUsageResult(result, resultWriter)

}

func processCPUUsageResult(result string, resultWriter *fileutils.ResultWriter) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	// 写入标题
	header := "### 远程系统CPU使用率:\n"
	resultWriter.WriteResult(header)

	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader([]string{"CPU使用率"})
	writer.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, line := range lines {
		if strings.Contains(line, "Cpu") {
			hasData = true
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				cpuUsage := fields[1]
				cpuUsage = strings.TrimSuffix(cpuUsage, "%")
				cpuUsage = strings.ReplaceAll(cpuUsage, "\n", "") + "%"
				writer.Append([]string{cpuUsage})
			}
			break // 假设我们只需要第一行的CPU使用率信息
		}
	}

	if !hasData {
		resultWriter.WriteResult("未查询到远程系统CPU使用率相关信息")
		return
	}

	writer.Render()
	resultWriter.WriteResult(buffer.String())
}
