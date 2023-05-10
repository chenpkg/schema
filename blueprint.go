package schema

import (
	"context"
	"errors"
	"strings"
)

type Blueprint struct {
	Prefix    string // database prefix
	Engine    string // engin, default InnoDB
	Charset   string // charset, default utf8mb4
	Collation string // collation, default utf8mb4_unicode_ci

	table    string
	columns  []*Column  // The columns that should be added to the table.
	commands []*Command // The commands that should be run for the table.
	config   *Config
	ctx      context.Context
}

// NewBlueprint generate blueprint
func NewBlueprint(schema *Schema, table string, callback ...func(b *Blueprint)) *Blueprint {
	blueprint := &Blueprint{
		Prefix:    schema.config.Prefix,
		Engine:    ternary(schema.config.Engine == "", DefaultEngine, schema.config.Engine),
		Charset:   ternary(schema.config.Charset == "", DefaultCharset, schema.config.Charset),
		Collation: ternary(schema.config.Collation == "", DefaultCollation, schema.config.Collation),

		table:  table,
		config: schema.config,
		ctx:    schema.ctx,
	}

	if len(callback) > 0 {
		callback[0](blueprint)
	}

	return blueprint
}

// AddColumn add new column
func (b *Blueprint) AddColumn(types string, name string, parameters ...Map) (column *Column) {
	params := varDef(parameters, Map{})

	column = &Column{
		Type:       types,
		Name:       name,
		Attributes: params,
	}

	b.columns = append(b.columns, column)

	return column
}

// DropColumn Indicate that the given columns should be dropped.
func (b *Blueprint) DropColumn(columns ...string) *Command {
	if len(columns) == 0 {
		return nil
	}

	return b.addCommand(commandDropColumn, Map{
		commandAttrColumns: columns,
	})
}

// Rename the table to given name
func (b *Blueprint) Rename(to string) *Command {
	return b.addCommand(commandRename, Map{
		commandAttrTo: to,
	})
}

// RenameColumn Indicate that the given columns should be renamed.
// func (b *Blueprint) RenameColumn(from, to string) *Command {
// 	return b.addCommand(commandRenameColumn, Map{
// 		commandAttrFrom: from,
// 		commandAttrTo:   to,
// 	})
// }

// Primary add primary index
func (b *Blueprint) Primary(columns interface{}, algorithm ...string) {
	b.indexCommand(commandPrimary, columns, algorithm...)
}

// Unique add unique column
func (b *Blueprint) Unique(columns interface{}, algorithm ...string) {
	b.indexCommand(commandUnique, columns, algorithm...)
}

// Index add index
func (b *Blueprint) Index(columns interface{}, algorithm ...string) {
	b.indexCommand(commandIndex, columns, algorithm...)
}

// Id id primary
func (b *Blueprint) Id(column ...string) *Column {
	return b.BigIncrements(varDef(column, "id"))
}

// Increments index
func (b *Blueprint) Increments(column string) *Column {
	return b.UnsignedInt(column, true)
}

// BigIncrements bigint index
func (b *Blueprint) BigIncrements(column string) *Column {
	return b.UnsignedBigInt(column, true)
}

// Char add char column
func (b *Blueprint) Char(column string, length ...int) *Column {
	return b.AddColumn(ColumnTypeChar, column, Map{
		ColumnAttrLength: varDef(length, DefaultStringLength),
	})
}

// Varchar add varchar column
func (b *Blueprint) Varchar(column string, length ...int) *Column {
	return b.AddColumn(ColumnTypeVarchar, column, Map{
		ColumnAttrLength: varDef(length, DefaultStringLength),
	})
}

// TinyText add tinytext column
func (b *Blueprint) TinyText(column string) *Column {
	return b.AddColumn(ColumnTypeTinyText, column)
}

// Text add text column
func (b *Blueprint) Text(column string) *Column {
	return b.AddColumn(ColumnTypeText, column)
}

// MediumText add medium column
func (b *Blueprint) MediumText(column string) *Column {
	return b.AddColumn(ColumnTypeMediumText, column)
}

// LongText add longtext column
func (b *Blueprint) LongText(column string) *Column {
	return b.AddColumn(ColumnTypeLongText, column)
}

// String is Varchar alias
func (b *Blueprint) String(column string, length ...int) *Column {
	return b.Varchar(column, length...)
}

// Int add int column
func (b *Blueprint) Int(column string, autoIncrementAndUnsigned ...bool) *Column {
	return b.AddColumn(ColumnTypeInt, column, b.getAutoIncrementAndUnsignedMap(autoIncrementAndUnsigned...))
}

// TinyInt add tinyint column
func (b *Blueprint) TinyInt(column string, autoIncrementAndUnsigned ...bool) *Column {
	return b.AddColumn(ColumnTypeTinyInt, column, b.getAutoIncrementAndUnsignedMap(autoIncrementAndUnsigned...))
}

// SmallInt add smallint column
func (b *Blueprint) SmallInt(column string, autoIncrementAndUnsigned ...bool) *Column {
	return b.AddColumn(ColumnTypeSmallInt, column, b.getAutoIncrementAndUnsignedMap(autoIncrementAndUnsigned...))
}

