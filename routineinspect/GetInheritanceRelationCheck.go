package routineinspect

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetInheritanceRelationCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 先获取所有非template数据库名称
	dbNamesCmd := exec.Command("psql", "--pset=pager=off", "-t", "-A", "-q", "-c", `select datname from pg_database where datname not in ('template0', 'template1')`)
	var dbNamesResult bytes.Buffer
	dbNamesCmd.Stdout = &dbNamesResult
	err := dbNamesCmd.Run()
	if err != nil {
		fmt.Printf("执行获取数据库名称命令失败: %s\n", err)
		return
	}
	dbNames := strings.Split(strings.TrimSpace(dbNamesResult.String()), "\n")

	// 用于存储所有继承关系检查信息的结果
	var allResult bytes.Buffer

	// 遍历每个数据库，获取继承关系信息并合并结果
	for _, db := range dbNames {
		cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select inhrelid::regclass,inhparent::regclass,inhseqno from pg_inherits order by 2,3`)
		var result bytes.Buffer
		cmd.Stdout = &result
		err := cmd.Run()
		if err != nil {
			fmt.Printf("执行获取继承关系信息命令失败（数据库：%s）: %s\n", db, err)
			continue
		}
		allResult.WriteString(result.String())
	}

	// 解析合并后的结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(allResult.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"继承关系表ID", "父表ID", "继承顺序号"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(allResult.String()), "\n")
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
		fmt.Println("未查询到继承关系检查相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 如果使用继承来实现分区表，注意分区表的触发器中逻辑是否正常，对于时间模式的分区表是否需要及时加分区，修改触发器函数。")
	fmt.Println("   建议继承表的权限统一，如果权限不一致，可能导致某些用户查询时权限不足。")
	fmt.Println()
}
