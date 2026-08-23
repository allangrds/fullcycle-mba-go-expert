# 💉 Injeção de Dependência em Go com Wire

Este módulo ensina o conceito de Injeção de Dependência (DI) em Go e como aplicá-lo usando o `google/wire`, uma ferramenta que gera código de injeção em **tempo de compilação**, sem reflection e sem container em runtime. Você vai entender por que DI existe, como o wire resolve o grafo de dependências do projeto e como esse código gerado se encaixa em uma arquitetura com Repository Pattern e Use Case.

---

## 📑 Sumário

- [🤔 O que é Injeção de Dependência?](#-o-que-é-injeção-de-dependência)
  - [A Analogia da Fábrica de Móveis](#a-analogia-da-fábrica-de-móveis)
  - [Sem DI vs Com DI](#sem-di-vs-com-di)
- [⚔️ DI Manual vs Container em Runtime vs Wire (Compile-time)](#️-di-manual-vs-container-em-runtime-vs-wire-compile-time)
  - [Tabela Comparativa](#tabela-comparativa)
- [📚 Conceitos Fundamentais do Wire](#-conceitos-fundamentais-do-wire)
  - [Provider](#provider)
  - [Provider Set — wire.NewSet](#provider-set--wirenewset)
  - [Bind — wire.Bind](#bind--wirebind)
  - [Injector — wire.Build](#injector--wirebuild)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
  - [product/entity.go — A Entidade](#productentitygo--a-entidade)
  - [product/repository.go — O Repository Pattern](#productrepositorygo--o-repository-pattern)
  - [product/usecase.go — O Use Case](#productusecasego--o-use-case)
  - [wire.go — O Molde do Injector](#wiego--o-molde-do-injector)
  - [wire_gen.go — O Código Gerado](#wire_gengo--o-código-gerado)
  - [main.go — Consumindo a Injeção](#maingo--consumindo-a-injeção)
- [🏷️ O Truque das Build Tags](#️-o-truque-das-build-tags)
- [▶️ Como Executar](#️-como-executar)
- [⚖️ Trade-offs: Wire vs Alternativas](#️-trade-offs-wire-vs-alternativas)
  - [Vantagens](#vantagens)
  - [Desvantagens](#desvantagens)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é Injeção de Dependência?

**Injeção de Dependência (DI)** é um padrão onde um objeto não cria as suas próprias dependências — ele as recebe prontas de fora, geralmente pelo construtor. Isso inverte o controle: em vez de o `ProductUseCase` decidir e instanciar qual `ProductRepository` usar, alguém de fora entrega essa dependência já pronta para ele.

### A Analogia da Fábrica de Móveis

```
Marceneiro Sem DI (faz tudo sozinho):
Recebe o pedido → corta a própria madeira → fabrica os próprios parafusos
                 → produz a própria cola → monta o móvel
Se precisar trocar o fornecedor de parafusos, precisa alterar o marceneiro inteiro.

Marceneiro Com DI (monta com peças entregues):
Recebe o pedido + madeira pronta + parafusos prontos + cola pronta
                 → apenas monta o móvel
Trocar o fornecedor de parafusos não muda em nada o trabalho do marceneiro.
```

O `ProductUseCase` é o marceneiro: ele só sabe **montar** a resposta a partir de um repositório que alguém já entregou pronto. Ele não sabe (nem precisa saber) se esse repositório usa SQLite, Postgres ou uma API externa.

### Sem DI vs Com DI

```go
// ❌ Sem DI: o usecase cria sua própria dependência (acoplamento forte)
type ProductUseCase struct {
    repository *ProductRepository
}

func NewProductUseCase(db *sql.DB) *ProductUseCase {
    repo := NewProductRepository(db) // decide e instancia sozinho
    return &ProductUseCase{repo}
}
```

```go
// ✅ Com DI: o usecase recebe a dependência já pronta, via interface
type ProductUseCase struct {
    repository ProductRepositoryInterface
}

func NewProductUseCase(repository ProductRepositoryInterface) *ProductUseCase {
    return &ProductUseCase{repository}
}
```

A segunda versão é a usada neste projeto. `ProductUseCase` depende apenas da **interface** `ProductRepositoryInterface` — quem decide qual implementação concreta entra ali é código externo (no nosso caso, o wire).

---

## ⚔️ DI Manual vs Container em Runtime vs Wire (Compile-time)

Existem três formas comuns de resolver dependências em Go:

1. **DI manual** — você mesmo escreve, à mão, as chamadas de construtor em cascata (`repo := New...(db); usecase := New...(repo)`), geralmente dentro do `main.go`.
2. **Container em runtime (reflection)** — bibliotecas como `uber-go/dig` e `uber-go/fx` mantêm um container que resolve o grafo de dependências em tempo de **execução**, usando reflection.
3. **Geração de código em tempo de compilação** — o `google/wire`, usado nesta aula, lê anotações no código-fonte e **gera** um arquivo `.go` comum com as mesmas chamadas em cascata que você escreveria manualmente — só que de forma automática e consistente.

### Tabela Comparativa

| Critério | DI Manual | Runtime (dig/fx) | Wire (Compile-time) |
|---|---|---|---|
| Quando resolve o grafo | Em compilação (é código normal) | Em execução, via reflection | Em compilação, via geração de código |
| Performance em runtime | Máxima (código puro) | Menor (overhead de reflection) | Máxima (código puro, igual ao manual) |
| Erros de dependência faltando | Compile-time | Runtime (panic ao rodar) | Compile-time |
| Legibilidade do código final | Alta, mas repetitiva em projetos grandes | Baixa (mágica escondida no container) | Alta (código gerado é legível, "parece" manual) |
| Ferramenta extra no fluxo | Não precisa | Não precisa (mas roda em runtime) | Precisa rodar `wire`/`go generate` |
| Escala bem para grafos grandes? | Fica repetitivo e propenso a erro | Sim, mas esconde a complexidade | Sim, e ainda mostra o resultado explícito |

O wire tenta ficar no melhor dos dois mundos: você declara as regras uma vez (quais providers existem, qual interface se liga a qual implementação) e ele gera o código repetitivo para você — sem pagar o custo de reflection em runtime.

---

## 📚 Conceitos Fundamentais do Wire

### Provider

Um **provider** é simplesmente uma função construtora comum — não tem nada de mágico. Qualquer função que recebe dependências e devolve um valor pronto pode ser um provider:

```go
func NewProductRepository(db *sql.DB) *ProductRepository {
    return &ProductRepository{db}
}

func NewProductUseCase(repository ProductRepositoryInterface) *ProductUseCase {
    return &ProductUseCase{repository}
}
```

O wire usa a assinatura dessas funções (o que elas recebem e o que devolvem) para descobrir como encadeá-las.

### Provider Set — `wire.NewSet`

`wire.NewSet` agrupa providers relacionados para que possam ser reutilizados em múltiplos injectors sem repetir a lista inteira:

```go
var setRepositoryDependency = wire.NewSet(
    product.NewProductRepository,
    wire.Bind(new(product.ProductRepositoryInterface), new(*product.ProductRepository)),
)
```

Aqui, `setRepositoryDependency` empacota "tudo que é preciso para produzir um `ProductRepositoryInterface`" em um único bloco reaproveitável.

### Bind — `wire.Bind`

Quando um provider retorna um tipo concreto (`*ProductRepository`) mas quem consome espera uma **interface** (`ProductRepositoryInterface`), o wire precisa ser instruído sobre essa relação — ele não infere isso sozinho:

```go
wire.Bind(new(product.ProductRepositoryInterface), new(*product.ProductRepository))
```

Isso diz: "sempre que algo pedir `ProductRepositoryInterface`, entregue o resultado do provider que produz `*ProductRepository`".

### Injector — `wire.Build`

O **injector** é a função que declara o que você quer construir. Dentro dela, `wire.Build` lista os providers/sets necessários:

```go
func NewUseCase(db *sql.DB) *product.ProductUseCase {
    wire.Build(
        setRepositoryDependency,
        product.NewProductUseCase,
    )
    return &product.ProductUseCase{}
}
```

O wire lê essa assinatura (`db *sql.DB` → `*product.ProductUseCase`), calcula o caminho `NewProductRepository(db)` → `NewProductUseCase(repo)`, e gera o código real que faz exatamente isso.

---

## 🗂️ Estrutura do Projeto

```
17-DI/
├── go.mod              # module github.com/devfullcycle/19-DI — deps: google/wire, mattn/go-sqlite3
├── go.sum
├── main.go             # ponto de entrada: abre o banco e chama NewUseCase(db) gerado pelo wire
├── wire.go             # molde do injector (build tag "wireinject") — nunca entra no binário final
├── wire_gen.go         # código gerado pelo wire (build tag "!wireinject") — o que de fato compila e roda
└── product/
    ├── entity.go        # struct Product — a entidade de domínio
    ├── repository.go    # ProductRepositoryInterface + ProductRepository (implementação)
    └── usecase.go        # ProductUseCase — depende apenas da interface do repository
```

Note que não há subpastas separando `entity/`, `repository/`, `usecase/` — todo o domínio de produto vive junto no pacote `product/`, um estilo mais simples e comum em exemplos didáticos.

---

## 🔍 Walkthrough do Código

A ordem abaixo segue o próprio grafo de dependências: da peça mais interna (entidade) até o ponto onde tudo é consumido (`main.go`).

### `product/entity.go` — A Entidade

```go
package product

type Product struct {
	ID   int
	Name string
}
```

A struct de domínio, sem nenhuma lógica. É o que trafega entre repository e usecase.

### `product/repository.go` — O Repository Pattern

```go
package product

import "database/sql"

type ProductRepositoryInterface interface {
	GetProduct(id int) (*Product, error)
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) GetProduct(id int) (*Product, error) {
	return &Product{
		ID:   id,
		Name: "Product Name",
	}, nil
}
```

Duas peças importantes aqui:
- `ProductRepositoryInterface` — o **contrato** que o resto da aplicação enxerga.
- `ProductRepository` — a implementação concreta, que guarda a `*sql.DB`. `NewProductRepository` é o **provider** que o wire vai usar para construir esse valor.

> `GetProduct` está propositalmente simplificado (retorna um produto fixo, sem consultar o SQLite de verdade) para manter o foco da aula em DI, não em persistência.

### `product/usecase.go` — O Use Case

```go
package product

type ProductUseCase struct {
	repository ProductRepositoryInterface
}

func NewProductUseCase(repository ProductRepositoryInterface) *ProductUseCase {
	return &ProductUseCase{repository}
}

// GetProduct returns a product by id
// This Product was not supposed to be returned. We should return a DTO instead.
// However, we will return it for now to keep the example simple.
func (u *ProductUseCase) GetProduct(id int) (*Product, error) {
	return u.repository.GetProduct(id)
}
```

`ProductUseCase` depende só da **interface**, nunca da struct concreta — é exatamente esse ponto de indireção que o `wire.Bind` precisa resolver. O comentário no código já assume, de forma didática, que devolver a entidade diretamente (em vez de um DTO) é uma simplificação proposital.

### `wire.go` — O Molde do Injector

```go
//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/devfullcycle/19-DI/product"
	"github.com/google/wire"
)

var setRepositoryDependency = wire.NewSet(
	product.NewProductRepository,
	wire.Bind(new(product.ProductRepositoryInterface), new(*product.ProductRepository)),
)

func NewUseCase(db *sql.DB) *product.ProductUseCase {
	wire.Build(
		setRepositoryDependency,
		product.NewProductUseCase,
	)
	return &product.ProductUseCase{}
}
```

Este arquivo **nunca é compilado no binário final** — ele existe só para a ferramenta `wire` ler e entender o que você quer construir. O `return &product.ProductUseCase{}` no final é morto (nunca executa); ele só está ali porque o compilador Go exige que toda função declare um retorno, mesmo que este arquivo jamais rode de verdade.

### `wire_gen.go` — O Código Gerado

```go
// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"database/sql"
	"github.com/devfullcycle/19-DI/product"
	"github.com/google/wire"
)

import (
	_ "github.com/mattn/go-sqlite3"
)

// Injectors from wire.go:

func NewUseCase(db *sql.DB) *product.ProductUseCase {
	productRepository := product.NewProductRepository(db)
	productUseCase := product.NewProductUseCase(productRepository)
	return productUseCase
}

// wire.go:

var setRepositoryDependency = wire.NewSet(product.NewProductRepository, wire.Bind(new(product.ProductRepositoryInterface), new(*product.ProductRepository)))
```

Este é o arquivo que **realmente compila e roda**. Repare que o corpo de `NewUseCase` é código Go comum, sem reflection, sem magia — é literalmente o mesmo encadeamento que você escreveria fazendo DI manual: `NewProductRepository(db)` → `NewProductUseCase(productRepository)`. A diferença é que o wire escreveu isso por você, de forma determinística, a partir das regras declaradas em `wire.go`.

### `main.go` — Consumindo a Injeção

```go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}

	usecase := NewUseCase(db)

	product, err := usecase.GetProduct(1)
	if err != nil {
		panic(err)
	}

	fmt.Println(product.Name)
}
```

Do ponto de vista de `main.go`, `NewUseCase(db)` é só mais uma função — não há nenhuma referência direta ao wire aqui. Essa é a ideia central: depois que o código é gerado, ele se torna indistinguível de DI manual escrita à mão.

---

## 🏷️ O Truque das Build Tags

`wire.go` e `wire_gen.go` declaram a **mesma função** `NewUseCase` e a **mesma variável** `setRepositoryDependency` — isso só é possível porque as build tags garantem que apenas um dos dois arquivos entra em cada build, nunca os dois ao mesmo tempo:

```
┌─────────────────────────────┐      ┌─────────────────────────────┐
│  wire.go                    │      │  wire_gen.go                │
│  //go:build wireinject      │      │  //go:build !wireinject      │
│                              │      │                              │
│  Lido apenas pela ferramenta │      │  Compilado no binário final  │
│  `wire` para gerar código.   │      │  em builds/execuções normais.│
│  NUNCA compilado no binário. │      │  É o que `go run .` executa. │
└─────────────────────────────┘      └─────────────────────────────┘
```

- Build normal (`go build`, `go run .`) → tag `wireinject` está **ausente** → só `wire_gen.go` entra.
- Comando `wire` → roda com a tag `wireinject` **ativa** → só `wire.go` é considerado, e o resultado é escrito em `wire_gen.go`.

É esse mecanismo que permite os dois arquivos coexistirem no mesmo pacote `main` sem colisão de símbolos.

---

## ▶️ Como Executar

```bash
# 1. Instalar a CLI do wire (uma vez só, globalmente)
go install github.com/google/wire/cmd/wire@latest

# 2. Entrar na pasta da aula
cd fullcycle-mba-go-expert/aulas/17-DI

# 3. (Re)gerar o wire_gen.go a partir de wire.go
wire
# ou, usando a diretiva go:generate presente em wire_gen.go:
go generate ./...

# 4. Rodar a aplicação
go run .
```

Saída esperada:
```
Product Name
```

O arquivo `test.db` é criado automaticamente pelo driver SQLite ao abrir a conexão — lembre que `GetProduct` não consulta esse banco de fato, apenas retorna um produto fixo com o `Name` "Product Name".

---

## ⚖️ Trade-offs: Wire vs Alternativas

### Vantagens

- ✅ Zero reflection em runtime — o código gerado é Go puro, com a mesma performance de DI manual.
- ✅ Dependência faltando ou tipo incompatível vira **erro de compilação** (via `wire`), não panic em produção.
- ✅ O código gerado é legível e depurável — você pode abrir `wire_gen.go` e ler exatamente o que acontece, sem "caixa preta".
- ✅ Reduz a repetição de escrever manualmente cada `New...()` em cascata conforme o grafo cresce.

### Desvantagens

- ❌ Exige uma ferramenta e um passo extra no fluxo de desenvolvimento (`wire` ou `go generate`) sempre que o grafo de dependências muda.
- ❌ Introduz dois arquivos por injector (`wire.go` + `wire_gen.go`) e a convenção de build tags, que é um conceito a mais para o time aprender.
- ❌ `wire_gen.go` é gerado e não deve ser editado manualmente — é preciso lembrar de sempre regenerar após alterar `wire.go`.
- ❌ Para grafos pequenos (como este exemplo), a DI manual seria quase tão simples de escrever à mão.

---

## 🎯 Casos de Uso Ideais

- **Aplicações com muitos serviços e repositórios**, onde escrever a cascata de construtores manualmente ficaria repetitivo e sujeito a erro.
- **Times que priorizam performance e querem evitar reflection em runtime**, mas ainda assim querem uma forma organizada de declarar dependências.
- **Projetos onde pegar erros de configuração de DI em tempo de compilação** (em vez de descobrir em produção que uma dependência não foi passada) é um requisito importante.
- Para protótipos pequenos ou aplicações com poucas dependências, DI manual direta no `main.go` costuma ser suficiente e mais simples que introduzir o wire.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **Injeção de Dependência (DI)** | Padrão onde um objeto recebe suas dependências prontas de fora, em vez de criá-las internamente. |
| **Provider** | Função construtora comum que o wire usa como peça para montar o grafo de dependências. |
| **Provider Set** (`wire.NewSet`) | Agrupamento reutilizável de providers e binds relacionados. |
| **Bind** (`wire.Bind`) | Instrução que diz ao wire qual implementação concreta satisfaz uma determinada interface. |
| **Injector** | Função anotada com `wire.Build` que declara o que deve ser construído; é a partir dela que o wire gera o código real. |
| **Build Tag** | Diretiva de compilação (`//go:build`) que inclui ou exclui um arquivo do build conforme a tag ativa. |
| **Repository Pattern** | Padrão que isola o acesso a dados atrás de uma interface, desacoplando as camadas superiores da fonte de dados. |
| **Use Case** | Camada que orquestra a lógica de aplicação, dependendo de abstrações (interfaces) em vez de implementações concretas. |
| **Compile-time DI** | Resolução do grafo de dependências feita em tempo de compilação (ex.: wire), gerando código estático. |
| **Runtime DI** | Resolução do grafo de dependências feita durante a execução, geralmente via reflection (ex.: dig, fx). |

---

## 🚀 Próximos Passos

**Imediato**
- [ ] Adicionar um segundo provider (ex.: um `Logger`) e observar como o wire resolve grafos maiores.
- [ ] Trocar a implementação de `GetProduct` para consultar o SQLite de verdade em vez de retornar um valor fixo.

**Intermediário**
- [ ] Criar múltiplos injectors no mesmo projeto (ex.: um para testes, com um repository em memória, e outro para produção).
- [ ] Aplicar a sugestão do comentário em `usecase.go`: introduzir um DTO em vez de retornar a entidade de domínio diretamente.
- [ ] Escrever testes para `ProductUseCase` usando um mock de `ProductRepositoryInterface`, aproveitando que a dependência já é injetada via interface.

**Avançado**
- [ ] Comparar esta abordagem com `uber-go/dig` ou `uber-go/fx` implementando o mesmo cenário, para sentir na prática a diferença entre DI em compile-time e em runtime.
- [ ] Explorar `wire.Value` e `wire.InterfaceValue` para injetar valores estáticos, além de providers.
- [ ] Integrar o wire a um projeto com múltiplas entidades e camadas (ex.: HTTP handlers → usecases → repositories), organizando provider sets por módulo.
