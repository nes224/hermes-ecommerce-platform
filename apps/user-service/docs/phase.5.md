☁️ Phase 5: DevOps, Monitoring & Cloud Deployment
เตรียมระบบให้พร้อมสำหรับนำเสนอเล่ม Senior Project และเปิดใช้งานจริงบน Cloud Provider (AWS / GCP / DigitalOcean):

Container Orchestration (Kubernetes / K8s)

แปลง docker-compose.yml ไปเป็น Kubernetes Manifests (Deployments, Services, ConfigMaps, Secrets)

ตั้งค่า Horizontal Pod Autoscaler (HPA) ปรับขยายจำนวน Pod อัตโนมัติเมื่อ Traffic สูง

CI/CD Pipeline (Automated Testing & Deployment)

ใช้ GitHub Actions สร้าง Pipeline:

Push / PR -> รัน Linter + Unit Tests + Build Docker Image

Merge to main -> Push Image ไปยัง Docker Hub / ECR -> Auto Deploy ไปยัง Kubernetes Cluster

Observability Stack (Monitoring & Tracing)

Prometheus + Grafana: เก็บและแสดงผล Dashboard ของ CPU, Memory, HTTP Request Rate, Latency, Error Rate

OpenTelemetry + Jaeger: ทำ Distributed Tracing เพื่อให้เห็น Flow ของ 1 Request ที่วิ่งข้ามจาก API Gateway -> Order Service -> Payment Service ว่าติด Bottleneck ตรงไหน