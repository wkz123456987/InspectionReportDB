package detection

import (
	"GoBasic/utils/fileutils"
	"bytes"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func DiskIOCheck(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始巡检远程磁盘IO使用率...")
	RemoteDiskIOCheck(GetSSHConfig(logWriter), logWriter, resultWriter)
}

// RemoteDiskIOCheck 获取远程磁盘IO情况并展示
func RemoteDiskIOCheck(sshConf SSHConfig, logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	firstResult, err := ExecuteRemoteCommand(sshConf, "iostat -mx 1 1")
	if err != nil {
		logWriter.WriteLog("执行远程命令失败: " + err.Error())
		return
	}
	diskDevices := parseDiskDevices(firstResult)
	if len(diskDevices) > 0 {
		header := "### 远程输出磁盘IO情况:\n"
		resultWriter.WriteResult(header)
		for _, disk := range diskDevices {
			ioResult, err := ExecuteRemoteCommand(sshConf, fmt.Sprintf("iostat -mx 1 2 %s", disk))
			if err != nil {
				logWriter.WriteLog(fmt.Sprintf("Failed to execute command for disk %s: %s", disk, err))
				continue
			}
			parseAndAppendIOResult(ioResult, disk, resultWriter)
		}
	} else {
		resultWriter.WriteResult("未查询到远程磁盘IO情况相关信息")
	}

	resultWriter.WriteResult("建议: ")
	resultWriter.WriteResult("   > 注意检查IO占用高的原因.")
}

func parseDiskDevices(result string) []string {
	var diskDevices []string
	lines := strings.Split(strings.TrimSpace(result), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 1 && (strings.Contains(line, "sd") || strings.Contains(line, "vd") || strings.Contains(line, "dm")) {
			diskDevices = append(diskDevices, fields[0])
		}
	}
	return diskDevices
}

func parseAndAppendIOResult(result string, disk string, resultWriter *fileutils.ResultWriter) {
	ioLines := strings.Split(strings.TrimSpace(result), "\n")
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader([]string{"设备名", "磁盘IO：%util"})
	writer.SetAlignment(tablewriter.ALIGN_LEFT)

	for index, ioLine := range ioLines {
		fields := strings.Fields(ioLine)
		if len(fields) >= 14 {
			ioUtil := fields[13]
			if index > 0 { // 跳过表头行，只处理数据行
				ioUtil = strings.TrimSpace(strings.TrimPrefix(ioUtil, "%util"))
				if ioUtil != "" {
					writer.Append([]string{disk, ioUtil}) // 将设备名和处理后的%util值添加到表格中
				}
			}
		}
	}
	writer.Render()
	resultWriter.WriteResult(buffer.String())
}
