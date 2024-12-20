package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetUsedDataTypeCounts函数用于获取用户使用的数据类型统计信息，并以表格形式展示，同时输出相关建议。
func GetUsedDataTypeCounts() {
	// 先获取所有非template数据库名称
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		fmt.Println("未查询到有效数据库名称")
		return
	}
	dbNames := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbNames[i] = row[0]
	}

	// 用于存储所有数据类型统计信息的结果
	var allResult [][]string

	// 遍历每个数据库，获取数据类型及数量信息并合并结果
	for _, db := range dbNames {
		dataTypeCountsResult := ConnectPostgreSQL("[QUERY_USED_DATA_TYPE_COUNTS]", db)
		if len(dataTypeCountsResult) > 0 {
			allResult = append(allResult, dataTypeCountsResult...)
		}
	}

	// 根据是否有数据决定输出内容
	if len(allResult) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"当前数据库", "数据类型名称", "数量"})

		for _, row := range allResult {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到用户使用的数据类型相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 关注常用的数据类型，对于使用频率极低的数据类型可考虑是否合理，必要时进行优化调整。")
	fmt.Println()
}
