# Go Gin Example [![rcard](https://goreportcard.com/badge/github.com/EDDYCJY/go-gin-example)](https://goreportcard.com/report/github.com/EDDYCJY/go-gin-example) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/EDDYCJY/go-gin-example) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/EDDYCJY/go-gin-example/master/LICENSE)

`gin` 的一个例子，包含许多有用特性

## 目录

本项目提供 [Gin实践](https://segmentfault.com/a/1190000013297625) 的连载示例代码

1. [Gin实践 连载一 Golang介绍与环境安装](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-02-16-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E4%B8%80-Golang%E4%BB%8B%E7%BB%8D%E4%B8%8E%E7%8E%AF%E5%A2%83%E5%AE%89%E8%A3%85.md)
2. [Gin实践 连载二 搭建Blog API's（一）](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-02-16-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E4%BA%8C-%E6%90%AD%E5%BB%BABlogAPIs-01.md)
3. [Gin实践 连载三 搭建Blog API's（二）](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-02-16-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E4%B8%89-%E6%90%AD%E5%BB%BABlogAPIs-02.md)
4. [Gin实践 连载四 搭建Blog API's（三）](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-02-16-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%9B%9B-%E6%90%AD%E5%BB%BABlogAPIs-03.md)
5. [Gin实践 连载五 使用JWT进行身份校验](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-02-16-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E4%BA%94-%E4%BD%BF%E7%94%A8JWT%E8%BF%9B%E8%A1%8C%E8%BA%AB%E4%BB%BD%E6%A0%A1%E9%AA%8C.md)
6. [Gin实践 连载六 编写一个简单的文件日志](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-02-16-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%85%AD-%E7%BC%96%E5%86%99%E4%B8%80%E4%B8%AA%E7%AE%80%E5%8D%95%E7%9A%84%E6%96%87%E4%BB%B6%E6%97%A5%E5%BF%97.md)
7. [Gin实践 连载七 Golang优雅重启HTTP服务](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-03-15-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E4%B8%83-Golang%E4%BC%98%E9%9B%85%E9%87%8D%E5%90%AFHTTP%E6%9C%8D%E5%8A%A1.md)
8. [Gin实践 连载八 为它加上Swagger](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-03-18-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%85%AB-%E4%B8%BA%E5%AE%83%E5%8A%A0%E4%B8%8ASwagger.md)
9. [Gin实践 连载九 将Golang应用部署到Docker](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-03-24-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E4%B9%9D-%E5%B0%86Golang%E5%BA%94%E7%94%A8%E9%83%A8%E7%BD%B2%E5%88%B0Docker.md)
10. [Gin实践 连载十 定制 GORM Callbacks](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-04-15-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%8D%81-%E5%AE%9A%E5%88%B6GORM-Callbacks.md)
11. [Gin实践 连载十一 Cron定时任务](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-04-29-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%8D%81%E4%B8%80-Cron%E5%AE%9A%E6%97%B6%E4%BB%BB%E5%8A%A1.md)
12. [Gin实践 连载十二 优化配置结构及实现图片上传](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-05-27-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%8D%81%E4%BA%8C-%E4%BC%98%E5%8C%96%E9%85%8D%E7%BD%AE%E7%BB%93%E6%9E%84%E5%8F%8A%E5%AE%9E%E7%8E%B0%E5%9B%BE%E7%89%87%E4%B8%8A%E4%BC%A0.md)
13. [Gin实践 连载十三 优化你的应用结构和实现Redis缓存](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-06-02-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%8D%81%E4%B8%89-%E4%BC%98%E5%8C%96%E4%BD%A0%E7%9A%84%E5%BA%94%E7%94%A8%E7%BB%93%E6%9E%84%E5%92%8C%E5%AE%9E%E7%8E%B0Redis%E7%BC%93%E5%AD%98.md)
14. [Gin实践 连载十四 实现导出、导入 Excel](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-06-14-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%8D%81%E5%9B%9B-%E5%AE%9E%E7%8E%B0%E5%AF%BC%E5%87%BA%E3%80%81%E5%AF%BC%E5%85%A5-Excel.md)
15. [Gin实践 连载十五 生成二维码、合并海报](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-07-04-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%8D%81%E4%BA%94-%E7%94%9F%E6%88%90%E4%BA%8C%E7%BB%B4%E7%A0%81-%E5%90%88%E5%B9%B6%E6%B5%B7%E6%8A%A5.md)
16. [Gin实践 连载十六 在图片上绘制文字](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-07-07-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%8D%81%E5%85%AD-%E5%9C%A8%E5%9B%BE%E7%89%87%E4%B8%8A%E7%BB%98%E5%88%B6%E6%96%87%E5%AD%97.md)
17. [Gin实践 连载十七 用 Nginx 部署 Go 应用](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-09-01-Gin%E5%AE%9E%E8%B7%B5-%E8%BF%9E%E8%BD%BD%E5%8D%81%E4%B8%83-%E7%94%A8%20Nginx%20%E9%83%A8%E7%BD%B2%20Go%20%E5%BA%94%E7%94%A8.md)
18. [Gin实践 番外 Golang交叉编译](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-03-26-Gin%E5%AE%9E%E8%B7%B5-%E7%95%AA%E5%A4%96-Golang%E4%BA%A4%E5%8F%89%E7%BC%96%E8%AF%91.md)
19. [Gin实践 番外 请入门 Makefile](https://github.com/EDDYCJY/blog/blob/master/golang/gin/2018-08-26-Gin%E5%AE%9E%E8%B7%B5-%E7%95%AA%E5%A4%96-%E8%AF%B7%E5%85%A5%E9%97%A8%20Makefile.md)

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
