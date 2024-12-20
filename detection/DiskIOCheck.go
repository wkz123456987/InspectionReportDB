package detection

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// DiskIOCheck 获取磁盘IO情况并展示
func DiskIOCheck() {
	// 执行iostat -mx 1 1命令获取磁盘设备名相关信息
	firstCmd := exec.Command("iostat", "-mx", "1", "1")
	var firstResult bytes.Buffer
	firstCmd.Stdout = &firstResult
	err := firstCmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute first command: %s\n", err)
		return
	}

	// 解析第一步命令结果，获取磁盘设备名列表（简单提取包含特定关键字的行的第一个字段作为设备名）
	var diskDevices []string
	lines := strings.Split(strings.TrimSpace(firstResult.String()), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 1 && (strings.Contains(line, "sd") || strings.Contains(line, "vd") || strings.Contains(line, "dm")) {
			diskDevices = append(diskDevices, fields[0])
		}
	}

	if len(diskDevices) > 0 {
		fmt.Println("### 输出磁盘IO情况:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"磁盘IO：%util"})
		writer.SetAlignment(tablewriter.ALIGN_LEFT)

		// 遍历每个磁盘设备名，执行相应命令获取并处理磁盘IO数据
		for _, disk := range diskDevices {
			// 构建并执行iostat -mx 1 2 [磁盘设备名]命令获取磁盘IO信息
			cmd := exec.Command("iostat", "-mx", "1", "2", disk)
			var result bytes.Buffer
			cmd.Stdout = &result
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Failed to execute command for disk %s: %s\n", disk, err)
				continue
			}

			// 解析磁盘IO命令结果，提取%util字段数据
			ioLines := strings.Split(strings.TrimSpace(result.String()), "\n")
			for _, ioLine := range ioLines {
				fields := strings.Fields(ioLine)
				if len(fields) >= 14 {
					ioUtil := fields[13]
					writer.Append([]string{ioUtil}) // 将%util值添加到表格中
				}
			}
		}

		writer.Render() // 确保表头和数据行被正确渲染
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到磁盘IO情况相关信息")
	}

	fmt.Println("建议: ")
	fmt.Println("   > 注意检查IO占用高的原因.")
}
