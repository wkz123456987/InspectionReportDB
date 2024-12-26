package detection

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/ini.v1"
)

func CPUUsageCheck() {
	// 读取配置文件获取数据库配置信息
	cfg, err := ini.Load("database_config.ini")
	if err != nil || cfg == nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}
	section := cfg.Section("Linux")
	user := section.Key("User").String()
	password := section.Key("Password").String()
	port, err := section.Key("Port").Int()
	if err != nil {
		log.Fatalf("无法转换端口号: %v", err)
	}
	host := section.Key("Host").String()
	CPUUsageCheck1(user, password, host, port)
}

// CPUUsageCheck 检查远程Linux系统的CPU使用率
func CPUUsageCheck1(user, password, host string, port int) {
	sshConf := SSHConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
	result, err := ExecuteRemoteCommand(sshConf, "top -b -n 1")
	if err != nil {
		fmt.Println(err)
		return
	}

	processCPUUsageResult(result)
}

func processCPUUsageResult(result string) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	fmt.Println("### 远程系统CPU使用率:")
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
		fmt.Println("未查询到远程系统CPU使用率相关信息")
		return
	}

	writer.Render()
	fmt.Println(buffer.String())
}
