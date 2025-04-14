package models

// TableInfo represents information about a database table
type TableInfo struct {
	TableName    string `json:"table_name"`    // Table name
	TableComment string `json:"table_comment"` // Table comment
}

// ColumnInfo represents information about a table column
type ColumnInfo struct {
	ColumnName    string `json:"column_name"`    // Column name
	ColumnType    string `json:"column_type"`    // Column type
	ColumnComment string `json:"column_comment"` // Column comment
	IsNullable    string `json:"is_nullable"`    // Whether the column is nullable
	ColumnDefault string `json:"column_default"` // Column default value
}

// TableDetail represents detailed information about a table including its columns
type TableDetail struct {
	TableInfo
	Columns []ColumnInfo `json:"columns"` // Table columns
}
