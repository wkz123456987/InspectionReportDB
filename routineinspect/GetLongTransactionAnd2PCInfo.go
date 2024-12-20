package routineinspect

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

func GetLongTransactionAnd2PCInfo() {
	// 获取长事务相关信息
	result1 := ConnectPostgreSQL("[QUERY_LONG_TRANSACTION_INFO]")
	if len(result1) > 0 {
		buffer1 := &bytes.Buffer{}
		writer1 := tablewriter.NewWriter(buffer1)
		writer1.SetAutoFormatHeaders(true)
		writer1.SetHeader([]string{"数据库名", "用户名", "查询语句", "事务开始时间", "事务持续时间", "查询开始时间", "查询持续时间", "状态"})

		for _, row := range result1 {
			writer1.Append(row)
		}

		writer1.Render()
		fmt.Println(buffer1.String())
	} else {
		fmt.Println("未查事务持续时长(长事务)超过30分钟的相关信息")
	}

	// 获取2PC相关信息
	result2 := ConnectPostgreSQL("[QUERY_2PC_INFO]")
	if len(result2) > 0 {
		buffer2 := &bytes.Buffer{}
		writer2 := tablewriter.NewWriter(buffer2)
		writer2.SetAutoFormatHeaders(true)
		writer2.SetHeader([]string{"数据库名", "用户名", "查询语句", "事务开始时间", "事务持续时间", "查询开始时间", "查询持续时间", "状态"})

		for _, row := range result2 {
			writer2.Append(row)
		}

		writer2.Render()
		fmt.Println(buffer2.String())
	} else {
		fmt.Println("未查询到2PC持续时长超过30 分钟的相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 长事务过程中产生的垃圾，无法回收，建议不要在数据库中运行LONG SQL，或者错开DML高峰时间去运行LONG SQL。2PC事务一定要记得尽快结束掉，否则可能会导致数据库膨胀。")
	fmt.Println()
}
