package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// SchemaStats 函数用于检查数据库中schema的统计情况，并以表格形式打印相关信息。
func SchemaStats() {
	// 标记是否有数据库存在有效数据，初始化为false
	hasAnyData := false

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
		// 检查当前数据库是否有有效数据
		hasData := printSchemaStatsTable(db)
		if hasData {
			hasAnyData = true
		}
	}

	// 根据整体是否有数据决定输出内容
	if hasAnyData {
		fmt.Println("###  schema统计:")
	} else {
		fmt.Println("未查询到schema统计相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 主要关注pg_catalog的大小，若pg_catalog太大，需要排查是哪个系统表出现膨胀导致的. ")
	fmt.Println()
}

// printSchemaStatsTable 打印指定数据库的schema统计情况表格
func printSchemaStatsTable(db string) bool {
	// 标记当前数据库是否获取到有效数据，初始化为false
	hasData := false

	// 打印当前数据库的schema标题
	fmt.Printf("【%s】库的schema: \n", db)

	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"schemaName", "Byte", "MB", "GB"})

	// 构建psql命令以获取schema统计信息
	cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `SELECT schemaName as "schemaName",sum(total_size) as "Byte",round(sum(total_size)/1024/1024,1) as "MB",round(sum(total_size)/1024/1024/1024,1) as "GB" from (SELECT nspname as schemaName,pg_total_relation_size(pg_class.oid) as total_size FROM pg_class JOIN  pg_namespace ON (pg_namespace.oid = pg_class.relnamespace) WHERE relkind IN ('r', 'v', 'm', 'S', 'f') ORDER BY total_size DESC) as aa group by schemaName order by 4 desc;`)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command for database %s: %s\n", db, err)
		return false
	}

	// 使用正则表达式解析结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 5 {
			hasData = true
			writer.Append([]string{
				strings.TrimSpace(matches[1]),
				strings.TrimSpace(matches[2]),
				strings.TrimSpace(matches[3]),
				strings.TrimSpace(matches[4]),
			})
		}
	}

	if hasData {
		writer.Render()
		fmt.Println(buffer.String())
	}

	return hasData
}