// MediumInt add mediumint column
func (b *Blueprint) MediumInt(column string, autoIncrementAndUnsigned ...bool) *Column {
	return b.AddColumn(ColumnTypeMediumInt, column, b.getAutoIncrementAndUnsignedMap(autoIncrementAndUnsigned...))
}

// BigInt add bigint column
func (b *Blueprint) BigInt(column string, autoIncrementAndUnsigned ...bool) *Column {
	return b.AddColumn(ColumnTypeBigInt, column, b.getAutoIncrementAndUnsignedMap(autoIncrementAndUnsigned...))
}

func (b *Blueprint) getAutoIncrementAndUnsignedMap(autoIncrementAndUnsigned ...bool) Map {
	var (
		autoIncrement bool
		unsigned      bool
	)

	if len(autoIncrementAndUnsigned) == 2 {
		autoIncrement = autoIncrementAndUnsigned[0]
		unsigned = autoIncrementAndUnsigned[1]
	} else if len(autoIncrementAndUnsigned) == 1 {
		autoIncrement = autoIncrementAndUnsigned[0]
	}

	return Map{
		ColumnAttrAutoIncrement: autoIncrement,
		ColumnAttrUnsigned:      unsigned,
	}
}

// Boolean add tinyint(1) column
func (b *Blueprint) Boolean(column string) *Column {
	return b.AddColumn(ColumnTypeBoolean, column)
}

// UnsignedInt add int column, autoIncrement default false
func (b *Blueprint) UnsignedInt(column string, autoIncrement ...bool) *Column {
	return b.Int(column, varDef(autoIncrement, false), true)
}

func (b *Blueprint) UnsignedTinyInt(column string, autoIncrement ...bool) *Column {
	return b.TinyInt(column, varDef(autoIncrement, false), true)
}

func (b *Blueprint) UnsignedSmallInt(column string, autoIncrement ...bool) *Column {
	return b.SmallInt(column, varDef(autoIncrement, false), true)
}

func (b *Blueprint) UnsignedMediumInt(column string, autoIncrement ...bool) *Column {
	return b.MediumInt(column, varDef(autoIncrement, false), true)
}

func (b *Blueprint) UnsignedBigInt(column string, autoIncrement ...bool) *Column {
	return b.BigInt(column, varDef(autoIncrement, false), true)
}

// Float float, unsigned default false
func (b *Blueprint) Float(column string, total, places int, unsigned ...bool) *Column {
	return b.AddColumn(ColumnTypeFloat, column, Map{
		ColumnAttrTotal:    total,
		ColumnAttrPlaces:   places,
		ColumnAttrUnsigned: varDef(unsigned, false),
	})
}

func (b *Blueprint) Double(column string, total, places int, unsigned ...bool) *Column {
	return b.AddColumn(ColumnTypeDouble, column, Map{
		ColumnAttrTotal:    total,
		ColumnAttrPlaces:   places,
		ColumnAttrUnsigned: varDef(unsigned, false),
	})
}

func (b *Blueprint) Decimal(column string, total, places int, unsigned ...bool) *Column {
	return b.AddColumn(ColumnTypeDecimal, column, Map{
		ColumnAttrTotal:    total,
		ColumnAttrPlaces:   places,
		ColumnAttrUnsigned: varDef(unsigned, false),
	})
}

// UnsignedFloat unsigned float
func (b *Blueprint) UnsignedFloat(column string, total, places int) *Column {
	return b.Float(column, total, places, true)
}

func (b *Blueprint) UnsignedDouble(column string, total, places int) *Column {
	return b.Double(column, total, places, true)
}

func (b *Blueprint) UnsignedDecimal(column string, total, places int) *Column {
	return b.Decimal(column, total, places, true)
}

// Enum enum
func (b *Blueprint) Enum(column string, allowed []string) *Column {
	return b.AddColumn(ColumnTypeEnum, column, Map{
		ColumnAttrAllowed: allowed,
	})
}

func (b *Blueprint) Set(column string, allowed []string) *Column {
	return b.AddColumn(ColumnTypeSet, column, Map{
		ColumnAttrAllowed: allowed,
	})
}

// Json json column
func (b *Blueprint) Json(column string) *Column {
	return b.AddColumn(ColumnTypeJson, column)
}

func (b *Blueprint) Date(column string) *Column {
	return b.AddColumn(ColumnTypeDate, column)
}

func (b *Blueprint) DateTime(column string) *Column {
	return b.AddColumn(ColumnTypeDateTime, column)
}

func (b *Blueprint) Time(column string) *Column {
	return b.AddColumn(ColumnTypeTime, column)
}

func (b *Blueprint) Timestamp(column string) *Column {
	return b.AddColumn(ColumnTypeTimestamp, column)
}

func (b *Blueprint) Timestamps() {
	b.Timestamp("created_at").Nullable()
	b.Timestamp("updated_at").Nullable()
}

func (b *Blueprint) SoftDeletes() *Column {
	return b.Timestamp("deleted_at").Nullable()
}

func (b *Blueprint) Year(column string) *Column {
	return b.AddColumn(ColumnTypeYear, column)
}

