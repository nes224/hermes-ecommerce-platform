# 📊 User Service Load Testing Guide & Benchmark Report

เอกสารนี้รวบรวมแนวทางการทดสอบประสิทธิภาพ (Performance & Load Testing) ของ `user-service` ในระบบ **Hermes E-Commerce Platform** เพื่อประเมินความสามารถในการรองรับ Request (Throughput), เวลาในการตอบสนอง (Latency), และความเสถียรของระบบ (System Stability) เมื่อทำงานร่วมกับ PostgreSQL และ Redis Session Storage

---

## 🏗️ 1. Test Architecture & Environment Specification

* **Target Microservice:** `user-service` (Go / Gin Framework / Clean Architecture)
* **Database Pool:** PostgreSQL 16 (via `pgxpool`)
* **Session Cache:** Redis 7 (In-Memory Key-Value Storage)
* **Load Testing Tool:** [Grafana k6](https://k6.io/)

---

## 🛠️ 2. Prerequisite & Installation

ติดตั้ง `k6` บนเครื่องส่วนตัวสำหรับการทดสอบผ่าน Homebrew (macOS):

```bash
brew install k6