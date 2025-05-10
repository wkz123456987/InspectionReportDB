package detection

import (
	"GoBasic/utils/fileutils"
	"fmt"
	"strings"
)

// FileSystemInodeUsageCheck 读取配置文件并执行远程文件系统Inode使用情况检查
func FileSystemInodeUsageCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始巡检远程文件系统Inode使用情况...")

	RemoteFileSystemInodeUsageCheck(GetSSHConfig(logWriter), logWriter, resultWriter)
}

// RemoteFileSystemInodeUsageCheck 获取远程文件系统Inode使用情况并展示
func RemoteFileSystemInodeUsageCheck(sshConf SSHConfig, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	result, err := ExecuteRemoteCommand(sshConf, "df -ih")
	if err != nil {
		logWriter.WriteLog("执行远程命令失败: " + err.Error())
		return
	}

	processFileSystemInodeResult(result, resultWriter)
}

func processFileSystemInodeResult(result string, resultWriter *fileutils.ResultWriter) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	// 写入标题
	header := "### 1.5、远程文件系统Inode使用情况:\n"
	resultWriter.WriteResult(header)

	// Markdown 表格的表头和分隔行
	tableHeader := "| 文件系统     | inode容量 | 已使用 | 剩余 | 使用占比 | 挂载路径   |"
	separator := "|------------|----------|------|------|---------|----------|"
	resultWriter.WriteResult(tableHeader)
	resultWriter.WriteResult(separator)

	for index, line := range lines {
		if index == 0 {
			continue // 跳过表头行
		}
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			hasData = true
			// 将数据行添加到Markdown表格中
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |",
				fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]))
		}
	}

	if !hasData {
		resultWriter.WriteResult("未查询到远程文件系统Inode使用情况相关信息")
		return
	}

	// 写入说明和建议
	explanation := "\n**说明：** 在一个文件系统中，每个文件和目录都需要占用一个inode。当inode耗尽时，即使磁盘空间还有剩余，也无法创建新的文件"
	suggestion := "**建议:** \n > 时刻关注inode使用情况，及时清理无用文件和目录，释放inode空间。"
	resultWriter.WriteResult(explanation)
	resultWriter.WriteResult(suggestion)
}
