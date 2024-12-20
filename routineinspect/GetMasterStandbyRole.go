package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetMasterStandbyRole函数用于获取主备库角色信息，并以表格形式展示。
func GetMasterStandbyRole() {
	// 获取主备库角色信息
	result := ConnectPostgreSQL("[QUERY_MASTER_STANDBY_ROLE]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库主备角色"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到主备库角色相关信息")
	}
}
