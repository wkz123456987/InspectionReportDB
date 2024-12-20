package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// UnusedIndexesSinceLastCheck函数用于获取各个数据库中未使用或使用较少的索引信息，并以表格形式展示，同时输出相关建议。
func UnusedIndexesSinceLastCheck() {
	// 获取非template数据库名称
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		fmt.Println("未查询到有效数据库名称")
		return
	}
	dbList := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbList[i] = row[0]
	}

	for _, db := range dbList {
		if db == "" {
			continue
		}
		printUnusedIndexes(db)
	}

	fmt.Println("\n建议:")
	fmt.Println("   > 建议和应用开发人员确认后, 删除不需要的索引.")
	fmt.Println()
}

// printUnusedIndexes 打印指定数据库中未使用或使用较少的索引
func printUnusedIndexes(db string) {
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "未使用数量"})

	// 获取指定数据库中未使用或使用较少的索引信息
	unusedIndexesInfoResult := ConnectPostgreSQL("[QUERY_UNUSED_INDEXES_INFO]", db)
	if len(unusedIndexesInfoResult) > 0 {
		for _, row := range unusedIndexesInfoResult {
			writer.Append(row)
		}
		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Printf("在数据库 %s 中未查询到上次巡检以来未使用或使用较少的索引信息\n", db)
	}
}
