package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// ReplicationSlotStatus 函数用于检查复制槽状态情况，并以表格形式打印相关信息。
func ReplicationSlotStatus() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建psql命令以获取复制槽状态信息
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select slot_name,slot_type,active from pg_replication_slots order by 3`)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 解析结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（可根据实际数据格式调整正则表达式）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 { // 第一个匹配项是完整的匹配项，后面是列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  复制槽状态:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"复制槽名称", "复制槽类型", "复制槽状态"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 4 {
				writer.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
					strings.TrimSpace(matches[3]),
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到复制槽状态相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 若复制槽状态出现f，要及时处理，保留的 WAL 记录会占用磁盘空间，如果订阅端长时间无法跟上，主数据库的 WAL 文件会堆积，这可能会影响主数据库的性能和磁盘空间使用. ")
	fmt.Println("    请检查是否是否网络问题、服务器资源、数据库日志是否有复制冲突的问题")
	fmt.Println()
}
