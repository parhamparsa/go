package sqlt

type SqlTable struct {
	name         string
	primaryField string
	columns      []string
}

func (tbl *SqlTable) PrimaryField() string {
	return tbl.primaryField
}

func NewTable(name string, primaryField string, columns []string) *SqlTable {
	return &SqlTable{name: name, columns: columns, primaryField: primaryField}
}
func (tbl *SqlTable) Name() string {
	return tbl.name
}
func (tbl *SqlTable) Columns() []string {
	return tbl.columns
}
