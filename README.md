# kbgo

A production-ready Go backend foundation package following **Hexagonal Architecture (Ports & Adapters)**.

Built by Bounkhong — designed so any Go backend project can plug in and go.

---

## Install the package

```bash
go get github.com/bounkhongdev/kbgo
```

## Install the CLI

```bash
go install github.com/bounkhongdev/kbgo/cmd/kbgo@latest
```

---

## CLI Usage

### Scaffold a new project

```bash
kbgo new myapp
kbgo new myapp --module github.com/yourname/myapp
```

Generates:
```
myapp/
├── cmd/api/main.go
├── internal/
├── migrations/
├── .env.example
├── docker-compose.yml
├── Makefile
├── .gitignore
└── go.mod
```

### Generate a full module

```bash
kbgo generate module user
kbgo g module product
kbgo g module orderItem   # supports camelCase / snake_case / kebab-case
```

Generates:
```
internal/user/
├── domain.go       ← entity + repository interface (Port)
├── usecase.go      ← business logic
├── handler.go      ← HTTP handler (Fiber)
└── repository.go   ← PostgreSQL implementation (Adapter)
```

### Generate individual files

```bash
kbgo generate handler    product
kbgo generate service    product
kbgo generate repository product
```

### Remove a module

```bash
kbgo remove module user          # deletes internal/user/ entirely
kbgo rm module user              # same (rm alias)
```

### Remove individual files

```bash
kbgo remove handler    product   # deletes internal/product/handler.go
kbgo remove service    product   # deletes internal/product/usecase.go
kbgo remove repository product   # deletes internal/product/repository.go
```

---

## Package Usage

### Config

```go
import "github.com/bounkhongdev/kbgo/config"

cfg, err := config.Load()          // reads from .env
cfg, err := config.Load(".env.prod") // custom file
```

### Adapters

```go
import (
    "github.com/bounkhongdev/kbgo/adapter/postgres"
    "github.com/bounkhongdev/kbgo/adapter/redis"
    "github.com/bounkhongdev/kbgo/adapter/minio"
    "github.com/bounkhongdev/kbgo/adapter/jwt"
)

db, err    := postgres.New(ctx, cfg.Postgres)
cache, err := redis.New(ctx, cfg.Redis)
store, err := minio.New(ctx, cfg.MinIO)
token      := jwt.New(cfg.JWT)
```

### Contracts (Ports)

Your business logic depends only on these interfaces — never on the adapters directly:

```go
import "github.com/bounkhongdev/kbgo/contract"

type UserRepository interface {
    // uses contract.Database, not *pgxpool.Pool
}

func NewUserUsecase(db contract.Database, cache contract.Cache) *UserUsecase { ... }
```

Swap PostgreSQL for MySQL → zero changes to your business logic.

### Response

```go
import "github.com/bounkhongdev/kbgo/response"

c.JSON(response.Success(data))
c.JSON(response.Paginated(list, page, limit, total))
c.JSON(response.Error("NOT_FOUND", "user not found"))
```

### Errors

```go
import "github.com/bounkhongdev/kbgo/errs"

return errs.NotFound("user not found")
return errs.BadRequest("invalid email")
return errs.Conflict("email already taken")

// In handler:
if ae, ok := errs.IsAppError(err); ok {
    return c.Status(ae.Status).JSON(response.Error(ae.Code, ae.Message))
}
```

### Logger

```go
import "github.com/bounkhongdev/kbgo/logger"

log := logger.Development()   // debug + text output
log := logger.Production()    // info  + JSON output
slog.SetDefault(log)
```

### Middleware

```go
import "github.com/bounkhongdev/kbgo/middleware"

// Global CORS
app.Use(middleware.CORS())
app.Use(middleware.CORS(middleware.CORSConfig{AllowOrigins: "https://myapp.com"}))

// Protect routes with JWT
api := app.Group("/api/v1", middleware.JWT(token))

// Require specific roles
api.Delete("/users/:id",
    middleware.RequireRole("admin"),
    handler.Delete,
)

// Read claims inside a handler
claims := middleware.Claims(c)
userID := claims["user_id"].(string)
```

### Validator

```go
import "github.com/bounkhongdev/kbgo/validator"

type CreateUserInput struct {
    Name  string `validate:"required,min=2"`
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18"`
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
    var input CreateUserInput
    c.BodyParser(&input)

    if errs := validator.Validate(input); errs != nil {
        return c.Status(422).JSON(response.Error("VALIDATION_FAILED", errs))
    }
    // ...
}

