# Padrões de Arquitetura Go - go_inventory

## Domínio 1 — Estrutura de Pastas e Arquitetura em Camadas

### Estrutura Identificada

```
go_inventory/
├── main.go                          # Entry point
├── Container/                       # Injeção de dependências (fx)
├── Helpers/                         # Helpers compartilhados
│   └── Errors/errors.go             # AppError customizado
├── SupplyInventory/
│   ├── Domain/
│   │   ├── Entities/                # Entidades do domínio
│   │   │   ├── PalletEntity.go
│   │   │   ├── ProductEntity.go
│   │   │   ├── UserEntity.go
│   │   │   └── ...
│   │   └── contracts/repositories/ # Interfaces de repositório
│   │       ├── Pallet/PalletRepository.go
│   │       ├── User/UserRepository.go
│   │       └── ...
│   ├── Application/
│   │   ├── Controllers/             # Handlers HTTP (Gin)
│   │   ├── Services/               # Lógica de negócio
│   │   ├── Middlewares/            # Middlewares HTTP
│   │   ├── Requests/               # Structs de request
│   │   ├── Responses/              # Structs de response
│   │   └── ApiContracts/           # Contratos de API (Swagger)
│   └── Infrastructure/
│       ├── Db/                      # Conexão e migrations
│       └── repositories/            # Implementações de repositório
│           ├── Pallet/PalletRepository.go
│           └── ...
```

### Padrão Observado

**Arquitetura em Camadas (Clean Architecture simplificado)**
- `Domain/Entities` — Structs que representam o modelo de dados
- `Domain/contracts/repositories` — Interfaces que definem contratos de acesso a dados
- `Application/Services` — Implementação da lógica de negócio (depende de interfaces)
- `Application/Controllers` — Handlers HTTP (depende de Services)
- `Infrastructure/repositories` — Implementação concreta das interfaces (depende de GORM)

**Convenção de nomenclatura**
- Pastas em PascalCase: `PalletEntity`, `PalletRepository`
- Arquivos em PascalCase: `PalletController.go`, `PalletService.go`
- Construtores: `NewXxx(...)` para criar instâncias

---

## Domínio 4 — Injeção de Dependências com fx

O projeto utiliza `go.uber.org/fx` (não dig como mencionado na skill, mas o padrão é similar).

### Container Principal

**Arquivo**: `Container/container.go`

```go
func BuildOptions(db *gorm.DB) fx.Option {
    return fx.Options(
        fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
        fx.Supply(db),

        // Repositories - providers com dependência de db
        fx.Provide(func(db *gorm.DB) contractPallet.PalletRepository {
            return repositoriesPallet.NewPalletRepository(dbadapter.NewGormAdapter(db))
        }),

        // Services - providers com múltiplas dependências
        fx.Provide(func(repo contractPallet.PalletRepository, qrService qrcode.QRCodeService, ...) pallet.PalletService {
            return pallet.NewPalletService(repo, qrService, ...)
        }),

        // Controllers
        fx.Provide(palletController.NewPalletController),

        // Middlewares
        fx.Provide(middlewares.NewAuthMiddleware),
    )
}
```

### Padrão de Construtor

Cada componente segue convenção de construtor que recebe dependências no construtor:

```go
// Service
func NewPalletService(
    palletRepository palletRepository.PalletRepository,
    qrService qrCodeService.QRCodeService,
    palletRackRepository palletRackRepository.PalletRackRepository,
    storage storage.Storage,
    exportService PalletExportService,
) pallet.PalletService {
    return &palletService{...}
}

// Controller
func NewPalletController(service pallet.PalletService) *PalletController {
    return &PalletController{service: service}
}
```

### Registro de Rotas (main.go:127-172)

```go
func registerRoutes(
    router *gin.Engine,
    palletController *palletControllerPkg.PalletController,
    // ... outras dependências injetadas automaticamente
) {
    apiV1 := router.Group("/api/v1")

    palletsGroup := apiV1.Group("/pallets")
    palletsGroup.Use(authMiddleware.Handler())
    palletsGroup.Use(inventoryMiddleware.Handler())
    palletController.Register(palletsGroup)
    // ...
}
```

---

## Domínio 5 — Tratamento de Erros e Respostas HTTP

### AppError Customizado

**Arquivo**: `Helpers/Errors/errors.go`

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

func (e *AppError) ErrorCode() int {
    return e.Code
}

