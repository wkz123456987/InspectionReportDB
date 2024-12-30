# InspectionReportDB

## 使用方法

**克隆数据库巡检GO语言代码**
```shell
git clone https://github.com/wkz123456987/InspectionReportDB.git
```
**根据实际业务修改配置文件**

```bash
vim config/database_config
```

**切换到bin目录下**
```shell
cd bin
```
**执行巡检命令**

```bash
./DatabaseInspection
```

## 版本更替

| 版本号 | 功能更新描述 |
| ---- | ---- |
| v 0.1 | 由shell脚本修改为基本能跑的go语言工具。 |
| v 0.3 | 添加远程巡检Linux系统、PostgreSQL数据库功能。 |
| v 0.5 | 重构代码结构，提高代码复用度，提高程序可扩展性，添加写入日志功能。 |
| v 0.7 | 添加远程巡检密码认证功能,支持多种常见加密方式，包括 MD5、SHA256 等 |






## 项目文件架构

```plaintext
[root@postgres delve_test]# tree
.
|-- bin
|   `-- DatabaseInspection
|-- config
|   `-- database_config.ini
|-- detection
|   |-- CPUUsageCheck.go
|   |-- DiskIOCheck.go
|   |-- FileSystemInodeUsageCheck.go
|   |-- getFileSystemUsage.go
|   |-- MemoryUsageCheck.go
|   `-- remote_executor.go
|-- go.mod
|-- go.sum
|-- inspection
|   |-- ConnectPostgreSQL.go
|   |-- DatabaseAgeCheck.go
|   |-- DatabasesRepeatIndex.go
|   |-- DatabasesSizeTop10.go
|   |-- DatabaseStats.go
|   |-- GarbageDataCheck.go
|   |-- IndexBloatCheck.go
|   |-- LockWaitCheck.go
|   |-- PasswordLeakCheck.go
|   |-- ReplicationSlotStatus.go
|   |-- SchemaStats.go
|   |-- TableAgeCheck.go
|   |-- TablesWithTooManyIndexes.go
|   `-- UnusedIndexesSinceLastCheck.go
|-- log
|   `-- inspection-2024-12-30_144736.log
|-- main.go
|-- README.md
|-- result
|   `-- inspect_result-2024-12-30_144736.txt
|-- routineinspect
|   |-- ConnectPostgreSQL.go
|   |-- GetArchiveAndAutoVacuumSettings.go
|   |-- GetCheckpointBgwriterStats.go
|   |-- GetCreatedObjectCounts.go
|   |-- GetCurrentActivityStatus.go
|   |-- GetDatabaseConnectionLimits.go
|   |-- GetDatabaseUsage.go
|   |-- GetDBVersion.go
|   |-- GetInheritanceRelationCheck.go
|   |-- GetInstalledPluginVersions.go
|   |-- GetLongTransactionAnd2PCInfo.go
|   |-- GetMasterStandbyRole.go
|   |-- GetStandbyInfo.go
|   |-- GetTablespaceUsage.go
|   |-- GetUsedDataTypeCounts.go
|   |-- GetUserConnectionLimits.go
|   |-- GetUserObjectSpaceInfo.go
|   |-- GetUserPasswordExpiration.go
|   `-- TotalRemainingConnections.go
|-- SQL
|   |-- inspection.sql
|   `-- routineinspect.sql
`-- utils
    `-- fileutils
        `-- fileutils.go

```