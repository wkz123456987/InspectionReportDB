package routineinspect

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func CheckDBConnections() {
	// 执行SQL命令获取数据
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c",
		"SELECT setting::int AS max_conn FROM pg_settings WHERE name='max_connections' "+
			"UNION ALL "+
			"SELECT count(*) AS used FROM pg_stat_activity "+
			"UNION ALL "+
			"SELECT setting::int AS res_for_super FROM pg_settings WHERE name='superuser_reserved_connections' "+
			"UNION ALL "+
			"SELECT max_conn - used AS res_for_normal FROM (SELECT setting::int AS max_conn FROM pg_settings WHERE name='max_connections') settings, (SELECT count(*) AS used FROM pg_stat_activity) activity")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("执行数据库查询失败: %v\n", err)
		fmt.Println("输出结果：")
		fmt.Println(string(output))
		return
	}

	// 解析输出结果
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 4 {
		fmt.Println("输出结果不完整")
		return
	}

	// 提取数据
	maxConn, _ := strconv.Atoi(strings.Fields(lines[0])[1])
	usedConn, _ := strconv.Atoi(strings.Fields(lines[1])[1])
	resForSuper, _ := strconv.Atoi(strings.Fields(lines[2])[1])
	resForNormal, _ := strconv.Atoi(strings.Fields(lines[3])[1])

	// 使用tablewriter创建表格并设置表头
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"总连接", "已使用连接", "剩余给超级用户连接", "剩余给普通用户连接"})

	// 添加数据行
	table.Append([]string{fmt.Sprintf("%d", maxConn), fmt.Sprintf("%d", usedConn), fmt.Sprintf("%d", resForSuper), fmt.Sprintf("%d", resForNormal)})

	// 渲染表格
	table.Render()
	// 输出建议内容
	fmt.Println("建议:")
	fmt.Println("   > 给超级用户和普通用户设置足够的连接, 以免不能登录数据库. ")
}
