**Storage Setup**

Este documento descreve como configurar e usar os providers de storage implementados no projeto: Local (filesystem) e MinIO (S3-compatible). Inclui variáveis de ambiente, exemplos de `docker compose` e instruções para alternar entre providers.

**Visão Geral**

- **Local**: arquivos são salvos em disco no diretório `storage/` (por padrão). Útil para desenvolvimento local sem custo extra.
- **MinIO (S3)**: backend compatível com S3. Pode ser executado localmente via Docker (já incluído no `docker-compose.yml`) ou apontado para um serviço S3 real.

As decisões de qual provider usar são feitas por variáveis de ambiente; o bootstrap do projeto (provider do `QRCodeService`) lê essas variáveis e instancia o storage apropriado.

**Arquivos relevantes**

- Código da fábrica/boot: [Container/container.go](Container/container.go#L1-L240)
- Serviço de QR code: [SupplyInventory/Application/Services/QrCode/QrCode.go](SupplyInventory/Application/Services/QrCode/QrCode.go#L1-L200)
- Adaptador Local: [SupplyInventory/Application/Services/Storage/local.go](SupplyInventory/Application/Services/Storage/local.go#L1-L200)
- Adaptador S3/MinIO: [SupplyInventory/Application/Services/Storage/s3.go](SupplyInventory/Application/Services/Storage/s3.go#L1-L200)

Leia os arquivos acima para detalhes de implementação; abaixo estão instruções práticas para uso.

**Variáveis de ambiente (gerais)**

- `STORAGE_PROVIDER` — define o provider: `minio` ou `local` (padrão `local` se não informado).

1) Usando Local (filesystem)

Variáveis opcionais:
- `STORAGE_BASE_DIR` — diretório onde os arquivos serão salvos. Default: `storage`.
- `STORAGE_BASE_URL` — URL pública base para construir links (ex.: `http://localhost:3000/`). Default: `http://localhost:3000/`.

Exemplo de uso local (bash):

```bash
export STORAGE_PROVIDER=local
export STORAGE_BASE_DIR=storage
export STORAGE_BASE_URL=http://localhost:3000/
# iniciar a aplicação normalmente (ex.: via docker compose ou local)
docker compose up --build app
```

Com `LocalStorage` os arquivos são salvos em `storage/<relative_path>` e o serviço de QR code retorna tanto o caminho no disco quanto a URL pública montada por `GetURL`.

2) Usando MinIO (S3-compatible)

Variáveis necessárias:
- `STORAGE_PROVIDER=minio`
- `MINIO_ENDPOINT` — ex.: `localhost:9000` (sem `http://`).
- `MINIO_ACCESS_KEY` — chave de acesso.
- `MINIO_SECRET_KEY` — segredo.
- `MINIO_BUCKET` — nome do bucket a usar (o adaptador tenta criar o bucket se não existir).
- `MINIO_REGION` — opcional.
- `MINIO_USE_SSL` — `true|false` (opcional)

Exemplo local com o MinIO fornecido pelo `docker-compose.yml` (já presente no repositório):

1) Subir MinIO via docker compose (é criado um volume `minio_data` automaticamente):

```bash
docker compose up -d minio
```

2) Exportar variáveis e rodar app:

```bash
export STORAGE_PROVIDER=minio
export MINIO_ENDPOINT=localhost:9000
export MINIO_ACCESS_KEY=minioadmin
export MINIO_SECRET_KEY=minioadmin
export MINIO_BUCKET=inventory
export MINIO_REGION=us-east-1
export MINIO_USE_SSL=false

# iniciar sua aplicação (ex.: via docker compose ou local)
docker compose up --build app
```

Observações:
- O provider S3 criado em `Container/container.go` usa `NewS3Storage(...)` que verifica se o bucket existe e tenta criá-lo caso não exista.
- `MINIO_ENDPOINT` no exemplo é `localhost:9000` porque o MiniO do docker é exposto nessa porta por padrão.

**Comportamento de fallback**

Se as variáveis de configuração necessárias para MinIO não estiverem presentes ou ocorrer erro ao criar o cliente, o bootstrap faz fallback para `LocalStorage`. Isso facilita desenvolvimento sem dependências externas.

**Como alternar entre providers**

- Para usar Local: `export STORAGE_PROVIDER=local` (ou não definir `STORAGE_PROVIDER`).
- Para usar MinIO: defina `STORAGE_PROVIDER=minio` e as variáveis `MINIO_*` listadas.

**Testes / Verificações**

- Testes unitários: rode `go test ./...` para validar que tudo compila e os testes unitários passam.
- Teste de MinIO local: depois de subir MinIO, ao iniciar a aplicação com `STORAGE_PROVIDER=minio` o adaptador tentará criar o bucket e você poderá ver os objetos via WebUI: `http://localhost:9001` (credenciais padrão no `docker-compose.yml`: `minioadmin` / `minioadmin`).

**Segurança e produção**

- Não use as credenciais padrão `minioadmin`/`minioadmin` em produção. Configure segredos via gestor de segredos ou variáveis de ambiente seguras.
- Para S3 real (AWS), prefira usar roles/credentials gerenciadas (IAM) em vez de chaves embutidas.

**Links rápidos**

- Arquivo de bootstrap que faz a decisão por env vars: [Container/container.go](Container/container.go#L1-L240)
- Serviço de QR code que grava arquivos via storage: [SupplyInventory/Application/Services/QrCode/QrCode.go](SupplyInventory/Application/Services/QrCode/QrCode.go#L1-L200)
- Local adapter: [SupplyInventory/Application/Services/Storage/local.go](SupplyInventory/Application/Services/Storage/local.go#L1-L200)
- MinIO (S3) adapter: [SupplyInventory/Application/Services/Storage/s3.go](SupplyInventory/Application/Services/Storage/s3.go#L1-L200)

Se quiser, eu adiciono um `README` curto na raiz com esses passos ou crio scripts (`scripts/start-minio.sh`) para automatizar as exportações de variáveis e o `docker compose up`.
