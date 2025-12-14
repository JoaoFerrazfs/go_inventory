# Pull Request Description Instructions

Use the following template for pull request descriptions:

## Descrição
Essa seção deverá conter o motivo pelo qual o Pull Request está sendo aberto. Quaisquer decisões de arquitetura/design de código poderão ser explicadas por aqui.

## Pré-visualização (opcional)
Essa seção deverá incluir algum insumo para evidenciar o que foi implementado no seu Pull Request

Por exemplo:

- Diagramas
- Gif's ou Prints
- Se julgar necessário, mostrar o antes e depois

## Observações
Incluir nesta seção quaisquer pré-requisitos para o merge desse PR

Por exemplo:

- Virar ENV_VAR X
- Esperar PR [BG676] Prevent empty form submission exploding on backend #6666 antes de realizar o merge deste atual
- Incluir nesta seção também, quaisquer pré-requisitos para rodar a branch localmente (opcional)

Por exemplo:

- Virar ENV_VAR X
- Fazer build do docker: $ docker compose build
- Possuir uma conta com fidelidade
- etc

## História
Link para o card

## Reviewers
Ao aprovar esse PR você concorda que:

- O código está de acordo com nossas boas práticas de desenvolvimento;
- O código está de acordo com nossas boas práticas de testes;
- É responsável solidário pelo código que está sendo entregue;
- Em cada PR é esperado um deploy, exceto em casos excepcionais.

### Exemplo de Pull Request Description

```
## Descrição
Implementar geração de QR code para pallets para melhorar o rastreamento de inventário, permitindo acesso rápido aos detalhes do pallet via escaneamento.

## Pré-visualização (opcional)
![QR Code Example](path/to/qr-code.png)

## Observações
- Certifique-se de que a ENV_VAR QR_ENABLED esteja ativada.
- Fazer build do docker: $ docker compose build

## História
[Link para o card no Jira/Trello/etc.]

## Reviewers
Ao aprovar esse PR você concorda que:
- O código está de acordo com nossas boas práticas de desenvolvimento;
- O código está de acordo com nossas boas práticas de testes;
- É responsável solidário pelo código que está sendo entregue;
- Em cada PR é esperado um deploy, exceto em casos excepcionais.
```
