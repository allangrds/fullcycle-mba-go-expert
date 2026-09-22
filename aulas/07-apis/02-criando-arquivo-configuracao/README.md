# Estrutura de Projeto Go e Sistema de Configuração

Este diretório demonstra a implementação de um sistema de configuração seguindo o **Go Standard Project Layout** e utilizando a biblioteca **Viper** para gerenciamento de configurações.

## 📁 Estrutura de Pastas - Go Standard Project Layout

A estrutura de pastas segue o padrão amplamente adotado pela comunidade Go:

```text
2-criando-arquivo-configuracao/
├── api/                    # Definições de API (OpenAPI/Swagger)
├── cmd/                    # Aplicações principais
│   └── server/            # Aplicação do servidor
│       ├── main.go        # Ponto de entrada da aplicação
│       └── .env           # Variáveis de ambiente
├── configs/               # Arquivos e templates de configuração
│   └── config.go          # Lógica de carregamento de configuração
├── internal/              # Código privado da aplicação
├── pkg/                   # Código reutilizável por aplicações externas
├── test/                  # Testes adicionais e dados de teste
├── go.mod                 # Definição do módulo Go
└── go.sum                 # Checksums das dependências
```

### Explicação Detalhada das Pastas

#### `/cmd` - Aplicações Principais
- **Propósito:** Contém os pontos de entrada (main packages) das aplicações
- **Convenção:** Cada subdiretório é uma aplicação executável
- **Exemplo:** `/cmd/server` para aplicação servidor, `/cmd/cli` para ferramenta CLI
- **Por que usar:** Organiza múltiplas aplicações em um único repositório

```text
cmd/
├── server/     # go run cmd/server/main.go
├── worker/     # go run cmd/worker/main.go
└── migrate/    # go run cmd/migrate/main.go
```

#### `/internal` - Código Privado
- **Propósito:** Código que NÃO pode ser importado por outros projetos
- **Go enforcement:** O compilador Go impede import de pacotes internal
- **Uso típico:** Business logic, handlers, repositories, services
- **Vantagem:** Encapsulamento e proteção da API interna

```text
internal/
├── handler/    # HTTP handlers
├── service/    # Business logic
├── repository/ # Data access layer
└── middleware/ # Custom middleware
```

#### `/pkg` - Código Público
- **Propósito:** Código que PODE ser importado por outros projetos
- **Uso típico:** Utilities, clients, libraries
- **Cuidado:** Tudo aqui é API pública, mudanças quebram compatibilidade

```text
pkg/
├── httpclient/ # HTTP client utilities
├── validator/  # Custom validators
└── logger/     # Logging utilities
```

#### `/configs` - Configurações
- **Propósito:** Templates e arquivos de configuração
- **Conteúdo:** Config structs, parsing logic, default values
- **Não confundir:** Não coloque arquivos .env aqui, use /cmd/app/

#### `/api` - Definições de API
- **Propósito:** Schemas OpenAPI, Protocol Buffers, JSON schemas
- **Ferramentas:** Swagger, gRPC definitions
- **Geração:** Código gerado a partir dessas definições

#### `/test` - Testes Auxiliares
- **Propósito:** Dados de teste, helpers, testes de integração
- **Diferença:** Testes unitários ficam junto com o código (package_test.go)

## 🔧 Sistema de Configuração com Viper

### Dependências Utilizadas

```go
require (
    github.com/go-chi/jwtauth v1.2.0  // JWT authentication
    github.com/spf13/viper v1.20.1    // Configuration management
)
```

### Estrutura de Configuração

**Arquivo:** `configs/config.go`

```go
package configs

import (
    "github.com/go-chi/jwtauth"
    "github.com/spf13/viper"
)
```

#### Struct de Configuração

