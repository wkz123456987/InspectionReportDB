

-- 获取非template数据库名称
[QUERY_NON_TEMPLATE_DBS]
SELECT datname FROM pg_database WHERE datname NOT IN ('template0', 'template1')

-- 获取指定数据库中符合条件的表大小信息（每个数据库取前10）
[QUERY_TABLE_SIZE_INFO]
SELECT current_database(), b.nspname, c.relname, c.relkind, pg_size_pretty(pg_relation_size(c.oid)) 
FROM pg_stat_all_tables a, pg_class c, pg_namespace b 
WHERE pg_relation_size(c.oid) >= 10 AND c.relnamespace = b.oid AND c.relkind = 'r' AND a.relid = c.oid 
ORDER BY pg_relation_size(c.oid) DESC LIMIT 10

-- 获取指定数据库中索引数超过4且SIZE大于10MB的表信息
[QUERY_TABLES_WITH_TOO_MANY_INDEXES]
SELECT current_database(), t2.nspname, t1.relname, pg_size_pretty(pg_relation_size(t1.oid)), t3.idx_cnt 
FROM pg_class t1, pg_namespace t2, (SELECT indrelid, COUNT(*) idx_cnt FROM pg_index GROUP BY 1 HAVING COUNT(*) > 4) t3 
WHERE pg_relation_size(t1.oid) >= 10000000 AND t1.oid = t3.indrelid AND t1.relnamespace = t2.oid AND pg_relation_size(t1.oid) / 1024 / 1024.0 > 10 
ORDER BY t3.idx_cnt DESC

-- 获取指定数据库中重复索引信息
[QUERY_REPEAT_INDEX_INFO]
SELECT current_database() as dbname,indrelid::regclass AS TableName, array_agg(indexrelid::regclass) AS Indexes 
FROM pg_index 
GROUP BY indrelid, indkey 
HAVING COUNT(*) > 1

-- 获取指定数据库中未使用或使用较少的索引信息
[QUERY_UNUSED_INDEXES_INFO]
SELECT current_database(), t2.schemaname, t2.relname, t2.indexrelname
FROM pg_stat_all_tables t1, pg_stat_all_indexes t2 
WHERE t1.relid = t2.relid 
AND t2.idx_scan < 10 
AND t2.schemaname NOT IN ('pg_toast', 'pg_catalog') 
AND indexrelid NOT IN (SELECT conindid FROM pg_constraint WHERE contype IN ('p', 'u', 'f')) 
AND pg_relation_size(indexrelid) > 65536 
ORDER BY pg_relation_size(indexrelid) DESC

-- 获取数据库统计信息
[QUERY_DATABASE_STATS]
SELECT 
    datname,
    ROUND(100 * (xact_rollback::numeric / (CASE WHEN xact_commit > 0 THEN xact_commit ELSE 1 END + xact_rollback)), 2) || ' %' AS rollback_ratio,
    ROUND(100 * (blks_hit::numeric / (CASE WHEN blks_read > 0 THEN blks_read ELSE 1 END + blks_hit)), 2) || ' %' AS hit_ratio,
    blk_read_time,
    blk_write_time,
    conflicts,
    deadlocks 
FROM pg_stat_database;

-- 获取指定数据库中索引膨胀信息
[QUERY_INDEX_BLOAT_INFO]
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
    ) AS sml order by wastedibytes desc limit 5 ) as aa;


-- 获取指定数据库中垃圾数据信息
[QUERY_GARBAGE_DATA_INFO]
SELECT current_database(), schemaname, relname, n_dead_tup
FROM pg_stat_all_tables
WHERE n_live_tup > 0 AND n_dead_tup / n_live_tup > 0.2 AND schemaname NOT IN ('pg_toast', 'pg_catalog')
ORDER BY n_dead_tup DESC
LIMIT 5

-- 获取数据库年龄信息
[QUERY_DATABASE_AGE_INFO]
SELECT datname, age(datfrozenxid), 2^31 - age(datfrozenxid) AS age_remain
FROM pg_database
ORDER BY age(datfrozenxid) DESC

-- 获取指定数据库中表年龄信息
[QUERY_TABLE_AGE_INFO]
SELECT current_database(), rolname, nspname, relkind, relname, age(relfrozenxid), 2^31 - age(relfrozenxid) AS age_remain
FROM pg_authid t1
JOIN pg_class t2 ON t1.oid = t2.relowner
JOIN pg_namespace t3 ON t2.relnamespace = t3.oid
WHERE t2.relkind IN ('t', 'r')
ORDER BY age(relfrozenxid) DESC
LIMIT 5

