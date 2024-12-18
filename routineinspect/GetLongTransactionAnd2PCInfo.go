package routineinspect

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetLongTransactionAnd2PCInfo() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建第一个psql命令获取长事务相关信息（第一部分表格数据）
	cmd1 := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select datname,usename,query,xact_start,now()-xact_start xact_duration,query_start,now()-query_start query_duration,state from pg_stat_activity where state<>'idle' and (backend_xid is not null or backend_xmin is not null) and now()-xact_start > interval '30 min' order by xact_start`)
	var result1 bytes.Buffer
	cmd1.Stdout = &result1
	err := cmd1.Run()
	if err != nil {
		fmt.Printf("执行获取长事务相关信息命令失败: %s\n", err)
		return
	}

	// 解析第一个命令的结果判断是否有有效数据（对应第一部分表格）
	lines1 := strings.Split(strings.TrimSpace(result1.String()), "\n")
	for _, line := range lines1 {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 9 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容（先处理第一部分表格数据）
	if hasData {
		buffer1 := &bytes.Buffer{}
		writer1 := tablewriter.NewWriter(buffer1)
		writer1.SetAutoFormatHeaders(true)
		writer1.SetHeader([]string{"数据库名", "用户名", "查询语句", "事务开始时间", "事务持续时间", "查询开始时间", "查询持续时间", "状态"})

		// 重新解析第一部分结果并添加数据到表格
		lines1 = strings.Split(strings.TrimSpace(result1.String()), "\n")
		for _, line := range lines1 {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 9 {
				writer1.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
					strings.TrimSpace(matches[3]),
					strings.TrimSpace(matches[4]),
					strings.TrimSpace(matches[5]),
					strings.TrimSpace(matches[6]),
					strings.TrimSpace(matches[7]),
					strings.TrimSpace(matches[8]),
				})
			}
		}

		writer1.Render()
		fmt.Println(buffer1.String())
	} else {
		fmt.Println("未查事务持续时长(长事务)超过30分钟的相关信息")
	}

	// 重置hasData标记，准备获取并处理第二部分数据
	hasData = false

	// 构建第二个psql命令获取2PC相关信息（第二部分表格数据，此处示例中SQL语句和前面获取长事务一样，实际可能需根据真实需求调整）
	cmd2 := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select datname,usename,query,xact_start,now()-xact_start xact_duration,query_start,now()-query_start query_duration,state from pg_stat_activity where state<>'idle' and (backend_xid is not null or backend_xmin is not null) and now()-xact_start > interval '30 min' order by xact_start`)
	var result2 bytes.Buffer
	cmd2.Stdout = &result2
	err = cmd2.Run()
	if err != nil {
		fmt.Printf("执行获取2PC相关信息命令失败: %s\n", err)
		return
	}

	// 解析第二个命令的结果判断是否有有效数据（对应第二部分表格）
	lines2 := strings.Split(strings.TrimSpace(result2.String()), "\n")
	for _, line := range lines2 {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则，和前面格式一致所以复用）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 9 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容（再处理第二部分表格数据）
	if hasData {
		buffer2 := &bytes.Buffer{}
		writer2 := tablewriter.NewWriter(buffer2)
		writer2.SetAutoFormatHeaders(true)
		writer2.SetHeader([]string{"数据库名", "用户名", "查询语句", "事务开始时间", "事务持续时间", "查询开始时间", "查询持续时间", "状态"})

		// 重新解析第二部分结果并添加数据到表格
		lines2 = strings.Split(strings.TrimSpace(result2.String()), "\n")
		for _, line := range lines2 {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 9 {
				writer2.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
					strings.TrimSpace(matches[3]),
					strings.TrimSpace(matches[4]),
					strings.TrimSpace(matches[5]),
					strings.TrimSpace(matches[6]),
					strings.TrimSpace(matches[7]),
					strings.TrimSpace(matches[8]),
				})
			}
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
