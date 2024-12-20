package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetDatabaseUsage函数用于获取数据库使用情况信息，并以表格形式展示，同时输出相关建议。
func GetDatabaseUsage() {
	// 获取数据库使用情况信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_USAGE]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库", "数据库大小"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到数据库使用情况相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   >  注意检查数据库的大小, 是否需要清理历史数据. ")
	fmt.Println()
}
