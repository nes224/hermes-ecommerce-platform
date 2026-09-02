🛍️ Phase 3: พัฒนา Core Services ที่เหลือ (Microservices Ecosystem)
เมื่อระบบจัดการผู้ใช้ (user-service) นิ่งแล้ว จะเริ่มสร้าง Service ถัดไปตาม Flow การทำงานจริงของเว็บ E-Commerce:

product-service (ระบบจัดการสินค้า & สต็อก)

CRUD สินค้า, หมวดหมู่ (Category), แบรนด์ (Brand)

ระบบค้นหา/กรองสินค้า (Search & Filter) และจัดการจำนวนสต็อก (Inventory Control)

cart-service (ระบบตะกร้าสินค้า)

ใช้ Redis เก็บข้อมูลตะกร้าสินค้าแบบ Temporary ก่อน Checkout เพื่อความรวดเร็ว

order-service (ระบบสั่งซื้อสินค้า)

จัดการ State ของออเดอร์ (PENDING_PAYMENT, PAID, SHIPPED, CANCELLED)

คำนวณราคารวม, ส่วนลด, และตัดสต็อกสินค้า

payment-service (ระบบชำระเงิน)

เชื่อมต่อ Payment Gateway (เช่น Stripe, Omise, PromptPay QR)

รองรับ Webhook จาก Gateway เพื่ออัปเดตสถานะการชำระเงินกลับมาที่ order-service