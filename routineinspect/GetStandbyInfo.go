package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetStandbyInfo函数用于获取备库信息，并以表格形式展示。
func GetStandbyInfo() {
	// 获取备库信息
	result := ConnectPostgreSQL("[QUERY_STANDBY_INFO]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"用户名", "应用程序名称", "客户端地址"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到备库信息相关信息")
	}
}
