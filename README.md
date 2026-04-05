# github.com/taouniverse/tao-postgres

[![Go Report Card](https://goreportcard.com/badge/github.com/taouniverse/tao-postgres)](https://goreportcard.com/report/github.com/taouniverse/tao-postgres)
[![GoDoc](https://pkg.go.dev/badge/github.com/taouniverse/tao-postgres?status.svg)](https://pkg.go.dev/github.com/taouniverse/tao-postgres?tab=doc)

Tao Universe 组件单元（Unit），基于泛型工厂模式封装 **PostgreSQL** 数据库。

## 安装

```bash
go get github.com/taouniverse/tao-postgres
```

## 使用

### 导入

```go
import _ "github.com/taouniverse/tao-postgres"
```

### 配置

```yaml
# 单实例配置
postgres:
  host: localhost
  port: 5432
  user: tao
  password: 123456qwe
  db: test
  ssl: disable
  time_zone: Asia/Shanghai

# 多实例配置（如主从分离）
postgres:
  default_instance: primary
  primary:
    host: localhost
    port: 5432
    user: tao
    password: 123456qwe
    db: mydb
  standby:
    host: backup.example.com
    port: 5432
    user: readonly
    password: ro_pass
    db: mydb
```

### 配置字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `host` | string | `localhost` | PostgreSQL 服务器地址 |
| `port` | int | `5432` | PostgreSQL 端口 |
| `user` | string | `postgres` | 用户名 |
| `password` | string | - | 密码 |
| `db` | string | - | 数据库名 |
| `ssl` | string | `disable` | SSL 模式 (disable/require/verify-ca/verify-full) |
| `time_zone` | string | `UTC` | 时区 |
| `max_idle` | int | `10` | 空闲连接池大小 |
| `max_open` | int | `100` | 最大打开连接数 |
| `max_lifetime` | duration | `1h` | 连接最大生命周期 |

## 工厂模式 API

| API | 说明 |
|-----|------|
| `postgres.M` | 配置实例 `*Config` |
| `postgres.Factory` | `*tao.BaseFactory[*gorm.DB]` 工厂实例 |
| `postgres.DB()` | 获取默认数据库连接 `(*gorm.DB, error)` |
| `postgres.GetDB(name)` | 获取指定名称的连接 `(*gorm.DB, error)` |

## 使用示例

### 获取连接并执行操作

```go
package main

import (
    "log"
    
    "github.com/taouniverse/tao-postgres"
)

func main() {
    // 获取默认实例
    db, err := postgres.DB()
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取底层 sql.DB 进行 Ping 测试
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal(err)
    }
    
    err = sqlDB.Ping()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("PostgreSQL 连接成功")
}
```

### GORM 操作

```go
db, _ := postgres.DB()

// 自动迁移
db.AutoMigrate(&User{})

// 创建记录
db.Create(&User{Name: "tao", Age: 18})

// 查询记录
var user User
db.First(&user, "name = ?", "tao")

// 更新记录
db.Model(&user).Update("age", 20)

// 删除记录
db.Delete(&user)
```

### 多实例使用

```go
// 获取主库连接（读写）
primary, _ := postgres.GetDB("primary")

// 获取备库连接（只读）
standby, _ := postgres.GetDB("standby")

// 主库写入
primary.Create(&Order{Amount: 100})

// 备库查询
var orders []Order
standby.Find(&orders)
```

## 单元测试

### 快速测试（无需 Docker）

```bash
# 仅运行配置相关测试
go test -v -run "TestConfig" ./...
```

### 完整集成测试（需要 Docker）

```bash
# 启动 PostgreSQL 并运行单实例测试
make test

# 启动 PostgreSQL 并运行多实例测试
make test-multi

# 启动 PostgreSQL 并运行所有测试
make test-all

# 生成覆盖率报告
make coverage

# 停止 PostgreSQL 服务
make down
```

### 手动测试

```bash
# 1. 启动 PostgreSQL
docker-compose up -d

# 2. 运行单实例测试
go test -v ./...

# 3. 运行多实例测试
TAO_TEST_MULTI_INSTANCE=true go test -v ./...

# 4. 停止 PostgreSQL
docker-compose down
```

## 开发指南

| 文件 | 说明 |
|------|------|
| `config.go` | InstanceConfig 字段 + ValidSelf 默认值 |
| `postgres.go` | NewPostgres 构造器 + 工厂注册 |
