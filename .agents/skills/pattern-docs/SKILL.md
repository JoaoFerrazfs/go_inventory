---
name: go-pattern-extractor
description: >
  Analisa projetos Go com Gin e extrai/documenta os padrões arquiteturais usados, gerando um documento Markdown completo de referência. Use esta skill sempre que o usuário quiser: documentar padrões de um projeto Go, entender como o projeto está organizado, verificar padrões de arquitetura em camadas (Clean Architecture), analisar como rotas e controllers Gin estão estruturados, mapear padrões de testes com build tags (unit/integration/e2e), documentar uso de injeção de dependências com dig, tratamento de erros HTTP, ou padrões de repositório com GORM. Acione esta skill mesmo que o usuário diga apenas "documenta os padrões do meu projeto Go", "como meu projeto está organizado", "gera um guia de padrões", ou faça upload de arquivos .go para análise.
---

# Go Pattern Extractor

Analisa código Go (Gin + Clean Architecture) e gera um documento Markdown documentando os padrões encontrados no projeto.
O documento deve ser criado na pasta `documents` dividido por contextos, `documents\tests\patterns.md` para padrões de testes, `documents\gin\patterns.md` para padrões de rotas/controllers, em casos de documentações gerais ele poderia ficar em um `documents\general\patterns.md`.

## Fluxo de Execução

### 1. Coleta de Código

O código esta disponivel no projeto e será analisado para extrair os padrões.
Você deve analisar os arquivos disponiveis no `\documents` e entender se o que você esta documentado ja existe para nesse caso atualizar ou se é um novo padrão a ser documentado.

### 2. Geração do Documento

**Regras de geração:**
- Use exemplos de código reais extraídos do projeto, nunca inventados
- Se um padrão não for encontrado no código fornecido, marque como `⚠️ Não identificado nos arquivos analisados`
- Inclua o nome dos arquivos de onde cada padrão foi extraído
- Mantenha linguagem técnica mas direta

---

### 3. Domínios para analise

### Domínio 1 — Estrutura de Pastas e Arquitetura em Camadas
Identifique como o projeto está dividido em camadas. Procure por:
- Separação entre `domain/`, `infra/`, `application/` ou equivalentes
- Onde ficam entidades, repositórios, serviços e controllers
- Convenção de nomenclatura de pastas

### Domínio 2 — Rotas e Controllers Gin
Identifique como as rotas e handlers HTTP estão organizados. Procure por:
- Como o router Gin é configurado e onde as rotas são registradas
- Assinatura dos handler functions (`func(c *gin.Context)`)
- Agrupamento de rotas (`router.Group(...)`)
- Uso de middlewares nas rotas
- Como os dados da request são lidos (binding, params, query)
- Como as responses são retornadas (JSON, status codes)

### Domínio 3 — Testes (build tags unit/integration/e2e)
Identifique como os testes estão organizados. Procure por:
- Uso de `//go:build unit`, `//go:build integration`, `//go:build e2e`
- Convenção de nomenclatura dos arquivos de teste
- Estrutura de setup/teardown nos testes
- Como mocks e stubs são usados nos testes unitários
- Como o banco de dados é configurado nos testes de integração
- Padrão de asserção usado (testify, stdlib, etc.)

### Domínio 4 — Injeção de Dependências com dig
Identifique como o `go.uber.org/dig` é usado. Procure por:
- Onde o container dig é criado e configurado
- Como os providers são registrados (`container.Provide(...)`)
- Como as dependências são invocadas (`container.Invoke(...)`)
- Convenção para funções construtoras (`NewXxx(deps) *Xxx`)

### Domínio 5 — Tratamento de Erros e Respostas HTTP
Identifique como erros são tratados e como as respostas HTTP são padronizadas. Procure por:
- Tipos de erro customizados (structs que implementam `error`)
- Como erros são propagados entre camadas
- Funções helper para responses (ex: `RespondError`, `RespondSuccess`)
- Mapeamento de erros de domínio para status HTTP
- Formato padrão do JSON de resposta (campos `message`, `error`, `data`, etc.)

### Domínio 6 — GORM: Repositórios e Migrations
Identifique como o GORM é usado. Procure por:
- Interface do repositório definida na camada de domínio
- Implementação do repositório na camada de infra
- Como a conexão com o banco é gerenciada
- Como as migrations são executadas (`AutoMigrate`, arquivos SQL, etc.)
- Uso de scopes, preloads e transações
- Convenção de nomenclatura de models GORM (tags `gorm:"..."`)