package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetCreatedObjectCounts 用于获取用户创建的对象及数量信息
func GetCreatedObjectCounts() {
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

	// 用于存储所有对象统计信息的结果
	var allResult [][]string

	// 遍历每个数据库，获取对象及数量信息并合并结果
	for _, db := range dbNames {
		objectCountsResult := ConnectPostgreSQL("[QUERY_CREATED_OBJECT_COUNTS]", db)
		if len(objectCountsResult) > 0 {
			allResult = append(allResult, objectCountsResult...)
		}
	}

	// 根据是否有数据决定输出内容
	if len(allResult) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"当前数据库", "角色名称", "命名空间名称", "对象类型", "数量"})

		for _, row := range allResult {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到用户创建的对象相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 定期查看用户创建对象的情况，对于过多或长期未使用的对象可考虑清理，以优化数据库空间和性能。")
	fmt.Println()
}
