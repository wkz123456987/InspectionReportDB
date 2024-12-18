// detection/diskusage.go
package detection

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v4/disk"
)

// getDiskUsage 获取磁盘使用情况并返回
func getDiskUsage() (*disk.UsageStat, error) {
	u, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	return u, nil
}

// PrintDiskUsage 打印磁盘使用情况表格
func PrintDiskUsage() {
	u, err := getDiskUsage()
	if err != nil {
		fmt.Println("获取磁盘使用情况失败:", err)
		return
	}

	// 创建一个新的表格 writer
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"磁盘项", "值"})

	// 添加数据行
	table.Append([]string{"磁盘总空间", fmt.Sprintf("%d GB", u.Total/1024/1024/1024)})
	table.Append([]string{"已用空间", fmt.Sprintf("%d GB", u.Used/1024/1024/1024)})
	table.Append([]string{"可用空间", fmt.Sprintf("%d GB", u.Free/1024/1024/1024)})
	table.Append([]string{"使用率", fmt.Sprintf("%.2f%%", u.UsedPercent)})

	// 刷新表格
	table.Render()
}
