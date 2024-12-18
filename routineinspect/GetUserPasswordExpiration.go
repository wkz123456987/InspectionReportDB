package routineinspect

import (
	"bytes"

	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetUserPasswordExpiration() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false
	sql := `SELECT 
    rolname,
    CASE 
        WHEN rolvaliduntil IS NULL THEN '无有效期'
        ELSE rolvaliduntil::text
    END AS rolvaliduntil
FROM 
    pg_authid
ORDER BY 
    CASE 
        WHEN rolvaliduntil IS NULL THEN '9999-12-31 23:59:59.999999+00'
        ELSE rolvaliduntil
    END;`

	// 构建psql命令获取用户密码到期时间信息
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", sql)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("执行获取用户密码到期时间命令失败: %s\n", err)
		return
	}

	// 解析结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 3 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"用户名", "密码到期时间"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 3 {
				writer.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到用户密码到期时间相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 到期后，用户将无法登陆，记得修改密码，同时将密码到期时间延长到某个时间或无限时间，alter role... VALID UNTIL 'timestamp'.")
	fmt.Println()
}
