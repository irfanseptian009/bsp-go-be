# BSP Backend (Go)

Backend API untuk sistem asuransi kebakaran BSP, dibangun dengan Go, Gin, GORM, dan PostgreSQL (Supabase).

Dokumen ini mencakup panduan setup lokal, konfigurasi environment, alur bisnis, endpoint API, testing, hingga deployment.

---

## Ringkasan Fitur

- JWT authentication (`register`, `login`)
- Role-based access control (`ADMIN`, `CUSTOMER`)
- Master data: branch & occupation type
- Insurance request workflow (create, approve, reject)
- Policy management
- Upload foto profil ke Supabase Storage
- Swagger API docs
- Integration testing dengan Jest

---

## Tech Stack

| Area | Teknologi |
|---|---|
| Language | Go 1.23 |
| Web Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL (Supabase) |
| Auth | JWT (`github.com/golang-jwt/jwt/v5`) |
| Security | bcrypt (`golang.org/x/crypto/bcrypt`) |
| API Docs | Swaggo (`swag`, `gin-swagger`) |
| Storage | Supabase Storage |
| Test | Jest + Axios |
| Container | Docker (multi-stage) |
| Deploy | Railway |

---

## Struktur Proyek

```text
bsp-go-be/
├── main.go
├── config/
├── database/
├── models/
├── dto/
├── middleware/
├── handlers/
├── services/
├── routes/
├── utils/
├── docs/
├── tests/
├── Dockerfile
└── railway.toml
```

Arsitektur menggunakan pola `handler -> service -> database/model`.

---

## Prasyarat

- Go 1.23+
- PostgreSQL (direkomendasikan Supabase)
- Node.js (untuk test Jest)

---

## Setup Lokal

1) Install dependency Go:

```bash
go mod tidy
```

2) Install dependency test:

```bash
npm install
```

3) Siapkan environment:

```bash
cp .env.example .env
```

---

## Konfigurasi Environment

| Variable | Wajib | Default | Deskripsi |
|---|---|---|---|
| `PORT` | Tidak | `3001` | Port server |
| `DATABASE_URL` | Ya | - | Koneksi PostgreSQL |
| `JWT_SECRET` | Ya | - | Secret JWT |
| `JWT_EXPIRY` | Tidak | `24h` | Masa berlaku token |
| `SUPABASE_URL` | Opsional* | - | URL project Supabase |
| `SUPABASE_SERVICE_ROLE_KEY` | Opsional* | - | Service role key Supabase |
| `SUPABASE_STORAGE_BUCKET` | Tidak | `profile` | Bucket foto profil |
| `APP_ENV` | Tidak | `development` | Mode aplikasi |
| `ENABLE_SWAGGER` | Tidak | Auto by env | Paksa on/off Swagger |
| `ALLOWED_ORIGIN` | Tidak | - | Satu origin CORS |
| `ALLOWED_ORIGINS` | Tidak | `http://localhost:3000,https://bsp-fe.netlify.app` | Banyak origin CORS |
| `SEED` | Tidak | `false` | Jalankan seed saat startup |
| `RESET_DB` | Tidak | `false` | Reset tabel saat startup |

\* Opsional jika tidak menggunakan fitur upload foto profil.

Contoh minimal:

```env
PORT=3001
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
JWT_SECRET=your-jwt-secret-key-minimum-32-characters-long
JWT_EXPIRY=24h
APP_ENV=development
ENABLE_SWAGGER=true
ALLOWED_ORIGINS=http://localhost:3000,https://bsp-fe.netlify.app
```

---

## Menjalankan Aplikasi

Jalankan server:

```bash
go run .
```

URL penting:

- Base URL: `http://localhost:3001`
- API: `http://localhost:3001/api`
- Health: `http://localhost:3001/api/health`
- Swagger: `http://localhost:3001/swagger/index.html`

Menjalankan dengan seed:

```bash
SEED=true go run .
```

Reset database saat startup:

```bash
RESET_DB=true go run .
```

Reset + seed:

```bash
RESET_DB=true SEED=true go run .
```

---

## Akun Seed Default

- `admin@bsp.com` / `admin123` (ADMIN)
- `customer@bsp.com` / `customer123` (CUSTOMER)
- `siti@bsp.com` / `customer123`
- `andi@bsp.com` / `customer123`
- `dewi@bsp.com` / `customer123`
- `rizky@bsp.com` / `customer123`

---

## Rumus Bisnis

Insurance request:

$$
basicPremium = \frac{buildingPrice \times premiumRate}{1000} \times duration
$$

