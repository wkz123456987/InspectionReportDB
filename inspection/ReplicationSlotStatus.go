package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// ReplicationSlotStatus函数用于检查复制槽状态情况，并以表格形式打印相关信息，同时输出相关建议。
func ReplicationSlotStatus() {

	// 获取复制槽状态信息
	result := ConnectPostgreSQL("[QUERY_REPLICATION_SLOT_STATUS_INFO]")
	if len(result) > 0 {
		fmt.Println("###  复制槽状态:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"复制槽名称", "复制槽类型", "复制槽状态"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到复制槽状态相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 若复制槽状态出现f，要及时处理，保留的 WAL 记录会占用磁盘空间，如果订阅端长时间无法跟上，主数据库的 WAL 文件会堆积，这可能会影响主数据库的性能和磁盘空间使用. ")
	fmt.Println("    请检查是否是否网络问题、服务器资源、数据库日志是否有复制冲突的问题")
	fmt.Println()
}
