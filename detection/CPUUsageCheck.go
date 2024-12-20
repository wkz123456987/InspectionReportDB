package detection

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// CPUUsageCheck 获取CPU使用率并展示
func CPUUsageCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建top -b -n 1命令获取CPU使用率相关信息
	cmd := exec.Command("top", "-b", "-n", "1")
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 解析结果，判断是否获取到有效数据（这里简单判断是否包含相关CPU使用率的文本行）
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cpu") {
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  CPU使用率:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"CPU使用率"})
		writer.SetAlignment(tablewriter.ALIGN_LEFT)

		// 重新解析结果提取CPU使用率数据并添加到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Cpu") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					cpuUsage := fields[1]
					cpuUsage = strings.TrimSuffix(cpuUsage, "%")
					cpuUsage = strings.ReplaceAll(cpuUsage, "\n", "") + "%"
					writer.Append([]string{cpuUsage})
				}
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到CPU使用率相关信息")
	}
}
