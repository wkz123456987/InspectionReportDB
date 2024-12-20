package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetDatabaseConnectionLimits函数用于获取数据库连接限制相关信息，并以表格形式展示，同时输出相关建议。
func GetDatabaseConnectionLimits() {
	// 获取数据库连接限制相关信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_CONNECTION_LIMITS]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库", "数据库连接限制", "数据库已使用连接"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到数据库连接限制相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 给数据库设置足够的连接数, alter database... CONNECTION LIMIT.")
	fmt.Println()
}
