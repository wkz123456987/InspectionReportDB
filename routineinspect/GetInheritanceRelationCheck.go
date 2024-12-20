package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetInheritanceRelationCheck函数用于获取继承关系检查相关信息，并以表格形式展示，同时输出相关建议。
func GetInheritanceRelationCheck() {
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

	// 用于存储所有继承关系检查信息的结果
	var allResult [][]string

	// 遍历每个数据库，获取继承关系信息并合并结果
	for _, db := range dbNames {
		inheritanceRelationInfoResult := ConnectPostgreSQL("[QUERY_INHERITANCE_RELATION_INFO]", db)
		if len(inheritanceRelationInfoResult) > 0 {
			allResult = append(allResult, inheritanceRelationInfoResult...)
		}
	}

	// 根据是否有数据决定输出内容
	if len(allResult) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"继承关系表ID", "父表ID", "继承顺序号"})

		for _, row := range allResult {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到继承关系检查相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 如果使用继承来实现分区表，注意分区表的触发器中逻辑是否正常，对于时间模式的分区表是否需要及时加分区，修改触发器函数。")
	fmt.Println("   建议继承表的权限统一，如果权限不一致，可能导致某些用户查询时权限不足。")
	fmt.Println()
}
