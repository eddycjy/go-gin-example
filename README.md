# Go Gin Example - Blog API[![rcard](https://goreportcard.com/badge/github.com/EDDYCJY/go-gin-example)](https://goreportcard.com/report/github.com/EDDYCJY/go-gin-example) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/EDDYCJY/go-gin-example) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/EDDYCJY/go-gin-example/master/LICENSE)

A production-ready RESTful blog API example built with Go and Gin framework, demonstrating real-world patterns and best practices.

[简体中文](https://github.com/EDDYCJY/go-gin-example/blob/master/README_ZH.md)

## Project Overview

This project is a comprehensive blog backend API system that provides complete article and tag management functionalities, along with features like JWT authentication, image upload, QR code generation, and Excel import/export.

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://github.com/jinzhu/gorm) |
| Database | MySQL |
| Cache | Redis (via [Redigo](https://github.com/gomodule/redigo)) |
| Authentication | JWT ([jwt-go](https://github.com/dgrijalva/jwt-go)) |
| Configuration | [go-ini](https://github.com/go-ini/ini) |
| API Documentation | [Swagger](https://github.com/swaggo/gin-swagger) |
| Excel Processing | [excelize](https://github.com/360EntSecGroup-Skylar/excelize), [xlsx](https://github.com/tealeg/xlsx) |
| Image Processing | [freetype](https://github.com/golang/freetype), [barcode](https://github.com/boombuler/barcode) |
| Validation | [beego/validation](https://github.com/astaxie/beego/validation) |

## Project Structure

```
go-gin-example/
├── conf/                       # Configuration files
│   └── app.ini                 # Application configuration
├── docs/                       # Documentation
│   ├── sql/                    # Database scripts
│   │   └── blog.sql            # Database schema
│   └── swagger/                # Swagger documentation
├── middleware/                 # Middleware
│   └── jwt/                    # JWT authentication middleware
│       └── jwt.go
├── models/                     # Data models (ORM)
│   ├── article.go              # Article model
│   ├── auth.go                 # Auth model
│   ├── models.go               # Database initialization
│   └── tag.go                  # Tag model
├── pkg/                        # Shared packages
│   ├── app/                    # Application utilities
│   │   ├── form.go             # Form binding
│   │   ├── request.go          # Request handling
│   │   └── response.go         # Response formatting
│   ├── e/                      # Error codes
│   │   ├── cache.go            # Cache key constants
│   │   ├── code.go             # Error code definitions
│   │   └── msg.go              # Error messages
│   ├── export/                 # Excel export utilities
│   │   └── excel.go
│   ├── file/                   # File utilities
│   │   └── file.go
│   ├── gredis/                 # Redis client
│   │   └── redis.go
│   ├── logging/                # Logging utilities
│   │   ├── file.go
│   │   └── log.go
│   ├── qrcode/                 # QR code generation
│   │   └── qrcode.go
│   ├── setting/                # Configuration management
│   │   └── setting.go
│   ├── upload/                 # Image upload utilities
│   │   └── image.go
│   └── util/                   # Common utilities
│       ├── jwt.go              # JWT utilities
│       ├── md5.go              # MD5 hashing
│       ├── pagination.go       # Pagination helper
│       └── util.go
├── routers/                    # Route definitions
│   ├── api/                    # API handlers
│   │   ├── v1/                 # API v1 handlers
│   │   │   ├── article.go      # Article endpoints
│   │   │   └── tag.go          # Tag endpoints
│   │   ├── auth.go             # Authentication endpoint
│   │   └── upload.go           # Image upload endpoint
│   └── router.go               # Route initialization
├── runtime/                    # Runtime resources
│   ├── fonts/                  # Font files
│   └── qrcode/                 # QR code resources
├── service/                    # Business logic layer
│   ├── article_service/        # Article services
│   │   ├── article.go          # Article CRUD operations
│   │   └── article_poster.go   # Poster generation
│   ├── auth_service/           # Auth services
│   │   └── auth.go
│   ├── cache_service/          # Cache key generation
│   │   ├── article.go
│   │   └── tag.go
│   └── tag_service/            # Tag services
│       └── tag.go
├── Dockerfile                  # Docker build file
├── Makefile                    # Build automation
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
└── main.go                     # Application entry point
```

## Architecture

The project follows a layered architecture pattern:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Requests                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Middleware Layer                         │
│                   (JWT Authentication)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Router Layer                            │
│                (routers/api/v1/*.go)                        │
│        - Request validation                                 │
│        - Parameter binding                                  │
│        - Response formatting                                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Service Layer                            │
│                   (service/*/*.go)                          │
│        - Business logic                                     │
│        - Cache management                                   │
│        - Cross-model operations                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Model Layer                             │
│                    (models/*.go)                            │
│        - Database operations                                │
│        - CRUD methods                                       │
│        - Data structures                                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Storage Layer                            │
│              MySQL (Primary) + Redis (Cache)                │
└─────────────────────────────────────────────────────────────┘
```

## Feature Tree

```
Go Gin Blog API
├── Authentication
│   └── JWT Login Validation
│       ├── Token Generation (3-hour expiry)
│       ├── Token Validation
│       └── Token Refresh
├── Article Management
│   ├── Create Article
│   ├── Read Article (with Redis caching)
│   ├── Update Article
│   ├── Delete Article (soft delete)
│   ├── List Articles (paginated)
│   ├── Count Articles
│   └── Generate Article Poster
│       ├── Embed QR Code
│       ├── Apply Background Image
│       ├── Render Text Overlay
│       └── Save Merged Image
├── Tag Management
│   ├── CRUD Operations
│   │   ├── Create Tag
│   │   ├── Read Tags (paginated, cached)
│   │   ├── Update Tag
│   │   └── Delete Tag (soft delete)
│   ├── Export Tags to Excel
│   └── Import Tags from Excel
├── File Upload
│   └── Image Upload
│       ├── Format validation (.jpg, .jpeg, .png)
│       ├── Size validation (max 5MB)
│       └── MD5-based naming
├── API Documentation
│   └── Swagger UI (/swagger/*any)
└── Static File Serving
    ├── Exported Excel files (/export)
    ├── Uploaded images (/upload/images)
    └── Generated QR codes (/qrcode)
```

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth` | User authentication, returns JWT token |
| GET | `/swagger/*any` | Swagger API documentation |
| POST | `/upload` | Image upload |
| POST | `/tags/export` | Export tags to Excel |
| POST | `/tags/import` | Import tags from Excel |

### Protected Endpoints (Require JWT Token)

#### Tags

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tags` | Get tag list (paginated) |
| POST | `/api/v1/tags` | Create new tag |
| PUT | `/api/v1/tags/:id` | Update tag by ID |
| DELETE | `/api/v1/tags/:id` | Delete tag by ID |

#### Articles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/articles` | Get article list (paginated) |
| GET | `/api/v1/articles/:id` | Get article by ID |
| POST | `/api/v1/articles` | Create new article |
| PUT | `/api/v1/articles/:id` | Update article by ID |
| DELETE | `/api/v1/articles/:id` | Delete article by ID |
| POST | `/api/v1/articles/poster/generate` | Generate article poster with QR code |

## Database Schema

### Tables

**blog_auth** - User authentication
```sql
- id: INT (PK, AUTO_INCREMENT)
- username: VARCHAR(50)
- password: VARCHAR(50)
```

**blog_tag** - Article tags
```sql
- id: INT (PK, AUTO_INCREMENT)
- name: VARCHAR(100) - Tag name
- created_on: INT - Creation timestamp
- created_by: VARCHAR(100) - Creator
- modified_on: INT - Modification timestamp
- modified_by: VARCHAR(100) - Modifier
- deleted_on: INT - Deletion timestamp (soft delete)
- state: TINYINT - Status (0: disabled, 1: enabled)
```

**blog_article** - Articles
```sql
- id: INT (PK, AUTO_INCREMENT)
- tag_id: INT (FK) - Associated tag ID
- title: VARCHAR(100) - Article title
- desc: VARCHAR(255) - Description
- content: TEXT - Article content
- cover_image_url: VARCHAR(255) - Cover image URL
- created_on: INT - Creation timestamp
- created_by: VARCHAR(100) - Creator
- modified_on: INT - Modification timestamp
- modified_by: VARCHAR(255) - Modifier
- deleted_on: INT - Deletion timestamp (soft delete)
- state: TINYINT - Status
```

## Configuration

Configuration is managed through `conf/app.ini`:

```ini
[app]
PageSize = 10                    # Pagination page size
JwtSecret = 233                  # JWT signing secret
PrefixUrl = http://127.0.0.1:8000
RuntimeRootPath = runtime/
ImageSavePath = upload/images/
ImageMaxSize = 5                 # Max image size in MB
ImageAllowExts = .jpg,.jpeg,.png
ExportSavePath = export/
QrCodeSavePath = qrcode/
FontSavePath = fonts/
LogSavePath = logs/

[server]
RunMode = debug                  # debug or release
HttpPort = 8000
ReadTimeout = 60                 # seconds
WriteTimeout = 60                # seconds

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

## Getting Started

### Prerequisites

- Go 1.13+
- MySQL 5.6+
- Redis

### Database Setup

1. Create a MySQL database named `blog`
2. Execute the SQL script:
```bash
mysql -u root -p blog < docs/sql/blog.sql
```

### Configuration

1. Edit `conf/app.ini` to match your environment
2. Update database credentials
3. Update Redis connection settings

### Running the Application

```bash
# Build
make build

# Run
./go-gin-example

# Or run directly
go run main.go
```

The server will start at `http://localhost:8000`

### Using Docker

```bash
# Build image
docker build -t go-gin-example .

# Run container
docker run -p 8000:8000 go-gin-example
```

## API Usage Examples

### 1. Get Authentication Token

```bash
curl -X POST http://localhost:8000/auth \
  -d "username=test&password=test123"
```

Response:
```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 2. Create a Tag

```bash
curl -X POST "http://localhost:8000/api/v1/tags?token=YOUR_TOKEN" \
  -d "name=Go&created_by=admin&state=1"
```

### 3. Get Tags List

```bash
curl "http://localhost:8000/api/v1/tags?token=YOUR_TOKEN"
```

### 4. Create an Article

```bash
curl -X POST "http://localhost:8000/api/v1/articles?token=YOUR_TOKEN" \
  -d "tag_id=1&title=Hello Gin&desc=Introduction to Gin&content=Article content...&created_by=admin&cover_image_url=http://example.com/image.jpg&state=1"
```

### 5. Upload an Image

```bash
curl -X POST http://localhost:8000/upload \
  -F "image=@/path/to/image.jpg"
```

### 6. Export Tags to Excel

```bash
curl -X POST http://localhost:8000/tags/export
```

## Key Design Patterns

### 1. Soft Delete
All models use soft delete by setting `deleted_on` timestamp instead of actual deletion.

### 2. Redis Caching
Articles and tags are cached in Redis with 1-hour TTL to reduce database load.

### 3. Service Layer Pattern
Business logic is separated into service layer, keeping handlers thin and focused on request/response handling.

### 4. Unified Response Format
All API responses follow consistent format:
```json
{
  "code": 200,
  "msg": "ok",
  "data": {}
}
```

### 5. Custom GORM Callbacks
Custom callbacks for automatic timestamp management:
- `CreatedOn` set on create
- `ModifiedOn` updated on modifications
- `DeletedOn` set on soft delete

## Error Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Invalid parameters |
| 500 | Internal server error |
| 10001 | Tag already exists |
| 10003 | Tag not found |
| 10011 | Article not found |
| 20001 | Token validation failed |
| 20002 | Token expired |
| 20003 | Token generation error |
| 20004 | Authentication failed |
| 30001 | Image save failed |
| 30002 | Image check failed |
| 30003 | Invalid image format |

## Development Commands

```bash
# Build
make build

# Run code analysis
make tool

# Run linter
make lint

# Clean build artifacts
make clean
```

## License

MIT License - See [LICENSE](LICENSE) for details.

## Credits

Project by [EDDYCJY](https://github.com/EDDYCJY/go-gin-example)
