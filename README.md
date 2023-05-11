## Go mysql schema

数据库迁移组件，可自动生成表与字段

## 使用

```shell
$ go get github.com/chenpkg/schema
```

```go
package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/chenpkg/schema"
)

func main() {
	// db handle
	var db *sql.DB
	
	// use goframe orm
	db, err := g.DB().Open(g.DB().GetConfig())
	
	// or use gorm
	// db, err := gorm.Open(mysql.open("..."), &gorm.Config{}).DB()
	
	// or use go-sql-driver/mysql
	// import _ "github.com/go-sql-driver/mysql"
	// db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/test")

	config := &schema.Config{
		DB:       db,
		Database: "test",
		Prefix:   "",
	}

	dbSchema := schema.NewSchema(context.Background(), config)

	// create users table
	err := dbSchema.Create("users", func(table *schema.Blueprint) {
		table.Id("id")
		table.String("account", 30).Comment("账号")
		table.String("password", 100).Comment("密码")
		table.String("name", 20).Default("").Comment("姓名")
		table.UnsignedInt("age").Default(18).Comment("年龄")
		table.Timestamp("birthday").Nullable().Comment("生日")

		// create created_at and updated_at column
		table.Timestamps()
		// create timestamp deleted_at column
		table.SoftDeletes()

		// add index
		table.Index("account")
		// add join index
		table.Index([]string{"account", "name"})
		// add unique index
		table.Unique("account")
		// or table.String("account", 30).Unique()
	})
	if err != nil {
		log.Fatal(err)
	}

	// check table exists 
	exists, err := dbSchema.HasTable("users")
}

```

`Config` 参数默认配置

```go
const (
	DefaultPrefix       = ""                   // 数据库表前缀
	DefaultEngine       = "InnoDB"             // 默认数据表引擎
	DefaultCharset      = "utf8mb4"            // 默认编码
	DefaultCollation    = "utf8mb4_unicode_ci" // 默认排序
	DefaultStringLength = 255                  // string 字段默认长度
)
```

## 数据表

可使用 `Engine` 指定表的储存引擎

    table.Engine = "InnoDB"

`Charset` 和 `Collation` 属性可用于在使用 MySQL 时为创建的表指定字符集和排序规则

    table.Charset = "utf8mb4"
    table.Collation = "utf8mb4_unicode_ci"

要删除已存在的表，可以使用 `Drop` 或 `DropIfExists` 方法

    dbSchema.Drop()

    dbSchema.DropIfExists()

使用 `HasTable` 判断表是否存在

    dbSchema.HasTable("users")

要重命名已存在的数据表，使用 rename 方法

    dbSchema.Rename("users", "new_users")

## 字段

下面列出了所有可用字段类型的方法：

```go
// 创建 big int 主键
table.Id()
// 同上
table.BigIncrements()
// int primary
table.Increments()

// char
table.Char()
table.Varchar()
table.TinyText()
table.Text()
table.MediumText()
table.LongText()
// string is Varchar alias
table.String()
// Uuid is char(36) alias
table.Uuid()

// integer
table.Int()
table.TinyInt()
// Boolean is tinyint(1) alias
table.Boolean()
table.SmallInt()
table.MediumInt()
table.BigInt()

// unsigned int 无符号 int
table.UnsignedInt()
table.UnsignedTinyInt()
table.UnsignedSmallInt()
table.UnsignedMediumInt()
table.UnsignedBigInt()

// float
table.Float("column", 10, 2)
table.Double()
table.Decimal()

// unsigned float
table.UnsignedFloat()
table.UnsignedDouble()
table.UnsignedDecimal()

// enum
table.Enum("column", []string{"one", "two", "three"})
table.Set("column", []string{"one", "two", "three"})

// json
table.Json()

// date
table.Date()
table.DateTime()
table.Time()
table.Timestamp()
// 同时创建 created_at, updated_at 创建时间与修改时间字段
table.Timestamps()
// 创建 deleted_at 删除时间字段
table.SoftDeletes()
table.Year()

// binary
table.Binary()
table.Blob()
```

### 字段修饰符

除了上面列出的列类型外，在向数据库表添加列时还有几个可以使用的「修饰符」。例如，如果要把列设置为要使列为「可空」，你可以使用 `Nullable` 方法：

```go
dbSchema.Create("users", func(table *schema.Blueprint) {
	table.String("email").Nullable()
})
```

下面是所有可用的列表修饰符，此列表不包含索引修饰符

```go
// 为该列指定字符集
Charset()

// 为该列指定排序规则
Collation()

// 为该列添加注释
Comment()

// 为该列指定一个「默认值」
Default()

// 运行 NULL 值插入到该列
Nullable()
```

### 修改字段

`Change` 方法可以将现有的字段类型修改为新的类型或修改属性。比如，你可能想增加 `string` 字段的长度，可以使用 `Change` 方法把 `name` 字段的长度从 25 增加到 50。所以，我们可以简单的更新字段属性然后调用 `Change` 方法：

```go
dbSchema.Create("users", func(table *schema.Blueprint) {
	table.String("name", 50).Change()
})
```

### 删除字段

```go
// 删除一个字段
dbSchema.DropColumns("users", "account")
// 或者删除多个
dbSchema.DropColumns("users", "account", "password", "age")
```

## 索引

### 创建索引

下面的例子中新建了一个值唯一的 email 字段。我们可以将 `unique` 方法链式地添加到字段定义上来创建索引：

```go
dbSchema.Create("users", func(table *schema.Blueprint) {
	table.String("email", 50).Unique()
})
```

或者，你也可以在定义完字段之后创建索引。为此，你应该调用结构生成器上的 unique 方法，此方法应该传入唯一索引的列名称：

    table.Unique("email")

或者创建复合（或合成）索引：

    table.Index([]string{"account", "name"})

下面是可用的索引类型:

```go
// 添加主键
Primary("name")

// 添加复合主键
Primary([]string{"id", "name"})

// 添加唯一索引
Unique()

// 添加普通索引
Index()
```