```go
type conf struct {
    DBDriver      string `mapstructure:"DB_DRIVER"`
    DBHost        string `mapstructure:"DB_HOST"`
    DBPort        string `mapstructure:"DB_PORT"`
    DBUser        string `mapstructure:"DB_USER"`
    DBPassword    string `mapstructure:"DB_PASSWORD"`
    DBName        string `mapstructure:"DB_NAME"`
    WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
    JWTSecret     string `mapstructure:"JWT_SECRET"`
    JWTExpiresIn  int    `mapstructure:"JWT_EXPIRES_IN"`
    TokenAuth     *jwtauth.JWTAuth
}
```

**Explicação das Tags `mapstructure`:**

- **`mapstructure:"DB_DRIVER"`:** Mapeia a variável de ambiente `DB_DRIVER` para o campo `DBDriver`
- **Conversão automática:** Viper converte string para int automaticamente (`JWT_EXPIRES_IN`)
- **Case sensitivity:** Tags devem corresponder exatamente aos nomes das variáveis
- **Flexibilidade:** Permite nomes diferentes entre struct e variável de ambiente

#### Variável Global

```go
var config *conf
```

**Por que usar variável global:**
- **Singleton pattern:** Garantia de uma única instância de configuração
- **Acesso global:** Qualquer pacote pode acessar as configurações
- **Performance:** Carregada uma vez na inicialização
- **Thread-safe:** Após carregamento, apenas leitura

### Função LoadConfig - Análise Detalhada

```go
func LoadConfig(path string) (*conf, error) {
```

**Parâmetro `path`:** Diretório onde buscar arquivos de configuração

#### 1. Configuração do Viper

```go
viper.SetConfigName("app_confi")  // Nome do arquivo (sem extensão)
viper.SetConfigType("env")        // Tipo do arquivo
viper.AddConfigPath(path)         // Diretório de busca
viper.SetConfigFile(".env")       // Arquivo específico
```

**Explicação:**
- **SetConfigName:** Nome base do arquivo de configuração
- **SetConfigType:** Formato do arquivo (env, json, yaml, toml)
- **AddConfigPath:** Adiciona diretório à lista de busca
- **SetConfigFile:** Força uso de arquivo específico (sobrescreve nome)

#### 2. Precedência de Configurações

```go
viper.AutomaticEnv()  // Variáveis de ambiente têm precedência
```

**Ordem de precedência (maior para menor):**
1. **Variáveis de ambiente do sistema**
2. **Flags de linha de comando** (se configurado)
3. **Arquivo de configuração** (.env)
4. **Valores padrão** (se definidos)

#### 3. Leitura do Arquivo

```go
err := viper.ReadInConfig()
if err != nil {
    panic(err)  // ❌ Para execução se arquivo não encontrado
}
```

**Comportamento:**
- **Sucesso:** Carrega variáveis do arquivo
- **Erro:** Arquivo não encontrado ou formato inválido
- **Panic:** Aplicação para completamente (adequado para configurações críticas)

#### 4. Unmarshaling para Struct

```go
err = viper.Unmarshal(&config)
if err != nil {
    panic(err)
}
```

**O que acontece:**
- **Viper mapeia** variáveis para campos da struct usando tags `mapstructure`
- **Conversão automática** de tipos (string → int, string → bool)
- **Validação implícita** de tipos durante conversão
- **Ponteiro (&config)** permite modificação da variável global

#### 5. Configuração do JWT

```go
config.TokenAuth = jwtauth.New("HS256", []byte(config.JWTSecret), nil)
```

**Parâmetros:**
- **"HS256":** Algoritmo de assinatura (HMAC SHA-256)
- **[]byte(config.JWTSecret):** Chave secreta convertida para bytes
- **nil:** Chave pública (não necessária para HMAC)

### Arquivo de Ambiente

**Arquivo:** `cmd/server/.env`

```properties
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=fullcycle
WEB_SERVER_PORT=8080
JWT_SECRET=secret
JWT_EXPIRES_IN=300
```

**Formato explicado:**
- **Chave=Valor:** Sintaxe simples sem espaços ao redor do `=`
- **Sem aspas:** Valores são tratados como strings por padrão
- **Tipos:** Viper converte automaticamente para tipos da struct
- **Comentários:** Linhas iniciadas com `#` são ignoradas

