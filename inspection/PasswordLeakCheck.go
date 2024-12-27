package inspection

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// PasswordLeakCheck 函数用于检查密码泄露情况，并以表格形式打印相关信息，同时输出相关建议。
func PasswordLeakCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始检查密码泄露情况...")
	// 标记是否获取到有效数据，初始化为false
	hasDataPgAuthid := false
	hasDataPgUserMappings := false
	hasDataPgViews := false

	// 打印整体标题
	resultWriter.WriteResult("### 密码泄露检查:")

	// 检查 pg_authid部分
	resultWriter.WriteResult("#### 检查 pg_authid：")
	hasDataPgAuthid = checkPgAuthid(logWriter, resultWriter)

	// 检查 pg_user_mappings, pg_views部分
	resultWriter.WriteResult("#### 检查 pg_user_mappings, pg_views：")
	hasDataPgUserMappings = checkPgUserMappings(logWriter, resultWriter)
	hasDataPgViews = checkPgViews(logWriter, resultWriter)

	// 根据是否有数据决定各部分输出内容
	if hasDataPgAuthid {
		resultWriter.WriteResult("以下是pg_authid相关检查结果：")
	} else {
		resultWriter.WriteResult("未查询到pg_authid中密码泄露相关信息")
	}

	if hasDataPgUserMappings {
		resultWriter.WriteResult("以下是pg_user_mappings相关检查结果：")
	} else {
		resultWriter.WriteResult("未查询到pg_user_mappings中密码泄露相关信息")
	}

	if hasDataPgViews {
		resultWriter.WriteResult("以下是pg_views相关检查结果：")
	} else {
		resultWriter.WriteResult("未查询到pg_views中密码泄露相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 如果以上输出显示密码已泄露, 尽快修改, 并通过参数避免密码又被记录到以上文件中(psql -n) (set log_statement='none'; set log_min_duration_statement=-1; set log_duration=off; set pg_stat_statements.track_utility=off;). 
        明文密码不安全, 建议使用create|alter role... encrypted password. 
        在fdw, dblink based view中不建议使用密码明文. 
        在recovery.*的配置中不要使用密码, 不安全, 可以使用.pgpass配置密码. 
	`
	resultWriter.WriteResult(suggestion)
}

// checkPgAuthid 检查pg_authid中密码相关情况并输出结果
func checkPgAuthid(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) bool {
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
		resultWriter.WriteResult(buffer.String())
		hasData = true
	} else {
		logWriter.WriteLog("未查询到pg_authid中密码泄露相关信息")
		resultWriter.WriteResult("未查询到pg_authid中密码泄露相关信息")
	}

	return hasData
}

// checkPgUserMappings 检查pg_user_mappings中密码相关情况并输出结果
func checkPgUserMappings(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) bool {
	// 标记当前部分是否获取到有效数据，初始化为false
	hasData := false

	// 获取非template数据库名称
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		logWriter.WriteLog("未查询到有效数据库名称")
		resultWriter.WriteResult("未查询到有效数据库名称\n")
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
		resultWriter.WriteResult(tableData)
	} else {
		logWriter.WriteLog("未查询到pg_user_mappings中密码泄露相关信息")
		resultWriter.WriteResult("未查询到pg_user_mappings中密码泄露相关信息")
	}

	return hasData
}

// checkPgViews 检查pg_views中密码相关情况并输出结果
func checkPgViews(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) bool {
	// 标记当前部分是否获取到有效数据，初始化为false
	hasData := false

	// 获取非template数据库名称
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		logWriter.WriteLog("未查询到有效数据库名称")
		resultWriter.WriteResult("未查询到有效数据库名称\n")
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
		resultWriter.WriteResult(tableData)
	} else {
		logWriter.WriteLog("未查询到pg_views中密码泄露相关信息")
		resultWriter.WriteResult("未查询到pg_views中密码泄露相关信息")
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
