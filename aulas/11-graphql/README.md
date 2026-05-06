# GraphQL em Go — Guia Didático Completo

Este módulo ensina como construir APIs GraphQL em Go usando a biblioteca `gqlgen`. Você aprenderá a definir um schema, implementar resolvers, conectar um banco de dados e expor uma API flexível que permite ao cliente escolher exatamente quais dados quer receber.

---

## 📑 Sumário

- [🤔 O que é GraphQL?](#-o-que-é-graphql)
  - [A Analogia do Restaurante](#a-analogia-do-restaurante)
  - [Como o GraphQL Funciona na Prática](#como-o-graphql-funciona-na-prática)
- [⚔️ REST vs GraphQL — O Grande Comparativo](#️-rest-vs-graphql--o-grande-comparativo)
  - [O Problema do Overfetching](#o-problema-do-overfetching)
  - [O Problema do Underfetching](#o-problema-do-underfetching)
  - [Mesma Operação: REST vs GraphQL](#mesma-operação-rest-vs-graphql)
  - [Tabela Comparativa](#tabela-comparativa)
- [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
  - [Schema — O Contrato da API](#schema--o-contrato-da-api)
  - [Types — Os Tipos de Dados](#types--os-tipos-de-dados)
  - [Query — Consultando Dados](#query--consultando-dados)
  - [Mutation — Modificando Dados](#mutation--modificando-dados)
  - [Resolver — A Lógica que Resolve Cada Campo](#resolver--a-lógica-que-resolve-cada-campo)
  - [Field Resolver — Resolvendo Relações](#field-resolver--resolvendo-relações)
  - [O Problema N+1](#o-problema-n1)
- [🛠️ A Ferramenta gqlgen](#️-a-ferramenta-gqlgen)
  - [O que é e Por que Usar](#o-que-é-e-por-que-usar)
  - [O Ciclo de Geração de Código](#o-ciclo-de-geração-de-código)
  - [O Arquivo gqlgen.yml Explicado](#o-arquivo-gqlyenyml-explicado)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
  - [O Schema (graph/schema.graphqls)](#o-schema-graphschemagraphqls)
  - [Os Modelos Go (graph/model/)](#os-modelos-go-graphmodel)
  - [O Resolver Raiz (graph/resolver.go)](#o-resolver-raiz-graphresolvergo)
  - [Os Resolvers em Ação (graph/schema.resolvers.go)](#os-resolvers-em-ação-graphschemaresolversgo)
  - [A Camada de Banco (internal/database/)](#a-camada-de-banco-internaldatabase)
  - [O Servidor (cmd/server/server.go)](#o-servidor-cmdserverservergo)
- [▶️ Como Executar](#️-como-executar)
- [🎮 GraphQL Playground](#-graphql-playground)
  - [Criando Dados com Mutations](#criando-dados-com-mutations)
  - [Consultando Dados com Queries](#consultando-dados-com-queries)
  - [Queries com Campos Aninhados](#queries-com-campos-aninhados)
- [⚖️ Trade-offs: GraphQL vs REST](#️-trade-offs-graphql-vs-rest)
  - [Vantagens do GraphQL](#vantagens-do-graphql)
  - [Desvantagens do GraphQL](#desvantagens-do-graphql)
  - [Quando Escolher GraphQL](#quando-escolher-graphql)
  - [Quando Escolher REST](#quando-escolher-rest)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é GraphQL?

**GraphQL** é uma linguagem de consulta para APIs criada pelo Facebook em 2012 e aberta ao público em 2015. Em vez de ter vários endpoints fixos (como no REST), você tem **um único endpoint** e o cliente decide exatamente quais dados quer receber.

### A Analogia do Restaurante

Imagine dois tipos de restaurante:

**Restaurante REST (Prato Fixo):**
```
Você pede o "Prato A" → recebe: filé, arroz, feijão, farofa, salada, suco
Você só queria o filé e o arroz, mas veio tudo junto.
Se quiser a bebida separada, precisa fazer um segundo pedido.
```

**Restaurante GraphQL (Cardápio à la Carte):**
```
Você olha o cardápio e diz: "Quero filé e arroz, só isso."
→ recebe: filé e arroz.

Na mesma ordem, você pode dizer: "Quero filé, arroz e o nome do chef."
→ recebe exatamente isso.
```

No GraphQL, **o cliente tem o controle** do que vem na resposta. Nada a mais, nada a menos.

### Como o GraphQL Funciona na Prática

```
                    ┌─────────────────────────────────┐
                    │        CLIENTE (App/Browser)     │
                    │                                  │
                    │  query {                         │
                    │    categories {                  │
                    │      id                          │
                    │      name                        │  ← Só quer esses 2 campos
                    │    }                             │
                    │  }                               │
                    └──────────────┬──────────────────┘
                                   │ POST /query
                                   ▼
                    ┌─────────────────────────────────┐
                    │       SERVIDOR GRAPHQL           │
                    │                                  │
                    │  1. Valida a query contra schema │
                    │  2. Executa os resolvers certos  │
                    │  3. Monta a resposta              │
                    └──────────────┬──────────────────┘
                                   │ JSON
                                   ▼
                    ┌─────────────────────────────────┐
                    │           RESPOSTA               │
                    │                                  │
                    │  {                               │
                    │    "data": {                     │
                    │      "categories": [             │
                    │        { "id": "1",              │
                    │          "name": "Backend" }     │
                    │      ]                           │
                    │    }                             │
                    │  }                               │
                    └─────────────────────────────────┘
```

---

## ⚔️ REST vs GraphQL — O Grande Comparativo

Para entender o valor do GraphQL, é fundamental entender os problemas que ele resolve.

### O Problema do Overfetching

**Overfetching** = receber mais dados do que você precisa.

Imagine um app mobile que só precisa exibir o **nome** das categorias numa lista. Com REST:

```
GET /categories

Resposta (você pediu nome, mas veio tudo):
[
  {
    "id": "abc-123",
    "name": "Backend",
    "description": "Cursos de desenvolvimento backend...",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-20T14:30:00Z",
    "total_courses": 42,
    "is_active": true
  },
  ...
]
```

Você precisava só do `name`, mas recebeu 7 campos. Em dispositivos móveis com conexão lenta, isso importa muito.

Com GraphQL:
```graphql
query {
  categories {
    name   # ← Só isso. Ponto.
  }
}

Resposta:
{
  "data": {
    "categories": [
      { "name": "Backend" },
      { "name": "Frontend" }
    ]
  }
}
```

### O Problema do Underfetching

**Underfetching** = precisar de múltiplas requisições para montar uma tela.

Imagine montar uma tela que mostra categorias e seus cursos. Com REST:

```
Requisição 1: GET /categories
→ [ { "id": "1", "name": "Backend" }, { "id": "2", "name": "Frontend" } ]

Requisição 2: GET /categories/1/courses
→ [ { "id": "10", "name": "Go Avançado" }, ... ]

Requisição 3: GET /categories/2/courses
→ [ { "id": "20", "name": "React" }, ... ]
```

São **3 requisições** para montar uma tela. Com GraphQL, é **1 só**:

```graphql
query {
  categories {
    name
    courses {
      name
    }
  }
}
```

### Mesma Operação: REST vs GraphQL

| Aspecto | REST | GraphQL |
|--------|------|---------|
| Endpoint para criar categoria | `POST /categories` | `mutation { createCategory(...) }` em `/query` |
| Endpoint para buscar cursos de uma categoria | `GET /categories/:id/courses` | campo `courses` dentro de `category` |
| Número de endpoints | 1 por recurso/ação | 1 único (`/query`) |
| Controle dos campos retornados | Nenhum (servidor decide) | Total (cliente decide) |
| Documentação integrada | Precisa de Swagger/OpenAPI | Introspection nativa |

### Tabela Comparativa

| Característica | REST | GraphQL |
|---------------|------|---------|
| **Número de endpoints** | Múltiplos | Um único |
| **Formato da requisição** | URL + verbo HTTP | Linguagem de query |
| **Controle pelo cliente** | Limitado | Total |
| **Overfetching** | Comum | Inexistente |
| **Underfetching** | Comum | Inexistente |
| **Cache HTTP nativo** | ✅ Sim (GET é cacheável) | ❌ Difícil (tudo é POST) |
| **Curva de aprendizado** | Baixa | Média-alta |
| **Ferramentas de exploração** | Swagger, Postman | Playground integrado |
| **Versioning de API** | `/v1/`, `/v2/` | Evolução sem versionamento |
| **Upload de arquivos** | ✅ Simples | ❌ Mais complexo |
| **Adoção no mercado** | Muito alta | Crescente |

---

## 📚 Conceitos Fundamentais

### Schema — O Contrato da API

O **schema** é a peça central do GraphQL. É um arquivo que define **tudo que a API pode fazer**: quais tipos existem, quais campos cada tipo tem, quais queries são possíveis e quais mutations existem.

Pense no schema como **o contrato entre servidor e cliente**. O servidor promete: "Tenho esses dados, com esses tipos, acessíveis dessa forma." O cliente sabe exatamente o que pode pedir.

```graphql
# schema.graphqls — Este arquivo é a "planta da API"

# Tipo Category: representa uma categoria de cursos
type Category {
  id: ID!           # ID único. O "!" significa obrigatório (nunca null)
  name: String!     # Nome da categoria. Sempre presente
  description: String  # Descrição. Sem "!" = pode ser null (opcional)
  courses: [Course!]!  # Lista de cursos. A lista em si nunca é null,
                       # e cada item dentro também não é null
}
```

### Types — Os Tipos de Dados

O GraphQL tem **tipos escalares** (valores simples) e **tipos de objeto** (structs):

**Tipos Escalares:**
| Tipo | Descrição | Exemplo |
|------|-----------|---------|
| `String` | Texto | `"GraphQL é incrível"` |
| `Int` | Número inteiro | `42` |
| `Float` | Número decimal | `3.14` |
| `Boolean` | Verdadeiro/falso | `true` |
| `ID` | Identificador único (string especial) | `"abc-123"` |

**A Exclamação `!` é Importante:**
```graphql
name: String    # Pode ser null  → campo opcional
name: String!   # Nunca é null   → campo obrigatório
courses: [Course!]!   # A lista nunca é null E cada curso dentro também não
```

**Tipos de Entrada (input):**

Mutations precisam receber dados. Para isso, usamos `input` (diferente de `type`):

```graphql
# "input" é usado para RECEBER dados (parâmetros)
input NewCategory {
  name: String!        # Nome é obrigatório na criação
  description: String  # Descrição é opcional
}

# "type" é usado para RETORNAR dados (respostas)
type Category {
  id: ID!
  name: String!
  description: String
  courses: [Course!]!
}
```

### Query — Consultando Dados

`Query` é o tipo especial que define **todas as formas de leitura** da API:

```graphql
type Query {
  categories: [Category!]!   # Retorna todos as categorias
  courses: [Course!]!        # Retorna todos os cursos
}
```

O cliente usa assim:
```graphql
# Buscar só id e name de todas as categorias
query {
  categories {
    id
    name
  }
}

# Buscar cursos com seus campos aninhados
query {
  courses {
    id
    name
    category {
      name   # Campo do tipo relacionado!
    }
  }
}
```

### Mutation — Modificando Dados

`Mutation` é o tipo especial para **criar, atualizar ou deletar** dados:

```graphql
type Mutation {
  createCategory(input: NewCategory!): Category!
  createCourse(input: NewCourse!): Course!
}
```

O cliente usa assim:
```graphql
# Criar uma categoria e receber de volta id e name
mutation {
  createCategory(input: {
    name: "Backend"
    description: "Desenvolvimento de APIs e serviços"
  }) {
    id      # ← Você ainda escolhe quais campos quer de volta!
    name
  }
}
```

### Resolver — A Lógica que Resolve Cada Campo

Cada campo de um tipo pode ter um **resolver** — uma função Go que sabe como buscar aquele dado.

Imagine assim:

```
Query { categories }
         │
         ▼
  queryResolver.Categories()   ← Resolver que vai ao banco buscar todas as categorias
         │
         ▼
  [ Category1, Category2 ]
```

Para cada categoria retornada, se o cliente pediu o campo `courses`, o GraphQL chama outro resolver:

```
Category.courses
    │
    ▼
categoryResolver.Courses(category)   ← Resolver que busca cursos desta categoria específica
    │
    ▼
[ Course1, Course2 ]
```

**Resolvers em Go seguem uma estrutura de interfaces:**

```go
// Interface gerada pelo gqlgen — você precisa implementar isso
type QueryResolver interface {
    Categories(ctx context.Context) ([]*model.Category, error)
    Courses(ctx context.Context) ([]*model.Course, error)
}

type MutationResolver interface {
    CreateCategory(ctx context.Context, input model.NewCategory) (*model.Category, error)
    CreateCourse(ctx context.Context, input model.NewCourse) (*model.Course, error)
}
```

### Field Resolver — Resolvendo Relações

Alguns campos são **relações** (um `Category` tem vários `Course`). O GraphQL resolve isso com **field resolvers** — funções específicas para campos relacionados:

```go
// Chamado quando o cliente pede category.courses
func (r *categoryResolver) Courses(ctx context.Context, obj *model.Category) ([]*model.Course, error) {
    // obj é a Category que já foi resolvida
    // Agora buscamos os cursos DESTA categoria específica
    return r.CourseDB.FindByCategoryID(obj.ID)
}

// Chamado quando o cliente pede course.category
func (r *courseResolver) Category(ctx context.Context, obj *model.Course) (*model.Category, error) {
    // obj é o Course já resolvido
    // Buscamos a categoria DESTE curso
    return r.CategoryDB.FindByCourseID(obj.ID)
}
```

### O Problema N+1

Field resolvers trazem um problema clássico de performance chamado **N+1**:

```
Query: { categories { courses { name } } }

1 query para buscar todas as categorias → retorna 10 categorias
10 queries para buscar os cursos de cada categoria → 1 por categoria

Total: 11 queries ao banco de dados!
```

```
           ┌─ FindAll categories ──────────────────── 1 query
           │
           ├─ FindByCategoryID("cat-1") ─────────────┐
           ├─ FindByCategoryID("cat-2") ─────────────┤
           ├─ FindByCategoryID("cat-3") ─────────────┤ N queries
           ├─ FindByCategoryID("cat-4") ─────────────┤ (uma por categoria)
           └─ FindByCategoryID("cat-N") ─────────────┘

Total = 1 + N queries
```

A solução para isso é o padrão **DataLoader** (carregador em lote), que agrupa as N queries em uma só. Nesta aula, o problema existe — e reconhecê-lo já é o primeiro passo para resolvê-lo em projetos reais.

---

## 🛠️ A Ferramenta gqlgen

### O que é e Por que Usar

**gqlgen** (`github.com/99designs/gqlgen`) é um gerador de código Go para GraphQL. Em vez de você escrever código de parsing, validação e execução de queries do zero, o gqlgen **gera tudo isso automaticamente** a partir do seu schema.

```
Você escreve:        gqlgen gera:              Você implementa:
─────────────        ─────────────────────     ─────────────────
schema.graphqls  →   graph/generated/      →   schema.resolvers.go
                     generated.go               (só a lógica de negócio)
```

**Por que não escrever na mão?**

O código de execução de queries GraphQL é extremamente complexo: precisa fazer parsing da query, validar campos contra o schema, iterar em árvore pelos campos pedidos, chamar o resolver certo para cada campo, e montar a resposta JSON. O arquivo gerado `generated.go` tem mais de 3000 linhas — nenhum desenvolvedor quer escrever isso.

### O Ciclo de Geração de Código

```
┌─────────────────┐
│  schema.graphqls │  ← Você define os tipos e operações
└────────┬────────┘
         │ go generate ./...
         ▼
┌─────────────────────────────────────┐
│  graph/generated/generated.go       │  ← gqlgen gera automaticamente
│  graph/model/models_gen.go          │     (NUNCA edite esses arquivos)
└─────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  graph/schema.resolvers.go          │  ← Você implementa a lógica
│  (interfaces definidas pelo gqlgen) │
└─────────────────────────────────────┘
```

**Regra de ouro:**

✅ Edite `schema.graphqls` → rode `go generate` → implemente em `schema.resolvers.go`

❌ Nunca edite `graph/generated/generated.go` (será sobrescrito na próxima geração)

### O Arquivo gqlgen.yml Explicado

```yaml
# Onde estão os arquivos de schema GraphQL
schema:
  - graph/*.graphqls         # Todos os .graphqls dentro de graph/

# Onde salvar o código gerado do servidor
exec:
  filename: graph/generated/generated.go   # Arquivo gerado (não toque!)
  package: generated

# Onde salvar os modelos gerados
model:
  filename: graph/model/models_gen.go      # Tipos de input gerados automaticamente
  package: model

# Onde salvar as implementações dos resolvers
resolver:
  layout: follow-schema    # Um arquivo .resolvers.go por arquivo .graphqls
  dir: graph
  package: graph

# Mapeamento: tipos GraphQL → structs Go
models:
  Category:
    model:
      - github.com/devfullcycle/13-GraphQL/graph/model.Category  # Usa nossa struct
  Course:
    model:
      - github.com/devfullcycle/13-GraphQL/graph/model.Course    # Usa nossa struct
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID    # ID pode ser string ou int
```

---

## 🗂️ Estrutura do Projeto

```
aulas/11-graphql/
│
├── cmd/
│   └── server/
│       └── server.go           # Ponto de entrada: inicia o servidor HTTP
│
├── graph/
│   ├── generated/
│   │   └── generated.go        # ⚠️ GERADO pelo gqlgen — não edite!
│   │
│   ├── model/
│   │   ├── category.go         # Struct Category (você escreveu)
│   │   ├── course.go           # Struct Course (você escreveu)
│   │   └── models_gen.go       # ⚠️ Inputs gerados pelo gqlgen (NewCategory, NewCourse)
│   │
│   ├── resolver.go             # Struct Resolver com injeção de dependência
│   ├── schema.graphqls         # O schema GraphQL — coração da API
│   └── schema.resolvers.go     # Implementação dos resolvers (você escreve aqui)
│
├── internal/
│   └── database/
│       ├── category.go         # Acesso ao banco para categorias
│       └── course.go           # Acesso ao banco para cursos
│
├── gqlgen.yml                  # Configuração do gerador de código
├── tools.go                    # Import do gqlgen para go generate
├── go.mod                      # Dependências do módulo
└── go.sum                      # Checksums das dependências
```

**Separação de responsabilidades:**

```
schema.graphqls         → O QUÊ a API expõe (contrato)
schema.resolvers.go     → COMO buscar cada dado (lógica)
internal/database/      → ONDE os dados estão (persistência)
cmd/server/server.go    → COMO tudo se conecta (infraestrutura)
```

---

## 🔍 Walkthrough do Código

### O Schema (`graph/schema.graphqls`)

```graphql
# Tipo principal: Categoria de cursos
type Category {
  id: ID!              # Identificador único (UUID gerado pelo servidor)
  name: String!        # Nome é obrigatório
  description: String  # Descrição é opcional (pode ser null)
  courses: [Course!]!  # Relação: lista de cursos desta categoria
                       # Field resolver vai buscar isso sob demanda
}

# Tipo principal: Curso
type Course {
  id: ID!
  name: String!
  description: String
  category: Category!  # Relação inversa: cada curso pertence a uma categoria
                       # Field resolver vai buscar isso sob demanda
}

# Tipos de entrada (para criar novos registros)
input NewCategory {
  name: String!
  description: String
}

input NewCourse {
  name: String!
  description: String
  categoryId: ID!   # Precisa informar a qual categoria o curso pertence
}

# Operações de leitura
type Query {
  categories: [Category!]!   # Lista todas as categorias
  courses: [Course!]!        # Lista todos os cursos
}

# Operações de escrita
type Mutation {
  createCategory(input: NewCategory!): Category!   # Cria e retorna a categoria criada
  createCourse(input: NewCourse!): Course!         # Cria e retorna o curso criado
}
```

**Por que `Category` tem `courses` e `Course` tem `category`?**

Isso é a **relação bidirecional** do GraphQL. Em REST você teria dois endpoints separados. No GraphQL, o cliente navega pelo grafo de dados como quiser:

```graphql
# Navegar de Category para Course:
{ categories { courses { name } } }

# Navegar de Course para Category:
{ courses { category { name } } }
```

### Os Modelos Go (`graph/model/`)

Os modelos são as structs Go que representam os dados em memória:

**`graph/model/category.go`** (escrito manualmente):
```go
package model

type Category struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description *string `json:"description"` // Ponteiro porque pode ser null no GraphQL
}
```

**Por que `*string` e não `string`?**

No schema, `description: String` (sem `!`) significa que pode ser `null`. Em Go, um `string` nunca é nulo — ele é vazio (`""`). Para representar `null` de verdade, usamos um **ponteiro** (`*string`):

```
GraphQL null  →  Go: *string = nil
GraphQL ""    →  Go: *string = &""
GraphQL "abc" →  Go: *string = &"abc"
```

**`graph/model/models_gen.go`** (gerado pelo gqlgen):
```go
// ARQUIVO GERADO — NÃO EDITE

// NewCategory é o tipo de input para createCategory
type NewCategory struct {
    Name        string  `json:"name"`
    Description *string `json:"description"`
}

// NewCourse é o tipo de input para createCourse
type NewCourse struct {
    Name        string  `json:"name"`
    Description *string `json:"description"`
    CategoryID  string  `json:"categoryId"`
}
```

### O Resolver Raiz (`graph/resolver.go`)

```go
package graph

import "github.com/devfullcycle/13-GraphQL/internal/database"

// Resolver é o ponto central de injeção de dependência.
// Todos os resolvers específicos vão embutir (embed) esta struct
// e terão acesso ao CategoryDB e CourseDB.
type Resolver struct {
    CategoryDB *database.Category  // Acesso ao banco de categorias
    CourseDB   *database.Course    // Acesso ao banco de cursos
}
```

**Padrão de injeção de dependência:**

```
Resolver (raiz)
    ├── CategoryDB ──→ banco de dados
    └── CourseDB   ──→ banco de dados

categoryResolver { *Resolver }   ← embutiu o Resolver, herda CategoryDB e CourseDB
courseResolver   { *Resolver }
mutationResolver { *Resolver }
queryResolver    { *Resolver }
```

Cada resolver específico é apenas um wrapper que aponta para o `Resolver` raiz. Assim, todos compartilham as mesmas conexões com o banco.

### Os Resolvers em Ação (`graph/schema.resolvers.go`)

**Resolver de Query — Categories:**
```go
// Chamado quando o cliente faz: query { categories { ... } }
func (r *queryResolver) Categories(ctx context.Context) ([]*model.Category, error) {
    // 1. Busca todas as categorias no banco
    categories, err := r.CategoryDB.FindAll()
    if err != nil {
        return nil, err   // Retorna o erro para o cliente GraphQL
    }

    // 2. Converte do tipo banco (database.Category) para o tipo GraphQL (model.Category)
    var categoriesModel []*model.Category
    for _, category := range categories {
        categoriesModel = append(categoriesModel, &model.Category{
            ID:          category.ID,
            Name:        category.Name,
            Description: &category.Description,  // string → *string
        })
    }
    return categoriesModel, nil
}
```

**Resolver de Mutation — CreateCategory:**
```go
// Chamado quando o cliente faz: mutation { createCategory(input: {...}) { ... } }
func (r *mutationResolver) CreateCategory(ctx context.Context, input model.NewCategory) (*model.Category, error) {
    // 1. Persiste no banco — o banco gera o UUID automaticamente
    category, err := r.CategoryDB.Create(input.Name, *input.Description)
    if err != nil {
        return nil, err
    }

    // 2. Retorna o objeto criado (com o ID gerado)
    return &model.Category{
        ID:          category.ID,
        Name:        category.Name,
        Description: &category.Description,
    }, nil
}
```

**Field Resolver — Courses de uma Category:**
```go
// Chamado APENAS quando o cliente pede o campo "courses" dentro de uma category
// Ex: query { categories { courses { name } } }
// obj é a Category já resolvida — usamos seu ID para buscar os cursos
func (r *categoryResolver) Courses(ctx context.Context, obj *model.Category) ([]*model.Course, error) {
    courses, err := r.CourseDB.FindByCategoryID(obj.ID)  // ← Busca só desta categoria
    if err != nil {
        return nil, err
    }

    var coursesModel []*model.Course
    for _, course := range courses {
        coursesModel = append(coursesModel, &model.Course{
            ID:          course.ID,
            Name:        course.Name,
            Description: &course.Description,
        })
    }
    return coursesModel, nil
}
```

**Como os resolvers específicos são conectados:**
```go
// No final do arquivo — conecta cada resolver ao Resolver raiz
func (r *Resolver) Category() generated.CategoryResolver { return &categoryResolver{r} }
func (r *Resolver) Course()   generated.CourseResolver   { return &courseResolver{r} }
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query()    generated.QueryResolver    { return &queryResolver{r} }

// Cada tipo de resolver é apenas um struct que embute o Resolver raiz
type categoryResolver struct{ *Resolver }
type courseResolver   struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver    struct{ *Resolver }
```

### A Camada de Banco (`internal/database/`)

A camada de banco é **independente do GraphQL** — são structs Go simples que falam com SQLite.

**`internal/database/category.go`:**
```go
type Category struct {
    db          *sql.DB  // Conexão com o banco (privada, uso interno)
    ID          string
    Name        string
    Description string
}

// Construtor — recebe a conexão e retorna o repositório
func NewCategory(db *sql.DB) *Category {
    return &Category{db: db}
}

// Create insere uma nova categoria e retorna ela com o ID gerado
func (c *Category) Create(name string, description string) (Category, error) {
    id := uuid.New().String()  // Gera UUID antes de inserir
    _, err := c.db.Exec(
        "INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)",
        id, name, description,
    )
    if err != nil {
        return Category{}, err
    }
    return Category{ID: id, Name: name, Description: description}, nil
}

// FindAll retorna todas as categorias do banco
func (c *Category) FindAll() ([]Category, error) {
    rows, err := c.db.Query("SELECT id, name, description FROM categories")
    // ...
}

// FindByCourseID busca a categoria de um curso específico (usado pelo field resolver)
func (c *Category) FindByCourseID(courseID string) (Category, error) {
    err := c.db.QueryRow(`
        SELECT c.id, c.name, c.description
        FROM categories c
        JOIN courses co ON c.id = co.category_id
        WHERE co.id = $1
    `, courseID).Scan(&id, &name, &description)
    // ...
}
```

**`internal/database/course.go`:**
```go
// FindByCategoryID retorna todos os cursos de uma categoria (usado pelo field resolver)
func (c *Course) FindByCategoryID(categoryID string) ([]Course, error) {
    rows, err := c.db.Query(
        "SELECT id, name, description, category_id FROM courses WHERE category_id = $1",
        categoryID,
    )
    // ...
}
```

### O Servidor (`cmd/server/server.go`)

```go
func main() {
    // 1. Abre o banco SQLite (cria o arquivo data.db se não existir)
    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        log.Fatalf("failed to open database: %v", err)
    }
    defer db.Close()  // Fecha a conexão quando main() terminar

    // 2. Cria os repositórios de banco (injeção de dependência)
    categoryDb := database.NewCategory(db)
    courseDb   := database.NewCourse(db)

    // 3. Pega a porta do ambiente (padrão: 8080)
    port := os.Getenv("PORT")
    if port == "" {
        port = defaultPort
    }

    // 4. Cria o servidor GraphQL com os resolvers
    srv := handler.NewDefaultServer(
        generated.NewExecutableSchema(generated.Config{
            Resolvers: &graph.Resolver{
                CategoryDB: categoryDb,  // Injeta o repositório de categorias
                CourseDB:   courseDb,    // Injeta o repositório de cursos
            },
        }),
    )

    // 5. Registra as rotas HTTP
    http.Handle("/", playground.Handler("GraphQL playground", "/query"))  // UI visual
    http.Handle("/query", srv)                                            // Endpoint da API

    // 6. Sobe o servidor
    log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

**O fluxo completo:**

```
Cliente HTTP
    │
    │ POST /query  { query: "{ categories { name } }" }
    ▼
http.Handle("/query", srv)
    │
    ▼
handler.NewDefaultServer (gqlgen)
    │ Faz parsing da query
    │ Valida contra o schema
    │ Chama os resolvers certos
    ▼
graph.Resolver
    │
    ├── queryResolver.Categories()
    │       └── categoryDb.FindAll()
    │               └── SQLite: SELECT id, name, description FROM categories
    │
    └── [field resolver se pediu courses]
            └── categoryResolver.Courses(category)
                    └── courseDb.FindByCategoryID(category.ID)
```

---

## ▶️ Como Executar

**Pré-requisitos:**
- Go 1.19+
- GCC instalado (necessário para o driver SQLite — `go-sqlite3` usa CGO)

**No macOS:**
```bash
# GCC já vem com Xcode Command Line Tools
xcode-select --install
```

**No Linux (Ubuntu/Debian):**
```bash
sudo apt-get install gcc
```

**Passos:**

```bash
# 1. Entre na pasta da aula
cd aulas/11-graphql

# 2. Instale as dependências
go mod download

# 3. (Opcional) Regenere o código se modificou o schema
go generate ./...

# 4. Execute o servidor
go run cmd/server/server.go

# Saída esperada:
# 2024/01/15 10:00:00 connect to http://localhost:8080/ for GraphQL playground
```

**Criando as tabelas no banco:**

O projeto usa SQLite mas não tem migrations automáticas. Antes de usar, crie as tabelas:

```bash
# Cria o banco e as tabelas manualmente
sqlite3 data.db "
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS courses (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    category_id TEXT NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
"
```

---

## 🎮 GraphQL Playground

O **GraphQL Playground** é uma interface visual interativa que o gqlgen disponibiliza automaticamente em `http://localhost:8080/`. É equivalente ao Swagger, mas para GraphQL — e muito mais poderoso.

```
http://localhost:8080/
    │
    ├── Painel esquerdo: Escreva queries e mutations
    ├── Painel direito:  Veja a resposta JSON
    ├── Botão "DOCS":    Documentação gerada do schema
    └── Botão "SCHEMA":  O schema completo
```

### Criando Dados com Mutations

**Criar uma categoria:**
```graphql
mutation {
  createCategory(input: {
    name: "Backend"
    description: "Desenvolvimento de APIs e serviços"
  }) {
    id
    name
    description
  }
}
```

Resposta:
```json
{
  "data": {
    "createCategory": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Backend",
      "description": "Desenvolvimento de APIs e serviços"
    }
  }
}
```

**Criar um curso (use o ID da categoria acima):**
```graphql
mutation {
  createCourse(input: {
    name: "Go Avançado"
    description: "Concorrência, módulos e microserviços em Go"
    categoryId: "550e8400-e29b-41d4-a716-446655440000"
  }) {
    id
    name
  }
}
```

### Consultando Dados com Queries

**Listar todas as categorias (só id e name):**
```graphql
query {
  categories {
    id
    name
  }
}
```

**Listar todos os cursos:**
```graphql
query {
  courses {
    id
    name
    description
  }
}
```

### Queries com Campos Aninhados

Aqui o GraphQL brilha — você pode navegar por relações numa única query:

**Categorias com seus cursos:**
```graphql
query {
  categories {
    id
    name
    description
    courses {          # ← Campo relacionado! Dispara o field resolver
      id
      name
    }
  }
}
```

Resposta:
```json
{
  "data": {
    "categories": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Backend",
        "description": "Desenvolvimento de APIs e serviços",
        "courses": [
          {
            "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
            "name": "Go Avançado"
          }
        ]
      }
    ]
  }
}
```

**Cursos com a categoria de cada um:**
```graphql
query {
  courses {
    name
    category {        # ← Navega no sentido inverso!
      name
    }
  }
}
```

**Comparativo: REST precisaria de 3 requisições, GraphQL faz em 1!**

```
REST:
  GET /categories          → [ {id, name, desc} ]
  GET /categories/1/courses → [ {id, name} ]
  GET /categories/2/courses → [ {id, name} ]

GraphQL:
  POST /query com { categories { name courses { name } } }
  → Tudo em uma requisição
```

---

## ⚖️ Trade-offs: GraphQL vs REST

Nenhuma tecnologia é melhor em todos os contextos. Conhecer os trade-offs é o que diferencia um desenvolvedor sênior.

### Vantagens do GraphQL

✅ **Sem overfetching/underfetching**: O cliente recebe exatamente o que pediu

✅ **Uma única requisição**: Dados relacionados chegam juntos, sem múltiplas chamadas

✅ **Evolução sem versionamento**: Adicionar campos ao schema não quebra clientes antigos. Nunca mais `/v1/` e `/v2/`

✅ **Documentação automática**: O Playground gera docs a partir do schema. Não precisa manter Swagger separado

✅ **Tipagem forte**: O schema é um contrato formal. Erros de tipo são detectados antes de chegar ao servidor

✅ **Ótimo para múltiplos clientes**: Web, mobile, TV — cada um pede só o que precisa

✅ **Introspection**: Clientes podem perguntar ao servidor quais operações existem programaticamente

### Desvantagens do GraphQL

❌ **Cache HTTP difícil**: REST usa cache nativo do HTTP (GET é cacheável por padrão). GraphQL usa POST para tudo, então o cache de rede não funciona sem bibliotecas extras (ex: `apollo-client`)

❌ **Curva de aprendizado**: Schema, resolvers, field resolvers, N+1, DataLoader — são muitos conceitos novos

❌ **Complexidade do servidor**: O servidor precisa executar queries arbitrárias. Um cliente pode fazer uma query muito cara sem que você perceba. É necessário implementar complexity limits

❌ **Upload de arquivos**: REST lida com multipart/form-data nativamente. GraphQL precisa de extensões específicas

❌ **Ferramentas de debug mais complexas**: Logs e rastreamento são mais difíceis porque cada campo tem seu próprio resolver

❌ **N+1 problem**: Field resolvers criam queries extras ao banco se não usar DataLoader

❌ **Overhead para APIs simples**: Se sua API tem 3 endpoints simples e nenhum dado relacionado, o GraphQL adiciona complexidade sem benefício

### Quando Escolher GraphQL

| Cenário | Por quê GraphQL |
|---------|-----------------|
| App mobile (iOS/Android) | Economia de banda — pede só o que precisa |
| Dashboard com dados de múltiplas fontes | Uma query busca tudo junto |
| API consumida por múltiplos clientes diferentes | Web quer mais campos, mobile quer menos |
| Dados com muitas relações | Navegação em grafo é natural |
| Equipe de front independente do back | Schema é contrato versionado |
| Prototipagem rápida com schema-first | Define o contrato antes de implementar |

### Quando Escolher REST

| Cenário | Por quê REST |
|---------|-------------|
| API pública para terceiros | REST é mais familiar e amplamente documentado |
| Operações simples de CRUD sem relações | Complexidade do GraphQL não compensa |
| Cache HTTP é crítico | GET requests são cacheáveis nativamente |
| Upload de arquivos | REST lida com multipart nativo |
| Integração com sistemas legados | REST tem mais bibliotecas e suporte |
| Equipe pequena sem experiência em GraphQL | Curva de aprendizado tem custo |
| Webhooks e notificações push | REST/HTTP é o padrão aqui |

---

## 🎯 Casos de Uso Ideais

### GraphQL faz sentido para:

**1. Aplicativos Mobile**
```
App iOS quer: id, name, thumbnail (imagem pequena)
App Android quer: id, name, description, thumbnail
App Web quer: id, name, description, thumbnail, category, author

Com REST:  3 endpoints diferentes OU overfetching em todos
Com GraphQL: 1 schema, cada app pede o que precisa
```

**2. Dashboards e Painéis Analíticos**
```
Um dashboard pode precisar de:
- Número de cursos por categoria
- Últimos cursos criados
- Categorias mais acessadas

Com REST: 3 chamadas separadas
Com GraphQL: 1 query busca os 3 conjuntos de dados
```

**3. BFF (Backend for Frontend)**
```
Frontend Web     ─┐
Frontend Mobile  ─┼── GraphQL API ──→ Microserviços REST internos
Smart TV App     ─┘

O GraphQL serve como camada de agregação e personalização
```

**4. Produtos SaaS com API Pública**
```
Notion, GitHub, Shopify — todos têm GraphQL API
porque clientes externos sempre querem dados de formas diferentes
```

### REST faz sentido para:

**1. APIs de pagamento e webhooks**
- `POST /payments` → simples, sem relações, sem necessidade de escolher campos

**2. Serviços internos simples**
- Um serviço que só expõe `GET /health` e `POST /process` não precisa de GraphQL

**3. Streaming e uploads**
- `POST /videos/upload` com multipart — REST é a escolha natural

**4. Integrações com parceiros externos**
- REST + OpenAPI/Swagger é o padrão de mercado para integrações B2B

---

## 📖 Glossário

| Termo | Definição |
|-------|-----------|
| **Schema** | Contrato da API GraphQL. Define tipos, queries e mutations disponíveis |
| **Type** | Tipo de dado no schema (ex: `Category`, `Course`, `String`, `ID`) |
| **Query** | Operação de leitura no GraphQL (equivalente ao GET do REST) |
| **Mutation** | Operação de escrita no GraphQL (equivalente ao POST/PUT/DELETE do REST) |
| **Resolver** | Função Go que sabe como buscar um campo específico |
| **Field Resolver** | Resolver específico para campos relacionados (ex: `category.courses`) |
| **gqlgen** | Biblioteca Go que gera código a partir do schema GraphQL |
| **Playground** | Interface visual de exploração da API GraphQL (como um Swagger interativo) |
| **Overfetching** | Receber mais dados do que o necessário numa resposta |
| **Underfetching** | Precisar de múltiplas requisições para obter todos os dados necessários |
| **N+1 Problem** | Problema de performance onde 1 query gera N queries adicionais |
| **DataLoader** | Padrão para resolver o N+1: agrupa N queries em uma só (em lote) |
| **Introspection** | Capacidade do GraphQL de descrever a si mesmo (base do Playground e docs) |
| **Scalar** | Tipo primitivo do GraphQL: `String`, `Int`, `Float`, `Boolean`, `ID` |
| **Input Type** | Tipo especial para receber parâmetros em mutations (ex: `NewCategory`) |
| **Code Generation** | Geração automática de código Go a partir do schema (o que o gqlgen faz) |
| **BFF** | Backend for Frontend — camada GraphQL que agrega e personaliza dados para clientes |
| **SQLite** | Banco de dados em arquivo, perfeito para desenvolvimento e estudo |
| **UUID** | Identificador único universal — string de 36 caracteres gerada aleatoriamente |

---

## 🚀 Próximos Passos

Agora que você entende GraphQL com gqlgen, explore estes tópicos:

**Imediato:**
- [ ] Adicionar migrations com `golang-migrate` para criar as tabelas automaticamente
- [ ] Implementar queries com filtros: `categories(name: "Backend")`
- [ ] Implementar `deleteCategory` e `updateCategory` nas mutations
- [ ] Adicionar validação de input (ex: nome não pode ser vazio)

**Intermediário:**
- [ ] **Subscriptions** — evento em tempo real (WebSocket): `subscription { onCategoryCreated { name } }`
- [ ] **DataLoader** — resolver o problema N+1 com `github.com/graph-gophers/dataloader`
- [ ] **Autenticação** — middleware JWT no contexto dos resolvers
- [ ] **Paginação** — cursor-based pagination (padrão Relay)
- [ ] **Complexity Limits** — limitar queries muito custosas

**Avançado:**
- [ ] **Federation** — juntar múltiplos serviços GraphQL num único schema (Apollo Federation)
- [ ] **gRPC internamente** — use a aula 12 (gRPC) para comunicação entre microserviços, e GraphQL como API pública
- [ ] **Persisted Queries** — cache e segurança em APIs GraphQL de produção

**Comparativo com a Aula 12:**

Aula 11 (GraphQL) e Aula 12 (gRPC) resolvem problemas parecidos de formas diferentes:

```
GraphQL                          gRPC
─────────────────────────        ─────────────────────────
HTTP/JSON                        HTTP/2 + Protocol Buffers
Ideal para frontend/clientes     Ideal para comunicação interna
Um endpoint flexível             Endpoints fixos e tipados
Playground interativo            Evans / grpcurl para debug
Streaming via Subscriptions      Streaming nativo (3 tipos)
```

Uma arquitetura madura costuma usar **os dois**: gRPC para comunicação entre microserviços (eficiente e tipado), e GraphQL como API pública para os clientes (flexível e autodocumentada).
