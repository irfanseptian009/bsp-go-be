# FIMS Backend (Go + Gin)

Fire Insurance Management System — Backend API built with **Go**, **Gin**, and **GORM**, connected to **Supabase PostgreSQL**.

## Tech Stack

| Technology | Purpose |
|---|---|
| [Go 1.22+](https://go.dev/) | Language |
| [Gin](https://gin-gonic.com/) | HTTP framework |
| [GORM](https://gorm.io/) | ORM |
| [PostgreSQL (Supabase)](https://supabase.com/) | Database |
| [JWT](https://github.com/golang-jwt/jwt) | Authentication |
| [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) | Password hashing |

## Project Structure

```
backend-go/
├── main.go                      # Entry point
├── config/config.go             # Environment config
├── database/
│   ├── database.go              # DB connection + migration
│   └── seed.go                  # Seed data
├── models/                      # GORM models
├── dto/                         # Request/Response structs
├── middleware/                   # Auth, CORS, Role guard
├── handlers/                    # HTTP handlers (controllers)
├── services/                    # Business logic
├── routes/routes.go             # Route registration
└── utils/response.go            # Response helpers
```

## Quick Start

### 1. Install Go

Download from [https://go.dev/dl/](https://go.dev/dl/) or use Homebrew:

```bash
brew install go
```

### 2. Setup Environment

```bash
cd backend-go
cp .env.example .env
```

Edit `.env` with your Supabase credentials:

```env
PORT=3001
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
JWT_SECRET=your-jwt-secret-key-minimum-32-characters-long
JWT_EXPIRY=24h
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Server

```bash
go run main.go
```

Swagger UI tersedia di:

```text
http://localhost:3001/swagger/index.html
```

### 5. Seed Database (First Time Only)

```bash
SEED=true go run main.go
```

This creates:
- 5 branches (Kuningan, Tebet, Harmoni, Sudirman, Kelapa Gading)
- 5 occupation types (Rumah, Ruko, Gedung Kantor, Gudang, Apartemen)
- 1 admin account: `admin@fims.com` / `admin123`
- 3 customer accounts: `customer@fims.com`, `siti@fims.com`, `andi@fims.com` / `customer123`

## API Documentation (Swagger)

1. Install generator (sekali saja):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate docs:

```bash
swag init -g main.go -o docs
```

3. Jalankan server lalu buka:

```text
http://localhost:3001/swagger/index.html
```

## Testing (Jest E2E)

Testing disiapkan di folder [tests](tests) untuk menguji endpoint backend-go via HTTP.

```bash
npm install
npm run test:e2e
```

Catatan:
- Secara default test akan mencoba menjalankan backend-go otomatis.
- Jika server sudah jalan manual, jalankan dengan env:

```bash
BACKEND_GO_MANUAL_SERVER=true npm run test:e2e
```

## API Endpoints

### Auth (Public)
- `POST /api/auth/register` — Register new account
- `POST /api/auth/login` — Login

### Users (Authenticated)
- `GET /api/users/me` — Get profile
- `PATCH /api/users/me` — Update profile

### Branches
- `GET /api/branches` — List all
- `GET /api/branches/:id` — Get by ID
- `POST /api/branches` — Create (Admin)
- `PATCH /api/branches/:id` — Update (Admin)
- `DELETE /api/branches/:id` — Delete (Admin)

### Occupation Types
- `GET /api/occupation-types` — List all
- `GET /api/occupation-types/:id` — Get by ID
- `POST /api/occupation-types` — Create (Admin)
- `PATCH /api/occupation-types/:id` — Update (Admin)
- `DELETE /api/occupation-types/:id` — Delete (Admin)

### Insurance Requests
- `POST /api/insurance-requests` — Create request (Customer)
- `GET /api/insurance-requests/my-requests` — My requests (Customer)
- `GET /api/insurance-requests` — List all (Admin)
- `GET /api/insurance-requests/:id` — Get by ID
- `PATCH /api/insurance-requests/:id/approve` — Approve (Admin)
- `PATCH /api/insurance-requests/:id/reject` — Reject (Admin)

### Policies
- `GET /api/policies` — List/search
- `GET /api/policies/:id` — Get by ID
- `POST /api/policies` — Create (Admin)
- `PATCH /api/policies/:id` — Update (Admin)
- `DELETE /api/policies/:id` — Delete (Admin)

## Connecting to Frontend

The frontend Next.js app already points to `http://localhost:3001/api`. Just start the Go backend on port 3001 (default) and the frontend will connect automatically.
