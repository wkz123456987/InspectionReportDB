package detection

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/ini.v1"
)

func FileSystemUsageCheck() {
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
	RemoteFileSystemUsageCheck(sshConf)
}

// RemoteFileSystemUsageCheck 获取远程文件系统使用情况并展示
func RemoteFileSystemUsageCheck(sshConf SSHConfig) {
	result, err := ExecuteRemoteCommand(sshConf, "df -h")
	if err != nil {
		fmt.Println(err)
		return
	}

	processFileSystemUsageResult(result)
}

func processFileSystemUsageResult(result string) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	fmt.Println("### 远程文件系统使用情况:")
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(false) // 关闭自动格式化表头，避免重复表头
	writer.SetHeader([]string{"文件系统", "总大小", "已用大小", "可用大小", "使用占比", "挂载点"})
	writer.SetAlignment(tablewriter.ALIGN_LEFT)

	// 重新解析结果并添加数据到表格，跳过第一行（表头行）
	for index, line := range lines {
		if index == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			hasData = true
			writer.Append([]string{
				fields[0],
				fields[1],
				fields[2],
				fields[3],
				fields[4],
				fields[5],
			})
		}
	}

	if !hasData {
		fmt.Println("未查询到远程文件系统使用情况相关信息")
		return
	}

	writer.Render()
	fmt.Println(buffer.String())
	fmt.Println("建议: ")
	fmt.Println("   > 注意预留足够的空间给数据库. ")
}
