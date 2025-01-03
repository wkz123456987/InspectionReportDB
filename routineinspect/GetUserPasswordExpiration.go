package routineinspect

import (
	"GoBasic/utils/fileutils"
	"bytes"

	"github.com/olekukonko/tablewriter"
)

// GetUserPasswordExpiration 用于获取用户密码到期时间信息，并以表格形式展示，同时输出相关建议。
func GetUserPasswordExpiration(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取用户密码到期时间信息...")
	resultWriter.WriteResult("\n###  用户密码到期时间信息:\n")
	// 获取用户密码到期时间信息
	result := ConnectPostgreSQL("[QUERY_USER_PASSWORD_EXPIRATION]")
	if len(result) > 0 {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"用户名", "密码到期时间"})

		for _, row := range result {
			writer.Append(row)
		}

		writer.Render()
		resultWriter.WriteResult(buffer.String())
	} else {
		logWriter.WriteLog("未查询到用户密码到期时间相关信息")
		resultWriter.WriteResult("未查询到用户密码到期时间相关信息")
	}

	// 打印建议
	suggestion := `
    建议:
        > 到期后，用户将无法登陆，记得修改密码，同时将密码到期时间延长到某个时间或无限时间，使用alter role... VALID UNTIL 'timestamp'。
	`
	resultWriter.WriteResult(suggestion)
}
