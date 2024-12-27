package inspection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// TableAgeCheck 函数用于检查表年龄情况，并以表格形式打印相关信息，同时输出相关建议。
func TableAgeCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始检查表年龄情况...")
	resultWriter.WriteResult("\n### 表年龄:\n")
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 获取非template数据库名称
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		logWriter.WriteLog("未查询到有效数据库名称")
		resultWriter.WriteResult("未查询到有效数据库名称")
		return
	}
	dbList := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbList[i] = row[0]
	}

	for _, db := range dbList {
		if db == "" {
			continue
		}
		// 调用函数处理每个数据库的表年龄情况，更新hasData的值
		hasDataForDb := printTableAgeTable(db, logWriter, resultWriter)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {

	} else {
		resultWriter.WriteResult("未查询到表年龄相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 表的年龄正常情况下应该小于vacuum_freeze_table_age, 如果剩余年龄小于5亿, 建议人为干预, 将LONG SQL或事务杀掉后, 执行vacuum freeze.
	`
	resultWriter.WriteResult(suggestion)
}

// printTableAgeTable 打印指定数据库的表年龄情况表格
func printTableAgeTable(db string, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) bool {
	// 创建用于当前数据库表格输出的对象并设置表头
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "rolname", "nspname", "relkind", "表名", "年龄", "年龄_剩余"})

	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 获取指定数据库中表年龄信息
	tableAgeInfoResult := ConnectPostgreSQL("[QUERY_TABLE_AGE_INFO]", db)
	if len(tableAgeInfoResult) > 0 {
		for _, row := range tableAgeInfoResult {
			writer.Append(row)
		}
		writer.Render()
		resultWriter.WriteResult(buffer.String())
		currentHasData = true
	} else {
		logWriter.WriteLog(fmt.Sprintf("在数据库 %s 中未查询到表年龄相关信息", db))
		resultWriter.WriteResult(fmt.Sprintf("在数据库 %s 中未查询到表年龄相关信息", db))
	}

	return currentHasData
}
