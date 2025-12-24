# Fleetify – Purchasing Management System

Sistem manajemen pembelian barang (Procurement) dengan **Backend Golang (Fiber)** dan **Frontend jQuery** untuk mengelola Supplier, Item, dan Transaksi Pembelian.

Project ini dibuat sebagai **Technical Test – Junior Fullstack Engineer**.

---

## 📋 Daftar Isi
- [Fitur Utama](#-fitur-utama)
- [Tech Stack](#-tech-stack)
- [Prasyarat](#-prasyarat)
- [Instalasi](#-instalasi)
- [Konfigurasi](#-konfigurasi)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [Menjalankan Seeder & Akun Default](#-menjalankan-seeder--akun-default)
- [Cara Login](#-cara-login)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)

---

## 🚀 Fitur Utama

### Backend (Golang)
- Authentication & Authorization menggunakan JWT
- CRUD Items & Suppliers
- Transaksi Purchasing dengan perhitungan **SubTotal & GrandTotal di server**
- Database Transaction (ACID) untuk create purchasing
- Update stock otomatis saat transaksi berhasil
- Webhook Notification setelah transaksi sukses
- Middleware authentication
- Validasi & error handling

### Frontend (jQuery)
- Login menggunakan AJAX + JWT
- Dashboard inventory
- CRUD Items & Suppliers
- Create Purchase dengan sistem cart (tanpa refresh)
- Reusable AJAX wrapper
- Event delegation untuk elemen dinamis
- Error handling menggunakan SweetAlert2
- Responsive UI (Bootstrap)

---

## 🛠 Tech Stack

### Backend
- Go (Golang) 1.21+
- Fiber v2
- GORM
- MySQL / PostgreSQL
- JWT
- bcrypt

### Frontend
- jQuery 3.7.1
- Bootstrap 5.3
- SweetAlert2

---

## 📋 Prasyarat

Pastikan sudah terinstall:
- Go 1.21+
- MySQL atau PostgreSQL
- Git
- Browser modern

---

## 📦 Instalasi

### 1. Clone Repository
```bash
git clone https://github.com/yourusername/fleetify.git
cd fleetify
```

### 2. Setup Backend
```bash
cd backend
go mod tidy
cp .env.example .env
```

### 3. Setup Database
Contoh MySQL:
```sql
CREATE DATABASE fleetify_db;
CREATE USER 'fleetify_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON fleetify_db.* TO 'fleetify_user'@'localhost';
```

---

## ⚙️ Konfigurasi

Edit file `backend/.env`:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=fleetify_user
DB_PASSWORD=password
DB_NAME=fleetify_db

JWT_SECRET=super_secret_key
SERVER_HOST=127.0.0.1
SERVER_PORT=8080

# Optional
WEBHOOK_URL=https://webhook.site/your-id
```

---

## ▶️ Menjalankan Aplikasi

### Jalankan Backend
```bash
cd backend
go run main.go
```

Server akan berjalan di:
```
http://127.0.0.1:8080
```

### Jalankan Frontend
Frontend **tidak perlu build**, cukup buka:
```
http://127.0.0.1:8080/frontend/login.html
```

---

## 🌱 Menjalankan Seeder & Akun Default

Seeder digunakan untuk mengisi **data awal** (user, supplier, item, purchasing).

### Jalankan Seeder
```bash
cd backend
go run seeder/main.go
```

### Data yang di-generate:
- User admin
- Suppliers
- Items
- Sample purchasing transactions

Seeder bersifat **idempotent** (aman dijalankan berkali-kali).

---

## 🔐 Akun Login Default

Setelah menjalankan seeder, gunakan akun berikut:

```
Username : admin
Password : admin123
```

---

## 🔑 Cara Login

1. Buka:
```
http://127.0.0.1:8080/frontend/login.html
```

2. Masukkan akun default:
```
admin / admin123
```

3. Setelah login:
- Token JWT disimpan otomatis
- User dapat mengakses dashboard, items, suppliers, dan purchasing

---

## 📚 API Documentation

### Auth
```
POST /auth/login
POST /auth/register
```

### Items
```
GET    /item
POST   /item        (Auth)
PUT    /item/:id    (Auth)
DELETE /item/:id    (Auth)
```

### Suppliers
```
GET    /supplier
POST   /supplier        (Auth)
PUT    /supplier/:id    (Auth)
DELETE /supplier/:id    (Auth)
```

### Purchasing
```
GET    /purchasing      (Auth)
POST   /purchasing      (Auth)
```

---

## 🗄️ Database Schema

Relasi utama:
- User → Purchasing
- Supplier → Purchasing
- Purchasing → PurchasingDetail
- Item → PurchasingDetail

---

## ✅ Catatan Penting

- Perhitungan **harga, subtotal, dan grand total dilakukan di backend**
- Stock item otomatis ter-update
- Transaksi purchasing menggunakan **database transaction (rollback jika gagal)**
- Frontend hanya mengirim `item_id` dan `qty`

---

## 📄 Lisensi
MIT License
