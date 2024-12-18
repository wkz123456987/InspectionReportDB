package detection

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/cpu"
)

// getCPUUsage 获取 CPU 使用情况并返回
func getCPUUsage() ([]float64, error) {
	percent, err := cpu.Percent(0, false) // 获取 CPU 使用率（百分比）
	if err != nil {
		return nil, err
	}
	return percent, nil
}

// PrintCPUUsage 打印 CPU 使用情况表格
func PrintCPUUsage() {
	percent, err := getCPUUsage()
	if err != nil {
		fmt.Println("获取CPU使用情况失败:", err)
		return
	}

	// 清屏（仅适用于Unix-like操作系统）
	//fmt.Print("\033[H\033[2J")

	// 创建一个新的表格 writer
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CPU项", "值"})

	// 添加数据行
	table.Append([]string{"CPU 使用率", fmt.Sprintf("%.2f%%", percent[0])})

	// 渲染表格
	table.Render()
}
