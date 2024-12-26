package detection

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/ini.v1"
)

func DiskIOCheck() {
	cfg, err := ini.Load("database_config.ini")
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
	RemoteDiskIOCheck(sshConf)
}

// RemoteDiskIOCheck 获取远程磁盘IO情况并展示
func RemoteDiskIOCheck(sshConf SSHConfig) {
	firstResult, err := ExecuteRemoteCommand(sshConf, "iostat -mx 1 1")
	if err != nil {
		fmt.Println(err)
		return
	}
	diskDevices := parseDiskDevices(firstResult)
	if len(diskDevices) > 0 {
		fmt.Println("### 远程输出磁盘IO情况:")
		for _, disk := range diskDevices {
			ioResult, err := ExecuteRemoteCommand(sshConf, fmt.Sprintf("iostat -mx 1 2 %s", disk))
			if err != nil {
				fmt.Printf("Failed to execute command for disk %s: %s\n", disk, err)
				continue
			}
			parseAndAppendIOResult(ioResult, disk)
		}
	} else {
		fmt.Println("未查询到远程磁盘IO情况相关信息")
	}

	fmt.Println("建议: ")
	fmt.Println("   > 注意检查IO占用高的原因.")
}

func parseDiskDevices(result string) []string {
	var diskDevices []string
	lines := strings.Split(strings.TrimSpace(result), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 1 && (strings.Contains(line, "sd") || strings.Contains(line, "vd") || strings.Contains(line, "dm")) {
			diskDevices = append(diskDevices, fields[0])
		}
	}
	return diskDevices
}

func parseAndAppendIOResult(result string, disk string) {
	ioLines := strings.Split(strings.TrimSpace(result), "\n")
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader([]string{"设备名", "磁盘IO：%util"})
	writer.SetAlignment(tablewriter.ALIGN_LEFT)

	for index, ioLine := range ioLines {
		fields := strings.Fields(ioLine)
		if len(fields) >= 14 {
			ioUtil := fields[13]
			if index > 0 { // 跳过表头行，只处理数据行
				ioUtil = strings.TrimSpace(strings.TrimPrefix(ioUtil, "%util"))
				if ioUtil != "" {
					writer.Append([]string{disk, ioUtil}) // 将设备名和处理后的%util值添加到表格中
				}
			}
		}
	}
	writer.Render()
	fmt.Println(buffer.String())
}
