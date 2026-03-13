# Go Gin Example - 博客 API [![rcard](https://goreportcard.com/badge/github.com/EDDYCJY/go-gin-example)](https://goreportcard.com/report/github.com/EDDYCJY/go-gin-example) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/EDDYCJY/go-gin-example) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/EDDYCJY/go-gin-example/master/LICENSE)

一个基于 Go 和 Gin 框架构建的生产级 RESTful 博客 API 示例，展示了真实项目的设计模式和最佳实践。


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


## 技术栈

| 分类 | 技术 |
|------|------|
| 编程语言 | Go |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://github.com/jinzhu/gorm) |
| 数据库 | MySQL |
| 缓存 | Redis (通过 [Redigo](https://github.com/gomodule/redigo)) |
| 认证 | JWT ([jwt-go](https://github.com/dgrijalva/jwt-go)) |
| 配置管理 | [go-ini](https://github.com/go-ini/ini) |
| API 文档 | [Swagger](https://github.com/swaggo/gin-swagger) |
| Excel 处理 | [excelize](https://github.com/360EntSecGroup-Skylar/excelize), [xlsx](https://github.com/tealeg/xlsx) |
| 图片处理 | [freetype](https://github.com/golang/freetype), [barcode](https://github.com/boombuler/barcode) |
| 参数校验 | [beego/validation](https://github.com/astaxie/beego/validation) |

## 项目结构

```
go-gin-example/
├── conf/                       # 配置文件目录
│   └── app.ini                 # 应用配置文件
├── docs/                       # 文档目录
│   ├── sql/                    # 数据库脚本
│   │   └── blog.sql            # 数据库表结构
│   └── swagger/                # Swagger 文档
├── middleware/                 # 中间件
│   └── jwt/                    # JWT 认证中间件
│       └── jwt.go
├── models/                     # 数据模型层 (ORM)
│   ├── article.go              # 文章模型
│   ├── auth.go                 # 认证模型
│   ├── models.go               # 数据库初始化
│   └── tag.go                  # 标签模型
├── pkg/                        # 公共包
│   ├── app/                    # 应用工具
│   │   ├── form.go             # 表单绑定
│   │   ├── request.go          # 请求处理
│   │   └── response.go         # 响应格式化
│   ├── e/                      # 错误码
│   │   ├── cache.go            # 缓存键常量
│   │   ├── code.go             # 错误码定义
│   │   └── msg.go              # 错误信息
│   ├── export/                 # Excel 导出工具
│   │   └── excel.go
│   ├── file/                   # 文件工具
│   │   └── file.go
│   ├── gredis/                 # Redis 客户端
│   │   └── redis.go
│   ├── logging/                # 日志工具
│   │   ├── file.go
│   │   └── log.go
│   ├── qrcode/                 # 二维码生成
│   │   └── qrcode.go
│   ├── setting/                # 配置管理
│   │   └── setting.go
│   ├── upload/                 # 图片上传工具
│   │   └── image.go
│   └── util/                   # 通用工具
│       ├── jwt.go              # JWT 工具
│       ├── md5.go              # MD5 哈希
│       ├── pagination.go       # 分页助手
│       └── util.go
├── routers/                    # 路由定义
│   ├── api/                    # API 处理器
│   │   ├── v1/                 # API v1 处理器
│   │   │   ├── article.go      # 文章接口
│   │   │   └── tag.go          # 标签接口
│   │   ├── auth.go             # 认证接口
│   │   └── upload.go           # 图片上传接口
│   └── router.go               # 路由初始化
├── runtime/                    # 运行时资源
│   ├── fonts/                  # 字体文件
│   └── qrcode/                 # 二维码资源
├── service/                    # 业务逻辑层
│   ├── article_service/        # 文章服务
│   │   ├── article.go          # 文章 CRUD 操作
│   │   └── article_poster.go   # 海报生成
│   ├── auth_service/           # 认证服务
│   │   └── auth.go
│   ├── cache_service/          # 缓存键生成
│   │   ├── article.go
│   │   └── tag.go
│   └── tag_service/            # 标签服务
│       └── tag.go
├── Dockerfile                  # Docker 构建文件
├── Makefile                    # 构建自动化
├── go.mod                      # Go 模块定义
├── go.sum                      # 依赖校验
└── main.go                     # 应用入口
```

## 架构设计

项目采用分层架构模式：

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP 请求                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       中间件层                               │
│                    (JWT 认证)                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       路由层                                 │
│                (routers/api/v1/*.go)                        │
│        - 请求校验                                           │
│        - 参数绑定                                           │
│        - 响应格式化                                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       服务层                                 │
│                   (service/*/*.go)                          │
│        - 业务逻辑                                           │
│        - 缓存管理                                           │
│        - 跨模型操作                                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       模型层                                 │
│                    (models/*.go)                            │
│        - 数据库操作                                         │
│        - CRUD 方法                                          │
│        - 数据结构                                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       存储层                                 │
│              MySQL (主存储) + Redis (缓存)                   │
└─────────────────────────────────────────────────────────────┘
```

## 功能树

```
Go Gin 博客 API
├── 认证模块
│   └── JWT 登录验证
│       ├── Token 生成（3小时有效期）
│       ├── Token 验证
│       └── Token 刷新
├── 文章管理
│   ├── 创建文章
│   ├── 查询文章（带 Redis 缓存）
│   ├── 更新文章
│   ├── 删除文章（软删除）
│   ├── 文章列表（分页）
│   ├── 文章统计
│   └── 生成文章海报
│       ├── 嵌入二维码
│       ├── 应用背景图
│       ├── 渲染文字叠加
│       └── 保存合成图片
├── 标签管理
│   ├── CRUD 操作
│   │   ├── 创建标签
│   │   ├── 查询标签（分页，带缓存）
│   │   ├── 更新标签
│   │   └── 删除标签（软删除）
│   ├── 导出标签到 Excel
│   └── 从 Excel 导入标签
├── 文件上传
│   └── 图片上传
│       ├── 格式校验 (.jpg, .jpeg, .png)
│       ├── 大小校验（最大 5MB）
│       └── MD5 命名
├── API 文档
│   └── Swagger UI (/swagger/*any)
└── 静态文件服务
    ├── 导出的 Excel 文件 (/export)
    ├── 上传的图片 (/upload/images)
    └── 生成的二维码 (/qrcode)
```

## API 接口

### 公开接口

| 方法 | 接口 | 描述 |
|------|------|------|
| POST | `/auth` | 用户认证，返回 JWT Token |
| GET | `/swagger/*any` | Swagger API 文档 |
| POST | `/upload` | 图片上传 |
| POST | `/tags/export` | 导出标签到 Excel |
| POST | `/tags/import` | 从 Excel 导入标签 |

### 受保护接口（需要 JWT Token）

#### 标签接口

| 方法 | 接口 | 描述 |
|------|------|------|
| GET | `/api/v1/tags` | 获取标签列表（分页） |
| POST | `/api/v1/tags` | 创建新标签 |
| PUT | `/api/v1/tags/:id` | 根据 ID 更新标签 |
| DELETE | `/api/v1/tags/:id` | 根据 ID 删除标签 |

#### 文章接口

| 方法 | 接口 | 描述 |
|------|------|------|
| GET | `/api/v1/articles` | 获取文章列表（分页） |
| GET | `/api/v1/articles/:id` | 根据 ID 获取文章 |
| POST | `/api/v1/articles` | 创建新文章 |
| PUT | `/api/v1/articles/:id` | 根据 ID 更新文章 |
| DELETE | `/api/v1/articles/:id` | 根据 ID 删除文章 |
| POST | `/api/v1/articles/poster/generate` | 生成带二维码的文章海报 |

## 数据库设计

### 数据表

**blog_auth** - 用户认证表
```sql
- id: INT (主键, 自增)
- username: VARCHAR(50) - 用户名
- password: VARCHAR(50) - 密码
```

**blog_tag** - 文章标签表
```sql
- id: INT (主键, 自增)
- name: VARCHAR(100) - 标签名称
- created_on: INT - 创建时间戳
- created_by: VARCHAR(100) - 创建人
- modified_on: INT - 修改时间戳
- modified_by: VARCHAR(100) - 修改人
- deleted_on: INT - 删除时间戳（软删除）
- state: TINYINT - 状态 (0: 禁用, 1: 启用)
```

**blog_article** - 文章表
```sql
- id: INT (主键, 自增)
- tag_id: INT (外键) - 关联标签 ID
- title: VARCHAR(100) - 文章标题
- desc: VARCHAR(255) - 文章简述
- content: TEXT - 文章内容
- cover_image_url: VARCHAR(255) - 封面图片地址
- created_on: INT - 创建时间戳
- created_by: VARCHAR(100) - 创建人
- modified_on: INT - 修改时间戳
- modified_by: VARCHAR(255) - 修改人
- deleted_on: INT - 删除时间戳（软删除）
- state: TINYINT - 状态
```

## 配置说明

配置文件位于 `conf/app.ini`：

```ini
[app]
PageSize = 10                    # 分页大小
JwtSecret = 233                  # JWT 签名密钥
PrefixUrl = http://127.0.0.1:8000
RuntimeRootPath = runtime/
ImageSavePath = upload/images/
ImageMaxSize = 5                 # 图片最大大小（MB）
ImageAllowExts = .jpg,.jpeg,.png
ExportSavePath = export/
QrCodeSavePath = qrcode/
FontSavePath = fonts/
LogSavePath = logs/

[server]
RunMode = debug                  # debug 或 release
HttpPort = 8000
ReadTimeout = 60                 # 秒
WriteTimeout = 60                # 秒

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
```

## 快速开始

### 环境要求

- Go 1.13+
- MySQL 5.6+
- Redis

### 数据库初始化

1. 创建名为 `blog` 的 MySQL 数据库
2. 执行 SQL 脚本：
```bash
mysql -u root -p blog < docs/sql/blog.sql
```

### 配置修改

1. 编辑 `conf/app.ini` 以匹配您的环境
2. 更新数据库连接信息
3. 更新 Redis 连接配置

### 运行应用

```bash
# 构建
make build

# 运行
./go-gin-example

# 或直接运行
go run main.go
```

服务将启动在 `http://localhost:8000`

### 使用 Docker

```bash
# 构建镜像
docker build -t go-gin-example .

# 运行容器
docker run -p 8000:8000 go-gin-example
```

## API 使用示例

### 1. 获取认证 Token

```bash
curl -X POST http://localhost:8000/auth \
  -d "username=test&password=test123"
```

响应：
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 2. 创建标签

```bash
curl -X POST "http://localhost:8000/api/v1/tags?token=YOUR_TOKEN" \
  -d "name=Go&created_by=admin&state=1"
```

### 3. 获取标签列表

```bash
curl "http://localhost:8000/api/v1/tags?token=YOUR_TOKEN"
```

### 4. 创建文章

```bash
curl -X POST "http://localhost:8000/api/v1/articles?token=YOUR_TOKEN" \
  -d "tag_id=1&title=Hello Gin&desc=Gin 入门介绍&content=文章内容...&created_by=admin&cover_image_url=http://example.com/image.jpg&state=1"
```

### 5. 上传图片

```bash
curl -X POST http://localhost:8000/upload \
  -F "image=@/path/to/image.jpg"
```

### 6. 导出标签到 Excel

```bash
curl -X POST http://localhost:8000/tags/export
```

## 核心设计模式

### 1. 软删除
所有模型使用软删除，通过设置 `deleted_on` 时间戳而非实际删除数据。

### 2. Redis 缓存
文章和标签数据缓存在 Redis 中，TTL 为 1 小时，以减少数据库负载。

### 3. 服务层模式
业务逻辑分离到服务层，保持处理器精简，专注于请求/响应处理。

### 4. 统一响应格式
所有 API 响应遵循统一格式：
```json
{
  "code": 200,
  "msg": "ok",
  "data": {}
}
```

### 5. 自定义 GORM 回调
自定义回调实现自动时间戳管理：
- `CreatedOn` 在创建时设置
- `ModifiedOn` 在修改时更新
- `DeletedOn` 在软删除时设置

## 错误码说明

| 错误码 | 描述 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数无效 |
| 500 | 服务器内部错误 |
| 10001 | 标签已存在 |
| 10003 | 标签不存在 |
| 10011 | 文章不存在 |
| 20001 | Token 验证失败 |
| 20002 | Token 已过期 |
| 20003 | Token 生成错误 |
| 20004 | 认证失败 |
| 30001 | 图片保存失败 |
| 30002 | 图片检查失败 |
| 30003 | 图片格式无效 |

## 开发命令

```bash
# 构建
make build

# 运行代码分析
make tool

# 运行代码检查
make lint

# 清理构建产物
make clean
```

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件。

## 联系我

![image](https://image.eddycjy.com/7074be90379a121746146bc4229819f8.jpg)
