-- 获取数据库当前活跃度状态信息
[QUERY_ACTIVITY_STATUS]
select now(),state,count(*) from pg_stat_activity group by 1,2;

-- 获取数据库连接相关信息
[QUERY_DB_CONNECTIONS]
SELECT 
    max_conn,
    used,
    res_for_super,
    max_conn - used - res_for_super AS res_for_normal
FROM 
    (SELECT count(*) AS used FROM pg_stat_activity) t1,
    (SELECT setting::int AS res_for_super FROM pg_settings WHERE name ='superuser_reserved_connections') t2,
    (SELECT setting::int AS max_conn FROM pg_settings WHERE name ='max_connections') t3

-- 查询数据库版本
[QUERY_DB_VERSION]
select version();

-- 获取非template数据库名称
[QUERY_NON_TEMPLATE_DBS]
select datname from pg_database where datname not in ('template0', 'template1')

-- 获取插件版本信息
[QUERY_PLUGIN_VERSIONS]
select 
    current_database(), 
    e.extname, 
    u.usename as extowner, 
    n.nspname as extnamespace, 
    e.extrelocatable, 
    e.extversion 
from 
    pg_extension e
left join 
    pg_user u on e.extowner = u.usesysid
left join 
    pg_namespace n on e.extnamespace = n.oid;


-- 获取数据类型统计信息
[QUERY_USED_DATA_TYPE_COUNTS]
select 
    current_database(),  -- 获取当前数据库名称
    b.typname,            -- 获取数据类型名称
    count(*)              -- 统计对应数据类型的数量
from 
    pg_attribute a,
    pg_type b
where 
    a.atttypid = b.oid
    and a.attrelid in (
        select oid 
        from pg_class 
        where relnamespace not in (
            select oid 
            from pg_namespace 
            where nspname ~ '^pg_' or nspname = 'information_schema'
        )
    )
group by 
    1, 2
order by 
    3 desc;

-- 获取用户创建对象统计信息
[QUERY_CREATED_OBJECT_COUNTS]
select current_database(),rolname,nspname,relkind,count(*) from pg_class a,pg_authid b,pg_namespace c where a.relnamespace=c.oid and a.relowner=b.oid and nspname!~ '^pg_' and nspname<>'information_schema' group by 1,2,3,4 order by 5 desc;

-- 获取用户对象占用空间信息
[QUERY_USER_OBJECT_SPACE_INFO]
select current_database(), buk this_buk_no, cnt rels_in_this_buk, pg_size_pretty(min) buk_min, pg_size_pretty(max) buk_max
from (
    select row_number() over (partition by buk order by tsize), tsize, buk, min(tsize) over (partition by buk), max(tsize) over (partition by buk), count(*) over (partition by buk) cnt
    from (
        select pg_relation_size(a.oid) tsize, width_bucket(pg_relation_size(a.oid), tmin - 1, tmax + 1, 10) buk
        from (
            select min(pg_relation_size(a.oid)) tmin, max(pg_relation_size(a.oid)) tmax
            from pg_class a, pg_namespace c
            where a.relnamespace = c.oid and nspname!~ '^pg_' and nspname <>'information_schema'
        ) t,
        pg_class a, pg_namespace c
        where a.relnamespace = c.oid and nspname!~ '^pg_' and nspname <>'information_schema'
    ) t
) t
where row_number = 1;

-- 获取表空间使用情况信息
[QUERY_TABLESPACE_USAGE]
select spcname,pg_tablespace_location(oid),pg_size_pretty(pg_tablespace_size(oid)) from pg_tablespace order by pg_tablespace_size(oid) desc;

-- 获取数据库使用情况
[QUERY_DATABASE_USAGE]
select datname,pg_size_pretty(pg_database_size(oid)) from pg_database order by pg_database_size(oid) desc;

-- 获取用户连接数限制相关信息
[QUERY_USER_CONNECTION_LIMITS]
SELECT a.rolname, a.rolconnlimit, b.connects
FROM pg_authid a
JOIN (
    SELECT usename, COUNT(*) AS connects
    FROM pg_stat_activity
    GROUP BY usename
) b ON a.rolname = b.usename
ORDER BY b.connects DESC;


-- 获取数据库连接限制相关信息
[QUERY_DATABASE_CONNECTION_LIMITS]
SELECT a.datname, a.datconnlimit, b.connects
FROM pg_database a
JOIN (
    SELECT datname, COUNT(*) AS connects
    FROM pg_stat_activity
    GROUP BY datname
) b ON a.datname = b.datname
ORDER BY b.connects DESC;

-- 获取检查点、bgwriter统计信息
[QUERY_CHECKPOINT_BGWRITER_STATS]
select * from pg_stat_bgwriter;

-- 获取长事务相关信息
[QUERY_LONG_TRANSACTION_INFO]
SELECT datname, usename, query, xact_start, (now() - xact_start) AS xact_duration, query_start, (now() - query_start) AS query_duration, state
FROM pg_stat_activity
WHERE state <> 'idle' AND (backend_xid IS NOT NULL OR backend_xmin IS NOT NULL) AND (now() - xact_start > interval '3 s')
ORDER BY xact_start;

-- 获取2PC相关信息
[QUERY_2PC_INFO]
SELECT transaction, gid, prepared, owner, database, (now() - prepared) AS duration
FROM pg_prepared_xacts
WHERE (now() - prepared) > INTERVAL '3 s'
ORDER BY prepared;

-- 获取用户密码到期时间信息
[QUERY_USER_PASSWORD_EXPIRATION]
SELECT rolname, 
       (CASE WHEN rolvaliduntil IS NULL THEN '无有效期' ELSE rolvaliduntil::text END) AS rolvaliduntil
FROM pg_authid
ORDER BY (CASE WHEN rolvaliduntil IS NULL THEN '9999-12-31 23:59:59.999999+00' ELSE rolvaliduntil END);

-- 获取继承关系信息
[QUERY_INHERITANCE_RELATION_INFO]
SELECT inhrelid::regclass, inhparent::regclass, inhseqno FROM pg_inherits ORDER BY 2, 3


-- 获取是否开启归档、自动垃圾回收设置信息
[QUERY_ARCHIVE_AND_AUTOVACUUM_SETTINGS]
SELECT name, setting FROM pg_settings WHERE name IN ('archive_mode', 'autovacuum', 'archive_command');


-- 获取主备库角色信息
[QUERY_MASTER_STANDBY_ROLE]
SELECT CASE WHEN pg_is_in_recovery() = 'f' THEN '主库' WHEN pg_is_in_recovery() = 't' THEN '备库' END AS database_role;

-- 获取备库信息
[QUERY_STANDBY_INFO]
SELECT usename, application_name, client_addr FROM pg_stat_replication;

