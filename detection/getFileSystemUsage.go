package detection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/ini.v1"
)

func FileSystemUsageCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始巡检远程文件系统使用情况...")
	cfg, err := ini.Load("database_config.ini")
	if cfg == nil || err != nil {
		logWriter.WriteLog("无法读取配置文件: " + err.Error())
		return
	}
	section := cfg.Section("Linux")
	user := section.Key("User").String()
	password := section.Key("Password").String()
	port, err := section.Key("Port").Int()
	if err != nil {
		logWriter.WriteLog("无法转换端口号: " + err.Error())
		return
	}
	host := section.Key("Host").String()

	sshConf := SSHConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}

	RemoteFileSystemUsageCheck(sshConf, logWriter, resultWriter)
}

// RemoteFileSystemUsageCheck 获取远程文件系统使用情况并展示
func RemoteFileSystemUsageCheck(sshConf SSHConfig, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	result, err := ExecuteRemoteCommand(sshConf, "df -h")
	if err != nil {
		logWriter.WriteLog("执行远程命令失败: " + err.Error())
		return
	}

	processFileSystemUsageResult(result, resultWriter)
}

func processFileSystemUsageResult(result string, resultWriter *fileutils.ResultWriter) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	// 写入标题
	header := "### 远程文件系统使用情况:\n"
	resultWriter.WriteResult(header)

	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(false) // 关闭自动格式化表头，避免重复表头
	writer.SetHeader([]string{"文件系统", "总大小", "已用大小", "可用大小", "使用占比", "挂载点"})
	writer.SetAlignment(tablewriter.ALIGN_LEFT)

	// 重新解析结果并添加数据到表格，跳过第一行（表头行）
	for _, line := range lines {
		if len(line) == 0 {
			continue // 跳过空行
		}
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			hasData = true
			writer.Append([]string{
				fields[0],
				fields[1],
				fields[2],
				fields[3],
				fields[4],
				fields[5],
			})
		}
	}

	if !hasData {
		resultWriter.WriteResult("未查询到远程文件系统使用情况相关信息")
		return
	}

	// 所有数据处理完毕后，渲染表格并写入结果
	writer.Render()
	resultWriter.WriteResult(buffer.String())

	// 写入建议
	resultWriter.WriteResult("建议: ")
	resultWriter.WriteResult("   > 注意预留足够的空间给数据库. ")
}
