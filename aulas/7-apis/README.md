# APIs em Go - Estudo Completo

Este módulo aborda o desenvolvimento de APIs REST em Go, cobrindo desde a estruturação básica de projetos até implementações avançadas com autenticação JWT e documentação Swagger.

## 📁 Estrutura do Módulo

```
7-apis/
├── 1-estruturando-diretorios/     # Organização inicial do projeto
├── 2-criando-arquivo-configuracao/ # Configuração com Viper
├── 3-outra-possibilidade-configuracao/ # Alternativas de configuração
├── 4-criando-entidade-user/       # Implementação da entidade User
├── 5-criando-entidade-product/    # Implementação da entidade Product
├── 6-criando-handler-produtos/    # Handlers para produtos
├── 7-ajustando-package-handlers/  # Refatoração dos handlers
└── 8-go-chi/                      # Implementação completa com Chi Router
```

## 🏗️ Arquitetura do Projeto (Clean Architecture)

A estrutura final segue os princípios da Clean Architecture:

```
projeto/
├── cmd/
│   └── server/
│       ├── main.go              # Ponto de entrada da aplicação
│       ├── .env                 # Variáveis de ambiente
│       └── test.db              # Banco SQLite
├── configs/
│   └── config.go                # Configurações da aplicação
├── docs/                        # Documentação Swagger gerada
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── dto/
│   │   └── dto.go               # Data Transfer Objects
│   ├── entity/
│   │   ├── product.go           # Entidade Product
│   │   ├── product_test.go      # Testes da entidade Product
│   │   ├── user.go              # Entidade User
│   │   └── user_test.go         # Testes da entidade User
│   └── infra/
│       ├── database/            # Camada de persistência
│       └── webserver/
│           └── handlers/        # Handlers HTTP
├── pkg/
│   └── entity/                  # Entidades compartilhadas
└── test/                        # Testes de integração
```

## 🔧 Tecnologias e Bibliotecas Utilizadas

### Dependências Principais

```go
// Router e middlewares
github.com/go-chi/chi v1.5.1           // Router HTTP minimalista
github.com/go-chi/jwtauth v1.2.0       // Middleware JWT para Chi

// Configuração
github.com/spf13/viper v1.20.1         // Gerenciamento de configurações

// Banco de dados
gorm.io/gorm v1.30.3                   // ORM para Go
gorm.io/driver/sqlite v1.6.0           // Driver SQLite para GORM

// Documentação (Swagger)
github.com/swaggo/swag v1.16.6         // Gerador de documentação Swagger
github.com/swaggo/http-swagger v1.3.4  // Handler Swagger para HTTP

// Utilitários
github.com/google/uuid v1.6.0          // Geração de UUIDs
golang.org/x/crypto v0.32.0            // Criptografia (bcrypt)

// Testes
github.com/stretchr/testify v1.10.0    // Framework de testes
```

## 📋 Conceitos Fundamentais Abordados

### 1. Estruturação de Projetos Go

- **Padrão de diretórios**: `cmd/`, `internal/`, `pkg/`, `configs/`
- **Separação de responsabilidades**: Entities, DTOs, Handlers, Database
- **Encapsulamento**: Uso do diretório `internal/` para código privado

### 2. Entidades de Negócio

