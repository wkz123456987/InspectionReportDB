package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// TablesWithTooManyIndexes 查找索引数超过4并且SIZE大于10MB的表
func TablesWithTooManyIndexes() {
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
		printTablesWithTooManyIndexes(db)
	}

	fmt.Println("\n建议:")
	fmt.Println("   > 索引数量太多, 影响表的增删改性能, 建议检查是否有不需要的索引.")
	fmt.Println()
}

// printTablesWithTooManyIndexes 打印指定数据库中索引数超过4并且SIZE大于10MB的表
func printTablesWithTooManyIndexes(db string) {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "模式", "表名", "表大小", "索引数量"})

	hasData := false // 标记是否有有效数据行

	// 构建psql命令以获取表大小和索引数量信息
	cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `
SELECT current_database(), t2.nspname, t1.relname, pg_size_pretty(pg_relation_size(t1.oid)), t3.idx_cnt 
FROM pg_class t1, pg_namespace t2, 
     (SELECT indrelid, COUNT(*) idx_cnt 
      FROM pg_index 
      GROUP BY 1 
      HAVING COUNT(*) > 4) t3 
WHERE pg_relation_size(t1.oid) >= 10000000 
  AND t1.oid = t3.indrelid 
  AND t1.relnamespace = t2.oid 
  AND pg_relation_size(t1.oid) / 1024 / 1024.0 > 10 
ORDER BY t3.idx_cnt DESC;`)
	var result bytes.Buffer
	cmd.Stdout = &result
	cmd.Stderr = &bytes.Buffer{} // 用于捕获错误信息

	err := cmd.Run()
	if err != nil {
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
			tableSize := strings.TrimSpace(matches[4])
			indexCount := strings.TrimSpace(matches[5])

			if database != "" || schema != "" || tableName != "" || tableSize != "" || indexCount != "" {
				writer.Append([]string{
					database,
					schema,
					tableName,
					tableSize,
					indexCount,
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
		fmt.Printf("在数据库 %s 中未查询到索引数超过4且SIZE大于10MB的表信息\n", db)
	}
}
