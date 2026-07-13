# 🛠️ SQLC — Gerando Código Go Type-Safe a partir de SQL Puro

> Aprenda como escrever apenas SQL e deixar uma ferramenta gerar todo o código Go de acesso ao banco de dados para você — sem ORM, sem reflection, sem scans manuais.

---

## 📑 Sumário

- [🤔 O que é SQLC?](#-o-que-é-sqlc)
  - [Por que não usar um ORM ou database/sql puro?](#por-que-não-usar-um-orm-ou-databasesql-puro)
  - [Quando usar SQLC?](#quando-usar-sqlc)
- [🌎 Panorama Atual: Formas de Trabalhar com Banco de Dados em Go](#-panorama-atual-formas-de-trabalhar-com-banco-de-dados-em-go)
  - [Tabela Comparativa](#tabela-comparativa)
  - [Abordagem de Baixo Nível: database/sql, sqlx e pgx](#abordagem-de-baixo-nível-databasesql-sqlx-e-pgx)
  - [Abordagem de Geração de Código: SQLC e Ent](#abordagem-de-geração-de-código-sqlc-e-ent)
  - [Abordagem de ORM Tradicional: GORM e Bun](#abordagem-de-orm-tradicional-gorm-e-bun)
  - [Uma Categoria à Parte: Ferramentas de Migration](#uma-categoria-à-parte-ferramentas-de-migration)
  - [Onde o SQLC se Encaixa](#onde-o-sqlc-se-encaixa)
- [🗂️ Arquitetura do Projeto](#️-arquitetura-do-projeto)
  - [Estrutura de Pastas](#estrutura-de-pastas)
  - [O Fluxo: de SQL a Código Go](#o-fluxo-de-sql-a-código-go)
- [📖 Conceitos Fundamentais](#-conceitos-fundamentais)
  - [Migrations: o Schema do Banco](#migrations-o-schema-do-banco)
  - [Queries Anotadas](#queries-anotadas)
  - [sqlc.yaml — o Arquivo de Configuração](#sqlcyaml--o-arquivo-de-configuração)
  - [Código Gerado: models.go](#código-gerado-modelsgo)
  - [Código Gerado: db.go e a Interface DBTX](#código-gerado-dbgo-e-a-interface-dbtx)
  - [Código Gerado: query.sql.go](#código-gerado-querysqlgo)
  - [sql.NullString — Lidando com Valores Nulos](#sqlnullstring--lidando-com-valores-nulos)
- [⚙️ Como Funciona o Projeto](#️-como-funciona-o-projeto)
  - [cmd/runSQLC — Uso Básico do Queries Gerado](#cmdrunsqlc--uso-básico-do-queries-gerado)
  - [cmd/runSQLCTX — Transações com callTx](#cmdrunsqlctx--transações-com-calltx)
- [✅ Boas Práticas Presentes no Projeto](#-boas-práticas-presentes-no-projeto)
  - [1. Código Gerado Nunca é Editado à Mão](#1-código-gerado-nunca-é-editado-à-mão)
  - [2. Interface DBTX para Reuso em Transações](#2-interface-dbtx-para-reuso-em-transações)
  - [3. Separação entre Schema, Queries e Configuração](#3-separação-entre-schema-queries-e-configuração)
  - [4. Um Tipo por Formato de Resultado](#4-um-tipo-por-formato-de-resultado)
- [🛡️ O que as Boas Práticas Evitaram](#️-o-que-as-boas-práticas-evitaram)
- [🚀 O que Poderia Ser Melhorado](#-o-que-poderia-ser-melhorado)
- [⚠️ Principais Problemas ao Trabalhar com SQLC](#️-principais-problemas-ao-trabalhar-com-sqlc)
  - [1. Esquecer de Rodar sqlc generate](#1-esquecer-de-rodar-sqlc-generate)
  - [2. Migrations e Queries Fora de Sincronia](#2-migrations-e-queries-fora-de-sincronia)
  - [3. Placeholders Diferem entre Bancos](#3-placeholders-diferem-entre-bancos)
  - [4. Nulidade Contamina Todo o Código](#4-nulidade-contamina-todo-o-código)
  - [5. DSN com Credenciais no Código](#5-dsn-com-credenciais-no-código)
  - [6. Código de Exemplo Comentado](#6-código-de-exemplo-comentado)
- [🔧 Como Usar o Projeto](#-como-usar-o-projeto)
  - [Pré-requisitos](#pré-requisitos)
  - [Subindo o Banco de Dados](#subindo-o-banco-de-dados)
  - [Rodando as Migrations](#rodando-as-migrations)
  - [Gerando o Código com SQLC](#gerando-o-código-com-sqlc)
  - [Executando os Exemplos](#executando-os-exemplos)
- [📖 Glossário](#-glossário)
- [🎯 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é SQLC?

Imagine que, em vez de escrever manualmente todo o código Go que abre uma conexão, monta uma query, executa e faz o "scan" de cada coluna do resultado para uma struct, você pudesse simplesmente escrever a query SQL... e ganhar de graça uma função Go pronta, com os tipos certos, para chamá-la.

É exatamente isso que o **SQLC** faz. Ele é um gerador de código: você escreve o **schema** do seu banco (as tabelas) e as **queries** que quer executar (SELECT, INSERT, UPDATE, DELETE) em arquivos `.sql` comuns, e o SQLC lê esses arquivos e gera automaticamente:

- Uma `struct` Go para cada tabela (o "modelo").
- Uma função Go para cada query, já com os parâmetros e o tipo de retorno corretos.

```
Você escreve isto (SQL):          Você ganha isto (Go), de graça:
┌─────────────────────────┐      ┌──────────────────────────────────────┐
│ -- name: GetCategory :one│  →   │ func (q *Queries) GetCategory(       │
│ SELECT * FROM categories  │      │     ctx context.Context, id string, │
│ WHERE id = ?;             │      │ ) (Category, error) { ... }          │
└─────────────────────────┘      └──────────────────────────────────────┘
```

Nenhum código Go é escrito à mão para acessar o banco — ele nasce a partir do SQL.

### Por que não usar um ORM ou database/sql puro?

Existem hoje três caminhos comuns em Go para conversar com um banco relacional:

```go
// ❌ database/sql puro — funciona, mas todo scan é manual e repetitivo
rows, err := db.Query("SELECT id, name, description FROM categories")
if err != nil {
    return nil, err
}
defer rows.Close()

var categories []Category
for rows.Next() {
    var c Category
    if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
        return nil, err
    }
    categories = append(categories, c)
}
```

```go
// ⚠️ ORM (como GORM) — menos código, mas usa reflection em runtime
// e pode gerar SQL implícito que você não vê nem controla totalmente
var categories []Category
db.Find(&categories)
```

```go
// ✅ SQLC — você escreve o SQL de verdade, e o Scan/loop acima
// é gerado automaticamente, sem reflection em runtime
categories, err := queries.ListCategories(ctx)
```

O SQLC fica no meio do caminho: você mantém controle total sobre o SQL que será executado (nada de "SQL mágico" gerado por um ORM), mas não escreve à mão o código repetitivo de `Scan`. E como o código é gerado **em tempo de build**, os erros de tipo aparecem em tempo de compilação — não em produção.

### Quando usar SQLC?

| Situação | Ferramenta recomendada |
|----------|------------------------|
| Projeto pequeno, poucas queries simples | `database/sql` puro |
| Quer controle total do SQL executado, mas sem repetir scans | **SQLC** ✅ |
| Precisa de recursos avançados de ORM (associations automáticas, soft delete, hooks) | GORM ou similar |
| Time já domina SQL e quer evitar "magia" de ORM | **SQLC** ✅ |
| Projeto com múltiplos bancos/dialetos SQL divergentes | Avaliar com cuidado — queries são específicas por dialeto |

> Essa tabela é um resumo rápido. Para uma visão mais ampla de **todas** as formas populares de trabalhar com banco de dados em Go hoje — não só SQLC vs GORM — veja a seção [🌎 Panorama Atual: Formas de Trabalhar com Banco de Dados em Go](#-panorama-atual-formas-de-trabalhar-com-banco-de-dados-em-go) logo abaixo.

---

## 🌎 Panorama Atual: Formas de Trabalhar com Banco de Dados em Go

Diferente de linguagens como Python (Django ORM) ou Ruby (ActiveRecord), Go **não tem um framework de banco de dados oficial**. A biblioteca padrão (`database/sql`) oferece só o essencial — abrir conexão, rodar query, ler resultado — e deixou um espaço enorme para a comunidade construir ferramentas com filosofias bem diferentes entre si. Por isso, ao entrar em um projeto Go novo, é comum encontrar abordagens bem diferentes da que você já conhece.

Esta seção mapeia as ferramentas mais relevantes do ecossistema atual, para você entender onde o SQLC (o assunto deste projeto) se encaixa entre elas.

### Tabela Comparativa

| Ferramenta | Categoria | Controle sobre o SQL | Performance | Curva de aprendizado | Segurança de tipos |
|------------|-----------|----------------------|-------------|----------------------|--------------------|
| `database/sql` (padrão) | Baixo nível | Total (você escreve tudo) | Máxima | Baixa API, mas alta verbosidade | Só em runtime (scans manuais) |
| `sqlx` | Baixo nível "com açúcar" | Total | Alta | Baixa | Só em runtime |
| `pgx` | Driver/toolkit (Postgres) | Total | Alta (a mais rápida p/ Postgres) | Média | Só em runtime |
| **SQLC** | Geração de código a partir de SQL | Total (você escreve o SQL) | Alta (sem reflection) | Média | **Compile-time** |
| Ent | Geração de código a partir de um DSL Go | Baixo/médio (schema declarado em Go) | Alta | Média/Alta | **Compile-time** |
| GORM | ORM tradicional | Baixo (SQL gerado pelo ORM) | Média (usa reflection) | Baixa (muito popular, muita doc) | Só em runtime |
| Bun | ORM leve | Médio | Média/Alta | Baixa/Média | Só em runtime |

> "Compile-time" aqui significa que erros de tipo (ex.: comparar uma coluna `string` com um `int`) aparecem quando o código é **compilado**, não quando a query roda em produção contra o banco.

### Abordagem de Baixo Nível: `database/sql`, `sqlx` e `pgx`

Essas três ferramentas têm em comum: **você mesmo escreve o SQL e mesmo assim ainda participa ativamente do "scan" do resultado** — nenhuma delas gera código automaticamente por você.

- **`database/sql`** — pacote da biblioteca padrão do Go. É a base sobre a qual quase todo o resto é construído (inclusive o código que o SQLC gera usa `database/sql` por baixo). Reúne o mínimo necessário: `Open`, `Query`, `Exec`, `Scan`. Você já viu esse padrão em ação neste próprio projeto, no arquivo gerado [`internal/db/query.sql.go`](internal/db/query.sql.go).
- **`sqlx`** (`github.com/jmoiron/sqlx`) — uma camada fina por cima de `database/sql`, que adiciona helpers como `StructScan`, `Get` e `Select` para reduzir a repetição de `rows.Scan(&a, &b, &c)` linha a linha. Continua sendo "baixo nível" (você escreve o SQL manualmente e não há geração de código), mas com bem menos boilerplate que o `database/sql` puro.
- **`pgx`** (`github.com/jackc/pgx/v5`) — driver e toolkit específico para PostgreSQL, considerado hoje o driver Postgres mais rápido e completo em Go. Vai além de um driver comum: tem seu próprio pool de conexões (`pgxpool`), suporta tipos nativos do Postgres (arrays, JSON/JSONB, UUID) de forma mais direta que o driver genérico via `database/sql`, e é frequentemente usado **junto** com o SQLC (o SQLC pode gerar código usando `pgx` como alvo em vez de `database/sql`, quando o projeto usa Postgres).

### Abordagem de Geração de Código: SQLC e Ent

Essas ferramentas têm em comum a ideia de **gerar código Go automaticamente antes da compilação**, eliminando reflection em runtime — mas partem de fontes diferentes:

- **SQLC** (`github.com/sqlc-dev/sqlc`) — como você viu ao longo deste README, parte de **SQL puro** (schema + queries anotadas) e gera as structs e funções Go. Você continua no controle total do SQL executado.
- **Ent** (`entgo.io/ent`) — criado originalmente na Meta/Facebook, parte de um **schema declarado em código Go** (não em SQL) — você descreve entidades, campos e relacionamentos usando um DSL Go, e o Ent gera todo o código de acesso, incluindo navegação de relacionamentos complexas (grafos de entidades) e até integração automática com GraphQL. É mais indicado quando o domínio tem muitas relações entre entidades e você prefere modelar isso em Go em vez de SQL.

### Abordagem de ORM Tradicional: GORM e Bun

ORMs mapeiam objetos/structs Go diretamente para tabelas, gerando o SQL "por trás das cenas" via reflection em tempo de execução:

- **GORM** (`gorm.io/gorm`) — o ORM mais popular e mais usado do ecossistema Go atualmente. Tem produtividade alta: migrations automáticas, hooks (before/after save), soft delete, eager/lazy loading de associations. O custo é menos controle explícito sobre o SQL final gerado, e overhead de performance por causa do uso de reflection.
- **Bun** (`github.com/uptrace/bun`) — um ORM mais leve, geralmente citado como alternativa ao GORM com melhor performance e uma API mais próxima do SQL (os métodos do query builder se parecem mais com a query SQL real do que os do GORM), mantendo boa parte da produtividade de um ORM.

### Uma Categoria à Parte: Ferramentas de Migration

`golang-migrate`, `goose` e `Atlas` não competem com as ferramentas acima — elas resolvem um problema ortogonal: **aplicar e versionar mudanças no schema do banco**. Você pode usar qualquer uma delas junto de qualquer ferramenta de acesso a dados (inclusive SQLC, como este projeto faz com `golang-migrate` — veja [Migrations: o Schema do Banco](#migrations-o-schema-do-banco)). A escolha de "como migrar o schema" é independente da escolha de "como consultar os dados".

### Onde o SQLC se Encaixa

O SQLC ocupa um espaço específico no meio dessas opções: ele te dá o **controle explícito de SQL** que `database/sql`, `sqlx` e `pgx` oferecem, mas sem o boilerplate manual de scans — porque o código é **gerado**, assim como no Ent, mas a partir de SQL puro em vez de um DSL Go. Isso o torna uma escolha popular para quem já sabe SQL e quer manter esse conhecimento útil, evitando a "caixa-preta" de um ORM tradicional como GORM, mas também evitando reescrever manualmente cada `Scan()`.

---

## 🗂️ Arquitetura do Projeto

### Estrutura de Pastas

```
15-SQLC/
├── Makefile                    ← Atalhos para criar/rodar/desfazer migrations
├── docker-compose.yaml         ← Sobe um MySQL local para testar o projeto
├── go.mod                      ← Declaração do módulo e dependências
├── sqlc.yaml                   ← Configuração do gerador SQLC
├── cmd/                        ← Pontos de entrada (uso do código gerado)
│   ├── runSQLC/
│   │   └── main.go             ← Exemplo de uso direto do Queries gerado
│   └── runSQLCTX/
│       └── main.go             ← Exemplo de uso com transações
├── internal/db/                ← Código GERADO pelo SQLC — nunca editar à mão
│   ├── db.go                   ← Interface DBTX, struct Queries, New(), WithTx()
│   ├── models.go                ← Structs Category e Course (uma por tabela)
│   └── query.sql.go             ← Uma função Go para cada query anotada
└── sql/
    ├── migrations/              ← Schema do banco (o que existe)
    │   ├── 000001_init.up.sql
    │   └── 000001_init.down.sql
    └── queries/
        └── query.sql            ← Queries anotadas (o que você quer executar)
```

**Por que `internal/db/` fica separado de `cmd/`?** Porque `internal/db/` é **gerado automaticamente** — todo arquivo ali começa com o comentário `// Code generated by sqlc. DO NOT EDIT.`. Já `cmd/` contém código escrito por você, que **usa** o que foi gerado. Separar as duas coisas deixa claro: "isso aqui eu não devo editar" vs. "isso aqui é meu".

### O Fluxo: de SQL a Código Go

```
sql/migrations/*.sql  ──┐
  (schema: tabelas)      │
                         ├──► sqlc.yaml ──► `sqlc generate` ──► internal/db/*.go
sql/queries/query.sql ──┘        │                                 │
  (queries anotadas)             │                                 │
                          (arquivo de config             (código Go pronto para
                           que diz onde ler                importar e usar)
                           schema/queries e
                           onde escrever o
                           código gerado)
```

O SQLC lê o **schema** (`sql/migrations`) para saber quais tabelas e colunas existem e quais tipos elas têm. Depois lê as **queries anotadas** (`sql/queries/query.sql`) e, cruzando as duas informações, sabe exatamente que tipos Go usar nos parâmetros e no retorno de cada função gerada.

Importante: o SQLC **não conecta em nenhum banco de dados** para gerar código — ele só lê os arquivos `.sql` estaticamente. Quem efetivamente cria as tabelas no banco é a ferramenta `golang-migrate` (explicada mais abaixo), rodando à parte.

---

## 📖 Conceitos Fundamentais

### Migrations: o Schema do Banco

O SQLC precisa saber a "forma" das tabelas para gerar tipos corretos. Essa forma vem dos arquivos de migration em `sql/migrations/`:

```sql
-- sql/migrations/000001_init.up.sql
CREATE TABLE categories (
  id   varchar(36)  NOT NULL PRIMARY KEY,
  name text    NOT NULL,
  description  text
);

CREATE TABLE courses (
  id   varchar(36)  NOT NULL PRIMARY KEY,
  category_id   varchar(36)  NOT NULL,
  name text    NOT NULL,
  description  text,
  price  decimal(10,2)  NOT NULL,
  FOREIGN KEY (category_id) REFERENCES categories(id)
);
```

Repare que `description` não tem `NOT NULL` em nenhuma das duas tabelas — ou seja, essa coluna **aceita valores nulos**. Isso é importante e vai aparecer de novo mais adiante, quando falarmos de `sql.NullString`.

Cada migration tem um par de arquivos: `.up.sql` (aplica a mudança) e `.down.sql` (desfaz a mudança):

```sql
-- sql/migrations/000001_init.down.sql
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS categories;
```

Isso segue a convenção da ferramenta [`golang-migrate`](https://github.com/golang-migrate/migrate), que é quem realmente executa esses arquivos contra o banco (veja a seção [Rodando as Migrations](#rodando-as-migrations)). O SQLC apenas **lê** esses mesmos arquivos para entender o schema — ele não os executa.

### Queries Anotadas

Em `sql/queries/query.sql`, cada query é precedida por um comentário especial no formato `-- name: NomeDaFuncao :tipo`:

```sql
-- name: ListCategories :many
SELECT * FROM categories;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = ?;

-- name: CreateCategory :exec
INSERT INTO categories (id, name, description)
VALUES (?,?,?);

-- name: UpdateCategory :exec
UPDATE categories SET name = ?, description = ?
WHERE id = ?;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = ?;

-- name: CreateCourse :exec
INSERT INTO courses (id, name, description, category_id, price)
VALUES (?,?,?,?,?);

-- name: ListCourses :many
SELECT c.*, ca.name as category_name
FROM courses c JOIN categories ca ON c.category_id = ca.id;
```

O `:tipo` no final do comentário diz ao SQLC **que formato de retorno** gerar:

| Anotação | Significa | Assinatura Go gerada |
|----------|-----------|------------------------|
| `:exec` | A query não retorna linhas (INSERT/UPDATE/DELETE) | `func(ctx, ...) error` |
| `:one` | A query retorna exatamente uma linha | `func(ctx, ...) (Struct, error)` |
| `:many` | A query retorna zero ou mais linhas | `func(ctx, ...) ([]Struct, error)` |

Note também que os parâmetros usam `?` como placeholder — essa é a sintaxe do MySQL (outros bancos, como PostgreSQL, usariam `$1`, `$2`, etc.). Isso é definido pelo `engine` configurado em `sqlc.yaml`.

### sqlc.yaml — o Arquivo de Configuração

```yaml
version: "2"
sql:
- schema: "sql/migrations"
  queries: "sql/queries"
  engine: "mysql"
  gen:
    go:
      package: "db"
      out: "internal/db"
      overrides:
        - db_type: "decimal"
          go_type: "float64"
```

Explicando campo a campo:

- **`schema: "sql/migrations"`** — onde o SQLC vai procurar os arquivos que definem as tabelas.
- **`queries: "sql/queries"`** — onde estão os arquivos com as queries anotadas.
- **`engine: "mysql"`** — qual dialeto de SQL usar (afeta a sintaxe de placeholders, tipos, etc.). Aqui é MySQL, por isso o driver usado no projeto é o `go-sql-driver/mysql`.
- **`gen.go.package: "db"`** — o nome do pacote Go que será gerado (`package db`).
- **`gen.go.out: "internal/db"`** — a pasta onde os arquivos `.go` gerados serão escritos.
- **`overrides`** — regras para forçar um tipo Go específico em vez do padrão que o SQLC escolheria. Aqui, toda coluna `decimal` (como `price`) vira `float64` em Go, em vez de outro tipo que o SQLC usaria por padrão.

### Código Gerado: models.go

```go
// internal/db/models.go
// Code generated by sqlc. DO NOT EDIT.

package db

import "database/sql"

type Category struct {
	ID          string
	Name        string
	Description sql.NullString
}

type Course struct {
	ID          string
	CategoryID  string
	Name        string
	Description sql.NullString
	Price       float64
}
```

Cada tabela do schema vira uma `struct` Go com o mesmo nome (no singular, com a primeira letra maiúscula). Os nomes de coluna (`snake_case` no SQL) viram campos em `PascalCase` (convenção Go). E, como vimos, `description` é nula no banco, então vira `sql.NullString` em vez de `string`.

### Código Gerado: db.go e a Interface DBTX

```go
// internal/db/db.go
// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
```

Este é o arquivo mais importante para entender **como o SQLC se conecta de fato ao banco**. A peça central é a interface `DBTX`:

```go
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
```

Em Go, uma interface é satisfeita por qualquer tipo que implemente os métodos exigidos — não é preciso declarar explicitamente "eu implemento essa interface". E acontece que **tanto `*sql.DB` quanto `*sql.Tx`** (uma transação) já possuem exatamente esses quatro métodos. Ou seja, ambos "são" um `DBTX` automaticamente.

Isso é o que permite a função `New(db DBTX) *Queries` aceitar tanto uma conexão normal quanto uma transação:

```go
queries := db.New(dbConn)      // dbConn é *sql.DB → funciona
queries := db.New(tx)          // tx é *sql.Tx     → também funciona!
```

E é exatamente essa flexibilidade que o método `WithTx` (e o padrão usado em `runSQLCTX/main.go`) explora para rodar operações dentro de uma transação, como veremos adiante.

### Código Gerado: query.sql.go

Cada query anotada em `query.sql` vira uma função (e, quando necessário, um tipo `...Params`) neste arquivo. Vamos olhar três exemplos que representam os três padrões (`:exec`, `:one`, `:many`):

```go
// :exec — só retorna erro, porque INSERT/UPDATE/DELETE não têm linhas para ler
const createCategory = `-- name: CreateCategory :exec
INSERT INTO categories (id, name, description) 
VALUES (?,?,?)
`

type CreateCategoryParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) error {
	_, err := q.db.ExecContext(ctx, createCategory, arg.ID, arg.Name, arg.Description)
	return err
}
```

```go
// :one — usa QueryRowContext + Scan, retorna uma única struct
const getCategory = `-- name: GetCategory :one
SELECT id, name, description FROM categories 
WHERE id = ?
`

func (q *Queries) GetCategory(ctx context.Context, id string) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategory, id)
	var i Category
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}
```

```go
// :many — usa QueryContext + loop rows.Next(), retorna um slice
const listCategories = `-- name: ListCategories :many
SELECT id, name, description FROM categories
`

func (q *Queries) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, listCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
```

Repare que o texto SQL original é preservado como uma constante Go (`const listCategories = ...`) — não há mágica nenhuma acontecendo em runtime, é só `database/sql` puro, gerado por você.

Um detalhe interessante é a query `ListCourses`, que faz um `JOIN` e seleciona uma coluna extra (`category_name`) que não pertence a nenhuma tabela isoladamente. Nesse caso, o SQLC não reaproveita a struct `Course` — ele gera um tipo específico só para essa consulta:

```go
type ListCoursesRow struct {
	ID           string
	CategoryID   string
	Name         string
	Description  sql.NullString
	Price        float64
	CategoryName string
}

func (q *Queries) ListCourses(ctx context.Context) ([]ListCoursesRow, error) {
	// ...
}
```

Isso é uma vantagem do SQLC sobre um ORM tradicional: o tipo de retorno reflete **exatamente** as colunas selecionadas pela query, mesmo quando ela mistura dados de várias tabelas.

### sql.NullString — Lidando com Valores Nulos

Go não tem um jeito nativo de representar "uma string, ou nada" (diferente de outras linguagens que têm `null`/`nil` para qualquer tipo). O valor zero de `string` em Go é `""` (string vazia) — que é diferente de "não tem valor".

Como a coluna `description` pode ser `NULL` no banco, o SQLC gera o campo como `sql.NullString`, um tipo do pacote padrão `database/sql`:

```go
type NullString struct {
    String string
    Valid  bool // false quando o valor é NULL no banco
}
```

Na prática, isso significa que você sempre acessa o valor através de `.String`, e checa `.Valid` se precisar saber se o valor realmente existe:

```go
// ❌ Isso não compila — Description não é uma string, é sql.NullString
fmt.Println(category.Description)

// ✅ Acesso correto
fmt.Println(category.Description.String)

// ✅ Construindo um valor não-nulo para inserir
db.CreateCategoryParams{
	Description: sql.NullString{String: "Backend description", Valid: true},
}
```

---

## ⚙️ Como Funciona o Projeto

### cmd/runSQLC — Uso Básico do Queries Gerado

```go
// cmd/runSQLC/main.go
func main() {
    ctx := context.Background()
    dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
    if err != nil {
        panic(err)
    }
    defer dbConn.Close()

    queries := db.New(dbConn)

    err = queries.DeleteCategory(ctx, "361d935e-4ed1-49c3-afad-bdf5da95bbc9")
    if err != nil {
        panic(err)
    }

    categories, err := queries.ListCategories(ctx)
    if err != nil {
        panic(err)
    }

    for _, category := range categories {
        println(category.ID, category.Name, category.Description.String)
    }
}
```

O fluxo é simples: abrir a conexão com `sql.Open`, criar um `*db.Queries` com `db.New(dbConn)`, e então chamar os métodos gerados normalmente, como se fossem qualquer outra função Go.

O arquivo também traz, comentados, exemplos prontos de como chamar `CreateCategory` e `UpdateCategory` — eles não rodam de fato (estão dentro de `//`), mas servem como referência de uso caso você queira testar essas operações:

```go
// err = queries.CreateCategory(ctx, db.CreateCategoryParams{
//     ID: uuid.New().String(),
//     Name: "Backend",
//     Description: sql.NullString{String: "Backend description", Valid: true},
// })
```

> **Nota:** o arquivo usa `println` (a função *builtin* do Go, de baixo nível, que escreve direto no `stderr`) em vez de `fmt.Println`. Isso funciona para um exemplo rápido, mas não é recomendado em código real — veja a seção de [problemas](#6-código-de-exemplo-comentado) mais abaixo.

### cmd/runSQLCTX — Transações com callTx

Este exemplo mostra como usar transações por cima do código gerado, para garantir que duas operações relacionadas aconteçam **atomicamente** (ou as duas dão certo, ou nenhuma é aplicada).

```go
// cmd/runSQLCTX/main.go
type CourseDB struct {
	dbConn *sql.DB
	*db.Queries          // embedding: CourseDB "herda" todos os métodos de *db.Queries
}

func NewCourseDB(dbConn *sql.DB) *CourseDB {
	return &CourseDB{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}
```

O `*db.Queries` é embutido (embedding) dentro de `CourseDB`, então qualquer instância de `CourseDB` também tem acesso direto a `ListCourses`, `CreateCategory`, etc., além de guardar a conexão bruta (`dbConn`) para poder abrir transações.

O coração do padrão é o método privado `callTx`:

```go
func (c *CourseDB) callTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := c.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := db.New(tx)              // aqui o truque do DBTX aparece: New aceita *sql.Tx também
	err = fn(q)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return fmt.Errorf("error on rollback: %v, original error: %w", errRb, err)
		}
		return err
	}
	return tx.Commit()
}
```

`callTx` recebe uma função (`fn`) que sabe o que fazer com um `*db.Queries` — e é essa função que decide se a transação deve ser aceita (retornando `nil`) ou desfeita (retornando um erro). `callTx` cuida de todo o "esqueleto" da transação: abrir (`BeginTx`), criar um `Queries` ligado à transação (`db.New(tx)`), rodar a função, e então decidir entre `Rollback` ou `Commit`.

Isso é usado em `CreateCourseAndCategory`, que precisa criar uma categoria **e** um curso vinculado a ela como uma única operação atômica:

```go
func (c *CourseDB) CreateCourseAndCategory(ctx context.Context, argsCategory CategoryParams, argsCourse CourseParams) error {
	err := c.callTx(ctx, func(q *db.Queries) error {
		err := q.CreateCategory(ctx, db.CreateCategoryParams{ /* ... */ })
		if err != nil {
			return err
		}
		err = q.CreateCourse(ctx, db.CreateCourseParams{ /* ..., CategoryID: argsCategory.ID, ... */ })
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
```

Se a criação da categoria funcionar mas a criação do curso falhar (por exemplo, por um erro de conexão), a transação inteira é desfeita — não fica um curso "quebrado" apontando para nada, nem uma categoria órfã sem curso. É exatamente para evitar esse tipo de inconsistência que transações existem.

```
Sem transação:                       Com transação (callTx):
┌────────────────────┐               ┌────────────────────┐
│ CreateCategory ✅   │               │ BEGIN               │
│ CreateCourse   ❌   │               │  CreateCategory ✅  │
│                     │               │  CreateCourse   ❌  │
│ Resultado: categoria│               │ ROLLBACK             │
│ criada, curso não — │               │                      │
│ inconsistência!     │               │ Resultado: nada foi  │
└────────────────────┘               │ salvo — consistente  │
                                      └────────────────────┘
```

---

## ✅ Boas Práticas Presentes no Projeto

### 1. Código Gerado Nunca é Editado à Mão

**O problema:** se alguém editasse manualmente `internal/db/query.sql.go` para "ajustar" algo, essa mudança seria **perdida** na próxima vez que `sqlc generate` rodasse — silenciosamente.

```go
// ❌ Editar código gerado diretamente
// internal/db/query.sql.go
func (q *Queries) ListCategories(ctx context.Context) ([]Category, error) {
    // "só uma pequena correção aqui..."   ← será apagada no próximo `sqlc generate`
}

// ✅ A mudança certa é sempre no SQL de origem
// sql/queries/query.sql
-- name: ListCategories :many
SELECT * FROM categories ORDER BY name;  -- editar aqui, e rodar sqlc generate de novo
```

Todo arquivo gerado carrega o aviso `// Code generated by sqlc. DO NOT EDIT.` logo no topo — um sinal claro de que qualquer ajuste deve ser feito na fonte (SQL), nunca no resultado (Go).

### 2. Interface DBTX para Reuso em Transações

**O problema:** sem um mecanismo compartilhado, seria necessário duplicar cada função gerada — uma versão que aceita `*sql.DB` e outra que aceita `*sql.Tx`.

```go
// ❌ Duplicação — uma função para conexão normal, outra para transação
func (q *Queries) ListCategoriesDB(ctx context.Context, db *sql.DB) ([]Category, error) { /* ... */ }
func (q *Queries) ListCategoriesTx(ctx context.Context, tx *sql.Tx) ([]Category, error) { /* ... */ }

// ✅ Uma única implementação, graças à interface DBTX
type DBTX interface {
    ExecContext(...) (sql.Result, error)
    QueryContext(...) (*sql.Rows, error)
    QueryRowContext(...) *sql.Row
    PrepareContext(...) (*sql.Stmt, error)
}
// tanto *sql.DB quanto *sql.Tx implementam DBTX automaticamente
```

Isso permite que o mesmo `*db.Queries` funcione idêntico dentro ou fora de uma transação — é a base do padrão `callTx` visto em `runSQLCTX/main.go`.

### 3. Separação entre Schema, Queries e Configuração

**O problema:** misturar "o que existe no banco" com "o que eu quero consultar" dificultaria manter e regenerar o código de forma previsível.

```
sql/migrations/   → "O que existe" (tabelas, colunas, chaves estrangeiras)
sql/queries/      → "O que eu quero fazer" (SELECT, INSERT, UPDATE, DELETE)
sqlc.yaml         → "Como gerar" (dialeto, pacote, overrides de tipo)
```

Essa separação deixa cada arquivo com uma única responsabilidade e facilita saber onde mexer quando uma mudança é necessária.

### 4. Um Tipo por Formato de Resultado

**O problema:** reutilizar sempre a mesma struct de modelo (`Category`, `Course`) para toda consulta obrigaria a struct a ter campos "opcionais" para caber em queries com `JOIN` ou colunas calculadas.

```go
// ✅ Uma query com JOIN ganha seu próprio tipo, refletindo exatamente as colunas
type ListCoursesRow struct {
	ID           string
	CategoryID   string
	Name         string
	Description  sql.NullString
	Price        float64
	CategoryName string  // não existe em Course — só faz sentido para essa query
}
```

Isso evita structs "genéricas demais", com campos que só fazem sentido em algumas consultas.

---

## 🛡️ O que as Boas Práticas Evitaram

### Divergência entre Código e Banco

Se o código de acesso ao banco fosse escrito à mão, seria fácil ele "andar sozinho" e ficar dessincronizado do schema real (por exemplo, esquecer de atualizar um `Scan` depois de adicionar uma coluna). Como o SQLC gera o código a partir do schema real, qualquer mudança no banco que quebre uma query aparece como **erro de compilação**, não como bug em produção.

### Duplicação de Lógica de Transação

Sem a interface `DBTX`, cada função precisaria de duas versões (uma para `*sql.DB`, outra para `*sql.Tx`) ou algum tipo de wrapper manual complicado. A interface elimina essa duplicação por completo.

### Inconsistência de Dados entre Tabelas Relacionadas

Sem o padrão `callTx`, uma falha ao criar o curso depois de já ter criado a categoria deixaria dados inconsistentes no banco (categoria sem curso, ou pior, tentativas de retry criando categorias duplicadas). A transação garante que a operação composta seja tudo-ou-nada.

---

## 🚀 O que Poderia Ser Melhorado

### 1. Mover a DSN (Connection String) para Variável de Ambiente

Hoje a string de conexão está fixa no código:

```go
// Estado atual — credenciais hardcoded
dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
```

Poderia vir de uma variável de ambiente, evitando credenciais no código-fonte:

```go
dsn := os.Getenv("DATABASE_URL")
dbConn, err := sql.Open("mysql", dsn)
```

### 2. Adicionar `sqlc generate` como Target do Makefile

O `Makefile` atual só tem targets para migrations. Adicionar um target dedicado tornaria o fluxo mais explícito:

```makefile
generate:
	sqlc generate

.PHONY: generate migrate migratedown createmigration
```

### 3. Usar Contextos com Timeout

Todos os exemplos usam `context.Background()`, que nunca expira:

```go
// ❌ Nunca expira — uma query travada pode bloquear para sempre
ctx := context.Background()

// ✅ Com timeout — a query é cancelada após 5 segundos
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### 4. Tratar Erros em Vez de panic

```go
// ❌ Atual — encerra o processo abruptamente
dbConn, err := sql.Open("mysql", dsn)
if err != nil {
    panic(err)
}

// ✅ Melhor — log estruturado e saída controlada
dbConn, err := sql.Open("mysql", dsn)
if err != nil {
    log.Fatalf("erro ao conectar no banco: %v", err)
}
```

### 5. Testes com Banco de Teste

Como `Queries` depende apenas da interface `DBTX`, é possível criar um "fake" que implemente essa interface para testes sem precisar de um banco real, ou usar um banco de teste dedicado (ex.: MySQL em container efêmero) para testes de integração.

---

## ⚠️ Principais Problemas ao Trabalhar com SQLC

### 1. Esquecer de Rodar `sqlc generate`

**O problema:** depois de editar `sql/queries/query.sql` ou uma migration, é preciso rodar `sqlc generate` manualmente. Se você esquecer, o código em `internal/db/` fica desatualizado — e você pode nem perceber, já que o Go continua compilando com a versão antiga.

```bash
# Depois de qualquer mudança em sql/migrations ou sql/queries:
sqlc generate
# Sem isso, internal/db/*.go continua com a versão antiga do código
```

### 2. Migrations e Queries Fora de Sincronia

**O problema:** o SQLC só lê os arquivos `.sql` — ele não verifica se o banco de dados **real** está de fato com esse schema aplicado. Se você editar uma migration mas esquecer de rodar `make migrate` no banco, o código gerado vai assumir colunas que não existem de verdade no banco, e as queries vão falhar em runtime.

```
sql/migrations (arquivo)  →  sqlc generate  →  código Go assume esse schema
        │
        └─── mas o BANCO REAL só muda quando `make migrate` roda!
```

### 3. Placeholders Diferem entre Bancos

**O problema:** a sintaxe de parâmetros SQL não é a mesma em todo banco. Este projeto usa MySQL (`?`), mas PostgreSQL usa `$1, $2, ...`. Copiar uma query de um projeto Postgres para um projeto MySQL (ou vice-versa) sem ajustar os placeholders quebra a geração de código.

```sql
-- MySQL (este projeto)
WHERE id = ?;

-- PostgreSQL
WHERE id = $1;
```

O campo `engine` em `sqlc.yaml` precisa sempre bater com o banco real usado.

### 4. Nulidade Contamina Todo o Código

**O problema:** uma vez que uma coluna é nula no schema, todo lugar que usa esse campo precisa lidar com `sql.NullString` (ou `sql.NullInt64`, `sql.NullFloat64`, etc.) em vez do tipo Go nativo — o que é mais verboso.

```go
// ❌ Fácil de esquecer o .String e imprimir o struct inteiro
fmt.Println(category.Description)
// Saída: {Backend description true}  ← não é isso que você queria

// ✅ Sempre acessar explicitamente
fmt.Println(category.Description.String)
```

Se uma coluna nunca deveria ser nula na prática, vale a pena marcar `NOT NULL` na migration — isso simplifica o tipo Go gerado para `string` puro.

### 5. DSN com Credenciais no Código

**O problema:** como visto na seção de melhorias, a string de conexão (usuário, senha, host) está escrita diretamente no código-fonte (`root:root@tcp(localhost:3306)/courses`). Isso funciona para aprendizado, mas nunca deve ser feito em projetos reais — credenciais versionadas em Git são um risco de segurança.

### 6. Código de Exemplo Comentado

**O problema:** os arquivos `cmd/runSQLC/main.go` e `cmd/runSQLCTX/main.go` têm vários blocos de código comentados (`// err = queries.CreateCategory(...)`). Isso é proposital neste projeto de estudo — servem como referência de uso — mas em um projeto real, código morto comentado deve ser removido (o histórico do Git já guarda essas versões, não é necessário duplicar como comentário).

---

## 🔧 Como Usar o Projeto

### Pré-requisitos

```bash
# Go 1.19 ou superior
go version

# Docker (para subir o MySQL local)
docker --version

# SQLC CLI (para gerar o código a partir do SQL)
# Instalação via Go:
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# golang-migrate CLI (para aplicar as migrations)
# macOS (via Homebrew):
brew install golang-migrate
```

### Subindo o Banco de Dados

```bash
cd aulas/15-SQLC
docker-compose up -d
# Sobe um MySQL 5.7 com o banco "courses", usuário root, senha "root", na porta 3306
```

### Rodando as Migrations

```bash
# Aplica o schema (cria as tabelas categories e courses)
make migrate

# Caso precise desfazer
make migratedown

# Caso precise criar uma nova migration no futuro
make createmigration
```

### Gerando o Código com SQLC

```bash
# A partir da raiz de aulas/15-SQLC, onde está o sqlc.yaml
sqlc generate
# Gera/atualiza os arquivos em internal/db/
```

### Executando os Exemplos

```bash
# Baixar dependências
go mod tidy

# Exemplo de uso básico (Create/List/Delete de categorias)
go run ./cmd/runSQLC

# Exemplo de uso com transações (criar categoria + curso atomicamente)
go run ./cmd/runSQLCTX
```

---

## 📖 Glossário

| Termo | Definição |
|-------|-----------|
| **SQLC** | Ferramenta que gera código Go type-safe a partir de arquivos SQL (schema + queries) |
| **Code generation** | Técnica de gerar código-fonte automaticamente a partir de outra fonte de informação (aqui, SQL) |
| **Migration** | Arquivo que descreve uma mudança incremental no schema do banco de dados (criar/alterar/remover tabelas) |
| **DSN** | Data Source Name — string com as informações de conexão a um banco (usuário, senha, host, porta, database) |
| **DBTX** | Interface gerada pelo SQLC, satisfeita tanto por `*sql.DB` quanto por `*sql.Tx`, que permite reusar o mesmo código gerado dentro ou fora de transações |
| **Transação** | Conjunto de operações no banco que devem ser aplicadas todas juntas ou nenhuma (atomicidade) |
| **sql.NullString** | Tipo do pacote `database/sql` usado para representar uma coluna `string` que pode ser `NULL` no banco |
| **Placeholder** | Marcador usado numa query SQL para representar um parâmetro (`?` no MySQL, `$1` no PostgreSQL) |
| **ORM** | Object-Relational Mapper — biblioteca que mapeia objetos/structs para tabelas do banco automaticamente (ex.: GORM) |
| **Driver** | Pacote Go responsável por implementar a comunicação de baixo nível com um banco específico (ex.: `go-sql-driver/mysql`) |
| **Engine (SQLC)** | Configuração que define o dialeto SQL alvo (mysql, postgresql, sqlite) |
| **:exec / :one / :many** | Anotações usadas nas queries do SQLC para indicar o formato de retorno esperado |

---

## 🎯 Próximos Passos

### Para consolidar o aprendizado:

1. **Descomente e teste os exemplos** em `cmd/runSQLC/main.go` e `cmd/runSQLCTX/main.go` — crie uma categoria, atualize-a, e depois liste para ver o resultado.

2. **Adicione um target `generate` ao Makefile** — automatize a chamada de `sqlc generate` para não depender de lembrar o comando manualmente.

3. **Mova a DSN para uma variável de ambiente** — pratique lendo configuração via `os.Getenv` em vez de deixá-la fixa no código.

4. **Adicione uma nova query** — por exemplo, `-- name: CountCoursesByCategory :one` — e rode `sqlc generate` para ver o código novo aparecer automaticamente.

5. **Escreva um teste** para `CreateCourseAndCategory` usando um banco MySQL de teste (ou container efêmero), validando que uma falha no meio da transação realmente desfaz tudo.

6. **Experimente marcar uma coluna como `NOT NULL`** em uma nova migration e rode `sqlc generate` de novo — observe como o tipo gerado muda de `sql.NullString` para `string` puro.

### Conceitos relacionados no curso:

- **Aulas anteriores de banco de dados** — compare esta abordagem (SQLC) com o uso de um ORM, para entender as trocas entre controle explícito de SQL e produtividade automática de um ORM.
- **Aula 14 (Cobra CLI)** — usa um padrão DAO escrito à mão sobre SQLite; compare com o código gerado pelo SQLC sobre MySQL neste projeto.
- **Aula 7 (APIs)** — os métodos gerados aqui (`ListCategories`, `CreateCourse`, etc.) são exatamente o tipo de código que alimentaria os handlers de uma API REST.
