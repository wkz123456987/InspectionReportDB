package routineinspect

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetInstalledPluginVersions() {
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

	// 用于存储所有插件版本信息的结果
	var allResult bytes.Buffer

	// 遍历每个数据库，获取插件版本信息并合并结果
	for _, db := range dbNames {
		// 构建带参数的SQL语句，这里以数据库名为例作为参数
		sql := `select 
    current_database(), 
    e.extname, 
    u.usename as extowner, 
    n.nspname as extnamespace, 
    e.extrelocatable, 
    e.extversion 
from 
    pg_extension e
left join 
    pg_user u on e.extowner = u.usesysid
left join 
    pg_namespace n on e.extnamespace = n.oid;`
		// 修改后的SQL语句，获取插件相关的核心信息列，这里可根据实际需求进一步调整列的选择
		cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", sql)
		var result bytes.Buffer
		cmd.Stdout = &result
		err := cmd.Run()
		if err != nil {
			fmt.Printf("执行获取插件版本信息命令失败（数据库：%s）: %s\n", db, err)
			continue
		}
		allResult.WriteString(result.String())
	}

	// 解析合并后的结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(allResult.String()), "\n")
	for _, line := range lines {
		// 根据实际返回的6列数据修改正则表达式
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 7 { // 第一个匹配项是完整的匹配项，后面是6列的数据，所以共7项
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"当前数据库", "插件名称", "插件所有者", "插件命名空间", "插件可重定位", "插件版本"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(allResult.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 7 {
				writer.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
					strings.TrimSpace(matches[3]),
					strings.TrimSpace(matches[4]),
					strings.TrimSpace(matches[5]),
					strings.TrimSpace(matches[6]),
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到用户已安装的插件版本相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 定期检查已安装插件的版本，及时更新插件以获取更好的功能支持、性能优化以及安全修复等。")
	fmt.Println()
}
