package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// TableAgeCheck函数用于检查表年龄情况，并以表格形式打印相关信息，同时输出相关建议。
func TableAgeCheck() {
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
		// 调用函数处理每个数据库的表年龄情况，更新hasData的值
		hasDataForDb := printTableAgeTable(db)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  表年龄:")
	} else {
		fmt.Println("未查询到表年龄相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 表的年龄正常情况下应该小于vacuum_freeze_table_age, 如果剩余年龄小于5亿, 建议人为干预, 将LONG SQL或事务杀掉后, 执行vacuum freeze. ")
	fmt.Println()
}

// printTableAgeTable 打印指定数据库的表年龄情况表格
func printTableAgeTable(db string) bool {
	// 创建用于当前数据库表格输出的对象并设置表头
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "rolname", "nspname", "relkind", "表名", "年龄", "年龄_剩余"})

	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 获取指定数据库中表年龄信息
	tableAgeInfoResult := ConnectPostgreSQL("[QUERY_TABLE_AGE_INFO]", db)
	if len(tableAgeInfoResult) > 0 {
		for _, row := range tableAgeInfoResult {
			writer.Append(row)
		}
		writer.Render()
		fmt.Println(buffer.String())
		currentHasData = true
	}

	return currentHasData
}
