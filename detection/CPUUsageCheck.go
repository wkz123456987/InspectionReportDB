package detection

import (
	"GoBasic/utils/fileutils"
	"fmt"
	"strings"
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

// processCPUUsageResult 处理CPU使用率结果
func processCPUUsageResult(result string, resultWriter *fileutils.ResultWriter) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	// 写入标题
	header := "### 1.1、系统CPU使用率:\n"
	resultWriter.WriteResult(header)
	// Markdown 表格的表头
	resultWriter.WriteResult("| CPU使用率 |\n|-----------|")
	for _, line := range lines {
		if strings.Contains(line, "Cpu") {
			hasData = true
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				cpuUsage := fields[1]
				cpuUsage = strings.TrimSuffix(cpuUsage, "%")
				cpuUsage = strings.ReplaceAll(cpuUsage, "\n", "") + "%"
				// Markdown 表格的单元格
				resultWriter.WriteResult(fmt.Sprintf("| %s |\n", cpuUsage))
			}
			break // 假设我们只需要第一行的CPU使用率信息
		}
	}
	if !hasData {
		resultWriter.WriteResult("未查询到远程系统CPU使用率相关信息")
	}
}
