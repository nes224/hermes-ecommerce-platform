## 👤 Author

**Tarawut Chaisri / GitHub @nes224**
* LinkedIn: https://www.linkedin.com/in/tarawut-chaisri-245036212/
* Email: nes224@hotmail.com

# 🏛️ Hermes E-Commerce Platform

A production-grade, high-concurrency event-driven e-commerce reference architecture built with Go, Next.js, gRPC, and PostgreSQL, deployed on Kubernetes (GKE).

## 🚀 Key Architectural Highlights

* **Microservices & gRPC:** Inter-service communication via Protocol Buffers for ultra-fast execution.
* **Event-Driven Architecture:** Asynchronous event processing backed by Kafka / NATS.
* **Distributed Transactions:** Saga Orchestration Pattern for multi-service order checkout flows.
* **Stripe Integration & Outbox Pattern:** Webhook processing featuring Idempotency Checks and Transactional Outbox Pattern to resolve dual-write issues.
* **High-Concurrency Resilience:** Race condition mitigation for flash sale scenarios using Redis Distributed Locks and PostgreSQL Atomic Transactions.
* **Observability:** End-to-end distributed tracing using OpenTelemetry & Jaeger.

## 📂 Repository Structure

```text
hermes-ecommerce-platform/
├── apps/                 # Microservices (Go) & Frontend (Next.js)
├── api/proto/            # Protocol Buffer definitions
├── pkg/                  # Shared Go libraries
├── deployments/          # Docker Compose & Kubernetes Manifests
└── .github/workflows/    # CI/CD Pipelines

---
*Developed as a full-scale reference architecture for high-concurrency event-driven microservices.*

hermes-ecommerce/
├── .github/
│   └── workflows/              # GitHub Actions (CI/CD Pipelines)
│       ├── ci-backend.yaml
│       └── ci-frontend.yaml
├── apps/                       # โฟลเดอร์รวม Application หลัก
│   ├── web/                    # [FE] Next.js App Router (UI & SSE Client)
│   ├── order-service/          # [BE] Go Service (Saga Orchestrator)
│   ├── inventory-service/      # [BE] Go Service (Stock & Distributed Locks)
│   ├── payment-service/        # [BE] Go Service (Stripe Integration & Outbox Worker)
│   └── notification-service/   # [BE] Go Service (WebSockets / SSE Engine)
├── api/
│   └── proto/                  # [Shared] Protobuf files (.proto)
│       ├── order/v1/
│       ├── inventory/v1/
│       └── payment/v1/
├── pkg/                        # [Shared Go Packages] โค้ดที่ใช้ร่วมกันใน Go Services
│   ├── logger/                 # Structured Logger (Zap / Zerolog)
│   ├── otel/                   # OpenTelemetry Helper Initializer
│   └── database/               # DB Connection Pool Helpers
├── deployments/                # Infrastructure & Operations
│   ├── docker/                 # Dockerfile ของแต่ละ Service & docker-compose.yaml
│   └── k8s/                    # Kubernetes Manifests / Helm Charts
│       ├── base/
│       └── overlays/
├── go.work                     # Go Workspace File (ช่วยให้จัดการ Go Modules ใน Monorepo ได้ง่าย)
├── Makefile                    # Command Shortcuts (make proto-gen, make run-dev, ฯลฯ)
├── README.md                   # System Architecture & Documentation
└── LICENSE                     # MIT License

Diagram สถาปัตยกรรมระบบ (System Architecture)

[ FE: Next.js App Router ]
                                  (Stripe Elements / SSE Client)
                                               │
                                               ▼
                                 [ API Gateway (Envoy / NGINX) ]
                                               │
               ┌───────────────────────────────┼───────────────────────────────┐
               │ (gRPC)                        │ (gRPC)                        │ (HTTP Webhook / gRPC)
               ▼                               ▼                               ▼
     [ Order Service ]                [ Inventory Service ]           [ Payment Service ]
     • Order Management               • Stock Reservation             • Stripe SDK Integration
     • Saga Orchestrator              • Distributed Lock (Redis)      • Webhook Handler
               │                               │                               │
               ▼                               ▼                               ▼
     ┌───────────────────┐           ┌───────────────────┐           ┌───────────────────┐
     │  PostgreSQL (DB)  │           │  PostgreSQL (DB)  │           │  PostgreSQL (DB)  │
     │ - orders          │           │ - inventories     │           │ - payments        │
     │ - saga_states     │           │ - stock_locks     │           │ - processed_events│
     └───────────────────┘           └───────────────────┘           │ - outbox_messages │
                                                                     └─────────┬─────────┘
                                                                               │
                                                                   (Background Outbox Worker)
                                                                               │
                                                                               ▼
                                                                     [ Event Bus: Kafka / NATS ]
                                                                               │
                                               ┌───────────────────────────────┴───────────────────────────────┐
                                               ▼                                                               ▼
                                     [ Order Service ]                                               [ Notification Service ]
                                   (Saga Event Consumer)                                            (WebSockets / SSE Engine)
                                               │                                                               │
                                               └───────────────────────────────┬───────────────────────────────┘
                                                                               │
                                                                               ▼
                                                                [ OpenTelemetry Collector ]
                                                                               │
                                                                               ▼
                                                                  [ Jaeger / GCP Cloud Trace ]

หน้าที่และการทำงานของแต่ละส่วน (Component Responsibilities)
## Frontend (Next.js - App Router)
    - UI/UX: แสดงผลหน้า Flash Sale, ปุ่มชำระเงินด้วย Stripe Elements และแสดงสถานะคำสั่งซื้อแบบ Real-time
    - Real-time Updates: เชื่อมต่อกับ Notification Service ผ่าน SSE (Server-Sent Events) เพื่อรับแจ้งเตือนสถานะ Saga Workflow (เช่น "กำลังตัดเงิน..." -> "ตัดเงินสำเร็จ" -> "กำลังเตรียมจัดส่ง")

## Order Service (Saga Orchestrator)
    - จุดประสงค์: จัดการ Lifecycle ของ Order และทำหน้าที่เป็น Saga Orchestrator
    - Database: บันทึก Order และ State ของ Saga Transaction
    - Communication:
        - เรียก Inventory Service (gRPC) เพื่อจองสต็อกสินค้า
        - เรียก Payment Service (gRPC) เพื่อขอ client_secret ของ Stripe ไปให้ Frontend
        - คอยดึง Event จาก Kafka/NATS มาอัปเดตสถานะ Saga (เช่น OrderPaid, OrderFailed) และสั่งทำ Compensating Transaction คืนสต็อกหากชำระเงินไม่ผ่าน

## Inventory Service (Concurrency & Lock Control)
    - จุดประสงค์: ตัด/กัก/คืน สต็อกสินค้า โดยเน้นความเร็วและป้องกัน Race Condition
    - Concurrency Control: ใช้ Redis Distributed Lock หรือ Atomic PostgreSQL Update (WHERE stock >= qty) 
        เพื่อรองรับ Traffic ระดับ Flash Sale
    - Saga Actions:
        - ReserveStock (ลดสต็อกชั่วคราว)
        - ReleaseStock (Compensating: คืนสต็อกเมื่อ Payment ล้มเหลว)
        - ConfirmStock (ยืนยันการตัดสต็อกเมื่อ Payment สำเร็จ)

## Payment Service (Stripe & Transactional Outbox)
    - จุดประสงค์: สื่อสารกับ Stripe API และจัดการ Webhook แบบมีความปลอดภัยและน่าเชื่อถือสูง
    - Key Patterns ที่ใช้:
        - Idempotency Check: บันทึก Stripe event_id ลงตาราง processed_events ใน PostgreSQL ป้องกันการประมวลผล Webhook ซ้ำ
        - Transactional Outbox Pattern: เขียน Event (OrderPaid) ลงตาราง outbox_messages ภายใน DB Transaction เดียวกันกับ Local Data ป้องกันปัญหา Dual-Write
        - Background Outbox Worker: รัน Goroutine ที่ใช้ Query แบบ FOR UPDATE SKIP LOCKED ดึง Event ไปยิงเข้า Kafka/NATS โดยการันตี At-Least-Once Delivery

## Message Bus & Event-Driven Layer (Kafka / NATS JetStream)
    - จุดประสงค์: ทำหน้าที่เป็นตัวกลางแบบ Asynchronous Decoupling ช่วยให้ Microservices ทำงานสอดคล้องกันโดยไม่ต้องรอกัน (Non-blocking)
    - Events: OrderCreated, StockReserved, PaymentSucceeded, PaymentFailed, OrderPaid, StockReleased

## Observability Layer (OpenTelemetry + Jaeger)
    - จุดประสงค์: ดู Distributed Tracing ข้าม Service
    - การทำงาน: สอดใส่ Trace ID ไปใน gRPC Metadata, HTTP Headers และ Kafka/NATS Headers ทำให้สามารถแกะรอย Request 
        ตั้งแต่ Next.js -> Order     Service -> Payment Service -> Outbox Worker -> Notification Service ได้ใน Jaeger Dashboard เดียว


[User]                 [Next.js]            [Order Svc]          [Inventory Svc]        [Payment Svc]           [Stripe]             [Kafka/NATS]        [Noti Svc]
  │                        │                     │                      │                     │                    │                    │                   │
  │── 1. Press Checkout ──>│                     │                      │                     │                    │                    │                   │
  │                        │── 2. CreateOrder ──>│                      │                     │                    │                    │                   │
  │                        │   (gRPC)            │── 3. ReserveStock ──>│                     │                    │                    │                   │
  │                        │                     │   (gRPC w/ DB Lock)  │                     │                    │                    │                   │
  │                        │                     │<── Stock Reserved ───│                     │                    │                    │                   │
  │                        │                     │                      │                     │                    │                    │                   │
  │                        │                     │── 4. CreateIntent ────────────────────────>│                    │                    │                   │
  │                        │                     │   (gRPC)             │                     │── 5. PaymentIntent ─>│                    │                   │
  │                        │                     │                      │                     │<── client_secret ──│                    │                   │
  │                        │<── client_secret ───│<── client_secret ────────────────────────────│                    │                    │                   │
  │                        │                     │                      │                     │                    │                    │                   │
  │── 6. Pay with Card ───>│                     │                      │                     │                    │                    │                   │
  │   (Stripe Elements)    │──────────────────────────────────────────────────────────────────────────────────────>│                    │                   │
  │                        │                     │                      │                     │                    │                    │                   │
  │                        │                     │                      │                     │<── 7. Webhook ─────│                    │                   │
  │                        │                     │                      │                     │   (payment.succeeded)                  │                   │
  │                        │                     │                      │                     │                    │                    │                   │
  │                        │                     │                      │                     │── 8. Process Event ─┐                   │                   │
  │                        │                     │                      │                     │   (Idempotent Check)│                   │                   │
  │                        │                     │                      │                     │   (Write Outbox DB) │                   │                   │
  │                        │                     │                      │                     │<────────────────────┘                   │                   │
  │                        │                     │                      │                     │                                         │                   │
  │                        │                     │                      │                     │── 9. Outbox Worker Publish ────────────>│                   │
  │                        │                     │                      │                     │   (OrderPaid Event)                     │                   │
  │                        │                     │                      │                     │                                         │                   │
  │                        │                     │<── 10. Consume OrderPaid Event ──────────────────────────────────────────────────────│                   │
  │                        │                     │   (Complete Saga & Update Order DB)        │                                         │                   │
  │                        │                     │                      │                     │                                         │                   │
  │                        │                     │                      │                     │<── 11. Consume OrderPaid Event ─────────│──────────────────>│
  │                        │                     │                      │                     │                                         │                   │
  │<── 12. Push Status ────│────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────│
  │   (SSE: Success)       │                     │                      │                     │                                         │                   │

