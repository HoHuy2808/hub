# Backend-Hub

Backend-Hub is a social media project I builded while learning **Golang**, **Gin**, **GORM** and **WebSockets**.

## Schema
<img width="868" height="730" alt="postgres 3 - hub - public" src="https://github.com/user-attachments/assets/8cc308fb-1d5a-4cc9-8dd8-058315ed7596" />


## 🚀 Tech Stack
- **Language:** Golang
- **Framework:** Gin Web Framework
- **ORM:** GORM
- **Database:** PostgreSQL
- **Real-time:** WebSockets
- **API Documentation:** Swagger (swaggo)

## 📦 Module Architecture
The project structure:

- **Auth Module:** Handles registration, login, and JWT token issuance.
- **Post Module:** Manages personal posts (CRUD operations, pagination, search).
- **Attachment Module:** Processes file attachments (images/videos) for posts and users.
- **Comment Module:** Manages post comments and nested interactions.
- **Reaction Module:** Handles user reactions (Like, Love, Haha, etc.) on posts.
- **Request Module:** Manages the friend request system (send, accept, reject).
- **Contact Module:** Manages the user's friend list (block, unfriend, list contacts).
- **Notification/WebSocket Module:** Broadcasts real-time notifications via WebSockets.

## ⚙️ Installation & Setup

### 1. Prerequisites
- Install Go
- Install PostgreSQL
- Install `swag` to generate API documentation:
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

### 2. Database Configuration
Create a `.env` file in the root directory of the project and provide your database credentials:
```env
DB_HOST=localhost
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=hub_db
DB_PORT=5432
```

### 3. Running the Server
Generate the Swagger documentation (run this whenever the API annotations change):
```bash
swag init -g cmd/server/main.go
```

Start the server (using Air for hot-reload):
```bash
air
# Or run manually:
go run cmd/server/main.go
```
*The server runs on port `2808` by default.*

## 📚 API Documentation (Swagger)
Once the server is running successfully, you can access the Swagger UI at:
👉 `http://localhost:2808/swagger/index.html`
