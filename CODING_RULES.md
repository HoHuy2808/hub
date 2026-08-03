# Coding Rules & Guidelines - ConnectHub API

> Tech Stack
>
> - Go 1.25+
> - Gin
> - PostgreSQL
> - GORM
> - Redis
> - Docker

---

# 📁 Project Structure

## Feature-based Architecture

```
cmd/
    api/
    worker/

configs/

database/
    migrations/
    seeds/

internal/

    auth/

    user/

    profile/

    follow/

    post/

    media/

    comment/

    reaction/

    feed/

    story/

    notification/

    conversation/

    message/

    websocket/

    search/

    upload/

    admin/

pkg/

docs/

scripts/
```

Không tổ chức project theo layer.

❌ Sai

```
handlers/
services/
repositories/
models/
```

✅ Đúng

```
post/

comment/

user/

auth/
```

Mỗi feature tự quản lý toàn bộ code của mình.

---

# 📦 Module Structure

Mỗi module phải có cấu trúc giống nhau.

```
post/

    dto/
        request/
        response/

    handler/

    service/

    repository/

    model/

    mapper/

    validator/

    routes.go
```

Nếu module nhỏ thì có thể bỏ `mapper` hoặc `validator`.

---

# 🏗 Architecture

Dependency chỉ được theo chiều sau

```
HTTP

↓

Handler

↓

Service

↓

Repository

↓

Database
```

Không được phép

```
Handler

↓

Database
```

Hoặc

```
Handler

↓

Redis
```

---

# 🎯 Handler Rules

Handler chỉ làm các việc sau

- Parse request
- Validate request
- Lấy UserID từ Context
- Gọi Service
- Trả Response

Không được

- Business Logic
- SQL
- Redis
- Kafka
- Upload File
- Permission Check

Ví dụ

```go
func (h *PostHandler) Create(c *gin.Context) {

    var req request.CreatePostRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err)
        return
    }

    res, err := h.service.Create(c.Request.Context(), middleware.UserID(c), req)

    if err != nil {
        response.Error(c, err)
        return
    }

    response.Created(c, res)
}
```

---

# 🧠 Service Rules

Service chứa toàn bộ Business Logic.

Ví dụ

```
Create Post

↓

Check User

↓

Check Permission

↓

Save Database

↓

Publish Event

↓

Return Response
```

Service được phép gọi

- Repository
- Redis
- Queue
- Storage
- External API

Service không được gọi Handler.

---

# 🗄 Repository Rules

Repository chỉ thao tác Database.

Repository không được

- Validate
- Permission
- Redis
- Kafka
- HTTP
- Response DTO

Repository chỉ có nhiệm vụ

```
Create()

Update()

Delete()

FindByID()

FindByUserID()

Exists()

Count()
```

---

# 📄 DTO Rules

Request

```
dto/request
```

Ví dụ

```
LoginRequest

RegisterRequest

CreatePostRequest

UpdateProfileRequest
```

Response

```
dto/response
```

Ví dụ

```
PostResponse

CommentResponse

UserResponse
```

Không return Model trực tiếp.

---

# 🧱 Model Rules

Model mapping 1-1 với Database.

Ví dụ

```
User

Post

Comment

Message

Conversation
```

Không viết Business Logic trong Model.

---

# 🔌 Interface Rules

Service luôn expose interface.

```go
type Service interface {

    Create(ctx context.Context, req request.CreatePostRequest) (*response.PostResponse, error)

    Update(...)

    Delete(...)
}
```

Repository cũng tương tự.

Inject interface thay vì concrete implementation.

---

# 💉 Dependency Injection

Luôn Constructor Injection.

Đúng

```go
repo := repository.New(db)

service := service.New(repo)

handler := handler.New(service)
```

Không dùng

```
Global Variable

Singleton

init()
```

---

# 🧾 Context Rules

Tất cả Service và Repository đều phải nhận

```go
context.Context
```

Ví dụ

```go
Create(ctx context.Context,...)
```

Không tạo context mới nếu không cần.

---

# ❌ Error Handling

Không dùng

```go
errors.New("error")
```

Tạo package

```
pkg/errors
```

Ví dụ

```go
var ErrPostNotFound = New(

404,

"POST_NOT_FOUND",

"post not found",

)
```

Service chỉ return custom error.

---

# 📤 Response Format

Tất cả API phải trả về

```json
{
    "success": true,
    "message": "Success",
    "data": {}
}
```

Pagination

```json
{
    "success": true,
    "data": [],
    "pagination": {
        "page": 1,
        "limit": 20,
        "total": 100
    }
}
```

