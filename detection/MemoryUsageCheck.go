package detection

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// MemoryUsageCheck 获取内存使用率并展示
func MemoryUsageCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建free命令获取内存相关信息
	cmd := exec.Command("free")
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 解析结果，判断是否获取到有效数据（这里简单判断是否包含Mem相关文本行）
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Mem") {
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  输出内存使用率:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"内存使用率"})
		writer.SetAlignment(tablewriter.ALIGN_LEFT)

		// 重新解析结果提取内存使用率数据并添加到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Mem") {
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					used := fields[2]
					total := fields[1]
					usageRate := fmt.Sprintf("%.0f%%", float64(atoi(used))/float64(atoi(total))*100)
					writer.Append([]string{usageRate})
				}
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到内存使用率相关信息")
	}

	fmt.Println("建议: ")
	fmt.Println("   > 注意检查业务中内存占用高的原因. ")
}

func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
