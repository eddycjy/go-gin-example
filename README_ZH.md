# Go Gin Example [![rcard](https://goreportcard.com/badge/github.com/EDDYCJY/go-gin-example)](https://goreportcard.com/report/github.com/EDDYCJY/go-gin-example) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/EDDYCJY/go-gin-example) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/EDDYCJY/go-gin-example/master/LICENSE)

`gin` 的一个例子，包含许多有用特性

## 目录

本项目提供 [Gin实践](https://segmentfault.com/a/1190000013297625) 的连载示例代码

1. [Gin实践 连载一 Golang介绍与环境安装](https://book.eddycjy.com/golang/gin/install.html)
2. [Gin实践 连载二 搭建Blog API's（一）](https://book.eddycjy.com/golang/gin/api-01.html)
3. [Gin实践 连载三 搭建Blog API's（二）](https://book.eddycjy.com/golang/gin/api-02.html)
4. [Gin实践 连载四 搭建Blog API's（三）](https://book.eddycjy.com/golang/gin/api-03.html)
5. [Gin实践 连载五 使用JWT进行身份校验](https://book.eddycjy.com/golang/gin/jwt.html)
6. [Gin实践 连载六 编写一个简单的文件日志](https://book.eddycjy.com/golang/gin/log.html)
7. [Gin实践 连载七 Golang优雅重启HTTP服务](https://book.eddycjy.com/golang/gin/reload-http.html)
8. [Gin实践 连载八 为它加上Swagger](https://book.eddycjy.com/golang/gin/swagger.html)
9. [Gin实践 连载九 将Golang应用部署到Docker](https://book.eddycjy.com/golang/gin/golang-docker.html)
10. [Gin实践 连载十 定制 GORM Callbacks](https://book.eddycjy.com/golang/gin/gorm-callback.html)
11. [Gin实践 连载十一 Cron定时任务](https://book.eddycjy.com/golang/gin/cron.html)
12. [Gin实践 连载十二 优化配置结构及实现图片上传](https://book.eddycjy.com/golang/gin/config-upload.html)
13. [Gin实践 连载十三 优化你的应用结构和实现Redis缓存](https://book.eddycjy.com/golang/gin/application-redis.html)
14. [Gin实践 连载十四 实现导出、导入 Excel](https://book.eddycjy.com/golang/gin/excel.html)
15. [Gin实践 连载十五 生成二维码、合并海报](https://book.eddycjy.com/golang/gin/image.html)
16. [Gin实践 连载十六 在图片上绘制文字](https://book.eddycjy.com/golang/gin/font.html)
17. [Gin实践 连载十七 用 Nginx 部署 Go 应用](https://book.eddycjy.com/golang/gin/nginx.html)
18. [Gin实践 番外 Golang交叉编译](https://book.eddycjy.com/golang/gin/cgo.html)
19. [Gin实践 番外 请入门 Makefile](https://book.eddycjy.com/golang/gin/makefile.html)

## 安装
```
$ go get github.com/EDDYCJY/go-gin-example
```

## 如何运行

### 必须

- Mysql
- Redis

### 准备

创建一个 `blog` 数据库，并且导入建表的 [SQL](https://github.com/EDDYCJY/go-gin-example/blob/master/docs/sql/blog.sql)

### 配置

你应该修改 `conf/app.ini` 配置文件

```
[database]
Type = mysql
User = root
Password = rootroot
Host = 127.0.0.1:3306
Name = blog
TablePrefix = blog_

[redis]
Host = 127.0.0.1:6379
Password =
MaxIdle = 30
MaxActive = 30
IdleTimeout = 200
...
```


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

![image](https://i.imgur.com/bVRLTP4.jpg)

## 特性

- RESTful API
- Gorm
- Swagger
- logging
- Jwt-go
- Gin
- Graceful restart or stop (fvbock/endless)
- App configurable
- Cron
- Redis

## 联系我

![image](https://image.eddycjy.com/7074be90379a121746146bc4229819f8.jpg)
