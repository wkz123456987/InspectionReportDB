package detection

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/ini.v1"
)

func MemoryUsageCheck() {
	cfg, err := ini.Load("database_config.ini")
	if cfg == nil || err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}
	section := cfg.Section("Linux")
	user := section.Key("User").String()
	password := section.Key("Password").String()
	port, err := section.Key("Port").Int()
	if err != nil {
		log.Fatalf("无法转换端口号: %v", err)
	}
	host := section.Key("Host").String()

	sshConf := SSHConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
	RemoteMemoryUsageCheck(sshConf)
}

// RemoteMemoryUsageCheck 获取远程内存使用率并展示
func RemoteMemoryUsageCheck(sshConf SSHConfig) {
	result, err := ExecuteRemoteCommand(sshConf, "free")
	if err != nil {
		fmt.Println(err)
		return
	}

	processMemoryUsageResult(result)
}

func processMemoryUsageResult(result string) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	fmt.Println("### 远程内存使用率:")
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader([]string{"内存使用率"})
	writer.SetAlignment(tablewriter.ALIGN_LEFT)

	// 重新解析结果提取内存使用率数据并添加到表格
	for _, line := range lines {
		if strings.Contains(line, "Mem") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				used := fields[2]
				total := fields[1]
				usageRate := fmt.Sprintf("%.0f%%", float64(atoi(used))/float64(atoi(total))*100)
				writer.Append([]string{usageRate})
				hasData = true
			}
		}
	}

	if !hasData {
		fmt.Println("未查询到远程内存使用率相关信息")
		return
	}

	writer.Render()
	fmt.Println(buffer.String())

	fmt.Println("建议: ")
	fmt.Println("   > 注意检查业务中内存占用高的原因. ")
}

func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
