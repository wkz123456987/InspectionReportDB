package main

import (
	"GoBasic/detection"
	"GoBasic/inspection"
	"GoBasic/routineinspect"
	"fmt"
)

func main() {

	// 定义一个简单的ASCII艺术字体
	art := `
		_   _      _ _         __        __        _   _
	   | \ | |    | | |  __  / _|  ___  \ \      / / | |
	   |  \| | ___| | | / _|| |_ / _ \   \ \ /\ / /| |_| |
	   | . \ |/ _ \ | | |  _||  _|  __/    \ V  V / |  _  |
	   |_| \_/\___/_|_| |_|  |_|  \___|      \_/\_/  |_| |_|
		`
	// 清屏（仅适用于Unix-like操作系统）
	//fmt.Print("\033[H\033[2J")
	fmt.Println("开始监控系统资源使用情况...")
	fmt.Println(art)
	//os_detection()        //调用操作系统巡检
	//database_inspection() //调用数据库巡检
	//database_routine_inspection() //调用数据库常规巡检
	test()
}

func test() {
	//detection.FileSystemUsageCheck()
	//detection.CPUUsageCheck()
	//detection.MemoryUsageCheck()
	detection.DiskIOCheck()
	//detection.FileSystemInodeUsageCheck()
}

// 操作系统巡检
func os_detection() {

}

// 数据库巡检
func database_inspection() {
	fmt.Println("###  TOP 10 size对象:")
	inspection.DatabasesTop10()
	fmt.Println("###  查找索引数超过4并且SIZE大于10MB的表")
	inspection.TablesWithTooManyIndexes()
	fmt.Println()
	fmt.Println("###  重复创建的索引:")
	inspection.DatabasesRepeatIndex()
	fmt.Println("###  上次巡检以来未使用或使用较少的索引:")
	inspection.UnusedIndexesSinceLastCheck()
	fmt.Println("###  获取数据库统计信息,回滚比例, 命中比例, 数据块读写时间, 死锁, 复制冲突:")
	inspection.DatabaseStats()
	fmt.Println("###  检查数据库中索引膨胀情况:")
	inspection.IndexBloatCheck()
	fmt.Println("###   检查数据库中垃圾数据情况:")
	inspection.GarbageDataCheck()
	fmt.Println("###   数据库年龄:")
	inspection.DatabaseAgeCheck()
	fmt.Println("###  表年龄:")
	inspection.TableAgeCheck()
	fmt.Println("###  锁等待:")
	inspection.LockWaitCheck()
	fmt.Println("###  密码泄露检查:")
	inspection.PasswordLeakCheck()
	fmt.Println("###  复制槽状态:")
	inspection.ReplicationSlotStatus()
	fmt.Println("###  schema统计:")
	inspection.GetSchemaStats()
}

// 数据库常规巡检
func database_routine_inspection() {
	fmt.Println("###  当前活跃度:")
	routineinspect.GetCurrentActivityStatus()
	fmt.Println("###  总剩余连接数:")
	routineinspect.CheckDBConnections()
	fmt.Println("###  数据库版本:")
	routineinspect.GetDBVersion()
	fmt.Println("###  数据库插件版本:")
	routineinspect.GetInstalledPluginVersions()
	fmt.Println("###  用户使用了多少种数据类型:")
	routineinspect.GetUsedDataTypeCounts()
	fmt.Println("###  用户创建了多少对象:	")
	routineinspect.GetCreatedObjectCounts()
	fmt.Println("###  用户对象占用空间的柱状图:")
	routineinspect.GetUserObjectSpaceInfo()
	fmt.Println("###  表空间使用情况:")
	routineinspect.GetTablespaceUsage()
	fmt.Println("###  数据库使用情况:")
	routineinspect.GetDatabaseUsage()
	fmt.Println("###  用户连接数限制:")
	routineinspect.GetUserConnectionLimits()
	fmt.Println("###  数据库连接数限制:")
	routineinspect.GetDatabaseConnectionLimits()
	fmt.Println("###  数据库检查点和bgwriter统计信息:")
	routineinspect.GetCheckpointBgwriterStats()
	fmt.Println("###  长事务和2PC相关信息:")
	routineinspect.GetLongTransactionAnd2PCInfo()
	fmt.Println("###  数据库的用户密码到期时间:")
	routineinspect.GetUserPasswordExpiration()
	fmt.Println("###  继承关系检查:")
	routineinspect.GetInheritanceRelationCheck()
	fmt.Println("###  是否开启归档, 自动垃圾回收:")
	routineinspect.GetArchiveAndAutoVacuumSettings()
	fmt.Println("###  数据库主备角色:")
	routineinspect.GetMasterStandbyRole()
	fmt.Println("###  备库信息:")
	routineinspect.GetStandbyInfo()
}
