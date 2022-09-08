# 凭证管理

通用的凭证管理组件，支持将凭证存储到 mysql 等。

## 安装

依赖 `go >= 1.18` ，初始化 go module 后直接安装

```bash
go get github.com/duolacloud/credentials-provider
```

## 使用

### 基于 Gorm 的凭证管理

```go

import (
  "github.com/duolacloud/crud-core/repositories"
  "github.com/duolacloud/crud-cache-redis"
  "github.com/duolacloud/credentials-provider/gorm"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

// 初始化数据库
db, err := gorm.Open(mysql.Open("root:root@(localhost)/test_db?charset=utf8mb4"))
err = db.AutoMigrate(&gorm.Credential{})

// 创建凭证管理
c, err := cache.NewRedisCache()
provider := NewGormCredentialProvider(
  // gorm实例
  db,

  // 启用缓存
  WithCache(cache),
  // 缓存超时时间
  WithCacheRepositoryOptions(
    repositories.WithExpiration(5*time.Second),
  ),
)

// 设置凭证
err = provider.Set(
  context.TODO(), 
  "douyin", 
  "123456", 
  "password", 
  map[string]any{"username": "root", "password": "root"},
)

// 查询凭证
options, err = provider.Get(context.TODO(), "douyin", "123456", "password")

```

### 基于 Gorm 的凭证管理

```go
import "github.com/duolacloud/credentials-provider/gorm"

provider := NewMemoryCredentialProvider()
```