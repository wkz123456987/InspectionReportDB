package detection

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/mem"
)

// ByteToMB 将字节转换为MB
func ByteToMB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024)
}

// getMemoryUsage 获取内存使用情况并返回所有相关参数
func getMemoryUsage() (*mem.VirtualMemoryStat, error) {
	vm, err := mem.VirtualMemory() // 获取虚拟内存使用情况
	if err != nil {
		return nil, err
	}
	return vm, nil
}

// getSwapUsage 获取交换空间的使用情况
func getSwapUsage() (*mem.SwapMemoryStat, error) {
	swap, err := mem.SwapMemory() // 获取交换内存使用情况
	if err != nil {
		return nil, err
	}
	return swap, nil
}

// PrintMemoryUsage 打印内存使用情况的表格
// PrintMemoryUsage 打印内存使用情况的表格
func PrintMemoryUsage() {
	// 获取虚拟内存信息
	vm, err := getMemoryUsage()
	if err != nil {
		fmt.Println("获取内存使用情况失败:", err)
		return
	}

	// 获取交换空间信息
	swap, err := getSwapUsage()
	if err != nil {
		fmt.Println("获取交换空间使用情况失败:", err)
		return
	}

	// 清屏（仅适用于Unix-like操作系统）
	//fmt.Print("\033[H\033[2J")

	// 创建内存使用情况的表格
	memTable := tablewriter.NewWriter(os.Stdout)
	memTable.SetHeader([]string{"内存项", "值"})

	// 添加内存使用情况数据
	memTable.Append([]string{"总内存", fmt.Sprintf("%.2f MB", ByteToMB(vm.Total))})
	memTable.Append([]string{"已用内存", fmt.Sprintf("%.2f MB", ByteToMB(vm.Used))})
	memTable.Append([]string{"剩余内存", fmt.Sprintf("%.2f MB", ByteToMB(vm.Free))})
	memTable.Append([]string{"缓存内存", fmt.Sprintf("%.2f MB", ByteToMB(vm.Cached))})
	memTable.Append([]string{"已用内存百分比", fmt.Sprintf("%.2f%%", vm.UsedPercent)})

	// 添加交换空间相关信息
	memTable.Append([]string{"交换空间总量", fmt.Sprintf("%.2f MB", ByteToMB(swap.Total))})
	memTable.Append([]string{"交换空间已用", fmt.Sprintf("%.2f MB", ByteToMB(swap.Used))})
	memTable.Append([]string{"交换空间剩余", fmt.Sprintf("%.2f MB", ByteToMB(swap.Free))})
	memTable.Append([]string{"交换空间使用百分比", fmt.Sprintf("%.2f%%", swap.UsedPercent)})

	// 刷新表格
	memTable.Render()
}
