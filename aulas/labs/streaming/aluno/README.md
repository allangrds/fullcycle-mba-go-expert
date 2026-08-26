# 📡 Streaming em Go — Guia Didático Completo

Este laboratório mostra, na prática, o que é **streaming de dados** usando apenas a biblioteca padrão do Go: um **socket TCP cru**, sem nenhum framework de mensageria (nada de Kafka, RabbitMQ ou gRPC aqui). O objetivo é entender a mecânica mais fundamental de todas: como dois programas trocam uma sequência de bytes através da rede, aos poucos, em vez de "tudo de uma vez".

Se você nunca ouviu falar de "streaming" fora do contexto de Netflix/Spotify, não se preocupe: cada conceito aqui é explicado primeiro com uma analogia do dia a dia, depois com a definição técnica, e por fim com o trecho de código real deste projeto que o implementa. Ao final deste README você vai entender por que TCP é "só um cano de bytes", o que é "length-prefix framing", por que `io.Copy`/`io.CopyN` são o coração do streaming em Go, e como um servidor atende várias conexões ao mesmo tempo usando goroutines.

Este lab tem duas partes que conversam entre si:

- **`server/server.go`** — um servidor TCP que fica escutando na porta `7777`, aceita conexões e lê os dados que chegam.
- **`client/client.go`** — um cliente que se conecta ao servidor e envia um arquivo (simulado com bytes aleatórios) em streaming.

---

## 📑 Sumário

