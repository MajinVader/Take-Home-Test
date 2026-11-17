# Booking Lapangan API
Backend Take Home Test

API backend untuk sistem booking lapangan olahraga.

---

## Daftar Isi
- [Fitur](#fitur)
- [Tech Stack](#tech-stack)
- [Struktur Project](#struktur-project)
- [Prasyarat](#prasyarat)
- [Cara Menjalankan](#cara-menjalankan)
  - [Opsi A: Jalan Lokal (Go)](#opsi-a-jalan-lokal-go)
  - [Opsi B: Pakai Docker](#opsi-b-pakai-docker)
- [Membuat Akun Admin](#membuat-akun-admin)
- [Testing dengan Postman](#testing-dengan-postman)
- [Daftar Endpoint](#daftar-endpoint)
- [Contoh Request](#contoh-request)
- [Troubleshooting](#troubleshooting)
- [Catatan](#catatan)

---

## Fitur
- Register & login dengan JWT (role `user` dan `admin`)
- CRUD lapangan (create, update, delete hanya untuk admin)
- Booking lapangan dengan pengecekan jadwal bentrok (overlap)
- Payment mock: mengubah status booking dari `pending` menjadi `paid`
- Role-based authorization

## Tech Stack
- Go (go.mod: 1.21, Dockerfile: 1.22)
- Fiber v2
- GORM + PostgreSQL (versi 14-18)
- JWT (`golang-jwt/jwt/v5`)
- Bcrypt untuk hash password
- godotenv untuk membaca file `.env`

Tools pendukung: Postman, Docker, pgAdmin.

## Struktur Project
```
.
├── cmd/main.go            # entry point (load .env, connect DB, migrate, start server)
├── database/postgres.go   # koneksi PostgreSQL (baca DATABASE_URL)
├── handlers/              # logic endpoint (auth, field, booking, payment)
├── middlewares/           # AuthRequired & AdminOnly
├── models/                # User, Field, Booking
├── routes/routes.go       # daftar semua route
├── utils/jwt.go           # generate JWT
├── Dockerfile
└── Take Home Test.postman_collection.json
```

---

## Prasyarat
- Go 1.22+ (kalau jalan lokal)
- PostgreSQL yang sudah berjalan, atau Docker untuk menjalankannya
- Postman (opsional, untuk testing)

## Cara Menjalankan

### 1. Siapkan database PostgreSQL
Kalau belum punya Postgres, cara tercepat pakai Docker:

```bash
docker run -d --name booking-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=booking \
  -p 5432:5432 \
  postgres:16
```

Kalau pakai Postgres sendiri, buat database kosong dulu:

```sql
CREATE DATABASE booking;
```

Tabel `users`, `fields`, dan `bookings` dibuat otomatis oleh GORM AutoMigrate saat aplikasi pertama kali jalan, jadi tidak perlu membuat tabel manual.

### 2. Buat file `.env`
Project ini membaca konfigurasi dari environment variable `DATABASE_URL`. Buat file `.env` di root project (sejajar dengan `go.mod`):

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/booking?sslmode=disable
```

Format: `postgres://<user>:<password>@<host>:<port>/<nama_database>?sslmode=disable`

> File `.env` tidak ikut di repository. Jangan di-commit ke Git.

### Opsi A: Jalan Lokal (Go)

```bash
# 1. Download dependency
go mod download

# 2. Jalankan server (WAJIB dari root project, supaya .env terbaca)
go run ./cmd
```

Kalau berhasil, akan muncul log:

```
✅ Database connected
Server running on :3000
```

Server berjalan di **http://localhost:3000**.

Opsional, build jadi binary:

```bash
go build -o app ./cmd
./app
```

### Opsi B: Pakai Docker

```bash
# 1. Build image
docker build -t booking-api .

# 2. Jalankan container
docker run --rm -p 3000:3000 \
  -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/booking?sslmode=disable" \
  booking-api
```

Catatan untuk Docker:
- File `.env` tidak ikut masuk ke image, jadi `DATABASE_URL` harus dikirim lewat `-e` atau `--env-file .env`.
- Di dalam container, `localhost` menunjuk ke container itu sendiri, bukan ke komputer kamu. Pakai `host.docker.internal` (Windows/Mac) untuk mengakses Postgres di host.
- Di Linux, tambahkan `--add-host=host.docker.internal:host-gateway` pada perintah `docker run`, atau jalankan app dan database dalam satu Docker network lalu pakai nama container (misal `booking-db`) sebagai host.

---

## Membuat Akun Admin
Endpoint `/register` selalu membuat user dengan role `user`. Tidak ada endpoint untuk membuat admin, jadi role harus diubah manual lewat database.

1. Register dulu lewat API (lihat [Contoh Request](#contoh-request)).
2. Ubah role-nya menjadi admin:

```sql
UPDATE users SET role = 'admin' WHERE email = 'steven@mail.com';
```

Kalau pakai container Postgres dari langkah di atas:

```bash
docker exec -it booking-db psql -U postgres -d booking \
  -c "UPDATE users SET role = 'admin' WHERE email = 'steven@mail.com';"
```

3. **Login ulang** untuk mendapatkan token baru. Role disimpan di dalam token, jadi token lama tetap berisi role `user`.

---

## Testing dengan Postman
1. Import file `Take Home Test.postman_collection.json` ke Postman.
2. Variabel collection sudah ada: `base_url` (default `http://localhost:3000`) dan `token`.
3. Jalankan request **Register**, lalu **Login**.
4. Salin nilai `token` dari response login ke variabel `token` (atau ke header `Authorization`).
5. Untuk endpoint yang butuh login, kirim header:
   ```
   Authorization: Bearer <token>
   ```
6. Untuk endpoint admin, pakai token dari akun yang rolenya sudah diubah ke `admin`.

Urutan test yang disarankan:
Register -> Login -> (jadikan admin lewat SQL, login ulang) -> Create Field -> Create Booking -> Pay Booking -> Get My Bookings.

---

## Daftar Endpoint

| Method | Endpoint        | Akses  | Keterangan                        |
|--------|-----------------|--------|-----------------------------------|
| POST   | `/register`     | Publik | Daftar akun baru (role `user`)    |
| POST   | `/login`        | Publik | Login, mengembalikan JWT          |
| GET    | `/fields`       | Publik | Daftar semua lapangan             |
| GET    | `/fields/:id`   | Publik | Detail lapangan                   |
| POST   | `/fields`       | Admin  | Tambah lapangan                   |
| PUT    | `/fields/:id`   | Admin  | Update lapangan                   |
| DELETE | `/fields/:id`   | Admin  | Hapus lapangan                    |
| POST   | `/bookings`     | Login  | Buat booking                      |
| GET    | `/my-bookings`  | Login  | Booking milik user yang login     |
| GET    | `/bookings`     | Admin  | Semua booking                     |
| POST   | `/payments`     | Login  | Bayar booking (mock)              |

## Contoh Request

**Register**
```bash
curl -X POST http://localhost:3000/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Steven","email":"steven@mail.com","password":"123456"}'
```

**Login**
```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"email":"steven@mail.com","password":"123456"}'
```

**Buat lapangan (admin)**
```bash
curl -X POST http://localhost:3000/fields \
  -H "Authorization: Bearer <token_admin>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Lapangan Futsal A","price_per_hour":150000,"location":"Jl. Batununggal No. 45, Bandung"}'
```

**Buat booking**
Format waktu harus RFC3339 dan `end_time` harus setelah `start_time`.
```bash
curl -X POST http://localhost:3000/bookings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"field_id":1,"start_time":"2026-12-01T10:00:00Z","end_time":"2026-12-01T12:00:00Z"}'
```

**Bayar booking (mock)**
```bash
curl -X POST http://localhost:3000/payments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"booking_id":1}'
```

---

## Troubleshooting

| Masalah | Penyebab / Solusi |
|---------|-------------------|
| `No .env file found` lalu gagal connect DB | `.env` tidak ada atau program tidak dijalankan dari root project. Pastikan `.env` ada dan jalankan `go run ./cmd` dari folder yang sama dengan `go.mod`. |
| `failed to connect database` | Postgres belum jalan, `DATABASE_URL` salah (user, password, port, nama database), atau database belum dibuat. |
| Di Docker tidak bisa connect ke DB | Jangan pakai `localhost` di dalam container, pakai `host.docker.internal` atau nama container Postgres di Docker network yang sama. |
| `Admin only` (403) padahal sudah jadi admin | Token masih yang lama. Login ulang untuk mendapat token baru. |
| `Invalid or expired token` (401) | Token salah, kedaluwarsa (berlaku 24 jam), atau format header bukan `Bearer <token>`. |
| `Invalid start_time format, use RFC3339` | Gunakan format seperti `2026-12-01T10:00:00Z`. |
| `Time slot already booked for this field` | Jadwal bentrok dengan booking lain di lapangan yang sama. Pilih jam lain. |
| Port 3000 sudah dipakai | Hentikan proses lain yang memakai port 3000. Port ditulis langsung di `cmd/main.go`. |

---

## Catatan
- `JWT_SECRET` saat ini ditulis langsung di `utils/jwt.go`, belum dibaca dari `.env`. Untuk production sebaiknya dipindah ke environment variable.
- Pengecekan bentrok jadwal menghitung semua booking, termasuk yang berstatus `cancelled` (kalau nanti ada fitur pembatalan, filter status perlu ditambahkan).
- Pembayaran hanya simulasi: setiap request valid langsung mengubah status booking menjadi `paid`.
- Versi Go di `go.mod` (1.21) dan `Dockerfile` (1.22) berbeda. Keduanya jalan, tapi sebaiknya diseragamkan.
