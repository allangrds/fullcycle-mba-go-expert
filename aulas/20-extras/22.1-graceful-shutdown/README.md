# 🛑 Graceful Shutdown em Go — Guia Didático

Todo servidor, em algum momento, precisa parar. Pode ser um deploy, um `Ctrl+C` acidental, um autoscaler do Kubernetes decidindo derrubar um pod, ou o próprio sistema operacional pedindo para o processo encerrar. A pergunta que este módulo responde é: **o que acontece com as requisições que estão em andamento quando isso acontece?**

Se você é iniciante no assunto, pense assim: existe uma diferença enorme entre "desligar" um servidor e "desligar educadamente" um servidor. Este exemplo é minimalista de propósito — um único arquivo `main.go` com pouco mais de 40 linhas — mas concentra conceitos centrais de Go usados em praticamente todo serviço em produção: goroutines, canais, sinais do sistema operacional e `context` com timeout.

---

## 📑 Sumário

- [🤔 O que é Graceful Shutdown?](#-o-que-é-graceful-shutdown)
- [⚔️ Graceful Shutdown vs Alternativas](#️-graceful-shutdown-vs-alternativas)
- [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
- [▶️ Como Executar](#️-como-executar)
- [⚖️ Trade-offs](#️-trade-offs)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [❓ Perguntas de Entrevista](#-perguntas-de-entrevista)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é Graceful Shutdown?

**Analogia:** imagine um restaurante que precisa fechar às 22h. Existem duas formas de fazer isso:

```
❌ SEM SHUTDOWN GRACIOSO (o jeito bruto)
┌─────────────────────────────────────────┐
│  22h em ponto:                           │
│  - Apaga as luzes                        │
│  - Tranca a porta                        │
│  - Manda todo mundo embora, inclusive    │
│    quem estava no meio do jantar         │
│  → Clientes ficam com fome e irritados   │
└─────────────────────────────────────────┘

✅ COM SHUTDOWN GRACIOSO (o jeito educado)
┌─────────────────────────────────────────┐
│  22h em ponto:                           │
│  - Para de aceitar clientes novos        │
│    (pendura a placa "fechado" na porta)  │
│  - Deixa quem já está sentado terminar   │
│    de comer (com um prazo razoável)      │
│  - Só aí apaga as luzes e fecha          │
│  → Ninguém fica no meio do prato         │
└─────────────────────────────────────────┘
```

Tecnicamente, **graceful shutdown** (encerramento gracioso) é a técnica de, ao receber um pedido para parar, um processo:

1. Parar de aceitar **novo** trabalho (novas conexões, novas mensagens de fila, etc.);
2. Dar tempo para o trabalho **em andamento** terminar;
3. Só então encerrar de fato — liberando portas, fechando conexões com banco de dados, etc.

Em código, a diferença fica bem clara:

```go
// ❌ "Shutdown" abrupto — mata o processo na marra
// Ctrl+C sem tratamento nenhum de sinal = o SO manda SIGINT,
// o runtime do Go simplesmente derruba tudo, requisição em andamento
// é cortada no meio e o cliente recebe uma conexão resetada.
func main() {
	http.ListenAndServe(":3000", handler)
}
```

```go
// ✅ Shutdown gracioso — este é o exemplo desta pasta
server := &http.Server{Addr: ":3000"}

go func() {
	server.ListenAndServe() // roda em background
}()

<-stop // espera o pedido de parada

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
server.Shutdown(ctx) // para de aceitar novas conexões e espera as atuais terminarem
```

---

## ⚔️ Graceful Shutdown vs Alternativas

| Abordagem | O que faz | Requisições em andamento | Quando usar |
|---|---|---|---|
| **Matar o processo (`SIGKILL` / `kill -9`)** | Encerra o processo imediatamente, sem chance de reação | Perdidas na hora, sem aviso | Nunca, a não ser que o processo esteja travado e não responda a mais nada |
| **`os.Exit(0)` direto** | Encerra o processo Go instantaneamente, sem rodar `defer`s nem fechar conexões | Perdidas | Scripts curtos de linha de comando que não têm estado a proteger |
| **Graceful shutdown com timeout** (este exemplo) | Para de aceitar conexões novas, espera as atuais terminarem até um prazo limite | Têm chance de terminar, dentro do prazo | APIs HTTP, workers, qualquer serviço de longa duração |
| **Orquestração (Kubernetes `preStop` + `terminationGracePeriodSeconds`)** | O orquestrador avisa a aplicação *antes* de mandar o sinal de término, dando um tempo extra de "descanso" antes do `SIGTERM` | Melhor ainda: dá tempo até para o load balancer parar de rotear tráfego novo | Serviços rodando em produção, dentro de containers/Kubernetes |

Este exemplo implementa a terceira linha da tabela — a base que qualquer uma das abordagens mais robustas (incluindo a de Kubernetes) usa por baixo dos panos.

---

## 📚 Conceitos Fundamentais

### 1. Sinais do Sistema Operacional

**Analogia:** um sinal é como alguém batendo na sua porta. Você pode ignorar, atender rápido, ou demorar para atender — mas a batida em si é só um aviso, não uma ordem que te obriga a sair correndo.

```go
// main.go — linha 30
signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
```

O sistema operacional usa **sinais** para se comunicar com processos. Os mais comuns em shutdown:

- `SIGINT` — enviado quando você aperta `Ctrl+C` no terminal;
- `SIGTERM` — enviado por orquestradores (Docker, Kubernetes, `systemd`) pedindo para o processo terminar "com educação";
- `SIGKILL` — não pode nem ser capturado; o SO simplesmente mata o processo (por isso ele não aparece no código acima).

> 💡 **Detalhe interessante:** no código, `os.Interrupt` e `syscall.SIGINT` são, na prática, **o mesmo sinal** em sistemas Unix (`os.Interrupt` é apenas um alias multiplataforma). Registrar os dois não quebra nada, mas é redundante — um resquício comum de copiar exemplos de diferentes fontes.

### 2. Canal como Ponte entre o SO e a Aplicação

**Analogia:** um canal (`chan`) em Go é como uma caixa de correio. O carteiro (o sistema operacional) deposita uma carta (o sinal) nela, e alguém dentro de casa (a goroutine principal) fica esperando essa carta chegar para agir.

```go
// main.go — linhas 29-31
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
<-stop
```

- `make(chan os.Signal, 1)` cria um canal **com buffer de tamanho 1** — importante porque o pacote `signal` nunca bloqueia ao entregar um sinal; se o canal estivesse sem buffer (ou cheio) e ninguém estivesse lendo, o sinal seria simplesmente descartado.
- `signal.Notify(...)` diz ao runtime do Go: "quando um desses sinais chegar, em vez do comportamento padrão do SO, jogue uma mensagem nesse canal".
- `<-stop` é uma **leitura bloqueante**: a `main` para exatamente aqui e não faz mais nada até que algo apareça no canal.

### 3. Goroutine para o Servidor Não Bloquear a `main`

**Analogia:** é como colocar alguém para atender o telefone enquanto você continua de olho na porta esperando uma visita. Se você mesmo atendesse o telefone, nunca perceberia quando a visita chegasse.

```go
// main.go — linhas 22-27
go func() {
	fmt.Println("Server is running at http://localhost:3000")
	if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}()
```

`server.ListenAndServe()` é uma chamada **bloqueante** — ela só retorna quando o servidor para. Se essa linha rodasse direto na `main` (sem `go func`), o código nunca chegaria até o `signal.Notify` e o `<-stop` logo abaixo. Rodando em uma **goroutine**, o servidor escuta requisições em paralelo enquanto a `main` fica livre para esperar o sinal de parada.

### 4. `context.WithTimeout` como "Prazo de Tolerância"

**Analogia:** é o "vocês têm até 22h15 para terminar de comer" do dono do restaurante. Depois desse prazo, mesmo quem ainda não terminou vai ter que ir embora.

```go
// main.go — linhas 33-34
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

Um `context.Context` com timeout carrega um prazo (*deadline*). Quando esse prazo expira, o contexto é automaticamente "cancelado" — e qualquer função que respeite esse contexto (como `server.Shutdown`, mais abaixo) é obrigada a parar o que está fazendo. O `defer cancel()` libera os recursos internos do contexto assim que a função onde ele foi criado termina, mesmo que o timeout nunca tenha sido atingido.

### 5. `server.Shutdown(ctx)` — o Coração do Graceful Shutdown

```go
// main.go — linhas 36-39
fmt.Println("Shutting down server...")
if err := server.Shutdown(ctx); err != nil {
	log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
}
fmt.Println("Server stopped")
```

`server.Shutdown(ctx)` faz três coisas, em ordem:

1. Fecha imediatamente todos os *listeners* abertos (ninguém mais consegue **abrir uma conexão nova**);
2. Espera todas as conexões ativas terminarem sozinhas (sem interromper à força);
3. Se o `ctx` expirar antes de todas terminarem, `Shutdown` retorna um erro e as conexões restantes **são então fechadas à força**.

Ou seja: o timeout do contexto não é o tempo que o servidor "aguenta rodando" — é o **prazo máximo** que ele dá para as requisições atuais se resolverem antes de desistir delas.

### 6. `http.ErrServerClosed` — um "Erro" que na Verdade é Sucesso

```go
// main.go — linha 24
if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
```

Quando `server.Shutdown()` é chamado, o `ListenAndServe()` que estava bloqueado na goroutine finalmente retorna — mas retorna um erro: `http.ErrServerClosed`. Esse erro **não representa uma falha real**, é apenas o jeito da biblioteca padrão avisar "parei porque alguém pediu, não porque quebrei". Por isso o código verifica explicitamente esse caso antes de decidir se deve chamar `log.Fatalf`.

---

## 🗂️ Estrutura do Projeto

```
22.1-graceful-shutdown/
├── main.go       → todo o exemplo: servidor HTTP + sinais + shutdown gracioso
├── go.mod        → módulo Go (sem dependências externas, só biblioteca padrão)
└── .air.toml     → configuração do "air" (ferramenta de live-reload, opcional)
```

Vale notar que **não há nenhuma dependência de terceiros** — todo o exemplo usa apenas os pacotes padrão do Go (`net/http`, `os/signal`, `syscall`, `context`), o que reforça que graceful shutdown é um conceito suportado nativamente pela linguagem, sem precisar de framework nenhum.

---

## 🔍 Walkthrough do Código

Seguindo a ordem de execução real do programa:

```go
// 1. Cria o servidor, mas ainda não o inicia
server := &http.Server{Addr: ":3000"}

// 2. Registra o handler da rota "/", que propositalmente demora 4 segundos
//    — isso existe só para facilitar observar o shutdown "no meio" de uma requisição
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	time.Sleep(4 * time.Second)
	w.Write([]byte("Hello World"))
})

// 3. Sobe o servidor em uma goroutine, para não travar a execução principal
go func() {
	fmt.Println("Server is running at http://localhost:3000")
	if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}()

// 4. Cria o canal e diz ao runtime para nos avisar quando um sinal de parada chegar
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

// 5. A main "trava" aqui até que Ctrl+C (ou outro sinal) aconteça
<-stop

// 6. Sinal recebido! Cria um prazo de 5s para o shutdown
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 7. Pede o shutdown gracioso e espera o resultado
fmt.Println("Shutting down server...")
if err := server.Shutdown(ctx); err != nil {
	log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
}
fmt.Println("Server stopped")
```

O ponto-chave é que os passos 3 e 4-5 acontecem **em paralelo**: uma goroutine cuida das requisições HTTP, enquanto a `main` fica bloqueada esperando um sinal. Só quando o sinal chega é que o fluxo principal "acorda" e inicia a sequência de encerramento.

---

## ▶️ Como Executar

```bash
# Dentro da pasta aulas/20-extras/22.1-graceful-shutdown

# Opção 1: rodar direto
go run .

# Opção 2: com live-reload (reinicia automaticamente ao salvar o arquivo)
air
```

Para observar o graceful shutdown na prática:

1. Rode o servidor com `go run .` — você verá `Server is running at http://localhost:3000`.
2. Em outro terminal, faça uma requisição: `curl http://localhost:3000/` (ela vai demorar 4 segundos para responder, de propósito).
3. **Enquanto a requisição ainda está pendente**, volte ao terminal do servidor e aperte `Ctrl+C`.
4. Observe a saída: aparecerá `Shutting down server...`, mas o processo **não morre imediatamente** — ele espera a requisição do passo 2 terminar (e você verá o `curl` receber "Hello World" normalmente) antes de imprimir `Server stopped` e encerrar de fato.
5. Experimente repetir o teste, mas trocando o `4 * time.Second` do handler para algo maior que os `5 * time.Second` do timeout do contexto (ex.: `10 * time.Second`) — dessa vez, o `Shutdown` vai retornar um erro porque a requisição não terminou a tempo.

---

## ⚖️ Trade-offs

**✅ Vantagens**

- Nenhuma requisição em andamento é interrompida bruscamente (sem conexões resetadas no meio do processo).
- Evita respostas de erro (ex.: `502 Bad Gateway`) para usuários que já estavam sendo atendidos no momento do deploy.
- Se comunica bem com orquestradores (Kubernetes, Docker Swarm, `systemd`), que esperam que os processos tratem `SIGTERM` antes de matar com `SIGKILL`.
- Implementado inteiramente com a biblioteca padrão — nenhuma dependência externa necessária.

**❌ Desvantagens**

- Um timeout mal calibrado pode ainda assim cortar requisições legítimas que só demoram um pouco mais (ex.: um upload grande).
- Adiciona complexidade ao código de inicialização/encerramento que um serviço "descartável" (como um script de linha de comando de vida curta) não precisa.
- Graceful shutdown por si só **não resolve tudo**: em produção, o load balancer/proxy também precisa saber que a instância está saindo (via *readiness probe*) para parar de mandar tráfego novo para ela — senão, novas conexões continuam chegando durante a janela de shutdown.
- Se o processo tiver outras tarefas em segundo plano além do servidor HTTP (workers, conexões de banco, filas), elas precisam de tratamento próprio — o `server.Shutdown` só cuida do servidor HTTP.

---

## 🎯 Casos de Uso Ideais

**Use graceful shutdown quando:**
- Você está construindo uma API HTTP ou gRPC que roda continuamente;
- O serviço processa mensagens de fila (Kafka, RabbitMQ, SQS) e não pode simplesmente descartar uma mensagem em processamento;
- A aplicação roda em containers/Kubernetes, onde deploys e reinícios são frequentes;
- Existem recursos que precisam ser fechados corretamente (conexões de banco, arquivos abertos, locks distribuídos).

**Pode ser dispensável quando:**
- É um script CLI de vida curta que roda, faz seu trabalho e termina sozinho;
- É uma ferramenta de uso local, sem múltiplos usuários dependendo de disponibilidade;
- O estado da aplicação é totalmente descartável e não há requisições concorrentes para proteger.

---

## ❓ Perguntas de Entrevista

**O que é graceful shutdown e por que ele é importante em produção?**
É a técnica de, ao receber um pedido de encerramento, parar de aceitar trabalho novo e dar tempo para o trabalho em andamento terminar antes de finalizar o processo. É importante porque evita interromper requisições de usuários reais no meio do caminho, algo que se torna crítico em ambientes com deploys frequentes (cada deploy derruba e sobe processos).

**Qual a diferença entre `SIGINT`, `SIGTERM` e `SIGKILL`?**
`SIGINT` é enviado normalmente por uma interação manual (`Ctrl+C`). `SIGTERM` é o sinal "educado" de pedido de término, usado por orquestradores e ferramentas de deploy — ele pode ser capturado e tratado pela aplicação, como faz este exemplo. `SIGKILL` não pode ser capturado nem ignorado: o sistema operacional mata o processo imediatamente, sem dar chance de nenhum código rodar — por isso ele deve ser o último recurso, usado apenas quando o processo não responde a mais nada.

**Por que o servidor HTTP precisa rodar em uma goroutine separada?**
Porque `server.ListenAndServe()` é uma chamada bloqueante: ela só retorna quando o servidor para de rodar. Se ela rodasse diretamente na função `main`, o programa nunca chegaria ao código que registra e espera por sinais do sistema operacional. Rodá-la em uma goroutine permite que o servidor atenda requisições em paralelo enquanto a `main` fica livre para monitorar o pedido de encerramento.

**O que acontece se o timeout do `context.WithTimeout` expirar antes de todas as requisições terminarem?**
`server.Shutdown(ctx)` retorna um erro (o contexto expirado) e as conexões que ainda estavam ativas são fechadas à força — ou seja, o timeout é o prazo máximo de tolerância, não uma garantia de que tudo vai terminar bem.

**Por que o código verifica `err != http.ErrServerClosed` em vez de tratar qualquer erro como falha?**
Porque quando `Shutdown()` é chamado, `ListenAndServe()` retorna exatamente esse erro como forma de sinalizar "parei porque pediram, não porque quebrei". Tratá-lo como uma falha real (chamando `log.Fatalf` sem essa verificação) faria o programa reportar erro mesmo em um encerramento bem-sucedido e esperado.

**Como esse conceito se relaciona com `readinessProbe` e `preStop hook` no Kubernetes?**
No Kubernetes, quando um pod é marcado para encerramento, o `preStop hook` roda antes do `SIGTERM` ser enviado, dando um tempo extra para o *service mesh*/load balancer parar de rotear tráfego novo para aquele pod. Depois disso, o `SIGTERM` chega até a aplicação, que deve reagir exatamente como este exemplo faz — sem isso, mesmo um graceful shutdown correto no código pode receber tráfego novo durante a janela de encerramento, pois o balanceador ainda não sabia que o pod estava saindo.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **Goroutine** | Uma unidade leve de execução concorrente gerenciada pelo runtime do Go, iniciada com a palavra-chave `go`. |
| **Canal (`chan`)** | Estrutura usada para comunicação e sincronização entre goroutines; pode ter ou não um buffer interno. |
| **Sinal (signal)** | Notificação assíncrona enviada pelo sistema operacional a um processo (ex.: `SIGINT`, `SIGTERM`, `SIGKILL`). |
| **`context.Context`** | Estrutura da biblioteca padrão usada para propagar prazos, cancelamento e valores entre chamadas de função. |
| **Deadline / Timeout** | Prazo máximo até o qual uma operação deve ser concluída antes de ser considerada expirada. |
| **Graceful shutdown** | Encerramento de um processo que primeiro para de aceitar trabalho novo e espera o trabalho em andamento terminar. |
| **`http.ErrServerClosed`** | Erro sentinela retornado por `ListenAndServe` quando o servidor foi parado intencionalmente via `Shutdown`/`Close`. |
| **Listener** | Objeto que "escuta" uma porta de rede aguardando novas conexões. |
| **Blocking call (chamada bloqueante)** | Uma chamada de função que não retorna até que sua tarefa esteja completa, "travando" a goroutine que a executa. |

---

## 🚀 Próximos Passos

- [ ] **Imediato:** altere o `time.Sleep` do handler e o timeout do `context.WithTimeout` para valores diferentes e observe como isso muda o comportamento do shutdown (tanto o caso de sucesso quanto o de timeout).
- [ ] **Imediato:** adicione um segundo endpoint (`/health`, por exemplo) e teste o shutdown gracioso com múltiplas requisições simultâneas em rotas diferentes.
- [ ] **Intermediário:** simule uma conexão de banco de dados (mesmo que fake) que precise ser fechada no shutdown, usando um `sync.WaitGroup` para garantir que ela só feche depois que o servidor HTTP parar.
- [ ] **Intermediário:** registre um log estruturado (ex.: com `slog`, da própria biblioteca padrão) em vez de `fmt.Println`, para ficar mais próximo de um cenário de produção.
- [ ] **Avançado:** use `golang.org/x/sync/errgroup` para orquestrar múltiplos processos concorrentes (ex.: um servidor HTTP e um consumidor de fila) que precisam encerrar graciosamente juntos.
- [ ] **Avançado:** empacote este exemplo em um container Docker e rode-o no Kubernetes, configurando `terminationGracePeriodSeconds` e um `preStop hook`, para observar o graceful shutdown funcionando em conjunto com o orquestrador.
