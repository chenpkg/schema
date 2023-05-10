package schema

import (
	"fmt"
	"strconv"
	"strings"
)

var localGrammar = &Grammar{}

type Grammar struct{}

// GetColumns get add columns
func (g *Grammar) GetColumns(blueprint *Blueprint) []string {
	var columns []string

	for _, column := range blueprint.getAddedColumns() {
		sql := "`" + column.Name + "` " + g.GetType(column)
		columns = append(columns, g.addModifiers(sql, blueprint, column))
	}

	return columns
}

// GetChangeColumns get change columns
func (g *Grammar) GetChangeColumns(blueprint *Blueprint) (columns []string) {
	for _, column := range blueprint.getChangedColumns() {
		sql := "`" + column.Name + "` " + g.GetType(column)
		columns = append(columns, g.addModifiers(sql, blueprint, column))
	}
	return
}

// GetType get column type
func (g *Grammar) GetType(column *Column) string {
	switch column.Type {
	case ColumnTypeChar:
		return column.Type + "(" + strconv.Itoa(column.Attributes[ColumnAttrLength].(int)) + ")"

	case ColumnTypeVarchar:
		return column.Type + "(" + strconv.Itoa(column.Attributes[ColumnAttrLength].(int)) + ")"

	case ColumnTypeTinyText, ColumnTypeText, ColumnTypeMediumText, ColumnTypeLongText, ColumnTypeBigInt, ColumnTypeInt, ColumnTypeMediumInt, ColumnTypeTinyInt, ColumnTypeSmallInt, ColumnTypeJson, ColumnTypeDate, ColumnTypeDateTime, ColumnTypeTime, ColumnTypeTimestamp, ColumnTypeYear, ColumnTypeBinary, ColumnTypeBlob:
		return column.Type

	case ColumnTypeFloat, ColumnTypeDouble, ColumnTypeDecimal:
		var (
			t      = column.Type
			total  = strconv.Itoa(column.Attributes[ColumnAttrTotal].(int))
			places = strconv.Itoa(column.Attributes[ColumnAttrPlaces].(int))
		)
		if t == ColumnTypeFloat {
			t = ColumnTypeDouble
		}
		return t + "(" + total + ", " + places + ")"

	case ColumnTypeEnum:
		return fmt.Sprintf("enum(%s)", g.quoteString(column.Attributes[ColumnAttrAllowed].([]string)))

	case ColumnTypeSet:
		return fmt.Sprintf("set(%s)", g.quoteString(column.Attributes[ColumnAttrAllowed].([]string)))

	case ColumnTypeUuid:
		return ColumnTypeChar + "(36)"

	case ColumnTypeBoolean:
		return ColumnTypeTinyInt + "(1)"
	}

	return ""
}

// quoteString Quote the given string literal.
func (g *Grammar) quoteString(value []string) string {
	value = arrMap(value, func(v string) string {
		return "'" + v + "'"
	})

	return strings.Join(value, ", ")
}

// addModifiers Add the column modifiers to the definition.
func (g *Grammar) addModifiers(sql string, blueprint *Blueprint, column *Column) string {
	// modifiers := []string{
	// 	"Unsigned", "Charset", "Collate", "Nullable", "Default", "Increment", "Comment",
	// }

	// Unsigned
	if column.Attributes[ColumnAttrUnsigned] == true {
		sql += " " + ColumnAttrUnsigned
	}

	// Charset
	if charset, ok := column.Attributes[ColumnAttrCharset]; ok {
		sql += " character set " + charset.(string)
	}

	// Collate
	if collate, ok := column.Attributes[ColumnAttrCollate]; ok {
		sql += " collate '" + collate.(string) + "'"
	}

	// Nullable
	if nullable, ok := column.Attributes[ColumnAttrNullable]; ok && nullable.(bool) == true {
		sql += " null"
	} else {
		sql += " not null"
	}

	// Default
	if def, ok := column.Attributes[ColumnAttrDefault]; ok {
		if def == "" {
			sql += " default ''"
		} else {
			sql += " default '" + convString(def) + "'"
		}
	}

	// Increment
	serials := []string{
		ColumnTypeBigInt, ColumnTypeInt, ColumnTypeMediumInt, ColumnTypeSmallInt, ColumnTypeTinyInt,
	}
	if inArray(column.Type, serials) && column.Attributes[ColumnAttrAutoIncrement] == true {
		sql += " auto_increment primary key"
	}

	// Comment
	if comment, ok := column.Attributes[ColumnAttrComment]; ok && comment.(string) != "" {
		sql += " comment '" + comment.(string) + "'"
	}

	return sql
}

func (g *Grammar) wrap(value string) string {
	return "`" + value + "`"
}

func (g *Grammar) wrapTable(blueprint *Blueprint) string {
	return g.wrap(blueprint.Prefix + blueprint.GetTable())
}

// prefixStrings
func (g *Grammar) prefixStrings(prefix string, values []string) []string {
	return arrMap(values, func(v string) string {
		return prefix + v
	})
}
