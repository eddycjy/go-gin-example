# Go Gin Example

An example of gin contains many useful features

[简体中文](https://github.com/EDDYCJY/go-gin-example/blob/master/README_ZH.md)

## Installation
```
$ go get github.com/EDDYCJY/go-gin-example
```

## How to run

### Ready

Create a **blog database** and import [SQL](https://github.com/EDDYCJY/go-gin-example/blob/master/docs/sql/blog.sql)

### Run
```
$ cd $GOPATH/src/go-gin-example

$ go run main.go 
```

Project information and existing API

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
Swagger doc

![image](https://sfault-image.b0.upaiyun.com/286/780/2867807553-5aae27c4ac806_articlex)

## Features

- RESTful API
- Gorm
- Swagger
- logging
- Jwt-go
- Gin
- Graceful restart or stop (fvbock/endless)
- App configurable
