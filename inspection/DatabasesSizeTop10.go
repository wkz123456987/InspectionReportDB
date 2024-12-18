package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// DatabasesTop10 函数用于检查数据库中最大的表，并以表格形式打印相关信息。
func DatabasesTop10() {
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
		printTable(db)
	}

	fmt.Println("\n建议:")
	fmt.Println("   > 经验值: 单表超过8GB, 并且这个表需要频繁更新 或 删除+插入的话, 建议对表根据业务逻辑进行合理拆分后获得更好的性能, 以及便于对膨胀索引进行维护; 如果是只读的表, 建议适当结合SQL语句进行优化.")
}

// printTable 打印指定数据库的表格
func printTable(db string) {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "模式", "表名", "类型", "大小"})

	hasData := false // 标记是否有有效数据行

	// 构建psql命令以获取表大小信息
	cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", "SELECT current_database(), b.nspname, c.relname, c.relkind, pg_size_pretty(pg_relation_size(c.oid)) FROM pg_stat_all_tables a, pg_class c, pg_namespace b WHERE pg_relation_size(c.oid) >= 10 AND c.relnamespace = b.oid AND c.relkind = 'r' AND a.relid = c.oid ORDER BY pg_relation_size(c.oid) DESC LIMIT 10")
	var result bytes.Buffer
	cmd.Stdout = &result
	cmd.Stderr = &bytes.Buffer{} // 用于捕获错误信息

	err := cmd.Run()
	if err != nil {
		// 打印错误信息
		fmt.Printf("Failed to execute command for database %s: %s\n", db, err)
		return
	}

	// 使用正则表达式提取每行的数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		// 正则表达式用于匹配每一列
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 6 { // 第一个匹配项是完整的匹配项，后面是列的数据
			database := strings.TrimSpace(matches[1])
			schema := strings.TrimSpace(matches[2])
			tableName := strings.TrimSpace(matches[3])
			tableType := strings.TrimSpace(matches[4])
			tableSize := strings.TrimSpace(matches[5])

			if database != "" || schema != "" || tableName != "" || tableType != "" || tableSize != "" {
				writer.Append([]string{
					database,
					schema,
					tableName,
					tableType,
					tableSize,
				})
				hasData = true
			}
		}
	}

	if hasData {
		// 渲染并输出格式化的表格
		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Printf("在数据库 %s 中未查询到符合条件的表信息\n", db)
	}
}
