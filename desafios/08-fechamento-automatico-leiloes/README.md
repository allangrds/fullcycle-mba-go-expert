# Fechamento Automático de Leilões

API de leilões em Go (Gin + MongoDB) construída sobre a base oficial do curso
([`devfullcycle/labs-auction-goexpert`](https://github.com/devfullcycle/labs-auction-goexpert)),
acrescentando o fechamento automático: assim que um leilão é criado, uma Goroutine em
background monitora o tempo e, quando a duração configurada expira, atualiza o status do leilão
para **fechado** no MongoDB — sem nenhuma intervenção manual e sem bloquear quem criou o leilão.

## Sumário

- [Como instalar](#como-instalar)
- [Como executar o projeto](#como-executar-o-projeto)
- [Testando manualmente o fechamento automático](#testando-manualmente-o-fechamento-automático)
- [Como rodar os testes](#como-rodar-os-testes)
- [Configuração](#configuração)
- [Arquitetura do projeto](#arquitetura-do-projeto)
- [Resumo de como o projeto foi feito](#resumo-de-como-o-projeto-foi-feito)

## Como instalar

O projeto foi pensado para viver **inteiramente dentro do Docker** — não é necessário ter Go
nem MongoDB instalados na máquina.

Pré-requisitos: **Docker Engine 20.10+** com o plugin **Docker Compose V2** (`docker compose`,
já incluso no Docker Desktop). O `Makefile` é opcional, só um atalho para os comandos do
compose.

Antes de subir o projeto, copie o `.env.example` para `.env` (veja
[Configuração](#configuração) para o que cada variável faz):

```bash
cp .env.example .env
```

## Como executar o projeto

```bash
make up
# equivalente a: docker compose up --build
```

Isso sobe o MongoDB e a API na porta `8080`.

| Comando | O que faz |
|---|---|
| `make up` | Sobe MongoDB + API em primeiro plano |
| `make stop` | Para os containers do projeto, sem removê-los |
| `make destroy` | Derruba containers, remove imagens locais (app/test), volumes e redes do projeto |
| `make logs` | Acompanha os logs do app |
| `make ps` | Lista os containers do projeto |
| `make test` | Sobe o MongoDB e roda toda a suíte de testes |

### Endpoints disponíveis

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/auction` | Cria um leilão (status inicial: aberto) |
| `GET` | `/auction` | Lista leilões (filtros `status`, `category`, `productName`) |
| `GET` | `/auction/:auctionId` | Busca um leilão por id |
| `GET` | `/auction/winner/:auctionId` | Leilão + lance vencedor |
| `POST` | `/bid` | Registra um lance |
| `GET` | `/bid/:auctionId` | Lista os lances de um leilão |
| `GET` | `/user/:userId` | Busca um usuário por id |

## Testando manualmente o fechamento automático

```bash
# 1. Cria um leilão (guarde o "id" da resposta em GET /auction para o próximo passo)
curl -i -X POST localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{
        "product_name": "Notebook Dell",
        "category": "Eletrônicos",
        "description": "Notebook seminovo em ótimo estado de conservação",
        "condition": 1
      }'

# 2. Confirma que o leilão está aberto (status 0 = Active)
curl -s "localhost:8080/auction?status=0" | jq

# 3. Espera o tempo configurado em AUCTION_DURATION no .env (o valor padrão do
#    .env.example é 5m; para testar rápido, use algo como AUCTION_DURATION=10s)

# 4. Confirma que o leilão fechou sozinho (status 1 = Completed/Closed), sem
#    nenhuma ação manual:
curl -s "localhost:8080/auction?status=1" | jq
```

## Como rodar os testes

```bash
make test
# equivalente a: docker compose run --build --rm test
```

Esse serviço sobe o MongoDB (com healthcheck, então os testes só começam quando ele está
pronto), builda até o estágio `builder` do `Dockerfile` (que ainda tem o código-fonte e o
toolchain do Go) e roda `go test ./... -v`. O teste que prova o fechamento automático
(`internal/infra/database/auction/create_auction_test.go`) cria um leilão com um
`AUCTION_DURATION` curto (2s, fixado no próprio teste via `os.Setenv`, independente do
`.env`), confirma que ele nasce com status aberto, espera a duração expirar e confirma que o
status virou fechado sozinho — reproduzindo exatamente o cenário pedido no desafio.

Se preferir rodar localmente com Go instalado (sem Docker), suba só o MongoDB
(`docker compose up -d mongodb`) e rode `MONGODB_URL=mongodb://admin:admin@localhost:27017/auctions?authSource=admin go test ./... -v`.

## Configuração

Todas as configurações vêm de variáveis de ambiente, definidas no arquivo `.env` na raiz do
projeto (não commitado — copie de `.env.example`):

| Variável | Padrão sugerido | Descrição |
|---|---|---|
| `AUCTION_DURATION` | `5m` | **A variável do desafio.** Duração de um leilão (formato de `time.ParseDuration`: `20s`, `5m`, `1h`...). Depois desse tempo a Goroutine de fechamento atualiza o leilão para status fechado. Se ausente ou inválida, cai no fallback de 5 minutos definido no código. |
| `AUCTION_INTERVAL` | `5m` | Já existia na base do projeto: usada em `create_bid.go` para o repositório de lances saber se um leilão ainda está aberto antes de aceitar um bid. É uma checagem independente da Goroutine de fechamento — mantenha-a igual a `AUCTION_DURATION` para os dois lados concordarem sobre quando um leilão fecha. |
| `BATCH_INSERT_INTERVAL` | `20s` | Intervalo de flush do lote de lances (já existia na base do projeto). |
| `MAX_BATCH_SIZE` | `4` | Tamanho máximo do lote de lances antes de gravar no Mongo (já existia na base do projeto). |
| `MONGO_INITDB_ROOT_USERNAME` / `MONGO_INITDB_ROOT_PASSWORD` | `admin` / `admin` | Credenciais do usuário root criado pela imagem oficial do MongoDB. |
| `MONGODB_URL` | `mongodb://admin:admin@mongodb:27017/auctions?authSource=admin` | String de conexão usada pela API. O host `mongodb` é o nome do serviço no `docker-compose.yaml`. |
| `MONGODB_DB` | `auctions` | Nome do database usado pela API. |

## Arquitetura do projeto

```
8-fechamento-automatico-leiloes/
├── cmd/auction/main.go              # monta as dependências, registra as rotas e sobe o servidor
├── configuration/
│   ├── database/mongodb/             # conexão com o MongoDB
│   ├── logger/                        # logger estruturado (zap)
│   └── rest_err/                      # erro HTTP padronizado
├── internal/
│   ├── entity/                        # regras de negócio puras (Auction, Bid, User)
│   ├── usecase/                       # orquestração (DTOs de entrada/saída, casos de uso)
│   ├── infra/database/                # implementação dos repositórios (Mongo)
│   │   └── auction/create_auction.go    # ONDE A MÁGICA ACONTECE — feature deste desafio
│   ├── infra/api/web/controller/      # handlers HTTP (Gin)
│   └── internal_error/                # erro de domínio
├── Dockerfile                        # multi-stage: builder (Go) -> imagem final enxuta
├── docker-compose.yaml               # app + mongodb + test
└── .env                              # configuração local (não commitado)
```

**Onde fica a feature:** toda a lógica nova está em
`internal/infra/database/auction/create_auction.go`, exatamente onde o enunciado pede.
`AuctionRepository` ganhou um campo `auctionDuration` (lido do `AUCTION_DURATION` uma única vez,
na criação do repositório) e, depois de inserir o leilão no Mongo, `CreateAuction` dispara
`go ar.scheduleAuctionClosure(auctionEntity.Id)` — uma Goroutine que dorme pela duração
configurada e então faz um `UpdateOne` mudando o status para `Completed` (o "Closed" do
enunciado). Nenhum outro arquivo da base mudou de comportamento: a validação que já impedia
lances em leilões fechados (`internal/infra/database/bid/create_bid.go`, guiada por
`AUCTION_INTERVAL`) continua exatamente como estava.

## Resumo de como o projeto foi feito

- A base é obrigatória pelo enunciado do desafio (repositório
  `devfullcycle/labs-auction-goexpert`), então toda a estrutura em camadas (`entity`,
  `usecase`, `infra/database`, `infra/api/web`), o Gin como router e o MongoDB como banco foram
  mantidos como estão — só o módulo Go foi renomeado para o padrão usado nos outros desafios
  deste repositório (`github.com/allangrds/fullcycle-mba-go-expert/desafios/<n>-<slug>`) e o
  `.env` passou a ficar na raiz do projeto em vez de `cmd/auction/.env`.
- **Como a Goroutine foi escrita:** o próprio código base já lê durações de tempo de env vars
  com o padrão `os.Getenv` + `time.ParseDuration` + fallback (veja `getAuctionInterval` e
  `getMaxBatchSizeInterval` em `bid_usecase`/`create_bid.go`) — `getAuctionDuration()` segue
  exatamente esse idioma para `AUCTION_DURATION`. Para o agendamento em si, foi usado o padrão
  mais simples possível — uma Goroutine com `time.Sleep` — em vez de `time.Timer`/`select`: não
  há necessidade de resetar ou cancelar o fechamento no meio do caminho, então a ferramenta mais
  simples resolve. A Goroutine cria seu próprio `context.Background()` (em vez de reaproveitar o
  `ctx` da requisição HTTP que criou o leilão) para não correr risco de ser cancelada quando essa
  requisição termina.
- **Por que `AUCTION_DURATION` além de `AUCTION_INTERVAL`:** o enunciado deixa explícito que a
  validação de lances em leilões fechados "já existe" e não deve ser mexida — ela usa
  `AUCTION_INTERVAL`. Em vez de reaproveitar essa env var para a nova Goroutine (o que exigiria
  tocar num arquivo fora do escopo pedido), foi criada `AUCTION_DURATION` dedicada à nova rotina,
  como o próprio enunciado sugere ("ex: AUCTION_DURATION"). As duas ficam com o mesmo valor no
  `.env` por padrão — documentado acima para não gerar confusão entre quem lê o projeto.
- **Teste automatizado:** o padrão usado nas aulas do curso para testes de repositório (UOW,
  Clean Architecture) sempre bate num banco real local via docker-compose com `testify`, nunca
  mock de driver ou testcontainers — o mesmo caminho foi seguido aqui. O teste conecta num
  MongoDB real, cria um leilão com `AUCTION_DURATION` curto, confirma o status inicial aberto,
  espera a duração expirar e confirma que o status virou fechado sozinho. O serviço `test` do
  `docker-compose.yaml` (baseado no mesmo padrão do desafio
  `6-clima-google-cloud-run`, que builda até o estágio `builder` do Dockerfile e roda
  `go test ./... -v`) sobe o MongoDB automaticamente antes de rodar a suíte, então `make test`
  funciona sem nenhuma dependência instalada localmente.
- **Bug corrigido do `.env` original:** o `.env` de exemplo do repositório base definia
  `MONGO_INITDB_ROOT_USERNAME`/`PASSWORD` com `:` em vez de `=`, o que nunca funcionaria como
  env file de verdade — corrigido no `.env.example` deste projeto.
- **Bug corrigido no carregamento do `.env`:** o `main.go` original dava `log.Fatal` se
  `godotenv.Load` não encontrasse o arquivo `.env` no disco. Isso quebra a própria forma como
  este projeto é distribuído: o `.dockerignore` (propositalmente) não copia o `.env` para dentro
  da imagem, já que o `docker-compose.yaml` injeta as variáveis diretamente via `env_file`. O
  container rodava e morria só com "Error trying to load env variables" mesmo com todas as
  variáveis corretas no `.env` da raiz. A correção loga um aviso em vez de encerrar o processo
  quando o arquivo não existe — `godotenv.Load` continua útil para quem roda `go run`/`go test`
  fora do Docker.
