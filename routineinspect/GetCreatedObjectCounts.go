package routineinspect

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// GetCreatedObjectCounts 用于获取用户创建的对象及数量信息
func GetCreatedObjectCounts() {
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

	// 用于存储所有对象统计信息的结果
	var allResult bytes.Buffer

	// 遍历每个数据库，获取对象及数量信息并合并结果
	for _, db := range dbNames {
		cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select current_database(),rolname,nspname,relkind,count(*) from pg_class a,pg_authid b,pg_namespace c where a.relnamespace=c.oid and a.relowner=b.oid and nspname!~ '^pg_' and nspname<>'information_schema' group by 1,2,3,4 order by 5 desc`)
		var result bytes.Buffer
		cmd.Stdout = &result
		err := cmd.Run()
		if err != nil {
			fmt.Printf("执行获取对象统计信息命令失败（数据库：%s）: %s\n", db, err)
			continue
		}
		allResult.WriteString(result.String())
	}

	// 解析合并后的结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(allResult.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 6 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"当前数据库", "角色名称", "命名空间名称", "对象类型", "数量"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(allResult.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 6 {
				writer.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
					strings.TrimSpace(matches[3]),
					strings.TrimSpace(matches[4]),
					strings.TrimSpace(matches[5]),
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到用户创建的对象相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 定期查看用户创建对象的情况，对于过多或长期未使用的对象可考虑清理，以优化数据库空间和性能。")
	fmt.Println()
}
