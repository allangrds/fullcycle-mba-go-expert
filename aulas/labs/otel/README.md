# 🔭 OpenTelemetry em Go — Guia Didático Completo

Este laboratório mostra, na prática, como instrumentar aplicações Go com **OpenTelemetry** para gerar **tracing distribuído**, exportar esses dados através de um **OpenTelemetry Collector** e visualizá-los no **Jaeger**, além de expor **métricas** para **Prometheus** e **Grafana**.

Se você nunca trabalhou com observabilidade antes, não se preocupe: cada conceito aqui é explicado primeiro com uma analogia do dia a dia, depois com a definição técnica, e por fim com o trecho de código real deste projeto que o implementa. Ao final deste README você vai entender o que é um trace, o que é um span, como o "contexto" de uma requisição viaja entre serviços diferentes, para que serve o Collector, e por que tudo isso importa em sistemas de microsserviços.

Este lab tem **dois demos**:

- **Raiz (`aulas/labs/otel/main.go`)** — o exemplo mais simples possível: um único serviço criando spans manualmente.
- **`comunicacao-ms/`** — um exemplo avançado com **3 microsserviços encadeados** (`goapp` → `goapp2` → `goapp3`), mostrando como o trace de uma requisição atravessa vários processos diferentes e ainda assim aparece como **um único trace** no Jaeger.

---

## 📑 Sumário

