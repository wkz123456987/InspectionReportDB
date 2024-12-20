package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GarbageDataCheck函数用于检查数据库中垃圾数据情况，并以表格形式打印相关信息，同时输出相关建议。
func GarbageDataCheck() {
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
		// 调用函数处理每个数据库的垃圾数据情况，更新hasData的值
		hasDataForDb := printGarbageDataTable(db)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  垃圾数据:")
	} else {
		fmt.Println("未查询到数据库中垃圾数据相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 通常垃圾过多, 可能是因为无法回收垃圾, 或者回收垃圾的进程繁忙或没有及时唤醒, 或者没有开启autovacuum, 或在短时间内产生了大量的垃圾. ")
	fmt.Println("    可以等待autovacuum进行处理, 或者手工执行vacuum table. ")
	fmt.Println()
}

// printGarbageDataTable 打印指定数据库的垃圾数据情况表格
func printGarbageDataTable(db string) bool {
	// 创建用于当前数据库表格输出的对象并设置表头
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "schema", "表名", "死元组数量"})

	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 获取指定数据库中垃圾数据信息
	garbageDataInfoResult := ConnectPostgreSQL("[QUERY_GARBAGE_DATA_INFO]", db)
	if len(garbageDataInfoResult) > 0 {
		for _, row := range garbageDataInfoResult {
			writer.Append(row)
		}
		writer.Render()
		fmt.Println(buffer.String())
		currentHasData = true
	}

	return currentHasData
}
