package main

import (
	"GoBasic/detection"
	"GoBasic/inspection"
	"GoBasic/routineinspect"
	"GoBasic/utils/fileutils"
	"fmt"
	"log"
	"time"

	"gopkg.in/ini.v1"
)

// 定义全局变量
var (
	resultWriter *fileutils.ResultWriter
	logWriter    *fileutils.LogWriter
)

func main() {

	// 读取配置文件获取数据库配置信息
	cfg, err := ini.Load("../config/database_config.ini")
	if err != nil || cfg == nil {
		log.Fatalf("无法加载配置文件: %v", err)
		logWriter.WriteLog("无法加载配置文件: " + err.Error())
	}
	section := cfg.Section("Config")
	LogDir := section.Key("LogDir").String()
	ResultDir := section.Key("ResultDir").String()
	// 获取当前时间，并格式化为2024-12-10_143023的格式
	currentTime := time.Now().Format("2006-01-02_150405")
	section2 := cfg.Section("Linux")
	hostname := section2.Key("Host")
	// 初始化全局变量resultWriter
	logFileName := fmt.Sprintf("%s-%s.log", hostname, currentTime) // 创建完整的日志文件名，包含时间戳
	// 使用文件名创建日志写入器
	logWriter, err = fileutils.NewLogWriter(LogDir, logFileName)
	if err != nil {
		fmt.Printf("无法创建日志文件: %v\n", err)
		logWriter.WriteLog("无法创建日志文件: " + err.Error())
		return
	}
	defer logWriter.Close()
	// 初始化全局变量resultWriter
	resultFileName := fmt.Sprintf("%s-%s.md", hostname, currentTime)
	resultWriter, err = fileutils.NewResultWriter(ResultDir, resultFileName)
	if err != nil {
		fmt.Printf("无法创建结果文件: %v\n", err)
		logWriter.WriteLog("无法创建结果文件: " + err.Error())
		return
	}
	defer resultWriter.Close()

	// 定义一个简单的ASCII艺术字体
	art := `
		_   _      _ _         __        __        _   _
	   | \ | |    | | |  __  / _|  ___  \ \      / / | |
	   |  \| | ___| | | / _|| |_ / _ \   \ \ /\ / /| |_| |
	   | . \ |/ _ \ | | |  _||  _|  __/    \ V  V / |  _  |
	   |_| \_/\___/_|_| |_|  |_|  \___|      \_/\_/  |_| |_|
		`
	fmt.Println("开始巡检系统资源使用情况...")
	fmt.Println(art)
	resultWriter.WriteResult("开始巡检系统资源使用情况...")
	resultWriter.WriteResult(art)
	os_detection()                //调用操作系统巡检
	database_inspection()         //调用数据库巡检
	database_routine_inspection() //调用数据库常规巡检
}

// 操作系统巡检
func os_detection() {
	resultWriter.WriteResult("\n## 一、操作系统巡检\n")
	detection.CPUUsageCheck(logWriter, resultWriter)
	detection.MemoryUsageCheck(logWriter, resultWriter)
	detection.DiskIOCheck(logWriter, resultWriter)
	detection.FileSystemUsageCheck(logWriter, resultWriter)
	detection.FileSystemInodeUsageCheck(logWriter, resultWriter)
}

