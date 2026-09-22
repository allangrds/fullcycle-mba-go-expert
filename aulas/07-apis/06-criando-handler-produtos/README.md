# Estrutura de Projeto Go e Sistema de Configuração

Este diretório demonstra a implementação de um sistema de configuração seguindo o **Go Standard Project Layout** e utilizando a biblioteca **Viper** para gerenciamento de configurações.

## 📋 Sumário

- [📋 Sumário](#-sumário)
- [🔄 Mudanças na Pasta configs/](#-mudanças-na-pasta-configs)
  - [Abordagem Anterior (LoadConfig)](#abordagem-anterior-loadconfig)
  - [Nova Abordagem (init function)](#nova-abordagem-init-function)
  - [Vantagens da Nova Abordagem](#vantagens-da-nova-abordagem)
  - [Métodos Getter](#métodos-getter)
- [📁 Estrutura de Pastas - Go Standard Project Layout](#-estrutura-de-pastas---go-standard-project-layout)
  - [Explicação Detalhada das Pastas](#explicação-detalhada-das-pastas)
    - [`/cmd` - Aplicações Principais](#cmd---aplicações-principais)
    - [`/internal` - Código Privado](#internal---código-privado)
    - [`/pkg` - Código Público](#pkg---código-público)
    - [`/configs` - Configurações](#configs---configurações)
    - [`/api` - Definições de API](#api---definições-de-api)
    - [`/test` - Testes Auxiliares](#test---testes-auxiliares)
- [🔧 Sistema de Configuração com Viper](#-sistema-de-configuração-com-viper)
  - [Dependências Utilizadas](#dependências-utilizadas)
  - [Estrutura de Configuração](#estrutura-de-configuração)
    - [Struct de Configuração](#struct-de-configuração)
    - [Variável Global](#variável-global)
  - [Função init() - Análise Detalhada](#função-init---análise-detalhada)
  - [Arquivo de Ambiente](#arquivo-de-ambiente)
- [🔤 Detalhes de Sintaxe Go](#-detalhes-de-sintaxe-go)
  - [Tags de Struct](#tags-de-struct)
  - [Ponteiros e Referências](#ponteiros-e-referências)
  - [Error Handling](#error-handling)
  - [Package-level Variables](#package-level-variables)
- [🚀 Uso da Configuração](#-uso-da-configuração)
  - [Importando e Usando Configurações](#importando-e-usando-configurações)
  - [Acessando Configurações](#acessando-configurações)
- [📊 Boas Práticas](#-boas-práticas)
  - [Estrutura de Projeto](#estrutura-de-projeto)
  - [Configuração](#configuração)
  - [Segurança](#segurança)
  - [Configuração por Ambiente](#configuração-por-ambiente)
- [🎯 Benefícios desta Abordagem](#-benefícios-desta-abordagem)
  - [Go Standard Layout](#go-standard-layout)
  - [Viper Configuration](#viper-configuration)

## 🔄 Mudanças na Pasta configs/

Esta seção demonstra uma **evolução significativa** na abordagem de configuração, mudando de uma função explícita `LoadConfig()` para uma inicialização automática usando a função `init()`.

### Abordagem Anterior (LoadConfig)

Na implementação anterior (`aulas/07-apis/02-criando-arquivo-configuracao`), o carregamento das configurações era feito através de uma função que precisava ser chamada explicitamente:

```go
// Abordagem anterior - Função explícita
func LoadConfig(path string) (*conf, error) {
    viper.SetConfigName("app_confi")
    viper.SetConfigType("env")
    viper.AddConfigPath(path)
    // ... resto da configuração
    return config, err
}

// Uso no main.go
func main() {
    cfg, err := configs.LoadConfig(".")
    if err != nil {
        log.Fatal(err)
    }
    // ... usar cfg
}
```

### Nova Abordagem (init function)

Na implementação atual, o carregamento acontece automaticamente através da função `init()`:

```go
// Nova abordagem - Inicialização automática
func init() {
    viper.SetConfigName("app_confi")
    viper.SetConfigType("env")
    viper.AddConfigPath(".")
    // ... configuração automática

    // Configuração é carregada automaticamente
    err := viper.ReadInConfig()
    if err != nil {
        panic(err)
    }

    err = viper.Unmarshal(&config)
    if err != nil {
        panic(err)
    }

    config.TokenAuth = jwtauth.New("HS256", []byte(config.JWTSecret), nil)
}

// Uso no main.go - Sem necessidade de chamar LoadConfig
func main() {
    // Configuração já está disponível!
    port := configs.GetWebServerPort()
    // ... usar diretamente
}
```

### Vantagens da Nova Abordagem

1. **🚀 Inicialização Automática**
   - Configurações são carregadas automaticamente na importação do pacote
   - Elimina a necessidade de chamar `LoadConfig()` explicitamente
   - Reduz boilerplate code no `main.go`

2. **🔒 Encapsulamento Melhorado**
   - Struct `Conf` agora é exportada (maiúscula)
   - Configuração interna permanece privada via variável `config`
   - Acesso controlado através de métodos getter

3. **🛡️ Fail-Fast Behavior**
   - Problemas de configuração são detectados na inicialização
   - Aplicação falha rapidamente se configurações são inválidas
   - Evita erros silenciosos durante execução

4. **📦 Simplicidade de Uso**
   - Interface mais limpa para consumers do pacote
   - Menos código necessário para usar as configurações
   - Padrão mais idiomático em Go

### Métodos Getter

A nova implementação introduz métodos getter para acesso controlado às configurações:

```go
// Métodos públicos para acessar configurações
func GetDBDriver() string       { return config.DBDriver }
func GetDBHost() string         { return config.DBHost }
func GetDBPort() string         { return config.DBPort }
func GetWebServerPort() string  { return config.WebServerPort }
func GetTokenAuth() *jwtauth.JWTAuth { return config.TokenAuth }
// ... outros getters
```

**Por que usar getters:**
- **Encapsulamento:** Protege dados internos de modificação acidental
- **Flexibilidade:** Permite adicionar validação ou transformação no futuro
- **Thread-safety:** Apenas leitura, sem riscos de concorrência
- **API limpa:** Interface consistente para acessar configurações

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
type Conf struct {
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

**Mudança importante:** A struct agora é **exportada** (`Conf` em vez de `conf`)

**Explicação das Tags `mapstructure`:**

- **`mapstructure:"DB_DRIVER"`:** Mapeia a variável de ambiente `DB_DRIVER` para o campo `DBDriver`
- **Conversão automática:** Viper converte string para int automaticamente (`JWT_EXPIRES_IN`)
- **Case sensitivity:** Tags devem corresponder exatamente aos nomes das variáveis
- **Flexibilidade:** Permite nomes diferentes entre struct e variável de ambiente

#### Variável Global

```go
var config *Conf
```

**Por que usar variável global:**

- **Singleton pattern:** Garantia de uma única instância de configuração
- **Acesso global:** Qualquer pacote pode acessar as configurações através dos getters
- **Performance:** Carregada uma vez na inicialização
- **Thread-safe:** Após carregamento, apenas leitura através dos getters

### Função init() - Análise Detalhada

```go
func init() {
```

**Função `init()`:** Executada automaticamente quando o pacote é importado, antes da função `main()`

#### 1. Configuração do Viper

```go
viper.SetConfigName("app_confi")  // Nome do arquivo (sem extensão)
viper.SetConfigType("env")        // Tipo do arquivo
viper.AddConfigPath(".")          // Diretório de busca (diretório atual)
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

### Importando e Usando Configurações

```go
package main

import (
    "fmt"
    "log"
    "path/to/configs"
)

func main() {
    // Configuração já está carregada automaticamente!
    // Não é necessário chamar LoadConfig()

    // Acessar configurações através dos getters
    port := configs.GetWebServerPort()
    dbHost := configs.GetDBHost()
    tokenAuth := configs.GetTokenAuth()

    fmt.Printf("Server will run on port: %s\n", port)
    fmt.Printf("Database host: %s\n", dbHost)
}
```

### Acessando Configurações

```go
// Database connection
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
    configs.GetDBUser(),
    configs.GetDBPassword(),
    configs.GetDBHost(),
    configs.GetDBPort(),
    configs.GetDBName())

// Web server
server := &http.Server{
    Addr: ":" + configs.GetWebServerPort(),
}

// JWT middleware
tokenAuth := configs.GetTokenAuth()
```

**Vantagens dos métodos getter:**

- **Simplicidade:** Não é necessário passar structs de configuração
- **Encapsulamento:** Dados internos ficam protegidos
- **Flexibilidade:** Facilita modificações futuras na implementação
- **Clareza:** Interface mais limpa e intuitiva

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

### Comparação: Antes vs. Agora

| Aspecto | Abordagem Anterior | Nova Abordagem |
|---------|-------------------|----------------|
| **Inicialização** | `cfg, err := configs.LoadConfig(".")` | Automática via `init()` |
| **Struct** | `conf` (privada) | `Conf` (exportada) |
| **Acesso** | `cfg.WebServerPort` | `configs.GetWebServerPort()` |
| **Error Handling** | Retorna erro para handling | Panic na inicialização |
| **Boilerplate** | Mais código no `main.go` | Menos código, mais limpo |
| **Encapsulamento** | Struct exposta diretamente | Acesso controlado via getters |

### Go Standard Layout

- **Familiarity:** Estrutura conhecida pela comunidade
- **Tooling:** Compatível com ferramentas Go
- **Maintenance:** Facilita manutenção e onboarding

### Viper Configuration

- **Flexibilidade:** Múltiplas fontes de configuração
- **Type Safety:** Conversão automática de tipos
- **Environment Support:** Precedência de variáveis de ambiente
- **Format Agnostic:** Suporte a JSON, YAML, TOML, ENV

### Nova Implementação com init()

- **Zero Configuration:** Funciona imediatamente após import
- **Fail-Fast:** Problemas detectados na inicialização
- **Clean API:** Interface consistente através de getters
- **Better Encapsulation:** Dados internos protegidos
- **Thread-Safe:** Apenas leitura após inicialização

Esta evolução representa uma **melhoria significativa** na arquitetura de configuração, tornando o código mais robusto, limpo e fácil de usar!
