# 数据库 MCP 服务

一个基于 MCP (Metoro Control Protocol) 的数据库服务，通过 GORM 支持多种数据库类型。

## 功能特点

- 支持多种数据库类型：
  - MySQL
  - PostgreSQL
  - SQLite
  - SQL Server
  - ClickHouse
- 多种配置方式：
  - 配置文件 (YAML)
  - 命令行参数
  - 环境变量
- MCP 协议集成
- GORM ORM 支持

## 安装

1. 克隆仓库
2. 安装依赖：
   ```bash
   go mod tidy
   ```

## 配置

### 配置文件 (config.yaml)

创建 `config.yaml` 文件，结构如下：

```yaml
database:
  type: "mysql"  # mysql, postgres, sqlite, sqlserver, clickhouse
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "mydb"
  ssl_mode: "disable"  # 用于 postgres
  file: "database.db"  # 用于 sqlite
```

### 命令行参数

可以通过命令行参数覆盖配置文件中的设置：

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

可用的命令行参数：
- `--config`: 配置文件路径（默认："config.yaml"）
- `--db-type`: 数据库类型（mysql, postgres, sqlite, sqlserver, clickhouse）
- `--db-host`: 数据库主机
- `--db-port`: 数据库端口
- `--db-user`: 数据库用户名
- `--db-pass`: 数据库密码
- `--db-name`: 数据库名称
- `--db-ssl-mode`: SSL 模式（用于 PostgreSQL）
- `--db-file`: 数据库文件（用于 SQLite）

## 使用方法

1. 启动服务：
   ```bash
   ./database-mcp --db-type=mysql --db-host=127.0.0.1 --db-port=3306 --db-user=dbuser --db-pass=passwd --db-name=dbname --db-ssl-mode=disable --mode=http
   ```

2. 服务将执行以下操作：
   - 从配置文件和/或命令行加载配置
   - 初始化数据库连接
   - 启动 MCP 服务器
   - 注册可用的工具和资源

## MCP 配置

要在 MCP 客户端中使用本服务，请配置以下 MCP 端点：

```
http://127.0.0.1:8080/sse
```

这将使您的 MCP 客户端连接到本地运行在 8080 端口的服务。

## MCP 工具

服务提供以下 MCP 工具：

1. `get_tables`: 获取数据库中的所有表
   - 返回包含表名和注释的表列表
2. `get_table_detail`: 获取特定表的详细信息
   - 参数：
     - `table_name`: 要获取详细信息的表名
   - 返回表信息，包括：
     - 表名和注释
     - 列信息（名称、类型、注释、是否可为空、默认值）
3. `execute_sql`: 执行 SQL 查询
   - 参数：
     - `query`: 要执行的 SQL 查询
   - 返回：
     - 对于 SELECT 查询：结果数组
     - 对于其他查询：受影响的行数

## 许可证

MIT 许可证 