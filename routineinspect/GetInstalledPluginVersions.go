package routineinspect

import (
	"GoBasic/utils/fileutils"
	"fmt"
)

// GetInstalledPluginVersions 函数用于获取已安装插件的版本信息，并以表格形式展示，同时输出相关建议。
func GetInstalledPluginVersions(logWriter *fileutils.LogWriter, resultWriter *fileutils.ResultWriter) {
	logWriter.WriteLog("开始获取已安装插件的版本信息...")
	resultWriter.WriteResult("\n### 3.4、数据库插件版本:\n")
	dbNamesResult := ConnectPostgreSQL("[QUERY_NON_TEMPLATE_DBS]")
	if len(dbNamesResult) == 0 {
		logWriter.WriteLog("未查询到有效数据库名称")
		resultWriter.WriteResult("未查询到有效数据库名称")
		return
	}
	dbNames := make([]string, len(dbNamesResult))
	for i, row := range dbNamesResult {
		dbNames[i] = row[0]
	}

	// 用于存储所有插件版本信息的结果
	var allResult [][]string

	// 遍历每个数据库，获取插件版本信息并合并结果
	for _, db := range dbNames {
		pluginVersionsResult := ConnectPostgreSQL("[QUERY_PLUGIN_VERSIONS]", db)
		if len(pluginVersionsResult) > 0 {
			allResult = append(allResult, pluginVersionsResult...)
		}
	}

	// 根据是否有数据决定输出内容
	if len(allResult) > 0 {
		// Markdown 表格的表头
		tableHeader := "| 当前数据库 | 插件名称 | 插件所有者 | 插件命名空间 | 插件可重定位 | 插件版本 |"
		resultWriter.WriteResult(tableHeader)

		// Markdown 表格的分隔行
		separator := "|------------|----------|------------|--------------|--------------|---------|"
		resultWriter.WriteResult(separator)

		for _, row := range allResult {
			// 假设row是一个包含所需字段的切片
			resultWriter.WriteResult(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |",
				row[0], row[1], row[2], row[3], row[4], row[5]))
		}
	} else {
		logWriter.WriteLog("未查询到用户已安装的插件版本相关信息")
		resultWriter.WriteResult("未查询到用户已安装的插件版本相关信息")
	}

	// 打印建议
	suggestion := "> 定期检查已安装插件的版本，及时更新插件以获取更好的功能支持、性能优化以及安全修复等。"
	resultWriter.WriteResult("\n**建议:** \n" + suggestion)
}
