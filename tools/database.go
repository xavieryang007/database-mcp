package tools

import (
	"database-mcp/models"
	"fmt"

	"gorm.io/gorm"
)

// DatabaseTool provides database-related operations
type DatabaseTool struct {
	db *gorm.DB
}

// NewDatabaseTool creates a new DatabaseTool instance
func NewDatabaseTool(db *gorm.DB) *DatabaseTool {
	return &DatabaseTool{db: db}
}

// GetTables retrieves all tables and their comments from the database
func (t *DatabaseTool) GetTables() ([]models.TableInfo, error) {
	var tables []models.TableInfo
	var err error

	// Get database type
	dbType := t.db.Dialector.Name()

	switch dbType {
	case "mysql":
		err = t.db.Raw(`
			SELECT table_name as table_name, 
				   table_comment as table_comment
			FROM information_schema.tables 
			WHERE table_schema = DATABASE()
		`).Scan(&tables).Error
	case "postgres":
		err = t.db.Raw(`
			SELECT table_name as table_name,
				   obj_description(relfilenode, 'pg_class') as table_comment
			FROM information_schema.tables t
			JOIN pg_class c ON c.relname = t.table_name
			WHERE table_schema = 'public'
		`).Scan(&tables).Error
	case "sqlite":
		err = t.db.Raw(`
			SELECT name as table_name,
				   '' as table_comment
			FROM sqlite_master
			WHERE type = 'table'
		`).Scan(&tables).Error
	case "sqlserver":
		err = t.db.Raw(`
			SELECT t.name as table_name,
				   ep.value as table_comment
			FROM sys.tables t
			LEFT JOIN sys.extended_properties ep ON ep.major_id = t.object_id
			WHERE ep.name = 'MS_Description'
		`).Scan(&tables).Error
	case "clickhouse":
		err = t.db.Raw(`
			SELECT name as table_name,
				   comment as table_comment
			FROM system.tables
			WHERE database = currentDatabase()
		`).Scan(&tables).Error
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %v", err)
	}

	return tables, nil
}

// GetTableDetail retrieves detailed information about a specific table
func (t *DatabaseTool) GetTableDetail(tableName string) (*models.TableDetail, error) {
	var tableDetail models.TableDetail
	var err error

	// Get database type
	dbType := t.db.Dialector.Name()

	// First get table info
	switch dbType {
	case "mysql":
		err = t.db.Raw(`
			SELECT table_name as table_name, 
				   table_comment as table_comment
			FROM information_schema.tables 
			WHERE table_schema = DATABASE()
			AND table_name = ?
		`, tableName).Scan(&tableDetail.TableInfo).Error
	case "postgres":
		err = t.db.Raw(`
			SELECT table_name as table_name,
				   obj_description(relfilenode, 'pg_class') as table_comment
			FROM information_schema.tables t
			JOIN pg_class c ON c.relname = t.table_name
			WHERE table_schema = 'public'
			AND table_name = ?
		`, tableName).Scan(&tableDetail.TableInfo).Error
	case "sqlite":
		err = t.db.Raw(`
			SELECT name as table_name,
				   '' as table_comment
			FROM sqlite_master
			WHERE type = 'table'
			AND name = ?
		`, tableName).Scan(&tableDetail.TableInfo).Error
	case "sqlserver":
		err = t.db.Raw(`
			SELECT t.name as table_name,
				   ep.value as table_comment
			FROM sys.tables t
			LEFT JOIN sys.extended_properties ep ON ep.major_id = t.object_id
			WHERE ep.name = 'MS_Description'
			AND t.name = ?
		`, tableName).Scan(&tableDetail.TableInfo).Error
	case "clickhouse":
		err = t.db.Raw(`
			SELECT name as table_name,
				   comment as table_comment
			FROM system.tables
			WHERE database = currentDatabase()
			AND name = ?
		`, tableName).Scan(&tableDetail.TableInfo).Error
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get table info: %v", err)
	}

	// Then get column info
	var columns []models.ColumnInfo
	switch dbType {
	case "mysql":
		err = t.db.Raw(`
			SELECT column_name as column_name,
				   column_type as column_type,
				   column_comment as column_comment,
				   is_nullable as is_nullable,
				   column_default as column_default
			FROM information_schema.columns
			WHERE table_schema = DATABASE()
			AND table_name = ?
			ORDER BY ordinal_position
		`, tableName).Scan(&columns).Error
	case "postgres":
		err = t.db.Raw(`
			SELECT column_name as column_name,
				   data_type as column_type,
				   col_description((table_schema||'.'||table_name)::regclass::oid, ordinal_position) as column_comment,
				   is_nullable as is_nullable,
				   column_default as column_default
			FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = ?
			ORDER BY ordinal_position
		`, tableName).Scan(&columns).Error
	case "sqlite":
		err = t.db.Raw(`
			SELECT name as column_name,
				   type as column_type,
				   '' as column_comment,
				   CASE WHEN "notnull" = 0 THEN 'YES' ELSE 'NO' END as is_nullable,
				   dflt_value as column_default
			FROM pragma_table_info(?)
		`, tableName).Scan(&columns).Error
	case "sqlserver":
		err = t.db.Raw(`
			SELECT c.name as column_name,
				   t.name as column_type,
				   ep.value as column_comment,
				   CASE WHEN c.is_nullable = 1 THEN 'YES' ELSE 'NO' END as is_nullable,
				   OBJECT_DEFINITION(c.default_object_id) as column_default
			FROM sys.columns c
			JOIN sys.types t ON c.user_type_id = t.user_type_id
			LEFT JOIN sys.extended_properties ep ON ep.major_id = c.object_id AND ep.minor_id = c.column_id
			WHERE c.object_id = OBJECT_ID(?)
		`, tableName).Scan(&columns).Error
	case "clickhouse":
		err = t.db.Raw(`
			SELECT name as column_name,
				   type as column_type,
				   comment as column_comment,
				   CASE WHEN is_nullable = 1 THEN 'YES' ELSE 'NO' END as is_nullable,
				   default_expression as column_default
			FROM system.columns
			WHERE database = currentDatabase()
			AND table = ?
		`, tableName).Scan(&columns).Error
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get column info: %v", err)
	}

	tableDetail.Columns = columns
	return &tableDetail, nil
}
