# gRPC em Go — Guia Didático Completo

Este módulo ensina como construir APIs **gRPC** em Go. Você vai aprender a definir contratos de serviço usando **Protocol Buffers**, implementar os quatro padrões de comunicação (unário, server streaming, client streaming e bidirecional), conectar um banco de dados SQLite e expor um servidor gRPC altamente performático — tudo isso entendendo *por que* cada decisão existe.

Se você nunca ouviu falar de gRPC, não se preocupe. Vamos começar do zero, construir o raciocínio passo a passo com analogias do mundo real e terminar com um projeto funcional que você poderá usar como referência.

Este guia também inclui um **comparativo profundo entre REST, GraphQL e gRPC** — com casos de uso certos, incorretos, trade-offs reais e exemplos de como grandes empresas escolhem cada tecnologia.

---

## 📑 Sumário

- [⚡ Quick Start — Rodando em 5 Minutos](#-quick-start--rodando-em-5-minutos)
- [🤔 O que é gRPC?](#-o-que-é-grpc)
  - [A Analogia do Telefone Corporativo](#a-analogia-do-telefone-corporativo)
  - [Como o gRPC Funciona na Prática](#como-o-grpc-funciona-na-prática)
  - [Por que o Google Criou o gRPC?](#por-que-o-google-criou-o-grpc)
  - [RPC — A Ideia que Veio Antes](#rpc--a-ideia-que-veio-antes)
- [🌐 HTTP/2 — A Base do gRPC](#-http2--a-base-do-grpc)
  - [O Problema com HTTP/1.1](#o-problema-com-http11)
  - [O que o HTTP/2 Resolveu](#o-que-o-http2-resolveu)
  - [Multiplexação de Streams](#multiplexação-de-streams)
  - [Compressão de Headers com HPACK](#compressão-de-headers-com-hpack)
  - [Server Push](#server-push)
  - [Conexão Persistente](#conexão-persistente)
  - [HTTP/2 vs HTTP/1.1 — Tabela Comparativa](#http2-vs-http11--tabela-comparativa)
- [📦 Protocol Buffers — A Linguagem dos Dados](#-protocol-buffers--a-linguagem-dos-dados)
  - [O Problema com JSON e XML](#o-problema-com-json-e-xml)
  - [O que são Protocol Buffers?](#o-que-são-protocol-buffers)
  - [A Analogia do Formulário](#a-analogia-do-formulário)
  - [Como Definir Mensagens em .proto](#como-definir-mensagens-em-proto)
  - [Tipos de Dados Disponíveis](#tipos-de-dados-disponíveis)
  - [Regras de Numeração dos Campos](#regras-de-numeração-dos-campos)
  - [Serialização Binária vs JSON](#serialização-binária-vs-json)
  - [Geração de Código Automática](#geração-de-código-automática)
- [🔄 Os 4 Padrões de Comunicação gRPC](#-os-4-padrões-de-comunicação-grpc)
  - [1. Unary RPC](#1-unary-rpc)
  - [2. Server-Side Streaming](#2-server-side-streaming)
  - [3. Client-Side Streaming](#3-client-side-streaming)
  - [4. Bidirectional Streaming](#4-bidirectional-streaming)
  - [Diagrama Comparativo dos 4 Padrões](#diagrama-comparativo-dos-4-padrões)
- [⚔️ REST vs GraphQL vs gRPC — O Grande Comparativo](#️-rest-vs-graphql-vs-grpc--o-grande-comparativo)
  - [Filosofias Diferentes](#filosofias-diferentes)
  - [Quando Cada Um Foi Criado e Por Quê](#quando-cada-um-foi-criado-e-por-quê)
  - [Como Cada Um Trafega Dados](#como-cada-um-trafega-dados)
  - [Tabela Comparativa Completa](#tabela-comparativa-completa)
  - [Performance e Overhead](#performance-e-overhead)
  - [Tipagem e Contratos](#tipagem-e-contratos)
  - [Suporte a Streaming](#suporte-a-streaming)
  - [Facilidade de Uso e Ecossistema](#facilidade-de-uso-e-ecossistema)
- [✅ Casos de Uso Certos e Incorretos](#-casos-de-uso-certos-e-incorretos)
  - [Quando Usar REST](#quando-usar-rest)
  - [Quando NÃO Usar REST](#quando-não-usar-rest)
  - [Quando Usar GraphQL](#quando-usar-graphql)
  - [Quando NÃO Usar GraphQL](#quando-não-usar-graphql)
  - [Quando Usar gRPC](#quando-usar-grpc)
  - [Quando NÃO Usar gRPC](#quando-não-usar-grpc)
  - [Tabela de Decisão Rápida](#tabela-de-decisão-rápida)
  - [Exemplos de Empresas e suas Escolhas](#exemplos-de-empresas-e-suas-escolhas)
- [🏗️ Arquitetura do Projeto](#️-arquitetura-do-projeto)
  - [Visão Geral](#visão-geral)
  - [Diagrama de Componentes](#diagrama-de-componentes)
  - [Estrutura de Diretórios](#estrutura-de-diretórios)
- [📄 O Arquivo Proto — Definindo o Contrato](#-o-arquivo-proto--definindo-o-contrato)
  - [Anatomia de um Arquivo .proto](#anatomia-de-um-arquivo-proto)
  - [As Mensagens do Projeto](#as-mensagens-do-projeto)
  - [O Serviço CategoryService](#o-serviço-categoryservice)
  - [Os 5 RPCs Definidos](#os-5-rpcs-definidos)
- [⚙️ Geração de Código](#️-geração-de-código)
  - [O que é Gerado Automaticamente](#o-que-é-gerado-automaticamente)
  - [Os Arquivos .pb.go e _grpc.pb.go](#os-arquivos-pbgo-e-_grpcpbgo)
  - [Por que Não Editar os Arquivos Gerados](#por-que-não-editar-os-arquivos-gerados)
- [🗄️ Camada de Banco de Dados](#️-camada-de-banco-de-dados)
  - [CategoryDB — Estrutura e Métodos](#categorydb--estrutura-e-métodos)
  - [Queries Parametrizadas](#queries-parametrizadas)
- [🔧 Implementação do Serviço](#-implementação-do-serviço)
  - [CategoryService — Estrutura Geral](#categoryservice--estrutura-geral)
  - [CreateCategory — Unary RPC](#createcategory--unary-rpc)
  - [ListCategories — Unary RPC](#listcategories--unary-rpc)
  - [GetCategory — Unary RPC](#getcategory--unary-rpc)
  - [CreateCategoryStream — Client Streaming](#createcategorystream--client-streaming)
  - [CreateCategoryStreamBidirectional — Bidirecional](#createcategorystreambidirectional--bidirecional)
- [🚀 Rodando o Projeto](#-rodando-o-projeto)
  - [Pré-requisitos](#pré-requisitos)
  - [Criando o Banco de Dados](#criando-o-banco-de-dados)
  - [Gerando o Código Proto](#gerando-o-código-proto)
  - [Iniciando o Servidor](#iniciando-o-servidor)
  - [Testando com Evans](#testando-com-evans)
  - [Testando com grpcurl](#testando-com-grpcurl)
- [🪞 gRPC Reflection](#-grpc-reflection)
  - [O que é Reflection](#o-que-é-reflection)
  - [Como Está Habilitado no Projeto](#como-está-habilitado-no-projeto)
- [🧪 Comparativo Prático — Mesma Feature em 3 Tecnologias](#-comparativo-prático--mesma-feature-em-3-tecnologias)
  - [Criar uma Categoria](#criar-uma-categoria)
  - [Listar Categorias](#listar-categorias)
  - [Streaming em Tempo Real](#streaming-em-tempo-real)
- [🏆 Boas Práticas](#-boas-práticas)
  - [Versionamento de Contratos .proto](#versionamento-de-contratos-proto)
  - [Tratamento de Erros gRPC](#tratamento-de-erros-grpc)
  - [Deadlines e Timeouts](#deadlines-e-timeouts)
  - [Interceptors — O Middleware do gRPC](#interceptors--o-middleware-do-grpc)
- [⚠️ Armadilhas Comuns](#️-armadilhas-comuns)
- [📖 Glossário](#-glossário)
- [🎓 Conceitos Aprendidos](#-conceitos-aprendidos)
- [🚀 Próximos Passos](#-próximos-passos)

---

## ⚡ Quick Start — Rodando em 5 Minutos

Esta seção tem todos os comandos para você subir o servidor e testar com Evans rapidamente. Os detalhes de cada passo estão na seção [Rodando o Projeto](#-rodando-o-projeto).

### Passo 1 — Instalar dependências (uma vez só)

```bash
# Instalar Evans (cliente gRPC interativo)
brew install evans          # macOS
# ou: go install github.com/ktr0731/evans@latest

# Instalar grpcurl (alternativa linha de comando)
brew install grpcurl         # macOS
# ou: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Instalar sqlite3
brew install sqlite          # macOS
# apt-get install sqlite3    # Linux
```

### Passo 2 — Criar o banco de dados

```bash
# Entre na pasta do projeto:
cd aulas/12-gRPC

# Crie o banco e as tabelas:
sqlite3 db.sqlite "
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

### Passo 3 — Subir o servidor gRPC

```bash
# Ainda em aulas/12-gRPC/:
go run cmd/grpcServer/main.go

# O servidor sobe silenciosamente na porta 50051.
# Deixe este terminal aberto.
```

### Passo 4 — Conectar com Evans (novo terminal)

```bash
# Abra outro terminal e execute (ainda em aulas/12-gRPC/):
evans --host localhost --port 50051 -r repl
```

Você verá o banner do Evans e o prompt `localhost:50051>`.

### Passo 5 — Navegar e chamar RPCs no Evans

```
# Selecionar o pacote:
localhost:50051> package pb

# Selecionar o serviço:
pb@localhost:50051> service CategoryService

# Criar uma categoria (Unary):
pb.CategoryService@localhost:50051> call CreateCategory
name (TYPE_STRING) => Programação Go
description (TYPE_STRING) => Aprenda Go do zero ao avançado
{
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "name": "Programação Go",
  "description": "Aprenda Go do zero ao avançado"
}

# Listar todas as categorias (Unary):
pb.CategoryService@localhost:50051> call ListCategories
{}
{
  "categories": [
    {
      "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "name": "Programação Go",
      "description": "Aprenda Go do zero ao avançado"
    }
  ]
}

# Buscar por ID (Unary) — cole o ID retornado acima:
pb.CategoryService@localhost:50051> call GetCategory
id (TYPE_STRING) => xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
{
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "name": "Programação Go",
  "description": "Aprenda Go do zero ao avançado"
}

# Client Streaming — envia várias categorias de uma vez:
pb.CategoryService@localhost:50051> call CreateCategoryStream
name (TYPE_STRING) => Docker
description (TYPE_STRING) => Containers e orquestração
# Evans pergunta se quer enviar mais (Ctrl+D para encerrar o stream):
name (TYPE_STRING) => Kubernetes
description (TYPE_STRING) => Orquestração de containers
name (TYPE_STRING) =>   ← pressione Ctrl+D aqui para encerrar
{
  "categories": [
    { "id": "...", "name": "Docker", "description": "Containers e orquestração" },
    { "id": "...", "name": "Kubernetes", "description": "Orquestração de containers" }
  ]
}

# Bidirectional Streaming — para cada item enviado, recebe resposta imediata:
pb.CategoryService@localhost:50051> call CreateCategoryStreamBidirectional
name (TYPE_STRING) => DevOps
description (TYPE_STRING) => CI/CD e automação
{
  "id": "...",
  "name": "DevOps",
  "description": "CI/CD e automação"
}
name (TYPE_STRING) => Cloud
description (TYPE_STRING) => AWS, GCP, Azure
{
  "id": "...",
  "name": "Cloud",
  "description": "AWS, GCP, Azure"
}
name (TYPE_STRING) =>   ← pressione Ctrl+D para encerrar
```

### Comandos úteis dentro do Evans

```
show service              # lista todos os serviços e RPCs disponíveis
show message              # lista todas as mensagens (tipos de dados)
desc CreateCategoryRequest  # descreve os campos de uma mensagem
exit                      # sai do Evans
```

### Alternativa rápida com grpcurl (sem modo interativo)

```bash
# Listar serviços:
grpcurl -plaintext localhost:50051 list

# Criar categoria:
grpcurl -plaintext \
  -d '{"name": "Go", "description": "Aprenda Go"}' \
  localhost:50051 pb.CategoryService/CreateCategory

# Listar categorias:
grpcurl -plaintext -d '{}' localhost:50051 pb.CategoryService/ListCategories

# Buscar por ID:
grpcurl -plaintext \
  -d '{"id": "COLE-O-ID-AQUI"}' \
  localhost:50051 pb.CategoryService/GetCategory
```

---

## 🤔 O que é gRPC?

**gRPC** (Google Remote Procedure Call) é um framework de comunicação entre serviços criado pelo Google em 2015 e aberto ao público em 2016. Ele permite que um programa em uma máquina chame funções em outro programa em outra máquina — como se fossem funções locais — de forma muito mais eficiente do que REST ou GraphQL.

O "g" significa Google, mas também existe uma tradição na comunidade de brincar que o "g" muda de significado a cada versão (gRPC Remote Procedure Calls, por exemplo).

### A Analogia do Telefone Corporativo

Imagine que você trabalha numa grande empresa com várias filiais espalhadas pelo Brasil. Para pedir um relatório para a filial de Recife, você tem duas opções:

**Opção 1 — E-mail (REST/GraphQL):**
- Você escreve um e-mail em português
- Descreve o que precisa com texto livre
- A filial lê, interpreta, monta o relatório
- Te envia de volta por e-mail
- Você lê a resposta
- Cada e-mail tem um cabeçalho enorme (destinatário, assunto, data, assinatura, disclaimer legal...)

**Opção 2 — Sistema de Chamada Direta (gRPC):**
- A empresa instalou um sistema interno onde cada filial tem funções pré-definidas: `GerarRelatorioVendas(mes, ano)`, `ListarFuncionarios(departamento)`, etc.
- Você simplesmente chama `filialRecife.GerarRelatorioVendas(5, 2024)`
- O sistema sabe exatamente quais parâmetros esperar, em que formato e o que retornar
- A comunicação é direta, sem ambiguidade, muito mais rápida

O gRPC é essa segunda opção. Você define as "funções disponíveis" (chamadas de **RPCs**) num arquivo de contrato e qualquer serviço pode chamá-las como se fossem funções locais.

### Como o gRPC Funciona na Prática

```
┌─────────────────────────────────────────────────────────────────┐
│                     VISÃO GERAL DO gRPC                         │
│                                                                  │
│   CLIENTE                              SERVIDOR                  │
│ ┌──────────┐                         ┌──────────────────────┐   │
│ │ Seu app  │                         │   Servidor gRPC      │   │
│ │          │  1. Chama função local  │                      │   │
│ │ stub     │ ──────────────────────► │  Implementação real  │   │
│ │ (gerado) │                         │  da função           │   │
│ │          │  2. Recebe resultado    │                      │   │
│ │          │ ◄────────────────────── │                      │   │
│ └──────────┘                         └──────────────────────┘   │
│      │                                        │                  │
│      └──── Protocol Buffers + HTTP/2 ─────────┘                 │
│            (serialização binária eficiente)                      │
└─────────────────────────────────────────────────────────────────┘
```

O cliente não sabe (e não precisa saber) que está falando com outro processo. Ele chama `categoryService.CreateCategory(...)` e o **stub gerado automaticamente** cuida de:
1. Serializar os parâmetros em formato binário (Protocol Buffers)
2. Enviar pela rede usando HTTP/2
3. Receber a resposta binária
4. Desserializar e retornar o objeto Go

### Por que o Google Criou o gRPC?

O Google opera com **bilhões** de chamadas entre microserviços por segundo. Em 2015, eles já usavam internamente um sistema chamado **Stubby** (o predecessor do gRPC) há mais de uma década. O problema com as alternativas era:

- **REST com JSON**: JSON é um texto legível por humanos, mas isso tem um custo. É lento de serializar/desserializar e ocupa muito espaço na rede.
- **HTTP/1.1**: Uma requisição por vez por conexão — gargalo enorme em alta escala.
- **Falta de contrato forte**: Com REST, não há garantia de que o cliente e servidor concordam no formato exato dos dados.

O gRPC resolve tudo isso de uma vez:
- Dados em **binário** (Protocol Buffers) — muito menores e mais rápidos
- Usa **HTTP/2** — múltiplas chamadas simultâneas na mesma conexão
- **Contrato obrigatório** (arquivo `.proto`) — cliente e servidor gerados do mesmo contrato

### RPC — A Ideia que Veio Antes

O conceito de **Remote Procedure Call** (chamada de procedimento remoto) existe desde os anos 1980. A ideia central é simples: fazer com que chamar uma função remota pareça igual a chamar uma função local.

Antes do gRPC, existiam outras implementações de RPC:
- **CORBA** (anos 90) — complexo, pesado
- **XML-RPC / SOAP** — usava XML, muito verboso
- **Thrift** (Facebook) — similar ao gRPC, mas menos adotado

O gRPC trouxe a ideia de RPC para a era moderna: simples, eficiente, multi-linguagem e com suporte nativo a streaming.

---

## 🌐 HTTP/2 — A Base do gRPC

O gRPC **obrigatoriamente** usa HTTP/2. Não é uma escolha: sem HTTP/2, gRPC não funciona. Para entender por que isso é importante, precisamos entender os problemas do HTTP/1.1.

### O Problema com HTTP/1.1

HTTP/1.1 foi criado em 1997 para uma web que mal existia. Em 1997, uma página HTML tinha algumas imagens e alguns links. Hoje, uma única página pode fazer 300+ requisições HTTP.

**Problema 1 — Uma requisição por vez por conexão:**

```
┌─────── HTTP/1.1 ──────────────────────────────────────────────┐
│                                                                │
│  Conexão TCP #1:  [req1]──►[resp1]──[req2]──►[resp2]          │
│                   ↑ bloqueado até resp1 chegar                 │
│                                                                │
│  Para 4 requisições simultâneas, você precisa de 4 conexões:  │
│  Conexão TCP #1:  [req1]──►[resp1]                            │
│  Conexão TCP #2:  [req2]──►[resp2]                            │
│  Conexão TCP #3:  [req3]──►[resp3]                            │
│  Conexão TCP #4:  [req4]──►[resp4]                            │
│                                                                │
│  Cada conexão TCP tem overhead de handshake (~100-300ms)       │
└────────────────────────────────────────────────────────────────┘
```

**Problema 2 — Headers repetitivos e não comprimidos:**

Cada requisição HTTP/1.1 carrega headers como:
```
Host: api.exemplo.com
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)...
Accept: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
Accept-Language: pt-BR,pt;q=0.9,en;q=0.8
Cache-Control: no-cache
```

Esses headers podem ter 1-2 KB. Se você faz 100 requisições por segundo, são 100-200 KB/s só de cabeçalhos — e o conteúdo deles raramente muda entre requisições!

**Problema 3 — Head-of-Line Blocking:**

Se a requisição #1 demorar muito (um arquivo grande, por exemplo), as requisições #2 e #3 ficam esperando na fila, mesmo que sejam pequenas e rápidas.

### O que o HTTP/2 Resolveu

HTTP/2 foi lançado em 2015 (RFC 7540) e resolveu todos esses problemas:

### Multiplexação de Streams

O maior avanço do HTTP/2: **múltiplas requisições e respostas simultâneas numa única conexão TCP**.

```
┌─────── HTTP/2 ────────────────────────────────────────────────┐
│                                                                │
│  Uma única conexão TCP:                                        │
│                                                                │
│  Stream 1: [req1-frame]──────────────────►[resp1-frame]       │
│  Stream 3: [req2-frame]──────►[resp2-frame]                   │
│  Stream 5: [req3-frame]──►[resp3-frame]                       │
│  Stream 7: [req4-f]──►[resp4-frame]                           │
│                                                                │
│  Tudo simultaneamente, numa só conexão!                        │
│  Sem esperar o anterior terminar.                             │
└────────────────────────────────────────────────────────────────┘
```

No gRPC, cada chamada RPC é um **stream HTTP/2**. Isso permite:
- Múltiplos RPCs simultâneos na mesma conexão
- Streaming bidirecional de dados
- Sem o overhead de criar novas conexões TCP

**O que é um Stream HTTP/2?**

Um stream é uma sequência bidirecional de frames dentro de uma conexão HTTP/2. Cada stream tem um ID numérico (sempre ímpar para streams iniciados pelo cliente: 1, 3, 5, 7...). Múltiplos streams podem existir ao mesmo tempo na mesma conexão TCP.

### Compressão de Headers com HPACK

HTTP/2 usa o algoritmo **HPACK** para comprimir os headers. A ideia é inteligente:

1. **Tabela de headers estáticos**: Os 61 headers mais comuns (`:method`, `:path`, `content-type`, etc.) têm índices pré-definidos. Em vez de enviar `content-type: application/grpc`, você envia um número.

2. **Tabela de headers dinâmicos**: Headers que você enviou antes são armazenados numa tabela compartilhada. Na próxima requisição, em vez de repetir o header, você envia apenas o índice.

Resultado: headers que ocupavam 1KB podem ser comprimidos para poucos bytes nas requisições subsequentes.

### Server Push

HTTP/2 permite que o servidor envie recursos ao cliente **antes mesmo que o cliente os solicite**. O servidor pode "prever" que se o cliente pediu a página HTML, ele vai precisar do CSS e JS logo em seguida.

No contexto do gRPC, o server push é menos relevante, mas a arquitetura de frames bidirecional do HTTP/2 é o que permite o **streaming bidirecional do gRPC**.

### Conexão Persistente

No HTTP/1.1, mesmo com `keep-alive`, havia limitações. No HTTP/2, a conexão é genuinamente persistente e multiplexada — você estabelece uma vez e usa indefinidamente para múltiplos streams.

Para microsserviços que se comunicam com alta frequência, isso elimina o overhead de handshake TCP (e TLS) em cada requisição.

### HTTP/2 vs HTTP/1.1 — Tabela Comparativa

| Característica | HTTP/1.1 | HTTP/2 |
|---------------|----------|--------|
| **Multiplexação** | ❌ Uma req por conexão | ✅ Múltiplas simultâneas |
| **Compressão de Headers** | ❌ Sem compressão | ✅ HPACK |
| **Formato** | Texto (legível) | Binário (frames) |
| **Server Push** | ❌ Não | ✅ Sim |
| **Prioridade de streams** | ❌ Não | ✅ Sim |
| **Head-of-Line Blocking** | ❌ Problema grave | ✅ Resolvido na camada HTTP |
| **Número de conexões TCP** | Muitas (6-8 por host) | Uma por host |
| **Streaming bidirecional** | ❌ Não nativo | ✅ Nativo |
| **Compatibilidade** | Universal | Requer suporte |
| **Debug direto no terminal** | ✅ Fácil (curl) | ❌ Difícil (binário) |

---

## 📦 Protocol Buffers — A Linguagem dos Dados

### O Problema com JSON e XML

JSON é o formato padrão do REST moderno. Ele é legível por humanos, fácil de debugar e amplamente suportado. Mas tem custos significativos:

**JSON (40 bytes):**
```json
{"id":"abc123","name":"Go","desc":"Curso"}
```

**XML (equivalente, 78 bytes):**
```xml
<category><id>abc123</id><name>Go</name><desc>Curso</desc></category>
```

Problemas com ambos:
1. **São texto**: precisam de parsing para virar estruturas de dados. Parsing de texto é lento.
2. **Sem tipagem estrita**: um campo pode ser string, número ou null sem aviso.
3. **Sem schema obrigatório**: o cliente pode receber campos que não esperava (ou não receber campos que esperava).
4. **Verbosos**: nomes dos campos se repetem em cada objeto da lista.

Se você tem um microsserviço que faz 1 milhão de chamadas por segundo, o overhead de serialização JSON pode consumir CPUs inteiras.

### O que são Protocol Buffers?

**Protocol Buffers** (ou **Protobuf**) é um sistema de serialização de dados criado pelo Google. Em vez de texto, ele serializa dados em **formato binário** usando um schema pré-definido.

A ideia central: **o schema já é conhecido por ambos os lados**. Então não precisamos enviar os nomes dos campos — apenas os valores, identificados por números.

**Protobuf equivalente ao JSON acima (8-12 bytes):**
```
0a 06 61 62 63 31 32 33 12 02 47 6f 1a 05 43 75 72 73 6f
(binário — não legível, mas compacto)
```

O campo `id` não precisa ser transmitido como a string "id" — é substituído pelo número `1`. O decodificador do outro lado sabe que o campo `1` é o `id` porque ambos têm o mesmo schema.

### A Analogia do Formulário

Imagine que você precisa enviar informações de uma categoria por correio para outro departamento:

**Sem Protobuf (JSON/XML):**
```
Campo "id": abc123
Campo "name": Programação em Go
Campo "description": Aprenda Go do zero ao avançado
```
Você escreve os nomes dos campos em cada envelope. Se enviar 1000 envelopes, escreve "id", "name", "description" 1000 vezes.

**Com Protobuf:**
Ambos os departamentos têm um formulário padrão com campos numerados:
```
[1] _________  [2] _________  [3] _________
```
Você preenche apenas os valores: `abc123`, `Programação em Go`, `Aprenda Go...`

Sem repetir os nomes. O outro lado sabe que o campo `[1]` é o `id` porque ambos têm o mesmo formulário.

### Como Definir Mensagens em .proto

Um arquivo `.proto` define o schema — o "formulário padrão":

```protobuf
syntax = "proto3";
package pb;
option go_package = "internal/pb";

// Uma "mensagem" é como uma struct
message Category {
    string id = 1;          // campo 1: id do tipo string
    string name = 2;        // campo 2: name do tipo string
    string description = 3; // campo 3: description do tipo string
}
```

**Linha por linha:**

- `syntax = "proto3"` — versão do Protobuf. Proto3 é a versão atual.
- `package pb` — namespace para evitar conflitos de nomes.
- `option go_package` — diz ao compilador onde colocar o código Go gerado.
- `message Category` — define uma estrutura de dados (equivalente a uma `struct` em Go).
- `string id = 1` — campo chamado `id`, tipo `string`, com o **número 1**. Esse número é o que vai no fio (wire), não o nome.

### Tipos de Dados Disponíveis

| Tipo Protobuf | Equivalente Go | Descrição |
|--------------|---------------|-----------|
| `double` | `float64` | Número decimal 64 bits |
| `float` | `float32` | Número decimal 32 bits |
| `int32` | `int32` | Inteiro 32 bits (signed) |
| `int64` | `int64` | Inteiro 64 bits (signed) |
| `uint32` | `uint32` | Inteiro 32 bits sem sinal |
| `uint64` | `uint64` | Inteiro 64 bits sem sinal |
| `bool` | `bool` | Verdadeiro/falso |
| `string` | `string` | UTF-8 |
| `bytes` | `[]byte` | Bytes brutos |
| `repeated T` | `[]T` | Lista/array do tipo T |
| `message` | `struct` | Objeto aninhado |

**Exemplo de tipos variados:**
```protobuf
message Produto {
    string id = 1;
    string nome = 2;
    double preco = 3;
    int32 estoque = 4;
    bool disponivel = 5;
    repeated string tags = 6;    // lista de strings
    Categoria categoria = 7;     // objeto aninhado
    bytes imagem_thumb = 8;      // imagem em bytes
}
```

### Regras de Numeração dos Campos

Os números dos campos são **permanentes** e têm regras importantes:

1. **Cada campo tem um número único** de 1 a 536.870.911
2. **Números 1-15** ocupam 1 byte no encoding — use para os campos mais frequentes
3. **Números 16-2047** ocupam 2 bytes — para campos menos frequentes
4. **Nunca reutilize** um número mesmo que o campo seja removido
5. **Números 19000-19999** são reservados pelo Protobuf internamente

```protobuf
message Category {
    string id = 1;           // 1 byte de overhead — campo frequente ✅
    string name = 2;         // 1 byte de overhead — campo frequente ✅
    string description = 3;  // 1 byte de overhead ✅
    // Se remover um campo, NUNCA reutilize o número:
    // reserved 4;           // assim você documenta que 4 foi usado antes
}
```

### Serialização Binária vs JSON

Vamos quantificar a diferença. Considere uma lista de 100 categorias:

**Em JSON:**
```json
[
  {"id": "550e8400-e29b-41d4-a716-446655440000", "name": "Go", "description": "Programação Go"},
  {"id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "name": "Docker", "description": "Containers"},
  ...
]
```
- Tamanho aproximado: ~10 KB para 100 itens
- Parsing: O Go precisa ler byte a byte, identificar strings, interpretar o JSON

**Em Protobuf:**
- Tamanho aproximado: ~3-4 KB (60-70% menor)
- Parsing: Leitura direta de offsets, sem interpretação de texto
- Velocidade de serialização: 3-10x mais rápido que JSON

```
┌─────────────────────────────────────────────────────┐
│           Benchmark: Serialização                   │
│                                                     │
│  JSON:    ████████████████████  ~10KB, 100μs        │
│  Protobuf: ██████               ~3.5KB, 15μs        │
│                                                     │
│  (valores ilustrativos — variam por payload)        │
└─────────────────────────────────────────────────────┘
```

### Geração de Código Automática

A grande vantagem do Protobuf: você escreve o schema uma vez e o compilador (`protoc`) gera código para **Go, Java, Python, Rust, C++, C#, JavaScript** e mais.

```bash
# Instalar o compilador protoc e plugins Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Gerar código Go a partir do .proto
protoc --go_out=. --go-grpc_out=. proto/course_category.proto
```

Isso gera dois arquivos:
- `course_category.pb.go` — structs Go para as mensagens
- `course_category_grpc.pb.go` — interfaces e implementações para o serviço gRPC

---

## 🔄 Os 4 Padrões de Comunicação gRPC

O gRPC suporta quatro padrões de comunicação. Cada um serve para um caso de uso diferente. Entender quando usar cada um é fundamental.

### 1. Unary RPC

O padrão mais simples: **uma requisição, uma resposta**. Igual ao REST.

```
Cliente ──── Request ────► Servidor
Cliente ◄─── Response ──── Servidor
```

**No arquivo .proto:**
```protobuf
// Sem "stream" = Unary
rpc CreateCategory(CreateCategoryRequest) returns (Category) {}
rpc ListCategories(blank) returns (CategoryList) {}
rpc GetCategory(CategoryGetRequest) returns (Category) {}
```

**Quando usar:**
- Criar, ler, atualizar ou deletar um recurso
- Queries simples que retornam uma resposta única
- A maioria dos casos do dia a dia

**Exemplo real:** Buscar os dados de um usuário, criar um pedido, autenticar login.

### 2. Server-Side Streaming

O cliente envia **uma** requisição, o servidor responde com **um fluxo** de mensagens.

```
Cliente ──── Request ────────────────────► Servidor
Cliente ◄─── Response 1 ──────────────── Servidor
Cliente ◄─── Response 2 ──────────────── Servidor
Cliente ◄─── Response 3 ──────────────── Servidor
...
Cliente ◄─── EOF ─────────────────────── Servidor
```

**No arquivo .proto:**
```protobuf
// "stream" no returns = Server Streaming
rpc ListarLogsEmTempoReal(FiltroRequest) returns (stream LogEntry) {}
```

**Quando usar:**
- Retornar uma lista grande de dados progressivamente
- Notificações em tempo real (o servidor "empurra" eventos)
- Downloads de arquivos grandes em chunks
- Monitoramento: CPU, memória, métricas

**Exemplo real:** "Me manda todos os logs do servidor de hoje" — o servidor vai enviando conforme processa, sem precisar acumular tudo na memória.

### 3. Client-Side Streaming

O cliente envia **um fluxo** de mensagens, o servidor responde com **uma** mensagem ao final.

```
Cliente ──── Request 1 ────────────────► Servidor
Cliente ──── Request 2 ────────────────► Servidor
Cliente ──── Request 3 ────────────────► Servidor
Cliente ──── EOF ──────────────────────► Servidor
                                         Servidor processa tudo
Cliente ◄─── Response único ─────────── Servidor
```

**No arquivo .proto:**
```protobuf
// "stream" no parâmetro = Client Streaming
rpc CreateCategoryStream(stream CreateCategoryRequest) returns (CategoryList) {}
```

**Quando usar:**
- Upload de arquivos grandes em chunks
- Envio de lote (batch) de dados para processamento
- Importação de registros em massa
- Telemetria: enviar métricas periodicamente, receber confirmação ao final

**Exemplo real no projeto:** O cliente envia várias categorias de uma vez (streaming), o servidor cria todas e ao final retorna a lista completa.

### 4. Bidirectional Streaming

Tanto o cliente quanto o servidor enviam **fluxos** de mensagens simultaneamente. Cada lado pode enviar na ordem e no ritmo que quiser.

```
Cliente ──── Request 1 ───────────────────────────► Servidor
Cliente ◄─── Response 1 ─────────────────────────── Servidor
Cliente ──── Request 2 ───────────────────────────► Servidor
Cliente ──── Request 3 ───────────────────────────► Servidor
Cliente ◄─── Response 2 ─────────────────────────── Servidor
Cliente ◄─── Response 3 ─────────────────────────── Servidor
...
```

Não há ordem obrigatória. Cliente e servidor são independentes.

**No arquivo .proto:**
```protobuf
// "stream" nos dois lados = Bidirectional Streaming
rpc CreateCategoryStreamBidirectional(stream CreateCategoryRequest) returns (stream Category) {}
```

**Quando usar:**
- Chat em tempo real entre usuários
- Jogos multiplayer (estado do jogo fluindo nos dois sentidos)
- Colaboração em tempo real (Google Docs, por exemplo)
- Sistemas de trading: ordens chegando, confirmações saindo
- Processamento de stream: o cliente manda dados crus, o servidor manda resultados conforme processa

**Exemplo real no projeto:** O cliente envia uma categoria → o servidor cria e já responde imediatamente com a categoria criada → cliente envia outra → servidor responde → e assim sucessivamente.

### Diagrama Comparativo dos 4 Padrões

```
┌──────────────────────────────────────────────────────────────────┐
│                    OS 4 PADRÕES DO gRPC                          │
│                                                                  │
│  1. UNARY                                                        │
│     Cliente:  ──[req]──────────────────────────────────────►    │
│     Servidor: ◄─────────────────────────────────[resp]──────    │
│                                                                  │
│  2. SERVER STREAMING                                             │
│     Cliente:  ──[req]──────────────────────────────────────►    │
│     Servidor: ◄──[r1]──[r2]──[r3]──[r4]──[EOF]─────────────    │
│                                                                  │
│  3. CLIENT STREAMING                                             │
│     Cliente:  ──[r1]──[r2]──[r3]──[EOF]────────────────────►   │
│     Servidor: ◄────────────────────────────────[resp]──────     │
│                                                                  │
│  4. BIDIRECTIONAL STREAMING                                      │
│     Cliente:  ──[r1]────[r2]──────────[r3]─────────────────►   │
│     Servidor: ◄──────[resp1]────[resp2]────[resp3]─────────     │
│              (independentes, sem ordem fixa)                     │
└──────────────────────────────────────────────────────────────────┘
```

| Padrão | Cliente envia | Servidor responde | Proto keyword |
|--------|--------------|-------------------|--------------|
| Unary | Uma mensagem | Uma mensagem | nenhum `stream` |
| Server Streaming | Uma mensagem | Fluxo de mensagens | `returns (stream T)` |
| Client Streaming | Fluxo de mensagens | Uma mensagem | `(stream T)` no parâmetro |
| Bidirectional | Fluxo de mensagens | Fluxo de mensagens | `stream` nos dois lados |

---

## ⚔️ REST vs GraphQL vs gRPC — O Grande Comparativo

Estas são as três principais tecnologias para comunicação entre serviços hoje. Cada uma tem uma filosofia diferente e casos de uso ideais distintos.

### Filosofias Diferentes

**REST** (Representational State Transfer, 2000):
> "Os recursos são o centro. Cada recurso tem uma URL. Operações padronizadas (GET, POST, PUT, DELETE) determinam o que fazer com eles."

REST pensa em **substantivos** (recursos): `/users`, `/products`, `/orders`.

**GraphQL** (Facebook, 2012/2015):
> "O cliente deve poder pedir exatamente o que precisa — nem mais, nem menos. Uma única requisição pode buscar dados de várias entidades relacionadas."

GraphQL pensa em **grafos de dados**: o cliente navega pelo grafo solicitando os campos que precisa.

**gRPC** (Google, 2015/2016):
> "Chamadas de função entre serviços devem ser tão eficientes quanto chamadas locais. O contrato deve ser estrito e o transporte, otimizado."

gRPC pensa em **funções** (verbos): `CreateUser`, `GetProduct`, `ProcessOrder`.

### Quando Cada Um Foi Criado e Por Quê

| Tecnologia | Ano | Criado por | Problema que resolvia |
|-----------|-----|------------|----------------------|
| REST | 2000 | Roy Fielding (dissertação) | Padronizar a web, substituir SOAP/XML-RPC |
| GraphQL | 2012 (interno) / 2015 (público) | Facebook | App mobile com conexão lenta, overfetching |
| gRPC | 2015 (interno como Stubby) / 2016 (público) | Google | Comunicação eficiente entre bilhões de microserviços |

### Como Cada Um Trafega Dados

**REST:**
```
POST /categories HTTP/1.1
Host: api.exemplo.com
Content-Type: application/json
Authorization: Bearer eyJhbGc...

{
  "name": "Programação Go",
  "description": "Aprenda Go"
}

─────── Resposta ───────

HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": "550e8400-...",
  "name": "Programação Go",
  "description": "Aprenda Go",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**GraphQL:**
```graphql
# Uma única requisição POST /graphql
mutation {
  createCategory(input: {
    name: "Programação Go"
    description: "Aprenda Go"
  }) {
    id
    name
    # não pedi "description" e "created_at" — não vão na resposta
  }
}
```

**gRPC** (representação conceitual — na prática é binário):
```
// O cliente chama como uma função:
response = categoryService.CreateCategory(ctx, {
  name: "Programação Go",
  description: "Aprenda Go"
})

// Na rede: bytes comprimidos do Protobuf via HTTP/2
// Nenhum texto — apenas binário
```

### Tabela Comparativa Completa

| Aspecto | REST | GraphQL | gRPC |
|---------|------|---------|------|
| **Protocolo** | HTTP/1.1 ou HTTP/2 | HTTP/1.1 ou HTTP/2 | HTTP/2 obrigatório |
| **Formato dos dados** | JSON/XML (texto) | JSON (texto) | Protocol Buffers (binário) |
| **Schema/Contrato** | Opcional (OpenAPI) | Obrigatório (SDL) | Obrigatório (.proto) |
| **Tipagem** | Fraca | Forte | Forte |
| **Overfetching** | Comum | ❌ Resolvido | Não aplicável |
| **Underfetching** | Comum | ❌ Resolvido | Não aplicável |
| **Streaming** | SSE / WebSocket (externo) | Subscriptions | Nativo (4 padrões) |
| **Performance** | Média | Média | Alta |
| **Cache HTTP** | ✅ Nativo (GET) | ❌ Difícil | ❌ Não aplicável |
| **Legibilidade** | ✅ Alta (JSON) | ✅ Alta (GraphQL) | ❌ Baixa (binário) |
| **Curva de aprendizado** | Baixa | Média | Média-Alta |
| **Tooling (navegador)** | ✅ Excelente | ✅ Bom | ❌ Limitado |
| **Multi-linguagem** | ✅ Universal | ✅ Boa | ✅ Excelente (codegen) |
| **Versionamento** | URL versioning | Deprecation de campos | Compatibilidade de campo |
| **Geração de código** | Sim (OpenAPI/Swagger) | Sim | Sim (obrigatório) |
| **Suporte a browser** | ✅ Nativo | ✅ Nativo | ❌ Requer gRPC-Web |
| **Tamanho do payload** | Grande | Médio | Pequeno |
| **Latência** | Alta (texto) | Média | Baixa |

### Performance e Overhead

```
┌──────────────────────────────────────────────────────────────┐
│            Performance Relativa (ilustrativa)                │
│                                                              │
│  Tamanho do Payload:                                         │
│  REST/JSON   ████████████████████████████  100%             │
│  GraphQL     ████████████████████          70-80%           │
│  gRPC        ██████████                    30-40%           │
│                                                              │
│  Velocidade de Serialização:                                 │
│  REST/JSON   ████████████████████████████  1x (base)        │
│  GraphQL     ████████████████████          1x (similar)     │
│  gRPC        ████                          5-10x mais rápido│
│                                                              │
│  Throughput (req/s numa conexão):                            │
│  REST/JSON   ████████████████████████████  1x (base)        │
│  gRPC        ████████████████████████████████████████  2-5x │
└──────────────────────────────────────────────────────────────┘
```

### Tipagem e Contratos

**REST sem OpenAPI:**
- Sem garantia de tipos
- Um campo pode mudar de `string` para `int` silenciosamente
- Cliente descobre erros em produção

**REST com OpenAPI:**
- Melhor, mas o schema é opcional e pode ficar desatualizado
- Nem todo time mantém o OpenAPI atualizado

**GraphQL:**
- Schema fortemente tipado obrigatório
- Alterações quebram consultas existentes
- Deprecation de campos é possível

**gRPC:**
- Schema no `.proto` é a única fonte de verdade
- Cliente e servidor gerados do mesmo arquivo
- Impossível ter divergência entre o que o cliente espera e o que o servidor retorna
- Campos adicionados com novos números são ignorados por versões antigas (compatibilidade retroativa)

### Suporte a Streaming

| Tecnologia | Streaming Servidor | Streaming Cliente | Bidirecional |
|-----------|-------------------|-------------------|-------------|
| REST | SSE (Server-Sent Events) | ❌ Não | WebSocket (protocolo diferente) |
| GraphQL | Subscriptions (WebSocket) | ❌ Não | ❌ Não nativo |
| gRPC | ✅ Nativo | ✅ Nativo | ✅ Nativo |

O streaming é onde o gRPC brilha mais. REST e GraphQL foram projetados para o modelo request-response e precisam de soluções externas para streaming.

### Facilidade de Uso e Ecossistema

**REST:**
- `curl`, Postman, Insomnia funcionam nativamente
- Qualquer linguagem tem suporte HTTP básico
- Testes direto no navegador
- Cache HTTP funciona para GET

**GraphQL:**
- GraphiQL, Apollo Studio, GraphQL Playground
- Introspection permite descobrir o schema
- N+1 problem requer DataLoader
- Configuração mais complexa no servidor

**gRPC:**
- Precisa de ferramentas especializadas: Evans, grpcurl, Postman (versões recentes)
- Não funciona diretamente no navegador (precisa de gRPC-Web com proxy)
- Reflection permite introspection (se habilitada)
- Tooling ainda menos maduro que REST

---

## ✅ Casos de Uso Certos e Incorretos

### Quando Usar REST

✅ **APIs públicas consumidas por terceiros**
> Desenvolvedores de terceiros preferem REST por ser familiar. OpenAPI/Swagger facilita a documentação.

✅ **APIs consumidas por browsers diretamente**
> Fetch API, axios, XMLHttpRequest funcionam nativamente com REST. gRPC precisaria de gRPC-Web + proxy.

✅ **Quando o cache HTTP importa**
> GET requests em REST são cacheáveis nativamente por CDNs, proxies e browsers. Reduz carga no servidor.

✅ **Equipes sem experiência com Protobuf**
> REST tem uma curva de aprendizado menor. Qualquer desenvolvedor web já entende HTTP + JSON.

✅ **CRUD simples sem alta performance**
> Um blog, um e-commerce pequeno, uma API admin — REST é perfeito.

✅ **Quando você precisa de simplicidade operacional**
> Debugar REST com Wireshark, curl ou logs é simples. gRPC em binário é mais difícil.

### Quando NÃO Usar REST

❌ **Comunicação interna entre microserviços com alta frequência**
> Se o serviço A chama o serviço B 10.000 vezes por segundo, o overhead do JSON começa a pesar.

❌ **Streaming de dados em tempo real**
> REST não tem streaming nativo. SSE serve para um sentido, WebSocket é um protocolo separado.

❌ **Quando contratos estritos são críticos**
> REST sem OpenAPI não garante que cliente e servidor concordam no formato.

❌ **Quando payload size importa muito**
> IoT, mobile com conexão ruim, sistemas embedded — JSON é pesado demais.

### Quando Usar GraphQL

✅ **Aplicativos mobile com dados variados por tela**
> Cada tela do app pode precisar de dados diferentes. GraphQL permite pedir só o que precisa.

✅ **Frontend com múltiplas fontes de dados**
> Uma query GraphQL pode agregar dados de múltiplos serviços backend.

✅ **Quando overfetching é um problema real**
> Se o app recebe 50 campos mas usa 5, você está desperdiçando banda (crítico no mobile).

✅ **Time de frontend quer autonomia**
> Com GraphQL, o frontend pode evoluir consultas sem depender do backend para cada mudança.

✅ **APIs com dados altamente relacionais**
> Grafos de dados onde o cliente quer navegar livremente pelas relações.

✅ **Múltiplos clientes com necessidades diferentes**
> App mobile, web, tablet — cada um pede o que precisa sem endpoints específicos.

### Quando NÃO Usar GraphQL

❌ **APIs simples com poucos endpoints**
> Para um serviço com 3 endpoints, GraphQL é over-engineering.

❌ **Quando cache HTTP é crítico**
> GraphQL usa POST por padrão — não é cacheável nativamente por CDNs.

❌ **Uploads de arquivos**
> GraphQL não foi projetado para isso — você acaba misturando GraphQL com REST para uploads.

❌ **Equipe pequena sem experiência**
> GraphQL tem curva de aprendizado: schema, resolvers, N+1, DataLoader, subscriptions...

❌ **Microsserviços comunicando entre si internamente**
> GraphQL traz overhead de texto + flexibilidade que microsserviços internos não precisam.

❌ **Quando você precisa de rastreamento simples de queries**
> Em REST, você sabe exatamente qual endpoint foi chamado. Em GraphQL, todas as queries passam por `/graphql`.

### Quando Usar gRPC

✅ **Comunicação interna entre microserviços**
> O caso de uso principal. Eficiente, contratos estritos, código gerado automaticamente.

✅ **Alta performance e baixa latência**
> Sistemas de trading, jogos, real-time analytics — onde cada milissegundo importa.

✅ **Streaming bidirecional**
> Chat, colaboração em tempo real, processamento de streams de dados.

✅ **Múltiplas linguagens no mesmo sistema**
> Um serviço em Go fala com um em Java via gRPC — o contrato `.proto` garante compatibilidade.

✅ **IoT e dispositivos com recursos limitados**
> Protobuf é compacto — ideal para dispositivos com pouca memória e banda limitada.

✅ **Pipelines de processamento de dados**
> Ingestão de telemetria, ETL em tempo real, processamento de eventos.

✅ **Quando você precisa de contratos estritos e versionamento**
> O `.proto` é a única fonte de verdade. Mudanças são controladas.

### Quando NÃO Usar gRPC

❌ **APIs públicas para terceiros**
> Desenvolvedores externos preferem REST. Documentar e consumir uma API gRPC externamente é mais difícil.

❌ **Quando browsers são o cliente principal**
> Browsers não suportam HTTP/2 raw — precisa de gRPC-Web + proxy (Envoy, por exemplo). Complexidade extra.

❌ **Equipes sem familiaridade com Protobuf**
> A curva de aprendizado do `.proto` + toolchain (`protoc`, plugins) pode ser alta para equipes júnior.

❌ **Quando você precisa de cache HTTP**
> gRPC não suporta cache HTTP nativo.

❌ **Debugging e monitoramento simples**
> Logs em JSON são fáceis de ler. Logs em Protobuf binário precisam de ferramentas especiais.

❌ **Projetos pequenos onde simplicidade importa mais**
> Se você tem uma startup com 2 desenvolvedores e 3 serviços, REST resolve mais rápido.

### Tabela de Decisão Rápida

| Situação | Recomendação |
|----------|-------------|
| API pública para devs externos | REST |
| App mobile com dados variados por tela | GraphQL |
| Microserviços comunicando internamente | gRPC |
| Browser como cliente direto | REST ou GraphQL |
| Streaming de dados em tempo real | gRPC |
| Aggregation de múltiplos backends | GraphQL |
| Alta performance e baixa latência | gRPC |
| Cache HTTP importante (CDN) | REST |
| Múltiplas linguagens de programação | gRPC |
| Upload de arquivos grandes | REST (multipart) |
| CRUD simples | REST |
| IoT / dispositivos com recursos limitados | gRPC |
| Equipe pequena, entrega rápida | REST |
| Real-time colaborativo (chat, jogos) | gRPC |

### Exemplos de Empresas e suas Escolhas

| Empresa | Tecnologia | Motivo |
|---------|-----------|--------|
| **Netflix** | gRPC (interno) + REST (externo) | Performance interna, compatibilidade externa |
| **Uber** | gRPC (interno) | Centenas de microserviços, alta performance |
| **Shopify** | GraphQL (API pública) | Parceiros precisam de flexibilidade |
| **GitHub** | GraphQL (v4 API) | Diferentes clientes com necessidades diferentes |
| **Twitter/X** | REST + gRPC (interno) | REST para API pública, gRPC internamente |
| **Google** | gRPC (tudo interno) | Criaram o gRPC para isso |
| **Stripe** | REST | API de pagamentos precisa ser simples e acessível |
| **Facebook** | GraphQL | Criaram o GraphQL para suas necessidades |

> **Conclusão prática**: A maioria das empresas usa **REST para APIs públicas** e **gRPC para comunicação interna**. GraphQL aparece quando o frontend precisa de muita flexibilidade.

---

## 🏗️ Arquitetura do Projeto

### Visão Geral

O projeto implementa um servidor gRPC de gerenciamento de categorias de cursos. Ele recebe requisições via gRPC, processa usando uma camada de serviço e persiste em um banco SQLite.

```
┌────────────────────────────────────────────────────────────────┐
│                     ARQUITETURA GERAL                          │
│                                                                │
│  ┌─────────────────┐         ┌─────────────────────────────┐  │
│  │   CLIENTE gRPC  │         │      SERVIDOR gRPC           │  │
│  │                 │         │      (porta :50051)          │  │
│  │  Evans / grpcurl│◄───────►│                             │  │
│  │  ou outro serv. │  HTTP/2 │  ┌─────────────────────┐   │  │
│  │                 │ Protobuf│  │  CategoryService     │   │  │
│  └─────────────────┘         │  │  (service layer)     │   │  │
│                               │  └──────────┬──────────┘   │  │
│                               │             │               │  │
│                               │  ┌──────────▼──────────┐   │  │
│                               │  │   CategoryDB         │   │  │
│                               │  │   (database layer)   │   │  │
│                               │  └──────────┬──────────┘   │  │
│                               └─────────────┼───────────────┘  │
│                                             │                   │
│                               ┌─────────────▼──────────┐       │
│                               │     SQLite (db.sqlite)  │       │
│                               └────────────────────────┘       │
└────────────────────────────────────────────────────────────────┘
```

### Diagrama de Componentes

```
┌─────────────────────────────────────────────────────────────────┐
│                    COMPONENTES DO PROJETO                        │
│                                                                  │
│  proto/course_category.proto                                     │
│    └── Define o CONTRATO (mensagens + serviço)                  │
│         │                                                        │
│         ▼ protoc gera automaticamente                            │
│                                                                  │
│  internal/pb/                                                    │
│    ├── course_category.pb.go      (structs das mensagens)        │
│    └── course_category_grpc.pb.go (interfaces do serviço)        │
│                                                                  │
│  cmd/grpcServer/main.go                                          │
│    └── Inicializa: DB → Service → Server → Listen               │
│         │                                                        │
│         ▼ instancia                                              │
│                                                                  │
│  internal/service/category.go                                    │
│    └── CategoryService (implementa os 5 RPCs)                   │
│         │                                                        │
│         ▼ usa                                                    │
│                                                                  │
│  internal/database/category.go                                   │
│    └── CategoryDB (CRUD no SQLite)                              │
└─────────────────────────────────────────────────────────────────┘
```

### Estrutura de Diretórios

```
12-gRPC/
├── cmd/
│   └── grpcServer/
│       └── main.go              # Ponto de entrada do servidor
├── internal/
│   ├── database/
│   │   ├── category.go          # Operações de banco para Category
│   │   └── course.go            # Operações de banco para Course
│   ├── pb/                      # Código GERADO pelo protoc (não editar!)
│   │   ├── course_category.pb.go
│   │   └── course_category_grpc.pb.go
│   └── service/
│       └── category.go          # Implementação dos RPCs gRPC
├── proto/
│   └── course_category.proto    # FONTE DA VERDADE — definição do contrato
├── db.sqlite                    # Banco de dados SQLite
├── go.mod                       # Definição do módulo Go
└── go.sum                       # Checksums das dependências
```

**Por que `internal/`?**

Em Go, a pasta `internal/` tem uma regra especial do compilador: apenas o código dentro do mesmo módulo pode importar pacotes de `internal/`. Isso cria um encapsulamento automático — nenhum código externo pode importar `internal/database` ou `internal/service`, garantindo que essas implementações sejam privadas ao módulo.

---

## 📄 O Arquivo Proto — Definindo o Contrato

O arquivo `.proto` é o coração do projeto gRPC. Ele define o que existe, o que pode ser enviado e o que pode ser recebido. É a **única fonte de verdade** — cliente e servidor são gerados a partir dele.

### Anatomia de um Arquivo .proto

```protobuf
// 1. Versão da sintaxe — sempre proto3 em projetos modernos
syntax = "proto3";

// 2. Package — namespace para evitar conflito de nomes entre protos
package pb;

// 3. Opção específica de linguagem — onde gerar o código Go
option go_package = "internal/pb";

// 4. Mensagens — estruturas de dados
message MinhaMsg {
    // tipo   nome   = número_do_campo;
    string id = 1;
}

// 5. Serviço — define as funções RPC disponíveis
service MeuServico {
    rpc MinhaFuncao(MinhaMsg) returns (MinhaMsg) {}
}
```

### As Mensagens do Projeto

Vamos analisar cada mensagem do projeto:

```protobuf
// Mensagem vazia — usada quando não há parâmetro de entrada
// Equivalente ao "void" ou "()" em Go
message blank{}

// Representa uma categoria no banco de dados
message Category {
    string id = 1;           // UUID gerado no servidor
    string name = 2;         // Nome da categoria
    string description = 3;  // Descrição
}

// Request para criar uma categoria
// Não inclui "id" pois o servidor gera o UUID
message CreateCategoryRequest {
    string name = 1;
    string description = 2;
}

// Resposta com múltiplas categorias
// "repeated" = array/slice em Go
message CategoryList {
    repeated Category categories = 1;
}

// Request para buscar uma categoria específica pelo ID
message CategoryGetRequest {
    string id = 1;
}
```

**Por que separar `Category` de `CreateCategoryRequest`?**

Boa prática: `Category` é a representação completa (inclui `id`). `CreateCategoryRequest` é o que o cliente envia para criar (sem `id` — o servidor gera). Isso segue o princípio de menor privilégio: o cliente não deveria enviar um `id` que vai ser ignorado de qualquer forma.

### O Serviço CategoryService

```protobuf
service CategoryService {
    rpc CreateCategory(CreateCategoryRequest) returns (Category) {}
    rpc CreateCategoryStream(stream CreateCategoryRequest) returns (CategoryList) {}
    rpc CreateCategoryStreamBidirectional(stream CreateCategoryRequest) returns (stream Category) {}
    rpc ListCategories(blank) returns (CategoryList) {}
    rpc GetCategory(CategoryGetRequest) returns (Category) {}
}
```

### Os 5 RPCs Definidos

| RPC | Padrão | Entrada | Saída | Descrição |
|-----|--------|---------|-------|-----------|
| `CreateCategory` | Unary | `CreateCategoryRequest` | `Category` | Cria uma categoria |
| `CreateCategoryStream` | Client Streaming | `stream CreateCategoryRequest` | `CategoryList` | Recebe várias, retorna todas |
| `CreateCategoryStreamBidirectional` | Bidirecional | `stream CreateCategoryRequest` | `stream Category` | Responde cada uma em tempo real |
| `ListCategories` | Unary | `blank` | `CategoryList` | Lista todas as categorias |
| `GetCategory` | Unary | `CategoryGetRequest` | `Category` | Busca por ID |

**Notação `stream`:**
- `stream T` no parâmetro = cliente envia múltiplas mensagens do tipo `T`
- `stream T` no returns = servidor envia múltiplas mensagens do tipo `T`
- `stream T` nos dois = bidirecional

---

## ⚙️ Geração de Código

### O que é Gerado Automaticamente

Quando você roda `protoc` com os plugins do Go, dois arquivos são gerados em `internal/pb/`:

**1. `course_category.pb.go` — As Structs**

Contém todas as structs Go correspondentes às mensagens do `.proto`:

```go
// Gerado automaticamente — NÃO EDITE
type Category struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
    Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
    Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

// Getters automáticos:
func (x *Category) GetId() string          { ... }
func (x *Category) GetName() string        { ... }
func (x *Category) GetDescription() string { ... }
```

As tags `protobuf:"bytes,1,..."` dizem ao runtime como serializar/desserializar o campo.

**2. `course_category_grpc.pb.go` — As Interfaces do Serviço**

Contém as interfaces que você precisa implementar (servidor) e os stubs que o cliente usa:

```go
// Interface que você IMPLEMENTA no servidor:
type CategoryServiceServer interface {
    CreateCategory(context.Context, *CreateCategoryRequest) (*Category, error)
    CreateCategoryStream(CategoryService_CreateCategoryStreamServer) error
    CreateCategoryStreamBidirectional(CategoryService_CreateCategoryStreamBidirectionalServer) error
    ListCategories(context.Context, *Blank) (*CategoryList, error)
    GetCategory(context.Context, *CategoryGetRequest) (*Category, error)
    mustEmbedUnimplementedCategoryServiceServer()
}

// Struct que o CLIENTE usa para chamar o servidor:
type CategoryServiceClient interface {
    CreateCategory(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*Category, error)
    CreateCategoryStream(ctx context.Context, opts ...grpc.CallOption) (CategoryService_CreateCategoryStreamClient, error)
    // ...
}
```

### Os Arquivos .pb.go e _grpc.pb.go

```
┌─────────────────────────────────────────────────────────────┐
│                    O QUE CADA ARQUIVO FAZ                   │
│                                                             │
│  course_category.pb.go                                      │
│    ├── type Category struct { ... }                         │
│    ├── type CreateCategoryRequest struct { ... }            │
│    ├── type CategoryList struct { ... }                     │
│    ├── type CategoryGetRequest struct { ... }               │
│    ├── type Blank struct { ... }                            │
│    └── Métodos de serialização/desserialização              │
│                                                             │
│  course_category_grpc.pb.go                                 │
│    ├── CategoryServiceServer (interface p/ implementar)     │
│    ├── UnimplementedCategoryServiceServer (embed obrig.)    │
│    ├── CategoryServiceClient (interface p/ chamar)          │
│    ├── categorySeviceClient (implementação do cliente)      │
│    ├── RegisterCategoryServiceServer() (registrar no server)│
│    └── Tipos de stream para cada RPC de streaming           │
└─────────────────────────────────────────────────────────────┘
```

### Por que Não Editar os Arquivos Gerados

1. **Serão sobrescritos**: Qualquer mudança no `.proto` + `protoc` apagará suas edições.
2. **São complexos**: O código gerado lida com serialização binária de baixo nível — você quase certamente vai quebrar algo.
3. **São corretos**: O gerador foi testado extensivamente. Editar manualmente introduz bugs.

Se você precisa de comportamento customizado, adicione métodos nas suas próprias structs ou use wrappers — nunca edite os arquivos `*.pb.go`.

---

## 🗄️ Camada de Banco de Dados

### CategoryDB — Estrutura e Métodos

O arquivo `internal/database/category.go` implementa o acesso ao banco de dados para a entidade `Category`:

```go
package database

import (
    "database/sql"
    "github.com/google/uuid"
)

// Category é tanto a struct de domínio quanto o repositório
// (o campo `db` é o repositório embutido na struct)
type Category struct {
    db          *sql.DB  // conexão com o banco — privado (minúscula)
    ID          string   // campos públicos para uso externo
    Name        string
    Description string
}

func NewCategory(db *sql.DB) *Category {
    return &Category{db: db}
}
```

**Métodos disponíveis:**

```go
// Create — insere nova categoria no banco e retorna com o ID gerado
func (c *Category) Create(name string, description string) (Category, error) {
    id := uuid.New().String()  // gera UUID v4
    _, err := c.db.Exec(
        "INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)",
        id, name, description,
    )
    if err != nil {
        return Category{}, err
    }
    return Category{ID: id, Name: name, Description: description}, nil
}

// FindAll — retorna todas as categorias
func (c *Category) FindAll() ([]Category, error) {
    rows, err := c.db.Query("SELECT id, name, description FROM categories")
    // ...
}

// Find — retorna uma categoria pelo ID
func (c *Category) Find(id string) (Category, error) {
    var name, description string
    err := c.db.QueryRow(
        "SELECT name, description FROM categories WHERE id = $1", id,
    ).Scan(&name, &description)
    // ...
}

// FindByCourseID — retorna a categoria associada a um curso (via JOIN)
func (c *Category) FindByCourseID(courseID string) (Category, error) {
    // SELECT c.* FROM categories c JOIN courses co ON c.id = co.category_id WHERE co.id = $1
}
```

### Queries Parametrizadas

Note o uso de `$1`, `$2`, `$3` nas queries:

```go
c.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)", id, name, description)
```

**Por que parametrizado?** Para evitar **SQL Injection**. Se o nome da categoria fosse concatenado diretamente:

```go
// PERIGOSO — nunca faça isso:
query := "INSERT INTO categories (name) VALUES ('" + name + "')"
// Se name = "Go'; DROP TABLE categories; --" → desastre
```

Com `$1`, o valor é escapado automaticamente pelo driver do banco.

---

## 🔧 Implementação do Serviço

### CategoryService — Estrutura Geral

```go
package service

import (
    "context"
    "io"

    "github.com/devfullcycle/14-gRPC/internal/database"
    "github.com/devfullcycle/14-gRPC/internal/pb"
)

type CategoryService struct {
    pb.UnimplementedCategoryServiceServer  // embed obrigatório
    CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
    return &CategoryService{
        CategoryDB: categoryDB,
    }
}
```

**O que é `pb.UnimplementedCategoryServiceServer`?**

É uma struct gerada pelo `protoc` que implementa todos os métodos da interface `CategoryServiceServer` com um comportamento padrão: retornar um erro `codes.Unimplemented`. 

Se você não implementar um método mas embeber `Unimplemented...`, o código ainda compila — e o cliente recebe um erro claro ("método não implementado") em vez de um panic ou comportamento indefinido. É uma proteção contra implementações parciais.

### CreateCategory — Unary RPC

```go
func (c *CategoryService) CreateCategory(
    ctx context.Context,
    in *pb.CreateCategoryRequest,
) (*pb.Category, error) {

    // 1. Chama o banco de dados
    category, err := c.CategoryDB.Create(in.Name, in.Description)
    if err != nil {
        return nil, err  // propaga erro para o cliente gRPC
    }

    // 2. Converte domain model → protobuf model
    categoryResponse := &pb.Category{
        Id:          category.ID,
        Name:        category.Name,
        Description: category.Description,
    }

    return categoryResponse, nil
}
```

**Por que converter entre dois tipos?**

`database.Category` é o modelo de domínio (representa o dado no banco). `pb.Category` é o modelo gRPC (gerado pelo protoc, com tags de serialização). Misturá-los criaria acoplamento entre a camada de banco e a camada de transporte — uma mudança no banco quebraria o contrato gRPC e vice-versa.

### ListCategories — Unary RPC

```go
func (c *CategoryService) ListCategories(
    ctx context.Context,
    in *pb.Blank,  // sem parâmetros — mensagem vazia
) (*pb.CategoryList, error) {

    categories, err := c.CategoryDB.FindAll()
    if err != nil {
        return nil, err
    }

    // Converte slice de domain model para slice de protobuf
    var categoriesResponse []*pb.Category
    for _, category := range categories {
        categoriesResponse = append(categoriesResponse, &pb.Category{
            Id:          category.ID,
            Name:        category.Name,
            Description: category.Description,
        })
    }

    return &pb.CategoryList{Categories: categoriesResponse}, nil
}
```

Note: `in *pb.Blank` — mesmo sem parâmetros, o gRPC exige uma mensagem de entrada. A convenção é usar uma mensagem `blank` (ou `Empty` do pacote `google.protobuf`).

### GetCategory — Unary RPC

```go
func (c *CategoryService) GetCategory(
    ctx context.Context,
    in *pb.CategoryGetRequest,
) (*pb.Category, error) {

    category, err := c.CategoryDB.Find(in.Id)
    if err != nil {
        return nil, err
    }

    return &pb.Category{
        Id:          category.ID,
        Name:        category.Name,
        Description: category.Description,
    }, nil
}
```

### CreateCategoryStream — Client Streaming

Este é onde o gRPC começa a mostrar seus superpoderes. O cliente envia várias requisições numa stream, o servidor coleta todas e retorna uma lista ao final.

```go
func (c *CategoryService) CreateCategoryStream(
    stream pb.CategoryService_CreateCategoryStreamServer,
) error {
    // Acumula as categorias criadas
    categories := &pb.CategoryList{}

    for {
        // Recebe a próxima mensagem do cliente
        category, err := stream.Recv()

        // io.EOF = cliente encerrou o stream (não é erro real)
        if err == io.EOF {
            // Envia a resposta final E fecha o stream
            return stream.SendAndClose(categories)
        }
        if err != nil {
            return err
        }

        // Processa cada mensagem recebida
        categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
        if err != nil {
            return err
        }

        // Acumula na lista
        categories.Categories = append(categories.Categories, &pb.Category{
            Id:          categoryResult.ID,
            Name:        categoryResult.Name,
            Description: categoryResult.Description,
        })
    }
}
```

**Fluxo de execução:**

```
Cliente                           Servidor
  │                                  │
  ├──[CreateCategoryRequest A]──────►│ Recv() → "A", nil
  │                                  │ DB.Create("A") → id1
  │                                  │ append(categories, A)
  ├──[CreateCategoryRequest B]──────►│ Recv() → "B", nil
  │                                  │ DB.Create("B") → id2
  │                                  │ append(categories, B)
  ├──[CreateCategoryRequest C]──────►│ Recv() → "C", nil
  │                                  │ DB.Create("C") → id3
  │                                  │ append(categories, C)
  ├──[EOF / CloseStream]────────────►│ Recv() → nil, io.EOF
  │                                  │ SendAndClose(categories)
  │◄─[CategoryList{A, B, C}]────────┤
  │                                  │
```

**Diferença entre `Recv()` retornar `io.EOF` vs erro real:**

- `io.EOF` = o cliente encerrou normalmente o stream (stream fechado com sucesso)
- Qualquer outro erro = algo deu errado (rede, bug no cliente, timeout)

Por isso verificamos `io.EOF` antes do erro genérico.

### CreateCategoryStreamBidirectional — Bidirecional

O padrão mais complexo: para cada mensagem recebida do cliente, o servidor responde imediatamente.

```go
func (c *CategoryService) CreateCategoryStreamBidirectional(
    stream pb.CategoryService_CreateCategoryStreamBidirectionalServer,
) error {
    for {
        // Aguarda próxima mensagem do cliente
        category, err := stream.Recv()
        if err == io.EOF {
            return nil  // cliente encerrou — servidor encerra também
        }
        if err != nil {
            return err
        }

        // Processa
        categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
        if err != nil {
            return err
        }

        // Responde IMEDIATAMENTE para este item específico
        err = stream.Send(&pb.Category{
            Id:          categoryResult.ID,
            Name:        categoryResult.Name,
            Description: categoryResult.Description,
        })
        if err != nil {
            return err
        }
    }
}
```

**Diferença chave em relação ao Client Streaming:**

| | Client Streaming | Bidirectional |
|-|-----------------|---------------|
| `stream.Recv()` | ✅ Usa | ✅ Usa |
| `stream.Send()` | ❌ Não usa | ✅ Usa (resposta por item) |
| `stream.SendAndClose()` | ✅ Usa (ao final) | ❌ Não usa |
| Quando responde | Uma vez, ao final | A cada item recebido |

**Fluxo de execução:**

```
Cliente                           Servidor
  │                                  │
  ├──[Category A]───────────────────►│ Recv() → "A"
  │                                  │ DB.Create("A") → id1
  │◄─[Category{id1, "A", ...}]───────┤ Send(A_com_id)
  ├──[Category B]───────────────────►│ Recv() → "B"
  │                                  │ DB.Create("B") → id2
  │◄─[Category{id2, "B", ...}]───────┤ Send(B_com_id)
  ├──[EOF]──────────────────────────►│ Recv() → io.EOF
  │                                  │ return nil
  │                                  │
```

O servidor não precisa acumular nada — ele responde on-the-fly para cada item.

---

## 🚀 Rodando o Projeto

### Pré-requisitos

Você precisa ter instalado:

1. **Go 1.19+**
   ```bash
   go version
   # go version go1.21.x darwin/arm64
   ```

2. **protoc** — compilador de Protocol Buffers
   ```bash
   # macOS com Homebrew:
   brew install protobuf

   # Linux (Ubuntu/Debian):
   apt-get install -y protobuf-compiler

   # Verificar:
   protoc --version
   # libprotoc 24.x
   ```

3. **Plugins Go para o protoc**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

   # Adicionar ao PATH (se necessário):
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

4. **Evans** — cliente gRPC interativo (opcional, mas muito útil)
   ```bash
   # macOS:
   brew tap ktr0731/evans
   brew install evans

   # Go install:
   go install github.com/ktr0731/evans@latest
   ```

5. **grpcurl** — curl para gRPC (alternativa ao Evans)
   ```bash
   # macOS:
   brew install grpcurl

   # Go install:
   go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
   ```

### Criando o Banco de Dados

O projeto usa SQLite. Crie o banco com as tabelas necessárias:

```bash
# Entre na pasta do projeto
cd aulas/12-gRPC

# Instale o sqlite3 se necessário:
# macOS: brew install sqlite
# Linux: apt-get install sqlite3

# Crie o banco e as tabelas:
sqlite3 db.sqlite "
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

Verificar se funcionou:
```bash
sqlite3 db.sqlite ".tables"
# categories  courses
```

### Gerando o Código Proto

Se você modificar o arquivo `.proto`, regenere o código:

```bash
# Da raiz do projeto (aulas/12-gRPC/)
protoc --go_out=. --go-grpc_out=. proto/course_category.proto

# Isso gera/atualiza:
# internal/pb/course_category.pb.go
# internal/pb/course_category_grpc.pb.go
```

Os arquivos já estão gerados no repositório, então você só precisa deste passo se modificar o `.proto`.

### Iniciando o Servidor

```bash
# Da pasta aulas/12-gRPC/
go run cmd/grpcServer/main.go

# Você verá... nada. O servidor sobe silenciosamente.
# Isso é intencional — o servidor está esperando conexões em :50051
```

Para confirmar que está rodando:
```bash
# Em outro terminal:
lsof -i :50051
# COMMAND   PID   ... TCP *:50051 (LISTEN)
```

O servidor vai ficar rodando até você pressionar `Ctrl+C`.

### Testando com Evans

Evans é um cliente gRPC interativo — como o `psql` para PostgreSQL, mas para gRPC.

```bash
# Com o servidor rodando, em outro terminal:
evans --host localhost --port 50051 -r repl

# Saída esperada:
#
#   ______
#  |  ____|
#  | |__    __  __   __ _   _ __    ___
#  |  __|  \ \/ /  / _' | | '_ \  / __|
#  | |____  >  <  | (_| | | | | | \__ \
#  |______| /_/\_\  \__,_| |_| |_| |___/
#
#  more expressive universal gRPC client
#
# localhost:50051> 
```

O `-r` significa "reflection" — o Evans vai descobrir os serviços disponíveis automaticamente.

**Comandos dentro do Evans:**

```
# Ver todos os serviços disponíveis:
localhost:50051> show service
+─────────────────────────────────────+───────────────────────────────────────────────────+───────────────────────────────────────────────────+
| SERVICE                             | RPC                                               | REQUEST TYPE                                      | RESPONSE TYPE
+─────────────────────────────────────+───────────────────────────────────────────────────+───────────────────────────────────────────────────+
| CategoryService                     | CreateCategory                                    | CreateCategoryRequest                             | Category
| CategoryService                     | CreateCategoryStream                              | CreateCategoryRequest                             | CategoryList
| CategoryService                     | CreateCategoryStreamBidirectional                 | CreateCategoryRequest                             | Category
| CategoryService                     | ListCategories                                    | blank                                             | CategoryList
| CategoryService                     | GetCategory                                       | CategoryGetRequest                                | Category

# Selecionar o pacote:
localhost:50051> package pb

# Selecionar o serviço:
pb@localhost:50051> service CategoryService

# Chamar CreateCategory:
pb.CategoryService@localhost:50051> call CreateCategory
name (TYPE_STRING) => Programação Go
description (TYPE_STRING) => Aprenda Go do zero ao avançado
{
  "id": "7f3a1e2c-4b5d-6789-abcd-ef0123456789",
  "name": "Programação Go",
  "description": "Aprenda Go do zero ao avançado"
}

# Listar categorias:
pb.CategoryService@localhost:50051> call ListCategories
{}
{
  "categories": [
    {
      "id": "7f3a1e2c-4b5d-6789-abcd-ef0123456789",
      "name": "Programação Go",
      "description": "Aprenda Go do zero ao avançado"
    }
  ]
}

# Buscar por ID:
pb.CategoryService@localhost:50051> call GetCategory
id (TYPE_STRING) => 7f3a1e2c-4b5d-6789-abcd-ef0123456789
{
  "id": "7f3a1e2c-4b5d-6789-abcd-ef0123456789",
  "name": "Programação Go",
  "description": "Aprenda Go do zero ao avançado"
}
```

**Testando Client Streaming com Evans:**

```bash
pb.CategoryService@localhost:50051> call CreateCategoryStream
# Evans vai pedir cada campo, você aperta Enter para enviar cada um
name (TYPE_STRING) => Docker
description (TYPE_STRING) => Containers e orquestração

# Para adicionar outro item, Evans pergunta se quer continuar
# Quando terminar, envie EOF (Ctrl+D no Linux/Mac):
name (TYPE_STRING) => Kubernetes
description (TYPE_STRING) => Orquestração de containers

# EOF → servidor retorna a lista:
{
  "categories": [
    { "id": "...", "name": "Docker", "description": "..." },
    { "id": "...", "name": "Kubernetes", "description": "..." }
  ]
}
```

### Testando com grpcurl

grpcurl é uma alternativa linha de comando, sem modo interativo:

```bash
# Listar serviços disponíveis (requer reflection):
grpcurl -plaintext localhost:50051 list
# pb.CategoryService
# grpc.reflection.v1alpha.ServerReflection

# Listar métodos de um serviço:
grpcurl -plaintext localhost:50051 list pb.CategoryService
# pb.CategoryService.CreateCategory
# pb.CategoryService.CreateCategoryStream
# pb.CategoryService.CreateCategoryStreamBidirectional
# pb.CategoryService.GetCategory
# pb.CategoryService.ListCategories

# Chamar CreateCategory:
grpcurl -plaintext -d '{"name": "DevOps", "description": "CI/CD e automação"}' \
  localhost:50051 pb.CategoryService/CreateCategory
# {
#   "id": "abc123...",
#   "name": "DevOps",
#   "description": "CI/CD e automação"
# }

# Listar todas as categorias:
grpcurl -plaintext -d '{}' localhost:50051 pb.CategoryService/ListCategories

# Buscar por ID:
grpcurl -plaintext -d '{"id": "abc123..."}' localhost:50051 pb.CategoryService/GetCategory
```

**`-plaintext`**: Sem TLS. Em produção você usaria certificados e não precisaria dessa flag.

---

## 🪞 gRPC Reflection

### O que é Reflection

**gRPC Reflection** é um serviço especial que permite que clientes descubram quais serviços e métodos um servidor gRPC expõe — sem precisar do arquivo `.proto` em mãos.

É o equivalente do **GraphQL Introspection** (que permite ao GraphiQL descobrir o schema) ou do **OpenAPI/Swagger** (que documenta os endpoints REST).

Sem reflection, o cliente precisa ter o arquivo `.proto` ou os arquivos `.pb.go` gerados. Com reflection, ferramentas como Evans e grpcurl podem descobrir tudo automaticamente.

```
┌─────────────────────────────────────────────────────────────┐
│                  COM REFLECTION                             │
│                                                             │
│  Evans: "Oi servidor, quais serviços você tem?"             │
│  Servidor: "Tenho CategoryService com CreateCategory,       │
│             ListCategories, GetCategory..."                 │
│  Evans: "Quais campos tem CreateCategoryRequest?"           │
│  Servidor: "name (string) e description (string)"          │
│  Evans: "Perfeito, vou montar a interface!"                 │
│                                                             │
│  SEM REFLECTION:                                            │
│  Evans não consegue descobrir — você precisa passar o       │
│  arquivo .proto manualmente: evans --proto ...              │
└─────────────────────────────────────────────────────────────┘
```

### Como Está Habilitado no Projeto

Em `cmd/grpcServer/main.go`:

```go
import "google.golang.org/grpc/reflection"

grpcServer := grpc.NewServer()
pb.RegisterCategoryServiceServer(grpcServer, categoryService)

// Esta linha habilita reflection:
reflection.Register(grpcServer)

lis, err := net.Listen("tcp", ":50051")
grpcServer.Serve(lis)
```

Apenas uma linha: `reflection.Register(grpcServer)`. Isso registra o serviço de reflection no servidor, que passa a responder a requisições de introspecção.

**Atenção em produção**: Reflection expõe a estrutura interna dos seus serviços. Em produção, você pode querer desabilitar (remove essa linha) ou proteger com autenticação para evitar que atacantes mapeiem seus endpoints.

---

## 🧪 Comparativo Prático — Mesma Feature em 3 Tecnologias

Vamos implementar "criar uma categoria" e "listar categorias" em REST, GraphQL e gRPC para você ver a diferença concreta.

### Criar uma Categoria

**REST (HTTP + JSON):**

```bash
# Requisição:
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Go", "description": "Aprenda Go"}'

# Resposta HTTP:
# Status: 201 Created
# Body:
{
  "id": "550e8400-...",
  "name": "Go",
  "description": "Aprenda Go",
  "created_at": "2024-01-15T10:30:00Z"
}
```

```go
// Servidor Go (REST com net/http ou Gin):
func createCategory(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name        string `json:"name"`
        Description string `json:"description"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    category := db.CreateCategory(req.Name, req.Description)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(category)
}
```

**GraphQL:**

```graphql
# Requisição (POST /graphql):
mutation {
  createCategory(input: {
    name: "Go"
    description: "Aprenda Go"
  }) {
    id
    name
    # description não foi pedida — não vai na resposta
  }
}

# Resposta:
{
  "data": {
    "createCategory": {
      "id": "550e8400-...",
      "name": "Go"
    }
  }
}
```

```go
// Resolver Go (gqlgen):
func (r *mutationResolver) CreateCategory(ctx context.Context, input model.NewCategory) (*model.Category, error) {
    category := r.DB.CreateCategory(input.Name, input.Description)
    return &model.Category{
        ID:   category.ID,
        Name: category.Name,
    }, nil
}
```

**gRPC:**

```bash
# Via grpcurl:
grpcurl -plaintext -d '{"name": "Go", "description": "Aprenda Go"}' \
  localhost:50051 pb.CategoryService/CreateCategory

# Resposta:
{
  "id": "550e8400-...",
  "name": "Go",
  "description": "Aprenda Go"
}
```

```go
// Servidor Go (gRPC):
func (c *CategoryService) CreateCategory(
    ctx context.Context,
    in *pb.CreateCategoryRequest,
) (*pb.Category, error) {
    category, err := c.CategoryDB.Create(in.Name, in.Description)
    if err != nil {
        return nil, err
    }
    return &pb.Category{
        Id:          category.ID,
        Name:        category.Name,
        Description: category.Description,
    }, nil
}
```

### Listar Categorias

**REST:**
```bash
# Requisição:
curl http://localhost:8080/categories

# Resposta — SEMPRE retorna todos os campos, mesmo que você não precise:
[
  {"id": "1", "name": "Go", "description": "...", "created_at": "...", "updated_at": "..."},
  {"id": "2", "name": "Docker", "description": "...", "created_at": "...", "updated_at": "..."}
]
```

**GraphQL:**
```graphql
# Você escolhe exatamente quais campos quer:
query {
  listCategories {
    id
    name
    # Não pediu description — não vai na resposta (sem overfetching)
  }
}
```

**gRPC:**
```bash
grpcurl -plaintext -d '{}' localhost:50051 pb.CategoryService/ListCategories
# Retorna todos os campos definidos no proto — sem customização pelo cliente
```

### Streaming em Tempo Real

Este é o caso onde o gRPC não tem rival:

**REST — Server-Sent Events (solução parcial):**
```javascript
// Cliente browser:
const eventSource = new EventSource('/stream/categories');
eventSource.onmessage = (event) => {
    console.log(JSON.parse(event.data));
};

// Servidor Go:
func streamCategories(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    // Apenas servidor → cliente. Cliente não pode enviar dados de volta!
    for _, cat := range getNewCategories() {
        fmt.Fprintf(w, "data: %s\n\n", toJSON(cat))
    }
}
```
❌ Só funciona em um sentido (servidor → cliente)
❌ Precisa de reconexão manual
❌ Não é HTTP nativo — é uma gambiarra em cima de HTTP

**GraphQL — Subscriptions (via WebSocket):**
```graphql
subscription {
  categoryCreated {
    id
    name
  }
}
```
✅ Funciona, mas requer WebSocket (protocolo diferente de HTTP)
❌ Não tem streaming bidirecional nativo
❌ Mais complexo de implementar e escalar

**gRPC — Bidirectional Streaming (nativo):**
```go
// É literalmente o que o projeto implementa:
// O cliente envia categorias, o servidor responde em tempo real
func (c *CategoryService) CreateCategoryStreamBidirectional(
    stream pb.CategoryService_CreateCategoryStreamBidirectionalServer,
) error {
    for {
        req, err := stream.Recv()  // recebe do cliente
        if err == io.EOF { return nil }
        
        result, _ := c.CategoryDB.Create(req.Name, req.Description)
        stream.Send(&pb.Category{...})  // responde imediatamente
    }
}
```
✅ Bidirecional nativo no mesmo protocolo (HTTP/2)
✅ Baixa latência
✅ Tipado e com contrato obrigatório
✅ Funciona com múltiplos streams simultâneos

---

## 🏆 Boas Práticas

### Versionamento de Contratos .proto

O Protobuf foi projetado para **compatibilidade retroativa e futura**. A regra de ouro: **nunca mude o número de um campo existente**.

**O que você PODE fazer sem quebrar clientes existentes:**
```protobuf
// v1 original:
message Category {
    string id = 1;
    string name = 2;
}

// v2 — adicionar campos novos no final:
message Category {
    string id = 1;
    string name = 2;
    string description = 3;  // novo — clientes velhos ignoram
    int32 course_count = 4;  // novo — clientes velhos ignoram
}
```

Clientes que não conhecem o campo `3` simplesmente ignoram ao desserializar. Clientes novos que leem dados antigos (sem o campo `3`) recebem o valor zero (`""` para string, `0` para int).

**O que você NÃO PODE fazer:**
```protobuf
// PROIBIDO — muda o tipo do campo 2:
message Category {
    string id = 1;
    int32 name = 2;  // ❌ era string, virou int — quebra tudo
}

// PROIBIDO — reutiliza o número de um campo removido:
message Category {
    string id = 1;
    // name foi removido, mas o número 2 ainda é "reservado"
    string description = 2;  // ❌ conflito com o "name" antigo
}
```

**Quando você remover um campo:**
```protobuf
message Category {
    string id = 1;
    reserved 2;              // reserva o número
    reserved "name";         // reserva o nome também
    string description = 3;
}
```

**Versionamento de serviços:**

Para mudanças que realmente quebram compatibilidade, use namespaces de versão:
```
proto/
├── v1/
│   └── course_category.proto  (opção go_package = "internal/pb/v1")
└── v2/
    └── course_category.proto  (opção go_package = "internal/pb/v2")
```

### Tratamento de Erros gRPC

O gRPC usa **status codes** próprios (não os códigos HTTP). São 16 códigos definidos:

| Código | Nome | Descrição | Equivalente HTTP |
|--------|------|-----------|-----------------|
| 0 | OK | Sucesso | 200 |
| 1 | CANCELLED | Operação cancelada pelo cliente | - |
| 2 | UNKNOWN | Erro desconhecido | 500 |
| 3 | INVALID_ARGUMENT | Parâmetro inválido | 400 |
| 4 | DEADLINE_EXCEEDED | Timeout | 504 |
| 5 | NOT_FOUND | Recurso não encontrado | 404 |
| 6 | ALREADY_EXISTS | Recurso já existe | 409 |
| 7 | PERMISSION_DENIED | Sem permissão | 403 |
| 8 | RESOURCE_EXHAUSTED | Limite atingido (rate limit) | 429 |
| 13 | INTERNAL | Erro interno | 500 |
| 14 | UNAVAILABLE | Serviço indisponível | 503 |
| 16 | UNAUTHENTICATED | Não autenticado | 401 |

**Como retornar erros com status correto:**

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
    category, err := c.CategoryDB.Find(in.Id)
    if err != nil {
        if err == sql.ErrNoRows {
            // Recurso não encontrado → NOT_FOUND (5)
            return nil, status.Errorf(codes.NotFound, "categoria com id %s não encontrada", in.Id)
        }
        // Erro interno → INTERNAL (13)
        return nil, status.Errorf(codes.Internal, "erro ao buscar categoria: %v", err)
    }
    return &pb.Category{...}, nil
}
```

**No cliente, você pode verificar o código:**
```go
resp, err := client.GetCategory(ctx, &pb.CategoryGetRequest{Id: "xyz"})
if err != nil {
    st, ok := status.FromError(err)
    if ok && st.Code() == codes.NotFound {
        fmt.Println("Categoria não encontrada")
    } else {
        fmt.Println("Erro inesperado:", err)
    }
}
```

### Deadlines e Timeouts

No gRPC, o cliente define um **deadline** — uma hora absoluta no futuro até quando a operação deve completar. O servidor respeita esse deadline automaticamente:

```go
// No cliente:
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.CreateCategory(ctx, &pb.CreateCategoryRequest{
    Name:        "Go",
    Description: "Aprenda Go",
})
// Se o servidor demorar mais de 5 segundos, err será codes.DeadlineExceeded
```

O deadline é propagado automaticamente entre serviços. Se o serviço A chama o serviço B com deadline de 5s, e B chama C, B vai passar o deadline restante para C automaticamente (se você usar o mesmo `ctx`).

**Por que deadline em vez de timeout?**

Timeout é relativo: "espera 5 segundos a partir de agora". Deadline é absoluto: "não passe das 15:30:05.000". Em sistemas distribuídos, deadline é mais correto — se o serviço A espera 5s, e depois passa para B, B não precisa esperar mais 5s, mas apenas o tempo que sobrou do deadline original.

### Interceptors — O Middleware do gRPC

Interceptors são o equivalente de middlewares HTTP (como o `middleware` do Express.js ou o `Middleware` do Gin). Eles interceptam as chamadas RPC antes de chegarem ao handler.

**Unary Interceptor (para RPCs normais):**

```go
// Logger de todas as chamadas:
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    
    // Antes do handler
    log.Printf("Chamando: %s", info.FullMethod)
    
    // Chama o handler real
    resp, err := handler(ctx, req)
    
    // Depois do handler
    log.Printf("Concluído: %s em %v (erro: %v)", info.FullMethod, time.Since(start), err)
    
    return resp, err
}

// Registrar no servidor:
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(loggingInterceptor),
)
```

**Stream Interceptor (para streaming RPCs):**

```go
func streamLoggingInterceptor(
    srv interface{},
    ss grpc.ServerStream,
    info *grpc.StreamServerInfo,
    handler grpc.StreamHandler,
) error {
    log.Printf("Stream iniciado: %s", info.FullMethod)
    err := handler(srv, ss)
    log.Printf("Stream concluído: %s (erro: %v)", info.FullMethod, err)
    return err
}
```

**Casos de uso comuns para interceptors:**
- Logging de todas as chamadas
- Autenticação (valida JWT antes de cada chamada)
- Rate limiting
- Tracing (OpenTelemetry)
- Recuperação de panics
- Validação de entradas

---

## ⚠️ Armadilhas Comuns

**1. Reutilizar números de campo ao modificar o .proto**

```protobuf
// ANTES:
message Category {
    string id = 1;
    string name = 2;
    string old_field = 3;  // vai ser removido
}

// DEPOIS (ERRADO):
message Category {
    string id = 1;
    string name = 2;
    // removeu old_field
    string new_field = 3;  // ❌ reutilizou o 3! Dados antigos vão corromper
}

// DEPOIS (CORRETO):
message Category {
    string id = 1;
    string name = 2;
    reserved 3;            // ✅ reserva o número
    reserved "old_field";  // ✅ reserva o nome
    string new_field = 4;  // ✅ usa número novo
}
```

**2. Esquecer o `io.EOF` no loop de streaming**

```go
// ERRADO — vai tratar EOF como erro:
for {
    req, err := stream.Recv()
    if err != nil {
        return err  // ❌ io.EOF também é um "erro" aqui — vai retornar cedo demais
    }
    // processa req...
}

// CORRETO:
for {
    req, err := stream.Recv()
    if err == io.EOF {
        return nil  // ✅ cliente encerrou normalmente
    }
    if err != nil {
        return err  // ✅ erro real
    }
    // processa req...
}
```

**3. Editar os arquivos gerados pelo protoc**

```
internal/pb/course_category.pb.go      ← NÃO EDITE
internal/pb/course_category_grpc.pb.go ← NÃO EDITE
```

Toda lógica customizada vai em `internal/service/`. O código gerado é regenerado toda vez que você roda `protoc`.

**4. Usar gRPC sem TLS em produção**

No projeto de exemplo, usamos `-plaintext` (sem TLS) para facilitar o desenvolvimento. Em produção, **sempre use TLS**:

```go
// Produção — com TLS:
creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
grpcServer := grpc.NewServer(grpc.Creds(creds))
```

```bash
# grpcurl com TLS (sem a flag -plaintext):
grpcurl -cacert ca.pem -d '{}' api.exemplo.com:443 pb.CategoryService/ListCategories
```

**5. Não usar `context.Context` para cancelamento**

```go
// RUIM — ignora o contexto:
func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
    category, err := c.CategoryDB.Create(in.Name, in.Description)  // não passa ctx
    // ...
}

// BOM — propaga o contexto:
func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
    category, err := c.CategoryDB.CreateWithContext(ctx, in.Name, in.Description)
    // Se o cliente cancelar a chamada, o ctx vai ser cancelado e o DB vai abortar
    // ...
}
```

O `context.Context` carrega: deadlines, tokens de cancelamento, metadados. Sempre propague-o para operações que podem bloquear (banco de dados, HTTP, etc.).

**6. Não usar `Unimplemented...` como embed**

```go
// ERRADO — vai ter problema de compilação se adicionar um novo RPC no .proto:
type CategoryService struct {
    CategoryDB database.Category
    // Não embebe UnimplementedCategoryServiceServer
}

// CORRETO — protege contra implementações parciais:
type CategoryService struct {
    pb.UnimplementedCategoryServiceServer  // ✅ obrigatório
    CategoryDB database.Category
}
```

Sem o embed, se você adicionar um novo RPC no `.proto` e esquecer de implementar, o código nem compila. Com o embed, compila mas retorna `UNIMPLEMENTED` para o método não implementado — melhor do que um panic.

**7. Confundir erros gRPC com erros Go**

```go
// Erros gRPC são verificados com status.FromError():
resp, err := client.GetCategory(ctx, req)
if err != nil {
    // NÃO faça:
    if err == someError { ... }  // ❌ não funciona com erros gRPC

    // FAÇA:
    st, ok := status.FromError(err)
    if ok {
        switch st.Code() {
        case codes.NotFound:
            // ...
        case codes.Internal:
            // ...
        }
    }
}
```

---

## 📖 Glossário

| Termo | Definição |
|-------|-----------|
| **gRPC** | Google Remote Procedure Call — framework para comunicação entre serviços usando HTTP/2 e Protocol Buffers |
| **RPC** | Remote Procedure Call — paradigma de chamar funções em outro processo como se fossem locais |
| **Protocol Buffers (Protobuf)** | Formato de serialização binária do Google, mais eficiente que JSON/XML |
| **proto3** | Versão atual do Protocol Buffers (2016+) |
| **protoc** | Compilador de arquivos `.proto` — gera código em múltiplas linguagens |
| **Stub** | Código gerado automaticamente que o cliente usa para chamar o servidor remotamente |
| **Unary RPC** | Padrão mais simples: uma requisição, uma resposta |
| **Server Streaming** | Servidor envia múltiplas respostas para uma única requisição |
| **Client Streaming** | Cliente envia múltiplas mensagens, servidor responde uma vez ao final |
| **Bidirectional Streaming** | Ambos os lados enviam fluxos de mensagens independentemente |
| **Stream** | Sequência bidirecional de frames dentro de uma conexão HTTP/2 |
| **Frame** | Unidade básica de comunicação no HTTP/2 (header frame, data frame, etc.) |
| **HTTP/2** | Segunda versão do protocolo HTTP, com multiplexação, compressão de headers e streaming nativo |
| **Multiplexação** | Múltiplas requisições simultâneas na mesma conexão TCP |
| **HPACK** | Algoritmo de compressão de headers usado pelo HTTP/2 |
| **TLS** | Transport Layer Security — criptografia da comunicação |
| **Reflection** | Serviço que permite a clientes descobrirem a estrutura de um servidor gRPC dinamicamente |
| **Deadline** | Hora absoluta até quando uma operação deve completar |
| **Interceptor** | Middleware do gRPC — intercepta chamadas antes de chegarem ao handler |
| **Status Code** | Código de resultado de uma chamada gRPC (OK, NOT_FOUND, INTERNAL, etc.) |
| **Serialização** | Processo de converter estruturas de dados em bytes para transmissão |
| **Desserialização** | Processo inverso — bytes → estruturas de dados |
| **Schema** | Definição da estrutura dos dados (o arquivo `.proto` é o schema do gRPC) |
| **`.pb.go`** | Arquivo Go gerado pelo protoc com as structs das mensagens |
| **`_grpc.pb.go`** | Arquivo Go gerado pelo protoc com as interfaces e stubs do serviço gRPC |
| **`Unimplemented...Server`** | Struct gerada que implementa todos os métodos com `UNIMPLEMENTED` — embeded para proteção |
| **`io.EOF`** | Sentinela em Go que indica fim de stream (normal, não é erro) |
| **`SendAndClose()`** | Método de stream do servidor que envia a resposta final e fecha o stream (client streaming) |
| **`Send()`** | Envia uma mensagem no stream sem fechá-lo (server streaming, bidirecional) |
| **`Recv()`** | Recebe a próxima mensagem do stream |
| **Evans** | Cliente gRPC interativo — permite testar RPCs como um REPL |
| **grpcurl** | Ferramenta CLI para chamar RPCs gRPC (como `curl` para REST) |
| **Wire format** | Formato binário dos dados transmitidos pela rede |
| **Field number** | Número inteiro que identifica um campo no Protobuf (substituindo o nome na transmissão) |
| **`repeated`** | Palavra-chave do Protobuf para listas/arrays |
| **`reserved`** | Palavra-chave que reserva números/nomes de campos removidos para evitar reutilização |
| **Overfetching** | Receber mais dados do que precisa (problema do REST) |
| **Underfetching** | Receber menos dados do que precisa, forçando múltiplas requisições |
| **UUID** | Universally Unique Identifier — identificador único gerado pelo servidor |
| **SQLite** | Banco de dados SQL embutido (arquivo único, sem servidor separado) |
| **Context** | Objeto Go que carrega deadlines, cancelamento e metadados — deve ser propagado |

---

## 🎓 Conceitos Aprendidos

Ao estudar este módulo, você aprendeu:

**Fundamentos:**
- [ ] O que é gRPC e por que foi criado pelo Google
- [ ] A diferença entre RPC e REST como paradigmas
- [ ] Por que o gRPC usa HTTP/2 obrigatoriamente
- [ ] Como o HTTP/2 resolve os problemas do HTTP/1.1 (multiplexação, HPACK, conexão persistente)

**Protocol Buffers:**
- [ ] O que são Protocol Buffers e sua vantagem sobre JSON/XML
- [ ] Como escrever um arquivo `.proto` (mensagens, tipos, serviços)
- [ ] A importância dos números de campo e por que nunca reutilizá-los
- [ ] Como gerar código Go automaticamente com `protoc`
- [ ] Por que nunca editar os arquivos `.pb.go` gerados

**Padrões de Comunicação:**
- [ ] Unary RPC — uma requisição, uma resposta
- [ ] Server-Side Streaming — uma requisição, múltiplas respostas
- [ ] Client-Side Streaming — múltiplas requisições, uma resposta
- [ ] Bidirectional Streaming — fluxos independentes nos dois sentidos
- [ ] Como identificar o padrão pelo `.proto` (keyword `stream`)

**Comparativo de Tecnologias:**
- [ ] As filosofias diferentes de REST, GraphQL e gRPC
- [ ] Quando cada tecnologia é a escolha certa
- [ ] Trade-offs de performance, simplicidade e ecossistema
- [ ] Por que empresas usam REST externamente e gRPC internamente

**Implementação:**
- [ ] Como estruturar um projeto gRPC em Go (camadas: proto, pb, service, database, cmd)
- [ ] Como implementar os 4 padrões de streaming em Go
- [ ] O papel do `UnimplementedCategoryServiceServer`
- [ ] Como tratar `io.EOF` nos loops de streaming
- [ ] Como usar `context.Context` para deadlines e cancelamento
- [ ] Como retornar erros gRPC com status codes corretos
- [ ] Como usar Interceptors como middleware

**Ferramentas:**
- [ ] Como habilitar gRPC Reflection
- [ ] Como testar com Evans (cliente interativo)
- [ ] Como testar com grpcurl (linha de comando)
- [ ] Como usar `protoc` para gerar código

---

## 🚀 Próximos Passos

Agora que você domina os fundamentos do gRPC, explore estes tópicos:

**Imediato — Expanda o projeto atual:**
- [ ] Implemente os RPCs de `Course` (o banco já tem a tabela e o `database/course.go` existe)
- [ ] Adicione validação de entrada: o que acontece se `name` for vazio? Retorne `codes.InvalidArgument`
- [ ] Adicione TLS ao servidor para simular produção
- [ ] Escreva um cliente Go que chama os RPCs programaticamente

**Intermediário — Qualidade e Observabilidade:**
- [ ] **Interceptors**: Adicione logging e métricas com interceptors
- [ ] **OpenTelemetry**: Adicione distributed tracing às chamadas gRPC
- [ ] **Health Check**: Implemente `grpc.health.v1.Health` (protocolo padrão de liveness)
- [ ] **Retry e Circuit Breaker**: Use os mecanismos built-in do gRPC para resiliência
- [ ] **Rate Limiting**: Limite requisições por cliente usando interceptors

**Avançado — Arquitetura:**
- [ ] **gRPC Gateway**: Exponha os mesmos endpoints como REST automaticamente via proxy (ideal para APIs públicas)
- [ ] **gRPC-Web**: Permita que browsers se conectem ao servidor gRPC via proxy Envoy
- [ ] **Balanceamento de carga**: Como fazer load balancing com múltiplas instâncias do servidor gRPC
- [ ] **Service Mesh**: Integre com Istio ou Linkerd para comunicação segura entre microserviços
- [ ] **Buf**: Ferramenta moderna para gerenciar schemas Protobuf em equipes (substitui scripts `protoc` manuais)

**Referências para Continuar:**
- Documentação oficial: [grpc.io](https://grpc.io/)
- Documentação Protocol Buffers: [protobuf.dev](https://protobuf.dev/)
- Evans (cliente gRPC): [github.com/ktr0731/evans](https://github.com/ktr0731/evans)
- grpcurl: [github.com/fullstorydev/grpcurl](https://github.com/fullstorydev/grpcurl)
- Buf (gerenciamento de schemas): [buf.build](https://buf.build/)
- gRPC-Gateway: [github.com/grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

---

> **Resumo final**: gRPC é a escolha certa quando você precisa de **alta performance, contratos estritos e streaming nativo** — especialmente para comunicação **interna entre microserviços**. Ele não substitui REST ou GraphQL, mas complementa: REST para APIs públicas, GraphQL para frontends com dados variados, gRPC para o "motor" interno do sistema.
