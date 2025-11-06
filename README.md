# Fynance - Sistema de Gestão Financeira Pessoal

API RESTful para gestão financeira pessoal desenvolvida em Go com foco em segurança, escalabilidade e facilidade de uso. A aplicação permite aos usuários gerenciar suas finanças pessoais, incluindo transações, categorias, metas financeiras, investimentos e relatórios.

## Índice

- [Funcionalidades](#funcionalidades)
- [Tecnologias](#tecnologias)
- [Arquitetura](#arquitetura)
- [Requisitos](#requisitos)
- [Instalação](#instalação)
- [Configuração](#configuração)
- [Execução](#execução)
- [Documentação da API](#documentação-da-api)
- [Endpoints](#endpoints)
- [Autenticação](#autenticação)
- [Segurança](#segurança)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Testes](#testes)
- [Contribuição](#contribuição)
- [Licença](#licença)

## Funcionalidades

### Autenticação e Autorização
- Registro de novos usuários
- Autenticação via JWT (JSON Web Tokens)
- Validação de propriedade de recursos
- Sistema de planos de usuário com controle de acesso

### Gestão de Transações
- Registro de receitas e despesas
- Categorização de transações
- Consulta e filtragem de transações
- Atualização e exclusão de transações

### Categorias de Transações
- Criação de categorias personalizadas
- Listagem de categorias
- Atualização e exclusão de categorias

### Metas Financeiras
- Criação de metas financeiras
- Acompanhamento de progresso
- Atualização e exclusão de metas

### Investimentos
- Registro de investimentos
- Controle de contribuições e saques
- Cálculo de retorno sobre investimentos
- Consulta de histórico de investimentos

### Dashboard e Relatórios
- Visão consolidada da situação financeira
- Relatórios e análises financeiras

## Tecnologias

### Core
- **Go 1.25+**: Linguagem de programação principal
- **Gin Framework**: Framework web para construção da API REST
- **GORM**: ORM para interação com banco de dados PostgreSQL
- **PostgreSQL**: Banco de dados relacional

### Segurança e Autenticação
- **golang-jwt/jwt**: Geração e validação de tokens JWT
- **golang.org/x/crypto**: Criptografia bcrypt para senhas

### Documentação
- **Swagger/OpenAPI**: Documentação interativa da API
- **swaggo/swag**: Geração automática de documentação Swagger

### Utilidades
- **oklog/ulid/v2**: Geração de identificadores únicos (ULID)

## Arquitetura

O projeto segue os princípios de **Clean Architecture** e **SOLID**, com separação clara de responsabilidades:

### Camadas

- **Domain Layer** (`internal/domain/`): Contém as entidades de negócio, interfaces de repositórios e serviços, e regras de domínio
  - `auth/`: Autenticação e autorização
  - `user/`: Gestão de usuários
  - `transaction/`: Transações financeiras
  - `goal/`: Metas financeiras
  - `investment/`: Investimentos
  - `dashboard/`: Dashboard e análises
  - `reports/`: Relatórios financeiros

- **Infrastructure Layer** (`internal/infrastructure/`): Implementações concretas de repositórios e conexão com banco de dados
  - Conexão PostgreSQL via GORM
  - Implementação de repositórios para todas as entidades
  - Migrações automáticas de banco de dados

- **Middleware Layer** (`internal/middleware/`): Componentes para processamento de requisições HTTP
  - Autenticação JWT
  - Validação de propriedade de recursos
  - Validação de planos de usuário

- **Routes Layer** (`internal/routes/`): Handlers HTTP que conectam as requisições às regras de negócio
  - Handlers para autenticação
  - Handlers para transações
  - Handlers para categorias
  - Handlers para metas
  - Handlers para investimentos

- **Contracts Layer** (`internal/contracts/`): DTOs (Data Transfer Objects) e contratos de API

- **Utils Layer** (`internal/utils/`): Utilitários e serviços compartilhados

### Princípios Aplicados

- **Single Responsibility Principle (SRP)**: Cada camada e componente tem uma responsabilidade única
- **Dependency Inversion**: Dependências apontam para abstrações (interfaces), não implementações
- **Open/Closed Principle**: Extensível através de interfaces, fechado para modificação
- **DRY (Don't Repeat Yourself)**: Código reutilizável e modularizado
- **KISS (Keep It Simple, Stupid)**: Soluções simples e diretas

## Requisitos

### Software
- **Go 1.25+**: [Download e instalação](https://go.dev/dl/)
- **PostgreSQL 12+**: [Download e instalação](https://www.postgresql.org/download/)
- **Git**: Para clonar o repositório

### Variáveis de Ambiente (Opcional)
Atualmente, a conexão com o banco de dados está configurada diretamente no código. Para produção, recomenda-se usar variáveis de ambiente:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=postgres
DB_SSL_MODE=disable
JWT_SECRET_KEY=sua_chave_secreta_aqui
JWT_ISSUER=fynance
```

## Instalação

### 1. Clonar o Repositório

```bash
git clone https://github.com/seu-usuario/fynance.git
cd fynance
```

### 2. Instalar Dependências

```bash
go mod download
```

### 3. Configurar Banco de Dados PostgreSQL

Certifique-se de que o PostgreSQL está instalado e rodando. Por padrão, a aplicação espera:

- **Host**: localhost
- **Porta**: 5432
- **Usuário**: postgres
- **Senha**: admin
- **Database**: postgres

Para alterar essas configurações, edite o arquivo `internal/infrastructure/db.go`.

### 4. Criar Banco de Dados (se necessário)

```sql
CREATE DATABASE postgres;
```

As migrações são executadas automaticamente na inicialização da aplicação.

## Configuração

### Configuração do Banco de Dados

O arquivo `internal/infrastructure/db.go` contém a configuração de conexão. Para produção, recomenda-se:

1. Usar variáveis de ambiente
2. Implementar connection pooling adequado
3. Configurar SSL/TLS para conexões seguras
4. Usar credenciais seguras

### Migrações

As migrações são executadas automaticamente via GORM AutoMigrate na inicialização da aplicação para as seguintes entidades:

- User
- Goal
- Transaction
- Category
- Investment

## Execução

### Modo Desenvolvimento

```bash
go run cmd/api/main.go
```

### Build e Execução

```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

### Execução com Swagger

A documentação Swagger estará disponível em:

```
http://localhost:8080/swagger/index.html
```

A API estará disponível em:

```
http://localhost:8080/api
```

## Documentação da API

A documentação completa da API está disponível via Swagger UI:

```
http://localhost:8080/swagger/index.html
```

### Gerar Documentação Swagger

Para regenerar a documentação Swagger após alterações:

```bash
swag init -g cmd/api/main.go
```

## Endpoints

### Rotas Públicas

#### Autenticação

- **POST** `/api/auth/register` - Registro de novo usuário
  - Body: `{ "email": "string", "password": "string", "name": "string" }`
  - Response: `{ "message": "string" }`

- **POST** `/api/auth/login` - Autenticação de usuário
  - Body: `{ "email": "string", "password": "string" }`
  - Response: `{ "token": "string" }`

### Rotas Privadas (Requerem Autenticação)

#### Transações

- **POST** `/api/transactions` - Criar nova transação
- **GET** `/api/transactions` - Listar transações do usuário
- **GET** `/api/transactions/:id` - Obter transação específica
- **PATCH** `/api/transactions/:id` - Atualizar transação
- **DELETE** `/api/transactions/:id` - Excluir transação

#### Categorias

- **POST** `/api/categories` - Criar nova categoria
- **GET** `/api/categories` - Listar categorias do usuário
- **PATCH** `/api/categories/:id` - Atualizar categoria
- **DELETE** `/api/categories/:id` - Excluir categoria

#### Metas

- **POST** `/api/goals` - Criar nova meta financeira
- **GET** `/api/goals` - Listar metas do usuário
- **GET** `/api/goals/:id` - Obter meta específica
- **PATCH** `/api/goals/:id` - Atualizar meta
- **DELETE** `/api/goals/:id` - Excluir meta

#### Investimentos

- **POST** `/api/investments` - Criar novo investimento
- **GET** `/api/investments` - Listar investimentos do usuário
- **GET** `/api/investments/:id` - Obter investimento específico
- **POST** `/api/investments/:id/contribution` - Realizar contribuição
- **POST** `/api/investments/:id/withdraw` - Realizar saque
- **GET** `/api/investments/:id/return` - Obter retorno do investimento
- **PATCH** `/api/investments/:id` - Atualizar investimento
- **DELETE** `/api/investments/:id` - Excluir investimento

## Autenticação

Todas as rotas privadas requerem autenticação via JWT. Para acessar essas rotas:

1. Faça login via `POST /api/auth/login` com suas credenciais
2. Receba o token JWT na resposta
3. Inclua o token no header `Authorization` de todas as requisições subsequentes:

```
Authorization: Bearer <seu_token_jwt>
```

### Exemplo de Requisição Autenticada

```bash
curl -X GET http://localhost:8080/api/transactions \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Validade do Token

Os tokens JWT têm validade de 24 horas. Após expirar, é necessário fazer login novamente.

## Segurança

### Medidas Implementadas

- **Criptografia de Senhas**: Senhas são armazenadas com hash bcrypt (cost >= 12)
- **JWT com Assinatura**: Tokens JWT são assinados e validados
- **Validação de Propriedade**: Middleware garante que usuários só acessem seus próprios recursos
- **Validação de Planos**: Sistema de planos permite controle granular de acesso
- **Validação de Entrada**: Validação de dados de entrada em todos os endpoints
- **HTTPS Recomendado**: Para produção, utilize HTTPS para proteger dados em trânsito

### Boas Práticas

- Nunca exponha credenciais em logs ou mensagens de erro
- Use variáveis de ambiente para configurações sensíveis
- Implemente rate limiting em produção
- Configure CORS adequadamente
- Monitore tentativas de acesso não autorizadas

## Estrutura do Projeto

```
Fynance/
├── cmd/
│   └── api/
│       └── main.go                    # Ponto de entrada da aplicação
├── internal/
│   ├── contracts/                     # DTOs e contratos de API
│   │   ├── auth.go
│   │   ├── common.go
│   │   ├── goal.go
│   │   ├── investment.go
│   │   ├── transaction.go
│   │   └── user.go
│   ├── domain/                        # Camada de domínio (regras de negócio)
│   │   ├── auth/                      # Autenticação e autorização
│   │   ├── dashboard/                 # Dashboard e análises
│   │   ├── goal/                      # Metas financeiras
│   │   ├── investment/                # Investimentos
│   │   ├── reports/                   # Relatórios
│   │   ├── transaction/               # Transações
│   │   └── user/                      # Usuários
│   ├── infrastructure/                # Camada de infraestrutura
│   │   ├── db.go                      # Conexão com banco de dados
│   │   ├── goal_repository.go
│   │   ├── investment_repository.go
│   │   ├── transaction_category_repository.go
│   │   ├── transaction_repository.go
│   │   └── user_repository.go
│   ├── middleware/                    # Middlewares HTTP
│   │   ├── auth.go                    # Middleware de autenticação
│   │   ├── jwt.go                     # Serviço JWT
│   │   ├── jwt_service.go
│   │   └── plan_validator.go          # Validação de planos
│   ├── routes/                        # Handlers HTTP
│   │   ├── authentication.go
│   │   ├── goal.go
│   │   ├── handler.go
│   │   ├── investment.go
│   │   ├── transaction_category.go
│   │   └── transaction.go
│   └── utils/                         # Utilitários
│       └── ulid_utils.go
├── docs/                              # Documentação Swagger
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod                             # Gerenciamento de dependências
├── go.sum
└── README.md                          # Este arquivo
```

## Testes

### Executar Testes

```bash
go test ./...
```

### Executar Testes com Coverage

```bash
go test -cover ./...
```

### Executar Testes com Coverage Detalhado

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Benchmarks

```bash
go test -bench=. ./...
```

## Contribuição

Contribuições são bem-vindas! Para contribuir:

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Diretrizes de Contribuição

- Siga os princípios SOLID, DRY, KISS
- Mantenha o código limpo e autoexplicativo
- Adicione testes para novas funcionalidades
- Atualize a documentação Swagger quando necessário
- Siga os padrões de código Go (gofmt, go vet)

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

---

**Desenvolvido por Bernardo**

Para mais informações, consulte a documentação Swagger em `http://localhost:8080/swagger/index.html` após iniciar a aplicação.
