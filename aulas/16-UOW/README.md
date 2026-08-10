# 🔄 Unit of Work — Garantindo Atomicidade entre Múltiplos Repositórios

> Aprenda como coordenar múltiplas operações de banco de dados de forma que elas aconteçam todas juntas ou nenhuma — sem que o seu use case precise conhecer diretamente transações SQL.

---

## 📑 Sumário

- [🤔 O que é Unit of Work?](#-o-que-é-unit-of-work)
  - [Por que não usar transações "na unha" no use case?](#por-que-não-usar-transações-na-unha-no-use-case)
  - [Quando usar Unit of Work?](#quando-usar-unit-of-work)
- [🌎 Quando Unit of Work Faz Sentido](#-quando-unit-of-work-faz-sentido)
  - [Unit of Work vs. Transações Manuais](#unit-of-work-vs-transações-manuais)
  - [Unit of Work vs. Passar `*sql.Tx` por Parâmetro](#unit-of-work-vs-passar-sqltx-por-parâmetro)
  - [Unit of Work vs. Saga Pattern](#unit-of-work-vs-saga-pattern)
- [🗂️ Arquitetura do Projeto](#️-arquitetura-do-projeto)
  - [Estrutura de Pastas](#estrutura-de-pastas)
  - [O Fluxo: de Use Case a Banco de Dados](#o-fluxo-de-use-case-a-banco-de-dados)
- [📖 Conceitos Fundamentais](#-conceitos-fundamentais)
  - [RepositoryFactory e o Registro de Fábricas](#repositoryfactory-e-o-registro-de-fábricas)
  - [O Truque do Campo Público `Queries` e a Interface DBTX](#o-truque-do-campo-público-queries-e-a-interface-dbtx)
  - [Transação Preguiçosa (Lazy)](#transação-preguiçosa-lazy)
  - [Ciclo de Vida da Transação: `Do`, `CommitOrRollback` e `Rollback`](#ciclo-de-vida-da-transação-do-commitorrrollback-e-rollback)
  - [Type Assertion e a Ausência de Generics](#type-assertion-e-a-ausência-de-generics)
- [⚙️ Como Funciona o Projeto](#️-como-funciona-o-projeto)
  - [Comparação: Sem UOW vs. Com UOW](#comparação-sem-uow-vs-com-uow)
  - [Walkthrough do Wiring: Teste Completo](#walkthrough-do-wiring-teste-completo)
- [✅ Boas Práticas Presentes no Projeto](#-boas-práticas-presentes-no-projeto)
  - [1. Use Case Depende de Interface, Nunca de Infraestrutura](#1-use-case-depende-de-interface-nunca-de-infraestrutura)
  - [2. Repositórios Não Sabem que Existem Transações](#2-repositórios-não-sabem-que-existem-transações)
  - [3. Registro Nomeado de Repositórios Desacopla Wiring](#3-registro-nomeado-de-repositórios-desacopla-wiring)
- [🛡️ O que as Boas Práticas Evitaram](#️-o-que-as-boas-práticas-evitaram)
- [🚀 O que Poderia Ser Melhorado](#-o-que-poderia-ser-melhorado)
- [⚠️ Principais Problemas ao Trabalhar com Unit of Work](#️-principais-problemas-ao-trabalhar-com-unit-of-work)
  - [1. Reentrância — `Do` Não é Reentrante](#1-reentrância--do-não-é-reentrante)
  - [2. Esquecer de Registrar um Repositório](#2-esquecer-de-registrar-um-repositório)
  - [3. Transações Longas Demais Travando o Banco](#3-transações-longas-demais-travando-o-banco)
  - [4. Vazamento de `*sql.Tx` em Caso de Pânico](#4-vazamento-de-sqltx-em-caso-de-pânico)
- [🔧 Como Usar o Projeto](#-como-usar-o-projeto)
  - [Pré-requisitos](#pré-requisitos)
  - [Subindo o Banco de Dados](#subindo-o-banco-de-dados)
  - [Rodando os Testes](#rodando-os-testes)
- [📖 Glossário](#-glossário)
- [🎯 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é Unit of Work?

Imagine que você está em um banco e precisa fazer duas operações relacionadas:

1. Criar uma **categoria** (ex.: "Backend") no seu banco de dados
2. Criar um **curso** vinculado àquela categoria

Tudo vai bem na primeira operação — a categoria é inserida com sucesso. Mas aí, na segunda operação, a conexão cai, ou há um erro de validação, ou qualquer coisa inesperada acontece. O resultado? Uma categoria criada, mas nenhum curso vinculado a ela — **dados inconsistentes**.

O **Unit of Work** resolve exatamente esse problema: ele garante que **as duas operações aconteçam juntas ou nenhuma delas seja aplicada**. Se algo der errado no meio do caminho, tudo é desfeito automaticamente (rollback).

```
Sem Unit of Work:                    Com Unit of Work:
┌──────────────────┐                 ┌──────────────────────┐
│ Criar categoria  │ ✅              │ BEGIN                │
│                  │                 │  Criar categoria  ✅ │
│ Criar curso      │ ❌ ERRO         │  Criar curso      ❌ │
│                  │                 │ ROLLBACK (tudo volta) │
│ Resultado: dados │                 │                       │
│ inconsistentes   │                 │ Resultado: nada foi   │
└──────────────────┘                 │ salvo — consistente   │
                                     └──────────────────────┘
```

Formalmente, **Unit of Work é um padrão que coordena múltiplas operações de repositório dentro de uma **única transação de banco de dados**, sem que a camada de aplicação (use case, handler, etc.) precise conhecer os detalhes técnicos de como abrir, commitar ou fazer rollback de uma transação**.

### Por que não usar transações "na unha" no use case?

Existem várias formas de gerenciar transações em Go. Comparemos três abordagens:

```go
// ❌ Transação manual espalhada pelo use case (tight coupling com SQL)
func (a *AddCourseUseCase) Execute(ctx context.Context, input InputUseCase) error {
    tx, err := a.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    // O use case agora "sabe" de transação — acoplado à infraestrutura
    categoryRepo := repository.NewCategoryRepository(tx)
    courseRepo := repository.NewCourseRepository(tx)
    
    category := entity.Category{Name: input.CategoryName}
    err = categoryRepo.Insert(ctx, category)
    if err != nil {
        tx.Rollback()
        return err
    }
    
    course := entity.Course{Name: input.CourseName}
    err = courseRepo.Insert(ctx, course)
    if err != nil {
        tx.Rollback()  // rollback manual em todo lugar
        return err
    }
    
    return tx.Commit()
}
```

```go
// ⚠️ Passar *sql.Tx por parâmetro em toda função (melhor, mas ainda acoplado)
func (a *AddCourseUseCase) Execute(ctx context.Context, tx *sql.Tx, input InputUseCase) error {
    categoryRepo := repository.NewCategoryRepository(tx)
    courseRepo := repository.NewCourseRepository(tx)
    
    category := entity.Category{Name: input.CategoryName}
    err := categoryRepo.Insert(ctx, category)
    if err != nil {
        return err  // quem chama que cuida do rollback
    }
    
    // Mas agora o use case precisa de *sql.Tx como parâmetro
    // — o que faz quando NÃO tem transação? Cria outra função?
}
```

```go
// ✅ Unit of Work — use case é desacoplado de infraestrutura SQL
func (a *AddCourseUseCaseUow) Execute(ctx context.Context, input InputUseCase) error {
    return a.Uow.Do(ctx, func(uow *uow.Uow) error {
        category := entity.Category{Name: input.CategoryName}
        categoryRepo := a.getCategoryRepository(ctx)  // transparente — usa a tx dentro de Do
        err := categoryRepo.Insert(ctx, category)
        if err != nil {
            return err  // Do cuida do rollback automaticamente
        }
        
        course := entity.Course{Name: input.CourseName}
        courseRepo := a.getCourseRepository(ctx)
        err = courseRepo.Insert(ctx, course)
        if err != nil {
            return err  // Do cuida do rollback automaticamente
        }
        
        return nil  // Do commita automaticamente
    })
}
```

O UOW encapsula toda a lógica de "abrir, commitar, fazer rollback" em um lugar só. O use case não precisa saber que transações existem — ele só sabe que "dentro de um `Do`, tudo vai de uma vez ou nada vai".

### Quando usar Unit of Work?

| Situação | Abordagem recomendada |
|----------|------------------------|
| Uma única operação em um repositório (ex.: inserir um usuário, sem dependências) | `database/sql` puro, ou transação manual simples |
| Múltiplas operações em **um mesmo repositório** que precisam ser atômicas | Transação manual, ou `*sql.Tx` passado explicitamente |
| Múltiplas operações em **repositórios diferentes** que precisam ser atômicas (ex.: criar categoria + curso) | **Unit of Work** ✅ |
| Operações que cruzam **múltiplos serviços/bancos** e precisam ser atômicas | Saga pattern ou outbox pattern (UOW local não resolve) |
| Projeto pequeno, sem muita complexidade transacional | Use cases simples, sem UOW |

---

## 🌎 Quando Unit of Work Faz Sentido

Antes de se aprofundar no código, é útil entender como o Unit of Work se relaciona com outras abordagens para o mesmo problema.

### Unit of Work vs. Transações Manuais

**Transações manuais** (abrir, commitar, rollback na mão) são simples e diretas para casos pequenos:

```go
tx, _ := db.BeginTx(ctx, nil)
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

queries1 := db.New(tx)
queries1.InsertSomething(ctx, ...)

queries2 := db.New(tx)
queries2.InsertOtherThing(ctx, ...)

tx.Commit()
```

Mas quando o projeto cresce e você tem vários métodos que precisam de transações, você acaba repetindo esse padrão try/catch/commit/rollback em toda parte — muito boilerplate. Unit of Work centraliza isso.

### Unit of Work vs. Passar `*sql.Tx` por Parâmetro

Outra abordagem é simplesmente adicionar um parâmetro `*sql.Tx` em cada função de repositório:

```go
func (r *CourseRepository) Insert(ctx context.Context, tx *sql.Tx, course Course) error {
    // ...
}
```

Funciona, mas:
- O repositório precisa lidar com dois casos: quando tem `*sql.Tx` (dentro de transação) e quando é `nil` (sem transação)
- O use case ainda "sabe" que transações existem — precisa de lógica para decidir quando abrir uma
- Não escala bem quando você tem 5+ repositórios

Unit of Work resolve isso automaticamente — o repositório não precisa nem pensar em qual é a fonte de sua conexão.

### Unit of Work vs. Saga Pattern

Saga é uma alternativa quando a atomicidade precisa cruzar **múltiplos serviços**. Exemplo:

- Serviço A: criar categoria
- Serviço B: criar curso
- Se B falhar, Serviço A precisa ser "desfeito" (compensação)

Saga é bem mais complexo (envolve message queues, orquestradores, etc.) e é overkill para um único banco local. Unit of Work é a solução simples quando tudo está em um banco só.

---

## 🗂️ Arquitetura do Projeto

### Estrutura de Pastas

```
16-UOW/
├── docker-compose.yaml         ← Sobe um MySQL local para testar
├── go.mod                       ← Declaração do módulo e dependências
├── go.sum
├── sqlc.yaml                    ← Configuração do gerador SQLC (mesmo da aula 15)
├── pkg/
│   └── uow/
│       └── uow.go               ← ⭐ A implementação central do padrão Unit of Work
├── internal/
│   ├── entity/
│   │   └── entity.go            ← Domain entities: Category, Course
│   ├── repository/
│   │   ├── category.go          ← CategoryRepository (conhece sqlc, não conhece UOW)
│   │   └── course.go            ← CourseRepository (idem)
│   ├── db/                       ← ⭐ Código GERADO pelo SQLC (aula 15) — nunca editar à mão
│   │   ├── db.go                ← Interface DBTX que liga tudo
│   │   ├── models.go
│   │   └── queries.sql.go
│   └── usecase/
│       ├── add_course.go         ← Versão SEM Unit of Work (para comparação)
│       ├── add_course_test.go    ← Teste da versão sem UOW
│       ├── add_course_uow.go     ← ⭐ Versão COM Unit of Work
│       └── add_course_uow_test.go ← Teste da versão com UOW (wiring completo)
└── sql/
    ├── migrations/              ← Schema do banco (CREATE TABLE)
    ├── schema.sql
    └── queries.sql              ← Queries anotadas para SQLC
```

**Por que `pkg/uow/uow.go` fica separado do `internal/`?** Porque é código reutilizável e potencialmente publicável — é o padrão Unit of Work em si, não específico deste projeto. Já `internal/` contém código de aplicação (use cases, repositórios específicos de negócio).

### O Fluxo: de Use Case a Banco de Dados

```
┌─────────────────────────────────────────────────────────────┐
│ Handler HTTP / CLI (não existe neste projeto, mas imagine)  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────┐
        │ AddCourseUseCaseUow.Execute    │
        │ (depende de uow.UowInterface)  │
        └────────────┬───────────────────┘
                     │
                     ▼
        ┌────────────────────────────────┐
        │ Uow.Do(ctx, callback)          │
        │ ├─ BeginTx()                   │ ← Abre transação
        │ └─ fn(uow)                     │   Executa callback
        └────────────┬───────────────────┘
                     │
         ┌───────────┼───────────┐
         │           │           │
         ▼           ▼           ▼
    GetRepository  GetRepository  ...
    ("Category")   ("Course")
         │           │           │
         ▼           ▼           ▼
    CategoryRepo  CourseRepo
    (Queries      (Queries
     bound to     bound to
     *sql.Tx)     *sql.Tx)
         │           │           │
         └───────────┼───────────┘
                     │
                     ▼
        ┌────────────────────────────────┐
        │ CommitOrRollback()             │
        │ (Se callback retornou erro,    │
        │  rollback; senão, commit)      │
        └────────────────────────────────┘
```

---

## 📖 Conceitos Fundamentais

### RepositoryFactory e o Registro de Fábricas

A linha mais importante de `pkg/uow/uow.go` é esta:

```go
// pkg/uow/uow.go
type RepositoryFactory func(tx *sql.Tx) interface{}
```

Uma `RepositoryFactory` é uma função que recebe uma transação (`*sql.Tx`) e retorna um repositório (qualquer tipo, por enquanto genérico como `interface{}`).

```go
// pkg/uow/uow.go — linha 34-35
func (u *Uow) Register(name string, fc RepositoryFactory) {
    u.Repositories[name] = fc
}
```

`Register` guarda essa fábrica em um mapa (`map[string]RepositoryFactory`), indexado por um nome (string) arbitrário. Você pode registrar quantos repositórios quiser:

```go
// internal/usecase/add_course_uow_test.go — linhas 29-39
uow.Register("CategoryRepository", func(tx *sql.Tx) interface{} {
    repo := repository.NewCategoryRepository(dbt)  // cria o repo "normalmente"
    repo.Queries = db.New(tx)                       // depois rebinda a parte SQL à transação
    return repo
})

uow.Register("CourseRepository", func(tx *sql.Tx) interface{} {
    repo := repository.NewCourseRepository(dbt)
    repo.Queries = db.New(tx)
    return repo
})
```

Quando você chama `uow.GetRepository(ctx, "CategoryRepository")`, o `Uow` invoca aquela fábrica com a transação atual:

```go
// pkg/uow/uow.go — linha 50
repo := u.Repositories[name](u.Tx)  // executa a fábrica com u.Tx
return repo, nil
```

Resultado: um repositório que está "amarrado" à transação.

### O Truque do Campo Público `Queries` e a Interface DBTX

Aqui está a magia que conecta o Unit of Work com o código gerado pelo SQLC (visto na aula 15).

Cada repositório tem um campo público `Queries` (gerado pelo SQLC):

```go
// internal/repository/category.go
type CategoryRepository struct {
    DB      *sql.DB
    Queries *db.Queries  ← campo público
}

func NewCategoryRepository(dtb *sql.DB) *CategoryRepository {
    return &CategoryRepository{
        DB:      dtb,
        Queries: db.New(dtb),  ← inicializado com *sql.DB
    }
}
```

O `db.Queries` (gerado pelo SQLC) é inicializado com `db.New(db DBTX)`, onde `DBTX` é uma interface:

```go
// internal/db/db.go — código gerado pelo SQLC
type DBTX interface {
    ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...interface{}) *sql.Row
    PrepareContext(context.Context, string) (*sql.Stmt, error)
}

func New(db DBTX) *Queries {
    return &Queries{db: db}
}
```

E aqui está o truque: **tanto `*sql.DB` quanto `*sql.Tx` implementam essa interface `DBTX`**. Então, após criar um repositório com a conexão normal, o Unit of Work simplesmente sobrescreve o campo `Queries`:

```go
repo := repository.NewCategoryRepository(dbt)   // Queries = db.New(*sql.DB)
repo.Queries = db.New(tx)                        // Queries = db.New(*sql.Tx) ← rebindado!
```

Agora, quando o repositório chamar `repo.Insert()`, internamente ele usa o `Queries` que está amarrado à transação, não à conexão comum. Mágica desacoplada!

```
Sem Unit of Work:        Com Unit of Work:
┌──────────────┐         ┌──────────────┐
│ *sql.DB      │         │ *sql.Tx      │
└──────┬───────┘         └──────┬───────┘
       │                        │
       ▼                        ▼
   db.New()    ────────>   db.New()
       │                        │
       ▼                        ▼
   Queries                  Queries
   (no transac)         (bound to tx)
```

### Transação Preguiçosa (Lazy)

O Unit of Work não abre uma transação no `NewUow`. Ele só abre quando alguém pede um repositório:

```go
// pkg/uow/uow.go — linhas 42-51
func (u *Uow) GetRepository(ctx context.Context, name string) (interface{}, error) {
    if u.Tx == nil {
        tx, err := u.Db.BeginTx(ctx, nil)  ← Só abre aqui
        if err != nil {
            return nil, err
        }
        u.Tx = tx
    }
    repo := u.Repositories[name](u.Tx)
    return repo, nil
}
```

Isso é eficiente: se seu callback no `Do` não precisar de repositórios (caso raro), nenhuma transação é aberta.

### Ciclo de Vida da Transação: `Do`, `CommitOrRollback` e `Rollback`

O `Do` é o "maestro" da transação:

```go
// pkg/uow/uow.go — linhas 54-71
func (u *Uow) Do(ctx context.Context, fn func(Uow *Uow) error) error {
    if u.Tx != nil {
        return fmt.Errorf("transaction already started")  ← Reentrância não permitida
    }
    tx, err := u.Db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    u.Tx = tx
    err = fn(u)                      ← Executa o callback
    if err != nil {                  ← Se o callback retornar erro
        errRb := u.Rollback()        ← Tenta rollback
        if errRb != nil {            ← Se o rollback falhar também
            return errors.New(fmt.Sprintf("original error: %s, rollback error: %s", ...))
        }
        return err                   ← Retorna o erro original (não o do rollback)
    }
    return u.CommitOrRollback()      ← Se tudo OK, commita
}
```

Se o callback retornar erro, `Rollback()` é chamado automaticamente (sem que você precise fazer nada). Se o callback retornar `nil`, `CommitOrRollback()` é chamado:

```go
// pkg/uow/uow.go — linhas 86-96
func (u *Uow) CommitOrRollback() error {
    err := u.Tx.Commit()
    if err != nil {
        errRb := u.Rollback()        ← Se o commit falhar, tenta desfazer
        if errRb != nil {
            return errors.New(fmt.Sprintf("original error: %s, rollback error: %s", ...))
        }
        return err
    }
    u.Tx = nil                       ← Limpa a transação
    return nil
}
```

### Type Assertion e a Ausência de Generics

Como `GetRepository` retorna `interface{}`, você precisa fazer um type assertion para obter o tipo concreto:

```go
// internal/usecase/add_course_uow.go — linhas 52-57
func (a *AddCourseUseCaseUow) getCategoryRepository(ctx context.Context) repository.CategoryRepositoryInterface {
    repo, err := a.Uow.GetRepository(ctx, "CategoryRepository")
    if err != nil {
        panic(err)
    }
    return repo.(repository.CategoryRepositoryInterface)  ← Type assertion
}
```

**Por que `interface{}`?** Go 1.17 (versão de quando esse projeto foi escrito) não tinha generics. A alternativa seria:

```go
// Com generics (Go 1.18+), seria:
func (u *Uow) GetRepository[T any](ctx context.Context, name string) (T, error) {
    // ...
}
// Uso:
repo := uow.GetRepository[*CategoryRepository](ctx, "CategoryRepository")
```

Mas sem generics, `interface{}` é a única forma de ter uma função genérica. O custo é que o type assertion só é verificado em runtime — se você fizer `repo.(string)` quando o tipo é na verdade um `*CategoryRepository`, vai ser um pânico em runtime.

---

## ⚙️ Como Funciona o Projeto

### Comparação: Sem UOW vs. Com UOW

A pasta `internal/usecase/` tem dois arquivos que você deve comparar lado a lado.

**`add_course.go` (sem Unit of Work):**

```go
// internal/usecase/add_course.go
func (a *AddCourseUseCase) Execute(ctx context.Context, input InputUseCase) error {
    category := entity.Category{
        Name: input.CategoryName,
    }

    err := a.CategoryRepository.Insert(ctx, category)
    if err != nil {
        return err
    }

    course := entity.Course{
        Name:       input.CourseName,
        CategoryID: input.CourseCategoryID,
    }

    err = a.CourseRepository.Insert(ctx, course)  // ← Se isso falhar, categoria fica criada!
    if err != nil {
        return err
    }

    return nil
}
```

As duas operações são **independentes**. Se a segunda falhar, a primeira já está committed no banco. Dados inconsistentes.

**`add_course_uow.go` (com Unit of Work):**

```go
// internal/usecase/add_course_uow.go
func (a *AddCourseUseCaseUow) Execute(ctx context.Context, input InputUseCase) error {
    return a.Uow.Do(ctx, func(uow *uow.Uow) error {  // ← Tudo dentro de uma transação
        category := entity.Category{
            Name: input.CategoryName,
        }
        repoCategory := a.getCategoryRepository(ctx)
        err := repoCategory.Insert(ctx, category)
        if err != nil {
            return err  // ← Do faz rollback automaticamente
        }

        course := entity.Course{
            Name:       input.CourseName,
            CategoryID: input.CourseCategoryID,
        }

        repoCourse := a.getCourseRepository(ctx)
        err = repoCourse.Insert(ctx, course)  // ← Se isso falhar, a categoria também é desfeita
        if err != nil {
            return err  // ← Do faz rollback automaticamente
        }
        return nil  // ← Se chegar aqui, Do commita automaticamente
    })
}
```

As duas operações estão dentro do callback de `Do`, então compartilham a mesma transação. Se qualquer uma falhar, tudo é desfeito.

### Walkthrough do Wiring: Teste Completo

O teste `internal/usecase/add_course_uow_test.go` mostra o fluxo ponta a ponta. Vamos quebrar em partes:

**Passo 1: Criar o Uow e Registrar Repositórios**

```go
// internal/usecase/add_course_uow_test.go — linhas 26-39
ctx := context.Background()
uow := uow.NewUow(ctx, dbt)

uow.Register("CategoryRepository", func(tx *sql.Tx) interface{} {
    repo := repository.NewCategoryRepository(dbt)
    repo.Queries = db.New(tx)  ← O truque: rebindar Queries à transação
    return repo
})

uow.Register("CourseRepository", func(tx *sql.Tx) interface{} {
    repo := repository.NewCourseRepository(dbt)
    repo.Queries = db.New(tx)
    return repo
})
```

Aqui você cria um `Uow`, depois registra dois repositórios com fábricas. Cada fábrica sabe como construir o repositório.

**Passo 2: Criar o Use Case com o Uow Injetado**

```go
// internal/usecase/add_course_uow_test.go — linha 47
useCase := NewAddCourseUseCaseUow(uow)
```

`NewAddCourseUseCaseUow` recebe apenas o `uow.UowInterface` — nada de `*sql.DB` ou detalhe de infraestrutura.

**Passo 3: Executar**

```go
// internal/usecase/add_course_uow_test.go — linhas 48-49
input := InputUseCase{
    CategoryName:     "Category 1",
    CourseName:       "Course 1",
    CourseCategoryID: 2,  // Este ID não existe! (será erro de FK)
}

err = useCase.Execute(ctx, input)
assert.NoError(t, err)  // ← Esperaria erro por FK, mas teste passa (curiosidade!)
```

Internamente, `Execute` chama `a.Uow.Do(ctx, callback)`, que:

1. Abre uma transação
2. Executa o callback
3. Dentro do callback, `getCategoryRepository` e `getCourseRepository` invocam as fábricas registradas com a transação atual
4. Ambos os `Insert` usam a mesma `*sql.Tx`
5. Se ambos retornarem `nil`, tudo é committed

---

## ✅ Boas Práticas Presentes no Projeto

### 1. Use Case Depende de Interface, Nunca de Infraestrutura

```go
// ✅ Bom — depende de UowInterface (abstração)
type AddCourseUseCaseUow struct {
    Uow uow.UowInterface
}

// ❌ Ruim — dependeria de *sql.DB (infraestrutura concreta)
type AddCourseUseCase struct {
    db *sql.DB
}
```

Isso permite testar o use case mockando o `UowInterface`, sem precisar de um banco real.

### 2. Repositórios Não Sabem que Existem Transações

```go
// internal/repository/category.go
type CategoryRepository struct {
    DB      *sql.DB
    Queries *db.Queries
}

func (r *CategoryRepository) Insert(ctx context.Context, category entity.Category) error {
    return r.Queries.CreateCategory(ctx, db.CreateCategoryParams{
        Name: category.Name,
    })
}
```

O repositório não sabe se `r.Queries` está ligado a `*sql.DB` ou `*sql.Tx` — ele só usa. Isso é uma separação clara de responsabilidades.

### 3. Registro Nomeado de Repositórios Desacopla Wiring

```go
uow.Register("CategoryRepository", factory)
uow.Register("CourseRepository", factory)
```

Usar string como chave permite que o bootstrap (teste, main, DI container) decida quais repos existem, sem que o `Uow` precise conhecê-los.

---

## 🛡️ O que as Boas Práticas Evitaram

### Inconsistência de Dados entre Tabelas Relacionadas

Sem atomicidade, um erro parcial deixa o banco em estado inconsistente (categoria criada, curso não). Isso pode causar:

- Queries que referenciam registros órfãos
- Tentativas de retry que criam duplicatas
- Auditoria e troubleshooting complexos

A transação garante tudo-ou-nada.

### Acoplamento do Use Case a Detalhes SQL

Se o use case tivesse de gerenciar `*sql.Tx` manualmente, ele estaria acoplado a detalhes de infraestrutura. Isso dificulta:

- Testar sem um banco real
- Trocar de banco sem reescrever a camada de aplicação
- Reusar o use case em contextos diferentes (CLI vs. API HTTP vs. worker)

---

## 🚀 O que Poderia Ser Melhorado

### 1. Trocar `panic` nos Helpers por Retorno de Erro

Hoje:

```go
// internal/usecase/add_course_uow.go — linhas 52-57
func (a *AddCourseUseCaseUow) getCategoryRepository(ctx context.Context) repository.CategoryRepositoryInterface {
    repo, err := a.Uow.GetRepository(ctx, "CategoryRepository")
    if err != nil {
        panic(err)  // ❌ Encerra o programa
    }
    return repo.(repository.CategoryRepositoryInterface)
}
```

Melhor seria:

```go
// ✅ Retornar erro normalmente
func (a *AddCourseUseCaseUow) getCategoryRepository(ctx context.Context) (repository.CategoryRepositoryInterface, error) {
    repo, err := a.Uow.GetRepository(ctx, "CategoryRepository")
    if err != nil {
        return nil, err
    }
    return repo.(repository.CategoryRepositoryInterface), nil
}

// E usar em Execute:
repoCategory, err := a.getCategoryRepository(ctx)
if err != nil {
    return err
}
```

### 2. Usar Generics (Go 1.18+) para Eliminar Type Assertion

Hoje:

```go
repo := u.Repositories[name](u.Tx)
return repo, nil  // retorna interface{}
// Depois: repo.(repository.CategoryRepositoryInterface)
```

Com generics:

```go
func (u *Uow) GetRepository[T any](ctx context.Context, name string) (T, error) {
    if u.Tx == nil {
        tx, err := u.Db.BeginTx(ctx, nil)
        if err != nil {
            var zero T
            return zero, err
        }
        u.Tx = tx
    }
    repo := u.Repositories[name](u.Tx)
    return repo.(T), nil
}

// Uso — type-safe:
repoCategory, err := uow.GetRepository[*repository.CategoryRepository](ctx, "CategoryRepository")
```

### 3. Limpar Repositórios entre Execuções

O mapa `u.Repositories` persiste entre chamadas de `Do`. Se quiser reussar o mesmo `Uow`, seria bom:

```go
// ✅ Limpar após commit/rollback
func (u *Uow) CommitOrRollback() error {
    err := u.Tx.Commit()
    // ...
    u.Tx = nil
    u.Repositories = make(map[string]RepositoryFactory)  // ← Limpar
    return nil
}
```

Ou simplesmente criar um novo `Uow` a cada request HTTP.

### 4. DSN Hardcoded → Variável de Ambiente

Hoje a string de conexão está no teste:

```go
dbt, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
```

Deveria vir do ambiente:

```go
dsn := os.Getenv("DATABASE_URL")
dbt, err := sql.Open("mysql", dsn)
```

### 5. Adicionar `defer dbt.Close()` nos Testes

```go
dbt, err := sql.Open("mysql", dsn)
assert.NoError(t, err)
defer dbt.Close()  // ← Fechar conexão ao fim do teste
```

---

## ⚠️ Principais Problemas ao Trabalhar com Unit of Work

### 1. Reentrância — `Do` Não é Reentrante

Se você chamar `Do` duas vezes no mesmo `Uow`, a segunda chamada falhará:

```go
err := uow.Do(ctx, func(u *uow.Uow) error {
    // ...
    err := uow.Do(ctx, func(u *uow.Uow) error {  // ❌ ERRO
        // "transaction already started"
    })
    return err
})
```

`Solução:` Sempre crie um novo `Uow` por operação atômica, ou estruture seu código para não reivindicar transações aninhadas.

### 2. Esquecer de Registrar um Repositório

Se você pedir um repositório que não foi registrado, vai virar `nil` em runtime:

```go
uow.Register("CategoryRepository", ...)
// Esquecer de registrar CourseRepository

repo, _ := uow.GetRepository(ctx, "CourseRepository")
repo.(repository.CourseRepositoryInterface)  // ❌ Pânico — repo é nil
```

Solução: Criar um helper que valida se todos os repos necessários foram registrados no bootstrap.

### 3. Transações Longas Demais Travando o Banco

Se seu callback fizer operações lentas (chamadas HTTP, processamento pesado), a transação fica aberta por muito tempo, travando linhas/tabelas no banco:

```go
return a.Uow.Do(ctx, func(uow *uow.Uow) error {
    err := repoCategory.Insert(ctx, category)
    if err != nil {
        return err
    }
    
    resp, err := http.Get("https://external-api.com/data")  // ❌ Transação aberta durante IO
    if err != nil {
        return err
    }
    
    course := parseCourse(resp)
    return repoCourse.Insert(ctx, course)
})
```

Solução: Fazer operações externas **fora** da transação, trazer os dados para dentro.

### 4. Vazamento de `*sql.Tx` em Caso de Pânico

Se uma função dentro do callback fizer `panic()`, o `Do` não consegue fazer rollback (o pânico pula por cima):

```go
err := uow.Do(ctx, func(uow *uow.Uow) error {
    panic("algo deu muito errado")  // ❌ Transação fica aberta, conexão vazada
})
```

Solução: Usar `recover()` em volta do `Do`, ou evitar `panic()` completamente (retornar erro normalmente).

---

## 🔧 Como Usar o Projeto

### Pré-requisitos

```bash
# Go 1.19 ou superior
go version

# Docker (para subir o MySQL)
docker --version

# MySQL CLI (opcional, para inspecionar dados manualmente)
# macOS:
brew install mysql-client
```

### Subindo o Banco de Dados

```bash
cd aulas/16-UOW
docker-compose up -d

# Verifica se o container está rodando
docker ps | grep mysql

# Você pode se conectar manualmente (opcional):
mysql -h 127.0.0.1 -u root -p root -e "SHOW DATABASES;"
```

### Rodando os Testes

```bash
# Dentro de aulas/16-UOW
go test ./... -v

# Ou teste específico:
go test -v ./internal/usecase -run TestAddCourseUow
```

Cada teste:
1. Conecta no banco
2. Cria as tabelas (DROP + CREATE)
3. Executa as operações
4. Verifica resultado com `assert`

---

## 📖 Glossário

| Termo | Definição |
|-------|-----------|
| **Unit of Work (UOW)** | Padrão que coordena múltiplas operações de repositório dentro de uma única transação, abstraindo a complexidade de begin/commit/rollback da camada de aplicação |
| **RepositoryFactory** | Função que recebe uma `*sql.Tx` e retorna uma instância de repositório amarrada àquela transação |
| **Transação** | Sequência de operações no banco que são atômicas: todas são aplicadas ou nenhuma é |
| **Atomicidade** | Garantia de que múltiplas operações em banco de dados ou desaparecem juntas (rollback) ou chegam juntas (commit) |
| **Rollback** | Desfazer todas as mudanças feitas em uma transação que ainda não foi commitada |
| **DBTX** | Interface gerada pelo SQLC que agrupa os métodos de execução de queries, satisfeita tanto por `*sql.DB` quanto por `*sql.Tx` |
| **Type Assertion** | Operação em Go que "afirma" o tipo concreto de um valor de interface (`valor.(ConcreteType)`) |
| **Lazy Loading / Transação Preguiçosa** | Adiar a criação de um recurso (aqui, a transação) até o primeiro uso, economizando overhead |
| **Repository Pattern** | Padrão que abstrai a fonte de dados, oferecendo uma interface de coleção em memória para persistência |
| **Domain Entity** | Objeto que representa um conceito do negócio (ex.: Category, Course) independente de tecnologia |
| **Generics** | Recurso do Go 1.18+ que permite escrever código parametrizado por tipo (`func F[T any]()`) |

---

## 🎯 Próximos Passos

### Para consolidar o aprendizado:

1. **Escreva um teste de falha** — modifique `add_course_uow_test.go` para que o segundo `Insert` lance um erro propositalmente, e verifique que a categoria é realmente desfeita (ficou sem registros no banco). Confirme que sem o UOW (testando `add_course_test.go`), a categoria fica criada.

2. **Remova os `panic`s** — refatore `getCategoryRepository` e `getCourseRepository` em `add_course_uow.go` para retornar erro em vez de fazer `panic`.

3. **Use generics (Go 1.18+)** — migre `Uow.GetRepository` para usar generics e elimine a type assertion manual em `add_course_uow.go`. Compare a ergonomia.

4. **Compare com a aula anterior** — leia o README de aulas/15-SQLC para entender o `DBTX` gerado pelo SQLC e como ele é fundamental para o padrão UOW funcionar.

5. **Implemente `UnRegister`** — hoje existe o método, mas não é usado. Crie um teste que registra, depois desregistra um repositório e tenta usá-lo (deve falhar ou retornar erro).

6. **Crie um novo use case transacional** — invente um segundo use case que precisa de 3+ repositórios (ex.: criar usuário + perfil + permissão atomicamente) para praticar o padrão com mais complexidade.

### Conceitos relacionados no curso:

- **Aula 15 (SQLC)** — entender como o `DBTX` interface permite que o mesmo `Queries` trabalhe com `*sql.DB` ou `*sql.Tx`.
- **Aulas anteriores de database/sql** — transações manuais com `BeginTx()`, `Commit()`, `Rollback()` que o UOW encapsula.
- **Aulas de arquitetura e padrões** — como Unit of Work se relaciona com Repository, Dependency Injection, e Clean Architecture.
- **Aulas futuras (se houver)** — transações distribuídas, Saga pattern para múltiplos bancos/serviços.
