# Entity Management System – Backend Service

> 👤 **Maintained by:** **Muhammad Rifky Janzani**  
> ⚙️ **Tech Stack:** Go (Gin) • PostgreSQL 16 • Docker Compose • Clean Architecture

Service berbasis **Go (Gin)** dan **PostgreSQL** untuk mengelola data entitas yang memiliki informasi lokasi geografis (koordinat peta seperti kendaraan, perangkat IoT, fasilitas, dan aset). Dibangun dengan arsitektur **Clean Architecture (Separation of Concerns)**.

---

## 1. Cara Menjalankan Program

### A. Menggunakan Docker Compose (Sangat Direkomendasikan - 1 Perintah)
Menjalankan database PostgreSQL (port 5432) dan Backend API (port 8080) secara otomatis:
```bash
docker compose up --build -d
```
- **Backend API**: `http://localhost:8080`
- **Database**: `localhost:5432`

Untuk melihat log:
```bash
docker compose logs -f api
```

Untuk menghentikan:
```bash
docker compose down
```

---

### B. Menggunakan Local Go + Docker Database
1. Jalankan PostgreSQL di Docker:
   ```bash
   docker compose up -d postgres
   ```
2. Jalankan server Go di lokal:
   ```bash
   go run ./cmd/server/main.go
   ```

---

## 2. Alasan Pemilihan Teknologi & Library

| Teknologi / Library | Alasan Pemilihan |
| :--- | :--- |
| **Go (Golang)** | Performa tinggi, konsumsi RAM sangat rendah, kompilasi cepat menjadi single binary, dan handal menangani konkurensi data lokasi. |
| **Gin (`gin-gonic/gin`)** | Framework HTTP berbasis Radix Tree yang sangat cepat, minim alokasi memori, serta memiliki ekosistem middleware (CORS, binding validator) yang matang. |
| **PostgreSQL 16** | Database relasional dengan integritas data ACID yang kuat dan performa indeks B-Tree yang optimal untuk query koordinat lokasi. |
| **`jackc/pgx/v5`** | Driver PostgreSQL murni tercepat di Go dengan native connection pool (`pgxpool`). Jauh lebih efisien dan modern dibanding driver legacy `lib/pq` atau `database/sql`. |
| **`Masterminds/squirrel`** | SQL query builder yang fleksibel dan type-safe. Dipilih dibanding ORM berat (seperti GORM) untuk menghindari *hidden performance overhead* dan menjaga query SQL tetap transparan serta mudah dioptimasi. |
| **`golang-migrate/migrate/v4`** | Standar industri untuk migrasi database berbasis file SQL versi (`up`/`down`), otomatis tereksekusi saat startup server. |
| **`go-playground/validator/v10`** | Validasi input struct yang ketat pada layer DTO (validasi rentang koordinat latitude `-90 s/d 90`, longitude `-180 s/d 180`, dan enum). |
| **`gin-contrib/cors`** | Penanganan CORS yang fleksibel dan aman untuk integrasi frontend lintas origin lokal (`localhost`, `127.0.0.1`, dsb). |
| **Docker Multi-Stage Build** | Menghasilkan image container runtime Alpine yang sangat ringan (< 30 MB), aman, dan konsisten di berbagai environment. |

---

## 3. Struktur Folder (Separation of Concerns)

```
be-entity-info/
├── cmd/server/main.go       # Entry point, dependency injection & graceful shutdown
├── internal/
│   ├── config/              # Manajemen environment variables (.env)
│   ├── database/            # Connection pool (pgxpool) & runner migrasi SQL
│   ├── handler/             # HTTP controller, request binding & response formatting
│   ├── model/               # Domain struct, DTO requests, & response envelope
│   ├── repository/          # Query database PostgreSQL (pgx + squirrel)
│   ├── response/            # Standardisasi envelope JSON (success & error)
│   ├── router/              # Routing Gin, middleware CORS, & custom validator
│   └── service/             # Business logic (generasi ID ENT-XXXXXX, timestamps)
├── migrations/              # File SQL migrasi database
├── docs/                    # PRD backend & panduan integrasi frontend
├── AGENT.md                 # Konteks & panduan kerja Agentic AI
├── Dockerfile               # Multi-stage Docker build
└── docker-compose.yml       # Orkestrasi container PostgreSQL & API
```

---

## 4. Workflow Penggunaan Agentic AI

Proyek ini dikembangkan secara kolaboratif menggunakan pendekatan **Human-in-the-Loop Agentic AI Pair Programming** (Antigravity):

1. **Requirement & Architecture Scoping**: Membedah kebutuhan sistem, menentukan arsitektur *Separation of Concerns* (Handler ➔ Service ➔ Repository ➔ DB), dan menyusun *Implementation Plan* sebelum penulisan kode.
2. **Incremental Implementation**: Mengimplementasikan layer per layer secara terstruktur (Config ➔ DB & Migrations ➔ Repository ➔ Service ➔ Handler ➔ Router).
3. **Troubleshooting & Root Cause Analysis**: Mengidentifikasi issue *CORS 403 Forbidden* dari log Gin saat pengujian integrasi web, lalu memperbarui konfigurasi CORS secara dinamis untuk environment development.
4. **Dockerization & Testing**: Membangun multi-stage `Dockerfile` ringan, menyusun `docker-compose.yml` dengan *healthcheck dependency*, serta menguji seluruh endpoint via `curl`.

---

## 5. Status Fitur & Pengembangan Lanjutan (Future Works)

### Status Fitur: 100% Selesai & Berfungsi Penuh ✅
- [✅] CRUD Data Entitas Lengkap (`GET`, `POST`, `PUT`, `DELETE`).
- [✅] Ringkasan Metrik Dashboard (`GET /api/v1/entities/metrics`).
- [✅] Validasi Input (Koordinat Latitude/Longitude, Kategori, Status, Atribut).
- [✅] Standardisasi Format Respon JSON (`success`, `data`, `error`).
- [✅] Dockerization (PostgreSQL + Go API via Docker Compose).

---

### Pengembangan Lanjutan (Future Works) 🚀
1. **Live Telemetry via Streaming / WebSocket**: Streaming pembaruan pergerakan koordinat entitas secara real-time ke peta frontend.
2. **Kueri Area Lokasi (Spatial Query)**: Pencarian entitas berdasarkan area tampilan peta (*bounding box*) atau radius jarak tertentu.
3. **Automated Testing & CI/CD**: Penambahan unit test service layer dan integration test dengan Testcontainers pada pipeline CI/CD.

---

## 6. Referensi Dokumentasi Lainnya

- 📑 **Panduan Integrasi Frontend**: [docs/frontend-integration-prd.md](file:///home/rifky_rjanzani/takehome-test-entity/be-entity-info/docs/frontend-integration-prd.md)
- 📑 **Spesifikasi Teknis Backend**: [docs/backend-prd.md](file:///home/rifky_rjanzani/takehome-test-entity/be-entity-info/docs/backend-prd.md)
- 🤖 **Dokumentasi Konteks Agentic AI**: [AGENT.md](file:///home/rifky_rjanzani/takehome-test-entity/be-entity-info/AGENT.md)

---

## 7. Kontak & Author

👨‍💻 **Muhammad Rifky Janzani**

- 📧 **Email:** [muhammadrifkyjanzani@gmail.com](mailto:muhammadrifkyjanzani@gmail.com)
- 💼 **LinkedIn:** [linkedin.com/in/mrjanzani](https://www.linkedin.com/in/mrjanzani/)
