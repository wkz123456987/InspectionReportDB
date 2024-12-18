package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// TableAgeCheck 函数用于检查表年龄情况，并以表格形式打印相关信息。
func TableAgeCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 执行psql命令获取数据库列表
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "-A", "-q", "-c", "SELECT datname FROM pg_database WHERE datname NOT IN ('template0', 'template1')")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 解析数据库列表并遍历
	dbList := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, db := range dbList {
		if db == "" {
			continue
		}
		// 调用函数处理每个数据库的表年龄情况，更新hasData的值
		hasDataForDb := printTableAgeTable(db)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  表年龄:")
	} else {
		fmt.Println("未查询到表年龄相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 表的年龄正常情况下应该小于vacuum_freeze_table_age, 如果剩余年龄小于5亿, 建议人为干预, 将LONG SQL或事务杀掉后, 执行vacuum freeze. ")
	fmt.Println()
}

// printTableAgeTable 打印指定数据库的表年龄情况表格
func printTableAgeTable(db string) bool {
	// 创建用于当前数据库表格输出的对象并设置表头
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "rolname", "nspname", "relkind", "表名", "年龄", "年龄_剩余"})

	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 构建psql命令以获取表年龄信息
	cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select current_database(),rolname,nspname,relkind,relname,age(relfrozenxid),2^31-age(relfrozenxid) age_remain from pg_authid t1 join pg_class t2 on t1.oid=t2.relowner join pg_namespace t3 on t2.relnamespace=t3.oid where t2.relkind in ('t','r') order by age(relfrozenxid) desc limit 5`)
	var result bytes.Buffer
	cmd.Stdout = &result
	cmd.Stderr = &bytes.Buffer{} // 用于捕获错误信息

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command for database %s: %s\n", db, err)
		return false
	}

	// 使用正则表达式提取每行的数据（可根据实际数据格式调整正则表达式）
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 8 { // 第一个匹配项是完整的匹配项，后面是列的数据
			database := strings.TrimSpace(matches[1])
			rolname := strings.TrimSpace(matches[2])
			nspname := strings.TrimSpace(matches[3])
			relkind := strings.TrimSpace(matches[4])
			tableName := strings.TrimSpace(matches[5])
			age := strings.TrimSpace(matches[6])
			ageRemain := strings.TrimSpace(matches[7])

			if database != "" || rolname != "" || nspname != "" || relkind != "" || tableName != "" || age != "" || ageRemain != "" {
				writer.Append([]string{
					database,
					rolname,
					nspname,
					relkind,
					tableName,
					age,
					ageRemain,
				})
				currentHasData = true
			}
		}
	}

	if currentHasData {
		writer.Render()
		fmt.Println(buffer.String())
	}

	return currentHasData
}
