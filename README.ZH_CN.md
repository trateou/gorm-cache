## 免责声明

这是原项目 [gorm-cache](https://github.com/Pacific73/gorm-cache) 的分支版本，原作者为 [Pacific73](https://github.com/Pacific73)。

### 为什么创建这个分支：
- 因为需要 Redis v9 支持，但原项目使用的是 Redis v8
- 原项目不再积极维护，无法及时更新依赖
- 更新了模块路径为 `github.com/trateou/gorm-cache` 以支持正常的 `go get` 安装

### 此分支的主要变更：
- 升级 Redis 客户端从 v8 到 v9 (`github.com/redis/go-redis/v9`)
- 更新了 Go 模块路径以便正确的依赖管理
- 保持与原始API的兼容性

### 致谢：
所有原始工作和设计功劳归属于 [Pacific73](https://github.com/Pacific73)。此分支旨在保持项目的活跃和可用性。

### 使用方法：
在你的项目中使用此分支：
```bash
go get github.com/trateou/gorm-cache
```


---
[![GoVersion](https://img.shields.io/github/go-mod/go-version/Pacific73/gorm-cache)](https://github.com/Pacific73/gorm-cache/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/Pacific73/gorm-cache)](https://github.com/Pacific73/gorm-cache/releases)
[![Apache-2.0 license](https://img.shields.io/badge/license-Apache2.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)

[English Version](./README.md) | [中文版本](./README.ZH_CN.md)

`gorm-cache` 旨在为gorm v2用户提供一个即插即用的旁路缓存解决方案。本缓存只适用于数据库表单主键时的场景。

我们提供2种存储介质：

1. 内存 (所有数据存储在单服务器的内存中)
2. Redis (所有数据存储在redis中，如果你有多个实例使用本缓存，那么他们不共享redis存储空间)

## 使用说明

```go
import (
    "context"
    "github.com/trateou/gorm-cache/cache"
    "github.com/redis/go-redis/v9"
)

func main() {
    dsn := "user:pass@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4"
    db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    cache, _ := cache.NewGorm2Cache(&config.CacheConfig{
        CacheLevel:           config.CacheLevelAll,
        CacheStorage:         config.CacheStorageRedis,
        RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
        InvalidateWhenUpdate: true, // when you create/update/delete objects, invalidate cache
        CacheTTL:             5000, // 5000 ms
        CacheMaxItemCnt:      5,    // if length of objects retrieved one single time
                                    // exceeds this number, then don't cache
    })
    // More options in `config.config.go`
    db.Use(cache)    // use gorm plugin
    // cache.AttachToDB(db)

    var users []User

    db.Where("value > ?", 123).Find(&users) // search cache not hit, objects cached
    db.Where("value > ?", 123).Find(&users) // search cache hit

    db.Where("id IN (?)", []int{1, 2, 3}).Find(&users) // primary key cache not hit, users cached
    db.Where("id IN (?)", []int{1, 3}).Find(&users) // primary key cache hit
}
```

在gorm中主要有5种操作（括号中是gorm中对应函数名）:

1. Query (First/Take/Last/Find/FindInBatches/FirstOrInit/FirstOrCreate/Count/Pluck)
2. Create (Create/CreateInBatches/Save)
3. Delete (Delete)
4. Update (Update/Updates/UpdateColumn/UpdateColumns/Save)
5. Row (Row/Rows/Scan)

我们不支持Row操作的缓存。
