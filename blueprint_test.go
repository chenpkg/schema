package schema

import (
	"context"
	"testing"
)

func TestBlueprint_ToSql(t *testing.T) {
	type sqlCase struct {
		name     string
		table    string
		sql      []string
		callback func(table *Blueprint)
	}
	cases := []sqlCase{
		{
			name:  "set engine charset collation",
			table: "users",
			sql: []string{
				"create table `users` (`id` bigint unsigned not null auto_increment primary key) default character set utf8 collate 'utf8_unicode_ci' engine = MyISAM",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Id("id")
				table.Engine = "MyISAM"
				table.Charset = "utf8"
				table.Collation = "utf8_unicode_ci"
			},
		},
		{
			name:  "add table comment",
			table: "users",
			sql: []string{
				"create table `users` (`id` bigint unsigned not null auto_increment primary key) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
				"alter table `users` comment = '用户表'",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Id("id")
				table.Comment("用户表")
			},
		},
		{
			name:  "drop table",
			table: "users",
			sql: []string{
				"drop table `users`",
			},
			callback: func(table *Blueprint) {
				table.drop()
			},
		},
		{
			name:  "dropIfExists table",
			table: "users",
			sql: []string{
				"drop table if exists `users`",
			},
			callback: func(table *Blueprint) {
				table.dropIfExists()
			},
		},
		{
			name:  "rename table",
			table: "users",
			sql: []string{
				"rename table `users` to `new_user`",
			},
			callback: func(table *Blueprint) {
				table.Rename("new_user")
			},
		},
		{
			name:  "Id",
			table: "users",
			sql: []string{
				"create table `users` (`id` bigint unsigned not null auto_increment primary key) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Id("id")
			},
		},
		{
			name:  "Increments",
			table: "users",
			sql: []string{
				"create table `users` (`id` int unsigned not null auto_increment primary key) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Increments("id")
			},
		},
		{
			name:  "String",
			table: "users",
			sql: []string{
				"create table `users` (`name` char(20) not null, `name2` varchar(20) not null, `name3` tinytext not null, `name4` mediumtext not null, `name5` longtext not null, `name6` varchar(255) not null, `uuid` char(36) not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Char("name", 20)
				table.Varchar("name2", 20)
				table.TinyText("name3")
				table.MediumText("name4")
				table.LongText("name5")
				table.String("name6", 255)
				table.Uuid("uuid")
			},
		},
		{
			name:  "Int",
			table: "users",
			sql: []string{
				"create table `users` (`age` int not null, `age2` tinyint not null, `age3` tinyint(1) not null, `age4` smallint not null, `age5` mediumint not null, `age6` bigint not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Int("age")
				table.TinyInt("age2")
				table.Boolean("age3")
				table.SmallInt("age4")
				table.MediumInt("age5")
				table.BigInt("age6")
			},
		},
		{
			name:  "UnsignedInt",
			table: "users",
			sql: []string{
				"create table `users` (`age` int unsigned not null, `age2` tinyint unsigned not null, `age3` smallint unsigned not null, `age4` mediumint unsigned not null, `age5` bigint unsigned not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.UnsignedInt("age")
				table.UnsignedTinyInt("age2")
				table.UnsignedSmallInt("age3")
				table.UnsignedMediumInt("age4")
				table.UnsignedBigInt("age5")
			},
		},
		{
			name:  "Float",
			table: "users",
			sql: []string{
				"create table `users` (`age` double(10, 2) not null, `age2` double(8, 2) not null, `age3` decimal(10, 1) not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Float("age", 10, 2)
				table.Double("age2", 8, 2)
				table.Decimal("age3", 10, 1)
			},
		},
		{
			name:  "UnsignedFloat",
			table: "users",
			sql: []string{
				"create table `users` (`age` double(10, 2) unsigned not null, `age2` double(8, 2) unsigned not null, `age3` decimal(10, 1) unsigned not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.UnsignedFloat("age", 10, 2)
				table.UnsignedDouble("age2", 8, 2)
				table.UnsignedDecimal("age3", 10, 1)
			},
		},
		{
			name:  "Enum",
			table: "users",
			sql: []string{
				"create table `users` (`type` enum('one', 'two', 'three') not null, `type2` set('first', 'second', 'third') not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Enum("type", []string{"one", "two", "three"})
				table.Set("type2", []string{"first", "second", "third"})
			},
		},
		{
			name:  "Json",
			table: "users",
			sql: []string{
				"create table `users` (`data` json not null, `binary` binary not null, `blob` binary not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Json("data")
				table.Binary("binary")
				table.Blob("blob")
			},
		},
		{
			name:  "Date",
			table: "users",
			sql: []string{
				"create table `users` (`date` date not null, `date_time` datetime not null, `time` time not null, `timestamp` timestamp not null, `created_at` timestamp null, `updated_at` timestamp null, `deleted_at` timestamp null, `year` year not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.Date("date")
				table.DateTime("date_time")
				table.Time("time")
				table.Timestamp("timestamp")
				table.Timestamps()
				table.SoftDeletes()
				table.Year("year")
			},
		},
		{
			name:  "Column_Charset_Collation",
			table: "users",
			sql: []string{
				"create table `users` (`name` varchar(255) character set utf8mb4 collate 'utf8mb4_unicode_ci' not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.String("name").Charset("utf8mb4").Collation("utf8mb4_unicode_ci")
			},
		},
		{
			name:  "Column_Comment_Default_Nullable",
			table: "users",
			sql: []string{
				"create table `users` (`name` varchar(30) not null default '' comment 'name', `data` text null comment 'comment data') default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.String("name", 30).Default("").Comment("name")
				table.Text("data").Nullable().Comment("comment data")
			},
		},
		{
			name:  "Column_Update_Delete",
			table: "users",
			sql: []string{
				"alter table `users` drop `age`, drop `account`",
				"alter table `users` modify `name` varchar(30) not null default 'a' comment '姓名'",
			},
			callback: func(table *Blueprint) {
				table.DropColumn("age", "account")
				table.String("name", 30).Default("a").Comment("姓名").Change()
			},
		},
		{
			name:  "Column_Index",
			table: "users",
			sql: []string{
				"create table `users` (`id` int unsigned not null, `account` varchar(50) not null, `name` varchar(30) not null, `age` int not null) default character set utf8mb4 collate 'utf8mb4_unicode_ci' engine = InnoDB",
				"alter table `users` add unique users_account_name_unique(`account`, `name`)",
				"alter table `users` add primary key users_id_primary(`id`)",
				"alter table `users` add index users_age_index(`age`)",
			},
			callback: func(table *Blueprint) {
				table.create()
				table.UnsignedInt("id").Primary()
				table.String("account", 50)
				table.String("name", 30)
				table.Int("age").Index()

				table.Unique([]string{"account", "name"})
			},
		},
	}

	newSchema := NewSchema(context.Background(), &Config{
		Engine:       "InnoDB",
		Charset:      "utf8mb4",
		Collation:    "utf8mb4_unicode_ci",
		StringLength: 255,
	})

	for _, item := range cases {
		blueprint := NewBlueprint(newSchema, item.table, item.callback)
		sql := blueprint.ToSql(localGrammar)
		if len(sql) != len(item.sql) {
			t.Fatal("ToSql err:", item.name, "\nsql:", item.sql, "\ngen:", sql)
		}
		for i, s := range item.sql {
			if sql[i] != s {
				t.Fatal("ToSql err:", item.name, "\nsql:", item.sql, "\ngen:", sql)
			}
		}
	}
}
