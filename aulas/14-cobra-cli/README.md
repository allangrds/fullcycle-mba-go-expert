# 🐍 Cobra CLI — Construindo Aplicações de Linha de Comando em Go

> Aprenda a criar CLIs profissionais em Go com o Cobra, a mesma biblioteca usada pelo Docker, Kubernetes e GitHub CLI.

---

## 📑 Sumário

- [🤔 O que é Cobra CLI?](#-o-que-é-cobra-cli)
  - [Por que não usar os.Args direto?](#por-que-não-usar-osargs-direto)
  - [Quando usar Cobra?](#quando-usar-cobra)
- [🗂️ Arquitetura do Projeto](#️-arquitetura-do-projeto)
  - [Estrutura de Pastas](#estrutura-de-pastas)
  - [Hierarquia de Comandos](#hierarquia-de-comandos)
  - [Como o init() conecta as peças](#como-o-init-conecta-as-peças)
- [📖 Conceitos Fundamentais](#-conceitos-fundamentais)
  - [Comando Raiz (Root Command)](#comando-raiz-root-command)
  - [Subcomandos](#subcomandos)
  - [Flags Locais vs Persistentes](#flags-locais-vs-persistentes)
  - [Run vs RunE](#run-vs-rune)
  - [MarkFlagsRequiredTogether](#markflagsrequiredtogether)
- [⚙️ Como Funciona o Projeto](#️-como-funciona-o-projeto)
  - [Fluxo de Execução](#fluxo-de-execução)
  - [Padrão Construtor de Comando](#padrão-construtor-de-comando)
- [🗄️ Camada de Banco de Dados](#️-camada-de-banco-de-dados)
  - [SQLite — banco embutido, sem servidor](#sqlite--banco-embutido-sem-servidor)
  - [go-sqlite3 e CGO](#go-sqlite3-e-cgo)
  - [UUID como identificador único](#uuid-como-identificador-único)
  - [Padrão DAO](#padrão-dao)
- [✅ Boas Práticas Presentes no Projeto](#-boas-práticas-presentes-no-projeto)
  - [1. Injeção de Dependência via Construtor](#1-injeção-de-dependência-via-construtor)
  - [2. RunE para Propagação Correta de Erros](#2-rune-para-propagação-correta-de-erros)
  - [3. MarkFlagsRequiredTogether para Validação Declarativa](#3-markflagsrequiredtogether-para-validação-declarativa)
  - [4. Separação cmd/ vs internal/](#4-separação-cmd-vs-internal)
  - [5. GetDb() Centralizado](#5-getdb-centralizado)
  - [6. defer rows.Close() no DAO](#6-defer-rowsclose-no-dao)
- [🛡️ O que as Boas Práticas Evitaram](#️-o-que-as-boas-práticas-evitaram)
- [🚀 O que Poderia Ser Adicionado](#-o-que-poderia-ser-adicionado)
- [⚠️ Principais Problemas ao Trabalhar com Cobra CLI](#️-principais-problemas-ao-trabalhar-com-cobra-cli)
  - [1. Dependência de Ordem no init()](#1-dependência-de-ordem-no-init)
  - [2. Flags Globais Acidentais](#2-flags-globais-acidentais)
  - [3. panic em vez de retornar erro](#3-panic-em-vez-de-retornar-erro)
  - [4. Conexão com Banco Não Fechada](#4-conexão-com-banco-não-fechada)
  - [5. Conflito de Nomes de Flags](#5-conflito-de-nomes-de-flags)
  - [6. args vs flags — quando usar cada um](#6-args-vs-flags--quando-usar-cada-um)
- [🔧 Como Usar o Projeto](#-como-usar-o-projeto)
  - [Pré-requisitos](#pré-requisitos)
  - [Criando o Schema do Banco](#criando-o-schema-do-banco)
  - [Build e Execução](#build-e-execução)
  - [Exemplos de Comandos](#exemplos-de-comandos)
- [📖 Glossário](#-glossário)
- [🎯 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é Cobra CLI?

Imagine que você vai a um caixa eletrônico. Ao chegar, você vê um menu principal:

```
1. Sacar dinheiro
2. Ver saldo
3. Transferir
4. Voltar
```

Você escolhe "1. Sacar dinheiro" e aparece um novo menu:

```
Informe o valor:
[ $50  $100  $200  Outro valor ]
```

Uma aplicação de linha de comando (CLI) funciona exatamente assim, mas no terminal. Em vez de clicar, você digita comandos:

```bash
$ minha-app category create --name "Tecnologia" --description "Cursos de tech"
$ minha-app category list
```

**Cobra** é uma biblioteca Go que ajuda a construir essas CLIs de forma organizada, com suporte a subcomandos aninhados, flags tipadas, mensagens de ajuda automáticas e validação de entrada — sem você precisar implementar nada disso do zero.

> Cobra é a mesma biblioteca usada por projetos como **Docker**, **Kubernetes**, **GitHub CLI** e **Hugo**.

### Por que não usar os.Args direto?

Em Go, você consegue capturar os argumentos da linha de comando com `os.Args`:

```go
// ❌ Problema — funciona, mas vira um pesadelo rápido
func main() {
    args := os.Args[1:]
    if len(args) == 0 {
        fmt.Println("Uso: minha-app <comando>")
        return
    }
    if args[0] == "create" {
        if len(args) < 3 {
            fmt.Println("Uso: create <nome> <descricao>")
            return
        }
        // onde está o --name? o --description?
        // e se o usuário inverter a ordem?
        // e a mensagem de ajuda? e os erros?
    }
}
```

Essa abordagem funciona para scripts simples, mas escala muito mal. Com Cobra você declara o que quer e a biblioteca cuida do resto:

```go
// ✅ Correto — declarativo, auto-documentado
createCmd := &cobra.Command{
    Use:   "create",
    Short: "Cria uma nova categoria",
    RunE:  runCreate(categoryDb),
}
createCmd.Flags().StringP("name", "n", "", "Nome da categoria")
createCmd.Flags().StringP("description", "d", "", "Descrição da categoria")
```

### Quando usar Cobra?

| Situação | Ferramenta recomendada |
|----------|------------------------|
| Script simples, 1-2 argumentos | `os.Args` ou `flag` padrão do Go |
| Ferramenta com múltiplos subcomandos | **Cobra** ✅ |
| Ferramenta de administração/DevOps | **Cobra** ✅ |
| Precisa de autocompletion no shell | **Cobra** ✅ (gera automaticamente) |
| CLI integrada com config file + env vars | **Cobra + Viper** ✅ |

---

## 🗂️ Arquitetura do Projeto

### Estrutura de Pastas

```
14-cobra-cli/
├── main.go                      ← Ponto de entrada, apenas chama cmd.Execute()
├── go.mod                       ← Declaração do módulo e dependências
├── go.sum                       ← Hash de verificação das dependências
├── data.db                      ← Banco de dados SQLite (arquivo local)
├── cmd/                         ← Camada de apresentação (comandos CLI)
│   ├── root.go                  ← Comando raiz, inicialização do banco
│   ├── category.go              ← Comando pai "category"
│   ├── create.go                ← Subcomando "category create"
│   └── list.go                  ← Subcomando "category list" (placeholder)
└── internal/                    ← Lógica de negócio (acesso ao banco)
    └── database/
        ├── category.go          ← DAO de categorias (Create, FindAll, Find...)
        └── course.go            ← DAO de cursos (Create, FindAll, FindByCategoryID...)
```

**Por que `internal/`?** No Go, o diretório `internal/` tem um significado especial: pacotes dentro dele só podem ser importados por código do mesmo módulo. Isso protege a camada de banco de acesso externo acidental.

### Hierarquia de Comandos

```
16-CLI (root)
│
└── category                     ← "Contexto": operações sobre categorias
    │
    ├── create                   ← Cria uma nova categoria
    │   ├── --name, -n           ← Nome da categoria (obrigatório com --description)
    │   └── --description, -d   ← Descrição da categoria (obrigatório com --name)
    │
    └── list                     ← Lista categorias (placeholder)
```

Essa hierarquia é a mesma lógica de ferramentas como `git`:

```
git
├── commit   -m "mensagem"
├── push     origin main
└── branch   -d minha-branch
```

### Como o `init()` conecta as peças

Em Go, a função `init()` é executada automaticamente quando o pacote é carregado — antes do `main()`. O Cobra usa esse mecanismo para registrar subcomandos:

```
Ordem de inicialização:
1. root.go    → define rootCmd
2. category.go → init() adiciona categoryCmd ao rootCmd
3. create.go   → init() adiciona createCmd ao categoryCmd
4. list.go     → init() adiciona listCmd ao categoryCmd
5. main.go     → chama cmd.Execute(), que usa a árvore montada
```

---

## 📖 Conceitos Fundamentais

### Comando Raiz (Root Command)

O comando raiz é o "menu principal" da aplicação. É o que aparece quando o usuário digita o nome do binário sem nenhum subcomando:

```go
// cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "16-CLI",          // nome do binário
    Short: "Descrição curta", // aparece em listas
    Long:  `Descrição longa`, // aparece no --help
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
```

```bash
$ ./16-CLI --help
A longer description that spans multiple lines...

Usage:
  16-CLI [command]

Available Commands:
  category    A brief description of your command
  help        Help about any command

Flags:
  -h, --help     help for 16-CLI
  -t, --toggle   Help message for toggle
```

O Cobra gera essa saída de ajuda **automaticamente** a partir dos campos `Short`, `Long` e das flags registradas.

### Subcomandos

Subcomandos são agrupamentos lógicos de funcionalidades. O comando `category` funciona como um "namespace":

```go
// cmd/category.go
var categoryCmd = &cobra.Command{
    Use:   "category",
    Short: "A brief description of your command",
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help() // exibe ajuda quando chamado sem subcomando
    },
}

func init() {
    rootCmd.AddCommand(categoryCmd) // registra como filho do root
}
```

Quando o usuário chama `./16-CLI category` sem especificar `create` ou `list`, o comando exibe sua própria ajuda — o que é um comportamento amigável.

### Flags Locais vs Persistentes

```
Flags:
  Locais     → só funcionam no comando onde foram definidas
  Persistentes → funcionam no comando e em todos os seus filhos
```

```go
// Exemplos

// Flag local — só funciona em rootCmd
rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

// Flag persistente — funciona em rootCmd E em todos os subcomandos
rootCmd.PersistentFlags().String("config", "", "Arquivo de configuração")

// Flag local de subcomando — só funciona em createCmd
createCmd.Flags().StringP("name", "n", "", "Nome da categoria")
```

**Regra de ouro:** use flags persistentes apenas para configurações que fazem sentido em todos os subcomandos (como `--config`, `--verbose`, `--output`). Para parâmetros específicos de um comando, use flags locais.

### Run vs RunE

Cobra oferece dois campos para definir a lógica de um comando:

```go
// ❌ Run — engole o erro
var cmd = &cobra.Command{
    Run: func(cmd *cobra.Command, args []string) {
        err := fazAlgoQuePoderFalhar()
        if err != nil {
            fmt.Println(err) // imprime, mas o processo sai com código 0
        }
    },
}

// ✅ RunE — propaga o erro corretamente
var cmd = &cobra.Command{
    RunE: func(cmd *cobra.Command, args []string) error {
        err := fazAlgoQuePoderFalhar()
        if err != nil {
            return err // Cobra imprime o erro e sai com código 1
        }
        return nil
    },
}
```

O código de saída importa muito em scripts e pipelines de CI/CD:

```bash
$ ./16-CLI category create --name "Test" --description "Desc"
$ echo $?   # código de saída: 0 = sucesso, 1 = erro
```

### MarkFlagsRequiredTogether

O projeto usa uma feature poderosa do Cobra para validar flags em conjunto:

```go
// cmd/create.go
createCmd.MarkFlagsRequiredTogether("name", "description")
```

Isso significa: se o usuário fornecer `--name`, precisa fornecer `--description` também, e vice-versa. Sem isso, o usuário poderia criar uma categoria sem descrição sem receber nenhum aviso.

```bash
# ❌ Erro claro — muito melhor que silenciosamente criar registro incompleto
$ ./16-CLI category create --name "Tecnologia"
Error: if any flags in the group [name description] are set they must all be set; missing [description]
```

---

## ⚙️ Como Funciona o Projeto

### Fluxo de Execução

```
Usuário digita:
$ ./16-CLI category create --name "Go" --description "Cursos de Go"

                    ┌─────────────┐
                    │   main.go   │
                    │             │
                    │ cmd.Execute()│
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  root.go    │
                    │             │
                    │  rootCmd    │
                    │  .Execute() │
                    └──────┬──────┘
                           │ "category" reconhecido
                    ┌──────▼──────┐
                    │ category.go │
                    │             │
                    │ categoryCmd │
                    └──────┬──────┘
                           │ "create" reconhecido
                    ┌──────▼──────┐
                    │  create.go  │
                    │             │
                    │  createCmd  │
                    │  .RunE()    │
                    └──────┬──────┘
                           │
                    ┌──────▼──────────────┐
                    │ internal/database/  │
                    │ category.go         │
                    │                     │
                    │ categoryDb.Create() │
                    │ INSERT INTO ...     │
                    └─────────────────────┘
```

### Padrão Construtor de Comando

O código do `create.go` usa um padrão elegante: em vez de declarar o comando como variável global, ele usa uma **função construtora** que recebe as dependências:

```go
// cmd/create.go

// ✅ Construtor — recebe dependência como parâmetro
func newCreateCmd(categoryDb database.Category) *cobra.Command {
    return &cobra.Command{
        Use:   "create",
        Short: "Create a new category",
        RunE:  runCreate(categoryDb),
    }
}

// Separa a lógica do comando da sua definição
func runCreate(categoryDb database.Category) RunEFunc {
    return func(cmd *cobra.Command, args []string) error {
        name, _ := cmd.Flags().GetString("name")
        description, _ := cmd.Flags().GetString("description")
        _, err := categoryDb.Create(name, description)
        if err != nil {
            return err
        }
        return nil
    }
}

func init() {
    createCmd := newCreateCmd(GetCategoryDB(GetDb())) // injeta a dependência
    categoryCmd.AddCommand(createCmd)
    createCmd.Flags().StringP("name", "n", "", "Name of the category")
    createCmd.Flags().StringP("description", "d", "", "Description of the category")
    createCmd.MarkFlagsRequiredTogether("name", "description")
}
```

A função `RunEFunc` é um tipo alias definido em `root.go` que torna o código mais legível:

```go
// cmd/root.go
type RunEFunc func(cmd *cobra.Command, args []string) error
```

O `root.go` também centraliza a criação das dependências:

```go
// cmd/root.go
func GetDb() *sql.DB {
    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        panic(err)
    }
    return db
}

func GetCategoryDB(db *sql.DB) database.Category {
    return *database.NewCategory(db)
}
```

---

## 🗄️ Camada de Banco de Dados

### SQLite — banco embutido, sem servidor

SQLite é um banco de dados que vive em um único arquivo. Não precisa de servidor, não precisa de instalação separada, não precisa de configuração de rede. O arquivo `data.db` na raiz do projeto **é** o banco de dados.

```
PostgreSQL/MySQL:          SQLite:
┌──────────────┐          ┌──────────────┐
│ Sua Aplicação│          │ Sua Aplicação│
└──────┬───────┘          │   +          │
       │ TCP/IP           │   SQLite     │
┌──────▼───────┐          │   Library    │
│   Servidor   │          └──────┬───────┘
│   de Banco   │                 │ leitura/escrita
└──────┬───────┘          ┌──────▼───────┐
       │                  │  data.db     │
┌──────▼───────┐          │  (arquivo)   │
│  Arquivos    │          └──────────────┘
│  de Dados    │
└──────────────┘
```

Isso o torna perfeito para CLIs, ferramentas de desktop e aplicações que precisam de persistência sem infraestrutura.

### go-sqlite3 e CGO

O driver `github.com/mattn/go-sqlite3` é especial: ele inclui o código C do SQLite e compila junto com o Go usando **CGO** (C Go — integração entre Go e C).

```go
// cmd/root.go
import _ "github.com/mattn/go-sqlite3"
// O _ significa "importar apenas pelos efeitos colaterais"
// Isso registra o driver "sqlite3" no pacote database/sql
```

**Implicação prática:** para compilar este projeto, você precisa de um compilador C instalado (`gcc` no Linux/Mac, `mingw` no Windows). O build fica assim:

```bash
# Compilação normal funciona
go build -o cli .

# Se não tiver GCC, vai receber:
# cgo: C compiler "gcc" not found
```

### UUID como identificador único

O projeto usa `github.com/google/uuid` para gerar identificadores únicos para cada registro:

```go
// internal/database/category.go
func (c *Category) Create(name string, description string) (Category, error) {
    id := uuid.New().String() // ex: "550e8400-e29b-41d4-a716-446655440000"
    _, err := c.db.Exec(
        "INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)",
        id, name, description,
    )
    // ...
}
```

**Por que UUID em vez de auto-increment (1, 2, 3...)?**

| Auto-increment | UUID |
|----------------|------|
| Previsível (enumerable) | Difícil de adivinhar |
| Conflito ao mesclar bancos | Globalmente único |
| Revela quantidade de registros | Não vaza informação |
| Simples e compacto | Mais longo (36 chars) |

Para uma ferramenta de gestão de cursos, UUID é a escolha certa — especialmente se os dados forem sincronizados entre ambientes.

### Padrão DAO

DAO significa **Data Access Object** — um objeto cuja única responsabilidade é conversar com o banco de dados. O projeto separa isso claramente:

```go
// internal/database/category.go

type Category struct {
    db          *sql.DB  // conexão (privada, só o DAO tem acesso)
    ID          string
    Name        string
    Description string
}

// Construtor — recebe o banco como dependência
func NewCategory(db *sql.DB) *Category {
    return &Category{db: db}
}

// Métodos de acesso ao banco
func (c *Category) Create(name, description string) (Category, error) { ... }
func (c *Category) FindAll() ([]Category, error) { ... }
func (c *Category) Find(id string) (Category, error) { ... }
func (c *Category) FindByCourseID(courseID string) (Category, error) { ... }
```

A camada `cmd/` nunca escreve SQL — ela só usa os métodos do DAO. Isso significa que se um dia você trocar SQLite por PostgreSQL, só precisa mudar o `internal/database/`, e os comandos continuam funcionando.

---

## ✅ Boas Práticas Presentes no Projeto

### 1. Injeção de Dependência via Construtor

**O problema:** se o comando criasse a conexão com o banco internamente, seria impossível testar sem ter um banco real disponível.

```go
// ❌ Acoplamento forte — impossível testar
var createCmd = &cobra.Command{
    RunE: func(cmd *cobra.Command, args []string) error {
        db, _ := sql.Open("sqlite3", "./data.db") // dependência hardcoded
        categoryDb := database.NewCategory(db)
        // ...
    },
}

// ✅ Injeção de dependência — testável e flexível
func newCreateCmd(categoryDb database.Category) *cobra.Command {
    return &cobra.Command{
        RunE: runCreate(categoryDb), // dependência vem de fora
    }
}
```

Com injeção, um teste pode passar um `categoryDb` falso (mock) com dados em memória, sem tocar em arquivo nenhum.

### 2. RunE para Propagação Correta de Erros

**O problema:** se o banco estiver indisponível ou a query falhar, o programa precisa sinalizar isso ao ambiente que o chamou.

```go
// ❌ Run — processo sempre sai com código 0 (sucesso)
Run: func(cmd *cobra.Command, args []string) {
    err := categoryDb.Create(name, description)
    if err != nil {
        fmt.Println("Erro:", err) // imprime, mas não sinaliza falha
    }
}

// ✅ RunE — processo sai com código 1 em caso de erro
RunE: func(cmd *cobra.Command, args []string) error {
    _, err := categoryDb.Create(name, description)
    if err != nil {
        return err // Cobra cuida de imprimir e sair com código correto
    }
    return nil
}
```

Isso é crítico para scripts de automação:

```bash
#!/bin/bash
./cli category create --name "Go" --description "Cursos de Go"
if [ $? -ne 0 ]; then
    echo "Falhou! Abortando..."
    exit 1
fi
echo "Criado com sucesso!"
```

### 3. MarkFlagsRequiredTogether para Validação Declarativa

**O problema:** criar uma categoria sem nome ou sem descrição deixaria um registro incompleto no banco — e o Cobra não sabe disso automaticamente.

```go
// ❌ Validação manual — verbosa e fácil de esquecer
RunE: func(cmd *cobra.Command, args []string) error {
    name, _ := cmd.Flags().GetString("name")
    description, _ := cmd.Flags().GetString("description")
    if name == "" && description != "" {
        return fmt.Errorf("--name é obrigatório quando --description é fornecido")
    }
    if description == "" && name != "" {
        return fmt.Errorf("--description é obrigatório quando --name é fornecido")
    }
    // ...
}

// ✅ Declarativo — o Cobra valida automaticamente antes de chamar RunE
createCmd.MarkFlagsRequiredTogether("name", "description")
```

### 4. Separação cmd/ vs internal/

**O problema:** misturar lógica de negócio dentro dos handlers de comandos torna o código difícil de reutilizar e testar.

```
cmd/          → "O que o usuário pediu e como responder"
               Responsabilidade: ler flags, chamar o serviço, formatar saída

internal/     → "Como fazer de verdade"
               Responsabilidade: SQL, regras de negócio, validações de domínio
```

O resultado é que o comando `create.go` tem apenas 7 linhas de lógica real — e toda a complexidade de banco fica no DAO.

### 5. GetDb() Centralizado

**O problema:** se cada comando abrisse sua própria conexão de formas diferentes, seria difícil mudar o caminho do banco ou adicionar configuração.

```go
// cmd/root.go — único lugar que sabe como abrir o banco
func GetDb() *sql.DB {
    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        panic(err)
    }
    return db
}
```

Se amanhã o caminho do banco vier de uma variável de ambiente, você muda só aqui.

### 6. defer rows.Close() no DAO

**O problema:** em Go, quando você consulta múltiplas linhas com `db.Query()`, o resultado fica aberto consumindo recursos até ser fechado.

```go
// internal/database/category.go
func (c *Category) FindAll() ([]Category, error) {
    rows, err := c.db.Query("SELECT id, name, description FROM categories")
    if err != nil {
        return nil, err
    }
    defer rows.Close() // ✅ garantido que vai fechar, mesmo se houver panic ou return antecipado

    categories := []Category{}
    for rows.Next() {
        // ...
    }
    return categories, nil
}
```

O `defer` garante que `rows.Close()` será chamado quando a função terminar, independentemente de como ela terminar.

---

## 🛡️ O que as Boas Práticas Evitaram

### Acoplamento Forte → Impossibilidade de Testar

Sem injeção de dependência, testar o comando `create` exigiria um arquivo `data.db` real, criado no diretório correto, com o schema correto. Com injeção, basta passar uma implementação do DAO que usa SQLite em memória (`:memory:`) — os testes rodam em milissegundos, sem efeitos colaterais.

### Erros Silenciosos que Quebram Pipelines

Um pipeline de CI/CD que chama `./cli category create` e não recebe código de saída 1 em caso de falha pode continuar executando etapas seguintes com dados inconsistentes. `RunE` garante que qualquer erro se propague corretamente para o ambiente de execução.

### Flags Inconsistentes → Dados Corrompidos

Sem `MarkFlagsRequiredTogether`, um usuário poderia criar uma categoria com nome mas sem descrição — e o registro iria ao banco com `description = ""` sem nenhum aviso. Validação declarativa garante integridade nos dados antes mesmo de chegar ao banco.

### Vazamento de Recursos → Degradação de Performance

Sem `defer rows.Close()`, cada chamada de `FindAll()` manteria um cursor do banco aberto. Em uma aplicação com muitas chamadas, isso esgotaria o pool de conexões disponíveis.

---

## 🚀 O que Poderia Ser Adicionado

### 1. Implementação Real do `list`

O arquivo `cmd/list.go` atual é um placeholder:

```go
// Estado atual — não faz nada útil
Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("list called")
},
```

A implementação real usaria o `FindAll()` que já existe no DAO:

```go
func newListCmd(categoryDb database.Category) *cobra.Command {
    return &cobra.Command{
        Use:   "list",
        Short: "Lista todas as categorias",
        RunE: func(cmd *cobra.Command, args []string) error {
            categories, err := categoryDb.FindAll()
            if err != nil {
                return err
            }
            for _, c := range categories {
                fmt.Printf("ID: %s | Nome: %s | Descrição: %s\n",
                    c.ID, c.Name, c.Description)
            }
            return nil
        },
    }
}
```

### 2. Comandos para `course`

O DAO de cursos (`internal/database/course.go`) já está implementado com `Create`, `FindAll`, `FindByCategoryID` e `Find`, mas nenhum comando CLI utiliza esses métodos. A estrutura seria:

```
category
├── create
└── list

course
├── create --name --description --category-id
└── list   --category-id (opcional, para filtrar)
```

### 3. Viper para Configuração

Hoje o caminho do banco está hardcoded em `GetDb()`. Com **Viper** (biblioteca irmã do Cobra, do mesmo autor), você poderia configurar isso via arquivo, variável de ambiente ou flag:

```bash
# Via env var
export DB_PATH=/var/lib/minha-app/data.db
./cli category list

# Via flag global
./cli --config /etc/minha-app.yaml category list

# Via arquivo de config (~/.minha-app.yaml)
db_path: /var/lib/minha-app/data.db
```

### 4. Flag `--format` para Output

Ferramentas profissionais de CLI permitem diferentes formatos de saída, facilitando integração com scripts:

```bash
./cli category list --format table   # tabela ASCII (padrão humano)
./cli category list --format json    # JSON (para scripts e APIs)
./cli category list --format csv     # CSV (para planilhas)
```

### 5. Migração Automática do Schema

Hoje o banco precisa ter as tabelas criadas manualmente. Um `PersistentPreRunE` no root command poderia criar as tabelas caso não existam:

```go
rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
    db := GetDb()
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS categories (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT NOT NULL
        );
        CREATE TABLE IF NOT EXISTS courses (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT NOT NULL,
            category_id TEXT NOT NULL,
            FOREIGN KEY (category_id) REFERENCES categories(id)
        );
    `)
    return err
}
```

### 6. Testes com Banco em Memória

SQLite suporta bancos em memória com a string de conexão `:memory:`:

```go
func TestCreate(t *testing.T) {
    db, _ := sql.Open("sqlite3", ":memory:")
    // cria schema...
    categoryDb := database.NewCategory(db)
    cmd := newCreateCmd(*categoryDb)
    cmd.SetArgs([]string{"--name", "Test", "--description", "Desc"})
    err := cmd.Execute()
    assert.NoError(t, err)
}
```

---

## ⚠️ Principais Problemas ao Trabalhar com Cobra CLI

### 1. Dependência de Ordem no `init()`

**O problema:** a ordem em que os `init()` de diferentes arquivos executam não é garantida em Go — exceto que cada `init()` espera todos os `init()` do mesmo pacote rodarem antes. Isso pode causar `nil pointer dereference` se um subcomando tentar se registrar em um `rootCmd` que ainda não foi inicializado.

```go
// ❌ Problema — se category.go init() rodar antes de root.go init()
// rootCmd pode ser nil quando categoryCmd tenta se registrar nele
func init() {
    rootCmd.AddCommand(categoryCmd) // panic se rootCmd for nil
}

// ✅ Solução — declarar rootCmd no nível do pacote garante inicialização antes
var rootCmd = &cobra.Command{...} // inicializa junto com o pacote

func init() {
    rootCmd.AddCommand(categoryCmd) // seguro: rootCmd já existe
}
```

Na prática, como `root.go` e `category.go` estão no mesmo pacote `cmd`, Go garante que todas as variáveis de nível de pacote são inicializadas antes de qualquer `init()` rodar. Mas ao trabalhar com pacotes separados, esse ordering pode ser uma armadilha.

### 2. Flags Globais Acidentais

**O problema:** usar `PersistentFlags()` quando deveria usar `Flags()` faz a flag aparecer em todos os subcomandos, poluindo o `--help` e criando confusão.

```go
// ❌ Flag global acidental
rootCmd.PersistentFlags().StringP("name", "n", "", "Nome")
// Aparece em TODOS os subcomandos, inclusive os que não usam --name

// ✅ Flag local — só aparece no comando correto
createCmd.Flags().StringP("name", "n", "", "Nome da categoria")
```

**Regra prática:** só use `PersistentFlags` para flags que fazem sentido em **todos** os subcomandos sem exceção (`--verbose`, `--config`, `--output`).

### 3. `panic` em vez de retornar erro

**O problema:** `panic` em Go encerra o processo abruptamente, sem dar chance de Cobra formatar o erro de forma amigável ou executar defers.

```go
// ❌ panic — termina o processo sem mensagem amigável
RunE: func(cmd *cobra.Command, args []string) error {
    db := GetDb()
    if db == nil {
        panic("banco não inicializado") // stacktrace assustador para o usuário
    }
    return nil
}

// ✅ retornar erro — Cobra formata e mostra o uso correto
RunE: func(cmd *cobra.Command, args []string) error {
    db := GetDb()
    if db == nil {
        return fmt.Errorf("banco de dados não disponível: verifique se data.db existe")
    }
    return nil
}
```

O `panic` em `GetDb()` do projeto (`panic(err)`) é um caso especialmente crítico: se o arquivo `data.db` não existir ou estiver corrompido, o programa termina com um stacktrace em vez de uma mensagem de erro clara.

### 4. Conexão com Banco Não Fechada

**O problema:** o projeto abre uma conexão com o banco em `GetDb()` mas nunca chama `db.Close()`. Para uma CLI que executa e termina rapidamente, o OS fecha os file descriptors automaticamente — mas é uma má prática que se torna problemática em aplicações de longa duração.

```go
// cmd/root.go — GetDb() abre, mas ninguém fecha
func GetDb() *sql.DB {
    db, _ := sql.Open("sqlite3", "./data.db")
    return db // quem vai chamar db.Close()?
}

// ✅ Solução — fechar a conexão após uso
func init() {
    db := GetDb()
    defer db.Close() // fecha quando init() terminar... mas init() termina cedo!

    // ✅ Melhor: fechar no final do comando
    createCmd := newCreateCmd(GetCategoryDB(db))
    createCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
        return db.Close()
    }
}
```

Na prática, para CLIs simples que executam uma ação e terminam, isso raramente causa problemas reais. Mas é importante entender a limitação.

### 5. Conflito de Nomes de Flags

**O problema:** se um comando pai e um filho definem flags com o mesmo nome, Cobra lança um erro em runtime.

```go
// ❌ Conflito de flags
categoryCmd.PersistentFlags().StringP("name", "n", "", "Nome")
createCmd.Flags().StringP("name", "n", "", "Nome da categoria") // CONFLITO!

// Erro: flag redefined: name

// ✅ Solução — use flags persistentes no pai OU locais no filho, nunca os dois
categoryCmd.PersistentFlags().StringP("name", "n", "", "Nome")
// createCmd herda --name automaticamente
```

**Atenção especial com `-h`/`--help`:** Cobra registra essa flag automaticamente em todos os comandos. Nunca tente registrar uma flag com esse nome.

### 6. `args` vs `flags` — quando usar cada um

Cobra suporta dois jeitos de passar dados para um comando: **flags** (com `--nome valor`) e **argumentos posicionais** (`args []string`).

```bash
# Flags — explícitas, autocompletion funciona, ordem não importa
./cli category create --name "Go" --description "Cursos"

# Args posicionais — implícitos, ordem importa
./cli category create "Go" "Cursos"
```

**Quando usar flags:**
- Quando há mais de 1-2 parâmetros
- Quando os parâmetros são opcionais
- Quando a ordem dos parâmetros pode confundir
- Quando quer validação automática com `MarkFlagsRequiredTogether`

**Quando usar args posicionais:**
- Comandos que operam sobre um "objeto" óbvio: `delete <id>`, `get <nome>`
- Quando a CLI tem que ser compatível com pipelines Unix: `cat arquivo | minha-cli processar`

```go
// Exemplo de uso correto de args posicionais
var deleteCmd = &cobra.Command{
    Use:   "delete <id>",
    Args:  cobra.ExactArgs(1),  // valida que exatamente 1 arg foi passado
    RunE: func(cmd *cobra.Command, args []string) error {
        id := args[0]
        return categoryDb.Delete(id)
    },
}
```

---

## 🔧 Como Usar o Projeto

### Pré-requisitos

```bash
# Go 1.19 ou superior
go version

# GCC (necessário para go-sqlite3 com CGO)
# macOS:
xcode-select --install

# Ubuntu/Debian:
sudo apt-get install build-essential

# Verificar instalação
gcc --version
```

### Criando o Schema do Banco

O projeto não cria as tabelas automaticamente. Execute este SQL antes de usar:

```bash
# Usando sqlite3 CLI
sqlite3 data.db <<EOF
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS courses (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    category_id TEXT NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
EOF
```

Ou use qualquer cliente SQLite de sua preferência (DB Browser for SQLite, TablePlus, etc.).

### Build e Execução

```bash
# Entrar na pasta do projeto
cd aulas/14-cobra-cli

# Baixar dependências
go mod tidy

# Compilar
go build -o cli .

# Verificar ajuda
./cli --help
./cli category --help
./cli category create --help
```

### Exemplos de Comandos

```bash
# Criar uma categoria
./cli category create --name "Programação" --description "Cursos de programação"
./cli category create -n "DevOps" -d "Infraestrutura e automação"

# Tentar criar sem um dos campos obrigatórios (erro esperado)
./cli category create --name "Incompleta"
# Error: if any flags in the group [name description] are set they must all be set; missing [description]

# Listar categorias (placeholder atual)
./cli category list
# list called

# Ver ajuda de qualquer comando
./cli --help
./cli category --help
./cli category create --help
```

---

## 📖 Glossário

| Termo | Definição |
|-------|-----------|
| **CLI** | Command-Line Interface — interface de texto onde o usuário digita comandos em vez de clicar em botões |
| **Flag** | Parâmetro nomeado passado a um comando (`--name "valor"` ou `-n "valor"`) |
| **Subcomando** | Comando filho de outro comando (`category create` — "create" é subcomando de "category") |
| **Root Command** | O comando base de uma CLI, o ponto de entrada quando o usuário digita o nome do binário |
| **DAO** | Data Access Object — camada de código cuja única responsabilidade é ler e escrever no banco de dados |
| **SQLite** | Banco de dados relacional que vive em um único arquivo, sem necessidade de servidor |
| **CGO** | Mecanismo do Go para chamar código escrito em C — necessário para o driver go-sqlite3 |
| **UUID** | Universally Unique Identifier — string de 36 caracteres globalmente única, usada como chave primária |
| **Injeção de Dependência** | Padrão onde um objeto recebe suas dependências de fora em vez de criá-las internamente |
| **defer** | Palavra-chave do Go que adia a execução de uma função até o final da função atual |
| **Código de Saída** | Número retornado por um processo ao terminar — 0 = sucesso, qualquer outro = erro |
| **PersistentPreRunE** | Hook do Cobra executado antes de qualquer subcomando, útil para validações globais |
| **Viper** | Biblioteca de configuração complementar ao Cobra — lê config de arquivos, env vars e flags |

---

## 🎯 Próximos Passos

### Para consolidar o aprendizado:

1. **Implemente o `list` de verdade** — use o `FindAll()` do DAO de categorias e exiba os resultados no terminal

2. **Adicione comandos para `course`** — o DAO já existe em `internal/database/course.go`, falta apenas criar `cmd/course.go` e os subcomandos `create` e `list`

3. **Substitua o `panic` por erro** em `GetDb()` — trate o erro de abertura do banco de forma amigável

4. **Adicione migração automática** — implemente um `PersistentPreRunE` no root que cria as tabelas se não existirem

5. **Explore o Viper** — adicione suporte a variável de ambiente `DB_PATH` para configurar o caminho do banco sem recompilar

6. **Escreva testes** — use SQLite em memória (`:memory:`) para testar os DAOs sem criar arquivos temporários

### Conceitos relacionados no curso:

- **Aula 13 (Upload S3)** — padrões de concorrência e workers pools que podem ser usados em processamentos batch via CLI
- **Aula 9 (Eventos)** — sistema de eventos que pode ser integrado a uma CLI para processamento assíncrono
- **Aula 7 (APIs)** — uma CLI pode ser um cliente para a API REST que você construiu, usando os mesmos domínios
