package routineinspect

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetUserObjectSpaceInfo() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 先获取所有非template数据库名称
	dbNamesCmd := exec.Command("psql", "--pset=pager=off", "-t", "-A", "-q", "-c", `select datname from pg_database where datname not in ('template0', 'template1')`)
	var dbNamesResult bytes.Buffer
	dbNamesCmd.Stdout = &dbNamesResult
	err := dbNamesCmd.Run()
	if err != nil {
		fmt.Printf("执行获取数据库名称命令失败: %s\n", err)
		return
	}
	dbNames := strings.Split(strings.TrimSpace(dbNamesResult.String()), "\n")

	// 用于存储所有用户对象空间信息结果
	var allResult bytes.Buffer

	// 遍历每个数据库，获取用户对象占用空间信息并合并结果
	for _, db := range dbNames {
		cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select current_database(),buk this_buk_no,cnt rels_in_this_buk,pg_size_pretty(min) buk_min,pg_size_pretty(max) buk_max from( select row_number() over (partition by buk order by tsize),tsize,buk,min(tsize) over (partition by buk),max(tsize) over (partition by buk),count(*) over (partition by buk) cnt from ( select pg_relation_size(a.oid) tsize, width_bucket(pg_relation_size(a.oid),tmin-1,tmax+1,10) buk from (select min(pg_relation_size(a.oid)) tmin,max(pg_relation_size(a.oid)) tmax from pg_class a,pg_namespace c where a.relnamespace=c.oid and nspname!~ '^pg_' and nspname<>'information_schema') t, pg_class a,pg_namespace c where a.relnamespace=c.oid and nspname!~ '^pg_' and nspname<>'information_schema') t)t where row_number=1;`)
		var result bytes.Buffer
		cmd.Stdout = &result
		err := cmd.Run()
		if err != nil {
			fmt.Printf("执行获取用户对象空间信息命令失败（数据库：%s）: %s\n", db, err)
			continue
		}
		allResult.WriteString(result.String())
	}

	// 解析合并后的结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(allResult.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 6 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"当前数据库", "桶编号", "此桶中关系数量", "桶最小值（格式化后）", "桶最大值（格式化后）"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(allResult.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 6 {
				writer.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
					strings.TrimSpace(matches[3]),
					strings.TrimSpace(matches[4]),
					strings.TrimSpace(matches[5]),
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到用户对象占用空间相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 关注用户对象占用空间情况，对于占用空间较大的对象可考虑优化存储结构或进行归档处理，以节省数据库空间。")
	fmt.Println()
}
