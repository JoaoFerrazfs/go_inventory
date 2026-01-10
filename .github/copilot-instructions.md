# Regras principais a seguir
---
## Geral
- Ao responder em conversas, responda sempre em português brasileiro, salvo pedido explícito para outro idioma.
- Comentários no código devem ser escritos em Inglês.
- Escreva código limpo, legível e sustentável seguindo padrões estabelecidos.
- Use nomes de variáveis e funções que expressem intenção (code identifiers devem permanecer em English nos exemplos).
- Mantenha funções e métodos pequenos e com responsabilidade única.
- Comente o código somente quando necessário para explicar lógica complexa; evite comentar em excesso.
- Use controle de versão (por exemplo, Git) com mensagens de commit claras.
- Realize code reviews para garantir qualidade e compartilhar conhecimento na equipe.
- Refatore continuamente para melhorar estrutura e desempenho.
- Mantenha-se atualizado sobre tecnologias e boas práticas.
- Priorize segurança em todas as etapas do desenvolvimento.
- Escreva testes cobrindo funcionalidades críticas e casos de borda.
- Não faça nada sem antes explicar o porquê e obter aprovação.
- Antes de qualquer alteração no codigo, explique a mudança proposta e aguarde aprovação.

## APIs
- Siga princípios RESTful no design das APIs.
- Use convenções consistentes para endpoints (por exemplo, substantivos no plural).
- Utilize códigos HTTP apropriados nas respostas.
- Documente endpoints com descrições, parâmetros e formatos de resposta.
- Valide todos os dados de entrada e trate erros de forma adequada.
- Implemente versionamento das APIs para gerenciar mudanças.
- Use autenticação e autorização para proteger endpoints.
- Otimize performance, incluindo paginação em coleções grandes.
- Escreva testes unitários e de integração para endpoints.
- Aplique práticas de segurança para prevenir vulnerabilidades comuns (ex.: SQL injection, XSS).

## Documentação
- Mantenha documentação abrangente (setup, guias de uso, referências de API).
- Use linguagem clara e concisa na documentação.
- Mantenha a documentação sincronizada com mudanças no código.
- Inclua exemplos e casos de uso para ilustrar funcionalidades.
- Organize documentação com índice para facilitar navegação.
- Use diagramas e elementos visuais quando fizer sentido.
- Estimule contribuições de toda a equipe à documentação.
- Revise e atualize a documentação regularmente.
- Armazene a documentação na pasta "/documents" do projeto.
- Cada documento deve seguir a mesma estrutura usada nos demais arquivos do projeto. Exemplo de pastas:
  - /documents/domain/Pallet
  - /documents/domain/PalletProduct
  - /documents/domain/PalletRackEntity
  - /documents/domain/User
- Forneça guias de troubleshooting para problemas comuns.


## Commits
- Use mensagens de commit claras e descritivas que resumam as mudanças realizadas.
- Siga um formato consistente para mensagens de commit, por exemplo:
  - feat: para novas funcionalidades
  - fix: para correções de bugs
  - docs: para atualizações de documentação
  - style: para mudanças de estilo (format, etc.)
  - refactor: para refatorações sem alteração de comportamento
  - test: para adição/atualização de testes
  - chore: para tarefas de manutenção (build, dependências, etc.)
- Mantenha commits pequenos e focados em uma mudança por vez.
- Evite commitar arquivos gerados ou artefatos de build.
- Use branches para gerenciar features, correções e releases.
- Mantenha sua branch atualizada com `main` regularmente.
- Revise mudanças antes de commitar para garantir qualidade.
- Use o imperativo nas mensagens de commit (ex.: "Add feature" em vez de "Added feature").
- Referencie números de issues quando aplicável.
