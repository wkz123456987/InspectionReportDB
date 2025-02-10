package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetUserPasswordExpiration 用于获取用户密码到期时间信息，并以表格形式展示，同时输出相关建议。
func GetUserPasswordExpiration(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取用户密码到期时间信息...")
	resultWriter.WriteResult("\n### 3.14、用户密码到期时间信息:\n")

	// 获取用户密码到期时间信息
	result := ConnectPostgreSQL("[QUERY_USER_PASSWORD_EXPIRATION]")
	if len(result) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 用户名 | 密码到期时间 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|--------|--------------|"
		resultWriter.WriteResult(separator)

		for _, row := range result {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s |",
				row[0], row[1]))
		}
	} else {
		logWriter.WriteLog("未查询到用户密码到期时间相关信息")
		resultWriter.WriteResult("未查询到用户密码到期时间相关信息")
	}

	// 打印建议
	suggestion := "> 到期后，用户将无法登陆，记得修改密码，同时将密码到期时间延长到某个时间或无限时间，使用alter role... VALID UNTIL 'timestamp'。"
	resultWriter.WriteResult("\n**建议:**\n " + suggestion)
}
