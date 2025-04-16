package main

import (
	"context"
	"database-mcp/config"
	"database-mcp/tools"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gorm.io/gorm"
	"log"
)

// DatabaseMCP represents our MCP service with database capabilities
type DatabaseMCP struct {
	db       *gorm.DB
	server   *server.MCPServer
	dbConfig *config.DatabaseConfig
}

// Content represents the content structure for our tools
type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}

// MyFunctionsArguments represents the arguments for our tools
type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

// TableDetailArgs represents arguments for table detail tool
type TableDetailArgs struct {
	TableName string `json:"table_name" jsonschema:"required,description=The name of the table to get details for"`
}

// SQLQueryArgs represents arguments for executing SQL queries
type SQLQueryArgs struct {
	Query string `json:"query" jsonschema:"required,description=Execute SQL statements"`
}

type DatabasesArgs struct {
	Databases string `json:"databases" jsonschema:"required,description=Name of the database to be operated on"`
}

func NewDatabaseMCP(dbConfig *config.DatabaseConfig) (*DatabaseMCP, error) {
	// Initialize database
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		return nil, err
	}

	// Initialize MCP server
	s := server.NewMCPServer(
		"database-mcp",
		"0.0.1",
	)

	return &DatabaseMCP{
		db:       db,
		server:   s,
		dbConfig: dbConfig,
	}, nil
}

func (m *DatabaseMCP) registerTools() error {
	// Initialize database tool
	dbTool := tools.NewDatabaseTool(m.db)

	// Register get tables tool
	getTablesTool := mcp.NewTool("get_tables",
		mcp.WithDescription("Get all tables in the database"),
	)

	m.server.AddTool(getTablesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tables, err := dbTool.GetTables()
		if err != nil {
			return nil, fmt.Errorf("failed to get tables: %v", err)
		}
		jsonData, err := json.Marshal(tables)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tables: %v", err)
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	})

	// Register get table detail tool
	getTableDetailTool := mcp.NewTool("get_table_detail",
		mcp.WithDescription("Get detailed information about a specific table"),
		mcp.WithString("table_name",
			mcp.Required(),
			mcp.Description("The name of the table to get details for"),
		),
	)

	m.server.AddTool(getTableDetailTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tableName := request.Params.Arguments["table_name"].(string)
		detail, err := dbTool.GetTableDetail(tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get table detail: %v", err)
		}
		jsonData, err := json.Marshal(detail)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal detail: %v", err)
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	})

	// Register execute_sql tool
	executeSQLTool := mcp.NewTool("execute_sql",
		mcp.WithDescription("Execute a SQL query"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	m.server.AddTool(executeSQLTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.Params.Arguments["query"].(string)
		result, err := dbTool.ExecuteSQL(query)
		if err != nil {
			return nil, fmt.Errorf("failed to execute SQL: %v", err)
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %v", err)
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	})

	return nil
}

func (m *DatabaseMCP) Start() error {
	// Register all tools
	if err := m.registerTools(); err != nil {
		return fmt.Errorf("failed to register tools: %v", err)
	}

	if m.dbConfig.Mode == "http" {
		if err := server.NewSSEServer(m.server).Start(m.dbConfig.Addr); err != nil {
			return fmt.Errorf("failed to start server: %v", err)
		}
	} else {
		// Default to stdio mode
		if err := server.ServeStdio(m.server); err != nil {
			return fmt.Errorf("failed to start stdio server: %v", err)
		}
	}
	return nil
}

func main() {
	// Load configuration
	dbConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	_mcp, err := NewDatabaseMCP(dbConfig)
	if err != nil {
		log.Fatalf("Failed to create DatabaseMCP: %v", err)
	}

	if err := _mcp.Start(); err != nil {
		log.Fatalf("Failed to start DatabaseMCP: %v", err)
	}
}
