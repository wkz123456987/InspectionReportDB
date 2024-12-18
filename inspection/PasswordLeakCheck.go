package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// PasswordLeakCheck 函数用于检查密码泄露情况，并以表格形式打印相关信息。
func PasswordLeakCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasDataPgAuthid := false
	hasDataPgUserMappings := false
	hasDataPgViews := false

	// 打印整体标题
	fmt.Println("###  密码泄露检查:")

	// 检查 pg_authid部分
	fmt.Println("#### 检查 pg_authid：")
	hasDataPgAuthid = checkPgAuthid()

	// 检查 pg_user_mappings, pg_views部分
	fmt.Println("#### 检查 pg_user_mappings, pg_views：")
	hasDataPgUserMappings = checkPgUserMappings()
	hasDataPgViews = checkPgViews()

	// 根据是否有数据决定各部分输出内容
	if hasDataPgAuthid {
		fmt.Println("以下是pg_authid相关检查结果：")
	} else {
		fmt.Println("未查询到pg_authid中密码泄露相关信息")
	}

	if hasDataPgUserMappings {
		fmt.Println("以下是pg_user_mappings相关检查结果：")
	} else {
		fmt.Println("未查询到pg_user_mappings中密码泄露相关信息")
	}

	if hasDataPgViews {
		fmt.Println("以下是pg_views相关检查结果：")
	} else {
		fmt.Println("未查询到pg_views中密码泄露相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 如果以上输出显示密码已泄露, 尽快修改, 并通过参数避免密码又被记录到以上文件中(psql -n) (set log_statement='none'; set log_min_duration_statement=-1; set log_duration=off; set pg_stat_statements.track_utility=off;). ")
	fmt.Println("    明文密码不安全, 建议使用create|alter role... encrypted password. ")
	fmt.Println("    在fdw, dblink based view中不建议使用密码明文. ")
	fmt.Println("    在recovery.*的配置中不要使用密码, 不安全, 可以使用.pgpass配置密码. ")
	fmt.Println()
}

// checkPgAuthid 检查pg_authid中密码相关情况并输出结果
func checkPgAuthid() bool {
	// 标记当前部分是否获取到有效数据，初始化为false
	hasData := false

	// 构建psql命令以获取pg_authid相关信息
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select count(*) from pg_authid where rolpassword!~ '^md5' or length(rolpassword)<>35`)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command for pg_authid check: %s\n", err)
		return false
	}

	// 解析结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		if line != "" {
			hasData = true
			break
		}
	}

	if hasData {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"计数"})

		writer.Append([]string{lines[0]})
		writer.Render()
		fmt.Println(buffer.String())
	}

	return hasData
}

// checkPgUserMappings 检查pg_user_mappings中密码相关情况并输出结果
func checkPgUserMappings() bool {
	// 标记当前部分是否获取到有效数据，初始化为false
	hasData := false

	// 执行psql命令获取数据库列表
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "-A", "-q", "-c", "SELECT datname FROM pg_database WHERE datname NOT IN ('template0', 'template1')")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command for getting database list: %s\n", err)
		return false
	}

	// 解析数据库列表并遍历
	dbList := strings.Split(strings.TrimSpace(out.String()), "\n")
	var tableData string
	for _, db := range dbList {
		if db == "" {
			continue
		}
		// 构建psql命令以获取pg_user_mappings相关信息
		cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select current_database(),* from pg_user_mappings where umoptions::text ~* 'password'`)
		var result bytes.Buffer
		cmd.Stdout = &result
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to execute command for pg_user_mappings in database %s: %s\n", db, err)
			continue
		}

		// 使用正则表达式提取每行的数据（可根据实际数据格式调整正则表达式）
		lines := strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 6 {
				if hasData == false {
					headers := []string{"数据库", "umid", "umuser", "usename", "umoptions"}
					tableData = buildTableRow(headers) + buildTableRow([]string{strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2]), strings.TrimSpace(matches[3]), strings.TrimSpace(matches[4]), strings.TrimSpace(matches[5])})
					hasData = true
				} else {
					tableData += buildTableRow([]string{strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2]), strings.TrimSpace(matches[3]), strings.TrimSpace(matches[4]), strings.TrimSpace(matches[5])})
				}
			}
		}
	}

	if hasData {
		fmt.Println(tableData)
	}

	return hasData
}

// checkPgViews 检查pg_views中密码相关情况并输出结果
func checkPgViews() bool {
	// 标记当前部分是否获取到有效数据，初始化为false
	hasData := false

	// 执行psql命令获取数据库列表
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "-A", "-q", "-c", "SELECT datname FROM pg_database WHERE datname NOT IN ('template0', 'template1')")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command for getting database list: %s\n", err)
		return false
	}

	// 解析数据库列表并遍历
	dbList := strings.Split(strings.TrimSpace(out.String()), "\n")
	var tableData string
	for _, db := range dbList {
		if db == "" {
			continue
		}
		// 构建psql命令以获取pg_views相关信息
		cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select current_database(),* from pg_views where definition ~* 'password' and definition ~* 'dblink'`)
		var result bytes.Buffer
		cmd.Stdout = &result
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to execute command for pg_views in database %s: %s\n", db, err)
			continue
		}

		// 使用正则表达式提取每行的数据（可根据实际数据格式调整正则表达式）
		lines := strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 6 {
				if hasData == false {
					headers := []string{"数据库", "schemaname", "viewname", "viewowner", "definition"}
					tableData = buildTableRow(headers) + buildTableRow([]string{strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2]), strings.TrimSpace(matches[3]), strings.TrimSpace(matches[4]), strings.TrimSpace(matches[5])})
					hasData = true
				} else {
					tableData += buildTableRow([]string{strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2]), strings.TrimSpace(matches[3]), strings.TrimSpace(matches[4]), strings.TrimSpace(matches[5])})
				}
			}
		}
	}

	if hasData {
		fmt.Println(tableData)
	}

	return hasData
}

// buildTableRow 辅助函数，用于构建表格中的一行数据（以制表符分隔各列）
func buildTableRow(data []string) string {
	var row string
	for i, col := range data {
		if i > 0 {
			row += "\t"
		}
		row += col
	}
	row += "\n"
	return row
}
