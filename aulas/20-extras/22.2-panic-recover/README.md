# 🚨 Panic e Recover em Go — Guia Didático

Go não tem `try/catch`. Isso é uma escolha de design deliberada: na esmagadora maioria dos casos, erros em Go são **valores comuns**, retornados como o último item de uma função (`error`) e tratados no fluxo normal do código — sem desviar a execução para "outro lugar" como uma exceção faria em Java, Python ou JavaScript.

Mas existe uma categoria diferente de problema: aquele que não é um erro "esperado" do domínio, e sim algo que quebrou uma invariante do próprio programa — acessar um índice fora dos limites de um slice, desreferenciar um ponteiro `nil`, dividir por zero. Para esses casos, Go tem um mecanismo separado: **`panic`** e **`recover`**. Este exemplo é minimalista de propósito — dois arquivos `main.go` pequenos — mas cobre o essencial: como um `panic` se propaga, como um `recover` o intercepta, e o padrão mais comum de uso em produção (recover em um middleware HTTP).

---

## 📑 Sumário

- [🤔 O que é Panic e Recover?](#-o-que-é-panic-e-recover)
- [⚔️ Panic/Recover vs Alternativas](#️-panicrecover-vs-alternativas)
- [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
- [▶️ Como Executar](#️-como-executar)
- [⚖️ Trade-offs](#️-trade-offs)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [⚠️ Principais Problemas ao Trabalhar com Panic/Recover](#️-principais-problemas-ao-trabalhar-com-panicrecover)
- [❓ Perguntas de Entrevista](#-perguntas-de-entrevista)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é Panic e Recover?

**Analogia:** pense num prédio com alarme de incêndio. Um erro comum (`error`) é como um funcionário avisando por e-mail "a impressora do 3º andar acabou o papel" — um aviso normal, que segue o fluxo de trabalho de sempre. Um `panic` é como alguém puxando o alarme de incêndio: interrompe tudo, todo mundo para o que está fazendo e começa a sair do prédio andar por andar. O `recover` é o painel de controle do alarme — se alguém chegar lá a tempo e desligá-lo, o prédio não precisa ser evacuado por completo; a operação retoma a partir de um ponto seguro.

```
❌ SEM RECOVER (o alarme toca até o fim)
┌─────────────────────────────────────────┐
│  panic("algo quebrou")                   │
│  → a função atual para                   │
│  → "sobe" para quem chamou ela           │
│  → e para quem chamou quem chamou...     │
│  → até não sobrar mais ninguém           │
│  → o processo inteiro morre (crash)      │
└─────────────────────────────────────────┘

✅ COM RECOVER (alguém desliga o alarme)
┌─────────────────────────────────────────┐
│  panic("algo quebrou")                   │
│  → a função atual para                   │
│  → "sobe" a pilha de chamadas...         │
│  → ...até encontrar um defer com         │
│    recover() esperando por ela           │
│  → recover() captura o valor do panic    │
│  → o processo continua vivo              │
└─────────────────────────────────────────┘
```

Tecnicamente:

- **`panic`** interrompe o fluxo normal de execução da goroutine atual. Em vez de continuar executando a próxima linha, a função para imediatamente, executa qualquer `defer` pendente, e propaga esse mesmo comportamento para quem a chamou — repetindo o processo função a função, "subindo" pela pilha de chamadas. Se ninguém interceptar, o programa inteiro termina com um crash e imprime o stack trace no terminal.
- **`recover`** é uma função embutida (builtin) que, quando chamada **dentro de uma função `defer`**, é capaz de capturar o valor passado para o `panic` e interromper esse processo de propagação. A partir daí, a função onde o `defer` estava registrado retorna normalmente — o resto do programa nem fica sabendo que um panic aconteceu.

Em código, a diferença fica bem clara:

```go
// ❌ Sem recover — o panic sobe até derrubar o processo inteiro
func main() {
	panic("algo quebrou")
	// esta linha nunca executa
}
// saída: panic: algo quebrou
//        goroutine 1 [running]: ...
//        exit status 2
```

```go
// ✅ Com recover — este é o padrão do exemplo 1/main.go
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recuperado:", r)
		}
	}()

	panic("algo quebrou")
	// o processo continua vivo, imprime "recuperado: algo quebrou" e termina normalmente
}
```

---

## ⚔️ Panic/Recover vs Alternativas

| Abordagem | O que faz | Quando usar |
|---|---|---|
| **Retornar `error`** | Devolve o erro como um valor comum, o chamador decide o que fazer | Praticamente sempre — é o padrão idiomático em Go para qualquer condição de erro **esperada** (arquivo não encontrado, validação falhou, timeout de rede) |
| **`panic` sem `recover`** | Interrompe a execução e propaga até derrubar o processo | Falhas de programação genuinamente irrecuperáveis, ou fase de inicialização onde não faz sentido continuar (ex.: configuração obrigatória ausente) |
| **`panic` + `recover`** (este exemplo) | Interrompe a execução, mas alguém "pega no ar" antes do processo morrer | Isolar falhas inesperadas em um limite conhecido (um handler HTTP, um worker de fila) sem derrubar o serviço inteiro |
| **`os.Exit(code)`** | Encerra o processo imediatamente, sem rodar `defer`s, sem chance de `recover` | Scripts CLI que precisam sair com um código de status específico e não têm estado/recursos para limpar |

`panic`/`recover` **não é** um substituto para tratamento de erro com `error` — é uma rede de segurança para o que o código de tratamento de erro normal não previu.

---

## 📚 Conceitos Fundamentais

### 1. `panic()` — Interrompendo o Fluxo Normal

**Analogia:** é puxar o alarme de incêndio — ninguém mais segue a rotina normal a partir dali.

```go
// 1/main.go
func panic2() {
	panic("panic2")
}
```

Quando `panic()` é chamado, a função para imediatamente. Qualquer `defer` registrado nela é executado (na ordem inversa em que foram declarados — LIFO), e então a execução "sobe" para quem chamou essa função, repetindo o mesmo processo: para, roda os `defer`s, sobe mais um nível. Isso se chama **stack unwinding** ("desenrolar a pilha"). Se esse processo chegar até o topo da goroutine sem que ninguém tenha chamado `recover()`, o runtime do Go imprime o stack trace e encerra o programa.

### 2. `defer` como Pré-requisito de `recover()`

**Analogia:** o painel de controle do alarme só funciona se alguém estiver de plantão *antes* do alarme tocar — não adianta correr para ele depois que o prédio já está sendo evacuado.

```go
// 1/main.go
defer func() {
	if r := recover(); r != nil {
		// ...
	}
}()

panic2()
```

`recover()` só tem efeito quando chamado **diretamente dentro de uma função executada via `defer`**. Isso não é coincidência: o `defer` é justamente o mecanismo que continua rodando mesmo quando a função está no meio do processo de "desenrolar" causado por um panic — é o único lugar em que ainda dá tempo de interceptar algo. Chamar `recover()` fora de um `defer` (por exemplo, direto no corpo de uma função, fora de qualquer panic em andamento) sempre retorna `nil` e não tem efeito nenhum.

> 💡 **Detalhe interessante:** o `defer` precisa estar registrado **antes** do `panic` acontecer. No exemplo acima, o `defer` é a primeira linha da `main`, então ele já está "de plantão" quando `panic2()` é chamado logo em seguida.

### 3. `recover()` — Capturando o Valor do Panic

```go
// 1/main.go
if r := recover(); r != nil {
	if r == "panic1" {
		fmt.Println("panic1 recovered")
	}
	if r == "panic2" {
		fmt.Println("panic2 recovered")
	}
}
```

`recover()` é uma função embutida que retorna `any` (`interface{}`):

- Se houver um panic em andamento na goroutine atual, `recover()` **captura** o valor passado para `panic(...)`, interrompe o processo de stack unwinding, e a função onde o `defer` foi declarado passa a retornar normalmente.
- Se **não** houver panic em andamento, `recover()` simplesmente retorna `nil` — por isso o `if r != nil` é o jeito idiomático de checar "houve um panic aqui?".
- Uma vez capturado, o valor original do panic "some": ele não continua subindo a pilha. Se você precisar repropagá-lo, precisa chamar `panic(r)` de novo manualmente (isso se chama **re-panic**, coberto mais abaixo).

No exemplo, o valor de `panic(...)` é uma `string` simples (`"panic1"`, `"panic2"`), o que permite comparar com `==`. Em código de produção, é mais comum usar um `error` (ou até um tipo customizado) em vez de string solta — assim dá para usar `errors.As`/`errors.Is` e ter mais contexto estruturado sobre o que quebrou.

### 4. Recover em Middleware HTTP — o Padrão Mais Usado em Produção

**Analogia:** é como ter um "extintor" embutido em cada porta do prédio, em vez de só um no térreo — cada requisição HTTP é tratada como um incêndio isolado, que não deve derrubar o prédio inteiro (o servidor).

```go
// 2/main.go
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered panic: %v\n", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

Esse é provavelmente o uso mais comum de `recover()` em código Go real: um **middleware** que envolve todos os handlers HTTP. Se qualquer handler causar um panic (um `nil pointer`, um índice inválido, etc.), o middleware o captura, loga o problema e responde `500 Internal Server Error` ao cliente — em vez de derrubar o processo inteiro e tirar do ar **todas** as outras requisições que estavam sendo atendidas em paralelo.

> 💡 **Detalhe interessante:** frameworks web populares em Go (Gin, Echo, Chi) já vêm com um middleware de recover pronto — mas por baixo dos panos, todos fazem exatamente isso: um `defer` + `recover()` envolvendo a chamada do handler.

### 5. Stack Traces com `debug.PrintStack()`

```go
// 2/main.go — linha comentada
// debug.PrintStack()
```

O pacote `runtime/debug` tem a função `debug.PrintStack()`, que imprime o stack trace completo do ponto onde o panic ocorreu — não apenas a mensagem do erro. Em produção, saber *onde* exatamente o panic aconteceu é muito mais valioso do que só saber *que* algo quebrou; por isso essa linha (hoje comentada no exemplo) é o primeiro passo natural para tornar esse middleware realmente útil para debugging — normalmente combinada com um logger estruturado em vez de `log.Printf`.

---

## 🗂️ Estrutura do Projeto

```
22.2-panic-recover/
├── 1/
│   └── main.go   → panic/recover básico: identifica qual das duas panics ocorreu
└── 2/
    └── main.go   → recover aplicado como middleware HTTP, protegendo o servidor inteiro
```

Assim como no exemplo de graceful shutdown, **nenhuma dependência externa** é usada — tudo é biblioteca padrão do Go (`fmt`, `net/http`, `log`).

---

## 🔍 Walkthrough do Código

### Exemplo 1 — `1/main.go`

```go
// 1. Duas funções que representam falhas diferentes
func panic1() {
	panic("panic1")
}

func panic2() {
	panic("panic2")
}

func main() {
	// 2. O defer é registrado ANTES de qualquer panic acontecer —
	//    é isso que permite que recover() funcione mais abaixo
	defer func() {
		if r := recover(); r != nil {
			// 4. Quando o panic acontece, a execução chega até aqui
			//    (dentro do defer), e recover() captura o valor
			if r == "panic1" {
				fmt.Println("panic1 recovered")
			}
			if r == "panic2" {
				fmt.Println("panic2 recovered")
			}
		}
	}()

	// 3. panic2() é chamado — dispara panic("panic2"),
	//    que começa a subir a pilha de chamadas
	panic2()
}
```

Saída esperada: `panic2 recovered` — e o programa termina normalmente (código de saída 0), sem crash.

### Exemplo 2 — `2/main.go`

```go
// 1. O middleware envolve QUALQUER handler passado para ele
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 2. O defer + recover fica "de plantão" para cada requisição
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered panic: %v\n", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// 3. Chama o handler real — se ele der panic, o defer acima captura
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	// rota normal, nunca dá panic
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// rota que sempre dá panic — só para demonstrar o recover em ação
	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	})

	// 4. O mux inteiro é envolvido pelo recoverMiddleware antes de subir o servidor
	log.Println("Listening on", ":3000")
	http.ListenAndServe(":3000", recoverMiddleware(mux))
}
```

O ponto-chave: o `recoverMiddleware` fica **por fora** de todas as rotas, então qualquer handler registrado no `mux` fica automaticamente protegido — nenhum precisa implementar seu próprio `recover`.

---

## ▶️ Como Executar

Nenhuma das duas subpastas tem `go.mod`, então rode com `go run` apontando direto para o arquivo:

```bash
# Exemplo 1 — dentro de aulas/20-extras/22.2-panic-recover/1
go run main.go
# saída: panic2 recovered
```

```bash
# Exemplo 2 — dentro de aulas/20-extras/22.2-panic-recover/2
go run main.go
# saída: Listening on :3000
```

Para observar o recover em ação no exemplo 2:

1. Rode `go run main.go` dentro da pasta `2/`.
2. Em outro terminal, chame a rota normal: `curl http://localhost:3000/` → responde `Hello World` normalmente.
3. Chame a rota que dá panic: `curl http://localhost:3000/panic` → responde `Internal Server Error` (HTTP 500), e no terminal do servidor aparece o log `recovered panic: panic`.
4. Repita o passo 2 — note que o servidor **continua no ar**, mesmo depois do panic da rota `/panic`. Sem o `recoverMiddleware`, o processo inteiro teria morrido no passo 3.

---

## ⚖️ Trade-offs

**✅ Vantagens**

- O processo sobrevive a falhas inesperadas em vez de derrubar tudo (crítico em servidores que atendem múltiplas requisições simultâneas).
- Centraliza o tratamento de erros verdadeiramente inesperados em um único lugar (ex.: um middleware), em vez de espalhar `recover()` por todo o código.
- Permite logar/observar falhas que, de outra forma, apenas apareceriam como o processo caindo sem explicação clara para quem está fora do terminal.

**❌ Desvantagens**

- Usado em excesso, `recover()` pode **esconder bugs reais** — um panic é, na maioria das vezes, sintoma de um erro de programação que deveria ser corrigido, não silenciado indefinidamente.
- Não substitui tratamento de erro idiomático: se uma falha é esperada (arquivo pode não existir, input pode ser inválido), a resposta correta é retornar `error`, não confiar em `panic`/`recover`.
- `recover()` só protege a **goroutine onde o `defer` está registrado** — um panic disparado em uma goroutine diferente (por exemplo, dentro de um `go func() {...}()` iniciado por um handler) **não** é capturado pelo `recoverMiddleware`, e ainda derruba o processo inteiro (veja a seção de problemas comuns abaixo).
- Adiciona uma camada extra de indireção que pode dificultar entender o fluxo real de erros do sistema, se abusado.

---

## 🎯 Casos de Uso Ideais

**Use panic/recover quando:**
- Você precisa de um limite (boundary) que isola falhas inesperadas sem derrubar o processo inteiro — o exemplo clássico é um middleware HTTP;
- Você está processando itens de forma independente (ex.: mensagens de uma fila, uma a uma) e quer que a falha em **um** item não pare o processamento dos demais;
- Você está escrevendo uma biblioteca cuja API pública deve retornar `error`, mas internamente usa `panic` como atalho de controle de fluxo, convertendo para `error` via `recover()` antes de retornar ao chamador (padrão às vezes usado em parsers).

**Evite panic/recover quando:**
- O erro é uma condição **esperada** do domínio (usuário não encontrado, campo obrigatório ausente, timeout de rede) — aí o caminho correto é retornar `error`;
- Você está tentado a usar `panic`/`recover` como substituto de `if err != nil` só para "encurtar" o código — isso quebra a legibilidade que Go valoriza;
- A falha realmente não deveria ter uma forma de recuperação sensata (ex.: configuração crítica ausente na inicialização) — nesse caso, deixar o `panic` derrubar o processo pode ser o comportamento correto.

---

## ⚠️ Principais Problemas ao Trabalhar com Panic/Recover

### 1. Chamar `recover()` fora de uma função deferida

```go
// ❌ Não funciona — recover() aqui sempre retorna nil
func main() {
	panic("algo quebrou")
	r := recover() // nunca executa: a linha acima já interrompeu o fluxo
	fmt.Println(r)
}
```

**Solução:** `recover()` só tem efeito dentro do corpo de uma função chamada via `defer`, registrada **antes** do panic acontecer:

```go
// ✅
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recuperado:", r)
		}
	}()
	panic("algo quebrou")
}
```

### 2. Panic em uma goroutine não é capturado pelo recover de outra

```go
// ❌ O recoverMiddleware NÃO protege este caso
mux.HandleFunc("/panic-async", func(w http.ResponseWriter, r *http.Request) {
	go func() {
		panic("panic dentro da goroutine")
	}()
	w.Write([]byte("ok"))
})
```

Cada goroutine tem sua própria pilha de execução. Um `recover()` registrado na goroutine principal (ou na goroutine da requisição HTTP) **não enxerga** panics de outras goroutines. Se essa goroutine interna panicar sem seu próprio `recover`, o processo inteiro morre — mesmo estando "protegido" por um middleware.

**Solução:** toda goroutine que pode panicar precisa do seu próprio `defer`+`recover()`:

```go
// ✅
go func() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered panic em goroutine: %v\n", r)
		}
	}()
	panic("panic dentro da goroutine")
}()
```

### 3. Engolir o panic silenciosamente

```go
// ❌ O erro desaparece sem deixar rastro
defer func() {
	recover() // captura e descarta — ninguém nunca vai saber que isso aconteceu
}()
```

**Solução:** sempre logar (ou de alguma forma registrar/observar) o valor capturado por `recover()`, como faz `2/main.go` com `log.Printf`. Um panic silenciado é um bug que vai ser muito mais difícil de diagnosticar depois.

### 4. Re-panic sem necessidade (ou a falta dele quando necessário)

```go
// Às vezes você quer logar e MESMO ASSIM deixar o panic continuar subindo,
// por exemplo se não souber tratar aquele tipo específico de falha:
defer func() {
	if r := recover(); r != nil {
		log.Printf("panic não tratado, repropagando: %v\n", r)
		panic(r) // re-panic: continua o stack unwinding a partir daqui
	}
}()
```

**Solução:** decida conscientemente se o `recover()` deve *absorver* o panic (deixar o processo seguir) ou apenas *observá-lo* antes de deixá-lo continuar (`panic(r)` de novo). Recuperar um panic que você não sabe como tratar corretamente pode deixar a aplicação em um estado inconsistente sem que ninguém perceba.

---

## ❓ Perguntas de Entrevista

**O que são `panic` e `recover` em Go, e para que servem?**
`panic` interrompe o fluxo normal de execução de uma goroutine, executando os `defer`s pendentes e propagando essa interrupção função a função até o topo da pilha de chamadas — derrubando o processo se ninguém intervier. `recover` é uma função embutida que, quando chamada dentro de uma função deferida, captura o valor do panic em andamento e interrompe essa propagação, permitindo que o programa continue rodando. Juntos, servem como uma rede de segurança para falhas inesperadas — não como um substituto do tratamento de erro comum com `error`.

**Qual a diferença entre `panic`/`recover` e exceptions (try/catch) de outras linguagens?**
Na superfície, parecem parecidos: ambos interrompem o fluxo normal e permitem capturar o problema em outro ponto do código. A diferença de filosofia é que, em Go, `panic`/`recover` é reservado para falhas **excepcionais e não previstas** — erros esperados do domínio (arquivo não encontrado, validação falhou) devem ser tratados com o padrão idiomático de retornar `error`, não com panic. Em linguagens com exceptions, é comum até erros esperados usarem o mecanismo de exceção; em Go, isso é considerado má prática.

**Por que `recover()` só funciona dentro de uma função chamada via `defer`?**
Porque o `defer` é o único mecanismo que continua sendo executado mesmo enquanto a pilha de chamadas está sendo "desenrolada" por causa de um panic. Se `recover()` fosse chamado em qualquer outro lugar, não haveria garantia de que ele ainda estaria "no caminho" da propagação do panic no momento certo — por isso o design da linguagem exige que ele esteja dentro de uma função deferida, registrada antes do panic acontecer.

**O que acontece se um panic ocorrer dentro de uma goroutine que não tem seu próprio `recover`?**
O processo inteiro é encerrado, mesmo que a goroutine principal (ou outras goroutines) tenham seus próprios blocos de `recover()`. Isso acontece porque `recover()` só enxerga panics da **mesma** goroutine em que está registrado — não existe propagação de panic entre goroutines diferentes. Por isso, qualquer goroutine que possa causar panic (inclusive as disparadas dentro de um handler HTTP já protegido por um middleware de recover) precisa ter seu próprio tratamento.

**Quando faz sentido usar `panic` em vez de retornar um `error`?**
Em situações onde a falha representa uma violação de uma invariante do programa — algo que, se está acontecendo, indica um bug, não uma condição de negócio esperada — ou em fases de inicialização onde não há um estado seguro para continuar (ex.: uma variável de configuração obrigatória ausente, tipicamente usando a função `Must`-prefixada, como em `regexp.MustCompile`). Fora desses casos, retornar `error` é o padrão idiomático e deve ser preferido.

**O que a linha comentada `debug.PrintStack()` adiciona ao exemplo do middleware?**
`log.Printf` sozinho registra apenas o valor passado para `panic(...)` (a mensagem), mas não mostra **onde** no código o panic ocorreu. `debug.PrintStack()`, do pacote `runtime/debug`, imprime o stack trace completo até o ponto do panic — informação essencial para diagnosticar a causa raiz em produção, onde não é possível simplesmente reproduzir o bug ao vivo com um debugger.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **`panic`** | Função embutida que interrompe o fluxo normal de execução da goroutine atual, iniciando o processo de stack unwinding. |
| **`recover`** | Função embutida que, chamada dentro de uma função deferida, captura o valor de um panic em andamento e interrompe sua propagação. |
| **`defer`** | Palavra-chave que agenda a execução de uma função para o momento em que a função atual retornar (inclusive por causa de um panic). |
| **Stack unwinding** | Processo de "desenrolar" a pilha de chamadas função a função, executando os `defer`s pendentes de cada uma, causado por um panic sem recover. |
| **Re-panic** | Chamar `panic(r)` novamente dentro de um bloco que já capturou o valor com `recover()`, repropagando a falha para quem chamou. |
| **Goroutine** | Unidade leve de execução concorrente do Go; cada uma tem sua própria pilha, e panics não atravessam entre goroutines diferentes. |
| **`interface{}` / `any`** | Tipo que pode representar qualquer valor; é o tipo de retorno de `recover()`, já que `panic()` aceita qualquer valor como argumento. |
| **Middleware** | Função que envolve um handler (por exemplo, HTTP), adicionando comportamento antes e/ou depois de chamá-lo — aqui, usado para centralizar o `recover()`. |

---

## 🚀 Próximos Passos

- [ ] **Imediato:** descomente `debug.PrintStack()` em `2/main.go`, dispare o panic via `curl http://localhost:3000/panic` e observe o stack trace completo no terminal.
- [ ] **Imediato:** em `1/main.go`, troque a chamada final de `panic2()` para `panic1()` e confirme que a saída muda para `panic1 recovered`.
- [ ] **Intermediário:** adicione um `go.mod` em cada subpasta (`go mod init exemplo1` / `go mod init exemplo2`) para que os exemplos funcionem como módulos independentes.
- [ ] **Intermediário:** troque as strings soltas (`"panic1"`, `"panic2"`) por um tipo de erro customizado (`type PanicError struct { ... }`) e use `errors.As` no `recover()` para identificar o tipo de forma mais robusta que comparação de string.
- [ ] **Avançado:** em `2/main.go`, adicione uma rota que dispara um panic **dentro de uma goroutine** (`go func() { panic("...") }()`) e confirme que o `recoverMiddleware` **não** protege esse caso — o servidor cai. Depois, implemente um `recover()` próprio dentro dessa goroutine para corrigir o problema.
- [ ] **Avançado:** substitua `log.Printf` por um logger estruturado (ex.: `log/slog` da biblioteca padrão) no `recoverMiddleware`, incluindo o stack trace como campo estruturado, para simular um cenário mais próximo de produção.
