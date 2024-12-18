package routineinspect

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetCheckpointBgwriterStats() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建psql命令获取检查点、bgwriter统计信息
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `select * from pg_stat_bgwriter`)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("执行获取检查点、bgwriter统计信息命令失败: %s\n", err)
		return
	}

	// 解析结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（根据实际格式调整正则）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 12 { // 第一个匹配项是完整的匹配项，后面是各列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"checkpoints_timed", "checkpoints_req", "checkpoint_write_time", "checkpoint_sync_time", "buffers_checkpoint", "buffers_clean", "maxwritten_clean", "buffers_backend", "buffers_backend_fsync", "buffers_alloc", "stats_reset"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 12 {
				writer.Append([]string{
					strings.TrimSpace(matches[1]),
					strings.TrimSpace(matches[2]),
					strings.TrimSpace(matches[3]),
					strings.TrimSpace(matches[4]),
					strings.TrimSpace(matches[5]),
					strings.TrimSpace(matches[6]),
					strings.TrimSpace(matches[7]),
					strings.TrimSpace(matches[8]),
					strings.TrimSpace(matches[9]),
					strings.TrimSpace(matches[10]),
					strings.TrimSpace(matches[11]),
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到检查点、bgwriter统计信息相关信息")
	}

	// 打印建议
	fmt.Println("建议: ")
	fmt.Println("   > 如果检测结果显示checkpoint_write_time多，说明检查点持续时间长，检查点过程中产生了较多的脏页。")
	fmt.Println("    checkpoint_sync_time代表检查点开始时的shared buffer中的脏页被同步到磁盘的时间，如果时间过长，并且数据库在检查点时性能较差，考虑一下提升块设备的IOPS能力。")
	fmt.Println("    buffers_backend_fsync太多说明需要加大shared buffer 或者 减小bgwriter_delay参数。")
	fmt.Println()
}
