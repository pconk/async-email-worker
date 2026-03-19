[![Go CI](https://github.com/pconk/async-email-worker/actions/workflows/ci.yml/badge.svg)](https://github.com/pconk/async-email-worker/actions/workflows/ci.yml)

# 📧 Async Email Worker (Golang + Redis)

Service worker asynchronous untuk memproses job email menggunakan queue berbasis Redis. Project ini mensimulasikan pengiriman email (tanpa SMTP) untuk menunjukkan arsitektur background worker yang scalable dan production-ready.

---

## 🚀 Overview

Project ini merupakan implementasi **asynchronous processing** dengan pola:

* API → push job ke Redis
* Worker → consume & process job di background

Saat ini, pengiriman email masih berupa **simulasi (logging)**, namun arsitektur sudah siap untuk diintegrasikan dengan SMTP atau email service lainnya.

---

## 🧠 Use Case

* 📩 Simulasi pengiriman email (low stock alert)
* 🔔 Event-driven notification system
* 📤 Background job processing

---

## 🏗️ Arsitektur

```text
Client → Warehouse API → Redis Queue → Email Worker → (Simulated Email Sender)
```

---

## 🛠️ Tech Stack

* Language: Go (Golang)
* Queue: Redis (LPUSH / BRPOP)
* Logging: slog (structured logging)
* Concurrency: Goroutine Worker Pool
* Configuration: Environment Variables (.env)

---

## 📁 Struktur Project

```
.
.
├── cmd/
│   ├── api/                   # Entry point untuk Main Web Server (REST API)
│   └── worker/                # Entry point untuk Background Worker
├── internal/
│   ├── config/                # Konfigurasi aplikasi & load .env
│   ├── entity/                # Domain entities (EmailJob struct & validation)
│   ├── handler/               # HTTP Handlers (Controller) untuk API
│   ├── helper/                # Fungsi bantuan (Response, Validator wrapper)
│   ├── middleware/            # Custom HTTP middleware (Logging, Auth, dll)
│   ├── queue/                 # Logic Redis Producer (Enqueue jobs)
│   ├── service/               # Business Logic & Simulated Email Sender
│   └── worker/                # Worker Pool engine & Job Processor logic
├── pkg/
│   └── redis/                 # Shared library/driver untuk koneksi Redis
├── api_test.http              # File testing untuk REST Client (VS Code)
├── docker-compose.yaml        # Orchestration untuk Redis, API, dan Worker
└── Dockerfile                 # Multi-stage build untuk API & Worker
```

---

## ⚙️ Environment Variables

Contoh `.env`:
```
REDIS_ADDR=localhost:6379
APP_PORT=8081
QUEUE_NAME=email_queue

WORKER_CONCURRENCY=5
MAX_RETRY=3
```

---

## ▶️ Cara Menjalankan

### 1. Jalankan Redis

Jika menggunakan Docker:

```bash
docker run -d -p 6379:6379 redis
```

---

### 2. Jalankan Worker

```bash
go mod tidy
go run cmd/worker/main.go
```

---

## 📦 Format Job (Email)

Contoh payload:

```json
{
  "id": "uuid",
  "to": "user@example.com",
  "subject": "Low Stock Alert",
  "message": "Stock item hampir habis"
}
```

---

## 🔁 Mekanisme Worker

* Worker mengambil job dari Redis menggunakan `BRPOP`
* Job diproses secara concurrent (worker pool)
* Jika gagal:

  * Retry hingga batas tertentu (`MAX_RETRY`)
  * Logging error untuk monitoring

---

## 📬 Simulasi Email

Saat ini, fungsi `SendEmail` hanya melakukan:

* Logging ke console
* Menampilkan isi email (to, subject, message)

---

## 🔒 Reliability Features

* ✅ Retry mechanism
* ✅ Structured logging
* ✅ Graceful shutdown
* ✅ Configurable concurrency
* ✅ Decoupled architecture (API ↔ Worker)

---

## 🔗 Integrasi dengan Warehouse API

Worker ini terintegrasi dengan:

👉 Warehouse Management API

Flow:

1. API mendeteksi kondisi (misalnya stock < 5)
2. API push job ke Redis
3. Worker memproses job
4. Email disimulasikan melalui logging

---

## ⭐ Highlights

* Async worker berbasis Redis
* Simulasi email (tanpa dependency eksternal)
* Clean architecture & modular code
* Siap di-upgrade ke real email service

---

---