-- 获取锁等待信息
[QUERY_LOCK_WAIT_INFO]
WITH t_wait AS (
    SELECT a.mode, a.locktype, a.database, a.relation, a.page, a.tuple, a.classid,
           a.objid, a.objsubid, a.pid, a.virtualtransaction, a.virtualxid, a.transactionid,
           b.query, b.xact_start, b.query_start, b.usename, b.datname 
    FROM pg_locks a, pg_stat_activity b 
    WHERE a.pid = b.pid AND NOT a.granted
),
t_run AS (
    SELECT a.mode, a.locktype, a.database, a.relation, a.page, a.tuple,
           a.classid, a.objid, a.objsubid, a.pid, a.virtualtransaction, a.virtualxid,
           a.transactionid, b.query, b.xact_start, b.query_start,
           b.usename, b.datname 
    FROM pg_locks a, pg_stat_activity b 
    WHERE a.pid = b.pid AND a.granted
)
SELECT r.locktype, r.mode AS r_mode, r.usename AS r_user, r.datname AS r_db,
       r.relation::regclass, r.pid AS r_pid,
       r.page AS r_page, r.tuple AS r_tuple, r.xact_start AS r_xact_start,
       r.query_start AS r_query_start,
       now() - r.query_start AS r_locktime, r.query AS r_query,
       w.mode AS w_mode, w.pid AS w_pid, w.page AS w_page,
       w.tuple AS w_tuple, w.xact_start AS w_xact_start, w.query_start AS w_query_start,
       now() - w.query_start AS w_locktime, w.query AS w_query  
FROM t_wait w, t_run r 
WHERE r.locktype IS NOT DISTINCT FROM w.locktype
  AND r.database IS NOT DISTINCT FROM w.database
  AND r.relation IS NOT DISTINCT FROM w.relation
  AND r.page IS NOT DISTINCT FROM w.page
  AND r.tuple IS NOT DISTINCT FROM w.tuple
  AND r.classid IS NOT DISTINCT FROM w.classid
  AND r.objid IS NOT DISTINCT FROM w.objid
  AND r.objsubid IS NOT DISTINCT FROM w.objsubid
  AND r.transactionid IS NOT DISTINCT FROM w.transactionid
  AND r.pid <> w.pid
ORDER BY ((
            CASE w.mode
                WHEN 'INVALID' THEN 0
                WHEN 'AccessShareLock' THEN 1
                WHEN 'RowShareLock' THEN 2
                WHEN 'RowExclusiveLock' THEN 3
                WHEN 'ShareUpdateExclusiveLock' THEN 4
                WHEN 'ShareLock' THEN 5
                WHEN 'ShareRowExclusiveLock' THEN 6
                WHEN 'ExclusiveLock' THEN 7
                WHEN 'AccessExclusiveLock' THEN 8
                ELSE 0
            END
          ) + (
            CASE r.mode
                WHEN 'INVALID' THEN 0
                WHEN 'AccessShareLock' THEN 1
                WHEN 'RowShareLock' THEN 2
                WHEN 'RowExclusiveLock' THEN 3
                WHEN 'ShareUpdateExclusiveLock' THEN 4
                WHEN 'ShareLock' THEN 5
                WHEN 'ShareRowExclusiveLock' THEN 6
                WHEN 'ExclusiveLock' THEN 7
                WHEN 'AccessExclusiveLock' THEN 8
                ELSE 0
            END
          )) DESC, r.xact_start;


-- 检查pg_authid中密码相关情况
[QUERY_PG_AUTHID_CHECK]
SELECT count(*) FROM pg_authid WHERE rolpassword!~ '^md5' OR length(rolpassword) <> 35

-- 获取非template数据库名称
[QUERY_NON_TEMPLATE_DBS]
SELECT datname FROM pg_database WHERE datname NOT IN ('template0', 'template1')

-- 检查pg_user_mappings中密码相关情况
[QUERY_PG_USER_MAPPINGS_CHECK]
SELECT current_database(), * FROM pg_user_mappings WHERE umoptions::text ~* 'password'

-- 检查pg_views中密码相关情况
[QUERY_PG_VIEWS_CHECK]
SELECT current_database(), * FROM pg_views WHERE definition ~* 'password' AND definition ~* 'dblink'

-- 获取复制槽状态信息
[QUERY_REPLICATION_SLOT_STATUS_INFO]
SELECT slot_name, slot_type, active
FROM pg_replication_slots
ORDER BY 3


-- 获取指定数据库中schema统计信息
[QUERY_SCHEMA_STATS]
SELECT schemaName as "schemaName",
       sum(total_size) as "Byte",
       round(sum(total_size)/1024/1024,1) as "MB",
       round(sum(total_size)/1024/1024/1024,1) as "GB" 
FROM (
    SELECT nspname as schemaName,
           pg_total_relation_size(pg_class.oid) as total_size 
    FROM pg_class 
    JOIN  pg_namespace ON (pg_namespace.oid = pg_class.relnamespace) 
    WHERE relkind IN ('r', 'v', 'm', 'S', 'f') 
    ORDER BY total_size DESC
) as aa 
GROUP BY schemaName 
ORDER BY 4 desc;











