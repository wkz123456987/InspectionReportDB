package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// GarbageDataCheck 函数用于检查数据库中垃圾数据情况，并以表格形式打印相关信息。
func GarbageDataCheck() {
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
		// 调用函数处理每个数据库的垃圾数据情况，更新hasData的值
		hasDataForDb := printGarbageDataTable(db)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  垃圾数据:")
	} else {
		fmt.Println("未查询到数据库中垃圾数据相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 通常垃圾过多, 可能是因为无法回收垃圾, 或者回收垃圾的进程繁忙或没有及时唤醒, 或者没有开启autovacuum, 或在短时间内产生了大量的垃圾. ")
	fmt.Println("    可以等待autovacuum进行处理, 或者手工执行vacuum table. ")
	fmt.Println()
}

// printGarbageDataTable 打印指定数据库的垃圾数据情况表格
func printGarbageDataTable(db string) bool {
	// 创建用于当前数据库表格输出的对象并设置表头
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "schema", "表名", "死元组数量"})

	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 构建psql命令以获取垃圾数据信息
	cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select current_database(),schemaname,relname,n_dead_tup from pg_stat_all_tables where n_live_tup>0 and n_dead_tup/n_live_tup>0.2 and schemaname not in ('pg_toast','pg_catalog') order by n_dead_tup desc limit 5`)
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
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 5 { // 第一个匹配项是完整的匹配项，后面是列的数据
			database := strings.TrimSpace(matches[1])
			schema := strings.TrimSpace(matches[2])
			tableName := strings.TrimSpace(matches[3])
			deadTupleCount := strings.TrimSpace(matches[4])

			if database != "" || schema != "" || tableName != "" || deadTupleCount != "" {
				writer.Append([]string{
					database,
					schema,
					tableName,
					deadTupleCount,
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
