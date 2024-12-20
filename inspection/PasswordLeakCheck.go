package inspection

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// PasswordLeakCheck函数用于检查密码泄露情况，并以表格形式打印相关信息，同时输出相关建议。
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

	// 获取pg_authid相关信息
	result := ConnectPostgreSQL("[QUERY_PG_AUTHID_CHECK]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"计数"})

		writer.Append(result[0])
		writer.Render()
		fmt.Println(buffer.String())
		hasData = true
	}

	return hasData
}

// checkPgUserMappings 检查pg_user_mappings中密码相关情况并输出结果
func checkPgUserMappings() bool {
	// 标记当前部分是否获取到有效数据，初始化为false
	hasData := false

	// 获取非template数据库名称
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		fmt.Printf("未查询到有效数据库名称\n")
		return false
	}
	dbList := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbList[i] = row[0]
	}

	var tableData string
	for _, db := range dbList {
		if db == "" {
			continue
		}
		// 获取pg_user_mappings相关信息
		pgUserMappingsResult := ConnectPostgreSQL("[QUERY_PG_USER_MAPPINGS_CHECK]", db)
		if len(pgUserMappingsResult) > 0 {
			if hasData == false {
				headers := []string{"数据库", "umid", "umuser", "usename", "umoptions"}
				tableData = buildTableRow(headers)
			}
			for _, row := range pgUserMappingsResult {
				tableData += buildTableRow(row)
			}
			hasData = true
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

	// 获取非template数据库名称
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		fmt.Printf("未查询到有效数据库名称\n")
		return false
	}
	dbList := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbList[i] = row[0]
	}

	var tableData string
	for _, db := range dbList {
		if db == "" {
			continue
		}
		// 获取pg_views相关信息
		pgViewsResult := ConnectPostgreSQL("[QUERY_PG_VIEWS_CHECK]", db)
		if len(pgViewsResult) > 0 {
			if hasData == false {
				headers := []string{"数据库", "schemaname", "viewname", "viewowner", "definition"}
				tableData = buildTableRow(headers)
			}
			for _, row := range pgViewsResult {
				tableData += buildTableRow(row)
			}
			hasData = true
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
