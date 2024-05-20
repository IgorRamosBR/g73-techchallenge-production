# g73-techchallenge-production

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

Este microsserviço é parte do sistema G73 TechChallenge e oferece funcionalidades relacionadas à gestão de pedidos de produção. Ele foi desenvolvido para facilitar a comunicação entre os clientes e o sistema de gestão de pedidos.


## Tecnologias Utilizadas

- Go (linguagem de programação)
- Gin (framework HTTP)
- Docker (opcional, para implantação)




## Funcionalidades
- **Recuperação de Pedidos:** O microsserviço permite aos clientes recuperar uma lista de todos os pedidos de produção registrados no sistema.

- **Atualização de Status de Pedido:** Os clientes podem atualizar o status de um pedido específico no sistema, indicando se está em processo de fabricação, concluído, ou qualquer outro status relevante.



## Como Executar
Para executar este microsserviço, siga estas etapas:

**1.** Clone este repositório:

**2.** Instale as dependências do Go:

```bash
git mod tudy

```

**3.** Defina as variáveis de ambiente necessárias, como **PORT**, **ORDER_URL**, e **ORDER_TIMEOUT**


**4.** Execute o Microsserviço:

```bash
go run main.go
```


## Endpoints

- **GET: /v1/orders:** Recupera todos os pedidos de produção.

- **PUT: /v1/orders/:id/status:** Atualiza o status de um pedido específico.

## Documentação e Coverage
[Documentation](https://github.com/IgorRamosBR/g73-techchallenge-production/tree/master/docs)

