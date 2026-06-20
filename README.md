# ☕ Kafe Omborxonasini Boshqarish Tizimi (Cafe Warehouse Management System)

Ushbu loyiha kafe omboridagi mahsulotlar harakatini to'liq nazorat qilish, moliyaviy hisobotlarni avtomatlashtirish, real vaqt rejimida ombor qoldig'ini va sof foydani hisoblash uchun ishlab chiqilgan zamonaviy full-stack veb-ilovadir.

---

## 🚀 Texnologiyalar Steki (Technology Stack)

Tizim yuqori unumdorlik, xavfsizlik va tezlikni ta'minlash maqsadida quyidagi texnologiyalar asosida qurilgan:

*   **Backend:**
    *   **Go (Golang) 1.22+**: Yuqori tezlik va tranzaksiyalar xavfsizligini ta'minlash uchun.
    *   **`net/http` standard library**: HTTP router sifatida ishlatilgan.
    *   **`pgx/v5`**: PostgreSQL bilan ishlash uchun eng tezkor va zamonaviy driver.
    *   **JWT (JSON Web Tokens)**: Stateless autentifikatsiya va sessiyalarni boshqarish uchun.
    *   **`bcrypt`**: Foydalanuvchilar parollarini xavfsiz heshlash (hashing) uchun.
