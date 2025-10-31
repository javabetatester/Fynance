# Fynance - Sistema de Gestão Financeira Pessoal

Fynance é uma API RESTful para gestão financeira pessoal, desenvolvida em Go com foco em segurança, escalabilidade e facilidade de uso. A aplicação permite aos usuários gerenciar suas finanças pessoais, incluindo transações, metas financeiras, investimentos e relatórios.

## Funcionalidades

- **Autenticação e Autorização**: Sistema seguro com JWT para autenticação e middleware de autorização baseado em propriedade
- **Gestão de Usuários**: Cadastro, atualização, consulta e exclusão de usuários
- **Transações Financeiras**: Registro e categorização de receitas e despesas
- **Metas Financeiras**: Definição e acompanhamento de objetivos financeiros
- **Investimentos**: Controle de aplicações financeiras e retornos
- **Dashboard**: Visão consolidada da situação financeira do usuário

## Tecnologias Utilizadas

- **Go (Golang)**: Linguagem de programação principal
- **Gin**: Framework web para construção da API REST
- **GORM**: ORM (Object-Relational Mapping) para interação com o banco de dados
- **PostgreSQL**: Banco de dados relacional
- **JWT**: Autenticação baseada em tokens
- **bcrypt**: Criptografia segura de senhas

## Arquitetura

O projeto segue uma arquitetura limpa (Clean Architecture) com separação clara de responsabilidades:

- **Domain**: Contém as entidades de negócio e regras de domínio
- **Infrastructure**: Implementações concretas de repositórios e serviços externos
- **Middleware**: Componentes para processamento de requisições HTTP
- **Routes**: Handlers HTTP que conectam as requisições às regras de negócio
- **Utils**: Utilitários e serviços compartilhados

## Requisitos

- Go 1.25+
- PostgreSQL
- Variáveis de ambiente (opcionais):
  - `JWT_SECRET_KEY`: Chave secreta para assinatura de tokens JWT
  - `JWT_ISSUER`: Emissor dos tokens JWT

## Instalação e Execução

1. Clone o repositório:
   ```
   git clone https://github.com/seu-usuario/fynance.git
   cd fynance
   ```

2. Instale as dependências:
   ```
   go mod download
   ```

3. Configure o banco de dados PostgreSQL:
   - O arquivo `internal/infrastructure/db.go` contém as configurações de conexão

4. Execute a aplicação:
   ```
   go run cmd/api/main.go
   ```

5. A API estará disponível em `http://localhost:8080`

## Endpoints da API

### Rotas Públicas

- `POST /api/login`: Autenticação de usuário e geração de token JWT
- `POST /api/users`: Criação de novo usuário

### Rotas Privadas (Requerem Autenticação)

- `GET /api/users/:id`: Obter detalhes de um usuário
- `GET /api/users/email`: Buscar usuário por email
- `PUT /api/users/:id`: Atualizar dados de um usuário
- `DELETE /api/users/:id`: Excluir um usuário

## Autenticação

Para acessar as rotas privadas, é necessário incluir o token JWT no cabeçalho de autorização:

```
Authorization: Bearer seu_token_jwt
```

O token é obtido através da rota de login e tem validade de 24 horas.

## Segurança

- Senhas são armazenadas com hash bcrypt
- Validação de propriedade para garantir que usuários só acessem seus próprios recursos
- Verificação de método de assinatura JWT para prevenir ataques
- Validação de tokens expirados

## Estrutura de Diretórios

```
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── auth/
│   │   ├── dashboard/
│   │   ├── goal/
│   │   ├── investment/
│   │   ├── reports/
│   │   ├── transaction/
│   │   └── user/
│   ├── infrastructure/
│   │   ├── db.go
│   │   └── user_repository.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── jwt.go
│   ├── routes/
│   │   ├── auth.go
│   │   ├── handler.go
│   │   └── user.go
│   └── utils/
│       └── jwt_service.go
├── go.mod
├── go.sum
└── README.md
```

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).