func (b *Blueprint) Binary(column string) *Column {
	return b.AddColumn(ColumnTypeBinary, column)
}

func (b *Blueprint) Blob(column string) *Column {
	return b.AddColumn(ColumnTypeBinary, column)
}

func (b *Blueprint) Uuid(column string) *Column {
	return b.AddColumn(ColumnTypeUuid, column)
}

// Comment add comment
func (b *Blueprint) Comment(comment string) {
	b.addCommand(commandComment, Map{
		commandAttrComment: comment,
	})
}

// GetTable get table name
func (b *Blueprint) GetTable() string {
	return b.table
}

// GetConfig get config
func (b *Blueprint) GetConfig() *Config {
	return b.config
}

// GetColumns get columns
func (b *Blueprint) GetColumns() []*Column {
	return b.columns
}

// GetCommands get commands
func (b *Blueprint) GetCommands() []*Command {
	return b.commands
}

// getAddedColumns get added columns
func (b *Blueprint) getAddedColumns() (columns []*Column) {
	return filter(b.columns, func(v *Column) bool {
		return v.Attributes[ColumnAttrChange] != true
	})
}

// getChangedColumns get need change columns
func (b *Blueprint) getChangedColumns() []*Column {
	return filter(b.columns, func(v *Column) bool {
		return v.Attributes[ColumnAttrChange] == true
	})
}

// build exec sql
func (b *Blueprint) build(grammar *Grammar) (err error) {
	if b.config.DB == nil {
		return errors.New("DB is nil")
	}

	statements := b.ToSql(grammar)

	for _, statement := range statements {
		_, err = b.config.DB.ExecContext(b.ctx, statement)
		if err != nil {
			return err
		}
	}

	return nil
}

// ToSql Get the raw SQL statements for the blueprint.
func (b *Blueprint) ToSql(grammar *Grammar) []string {
	b.addImpliedCommands()

	var statements []string

	for _, command := range b.commands {
		statements = append(statements, grammar.Compile(b, command))
	}

	statements = unique(statements)

	return statements
}

// Create Indicate that the table needs to be created.
func (b *Blueprint) create() *Command {
	return b.addCommand(commandCreate)
}

// Drop Indicate that the table should be dropped.
func (b *Blueprint) drop() *Command {
	return b.addCommand(commandDrop)
}

// DropIfExists Indicate that the table should be dropped if it exists.
func (b *Blueprint) dropIfExists() *Command {
	return b.addCommand(commandDropIfExists)
}

// addImpliedCommands Add the commands that are implied by the blueprint's state.
func (b *Blueprint) addImpliedCommands() {
	if !b.creating() && len(b.getAddedColumns()) > 0 {
		// b.commands = append([]*Command{b.createCommand(commandAdd)}, b.commands...)
		b.addCommand(commandAdd)
	}

	if !b.creating() && len(b.getChangedColumns()) > 0 {
		b.addCommand(commandChange)
	}

	b.addFluentIndexes()
}

// addFluentIndexes Add the index commands fluently specified on columns.
func (b *Blueprint) addFluentIndexes() {
	indexes := []string{commandPrimary, commandUnique, commandIndex}

	for _, column := range b.columns {
		continue2 := false
		for _, index := range indexes {
			if continue2 {
				continue
			}
			if _, ok := column.Attributes[index]; ok {
				switch index {
				case commandPrimary:
					b.Primary(column.Name)
				case commandUnique:
					b.Unique(column.Name)
				case commandIndex:
					b.Index(column.Name)
				}
				column.Attributes[index] = false
				continue2 = true
			}
		}
	}
}

// creating check has create command
func (b *Blueprint) creating() bool {
	for _, command := range b.commands {
		if command.Name == commandCreate {
			return true
		}
	}
	return false
}

// addCommand add new command
func (b *Blueprint) addCommand(name string, parameters ...Map) (command *Command) {
	command = b.createCommand(name, parameters...)
	b.commands = append(b.commands, command)
	return
}

// createCommand create command
func (b *Blueprint) createCommand(name string, parameters ...Map) (command *Command) {
	params := varDef(parameters, Map{})
	command = &Command{
		Name:       name,
		Attributes: params,
	}
	return
}

// indexCommand add index command
func (b *Blueprint) indexCommand(t string, columns interface{}, algorithm ...string) {
	if column, ok := columns.(string); ok {
		b.addCommand(t, Map{
			commandAttrIndex:     b.createIndexName(t, []string{column}),
			commandAttrAlgorithm: varDef(algorithm, ""),
			commandAttrColumns:   []string{column},
		})
	}

	if column, ok := columns.([]string); ok {
		b.addCommand(t, Map{
			commandAttrIndex:     b.createIndexName(t, column),
			commandAttrAlgorithm: varDef(algorithm, ""),
			commandAttrColumns:   column,
		})
	}
}

// createIndexName create index name
func (b *Blueprint) createIndexName(t string, columns []string) string {
	index := strings.ToLower(b.config.Prefix + b.table + "_" + strings.Join(columns, "_") + "_" + t)
	return replaceByArray(index, []string{"-", "_", ".", "_"})
}