Không tự tạo response riêng.

---

# 🔐 Authentication

JWT Middleware

↓

Set UserID vào Context

Handler lấy

```go
userID := middleware.UserID(c)
```

Không parse JWT nhiều lần.

---

# 🛡 Authorization

Authorization luôn ở Service.

Ví dụ

```
Delete Post

↓

Owner ?

↓

Admin ?

↓

Delete
```

Không check Permission trong Handler.

---

# ✅ Validation

Request Validation

```
binding:"required"

binding:"email"

binding:"max=255"
```

Business Validation

```
Service
```

Ví dụ

Email format

↓

Handler

Email đã tồn tại

↓

Service

---

# 🗃 Database Rules

Database chỉ được truy cập qua Repository.

Không viết SQL trong Handler.

Không viết SQL trong Service.

---

# 🔄 Transaction Rules

Transaction chỉ được mở ở Service.

Ví dụ

```
Create Post

↓

Create Media

↓

Create Hashtag

↓

Update Feed

↓

Commit
```

Repository không tự mở Transaction.

---

# ⚡ Redis Rules

Redis chỉ dùng cho

- Cache
- Session
- Rate Limit
- Online Status

Không lưu dữ liệu chính.

---

# 📬 Event Rules

Các tác vụ lâu phải dùng Queue.

Ví dụ

```
Create Post

↓

Commit DB

↓

Publish Event

↓

Worker

↓

Send Notification

↓

Resize Image

↓

Generate Feed
```

Không chạy Goroutine trong Handler.

---

# 🌐 WebSocket Rules

Connection Manager chịu trách nhiệm

- Connect
- Disconnect
- Heartbeat
- Broadcast
- Online Status

Business Logic vẫn nằm trong Service.

---

# 📜 Logging

Không dùng

```go
fmt.Println()
```

Luôn dùng

```
zap.Logger
```

Log Levels

```
Debug

Info

Warn

Error
```

Không log Password, JWT, Refresh Token.

---

# ⚙ Configuration

Config chỉ load bằng Viper.

```
configs/

app.yaml

database.yaml

redis.yaml

jwt.yaml
```

Không hardcode.

---

# 🧩 Mapper Rules

Entity

↓

Mapper

↓

Response DTO

Không map trong Handler.

---

# 📦 Package Rules

Không tạo package

```
utils
```

Hãy chia nhỏ

```
pkg/

hash/

jwt/

pagination/

validator/

storage/

cache/

response/
```

---

# 📛 Naming Convention

## File

```
post_handler.go

post_service.go

post_repository.go

post_mapper.go
```

## Constructor

```
NewHandler()

NewService()

NewRepository()
```

## Interface

```
Repository

Service
```

## Variables

```
camelCase
```

## Struct

```
PascalCase
```

## Constants

```
UPPER_CASE
```

---

# 🧪 Testing

Service phải có Unit Test.

Repository có Integration Test.

Handler chỉ test HTTP.

Mock Repository.

Không mock Database.

---

# 🚀 Performance Rules

Không SELECT *

Chỉ preload khi cần.

Pagination mặc định.

Index cho

- user_id
- created_at
- conversation_id
- post_id

Không N+1 Query.

---

# 📋 Code Review Checklist

## Architecture

- [ ] Feature-based structure
- [ ] Clean dependency flow
- [ ] No circular dependency

## Handler

- [ ] Không có Business Logic
- [ ] Chỉ parse request
- [ ] Chỉ gọi Service

## Service

- [ ] Có Business Logic
- [ ] Authorization đúng
- [ ] Transaction đúng

## Repository

- [ ] Chỉ thao tác Database
- [ ] Không có Business Logic

## API

- [ ] Response format chuẩn
- [ ] Request validation
- [ ] Error handling

## Database

- [ ] Có Index
- [ ] Không N+1 Query
- [ ] Pagination

## Security

- [ ] JWT
- [ ] Authorization
- [ ] Không log dữ liệu nhạy cảm

## Logging

- [ ] Zap Logger
- [ ] Structured Logging

## Testing

- [ ] Unit Test
- [ ] Integration Test

---

# 📚 Best Practices

1. Handler phải mỏng.
2. Business Logic chỉ nằm ở Service.
3. Repository chỉ thao tác Database.
4. Không dùng Global Variable.
5. Không dùng package utils chung chung.
6. Inject Interface thay vì Concrete.
7. Luôn truyền context.Context.
8. Không panic trong Business Logic.
9. Không tạo goroutine trong Handler.
10. Luôn nghĩ đến khả năng scale trước khi tối ưu.