#### User Entity
```go
type User struct {
    ID       entity.ID `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Password string    `json:"-"`  // Campo omitido na serialização JSON
}
```

**Conceitos aplicados:**
- **Validação de dados** no próprio modelo
- **Encriptação de senhas** com bcrypt
- **Uso de ponteiros** para evitar cópia desnecessária
- **Tags JSON** para controle de serialização

#### Product Entity
```go
type Product struct {
    ID        entity.ID `json:"id"`
    Name      string    `json:"name"`
    Price     int       `json:"price"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Conceitos aplicados:**
- **Validação de negócio** (preço, nome obrigatórios)
- **Timestamps automáticos**
- **Tratamento de erros customizados**

### 3. DTOs (Data Transfer Objects)

```go
type CreateProductInput struct {
    Name  string `json:"name"`
    Price int    `json:"price"`
}

type GetJWTInput struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

**Finalidade:**
- Separar a representação externa dos dados das entidades internas
- Controlar quais campos são expostos na API
- Facilitar validação de entrada

### 4. Configuração com Viper

```go
type Conf struct {
    DBDriver      string `mapstructure:"DB_DRIVER"`
    DBHost        string `mapstructure:"DB_HOST"`
    WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
    JWTSecret     string `mapstructure:"JWT_SECRET"`
    JWTExpiresIn  int    `mapstructure:"JWT_EXPIRES_IN"`
}
```

**Características:**
- Carregamento de variáveis de ambiente
- Arquivo `.env` para configurações locais
- `AutomaticEnv()` dá prioridade às variáveis de ambiente

### 5. Router Chi

O **Chi** é um router HTTP leve e idiomático para Go:

```go
router := chi.NewRouter()
router.Use(middleware.Logger)      // Log de requisições
router.Use(middleware.Recoverer)   // Recuperação de panics

// Agrupamento de rotas com middleware JWT
router.Route("/products", func(router chi.Router) {
    router.Use(jwtauth.Verifier(configs.TokenAuth))
    router.Use(jwtauth.Authenticator)
    
    router.Post("/", productHandler.CreateProduct)
    router.Get("/{id}", productHandler.GetProduct)
    router.Put("/{id}", productHandler.UpdateProduct)
    router.Delete("/{id}", productHandler.DeleteProduct)
})
```

**Vantagens do Chi:**
- **Lightweight**: Menor overhead que outros routers
- **Compatível com net/http**: Usa http.Handler padrão
- **Middleware integrado**: Fácil aplicação de middlewares
- **Sub-routers**: Organização hierárquica de rotas

### 6. Autenticação JWT

```go
// Middleware JWT
router.Use(jwtauth.Verifier(configs.TokenAuth))
router.Use(jwtauth.Authenticator)

// Geração de token
token, _ := configs.TokenAuth.Encode(claims)
```

**Implementação:**
- **Token JWT** com claims customizados
- **Middleware de verificação** automática
- **Proteção de rotas** sensíveis

### 7. GORM (ORM)

```go
// Auto-migrate
db.AutoMigrate(&entity.Product{}, &entity.User{})

// Operações CRUD
productDB := database.NewProduct(db)
product, err := productDB.Create(product)
```

**Características:**
- **Auto-migration**: Criação automática de tabelas
- **Repository Pattern**: Abstração da camada de dados
- **Interface-based**: Facilita testes e mocking

## 📖 Swagger/OpenAPI Documentation

### Instalação e Configuração

```bash
# Instalar o Swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# O binário fica em $GOPATH/bin ou $HOME/go/bin
# Adicione ao PATH se necessário
export PATH=$PATH:$(go env GOPATH)/bin
```

### Anotações Swagger no Código

#### Configuração Global (main.go)
```go
//@title Go Expert API Example
//@version 1.0
//@description This is a sample server for a Go Expert API Example.
//@termsOfService http://swagger.io/terms/

//@contact.name API Support
//@contact.url http://www.swagger.io/support
//@contact.email

//@license.name Apache 2.0
//@license.url http://www.apache.org/licenses/LICENSE-2.0.html

//@host localhost:8000
//@BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
```

#### Documentação de Endpoints
```go
// CreateProduct godoc
// @Summary Create product
// @Description Create product
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductInput true "Product data"
// @Success 201 {object} entity.Product
// @Failure 400 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    // implementação...
}
```

### Comandos Swag

```bash
# Gerar documentação Swagger
swag init -g cmd/server/main.go

# Especificar diretório de saída
swag init -g cmd/server/main.go -o docs/

# Incluir tipos de outros pacotes
swag init -g cmd/server/main.go --parseInternal
```

### Estrutura Gerada

Após executar `swag init`, são criados:

```
docs/
├── docs.go          # Código Go com a documentação embarcada
├── swagger.json     # Especificação OpenAPI em JSON
└── swagger.yaml     # Especificação OpenAPI em YAML
```

### Acessando a Documentação

```go
// Handler para servir a documentação
router.Get("/docs/*", httpSwagger.Handler(
    httpSwagger.URL("http://localhost:8000/docs/doc.json")
))
```

**URLs de acesso:**
- Swagger UI: `http://localhost:8000/docs/index.html`
- JSON Spec: `http://localhost:8000/docs/doc.json`
- YAML Spec: `http://localhost:8000/docs/doc.yaml`

## 🚀 Comandos para Execução

### Instalação de Dependências

```bash
# Inicializar módulo Go
go mod init github.com/seu-usuario/projeto-api

# Instalar dependências
go mod tidy

# Instalar Swag CLI
go install github.com/swaggo/swag/cmd/swag@latest
```

### Desenvolvimento

```bash
# Gerar documentação Swagger
swag init -g cmd/server/main.go

# Executar a aplicação
go run cmd/server/main.go

# Executar testes
go test ./...

# Executar testes com verbose
go test -v ./...

# Executar testes de uma entidade específica
go test -v ./internal/entity/
```

### Build e Deploy

```bash
# Build da aplicação
go build -o api cmd/server/main.go

# Executar o binário
./api

# Build para diferentes plataformas
GOOS=linux GOARCH=amd64 go build -o api-linux cmd/server/main.go
```

## 🔧 Variáveis de Ambiente (.env)

```env
DB_DRIVER=sqlite3
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=goexpert
WEB_SERVER_PORT=8000
JWT_SECRET=secret
JWT_EXPIRES_IN=3600
```

## 📚 Padrões e Boas Práticas

### 1. Repository Pattern
- Abstração da camada de dados
- Facilita testes com mocks
- Separação de responsabilidades

### 2. Dependency Injection
- Handlers recebem dependências via construtor
- Facilita testes unitários
- Código mais modular

### 3. Error Handling
- Erros customizados por domínio
- Tratamento consistente de erros HTTP
- Logs estruturados

### 4. Security
- Senhas sempre criptografadas
- JWT para autenticação stateless
- Middleware de autenticação

### 5. Testing
- Testes unitários para entidades
- Testes de integração para handlers
- Mocks para dependências externas

## 🧪 Testando a API

### Criar Usuário
```bash
curl -X POST http://localhost:8000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"João","email":"joao@email.com","password":"123456"}'
```

### Gerar Token JWT
```bash
curl -X POST http://localhost:8000/users/generate-token \
  -H "Content-Type: application/json" \
  -d '{"email":"joao@email.com","password":"123456"}'
```

### Criar Produto (com autenticação)
```bash
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_JWT" \
  -d '{"name":"Produto Teste","price":1000}'
```

### Listar Produtos
```bash
curl -X GET http://localhost:8000/products \
  -H "Authorization: Bearer SEU_TOKEN_JWT"
```

## 📋 Checklist de Estudo

- [ ] Entender a estrutura de diretórios Go
- [ ] Implementar entidades com validação
- [ ] Configurar Viper para gerenciar configurações
- [ ] Usar GORM para persistência
- [ ] Implementar handlers com Chi router
- [ ] Adicionar autenticação JWT
- [ ] Documentar API com Swagger
- [ ] Escrever testes unitários
- [ ] Testar endpoints manualmente
- [ ] Entender middleware e sua aplicação

## 🎯 Próximos Passos

1. **Melhorias na API**: Paginação, filtros, sorting
2. **Testes**: Cobertura completa de testes
3. **Docker**: Containerização da aplicação
4. **CI/CD**: Pipeline de deploy automático
5. **Monitoring**: Logs estruturados e métricas
6. **Database**: Migração para PostgreSQL/MySQL

Este módulo fornece uma base sólida para desenvolvimento de APIs REST em Go, seguindo boas práticas da comunidade e padrões da indústria.
