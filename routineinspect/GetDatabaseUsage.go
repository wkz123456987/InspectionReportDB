package routineinspect

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetDatabaseUsage() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建psql命令获取数据库使用情况
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select datname,pg_size_pretty(pg_database_size(oid)) from pg_database order by pg_database_size(oid) desc`)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("执行获取数据库使用情况命令失败: %s\n", err)
		return
	}

	// 解析结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 3 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库", "数据库大小"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 3 {
				writer.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到数据库使用情况相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   >  注意检查数据库的大小, 是否需要清理历史数据. ")
	fmt.Println()
}