// Register a custom rule
validator.RegisterTag("lao_phone", func(fl validator.FieldLevel) bool {
    return strings.HasPrefix(fl.Field().String(), "020")
})
```

### Transactions

```go
// Your repository depends on contract.Transactional instead of contract.Database
type orderRepo struct {
    db contract.Transactional
}

func (r *orderRepo) CreateWithInventory(ctx context.Context, order *Order) error {
    tx, err := r.db.BeginTx(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    if err := tx.Exec(ctx, `INSERT INTO orders ...`, ...); err != nil {
        return err
    }
    if err := tx.Exec(ctx, `UPDATE inventories SET stock_qty = stock_qty - $1 ...`, ...); err != nil {
        return err
    }

    return tx.Commit(ctx)
}
```

### Hash

```go
import "github.com/bounkhongdev/kbgo/hash"

hashed, err := hash.Password("mysecret")
ok          := hash.CheckPassword("mysecret", hashed)
```

### i18n (Multi-language error messages)

Built-in locales: `en`, `lo` (Lao), `th` (Thai), `zh` (Chinese).

```go
import "github.com/bounkhongdev/kbgo/i18n"

// Auto-detect from Accept-Language header
locale := i18n.FromHeader(c.Get("Accept-Language"))

// Translate a standard error code
msg := i18n.Translate(locale, "NOT_FOUND")
// "lo" → "ບໍ່ພົບຂໍ້ມູນ"
// "th" → "ไม่พบข้อมูล"
// "en" → "Resource not found"

// Add a custom language at startup
i18n.Register("fr", map[string]string{
    "NOT_FOUND":   "Ressource introuvable",
    "BAD_REQUEST": "Mauvaise requête",
    // ... your custom app codes too
    "USER_EMAIL_TAKEN": "Cet email est déjà utilisé",
})

// Add app-specific codes to an existing locale
i18n.Register(i18n.LO, map[string]string{
    "USER_NOT_FOUND":    "ບໍ່ພົບຜູ້ໃຊ້",
    "EMAIL_TAKEN":       "ອີເມວນີ້ຖືກໃຊ້ງານແລ້ວ",
    "INVALID_PASSWORD":  "ລະຫັດຜ່ານບໍ່ຖືກຕ້ອງ",
})
```

The generated `handler.go` wires this automatically — every error response reads `Accept-Language` and returns the message in the right language.

---

### Paginate

```go
import "github.com/bounkhongdev/kbgo/paginate"

p := paginate.Params{Page: 1, Limit: 20}
p.Normalize()
offset := p.Offset()  // → 0
```

---

## Architecture

```
Your Project
    │
    ├── internal/user/
    │   ├── domain.go        ← entity + Repository interface  (PORT)
    │   ├── usecase.go       ← business logic (depends on PORT only)
    │   ├── handler.go       ← HTTP delivery
    │   └── repository.go    ← PostgreSQL implementation      (ADAPTER)
    │
    └── cmd/api/main.go      ← wire everything together
            │
            ├── kbgo/config        ← load env
            ├── kbgo/adapter/postgres  ← satisfies contract.Database
            ├── kbgo/adapter/redis     ← satisfies contract.Cache
            ├── kbgo/adapter/minio     ← satisfies contract.Storage
            └── kbgo/adapter/jwt       ← satisfies contract.Token
```

---

## Package Map

| Package | Purpose |
|---|---|
| `contract` | Ports — interfaces for Database, Cache, Storage, Token |
| `config` | Load env vars into typed structs |
| `errs` | AppError type + common HTTP errors |
| `response` | Standard JSON response envelope |
| `logger` | slog-based structured logger factory |
| `hash` | bcrypt password hashing |
| `paginate` | Pagination params + offset helper |
| `i18n` | Multi-language error messages (EN, LO, TH, ZH + custom) |
| `adapter/postgres` | PostgreSQL adapter (pgx v5) + transaction support |
| `adapter/redis` | Redis adapter (go-redis v9) |
| `adapter/minio` | MinIO adapter (minio-go v7) |
| `adapter/jwt` | JWT adapter (golang-jwt v5) |
| `middleware` | Fiber middleware — JWT auth, RBAC, CORS |
| `validator` | Request body validation (go-playground/validator) |
| `mock` | Test doubles — Database, Cache, Storage, Token, Tx |
| `cmd/kbgo` | CLI tool — scaffold & generate |