1. [🤔 O que é Observabilidade e por que ela importa?](#-o-que-é-observabilidade-e-por-que-ela-importa)
2. [🤔 O que é o OpenTelemetry?](#-o-que-é-o-opentelemetry)
3. [🗂️ Estrutura deste projeto](#️-estrutura-deste-projeto)
4. [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
5. [⚙️ Como funciona o fluxo (passo a passo)](#️-como-funciona-o-fluxo-passo-a-passo)
6. [🔍 Demo simples vs demo avançado](#-demo-simples-vs-demo-avançado)
7. [✅ Boas práticas presentes no projeto](#-boas-práticas-presentes-no-projeto)
8. [⚠️ Problemas e pegadinhas encontrados no próprio código](#️-problemas-e-pegadinhas-encontrados-no-próprio-código)
9. [⚖️ Tradeoffs importantes em observabilidade](#️-tradeoffs-importantes-em-observabilidade)
10. [🔧 Como rodar o projeto](#-como-rodar-o-projeto)
11. [📖 Glossário](#-glossário)
12. [💼 Perguntas de Entrevista Respondidas](#-perguntas-de-entrevista-respondidas)
13. [🚀 Próximos passos / exercícios sugeridos](#-próximos-passos--exercícios-sugeridos)

---

## 🤔 O que é Observabilidade e por que ela importa?

Imagine o painel de instrumentos de um carro: velocímetro, conta-giros, luz de óleo, temperatura do motor. Você não precisa abrir o motor para saber se algo está errado — os instrumentos **contam a história** do que está acontecendo lá dentro. Observabilidade é exatamente isso, só que para sistemas de software: um conjunto de sinais que permitem entender o que está acontecendo **dentro** de um sistema só olhando o que sai dele "de fora" (logs, métricas, traces).

Em um monólito, quando algo dá errado, basta olhar um log. Mas em um sistema de microsserviços — como o `comunicacao-ms/` deste lab, onde uma requisição passa por `goapp → goapp2 → goapp3` — uma única ação do usuário pode gerar dezenas de chamadas entre serviços diferentes, cada um com seu próprio log. Sem uma forma de **conectar** esses logs, descobrir "por que essa requisição demorou 3 segundos?" ou "qual serviço causou o erro?" vira uma caça ao tesouro.

Observabilidade moderna se apoia em **três pilares**:

| Pilar | O que é | Pergunta que responde | Exemplo neste lab |
|---|---|---|---|
| **Logs** | Registros de eventos discretos, com timestamp | "O que aconteceu?" | `log.Println("Starting server on port...")` |
| **Métricas** | Números agregados ao longo do tempo (contadores, histogramas) | "Quão saudável está o sistema?" | Endpoint `/metrics` exposto via `promhttp.Handler()` |
| **Traces** | O caminho completo de uma requisição através de vários serviços | "Por onde essa requisição passou e onde ela demorou?" | Os spans criados em `internal/web/server.go` |

Este lab foca principalmente em **traces** (o pilar mais difícil de implementar sozinho) e toca também em métricas.

---

## 🤔 O que é o OpenTelemetry?

Antes do OpenTelemetry, cada ferramenta de observabilidade (Jaeger, Zipkin, Datadog, New Relic...) tinha seu próprio SDK de instrumentação. Se você instrumentasse seu código com o SDK do Jaeger e um dia quisesse trocar para outra ferramenta, teria que reescrever toda a instrumentação. Isso é **vendor lock-in**.

O **OpenTelemetry (OTel)** é um padrão aberto, mantido pela CNCF, que resolve esse problema separando duas coisas:

1. **Instrumentação** — o código que cria spans, métricas e propaga contexto (é o que a aplicação Go faz usando o pacote `go.opentelemetry.io/otel`).
2. **Backend de armazenamento/visualização** — Jaeger, Zipkin, Prometheus, Grafana, Datadog, etc.

> ⚠️ **Importante**: o OpenTelemetry **não é** um backend. Ele não armazena nem visualiza nada — ele só coleta e exporta os dados no formato padrão **OTLP** (OpenTelemetry Protocol). Quem armazena e mostra os dados, neste lab, é o Jaeger (traces) e o Prometheus/Grafana (métricas).

---

## 🗂️ Estrutura deste projeto

```
aulas/labs/otel/
├── main.go                          # Demo 1: serviço único, cria spans manualmente
├── docker-compose.yaml              # jaeger, zipkin, prometheus, otel-collector, goapp
├── .docker/
│   ├── prometheus.yaml              # Prometheus faz scrape do otel-collector
│   └── otel-collector-config.yaml   # pipeline só de TRACES -> jaeger + logging
│
└── comunicacao-ms/                  # Demo 2: 3 microsserviços encadeados
    ├── cmd/microservice/main.go     # initProvider() -> configura o SDK do OTel
    ├── internal/web/
    │   ├── server.go                # Extract (entrada) + Inject (saída) do trace context
    │   └── template/index.html      # renderiza o HTML recursivamente (goapp embute goapp2 embute goapp3)
    ├── Dockerfile
    ├── docker-compose.yaml          # goapp(8080) -> goapp2(8181) -> goapp3(8282) + grafana(3001)
    └── .docker/
        ├── prometheus.yaml          # scrape do goapp:8080 + otel-collector
        └── otel-collector-config.yaml  # pipelines de TRACES **e** METRICS
```

A pasta raiz é o "hello world" do OpenTelemetry: um processo Go criando spans manualmente, só para você ver a mecânica básica. A pasta `comunicacao-ms/` é o exemplo "de verdade": três serviços HTTP independentes, cada um rodando em seu próprio container, e o desafio real de observabilidade — fazer o trace de uma única requisição do usuário atravessar os três e aparecer como **uma única linha do tempo** no Jaeger.

---

## 📚 Conceitos Fundamentais

### 1. Trace e Span

Um **Trace** representa a jornada completa de uma requisição através do sistema — do início ao fim, mesmo que ela passe por vários serviços. Um **Span** é uma unidade de trabalho dentro desse trace: por exemplo, "processar a requisição no goapp", "chamar o goapp2", "consultar o banco de dados". Um trace é literalmente uma árvore de spans (um span pai pode ter vários spans filhos).

Veja como `comunicacao-ms/internal/web/server.go` cria dois spans aninhados dentro de uma mesma requisição:

```go
// internal/web/server.go
ctx, spanInicial := h.TemplateData.OTELTracer.Start(ctx, "SPAN_INICIAL"+h.TemplateData.RequestNameOTEL)
time.Sleep(time.Second)
spanInicial.End()

ctx, span := h.TemplateData.OTELTracer.Start(ctx, "Chama externa"+h.TemplateData.RequestNameOTEL)
defer span.End()
```

Cada `tracer.Start(ctx, "nome")` cria um novo span **filho** do span presente no `ctx` atual (se houver um) e devolve um novo `ctx` que carrega esse span — por isso o `ctx` é sempre reatribuído. Um span guarda, entre outras coisas:

- **Trace ID** — identificador único do trace inteiro (compartilhado por todos os spans da mesma requisição).
- **Span ID** — identificador único daquele span específico.
- **Parent Span ID** — aponta para o span pai (permite montar a árvore).
- **Nome, timestamps de início/fim, atributos e eventos**.

### 2. TracerProvider e Resource

O `TracerProvider` é a "fábrica" central que cria tracers e sabe para onde exportar os spans. O `Resource` identifica **quem** está gerando aquele dado (o nome do serviço, por exemplo), o que é essencial quando você tem vários serviços (como `goapp`, `goapp2`, `goapp3`) e precisa saber qual span veio de qual processo:

```go
// main.go (raiz)
res, err := resource.New(ctx,
    resource.WithAttributes(
        semconv.ServiceName("mysrv1"),
    ),
)
```

Em `comunicacao-ms/`, o nome do serviço vem de uma variável de ambiente (`OTEL_SERVICE_NAME`), o que permite reaproveitar o mesmo binário para `goapp`, `goapp2` e `goapp3`, cada um se identificando com um nome diferente.

### 3. Sampler

Nem sempre você quer gravar **todos** os traces (imagine um sistema com 1 milhão de requisições por segundo — armazenar 100% dos traces seria caríssimo). O `Sampler` decide quais traces são efetivamente gravados. Este lab usa a opção mais simples:

```go
tracerProvider := sdktrace.NewTracerProvider(
    sdktrace.WithSampler(sdktrace.AlwaysSample()),
    ...
)
```

`AlwaysSample()` grava 100% dos traces — ótimo para aprender e depurar, mas [inviável em produção de alto tráfego](#️-tradeoffs-importantes-em-observabilidade).

### 4. SpanProcessor (Batch vs Simple)

O `SpanProcessor` decide **quando** os spans são enviados ao exportador. O projeto usa o `BatchSpanProcessor`:

```go
bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
```

Ele acumula spans em memória e os envia em lotes periodicamente, em vez de fazer uma chamada de rede por span. A alternativa seria o `SimpleSpanProcessor`, que exporta cada span imediatamente (mais simples, mas gera muito mais chamadas de rede — normalmente só usado em testes/debug).

### 5. Exporter (OTLP/gRPC)

O `Exporter` é responsável por **serializar e enviar** os spans para fora do processo. Aqui, os spans são enviados via gRPC no protocolo padrão **OTLP** para o OpenTelemetry Collector:

```go
conn, err := grpc.DialContext(ctx, "otel-collector:4317",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithBlock(),
)
traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
```

### 6. Propagator e Context Propagation — o conceito mais importante

Esta é a peça que faz a mágica de **conectar** spans de processos diferentes em um único trace. Cada processo Go tem seu próprio `ctx` local — quando `goapp` faz uma requisição HTTP para `goapp2`, o `ctx` do Go **não atravessa a rede sozinho**. É preciso serializá-lo em algum lugar (os headers HTTP) e desserializá-lo do outro lado. É para isso que servem `Inject` e `Extract`:

```go
// internal/web/server.go — RECEBE uma requisição
carrier := propagation.HeaderCarrier(r.Header)
ctx := r.Context()
ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)   // lê o trace context dos headers recebidos

// ... cria os spans locais usando esse ctx ...

// ANTES de chamar o próximo serviço:
otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))  // escreve o trace context nos headers de saída
resp, err := http.DefaultClient.Do(req)
```

O propagador usado é o padrão **W3C Trace Context** (`propagation.TraceContext{}`), configurado uma vez globalmente:

```go
otel.SetTextMapPropagator(propagation.TraceContext{})
```

Esse padrão injeta um header HTTP chamado `traceparent` (algo como `00-<trace-id>-<span-id>-01`) — é literalmente esse header que viaja de `goapp` para `goapp2` para `goapp3` e permite que o Jaeger monte a árvore completa do trace, mesmo que os três sejam processos totalmente independentes.

> No demo simples da raiz (`main.go`), o propagador é configurado mas **nunca usado** (`Extract`/`Inject`), porque há apenas um processo — não existe "rede" para atravessar. É só em `comunicacao-ms/` que essa peça entra em ação de verdade.

### 7. OpenTelemetry Collector

O Collector é um processo intermediário — um "hub" — entre suas aplicações e os backends de observabilidade (Jaeger, Prometheus, etc). Ele é configurado como um **pipeline** de três estágios:

```
receivers  →  processors  →  exporters
(recebe)      (transforma)     (envia)
```

Compare os dois arquivos de configuração deste lab:

**Raiz — só traces** (`.docker/otel-collector-config.yaml`):
```yaml
receivers:
  otlp:
    protocols:
      grpc:
exporters:
  otlp:
    endpoint: jaeger-all-in-one:4317
    tls:
      insecure: true
  logging:
processors:
  batch:
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp, logging]
```

**`comunicacao-ms/` — traces E metrics** (`.docker/otel-collector-config.yaml`):
```yaml
exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
  logging:
  otlp:
    endpoint: jaeger-all-in-one:4317
    tls:
      insecure: true
processors:
  batch:
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, otlp]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, prometheus]
```

Note que o segundo arquivo define **dois pipelines** — um para `traces` (que exporta pro Jaeger) e outro para `metrics` (que expõe um endpoint Prometheus na porta `8889`, de onde o Prometheus faz o scrape). Isso mostra que um único Collector pode processar múltiplos tipos de sinal (traces, métricas e até logs) ao mesmo tempo, cada um com seu próprio pipeline.

---

## ⚙️ Como funciona o fluxo (passo a passo)

Veja o que acontece quando você acessa `http://localhost:8080` (o `goapp` de `comunicacao-ms/`) no navegador:

1. **Navegador → `goapp` (porta 8080)**: chega uma requisição HTTP sem nenhum header `traceparent` (é a primeira parada).
2. **`goapp` faz `Extract`**: como não há header de trace, um **novo trace** é iniciado.
3. **`goapp` cria 2 spans**: `SPAN_INICIAL...` (simula processamento com `time.Sleep(1s)`) e `Chama externa...` (span que vai envolver a chamada ao próximo serviço).
4. **`goapp` faz `Inject`** no header do request que vai enviar para `goapp2`, escrevendo o `traceparent` com o trace ID gerado no passo 2 e o span ID do span atual como "pai".
5. **`goapp` → `goapp2` (porta 8181)**: a requisição HTTP chega já com o header `traceparent` preenchido.
6. **`goapp2` faz `Extract`**: em vez de criar um trace novo, ele **reconhece** o trace ID recebido e cria seus próprios spans como **filhos** do span de `goapp`.
7. Os passos 3–6 se repetem entre `goapp2` e `goapp3` (porta 8282).
8. **`goapp3` responde** com seu HTML; a resposta sobe a cadeia (`goapp3` → `goapp2` → `goapp`), e cada serviço embute o HTML do próximo dentro do seu próprio template — por isso, ao abrir `localhost:8080` no navegador, você vê visualmente os três serviços "aninhados" na página.
9. Em paralelo, cada serviço já enviou seus spans (via `BatchSpanProcessor` → OTLP/gRPC) para o **otel-collector**, que os encaminha para o **Jaeger**.
10. No Jaeger UI (`localhost:16686`), ao buscar pelo trace, você vê **uma única linha do tempo** com os spans de `goapp`, `goapp2` e `goapp3` organizados hierarquicamente — mesmo sendo três processos, containers e binários completamente diferentes.

---

## 🔍 Demo simples vs demo avançado

| | Raiz (`main.go`) | `comunicacao-ms/` |
|---|---|---|
| Nº de serviços | 1 | 3 (`goapp`, `goapp2`, `goapp3`) |
| Comunicação entre serviços | Não há | HTTP, com propagação de contexto |
| Usa `Extract`/`Inject`? | Não (configurado mas não usado) | Sim — é o ponto central do exemplo |
| Métricas expostas (`/metrics`) | Não | Sim, via `promhttp.Handler()` |
| Pipeline do Collector | Só `traces` | `traces` + `metrics` |
| Grafana | Não incluso | Incluso (`localhost:3001`) |
| Configuração | Hardcoded no código (`"mysrv1"`) | Via variáveis de ambiente (Viper) |
| Como roda | Manual, dentro do container `goapp` | `docker-compose up` builda e sobe tudo |
| Objetivo didático | Ver a mecânica básica do SDK | Ver tracing distribuído "de verdade" |

---

## ✅ Boas práticas presentes no projeto

- **Configuração via variáveis de ambiente** (`comunicacao-ms`): usar `viper.AutomaticEnv()` para ler `OTEL_SERVICE_NAME`, `OTEL_EXPORTER_OTLP_ENDPOINT`, etc., em vez de hardcoded, permite reaproveitar o mesmo binário/imagem Docker para os três serviços (`goapp`, `goapp2`, `goapp3`) apenas variando o `docker-compose.yaml`.
- **`initProvider` como função reutilizável e testável** (`cmd/microservice/main.go`): toda a inicialização do SDK do OTel fica isolada em uma função que devolve `(shutdown func(context.Context) error, err error)`, em vez de espalhada pelo `main`.
- **Shutdown gracioso do TracerProvider**:
  ```go
  defer func() {
      if err := shutdown(ctx); err != nil {
          log.Fatal("failed to shutdown TracerProvider: %w", err)
      }
  }()
  ```
  Isso garante que spans pendentes no buffer do `BatchSpanProcessor` sejam enviados antes do processo encerrar — sem isso, você pode perder os últimos segundos de dados a cada deploy/reinício.
- **Uso do `BatchSpanProcessor`** em vez de `SimpleSpanProcessor`, reduzindo a quantidade de chamadas de rede ao Collector (ver [tradeoffs](#️-tradeoffs-importantes-em-observabilidade)).

## ⚠️ Problemas e pegadinhas encontrados no próprio código

**1. Erros ignorados silenciosamente no `main.go` da raiz**

```go
// ❌ main.go (raiz) — o erro é formatado, mas nunca usado
res, err := resource.New(ctx, ...)
if err != nil {
    fmt.Errorf("failed to create resource: %w", err) // não retorna, não loga, não faz nada!
}
```

`fmt.Errorf` apenas **cria** um valor de erro — ele não imprime nada e não interrompe a execução. Se `resource.New` falhar, o programa continua rodando silenciosamente com um `res` inválido. O jeito certo é o que `comunicacao-ms/cmd/microservice/main.go` já faz:

```go
// ✅ comunicacao-ms/cmd/microservice/main.go
res, err := resource.New(ctx, ...)
if err != nil {
    return nil, fmt.Errorf("failed to create resource: %w", err) // propaga o erro de verdade
}
```

**2. `zipkin-all-in-one` definido mas nunca conectado**

O `docker-compose.yaml` da raiz sobe um container do Zipkin (`localhost:9411`), mas o `otel-collector-config.yaml` só tem um exporter `otlp` apontando para o Jaeger — nenhum dado chega ao Zipkin. Ele fica ocioso. Bom exercício: adicionar um exporter `zipkin` ao Collector e comparar a mesma trace nas duas UIs.

**3. Prometheus só faz scrape de `goapp`, não de `goapp2`/`goapp3`**

```yaml
# comunicacao-ms/.docker/prometheus.yaml
- job_name: 'goapp'
  static_configs:
    - targets: ['goapp:8080']   # falta goapp2:8181 e goapp3:8282!
```

Mesmo os três serviços expondo `/metrics` via `promhttp.Handler()`, só as métricas do `goapp` chegam ao Prometheus/Grafana. É um exercício sugerido corrigir isso (veja [Próximos passos](#-próximos-passos--exercícios-sugeridos)).

---

## ⚖️ Tradeoffs importantes em observabilidade

| Decisão | Opção A | Opção B | Quando escolher cada uma |
|---|---|---|---|
| **Sampling** | `AlwaysSample()` (100%) | Sampling probabilístico (ex: 10%) | 100% é ótimo em dev/staging e sistemas de baixo tráfego; em produção de alto volume, gravar tudo é caro (armazenamento) e pode sobrecarregar o Collector — usa-se sampling parcial ou baseado em regras (ex: sempre gravar erros, amostrar o resto). |
| **SpanProcessor** | `SimpleSpanProcessor` | `BatchSpanProcessor` | Simple exporta cada span na hora — útil só para debug local, pois cada span vira uma chamada de rede. Batch acumula e envia em lote, reduzindo I/O, mas com uma pequena latência entre o span terminar e ele aparecer no backend. |
| **Exportar direto vs via Collector** | App → Jaeger diretamente | App → Collector → Jaeger | Exportar direto é mais simples (menos peças), mas acopla a aplicação ao backend específico e multiplica configuração se houver múltiplos backends. O Collector centraliza essa configuração, permite trocar de backend sem tocar no código da aplicação, e pode fazer transformação/enriquecimento de dados no caminho. |
| **Instrumentação manual vs automática** | `tracer.Start()`/`span.End()` explícitos (como neste lab) | Bibliotecas de auto-instrumentação (ex: middlewares para `net/http`, `database/sql`) | Manual dá controle fino sobre nomes e granularidade dos spans, mas exige disciplina do time. Automática cobre rapidamente casos comuns (HTTP, banco de dados) com pouco código, mas é menos flexível para lógica de negócio específica. |
| **Overhead de performance vs visibilidade** | Instrumentar tudo | Instrumentar só pontos críticos | Mais spans = mais visibilidade, mas também mais overhead de CPU/memória e mais volume de dados para armazenar. Bom senso: instrumentar fronteiras de serviço, chamadas de rede/banco, e operações custosas — não cada função interna. |

---

## 🔧 Como rodar o projeto

### Demo simples (raiz)

```bash
cd aulas/labs/otel
docker-compose up -d
```

O serviço `goapp` sobe com uma imagem `golang:latest` genérica e fica ocioso (`tail -f /dev/null`) — é necessário entrar no container e rodar o programa manualmente:

```bash
docker exec -it goapp bash
cd /go/src/app
go run main.go
```

Depois, acesse:
- **Jaeger UI**: http://localhost:16686 (procure pelo serviço `mysrv1`)
- **Prometheus**: http://localhost:9090
- **Zipkin UI**: http://localhost:9411 (sobe, mas não recebe dados — ver [pegadinhas](#️-problemas-e-pegadinhas-encontrados-no-próprio-código))

### Demo avançado (`comunicacao-ms/`)

```bash
cd aulas/labs/otel/comunicacao-ms
docker-compose up -d --build
```

Isso builda a imagem (via `Dockerfile`) e sobe os três serviços (`goapp`, `goapp2`, `goapp3`) mais `jaeger`, `prometheus`, `grafana` e `otel-collector`. Depois, acesse:

- **Aplicação**: http://localhost:8080 (dispara a cadeia completa `goapp → goapp2 → goapp3`)
- **Jaeger UI**: http://localhost:16686 (busque pelo serviço `microservice-demo` para ver o trace completo)
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001 (usuário/senha padrão `admin`/`admin` na primeira vez — adicione o Prometheus como data source em `http://prometheus:9090`)

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **Observabilidade** | Capacidade de entender o estado interno de um sistema a partir dos dados que ele expõe externamente (logs, métricas, traces). |
| **Trace** | O caminho completo de uma requisição através de um ou mais serviços; um conjunto de spans relacionados. |
| **Span** | Uma unidade de trabalho dentro de um trace (uma operação, com início e fim). |
| **Trace ID** | Identificador único de um trace inteiro, compartilhado por todos os spans que o compõem. |
| **Span ID** | Identificador único de um span específico dentro de um trace. |
| **Parent Span** | O span "pai" de outro span, formando a árvore hierárquica do trace. |
| **Context Propagation** | O mecanismo de serializar o contexto de um trace (trace ID, span ID) para que ele atravesse a rede entre serviços diferentes. |
| **`traceparent`** | Header HTTP padronizado pelo W3C Trace Context que carrega o trace ID e span ID entre serviços. |
| **Propagator** | Componente do SDK responsável por fazer `Inject` (escrever) e `Extract` (ler) o contexto de trace em/de um carrier (ex: headers HTTP). |
| **Resource** | Metadados que identificam a origem dos dados de telemetria (ex: nome do serviço). |
| **TracerProvider** | Componente central do SDK que cria `Tracer`s e sabe configurar sampler, processor e exporter. |
| **Sampler** | Decide quais traces são efetivamente gravados/exportados. |
| **SpanProcessor** | Decide como e quando os spans terminados são enviados ao exportador (`Simple` = imediato, `Batch` = em lotes). |
| **Exporter** | Serializa e envia os dados de telemetria para fora do processo, em um protocolo específico (ex: OTLP). |
| **OTLP** | OpenTelemetry Protocol — o formato/protocolo padrão de transporte de dados de telemetria (traces, métricas, logs). |
| **OpenTelemetry Collector** | Processo intermediário que recebe, processa (transforma/agrega) e reenvia dados de telemetria para um ou mais backends. |
| **Pipeline (no Collector)** | Uma cadeia `receivers → processors → exporters` configurada para um tipo de sinal (traces, metrics ou logs). |
| **Baggage** | Mecanismo do OpenTelemetry para propagar pares chave-valor arbitrários junto com o contexto do trace (não usado ativamente neste lab, mas comentado no `main.go`). |
| **Vendor lock-in** | Dependência excessiva de um fornecedor/ferramenta específica, dificultando a troca no futuro — problema que o OpenTelemetry busca evitar. |

---

## 💼 Perguntas de Entrevista Respondidas

**1. O que é distributed tracing e por que ele é necessário em arquiteturas de microsserviços?**
Distributed tracing é a técnica de rastrear o caminho completo de uma requisição conforme ela atravessa múltiplos serviços, registrando quanto tempo cada etapa levou. Em um monólito, um log já mostra a sequência de execução; em microsserviços, uma requisição do usuário pode acionar dezenas de chamadas entre processos diferentes, e sem tracing distribuído fica praticamente impossível saber qual serviço causou uma lentidão ou um erro — você teria que correlacionar logs de vários sistemas manualmente.

**2. Qual a diferença entre logs, métricas e traces?**
Logs são eventos discretos e detalhados ("aconteceu X às 14:32:01"); métricas são números agregados ao longo do tempo, ótimos para dashboards e alertas ("latência média = 200ms nos últimos 5 min"); traces mostram o caminho e a duração de uma requisição específica através de vários serviços. Eles se complementam: uma métrica pode indicar que algo está lento, um trace mostra *onde* está lento, e um log mostra o *porquê* exato daquele ponto.

**3. O que é um Span e o que ele contém?**
Um Span representa uma unidade de trabalho dentro de um trace — por exemplo, uma chamada HTTP ou uma consulta ao banco. Ele contém, no mínimo: um Trace ID (compartilhado com todo o trace), um Span ID próprio, uma referência ao Span pai (se houver), nome da operação, timestamp de início e fim, e opcionalmente atributos (metadados chave-valor) e eventos (marcos dentro do span).

**4. Como o contexto de um trace é propagado entre serviços diferentes?**
Através de **context propagation**: o serviço que inicia (ou está no meio de) uma chamada usa um `Propagator` para fazer `Inject` do trace ID e span ID atual em um "carrier" — normalmente os headers HTTP, usando o padrão `traceparent` do W3C Trace Context. O serviço que recebe a requisição faz `Extract` desses headers para reconstituir o contexto e continuar o mesmo trace, em vez de iniciar um novo. É exatamente esse mecanismo que está implementado em `comunicacao-ms/internal/web/server.go`.

**5. Para que serve o OpenTelemetry Collector? Por que não exportar direto para o Jaeger?**
O Collector centraliza o recebimento, processamento (ex: batching, filtragem, enriquecimento) e reenvio de dados de telemetria para um ou mais backends. Sem ele, cada aplicação precisaria saber o endereço e protocolo específico de cada backend (Jaeger, Prometheus, etc.), e trocar de backend exigiria alterar e reimplantar todas as aplicações. Com o Collector no meio, as aplicações só falam OTLP com ele, e a troca/adição de backends é só uma mudança de configuração do Collector.

**6. O que é sampling e por que não se usa 100% em produção?**
Sampling é a decisão de quais traces são efetivamente gravados e exportados. Gravar 100% dos traces (`AlwaysSample()`) é viável em desenvolvimento, mas em produção de alto tráfego gera um volume enorme de dados, aumentando custo de armazenamento e carga no Collector/backend. Por isso, em produção normalmente se usa sampling probabilístico (ex: 5-10% das requisições) ou baseado em regras (ex: sempre gravar requisições com erro, mas amostrar as bem-sucedidas).

**7. Qual a diferença entre `SimpleSpanProcessor` e `BatchSpanProcessor`?**
`SimpleSpanProcessor` exporta cada span assim que ele termina — simples, mas gera uma chamada de rede por span, o que é ineficiente em produção. `BatchSpanProcessor` acumula spans em um buffer e os envia periodicamente em lotes, reduzindo drasticamente o número de chamadas de rede, ao custo de um pequeno atraso entre o span terminar e ele aparecer no backend.

**8. O OpenTelemetry é, ele mesmo, um backend de observabilidade (tipo o Jaeger)?**
Não. OpenTelemetry é um padrão e um conjunto de SDKs/APIs para **instrumentar** código e **coletar** dados de telemetria de forma vendor-neutral. Ele não armazena nem visualiza dados — isso é responsabilidade de backends como Jaeger, Zipkin, Prometheus, Grafana, Datadog, etc. O OTel é a camada de instrumentação e transporte; o backend é quem persiste e exibe.

**9. O que é o W3C Trace Context e o header `traceparent`?**
É um padrão da W3C que define um formato comum de header HTTP (`traceparent`, no formato `00-<trace-id>-<span-id>-<flags>`) para propagar o contexto de um trace entre serviços, independentemente da linguagem ou framework usado em cada um. Antes desse padrão, cada ferramenta de tracing (Zipkin, Jaeger nativo, etc.) usava seu próprio formato de header, dificultando interoperabilidade entre serviços instrumentados com ferramentas diferentes.

**10. Como o OpenTelemetry ajuda a evitar vendor lock-in?**
Como a instrumentação (código da aplicação) fala apenas o protocolo OTLP, trocar de backend de observabilidade — por exemplo, sair do Jaeger self-hosted para o Datadog — não exige alterar o código instrumentado. Basta apontar o Exporter (ou a configuração do Collector) para o novo destino. Isso desacopla a decisão de "qual ferramenta de observabilidade usar" do código da aplicação.

---

## 🚀 Próximos passos / exercícios sugeridos

- **Conectar o Zipkin**: adicione um exporter `zipkin` no `otel-collector-config.yaml` da raiz e compare o mesmo trace nas UIs do Jaeger e do Zipkin.
- **Corrigir o scrape do Prometheus**: adicione `goapp2:8181` e `goapp3:8282` como targets em `comunicacao-ms/.docker/prometheus.yaml` e confirme no Prometheus que as três séries de métricas aparecem.
- **Adicionar métricas customizadas**: crie um `Counter` (ex: total de requisições) ou `Histogram` (ex: duração da chamada externa) usando `go.opentelemetry.io/otel/metric` em `comunicacao-ms`, e visualize no Grafana.
- **Testar sampling parcial**: troque `sdktrace.AlwaysSample()` por `sdktrace.TraceIDRatioBased(0.1)` (10%) e observe como menos traces aparecem no Jaeger.
- **Importar um dashboard no Grafana**: configure o Prometheus como data source (`http://prometheus:9090`) e monte um painel simples com as métricas expostas em `/metrics`.
- **Corrigir o tratamento de erros do `main.go` raiz**: siga o padrão de `comunicacao-ms/cmd/microservice/main.go`, fazendo cada `fmt.Errorf` virar um `return`/`log.Fatal` de verdade.
- **Adicionar atributos aos spans**: use `span.SetAttributes(...)` para anexar dados de negócio (ex: ID do usuário, rota chamada) e veja como eles aparecem no Jaeger UI.
