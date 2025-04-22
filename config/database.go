package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// DatabaseConfig represents the configuration for database connection
type DatabaseConfig struct {
	Type     string `json:"type"`     // mysql, postgres, sqlite, sqlserver, clickhouse
	Host     string `json:"host"`     // database host
	Port     int    `json:"port"`     // database port
	Username string `json:"username"` // database username
	Password string `json:"password"` // database password
	Database string `json:"database"` // database name
	SSLMode  string `json:"ssl_mode"` // for postgres
	File     string `json:"file"`     // for sqlite
	Mode     string `json:"mode"`     // server mode (stdio or http)
	Addr     string `json:"addr"`     // server mode (stdio or http)
}

// LoadConfig loads configuration from file and command line
func LoadConfig() (*DatabaseConfig, error) {
	// Define command line flags
	configFile := flag.String("config", "config.yaml", "Path to config file")
	dbType := flag.String("db-type", "", "Database type (mysql, postgres, sqlite, sqlserver, clickhouse)")
	dbHost := flag.String("db-host", "", "Database host")
	dbPort := flag.Int("db-port", 0, "Database port")
	dbUser := flag.String("db-user", "", "Database username")
	dbPass := flag.String("db-pass", "", "Database password")
	dbName := flag.String("db-name", "", "Database name")
	dbSSLMode := flag.String("db-ssl-mode", "", "Database SSL mode (for postgres)")
	dbFile := flag.String("db-file", "", "Database file (for sqlite)")
	mode := flag.String("mode", "stdio", "Server mode (stdio or http)")
	addr := flag.String("addr", ":8080", "http server listen address")

	// Check if any arguments are provided
	if len(os.Args) == 1 {
		fmt.Println("Command Line Arguments:")
		fmt.Println("  --config string     Path to config file (default \"config.yaml\")")
		fmt.Println("  --db-type string    Database type (mysql, postgres, sqlite, sqlserver, clickhouse)")
		fmt.Println("  --db-host string    Database host")
		fmt.Println("  --db-port int       Database port")
		fmt.Println("  --db-user string    Database username")
		fmt.Println("  --db-pass string    Database password")
		fmt.Println("  --db-name string    Database name")
		fmt.Println("  --db-ssl-mode string Database SSL mode (for postgres)")
		fmt.Println("  --db-file string    Database file (for sqlite)")
		fmt.Println("  --mode string       Server mode (stdio or http) (default \"stdio\")")
		fmt.Println("  --addr string       HTTP server listen address (default \":8080\")")
		fmt.Println("\nUsage examples:")
		fmt.Println("  Basic usage with config file:")
		fmt.Println("    ./database-mcp -config config.yaml")
		fmt.Println("  MySQL example:")
		fmt.Println("    ./database-mcp -db-type mysql -db-host localhost -db-port 3306 -db-user root -db-pass password -db-name mydb")
		fmt.Println("  PostgreSQL example:")
		fmt.Println("    ./database-mcp -db-type postgres -db-host localhost -db-port 5432 -db-user postgres -db-pass password -db-name mydb -db-ssl-mode disable")
		fmt.Println("  SQLite example:")
		fmt.Println("    ./database-mcp -db-type sqlite -db-file database.db")
		fmt.Println("  Server mode example:")
		fmt.Println("    ./database-mcp -mode http -addr :8080")
		os.Exit(0)
	}

	flag.Parse()

	// Initialize Viper
	v := viper.New()
	v.SetConfigFile(*configFile)
	v.SetConfigType("yaml")

	// Set default values
	v.SetDefault("database.type", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.username", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.database", "mydb")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.file", "database.db")
	v.SetDefault("database.mode", "stdio")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	// Override with command line flags if provided
	if *dbType != "" {
		v.Set("database.type", *dbType)
	}
	if *dbHost != "" {
		v.Set("database.host", *dbHost)
	}
	if *dbPort != 0 {
		v.Set("database.port", *dbPort)
	}
	if *dbUser != "" {
		v.Set("database.username", *dbUser)
	}
	if *dbPass != "" {
		v.Set("database.password", *dbPass)
	}
	if *dbName != "" {
		v.Set("database.database", *dbName)
	}
	if *dbSSLMode != "" {
		v.Set("database.ssl_mode", *dbSSLMode)
	}
	if *dbFile != "" {
		v.Set("database.file", *dbFile)
	}
	if *mode != "" {
		v.Set("database.mode", *mode)
	}

	if *addr != "" {
		v.Set("database.addr", *addr)
	}
	// Create config struct
	config := &DatabaseConfig{
		Type:     v.GetString("database.type"),
		Host:     v.GetString("database.host"),
		Port:     v.GetInt("database.port"),
		Username: v.GetString("database.username"),
		Password: v.GetString("database.password"),
		Database: v.GetString("database.database"),
		SSLMode:  v.GetString("database.ssl_mode"),
		File:     v.GetString("database.file"),
		Mode:     v.GetString("database.mode"),
		Addr:     v.GetString("database.addr"),
	}

	return config, nil
}

// NewDatabase creates a new database connection based on the configuration
func NewDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(config.Type) {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			config.Host, config.Username, config.Password, config.Database, config.Port, config.SSLMode)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(config.File)
	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			config.Username, config.Password, config.Host, config.Port, config.Database)
		dialector = sqlserver.Open(dsn)
	case "clickhouse":
		dsn := fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s",
			config.Host, config.Port, config.Database, config.Username, config.Password)
		dialector = clickhouse.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
