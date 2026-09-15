# 📌 Project Roadmap & Status: User Service (`user-service`)

> **Architecture:** Clean Architecture / Hexagonal Architecture  
> **Tech Stack:** Go (Gin), PostgreSQL (pgx / sqlc), Redis (go-redis/v9), JWT, Viper, Docker

---

## 📊 1. Current Progress Status (สถานะปัจจุบัน)

| Layer / Component | File Location | Status | Details |
| :--- | :--- | :---: | :--- |
| **Config** | `pkg/config/config.go` | ✅ Completed | อ่านค่า `app.env` ด้วย Viper (DB, Redis, JWT, App Config) |
| **Domain Ports** | `internal/core/ports/auth.go` | ✅ Completed | นิยาม Structs, DTOs, `RedisRepository` และ `AuthService` Interface |
| **Postgres Repository** | `internal/adapters/repository/db/` | ✅ Completed | เจน Code SQL Queries (CreateUser, GetUser) ด้วย `sqlc` |
| **Redis Repository** | `internal/adapters/repository/redis_repository.go` | ✅ Completed | จัดการ Session (CreateSession, GetSession, DeleteSession) |
| **Core Business Logic** | `internal/core/services/auth_service.go` | ✅ Completed | พัฒนา `SignUp`, `SignIn`, `RefreshToken`, และ `Logout` |
| **HTTP Router** | `internal/adapters/handler/router.go` | ✅ Completed | ตั้งค่า Gin Routes (`/signup`, `/signin`, `/refresh`, `/logout`, `/health`) |
| **Main Entrypoint** | `apps/user-service/cmd/api/main.go` | ✅ Completed | ต่อ Dependency Injection (DI) + Graceful Shutdown |
| **HTTP DTOs** | `internal/adapters/handler/dto/auth_dto.go` | ⏳ Pending | แยก Request/Response Binding Structs ออกมาจาก Handler |
| **HTTP Handler** | `internal/adapters/handler/auth_handler.go` | ⏳ Pending | เชื่อม Gin Context เข้ากับ `AuthService` |

---

## 🎯 2. Next Action Items (สิ่งที่ต้องทำต่อไป)

### 📌 Task 1: สร้าง HTTP DTO (`internal/adapters/handler/dto/auth_dto.go`)
- [✅] สร้าง `SignUpRequest` (email, password, first_name, last_name, phone_number)
- [✅] สร้าง `SignInRequest` (email, password)
- [✅] สร้าง `RefreshTokenRequest` (refresh_token)
- [✅] สร้าง `LogoutRequest` (session_id)

### 📌 Task 2: เติมโค้ด HTTP Handler (`internal/adapters/handler/auth_handler.go`)
- [✅] เขียน `SignUp` -> Bind JSON -> เรียก `authService.SignUp`
- [✅] เขียน `SignIn` -> Bind JSON -> ดึง IP/UserAgent -> เรียก `authService.SignIn`
- [✅] เขียน `RefreshToken` -> Bind JSON -> เรียก `authService.RefreshToken`
- [✅] เขียน `Logout` -> Bind JSON -> Parse UUID -> เรียก `authService.Logout`

### 📌 Task 3: Local Integration & Testing
- [✅] ตรวจสอบ `docker-compose.yml` (Postgres 16, Redis 7)
- [✅] รัน Database Migration (`golang-migrate`)
- [✅] ทดสอบสั่ง `go build ./...` และ `go run apps/user-service/cmd/api/main.go`
- [✅] ยิง End-to-End Test ด้วย Postman / cURL (Signup -> Signin -> Refresh -> Logout)

---

## 🛠️ 3. Key Commands Reference (คำสั่งสั้นไว้เตือนความจำ)

```bash
# 1. เช็ค Syntax และ Compile ทั้งโปรเจกต์
go build ./...

# 2. รัน Postgres และ Redis บน Docker
docker compose up -d

# 3. รัน API Service
go run apps/user-service/cmd/api/main.go