package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// LockWaitCheck 函数用于检查锁等待情况，并以表格形式打印相关信息。
func LockWaitCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 构建psql命令以获取锁等待信息
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `
    with t_wait as                     
    (select a.mode,a.locktype,a.database,a.relation,a.page,a.tuple,a.classid,
    a.objid,a.objsubid,a.pid,a.virtualtransaction,a.virtualxid,a.transactionid,b.query,b.xact_start,b.query_start,b.usename,b.datname 
    from pg_locks a,pg_stat_activity b where a.pid=b.pid and not a.granted),
    t_run as 
    (select a.mode,a.locktype,a.database,a.relation,a.page,a.tuple,
    a.classid,a.objid,a.objsubid,a.pid,a.virtualtransaction,a.virtualxid,
    a.transactionid,b.query,b.xact_start,b.query_start,
    b.usename,b.datname from pg_locks a,pg_stat_activity b where 
    a.pid=b.pid and a.granted) 
    select r.locktype,r.mode r_mode,r.usename r_user,r.datname r_db,
    r.relation::regclass,r.pid r_pid,
    r.page r_page,r.tuple r_tuple,r.xact_start r_xact_start,
    r.query_start r_query_start,
    now()-r.query_start r_locktime,r.query r_query,w.mode w_mode,
    w.pid w_pid,w.page w_page,
    w.tuple w_tuple,w.xact_start w_xact_start,w.query_start w_query_start,
    now()-w.query_start w_locktime,w.query w_query  
    from t_wait w,t_run r where
    r.locktype is not distinct from w.locktype and
    r.database is not distinct from w.database and
    r.relation is not distinct from w.relation and
    r.page is not distinct from w.page and
    r.tuple is not distinct from w.tuple and
    r.classid is not distinct from w.classid and
    r.objid is not distinct from w.objid and
    r.objsubid is not distinct from w.objsubid and
    r.transactionid is not distinct from w.transactionid and
    r.pid <> w.pid
    order by 
    ((  case w.mode
        when 'INVALID' then 0
        when 'AccessShareLock' then 1
        when 'RowShareLock' then 2
        when 'RowExclusiveLock' then 3
        when 'ShareUpdateExclusiveLock' then 4
        when 'ShareLock' then 5
        when 'ShareRowExclusiveLock' then 6
        when 'ExclusiveLock' then 7
        when 'AccessExclusiveLock' then 8
        else 0
    end  ) + 
    (  case r.mode
        when 'INVALID' then 0
        when 'AccessShareLock' then 1
        when 'RowShareLock' then 2
        when 'RowExclusiveLock' then 3
        when 'ShareUpdateExclusiveLock' then 4
        when 'ShareLock' then 5
        when 'ShareRowExclusiveLock' then 6
        when 'ExclusiveLock' then 7
        when 'AccessExclusiveLock' then 8
        else 0
    end  )) desc,r.xact_start;`)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 解析结果判断是否有有效数据
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		// 使用正则表达式提取每行的数据（可根据实际数据格式调整正则表达式）
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 21 { // 第一个匹配项是完整的匹配项，后面是列的数据
			hasData = true
			break
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("###  锁等待:")

		buffer := &bytes.Buffer{}
		writer := tablewriter.NewWriter(buffer)
		writer.SetAutoFormatHeaders(true)
		writer.SetHeader([]string{"锁类型", "读锁模式", "读锁用户", "读锁数据库", "关联关系", "读锁进程ID", "读锁页面", "读锁元组", "读锁事务开始时间", "读锁查询开始时间", "读锁锁定时长", "读锁查询语句", "写锁模式", "写锁进程ID", "写锁页面", "写锁元组", "写锁事务开始时间", "写锁查询开始时间", "写锁锁定时长", "写锁查询语句"})

		// 重新解析结果并添加数据到表格
		lines = strings.Split(strings.TrimSpace(result.String()), "\n")
		for _, line := range lines {
			re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
			matches := re.FindStringSubmatch(line)

			if len(matches) == 21 {
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
					strings.TrimSpace(matches[12]),
					strings.TrimSpace(matches[13]),
					strings.TrimSpace(matches[14]),
					strings.TrimSpace(matches[15]),
					strings.TrimSpace(matches[16]),
					strings.TrimSpace(matches[17]),
					strings.TrimSpace(matches[18]),
					strings.TrimSpace(matches[19]),
					strings.TrimSpace(matches[20]),
				})
			}
		}

		writer.Render()
		fmt.Println(buffer.String())
	} else {
		fmt.Println("未查询到锁等待相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 锁等待状态, 反映业务逻辑的问题或者SQL性能有问题, 建议深入排查持锁的SQL.")
	fmt.Println()
}
