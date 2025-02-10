package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
	"regexp"
	"strings"
)

// GetTablespaceUsage 函数用于获取表空间使用情况信息，并以表格形式展示，同时输出相关建议。
func GetTablespaceUsage(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取表空间使用情况信息...")
	resultWriter.WriteResult("\n### 3.8、表空间使用情况:\n")

	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建psql命令以获取表空间使用情况信息
	//	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select spcname,pg_tablespace_location(oid),pg_size_pretty(pg_tablespace_size(oid)) from pg_tablespace order by pg_tablespace_size(oid) desc`)

	// cmd.Stdout = &result
	// err := cmd.Run()
	// if err != nil {
	// 	logWriter.WriteLog(fmt.Sprintf("执行获取表空间使用情况命令失败: %s", err))
	// 	resultWriter.WriteResult(fmt.Sprintf("执行获取表空间使用情况命令失败: %s", err))
	// 	return
	// }
	// 获取主备库角色信息
	// 获取主备库角色信息
	result := ConnectPostgreSQL("[QUERY_MASTER_STANDBY_ROLE]")
	// 解析结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		// Markdown 表格的表头
		tableHeader := "| 表空间名 | 表空间路径 | 表空间大小 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|------------|--------------|------------|"
		resultWriter.WriteResult(separator)

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 4 {
				resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s |",
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
					strings.TrimSpace(matches[3])))
			}
		}
	} else {
		resultWriter.WriteResult("未查询到表空间使用情况相关信息")
	}

	// 打印建议
	suggestion := "> 注意检查表空间所在文件系统的剩余空间, (默认表空间在$PGDATA/base目录下), IOPS分配是否均匀, OS的sysstat包可以观察IO使用率."
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
