package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// DatabaseStats 获取数据库统计信息,回滚比例, 命中比例, 数据块读写时间, 死锁, 复制冲突:
func DatabaseStats() {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "回滚比例", "命中比例", "数据块读时间", "数据块写时间", "复制冲突", "死锁"})

	hasData := false // 标记是否有有效数据行

	// 构建psql命令以获取数据库统计信息
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `
SELECT 
    datname,
    ROUND(100 * (xact_rollback::numeric / (CASE WHEN xact_commit > 0 THEN xact_commit ELSE 1 END + xact_rollback)), 2) || ' %' AS rollback_ratio,
    ROUND(100 * (blks_hit::numeric / (CASE WHEN blks_read > 0 THEN blks_read ELSE 1 END + blks_hit)), 2) || ' %' AS hit_ratio,
    blk_read_time,
    blk_write_time,
    conflicts,
    deadlocks 
FROM pg_stat_database;`)
	var result bytes.Buffer
	cmd.Stdout = &result
	cmd.Stderr = &bytes.Buffer{} // 用于捕获错误信息

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 使用正则表达式提取每行的数据（可根据实际数据格式调整正则表达式）
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 8 { // 第一个匹配项是完整的匹配项，后面是列的数据
			database := strings.TrimSpace(matches[1])
			rollbackRatio := strings.TrimSpace(matches[2])
			hitRatio := strings.TrimSpace(matches[3])
			readTime := strings.TrimSpace(matches[4])
			writeTime := strings.TrimSpace(matches[5])
			conflicts := strings.TrimSpace(matches[6])
			deadlocks := strings.TrimSpace(matches[7])

			if database != "" || rollbackRatio != "" || hitRatio != "" || readTime != "" || writeTime != "" || conflicts != "" || deadlocks != "" {
				writer.Append([]string{
					database,
					rollbackRatio,
					hitRatio,
					readTime,
					writeTime,
					conflicts,
					deadlocks,
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
		fmt.Println("未查询到相关的数据库统计信息")
	}

	// 打印建议
	fmt.Println("\n建议:")
	fmt.Println("   > 回滚比例大说明业务逻辑可能有问题, 命中率小说明shared_buffer要加大, 数据块读写时间长说明块设备的IO性能要提升, 死锁次数多说明业务逻辑有问题, 复制冲突次数多说明备库可能在跑LONG SQL.")
	fmt.Println()
}
