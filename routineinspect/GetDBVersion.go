package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetDBVersion函数用于获取数据库版本信息，并进行展示（这里简单打印版本信息，可按需调整展示方式）
func GetDBVersion() {
	result := ConnectPostgreSQL("[QUERY_DB_VERSION]")
	if len(result) > 0 {
		// 因为查询版本信息一般是单行数据，这里取第一行第一列作为版本内容
		version := result[0][0]
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetHeader([]string{"数据库版本"})
		writer.Append([]string{version})
		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到有效数据")
	}
}
