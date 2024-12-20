package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// IndexBloatCheck函数用于检查数据库中索引膨胀情况，并以表格形式打印相关信息，同时输出相关建议。
func IndexBloatCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

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
		// 调用函数处理每个数据库的索引膨胀情况，更新hasData的值
		hasDataForDb := printIndexBloatTable(db)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("以下是数据库中索引膨胀相关信息：")
	} else {
		fmt.Println("未查询到数据库中索引膨胀相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 如果索引膨胀太大, 会影响性能, 建议重建索引, create index CONCURRENTLY.... ")
	fmt.Println()
}

// printIndexBloatTable 打印指定数据库的索引膨胀情况表格
func printIndexBloatTable(db string) bool {
	// 创建用于当前数据库表格输出的对象并设置表头
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "schema", "表名", "表膨胀系数", "索引名", "索引膨胀系数"})

	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 获取指定数据库中索引膨胀信息
	indexBloatInfoResult := ConnectPostgreSQL("[QUERY_INDEX_BLOAT_INFO]", db)
	if len(indexBloatInfoResult) > 0 {
		for _, row := range indexBloatInfoResult {
			writer.Append(row)
		}
		writer.Render()
		fmt.Println(buffer.String())
		currentHasData = true
	}

	return currentHasData
}
