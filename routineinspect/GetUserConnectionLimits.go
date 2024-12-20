package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetUserConnectionLimits函数用于获取用户连接数限制相关信息，并以表格形式展示，同时输出相关建议。
func GetUserConnectionLimits() {
	// 获取用户连接数限制相关信息
	result := ConnectPostgreSQL("[QUERY_USER_CONNECTION_LIMITS]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"用户名", "用户连接限制", "当前用户已使用的连接数"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到用户连接数限制相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 给用户设置足够的连接数, alter role... CONNECTION LIMIT.")
	fmt.Println()
}
