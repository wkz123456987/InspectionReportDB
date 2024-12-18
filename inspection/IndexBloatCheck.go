package inspection

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// IndexBloatCheck 函数用于检查数据库中索引膨胀情况，并以表格形式打印相关信息。
func IndexBloatCheck() {
	// 标记是否获取到有效数据，初始化为false
	hasData := false

	// 执行psql命令获取数据库列表
	cmd := exec.Command("psql", "--pset=pager=off", "-t", "-A", "-q", "-c", "SELECT datname FROM pg_database WHERE datname NOT IN ('template0', 'template1')")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return
	}

	// 解析数据库列表并遍历
	dbList := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, db := range dbList {
		if db == "" {
			continue
		}
		// 调用函数处理每个数据库的索引膨胀情况，更新hasData的值
		hasDataForDb := printIndexBloatTable(db)
		if hasDataForDb {
			hasData = true
		}
	}

	// 根据是否有数据决定输出内容
	if hasData {
		fmt.Println("以下是数据库中索引膨胀相关信息：")
	} else {
		fmt.Println("未查询到数据库中索引膨胀相关信息")
	}

	// 打印建议
	fmt.Println("\n建议: ")
	fmt.Println("   > 如果索引膨胀太大, 会影响性能, 建议重建索引, create index CONCURRENTLY.... ")
	fmt.Println()
}

