# Upload para Amazon S3 em Go - Guia Didático Completo

Este guia explica de forma didática como fazer upload de arquivos para o Amazon S3 usando Go, com controle de concorrência, tratamento de erros e as boas práticas que separam um código de produção de um código de laboratório.

## Sumário

### [O que é o Amazon S3?](#o-que-é-o-amazon-s3)
- [Buckets e Objetos](#buckets-e-objetos)
- [Chaves (Keys)](#chaves-keys)
- [Regiões](#regiões)

### [Conceitos de Go Usados no Projeto](#conceitos-de-go-usados-no-projeto)
- [A função init()](#a-função-init)
- [Variáveis Globais de Pacote](#variáveis-globais-de-pacote)
- [O AWS SDK para Go](#o-aws-sdk-para-go)

### [Como o Projeto Funciona](#como-o-projeto-funciona)
- [O Gerador de Arquivos (generator)](#o-gerador-de-arquivos-generator)
- [O Uploader (uploader)](#o-uploader-uploader)

### [Concorrência no Projeto](#concorrência-no-projeto)
- [sync.WaitGroup — Coordenando Goroutines](#syncwaitgroup--coordenando-goroutines)
- [Channels Bufferizados como Semáforo](#channels-bufferizados-como-semáforo)
- [Canal de Erros e Retry Automático](#canal-de-erros-e-retry-automático)
- [Diagrama do Fluxo de Concorrência](#diagrama-do-fluxo-de-concorrência)

### [Boas Práticas Presentes no Projeto](#boas-práticas-presentes-no-projeto)
- [Limitação de Concorrência](#limitação-de-concorrência)
- [Retry Automático via Canal Separado](#retry-automático-via-canal-separado)
- [Shutdown Gracioso com wg.Wait()](#shutdown-gracioso-com-wgwait)
- [Upload por Streaming](#upload-por-streaming)
- [Liberação de Recursos com defer](#liberação-de-recursos-com-defer)

### [O que as Boas Práticas Evitaram](#o-que-as-boas-práticas-evitaram)
- [Goroutine Explosion](#goroutine-explosion)
- [Uploads Perdidos Silenciosamente](#uploads-perdidos-silenciosamente)
- [Race Conditions](#race-conditions)
- [Out of Memory por Arquivos Grandes](#out-of-memory-por-arquivos-grandes)

### [O que Poderia Ser Melhorado](#o-que-poderia-ser-melhorado)
- [Credenciais Seguras](#credenciais-seguras)
- [Configuração via Variáveis de Ambiente](#configuração-via-variáveis-de-ambiente)
- [Timeout com Context](#timeout-com-context)
- [Retry com Limite e Backoff Exponencial](#retry-com-limite-e-backoff-exponencial)
- [Graceful Shutdown com Sinais do OS](#graceful-shutdown-com-sinais-do-os)
- [Logging Estruturado](#logging-estruturado)
- [Multipart Upload para Arquivos Grandes](#multipart-upload-para-arquivos-grandes)
- [Problemas no Código Atual](#problemas-no-código-atual)

### [Principais Problemas ao Trabalhar com S3](#principais-problemas-ao-trabalhar-com-s3)
- [1. Gerenciamento de Credenciais](#1-gerenciamento-de-credenciais)
- [2. Rate Limiting e Throttling](#2-rate-limiting-e-throttling)
- [3. Arquivos Grandes e Multipart Upload](#3-arquivos-grandes-e-multipart-upload)
- [4. Falhas de Rede e Idempotência](#4-falhas-de-rede-e-idempotência)
- [5. Custos Inesperados](#5-custos-inesperados)
- [6. Permissões IAM e Bucket Policies](#6-permissões-iam-e-bucket-policies)
- [7. Inconsistência de Região](#7-inconsistência-de-região)
- [8. Nomes de Objetos e Performance](#8-nomes-de-objetos-e-performance)

### [Glossário](#glossário)

### [Próximos Passos](#próximos-passos)

---

## O que é o Amazon S3?

Imagine que você tem uma empresa e precisa guardar milhões de fotos, vídeos, documentos e arquivos de backup. Onde você coloca tudo isso? Em um servidor local? E quando o HD enche? E se o servidor pegar fogo?

O **Amazon S3 (Simple Storage Service)** resolve exatamente esse problema. Pense nele como um **armazém infinito na nuvem** — você envia seus arquivos, a AWS cuida do armazenamento, da redundância e da disponibilidade. Você só paga pelo que usar.

```
Sem S3 (servidor local)          Com S3 (nuvem)
─────────────────────            ──────────────────────────
🖥️ Servidor                      ☁️  AWS S3
   ├── HD cheio = problema          ├── 🗂️ Bucket A (fotos)
   ├── Se cair = tudo perdido       │    ├── foto1.jpg
   └── Backup manual = trabalhoso   │    ├── foto2.jpg
                                    │    └── video.mp4
                                    └── 🗂️ Bucket B (backups)
                                         ├── backup-2024-01.zip
                                         └── backup-2024-02.zip
```

**Casos de uso reais:**
- Netflix armazena vídeos e thumbnails
- Dropbox usa S3 como backend de armazenamento
- Sistemas de backup automático
- Hospedagem de sites estáticos
- Armazenamento de logs e arquivos de dados

### Buckets e Objetos

O S3 tem dois conceitos centrais que você precisa entender:

**Bucket** — É como uma "pasta raiz" no S3. Cada bucket tem um nome único **no mundo inteiro** (não só na sua conta). É dentro do bucket que você guarda seus arquivos.

**Objeto** — É qualquer arquivo que você armazena no S3: imagens, vídeos, PDFs, arquivos zip, logs, tudo. Cada objeto pode ter até 5 TB de tamanho.

```
☁️ AWS S3
└── 🗂️ minha-empresa-producao (Bucket)
    ├── 📷 imagens/avatar-usuario-123.jpg (Objeto)
    ├── 📄 documentos/contrato-2024.pdf   (Objeto)
    └── 🎬 videos/tutorial-intro.mp4       (Objeto)
```

### Chaves (Keys)

No S3, o "nome" de um objeto é chamado de **Key** (chave). A key é o caminho completo do arquivo dentro do bucket:

```
Bucket: minha-empresa-producao
Key:    imagens/avatar-usuario-123.jpg

URL resultante:
https://minha-empresa-producao.s3.amazonaws.com/imagens/avatar-usuario-123.jpg
```

> **Importante:** O S3 não tem pastas de verdade. O que parece ser uma pasta (`imagens/`) é apenas parte da key. É uma convenção para organizar os arquivos.

### Regiões

A AWS tem datacenters espalhados pelo mundo, chamados de **regiões**. Ao criar um bucket, você escolhe em qual região ele ficará:

| Região | Localização |
|--------|-------------|
| `us-east-1` | Norte da Virgínia, EUA |
| `us-west-2` | Oregon, EUA |
| `sa-east-1` | São Paulo, Brasil |
| `eu-west-1` | Irlanda |

Escolha a região mais próxima dos seus usuários para menor latência. O código desta aula usa `us-east-1`.

---

## Conceitos de Go Usados no Projeto

### A função init()

Em Go, existe uma função especial chamada `init()`. Ela é executada **automaticamente** antes do `main()`, sem precisar ser chamada por ninguém. É como um "preparador de palco" que configura tudo antes do show começar.

```go
package main

func init() {
    // Este código roda ANTES do main()
    // Ideal para: configurar conexões, inicializar clientes, validar variáveis de ambiente
    fmt.Println("Preparando tudo...")
}

func main() {
    // Este código roda DEPOIS do init()
    fmt.Println("Show começando!")
}

// Saída:
// Preparando tudo...
// Show começando!
```

**Quando usar `init()`?**
- Inicializar clientes de banco de dados ou APIs externas (como o cliente S3)
- Validar configurações obrigatórias antes de rodar o programa
- Registrar drivers ou plugins

**Quando NÃO usar?**
- Para lógica de negócio (dificulta testes)
- Quando a inicialização pode falhar de forma recuperável (prefira retornar erros)

No projeto, o `init()` configura a sessão AWS e cria o cliente S3:

```go
func init() {
    sess, err := session.NewSession(
        &aws.Config{
            Region: aws.String("us-east-1"),
            Credentials: credentials.NewStaticCredentials("---", "---", ""),
        },
    )
    if err != nil {
        panic(err)
    }
    s3Client = s3.New(sess)
    s3Bucket = "goexpert-bucket-exemplo"
}
```

### Variáveis Globais de Pacote

O projeto usa variáveis declaradas fora de qualquer função, no nível do pacote:

```go
var (
    s3Client *s3.S3      // cliente do S3 — compartilhado entre todas as goroutines
    s3Bucket string      // nome do bucket — lido por todas as goroutines
    wg       sync.WaitGroup  // coordenador de goroutines
)
```

Variáveis globais são acessíveis por todas as funções do pacote. Aqui faz sentido porque `s3Client` e `s3Bucket` são configurados uma vez no `init()` e depois apenas lidos (nunca escritos), o que é seguro para uso concorrente.

> **Regra prática:** variáveis globais que são apenas lidas após a inicialização são seguras para concorrência. Variáveis globais que são escritas concorrentemente precisam de proteção (Mutex, atomic, etc.).

### O AWS SDK para Go

O projeto usa o **AWS SDK para Go v1** (`github.com/aws/aws-sdk-go`). O SDK é uma biblioteca oficial da Amazon que cuida de toda a comunicação com os serviços AWS: autenticação, serialização das requisições, retry automático e muito mais.

Para usar o S3, você precisa de três componentes:

```
1. Session (sessão)     → Autenticação e configuração regional
2. S3 Client (cliente)  → Interface para chamar as APIs do S3
3. Input (parâmetros)   → Dados específicos de cada operação (bucket, key, body)
```

---

## Como o Projeto Funciona

O projeto é dividido em dois programas independentes:

```
┌─────────────────────────────────────────────────────────────┐
│                      FLUXO DO PROJETO                       │
│                                                             │
│  ┌───────────┐      ┌──────────────┐      ┌─────────────┐  │
│  │ generator │─────▶│  ./tmp/      │─────▶│  uploader   │  │
│  │           │      │  file0.txt   │      │             │  │
│  │ Cria      │      │  file1.txt   │      │ Lê arquivos │  │
│  │ arquivos  │      │  file2.txt   │      │ e faz upload│  │
│  │ de teste  │      │  ...         │      │             │  │
│  └───────────┘      └──────────────┘      └──────┬──────┘  │
│                                                  │         │
│                                                  ▼         │
│                                           ☁️  AWS S3        │
│                                           🗂️ goexpert-bucket│
└─────────────────────────────────────────────────────────────┘
```

### O Gerador de Arquivos (generator)

O `generator` é um programa simples que cria arquivos de teste em `./tmp/`:

```go
// cmd/generator/main.go
func main() {
    i := 0
    for {  // loop infinito
        f, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", i))
        if err != nil {
            panic(err)
        }
        defer f.Close()
        f.WriteString("Hello, World!")
        i++
    }
}
```

**O que ele faz:** cria `file0.txt`, `file1.txt`, `file2.txt`, ... indefinidamente. É apenas uma ferramenta para ter arquivos para testar o uploader.

> **Atenção:** Este código tem um problema sério que veremos na seção [Problemas no Código Atual](#problemas-no-código-atual).

### O Uploader (uploader)

O `uploader` é o programa principal. Ele:
1. Abre o diretório `./tmp/`
2. Cria canais para controle de concorrência e tratamento de erros
3. Para cada arquivo, dispara uma goroutine para fazer o upload
4. Aguarda todos os uploads concluírem

Vamos analisar bloco a bloco:

**Bloco 1: Abertura do diretório e criação dos canais**
```go
func main() {
    dir, err := os.Open("./tmp")
    if err != nil {
        panic(err)
    }
    defer dir.Close()

    uploadControl := make(chan struct{}, 100)    // semáforo: máximo 100 uploads simultâneos
    errorFileUpload := make(chan string, 10)     // fila de arquivos que falharam
```

**Bloco 2: Goroutine de retry**
```go
    go func() {
        for {
            select {
            case filename := <-errorFileUpload:  // recebe nome do arquivo que falhou
                uploadControl <- struct{}{}       // reserva uma vaga
                wg.Add(1)
                go uploadFile(filename, uploadControl, errorFileUpload)  // tenta de novo
            }
        }
    }()
```

**Bloco 3: Loop principal de upload**
```go
    for {
        files, err := dir.ReadDir(1)   // lê 1 arquivo por vez
        if err != nil {
            if err == io.EOF {         // chegou ao fim do diretório
                break
            }
            fmt.Printf("Error reading directory: %s\n", err)
            continue
        }
        wg.Add(1)
        uploadControl <- struct{}{}    // reserva uma vaga (bloqueia se já tiver 100)
        go uploadFile(files[0].Name(), uploadControl, errorFileUpload)
    }
    wg.Wait()  // aguarda todos os uploads terminarem
}
```

**Função uploadFile:**
```go
func uploadFile(filename string, uploadControl <-chan struct{}, errorFileUpload chan<- string) {
    defer wg.Done()  // sinaliza ao WaitGroup que esta goroutine terminou

    completeFileName := fmt.Sprintf("./tmp/%s", filename)
    f, err := os.Open(completeFileName)
    if err != nil {
        <-uploadControl              // libera a vaga
        errorFileUpload <- filename  // envia para retry
        return
    }
    defer f.Close()

    _, err = s3Client.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(s3Bucket),
        Key:    aws.String(filename),
        Body:   f,               // envia o arquivo como stream
    })
    if err != nil {
        <-uploadControl              // libera a vaga
        errorFileUpload <- filename  // envia para retry
        return
    }
    fmt.Printf("File %s uploaded successfully\n", completeFileName)
    <-uploadControl  // libera a vaga para o próximo upload
}
```

---

## Concorrência no Projeto

Este projeto é um ótimo exemplo de concorrência em Go. Vamos entender cada peça.

### sync.WaitGroup — Coordenando Goroutines

Imagine que você é um gerente de projeto e precisa saber quando **todos** os seus funcionários terminaram o trabalho antes de fechar o escritório. O `sync.WaitGroup` é exatamente esse gerente.

```
Gerente (main)                 Funcionários (goroutines)
──────────────                 ─────────────────────────
wg.Add(1) → "contrato 1"      goroutine 1 começa
wg.Add(1) → "contrato 2"      goroutine 2 começa
wg.Add(1) → "contrato 3"      goroutine 3 começa
wg.Wait() → "aguardando..."
                               goroutine 2 termina → wg.Done() → "dispensei"
                               goroutine 1 termina → wg.Done() → "dispensei"
                               goroutine 3 termina → wg.Done() → "dispensei"
wg.Wait() → "todos prontos!"  (contador zerou, main continua)
```

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)           // avisa que mais 1 goroutine vai rodar
    go func(n int) {
        defer wg.Done() // avisa que esta goroutine terminou
        fmt.Printf("Goroutine %d concluída\n", n)
    }(i)
}

wg.Wait()  // bloqueia aqui até todas chamarem Done()
fmt.Println("Todas as goroutines terminaram!")
```

### Channels Bufferizados como Semáforo

O `uploadControl` com capacidade 100 age como um **semáforo** — um mecanismo clássico de controle de concorrência.

```
uploadControl := make(chan struct{}, 100)
```

Pense em um estacionamento com 100 vagas:

```
🅿️ ESTACIONAMENTO (uploadControl, cap=100)
┌────────────────────────────────────────┐
│ [✓][✓][✓][ ][ ][ ]...[ ][ ][ ][ ][ ] │
│  98 vagas ocupadas   2 vagas livres    │
└────────────────────────────────────────┘

Novo carro chega → uploadControl <- struct{}{} → entra se tiver vaga, espera se não tiver
Carro sai        → <-uploadControl              → libera a vaga para o próximo
```

**Por que `struct{}`?** É o tipo que ocupa **zero bytes** de memória em Go. Quando você não precisa do valor, só do efeito de bloqueio do canal, `struct{}` é a escolha certa.

```go
// ✅ Correto — sem desperdício de memória
uploadControl := make(chan struct{}, 100)
uploadControl <- struct{}{}  // ocupa vaga
<-uploadControl              // libera vaga

// ❌ Funciona mas desperdiça memória
uploadControl := make(chan bool, 100)
uploadControl <- true        // bool ocupa 1 byte desnecessariamente
<-uploadControl
```

### Canal de Erros e Retry Automático

O `errorFileUpload` é um canal onde goroutines que falharam colocam o nome do arquivo para ser tentado novamente:

```
goroutine A falhou → errorFileUpload <- "file42.txt"

                         ┌──────────────────────────┐
goroutine de retry       │ for { select {            │
(roda para sempre) ──────│   case f := <-errorFail:  │
                         │     go uploadFile(f, ...) │
                         │ }}                        │
                         └──────────────────────────┘
```

Este padrão é chamado de **error channel pattern** e é uma forma elegante de separar a lógica de retry do fluxo principal.

### Diagrama do Fluxo de Concorrência

```
main()
  │
  ├─ go retryWorker() ─────────────────────────────────────────┐
  │                                                            │
  ├─ loop: ReadDir(1) ─────────────────────────────────────────┼──────────────────────┐
  │         │                                                  │                      │
  │         ├─ wg.Add(1)                                       │                      │
  │         ├─ uploadControl <- struct{}{}  ←── bloqueia se 100 goroutines ativas     │
  │         └─ go uploadFile() ──────────────┐                 │                      │
  │                                          │                 │                      │
  │                                    sucesso? ──── ✅ <-uploadControl                │
  │                                          │                 │                      │
  │                                    falhou? ───── ❌ <-uploadControl                │
  │                                          │       errorFileUpload <- filename ──────┘
  │                                          └─ wg.Done()      │
  │                                                            │
  └─ wg.Wait() ─────────────────────────── aguarda todos Done() ┘
```

---

## Boas Práticas Presentes no Projeto

### Limitação de Concorrência

O projeto limita a **100 uploads simultâneos** usando o channel bufferizado. Sem esse limite, o programa poderia disparar dezenas de milhares de goroutines ao mesmo tempo, consumindo toda a memória do sistema.

```go
// 💡 Padrão: channel bufferizado como semáforo
uploadControl := make(chan struct{}, 100)

// Antes de cada goroutine:
uploadControl <- struct{}{}  // bloqueia se já tiver 100 goroutines ativas

// Dentro da goroutine (sempre, seja sucesso ou erro):
<-uploadControl  // libera a vaga para a próxima
```

Este padrão é chamado de **bounded goroutine pool** ou **semaphore pattern** e é amplamente usado em produção.

### Retry Automático via Canal Separado

Uploads que falham não são descartados — eles são colocados em um canal e tentados novamente por uma goroutine dedicada. Isso garante que arquivos importantes não sejam perdidos por falhas de rede transientes.

```go
// Goroutine dedicada ao retry — roda em paralelo com o loop principal
go func() {
    for {
        select {
        case filename := <-errorFileUpload:
            uploadControl <- struct{}{}
            wg.Add(1)
            go uploadFile(filename, uploadControl, errorFileUpload)
        }
    }
}()
```

### Shutdown Gracioso com wg.Wait()

O `wg.Wait()` no final do `main()` garante que o programa **não encerre antes de todos os uploads terminarem**. Sem isso, o programa poderia sair no meio de uploads em andamento, corrompendo os dados ou deixando operações incompletas.

```go
wg.Wait()  // aguarda todas as goroutines chamarem wg.Done()
// só aqui o programa termina
```

### Upload por Streaming

O arquivo é enviado ao S3 **como um stream** — sem carregar tudo na memória de uma vez. Isso é crítico para arquivos grandes.

```go
f, err := os.Open(completeFileName)
// ...
_, err = s3Client.PutObject(&s3.PutObjectInput{
    Bucket: aws.String(s3Bucket),
    Key:    aws.String(filename),
    Body:   f,  // io.Reader — lê do disco conforme envia pela rede
})
```

```
❌ Sem streaming (perigoso para arquivos grandes):
   disco → carrega 10 GB na RAM → envia pela rede
   (mata o processo com OOM kill)

✅ Com streaming (correto):
   disco → lê 64 KB → rede → lê mais 64 KB → rede → ...
   (usa memória constante independente do tamanho do arquivo)
```

### Liberação de Recursos com defer

O código usa `defer` para garantir que arquivos e o diretório sejam sempre fechados, mesmo em caso de erro:

```go
dir, err := os.Open("./tmp")
defer dir.Close()  // fecha quando a função retornar — com sucesso ou com erro

f, err := os.Open(completeFileName)
defer f.Close()    // idem para cada arquivo
```

---

## O que as Boas Práticas Evitaram

### Goroutine Explosion

Sem o limitador de concorrência, o código ficaria assim:

```go
// ❌ SEM LIMITADOR — perigoso
for {
    files, _ := dir.ReadDir(1)
    go uploadFile(files[0].Name())  // cria goroutine sem limite
}
// Se tiver 1.000.000 de arquivos → 1.000.000 de goroutines simultâneas
// Cada goroutine usa ~2 KB de stack → 2 GB de RAM só de goroutines
// Resultado: OOM Kill ou sistema travado
```

Com o limitador de concorrência, o sistema processa no máximo 100 arquivos por vez, independente de quantos arquivos existam no diretório.

### Uploads Perdidos Silenciosamente

Sem o canal de erros, um upload que falha simplesmente seria descartado:

```go
// ❌ SEM RETRY — dados podem ser perdidos
_, err = s3Client.PutObject(...)
if err != nil {
    fmt.Println("Erro!") // imprime e... tchau, arquivo perdido
    return
}
```

Com o `errorFileUpload`, o arquivo entra numa fila de retry e será tentado novamente automaticamente.

### Race Conditions

Sem o `sync.WaitGroup`, o `main()` poderia terminar enquanto goroutines ainda estão em execução:

```go
// ❌ SEM WAITGROUP — comportamento imprevisível
for i := 0; i < 1000; i++ {
    go uploadFile(...)
}
// main() termina aqui — matando todas as goroutines no meio do upload!

// ✅ COM WAITGROUP — correto
for i := 0; i < 1000; i++ {
    wg.Add(1)
    go uploadFile(...)
}
wg.Wait()  // aguarda todas terminarem
```

### Out of Memory por Arquivos Grandes

Sem o upload por streaming, arquivos de 1 GB carregariam 1 GB na RAM antes de enviar. Com 100 uploads simultâneos, isso seria 100 GB de RAM — impossível na maioria dos sistemas. O streaming mantém o uso de memória constante e baixo.

---

## O que Poderia Ser Melhorado

### Credenciais Seguras

O código atual usa credenciais hardcoded (literalmente `"---"` como placeholder):

```go
// ❌ NUNCA faça isso em produção
Credentials: credentials.NewStaticCredentials("AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/...", ""),
```

**Como fazer corretamente:**

**Opção 1: Variáveis de ambiente (mais simples)**
```go
// ✅ O SDK lê automaticamente AWS_ACCESS_KEY_ID e AWS_SECRET_ACCESS_KEY
sess, err := session.NewSession(&aws.Config{
    Region: aws.String(os.Getenv("AWS_REGION")),
})
// Configure no terminal:
// export AWS_ACCESS_KEY_ID="sua-chave"
// export AWS_SECRET_ACCESS_KEY="sua-chave-secreta"
// export AWS_REGION="us-east-1"
```

**Opção 2: IAM Roles (melhor para produção em AWS)**
```go
// ✅ Em instâncias EC2, ECS ou Lambda — sem credenciais no código
// O SDK detecta automaticamente as permissões da Role anexada ao serviço
sess, err := session.NewSession(&aws.Config{
    Region: aws.String("us-east-1"),
})
// Nenhuma credencial explícita — a Role IAM do servidor tem a permissão
```

**Opção 3: AWS Secrets Manager (para ambientes complexos)**
```go
// Busca credenciais de forma segura e com rotação automática
svc := secretsmanager.New(sess)
result, _ := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
    SecretId: aws.String("minha-aplicacao/s3-credentials"),
})
```

### Configuração via Variáveis de Ambiente

Nome do bucket e região hardcoded são uma má prática — o mesmo binário não pode ser usado em ambientes diferentes (dev, staging, prod):

```go
// ❌ Atual — obriga a recompilar para mudar de ambiente
s3Bucket = "goexpert-bucket-exemplo"

// ✅ Melhorado — configurável sem recompilação
func init() {
    bucket := os.Getenv("S3_BUCKET")
    if bucket == "" {
        log.Fatal("variável de ambiente S3_BUCKET não definida")
    }
    s3Bucket = bucket

    region := os.Getenv("AWS_REGION")
    if region == "" {
        region = "us-east-1" // valor padrão razoável
    }

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region),
    })
    // ...
}
```

### Timeout com Context

Operações de rede podem travar indefinidamente. Um `context` com timeout garante que o programa não fique preso esperando para sempre:

```go
// ✅ Com timeout
import "context"

func uploadFile(filename string, ...) {
    // Cancela o upload se demorar mais de 30 segundos
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    _, err = s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
        Bucket: aws.String(s3Bucket),
        Key:    aws.String(filename),
        Body:   f,
    })
    // Se demorar mais de 30s: ctx.Err() == context.DeadlineExceeded
}
```

### Retry com Limite e Backoff Exponencial

O retry atual é potencialmente infinito — se um arquivo nunca conseguir ser enviado, a goroutine de retry ficará tentando para sempre, bloqueando recursos.

```go
// ❌ Atual — retry infinito sem pausa
errorFileUpload <- filename // tenta de novo imediatamente, para sempre

// ✅ Melhorado — retry com limite e backoff exponencial
type uploadTask struct {
    filename string
    attempts int
}

// Na goroutine de retry:
case task := <-errorFileUpload:
    if task.attempts >= 3 {
        log.Printf("FALHA DEFINITIVA: %s após %d tentativas", task.filename, task.attempts)
        continue  // desiste após 3 tentativas
    }
    // Backoff exponencial: espera 1s, 2s, 4s entre tentativas
    wait := time.Duration(1<<task.attempts) * time.Second
    time.Sleep(wait)
    task.attempts++
    go uploadFileWithRetry(task, uploadControl, errorFileUpload)
```

**Backoff exponencial** significa aumentar o tempo de espera progressivamente:
```
Tentativa 1 → espera 1 segundo
Tentativa 2 → espera 2 segundos
Tentativa 3 → espera 4 segundos
Tentativa 4 → desiste
```

Isso evita sobrecarregar a API da AWS quando ela está com problemas.

### Graceful Shutdown com Sinais do OS

Se o usuário apertar `Ctrl+C` durante o upload, o programa termina abruptamente, potencialmente no meio de um envio:

```go
// ✅ Graceful shutdown — aguarda uploads em andamento antes de sair
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // ... inicia uploads normalmente ...

    // Aguarda sinal de parada OU conclusão de todos os uploads
    select {
    case <-ctx.Done():
        fmt.Println("Sinal recebido, aguardando uploads em andamento...")
        wg.Wait()
        fmt.Println("Encerrado com segurança.")
    }
}
```

### Logging Estruturado

O projeto usa `fmt.Printf` para logs, o que dificulta filtragem e análise em produção:

```go
// ❌ Atual — difícil de filtrar e analisar
fmt.Printf("Uploading file %s to bucket %s\n", completeFileName, s3Bucket)

// ✅ Com slog (disponível a partir do Go 1.21, sem dependências externas)
import "log/slog"

slog.Info("iniciando upload",
    "arquivo", completeFileName,
    "bucket", s3Bucket,
)
slog.Error("falha no upload",
    "arquivo", completeFileName,
    "erro", err.Error(),
    "tentativa", attempts,
)
// Saída JSON estruturada:
// {"time":"...","level":"INFO","msg":"iniciando upload","arquivo":"./tmp/file1.txt","bucket":"meu-bucket"}
```

### Multipart Upload para Arquivos Grandes

O `PutObject` tem um limite de **5 GB** por objeto. Para arquivos maiores, é necessário o **Multipart Upload** — que divide o arquivo em partes e as envia em paralelo:

```go
// ✅ Multipart upload — para arquivos > 5 GB ou quando quer upload paralelo
uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
    u.PartSize = 5 * 1024 * 1024  // partes de 5 MB
    u.Concurrency = 5             // 5 partes em paralelo
})

_, err := uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(s3Bucket),
    Key:    aws.String(filename),
    Body:   f,
})
// O SDK cuida automaticamente de dividir e enviar as partes
```

O `s3manager.Uploader` do SDK faz multipart upload automaticamente quando o arquivo é grande, e envia as partes em paralelo — mais rápido que o `PutObject` para arquivos grandes.

### Problemas no Código Atual

O código desta aula tem dois problemas que valem a pena entender:

**Problema 1: Conflito de merge não resolvido**

O arquivo `uploader/main.go` contém marcadores de conflito de git nas linhas 83-91 e 103-111:

```go
// ❌ Este código NÃO COMPILA
<<<<<<< HEAD
errorFileUpload <- filename
=======
<<<<<<< HEAD
errorFileUpload <- completeFileName
=======
errorFileUpload <- filename
>>>>>>> d86f2c7
>>>>>>> 07cc150
```

O conflito é sobre qual valor enviar para o canal de erros: o `filename` (só o nome) ou o `completeFileName` (caminho completo). A resposta correta é `filename`, porque a função `uploadFile` adiciona o prefixo `./tmp/` internamente.

**Problema 2: defer dentro de loop infinito no generator**

```go
// ❌ generator/main.go — problema sutil
for {
    f, err := os.Create(...)
    defer f.Close()  // ← este defer só executa quando a FUNÇÃO retorna
    // Como é um loop infinito, a função NUNCA retorna
    // → todos os arquivos ficam abertos até o programa ser morto
    // → File Descriptor Leak — o OS tem limite de arquivos abertos
}

// ✅ Correto — fecha explicitamente na mesma iteração
for {
    f, err := os.Create(...)
    f.WriteString("Hello, World!")
    f.Close()  // fecha imediatamente, sem defer
}
```

---

## Principais Problemas ao Trabalhar com S3

### 1. Gerenciamento de Credenciais

**Problema:** Credenciais AWS têm acesso a serviços críticos. Se vazarem, um invasor pode ler/deletar todos os seus arquivos, ou pior, criar serviços na sua conta gerando uma fatura enorme.

**Como abordar:**

| Abordagem | Segurança | Complexidade | Quando usar |
|-----------|-----------|--------------|-------------|
| Hardcoded no código | ❌ Péssima | Trivial | Nunca em produção |
| `.env` local (nunca no git) | ⚠️ Aceitável | Baixa | Desenvolvimento local |
| Variáveis de ambiente no CI/CD | ✅ Boa | Média | Build/deploy pipelines |
| IAM Role (EC2, ECS, Lambda) | ✅✅ Excelente | Baixa | Produção na AWS |
| AWS Secrets Manager | ✅✅ Excelente | Alta | Multi-nuvem, rotação automática |

**Regra de ouro:** Nunca commite credenciais no git. Adicione ao `.gitignore`:
```
.env
*.pem
*credentials*
```

### 2. Rate Limiting e Throttling

**Problema:** O S3 tem limites de requisições por segundo por prefix. Se você fizer muitas requisições muito rápido, começa a receber `503 SlowDown` errors.

**Como abordar:**
- Implementar retry com **exponential backoff** (já discutido acima)
- O AWS SDK v1 tem retry automático configurável:
```go
sess, _ := session.NewSession(&aws.Config{
    MaxRetries: aws.Int(5),  // tenta até 5 vezes automaticamente
})
```
- Para casos de alta carga: distribuir os objetos em diferentes prefixes (ver seção de naming)

### 3. Arquivos Grandes e Multipart Upload

**Problema:** O `PutObject` tem limite de 5 GB e fica lento para arquivos grandes por enviar em uma única requisição (sem paralelismo).

**Como abordar:**
- Usar `s3manager.Uploader` para qualquer arquivo > 100 MB (já mostrado acima)
- Para arquivos entre 5 MB e 5 GB: multipart upload com partes paralelas é ~3-5x mais rápido
- Para arquivos > 5 GB: multipart upload é obrigatório

**Limites do S3:**
| Método | Tamanho máximo | Paralelismo |
|--------|---------------|-------------|
| `PutObject` | 5 GB | Não |
| Multipart Upload | 5 TB | Sim |
| Parte de Multipart | 5 GB | — |

### 4. Falhas de Rede e Idempotência

**Problema:** Redes falham. Se uma conexão cair no meio de um upload, o S3 pode ter recebido parte do arquivo — ou nenhuma parte. Você não sabe.

**Como abordar:**
- S3 `PutObject` é **idempotente**: enviar o mesmo arquivo com a mesma Key duas vezes é seguro — o segundo sobrescreve o primeiro sem problema.
- Sempre que um upload falha, tente de novo com o arquivo completo (nunca tente "continuar de onde parou" com PutObject — use Multipart para isso).
- Use **checksums** para verificar integridade:
```go
// O S3 pode verificar MD5 automaticamente
h := md5.New()
io.Copy(h, f)
checksum := base64.StdEncoding.EncodeToString(h.Sum(nil))

_, err = s3Client.PutObject(&s3.PutObjectInput{
    ContentMD5: aws.String(checksum),  // S3 rejeita se o conteúdo não bater
    // ...
})
```

### 5. Custos Inesperados

**Problema:** S3 cobra por armazenamento, requisições e tráfego de saída (egress). Erros de design podem gerar faturas inesperadas.

**Como abordar:**

| Item | Dica de economia |
|------|-----------------|
| **Armazenamento** | Use Storage Classes (Standard → IA → Glacier) para dados antigos |
| **Requisições** | Agrupe operações pequenas — evite muitas requisições de objetos pequenos |
| **Egress** | Dados saindo da AWS para a internet são pagos. CloudFront reduz egress |
| **Multipart incompleto** | Configure Lifecycle Rule para deletar multipart uploads incompletos |
| **Versioning** | S3 Versioning mantém todas as versões — pode acumular custo silenciosamente |

```json
// Lifecycle Rule para limpar multiparts incompletos (AWS Console ou IaC)
{
    "Rules": [{
        "AbortIncompleteMultipartUpload": { "DaysAfterInitiation": 7 },
        "Status": "Enabled"
    }]
}
```

### 6. Permissões IAM e Bucket Policies

**Problema:** Permissões muito amplas (ex: `s3:*` em `*`) criam riscos de segurança. Permissões muito restritas fazem o programa não funcionar.

**Como abordar — princípio do menor privilégio:**

```json
// ✅ Permissão mínima para um uploader — só o necessário
{
    "Version": "2012-10-17",
    "Statement": [{
        "Effect": "Allow",
        "Action": [
            "s3:PutObject",        // pode criar/sobrescrever objetos
            "s3:GetObject"         // pode ler objetos (para verificação)
        ],
        "Resource": "arn:aws:s3:::meu-bucket-producao/*"
        //                                              ^ só neste bucket
    }]
}

// ❌ Perigoso — acesso total a todos os buckets
{
    "Effect": "Allow",
    "Action": "s3:*",
    "Resource": "*"
}
```

**Dica:** Use o [IAM Policy Simulator](https://policysim.aws.amazon.com/) para testar suas políticas antes de aplicar.

### 7. Inconsistência de Região

**Problema:** Tentar acessar um bucket em `sa-east-1` usando um cliente configurado para `us-east-1` resulta em erro `301 Moved Permanently` ou `PermanentRedirect`.

**Como abordar:**
```go
// ✅ Sempre configure a região correta
sess, _ := session.NewSession(&aws.Config{
    Region: aws.String(os.Getenv("AWS_REGION")),  // vem da variável de ambiente
})

// ✅ Ou descubra a região do bucket automaticamente
region, _ := s3manager.GetBucketRegion(context.Background(), sess, bucketName, "us-east-1")
sess.Config.Region = aws.String(region)
```

**Lembre-se:** Dados em diferentes regiões têm custo de transferência entre si. Mantenha o cliente e o bucket na mesma região.

### 8. Nomes de Objetos e Performance

**Problema:** O S3 usa os primeiros caracteres da Key para distribuir os objetos em servidores internos. Se todas as Keys começam com o mesmo prefixo (ex: `2024-01-15-file1`, `2024-01-15-file2`, ...), todos os objetos ficam no mesmo "shard" interno, criando um hotspot.

**Como abordar:**

```
// ❌ Hotspot — todas as keys começam com a mesma data
2024-01-15-arquivo1.jpg
2024-01-15-arquivo2.jpg
2024-01-15-arquivo3.jpg

// ✅ Distribuído — hash no início distribui entre shards
a1b2c3d4-2024-01-15-arquivo1.jpg  (hash do conteúdo ou UUID)
f7e8a9b0-2024-01-15-arquivo2.jpg
3c4d5e6f-2024-01-15-arquivo3.jpg
```

> **Nota:** A partir de julho de 2018, o S3 suporta automaticamente prefixos sequenciais com alto throughput (3.500 PUT/s e 5.500 GET/s por prefixo). Para a maioria dos casos, não é mais necessário randomizar prefixos — mas ainda é uma boa prática para cargas extremamente altas.

---

## Glossário

| Termo | Definição |
|-------|-----------|
| **S3** | Simple Storage Service — serviço de armazenamento de objetos da AWS |
| **Bucket** | Container raiz no S3; nome único global; análogo a um sistema de arquivos |
| **Object** | Arquivo armazenado no S3; composto por dados + metadados + key |
| **Key** | Nome/caminho do objeto dentro do bucket (ex: `imagens/foto.jpg`) |
| **IAM** | Identity and Access Management — sistema de permissões da AWS |
| **IAM Role** | Conjunto de permissões temporárias que podem ser assumidas por serviços AWS |
| **Região** | Localização geográfica dos datacenters AWS (ex: `us-east-1`, `sa-east-1`) |
| **Goroutine** | Thread leve do Go, gerenciada pelo runtime, muito mais barata que threads do OS |
| **Channel** | Mecanismo de comunicação segura entre goroutines em Go |
| **Semáforo** | Padrão de controle de concorrência que limita o número de operações simultâneas |
| **WaitGroup** | Tipo em Go (`sync.WaitGroup`) para aguardar conclusão de múltiplas goroutines |
| **Backoff exponencial** | Estratégia de retry que aumenta o tempo de espera progressivamente (1s, 2s, 4s...) |
| **Streaming** | Transferência de dados em blocos pequenos, sem carregar tudo na memória de uma vez |
| **Multipart Upload** | Método de upload que divide arquivos grandes em partes menores enviadas em paralelo |
| **Idempotente** | Operação que pode ser repetida múltiplas vezes com o mesmo resultado |
| **Egress** | Tráfego de dados saindo da AWS para a internet; gera custo |
| **Throttling** | Limitação de requisições por segundo imposta pela AWS para proteger a infraestrutura |
| **Race condition** | Bug onde o resultado de um programa depende da ordem não-determinística de execução de threads |
| **OOM Kill** | Out of Memory Kill — quando o OS mata um processo por consumo excessivo de memória |
| **Prefix (S3)** | Parte inicial da key usada internamente pelo S3 para distribuição de carga |
| **init()** | Função especial do Go executada automaticamente antes do `main()` |
| **defer** | Palavra-chave do Go que adia a execução de uma função até o retorno da função atual |

---

## Próximos Passos

Agora que você entendeu como fazer upload para o S3 com controle de concorrência em Go, os próximos passos naturais são:

1. **Implementar as melhorias discutidas neste guia**
   - Substituir credenciais hardcoded por variáveis de ambiente
   - Adicionar timeout com `context.WithTimeout`
   - Implementar retry com limite e backoff exponencial

2. **Explorar outras operações do S3**
   - `GetObject` — download de arquivos
   - `DeleteObject` — remoção de objetos
   - `ListObjectsV2` — listar objetos com paginação
   - `CopyObject` — copiar objetos entre buckets

3. **Estudar Storage Classes**
   - S3 Standard (acesso frequente)
   - S3 Standard-IA (acesso infrequente, mais barato)
   - S3 Glacier (arquivamento, muito barato)

4. **S3 com outros serviços AWS**
   - S3 + Lambda (processar arquivos automaticamente no upload)
   - S3 + CloudFront (CDN para servir arquivos com baixa latência)
   - S3 + SQS (processar uploads de forma assíncrona)

5. **Observabilidade**
   - S3 Server Access Logging
   - AWS CloudTrail para auditoria de operações
   - Métricas com Amazon CloudWatch

6. **Segurança avançada**
   - Server-Side Encryption (SSE-S3, SSE-KMS)
   - Bucket Policies para controle de acesso granular
   - S3 Block Public Access settings
   - VPC Endpoints para acesso privado ao S3

---

*Este projeto é parte do curso Go Expert da Full Cycle MBA.*
