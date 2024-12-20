package detection

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// FileSystemInodeUsageCheck 获取文件系统Inode使用情况并展示
func FileSystemInodeUsageCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建df -ih命令获取文件系统Inode使用情况信息
	cmd := exec.Command("df", "-ih")
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 解析结果，简单判断是否有非空数据行，可根据实际格式更精细判断
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  输出文件系统Inode使用情况:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"文件系统", "容量", "已使用", "剩余", "使用占比", "挂载路径"})
		writer.SetAlignment(tablewriter.ALIGN_LEFT)

		// 重新解析结果并添加数据到表格，跳过第一行（表头行）
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for index, line := range lines {
			if index == 0 {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 6 {
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

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到文件系统Inode使用情况相关信息")
	}

	fmt.Println("建议: ")
	fmt.Println("   > 注意预留足够的空间给数据库. ")
}
