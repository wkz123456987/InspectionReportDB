package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

func GetUserObjectSpaceInfo() {
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

	// 用于存储所有用户对象空间信息结果
	var allResult [][]string

	// 遍历每个数据库，获取用户对象占用空间信息并合并结果
	for _, db := range dbNames {
		userObjectSpaceInfoResult := ConnectPostgreSQL("[QUERY_USER_OBJECT_SPACE_INFO]", db)
		if len(userObjectSpaceInfoResult) > 0 {
			allResult = append(allResult, userObjectSpaceInfoResult...)
		}
	}

	// 根据是否有数据决定输出内容
	if len(allResult) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"当前数据库", "桶编号", "此桶中关系数量", "桶最小值（格式化后）", "桶最大值（格式化后）"})

		for _, row := range allResult {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到用户对象占用空间相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 关注用户对象占用空间情况，对于占用空间较大的对象可考虑优化存储结构或进行归档处理，以节省数据库空间。")
	fmt.Println()
}
