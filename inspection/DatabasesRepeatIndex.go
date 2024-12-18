package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// DatabasesRepeatIndex 函数用于检查数据库中重复创建的索引，并以表格形式打印相关信息。
func DatabasesRepeatIndex() {
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
		printRepeatIndexTable(db)
	}

	fmt.Println("\n建议:")
	fmt.Println("   > 当创建重复索引后，不会对数据库的性能产生优化作用，反而会产生一些维护上的成本，请删除重复索引")
}

// printRepeatIndexTable 打印指定数据库的重复索引表格
func printRepeatIndexTable(db string) {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"表名", "索引名"})

	hasData := false // 用于标记是否有有效数据行

	// 构建psql命令以获取重复索引信息
	cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", "SELECT indrelid::regclass AS TableName,array_agg(indexrelid::regclass) AS Indexes FROM pg_index GROUP BY indrelid,indkey HAVING COUNT(*) > 1;")
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
		// 正则表达式用于匹配每一列（这里根据实际数据格式调整正则表达式）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 3 { // 第一个匹配项是完整的匹配项，后面是列的数据
			tableName := strings.TrimSpace(matches[1])
			indexName := strings.TrimSpace(matches[2])
			if tableName != "" || indexName != "" {
				writer.Append([]string{tableName, indexName})
				hasData = true
			}
		}
	}

	if hasData {
		// 渲染并输出格式化的表格
		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("当前数据库中未检测到重复创建的索引信息")
	}
}
