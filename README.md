# Bazar Backend API (Golang + Fiber v3 + GORM)

Backend API untuk autentikasi (Login & Register) menggunakan:
- **Language**: Go (Golang) 1.27+
- **Framework**: [Fiber v3](https://github.com/gofiber/fiber)
- **ORM**: [GORM](https://gorm.io) dengan driver MySQL
- **Security**: Password hashing dengan `bcrypt` & Autentikasi sesi berbasis `JWT` (JSON Web Token)

---

## 📁 Struktur Folder

```text
.
├── config/
│   └── database.go          # Konfigurasi GORM MySQL & AutoMigrate
├── controllers/
│   └── auth_controller.go   # Handler HTTP (Register, Login, Me)
├── dto/
│   └── auth_dto.go          # Request & Response structs
├── middleware/
│   └── auth_middleware.go   # Validasi JWT Token Middleware
├── models/
│   └── user.go              # Model GORM entitas User
├── routes/
│   └── routes.go            # Routing endpoint API
├── utils/
│   ├── jwt.go               # JWT Generator & Validator
│   └── password.go          # Bcrypt Hash & Compare
├── .env                     # File konfigurasi lokal
├── .env.example             # Template konfigurasi
├── api_test.http            # File testing request HTTP (REST Client)
├── go.mod
├── go.sum
└── main.go                  # Entry point aplikasi
```

---

## ⚙️ Cara Menjalankan

### 1. Buat Database MySQL
Pastikan MySQL service aktif (XAMPP / Laragon / Docker / native MySQL), lalu buat database baru:
```sql
CREATE DATABASE bazar_db;
```

### 2. Sesuaikan Konfigurasi `.env`
Periksa file `.env` di root direktori:
```env
APP_PORT=8080

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=bazar_db

JWT_SECRET=super-secret-jwt-key-change-this-in-production
JWT_EXP_HOURS=24
```

### 3. Jalankan Aplikasi
```bash
go run main.go
```
> Tabel `users` akan otomatis dibuat oleh fitur **AutoMigrate** GORM saat aplikasi pertama kali dijalankan.

---

## 📌 Endpoint API

| Method | Endpoint | Auth Required | Deskripsi |
| :--- | :--- | :---: | :--- |
| `GET` | `/` | No | Health check server |
| `POST` | `/api/auth/register` | No | Pendaftaran user baru |
| `POST` | `/api/auth/login` | No | Login & mendapatkan JWT token |
| `GET` | `/api/auth/me` | **Yes (Bearer Token)** | Mengambil profil user yang sedang login |

---

### Contoh Payload & Response

#### 1. Register (`POST /api/auth/register`)
**Request Body:**
```json
{
  "name": "Nando Developer",
  "email": "nando@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "Nando Developer",
    "email": "nando@example.com",
    "created_at": "2026-08-30T08:00:00Z",
    "updated_at": "2026-08-30T08:00:00Z"
  }
}
```

#### 2. Login (`POST /api/auth/login`)
**Request Body:**
```json
{
  "email": "nando@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "Nando Developer",
      "email": "nando@example.com",
      "created_at": "2026-08-30T08:00:00Z",
      "updated_at": "2026-08-30T08:00:00Z"
    }
  }
}
```

#### 3. Me / Profile (`GET /api/auth/me`)
**Header:**
```http
Authorization: Bearer <TOKEN_DARI_LOGIN>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User profile retrieved successfully",
  "data": {
    "id": 1,
    "name": "Nando Developer",
    "email": "nando@example.com",
    "created_at": "2026-08-30T08:00:00Z",
    "updated_at": "2026-08-30T08:00:00Z"
  }
}
```
