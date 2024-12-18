package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// DatabaseAgeCheck 函数用于检查数据库年龄情况，并以表格形式打印相关信息。
func DatabaseAgeCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建psql命令以获取数据库年龄信息
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select datname,age(datfrozenxid),2^31-age(datfrozenxid) age_remain from pg_database order by age(datfrozenxid) desc`)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 解析结果
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（可根据实际数据格式调整正则表达式）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 { // 第一个匹配项是完整的匹配项，后面是列的数据
			database := strings.TrimSpace(matches[1])
			age := strings.TrimSpace(matches[2])
			ageRemain := strings.TrimSpace(matches[3])

			if database != "" || age != "" || ageRemain != "" {
				hasData = true
				break
			}
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  数据库年龄:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库", "年龄", "年龄_剩余"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 4 {
				database := strings.TrimSpace(matches[1])
				age := strings.TrimSpace(matches[2])
				ageRemain := strings.TrimSpace(matches[3])

				writer.Append([]string{
					database,
					age,
					ageRemain,
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到数据库年龄相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 数据库的年龄正常情况下应该小于vacuum_freeze_table_age, 如果剩余年龄小于5亿, 建议人为干预, 将LONG SQL或事务杀掉后, 执行vacuum freeze. ")
	fmt.Println()
}
