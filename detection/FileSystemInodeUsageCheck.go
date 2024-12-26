package detection

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/ini.v1"
)

func FileSystemInodeUsageCheck() {
	cfg, err := ini.Load("database_config.ini")
	if cfg == nil || err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}
	section := cfg.Section("Linux")
	user := section.Key("User").String()
	password := section.Key("Password").String()
	port, err := section.Key("Port").Int()
	if err != nil {
		log.Fatalf("无法转换端口号: %v", err)
	}
	host := section.Key("Host").String()

	sshConf := SSHConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
	RemoteFileSystemInodeUsageCheck(sshConf)

}

// RemoteFileSystemInodeUsageCheck 获取远程文件系统Inode使用情况并展示
func RemoteFileSystemInodeUsageCheck(sshConf SSHConfig) {
	result, err := ExecuteRemoteCommand(sshConf, "df -ih")
	if err != nil {
		fmt.Println(err)
		return
	}

	processFileSystemInodeResult(result)
}

func processFileSystemInodeResult(result string) {
	lines := strings.Split(strings.TrimSpace(result), "\n")
	hasData := false

	fmt.Println("### 远程文件系统Inode使用情况:")
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(false)
	writer.SetHeader([]string{"文件系统", "inode容量", "已使用", "剩余", "使用占比", "挂载路径"})
	writer.SetAlignment(tablewriter.ALIGN_LEFT)

	for index, line := range lines {
		if index == 0 {
			continue // 跳过表头行
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
		fmt.Println("未查询到远程文件系统Inode使用情况相关信息")
		return
	}

	writer.Render()
	fmt.Println(buffer.String())
	fmt.Println("说明：在一个文件系统中，每个文件和目录都需要占用一个inode。当inode耗尽时，即使磁盘空间还有剩余，也无法创建新的文件")
	fmt.Println("建议: ")
	fmt.Println("   > 时刻关注inode使用情况，及时清理无用文件和目录，释放inode空间。 ")
}
