# 🔁 For em Go — Guia Didático

Go é uma das poucas linguagens modernas que tem **apenas uma palavra-chave** para todos os tipos de laço: `for`. Não existe `while`, não existe `do-while` — tudo é `for`, em formatos diferentes. Isso é uma escolha deliberada dos criadores da linguagem: menos palavras-chave, menos formas de fazer a mesma coisa, mais previsibilidade ao ler código de outra pessoa.

Este exemplo, com pouco mais de 20 linhas, concentra dois usos muito comuns do `for` no dia a dia de um desenvolvedor Go: **iterar sobre um número** (`range` sobre um `int`, novidade do Go 1.22) e **disparar trabalho concorrente dentro de um loop**, usando goroutines e um canal para esperar que tudo termine. Se você é iniciante, pense neste README como um mapa: ele explica cada linha do `main.go`, compara as formas de `for` disponíveis, mostra armadilhas clássicas (inclusive uma que já foi corrigida pela própria linguagem) e fecha com perguntas típicas de entrevista.

---

## 📑 Sumário

- [🤔 O que é o `for` em Go?](#-o-que-é-o-for-em-go)
- [⚔️ Formas de `for` em Go](#️-formas-de-for-em-go)
- [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
- [▶️ Como Executar](#️-como-executar)
- [⚖️ Trade-offs](#️-trade-offs)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [⚠️ Principais Problemas ao Trabalhar com For e Goroutines](#️-principais-problemas-ao-trabalhar-com-for-e-goroutines)
- [❓ Perguntas de Entrevista](#-perguntas-de-entrevista)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é o `for` em Go?

**Analogia:** em muitas linguagens, escolher um laço é como escolher entre várias ferramentas parecidas na caixa — um "martelo para pregos grandes", outro "para pregos pequenos". Em Go, existe **um único martelo** que se ajusta ao tamanho do prego dependendo de como você o segura.

```
❌ OUTRAS LINGUAGENS (várias palavras-chave)
┌─────────────────────────────────────────┐
│  for (int i = 0; i < 10; i++) { ... }    │  ← laço contado
│  while (condição) { ... }                │  ← laço condicional
│  do { ... } while (condição);            │  ← executa ao menos 1 vez
│  for (item of lista) { ... }             │  ← laço por elementos
│  → 4 conceitos diferentes para aprender  │
└─────────────────────────────────────────┘

✅ GO (uma palavra-chave, quatro formatos)
┌─────────────────────────────────────────┐
│  for i := 0; i < 10; i++ { ... }         │  ← laço contado
│  for condição { ... }                    │  ← "while" disfarçado
│  for { ... }                             │  ← laço infinito
│  for i, v := range lista { ... }         │  ← laço por elementos
│  → 1 palavra-chave, 4 formas de usá-la   │
└─────────────────────────────────────────┘
```

Tecnicamente, o `for` em Go é a **única** construção de repetição da linguagem. O que muda entre os formatos acima não é a palavra-chave, mas o que vem entre `for` e o bloco `{ }`. Isso simplifica a leitura de código: ao ver `for`, você já sabe que está diante de um laço — só precisa entender qual das quatro variações está sendo usada.

```go
// ❌ Não existe em Go — isso é pseudocódigo de outra linguagem
while (x < 10) {
    x++
}
```

```go
// ✅ Em Go, a mesma ideia usa a palavra-chave for
for x < 10 {
    x++
}
```

---

## ⚔️ Formas de `for` em Go

| Formato | Sintaxe | Quando usar |
|---|---|---|
| **Clássico (3 cláusulas)** | `for i := 0; i < n; i++ { ... }` | Quando você precisa de controle total sobre início, condição de parada e incremento — ex.: percorrer de trás para frente, pular de 2 em 2. |
| **Condicional (estilo `while`)** | `for condição { ... }` | Quando você não sabe quantas iterações vai precisar, só sabe a condição de parada — ex.: ler de um canal até ele fechar. |
| **Infinito** | `for { ... }` | Quando o laço só termina por `break`, `return` ou `os.Exit` — ex.: o loop principal de um servidor ou de um worker. |
| **`range`** | `for i, v := range coleção { ... }` | Quando você quer percorrer os elementos de uma coleção (slice, array, map, string, canal) ou, desde o Go 1.22, um número inteiro — é o formato usado nas duas partes deste exemplo. |

Este exemplo usa exclusivamente a forma `range`, nas suas duas variações mais relevantes para o dia a dia: `range` sobre um `int` (linha 7) e `range` sobre um `slice` (linhas 14 e 21).

---

## 📚 Conceitos Fundamentais

### 1. `range` sobre um `int` — novidade do Go 1.22

**Analogia:** antes, se você queria repetir algo "10 vezes", precisava contar nos dedos manualmente (`for i := 0; i < 10; i++`). Desde o Go 1.22, você pode simplesmente dizer "repita isso dentro do intervalo de 10" — a contagem fica implícita.

```go
// main.go — linhas 6-8
x := 10
for i := range x {
	fmt.Println(i)
}
```

Até o Go 1.21, `range` só funcionava sobre coleções (slices, arrays, maps, strings, canais). A partir do Go 1.22, `range` também aceita um valor inteiro: `for i := range x` é equivalente a `for i := 0; i < x; i++`, mas mais curto e com menos chance de erro (não tem como esquecer o incremento ou errar o operador de comparação). O laço acima imprime `0, 1, 2, ..., 9` — dez números, começando em zero, nunca incluindo o `10`.

> 💡 **Detalhe interessante:** o `i` aqui é opcional. Se você só precisa repetir algo N vezes sem usar o índice, pode escrever `for range x { ... }`, descartando completamente a variável — muito comum em testes de carga ou repetições simples ("faça isso 5 vezes").

### 2. Goroutines Dentro de um Loop

**Analogia:** é como despachar vários garçons ao mesmo tempo, um para cada mesa, em vez de um garçom único atendendo mesa por mesa em sequência.

```go
// main.go — linhas 14-19
for _, v := range values {
	go func() {
		fmt.Println(v)
		done <- true
	}()
}
```

Cada iteração do `range` sobre `values` (`[]string{"a", "b", "c"}`) dispara uma **goroutine** — uma unidade de execução concorrente e independente. As três goroutines começam a rodar "ao mesmo tempo" (na prática, o scheduler do Go decide a ordem real), o que significa que a ordem de impressão de `"a"`, `"b"`, `"c"` **não é garantida** — pode sair `c, a, b`, `b, c, a`, ou qualquer outra combinação a cada execução.

### 3. Captura de Variável de Loop: Antes e Depois do Go 1.22

**Analogia:** imagine três pessoas anotando um recado numa única folha de papel compartilhada, uma de cada vez. Se todo mundo só olha a folha *depois* que a última pessoa escreveu, todas veem o mesmo (último) recado — mesmo que cada uma pensasse estar guardando o seu próprio.

```go
// main.go — linha 15 (dentro do for _, v := range values)
go func() {
	fmt.Println(v) // "v" é capturado por closure
	done <- true
}()
```

Esse é um dos bugs mais famosos e mais perguntados em entrevistas sobre Go. **Antes do Go 1.22**, a variável `v` do `for _, v := range values` era **uma única variável reutilizada** a cada iteração — todas as goroutines fechavam (faziam *closure*) sobre a mesma variável na memória. Como as goroutines rodam de forma assíncrona, era comum que o loop terminasse (e `v` ficasse com o último valor, `"c"`) antes mesmo de qualquer goroutine executar o `fmt.Println(v)` — resultado: `"c", "c", "c"` impresso três vezes, em vez de `"a", "b", "c"`.

> 💡 **Detalhe interessante:** o Go 1.22 (lançado em fevereiro de 2024) mudou esse comportamento na própria linguagem — agora `v` (e o índice, quando existe) é uma **variável nova a cada iteração**. O código deste exemplo só funciona corretamente "de graça" porque roda em Go 1.22 ou superior; em versões anteriores, ele teria o bug descrito acima. Veja a seção [⚠️ Principais Problemas](#️-principais-problemas-ao-trabalhar-com-for-e-goroutines) para o código do "jeito antigo" de corrigir isso manualmente.

### 4. Canal Sem Buffer como Mecanismo de Sincronização

**Analogia:** um canal sem buffer é como um aperto de mão — só acontece quando as duas pessoas estão prontas ao mesmo tempo. Quem estende a mão primeiro (`done <- true`) fica esperando até que a outra pessoa também estenda a mão (`<-done`).

```go
// main.go — linha 12
done := make(chan bool)
```

`make(chan bool)` cria um canal **sem buffer** (capacidade zero). Isso é diferente do canal usado em outros exemplos deste repositório (como `chan os.Signal, 1`, que tem buffer 1). Um canal sem buffer é **bloqueante nos dois lados**: quem envia (`done <- true`) fica parado até que alguém leia; quem lê (`<-done`) fica parado até que alguém envie. Essa característica é usada de propósito aqui — não para guardar dados, mas para **sincronizar** o momento em que cada goroutine termina.

### 5. `for range values { <-done }` — Esperando Todas as Goroutines Terminarem

```go
// main.go — linhas 21-23
for range values {
	<-done
}
```

Esse é o "recibo" de cada goroutine disparada. Como `values` tem 3 elementos, esse laço executa exatamente 3 vezes — uma leitura de `<-done` por goroutine esperada. Cada vez que uma goroutine termina seu trabalho e executa `done <- true`, uma dessas leituras é desbloqueada. Só depois que as 3 goroutines tiverem enviado seu sinal é que esse laço termina, e só então o `main` (e o programa) pode encerrar. Sem esse mecanismo, o `main` poderia terminar (e o processo inteiro morrer) antes mesmo de as goroutines chegarem a imprimir qualquer coisa — goroutines não impedem o programa de encerrar sozinhas.

---

## 🗂️ Estrutura do Projeto

```
22.6-For/
└── main.go       → todo o exemplo: range sobre int + goroutines sincronizadas por canal
```

Assim como os outros exemplos desta pasta, não há nenhuma dependência externa — tudo usa apenas a biblioteca padrão (`fmt`), reforçando que `for`, `range`, goroutines e canais são recursos nativos da linguagem, sem necessidade de nenhum pacote de terceiros.

---

## 🔍 Walkthrough do Código

Seguindo a ordem de execução real do programa:

```go
// 1. Repete de 0 até 9 usando range sobre um int (novidade do Go 1.22)
x := 10
for i := range x {
	fmt.Println(i)
}

// 2. Cria um canal sem buffer, usado só para sincronização (não carrega dado útil)
done := make(chan bool)
values := []string{"a", "b", "c"}

// 3. Dispara uma goroutine para cada elemento de values.
//    Desde o Go 1.22, "v" é uma variável nova a cada iteração — seguro por padrão.
for _, v := range values {
	go func() {
		fmt.Println(v)
		done <- true // avisa que esta goroutine terminou
	}()
}

// 4. Bloqueia a main até receber um sinal de cada uma das 3 goroutines
for range values {
	<-done
}
```

O ponto-chave é que os passos 3 e 4 formam um padrão de **fan-out / fan-in**: o passo 3 "espalha" trabalho em paralelo (fan-out), e o passo 4 "recolhe" a conclusão de cada pedaço de trabalho antes de seguir em frente (fan-in). É o mesmo princípio por trás de ferramentas mais robustas como `sync.WaitGroup` e `errgroup.Group`, só que implementado manualmente com canal puro — útil para entender o que essas ferramentas fazem por baixo dos panos.

---

## ▶️ Como Executar

```bash
# Dentro da pasta aulas/20-extras/22.6-For
go run main.go
```

Saída esperada (a primeira parte é sempre igual; a segunda parte muda de ordem a cada execução):

```
0
1
2
3
4
5
6
7
8
9
c        ← ordem de "a", "b", "c" varia a cada execução
b
a
```

Para observar a concorrência na prática:

1. Rode `go run main.go` várias vezes seguidas e repare que os números `0` a `9` sempre saem na mesma ordem (é um laço sequencial, sem concorrência), mas a ordem de `"a"`, `"b"`, `"c"` muda entre execuções (são goroutines concorrentes, sem ordem garantida).
2. Rode com a flag de detecção de race conditions: `go run -race main.go`. Como este exemplo não compartilha estado mutável entre as goroutines (cada uma só lê seu próprio `v` e escreve no canal), não deve aparecer nenhum aviso — um bom sinal de que o padrão está correto.
3. Experimente comentar temporariamente o laço final (`for range values { <-done }`) e rode várias vezes: você vai notar que, ocasionalmente, o programa termina **sem imprimir `a`, `b` ou `c` nenhuma vez** — prova de que goroutines não seguram sozinhas a execução do `main`.

---

## ⚖️ Trade-offs

**✅ Vantagens do canal manual de sincronização**

- Não exige nenhuma dependência além da biblioteca padrão.
- Deixa explícito, linha a linha, o que está acontecendo — bom para fins didáticos e para depurar o próprio funcionamento de primitivas de concorrência.
- Funciona igualmente bem para sincronizar goroutines que produzem um resultado (poderia enviar dados pelo canal, não só um sinal `bool`).

**❌ Desvantagens do canal manual de sincronização**

- Mais verboso do que `sync.WaitGroup` para o caso simples de "só espere todas terminarem".
- Fácil de errar a contagem: se `for range values { <-done }` executar mais vezes do que goroutines existem, o programa trava (*deadlock*) esperando um sinal que nunca chega.
- Não captura erros das goroutines automaticamente — se uma goroutine falhar (`panic`), esse loop de espera nunca recebe seu sinal, e o programa também trava, sem nenhuma mensagem clara do motivo.
- Para casos mais robustos (parar todas as goroutines se uma falhar, propagar erros), `errgroup.Group` costuma ser uma escolha melhor do que reimplementar isso na mão.

---

## 🎯 Casos de Uso Ideais

**Use `range` sobre `int` quando:**
- Você precisa repetir algo um número fixo de vezes e não precisa de controle refinado sobre o passo (`i += 2`, por exemplo) — nesse caso, o `for` clássico com 3 cláusulas ainda é mais apropriado.
- Você quer gerar índices, IDs sequenciais de teste, ou repetir uma chamada N vezes de forma legível.

**Use goroutines dentro de um `for range` quando:**
- Cada iteração representa uma unidade de trabalho **independente** das demais (chamadas de rede, processamento de itens de uma lista, leitura de múltiplos arquivos).
- O volume de trabalho é pequeno/conhecido o suficiente para não precisar de um *worker pool* com limite de concorrência (para volumes grandes, seria melhor limitar quantas goroutines rodam ao mesmo tempo).

**Prefira `sync.WaitGroup` ou `errgroup` em vez do canal manual quando:**
- Você só precisa esperar todas as goroutines terminarem, sem trocar dados entre elas (`WaitGroup` é mais direto para isso).
- Você precisa propagar o primeiro erro encontrado e cancelar as demais goroutines (`errgroup.Group` com `context`).

---

## ⚠️ Principais Problemas ao Trabalhar com For e Goroutines

### 1. Captura da Variável de Loop em Versões Anteriores ao Go 1.22

```go
// ❌ Bug clássico — só acontece em Go < 1.22
// (o código deste exemplo NÃO tem esse problema, pois roda em Go 1.22+)
for _, v := range values {
	go func() {
		fmt.Println(v) // "v" é a MESMA variável em todas as goroutines
	}()
}
// Saída provável em Go < 1.22: "c" "c" "c" (ou qualquer repetição do último valor)
```

Antes do Go 1.22, todas as goroutines fechavam sobre a mesma variável `v`, que continuava sendo sobrescrita a cada iteração do loop. Como as goroutines rodam de forma assíncrona, era comum que todas "enxergassem" o valor final de `v` no momento em que finalmente executavam.

**Solução (necessária apenas em Go < 1.22):** passar `v` como parâmetro da função, criando uma cópia local para cada goroutine.

```go
// ✅ Correção manual, necessária apenas em Go anterior ao 1.22
for _, v := range values {
	go func(v string) { // "v" agora é uma cópia, exclusiva desta goroutine
		fmt.Println(v)
	}(v)
}
```

### 2. Deadlock por Contagem Errada de Leituras do Canal

```go
// ❌ Espera um sinal a mais do que o número de goroutines disparadas
done := make(chan bool)
values := []string{"a", "b", "c"} // 3 elementos

for _, v := range values {
	go func() {
		fmt.Println(v)
		done <- true
	}()
}

for i := 0; i < len(values)+1; i++ { // ❌ espera 4 sinais, só existem 3
	<-done
}
// Resultado: o programa trava para sempre na última iteração (deadlock)
```

Como o canal não tem buffer e ninguém mais vai enviar um quarto `true`, a última chamada `<-done` fica bloqueada indefinidamente — o runtime do Go até detecta esse tipo de deadlock (quando todas as goroutines estão paradas) e encerra o programa com `fatal error: all goroutines are asleep - deadlock!`.

**Solução:** garantir que o número de leituras seja exatamente igual ao número de goroutines disparadas — por isso, no código original, o laço de espera usa `range values` (a mesma fonte usada para disparar as goroutines), em vez de um número "mágico" digitado à mão.

```go
// ✅ Número de leituras sempre igual ao número de goroutines disparadas
for range values {
	<-done
}
```

### 3. Esquecer que Goroutines Não Seguram Sozinhas a Execução do `main`

```go
// ❌ Dispara goroutines e não espera por nenhuma delas
for _, v := range values {
	go func() {
		fmt.Println(v)
	}()
}
// main() termina aqui — o processo pode morrer antes de qualquer goroutine rodar
```

Diferente de outras linguagens onde threads em segundo plano podem manter o processo "vivo", em Go, assim que a função `main` retorna, o processo é encerrado imediatamente — **mesmo que existam goroutines ainda rodando**. É comum, ao aprender Go, escrever um `go func(){...}()` e não ver nenhuma saída no terminal, justamente por esse motivo.

**Solução:** sempre ter um mecanismo explícito de espera — canal (como neste exemplo), `sync.WaitGroup`, ou `errgroup.Group`.

```go
// ✅ Espera explícita antes de encerrar o main
for range values {
	<-done
}
```

---

## ❓ Perguntas de Entrevista

**Por que Go tem apenas a palavra-chave `for`, sem `while` ou `do-while`?**
É uma escolha de design da linguagem para reduzir a superfície de sintaxe: qualquer forma de repetição — contada, condicional, infinita ou por elementos — usa `for`, mudando apenas o que vem entre a palavra-chave e o bloco `{ }`. Isso torna o código mais previsível de ler, já que existe só um jeito de reconhecer um laço.

**O que mudou no `range` a partir do Go 1.22?**
Duas coisas relevantes: (1) `range` passou a aceitar um valor inteiro diretamente (`for i := range 10`), equivalente a `for i := 0; i < 10; i++`; (2) a variável de iteração do `range` (em qualquer coleção) passou a ser **recriada a cada iteração**, em vez de reutilizada — o que corrigiu, na própria linguagem, o bug histórico de captura de variável em closures/goroutines.

**Explique o bug clássico de "captura de variável de loop" em goroutines. Ele ainda existe?**
Antes do Go 1.22, a variável declarada no `for ... range` era uma única variável compartilhada, sobrescrita a cada iteração. Goroutines que fechavam sobre essa variável (via closure) frequentemente viam o valor da última iteração, em vez do valor da iteração em que foram criadas — porque a leitura da goroutine acontecia depois que o loop já tinha avançado. Desde o Go 1.22, cada iteração cria uma variável nova, então esse bug específico não ocorre mais por padrão — mas é importante saber identificá-lo em código legado ou em código compilado com uma versão anterior.

**Qual a diferença entre um canal com buffer e um canal sem buffer, e por que este exemplo usa um sem buffer?**
Um canal sem buffer (`make(chan bool)`) bloqueia tanto quem envia quanto quem recebe até que os dois lados estejam prontos ao mesmo tempo — funciona como um ponto de encontro (*rendezvous*). Um canal com buffer (`make(chan bool, N)`) permite até N envios sem que haja alguém lendo no momento. Neste exemplo, o canal sem buffer é proposital: o objetivo não é armazenar dados, é **sincronizar** o exato momento em que cada goroutine sinaliza sua conclusão.

**Por que o `main` precisa esperar explicitamente pelas goroutines, em vez de elas continuarem rodando sozinhas?**
Porque, em Go, o processo inteiro termina assim que a função `main` retorna — independentemente de existirem goroutines ainda em execução. Goroutines não têm vida própria fora do processo que as criou; se `main` não esperar por elas de alguma forma (canal, `WaitGroup`, etc.), o programa pode encerrar antes que elas cheguem a produzir qualquer efeito observável.

**Quando você usaria `sync.WaitGroup` em vez de um canal manual para esperar goroutines terminarem?**
Quando o único objetivo é "esperar todas terminarem", sem precisar trocar nenhum dado entre a goroutine e quem está esperando — `WaitGroup` expressa essa intenção de forma mais direta (`Add`/`Done`/`Wait`) e evita ter que contar manualmente quantas leituras fazer no canal. Um canal continua sendo a escolha certa quando, além de sincronizar, você também precisa transportar um resultado (um valor, um erro) de volta da goroutine.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **Goroutine** | Uma unidade leve de execução concorrente gerenciada pelo runtime do Go, iniciada com a palavra-chave `go`. |
| **`range`** | Construção usada com `for` para iterar sobre coleções (slice, array, map, string, canal) ou, desde o Go 1.22, sobre um número inteiro. |
| **Closure** | Uma função que "fecha" sobre variáveis do escopo em que foi definida, mantendo acesso a elas mesmo depois que esse escopo termina. |
| **Canal (`chan`)** | Estrutura da biblioteca padrão usada para comunicação e sincronização entre goroutines; pode ter ou não um buffer interno. |
| **Canal sem buffer** | Canal de capacidade zero, onde envio e recebimento só acontecem quando as duas pontas estão prontas ao mesmo tempo (bloqueante nos dois lados). |
| **Fan-out / Fan-in** | Padrão de concorrência em que trabalho é distribuído entre várias goroutines (fan-out) e depois seus resultados/conclusões são recolhidos (fan-in). |
| **Deadlock** | Situação em que uma ou mais goroutines ficam bloqueadas esperando umas pelas outras indefinidamente, sem nenhuma conseguir prosseguir. |
| **Race condition (condição de corrida)** | Bug que ocorre quando duas ou mais goroutines acessam o mesmo dado concorrentemente sem sincronização adequada, e o resultado passa a depender da ordem de execução. |
| **`sync.WaitGroup`** | Tipo da biblioteca padrão usado para esperar um grupo de goroutines terminar, sem precisar trocar dados entre elas. |

---

## 🚀 Próximos Passos

- [ ] **Imediato:** rode `go run main.go` várias vezes seguidas e observe que a ordem de `"a"`, `"b"`, `"c"` muda, enquanto a ordem de `0` a `9` nunca muda.
- [ ] **Imediato:** troque `done := make(chan bool)` por `done := make(chan bool, 3)` (canal com buffer) e reflita sobre o que muda no comportamento de bloqueio (dica: as goroutines não vão mais esperar que alguém leia para poder enviar).
- [ ] **Intermediário:** reescreva o exemplo usando `sync.WaitGroup` no lugar do canal `done`, e compare a legibilidade das duas versões.
- [ ] **Intermediário:** modifique o código para que cada goroutine envie, além do sinal de conclusão, o valor processado (ex.: a string em maiúsculas) de volta por um canal com buffer, simulando um pipeline produtor-consumidor.
- [ ] **Avançado:** reescreva o exemplo usando `golang.org/x/sync/errgroup`, fazendo uma das goroutines retornar um erro proposital, e observe como o `errgroup.Group` propaga esse erro para quem chamou `Wait()`.
- [ ] **Avançado:** rode `go run -race main.go` neste exemplo e, em seguida, provoque uma race condition de propósito (por exemplo, todas as goroutines incrementando uma mesma variável `int` sem `sync.Mutex` ou `atomic`) para ver o detector de race conditions em ação.
