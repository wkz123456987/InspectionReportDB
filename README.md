# InspectionReportDB

## 使用方法

### 1、克隆数据库巡检GO语言代码
```shell
git clone https://github.com/wkz123456987/InspectionReportDB.git
```
### 2、根据实际业务修改配置文件

```bash
[root@postgres bin]# cat database_config.ini 
[Config]
LogDir = ./log/
ResultDir = ./result1/ 

[Linux]
User = "root" #用户名
Encryption_method = "AES-256" #明文（plaintext）/密文（AES-256）
#Password = "linux"
Password = "AAAAAAAAAAAAAAAAAAAAADwuyD63xJzrc88GR5F873Q=" 
Host = 192.168.6.1    # 主机名
Port = 22     # 端口

[Database]
DBName = postgres # 数据库名
Hostname = 192.168.6.1 # 主机名
Port = 8432  # 端口
Username = fbase  # 用户名
# 密码
Encryption_method = "AES-256" #明文（plaintext）/密文（AES-256）
Password = "AAAAAAAAAAAAAAAAAAAAAEmrZGR98OqLTqoVqDxku/k="
#Password = "linux123!@#"
```
### 3、如果要使用密文认证

**在bin目录下使用./encrypt_password 程序**
```bash
[root@postgres bin]# ./encrypt_password 


███████╗██╗   ██╗██████╗  ██████╗ ███╗   ███╗███████╗
██╔════╝╚██╗ ██╔╝██╔══██╗██╔═══██╗████╗ ████║██╔════╝
█████╗   ╚████╔╝ ██████╔╝██║   ██║██╔████╔██║█████╗  
██╔══╝    ╚██╔╝  ██╔═══╝ ██║   ██║██║╚██╔╝██║██╔══╝  
███████╗   ██║   ███████╗╚██████╔╝██║ ╚═╝ ██║███████╗
╚══════╝   ╚═╝   ╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝


Usage: ./encrypt_password [options] <plaintext>
Options:
  --help, -h    Show this help message and exit

Examples:
  $ ./encrypt_password "myPassword123"
  $ ./encrypt_password --help
```
```bash
[root@postgres bin]# ./encrypt_password 'linux123!@#'

███████╗██╗   ██╗██████╗  ██████╗ ███╗   ███╗███████╗
██╔════╝╚██╗ ██╔╝██╔══██╗██╔═══██╗████╗ ████║██╔════╝
█████╗   ╚████╔╝ ██████╔╝██║   ██║██╔████╔██║█████╗  
██╔══╝    ╚██╔╝  ██╔═══╝ ██║   ██║██║╚██╔╝██║██╔══╝  
███████╗   ██║   ███████╗╚██████╔╝██║ ╚═╝ ██║███████╗
╚══════╝   ╚═╝   ╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝

加密后的密码:
AAAAAAAAAAAAAAAAAAAAAEmrZGR98OqLTqoVqDxku/k=
```
**把加密后的密码复制到配置文件**

### 4、切换到bin目录下开始使用巡检程序

```shell
[root@postgres bin]# ./DatabaseInspection --help
Usage of ./DatabaseInspection:
  -config string
        Path to the configuration file (default "./config/database_config.ini")
```
**执行巡检命令**

```bash
[root@postgres bin]# ./DatabaseInspection -config database_config.ini
开始巡检系统资源使用情况...

███████╗██╗   ██╗██████╗  ██████╗ ███╗   ███╗███████╗
██╔════╝╚██╗ ██╔╝██╔══██╗██╔═══██╗████╗ ████║██╔════╝
█████╗   ╚████╔╝ ██████╔╝██║   ██║██╔████╔██║█████╗  
██╔══╝    ╚██╔╝  ██╔═══╝ ██║   ██║██║╚██╔╝██║██╔══╝  
███████╗   ██║   ███████╗╚██████╔╝██║ ╚═╝ ██║███████╗
╚══════╝   ╚═╝   ╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝

2025-03-01 11:56:27.071 CST 开始巡检远程系统CPU使用率...
2025-03-01 11:56:29.351 CST 开始巡检远程内存使用率...
2025-03-01 11:56:31.478 CST 开始巡检远程磁盘IO使用率...

```

### 4、巡检多个数据库

**新建shell/bash脚本**
```bash
[root@pg-2 bin]# cat start.sh 
/root/delve_test/bin/DatabaseInspection -config /root/delve_test/bin/node1.ini
/root/delve_test/bin/DatabaseInspection -config /root/delve_test/bin/node2.ini
/root/delve_test/bin/DatabaseInspection -config /root/delve_test/bin/node3.ini
```
**node1.ini**

```bash
[root@pg-2 bin]# cat node1.ini 
[Config]
LogDir = ./node1/log/
ResultDir = ./node1/ 

[Linux]
User = "root" #用户名
Password = "linux123!@#"
Host = 192.168.101.149    # 主机名
Port = 22     # 端口

[Database]
DBName = postgres # 数据库名
Hostname = 192.168.101.149 # 主机名
Port = 10432  # 端口
Username = fbase  # 用户名
# 密码
Password = "linux123!@#"
```

**node2.ini**
```bash
[root@pg-2 bin]# cat node2.ini 
[Config]
LogDir = ./node2/log/
ResultDir = ./node2/ 

[Linux]
User = "root" #用户名
Password = "linux123!@#"
Host = 192.168.101.150    # 主机名
Port = 22     # 端口

[Database]
DBName = postgres # 数据库名
Hostname = 192.168.101.150 # 主机名
Port = 8432  # 端口
Username = fbase  # 用户名
# 密码
Password = "linux123!@#"
```

## 版本更替

| 版本号 | 功能更新描述 |
| ---- | ---- |
| v 0.1 | 由shell脚本修改为基本能跑的go语言工具。 |
| v 0.3 | 添加远程巡检Linux系统、PostgreSQL数据库功能。 |
| v 0.5 | 重构代码结构，提高代码复用度，提高程序可扩展性，添加写入日志功能。 |
| v 0.7 | 添加远程巡检密码认证功能,支持多种常见加密方式，包括 明文,AES256加密 等 |
| v 1.0 | 添加对多台机器同事巡检的功能，支持巡检整个集群功能






## 项目文件架构

```plaintext
[root@postgres delve_test]# tree
.
|-- bin
|   |-- database_config.ini
|   |-- DatabaseInspection
|   `-- encrypt_password
|-- config
|   |-- config.go
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
|-- main.go
|-- README.md
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

8 directories, 51 files

```