# Clean Architecture em Go — Guia Didático Completo

Este módulo ensina Clean Architecture (Arquitetura Limpa) usando um sistema de pedidos real como estudo de caso. Você vai entender por que separar código em camadas, como inverter dependências para proteger sua regra de negócio, e como o mesmo núcleo de aplicação consegue ser exposto simultaneamente por REST, gRPC e GraphQL — sem duplicar uma linha de lógica de negócio.

Se você é iniciante em programação e nunca ouviu falar em "arquitetura de software", não se preocupe: cada conceito aqui começa com uma analogia do dia a dia antes de qualquer termo técnico.

---

## 📑 Sumário

- [🤔 O que é Clean Architecture?](#-o-que-é-clean-architecture)
  - [A Analogia da Casa Bem Projetada](#a-analogia-da-casa-bem-projetada)
  - [O Problema que Isso Resolve](#o-problema-que-isso-resolve)
  - [Os Círculos Concêntricos de Uncle Bob](#os-círculos-concêntricos-de-uncle-bob)
  - [A Regra de Dependência](#a-regra-de-dependência)
- [⚔️ Camadas: da Mais Interna à Mais Externa](#️-camadas-da-mais-interna-à-mais-externa)
- [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
  - [Entity — A Regra de Negócio Mais Pura](#entity--a-regra-de-negócio-mais-pura)
  - [Interface — O Contrato que Inverte a Dependência](#interface--o-contrato-que-inverte-a-dependência)
  - [Use Case — A Regra de Negócio da Aplicação](#use-case--a-regra-de-negócio-da-aplicação)
  - [Repository — Onde os Dados Realmente Moram](#repository--onde-os-dados-realmente-moram)
  - [Controllers/Adapters — Três Portas para a Mesma Casa](#controllersadapters--três-portas-para-a-mesma-casa)
  - [Dependency Injection — Manual vs. Google Wire](#dependency-injection--manual-vs-google-wire)
  - [Eventos de Domínio — O Event Dispatcher Dentro da Arquitetura](#eventos-de-domínio--o-event-dispatcher-dentro-da-arquitetura)
  - [Context.Context — o Passageiro Invisível que Atravessa Todas as Camadas](#contextcontext--o-passageiro-invisível-que-atravessa-todas-as-camadas)
- [🛠️ Composition Root — Onde Tudo se Conecta](#️-composition-root--onde-tudo-se-conecta)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
  - [Os Três Controllers Lado a Lado](#os-três-controllers-lado-a-lado)
  - [Quem Importa Quem — O Fluxo Real de Dependências](#quem-importa-quem--o-fluxo-real-de-dependências)
- [▶️ Como Executar](#️-como-executar)
- [🎮 Testando as Três Portas de Entrada](#-testando-as-três-portas-de-entrada)
  - [REST](#rest)
  - [gRPC](#grpc)
  - [GraphQL](#graphql)
- [🎨 Padrões de Design](#-padrões-de-design)
- [🧬 Clean Architecture vs. Outras Abordagens](#-clean-architecture-vs-outras-abordagens)
  - [Clean Architecture vs. DDD (Domain-Driven Design)](#clean-architecture-vs-ddd-domain-driven-design)
  - [Clean Architecture vs. Arquitetura Hexagonal (Ports & Adapters)](#clean-architecture-vs-arquitetura-hexagonal-ports--adapters)
  - [Clean Architecture vs. MVC/Layered (Arquitetura em Camadas Tradicional)](#clean-architecture-vs-mvclayered-arquitetura-em-camadas-tradicional)
  - [Resumo: as Quatro Abordagens Lado a Lado](#resumo-as-quatro-abordagens-lado-a-lado)
- [🧩 Adotando Só Partes da Clean Architecture](#-adotando-só-partes-da-clean-architecture)
- [⚖️ Trade-offs: Convencional vs. Fora do Convencional](#️-trade-offs-convencional-vs-fora-do-convencional)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [💼 Clean Architecture no Mercado de Trabalho](#-clean-architecture-no-mercado-de-trabalho)
- [🧭 Inconsistências Reais do Projeto — Material de Estudo](#-inconsistências-reais-do-projeto--material-de-estudo)
- [🛠️ Comandos Úteis](#️-comandos-úteis)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é Clean Architecture?

**Clean Architecture** (Arquitetura Limpa) é uma forma de organizar o código de um sistema em **camadas concêntricas**, onde a regra de negócio fica isolada no centro e tudo que é "detalhe técnico" (banco de dados, protocolo HTTP, framework web) fica nas bordas. Foi popularizada pelo livro *Clean Architecture*, de Robert C. Martin (conhecido como "Uncle Bob").

### A Analogia da Casa Bem Projetada

Imagine duas formas de projetar uma casa:

**Casa Malfeita (tudo misturado):**
```
Você quer trocar a fiação elétrica do quarto...
→ mas os fios passam por dentro da parede da cozinha,
→ que está colada na tubulação de água do banheiro,
→ que compartilha a mesma caixa do gás.

Resultado: trocar UM fio de UM quarto exige mexer na casa inteira.
```

**Casa Bem Projetada (organizada por função):**
```
Elétrica, hidráulica e gás têm dutos próprios, independentes.
Quer trocar a fiação do quarto? Só mexe na elétrica daquele cômodo.
Quer trocar o material do piso? A estrutura da casa nem percebe.
```

No software, a "casa malfeita" é aquele código em que a regra de negócio ("como calculamos o preço final de um pedido") está espalhada dentro do handler HTTP, junto com SQL cru, junto com a lógica de validação — tudo no mesmo arquivo, tudo dependendo de tudo.

A "casa bem projetada" é o que a Clean Architecture propõe: cada responsabilidade no seu devido lugar, com fronteiras claras entre elas.

### O Problema que Isso Resolve

Todo sistema começa pequeno e simples. O problema aparece com o tempo, quando:

- Você precisa **trocar o banco de dados** (de MySQL para PostgreSQL) e descobre que SQL está espalhado em 40 arquivos diferentes.
- Você precisa **adicionar uma segunda forma de expor a mesma funcionalidade** (por exemplo, além de REST, agora precisa de gRPC) e percebe que teria que reescrever toda a lógica de negócio de novo, porque ela estava colada no handler HTTP.
- Um teste unitário simples da regra "o pedido não pode ter preço zero" exige subir um banco de dados real, porque a validação está misturada com o código que salva no banco.

Esse problema tem nome: **acoplamento**. Acoplamento é o quanto uma parte do código depende de detalhes de outra parte para funcionar. Quanto mais acoplado, mais caro e arriscado é mudar qualquer coisa — porque uma mudança pequena se propaga por lugares que você nem imaginava.

A Clean Architecture ataca o acoplamento organizando o código de forma que a parte mais importante (a regra de negócio) **não saiba nada** sobre os detalhes técnicos que a cercam.

### Os Círculos Concêntricos de Uncle Bob

O diagrama mais famoso da Clean Architecture é uma série de círculos, um dentro do outro:

```
        ┌───────────────────────────────────────────────┐
        │     FRAMEWORKS & DRIVERS (mais externo)        │
        │   Banco de dados, servidor HTTP, gRPC,         │
        │   GraphQL engine, RabbitMQ, bibliotecas         │
        │                                                  │
        │   ┌───────────────────────────────────────┐    │
        │   │      INTERFACE ADAPTERS                │    │
        │   │   Controllers, Gateways, Repositories  │    │
        │   │   (traduzem entre o mundo externo       │    │
        │   │    e o núcleo da aplicação)             │    │
        │   │                                          │    │
        │   │   ┌───────────────────────────────┐    │    │
        │   │   │       USE CASES                │    │    │
        │   │   │  Regras de negócio da           │    │    │
        │   │   │  aplicação (o que o sistema      │    │    │
        │   │   │  FAZ)                            │    │    │
        │   │   │                                    │    │    │
        │   │   │   ┌───────────────────────┐      │    │    │
        │   │   │   │      ENTITIES          │      │    │    │
        │   │   │   │  Regras de negócio      │      │    │    │
        │   │   │   │  mais gerais (o que o    │      │    │    │
        │   │   │   │  sistema É)              │      │    │    │
        │   │   │   └───────────────────────┘      │    │    │
        │   │   └───────────────────────────────┘    │    │
        │   └───────────────────────────────────────┘    │
        └───────────────────────────────────────────────┘
```

Do centro para fora:

1. **Entities** — as regras de negócio mais fundamentais, que fariam sentido mesmo se este software não existisse (ex.: "um pedido sem preço é inválido" é uma regra do *negócio*, não do *software*).
2. **Use Cases** — as regras de negócio específicas *desta aplicação* (ex.: "quando um pedido é criado, calculamos o preço final e disparamos uma notificação").
3. **Interface Adapters** — tradutores entre o núcleo da aplicação e o mundo externo (um controller HTTP, um repositório que fala com o banco).
4. **Frameworks & Drivers** — os detalhes técnicos concretos: qual banco de dados, qual framework web, qual biblioteca de mensageria.

### A Regra de Dependência

A regra mais importante de todo esse desenho: **o código-fonte só pode apontar para dentro**.

```
Frameworks & Drivers ──depende de──▶ Interface Adapters
                                              │
                                    depende de▼
                                         Use Cases
                                              │
                                    depende de▼
                                          Entities
```

Ou seja: um `import` de código em Go só pode ir de fora para dentro, nunca de dentro para fora. A camada de `Entities` não sabe que existe um banco de dados. A camada de `Use Cases` não sabe se a aplicação está sendo acessada via REST ou gRPC. Só as camadas externas sabem dos detalhes internos — nunca o contrário.

É exatamente essa regra que permite trocar um banco de dados, adicionar um novo protocolo de entrega, ou testar a regra de negócio isoladamente sem precisar subir infraestrutura nenhuma. Vamos ver como isso é aplicado, na prática, neste projeto.

---

## ⚔️ Camadas: da Mais Interna à Mais Externa

Este projeto implementa um sistema de pedidos (`Order`): você envia um `id`, um `price` e um `tax`, o sistema calcula o `final_price` (`price + tax`) e salva o pedido — expondo essa mesma operação por REST, gRPC e GraphQL, e notificando a criação via RabbitMQ.

Aqui está o mapeamento direto entre a teoria e as pastas reais deste projeto:

| Camada da Clean Architecture | O que é | Pasta neste projeto |
|---|---|---|
| **Entities** | Regras de negócio mais gerais e estáveis | `internal/entity/` |
| **Use Cases** | Regras de negócio específicas da aplicação | `internal/usecase/` |
| **Interface Adapters** | Controllers (REST/gRPC/GraphQL) e Gateways (Repository) | `internal/infra/web/`, `internal/infra/grpc/service/`, `internal/infra/graph/`, `internal/infra/database/` |
| **Frameworks & Drivers** | Servidor HTTP, driver de banco, engine GraphQL, RabbitMQ | `internal/infra/web/webserver/`, `internal/infra/grpc/pb/`, `database/sql`, `amqp`, `cmd/ordersystem/main.go` |

Existe ainda um pacote transversal, `pkg/events/`, que implementa um Event Dispatcher genérico — infraestrutura de domínio reutilizável, independente de `Order` (voltamos a ele mais adiante).

---

## 📚 Conceitos Fundamentais

Cada conceito abaixo segue o mesmo roteiro: uma analogia, a definição, o código real do projeto comentado, e uma discussão de trade-offs — incluindo como as coisas poderiam ser feitas de forma diferente.

### Entity — A Regra de Negócio Mais Pura

Pense na *entity* como a "certidão de nascimento" das regras do seu negócio: são fatos que seriam verdade mesmo se você decidisse reescrever o sistema inteiro em outra linguagem, com outro banco de dados, outro framework. "Um pedido sem preço não existe" é uma regra assim — não é uma regra de banco de dados, é uma regra do negócio.

Em Go, uma entity é normalmente uma `struct` com os dados e os métodos que garantem suas próprias regras. Veja a entity `Order` deste projeto, em [`internal/entity/order.go`](internal/entity/order.go):

```go
package entity

import "errors"

type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

// NewOrder é o "construtor": a única forma recomendada de criar
// um Order válido de fora do pacote entity.
func NewOrder(id string, price float64, tax float64) (*Order, error) {
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}
	err := order.IsValid() // valida antes de devolver — nunca retorna um Order inválido
	if err != nil {
		return nil, err
	}
	return order, nil
}

// IsValid garante os invariantes: sem essas condições, um Order não faz sentido.
func (o *Order) IsValid() error {
	if o.ID == "" {
		return errors.New("invalid id")
	}
	if o.Price <= 0 {
		return errors.New("invalid price")
	}
	if o.Tax <= 0 {
		return errors.New("invalid tax")
	}
	return nil
}

// CalculateFinalPrice é a regra de negócio central: preço final = preço + taxa.
func (o *Order) CalculateFinalPrice() error {
	o.FinalPrice = o.Price + o.Tax
	err := o.IsValid()
	if err != nil {
		return err
	}
	return nil
}
```

Repare que `Order` não importa nada do projeto além do pacote padrão `errors` do Go. Ela não sabe que existe um banco MySQL, nem que existe um endpoint HTTP. Ela só sabe suas próprias regras.

O teste dessa entity, em [`internal/entity/order_test.go`](internal/entity/order_test.go), reflete exatamente essa pureza:

```go
func TestGivenAnEmptyID_WhenCreateANewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{}
	assert.Error(t, order.IsValid(), "invalid id")
}

func TestGivenAPriceAndTax_WhenICallCalculatePrice_ThenIShouldSetFinalPrice(t *testing.T) {
	order, err := NewOrder("123", 10.0, 2.0)
	assert.Nil(t, err)
	assert.Nil(t, order.CalculateFinalPrice())
	assert.Equal(t, 12.0, order.FinalPrice)
}
```

**Por que isso importa:** este teste roda em milissegundos, não precisa de Docker, não precisa de rede, não precisa de mock nenhum. Ele testa 100% de lógica de negócio pura. Esse é um dos maiores ganhos práticos da Clean Architecture: a parte mais importante do seu sistema é também a mais fácil e rápida de testar.

#### E se fosse diferente? Entidade Rica vs. Entidade Anêmica

Este projeto usa o que se chama de **entidade rica** — uma struct que carrega tanto dados quanto comportamento (`IsValid`, `CalculateFinalPrice` vivem dentro de `Order`).

A alternativa é a **entidade anêmica**: uma struct que só guarda dados (`ID`, `Price`, `Tax`, `FinalPrice` como campos públicos, sem métodos), com toda a validação e cálculo vivendo em funções soltas ou em um "validator" separado.

| Aspecto | Entidade Rica (usada aqui) | Entidade Anêmica |
|---|---|---|
| Onde fica a regra "preço não pode ser zero" | Dentro do método `IsValid()` do `Order` | Numa função/validador externo |
| Risco de criar um objeto inválido "sem querer" | Baixo — `NewOrder` já valida | Maior — qualquer código pode montar a struct sem validar |
| Facilidade de achar "quem calcula o quê" | Alta — está tudo dentro da entidade | Menor — a lógica fica espalhada |
| Ponto fraco | Structs podem crescer muito se acumularem muita lógica | Fácil de virar um "saco de dados" sem regra nenhuma (perde a proteção do domínio) |

Não existe "certo" absoluto aqui — mas a entidade anêmica é considerada por muitos um *anti-pattern* quando usada sem cuidado, porque devolve o software ao problema original: a regra de negócio se esparrama por qualquer lugar que manipule aqueles dados.

#### Erros Tipados: uma Melhoria Real sobre o `errors.New` Genérico

Olhe de novo o `IsValid()` real deste projeto, já mostrado acima: cada regra quebrada devolve um `errors.New("invalid price")` — uma string solta. Quem chama `IsValid()` só consegue reagir de forma diferente a cada erro comparando o **texto exato** da mensagem (`err.Error() == "invalid price"`), o que quebra silenciosamente se alguém um dia reescrever essa string.

O padrão idiomático de Go para isso é o **erro sentinela**: declarar cada erro possível como uma variável exportada, uma vez só, e devolver essa mesma variável em todo lugar:

```go
package entity

import "errors"

var (
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidPrice = errors.New("invalid price")
	ErrInvalidTax   = errors.New("invalid tax")
)

func (o *Order) IsValid() error {
	if o.ID == "" {
		return ErrInvalidID
	}
	if o.Price <= 0 {
		return ErrInvalidPrice
	}
	if o.Tax <= 0 {
		return ErrInvalidTax
	}
	return nil
}
```

O ganho concreto: agora quem chama `IsValid()` faz `if errors.Is(err, entity.ErrInvalidPrice)` — comparando a **identidade** do erro, não o texto — e pode reagir de forma diferente por tipo (por exemplo, devolver HTTP `422 Unprocessable Entity` para `ErrInvalidPrice`, mas `400 Bad Request` para `ErrInvalidID`). Se algum dia a mensagem "invalid price" virar "preço inválido", nenhum `errors.Is` quebra, porque a comparação nunca dependeu do texto.

Se fosse ainda mais longe, e o erro precisasse carregar dados extras (ex.: qual foi o valor inválido recebido), o próximo passo seria um erro customizado com `errors.As`:

```go
import "fmt"

type ValidationError struct {
	Field string
	Value any
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("campo %s inválido: %v", e.Field, e.Value)
}
```

**Trade-off:** para validações triviais e projetos pequenos, `errors.New` genérico (como este projeto usa hoje) já resolve e é mais rápido de escrever — declarar uma variável sentinela para cada erro possível é código a mais. Mas assim que algum chamador precisa **distinguir** tipos de erro (decidir o código HTTP certo, decidir se vale tentar de novo, logar de forma diferente), o investimento em erros sentinela se paga rápido — e é um dos primeiros ajustes que revisores experientes de Go sugerem em code review.

#### Fora do Domínio Order: um Exemplo Ainda Mais Simples

Se `Order` ainda parece complexo demais pra "sentir" o conceito, veja este exemplo isolado, fora deste projeto — só pra ilustrar a ideia, não faz parte do código real:

```go
package domain

import "errors"

// Idade é uma entity minúscula: só um valor e uma regra.
type Idade struct {
	anos int
}

func NovaIdade(anos int) (Idade, error) {
	if anos < 0 || anos > 130 {
		return Idade{}, errors.New("idade inválida")
	}
	return Idade{anos: anos}, nil
}

func (i Idade) EhMaiorDeIdade() bool {
	return i.anos >= 18
}
```

Repare: nenhum `import` de banco, nenhum `http`, nenhum framework. `Idade` protege sozinha a regra "não existe idade negativa nem maior que 130" — é exatamente essa independência total de infraestrutura que faz de algo uma *entity*, seja ela um `Order` inteiro ou um único valor como esse.

**Não confunda com:** o "Model" que a maioria dos frameworks web gera automaticamente (Django Model, ActiveRecord do Rails, Model do Eloquent/Laravel) geralmente é uma **entidade anêmica** amarrada ao ORM — ela sabe salvar a si mesma no banco (`model.save()`), o que já é uma responsabilidade de infraestrutura vazando pra dentro do domínio. A *entity* da Clean Architecture não sabe salvar nada; ela só sabe se validar. Também não confunda com "Entity" no sentido específico do DDD (Domain-Driven Design) — lá, uma Entity é definida estritamente por ter **identidade única e persistente ao longo do tempo** (duas instâncias com os mesmos dados mas IDs diferentes são objetos diferentes), uma distinção mais rígida do que a Clean Architecture exige. Voltamos a esse ponto na seção [Clean Architecture vs. DDD](#clean-architecture-vs-ddd-domain-driven-design), mais adiante.

---

### Interface — O Contrato que Inverte a Dependência

Esta é, sem exagero, a peça mais importante para entender Clean Architecture — e também a mais contraintuitiva para quem está começando.

**Analogia:** pense numa tomada elétrica padrão. A tomada da parede não sabe, e não precisa saber, se você vai plugar nela uma TV, um carregador de celular ou uma cafeteira. Ela só define um *contrato* (o formato dos pinos, a voltagem). Qualquer aparelho que respeite esse contrato funciona ali — a parede nunca precisa ser refeita quando você troca de aparelho.

No código, esse "contrato" é uma **interface**. Em Go, uma interface é uma lista de métodos que um tipo precisa implementar — e, diferente de outras linguagens, em Go a implementação é *implícita*: você nunca escreve "esta struct implementa esta interface", o compilador simplesmente reconhece quando os métodos batem.

Veja o arquivo [`internal/entity/interface.go`](internal/entity/interface.go):

```go
package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	// GetTotal() (int, error)
}
```

Repare onde esse arquivo mora: dentro de `internal/entity`, a camada **mais interna** de todas. Isso é o coração do **Dependency Inversion Principle (DIP, Princípio da Inversão de Dependência)**: quem *declara* o contrato ("preciso de algo que salve um Order") é a camada de negócio — não a camada que vai implementar esse contrato de verdade.

Quem implementa essa interface é `internal/infra/database`, uma camada bem mais externa:

```go
// internal/infra/database/order_repository.go
type OrderRepository struct {
	Db *sql.DB
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	// ...
}
```

A "mágica" aqui é: o pacote `usecase` (que veremos a seguir) vai depender apenas de `entity.OrderRepositoryInterface` — nunca de `*database.OrderRepository` diretamente. Isso significa que:

- O use case não sabe se os dados vão para MySQL, PostgreSQL, um arquivo, ou a memória.
- Você pode trocar `OrderRepository` por qualquer outra implementação (inclusive uma "fake" em testes) sem tocar em uma linha do use case.
- A seta de **import** vai de fora para dentro (`infra/database` importa `entity`), mas a seta de **dependência conceitual** vai de dentro para fora (é a camada externa que se adapta ao contrato definido pela camada interna) — essa inversão é o próprio nome do princípio.

#### Um exemplo real de "dívida técnica" para aprender a enxergar

Olhando com atenção, você vai notar que `OrderRepository` (a implementação) tem um método `GetTotal()`:

```go
func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	// ...
}
```

Mas na interface `OrderRepositoryInterface`, esse mesmo método está **comentado**:

```go
type OrderRepositoryInterface interface {
	Save(order *Order) error
	// GetTotal() (int, error)
}
```

Isso é um exemplo real (não didaticamente forjado) de um "método órfão": ele existe na implementação, mas ninguém que depende só da interface consegue chamá-lo — porque o contrato não promete esse método. Não é um bug grave, mas é um ótimo exercício: pare e pense — se você precisasse usar `GetTotal()` a partir de um use case, o que faltaria fazer? (Resposta: descomentar o método na interface, garantindo que o contrato reflita o que realmente pode ser usado de fora.)

#### E se fosse diferente? Sem interfaces, com import direto

A alternativa "fora do convencional" (e desaconselhada pela Clean Architecture) seria o `usecase` importar `internal/infra/database` diretamente e usar `*database.OrderRepository` no lugar da interface. Funcionaria — Go compilaria sem problema — mas você perderia justamente a proteção que buscamos:

| Aspecto | Com interface (usado aqui) | Sem interface (import direto) |
|---|---|---|
| Trocar MySQL por outro banco | Só criar uma nova implementação da interface | Reescrever o use case inteiro |
| Testar o use case sem banco real | Fácil — injeta um "fake"/mock que implementa a interface | Difícil ou impossível sem subir infraestrutura |
| Camada `usecase` "sabe" que existe MySQL | Não | Sim — acoplamento direto |
| Custo de escrever | Um arquivo extra (`interface.go`) | Nenhum arquivo extra |

O "custo" de manter interfaces é real — mais um arquivo, mais um nível de indireção para ler. Mas em qualquer sistema que precise durar anos, trocar de banco, ou ser testado de forma confiável, esse custo se paga rápido.

#### Fora do Domínio Order: o Exemplo Clássico do Notifier

Este é provavelmente o exemplo mais citado em qualquer material sobre DIP, em qualquer linguagem — vale ver a versão mínima, fora deste projeto:

```go
package notification

// Notifier é o "port": o contrato que a camada de negócio conhece.
type Notifier interface {
	Send(destinatario, mensagem string) error
}

// EmailNotifier e SMSNotifier são "adapters": implementações concretas,
// intercambiáveis, que a camada de negócio nunca importa diretamente.
type EmailNotifier struct{}

func (e *EmailNotifier) Send(destinatario, mensagem string) error {
	// lógica de enviar e-mail de verdade aqui
	return nil
}

type SMSNotifier struct{}

func (s *SMSNotifier) Send(destinatario, mensagem string) error {
	// lógica de enviar SMS de verdade aqui
	return nil
}

// AlertService só conhece Notifier — pode receber Email, SMS,
// ou um NotifierFake em teste, sem nunca mudar uma linha aqui.
type AlertService struct {
	notifier Notifier
}

func (a *AlertService) AlertarUsuario(destinatario string) error {
	return a.notifier.Send(destinatario, "Algo importante aconteceu!")
}
```

Troque `EmailNotifier` por `SMSNotifier` (ou por um `NotifierFake` de teste) na hora de construir `AlertService`, e nada dentro de `AlertarUsuario` muda — essa é a prova mais direta e minimalista de DIP funcionando.

**Não confunda com:** "Dependency Inversion" (o princípio, sobre **quem declara o contrato** — aqui, é o pacote `notification` que define `Notifier`, não quem implementa) não é sinônimo de "Dependency Injection" (a técnica, sobre **como a implementação concreta chega** até quem usa — normalmente via construtor). São conceitos que andam juntos o tempo todo (como você vai ver na seção de [Dependency Injection](#dependency-injection--manual-vs-google-wire) logo mais abaixo), mas resolvem perguntas diferentes: DIP responde "quem é dono do contrato?"; DI responde "como a peça concreta é entregue?". É perfeitamente possível ter DIP sem DI automatizada (você pode construir tudo manualmente, como faz o `wire_gen.go` deste projeto) — os dois não são a mesma coisa.

---

### Use Case — A Regra de Negócio da Aplicação

Se a *entity* responde "o que é um pedido válido?", o *use case* responde "o que a aplicação **faz** quando alguém quer criar um pedido?". É a diferença entre uma regra do negócio em si e uma regra de como o seu sistema específico usa aquele negócio.

Veja [`internal/usecase/create_order.go`](internal/usecase/create_order.go) completo:

```go
package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
)

// DTOs (Data Transfer Objects): formatos de entrada e saída do use case.
// Note que NENHUM dos dois é a entity.Order — são formatos próprios.
type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface   // interface, não implementação concreta
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice() // delega o cálculo para a entity — o use case não recalcula nada sozinho

	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated) // efeito colateral delegado, não feito aqui dentro

	return dto, nil
}
```

Três pontos centrais para entender este arquivo:

1. **`CreateOrderUseCase` depende só de interfaces** (`entity.OrderRepositoryInterface`, `events.EventInterface`, `events.EventDispatcherInterface`) — nunca de uma implementação concreta. É a mesma inversão de dependência que vimos na seção anterior, aplicada de novo.
2. **DTOs, não a entity, cruzam a fronteira do use case.** `Execute` recebe um `OrderInputDTO` e devolve um `OrderOutputDTO` — nunca um `*entity.Order` diretamente.
3. **O use case não fala com RabbitMQ.** Ele só dispara um evento nomeado (`OrderCreated`) através do `EventDispatcher`. Quem decide o que acontece quando esse evento é disparado é decidido em outro lugar (veremos no Composition Root).

#### Por que DTOs e não a entity direto?

**Analogia:** pense num DTO como o "formulário de pedido" de uma loja, diferente do "processo interno de fabricação". O cliente preenche um formulário simples (nome, quantidade). A fábrica por dentro tem um processo bem mais detalhado, com etapas que o cliente nunca vê. Se você mudar o processo interno da fábrica, o formulário do cliente não precisa mudar — e vice-versa.

| Aspecto | DTO explícito (usado aqui) | Reaproveitar a entity como contrato externo |
|---|---|---|
| Acoplamento entre "o que entra/sai da API" e "o modelo de domínio" | Baixo — são tipos diferentes | Alto — qualquer mudança na entity muda o contrato externo |
| Velocidade para escrever o primeiro protótipo | Um pouco mais lenta (precisa declarar structs extras) | Mais rápida no início |
| Risco de vazar campos internos sensíveis | Baixo — você escolhe exatamente o que expor | Maior — fácil expor um campo interno sem querer |
| Facilidade de ter formatos diferentes por protocolo (REST vs gRPC vs GraphQL) | Alta — cada protocolo converte para o mesmo DTO comum | Baixa — a entity vira o "mínimo múltiplo comum" de todos os protocolos |

Este projeto claramente escolheu DTOs explícitos — e você vai ver esse padrão se repetir nos três controllers (REST, gRPC, GraphQL) mais adiante: todos convertem seu formato específico para `OrderInputDTO` antes de chamar `Execute`.

#### Um ponto de atenção: não há teste isolado do use case

Diferente da entity (que tem `order_test.go`) e do event dispatcher (que tem testes completos), não existe um `create_order_test.go` neste projeto. Isso significa que o comportamento de `CreateOrderUseCase.Execute` — "deve chamar `Save`, deve disparar o evento, deve montar corretamente o DTO de saída" — não está coberto por um teste automatizado isolado.

Como você testaria isso, se fosse escrever esse teste? A resposta usando os mesmos padrões já vistos no projeto: criar mocks de `entity.OrderRepositoryInterface` e `events.EventDispatcherInterface` (a mesma técnica usada em `pkg/events/event_dispatcher_test.go`, que veremos adiante), injetá-los no `CreateOrderUseCase`, e verificar se `Save` e `Dispatch` foram chamados com os dados certos — sem precisar de um banco de dados real rodando.

#### Testando o Use Case com Mocks

Isto não existe no repositório real — é a proposta de como `create_order_test.go` ficaria, seguindo exatamente a técnica descrita acima. Repare que o nome de teste segue o mesmo padrão `TestGivenX_WhenY_ThenZ` já usado em [`order_test.go`](internal/entity/order_test.go):

```go
package usecase

import (
	"testing"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/internal/event"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// OrderRepositoryMock implementa entity.OrderRepositoryInterface para o teste.
type OrderRepositoryMock struct{ mock.Mock }

func (m *OrderRepositoryMock) Save(order *entity.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

// EventDispatcherMock implementa events.EventDispatcherInterface.
// Register/Remove/Has/Clear ficam com implementação trivial abaixo —
// só Dispatch importa para este teste específico.
type EventDispatcherMock struct{ mock.Mock }

func (m *EventDispatcherMock) Dispatch(evt events.EventInterface) error {
	args := m.Called(evt)
	return args.Error(0)
}
func (m *EventDispatcherMock) Register(eventName string, handler events.EventHandlerInterface) error {
	return nil
}
func (m *EventDispatcherMock) Remove(eventName string, handler events.EventHandlerInterface) error {
	return nil
}
func (m *EventDispatcherMock) Has(eventName string, handler events.EventHandlerInterface) bool {
	return false
}
func (m *EventDispatcherMock) Clear() {}

func TestGivenAValidInput_WhenExecuteCreateOrder_ThenShouldSaveAndDispatch(t *testing.T) {
	// Arrange
	repoMock := &OrderRepositoryMock{}
	repoMock.On("Save", mock.Anything).Return(nil)

	dispatcherMock := &EventDispatcherMock{}
	dispatcherMock.On("Dispatch", mock.Anything).Return(nil)

	useCase := NewCreateOrderUseCase(repoMock, event.NewOrderCreated(), dispatcherMock)

	// Act
	output, err := useCase.Execute(OrderInputDTO{ID: "123", Price: 10.0, Tax: 2.0})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 12.0, output.FinalPrice)
	repoMock.AssertExpectations(t)       // prova que Save foi chamado
	dispatcherMock.AssertExpectations(t) // prova que Dispatch foi chamado
}
```

O ganho concreto: este teste prova que `Execute` orquestra corretamente `Save` + `Dispatch`, com o `FinalPrice` calculado certo — sem subir MySQL nem RabbitMQ, em milissegundos, do mesmo jeito que o teste da entity já faz. É a mesma ideia de "teste de integração leve" da seção de Repository, aplicada agora à orquestração do use case, só que aqui inteiramente com *mocks* (nada de banco real, nem fake) — porque o que se quer provar não é "os dados persistem corretamente", e sim "o use case chama as peças certas, na ordem certa, com os dados certos". *(veja também: [item 7 de Inconsistências Reais do Projeto](#-inconsistências-reais-do-projeto--material-de-estudo) e o checklist de [Próximos Passos](#-próximos-passos), que apontam para este mesmo teste faltando)*

#### Fora do Domínio Order: um Use Case Ainda Mais Enxuto

`CreateOrderUseCase` já tem DTOs, repositório e evento. Veja a forma mais reduzida possível da mesma ideia, fora deste projeto:

```go
package usecase

// User e UserRepository são o mínimo necessário para este exemplo
// (numa aplicação real, viriam do pacote entity).
type User struct {
	Nome  string
	Email string
}

type UserRepository interface {
	Save(u User) error
}

// RegisterUserUseCase: um caso de uso, uma única intenção do usuário.
type RegisterUserUseCase struct {
	repo UserRepository // interface, não implementação concreta
}

func NewRegisterUserUseCase(repo UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{repo: repo}
}

// Execute é o único método público — a assinatura clássica de um use case.
func (uc *RegisterUserUseCase) Execute(nome, email string) error {
	user := User{Nome: nome, Email: email}
	return uc.repo.Save(user)
}
```

Sem DTO explícito, sem evento, sem nada além do essencial: recebe dados simples, monta a entity, delega a persistência a uma interface. É o mesmo esqueleto de `CreateOrderUseCase`, só que sem os acréscimos que o domínio `Order` (múltiplos protocolos, notificação assíncrona) exige.

**Não confunda com:** um "Service" genérico, comum em muitas arquiteturas em camadas (`OrderService`, `UserService`), que costuma reunir **vários métodos relacionados** numa única classe (`Create`, `Update`, `Delete`, `List`, todos dentro do mesmo `UserService`). Um *Use Case* na Clean Architecture, no sentido mais estrito do termo, é o oposto: **uma única operação por classe/struct**, geralmente com um único método público (`Execute`) — cada intenção distinta do usuário vira seu próprio use case (`RegisterUserUseCase`, `DeactivateUserUseCase`, cada um em seu próprio arquivo). Essa granularidade fina é proposital: facilita achar exatamente onde uma regra de negócio específica vive, e evita que uma classe "guarda-chuva" cresça sem controle.

---

### Repository — Onde os Dados Realmente Moram

O **Repository Pattern** é a forma de isolar "como os dados são persistidos" de "o que a aplicação faz com os dados". Você já viu a metade "contrato" dele (`OrderRepositoryInterface`, na seção de Interfaces); agora vamos ver a implementação real, em [`internal/infra/database/order_repository.go`](internal/infra/database/order_repository.go):

```go
package database

import (
	"database/sql"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}
```

Repare que este projeto usa `database/sql` **puro**, sem ORM (Object-Relational Mapper). O SQL é escrito à mão (`INSERT INTO orders...`), com `Prepare`/`Exec` para evitar SQL injection através de parâmetros (`?`).

#### E se fosse diferente? SQL puro vs. ORM

| Aspecto | SQL puro (usado aqui) | ORM (ex.: GORM) |
|---|---|---|
| Controle sobre a query exata executada | Total | Parcial — o ORM decide boa parte do SQL gerado |
| Velocidade para fazer um CRUD simples | Mais lenta (você escreve cada query) | Mais rápida (métodos prontos tipo `db.Create(&order)`) |
| Curva de aprendizado | Baixa se você já sabe SQL | Precisa aprender a API do ORM |
| Performance em queries complexas | Você otimiza exatamente como quiser | Pode gerar queries subótimas sem perceber |
| Risco de "vazamento de abstração" | Baixo | Médio-alto (comportamentos "mágicos" do ORM podem surpreender) |

Não existe escolha universalmente certa. Projetos didáticos e sistemas com queries muito específicas tendem a preferir SQL puro (como aqui); sistemas com muito CRUD repetitivo costumam ganhar velocidade de desenvolvimento com um ORM.

#### O teste do Repository — integração com banco leve

Veja [`internal/infra/database/order_repository_test.go`](internal/infra/database/order_repository_test.go):

```go
func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:") // SQLite em memória, não o MySQL real
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	// ...
}
```

Esse é um **teste de integração leve**: ele não usa mock nenhum — testa o `Save` real, gravando e lendo de volta de um banco de verdade. Só que, em vez de exigir o MySQL do `docker-compose.yaml` rodando, ele usa SQLite `:memory:` — um banco que existe só durante a execução do teste e desaparece depois. Isso dá confiança real (é SQL de verdade rodando) sem o custo de subir infraestrutura pesada a cada `go test`.

**Ponto de atenção real:** este projeto não tem nenhum arquivo de *migration* SQL versionado. O `CREATE TABLE` do teste é definido só ali, manualmente; para o MySQL real (do `docker-compose.yaml`), é preciso criar a tabela `orders` manualmente antes de usar a aplicação — não há script automatizado para isso no repositório.

#### Fora do Domínio Order: o Exemplo Clássico de Trocar o Banco por um Fake

O ganho mais citado do Repository Pattern fica óbvio quando você vê a mesma interface satisfeita por duas implementações completamente diferentes, lado a lado:

```go
package user

import "database/sql"

type User struct {
	ID   string
	Name string
}

type UserRepository interface {
	Save(u User) error
	FindByID(id string) (User, error)
}

// InMemoryUserRepository: um mapa em memória — perfeito para testes rápidos.
type InMemoryUserRepository struct {
	data map[string]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{data: make(map[string]User)}
}

func (r *InMemoryUserRepository) Save(u User) error {
	r.data[u.ID] = u
	return nil
}

func (r *InMemoryUserRepository) FindByID(id string) (User, error) {
	return r.data[id], nil
}

// PostgresUserRepository: a implementação real de produção.
type PostgresUserRepository struct {
	db *sql.DB
}

func (r *PostgresUserRepository) Save(u User) error {
	_, err := r.db.Exec("INSERT INTO users (id, name) VALUES ($1, $2)", u.ID, u.Name)
	return err
}

func (r *PostgresUserRepository) FindByID(id string) (User, error) {
	// query real ao Postgres aqui
	return User{}, nil
}
```

Qualquer código que dependa só de `UserRepository` funciona igualmente bem com as duas — em teste, você usa `InMemoryUserRepository` (rápido, sem Docker); em produção, `PostgresUserRepository`. É exatamente o mesmo princípio que o SQLite `:memory:` do teste real deste projeto (`order_repository_test.go`) já aplica.

**Não confunda com:** o "Repository" genérico de frameworks como Spring Data (Java) ou Entity Framework (.NET), que muitas vezes já vem com métodos prontos de paginação, filtros complexos e até *query builders* embutidos — uma abstração bem mais rica (e mais acoplada ao ORM por trás) do que a Clean Architecture costuma recomendar. No sentido "puro" da Clean Architecture, como este projeto demonstra com `Save` sendo praticamente o único método do contrato, o repository é deliberadamente **burro**: só guarda e recupera, sem regra de negócio nem lógica de consulta sofisticada — qualquer coisa mais elaborada que isso pertence ao *use case*, não ao repository.

---

### Controllers/Adapters — Três Portas para a Mesma Casa

Um *controller* (ou *adapter*, no vocabulário da Clean Architecture) é o tradutor entre um protocolo externo específico (HTTP, gRPC, GraphQL) e o `usecase.OrderInputDTO`/`OrderOutputDTO` que já vimos. A ideia central: **o mesmo `CreateOrderUseCase` é chamado por três controllers diferentes**, cada um só cuidando da tradução de/para seu protocolo.

**REST**, em [`internal/infra/web/order_handler.go`](internal/infra/web/order_handler.go):

```go
func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto) // JSON HTTP → DTO
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output) // DTO → JSON HTTP
	// ...
}
```

**gRPC**, em [`internal/infra/grpc/service/order_service.go`](internal/infra/grpc/service/order_service.go):

```go
func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{ // protobuf → DTO
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{ // DTO → protobuf
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}
```

**GraphQL**, em [`internal/infra/graph/schema.resolvers.go`](internal/infra/graph/schema.resolvers.go):

```go
func (r *mutationResolver) CreateOrder(ctx context.Context, input *model.OrderInput) (*model.Order, error) {
	dto := usecase.OrderInputDTO{ // input GraphQL → DTO
		ID:    input.ID,
		Price: float64(input.Price),
		Tax:   float64(input.Tax),
	}
	output, err := r.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &model.Order{ // DTO → tipo GraphQL
		ID:         output.ID,
		Price:      float64(output.Price),
		Tax:        float64(output.Tax),
		FinalPrice: float64(output.FinalPrice),
	}, nil
}
```

Reparou no padrão se repetindo três vezes? **Formato externo → `OrderInputDTO` → `Execute` → `OrderOutputDTO` → formato externo.** Essa repetição não é acidente — é a prova visível de que a lógica de negócio (dentro de `Execute`) está genuinamente isolada dos protocolos. Trocar, adicionar ou remover um desses três controllers nunca exige tocar em `entity` nem em `usecase`.

#### Um exercício real: qual dos três está "mais certo"?

Repare numa diferença sutil entre eles. No REST, o handler **recria** o use case a cada requisição:

```go
createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
```

Já no gRPC, o `OrderService` **recebe o use case já pronto** no construtor, e reutiliza a mesma instância em toda chamada:

```go
type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase // já vem pronto
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase) *OrderService {
	return &OrderService{CreateOrderUseCase: createOrderUseCase}
}
```

O GraphQL segue o mesmo padrão do gRPC (o `Resolver` recebe o use case pronto).

Isso é uma inconsistência real deste projeto, não didaticamente inventada. Both funcionam, mas por quê a diferença importa? Recriar o use case a cada request (REST) tem um custo pequeno, mas desnecessário, de alocação — e quebra o padrão de injeção de dependência usado em todo o resto do projeto (você constrói a dependência uma vez, no lugar certo, e a reutiliza). Pare e pense: se você fosse arrumar isso, o que mudaria no `WebOrderHandler`? (Dica: ele passaria a receber `*usecase.CreateOrderUseCase` já pronto no construtor, assim como `OrderService` recebe.)

#### Fora do Domínio Order: um Adapter no seu Sentido Mais Puro

Tirando os três protocolos e todo o contexto de rede, um *adapter* é só isto — traduzir um formato para outro, sem lógica de negócio:

```go
package output

import "fmt"

type Resultado struct {
	Nome  string
	Total float64
}

// ToJSONOutput adapta Resultado para o formato que uma API espera.
func ToJSONOutput(r Resultado) map[string]interface{} {
	return map[string]interface{}{
		"nome":  r.Nome,
		"total": r.Total,
	}
}

// ToCLIOutput adapta o MESMO Resultado para o formato que um terminal espera.
func ToCLIOutput(r Resultado) string {
	return fmt.Sprintf("%s: R$ %.2f", r.Nome, r.Total)
}
```

`Resultado` nem sabe que existem esses dois formatos de saída — cada `ToXOutput` é um adapter dedicado, exatamente como `WebOrderHandler`, `OrderService` e `mutationResolver` são, cada um, um adapter dedicado ao redor do mesmo `CreateOrderUseCase.Execute`.

**Não confunda com:** o "Controller" do MVC (Model-View-Controller) tradicional carrega, historicamente, uma responsabilidade mais ampla — em implementações mais simples de MVC, é comum o controller já conter regra de validação e até lógica de negócio direto, além da tradução de protocolo. O "Controller"/"Adapter" da Clean Architecture é deliberadamente mais magro: sua única responsabilidade é traduzir formato de entrada/saída, nunca decidir regra de negócio — isso sempre fica no *use case*. Mesma palavra ("Controller"), responsabilidades bem diferentes dependendo da arquitetura; a comparação completa está na seção [Clean Architecture vs. MVC/Layered](#clean-architecture-vs-mvclayered-arquitetura-em-camadas-tradicional), mais adiante.

---

### Dependency Injection — Manual vs. Google Wire

**Dependency Injection (DI, Injeção de Dependência)** é um nome chique para uma ideia simples: em vez de um componente **criar** suas próprias dependências internamente, ele **recebe** essas dependências já prontas, de fora.

**Analogia:** pense num restaurante. O cozinheiro não planta os próprios vegetais nem cria a própria panela — ele recebe os ingredientes e os utensílios prontos, entregues por fornecedores. Isso permite trocar o fornecedor de vegetais sem reformar a cozinha inteira, e permite que, num teste (um "treinamento" do cozinheiro), você use ingredientes de mentira sem afetar a cozinha real.

Praticamente todo construtor que vimos até aqui já faz DI manual — por exemplo, `NewCreateOrderUseCase(OrderRepository, OrderCreated, EventDispatcher)` recebe as três dependências prontas em vez de criá-las internamente.

#### DI manual, escrita à mão

Fazer isso manualmente significa escrever, num único lugar, todo o código que instancia e conecta as peças:

```go
orderRepository := database.NewOrderRepository(db)
orderCreated := event.NewOrderCreated()
createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreated, eventDispatcher)
```

Isso funciona bem em projetos pequenos. Em projetos grandes, com muitas dependências e muitas combinações, esse código de "fiação" (conectar tudo) cresce e fica repetitivo e chato de manter à mão.

#### DI via geração de código: Google Wire

Este projeto usa o [Google Wire](https://github.com/google/wire) para gerar automaticamente esse código de fiação. Veja [`cmd/ordersystem/wire.go`](cmd/ordersystem/wire.go) — o arquivo que você escreve manualmente, mas que **não é compilado normalmente** (só quando a tag de build `wireinject` está ativa; ele existe só para o Wire ler):

```go
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/internal/event"
	"github.com/devfullcycle/20-CleanArch/internal/infra/database"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}
```

A linha mais importante aqui é `wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository))`. Ela diz, em código: *"sempre que alguém pedir a interface `OrderRepositoryInterface`, injete um `*database.OrderRepository` de verdade."* É exatamente aqui, na borda mais externa do sistema, que a abstração (a interface declarada dentro de `entity`) se casa com a implementação concreta (a struct declarada dentro de `infra/database`).

O comando `wire` (uma ferramenta de linha de comando) lê esse arquivo e **gera** [`cmd/ordersystem/wire_gen.go`](cmd/ordersystem/wire_gen.go) — o arquivo que de fato é compilado e roda em produção:

```go
// Code generated by Wire. DO NOT EDIT.

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreated, eventDispatcher)
	return createOrderUseCase
}
```

Note bem: o código gerado é **Go comum, sem mágica nenhuma** — nada de reflection em tempo de execução, nada de container carregando tipos dinamicamente. O Wire só automatiza, em tempo de compilação (build time), a escrita desse código repetitivo que você mesmo poderia ter digitado à mão. É por isso que muita gente descreve o Wire como "DI sem custo de performance": o resultado final é indistinguível, para o Go runtime, de você ter escrito tudo manualmente.

#### Trade-offs: manual vs. Wire vs. container com reflection

| Abordagem | Como funciona | Vantagem | Desvantagem |
|---|---|---|---|
| **DI manual** | Você escreve à mão todo o código de fiação | Simples de entender, zero ferramentas extras | Repetitivo e cansativo em projetos grandes |
| **DI com codegen (Wire, usado aqui)** | Uma ferramenta gera o código de fiação em build time | Reduz repetição, mantém performance de Go puro | Mais uma ferramenta para aprender; erros de configuração só aparecem ao rodar `wire` |
| **DI com container/reflection** (comum em outras linguagens/frameworks) | Um container resolve dependências em tempo de execução, usando reflection | Muito flexível, pouca configuração explícita | Custo de performance em runtime; erros de dependência faltando só aparecem quando o código roda, não em tempo de compilação |

Este projeto claramente prioriza performance e checagem em tempo de compilação (Wire), em vez de flexibilidade dinâmica (container com reflection) — uma escolha comum no ecossistema Go, que valoriza simplicidade e previsibilidade.

#### Fora do Domínio Order: o Exemplo Clássico do "Clock Injetável"

Antes de qualquer Wire ou container, veja a forma mais reduzida possível de DI — e um dos exemplos mais citados em Go, porque resolve um problema real de testabilidade:

```go
package billing

import "time"

// Clock é o contrato: "algo que sabe que horas são agora".
type Clock interface {
	Now() time.Time
}

// RealClock é a implementação de produção.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// Sem DI: InvoiceService cria sua própria dependência internamente.
// Problema: todo teste depende do relógio real da máquina — impossível
// testar "o que acontece no dia 31 de dezembro" de forma determinística.
func NewInvoiceServiceSemDI() *InvoiceService {
	return &InvoiceService{clock: RealClock{}}
}

// Com DI: InvoiceService recebe o Clock pronto, de fora.
func NewInvoiceServiceComDI(clock Clock) *InvoiceService {
	return &InvoiceService{clock: clock}
}

type InvoiceService struct {
	clock Clock
}
```

Em teste, você injeta um `FakeClock` que sempre devolve uma data fixa — e de repente "testar o que acontece no fim do ano" deixa de ser um problema de esperar o calendário virar e vira uma linha de código.

**Não confunda com:** o **Service Locator**, um padrão relacionado mas geralmente considerado um anti-pattern quando comparado à DI — em vez de receber a dependência pronta no construtor (como `NewInvoiceServiceComDI(clock)` faz), o componente **pede** a dependência a um registro global (`locator.Get("clock")`). A diferença parece pequena, mas importa: com Service Locator, as dependências reais de um componente ficam escondidas dentro do seu código (você só descobre lendo linha a linha o que ele pede ao locator); com DI via construtor, como este projeto usa em todo lugar, as dependências são explícitas na assinatura da função — visíveis de fora, sem precisar ler a implementação.

---

### Eventos de Domínio — O Event Dispatcher Dentro da Arquitetura

Se você já viu a [Aula 9 (Eventos)](../9-eventos/README.md) deste curso, vai reconhecer imediatamente o pacote [`pkg/events/`](pkg/events/) deste projeto: é o mesmo Event Dispatcher construído lá, peça por peça — as interfaces `EventInterface`/`EventHandlerInterface`/`EventDispatcherInterface`, o `map` usado como registry, as goroutines com `sync.WaitGroup` no `Dispatch`, e a publicação em RabbitMQ. Se você ainda não viu a Aula 9, vale a pena revisar por lá a construção completa, passo a passo — aqui vamos direto ao ponto: **como esse mecanismo já pronto se encaixa dentro das camadas da Clean Architecture.**

O pacote `pkg/events` define só os contratos (abstrações puras, sem saber nada de `Order` nem de RabbitMQ):

```go
// pkg/events/interface.go
type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}
```

Já as peças específicas do domínio `Order` vivem em `internal/event/` — uma camada mais externa que *usa* as abstrações de `pkg/events`:

```go
// internal/event/order_created.go — o evento concreto
type OrderCreated struct {
	Name    string
	Payload interface{}
}

func NewOrderCreated() *OrderCreated {
	return &OrderCreated{Name: "OrderCreated"}
}
```

```go
// internal/event/handler/order_created_handler.go — o handler concreto, que fala com RabbitMQ
func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	jsonOutput, _ := json.Marshal(event.GetPayload())

	h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // routing key
		false, false,
		amqp.Publishing{ContentType: "application/json", Body: jsonOutput},
	)
}
```

O ponto-chave arquitetural: relembrando o `CreateOrderUseCase.Execute` que já vimos, ele só conhece `events.EventDispatcherInterface` e `events.EventInterface` — **nunca** importa `internal/event/handler` nem o pacote `amqp` do RabbitMQ. O use case dispara um evento chamado `"OrderCreated"` e não faz ideia do que vai acontecer depois — se vai virar uma mensagem no RabbitMQ, um e-mail, ou nada.

Quem decide o que realmente acontece quando `"OrderCreated"` é disparado é o **Composition Root** (`main.go`), amarrando o handler concreto ao dispatcher:

```go
eventDispatcher := events.NewEventDispatcher()
eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
	RabbitMQChannel: rabbitMQChannel,
})
```

Essa é a mesma inversão de dependência que já vimos com o Repository, agora aplicada a efeitos colaterais: o "o quê" (disparar um evento) fica dentro; o "como reagir a isso" (publicar em fila) fica fora, e só é amarrado na borda mais externa do sistema.

#### Trade-off real do projeto: sem garantia transacional

O `Dispatch` deste projeto acontece **depois** que `OrderRepository.Save` já terminou com sucesso:

```go
if err := c.OrderRepository.Save(&order); err != nil {
	return OrderOutputDTO{}, err
}
// ...
c.EventDispatcher.Dispatch(c.OrderCreated)
```

Isso significa que é possível salvar o pedido no banco com sucesso e, na sequência, falhar ao publicar no RabbitMQ (por exemplo, se a conexão cair naquele instante) — sem nenhum mecanismo de rollback ou nova tentativa. O pedido fica salvo, mas ninguém é notificado.

Esse é um trade-off real e comum: um **Event Dispatcher em memória, síncrono por chamada** (como este) é simples de entender e rápido de implementar, mas não garante que "salvar" e "notificar" aconteçam como uma coisa só. A alternativa mais robusta para isso, em sistemas de produção sérios, é o **Outbox Pattern** — salvar o evento na mesma transação de banco que salva o dado, e ter um processo separado que garante a entrega da mensagem, com retries, a partir dessa tabela. Fica como tópico "para ir além" no glossário e nos próximos passos.

#### Fora do Domínio Order: um Broadcaster de Dez Linhas

Sem goroutine, sem `WaitGroup`, sem RabbitMQ — só a ideia pura do padrão Observer, isolada:

```go
package pubsub

type Subscriber func(mensagem string)

type Broadcaster struct {
	subscribers []Subscriber
}

func (b *Broadcaster) Subscribe(s Subscriber) {
	b.subscribers = append(b.subscribers, s)
}

func (b *Broadcaster) Publish(mensagem string) {
	for _, s := range b.subscribers {
		s(mensagem) // chamada direta, síncrona — sem goroutine aqui
	}
}
```

`Broadcaster` não sabe (nem precisa saber) o que cada `Subscriber` faz — poderia ser logar, enviar e-mail, atualizar uma métrica. Essa é a essência do Event Dispatcher deste projeto (`pkg/events.EventDispatcher`), só que sem as camadas extras de concorrência que a versão real acrescenta pra rodar os handlers em paralelo.

**Não confunda com:** um **Command** (Comando) — em CQRS e em DDD, essa é uma distinção clássica. Um *Domain Event* como `OrderCreated` descreve **algo que já aconteceu**, um fato consumado, sempre nomeado no passado ("Pedido Criado") — quem recebe o evento só pode reagir a ele, nunca rejeitá-lo. Um *Command* como `CreateOrder` descreve **uma intenção**, algo que ainda vai acontecer, nomeado no imperativo ("Criar Pedido") — e pode ser rejeitado (validação falhar, regra de negócio impedir). Neste projeto, o próprio nome `OrderCreated` (não `CreateOrder`) já segue essa convenção corretamente: é um evento, não um comando.

---

### Context.Context — o Passageiro Invisível que Atravessa Todas as Camadas

**Analogia:** pense num crachá de visitante numa empresa. Desde a portaria, o crachá carrega quem é o visitante e até que horas a visita é válida — qualquer sala que ele entrar pode checar o crachá e decidir "sua visita já venceu, não posso mais te atender" sem precisar telefonar de volta pra portaria. `context.Context`, em Go, é esse crachá: um valor que atravessa toda a cadeia de chamadas de uma requisição, carregando prazos, sinais de cancelamento, e pequenos metadados — sem que cada função no meio do caminho precise saber os detalhes de onde ele veio.

Você já viu `context.Context` passando batido, sem comentário nenhum, em dois lugares deste próprio README. No controller gRPC:

```go
func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
```

E no resolver GraphQL:

```go
func (r *mutationResolver) CreateOrder(ctx context.Context, input *model.OrderInput) (*model.Order, error) {
```

Repare: os dois controllers **já recebem** um `ctx` pronto, entregue automaticamente pelos frameworks gRPC e GraphQL. Mas nenhum dos dois faz nada com ele — e olhando a assinatura de `CreateOrderUseCase.Execute` (já vista na seção de [Use Case](#use-case--a-regra-de-negócio-da-aplicação)):

```go
func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
```

...não existe parâmetro `ctx` nenhum. O "crachá" chega na portaria (o controller), mas ninguém o carrega para dentro do prédio (o use case, o repository). Se o cliente gRPC cancelar a chamada no meio do caminho, ou um timeout de API gateway estourar, o `CreateOrderUseCase` segue rodando até o fim, sem jeito nenhum de saber que ninguém mais está esperando pela resposta.

#### Como ficaria propagado

```go
func (c *CreateOrderUseCase) Execute(ctx context.Context, input OrderInputDTO) (OrderOutputDTO, error) {
	select {
	case <-ctx.Done():
		return OrderOutputDTO{}, ctx.Err() // cliente desistiu ou o prazo já estourou
	default:
	}

	order := entity.Order{ID: input.ID, Price: input.Price, Tax: input.Tax}
	order.CalculateFinalPrice()

	if err := c.OrderRepository.Save(&order); err != nil { // idealmente: c.OrderRepository.Save(ctx, &order)
		return OrderOutputDTO{}, err
	}
	// ... resto igual
}
```

Numa versão completa, `ctx` seguiria até dentro do repository, chegando em `r.Db.ExecContext(ctx, "INSERT INTO orders ...")` no lugar do `r.Db.Exec(...)` atual (visto na seção de [Repository](#repository--onde-os-dados-realmente-moram)) — aí sim a cadeia inteira, do gRPC/GraphQL até o banco, respeita um cancelamento vindo de fora.

**Não confunda com:** `context.Context` **não** é um lugar para passar dados de negócio. É tentador usar `ctx.WithValue(ctx, "input", dto)` para "economizar" um parâmetro — mas isso quebra a tipagem estática do Go (o valor sai como `interface{}`, exige type assertion pra recuperar) e esconde uma dependência real dentro de um mecanismo pensado para outra coisa. O uso correto de `ctx.Value()` fica restrito a metadados de escopo da requisição (um `request ID` para rastreamento/tracing, por exemplo) — nunca para os parâmetros de negócio normais, que continuam trafegando explicitamente, exatamente como `OrderInputDTO` já faz.

**Trade-off:** propagar `ctx` por toda a cadeia (controller → use case → repository) é mecânico, mas tem um custo real: toda assinatura de método no caminho ganha um parâmetro a mais, e é fácil esquecer de repassá-lo em algum ponto no meio. Compensa a partir do momento em que a aplicação precisa reagir a cancelamento/timeout de verdade — uma API pública exposta à internet, por exemplo. Em código Go revisado por outros desenvolvedores, um parâmetro `context.Context` ausente numa cadeia de chamadas de I/O (banco, rede, fila) costuma ser um dos primeiros sinais que um revisor mais experiente aponta.

---

## 🛠️ Composition Root — Onde Tudo se Conecta

**Composition Root** é o nome técnico para "o único lugar do sistema onde tudo é instanciado e conectado". É a camada mais externa de todas — `Frameworks & Drivers` — materializada em [`cmd/ordersystem/main.go`](cmd/ordersystem/main.go):

```go
func main() {
	configs, err := configs.LoadConfig(".")           // 1. carrega configuração (Viper)
	// ...

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf(  // 2. abre conexão MySQL
		"%s:%s@tcp(%s:%s)/%s",
		configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()             // 3. abre canal RabbitMQ

	eventDispatcher := events.NewEventDispatcher()       // 4. cria dispatcher e registra o handler concreto
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher) // gerado pelo Wire

	webserver := webserver.NewWebServer(configs.WebServerPort)      // 5. sobe REST
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	go webserver.Start()

	grpcServer := grpc.NewServer()                                  // 6. sobe gRPC
	createOrderService := service.NewOrderService(*createOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{ // 7. sobe GraphQL (bloqueante)
		Resolvers: &graph.Resolver{CreateOrderUseCase: *createOrderUseCase},
	}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}
```

Repare que os três servidores rodam **no mesmo processo**, simultaneamente: REST e gRPC são iniciados com `go` (em goroutines, não bloqueando o resto de `main`), e o GraphQL roda de forma bloqueante por último — mas isso é só uma escolha de ordem, não uma dependência real entre eles. Todos os três, ao final, chamam o mesmo `createOrderUseCase`.

`main.go` é o **único** arquivo do projeto que importa simultaneamente `entity`, `event`, `infra/database`, `infra/web`, `infra/grpc`, `infra/graph`, `usecase`, `pkg/events` e `configs`. Isso é esperado e correto: é exatamente o papel do Composition Root amarrar tudo — e é exatamente por isso que nenhuma outra camada precisa fazer isso.

### Configuração via Viper

[`configs/config.go`](configs/config.go) usa a biblioteca [Viper](https://github.com/spf13/viper) para ler o arquivo `.env`:

```go
type conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	// ...
}
```

E o `.env` real do projeto ([`cmd/ordersystem/.env`](cmd/ordersystem/.env)):

```
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders
WEB_SERVER_PORT=:8000
GRPC_SERVER_PORT=50051
GRAPHQL_SERVER_PORT=8080
```

### Infraestrutura externa: docker-compose

[`docker-compose.yaml`](docker-compose.yaml) sobe MySQL e RabbitMQ (mas não a aplicação Go em si — ela roda localmente via `go run`, apontando para `localhost`):

```yaml
services:
  mysql:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: orders
    ports:
      - 3306:3306

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - 5672:5672    # protocolo AMQP
      - 15672:15672  # painel de management (web)
```

### Ponto de atenção real: uma inconsistência de configuração

A conexão MySQL usa corretamente os valores vindos de `configs` (Viper, lidos do `.env`). Mas a conexão RabbitMQ está **hardcoded** direto no código:

```go
func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // string fixa, não vem de configs
	// ...
}
```

Isso é um exemplo real de inconsistência entre "como configuramos uma dependência externa" dentro do mesmo arquivo — vale notar como um ponto de atenção, não como algo a ser temido: acontece o tempo todo em projetos reais, e identificar esse tipo de coisa é parte do amadurecimento como leitor de código.

---

## 🗂️ Estrutura do Projeto

```
18-clean-architecture/
│
├── api/
│   └── create_order.http        # requisição REST de exemplo, pronta pra usar
│
├── cmd/
│   └── ordersystem/
│       ├── .env                 # variáveis de configuração
│       ├── main.go              # Composition Root — sobe os 3 servidores
│       ├── wire.go              # definição de DI para o Google Wire (não compilado normalmente)
│       └── wire_gen.go          # código de DI gerado pelo Wire (este SIM é compilado)
│
├── configs/
│   └── config.go                # leitura de configuração via Viper
│
├── internal/
│   ├── entity/                  # camada mais interna: regras de negócio puras
│   │   ├── interface.go         #   contrato do repositório (DIP)
│   │   ├── order.go             #   a entity Order
│   │   └── order_test.go        #   teste unitário puro
│   │
│   ├── usecase/
│   │   └── create_order.go      # regra de negócio da aplicação + DTOs
│   │
│   ├── event/                   # peças concretas de evento do domínio Order
│   │   ├── order_created.go
│   │   └── handler/
│   │       └── order_created_handler.go   # publica no RabbitMQ
│   │
│   └── infra/                   # Interface Adapters + Frameworks & Drivers
│       ├── database/
│       │   ├── order_repository.go        # implementa OrderRepositoryInterface
│       │   └── order_repository_test.go   # teste de integração com SQLite in-memory
│       │
│       ├── web/                 # controller REST
│       │   ├── order_handler.go
│       │   └── webserver/
│       │       └── webserver.go
│       │
│       ├── grpc/                # controller gRPC
│       │   ├── protofiles/order.proto
│       │   ├── pb/               # gerado por protoc — não editar
│       │   └── service/order_service.go
│       │
│       └── graph/                # controller GraphQL
│           ├── schema.graphqls   # contrato GraphQL
│           ├── resolver.go       # injeção de dependência do gqlgen
│           ├── schema.resolvers.go  # implementação manual do resolver
│           ├── model/            # gerado por gqlgen — não editar
│           └── generated.go      # gerado por gqlgen — não editar (3470 linhas)
│
├── pkg/
│   └── events/                   # Event Dispatcher genérico (independente de Order)
│       ├── interface.go
│       ├── event_dispatcher.go
│       └── event_dispatcher_test.go
│
├── docker-compose.yaml           # MySQL + RabbitMQ
├── gqlgen.yml                    # configuração do gerador GraphQL
├── go.mod
└── tools.go                      # garante que ferramentas fiquem no go.sum
```

### Separação de responsabilidades, em uma linha cada

```
entity/    → O QUÊ é uma regra de negócio válida (independe de tudo)
usecase/   → COMO esta aplicação usa essa regra (orquestra, sem saber de protocolo/infra)
infra/     → ONDE os dados moram e COM QUE PROTOCOLO o mundo externo entra
cmd/       → COMO tudo se conecta (Composition Root)
```

---

## 🔍 Walkthrough do Código

### Os Três Controllers Lado a Lado

Já vimos o código de cada um separadamente na seção de Conceitos Fundamentais. Aqui está o resumo visual do fluxo idêntico que os três seguem:

```
┌─────────────┐     ┌──────────────────┐     ┌───────────────────┐     ┌─────────────┐
│  Requisição  │────▶│  Controller       │────▶│  CreateOrderUseCase │────▶│  Resposta   │
│  (protocolo  │     │  traduz formato   │     │      .Execute()     │     │ (protocolo  │
│  específico) │     │  → OrderInputDTO  │     │                     │     │ específico) │
└─────────────┘     └──────────────────┘     └───────────────────┘     └─────────────┘

REST:    JSON HTTP           → WebOrderHandler.Create()
gRPC:    *pb.CreateOrderRequest → OrderService.CreateOrder()
GraphQL: *model.OrderInput   → mutationResolver.CreateOrder()

Todos os três chamam exatamente: usecase.CreateOrderUseCase.Execute(dto)
```

### Quem Importa Quem — O Fluxo Real de Dependências

Analisando os `import` de cada arquivo do projeto, esta é a direção real das setas de dependência:

```
                         ┌───────────────────────────────┐
                         │        cmd/ordersystem          │  (Composition Root)
                         │  main.go, wire.go, wire_gen.go  │
                         └────────────────┬────────────────┘
                                          │ importa tudo abaixo
        ┌──────────────────────────────────┼──────────────────────────────────┐
        │                                  │                                  │
┌───────▼────────┐              ┌──────────▼──────────┐              ┌────────▼────────┐
│  infra/web       │              │  infra/grpc/service  │              │  infra/graph      │
│  (chi + handler)  │              │  (pb + OrderService)  │              │  (gqlgen + resolver)│
└───────┬────────┘              └──────────┬──────────┘              └────────┬────────┘
        │  importa                          │  importa                          │  importa
        └────────────────┬───────────────────┴────────────────┬─────────────────┘
                          │                                    │
                   ┌──────▼──────┐                     ┌───────▼───────┐
                   │   usecase    │◄────────────────────┤ event/handler  │
                   │  (DTOs +      │   usecase importa    │  (RabbitMQ)    │
                   │   Execute)    │   events, não event   └───────┬───────┘
                   └──────┬──────┘                                │ importa
                          │  importa                                │
                 ┌────────┴────────┐                       ┌────────▼────────┐
                 │                 │                       │   pkg/events     │
          ┌──────▼──────┐   ┌──────▼───────┐               │ (interfaces +    │
          │   entity     │   │  pkg/events   │◄──────────────┤  dispatcher)     │
          │ (Order +     │   │ (interfaces + │               └──────────────────┘
          │  interface)  │   │  dispatcher)  │
          └──────▲──────┘   └──────────────┘
                 │  importa (implementa a interface declarada por entity)
          ┌──────┴──────┐
          │infra/database│
          │(OrderRepository)│
          └─────────────┘
```

O que confirmar visualmente neste diagrama, e que é a prova de que a Regra de Dependência está sendo respeitada:

- **`entity`** só importa `errors` (stdlib) — nenhuma dependência do projeto.
- **`pkg/events`** só importa `sync`, `time`, `errors` (stdlib) — infraestrutura de domínio genérica, reaproveitável fora deste contexto de pedidos.
- **`usecase`** importa `entity` e `pkg/events` — nunca `infra/*`.
- **`infra/database`** importa `entity` (para implementar a interface e usar o tipo `Order`) — a única seta de import "de fora para dentro" no sentido usual, mas semanticamente é a camada externa se adaptando ao contrato definido pela interna, exatamente a inversão de dependência funcionando.
- Nenhuma das três camadas mais internas (`entity`, `usecase`, `pkg/events`) importa qualquer coisa de `internal/infra/*`.

---

## ▶️ Como Executar

**Pré-requisitos:**
- Go 1.19+
- Docker (para MySQL e RabbitMQ)

**Passos:**

```bash
# 1. Entre na pasta do projeto
cd aulas/18-clean-architecture

# 2. Suba MySQL e RabbitMQ
docker-compose up -d

# 3. Crie a tabela "orders" no MySQL (não há migration automatizada)
docker exec -it mysql mysql -uroot -proot orders -e \
  "CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id));"

# 4. Instale as dependências
go mod download

# 5. Execute a aplicação
go run cmd/ordersystem/main.go

# Saída esperada:
# Starting web server on port :8000
# Starting gRPC server on port 50051
# Starting GraphQL server on port 8080
```

Ao final, você terá três servidores rodando simultaneamente no mesmo processo:

| Protocolo | Porta | Endpoint |
|---|---|---|
| REST | `8000` | `POST http://localhost:8000/order` |
| gRPC | `50051` | serviço `OrderService`, RPC `CreateOrder` |
| GraphQL | `8080` | Playground em `http://localhost:8080/`, queries em `/query` |
| RabbitMQ Management | `15672` | `http://localhost:15672` (usuário/senha: `guest`/`guest`) |

---

## 🎮 Testando as Três Portas de Entrada

### REST

O projeto já traz um exemplo pronto em [`api/create_order.http`](api/create_order.http) (compatível com a extensão "HTTP Client" de IDEs como VS Code):

```http
POST http://localhost:8000/order HTTP/1.1
Host: localhost:8000
Content-Type: application/json

{
    "id":"a",
    "price": 100.5,
    "tax": 0.5
}
```

Resposta esperada:
```json
{
  "id": "a",
  "price": 100.5,
  "tax": 0.5,
  "final_price": 101
}
```

Ou, via `curl`:
```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"a","price":100.5,"tax":0.5}'
```

### gRPC

O servidor gRPC deste projeto tem **reflection habilitada** (`reflection.Register(grpcServer)` em `main.go`), o que significa que ferramentas como [Evans](https://github.com/ktr0731/evans) ou `grpcurl` conseguem descobrir o serviço sem precisar do arquivo `.proto` em mãos:

```bash
grpcurl -plaintext -d '{"id": "b", "price": 100.5, "tax": 0.5}' \
  localhost:50051 pb.OrderService/CreateOrder
```

O contrato usado é o [`order.proto`](internal/infra/grpc/protofiles/order.proto):
```proto
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
}
```

### GraphQL

Abra `http://localhost:8080/` no navegador para acessar o Playground, e rode a mutation do schema real ([`schema.graphqls`](internal/infra/graph/schema.graphqls)):

```graphql
mutation {
  createOrder(input: {
    id: "c"
    Price: 100.5
    Tax: 0.5
  }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

Resposta:
```json
{
  "data": {
    "createOrder": {
      "id": "c",
      "Price": 100.5,
      "Tax": 0.5,
      "FinalPrice": 101
    }
  }
}
```

Note que, nos três casos, você está testando **a mesma regra de negócio** por três portas de entrada diferentes — o próprio `CreateOrderUseCase.Execute` sendo executado.

---

## 🎨 Padrões de Design

Este projeto usa, na prática, bem mais padrões de design clássicos do que os nomes "Clean Architecture" ou "Entity"/"Use Case" deixam óbvio à primeira vista. Alguns você já viu em detalhe nas seções anteriores — aqui eles ganham um resumo rápido com link de volta, para servir de índice de consulta. Outros ainda não tinham aparecido em código nenhum — são detalhados agora, pela primeira vez, com exemplo didático completo.

### Os padrões que você já viu em detalhe

Cada um destes já teve seção própria, com código real do projeto e um exemplo didático isolado — aqui é só o resumo, para consulta rápida.

**Repository Pattern** — isola "como os dados são persistidos" atrás de uma interface; o use case só chama `Save(order)`, sem saber se é MySQL, Postgres ou memória. Veja o exemplo completo em [Repository](#repository--onde-os-dados-realmente-moram).

**Dependency Inversion Principle (DIP)** — a camada de negócio declara o contrato (`entity.OrderRepositoryInterface`); a camada externa o implementa. Veja o exemplo completo em [Interface](#interface--o-contrato-que-inverte-a-dependência).

**Dependency Injection (DI)** — todo construtor `NewX(...)` do projeto recebe dependências prontas, em vez de criá-las internamente; a variante "DI via geração de código" (Google Wire) automatiza esse código de fiação em build time. Veja o exemplo completo em [Dependency Injection](#dependency-injection--manual-vs-google-wire).

**DTO (Data Transfer Object)** — `OrderInputDTO`/`OrderOutputDTO`, `model.Order`/`model.OrderInput`, `pb.CreateOrderRequest`/`pb.CreateOrderResponse`: cada protocolo tem seu próprio formato, convertido para o DTO comum antes de `Execute`. Veja a discussão completa em [Por que DTOs e não a entity direto?](#use-case--a-regra-de-negócio-da-aplicação), dentro da seção de Use Case.

**Observer / Event Dispatcher (Pub-Sub em memória)** — `pkg/events` + `internal/event` desacoplam o efeito colateral (publicar em fila) da lógica principal do use case. Veja o exemplo completo em [Eventos de Domínio](#eventos-de-domínio--o-event-dispatcher-dentro-da-arquitetura).

**Adapter** — `WebOrderHandler.Create`, `OrderService.CreateOrder`, `mutationResolver.CreateOrder`: cada um converte seu formato específico de entrada/saída para o contrato comum do use case. Veja o exemplo completo em [Controllers/Adapters](#controllersadapters--três-portas-para-a-mesma-casa).

**Composition Root** — `cmd/ordersystem/main.go` é o único lugar do sistema onde tudo é instanciado e amarrado. Veja a seção completa em [Composition Root](#️-composition-root--onde-tudo-se-conecta).

### Test Double / Mock — detalhado agora pela primeira vez

Você já viu esse padrão em ação, sem nenhum comentário sobre ele, em [`pkg/events/event_dispatcher_test.go`](pkg/events/event_dispatcher_test.go): a struct `MockHandler`, construída com `testify/mock`, verifica que `Dispatch` chama cada handler registrado, sem depender de lógica real de handler nenhuma. É a mesma técnica que já usamos no [teste do Use Case com Mocks](#use-case--a-regra-de-negócio-da-aplicação).

Veja a versão mais enxuta possível da mesma ideia, isolada, fora do domínio Order:

```go
package billing

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type PaymentGateway interface {
	Charge(amount float64) error
}

type PaymentGatewayMock struct{ mock.Mock }

func (m *PaymentGatewayMock) Charge(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func TestGivenAnAmount_WhenCharge_ThenShouldCallGatewayOnce(t *testing.T) {
	gatewayMock := &PaymentGatewayMock{}
	gatewayMock.On("Charge", 100.0).Return(nil)

	gatewayMock.Charge(100.0)

	gatewayMock.AssertExpectations(t) // prova que Charge(100.0) foi realmente chamado
}
```

O que importa aqui não é o valor devolvido — é a **prova de que a chamada aconteceu**, com o argumento certo. Essa é a essência de um Mock: ele verifica *comportamento*, não só *dado*.

**Não confunda com:** "Test Double" é o termo guarda-chuva para qualquer substituto de uma dependência real em teste — **Mock**, **Stub** e **Fake** são três categorias diferentes dentro desse guarda-chuva, e é comum iniciante usar os três nomes como sinônimos:
- **Mock** (visto acima): verifica se um método foi chamado, com quais argumentos, quantas vezes — a pergunta é "isso aconteceu do jeito certo?".
- **Stub**: só devolve um valor fixo, programado de antemão, sem nenhuma verificação de chamada — a pergunta é só "o que isso devolve?", nunca "isso foi chamado?".
- **Fake**: uma implementação funcional simplificada, com lógica de verdade (ainda que mais simples que a real) — é o próximo padrão desta seção.

### Fake/In-Memory Infra para Teste — detalhado agora pela primeira vez

Você já viu o exemplo real deste padrão duas vezes neste README, sem que ele fosse nomeado como tal: o teste de [`order_repository_test.go`](internal/infra/database/order_repository_test.go) usa SQLite `:memory:` no lugar do MySQL real (visto na seção de [Repository](#repository--onde-os-dados-realmente-moram)), e o próprio `InMemoryUserRepository` mostrado como ["Exemplo Clássico de Trocar o Banco por um Fake"](#repository--onde-os-dados-realmente-moram) — o nome do próprio bloco já entregava o padrão.

A diferença que separa um Fake de um Mock é justamente o que faz dele uma categoria própria: um Fake **tem lógica de verdade** por dentro — um mapa que realmente guarda e busca dados, uma estrutura que realmente se comporta como a coisa real, só que mais simples (sem rede, sem disco, sem SQL). Um Mock não tem lógica nenhuma: ele só registra "fui chamado com X" e devolve o que foi programado para devolver.

**Trade-off:** um Fake dá mais confiança que um Mock, porque testa um comportamento de verdade (inserir e depois buscar precisa realmente funcionar dentro do Fake) — não só "o método foi chamado". Em compensação, um Fake dá mais trabalho para escrever e manter do que um Mock (que costuma ser gerado quase automaticamente por uma biblioteca como `testify/mock`). Este projeto usa os dois: Mock para o Event Dispatcher (o que importa é só "o handler foi notificado?"), Fake/SQLite para o Repository (o que importa é "os dados realmente persistem e voltam corretos?").

### Os Quatro Padrões que Este Projeto Não Precisou

Estes quatro não aparecem em nenhum lugar do código real — não porque a Clean Architecture os proíba, mas porque o escopo deste projeto (criar um único tipo de entidade, sem consulta customizada, sem múltiplas escritas coordenadas) nunca criou a necessidade deles. Veja o que cada um resolveria, e por que não fez falta aqui:

**Factory Method** — na prática, `entity.NewOrder(id, price, tax)` (já visto na seção de [Entity](#entity--a-regra-de-negócio-mais-pura)) *é* um Factory Method simples: uma função dedicada a criar o objeto, validando antes de devolver. O que falta para ser um Factory Method "clássico" completo é a **decisão de qual tipo concreto instanciar**, baseada em algum parâmetro:

```go
func NewNotifier(tipo string) Notifier {
	if tipo == "email" {
		return &EmailNotifier{}
	}
	return &SMSNotifier{}
}
```

Como este projeto só tem um tipo de entidade (`Order`), nunca precisou decidir "qual struct concreta criar" — só "criar `Order`, validado".

**CQRS (Command Query Responsibility Segregation)** — separaria o caminho de **escrita** (`CreateOrderUseCase`, que já existe) do caminho de **leitura**, com estruturas dedicadas a consultar dados, sem passar pelas mesmas regras de validação da escrita:

```go
type ListOrdersQuery struct{ Repository OrderReadRepository }

func (q *ListOrdersQuery) Execute() ([]OrderOutputDTO, error) {
	return q.Repository.FindAll() // caminho de leitura, sem nenhuma regra de escrita no meio
}
```

Este projeto só implementa o lado de escrita (criar um pedido) — nunca teve uma segunda operação de leitura customizada que justificasse separar os dois caminhos. Não é uma rejeição consciente de CQRS, é o escopo nunca ter pedido.

**Unit of Work** — coordenaria múltiplas operações de escrita como uma transação só, garantindo que todas aconteçam ou nenhuma aconteça:

```go
tx, _ := db.Begin()
orderRepo.SaveWithTx(tx, order)
inventoryRepo.DecrementWithTx(tx, item) // uma segunda escrita, na MESMA transação
tx.Commit() // as duas só "valem" juntas
```

Este projeto salva uma única entidade por operação (`Save(order)`) — nunca precisou coordenar duas escritas diferentes dentro da mesma transação.

**Saga / Outbox** — já foi discutido em detalhe, com o trade-off real deste projeto, na seção de [Eventos de Domínio](#trade-off-real-do-projeto-sem-garantia-transacional): existe mensageria (RabbitMQ), mas de forma *fire-and-forget*, sem garantia transacional entre salvar o pedido e publicar o evento. O Outbox Pattern é justamente a técnica que resolveria isso — vale reler aquela seção para os detalhes completos.

Nenhum desses quatro é "mais avançado" ou "melhor" que os padrões já usados no projeto — eles resolvem problemas que este projeto específico, do tamanho que tem, nunca teve.

---

## 🧬 Clean Architecture vs. Outras Abordagens

Até aqui você viu como a Clean Architecture organiza este projeto específico. Mas Clean Architecture não é a única forma de pensar arquitetura — e entender como ela se relaciona com outras abordagens ajuda a enxergar o que é realmente específico dela, e o que é uma ideia mais ampla que aparece sob nomes diferentes.

### Clean Architecture vs. DDD (Domain-Driven Design)

**DDD (Domain-Driven Design, ou Design Orientado ao Domínio)** e Clean Architecture **não são concorrentes** — respondem perguntas diferentes:

- **DDD** é focado em *como modelar o domínio*: como representar as regras de negócio complexas em código, com vocabulário rico (Bounded Context, Aggregate, Value Object, Domain Event, Ubiquitous Language).
- **Clean Architecture** é focado em *como organizar o código em camadas*: onde cada peça mora, e em que direção as dependências apontam.

Na prática, é muito comum ver os dois juntos: DDD desenha **o que** vai dentro dos círculos de `Entities`/`Use Cases`; Clean Architecture desenha **as fronteiras** entre esse núcleo e o resto do sistema. Este projeto, aliás, já usa vocabulário de DDD sem nomear explicitamente — `internal/event/order_created.go` é literalmente um **Domain Event**, termo que vem do DDD.

| Conceito DDD | Equivalente neste projeto | Observação |
|---|---|---|
| **Entity** (tem identidade própria) | `entity.Order` (tem `ID`) | Bate diretamente |
| **Value Object** (sem identidade, imutável) | Não existe um exemplo dedicado — `Price`/`Tax` são só `float64` soltos | Veja "e se fosse diferente?" logo abaixo |
| **Aggregate / Aggregate Root** | Não modelado — o projeto só tem uma entity isolada, sem composição | O domínio é simples demais pra precisar disso |
| **Domain Event** | `internal/event/order_created.go` (`OrderCreated`) | Bate diretamente — é literalmente o vocabulário do DDD |
| **Bounded Context** | Não modelado — o projeto inteiro é um único módulo | Faria sentido num sistema com múltiplos domínios (ex.: Pedidos + Pagamentos + Estoque, cada um com seu próprio `Order`) |
| **Ubiquitous Language** (vocabulário compartilhado entre código e negócio) | Os nomes `Order`, `Price`, `Tax`, `FinalPrice` já seguem essa prática informalmente | Nenhum termo técnico inventado escondendo o vocabulário do negócio |

**E se fosse diferente?** Um `Price` modelado como Value Object, em vez de `float64` solto, poderia ficar assim:

```go
package entity

import "errors"

// Money é um Value Object: não tem identidade própria (duas instâncias
// com o mesmo valor SÃO iguais), e é imutável (não existe um "SetValue").
type Money struct {
	cents int64
}

func NewMoney(reais float64) (Money, error) {
	if reais < 0 {
		return Money{}, errors.New("valor não pode ser negativo")
	}
	return Money{cents: int64(reais * 100)}, nil
}

func (m Money) Add(other Money) Money {
	return Money{cents: m.cents + other.cents} // nunca modifica m ou other
}
```

Isso evitaria, por exemplo, erros de arredondamento com `float64` em cálculos financeiros — um problema real que `Money` como Value Object resolve de forma mais robusta do que dois `float64` soltos.

**Trade-off:** usar só Clean Architecture (organização em camadas) sem DDD completo (modelagem rica) é perfeitamente válido — foi exatamente o que este projeto fez, com uma entity relativamente simples. DDD completo (Bounded Contexts, Aggregates, Value Objects por toda parte) só compensa o esforço extra em domínios genuinamente complexos, com regras de negócio elaboradas e múltiplas equipes — o mesmo tipo de raciocínio "compensa o esforço?" que você já viu na seção [Trade-offs](#️-trade-offs-convencional-vs-fora-do-convencional), agora aplicado à escolha de quanto DDD vale a pena adotar.

### Clean Architecture vs. Arquitetura Hexagonal (Ports & Adapters)

Esta é a comparação mais próxima de todas: Clean Architecture e **Arquitetura Hexagonal** (também chamada de **Ports & Adapters**, criada por Alistair Cockburn) são, na prática, **quase a mesma ideia**, com vocabulário diferente.

```
HEXAGONAL (Ports & Adapters)              CLEAN ARCHITECTURE (já vista neste README)

        ╱‾‾‾‾‾‾‾‾╲                        ┌───────────────────────────┐
       ╱            ╲                      │   Frameworks & Drivers     │
      │   NÚCLEO      │  ← "dentro"         │  ┌───────────────────┐   │
      │  (regras de   │                     │  │ Interface Adapters │   │
      │   negócio)    │                     │  │  ┌─────────────┐  │   │
       ╲   PORT       ╱  ← contrato          │  │  │  Use Cases  │  │   │
        ╲____________╱                       │  │  │ ┌─────────┐ │  │   │
             │                               │  │  │ │Entities │ │  │   │
         ADAPTER  ← "fora"                   │  │  │ └─────────┘ │  │   │
                                              │  │  └─────────────┘  │   │
      Só 2 zonas: dentro / fora              │  └───────────────────┘   │
                                              └───────────────────────────┘
                                              4 camadas concêntricas nomeadas
```

| Conceito | Nome em Clean Architecture | Nome em Hexagonal |
|---|---|---|
| O contrato que a camada de negócio declara | Interface (ex.: `entity.OrderRepositoryInterface`) | **Port** |
| A implementação concreta que satisfaz o contrato | Implementação (ex.: `infra/database.OrderRepository`) | **Adapter** (mais especificamente, "driven adapter"/"adapter secundário") |
| O tradutor de entrada (ex.: `WebOrderHandler`) | Controller / Interface Adapter | **Adapter** (mais especificamente, "driving adapter"/"adapter primário") |
| Tudo que não é regra de negócio | Frameworks & Drivers | "Fora do hexágono" |
| A regra central protegida | Entities + Use Cases | "Núcleo"/"Domínio" |

A diferença mais citada entre as duas: Hexagonal desenha só **duas zonas** (dentro do hexágono / fora dele) — mais simples de explicar, menos granular. Clean Architecture desenha explicitamente **quatro camadas concêntricas** nomeadas (Entities, Use Cases, Interface Adapters, Frameworks & Drivers) — mais granular, dá nomes mais específicos pra cada responsabilidade. Na prática, o projeto deste README poderia ser descrito com qualquer um dos dois vocabulários sem mudar uma linha de código — `entity.OrderRepositoryInterface` é, ao mesmo tempo, uma interface no sentido Clean Architecture e um Port no sentido Hexagonal.

### Clean Architecture vs. MVC/Layered (Arquitetura em Camadas Tradicional)

Este é o contraponto mais familiar para quem já mexeu com algum framework web (Spring, Rails, Laravel, ASP.NET, Django) — o clássico `Controller → Service → Repository → Model`, muitas vezes já vindo pronto do "scaffold" do framework.

A pegadinha: MVC/Layered **também tem camadas**! A diferença crucial não é "ter camadas ou não" — é a **direção da dependência**:

```
LAYERED TRADICIONAL (pilha simples)          CLEAN ARCHITECTURE (círculos, já visto)

┌─────────────────┐                          ┌───────────────────────────┐
│    Controller     │                          │  Frameworks & Drivers      │
└─────────┬─────────┘                          │  ┌───────────────────┐   │
          │ importa                            │  │ Interface Adapters │   │
┌─────────▼─────────┐                          │  │  ┌─────────────┐  │   │
│      Service        │                          │  │  │  Use Cases  │  │   │
└─────────┬─────────┘                          │  │  │ ┌─────────┐ │  │   │
          │ importa                            │  │  │ │Entities │ │  │   │
┌─────────▼─────────┐                          │  │  │ └────┬────┘ │  │   │
│    Repository        │                          │  │  └──────┼──────┘  │   │
└─────────┬─────────┘                          │  └──────────┼──────────┘   │
          │ importa                            └─────────────┼──────────────┘
┌─────────▼─────────┐                                          │
│  Model / ORM        │                          Nada aqui dentro sabe que
└─────────────────┘                          "Model/ORM" existe — a seta
                                              de import NUNCA vai daqui
Seta de import desce direto,                 pra fora sem passar por
sem interface no meio                        uma interface primeiro
```

No Layered tradicional, é comum a camada de baixo (o Model/ORM) ser importada **livremente** por todas as camadas de cima, sem nenhuma interface no meio — o `Service` frequentemente conhece o tipo concreto do ORM diretamente. Na Regra de Dependência da Clean Architecture (já vista na seção [O que é Clean Architecture?](#a-regra-de-dependência)), isso é proibido: a camada de negócio nunca importa a de infra diretamente, sempre através de uma interface que a própria camada de negócio declara.

**Trade-off honesto:** Layered tradicional é mais rápido de aprender e escrever para um CRUD simples — é literalmente o que a maioria dos frameworks web gera por padrão ao rodar um "scaffold"/"generate". Clean Architecture exige mais disciplina (mais arquivos, mais interfaces, mais indireção), mas paga esse investimento de volta quando o sistema cresce, precisa trocar peças, ou precisa ser testado de forma confiável sem subir infraestrutura — o mesmo raciocínio já visto na seção [Casos de Uso Ideais](#-casos-de-uso-ideais), agora aplicado à escolha entre arquiteturas inteiras, não só a decisões dentro da Clean Architecture.

### Resumo: as Quatro Abordagens Lado a Lado

| Abordagem | Foco principal | Quando escolher |
|---|---|---|
| **Clean Architecture** | Organizar o código em camadas, com a Regra de Dependência protegendo o núcleo | Sistemas que precisam durar, trocar peças (banco, protocolo), ou ser testados com confiança |
| **Hexagonal (Ports & Adapters)** | O mesmo objetivo da Clean Architecture, com vocabulário mais simples (dentro/fora) | Times que preferem um modelo mental mais enxuto, sem precisar nomear 4 camadas |
| **DDD (Domain-Driven Design)** | Modelar bem um domínio de negócio complexo, com vocabulário rico e compartilhado | Domínios genuinamente complexos, com regras de negócio elaboradas e múltiplas equipes |
| **MVC/Layered tradicional** | Produtividade imediata num CRUD simples, aproveitando o que o framework já entrega pronto | Protótipos, CRUDs simples, times pequenos, sistemas sem expectativa de crescer muito |

---

## 🧩 Adotando Só Partes da Clean Architecture

Tudo que você viu até aqui descreve a Clean Architecture "completa" — as quatro camadas, interfaces em todo lugar, DTOs explícitos. Mas isso **não é tudo ou nada**. Adotar só pedaços dela, em projetos menores ou misturado com outras abordagens, é uma prática real e reconhecida no mercado — não é uma "versão errada" da arquitetura.

O próprio Robert C. Martin, criador do termo, sempre defendeu que **os princípios importam mais que a estrutura exata de pastas**: o que protege seu sistema é a Regra de Dependência (a camada de negócio nunca depender de detalhes de infraestrutura), não necessariamente ter as quatro camadas nomeadas exatamente como neste projeto.

### Cenário 1: Só o Repository Pattern

Um projeto pequeno (um CLI, um script de automação, uma API bem simples) pode não ter `usecase`/`entity` separados como camadas formais, mas ainda assim isolar a persistência atrás de uma interface — só para poder trocar por um fake em teste, sem o resto do aparato:

```go
package main

type TaskRepository interface {
	Save(task string) error
}

type fileRepository struct{ path string }

func (r *fileRepository) Save(task string) error {
	// grava num arquivo de verdade
	return nil
}

// O "programa principal" só depende da interface — não do arquivo.
func processarTarefa(repo TaskRepository, task string) error {
	return repo.Save(task)
}
```

Sem `entity`, sem `usecase`, sem DTO — só o Repository Pattern isolado, já suficiente para testar `processarTarefa` com um repositório fake, sem tocar no disco.

### Cenário 2: Só a Entity com Validação, sem as Quatro Camadas Completas

Um handler HTTP simples pode chamar diretamente uma entity validada, sem um use case dedicado nem DTOs formais — ganhando *alguma* proteção de negócio, mesmo sem o pacote completo:

```go
func handlerCriarPedido(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID    string
		Price float64
		Tax   float64
	}
	json.NewDecoder(r.Body).Decode(&input)

	// Chama a entity diretamente — sem usecase.CreateOrderUseCase no meio.
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(order)
}
```

Não tem a flexibilidade de múltiplos protocolos nem o desacoplamento total deste projeto real — mas já herda a proteção de `NewOrder` recusar dados inválidos, sem precisar montar as quatro camadas completas para um endpoint simples.

### Cenário 3: Clean Architecture "Lite" Dentro de um Monólito Layered Maior

Um cenário de mercado real e comum: um sistema inteiro segue um Layered/MVC tradicional mais simples — mas uma parte **específica e crítica** (por exemplo, o motor de cálculo de preços e impostos) ganha as quatro camadas completas, isolada e testável, enquanto o resto do sistema (telas administrativas, relatórios simples) continua no modelo mais direto:

```
┌─────────────────────────────────────────────────────────┐
│                    SISTEMA MAIOR (Layered/MVC)             │
│                                                              │
│  Telas administrativas  ──▶  Service  ──▶  Model/ORM       │
│  Relatórios simples      ──▶  Service  ──▶  Model/ORM       │
│                                                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │   Motor de Preços/Impostos (Clean Architecture)       │   │
│  │   Entity → Use Case → Interface Adapters isolados      │   │
│  │   (testado sem banco, protegido, trocável)             │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

A lógica que muda com frequência, tem muitas regras, ou precisa de confiança alta ganha o investimento completo; o resto do sistema, que raramente muda e não justifica o esforço, continua simples.

### O que Você Ganha e o que Abre Mão

| Adoção | O que você ganha | O que você abre mão (vs. Clean Architecture completa) |
|---|---|---|
| Só Repository Pattern | Testabilidade da persistência, sem reescrever tudo | Regra de negócio ainda pode vazar para fora da entity/use case |
| Só Entity com validação | Proteção contra dados inválidos, com pouquíssimo código extra | Nenhum desacoplamento de protocolo — trocar de REST pra gRPC ainda exige reescrever o handler inteiro |
| Clean Architecture "lite" dentro de um Layered maior | Investimento concentrado só onde compensa | Inconsistência de estilo dentro do mesmo repositório — exige documentação clara de onde é onde, para não confundir o time |

Adotar só uma parte é sempre um trade-off consciente — nunca uma versão "grátis" da arquitetura completa sem perdas.

---

## ⚖️ Trade-offs: Convencional vs. Fora do Convencional

### Clean Architecture: quando o esforço extra compensa

✅ **Vantagens**

✅ Regra de negócio isolada e fácil de testar sem infraestrutura

✅ Trocar banco de dados, framework web, ou biblioteca externa não exige reescrever regra de negócio

✅ Múltiplos protocolos de entrega (REST/gRPC/GraphQL) coexistem sem duplicar lógica

✅ Times grandes conseguem trabalhar em camadas diferentes com menos conflito

✅ Onboarding de novos desenvolvedores fica mais previsível (cada coisa tem seu lugar)

❌ **Desvantagens**

❌ Mais arquivos, mais indireção — um CRUD simples vira vários arquivos pequenos

❌ Curva de aprendizado real para iniciantes (é este README inteiro, basicamente)

❌ Overhead desnecessário para protótipos descartáveis ou scripts pequenos

❌ Fácil "over-engenheirar" — aplicar todo o rigor a um sistema que nunca vai crescer

### Tabela: quando usar Clean Architecture rigorosa vs. uma abordagem mais simples

| Cenário | Clean Architecture compensa? | Por quê |
|---|---|---|
| Startup validando uma ideia (MVP) | Provavelmente não | Velocidade de iteração importa mais que isolamento — o código pode nem sobreviver ao pivot |
| Sistema com múltiplos protocolos de entrada (como este projeto) | Sim | O ganho de não duplicar lógica entre REST/gRPC/GraphQL já paga o investimento |
| Sistema que vai durar anos, com múltiplos times | Sim | O custo de mudança composto ao longo do tempo é o que a arquitetura mais reduz |
| Script utilitário de uso único | Não | O overhead de camadas não tem tempo de se pagar |
| Domínio de negócio complexo, com muitas regras que mudam com frequência | Sim | Isolar essas regras é exatamente o que a arquitetura protege melhor |

### Múltiplos protocolos no mesmo processo vs. serviços separados

Este projeto sobe REST, gRPC e GraphQL **no mesmo binário/processo**. A alternativa "fora do convencional" seria cada protocolo ser um serviço/processo separado (ex.: um microsserviço REST, outro gRPC), compartilhando o mesmo `usecase` como uma biblioteca interna comum.

| Aspecto | Mesmo processo (usado aqui) | Serviços separados |
|---|---|---|
| Simplicidade de deploy | Um único binário para subir | Múltiplos deploys, mais infraestrutura |
| Isolamento de falhas | Um crash derruba os três protocolos juntos | Um serviço cai sem afetar os outros |
| Escalar cada protocolo de forma independente | Não — escalam juntos | Sim — escala só o que precisa |
| Adequado para | Times pequenos, sistemas didáticos, MVPs internos | Sistemas com tráfego desigual entre protocolos, times grandes |

---

## 🎯 Casos de Uso Ideais

### Quando este tipo de arquitetura faz muito sentido

**1. Backend que serve múltiplos tipos de cliente**
```
App mobile    → prefere gRPC (rápido, tipado)
Frontend web  → prefere GraphQL (flexível)
Parceiro B2B  → prefere REST (universal, fácil de integrar)

Com Clean Architecture: 1 usecase, 3 adapters — sem duplicar regra.
```

**2. Domínio com regras de negócio que mudam com frequência**
- Sistemas financeiros, de precificação, de regras fiscais — onde isolar e testar a regra isoladamente (sem precisar de banco) acelera muito o ciclo de desenvolvimento.

**3. Sistemas de longa duração, mantidos por times que mudam ao longo dos anos**
- Um novo desenvolvedor consegue entender "onde mexer" olhando a estrutura de pastas, mesmo sem conhecer todo o histórico do projeto.

### Quando é overkill

**1. Protótipo para validar uma ideia rapidamente**
- O tempo gasto criando interfaces e DTOs pode ser maior que o tempo de vida útil do próprio protótipo.

**2. Script de automação de uso único**
- Ninguém vai trocar o "banco de dados" de um script que roda uma vez e é descartado.

**3. Equipe muito pequena, sistema muito simples, sem planos de crescer**
- O overhead de camadas pode atrapalhar mais do que ajudar quando não há complexidade real para gerenciar.

---

## 💼 Clean Architecture no Mercado de Trabalho

Nenhuma vaga de emprego tem "Clean Architecture" no título. Mas se você reparar nas entrelinhas de praticamente qualquer vaga de Go pleno ou sênior, os conceitos que este README já ensinou aparecem descritos com outras palavras — e as perguntas de entrevista técnica mais clássicas para essas vagas são, no fundo, pedidos para explicar o que você acabou de aprender aqui.

| O que a vaga/entrevista pede | Onde isso já foi ensinado neste README |
|---|---|
| "Desenho de software testável, baixo acoplamento" (comum em quase toda job description de backend) | [Interface — O Contrato que Inverte a Dependência](#interface--o-contrato-que-inverte-a-dependência) |
| "Experiência com arquitetura hexagonal / ports and adapters" (comum em vagas de fintech, bancos, pagamentos) | [Clean Architecture vs. Arquitetura Hexagonal](#clean-architecture-vs-arquitetura-hexagonal-ports--adapters) |
| "Domain-Driven Design" (comum em domínios complexos: pagamentos, logística, saúde, seguros) | [Clean Architecture vs. DDD](#clean-architecture-vs-ddd-domain-driven-design) |
| "Como você testaria uma regra de negócio sem subir um banco de dados?" (pergunta clássica de entrevista) | [Entity](#entity--a-regra-de-negócio-mais-pura) e [Testando o Use Case com Mocks](#use-case--a-regra-de-negócio-da-aplicação) |
| "Como você trocaria de banco de dados sem reescrever a aplicação inteira?" (pergunta clássica de entrevista) | [Repository — Onde os Dados Realmente Moram](#repository--onde-os-dados-realmente-moram) |
| "Explique Dependency Injection" / "qual a diferença entre DI e DIP?" (confusão que entrevistadores gostam de testar) | [Dependency Injection — Manual vs. Google Wire](#dependency-injection--manual-vs-google-wire) |

Uma ressalva importante para gerenciar expectativa: entender esses conceitos é **necessário, mas não suficiente**. O que costuma ser avaliado de verdade numa entrevista técnica ou num code review sênior não é "você sabe definir Repository Pattern", e sim "você sabe **quando** vale a pena usar isso, e quando seria over-engineering" — exatamente o raciocínio das seções [Trade-offs](#️-trade-offs-convencional-vs-fora-do-convencional) e [Casos de Uso Ideais](#-casos-de-uso-ideais), já vistas antes desta. Recitar a definição de Clean Architecture de cor impressiona pouco; justificar por que ela compensa (ou não) num cenário específico é o que realmente diferencia um candidato.

---

## 🧭 Inconsistências Reais do Projeto — Material de Estudo

Um dos jeitos mais eficazes de aprender arquitetura é treinar o olho para notar esse tipo de detalhe em código real — não é uma crítica ao curso, é prática de leitura crítica, algo que todo desenvolvedor experiente faz o tempo todo.

1. **`WebOrderHandler.Create` recria o use case a cada requisição**, em vez de recebê-lo pronto via injeção de dependência — diferente do padrão usado no `OrderService` (gRPC) e no `Resolver` (GraphQL), que recebem o use case já pronto no construtor. *(veja [Controllers/Adapters](#controllersadapters--três-portas-para-a-mesma-casa))*

2. **`internal/infra/web/webserver/starter.go` define `WebServerStarter`, mas ele não é usado em nenhum outro lugar do projeto** — nem em `main.go`, nem em `wire.go`. Parece um resquício de refatoração — código morto no estado atual.

3. **`entity.OrderRepositoryInterface.GetTotal()` está comentado na interface, mas implementado de verdade em `OrderRepository`** — um método "órfão", inacessível para quem depende só da interface. *(veja [Interface](#interface--o-contrato-que-inverte-a-dependência))*

4. **A string de conexão do RabbitMQ está hardcoded em `main.go`**, enquanto a conexão MySQL usa corretamente os valores vindos de `configs` (Viper, lidos do `.env`). *(veja [Composition Root](#️-composition-root--onde-tudo-se-conecta))*

5. **Não há garantia transacional entre salvar o pedido no banco e publicar o evento no RabbitMQ** — é possível salvar com sucesso e falhar a publicação, sem rollback nem nova tentativa automática. *(veja [Eventos de Domínio](#eventos-de-domínio--o-event-dispatcher-dentro-da-arquitetura))*

6. **Não há migration SQL versionada** — o schema da tabela `orders` no MySQL real precisa ser criado manualmente; só o teste com SQLite tem seu `CREATE TABLE` embutido no código. *(veja [Repository](#repository--onde-os-dados-realmente-moram))*

7. **Não há teste automatizado isolado do `CreateOrderUseCase`** — diferente da entity e do event dispatcher, que têm suites de teste completas. *(veja [Use Case](#use-case--a-regra-de-negócio-da-aplicação))*

---

## 🛠️ Comandos Úteis

### Rodando a aplicação

```bash
# Subir infraestrutura (MySQL + RabbitMQ)
docker-compose up -d

# Ver logs da infraestrutura
docker-compose logs -f

# Parar a infraestrutura
docker-compose down

# Rodar a aplicação
go run cmd/ordersystem/main.go
```

### Testes

```bash
# Rodar todos os testes do projeto
go test ./...

# Com saída detalhada
go test -v ./...

# Só a camada de entity (rápido, sem infra)
go test ./internal/entity

# Só a camada de repository (sobe SQLite in-memory)
go test ./internal/infra/database

# Só o event dispatcher
go test ./pkg/events

# Com cobertura
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### RabbitMQ

```bash
# Acessar o painel de management
open http://localhost:15672
# usuário: guest / senha: guest
```

### Regerando código (opcional — só se você editar os contratos)

Estes comandos exigem as respectivas ferramentas instaladas (`gqlgen`, `wire`, `protoc`) e **não são necessários apenas para rodar o projeto** — só se você alterar `schema.graphqls`, `wire.go` ou `order.proto`:

```bash
# Regenerar código GraphQL após editar schema.graphqls
go run github.com/99designs/gqlgen generate

# Regenerar código de DI após editar wire.go
go run github.com/google/wire/cmd/wire

# Regenerar código gRPC após editar order.proto
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto
```

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **Clean Architecture** | Estilo de organização de código em camadas concêntricas, com a regra de negócio isolada no centro |
| **Entity** | Regra de negócio mais geral e estável, independente de qualquer detalhe técnico |
| **Use Case** | Regra de negócio específica de uma aplicação — o que o sistema faz com as entities |
| **DTO (Data Transfer Object)** | Struct usada só para transportar dados entre camadas/protocolos, sem lógica de negócio |
| **Repository Pattern** | Isola "como os dados são persistidos" atrás de uma interface |
| **DIP (Dependency Inversion Principle)** | Camadas internas declaram interfaces; camadas externas as implementam — a dependência conceitual aponta para dentro |
| **DI (Dependency Injection)** | Um componente recebe suas dependências prontas de fora, em vez de criá-las internamente |
| **Composition Root** | O único lugar do sistema onde todas as dependências são instanciadas e conectadas |
| **Google Wire** | Ferramenta que gera, em tempo de compilação, o código de injeção de dependência |
| **Event Dispatcher** | Componente que distribui um evento para todos os handlers interessados nele |
| **Observer / Pub-Sub** | Padrão onde "assinantes" são notificados automaticamente quando algo acontece, sem acoplamento direto |
| **Goroutine** | Thread leve gerenciada pelo runtime do Go |
| **sync.WaitGroup** | Contador que permite esperar até que várias goroutines terminem |
| **Adapter** | Componente que traduz entre um formato externo (protocolo) e o contrato interno da aplicação |
| **gRPC** | Framework de comunicação de alta performance baseado em Protocol Buffers e HTTP/2 |
| **Protocol Buffers (protobuf)** | Formato binário de serialização usado pelo gRPC, definido em arquivos `.proto` |
| **gqlgen** | Biblioteca Go que gera código de servidor GraphQL a partir de um schema |
| **Reflection (gRPC)** | Recurso que permite a um cliente descobrir os serviços de um servidor gRPC sem ter o `.proto` em mãos |
| **Acoplamento** | O quanto uma parte do código depende de detalhes de outra parte para funcionar |
| **Interface** | Contrato de métodos que um tipo deve implementar; em Go, a implementação é implícita |
| **RabbitMQ / AMQP** | Message broker (intermediário de mensagens) e o protocolo (Advanced Message Queuing Protocol) que ele usa |
| **Outbox Pattern** | Técnica para garantir consistência entre salvar um dado e publicar um evento sobre ele, usando a mesma transação de banco |

---

## 🚀 Próximos Passos

Agora que você entende Clean Architecture com um caso real de múltiplos protocolos, explore estes tópicos:

**Imediato:**
- [ ] Escrever `create_order_test.go`, testando `CreateOrderUseCase.Execute` com mocks de `OrderRepositoryInterface` e `EventDispatcherInterface`
- [ ] Descomentar `GetTotal()` na interface `OrderRepositoryInterface` e usá-lo em algum lugar (ex.: um novo endpoint `GET /order/total`)
- [ ] Corrigir `WebOrderHandler` para receber o use case já pronto via construtor, como fazem `OrderService` e `Resolver`
- [ ] Mover a string de conexão do RabbitMQ para `configs`, junto com os demais valores do `.env`
- [ ] Adicionar uma migration SQL versionada para a tabela `orders` (ex.: com `golang-migrate`)

**Intermediário:**
- [ ] Implementar um consumer real que leia da fila publicada pelo `OrderCreatedHandler` (hoje só existe o lado produtor)
- [ ] Adicionar um segundo use case (ex.: `ListOrders`) e observar como ele se encaixa nas mesmas camadas
- [ ] Escrever testes de integração ponta a ponta, subindo o servidor real e fazendo requisições HTTP contra `/order`
- [ ] Explorar `context.Context` propagado do controller até o repository, para cancelamento e timeouts

**Avançado:**
- [ ] Implementar o **Outbox Pattern** para garantir consistência entre `Save` e `Dispatch` (ver [Eventos de Domínio](#eventos-de-domínio--o-event-dispatcher-dentro-da-arquitetura))
- [ ] Extrair cada protocolo (REST/gRPC/GraphQL) para um serviço separado, compartilhando `usecase` como módulo interno
- [ ] Adicionar observabilidade (logs estruturados, métricas, tracing) nas fronteiras entre camadas

**Comparativo com a Aula 17 (DI):**

A Aula 17 é inteiramente dedicada a Dependency Injection — o que este projeto usa em `wire.go`/`wire_gen.go`. Se a seção [Dependency Injection — Manual vs. Google Wire](#dependency-injection--manual-vs-google-wire) te deixou com dúvidas sobre os fundamentos de DI em si (por que injetar, como escrever isso manualmente antes de usar uma ferramenta), vale revisar a Aula 17 — ela constrói esse raciocínio do zero, sem a complexidade adicional dos múltiplos protocolos que este projeto acrescenta.

**Comparativo com a Aula 9 (Eventos):**

Da mesma forma, se a seção [Eventos de Domínio](#eventos-de-domínio--o-event-dispatcher-dentro-da-arquitetura) passou rápido demais pelas goroutines e pelo `sync.WaitGroup`, a Aula 9 constrói o Event Dispatcher inteiro peça por peça, com testes a cada passo — é o complemento natural para entender a fundo o mecanismo que a Aula 18 já usa pronto.

```
Aula 9 (Eventos)                 Aula 17 (DI)                    Aula 18 (Clean Architecture)
─────────────────────            ─────────────────────           ─────────────────────────────
Constrói o Event Dispatcher      Constrói DI manual e com Wire   Usa os dois prontos, dentro
do zero, com testes               isoladamente                    de um sistema completo com
                                                                    3 protocolos de entrega
```

Uma arquitetura madura como a deste projeto é, no fundo, a combinação de várias peças menores — cada uma delas dominável isoladamente antes de serem combinadas.
