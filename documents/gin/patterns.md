# Padrões Gin (Routes & Controllers) - go_inventory

## Visão Geral

O projeto usa **Gin** como framework HTTP com padrões consistentes para controllers, middlewares e rotas.

## Estrutura de Controller

### Padrão Básico

Cada controller segue o padrão:

```go
type PalletController struct {
    service pallet.PalletService
}

func NewPalletController(service pallet.PalletService) *PalletController {
    return &PalletController{service: service}
}
```

**Construtor**: Recebe dependências (services) via injeção de dependência (fx).

### Registro de Rotas

Método `Register` que recebe um `*gin.RouterGroup` para definir as rotas:

```go
func (controller *PalletController) Register(group *gin.RouterGroup) {
    group.GET("", controller.ListPallets)
    group.GET("/export", controller.ExportPalletsCsv)
    group.GET("/:id", controller.FindPalletById)
    group.PATCH("/:id", controller.UpdatePallet)
    group.POST("", controller.CreatePallet)
    group.DELETE("/:id", controller.DeletePalletById)
}
```

**Convenção**: O grupo já inclui o path base (ex: `/pallets`), então cada rota adiciona sua extensão.

### Handler Functions

Assinatura padrão:

```go
func (controller *PalletController) ListPallets(c *gin.Context) {
    // Lógica do handler
}
```

## Leitura de Dados da Request

### Binding JSON

```go
var req palletRequests.PalletRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
    return
}
```

### Query Params

```go
func (controller *PalletController) parseFilterParams(c *gin.Context) (*uint, *int) {
    var palletRackId *uint
    var productEan *int

    if rackIdStr := c.Query("palletRackId"); rackIdStr != "" {
        if rackId, err := strconv.ParseUint(rackIdStr, 10, 32); err == nil {
            rackIdUint := uint(rackId)
            palletRackId = &rackIdUint
        }
    }
    // ...
    return palletRackId, productEan
}
```

### Path Params

```go
id, err := requestsHelper.GetIDParam(c, "id")
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
    return
}
```

### Headers

```go
authHeader := c.GetHeader("Authorization")
inventoryID := middlewares.GetInventoryID(c)
```

## Retorno de Respostas

### Sucesso

```go
// Retorno direto do objeto
c.JSON(http.StatusOK, pallets)

// Retorno com wrapper
c.JSON(http.StatusOK, gin.H{"data": gin.H{"url": url}})

// Retorno Created
c.JSON(http.StatusCreated, newPallet)

// Retorno sem conteúdo
c.JSON(http.StatusNoContent, result)
```

### Erro

```go
if appErr != nil {
    c.JSON(appErr.ErrorCode(), gin.H{"message": appErr.Error()})
    return
}
```

### Helper de Validação

```go
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusUnprocessableEntity, gin.H{"error": requestsHelper.FormatValidationErrors(err)})
    return
}
```

## Middlewares

### Auth Middleware

**Arquivo**: `SupplyInventory/Application/Middlewares/auth.go`

```go
type AuthMiddleware struct {
    JWTService jwt.JWTService
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        // Valida token Bearer
        // Define c.Set("userID", ...)
        c.Next()
    }
}
```

**Uso nas rotas**:
```go
palletsGroup.Use(authMiddleware.Handler())
```

### Inventory Middleware

Responsável por extrair o `inventory_id` do header `x-inventory-id` e disponibilizar no contexto.

### Middleware de CORS

Configurado no `main.go`:

```go
router.Use(cors.New(cors.Config{
    AllowOriginFunc: func(origin string) bool { return true },
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-inventory-id"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

## Agrupamento de Rotas

### Estrutura em main.go

```go
func registerRoutes(
    router *gin.Engine,
    palletController *palletControllerPkg.PalletController,
    // ... outras dependências
) {
    apiV1 := router.Group("/api/v1")

    // pallets
    palletsGroup := apiV1.Group("/pallets")
    palletsGroup.Use(authMiddleware.Handler())
    palletsGroup.Use(inventoryMiddleware.Handler())
    palletController.Register(palletsGroup)

    // auth (sem middleware)
    authController.RegisterLogin(apiV1.Group("/auth"))

    // users
    userController.RegisterUserRoutes(apiV1.Group("/users"))
}
```

### Exemplo de rota específica (Admin)

```go
adminRacksGroup := apiV1.Group("/admin/racks")
adminRacksGroup.Use(authMiddleware.Handler())
adminPalletRackController.RegisterAdminPalletRack(adminRacksGroup)
```

## Swagger

Anotações nos handlers para documentação automática:

```go
// @Summary List pallets
// @Tags Pallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.PalletEntity
// @Failure 404 "Not Found"
// @Param palletRackId query uint false "Filter by pallet rack ID"
// @Router /api/v1/pallets [get]
func (controller *PalletController) ListPallets(c *gin.Context) {
    // ...
}
```

## Helpers

### FormatValidationErrors

Disponível em `Helpers/RequestsHelper` para formatar erros de validação do Gin.

### GetIDParam

Helper para extrair e validar ID da URL.

---

## Arquivos de Referência

- `main.go` — Setup do router e registro de rotas
- `Application/Controllers/Pallet/PalletController.go` — Exemplo completo de controller
- `Application/Controllers/Auth/AuthController.go` — Controller sem middleware
- `Application/Middlewares/auth.go` — Middleware de autenticação
- `Application/Middlewares/inventory.go` — Middleware de inventário