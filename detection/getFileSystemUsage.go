package detection

import (
	"GoBasic/utils/fileutils"
	"fmt"
	"strings"
)

func FileSystemUsageCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始巡检远程文件系统使用情况...")
	RemoteFileSystemUsageCheck(GetSSHConfig(logWriter), logWriter, resultWriter)
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
	header := "### 1.4、文件系统使用情况:\n"
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	resultWriter.WriteResult("| 文件系统       | 总大小     | 已用大小     | 可用大小     | 使用占比     | 挂载点   |")
	resultWriter.WriteResult("|--------------|----------|----------|----------|----------|--------|")

	// 重新解析结果并添加数据到表格，跳过第一行（表头行）
	for _, line := range lines[1:] { // 跳过第一行表头
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			hasData = true
			// 将数据行添加到Markdown表格中
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |",
				fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]))
		}
	}

	if !hasData {
		resultWriter.WriteResult("未查询到远程文件系统使用情况相关信息")
		return
	}

	// 写入建议
	suggestion := "\n**建议:** \n   > 注意预留足够的空间给数据库. "
	resultWriter.WriteResult(suggestion)
}
