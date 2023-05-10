package schema

type Column struct {
	Type       string
	Name       string
	Attributes Map
}

// Comment add column comment
func (c *Column) Comment(comment string) *Column {
	c.Attributes[ColumnAttrComment] = comment
	return c
}

// Default add column default value
func (c *Column) Default(value interface{}) *Column {
	c.Attributes[ColumnAttrDefault] = value
	return c
}

// Nullable Can it be empty, default to true
func (c *Column) Nullable(value ...bool) *Column {
	c.Attributes[ColumnAttrNullable] = varDef(value, true)
	return c
}

// Primary add primary
func (c *Column) Primary() *Column {
	c.Attributes[ColumnAttrPrimary] = true
	return c
}

// Unique add unique index
func (c *Column) Unique() *Column {
	c.Attributes[ColumnAttrUnique] = true
	return c
}

// Index add index
func (c *Column) Index() *Column {
	c.Attributes[ColumnAttrIndex] = true
	return c
}

// Change column
func (c *Column) Change() *Column {
	c.Attributes[ColumnAttrChange] = true
	return c
}

// Charset set column charset
func (c *Column) Charset(charset string) *Column {
	c.Attributes[ColumnAttrCharset] = charset
	return c
}

// Collation set column collation
func (c *Column) Collation(collation string) *Column {
	c.Attributes[ColumnAttrCollate] = collation
	return c
}
