package routineinspect

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetDBVersion() {
	// 执行SQL命令获取数据库版本数据
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", "select version()")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("执行获取数据库版本命令失败: %v\n", err)
		fmt.Println("输出结果：")
		fmt.Println(string(output))
		return
	}

	// 解析输出结果
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		fmt.Println("输出结果为空，无法获取数据库版本信息")
		return
	}
	version := lines[0]

	// 使用tablewriter创建表格并设置表头
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"数据库版本"})
	// 添加数据行
	table.Append([]string{version})

	// 渲染表格
	table.Render()

	// 输出建议内容
	fmt.Println("建议:")
	fmt.Println("   > 定期查看数据库版本，以便及时了解是否有可用的更新，更新数据库版本可能会带来性能提升、安全修复以及新功能特性等好处。")
}
