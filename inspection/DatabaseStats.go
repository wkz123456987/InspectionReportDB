package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// DatabaseStats 获取数据库统计信息,回滚比例, 命中比例, 数据块读写时间, 死锁, 复制冲突:
func DatabaseStats() {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "回滚比例", "命中比例", "数据块读时间", "数据块写时间", "复制冲突", "死锁"})

	// 获取数据库统计信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_STATS]")
	if len(result) > 0 {
		for _, row := range result {
			writer.Append(row)
		}
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
