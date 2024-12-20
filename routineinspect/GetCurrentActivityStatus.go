package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetCurrentActivityStatus 函数用于获取数据库当前活跃度状态信息，并以表格形式打印相关信息，同时输出相关建议。
func GetCurrentActivityStatus() {
	result := ConnectPostgreSQL("[QUERY_ACTIVITY_STATUS]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetHeader([]string{"当前时间", "状态", "count"})
		for _, row := range result {
			writer.Append(row)
		}
		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到有效数据")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 如果active状态很多, 说明数据库比较繁忙. 如果idle in transaction很多, 说明业务逻辑设计可能有问题. 如果idle很多, 可能使用了连接池, 并且可能没有自动回收连接到连接池的最小连接数. ")
	fmt.Println()
}
