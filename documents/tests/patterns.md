# Padrões de Testes - go_inventory

## Visão Geral

O projeto mantém dois tipos de testes:
- **Unitários**: rápidos, sem dependência externa (usa `FakeDBAdapter`)
- **Integração**: com banco de dados real (MySQL)

## Estrutura de Testes

### Unit Tests
- **Localização**: mesmo diretório do código testado
- **Nomeação**: `{Nome}_test.go`
- **Execução**: `go test ./...`

### Integration Tests  
- **Localização**: mesmo diretório do código testado
- **Nomeação**: `{Nome}_integration_test.go`
- **Build Tags**: `//go:build integration`
- **Execução**: `go test -tags integration ./...`

## Padrões de Testes Unitários

### FakeDBAdapter

**Arquivo**: `SupplyInventory/tests/testutils/fake_db_adapter.go`

Implementa a interface `DBAdapter` com funções hook para controle em testes:

```go
type FakeDBAdapter struct {
    CreateFn            func(value interface{}) error
    FirstByIDFn         func(out interface{}, id uint) error
    WhereFirstFn        func(out interface{}, query string, args ...interface{}) error
    FindAllFn           func(out interface{}) error
    DeleteByIDFn        func(model interface{}, id uint) (int64, error)
    SaveFn              func(value interface{}) error
    PreloadFindFn       func(out interface{}, preload string, id ...uint) error
    // ...
}
```

### Exemplo de Uso em Repositório

```go
func TestPalletRepository_FindByID(t *testing.T) {
    fake := &testutils.FakeDBAdapter{}
    
    fake.FirstByIDFn = func(out interface{}, id uint) error {
        pallet := out.(*entities.PalletEntity)
        pallet.ID = id
        pallet.Name = "Test Pallet"
        return nil
    }
    
    repo := repositories.NewPalletRepository(fake)
    got, err := repo.FindByID(1)
    
    assert.NoError(t, err)
    assert.Equal(t, "Test Pallet", got.Name)
}
```

### Exemplo de Uso em Service

```go
func TestPalletService_CreatePallet(t *testing.T) {
    // Setup com mocks
    mockRepo := &mocks.PalletRepository{}
    mockPalletRackRepo := &mocks.PalletRackRepository{}
    mockQrService := &mocks.QRCodeService{}
    
    mockPalletRackRepo.On("FindPalletByID", uint(1)).Return(&entities.PalletRackEntity{
        ID:          1,
        InventoryID: 1,
        Name:        "Rack 1",
    }, nil)
    
    mockRepo.On("AddSupply", "New Pallet", uint(1)).Return(&entities.PalletEntity{
        ID:   1,
        Name: "New Pallet",
    }, nil)
    
    service := pallet.NewPalletService(mockRepo, mockQrService, mockPalletRackRepo, nil, nil)
    
    // Teste
    result, err := service.CreatePallet("New Pallet", 1, 1)
    
    assert.NoError(t, err)
    assert.Equal(t, "New Pallet", result.Name)
}
```

## Padrões de Testes de Integração

### IntegrationTestHelper

**Arquivo**: `SupplyInventory/tests/integration/helpers.go`

Fornece setup consistente para testes de API:

```go
type IntegrationTestHelper struct {
    DB      *gorm.DB
    Router  *gin.Engine
    // ...
}

func NewIntegrationTestHelper() *IntegrationTestHelper {
    // conecta ao banco de testes
}

func (h *IntegrationTestHelper) TruncateTables(db *gorm.DB) {
    // limpa todas as tabelas
}

func (h *IntegrationTestHelper) SetupRouterForPallet(db *gorm.DB) *gin.Engine {
    // cria router com dependências injetadas usando db de teste
}
```

### Exemplo de Teste de Integração

```go
//go:build integration

func TestIntegration_Pallet_Create(t *testing.T) {
    h := integration.NewIntegrationTestHelper()
    h.TruncateTables(h.DB)
    
    h.DB.Transaction(func(tx *gorm.DB) error {
        r := h.SetupRouterForPallet(tx)
        
        // Criar request
        body := `{"name": "Test Pallet", "palletRackId": 1}`
        req, _ := http.NewRequest("POST", "/api/v1/pallets", bytes.NewBufferString(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+testToken)
        
        // Executar
        w := httptest.NewRecorder()
        r.ServeHTTP(w, req)
        
        // Assertions
        assert.Equal(t, http.StatusCreated, w.Code)
        return nil
    })
}
```

## Mocks

**Localização**: `SupplyInventory/tests/mocks/`

O projeto gera mocks automaticamente para interfaces usando `mockery` ou manualmente:

```go
// SupplyInventory/tests/mocks/UserRepository.go
type UserRepositoryMock struct {
    mock.Mock
}

func (m *UserRepositoryMock) FindByEmail(email string) (*entities.UserEntity, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.UserEntity), args.Error(1)
}
```

## Comandos de Execução

```bash
# Unit tests
go test ./... -v

# Integration tests (requer MySQL)
go test -tags integration ./... -v

# Com coverage
go test ./... -cover
```

## Boas Práticas

1. **Nomes descritivos**: `Test{Nome}_{Cenario}`
2. **Um comportamento por teste**: evite testar múltiplas coisas
3. **Isolamento**: cada teste deve ser independente
4. **Setup/Teardown**: use funções auxiliares quando possível
5. **Tabelas de testes**: para múltiplos cenários similares

---

## Referências

- `documents/Testing/TestingStandards.md` — Padrões detalhados do projeto
- `SupplyInventory/tests/testutils/fake_db_adapter.go` — Implementação do FakeDBAdapter
- `SupplyInventory/tests/integration/helpers.go` — Helper de integração