# Features do Projeto go_inventory

## Visão Geral
O **go_inventory** é um sistema de gestão de armazém focado em rastreabilidade via Pallets e QR Codes. O sistema suporta múltiplos inventários isolados e gestão física de racks.

---

## 1. Autenticação e Segurança
Implementa autenticação JWT (JSON Web Token) para proteger os recursos.

### Fluxo de Login
- **Endpoint**: `POST /api/v1/auth/login`
- **Payload**: `{"email": "admin@example.com", "password": "..."}`
- **Uso**: Retorna um `token` e um `refreshToken`. Todas as rotas subsequentes (exceto login) exigem o header `Authorization: Bearer {token}`.

---

## 2. Gestão de Inventários
Permite a existência de múltiplos depósitos/clientes no mesmo sistema.

- **Isolamento**: O sistema exige o header `x-inventory-id` em rotas de operação para garantir que você está manipulando dados do inventário correto.
- **Endpoints**: CRUD completo em `/api/v1/inventories`.

---

## 3. Gestão de Pallet Racks (Estruturas de Armazenagem)
Define onde os pallets podem ser colocados fisicamente.

### Operações de Usuário
- **Listagem**: `GET /api/v1/racks` (Suporta paginação e filtros).
- **Criação**: `POST /api/v1/racks` - Define nome, localização e capacidade.

### Operações Administrativas
- **Endpoint**: `/api/v1/admin/racks`
- Permite controle total sobre as estruturas, incluindo ajustes de capacidade e remoção forçada.

---

## 4. Gestão de Pallets e QR Codes
O núcleo do sistema. Cada pallet é uma unidade rastreável.

### Ciclo de Vida do Pallet
1. **Criação**: `POST /api/v1/pallets` associando-o a um Rack.
2. **Identificação**: O sistema gera automaticamente um **QR Code** único para o pallet.
3. **Localização**: O pallet fica vinculado a um `PalletRack`.
4. **Consulta**: `GET /api/v1/pallets/:id` retorna os detalhes e a URL do QR Code.

### Filtros e Exportação
- **Busca por EAN**: `GET /api/v1/pallets?ean=123456` busca pallets que contenham um produto específico.
- **Busca por Rack**: `GET /api/v1/pallets?palletRackId=5`.
- **Exportação CSV**: `GET /api/v1/pallets/export` gera um relatório dos pallets filtrados.

---

## 5. Palletized Products (Conteúdo dos Pallets)
Gerencia o que está dentro de cada pallet.

- **Adição**: `POST /api/v1/pallet/products` - Adiciona uma quantidade de um produto (EAN) com data de validade.
- **Remoção**: `DELETE /api/v1/pallet/products/:id`.
- **Validade**: Permite rastrear produtos por data de expiração dentro dos pallets.

---

## 6. Catálogo de Produtos
Cadastro base de produtos.

- **EAN Único**: Cada produto é identificado pelo seu código de barras global.
- **Endpoints**: CRUD em `/api/v1/products`.

---

## Como Utilizar (Exemplo Prático)

### Passo 1: Login
```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -d '{"email":"admin@inventory.com", "password":"admin"}'
```

### Passo 2: Criar um Pallet (Exige Token e Inventory ID)
```bash
curl -X POST http://localhost:8000/api/v1/pallets \
  -H "Authorization: Bearer {TOKEN}" \
  -H "x-inventory-id: 1" \
  -d '{"name": "PALLET-A1", "palletRackId": 10}'
```

### Passo 3: Adicionar Produto ao Pallet
```bash
curl -X POST http://localhost:8000/api/v1/pallet/products \
  -H "Authorization: Bearer {TOKEN}" \
  -H "x-inventory-id: 1" \
  -d '{"palletId": 1, "productId": 5, "quantity": 100, "expirationDate": "2026-12-31"}'
```
