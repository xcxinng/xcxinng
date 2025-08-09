[现状]
一亿多条数据，无任何索引，大约12G磁盘空间，大约是18年至今没清理过任何数据

```sql
create table device_history (
	id varchar(16) not null default uuid_v2_gen();
	device_id varchar(16) not null,
	category smallint not null,
	content_before text,
	content_after text,
	created_at timestamp(6) without time zone default now(),
	created_by varchar(32)
)
```

查询需求: 根据指定device_id 过滤出来所有日志，然后根据created_by倒序排列,
```sql
select * from device_history where device_id='xx' order by created_at desc limit 20 offset xx
```

[优化方法]