## 🔤 Detalhes de Sintaxe Go

### Tags de Struct

```go
type conf struct {
    DBDriver string `mapstructure:"DB_DRIVER"`
    //              └── tag literal string
}
```

**Sintaxe das tags:**
- **Backticks (`):** Delimitam tag literal string
- **Formato:** `key:"value"` ou `key:"value,option"`
- **Múltiplas tags:** `json:"name" mapstructure:"NAME" validate:"required"`

### Ponteiros e Referências

```go
var config *conf                    // Ponteiro para conf
err = viper.Unmarshal(&config)      // Endereço de config
return config, err                  // Retorna ponteiro
```

**Uso de ponteiros:**
- **`*conf`:** Tipo ponteiro para struct conf
- **`&config`:** Operador de endereço, passa referência
- **Nil checking:** Ponteiros podem ser nil (útil para campos opcionais)

### Error Handling

```go
if err != nil {
    panic(err)  // Para execução imediatamente
}
```

**Alternativas ao panic:**
```go
// Retornar erro (mais comum)
if err != nil {
    return nil, fmt.Errorf("failed to load config: %w", err)
}

// Log e continuar
if err != nil {
    log.Printf("Config warning: %v", err)
}
```

### Package-level Variables

```go
var config *conf  // Variável global do pacote
```

**Características:**
- **Escopo:** Acessível por todo o pacote
- **Inicialização:** Zero value (nil para ponteiros)
- **Thread-safety:** Cuidado com modificações concorrentes

## 🚀 Uso da Configuração

### Carregando Configuração

```go
package main

import "path/to/configs"

func main() {
    // Carrega configuração do diretório atual
    cfg, err := configs.LoadConfig(".")
    if err != nil {
        log.Fatal(err)
    }

    // Usar configuração
    fmt.Printf("Server will run on port: %s\n", cfg.WebServerPort)
}
```

### Acessando Configurações

```go
// Database connection
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
    cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

// Web server
server := &http.Server{
    Addr: ":" + cfg.WebServerPort,
}

// JWT middleware
tokenAuth := cfg.TokenAuth
```

## 📊 Boas Práticas

### Estrutura de Projeto

- **Consistência:** Siga o Go Standard Layout
- **Separação:** Mantenha código público (pkg) e privado (internal) separados
- **Organização:** Um executável por subdiretório em cmd/

### Configuração

- **Validação:** Valide configurações críticas após carregamento
- **Defaults:** Defina valores padrão para configurações opcionais
- **Secrets:** Use variáveis de ambiente para dados sensíveis
- **Documentação:** Documente todas as configurações disponíveis

### Segurança

```go
// ❌ Nunca commitar secrets
JWT_SECRET=secret123

// ✅ Usar variáveis de ambiente em produção
JWT_SECRET=${RANDOM_SECRET_FROM_ENV}
```

### Configuração por Ambiente

```bash
# Development
APP_ENV=development go run cmd/server/main.go

# Production
APP_ENV=production go run cmd/server/main.go
```

```go
// Carregar configuração baseada no ambiente
env := os.Getenv("APP_ENV")
configPath := fmt.Sprintf("configs/%s", env)
cfg, err := configs.LoadConfig(configPath)
```

## 🎯 Benefícios desta Abordagem

### Go Standard Layout
- **Familiarity:** Estrutura conhecida pela comunidade
- **Tooling:** Compatível com ferramentas Go
- **Maintenance:** Facilita manutenção e onboarding

### Viper Configuration
- **Flexibilidade:** Múltiplas fontes de configuração
- **Type Safety:** Conversão automática de tipos
- **Environment Support:** Precedência de variáveis de ambiente
- **Format Agnostic:** Suporte a JSON, YAML, TOML, ENV

Esta estrutura fornece uma base sólida para aplicações Go escaláveis e bem organizadas!