1. [🤔 O que é Streaming e por que ele importa?](#-o-que-é-streaming-e-por-que-ele-importa)
2. [🗂️ Estrutura deste projeto](#️-estrutura-deste-projeto)
3. [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
4. [⚙️ Como funciona o fluxo (passo a passo)](#️-como-funciona-o-fluxo-passo-a-passo)
5. [✅ Boas práticas presentes no projeto](#-boas-práticas-presentes-no-projeto)
6. [⚠️ Problemas e pegadinhas encontrados no próprio código](#️-problemas-e-pegadinhas-encontrados-no-próprio-código)
7. [⚖️ Tradeoffs importantes em streaming](#️-tradeoffs-importantes-em-streaming)
8. [🔧 Como rodar o projeto](#-como-rodar-o-projeto)
9. [📖 Glossário](#-glossário)
10. [💼 Perguntas de Entrevista Respondidas](#-perguntas-de-entrevista-respondidas)
11. [🚀 Próximos passos / exercícios sugeridos](#-próximos-passos--exercícios-sugeridos)

---

## 🤔 O que é Streaming e por que ele importa?

Imagine que você precisa mandar uma mudança de casa inteira para outra cidade. Existem duas formas de fazer isso:

1. **Colocar tudo em um caminhão só**, esperar ele encher, e só então despachar — nada chega do outro lado até o caminhão inteiro estar pronto e percorrer a viagem completa.
2. **Usar uma esteira contínua** que vai carregando caixa por caixa, o tempo todo, e do outro lado alguém já vai recebendo e organizando as caixas conforme elas chegam — sem esperar a mudança inteira estar "pronta".

**Streaming** é a segunda abordagem, só que para dados: em vez de montar a mensagem inteira na memória e só depois enviá-la de uma vez (**buffering completo**), você envia (e processa) os dados **em pedaços, continuamente**, conforme eles ficam disponíveis. O receptor também processa em pedaços, sem precisar esperar "tudo" chegar para começar a fazer algo útil.

Isso importa porque:

- **Memória**: se você tentasse carregar um arquivo de 10GB inteiro na RAM antes de enviar, precisaria de 10GB livres só para isso. Com streaming, você processa em blocos pequenos (ex: 32KB por vez), então o consumo de memória fica praticamente constante, não importa o tamanho do arquivo.
- **Latência**: o receptor pode começar a processar (ex: tocar o vídeo, gravar no disco) antes que a transferência inteira termine.
- **Fluidez natural do TCP**: como você vai ver a seguir, o TCP já funciona como um fluxo contínuo de bytes — streaming é, na verdade, a forma "nativa" de se comunicar em rede; buffering completo é que é o passo extra (e custoso) que você adiciona por cima.

Este lab demonstra o caso mais simples possível de streaming: enviar um arquivo do cliente para o servidor através de um socket TCP, lendo e escrevendo em blocos, sem carregar a rede inteira de uma vez em uma única chamada.

---

## 🗂️ Estrutura deste projeto

```
aulas/labs/streaming/aluno/
├── client/
│   └── client.go   # Conecta ao servidor e envia um payload em streaming
└── server/
    └── server.go   # Escuta conexões TCP e recebe o payload em streaming
```

São dois programas Go independentes (`package main` cada um). Para testá-los, você roda o `server` primeiro (ele fica esperando conexões) e depois o `client` (ele conecta, manda os dados e termina).

---

## 📚 Conceitos Fundamentais

### 1. TCP é só um cano de bytes (não tem "mensagens")

Uma pegadinha comum para quem está começando: o TCP **não sabe o que é uma "mensagem"**. Ele garante que os bytes cheguem **na ordem certa** e **sem perda**, mas não sabe onde uma mensagem termina e outra começa — é literalmente um fluxo (stream) contínuo de bytes, como um cano d'água. Se você mandar `"OI"` e depois `"TCHAU"`, o outro lado pode receber `"OITCHAU"` de uma vez, ou `"OI"` e `"TCHAU"` separados, ou até `"OIT"` e `"CHAU"` — depende de como a rede fragmentou os pacotes.

Isso significa que **é responsabilidade da aplicação** decidir onde uma mensagem começa e termina. É exatamente o problema que o próximo conceito resolve.

### 2. Length-prefixed framing (o prefixo de tamanho)

Este projeto resolve o problema acima com uma técnica clássica chamada **length-prefix framing**: antes de mandar o conteúdo em si, o remetente manda **o tamanho** desse conteúdo. Assim, o receptor sabe exatamente quantos bytes ler para formar "uma mensagem completa".

No cliente, o tamanho é escrito primeiro, como um `int64`:

```go
// client/client.go
binary.Write(conn, binary.LittleEndian, int64(len(file)))

qtsBytes, err := io.CopyN(conn, bytes.NewReader(file), int64(len(file)))
```

No servidor, a leitura espelha exatamente essa ordem — primeiro lê o tamanho, depois lê exatamente aquele número de bytes:

```go
// server/server.go
var size int64
binary.Read(conn, binary.LittleEndian, &size)

qtdBytes, err := io.CopyN(buf, conn, size)
```

Sem esse prefixo, o servidor não teria como saber quando parar de ler — ele ficaria lendo bytes do socket indefinidamente, sem saber se o arquivo já terminou ou se ainda vem mais dado por aí. Esse é o mesmo princípio usado (com mais sofisticação) por protocolos como HTTP/2 e gRPC, que também usam frames com tamanho conhecido.

### 3. `encoding/binary`: convertendo números em bytes (e vice-versa)

Um `int64` na memória do Go ocupa 8 bytes, mas a **ordem** em que esses 8 bytes são escritos importa — isso se chama **endianness** (ordem de bytes). O pacote `encoding/binary` cuida dessa conversão:

```go
binary.Write(conn, binary.LittleEndian, int64(len(file))) // escreve o int64 como 8 bytes, do menos significativo pro mais significativo
binary.Read(conn, binary.LittleEndian, &size)              // lê 8 bytes e reconstrói o int64
```

O importante aqui é que **as duas pontas precisam concordar** na ordem de bytes usada (`LittleEndian` neste caso) — se o cliente escrevesse em `BigEndian` e o servidor lesse como `LittleEndian`, o número reconstruído sairia completamente errado.

### 4. `io.Copy` / `io.CopyN`: o coração do streaming em Go

Esta é a parte mais importante do lab. `io.CopyN(dst, src, n)` copia exatamente `n` bytes de um `io.Reader` (`src`) para um `io.Writer` (`dst`) — mas **não faz isso de uma vez só**. Internamente, ele lê em blocos (por padrão, um buffer de 32KB) e vai escrevendo cada bloco lido, em um loop, até completar os `n` bytes:

```go
// client/client.go — envia em streaming, não em uma única Write gigante
qtsBytes, err := io.CopyN(conn, bytes.NewReader(file), int64(len(file)))
```

```go
// server/server.go — recebe em streaming, não em um único Read gigante
qtdBytes, err := io.CopyN(buf, conn, size)
```

Compare com a alternativa "sem streaming", que seria carregar tudo em um `[]byte` só (ex: com `io.ReadAll`) antes de fazer qualquer coisa com o dado — isso funcionaria para arquivos pequenos, mas escalaria muito mal para arquivos grandes, porque exigiria memória suficiente para conter o arquivo inteiro de uma vez.

> 💡 Qualquer tipo em Go que implemente as interfaces `io.Reader` e `io.Writer` pode ser usado com `io.Copy`/`io.CopyN` — um arquivo, um socket TCP (`net.Conn`), um `bytes.Buffer`, a entrada padrão... É essa composabilidade que torna o modelo de streaming do Go tão poderoso.

### 5. Goroutine por conexão no servidor

O servidor precisa conseguir atender **várias conexões ao mesmo tempo** sem que uma trave a outra. A solução usada aqui é o padrão clássico de **goroutine por conexão**:

```go
// server/server.go
func (ss *StreamServer) ConnectAndReadFile() {
	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go ss.Process(conn)
	}
}
```

O `for` fica bloqueado em `ln.Accept()` esperando a próxima conexão chegar. Assim que uma chega, em vez de processá-la ali mesmo (o que bloquearia a próxima conexão de ser aceita), o servidor dispara `go ss.Process(conn)` — uma nova goroutine cuida daquela conexão específica, enquanto o loop principal já volta a esperar a próxima. Como goroutines são extremamente leves (poucos KB de stack inicial) comparadas a threads de sistema operacional, esse padrão escala para milhares de conexões simultâneas sem problema.

---

## ⚙️ Como funciona o fluxo (passo a passo)

1. **Servidor inicia**: `go run server/server.go` chama `ConnectAndReadFile()`, que abre um listener TCP na porta `7777` e fica bloqueado esperando conexões.
2. **Cliente gera um payload**: `go run client/client.go` cria um `[]byte` preenchido com dados aleatórios (`crypto/rand`), simulando o conteúdo de um "arquivo".
3. **Cliente conecta**: `net.Dial("tcp", ":7777")` abre uma conexão TCP com o servidor.
4. **Servidor aceita a conexão**: `ln.Accept()` retorna, e uma nova goroutine (`ss.Process(conn)`) é disparada para cuidar dela.
5. **Cliente escreve o prefixo de tamanho**: `binary.Write` manda 8 bytes representando `len(file)` como `int64`.
6. **Servidor lê o prefixo de tamanho**: `binary.Read` reconstrói esse `int64` a partir dos 8 bytes recebidos, guardando em `size`.
7. **Cliente transmite o conteúdo em streaming**: `io.CopyN(conn, bytes.NewReader(file), int64(len(file)))` envia os bytes do payload em blocos.
8. **Servidor recebe o conteúdo em streaming**: `io.CopyN(buf, conn, size)` lê exatamente `size` bytes do socket, em blocos, e acumula em um `bytes.Buffer`.
9. **Ambos os lados imprimem quantos bytes foram transferidos** (`fmt.Println`), confirmando que o número enviado bate com o número recebido.
10. **A goroutine do servidor termina** (o `Process` dá `break` após processar um frame) e o cliente também encerra, fechando a conexão.

---

## ✅ Boas práticas presentes no projeto

- **Uso de `io.CopyN` em vez de carregar tudo com `io.ReadAll`**: tanto o envio quanto o recebimento são feitos em streaming/blocos, não em uma única operação monolítica — é o padrão certo para lidar com payloads de tamanho variável (e potencialmente grande) sem estourar memória.
- **Length-prefix framing**: em vez de tentar adivinhar onde a mensagem termina, o protocolo é explícito sobre o tamanho esperado, tornando a leitura no servidor determinística.
- **Goroutine por conexão**: o servidor não bloqueia novas conexões enquanto processa uma conexão existente, graças ao `go ss.Process(conn)`.
- **Separação clara client/server em pacotes distintos**, cada um com sua própria responsabilidade e seu próprio `main`.

## ⚠️ Problemas e pegadinhas encontrados no próprio código

**1. Buffer gigantesco no cliente (~2 TB!)**

```go
// ❌ client/client.go
file := make([]byte, 2048000000000)
```

`2048000000000` bytes equivalem a quase **2 terabytes**. Na prática, essa alocação provavelmente vai falhar ou travar a máquina antes mesmo de chegar à parte de rede — é quase certamente um erro de digitação (faltam ou sobram zeros). Para testar o lab de verdade, vale trocar por um valor realista, por exemplo:

```go
// ✅ um tamanho de teste razoável (2 MB, por exemplo)
file := make([]byte, 2_048_000)
```

Esse bug é, na verdade, um ótimo lembrete do motivo pelo qual streaming importa: mesmo com um valor "realista" de alguns GB, tentar colocar tudo em um único `[]byte` antes de enviar já seria um desperdício de memória — o ideal seria ler o arquivo de origem (ex: de disco) também em streaming, em vez de materializá-lo inteiro em RAM antes de começar a enviar.

**2. Erros de `binary.Write`/`binary.Read` sendo ignorados**

```go
// ❌ client/client.go e server/server.go
binary.Write(conn, binary.LittleEndian, int64(len(file))) // erro descartado
binary.Read(conn, binary.LittleEndian, &size)              // erro descartado
```

Se a conexão cair exatamente durante a escrita/leitura do prefixo, esse erro é silenciosamente ignorado, e o código segue adiante com um `size` possivelmente inválido (zero, ou lixo de memória), o que pode causar um `io.CopyN` que trava esperando bytes que nunca vão chegar. O correto é sempre checar:

```go
// ✅
if err := binary.Write(conn, binary.LittleEndian, int64(len(file))); err != nil {
	panic(err)
}
```

**3. O servidor só lê um frame por conexão, apesar do `for`**

```go
// server/server.go
func (ss *StreamServer) Process(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)

		qtdBytes, err := io.CopyN(buf, conn, size)
		// ...
		break // <- sempre sai do loop após o primeiro frame!
	}
}
```

O `for` sugere a intenção de processar **múltiplos frames** na mesma conexão (por exemplo, o cliente enviando vários arquivos, um atrás do outro, sem reconectar), mas o `break` no final garante que só o primeiro frame é lido antes de a função retornar (e a goroutine morrer). Se essa era mesmo a intenção, o `break` deveria sair apenas quando a conexão fechar (`err == io.EOF`), não incondicionalmente — veja os [próximos passos](#-próximos-passos--exercícios-sugeridos).

---

## ⚖️ Tradeoffs importantes em streaming

| Decisão | Opção A | Opção B | Quando escolher cada uma |
|---|---|---|---|
| **Streaming vs. buffering completo** | Processar em blocos conforme os dados chegam (`io.Copy`/`io.CopyN`) | Carregar tudo em memória antes de processar (`io.ReadAll`) | Streaming é essencial para arquivos grandes ou de tamanho desconhecido, e reduz o pico de uso de memória. Buffering completo é mais simples de programar e aceitável quando você sabe que o payload é pequeno e cabe tranquilamente em memória. |
| **TCP cru vs. protocolo de mais alto nível** | Implementar seu próprio framing sobre TCP (como este lab) | Usar HTTP, gRPC, Kafka, RabbitMQ, etc. | TCP cru dá controle total e overhead mínimo, mas você reimplementa framing, reconexão, autenticação, compressão etc. do zero. Protocolos de mais alto nível já resolvem esses problemas (e são testados em produção por milhões de sistemas), ao custo de mais dependências e menos controle fino. |
| **Length-prefix framing vs. delimitador** | Prefixar cada mensagem com seu tamanho (como aqui) | Usar um caractere/sequência delimitadora (ex: `\n`, `\0`) para marcar o fim da mensagem | Length-prefix é mais robusto para dados binários (que podem conter qualquer byte, inclusive o delimitador escolhido) e permite alocar o buffer de leitura exato. Delimitador é mais simples de debugar visualmente (ex: em texto/JSON linha a linha), mas exige "escapar" o delimitador caso ele apareça dentro dos dados. |
| **Goroutine por conexão vs. worker pool** | Uma goroutine nova para cada conexão aceita (como este lab) | Um pool fixo de goroutines "workers" consumindo conexões de uma fila | Goroutine-per-connection é simples e funciona bem até dezenas de milhares de conexões simultâneas, graças ao custo baixo das goroutines. Um worker pool limita o paralelismo máximo, o que ajuda a proteger o sistema contra um pico de conexões que poderia esgotar CPU/memória/conexões de banco. |
| **Um frame por conexão vs. múltiplos frames (multiplexação)** | Conectar, enviar um frame, desconectar (como o comportamento atual do server, por causa do `break`) | Manter a conexão aberta e enviar vários frames em sequência | Uma conexão por frame é mais simples, mas paga o custo do handshake TCP (three-way handshake) a cada envio. Reaproveitar a conexão para múltiplos frames amortiza esse custo, mas exige um protocolo mais cuidadoso para saber quando um frame termina e o próximo começa (e quando a conexão deve, de fato, ser encerrada). |

---

## 🔧 Como rodar o projeto

Abra dois terminais.

**Terminal 1 — suba o servidor primeiro** (ele precisa estar escutando antes do cliente conectar):

```bash
cd aulas/labs/streaming/aluno
go run server/server.go
```

**Terminal 2 — rode o cliente**:

```bash
cd aulas/labs/streaming/aluno
go run client/client.go
```

> ⚠️ Antes de rodar, ajuste o tamanho do buffer em `client/client.go` (veja a [pegadinha #1](#️-problemas-e-pegadinhas-encontrados-no-próprio-código)) para um valor realista, como `2_048_000` (2 MB), para evitar que a alocação de quase 2 TB trave sua máquina.

Se tudo der certo, o terminal do cliente vai imprimir algo como `Sent 2048000 bytes to server`, e o terminal do servidor vai imprimir os bytes recebidos seguidos de `Received 2048000 bytes from client`.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **Streaming** | Processar/transmitir dados em blocos contínuos, conforme ficam disponíveis, em vez de esperar o conjunto inteiro estar pronto. |
| **Socket** | O "ponto de conexão" que um programa usa para enviar/receber dados pela rede; em Go, representado por `net.Conn` (cliente) e `net.Listener` (servidor escutando). |
| **TCP** | Protocolo de transporte que garante entrega ordenada e sem perdas de um fluxo de bytes entre duas máquinas, mas sem noção de "mensagens" — apenas um fluxo contínuo. |
| **Framing** | A técnica/convenção usada pela aplicação para decidir onde uma mensagem começa e termina dentro do fluxo contínuo de bytes do TCP. |
| **Length-prefixed framing** | Um tipo de framing em que o tamanho da mensagem é enviado antes do conteúdo, para que o receptor saiba exatamente quantos bytes ler. |
| **`io.Reader` / `io.Writer`** | Interfaces do Go que representam, respectivamente, "algo de onde dá para ler bytes" e "algo onde dá para escrever bytes" — a base de toda composabilidade de streaming em Go. |
| **`io.Copy` / `io.CopyN`** | Funções que copiam dados de um `io.Reader` para um `io.Writer` em blocos (streaming), sem carregar tudo em memória de uma vez. |
| **Endianness** | A ordem em que os bytes de um valor multi-byte (como um `int64`) são armazenados/transmitidos — `LittleEndian` (menos significativo primeiro) ou `BigEndian` (mais significativo primeiro). |
| **Goroutine** | Uma unidade de execução concorrente e leve gerenciada pelo runtime do Go (não é uma thread de SO); criada com a palavra-chave `go`. |
| **Goroutine per connection** | Padrão de servidor concorrente em que cada conexão aceita recebe sua própria goroutine dedicada para processamento. |
| **Buffering completo** | Estratégia oposta ao streaming: carregar o conjunto de dados inteiro em memória antes de processá-lo ou enviá-lo. |
| **Backpressure** | Mecanismo (não implementado neste lab simples) pelo qual um receptor mais lento sinaliza ao remetente para desacelerar o envio, evitando sobrecarga. |

---

## 💼 Perguntas de Entrevista Respondidas

**1. O que é streaming de dados e qual a diferença para uma transferência "tudo de uma vez"?**
Streaming é processar ou transmitir dados em blocos contínuos, conforme eles ficam disponíveis, em vez de esperar o conjunto completo estar pronto para só então agir sobre ele. Isso reduz o pico de uso de memória (você nunca precisa manter o conjunto inteiro em RAM) e permite que o consumidor comece a processar antes mesmo do envio terminar — como assistir um vídeo enquanto ele ainda está sendo baixado, em vez de esperar o arquivo inteiro chegar.

**2. Por que o TCP é descrito como um "stream de bytes" e não como um "protocolo de mensagens"?**
Porque o TCP garante entrega ordenada e sem perdas de uma sequência de bytes, mas não tem nenhum conceito nativo de "onde uma mensagem termina e outra começa". Dois `Write`s consecutivos do lado do remetente podem chegar juntos, separados, ou fragmentados de forma diferente no lado do receptor. Cabe à aplicação implementar sua própria forma de delimitar mensagens dentro desse fluxo contínuo — é isso que se chama de "framing".

**3. O que é length-prefixed framing e por que ele é útil?**
É uma técnica de framing em que o remetente envia o tamanho da mensagem antes do conteúdo em si, tipicamente como um inteiro de tamanho fixo. O receptor lê esse tamanho primeiro e então sabe exatamente quantos bytes precisa ler para reconstituir a mensagem completa. É especialmente útil para dados binários, onde não existe um caractere "seguro" para usar como delimitador (qualquer byte pode aparecer nos dados).

**4. Qual a diferença entre `io.Reader`/`io.Writer` e por que essas interfaces são tão centrais em Go?**
`io.Reader` define um método `Read(p []byte) (n int, err error)` — qualquer tipo que saiba "produzir bytes sob demanda" pode implementá-la (um arquivo, um socket, um buffer em memória). `io.Writer` define `Write(p []byte) (n int, err error)` de forma simétrica. Como funções como `io.Copy` operam apenas sobre essas interfaces, qualquer combinação de fonte e destino que as implemente pode ser conectada entre si sem código adicional — é essa composabilidade que torna o modelo de streaming do Go tão flexível.

**5. O que `io.Copy` faz internamente, e por que ele é preferível a ler tudo com `io.ReadAll` antes de escrever?**
`io.Copy` (e sua variante `io.CopyN`, que copia um número exato de bytes) lê do `Reader` de origem em blocos de tamanho fixo (por padrão cerca de 32KB) e escreve cada bloco no `Writer` de destino, em loop, até terminar. Isso mantém o uso de memória praticamente constante, independente do tamanho total dos dados. Usar `io.ReadAll` primeiro exigiria alocar um `[]byte` do tamanho do dado inteiro antes de sequer começar a escrever, o que não escala para arquivos grandes.

**6. O que é o padrão "goroutine por conexão" e quais são seus limites?**
É um padrão de servidor concorrente em que, a cada conexão aceita (`Accept()`), o servidor dispara uma nova goroutine dedicada a processá-la, enquanto a goroutine principal volta imediatamente a aceitar a próxima conexão. Como goroutines custam pouca memória (poucos KB de stack inicial, crescendo sob demanda), esse padrão escala bem até dezenas de milhares de conexões simultâneas. O limite prático aparece quando o número de conexões é tão alto que o overhead agregado de memória/scheduling das goroutines, ou de recursos compartilhados (banco de dados, file descriptors), começa a degradar o sistema — cenário em que padrões como worker pools com limite fixo de concorrência passam a fazer mais sentido.

**7. O que é "endianness" e por que ele importa ao serializar um número para enviar pela rede?**
Endianness é a ordem em que os bytes que compõem um valor multi-byte (como um `int64` de 8 bytes) são armazenados ou transmitidos: `LittleEndian` grava o byte menos significativo primeiro, `BigEndian` grava o mais significativo primeiro. Importa porque, ao serializar um número em bytes para mandar pela rede, o lado que recebe precisa saber exatamente qual convenção foi usada para reconstituir o valor corretamente — se as duas pontas usarem convenções diferentes, o número decodificado sai completamente errado, mesmo que os bytes tenham chegado intactos.

**8. Por que ignorar o valor de erro retornado por uma operação de I/O (como `binary.Write` ou `conn.Read`) é perigoso?**
Porque operações de I/O podem falhar parcialmente ou totalmente a qualquer momento — a conexão pode cair, o timeout pode expirar, o disco pode ficar cheio. Se o erro é ignorado, o código continua executando como se a operação tivesse tido sucesso, usando dados possivelmente incompletos, zerados ou corrompidos (por exemplo, um `size` inválido lido de uma conexão que caiu no meio da leitura), o que pode causar comportamento incorreto silencioso ou até travar a aplicação esperando por bytes que nunca vão chegar.

**9. Quando faz sentido implementar seu próprio protocolo sobre TCP cru, em vez de usar algo como HTTP, gRPC ou uma fila de mensagens (Kafka/RabbitMQ)?**
TCP cru faz sentido quando você precisa do controle mais fino possível sobre performance e overhead — por exemplo, protocolos internos de altíssima performance, ou como exercício de aprendizado para entender a camada por baixo dos frameworks. Na grande maioria dos sistemas em produção, porém, vale mais usar um protocolo de mais alto nível, porque ele já resolve (de forma testada e madura) problemas como negociação de conexão, autenticação, compressão, reconexão, e framing padronizado — problemas que, ao implementar TCP cru, você acaba tendo que resolver na mão, com maior risco de introduzir bugs sutis (como os vistos nas [pegadinhas](#️-problemas-e-pegadinhas-encontrados-no-próprio-código) deste próprio lab).

**10. O que é "backpressure" em um sistema de streaming, e por que este lab simples não precisa lidar com isso?**
Backpressure é o mecanismo pelo qual um consumidor mais lento sinaliza ao produtor para desacelerar o envio de dados, evitando que o consumidor seja sobrecarregado ou que dados se acumulem indefinidamente em buffers. O TCP já oferece uma forma básica disso no nível de transporte (controle de fluxo da própria pilha TCP, que pode fazer o `Write` do remetente bloquear se o receptor não estiver lendo rápido o suficiente), mas sistemas mais sofisticados de streaming (como Kafka) implementam backpressure explícito na camada de aplicação. Este lab não lida com isso explicitamente porque envia um único payload de tamanho conhecido por conexão — não há um fluxo contínuo e prolongado de dados onde a diferença de velocidade entre produtor e consumidor se tornaria um problema visível.

---

## 🚀 Próximos passos / exercícios sugeridos

- **Corrigir o tamanho do buffer no cliente**: troque `2048000000000` por um valor realista (ex: `2_048_000`, 2MB) e confirme que o envio/recebimento funciona de ponta a ponta.
- **Tratar os erros ignorados**: adicione checagem de erro em todas as chamadas de `binary.Write`/`binary.Read`, tanto no client quanto no server.
- **Suportar múltiplos frames por conexão**: altere `Process` para só sair do loop quando a conexão for fechada (`err == io.EOF`) em vez de sempre dar `break` após o primeiro frame — e ajuste o cliente para enviar dois ou três "arquivos" na mesma conexão, um atrás do outro.
- **Adicionar uma confirmação (ACK) do servidor**: depois de receber o payload completo, o servidor pode escrever de volta um byte ou uma pequena mensagem confirmando o recebimento, e o cliente pode aguardar essa confirmação antes de encerrar.
- **Trocar `io.CopyN` fixo por leitura em `bufio`**: experimente envolver `conn` com `bufio.NewReader`/`bufio.NewWriter` e comparar o comportamento — isso introduz um buffer adicional do lado da aplicação, reduzindo o número de syscalls.
- **Medir e limitar a velocidade de transferência**: adicione um contador de bytes/segundo no servidor, para visualizar o streaming acontecendo em tempo real (em vez de só ver o resultado final).
- **Comparar com um `io.ReadAll` no servidor**: implemente uma versão alternativa que lê a conexão inteira com `io.ReadAll` antes de processar, e discuta por que essa abordagem seria pior para arquivos muito grandes.
