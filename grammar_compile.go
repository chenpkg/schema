package schema

import (
	"fmt"
	"reflect"
	"strings"
)

// Compile compile command
func (g *Grammar) Compile(blueprint *Blueprint, command *Command) string {
	var (
		value   = reflect.ValueOf(g)
		name    = "Compile" + ucFirst(command.Name)
		compile = value.MethodByName(name)
	)

	return compile.Call([]reflect.Value{
		reflect.ValueOf(blueprint),
		reflect.ValueOf(command),
	})[0].Interface().(string)
}

// CompileCreate Compile a create table command.
func (g *Grammar) CompileCreate(blueprint *Blueprint, command *Command) string {
	var sql string

	sql = g.CompileCreateTable(blueprint)
	sql = g.CompileCreateEncoding(sql, blueprint)
	sql = g.CompileCreateEngine(sql, blueprint)

	return sql
}

// CompileCreateTable Create the main create table clause.
func (g *Grammar) CompileCreateTable(blueprint *Blueprint) string {
	return trim(fmt.Sprintf(
		"%s table %s (%s)", "create",
		g.wrapTable(blueprint),
		strings.Join(g.GetColumns(blueprint), ", ")))
}

// CompileCreateEncoding Append the character set specifications to a command.
func (g *Grammar) CompileCreateEncoding(sql string, blueprint *Blueprint) string {
	if blueprint.Charset != "" {
		sql += " default character set " + blueprint.Charset
	} else {
		sql += " default character set " + DefaultCharset
	}

	if blueprint.Collation != "" {
		sql += " collate '" + blueprint.Collation + "'"
	} else {
		sql += " collate '" + DefaultCollation + "'"
	}

	return sql
}

// CompileCreateEngine Append the engine specifications to a command.
func (g *Grammar) CompileCreateEngine(sql string, blueprint *Blueprint) string {
	if blueprint.Engine != "" {
		sql += " engine = " + blueprint.Engine
	} else {
		sql += " engine = " + DefaultEngine
	}

	return sql
}

// CompileAdd Compile an add column command.
func (g *Grammar) CompileAdd(blueprint *Blueprint, command *Command) string {
	columns := g.prefixStrings("add ", g.GetColumns(blueprint))
	return "alter table " + g.wrapTable(blueprint) + " " + strings.Join(columns, ", ")
}

// CompileChange Compile a change column command into a series of SQL statements.
func (g *Grammar) CompileChange(blueprint *Blueprint, command *Command) string {
	columns := g.prefixStrings("modify ", g.GetChangeColumns(blueprint))
	return "alter table " + g.wrapTable(blueprint) + " " + strings.Join(columns, ", ")
}

// CompileRename Rename the table to given name
func (g *Grammar) CompileRename(blueprint *Blueprint, command *Command) string {
	form := g.wrapTable(blueprint)
	to := g.wrap(blueprint.Prefix + command.Attributes[commandAttrTo].(string))
	return fmt.Sprintf("rename table %s to %s", form, to)
}

// func (g *Grammar) CompileRenameColumn(blueprint *Blueprint, command *Command) string {
// }

// CompilePrimary Compile a primary key command.
func (g *Grammar) CompilePrimary(blueprint *Blueprint, command *Command) string {
	command.Name = ""
	return g.CompileKey(blueprint, command, "primary key")
}

// CompileUnique Compile a unique key command.
func (g *Grammar) CompileUnique(blueprint *Blueprint, command *Command) string {
	return g.CompileKey(blueprint, command, "unique")
}

// CompileIndex Compile a plain index key command
func (g *Grammar) CompileIndex(blueprint *Blueprint, command *Command) string {
	return g.CompileKey(blueprint, command, "index")
}

// CompileKey Compile an index creation command.
func (g *Grammar) CompileKey(blueprint *Blueprint, command *Command, types string) string {
	var algorithm string

	if algo, ok := command.Attributes[commandAttrAlgorithm]; ok && algo.(string) != "" {
		algorithm = " using " + algo.(string)
	}

	columns := arrMap(command.Attributes[commandAttrColumns].([]string), g.wrap)
	columnize := strings.Join(columns, ", ")

	return trim(fmt.Sprintf(
		"alter table %s add %s %s%s(%s)",
		g.wrapTable(blueprint),
		types,
		command.Attributes[commandIndex].(string),
		algorithm,
		columnize))
}

// CompileDrop Compile a drop table command.
func (g *Grammar) CompileDrop(blueprint *Blueprint, command *Command) string {
	return "drop table " + g.wrapTable(blueprint)
}

// CompileDropIfExists Compile a drop table (if exists) command.
func (g *Grammar) CompileDropIfExists(blueprint *Blueprint, command *Command) string {
	return "drop table if exists " + g.wrapTable(blueprint)
}

// CompileDropColumn Compile a drop column command.
func (g *Grammar) CompileDropColumn(blueprint *Blueprint, command *Command) string {

	if cols, ok := command.Attributes[commandAttrColumns]; ok {
		var columns []string

		for _, col := range cols.([]string) {
			columns = append(columns, "drop "+g.wrap(col))
		}

		return "alter table " + g.wrapTable(blueprint) + " " + strings.Join(columns, ", ")
	}

	return ""
}

// CompileDropPrimary Compile a drop primary key command.
func (g *Grammar) CompileDropPrimary(blueprint *Blueprint, command *Command) string {
	return "alter table " + g.wrapTable(blueprint) + " drop primary key"
}

// CompileDropUnique Compile a drop unique key command.
func (g *Grammar) CompileDropUnique(blueprint *Blueprint, command *Command) string {
	return "alter table " + g.wrapTable(blueprint) + " drop index" + command.Attributes[commandAttrIndex].(string)
}

// CompileDropIndex Compile a drop index command.
func (g *Grammar) CompileDropIndex(blueprint *Blueprint, command *Command) string {
	return "alter table " + g.wrapTable(blueprint) + " drop index" + command.Attributes[commandAttrIndex].(string)
}

// CompileTableComment Compile a table comment command.
func (g *Grammar) CompileTableComment(blueprint *Blueprint, command *Command) string {
	comment := command.Attributes[commandAttrComment].(string)

	return fmt.Sprintf(
		"alter table %s comment = %s",
		g.wrapTable(blueprint),
		"'"+strings.Replace(comment, "'", "''", -1)+"'",
	)
}

// CompileTableExists Compile the query to determine the list of tables
func (g *Grammar) CompileTableExists() string {
	return "select * from information_schema.tables where table_schema = ? and table_name = ? and table_type = 'BASE TABLE'"
}