// printIndexBloatTable 打印指定数据库的索引膨胀情况表格
func printIndexBloatTable(db string) bool {
	// 创建用于当前数据库表格输出的对象并设置表头
	buffer := &bytes.Buffer{}
	writer := tablewriter.NewWriter(buffer)
	writer.SetAutoFormatHeaders(true)
	writer.SetHeader([]string{"数据库", "schema", "表名", "表膨胀系数", "索引名", "索引膨胀系数"})

	// 标记当前数据库是否获取到有效数据，初始化为false
	currentHasData := false

	// 构建psql命令以获取索引膨胀信息
	cmd := exec.Command("psql", "-d", db, "--pset=pager=off", "-t", "--pset=border=2", "-q", "-c", `
    select db,schemaname,tablename,tbloat,iname,ibloat from (
    SELECT
    current_database() AS db, schemaname, tablename, reltuples::bigint AS tups, relpages::bigint AS pages, otta,
    ROUND(CASE WHEN otta=0 OR sml.relpages=0 OR sml.relpages=otta THEN 0.0 ELSE sml.relpages/otta::numeric END,1) AS tbloat,
    CASE WHEN relpages < otta THEN 0 ELSE relpages::bigint - otta END AS wastedpages,
    CASE WHEN relpages < otta THEN 0 ELSE bs*(sml.relpages-otta)::bigint END AS wastedbytes,
    CASE WHEN relpages < otta THEN '0 bytes'::text ELSE (bs*(relpages-otta))::bigint || ' bytes' END AS wastedsize,
    iname, ituples::bigint AS itups, ipages::bigint AS ipages, iotta,
    ROUND(CASE WHEN iotta=0 OR ipages=0 OR ipages=iotta THEN 0.0 ELSE ipages/iotta::numeric END,1) AS ibloat,
    CASE WHEN ipages < iotta THEN 0 ELSE ipages::bigint - iotta END AS wastedipages,
    CASE WHEN ipages < iotta THEN 0 ELSE bs*(ipages-iotta) END AS wastedibytes,
    CASE WHEN ipages < iotta THEN '0 bytes' ELSE (bs*(ipages-iotta))::bigint || ' bytes' END AS wastedisize,
    CASE WHEN relpages < otta THEN
        CASE WHEN ipages < iotta THEN 0 ELSE bs*(ipages-iotta::bigint) END
        ELSE CASE WHEN ipages < iotta THEN bs*(relpages-otta::bigint)
        ELSE bs*(relpages-otta::bigint + ipages-iotta::bigint) END
    END AS totalwastedbytes
    FROM (
    SELECT
        nn.nspname AS schemaname,
        cc.relname AS tablename,
        COALESCE(cc.reltuples,0) AS reltuples,
        COALESCE(cc.relpages,0) AS relpages,
        COALESCE(bs,0) AS bs,
        COALESCE(CEIL((cc.reltuples*((datahdr+ma-
        (CASE WHEN datahdr%ma=0 THEN ma ELSE datahdr%ma END))+nullhdr2+4))/(bs-20::float)),0) AS otta,
        COALESCE(c2.relname,'?') AS iname, COALESCE(c2.reltuples,0) AS ituples, COALESCE(c2.relpages,0) AS ipages,
        COALESCE(CEIL((c2.reltuples*(datahdr-12))/(bs-20::float)),0) AS iotta -- very rough approximation, assumes all cols
    FROM
        pg_class cc
    JOIN pg_namespace nn ON cc.relnamespace = nn.oid AND nn.nspname <> 'information_schema'
    LEFT JOIN
    (
        SELECT
        ma,bs,foo.nspname,foo.relname,
        (datawidth+(hdr+ma-(case when hdr%ma=0 THEN ma ELSE hdr%ma END)))::numeric AS datahdr,
        (maxfracsum*(nullhdr+ma-(case when nullhdr%ma=0 THEN ma ELSE nullhdr%ma END))) AS nullhdr2
        FROM (
        SELECT
            ns.nspname, tbl.relname, hdr, ma, bs,
            SUM((1-coalesce(null_frac,0))*coalesce(avg_width, 2048)) AS datawidth,
            MAX(coalesce(null_frac,0)) AS maxfracsum,
            hdr+(
            SELECT 1+count(*)/8
            FROM pg_stats s2
            WHERE null_frac<>0 AND s2.schemaname = ns.nspname AND s2.tablename = tbl.relname
            ) AS nullhdr
        FROM pg_attribute att 
        JOIN pg_class tbl ON att.attrelid = tbl.oid
        JOIN pg_namespace ns ON ns.oid = tbl.relnamespace 
        LEFT JOIN pg_stats s ON s.schemaname=ns.nspname
        AND s.tablename = tbl.relname
        AND s.inherited=false
        AND s.attname=att.attname,
        (
            SELECT
            (SELECT current_setting('block_size')::numeric) AS bs,
                CASE WHEN SUBSTRING(SPLIT_PART(v, ' ', 2) FROM '#"[0-9]+.[0-9]+#"%' for '#')
                IN ('8.0','8.1','8.2') THEN 27 ELSE 23 END AS hdr,
            CASE WHEN v ~ 'mingw32' OR v ~ '64-bit' THEN 8 ELSE 4 END AS ma
            FROM (SELECT version() AS v) AS foo
        ) AS constants
        WHERE att.attnum > 0 AND tbl.relkind='r'
        GROUP BY 1,2,3,4,5
        ) AS foo
    ) AS rs
    ON cc.relname = rs.relname AND nn.nspname = rs.nspname
    LEFT JOIN pg_index i ON indrelid = cc.oid
    LEFT JOIN pg_class c2 ON c2.oid = i.indexrelid
    ) AS sml order by wastedibytes desc limit 5 ) as aa;`)
	var result bytes.Buffer
	cmd.Stdout = &result
	cmd.Stderr = &bytes.Buffer{} // 用于捕获错误信息

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute command for database %s: %s\n", db, err)
		return false
	}

	// 使用正则表达式提取每行的数据（可根据实际数据格式调整正则表达式）
	lines := strings.Split(strings.TrimSpace(result.String()), "\n")
	for _, line := range lines {
		re := regexp.MustCompile(`\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 7 { // 第一个匹配项是完整的匹配项，后面是列的数据
			database := strings.TrimSpace(matches[1])
			schema := strings.TrimSpace(matches[2])
			tableName := strings.TrimSpace(matches[3])
			tableBloatFactor := strings.TrimSpace(matches[4])
			indexName := strings.TrimSpace(matches[5])
			indexBloatFactor := strings.TrimSpace(matches[6])

			if database != "" || schema != "" || tableName != "" || tableBloatFactor != "" || indexName != "" || indexBloatFactor != "" {
				writer.Append([]string{
					database,
					schema,
					tableName,
					tableBloatFactor,
					indexName,
					indexBloatFactor,
				})
				currentHasData = true
			}
		}
	}

	if currentHasData {
		writer.Render()
		fmt.Println(buffer.String())
	}

	return currentHasData
}
