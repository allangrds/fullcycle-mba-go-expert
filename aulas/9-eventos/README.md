# Eventos e Arquitetura Orientada a Eventos em Go

Este módulo ensina como construir sistemas orientados a eventos em Go, desde conceitos básicos até sistemas distribuídos com RabbitMQ. Você aprenderá a criar um Event Dispatcher completo, utilizar concorrência com goroutines e integrar message brokers para comunicação entre microserviços.

## 📑 Sumário

- [O que é Arquitetura Orientada a Eventos?](#o-que-é-arquitetura-orientada-a-eventos)
- [Por que usar Eventos?](#por-que-usar-eventos)
- [Conceitos Fundamentais em Go](#conceitos-fundamentais-em-go)
  - [Interfaces](#interfaces)
  - [Maps como Registry](#maps-como-registry)
  - [Slices e Manipulação](#slices-e-manipulação)
  - [Goroutines](#goroutines)
  - [Channels](#channels)
  - [sync.WaitGroup](#syncwaitgroup)
  - [Defer Statement](#defer-statement)
- [Progressão do Aprendizado](#progressão-do-aprendizado)
  - [1. Criando Interfaces da Solução](#1-criando-interfaces-da-solução)
  - [2-5. Criando Suite de Testes](#2-5-criando-suite-de-testes)
  - [6. Implementando o Método Dispatch](#6-implementando-o-método-dispatch)
  - [7. Revisitando Slices](#7-revisitando-slices)
  - [8. Removendo Handlers](#8-removendo-handlers)
  - [9. Adicionando Goroutines ao Dispatcher](#9-adicionando-goroutines-ao-dispatcher)
  - [10. Utilizando WaitGroup para Sincronização](#10-utilizando-waitgroup-para-sincronização)
  - [11. Integração com RabbitMQ](#11-integração-com-rabbitmq)
- [Padrões de Design](#padrões-de-design)
- [Aplicações Práticas](#aplicações-práticas)
- [Comandos Úteis](#comandos-úteis)

---

## 🎯 O que é Arquitetura Orientada a Eventos?

**Arquitetura Orientada a Eventos (Event-Driven Architecture)** é um padrão de design onde o fluxo do programa é determinado por eventos - ocorrências ou mudanças de estado que são significativas para o sistema.

### Analogia do Mundo Real

Imagine uma pizzaria:

- **Sem Eventos (Tradicional)**: O cliente pede → atendente vai até a cozinha → espera o pizzaiolo fazer → volta com a pizza
- **Com Eventos**: Cliente pede → atendente grita "PEDIDO!" (evento) → pizzaiolo ouve e começa a fazer → entregador ouve quando fica pronta e entrega

No segundo caso:
- O atendente não fica bloqueado esperando
- Múltiplos funcionários podem reagir ao mesmo evento
- Cada um faz sua parte independentemente

### Na Programação

```
Produtor → [Evento] → Dispatcher → Handler 1
                                 → Handler 2
                                 → Handler 3
```

- **Produtor**: Gera o evento (ex: "Pedido Criado")
- **Dispatcher**: Distribui o evento para interessados
- **Handlers**: Executam ações em resposta (ex: enviar email, atualizar estoque, processar pagamento)

---

## 💡 Por que usar Eventos?

### Vantagens

✅ **Desacoplamento**: Componentes não precisam conhecer uns aos outros diretamente

✅ **Escalabilidade**: Handlers podem ser executados em paralelo ou em servidores diferentes

✅ **Extensibilidade**: Adicionar novos comportamentos sem modificar código existente

✅ **Resiliência**: Se um handler falha, outros continuam funcionando

✅ **Auditoria**: Eventos formam um log natural das ações do sistema

### Quando Usar?

- Sistemas com múltiplos efeitos colaterais (ex: criar pedido → enviar email + atualizar estoque + gerar nota fiscal)
- Microserviços que precisam se comunicar
- Sistemas que precisam de histórico de mudanças
- Aplicações que requerem processamento assíncrono

---

## 🔤 Conceitos Fundamentais em Go

### Interfaces

Interfaces em Go definem **contratos** - um conjunto de métodos que um tipo deve implementar. Diferente de outras linguagens, em Go a implementação é **implícita** (não precisa declarar que implementa).

```go
// Interface define o contrato
type EventInterface interface {
    GetName() string
    GetDateTime() time.Time
    GetPayload() interface{}
}

// Struct que implementa a interface
type OrderCreated struct {
    Name     string
    DateTime time.Time
    Payload  interface{}
}

// Implementando os métodos (implícito)
func (e *OrderCreated) GetName() string {
    return e.Name
}

func (e *OrderCreated) GetDateTime() time.Time {
    return e.DateTime
}

func (e *OrderCreated) GetPayload() interface{} {
    return e.Payload
}

// Agora OrderCreated implementa EventInterface automaticamente!
```

**Por que usar interfaces?**
- **Abstração**: Trabalhe com contratos, não implementações concretas
- **Testabilidade**: Fácil criar mocks e test doubles
- **Flexibilidade**: Trocar implementações sem quebrar código

### Maps como Registry

Maps em Go são estruturas de dados de chave-valor (como dicionários em Python ou objetos em JavaScript).

```go
// Criando um map
handlers := make(map[string][]EventHandlerInterface)

// string: nome do evento (chave)
// []EventHandlerInterface: slice de handlers (valor)

// Adicionando valores
handlers["order.created"] = []EventHandlerInterface{handler1, handler2}

// Acessando valores (forma segura)
if value, ok := handlers["order.created"]; ok {
    // encontrou! 'value' contém o slice de handlers
    fmt.Println("Achei os handlers:", value)
} else {
    // não encontrou
    fmt.Println("Evento não tem handlers registrados")
}
```

**Padrão Registry**: Usar maps para armazenar e recuperar objetos por chave é chamado de **Registry Pattern**. Perfeito para nosso Event Dispatcher!

### Slices e Manipulação

Slices são arrays dinâmicos em Go - podem crescer e diminuir.

```go
// Criando um slice vazio
handlers := []EventHandlerInterface{}

// Adicionando elementos
handlers = append(handlers, newHandler) // agora tem 1 elemento
handlers = append(handlers, handler2)   // agora tem 2 elementos

// Removendo elemento na posição i
i := 1 // índice do elemento a remover
handlers = append(handlers[:i], handlers[i+1:]...)

// Como funciona:
// handlers[:i]     → pega tudo ANTES de i
// handlers[i+1:]   → pega tudo DEPOIS de i
// ... → expande o slice
// append junta as duas partes, pulando o elemento em i
```

**Visualização da Remoção**:
```
Original: [A, B, C, D, E]
Remover índice 2 (C):
  handlers[:2]  = [A, B]
  handlers[3:]  = [D, E]
  append = [A, B, D, E]
```

### Goroutines

Goroutines são **threads leves** gerenciadas pelo runtime do Go. Extremamente baratas em memória (~2KB) comparadas a threads tradicionais (~1-2MB).

```go
// Função normal (bloqueante)
func processaPedido() {
    time.Sleep(2 * time.Second)
    fmt.Println("Pedido processado!")
}

func main() {
    processaPedido() // main espera 2 segundos aqui
    fmt.Println("Fim") // só executa depois
}

// Com goroutine (não-bloqueante)
func main() {
    go processaPedido() // executa em paralelo
    fmt.Println("Fim") // executa imediatamente!

    // Problema: programa pode terminar antes da goroutine!
    time.Sleep(3 * time.Second) // hack para esperar
}
```

**Uso no Event Dispatcher**:
```go
for _, handler := range handlers {
    go handler.Handle(event) // cada handler em sua própria goroutine
}
```

### Channels

Channels são **pipes** para comunicação entre goroutines. Pense neles como filas thread-safe.

```go
// Criando um channel
messages := make(chan string)

// Enviando valor (bloqueia até alguém receber)
go func() {
    messages <- "ping" // operador <- envia
}()

// Recebendo valor (bloqueia até receber)
msg := <-messages // operador <- recebe
fmt.Println(msg) // "ping"

// Channel com buffer (não bloqueia até encher)
buffered := make(chan string, 2)
buffered <- "hello"
buffered <- "world"
// buffered <- "!" // bloquearia aqui (buffer cheio)
```

**Uso com RabbitMQ**:
```go
msgs := make(chan amqp.Delivery) // channel de mensagens

go rabbitmq.Consume(ch, msgs, "orders") // goroutine produz

for msg := range msgs { // loop consome
    fmt.Println(string(msg.Body))
}
```

### sync.WaitGroup

WaitGroup é um **contador de goroutines** que permite esperar até que todas completem.

```go
import "sync"

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // decrementa contador ao finalizar

    fmt.Printf("Worker %d começou\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d terminou\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1) // incrementa contador
        go worker(i, &wg)
    }

    wg.Wait() // bloqueia até contador chegar a 0
    fmt.Println("Todos workers terminaram!")
}
```

**Fluxo**:
1. `Add(n)`: Adiciona n ao contador
2. Goroutine executa
3. `Done()`: Decrementa 1 do contador
4. `Wait()`: Bloqueia até contador = 0

### Defer Statement

`defer` agenda uma função para executar **quando a função atual terminar**, mesmo se houver erro ou panic.

```go
func processar() error {
    arquivo, err := os.Open("dados.txt")
    if err != nil {
        return err
    }
    defer arquivo.Close() // garante que fecha, aconteça o que acontecer

    // processa arquivo...

    // arquivo.Close() será chamado aqui automaticamente
    return nil
}
```

**Ordem de Execução** (LIFO - Last In First Out):
```go
func exemplo() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    fmt.Println("início")
}

// Saída:
// início
// 3
// 2
// 1
```

**Uso comum**:
- Fechar arquivos, conexões de banco
- Desbloquear mutexes
- Liberar recursos
- Cleanup em geral

---

## 📚 Progressão do Aprendizado

### 1. Criando Interfaces da Solução

**Localização**: `1-criando-interfaces-da-solucao/`

**O que você aprende**: Desenhar contratos para um sistema orientado a eventos usando interfaces.

#### Interfaces Criadas

```go
// EventInterface: Define o que é um evento
type EventInterface interface {
    GetName() string        // Nome do evento (ex: "order.created")
    GetDateTime() time.Time // Quando aconteceu
    GetPayload() interface{} // Dados do evento
}

// EventHandlerInterface: Define quem processa eventos
type EventHandlerInterface interface {
    Handle(event EventInterface) // Lógica de processamento
}

// EventDispatcherInterface: Define o orquestrador
type EventDispatcherInterface interface {
    Register(eventName string, handler EventHandlerInterface) error
    Dispatch(event EventInterface) error
    Remove(eventName string, handler EventHandlerInterface) error
    Has(eventName string, handler EventHandlerInterface) bool
    Clear() error
}
```

#### Implementação Básica

```go
// EventDispatcher armazena handlers organizados por nome do evento
type EventDispatcher struct {
    handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
    return &EventDispatcher{
        handlers: make(map[string][]EventHandlerInterface),
    }
}

// Register adiciona um handler para um evento específico
func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
    // Verifica se já existe handlers para este evento
    if _, ok := ed.handlers[eventName]; ok {
        // Percorre handlers existentes
        for _, h := range ed.handlers[eventName] {
            if h == handler {
                // Handler duplicado!
                return ErrHandlerAlreadyRegistered
            }
        }
    }

    // Adiciona o handler ao slice
    ed.handlers[eventName] = append(ed.handlers[eventName], handler)
    return nil
}
```

**Conceitos-chave**:
- **Interfaces** definem comportamento, não implementação
- **Map** organiza handlers por nome do evento
- **Slice** permite múltiplos handlers por evento
- **Error customizado** para casos específicos

```go
var ErrHandlerAlreadyRegistered = errors.New("handler already registered")
```

---

### 2-5. Criando Suite de Testes

**Localização**: `2-criando-suite-testes/` até `5-testando-metodo-has/`

**O que você aprende**: Test-Driven Development (TDD) usando testify/suite para testes organizados.

#### 2. Configurando a Suite

```go
import (
    "github.com/stretchr/testify/suite"
    "testing"
)

// Suite agrupa testes relacionados
type EventDispatcherTestSuite struct {
    suite.Suite
    event      TestEvent
    event2     TestEvent
    handler    TestEventHandler
    handler2   TestEventHandler
    handler3   TestEventHandler
    dispatcher *EventDispatcher
}

// SetupTest roda antes de CADA teste (fresh start)
func (suite *EventDispatcherTestSuite) SetupTest() {
    suite.dispatcher = NewEventDispatcher()
    suite.handler = TestEventHandler{}
    suite.handler2 = TestEventHandler{}
    suite.handler3 = TestEventHandler{}
    suite.event = TestEvent{Name: "test", Payload: "test"}
    suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
}

// Executar a suite
func TestSuite(t *testing.T) {
    suite.Run(t, new(EventDispatcherTestSuite))
}
```

**Test Doubles** (objetos para testes):

```go
type TestEvent struct {
    Name    string
    Payload interface{}
}

func (e *TestEvent) GetName() string        { return e.Name }
func (e *TestEvent) GetDateTime() time.Time { return time.Now() }
func (e *TestEvent) GetPayload() interface{} { return e.Payload }

type TestEventHandler struct{}

func (h *TestEventHandler) Handle(event EventInterface) {
    // Handler vazio para testes
}
```

#### 3. Testando Register

```go
func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
    err := suite.dispatcher.Register(suite.event.GetName(), &suite.handler)

    suite.Nil(err) // Não deve dar erro
    suite.Equal(1, len(suite.dispatcher.handlers[suite.event.GetName()]))

    // Tentar registrar duplicado
    err = suite.dispatcher.Register(suite.event.GetName(), &suite.handler)
    suite.Equal(ErrHandlerAlreadyRegistered, err)
}
```

#### 4. Testando Clear

```go
func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
    // Registra vários handlers
    suite.dispatcher.Register(suite.event.GetName(), &suite.handler)
    suite.dispatcher.Register(suite.event2.GetName(), &suite.handler2)

    // Clear deve limpar tudo
    err := suite.dispatcher.Clear()

    suite.Nil(err)
    suite.Equal(0, len(suite.dispatcher.handlers))
}

// Implementação do Clear
func (ed *EventDispatcher) Clear() error {
    ed.handlers = make(map[string][]EventHandlerInterface)
    return nil
}
```

#### 5. Testando Has

```go
func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
    suite.dispatcher.Register(suite.event.GetName(), &suite.handler)

    // Deve encontrar o handler registrado
    has := suite.dispatcher.Has(suite.event.GetName(), &suite.handler)
    suite.True(has)

    // Não deve encontrar handler não registrado
    has = suite.dispatcher.Has(suite.event.GetName(), &suite.handler2)
    suite.False(has)
}

// Implementação do Has
func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
    if _, ok := ed.handlers[eventName]; ok {
        for _, h := range ed.handlers[eventName] {
            if h == handler {
                return true
            }
        }
    }
    return false
}
```

**Por que usar testify/suite?**
- **SetupTest**: Garante ambiente limpo para cada teste
- **Assertions**: Sintaxe clara (`suite.Equal`, `suite.True`, etc.)
- **Organização**: Agrupa testes relacionados
- **Reutilização**: Compartilha setup comum

---

### 6. Implementando o Método Dispatch

**Localização**: `6-implementando-metodo-dispatch/`

**O que você aprende**: Implementar o núcleo do Event Dispatcher e usar mocks para verificar comportamento.

#### Implementação Básica (Síncrona)

```go
func (ed *EventDispatcher) Dispatch(event EventInterface) error {
    // Verifica se existe handlers para este evento
    if handlers, ok := ed.handlers[event.GetName()]; ok {
        // Chama cada handler sequencialmente
        for _, handler := range handlers {
            handler.Handle(event) // Bloqueia até terminar
        }
    }
    return nil
}
```

**Fluxo**:
1. Recebe um evento
2. Busca handlers registrados para aquele evento
3. Executa cada handler em ordem
4. Retorna (sem esperar handlers terminarem)

#### Testando com Mocks

```go
import "github.com/stretchr/testify/mock"

// Mock substitui implementação real
type MockHandler struct {
    mock.Mock
}

func (m *MockHandler) Handle(event EventInterface) {
    m.Called(event) // Registra que foi chamado
}

// Teste verifica se handler foi chamado
func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
    eh := &MockHandler{}
    eh.On("Handle", &suite.event).Return() // Configura expectativa

    suite.dispatcher.Register(suite.event.GetName(), eh)
    suite.dispatcher.Dispatch(&suite.event)

    eh.AssertExpectations(suite.T()) // Verifica se foi chamado conforme esperado
    eh.AssertNumberOfCalls(suite.T(), "Handle", 1) // Exatamente 1 vez
}
```

**Por que Mocks?**
- **Verificação de comportamento**: Confirma que handlers foram chamados
- **Isolamento**: Testa dispatcher sem depender de implementações reais
- **Flexibilidade**: Simula diferentes cenários facilmente

---

### 7. Revisitando Slices

**Localização**: `7-revisitando-slices/`

**O que você aprende**: Operações avançadas com slices usadas no Event Dispatcher.

#### Operações Fundamentais

```go
// Inicialização
handlers := []EventHandlerInterface{}
handlers := make([]EventHandlerInterface, 0) // equivalente

// Append (adicionar ao final)
handlers = append(handlers, handler1)
handlers = append(handlers, handler2, handler3) // múltiplos

// Iteração
for index, handler := range handlers {
    fmt.Printf("Handler %d: %v\n", index, handler)
}

// Iteração ignorando índice
for _, handler := range handlers {
    handler.Handle(event)
}

// Slicing (fatiar)
first := handlers[0]           // primeiro elemento
last := handlers[len(handlers)-1] // último elemento
middle := handlers[1:3]        // elementos índice 1 e 2 (não inclui 3)
fromStart := handlers[:2]      // primeiros 2 elementos
toEnd := handlers[2:]          // do índice 2 até o final
```

#### Removendo Elementos

```go
// Remover elemento no índice i
i := 2
handlers = append(handlers[:i], handlers[i+1:]...)

// Exemplo visual:
// [A, B, C, D, E] → remover índice 2 (C)
// handlers[:2] = [A, B]
// handlers[3:] = [D, E]
// append(...) = [A, B, D, E]
```

**Por que `...` (spread operator)?**
```go
// handlers[i+1:] retorna um slice []EventHandlerInterface
// append espera elementos individuais, não um slice
// ... "desempacota" o slice em elementos individuais

// Sem ...:
append(handlers[:i], handlers[i+1:])   // ❌ Erro de tipo

// Com ...:
append(handlers[:i], handlers[i+1:]...) // ✅ Correto
```

---

### 8. Removendo Handlers

**Localização**: `8-removendo-handlers/`

**O que você aprende**: Implementar remoção dinâmica de handlers.

#### Implementação

```go
func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
    // Verifica se existe handlers para este evento
    if _, ok := ed.handlers[eventName]; ok {
        // Procura o handler específico
        for i, h := range ed.handlers[eventName] {
            if h == handler {
                // Remove usando slice tricks
                ed.handlers[eventName] = append(
                    ed.handlers[eventName][:i],
                    ed.handlers[eventName][i+1:]...,
                )
                return nil
            }
        }
    }
    return nil
}
```

#### Teste

```go
func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
    // Registra handlers
    suite.dispatcher.Register(suite.event.GetName(), &suite.handler)
    suite.dispatcher.Register(suite.event.GetName(), &suite.handler2)

    // Remove um handler
    err := suite.dispatcher.Remove(suite.event.GetName(), &suite.handler)

    suite.Nil(err)
    suite.Equal(1, len(suite.dispatcher.handlers[suite.event.GetName()]))

    // Verifica que o handler correto foi removido
    suite.False(suite.dispatcher.Has(suite.event.GetName(), &suite.handler))
    suite.True(suite.dispatcher.Has(suite.event.GetName(), &suite.handler2))
}
```

**Caso de Uso Real**:
Imagine um handler temporário que se registra apenas durante uma sessão do usuário. Quando o usuário faz logout, remove o handler.

---

### 9. Adicionando Goroutines ao Dispatcher

**Localização**: `9-adicionando-go-routine-no-event-dispatcher/`

**O que você aprende**: Transformar dispatcher síncrono em assíncrono usando goroutines.

#### Mudança Chave

```go
// ANTES (síncrono)
func (ed *EventDispatcher) Dispatch(event EventInterface) error {
    if handlers, ok := ed.handlers[event.GetName()]; ok {
        for _, handler := range handlers {
            handler.Handle(event) // Bloqueia até terminar
        }
    }
    return nil
}

// DEPOIS (assíncrono)
func (ed *EventDispatcher) Dispatch(event EventInterface) error {
    if handlers, ok := ed.handlers[event.GetName()]; ok {
        for _, handler := range handlers {
            go handler.Handle(event) // Não bloqueia!
        }
    }
    return nil
}
```

#### Impacto

**Execução Síncrona**:
```
Dispatch chamado
  → Handler 1 executa (2s)
  → Handler 2 executa (2s)
  → Handler 3 executa (2s)
Dispatch retorna após 6s
```

**Execução com Goroutines**:
```
Dispatch chamado
  → lança Handler 1 (goroutine)
  → lança Handler 2 (goroutine)
  → lança Handler 3 (goroutine)
Dispatch retorna imediatamente!
Handlers executam em paralelo (2s total)
```

#### Problema

```go
func main() {
    dispatcher.Dispatch(event)
    fmt.Println("Fim") // Pode imprimir antes dos handlers terminarem!
}
```

O programa pode terminar antes das goroutines completarem porque `Dispatch` retorna imediatamente.

**Quando usar**:
- Handlers independentes (não dependem um do outro)
- Operações que podem demorar (enviar email, chamar API)
- Necessidade de alta performance

---

### 10. Utilizando WaitGroup para Sincronização

**Localização**: `10-utilizando-go-routines-no-dispatcher/`

**O que você aprende**: Sincronizar goroutines usando `sync.WaitGroup` para garantir que todos os handlers completem.

#### Mudança na Interface

```go
// Interface ANTIGA
type EventHandlerInterface interface {
    Handle(event EventInterface)
}

// Interface NOVA (com WaitGroup)
type EventHandlerInterface interface {
    Handle(event EventInterface, wg *sync.WaitGroup)
}
```

#### Implementação Atualizada

```go
import "sync"

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
    if handlers, ok := ed.handlers[event.GetName()]; ok {
        // Cria um WaitGroup
        wg := &sync.WaitGroup{}

        for _, handler := range handlers {
            wg.Add(1) // Incrementa contador antes de lançar goroutine
            go handler.Handle(event, wg) // Passa WaitGroup
        }

        wg.Wait() // Bloqueia até todos handlers chamarem Done()
    }
    return nil
}
```

#### Handler Atualizado

```go
type EmailHandler struct{}

func (h *EmailHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
    defer wg.Done() // Garante que Done() é chamado ao finalizar

    // Lógica do handler
    fmt.Println("Enviando email...")
    time.Sleep(2 * time.Second)
    fmt.Println("Email enviado!")
}
```

**Fluxo Completo**:
```
1. Dispatch chamado
2. WaitGroup criado com contador = 0
3. Para cada handler:
   a. wg.Add(1) → contador = 1, 2, 3...
   b. go handler.Handle(event, wg) → goroutine lançada
4. wg.Wait() → bloqueia
5. Handlers executam em paralelo
6. Cada handler chama wg.Done() → contador decrementa
7. Quando contador = 0, wg.Wait() desbloqueia
8. Dispatch retorna
```

#### Teste com Mock Atualizado

```go
type MockHandler struct {
    mock.Mock
}

func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
    m.Called(event)
    wg.Done() // IMPORTANTE: mock também precisa chamar Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
    eh := &MockHandler{}
    eh.On("Handle", &suite.event).Return()

    suite.dispatcher.Register(suite.event.GetName(), eh)
    suite.dispatcher.Dispatch(&suite.event)

    eh.AssertExpectations(suite.T())
    eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
}
```

**Vantagens**:
- ✅ Handlers executam em paralelo (performance)
- ✅ Dispatch espera todos completarem (garantia)
- ✅ Não há race conditions
- ✅ Código mais confiável

---

### 11. Integração com RabbitMQ

**Localização**: `11-rabbit-mq/`

**O que você aprende**: Arquitetura distribuída usando message broker para comunicação entre serviços.

#### O que é RabbitMQ?

**RabbitMQ** é um **message broker** - um intermediário que recebe, armazena e entrega mensagens entre aplicações.

**Analogia**: Pense no RabbitMQ como os Correios:
- **Producer (Remetente)**: Envia uma carta
- **RabbitMQ (Correios)**: Armazena e entrega
- **Consumer (Destinatário)**: Recebe a carta

**Por que usar?**
- **Desacoplamento**: Producer e Consumer não precisam estar online ao mesmo tempo
- **Resiliência**: Mensagens não se perdem se Consumer cair
- **Escalabilidade**: Múltiplos Consumers podem processar em paralelo
- **Distribuição**: Serviços podem estar em servidores diferentes

#### Arquitetura

```
Producer → [Exchange] → [Queue] → Consumer
```

- **Exchange**: Roteador de mensagens (decide para qual fila enviar)
- **Queue**: Fila que armazena mensagens
- **Binding**: Liga Exchange a Queue
- **Message**: Dados sendo transmitidos

#### Configuração com Docker

```yaml
# docker-compose.yaml
version: '3'

services:
  rabbitmq:
    image: rabbitmq:3.8.16-management
    ports:
      - "5672:5672"   # Porta AMQP (protocolo)
      - "15672:15672" # Porta Management UI
    environment:
      - RABBITMQ_DEFAULT_USER=guest
      - RABBITMQ_DEFAULT_PASS=guest
```

**Comandos**:
```bash
# Iniciar RabbitMQ
docker compose up -d

# Ver logs
docker compose logs -f rabbitmq

# Parar RabbitMQ
docker compose down

# Acessar Management UI
open http://localhost:15672
# Login: guest / guest
```

#### Dependências

```bash
go get github.com/rabbitmq/amqp091-go
```

```go
// go.mod
require github.com/rabbitmq/amqp091-go v1.5.0
```

#### Funções Helper

```go
package rabbitmq

import (
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
)

// OpenChannel abre conexão e cria canal
func OpenChannel() (*amqp.Channel, error) {
    // Conecta ao RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        panic(err)
    }

    // Cria canal de comunicação
    ch, err := conn.Channel()
    if err != nil {
        panic(err)
    }

    return ch, nil
}

// Publish envia mensagem para um exchange
func Publish(ch *amqp.Channel, body string, exName string) error {
    err := ch.Publish(
        exName,  // exchange
        "",      // routing key (vazio = broadcast)
        false,   // mandatory
        false,   // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body), // conteúdo da mensagem
        },
    )

    if err != nil {
        return err
    }

    fmt.Println("Mensagem publicada:", body)
    return nil
}

// Consume recebe mensagens de uma fila
func Consume(ch *amqp.Channel, out chan<- amqp.Delivery, queue string) error {
    // Registra consumer na fila
    msgs, err := ch.Consume(
        queue,        // nome da fila
        "go-consumer", // consumer tag
        false,        // auto-ack (false = manual)
        false,        // exclusive
        false,        // no-local
        false,        // no-wait
        nil,          // args
    )
    if err != nil {
        return err
    }

    // Loop infinito consumindo mensagens
    for msg := range msgs {
        out <- msg // Envia para Go channel
    }

    return nil
}
```

#### Producer (Publicador de Eventos)

```go
// cmd/producer/main.go
package main

import (
    "fullcycle-mba-go-expert/aulas/9-eventos/11-rabbit-mq/pkg/rabbitmq"
)

func main() {
    // Abre conexão com RabbitMQ
    ch, err := rabbitmq.OpenChannel()
    if err != nil {
        panic(err)
    }
    defer ch.Close() // Garante fechamento

    // Publica mensagem
    err = rabbitmq.Publish(ch, "Hello World!", "amq.direct")
    if err != nil {
        panic(err)
    }

    println("Mensagem enviada com sucesso!")
}
```

**Executar**:
```bash
cd cmd/producer
go run main.go
```

#### Consumer (Processador de Eventos)

```go
// cmd/consumer/main.go
package main

import (
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
    "fullcycle-mba-go-expert/aulas/9-eventos/11-rabbit-mq/pkg/rabbitmq"
)

func main() {
    // Abre conexão
    ch, err := rabbitmq.OpenChannel()
    if err != nil {
        panic(err)
    }
    defer ch.Close()

    // Cria channel Go para receber mensagens
    msgs := make(chan amqp.Delivery)

    // Inicia consumer em goroutine
    go rabbitmq.Consume(ch, msgs, "orders")

    // Loop infinito processando mensagens
    for msg := range msgs {
        fmt.Println("Mensagem recebida:", string(msg.Body))
        msg.Ack(false) // Confirma processamento (importante!)
    }
}
```

**Executar** (em outro terminal):
```bash
cd cmd/consumer
go run main.go
```

#### Fluxo Completo

```
1. Producer conecta ao RabbitMQ
2. Producer publica mensagem no exchange "amq.direct"
3. Exchange roteia para fila "orders"
4. Mensagem fica armazenada na fila
5. Consumer conecta ao RabbitMQ
6. Consumer subscreve na fila "orders"
7. RabbitMQ entrega mensagem ao Consumer
8. Consumer processa e envia ACK
9. RabbitMQ remove mensagem da fila
```

#### Conceitos Avançados

**Manual Acknowledgment**:
```go
msg.Ack(false) // Confirma processamento
```

- **false**: Confirma apenas esta mensagem
- Se não chamar `Ack()`: RabbitMQ reenvia a mensagem (útil se Consumer falhar)

**Auto-Ack vs Manual-Ack**:
```go
// Auto-Ack: RabbitMQ remove mensagem assim que envia
ch.Consume(queue, "consumer", true, ...)

// Manual-Ack: RabbitMQ espera confirmação
ch.Consume(queue, "consumer", false, ...)
msg.Ack(false) // precisa chamar manualmente
```

**Uso Prático**:
```go
for msg := range msgs {
    err := processarPedido(msg.Body)
    if err != nil {
        msg.Nack(false, true) // Rejeita e reinsere na fila
    } else {
        msg.Ack(false) // Sucesso! Remove da fila
    }
}
```

**Exchanges**:

| Tipo | Comportamento |
|------|---------------|
| **direct** | Roteia por routing key exata |
| **fanout** | Broadcast para todas filas |
| **topic** | Roteia por padrão (ex: "order.*") |
| **headers** | Roteia por headers HTTP-like |

**Exemplo com Exchange direto**:
```go
// Exchange "amq.direct" já existe no RabbitMQ
// Para criar seu próprio:
ch.ExchangeDeclare(
    "logs",   // nome
    "fanout", // tipo
    true,     // durable
    false,    // auto-deleted
    false,    // internal
    false,    // no-wait
    nil,      // arguments
)
```

---

## 🎨 Padrões de Design

### Observer Pattern (Pub/Sub)

**Definição**: Permite que objetos sejam notificados automaticamente quando algo acontece, sem acoplamento direto.

```go
// Subject (EventDispatcher) mantém lista de Observers (Handlers)
dispatcher.Register("order.created", emailHandler)
dispatcher.Register("order.created", smsHandler)
dispatcher.Register("order.created", logHandler)

// Quando evento acontece, todos Observers são notificados
dispatcher.Dispatch(orderCreatedEvent)
// → emailHandler é notificado
// → smsHandler é notificado
// → logHandler é notificado
```

### Registry Pattern

**Definição**: Armazena e recupera objetos por chave.

```go
// Map funciona como Registry
handlers := map[string][]EventHandlerInterface{
    "order.created": {emailHandler, smsHandler},
    "order.cancelled": {refundHandler},
}

// Recupera por chave
handlersForOrder := handlers["order.created"]
```

### Worker Pool Pattern

**Definição**: Pool de workers (goroutines) processam tarefas concorrentemente.

```go
// Cada handler é um worker
for _, handler := range handlers {
    wg.Add(1)
    go handler.Handle(event, wg) // Worker executa tarefa
}
wg.Wait() // Espera todos workers terminarem
```

### Producer-Consumer Pattern

**Definição**: Producers geram trabalho, Consumers processam.

```go
// Producer
msgs := make(chan amqp.Delivery)
go rabbitmq.Consume(ch, msgs, "orders") // Produz mensagens no channel

// Consumer
for msg := range msgs { // Consome mensagens
    processar(msg)
}
```

---

## 💼 Aplicações Práticas

### E-commerce: Processamento de Pedidos

```go
// Evento: Pedido criado
type OrderCreated struct {
    OrderID   string
    CustomerID string
    Total     float64
}

// Handlers independentes
type SendEmailHandler struct{}
func (h *SendEmailHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
    defer wg.Done()
    order := event.GetPayload().(OrderCreated)
    sendEmail(order.CustomerID, "Pedido confirmado!")
}

type UpdateInventoryHandler struct{}
func (h *UpdateInventoryHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
    defer wg.Done()
    order := event.GetPayload().(OrderCreated)
    decrementStock(order.OrderID)
}

type NotifyShippingHandler struct{}
func (h *NotifyShippingHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
    defer wg.Done()
    order := event.GetPayload().(OrderCreated)
    createShippingLabel(order.OrderID)
}

// Uso
dispatcher.Register("order.created", &SendEmailHandler{})
dispatcher.Register("order.created", &UpdateInventoryHandler{})
dispatcher.Register("order.created", &NotifyShippingHandler{})

// Quando pedido é criado
dispatcher.Dispatch(orderCreatedEvent)
// → Email enviado (em paralelo)
// → Estoque atualizado (em paralelo)
// → Etiqueta criada (em paralelo)
```

### Microserviços: Comunicação Assíncrona

```
[Order Service] → RabbitMQ → [Email Service]
                           → [Inventory Service]
                           → [Shipping Service]
```

- **Order Service** publica evento "order.created"
- **Email Service** consome e envia email
- **Inventory Service** consome e atualiza estoque
- **Shipping Service** consome e cria label

**Vantagens**:
- Serviços independentes (podem ser em Go, Python, Java...)
- Não bloqueia Order Service esperando processamento
- Se Email Service cai, mensagens ficam na fila esperando

### Sistema de Auditoria

```go
type AuditHandler struct {
    db *sql.DB
}

func (h *AuditHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
    defer wg.Done()

    // Registra todos eventos no banco
    h.db.Exec(`
        INSERT INTO audit_log (event_name, payload, timestamp)
        VALUES (?, ?, ?)
    `, event.GetName(), event.GetPayload(), event.GetDateTime())
}

// Registra em TODOS os eventos
dispatcher.Register("order.created", auditHandler)
dispatcher.Register("user.registered", auditHandler)
dispatcher.Register("payment.processed", auditHandler)
// → Histórico completo de tudo que acontece no sistema
```

---

## 🛠️ Comandos Úteis

### Testes

```bash
# Executar todos os testes
go test ./...

# Executar com verbose
go test -v ./...

# Executar testes de um pacote específico
cd 1-criando-interfaces-da-solucao
go test

# Executar teste específico
go test -run TestEventDispatcher_Register

# Cobertura de testes
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### RabbitMQ

```bash
# Subir RabbitMQ com Docker
docker compose up -d

# Ver logs
docker compose logs -f rabbitmq

# Parar RabbitMQ
docker compose down

# Acessar Management UI
open http://localhost:15672
# Usuário: guest
# Senha: guest

# Executar Producer
cd cmd/producer
go run main.go

# Executar Consumer (em outro terminal)
cd cmd/consumer
go run main.go
```

### Go Modules

```bash
# Inicializar módulo
go mod init github.com/seu-usuario/seu-projeto

# Baixar dependências
go mod download

# Adicionar dependência
go get github.com/rabbitmq/amqp091-go

# Limpar dependências não usadas
go mod tidy

# Verificar dependências
go list -m all
```

### Goroutines

```bash
# Detectar race conditions
go test -race

# Build com race detector
go build -race

# Executar com race detector
go run -race main.go
```

---

## 📖 Glossário

- **Event (Evento)**: Ocorrência significativa no sistema
- **Handler (Manipulador)**: Função/objeto que reage a um evento
- **Dispatcher (Despachador)**: Orquestrador que distribui eventos para handlers
- **Publisher/Producer (Publicador/Produtor)**: Quem gera eventos
- **Subscriber/Consumer (Assinante/Consumidor)**: Quem processa eventos
- **Message Broker**: Intermediário para mensagens (ex: RabbitMQ)
- **Queue (Fila)**: Armazena mensagens aguardando processamento
- **Exchange**: Roteador de mensagens no RabbitMQ
- **AMQP**: Advanced Message Queuing Protocol
- **Goroutine**: Thread leve do Go
- **Channel**: Pipe para comunicação entre goroutines
- **WaitGroup**: Sincronizador de goroutines
- **Mock**: Objeto falso para testes
- **Test Suite**: Conjunto organizado de testes relacionados

---

## 🎓 Conceitos Aprendidos

Ao completar este módulo, você domina:

✅ **Arquitetura Orientada a Eventos**
- Princípios e vantagens
- Observer/Pub-Sub pattern
- Event Dispatcher implementation

✅ **Interfaces em Go**
- Definição de contratos
- Implementação implícita
- Polimorfismo

✅ **Concorrência**
- Goroutines para paralelismo
- sync.WaitGroup para sincronização
- Channels para comunicação

✅ **Test-Driven Development**
- Test suites com testify
- Mocks para verificação
- Test doubles

✅ **Mensageria Distribuída**
- RabbitMQ
- Producer/Consumer pattern
- AMQP protocol
- Message acknowledgment

✅ **Boas Práticas Go**
- Error handling
- Defer para cleanup
- Maps e Slices
- Organização de código

---

## 📚 Próximos Passos

1. **Explore variações**: Tente implementar diferentes tipos de Exchanges no RabbitMQ
2. **Adicione features**: Implemente retry logic para handlers que falham
3. **Métricas**: Adicione logging e observabilidade aos eventos
4. **Persistência**: Armazene eventos em banco de dados (Event Sourcing)
5. **Streaming**: Explore Apache Kafka como alternativa ao RabbitMQ
6. **CQRS**: Combine eventos com Command Query Responsibility Segregation

---

**Dica final**: Eventos são poderosos, mas não use em tudo! Para operações simples e síncronas, chamadas diretas de função são mais apropriadas. Use eventos quando precisar de:
- Desacoplamento entre componentes
- Processamento assíncrono
- Múltiplas reações ao mesmo acontecimento
- Comunicação entre microserviços
- Histórico de mudanças (audit log)