*   **Frontend:**
    *   **Next.js 14 (App Router, TypeScript)**: Server-side rendering (SSR), qulay route guruhlari va tezkor ishlash uchun.
    *   **Vanilla CSS**: Premium dark-theme (qorong'u rejim) va glassmorphism dizayn tizimi.
*   **Ma'lumotlar Bazasi (Database):**
    *   **PostgreSQL 16**: Tranzaksiyaviy ma'lumotlar butunligi va ACID talablariga javob beruvchi relyatsion baza.
*   **Infratuzilma (Infrastructure):**
    *   **Docker & Docker Compose**: Loyihani barcha platformalarda oson ishga tushirish uchun.
    *   **Nginx**: Reverse proxy va statik fayllarni tarqatish uchun.

---

## 📂 Loyiha Strukturasi (Project Structure)

Loyiha ikkita asosiy modulga va sozlamalarga bo'lingan:

```text
kafe omborxona/
├── backend/                  # backend (Go API)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go       # Dasturni ishga tushirish nuqtasi (Dependency Injection)
│   ├── internal/
│   │   ├── config/           # .env sozlamalarni yuklovchi modul
│   │   ├── domain/           # Ma'lumotlar turlari (Entities) va interfeyslar
│   │   ├── handler/          # HTTP so'rovlarni qabul qiluvchi handlerlar (helpers, auth, product, va b.)
│   │   ├── middleware/       # Auth, RBAC (rolni tekshirish) va CORS filtrlari
│   │   ├── repository/       # Ma'lumotlar bazasi bilan ishlovchi SQL so'rovlar
│   │   └── service/          # Biznes logikasi (Auth logic, foydalanuvchilar boshqaruvi)
│   ├── migrations/           # SQL ma'lumotlar bazasi migratsiyalari
│   ├── .env                  # Backend uchun maxfiy sozlamalar
│   └── Dockerfile            # Go uchun Docker tasviri
│
├── frontend/                 # frontend (Next.js)
│   ├── src/
│   │   ├── app/              # Next.js App Router sahifalari
│   │   │   ├── (auth)/       # Tizimga kirish (Login) sahifasi
│   │   │   ├── (dashboard)/  # Ichki dashboard sahifalari (products, categories, transactions, va hk.)
│   │   │   ├── globals.css   # Premium dark mode dizayn va vizual uslublar
│   │   │   └── layout.tsx    # Root layout sozlamalari
│   │   ├── components/       # Umumiy UI komponentlar (Sidebar, Modal)
│   │   └── lib/              # API mijoz, turlar (types) va yordamchi funksiyalar
│   ├── .env.local            # Frontend API URL sozlamasi
│   └── Dockerfile            # Next.js uchun Docker tasviri
│
├── nginx/
│   └── nginx.conf            # Nginx proxy konfiguratsiyasi
└── docker-compose.yml        # Docker Compose orqali butun tizimni ishga tushirish fayli
```

---

## 🗄️ Ma'lumotlar Bazasi Sxemasi (Database Schema)

Ma'lumotlar relyatsion tuzilmaga ega bo'lib, quyidagi jadvallardan iborat:

1.  **`users`**: Tizim foydalanuvchilari.
    *   `id` (Serial, PK), `username` (Unique), `password_hash`, `full_name`, `role` (`admin` yoki `warehouse_manager`).
2.  **`categories`**: Mahsulotlar uchun kategoriyalar (masalan: Go'sht mahsulotlari, Ichimliklar, Sabzavotlar).
    *   `id` (Serial, PK), `name` (Unique).
3.  **`suppliers`**: Yetkazib beruvchilar ro'yxati.
    *   `id` (Serial, PK), `name`, `phone`, `address`.
4.  **`products`**: Mahsulotlarning asosiy kartotekasi.
    *   `id` (Serial, PK), `name`, `unit` (o'lchov birligi: kg, litr, dona), `category_id` (FK), `cost_price` (tan narxi), `sale_price` (sotuv narxi).
5.  **`transactions`**: Ombordagi barcha kirim, chiqim va hisobdan chiqarish harakatlari.
    *   `id` (Serial, PK), `product_id` (FK), `supplier_id` (FK, nullable), `user_id` (FK), `type` (`purchase` - kirim, `return` - qaytarish, `sale` - sotuv, `write_off` - yaroqsiz deb hisobdan chiqarish), `quantity` (miqdor), `unit_price` (tranzaksiya narxi), `total_amount` (umumiy summa), `note` (izoh).

---

## 🔐 Foydalanuvchilar Rollari (RBAC)

Tizimda xavfsizlik va vakolatlarni cheklash maqsadida Role-Based Access Control (RBAC) tizimi o'rnatilgan:

| Sahifa / Funksiya | Admin (Boshqaruvchi) | Ombor Mudiri (Warehouse Manager) |
| :--- | :---: | :---: |
| **Dashboard Ko'rsatkichlari** | Jami kirim, chiqim, foyda va ombor qiymati | Faqat qoldiq va tranzaksiyalar soni (Moliya yopiq) |
| **Ombor Qoldig'i (`/inventory`)** | ✅ Ko'ra oladi | ✅ Ko'ra oladi |
| **Mahsulotlar va Kategoriyalar** | ✅ To'liq boshqaruv (CRUD) | ✅ To'liq boshqaruv (CRUD) |
| **Tranzaksiyalar yaratish** | ✅ Kirim, chiqim, vozvrat qilish | ✅ Kirim, chiqim, vozvrat qilish |
| **Tranzaksiyani o'chirish/bekor qilish** | ✅ Ruxsat berilgan | ❌ Ruxsat berilmagan |
| **Moliyaviy Hisobotlar (`/reports`)** | ✅ Kirish huquqi bor | ❌ Kirish huquqi yopiq |
| **Foydalanuvchilar Boshqaruvi (`/users`)** | ✅ Yangi foydalanuvchilar yaratish va tahrirlash | ❌ Kirish huquqi yopiq |

---

## 📊 Moliyaviy Foyda Formulasi (Net Profit Formula)

Sof foyda tizimda avtomatik ravishda quyidagi formula asosida hisoblanadi:

$$\text{Sof Foyda} = \text{Jami Sotuv Summasi} - (\text{Sotilgan Mahsulotlar Tan Narxi} + \text{Yo'qotishlar (Spisaniye) Summasi})$$

Go backend-da bu quyidagicha ishlaydi:
*   **Sotuv summasi**: `type = 'sale'` bo'lgan barcha tranzaksiyalarning `total_amount` yig'indisi.
*   **Tan narxi (Cost of Sold)**: Har bir sotilgan mahsulotning miqdori ko'paytirilgan uning ayni vaqtdagi tan narxiga (`quantity * cost_price`).
*   **Yo'qotishlar (Write-off losses)**: Yaroqsiz deb topilgan yoki muddati o'tgan mahsulotlarning tan narx bo'yicha hisoblangan umumiy qiymati.

---

## 🛠️ Loyihani Ishga Tushirish Qo'llanmasi

### 🐋 Docker yordamida (Tavsiya etiladi)
Barcha sozlamalar (PostgreSQL, Go Backend, Next.js, Nginx) Docker Compose yordamida avtomatlashtirilgan. Sizdan faqat bitta buyruq talab qilinadi:

```bash
docker compose up --build
```
Ushbu buyruq ishga tushgandan so'ng, brauzerda **`http://localhost`** manziliga kiring.

*   **Tizimga kirish ma'lumotlari:**
    *   **Login:** `admin`
    *   **Parol:** `admin123`

---

### 💻 Lokal Tarzda Ishga Tushirish (Development Mode)

Agar loyihani ishlab chiquvchi (development) rejimida alohida ishga tushirmoqchi bo'lsangiz:

#### 1. PostgreSQL bazasini sozlash
PostgreSQL-ga ulanib, ma'lumotlar bazasini yarating:
```bash
PGPASSWORD=postgres psql -U postgres -h localhost -c "CREATE DATABASE kafe_omborxona;"
```

#### 2. Backend-ni ishga tushirish
```bash
cd backend
go run ./cmd/api/
```
*Backend `http://localhost:8080` portida ishga tushadi, migratsiyalarni avtomatik amalga oshiradi va birinchi marta ishga tushganida `admin` foydalanuvchisini yaratadi.*

#### 3. Frontend-ni ishga tushirish
```bash
cd frontend
npm install
npm run dev
```
*Frontend loyihasi `http://localhost:3000` portida ishga tushadi.*