// 数据库巡检
func database_inspection() {
	resultWriter.WriteResult("\n## 二、数据库巡检重点关注的巡检项\n")
	fmt.Println("###  TOP 10 size对象:")
	inspection.DatabasesTop10(logWriter, resultWriter)
	fmt.Println("###  查找索引数超过4并且SIZE大于10MB的表")
	inspection.TablesWithTooManyIndexes(logWriter, resultWriter)
	fmt.Println("###  重复创建的索引:")
	inspection.DatabasesRepeatIndex(logWriter, resultWriter)
	fmt.Println("###  上次巡检以来未使用或使用较少的索引:")
	inspection.UnusedIndexesSinceLastCheck(logWriter, resultWriter)
	fmt.Println("###  获取数据库统计信息,回滚比例, 命中比例, 数据块读写时间, 死锁, 复制冲突:")
	inspection.DatabaseStats(logWriter, resultWriter)
	fmt.Println("###  检查数据库中索引膨胀情况:")
	inspection.IndexBloatCheck(logWriter, resultWriter)
	fmt.Println("###   检查数据库中垃圾数据情况:")
	inspection.GarbageDataCheck(logWriter, resultWriter)
	fmt.Println("###   数据库年龄:")
	inspection.DatabaseAgeCheck(logWriter, resultWriter)
	fmt.Println("###  表年龄:")
	inspection.TableAgeCheck(logWriter, resultWriter)
	fmt.Println("###  锁等待:")
	inspection.LockWaitCheck(logWriter, resultWriter)
	fmt.Println("###  密码泄露检查:")
	inspection.PasswordLeakCheck(logWriter, resultWriter)
	fmt.Println("###  复制槽状态:")
	inspection.ReplicationSlotStatus(logWriter, resultWriter)
	fmt.Println("###  schema统计:")
	inspection.GetSchemaStats(logWriter, resultWriter)
}

// 数据库常规巡检
func database_routine_inspection() {
	resultWriter.WriteResult("\n## 三、数据库巡检常规巡检项\n")
	fmt.Println("###  当前活跃度:")
	routineinspect.GetCurrentActivityStatus(logWriter, resultWriter)
	fmt.Println("###  总剩余连接数:")
	routineinspect.CheckDBConnections(logWriter, resultWriter)
	fmt.Println("###  数据库版本:")
	routineinspect.GetDBVersion(logWriter, resultWriter)
	fmt.Println("###  数据库插件版本:")
	routineinspect.GetInstalledPluginVersions(logWriter, resultWriter)
	fmt.Println("###  用户使用了多少种数据类型:")
	routineinspect.GetUsedDataTypeCounts(logWriter, resultWriter)
	fmt.Println("###  用户创建了多少对象:	")
	routineinspect.GetCreatedObjectCounts(logWriter, resultWriter)
	fmt.Println("###  用户对象占用空间的柱状图:")
	routineinspect.GetUserObjectSpaceInfo(logWriter, resultWriter)
	fmt.Println("###  表空间使用情况:")
	routineinspect.GetTablespaceUsage(logWriter, resultWriter)
	fmt.Println("###  数据库使用情况:")
	routineinspect.GetDatabaseUsage(logWriter, resultWriter)
	fmt.Println("###  用户连接数限制:")
	routineinspect.GetUserConnectionLimits(logWriter, resultWriter)
	fmt.Println("###  数据库连接数限制:")
	routineinspect.GetDatabaseConnectionLimits(logWriter, resultWriter)
	fmt.Println("###  数据库检查点和bgwriter统计信息:")
	routineinspect.GetCheckpointBgwriterStats(logWriter, resultWriter)
	fmt.Println("###  长事务和2PC相关信息:")
	routineinspect.GetLongTransactionAnd2PCInfo(logWriter, resultWriter)
	fmt.Println("###  数据库的用户密码到期时间:")
	routineinspect.GetUserPasswordExpiration(logWriter, resultWriter)
	fmt.Println("###  表的继承关系检查:")
	routineinspect.GetInheritanceRelationCheck(logWriter, resultWriter)
	fmt.Println("###  是否开启归档, 自动垃圾回收:")
	routineinspect.GetArchiveAndAutoVacuumSettings(logWriter, resultWriter)
	fmt.Println("###  数据库主备角色:")
	routineinspect.GetMasterStandbyRole(logWriter, resultWriter)
	fmt.Println("###  备库信息:")
	routineinspect.GetStandbyInfo(logWriter, resultWriter)
}
