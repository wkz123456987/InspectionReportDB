package detection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/ini.v1"
)

// CPUUsageCheck 读取配置文件并执行远程CPU使用率检查
func CPUUsageCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始巡检远程系统CPU使用率...")
	cfg, err := ini.Load("database_config.ini")
	if cfg == nil || err != nil {
		logWriter.WriteLog("无法读取配置文件: " + err.Error())
		return
	}
	section := cfg.Section("Linux")
	user := section.Key("User").String()
	password := section.Key("Password").String()
	port, err := section.Key("Port").Int()
	if err != nil {
		logWriter.WriteLog("无法转换端口号: " + err.Error())
		return
	}
	host := section.Key("Host").String()
	CPUUsageCheck1(user, password, host, port, logWriter, resultWriter)
}

// CPUUsageCheck1 检查远程Linux系统的CPU使用率
func CPUUsageCheck1(user, password, host string, port int, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	sshConf := SSHConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
	result, err := ExecuteRemoteCommand(sshConf, "top -b -n 1")
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
