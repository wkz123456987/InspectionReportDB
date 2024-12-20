package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// LockWaitCheck函数用于检查锁等待情况，并以表格形式打印相关信息，同时输出相关建议。
func LockWaitCheck() {

	// 获取锁等待信息
	result := ConnectPostgreSQL("[QUERY_LOCK_WAIT_INFO]")
	if len(result) > 0 {
		fmt.Println("###  锁等待:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"锁类型", "读锁模式", "读锁用户", "读锁数据库", "关联关系", "读锁进程ID", "读锁页面", "读锁元组", "读锁事务开始时间", "读锁查询开始时间", "读锁锁定时长", "读锁查询语句", "写锁模式", "写锁进程ID", "写锁页面", "写锁元组", "写锁事务开始时间", "写锁查询开始时间", "写锁锁定时长", "写锁查询语句"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到锁等待相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 锁等待状态, 反映业务逻辑的问题或者SQL性能有问题, 建议深入排查持锁的SQL.")
	fmt.Println()
}
