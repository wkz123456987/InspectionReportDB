package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// DatabaseAgeCheck函数用于检查数据库年龄情况，并以表格形式打印相关信息，同时输出相关建议。
func DatabaseAgeCheck() {
	// 标记是否获取到有效数据，初始化为false

	// 获取数据库年龄信息
	result := ConnectPostgreSQL("[QUERY_DATABASE_AGE_INFO]")
	if len(result) > 0 {
		fmt.Println("###  数据库年龄:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"数据库", "年龄", "年龄_剩余"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到数据库年龄相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 数据库的年龄正常情况下应该小于vacuum_freeze_table_age, 如果剩余年龄小于5亿, 建议人为干预, 将LONG SQL或事务杀掉后, 执行vacuum freeze. ")
	fmt.Println()
}
