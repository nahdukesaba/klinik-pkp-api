# Klinik PKP API

Backend Go/Fiber untuk aplikasi Klinik PKP. API menyediakan data wilayah dan
program perumahan, pengelolaan pengguna Klinik, upload berkas, autentikasi
kompatibilitas, serta Browser SSO melalui Central SSO.

## Menjalankan aplikasi

1. Salin `.env.example` menjadi `.env` dan isi konfigurasi database, JWT, CORS,
   dan Central SSO yang sesuai dengan lingkungan.
2. Jalankan migrasi yang diperlukan dengan perintah di bawah.
3. Jalankan API dari root repository.

```bash
go run cmd/migrate/main.go --help
go run cmd/seed/main.go --help
go run .
```

Secara default API mendengarkan `APP_PORT=8000`. Endpoint pemeriksaan tersedia
pada `GET /api/v1/health`.

## Model autentikasi

Browser SSO adalah jalur authorization-code + PKCE yang dimiliki backend. Alur
aktifnya adalah:

```text
Klinik frontend -> GET /api/v1/auth/sso/start
-> Central /authorize
-> Klinik frontend callback
-> POST /api/v1/auth/sso/exchange
```

Backend menyimpan hash state dan verifier PKCE dalam flow berumur pendek,
menukar code ke Central secara server-to-server, kemudian memvalidasi konteks
`klinik_pkp`, status membership, peran, serta identity link sebelum membuat
sesi aplikasi.

Konfigurasi Browser SSO:

```env
SSO_BROWSER_ENABLED=true
SSO_BASE_URL=http://127.0.0.1:18082
SSO_CLIENT_ID=klinik-pkp-web
SSO_REDIRECT_URI=http://localhost:3000/auth/sso/callback
SSO_POST_LOGOUT_REDIRECT_URI=http://localhost:3000/login
SSO_FLOW_TTL=5m
SSO_FLOW_HMAC_KEY=
SSO_INTROSPECTION_SECRET=
```

`SSO_FLOW_HMAC_KEY` harus berupa base64 dari sedikitnya 32 byte acak.
`SSO_INTROSPECTION_SECRET` harus sama dengan secret Klinik pada Central dan
panjangnya sedikitnya 32 karakter. Jangan menyimpan nilai produksi dalam Git.

## Jalur kompatibilitas yang masih aktif

Route `POST`, `PUT`, dan `DELETE /api/v1/authentications` masih digunakan oleh
Klinik frontend/BFF untuk login langsung, refresh, dan logout. JWT HS256 lokal,
tabel authentication, serta `ACCESS_TOKEN_KEY` dan `REFRESH_TOKEN_KEY` karena
itu belum boleh dihapus.

Setiap request terlindungi tetap memuat ulang membership Central. Role dari
JWT atau tabel pengguna lokal bukan sumber otorisasi akhir. Role aplikasi yang
diterima hanya `user` dan `admin` dari membership `klinik_pkp` yang approved.

Logout mengembalikan `data.centralLogoutRequired`. Nilai tersebut dihitung
dari konteks autentikasi yang telah diverifikasi, bukan dari input browser.
Jika bernilai `true`, frontend melanjutkan ke bridge
`GET /api/v1/auth/sso/logout`.

## Endpoint utama

- `GET /api/v1/` — informasi API.
- `GET /api/v1/health` — health check proses dan koneksi awal database.
- `/api/v1/authentications` — login, refresh, dan logout kompatibilitas.
- `/api/v1/auth/sso/*` — start, exchange, dan logout Browser SSO.
- `/api/v1/users` — direktori pengguna Klinik.
- `/api/v1/provinces`, `/regions`, `/districts`, `/villages` — data wilayah.
- `/api/v1/bsps`, `/kumuh`, `/rusun`, `/bank-desain`, `/sosialisasi`, `/faq`
  — data domain Klinik.
- `/api/v1/uploads` — upload gambar dan PDF tervalidasi.

Detail method, parameter, dan pembatasan role mengikuti route pada masing-masing
package di `internal/api/v1`.

## Migrasi dan seeder

CLI migrasi mendukung migrasi model, modifikasi bernama, rollback, drop, dan
truncate. Lihat daftar dan sintaks terkini melalui:

```bash
go run cmd/migrate/main.go --help
```

Migration `20261023_create_klinik_sso_flow` membuat tabel flow Browser SSO.
Jalankan migrasi hanya pada database lingkungan yang memang dituju. Migration
lama tetap dipertahankan sebagai histori skema dan tidak boleh dihapus hanya
karena model aplikasinya berubah.

Seeder dijalankan secara eksplisit, bukan otomatis dari `DB_SEED`:

```bash
go run cmd/seed/main.go --help
```

## Verifikasi perubahan

```bash
go test ./...
go build ./...
go vet ./...
```

Pengujian otomatis mencakup konfigurasi Browser SSO, PKCE/state replay,
otorisasi membership, middleware JWT/SSO, kontrak logout, serta direktori
pengguna canonical.
