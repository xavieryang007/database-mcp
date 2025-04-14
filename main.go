package main

import (
	"database-mcp/config"
	"database-mcp/tools"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"gorm.io/gorm"
)

// DatabaseMCP represents our MCP service with database capabilities
type DatabaseMCP struct {
	db     *gorm.DB
	server *mcp.Server
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
	Query string `json:"query"`
}

func NewDatabaseMCP() (*DatabaseMCP, error) {
	// Load configuration
	dbConfig, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		return nil, err
	}

	// Initialize MCP server
	server := mcp.NewServer(stdio.NewStdioServerTransport())

	return &DatabaseMCP{
		db:     db,
		server: server,
	}, nil
}

func (m *DatabaseMCP) registerTools() error {
	// Initialize database tool
	dbTool := tools.NewDatabaseTool(m.db)
	var err error
	// Register get tables tool
	err = m.server.RegisterTool("get_tables", "Get all tables in the database", func(arguments struct{}) (*mcp.ToolResponse, error) {
		tables, err := dbTool.GetTables()
		if err != nil {
			return nil, fmt.Errorf("failed to get tables: %v", err)
		}
		jsonData, err := json.Marshal(tables)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tables: %v", err)
		}
		return mcp.NewToolResponse(&mcp.Content{
			Type: "application/json",
			TextContent: &mcp.TextContent{
				Text: string(jsonData),
			},
		}), nil
	})
	if err != nil {
		return err
	}

	// Register get table detail tool
	err = m.server.RegisterTool("get_table_detail", "Get detailed information about a specific table", func(arguments TableDetailArgs) (*mcp.ToolResponse, error) {
		detail, err := dbTool.GetTableDetail(arguments.TableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get table detail: %v", err)
		}
		jsonData, err := json.Marshal(detail)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal detail: %v", err)
		}

		return mcp.NewToolResponse(&mcp.Content{
			Type: "application/json",
			TextContent: &mcp.TextContent{
				Text: string(jsonData),
			},
		}), nil
	})
	if err != nil {
		return err
	}

	// Register execute_sql tool
	err = m.server.RegisterTool("execute_sql", "Execute a SQL query", func(arguments SQLQueryArgs) (*mcp.ToolResponse, error) {
		result, err := dbTool.ExecuteSQL(arguments.Query)
		if err != nil {
			return nil, fmt.Errorf("failed to execute SQL: %v", err)
		}
		// Convert result to string for text content
		resultStr := fmt.Sprintf("%v", result)

		return mcp.NewToolResponse(&mcp.Content{
			Type: "application/json",
			TextContent: &mcp.TextContent{
				Text: resultStr,
			},
		}), nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *DatabaseMCP) Start() error {
	// Register all tools
	if err := m.registerTools(); err != nil {
		return fmt.Errorf("failed to register tools: %v", err)
	}

	// Start the server
	if err := m.server.Serve(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	return nil
}

func main() {
	mcp, err := NewDatabaseMCP()
	if err != nil {
		log.Fatalf("Failed to create DatabaseMCP: %v", err)
	}

	if err := mcp.Start(); err != nil {
		log.Fatalf("Failed to start DatabaseMCP: %v", err)
	}
}
