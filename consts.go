package schema

const (
	commandCreate       = "create"
	commandAdd          = "add"
	commandChange       = "change"
	commandDrop         = "drop"
	commandDropColumn   = "dropColumn"
	commandRename       = "rename"
	commandRenameColumn = "renameColumn"
	commandDropIfExists = "dropIfExists"
	commandPrimary      = "primary"
	commandUnique       = "unique"
	commandIndex        = "index"
	commandComment      = "tableComment"

	commandAttrIndex     = "index"
	commandAttrAlgorithm = "algorithm"
	commandAttrColumns   = "columns" // []string
	commandAttrComment   = "comment"
	commandAttrFrom      = "from" // rename from
	commandAttrTo        = "to"   // rename to
)

const (
	ColumnTypeChar       = "char"
	ColumnTypeVarchar    = "varchar"
	ColumnTypeTinyText   = "tinytext"
	ColumnTypeText       = "text"
	ColumnTypeMediumText = "mediumtext"
	ColumnTypeLongText   = "longtext"
	ColumnTypeInt        = "int"
	ColumnTypeTinyInt    = "tinyint"
	ColumnTypeBoolean    = "boolean" // tinyint(1)
	ColumnTypeSmallInt   = "smallint"
	ColumnTypeMediumInt  = "mediumint"
	ColumnTypeBigInt     = "bigint"
	ColumnTypeFloat      = "float"
	ColumnTypeDouble     = "double"
	ColumnTypeDecimal    = "decimal"
	ColumnTypeEnum       = "enum"
	ColumnTypeSet        = "set"
	ColumnTypeJson       = "json"
	ColumnTypeDate       = "date"
	ColumnTypeDateTime   = "datetime"
	ColumnTypeTime       = "time"
	ColumnTypeTimestamp  = "timestamp"
	ColumnTypeYear       = "year"
	ColumnTypeBinary     = "binary"
	ColumnTypeBlob       = "blob"
	ColumnTypeUuid       = "uuid"

	ColumnAttrPrimary       = "primary"
	ColumnAttrUnique        = "unique"
	ColumnAttrIndex         = "index"
	ColumnAttrComment       = "comment"       // 字段注释
	ColumnAttrDefault       = "default"       // 字段默认值
	ColumnAttrNullable      = "nullable"      // 是否可为空 bool
	ColumnAttrLength        = "length"        // 字段长度
	ColumnAttrTotal         = "total"         // 浮点数字段小数
	ColumnAttrPlaces        = "places"        // 浮点数小数位
	ColumnAttrAllowed       = "allowed"       // enum allowed []string
	ColumnAttrChange        = "change"        // 是否是修改字段 bool
	ColumnAttrAutoIncrement = "autoIncrement" // 是否自动递增 bool
	ColumnAttrUnsigned      = "unsigned"      // 是否无符号 bool
	ColumnAttrCharset       = "charset"       // 字符集
	ColumnAttrCollate       = "collate"       // 排序规则
)

const (
	DefaultPrefix       = ""                   // 数据库表前缀
	DefaultEngine       = "InnoDB"             // 默认数据表引擎
	DefaultCharset      = "utf8mb4"            // 默认编码
	DefaultCollation    = "utf8mb4_unicode_ci" // 默认排序
	DefaultStringLength = 255                  // string 字段默认长度
)

// Map alias
type Map = map[string]interface{}

type Command struct {
	Name       string
	Attributes Map
}
