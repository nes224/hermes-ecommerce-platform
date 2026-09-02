🎯 Roadmap Checklist (แผนงานสืบเนื่อง)
🚀 Phase 2.1: User Profile & RBAC (Role-Based Access Control)
[ ] User Profile Management

[ ] GET /api/v1/users/me (ดึงข้อมูลส่วนตัวจาก Access Token)

[ ] PUT /api/v1/users/me (อัปเดตข้อมูล เช่น First Name, Last Name, Phone Number)

[ ] Authorization & Middleware Setup

[ ] เขียน AuthMiddleware ใน Gin เพื่อ Validate JWT Access Token

[ ] ทำ Role Checking (เช่น USER, ADMIN, SELLER) เพื่อจำกัดสิทธิ์การเข้าถึง Endpoint

🛡️ Phase 2.2: Security & Production Hardening
[ ] Input Validation & Sanitization

[ ] เพิ่ม Struct Tags Validation ใน DTOs (เช่น binding:"required,email", min=8)

[ ] Rate Limiting

[ ] ทำ Rate Limiter ด้วย Redis เพื่อป้องกัน Brute Force Attack บน Endpoint /signin และ /signup

[ ] Observability & Logging

[ ] เปลี่ยนจาก log.Printf เป็น Structured Logging (เช่น uber-go/zap หรือ rs/zerolog)

[ ] เพิ่ม LoggerMiddleware สำหรับบันทึก HTTP Requests (Method, Path, Status, Latency, Client IP)

🧪 Phase 2.3: Unit Testing & Integration Testing
[ ] Unit Tests (Core Business Logic)

[ ] เขียน Unit Test สำหรับ AuthService โดยใช้ Mock Repository (gomock หรือ testify/mock)

[ ] Integration Tests (Repository & Database)

[ ] เขียน Test สำหรับ Postgres Queries และ Redis Operations โดยใช้ testcontainers-go

🐳 Phase 2.4: Containerization & Deployment Setup
[ ] Multi-stage Dockerfile

[ ] เขียน Dockerfile สำหรับ user-service (ใช้วิธี Multi-stage build เพื่อให้ Image มีขนาดเล็กรวมถึงปลอดภัยที่สุด)

[ ] API Documentation

[ ] เขียน API Specification ด้วย Swagger / Open API 3.0 (swaggo/swag) เพื่อให้ทีม Frontend/Mobile เรียกใช้งานได้ง่าย

🔗 Phase 2.5: Platform Integration (อนาคต)
[ ] gRPC Support (Optional)

[ ] เพิ่ม gRPC Server สำหรับให้ Microservices อื่นๆ (เช่น Order Service) ยิงมา Verify Token หรือดึง User Profile ผ่าน internal network ได้อย่างรวดเร็ว

[ ] Event-Driven Integration

[ ] ส่ง Event ออกไปที่ Message Broker (เช่น Kafka หรือ RabbitMQ) เมื่อเกิดเหตุการณ์สำคัญ เช่น UserRegisteredEvent เพื่อให้ Notification Service ส่ง Email ต้อนรับ