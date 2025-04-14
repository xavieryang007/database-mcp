# Database MCP Service

A MCP (Metoro Control Protocol) service with database capabilities, supporting multiple database types through GORM.

## Features

- Support for multiple database types:
  - MySQL
  - PostgreSQL
  - SQLite
  - SQL Server
  - ClickHouse
- Configuration through:
  - Configuration file (YAML)
  - Command line arguments
  - Environment variables
- MCP protocol integration
- GORM ORM support

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

### Configuration File (config.yaml)

Create a `config.yaml` file with the following structure:

```yaml
database:
  type: "mysql"  # mysql, postgres, sqlite, sqlserver, clickhouse
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "mydb"
  ssl_mode: "disable"  # for postgres
  file: "database.db"  # for sqlite
```

### Command Line Arguments

You can override configuration file settings using command line arguments:

```bash
./database-mcp --config=config.yaml \
  --db-type=mysql \
  --db-host=localhost \
  --db-port=3306 \
  --db-user=root \
  --db-pass=password \
  --db-name=mydb \
  --db-ssl-mode=disable \
  --db-file=database.db
```

Available command line arguments:
- `--config`: Path to config file (default: "config.yaml")
- `--db-type`: Database type (mysql, postgres, sqlite, sqlserver, clickhouse)
- `--db-host`: Database host
- `--db-port`: Database port
- `--db-user`: Database username
- `--db-pass`: Database password
- `--db-name`: Database name
- `--db-ssl-mode`: SSL mode (for PostgreSQL)
- `--db-file`: Database file (for SQLite)

## Usage

1. Start the service:
   ```bash
   ./database-mcp
   ```

2. The service will:
   - Load configuration from file and/or command line
   - Initialize database connection
   - Start MCP server
   - Register available tools and resources

## MCP Tools

The service provides the following MCP tools:

1. `hello`: A simple greeting tool
2. `get_tables`: Get all tables in the database
   - Returns a list of tables with their names and comments
3. `get_table_detail`: Get detailed information about a specific table
   - Arguments:
     - `table_name`: The name of the table to get details for
   - Returns table information including:
     - Table name and comment
     - Column information (name, type, comment, nullable, default value)
4. `prompt_test`: A test prompt handler
5. `resource_test`: A test resource handler

## License

MIT License 