func NewAppError(message string, code ...int) *AppError {
    status := DEFAULT_ERROR_CODE // 500
    if len(code) > 0 {
        status = code[0]
    }
    return &AppError{Code: status, Message: message}
}
```

### Uso nos Controllers

```go
// Retorno de erro com status code
if appErr != nil {
    c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
    return
}

// Retorno de sucesso
c.JSON(http.StatusOK, pallet)
c.JSON(http.StatusCreated, newPallet)
c.JSON(http.StatusNoContent, result)
```

### Formato de Resposta

- **Erro**: `gin.H{"message": "..."}` ou `gin.H{"error": "..."}`
- **Sucesso com dados**: objeto direto (ex: `pallet`)
- **Sucesso com wrapper**: `gin.H{"data": gin.H{"url": url}}`

### Validação de Request

```go
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
    return
}
```

---

## Domínio 6 — GORM: Repositórios e Migrations

### DBAdapter Interface

**Arquivo**: `SupplyInventory/Infrastructure/repositories/db/adapter.go`

Padrão de abstração que permite troca de implementação (útil para testes).

```go
type DBAdapter interface {
    Create(value interface{}) error
    FirstByID(dest interface{}, id uint) error
    FindAll(dest interface{}) error
    DeleteByID(dest interface{}, id uint) (int64, error)
    Save(value interface{}) error
    GetDB() *gorm.DB
    // ... outros métodos
}

func NewGormAdapter(db *gorm.DB) DBAdapter {
    return &gormAdapter{db: db}
}
```

### Implementação de Repositório

**Arquivo**: `SupplyInventory/Infrastructure/repositories/Pallet/PalletRepository.go`

```go
type palletRepository struct {
    db dbadapter.DBAdapter
}

func NewPalletRepository(db dbadapter.DBAdapter) repositories.PalletRepository {
    return &palletRepository{db: db}
}

// Métodos usando dbadapter
func (repository *palletRepository) GetSupplyById(id uint) (*entities.PalletEntity, *errors.AppError) {
    var pallet entities.PalletEntity
    if err := repository.db.PreloadFind(&pallet, "PalletizedProduct", id); err != nil {
        return nil, errors.NewAppError(err.Error(), 500)
    }
    return &pallet, nil
}
```

### Migrations

**Arquivo**: `SupplyInventory/Infrastructure/Db/Migrate.go`

```go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &entities.PalletEntity{},
        &entities.ProductEntity{},
        &entities.UserEntity{},
        &entities.InventoryEntity{},
        &entities.PalletizedProductEntity{},
        &entities.PalletRackEntity{},
    )
}
```

### Entidades GORM

**Arquivo**: `SupplyInventory/Domain/Entities/PalletEntity.go`

```go
type PalletEntity struct {
    ID                uint                      `gorm:"primaryKey" json:"id"`
    InventoryID       uint                      `gorm:"index;not null" json:"inventory_id"`
    Name              string                    `gorm:"unique;not null" json:"name"`
    PalletizedProduct []PalletizedProductEntity `gorm:"constraint:OnDelete:CASCADE;foreignKey:PalletID"`
    PalletRackID      uint                      `gorm:"not null" json:"palletRackId"`
    PalletRackName    string                    `gorm:"not null" json:"palletRackName"`
    QrCode            string                    `json:"qr_code"`
    QrCodeUrl         string                    `json:"qr_code_url"`
    CreatedAt         time.Time                 `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt         time.Time                 `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PalletEntity) TableName() string {
    return "pallets"
}
```

---

## Resumo dos Padrões

| Domínio | Padrão | Arquivo de Referência |
|---------|--------|----------------------|
| Arquitetura | Clean Architecture (Domain/Application/Infrastructure) | Estrutura de pastas |
| DI | fx com providers e invoke | `Container/container.go` |
| Controllers | Gin handlers com método Register | `Application/Controllers/Pallet/PalletController.go` |
| Erros | AppError com código HTTP | `Helpers/Errors/errors.go` |
| Repositórios | Interface no Domain, implementação no Infrastructure | `Domain/contracts/repositories/Pallet/PalletRepository.go` |
| GORM | DBAdapter wrapper + AutoMigrate | `Infrastructure/repositories/db/adapter.go` |
| Testes | FakeDBAdapter + IntegrationTestHelper | `documents/Testing/TestingStandards.md` |