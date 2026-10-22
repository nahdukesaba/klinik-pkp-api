## **Clean Code Best Practices**

#### ✅ ImageService.go
- Single responsibility: setiap function 1 tugas
- Short functions (< 30 lines)
- Clear naming: `ExtractImageDimensions`, `UpdateImageDimensions`
- Error handling yang konsisten
- No magic numbers (pakai config)

#### ✅ ImageController.go
- Thin controller, logic di service
- Consistent response format
- Proper HTTP status codes
- Rollback pattern saat error

#### ✅ RusunService.go
- XSS sanitization
- SQL injection prevention (GORM parameterized queries)
- payload validation
- Foreign key validation

#### ✅ Seed Files
- Idempotent (bisa dijalankan berkali-kali)
- Check existence sebelum insert
- Clear logging
- Foreign key dependencies handled

## **API Endpoints**

Base URL: `http://localhost:3000/api/v1`

### Central Browser SSO

Browser SSO is an additional authorization-code + PKCE login path. It does not
remove or change the existing `/authentications` login, refresh, and logout
routes.

- `GET /api/v1/auth/sso/start` creates a short-lived flow and redirects to the
  trusted Central `/authorize` endpoint.
- `POST /api/v1/auth/sso/exchange` accepts `code`, `state`, and the optional
  safe Central `status`. The backend atomically consumes the hashed state,
  derives the PKCE verifier, exchanges the code server-to-server, validates
  the `klinik_pkp` context, and resolves `klinikpkp_user` by `app_user_id`.
- `GET /api/v1/auth/sso/logout` redirects to Central `/logout` using only the
  configured client ID and post-logout URI.

Local configuration:

```env
SSO_BROWSER_ENABLED=true
SSO_BASE_URL=http://127.0.0.1:18082
SSO_CLIENT_ID=klinik-pkp-web
SSO_REDIRECT_URI=http://localhost:3000/auth/sso/callback
SSO_POST_LOGOUT_REDIRECT_URI=http://localhost:3000/login
SSO_FLOW_TTL=5m
SSO_FLOW_HMAC_KEY=
```

`SSO_FLOW_HMAC_KEY` must be a base64-encoded random key of at least 32 bytes
and must not be committed. Production URLs are deployment-specific and are
not finalized. Migration `20261023_create_klinik_sso_flow` creates
`public.klinikpkp_sso_flow`; run it only against the intended environment.

### Klinik authorization and logout

Klinik authorization is sourced from the shared Central SSO tables. Every
protected request, including requests using a legacy local HS256 JWT, resolves
the linked `app_user_id` and requires an approved `app_memberships` row for
`klinik_pkp`. The only accepted application roles are `user` and `admin`.
`klinikpkp_user.role` and roles embedded in legacy JWTs are retained only for
compatibility and are not authorization sources.

The legacy `/authentications` login and refresh endpoints remain compatibility
paths. Both revalidate the current Central membership before issuing tokens,
and issued roles come from `app_memberships`. Consequently, an existing token
is denied on its next protected request after its membership is rejected or
disabled.

Logout returns `data.centralLogoutRequired`. Direct password and legacy local
sessions return `false`; a verified Browser SSO session returns `true`. The
value is derived server-side from the verified token's authentication method.
The existing `GET /api/v1/auth/sso/logout` bridge remains the trusted Central
logout target when that value is `true`.

#### Rusun:
- `GET /rusun` - List all rusun
- `GET /rusun/:id` - Get rusun detail
- `GET /rusun/map` - Get map data
- `GET /rusun/statistics` - Get statistics
- `POST /rusun` - Create rusun (admin only)
- `PUT /rusun/:id` - Update rusun (admin only)
- `DELETE /rusun/:id` - Delete rusun (admin only)

## **Security Best Practices**

✅ **SQL Injection Prevention:**
- Semua query menggunakan GORM parameterized queries
- Tidak ada string interpolation manual

✅ **XSS Prevention:**
- HTML escape semua input
- Remove script tags
- Remove event handlers (onclick, dll)

✅ **File Upload Security:**
- Whitelist MIME types
- File size limitation
- Unique filename generation
- No executable file upload

✅ **Payload Validation:**
- Sanitize semua string input
- Validate foreign keys
- Validate ranges (lat/long, units, years)
- Validate enum values (status)

## 📝 Cara Menggunakan

### 1. Run Migrations & Seeder
```bash
go run cmd/migrate/main.go --help
go run cmd/seed/main.go --help 
go run cmd/migrate/main.go villages
go run cmd/seed/main.g villages
```

### 2. Run Server
```bash
go run cmd/server/main.go
# atau
go build -o klinik-pkp-api.exe .
./klinik-pkp-api.exe
```

## ✨ Clean Code Principles Applied

1. **Single Responsibility Principle**
   - Setiap function melakukan 1 task
   - Service untuk business logic
   - Controller hanya handle HTTP
   - Utils untuk helper functions

2. **DRY (Don't Repeat Yourself)**
   - Reusable validation functions
   - Shared sanitization logic
   - Common response helpers

3. **KISS (Keep It Simple, Stupid)**
   - Short functions
   - Clear naming
   - No over-engineering

4. **Error Handling**
   - Consistent error messages
   - Proper error propagation
   - Rollback on failure

5. **Security First**
   - payload validation
   - XSS prevention
   - SQL injection prevention
   - File upload security

6. **Maintainability**
   - Clear comments
   - Logical file structure
   - Easy to test
   - Easy to extend

## 🎯 Testing Checklist

- [x] Build berhasil tanpa error
- [ ] Migrations berhasil
- [ ] Seeder berhasil (8 rusun, 2 users)
- [ ] Upload single image
- [ ] Upload multiple images
- [ ] Image width/height extracted
- [ ] GET /rusun menampilkan data lengkap dengan district & village
- [ ] Login admin berhasil
- [ ] Login user berhasil

## 🔧 Next Steps (Optional)

1. Image compression sebelum simpan
2. Thumbnail generation
3. Cloud storage integration (S3, GCS)
4. Image watermark
5. Rate limiting untuk upload
6. Virus scanning untuk file upload
