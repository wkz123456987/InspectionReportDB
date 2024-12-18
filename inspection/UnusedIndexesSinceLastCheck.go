package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// UnusedIndexesSinceLastCheck 查找上次巡检以来未使用或使用较少的索引
func UnusedIndexesSinceLastCheck() {
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
		printUnusedIndexes(db)
	}

	fmt.Println("\n建议:")
	fmt.Println("   > 建议和应用开发人员确认后, 删除不需要的索引.")
	fmt.Println()
}

// printUnusedIndexes 打印指定数据库中未使用或使用较少的索引
func printUnusedIndexes(db string) {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "未使用数量"})

	hasData := false // 标记是否有有效数据行

	// 构建psql命令以获取未使用或使用较少的索引信息
	cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `
SELECT current_database, COUNT(*) 
FROM (
    SELECT current_database(), t2.schemaname, t2.relname, t2.indexrelname, t2.idx_scan, t2.idx_tup_read, t2.idx_tup_fetch, pg_size_pretty(pg_relation_size(indexrelid))
    FROM pg_stat_all_tables t1, pg_stat_all_indexes t2 
    WHERE t1.relid = t2.relid 
    AND t2.idx_scan < 10 
    AND t2.schemaname NOT IN ('pg_toast', 'pg_catalog') 
    AND indexrelid NOT IN (SELECT conindid FROM pg_constraint WHERE contype IN ('p', 'u', 'f')) 
    AND pg_relation_size(indexrelid) > 65536 
    ORDER BY pg_relation_size(indexrelid) DESC
) aa 
GROUP BY current_database 
ORDER BY COUNT(*);`)
	var result bytes.Buffer
	cmd.Stdout = &result
	cmd.Stderr = &bytes.Buffer{} // 用于捕获错误信息

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command for database %s: %s\n", db, err)
		return
	}

	// 使用正则表达式提取每行的数据（这里假设数据格式符合以|分隔的形式，可根据实际调整正则）
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 3 { // 第一个匹配项是完整的匹配项，后面是列的数据
			database := strings.TrimSpace(matches[1])
			unusedCount := strings.TrimSpace(matches[2])

			if database != "" || unusedCount != "" {
				writer.Append([]string{database, unusedCount})
				hasData = true
			}
		}
	}

	if hasData {
		// 渲染并输出格式化的表格
		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Printf("在数据库 %s 中未查询到上次巡检以来未使用或使用较少的索引信息\n", db)
	}
}
