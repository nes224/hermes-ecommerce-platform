🚪 Phase 4: API Gateway & Event-Driven Architecture
เมื่อเริ่มมีหลาย Microservices ต้องทำให้ระบบคุยกันได้อย่างเป็นระเบียบ และมีทางเข้าเดียวสำหรับ Client:

API Gateway (Centralized Entrance)

ตั้ง API Gateway (เช่น Kong, Envoy หรือเขียนด้วย Go) ทำหน้าที่เป็น ประตูทางเข้าเดียว

ทำ Authentication Offloading: ให้ Gateway ตรวจสอบ JWT Access Token จาก user-service ก่อนจะ Forward Request ไปยัง Services อื่นๆ

ทำ Centralized Rate Limiting & CORS

Event-Driven Messaging (Async Communication)

ตั้งค่า Apache Kafka หรือ RabbitMQ เป็น Message Broker

Publish / Subscribe Pattern:

เมื่อ order-service สร้างออเดอร์สำเร็จ → ส่ง Event OrderCreated

inventory-service คอยฟังแล้วไป ตัดสต็อก

notification-service คอยฟังแล้ว ส่ง Email / SMS หาผู้ใช้