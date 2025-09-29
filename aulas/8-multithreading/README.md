# Multithreading em Go - Guia Didático Completo

Este guia explica de forma didática os conceitos de processos, threads e como o Go revoluciona a programação concorrente com suas goroutines.

## Sumário

### [Conceitos Fundamentais](#conceitos-fundamentais)
- [O que são Processos?](#o-que-são-processos)
- [O que são Threads?](#o-que-são-threads)
- [Problemas Clássicos do Multithreading](#problemas-clássicos-do-multithreading)

### [A Revolução do Go: Goroutines](#a-revolução-do-go-goroutines)
- [O que são Goroutines?](#o-que-são-goroutines)
- [O Scheduler do Go (M:N Threading)](#o-scheduler-do-go-mn-threading)
- [Channels: A Comunicação Segura](#channels-a-comunicação-segura)

### [Arquitetura Completa do Go Runtime](#arquitetura-completa-do-go-runtime)
- [Componentes do GPM Model](#componentes-do-gpm-model)

### [Exemplos Práticos](#exemplos-práticos)
- [1. Goroutine Básica](#1-goroutine-básica)
- [2. Múltiplas Goroutines](#2-múltiplas-goroutines)
- [3. Channels para Comunicação](#3-channels-para-comunicação)
- [4. Worker Pool](#4-worker-pool-pool-de-trabalhadores)

### [Padrões Comuns em Go](#padrões-comuns-em-go)
- [1. Fan-Out (Distribuir trabalho)](#1-fan-out-distribuir-trabalho)
- [2. Fan-In (Coletar resultados)](#2-fan-in-coletar-resultados)
- [3. Pipeline (Cadeia de processamento)](#3-pipeline-cadeia-de-processamento)

### [Vantagens do Modelo Go](#vantagens-do-modelo-go)
- [Performance](#performance)
- [Simplicidade](#simplicidade)
- [Segurança](#segurança)

### [Ferramentas e Debugging](#ferramentas-e-debugging)
- [1. Race Detector](#1-race-detector)
- [2. Profiling](#2-profiling)
- [3. Trace](#3-trace)

### [Resumo dos Conceitos](#resumo-dos-conceitos)

### [Próximos Passos](#próximos-passos)

### [Ferramentas de Sincronização](#ferramentas-de-sincronização)
- [WaitGroups - Coordenando Goroutines](#waitgroups---coordenando-goroutines)
- [Mutex - Exclusão Mútua](#mutex---exclusão-mútua)
- [Operações Atômicas](#operações-atômicas)

### [Channels - Comunicação Segura](#channels---comunicação-segura)
- [Channels Básicos](#channels-básicos)
- [Buffered Channels](#buffered-channels)
- [Channel Directions](#channel-directions)
- [Select Statement](#select-statement)
- [Forever Channels](#forever-channels)

### [Padrões Avançados](#padrões-avançados)
- [Worker Pools](#worker-pools)
- [Load Balancer](#load-balancer)
- [Pipeline Pattern](#pipeline-pattern)
- [Fan-Out / Fan-In](#fan-out--fan-in)

### [Casos de Uso Reais](#casos-de-uso-reais)
- [Servidores Web](#servidores-web)
- [Processamento de Dados](#processamento-de-dados)
- [Message Queues](#message-queues)

### [Trade-offs e Considerações](#trade-offs-e-considerações)
- [Performance vs Complexidade](#performance-vs-complexidade)
- [Segurança vs Velocidade](#segurança-vs-velocidade)
- [Escalabilidade vs Recursos](#escalabilidade-vs-recursos)

---

##  Conceitos Fundamentais

## O que são Processos?

Imagine que seu computador é como uma **grande fábrica**. Cada programa que você executa (como o navegador, editor de texto, ou um jogo) é como uma **linha de produção independente** dentro dessa fábrica.

```
🏭 COMPUTADOR (Fábrica)
├── 📱 Processo Chrome (Linha de Produção A)
├── 📝 Processo VSCode (Linha de Produção B)
├── 🎵 Processo Spotify (Linha de Produção C)
└── 🎮 Processo Jogo (Linha de Produção D)
```

**Características dos Processos:**
- **Isolamento total**: Cada processo tem sua própria área de memória
- **Comunicação custosa**: Para um processo "falar" com outro, precisa de mecanismos especiais
- **Proteção**: Se um processo trava, não afeta os outros
- **Overhead alto**: Criar um novo processo consome muitos recursos

**Analogia Real:**
É como ter departamentos separados em uma empresa. O RH não pode acessar diretamente os arquivos da Contabilidade - eles precisam se comunicar através de protocolos específicos.

## O que são Threads?

Dentro de cada processo (linha de produção), você pode ter várias **threads** - que são como **trabalhadores especializados** operando na mesma linha.

```
📱 PROCESSO CHROME
├── 🧵 Thread Interface (Renderização da tela)
├── 🧵 Thread Network (Downloads e uploads)
├── 🧵 Thread JavaScript (Execução de scripts)
└── 🧵 Thread Database (Gerenciar histórico/cookies)
```

**Características das Threads:**
- **Compartilham memória**: Todas as threads de um processo acessam a mesma área de memória
- **Comunicação rápida**: Podem "conversar" facilmente entre si
- **Risco de conflito**: Se duas threads modificam o mesmo dado simultaneamente, pode dar problema
- **Menos overhead**: Criar uma thread é mais barato que criar um processo

**Analogia Real:**
É como funcionários do mesmo departamento trabalhando na mesma sala, compartilhando os mesmos arquivos e recursos. Eles podem colaborar facilmente, mas precisam se coordenar para não atrapalhar uns aos outros.

## Problemas Clássicos do Multithreading

### 1. Race Condition (Condição de Corrida)
```
👩‍💻 Thread A: "Vou ler o saldo da conta: R$ 1000"
👨‍💻 Thread B: "Vou ler o saldo da conta: R$ 1000"
👩‍💻 Thread A: "Vou debitar R$ 200, novo saldo: R$ 800"
👨‍💻 Thread B: "Vou debitar R$ 300, novo saldo: R$ 700"

❌ RESULTADO: Saldo final R$ 700 (ERRADO!)
✅ DEVERIA SER: R$ 500
```

### 2. Deadlock (Impasse)
```
🧵 Thread A: "Preciso do Recurso X e depois do Y"
🧵 Thread B: "Preciso do Recurso Y e depois do X"

Thread A pega X, Thread B pega Y
Agora ambas ficam esperando para sempre! 🔄
```

### 3. Starvation (Inanição)
```
🧵 Thread VIP: Sempre consegue recursos
🧵 Thread Normal: Nunca consegue executar
```

##  A Revolução do Go: Goroutines

O Go criou uma abordagem revolucionária para resolver os problemas do multithreading tradicional. Em vez de usar threads pesadas do sistema operacional, o Go usa **goroutines**.

## O que são Goroutines?

Pense nas goroutines como **assistentes virtuais super eficientes**:

```
🖥️ SISTEMA OPERACIONAL
└── 🔧 Thread OS (Pesada - 2MB de memória)

🐹 GO RUNTIME
├── 🧚‍♀️ Goroutine 1 (Leve - 2KB de memória)
├── 🧚‍♂️ Goroutine 2 (Leve - 2KB de memória)
├── 🧚‍♀️ Goroutine 3 (Leve - 2KB de memória)
└── 🧚‍♂️ ... milhares de outras goroutines
```

**Vantagens das Goroutines:**
- **Ultra leves**: Começam com apenas 2KB de memória
- **Escaláveis**: Você pode ter milhões delas
- **Gerenciamento automático**: O Go decide quando e onde executá-las
- **Sintaxe simples**: Basta adicionar `go` antes de uma função

## O Scheduler do Go (M:N Threading)

O Go usa um modelo inteligente chamado **M:N Threading**:

```
🏭 GO RUNTIME (Fábrica Inteligente)

📋 GOROUTINES (Tarefas a fazer)
├── 📝 Tarefa 1: Baixar arquivo
├── 🔄 Tarefa 2: Processar dados
├── 📊 Tarefa 3: Gerar relatório
└── 🌐 Tarefa 4: Responder HTTP

👥 THREADS OS (Trabalhadores Físicos)
├── 🧑‍💼 Thread 1 ← Executando Tarefa 1
├── 👩‍💼 Thread 2 ← Executando Tarefa 3
└── 👨‍💼 Thread 3 ← Executando Tarefa 4

🎯 SCHEDULER (Gerente Inteligente)
"Vou distribuir as tarefas entre os trabalhadores de forma otimizada!"
```

**Como funciona:**
- **M Goroutines** são mapeadas para **N Threads OS**
- O scheduler do Go decide qual goroutine executa em qual thread
- Quando uma goroutine "dorme" (esperando I/O), outra pode usar a thread
- Balanceamento automático de carga entre threads

## Channels: A Comunicação Segura

O Go resolve o problema de comunicação entre goroutines com **channels** - canais seguros de comunicação:

```go
// Canal é como um tubo que conecta duas goroutines
canal := make(chan string)

// Goroutine 1: Produtora
go func() {
    canal <- "Olá do produtor!" // Envia dados
}()

// Goroutine 2: Consumidora
mensagem := <-canal // Recebe dados
fmt.Println(mensagem)
```

**Analogia dos Channels:**
É como um **sistema de tubos pneumáticos** em um banco antigo. Você coloca a mensagem no tubo, e ela chega com segurança do outro lado, sem risco de interferência.

##  Arquitetura Completa do Go Runtime

```
┌─────────────────────────────────────────────────────────┐
│                    PROGRAMA GO                          │
├─────────────────────────────────────────────────────────┤
│  🧚‍♀️ Goroutine 1    🧚‍♂️ Goroutine 2    🧚‍♀️ Goroutine N  │
│       ↓                   ↓                   ↓       │
├─────────────────────────────────────────────────────────┤
│               🎯 GO SCHEDULER (GPM Model)               │
│                                                         │
│  📋 G (Goroutines) - Tarefas a executar               │
│  🧑‍💼 P (Processors) - Contextos de execução           │
│  🏃‍♂️ M (Threads)    - Threads do OS                   │
├─────────────────────────────────────────────────────────┤
│            🖥️ SISTEMA OPERACIONAL                      │
│                                                         │
│  Thread 1    Thread 2    Thread 3    Thread N         │
└─────────────────────────────────────────────────────────┘
```

## Componentes do GPM Model

### G (Goroutines)
- **O que são**: As tarefas/funções que você quer executar
- **Características**: Leves, têm seu próprio stack
- **Estado**: Podem estar executando, esperando, ou dormindo

### P (Processors)
- **O que são**: Contextos de execução (geralmente = número de CPUs)
- **Função**: Mantêm filas de goroutines prontas para executar
- **Configuração**: `GOMAXPROCS` define quantos P's existem

### M (Machine/Threads)
- **O que são**: Threads reais do sistema operacional
- **Função**: Executam as goroutines
- **Flexibilidade**: Podem ser criadas/destruídas conforme necessário

##  Exemplos Práticos

## 1. Goroutine Básica
```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Função normal (sequencial)
    fmt.Println("Iniciando programa...")

    // Goroutine (paralela)
    go func() {
        fmt.Println("Executando em paralelo!")
    }()

    // Aguarda um pouco para ver o resultado
    time.Sleep(time.Second)
    fmt.Println("Programa finalizado")
}
```

## 2. Múltiplas Goroutines
```go
func main() {
    // Criando 5 "trabalhadores" paralelos
    for i := 1; i <= 5; i++ {
        go trabalhador(i)
    }

    time.Sleep(2 * time.Second)
}

func trabalhador(id int) {
    fmt.Printf("Trabalhador %d iniciado\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Trabalhador %d finalizado\n", id)
}
```

## 3. Channels para Comunicação
```go
func main() {
    // Canal para comunicação
    resultado := make(chan string)

    // Goroutine que faz uma "busca"
    go buscarDados(resultado)

    // Aguarda o resultado
    dados := <-resultado
    fmt.Println("Recebido:", dados)
}

func buscarDados(canal chan string) {
    // Simula uma operação demorada
    time.Sleep(2 * time.Second)
    canal <- "Dados importantes!"
}
```

## 4. Worker Pool (Pool de Trabalhadores)
```go
func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Criando 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Enviando 5 jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    // Coletando resultados
    for a := 1; a <= 5; a++ {
        <-results
    }
}

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d processando job %d\n", id, j)
        time.Sleep(time.Second)
        results <- j * 2
    }
}
```

##  Padrões Comuns em Go

## 1. Fan-Out (Distribuir trabalho)
```go
// Uma goroutine distribui trabalho para várias
func fanOut(input <-chan int) (<-chan int, <-chan int) {
    out1 := make(chan int)
    out2 := make(chan int)

    go func() {
        for val := range input {
            out1 <- val
            out2 <- val
        }
        close(out1)
        close(out2)
    }()

    return out1, out2
}
```

## 2. Fan-In (Coletar resultados)
```go
// Várias goroutines enviam para uma
func fanIn(input1, input2 <-chan string) <-chan string {
    output := make(chan string)

    go func() {
        for {
            select {
            case msg := <-input1:
                output <- msg
            case msg := <-input2:
                output <- msg
            }
        }
    }()

    return output
}
```

## 3. Pipeline (Cadeia de processamento)
```go
func pipeline() {
    numbers := make(chan int)
    squares := make(chan int)

    // Stage 1: Gerar números
    go func() {
        for i := 1; i <= 10; i++ {
            numbers <- i
        }
        close(numbers)
    }()

    // Stage 2: Calcular quadrados
    go func() {
        for num := range numbers {
            squares <- num * num
        }
        close(squares)
    }()

    // Stage 3: Imprimir resultados
    for square := range squares {
        fmt.Println(square)
    }
}
```

##  Vantagens do Modelo Go

## Performance
- **Milhões de goroutines**: Em vez de centenas de threads
- **Baixo overhead**: 2KB vs 2MB por thread
- **Context switching rápido**: Gerenciado pelo runtime Go

## Simplicidade
- **Sintaxe clara**: `go funcao()` vs configuração complexa de threads
- **Sem locks explícitos**: Use channels para comunicação segura
- **Garbage collector**: Gerenciamento automático de memória

## Segurança
- **Memory safety**: Previne corrupção de memória
- **Race detector**: Ferramenta integrada para detectar race conditions
- **Deadlock detection**: Runtime detecta alguns tipos de deadlock

##  Ferramentas e Debugging

## 1. Race Detector
```bash
# Compila e executa com detector de race conditions
go run -race main.go

# Build com race detector
go build -race
```

## 2. Profiling
```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // Seu código aqui...
}
```

## 3. Trace
```bash
# Gera trace de execução
go run main.go 2> trace.out
go tool trace trace.out
```

##  Resumo dos Conceitos

## Evolução do Paralelismo
```
1️⃣ Sequencial: Uma coisa por vez
2️⃣ Threads: Paralelo, mas complexo e perigoso
3️⃣ Go: Paralelo, simples e seguro
```

## Por que Go é Especial
- **Green Threads**: Goroutines são gerenciadas pelo runtime, não pelo OS
- **CSP (Communicating Sequential Processes)**: Modelo de comunicação por mensagens
- **Built-in Concurrency**: Concorrência é parte da linguagem, não uma biblioteca

## Filosofia Go
> "Don't communicate by sharing memory; share memory by communicating"
>
> "Não comunique compartilhando memória; compartilhe memória comunicando"

Esta frase resume a filosofia do Go: em vez de usar locks e compartilhar variáveis, use channels para trocar informações de forma segura.

##  Próximos Passos

1. **Pratique**: Implemente os exemplos deste guia
2. **Experimente**: Crie seus próprios padrões de concorrência
3. **Meça**: Use as ferramentas de profiling para otimizar
4. **Estude**: Explore bibliotecas como `sync`, `context` e `errgroup`

A concorrência em Go transforma problemas complexos em soluções elegantes e performáticas! 🎯

---

##  Ferramentas de Sincronização

## WaitGroups - Coordenando Goroutines

O `WaitGroup` é como um **contador inteligente** que espera todas as goroutines terminarem antes de continuar.

**Analogia Real:**
É como um chef esperando todos os cozinheiros terminarem suas tarefas antes de servir o prato final.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // ✅ Sempre use defer para garantir que Done() seja chamado

    fmt.Printf("Worker %d iniciado\n", id)
    time.Sleep(time.Second * 2)
    fmt.Printf("Worker %d finalizado\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1) // Incrementa o contador
        go worker(i, &wg)
    }

    wg.Wait() // Espera até o contador chegar a zero
    fmt.Println("Todos os workers terminaram!")
}
```

**Caso de Uso Real:**
- **Processamento de Lotes**: Processar 1000 arquivos em paralelo e esperar todos terminarem
- **Web Scraping**: Fazer múltiplas requisições HTTP e aguardar todas as respostas
- **Testes Paralelos**: Executar suítes de teste em paralelo

**⚠️ Armadilhas Comuns:**
```go
// ❌ ERRO: Deadlock - Add(3) mas só 2 Done()
wg.Add(3)
go func() { defer wg.Done() }()
go func() { defer wg.Done() }()
// Faltou uma goroutine!
wg.Wait() // Trava para sempre

// ✅ CORRETO: Add = Done
wg.Add(2)
go func() { defer wg.Done() }()
go func() { defer wg.Done() }()
wg.Wait()
```

**Prós:**
- ✅ Simples de usar
- ✅ Sincronização precisa
- ✅ Detecta deadlocks em runtime

**Contras:**
- ❌ Só conta execuções, não resultados
- ❌ Não permite cancelamento
- ❌ Reutilização requer cuidado

## Mutex - Exclusão Mútua

O `Mutex` (Mutual Exclusion) é como uma **chave única** para uma sala - apenas uma goroutine pode entrar por vez.

**Analogia Real:**
É como o banheiro de um avião - tem uma tranca, e apenas uma pessoa pode usar por vez.

```go
package main

import (
    "fmt"
    "net/http"
    "sync"
)

var (
    counter int
    mutex   sync.Mutex
)

func incrementCounter(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()   // 🔒 Tranca a "sala"
    counter++      // Área crítica - apenas uma goroutine por vez
    current := counter
    mutex.Unlock() // 🔓 Destranca a "sala"

    fmt.Fprintf(w, "Counter: %d", current)
}

func main() {
    http.HandleFunc("/", incrementCounter)
    http.ListenAndServe(":8080", nil)
}
```

**Problema sem Mutex:**
```go
// ❌ RACE CONDITION
var counter int = 0

// 1000 goroutines incrementando simultaneamente
for i := 0; i < 1000; i++ {
    go func() {
        counter++ // Não é thread-safe!
    }()
}
// Resultado: Pode ser qualquer valor < 1000 😱
```

**Solução com Mutex:**
```go
// ✅ THREAD-SAFE
var (
    counter int
    mu sync.Mutex
)

for i := 0; i < 1000; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
// Resultado: Sempre 1000 ✅
```

**Caso de Uso Real:**
- **Contadores Globais**: Views de página, estatísticas
- **Cache Compartilhado**: Múltiplas goroutines acessando cache
- **Arquivos de Log**: Escrever logs de forma sincronizada

**⚠️ Armadilhas do Mutex:**
```go
// ❌ DEADLOCK: Esqueceu de fazer Unlock
mu.Lock()
if condition {
    return // ❌ Saiu sem Unlock!
}
mu.Unlock()

// ✅ CORRETO: Sempre use defer
mu.Lock()
defer mu.Unlock()
if condition {
    return // ✅ defer garante o Unlock
}
```

**Prós:**
- ✅ Proteção garantida contra race conditions
- ✅ Implementação simples
- ✅ RWMutex permite múltiplos leitures

**Contras:**
- ❌ Pode causar deadlocks
- ❌ Reduz paralelismo
- ❌ Performance menor que operações atômicas

## Operações Atômicas

As operações atômicas são como **movimentos indivisíveis** - acontecem "de uma vez só", sem interrupção.

**Analogia Real:**
É como sacar dinheiro no caixa eletrônico - a operação acontece completamente ou não acontece.

```go
package main

import (
    "fmt"
    "net/http"
    "sync/atomic"
)

var counter uint64

func incrementCounter(w http.ResponseWriter, r *http.Request) {
    // ⚛️ Operação atômica - mais rápida que Mutex
    newValue := atomic.AddUint64(&counter, 1)
    fmt.Fprintf(w, "Counter: %d", newValue)
}

func main() {
    http.HandleFunc("/", incrementCounter)
    http.ListenAndServe(":8080", nil)
}
```

**Comparação de Performance:**
```go
// 🐌 MUTEX (mais lento, mais flexível)
mu.Lock()
counter++
value := counter
mu.Unlock()

// ⚡ ATÔMICA (mais rápido, menos flexível)
value := atomic.AddUint64(&counter, 1)
```

**Operações Atômicas Disponíveis:**
```go
var value uint64

// Adição
atomic.AddUint64(&value, 1)

// Troca (swap)
oldValue := atomic.SwapUint64(&value, 42)

// Comparar e trocar (CAS)
swapped := atomic.CompareAndSwapUint64(&value, 42, 100)

// Carregar (load)
current := atomic.LoadUint64(&value)

// Armazenar (store)
atomic.StoreUint64(&value, 200)
```

**Caso de Uso Real:**
- **Contadores de alta frequência**: Métricas, estatísticas
- **Flags de controle**: Estados simples (ligado/desligado)
- **IDs únicos**: Geradores de ID thread-safe

**Prós:**
- ✅ Performance máxima
- ✅ Livre de deadlocks
- ✅ Menos overhead que Mutex

**Contras:**
- ❌ Limitado a tipos primitivos
- ❌ Não serve para estruturas complexas
- ❌ Menos legível que Mutex

**Trade-off: Mutex vs Atômicas**
```
Operações Simples (contador, flags):
  Atômicas > Mutex

Operações Complexas (estruturas, múltiplas variáveis):
  Mutex > Atômicas

Alta Frequência:
  Atômicas > Mutex

Legibilidade do Código:
  Mutex > Atômicas
```

---

##  Channels - Comunicação Segura

## Channels Básicos

Channels são **tubos de comunicação** entre goroutines, seguindo o princípio: *"Don't communicate by sharing memory; share memory by communicating"*.

**Analogia Real:**
É como um **tubo pneumático** em bancos antigos - você coloca a mensagem, ela viaja pelo tubo e chega do outro lado com segurança.

```go
package main

import (
    "fmt"
    "time"
)

func producer(ch chan<- string) {
    for i := 0; i < 5; i++ {
        message := fmt.Sprintf("Mensagem %d", i)
        ch <- message // Envia para o channel
        time.Sleep(time.Second)
    }
    close(ch) // ✅ Sempre feche o channel quando terminar
}

func consumer(ch <-chan string) {
    for message := range ch { // Range automaticamente para quando o channel é fechado
        fmt.Println("Recebido:", message)
    }
}

func main() {
    ch := make(chan string)

    go producer(ch)
    consumer(ch) // Roda na main goroutine

    fmt.Println("Processamento concluído!")
}
```

**Estados de um Channel:**
```go
ch := make(chan string)

// 1. Channel vazio - receptor bloqueia
message := <-ch // ⏳ Bloqueia até alguém enviar

// 2. Channel com dados - receptor não bloqueia
ch <- "Hello"
message := <-ch // ✅ Recebe imediatamente

// 3. Channel fechado - receptor recebe valor zero
close(ch)
message, ok := <-ch // message="", ok=false
```

**Caso de Uso Real:**
- **Pipeline de Processamento**: Dados fluem de uma etapa para outra
- **Notificações**: Uma goroutine notifica outra sobre eventos
- **Resultados Assíncronos**: Coletar resultados de operações paralelas

## Buffered Channels

Channels com buffer são como uma **fila limitada** - podem armazenar várias mensagens antes de bloquear.

**Analogia Real:**
É como uma caixa de correio - pode guardar várias cartas até ficar cheia.

```go
package main

import "fmt"

func main() {
    // Channel com buffer de 3 mensagens
    ch := make(chan string, 3)

    // Pode enviar 3 mensagens sem bloquear
    ch <- "Primeira"
    ch <- "Segunda"
    ch <- "Terceira"
    // ch <- "Quarta" // ❌ Bloquearia aqui (buffer cheio)

    // Lendo as mensagens
    fmt.Println(<-ch) // "Primeira"
    fmt.Println(<-ch) // "Segunda"
    fmt.Println(<-ch) // "Terceira"
}
```

**Unbuffered vs Buffered:**
```go
// Unbuffered (Síncrono)
ch1 := make(chan int)    // Buffer = 0
ch1 <- 42               // ⏳ Bloqueia até alguém ler

// Buffered (Assíncrono até encher)
ch2 := make(chan int, 5) // Buffer = 5
ch2 <- 42               // ✅ Não bloqueia (buffer tem espaço)
```

**Caso de Uso Real:**
- **Rate Limiting**: Controlar quantas operações simultâneas
- **Batch Processing**: Acumular dados antes de processar
- **Buffering de Logs**: Evitar bloqueios em logging

**Trade-offs:**
```
Unbuffered Channels:
  ✅ Sincronização perfeita
  ✅ Menor uso de memória
  ❌ Mais bloqueios

Buffered Channels:
  ✅ Menos bloqueios
  ✅ Melhor throughput
  ❌ Usa mais memória
  ❌ Pode mascarar problemas de design
```

## Channel Directions

Go permite restringir channels para **apenas envio** ou **apenas recebimento**, melhorando a segurança do código.

```go
package main

import "fmt"

// Função que só pode ENVIAR para o channel
func sender(name string, ch chan<- string) {
    ch <- fmt.Sprintf("Hello from %s", name)
    // message := <-ch // ❌ ERRO: não pode receber de um send-only channel
}

// Função que só pode RECEBER do channel
func receiver(ch <-chan string) {
    message := <-ch
    fmt.Println("Received:", message)
    // ch <- "response" // ❌ ERRO: não pode enviar para um receive-only channel
}

func main() {
    ch := make(chan string) // Channel bidirecional

    go sender("Producer", ch)
    receiver(ch)
}
```

**Conversões Automáticas:**
```go
ch := make(chan string)        // Bidirecional

var sendOnly chan<- string = ch    // ✅ OK: bidirecional → send-only
var recvOnly <-chan string = ch    // ✅ OK: bidirecional → receive-only

// var bidirectional chan string = sendOnly // ❌ ERRO: não pode voltar
```

**Caso de Uso Real:**
- **APIs Claras**: Interface deixa claro quem faz o quê
- **Prevenção de Erros**: Compilador evita uso incorreto
- **Documentação Viva**: Tipo da função documenta o comportamento

## Select Statement

O `select` é como um **switch para channels** - permite lidar com múltiplos channels simultaneamente.

**Analogia Real:**
É como um **porteiro de hotel** que monitora várias portas ao mesmo tempo e atende a primeira que toca.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    // Simulando diferentes fontes de dados
    go func() {
        time.Sleep(2 * time.Second)
        ch1 <- "Dados do Banco de Dados"
    }()

    go func() {
        time.Sleep(1 * time.Second)
        ch2 <- "Dados da API"
    }()

    // Select pega o primeiro que chegar
    select {
    case msg1 := <-ch1:
        fmt.Println("Recebido de ch1:", msg1)
    case msg2 := <-ch2:
        fmt.Println("Recebido de ch2:", msg2) // ✅ Este será executado (mais rápido)
    case <-time.After(3 * time.Second):
        fmt.Println("Timeout! Nenhum channel respondeu a tempo")
    }
}
```

**Select com Default (Non-blocking):**
```go
select {
case msg := <-ch:
    fmt.Println("Mensagem recebida:", msg)
default:
    fmt.Println("Nenhuma mensagem disponível") // ✅ Executa imediatamente se ch estiver vazio
}
```

**Padrão de Timeout:**
```go
select {
case result := <-slowOperation():
    return result
case <-time.After(5 * time.Second):
    return errors.New("operação muito lenta")
}
```

**Caso de Uso Real:**
- **Multiple Data Sources**: Primeira API que responder
- **Timeouts**: Evitar esperar para sempre
- **Graceful Shutdown**: Aguardar finalização ou timeout
- **Message Routing**: Rotear mensagens baseado na disponibilidade

## Forever Channels

Padrão para manter goroutines rodando até receber sinal de parada.

```go
package main

import (
    "fmt"
    "time"
)

func worker(done chan bool) {
    for {
        select {
        case <-done:
            fmt.Println("Worker parando...")
            return
        default:
            fmt.Println("Trabalhando...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    done := make(chan bool)

    go worker(done)

    // Deixa trabalhar por 3 segundos
    time.Sleep(3 * time.Second)

    // Sinaliza para parar
    done <- true

    // Aguarda um pouco para ver a mensagem
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Programa finalizado")
}
```

---

##  Padrões Avançados

## Worker Pools

Worker Pool é como uma **equipe de trabalhadores especializados** que processam tarefas de uma fila comum.

**Analogia Real:**
É como um call center - várias pessoas (workers) atendem ligações (jobs) de uma fila única.

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Job struct {
    ID   int
    Data string
}

type Result struct {
    Job    Job
    Output string
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        fmt.Printf("Worker %d processando job %d\n", id, job.ID)

        // Simula trabalho variável
        time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

        result := Result{
            Job:    job,
            Output: fmt.Sprintf("Processado por worker %d", id),
        }

        results <- result
    }
}

func main() {
    const numWorkers = 3
    const numJobs = 9

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)

    // Iniciando workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // Enviando jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j, Data: fmt.Sprintf("dados-%d", j)}
    }
    close(jobs)

    // Coletando resultados
    for r := 1; r <= numJobs; r++ {
        result := <-results
        fmt.Printf("Job %d finalizado: %s\n", result.Job.ID, result.Output)
    }
}
```

**Caso de Uso Real:**
- **Processamento de Imagens**: Redimensionar milhares de fotos
- **Envio de Emails**: Pool de workers enviando emails
- **Web Scraping**: Workers fazendo requests para sites diferentes
- **ETL**: Extract, Transform, Load de dados

**Vantagens:**
- ✅ Limita uso de recursos (CPU, memória, conexões)
- ✅ Balanceamento automático de carga
- ✅ Fácil de escalar (ajustar número de workers)

## Load Balancer

Distribuição inteligente de trabalho entre múltiplos workers.

```go
package main

import (
    "fmt"
    "time"
)

func server(name string, requests <-chan string) {
    for req := range requests {
        fmt.Printf("Servidor %s processando: %s\n", name, req)
        time.Sleep(time.Second) // Simula processamento
        fmt.Printf("Servidor %s finalizou: %s\n", name, req)
    }
}

func loadBalancer(requests <-chan string, servers ...chan<- string) {
    currentServer := 0
    for req := range requests {
        // Round-robin: distribui sequencialmente
        servers[currentServer] <- req
        currentServer = (currentServer + 1) % len(servers)
    }

    // Fecha todos os servers
    for _, server := range servers {
        close(server)
    }
}

func main() {
    requests := make(chan string, 10)

    // Criando 3 servidores
    server1 := make(chan string)
    server2 := make(chan string)
    server3 := make(chan string)

    go server("Alpha", server1)
    go server("Beta", server2)
    go server("Gamma", server3)

    go loadBalancer(requests, server1, server2, server3)

    // Enviando requests
    for i := 1; i <= 9; i++ {
        requests <- fmt.Sprintf("Request-%d", i)
    }
    close(requests)

    time.Sleep(5 * time.Second) // Aguarda processamento
}
```

**Estratégias de Load Balancing:**
- **Round Robin**: Sequencial (como no exemplo)
- **Least Connections**: Servidor com menos conexões
- **Weighted**: Baseado na capacidade do servidor
- **Health Check**: Apenas servidores saudáveis

## Pipeline Pattern

Pipeline processa dados em **etapas sequenciais**, onde cada etapa é uma goroutine especializada.

**Analogia Real:**
É como uma **linha de montagem de carros** - cada estação faz uma parte específica do trabalho.

```go
package main

import (
    "fmt"
    "strings"
)

// Etapa 1: Gerar números
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Etapa 2: Elevar ao quadrado
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Etapa 3: Converter para string
func toString(in <-chan int) <-chan string {
    out := make(chan string)
    go func() {
        for n := range in {
            out <- fmt.Sprintf("Número: %d", n)
        }
        close(out)
    }()
    return out
}

func main() {
    // Pipeline: números → quadrados → strings
    numbers := generator(1, 2, 3, 4, 5)
    squares := square(numbers)
    strings := toString(squares)

    // Resultado final
    for result := range strings {
        fmt.Println(result)
    }
}
```

**Vantagens do Pipeline:**
- ✅ Processamento paralelo de diferentes etapas
- ✅ Fácil de testar cada etapa individualmente
- ✅ Escalável (pode adicionar mais etapas)
- ✅ Reutilizável

**Caso de Uso Real:**
- **Processamento de Logs**: Parse → Filtro → Agregação → Storage
- **Pipeline de CI/CD**: Build → Test → Deploy
- **Streaming de Dados**: Ingest → Transform → Validate → Store

## Fan-Out / 📥 Fan-In

**Fan-Out**: Distribuir trabalho de uma fonte para múltiplos workers
**Fan-In**: Coletar resultados de múltiplos workers em um ponto

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Fan-Out: Distribui jobs para múltiplos workers
func fanOut(jobs <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)

    for i := 0; i < workers; i++ {
        output := make(chan int)
        outputs[i] = output

        go func(out chan<- int) {
            defer close(out)
            for job := range jobs {
                // Simula processamento
                time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
                out <- job * job
            }
        }(output)
    }

    return outputs
}

// Fan-In: Coleta resultados de múltiplos channels
func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup

    wg.Add(len(inputs))

    for _, input := range inputs {
        go func(in <-chan int) {
            defer wg.Done()
            for value := range in {
                output <- value
            }
        }(input)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}

func main() {
    // Criando jobs
    jobs := make(chan int, 10)
    go func() {
        for i := 1; i <= 10; i++ {
            jobs <- i
        }
        close(jobs)
    }()

    // Fan-Out: 3 workers processando
    outputs := fanOut(jobs, 3)

    // Fan-In: Coletando todos os resultados
    results := fanIn(outputs...)

    // Processando resultados
    var allResults []int
    for result := range results {
        allResults = append(allResults, result)
    }

    fmt.Printf("Resultados coletados: %v\n", allResults)
}
```

---

##  Casos de Uso Reais

## Servidores Web

```go
// Problema: Race condition em contador de visitas
var visitCount int // ❌ Não thread-safe

func handler(w http.ResponseWriter, r *http.Request) {
    visitCount++ // ❌ Race condition!
    fmt.Fprintf(w, "Visitas: %d", visitCount)
}

// Solução 1: Mutex
var (
    visitCount int
    mu sync.Mutex
)

func handlerMutex(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    visitCount++
    current := visitCount
    mu.Unlock()
    fmt.Fprintf(w, "Visitas: %d", current)
}

// Solução 2: Operações Atômicas (mais rápida)
var visitCount uint64

func handlerAtomic(w http.ResponseWriter, r *http.Request) {
    current := atomic.AddUint64(&visitCount, 1)
    fmt.Fprintf(w, "Visitas: %d", current)
}
```

## Processamento de Dados

```go
// Processar milhões de registros em paralelo
func processLargeDataset(data []Record) []Result {
    const numWorkers = 8
    jobs := make(chan Record, len(data))
    results := make(chan Result, len(data))

    // Iniciar workers
    for i := 0; i < numWorkers; i++ {
        go func() {
            for record := range jobs {
                result := processRecord(record) // Operação pesada
                results <- result
            }
        }()
    }

    // Enviar jobs
    for _, record := range data {
        jobs <- record
    }
    close(jobs)

    // Coletar resultados
    var finalResults []Result
    for i := 0; i < len(data); i++ {
        finalResults = append(finalResults, <-results)
    }

    return finalResults
}
```

## Message Queues

```go
// Simulando RabbitMQ vs Kafka com Select
func messageRouter() {
    rabbitMQ := make(chan Message)
    kafka := make(chan Message)

    go func() {
        for {
            select {
            case msg := <-rabbitMQ:
                fmt.Printf("Processando mensagem do RabbitMQ: %s\n", msg.Content)
            case msg := <-kafka:
                fmt.Printf("Processando mensagem do Kafka: %s\n", msg.Content)
            case <-time.After(30 * time.Second):
                fmt.Println("Timeout - nenhuma mensagem em 30s")
                return
            }
        }
    }()
}
```

---

##  Trade-offs e Considerações

## Performance vs Complexidade

```
SIMPLES                    COMPLEXO
   ↓                          ↓
Sequential  →  Goroutines  →  Worker Pools  →  Actor Model
   ↓                          ↓                    ↓
 Lento        Médio          Rápido             Muito Rápido
```

**Quando usar cada abordagem:**

**Sequential (sem concorrência):**
- ✅ Operações simples e rápidas
- ✅ Código legível e fácil de debugar
- ❌ Não aproveita múltiplos cores

**Goroutines simples:**
- ✅ I/O bound operations (rede, disco)
- ✅ Paralelização básica
- ❌ Difícil controlar número de goroutines

**Worker Pools:**
- ✅ CPU intensive operations
- ✅ Controle de recursos
- ✅ Escalabilidade
- ❌ Mais código para gerenciar

## Segurança vs Velocidade

| Abordagem | Velocidade | Segurança | Use Quando |
|-----------|-----------|-----------|------------|
| **Channels** | 🐌 Média | 🛡️ Alta | Comunicação complexa |
| **Mutex** | 🐌 Baixa | 🛡️ Alta | Estruturas complexas |
| **Atômicas** | ⚡ Alta | 🛡️ Média | Operações simples |
| **Lock-free** | ⚡ Muito Alta | ⚠️ Baixa | Especialistas apenas |

## Escalabilidade vs Recursos

```go
// ❌ Não escalável: Uma goroutine por tarefa
for i := 0; i < 1000000; i++ {
    go processTask(i) // 1M de goroutines!
}

// ✅ Escalável: Worker pool limitado
func scalableProcessing(tasks []Task) {
    const maxWorkers = runtime.NumCPU()
    tasksChan := make(chan Task, len(tasks))

    // Limitado pelo número de CPUs
    for i := 0; i < maxWorkers; i++ {
        go worker(tasksChan)
    }

    for _, task := range tasks {
        tasksChan <- task
    }
}
```

**Métricas para Monitorar:**
- **Goroutines ativas**: `runtime.NumGoroutine()`
- **Uso de CPU**: Deve estar próximo de 100% em tarefas CPU-intensive
- **Uso de memória**: Cuidado com memory leaks em channels
- **Latência**: Tempo de resposta das operações

**Regras Práticas:**
- 🎯 **I/O bound**: Muitas goroutines (milhares)
- 🎯 **CPU bound**: Goroutines = número de CPUs
- 🎯 **Mixed workload**: Worker pool com buffer
- 🎯 **Memory limited**: Channels com buffer pequeno

A concorrência em Go é poderosa, mas requer **análise cuidadosa** do trade-off entre performance, complexidade e recursos. Sempre meça e profile seu código! 📈

