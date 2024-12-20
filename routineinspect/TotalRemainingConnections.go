package routineinspect

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// CheckDBConnections 函数用于获取数据库连接相关信息，并以表格形式展示，同时输出相关建议。
func CheckDBConnections() {
	result := ConnectPostgreSQL("[QUERY_DB_CONNECTIONS]")
	if len(result) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"总连接", "已使用连接", "剩余给超级用户连接", "剩余给普通用户连接"})
		table.Append(result[0])
		table.Render()
	} else {
		fmt.Println("未查询到有效数据")
	}

	// 输出建议内容
	fmt.Println("建议:")
	fmt.Println("   > 给超级用户和普通用户设置足够的连接, 以免不能登录数据库. ")
}