$$
totalAmount = basicPremium + 10000
$$

Policy:

$$
premium = \frac{buildingPrice \times premiumRate}{1000} \times duration + 2500
$$

Format nomor otomatis:

- Invoice: `K.001.XXXXX`
- Policy approval request: `K.01.001.XXXXX`
- Policy module: `001.{year}.XXXXX`
- Application: `00001{year}XXXXXX`

---

## API Endpoints

Semua endpoint non-public wajib header:

`Authorization: Bearer <token>`

### Public

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/health`

### Users

- `GET /api/users/me`
- `PATCH /api/users/me`
- `POST /api/users/me/photo` (multipart field: `photo`)

### Branches

- `GET /api/branches`
- `GET /api/branches/:id`
- `POST /api/branches` (ADMIN)
- `PATCH /api/branches/:id` (ADMIN)
- `DELETE /api/branches/:id` (ADMIN)

### Occupation Types

- `GET /api/occupation-types`
- `GET /api/occupation-types/:id`
- `POST /api/occupation-types` (ADMIN)
- `PATCH /api/occupation-types/:id` (ADMIN)
- `DELETE /api/occupation-types/:id` (ADMIN)

### Insurance Requests

- `POST /api/insurance-requests` (CUSTOMER)
- `GET /api/insurance-requests/my-requests` (CUSTOMER)
- `GET /api/insurance-requests/invoice/:invoiceNumber`
- `GET /api/insurance-requests` (ADMIN)
- `GET /api/insurance-requests/:id`
- `PATCH /api/insurance-requests/:id/approve` (ADMIN)
- `PATCH /api/insurance-requests/:id/reject` (ADMIN)

### Policies

- `GET /api/policies` (mendukung query: `name`, `branchId`, `occupationTypeId`)
- `GET /api/policies/:id`
- `POST /api/policies` (ADMIN)
- `PATCH /api/policies/:id` (ADMIN)
- `DELETE /api/policies/:id` (ADMIN)

---

## Validasi Payload (Ringkas)

- Register: `email` valid, `password` min 6
- Insurance request: `duration` 1..10, `buildingPrice` > 0, `constructionClass` salah satu `KELAS_1|KELAS_2|KELAS_3`
- Policy: `birthDate` format `YYYY-MM-DD` atau RFC3339, `duration` 1..10, `buildingPrice` > 0

---

## Swagger Documentation

Install generator (sekali):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate docs:

```bash
swag init -g main.go -o docs
```

---

## Testing (Jest Integration)

Jalankan:

```bash
npm run test:e2e
```

Catatan:

- Secara default test mencoba menyalakan backend otomatis.
- Jika backend sudah berjalan manual:

```bash
BACKEND_GO_MANUAL_SERVER=true npm run test:e2e
```

Environment test yang didukung:

- `BACKEND_TEST_PORT`
- `BACKEND_TEST_BASE_URL`

---

## Docker

Build image:

```bash
docker build -t bsp-go-be .
```

Run container:

```bash
docker run --rm -p 3001:3001 --env-file .env bsp-go-be
```

---

## Deployment (Railway)

Konfigurasi deploy tersedia di [railway.toml](railway.toml):

- Builder: Dockerfile
- Start command: `./main`
- Healthcheck: `/api/health`
- Restart policy: `ON_FAILURE`

Pastikan environment production sudah diisi (minimal `DATABASE_URL` dan `JWT_SECRET`).

---

## Integrasi Frontend

Frontend `bsp-fe` (Next.js) menggunakan base API `http://localhost:3001/api` pada local development.

---

## Troubleshooting

### `DATABASE_URL is not set`

- Periksa `.env`
- Pastikan `DATABASE_URL` valid

### `JWT_SECRET is not set`

- Isi `JWT_SECRET` yang kuat (disarankan 32+ karakter)

### CORS error

- Tambahkan origin frontend ke `ALLOWED_ORIGINS`

### Upload foto gagal

- Periksa `SUPABASE_URL`, `SUPABASE_SERVICE_ROLE_KEY`, dan bucket
- Format didukung: JPG/JPEG/PNG/WEBP

---

## Referensi File Penting

- [main.go](main.go)
- [routes/routes.go](routes/routes.go)
- [config/config.go](config/config.go)
- [database/database.go](database/database.go)
- [database/seed.go](database/seed.go)
- [middleware/auth.go](middleware/auth.go)
- [middleware/cors.go](middleware/cors.go)
- [services/insurance_request.go](services/insurance_request.go)
- [services/policy.go](services/policy.go)
- [services/storage.go](services/storage.go)

---

