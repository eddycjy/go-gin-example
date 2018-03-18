# Go Gin Example

`gin` 的一个例子，包含许多有用特性

## 目录

本项目提供 [Gin实践](https://segmentfault.com/a/1190000013297625) 的连载示例代码

- [Gin实践 连载一 Golang介绍与环境安装](https://segmentfault.com/a/1190000013297625)
- [Gin实践 连载二 搭建Blog API's（一）](https://segmentfault.com/a/1190000013297683)
- [Gin实践 连载三 搭建Blog API's（二）](https://segmentfault.com/a/1190000013297705)
- [Gin实践 连载四 搭建Blog API's（三）](https://segmentfault.com/a/1190000013297747)
- [Gin实践 连载五 搭建Blog API's（四）](https://segmentfault.com/a/1190000013297828)
- [Gin实践 连载六 搭建Blog API's（五）](https://segmentfault.com/a/1190000013297850)
- [Gin实践 连载七 Golang优雅重启HTTP服务](https://segmentfault.com/a/1190000013757098)
- [Gin实践 连载八 为它加上Swagger](https://segmentfault.com/a/1190000013808421)

## 安装
```
$ go get github.com/EDDYCJY/go-gin-example
```

## 如何运行

### 准备

创建一个 `blog` 数据库，并且导入建表的 [SQL](https://github.com/EDDYCJY/go-gin-example/blob/master/docs/sql/blog.sql)

### 运行
```
$ cd $GOPATH/src/go-gin-example

$ go run main.go 
```

项目的运行信息和已存在的 API's

```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /auth                     --> github.com/EDDYCJY/go-gin-example/routers/api.GetAuth (3 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/EDDYCJY/go-gin-example/vendor/github.com/swaggo/gin-swagger.WrapHandler.func1 (3 handlers)
[GIN-debug] GET    /api/v1/tags              --> github.com/EDDYCJY/go-gin-example/routers/api/v1.GetTags (4 handlers)
[GIN-debug] POST   /api/v1/tags              --> github.com/EDDYCJY/go-gin-example/routers/api/v1.AddTag (4 handlers)
[GIN-debug] PUT    /api/v1/tags/:id          --> github.com/EDDYCJY/go-gin-example/routers/api/v1.EditTag (4 handlers)
[GIN-debug] DELETE /api/v1/tags/:id          --> github.com/EDDYCJY/go-gin-example/routers/api/v1.DeleteTag (4 handlers)
[GIN-debug] GET    /api/v1/articles          --> github.com/EDDYCJY/go-gin-example/routers/api/v1.GetArticles (4 handlers)
[GIN-debug] GET    /api/v1/articles/:id      --> github.com/EDDYCJY/go-gin-example/routers/api/v1.GetArticle (4 handlers)
[GIN-debug] POST   /api/v1/articles          --> github.com/EDDYCJY/go-gin-example/routers/api/v1.AddArticle (4 handlers)
[GIN-debug] PUT    /api/v1/articles/:id      --> github.com/EDDYCJY/go-gin-example/routers/api/v1.EditArticle (4 handlers)
[GIN-debug] DELETE /api/v1/articles/:id      --> github.com/EDDYCJY/go-gin-example/routers/api/v1.DeleteArticle (4 handlers)

Listening port is 8000
Actual pid is 4393
```
Swagger 文档

![image](https://sfault-image.b0.upaiyun.com/286/780/2867807553-5aae27c4ac806_articlex)

## 特性

- RESTful API
- Gorm
- Swagger
- logging
- Jwt-go
- Gin
- Graceful restart or stop (fvbock/endless)
- App configurable
