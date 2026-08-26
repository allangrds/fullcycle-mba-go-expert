# 🧭 http.ServeMux (Go 1.22+) em Go — Guia Didático

Até a versão 1.21, o `http.ServeMux` da biblioteca padrão do Go era propositalmente simples: ele só sabia comparar o caminho da URL contra prefixos ou caminhos exatos, sem nenhuma noção de método HTTP (`GET`, `POST`...) e sem capturar partes variáveis do caminho (como um `id` no meio da URL). Por isso, qualquer projeto Go que precisasse de rotas como `GET /books/{id}` acabava dependendo de um pacote de terceiros — `chi`, `gorilla/mux`, `httprouter` — só para resolver roteamento. A partir do **Go 1.22**, o próprio `http.ServeMux` da stdlib ganhou três superpoderes: roteamento por método HTTP, wildcards de path e uma regra de precedência bem definida entre padrões concorrentes. Este exemplo é minimalista de propósito — um único arquivo `main.go` com pouco mais de 40 linhas — e reúne, em código real (parte ativo, parte comentado como material de estudo), os principais recursos desse novo roteador nativo.

---

## 📑 Sumário

- [🤔 O que é o novo ServeMux?](#-o-que-é-o-novo-servemux)
- [⚔️ ServeMux vs Alternativas](#️-servemux-vs-alternativas)
- [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
- [▶️ Como Executar](#️-como-executar)
- [⚖️ Trade-offs](#️-trade-offs)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [⚠️ Principais Problemas ao Trabalhar com ServeMux 1.22+](#️-principais-problemas-ao-trabalhar-com-servemux-122)
- [❓ Perguntas de Entrevista](#-perguntas-de-entrevista)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é o novo ServeMux?

**Analogia:** pense em um porteiro de prédio que, até pouco tempo atrás, só sabia dizer "esse elevador vai até o corredor X" — sem se importar se você ia entregar uma encomenda (`POST`) ou só visitar um morador (`GET`), e sem conseguir anotar "o apartamento é o número que vier depois do corredor". A partir do Go 1.22, esse porteiro foi promovido: agora ele pergunta o motivo da visita (o método HTTP), sabe anotar números de apartamento variáveis (`{id}`) e, quando duas instruções concorrem pela mesma pessoa, sabe qual delas é mais específica e deve prevalecer.

```
❌ SEM o ServeMux 1.22 (Go ≤ 1.21)
┌───────────────────────────────────────────────┐
│  mux.HandleFunc("/books/", handler)            │
│  → dentro do handler, você mesmo precisa:      │
│    - checar r.Method ("GET"? "POST"?)          │
│    - fazer parsing manual de r.URL.Path        │
│      (strings.Split, regex, ou lib externa)    │
│  → sem regra clara de precedência entre rotas  │
└───────────────────────────────────────────────┘

✅ COM o ServeMux 1.22+ (este exemplo)
┌───────────────────────────────────────────────┐
│  mux.HandleFunc("GET /books/{id}", handler)    │
│  → método HTTP já filtrado pelo próprio mux    │
│  → r.PathValue("id") entrega o valor pronto    │
│  → padrão mais específico sempre vence          │
│    (regra de precedência definida pela stdlib) │
└───────────────────────────────────────────────┘
```

Tecnicamente, o `http.ServeMux` continua sendo o mesmo tipo (`*http.ServeMux`), criado da mesma forma (`http.NewServeMux()`). O que mudou foi a **sintaxe dos padrões** aceitos por `Handle`/`HandleFunc`: agora um padrão pode opcionalmente começar com um método HTTP seguido de espaço, e o caminho pode conter segmentos entre chaves (`{nome}`) que funcionam como variáveis capturadas. Nenhuma dependência nova é necessária — é 100% biblioteca padrão.

Em código, a diferença fica bem clara:

```go
// ❌ Go ≤ 1.21 — sem method matching nem wildcards nativos
mux.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    id := strings.TrimPrefix(r.URL.Path, "/books/") // parsing manual
    w.Write([]byte("Book " + id))
})
```

```go
// ✅ Go 1.22+ — este é o padrão do exemplo desta pasta
mux.HandleFunc("GET /books/{id}", GetBookHandler)

func GetBookHandler(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // já vem pronto
    w.Write([]byte("Book " + id))
}
```

---

## ⚔️ ServeMux vs Alternativas

| Abordagem | O que faz | Quando usar |
|---|---|---|
| **`http.ServeMux` pré-1.22** | Casa apenas prefixos ou caminhos exatos; sem método HTTP nem wildcards — tudo isso fica por conta do handler | Código legado, ou quando o projeto precisa rodar em versões de Go anteriores à 1.22 |
| **`http.ServeMux` 1.22+** (este exemplo) | Roteamento por método HTTP, wildcards de path (`{nome}`, `{nome...}`), correspondência exata (`{$}`) e precedência automática entre padrões, tudo nativo na stdlib | APIs HTTP de porte pequeno a médio que não precisam de middleware embutido nem de grupos de rotas sofisticados |
| **`chi`** (`github.com/go-chi/chi`) | Router de terceiros com suporte a middleware encadeado, grupos de rotas (`r.Route`), subrotas montáveis | Projetos que crescem e precisam de middleware reutilizável, versionamento de rotas ou organização em módulos de rota |
| **`gorilla/mux`** | Router maduro e muito usado historicamente, com matching por regex em variáveis de path, subrotas e middleware | Projetos legados que já o utilizam, ou quando é necessário matching de path via regex, algo que o ServeMux nativo não oferece |
| **`httprouter`** (`julienschmidt/httprouter`) | Router extremamente otimizado para performance, baseado em árvore de prefixos (radix tree) | Cenários de altíssima performance/baixa latência onde cada alocação e cada nanosegundo de roteamento importam |

O `http.ServeMux` 1.22+ ocupa hoje o espaço que antes só um router de terceiros preenchia: para a maioria das APIs simples, ele elimina a necessidade de uma dependência externa só para roteamento. Ele fica "no meio do caminho" — mais capaz que o ServeMux antigo, porém ainda mais limitado que `chi`/`gorilla/mux` em recursos como middleware embutido e agrupamento de rotas.

---

## 📚 Conceitos Fundamentais

### 1. Roteamento por Método HTTP — `"MÉTODO /caminho"`

**Analogia:** é como colocar duas placas diferentes na mesma porta — "entregas aqui" e "visitantes aqui" — em vez de uma única placa genérica "entrada", deixando quem está do lado de dentro decidir o que fazer com cada um.

```go
// main.go — linha 12 (ativa)
mux.HandleFunc("GET /books/{s}", BooksPrecedenceHandler)
```

Antes do Go 1.22, o padrão passado para `HandleFunc` era só um caminho (`"/books/"`). Agora, o padrão pode começar com um verbo HTTP (`GET`, `POST`, `PUT`, `DELETE`...) seguido de um espaço e do caminho. Se o método da requisição não bater, o `ServeMux` já responde automaticamente com `405 Method Not Allowed` — sem que o handler precise checar `r.Method` manualmente.

> 💡 **Detalhe interessante:** o prefixo de método é **opcional**. Um padrão sem método (ex.: `"/books/"`) continua casando com qualquer verbo HTTP, exatamente como no comportamento pré-1.22 — isso mantém compatibilidade com código antigo.

### 2. Wildcard de Segmento Único — `{nome}` e `r.PathValue`

**Analogia:** é como preencher um formulário com um campo em branco chamado "id" — o valor pode variar a cada preenchimento, mas o formulário sempre sabe onde ir buscar esse valor.

```go
// main.go — linha 8 (comentada) e linhas 22-25
// mux.HandleFunc("GET /books/{id}", GetBookHandler)

func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("Book " + id))
}
```

Um segmento entre chaves, como `{id}`, casa com **exatamente um** segmento do caminho (tudo entre duas barras `/`) e dá um nome a esse valor. Dentro do handler, `r.PathValue("id")` devolve o texto capturado — por exemplo, uma requisição para `/books/42` faria `r.PathValue("id")` retornar `"42"`. Isso substitui o que antes exigia `strings.Split(r.URL.Path, "/")` ou expressões regulares manuais.

### 3. Wildcard "Resto do Caminho" — `{nome...}`

**Analogia:** enquanto `{id}` é como perguntar "qual é o seu apartamento?" (uma resposta curta e única), `{d...}` é como perguntar "qual é o caminho completo até você, incluindo todos os corredores no meio?" — a resposta pode ter vários "segmentos" concatenados.

```go
// main.go — linha 9 (comentada) e linhas 27-30
// mux.HandleFunc("GET /books/dir/{d...}", BooksPathHandler)

func BooksPathHandler(w http.ResponseWriter, r *http.Request) {
	dirpath := r.PathValue("d") // Access captured directory path segments as slice
	fmt.Fprintf(w, "Accessing directory path: %s\n", dirpath)
}
```

Quando o nome da variável termina com `...` (reticências), o wildcard deixa de capturar um único segmento e passa a capturar **tudo o que sobrar** do caminho a partir daquele ponto — incluindo barras adicionais. Uma requisição para `/books/dir/a/b/c` faria `r.PathValue("d")` retornar `"a/b/c"`. É o equivalente, em espírito, ao antigo padrão com prefixo (`"/books/dir/"`), mas com o valor capturado já disponível por nome, sem parsing manual.

> 💡 **Detalhe interessante:** um padrão com `{nome...}` só é válido se esse wildcard estiver **no final** do caminho — não é possível ter segmentos fixos depois dele.

### 4. Correspondência Exata — `{$}`

**Analogia:** é a diferença entre dizer "qualquer pessoa que passar por essa rua" e "só a pessoa que parar exatamente nesse endereço, sem seguir adiante".

```go
// main.go — linha 10 (comentada) e linhas 32-34
// mux.HandleFunc("GET /books/{$}", BooksHandler) // exato

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Books"))
}
```

Por padrão, um caminho terminado em `/` (como `"/books/"`) casa com o caminho exato **e** com qualquer coisa depois dele (`/books/`, `/books/1`, `/books/1/2`...) — esse é o comportamento histórico de "prefixo" do `ServeMux`. O `{$}` no final de um padrão (`"/books/{$}"`) restringe esse casamento para **apenas** o caminho exato `/books/`, sem aceitar nada depois da barra. É a forma explícita de dizer "aqui, e só aqui, sem subcaminhos".

### 5. Precedência entre Padrões Concorrentes

**Analogia:** imagine duas placas na mesma rua: uma diz "toda loja é por aqui" (uma regra geral) e outra diz "a Padaria do João é logo ali" (uma regra específica, só para uma loja). Quando alguém pergunta pela Padaria do João, faz sentido seguir a placa específica, não a genérica — é exatamente essa lógica que o `ServeMux` aplica.

```go
// main.go — linhas 13-14 (comentadas), ilustrando o conceito
// mux.HandleFunc("GET /books/precedence/latest", BooksPrecedenceHandler)
// mux.HandleFunc("GET /books/precedence/{x}", BooksPrecedence2Handler)
```

Quando dois padrões registrados poderiam, em teoria, casar com a mesma URL, o Go 1.22+ define uma regra determinística: o padrão **mais específico** vence. Um segmento literal (`precedence/latest`) é mais específico que um wildcard na mesma posição (`precedence/{x}`), então uma requisição para `/books/precedence/latest` sempre cairia no handler do padrão literal, mesmo que o padrão com wildcard também "pudesse" casar. Essa regra existe justamente para eliminar ambiguidade — sem ela, seria preciso confiar na ordem de registro das rotas (como acontecia em vários routers de terceiros mais simples).

> 💡 **Detalhe interessante:** existem casos em que dois padrões são **genuinamente ambíguos** (nenhum é mais específico que o outro para todas as URLs possíveis) — nesse caso o Go não escolhe "por acaso": ele gera um **panic em tempo de registro** (`mux.HandleFunc`), forçando o desenvolvedor a resolver a ambiguidade explicitamente. Veja mais em [⚠️ Principais Problemas](#️-principais-problemas-ao-trabalhar-com-servemux-122).

### 6. As Rotas Realmente Ativas no Exemplo

```go
// main.go — linhas 15-16 (ativas)
mux.HandleFunc("GET /books/{s}", BooksPrecedenceHandler)
mux.HandleFunc("GET /{s}/latest", BooksPrecedence2Handler)
```

Diferente de todo o resto do arquivo (que está comentado como material de estudo), essas duas linhas são as únicas rotas realmente registradas quando o programa roda. Repare que ambas usam o **mesmo nome** de variável, `{s}`, mas em posições diferentes do caminho (`/books/{s}` vs `/{s}/latest`) — isso é permitido porque o nome do wildcard é local a cada padrão; não há conflito entre eles, e o `ServeMux` os trata como dois padrões completamente independentes. Uma requisição para `/books/42` cai no primeiro; uma para `/xyz/latest` cai no segundo.

---

## 🗂️ Estrutura do Projeto

```
22.5-ServerMux-1.22/
├── main.go     → todo o exemplo: rotas ativas + handlers + variações comentadas para estudo
└── test.http   → requisições de exemplo prontas para disparar contra o servidor (extensão REST Client / VSCode)
```

Este exemplo usa **apenas a biblioteca padrão** (`net/http`, `fmt`) — não há `go.mod` próprio nesta subpasta porque ela é compilada como parte do módulo Go definido em um diretório pai (`aulas/20-extras` ou raiz do repositório). Diferente de exemplos como `22.3-fsnotify` (que depende de um pacote de terceiros), aqui nenhuma dependência externa é necessária: tudo o que este README descreve já vem embutido no Go 1.22+.

---

## 🔍 Walkthrough do Código

Seguindo a ordem de execução real do programa (com uma nota sobre as linhas comentadas):

```go
func main() {
	// 1. Cria o multiplexador de rotas (o "roteador" nativo do Go)
	mux := http.NewServeMux()

	// 2. Linhas comentadas: cada uma é uma variação de conceito para
	//    ser descomentada e testada isoladamente (veja a seção
	//    "Como Executar" para o passo a passo de cada uma).
	// mux.HandleFunc("GET /books/{id}", GetBookHandler)
	// mux.HandleFunc("GET /books/dir/{d...}", BooksPathHandler)
	// mux.HandleFunc("GET /books/{$}", BooksHandler) // exato
	// mux.HandleFunc("GET /books/precedence/latest", BooksPrecedenceHandler)
	// mux.HandleFunc("GET /books/precedence/{x}", BooksPrecedence2Handler)

	// 3. As duas únicas rotas realmente registradas nesta execução
	mux.HandleFunc("GET /books/{s}", BooksPrecedenceHandler)
	mux.HandleFunc("GET /{s}/latest", BooksPrecedence2Handler)

	// 4. Sobe o servidor HTTP na porta 9000, usando o mux como handler raiz
	http.ListenAndServe(":9000", mux)
}
```

O ponto-chave: o arquivo foi montado deliberadamente como um "laboratório" — a lógica de negócio de cada handler é trivial (só escreve uma string na resposta), porque o que importa aqui é o **padrão de rota** registrado em `HandleFunc`, não o corpo do handler. Isso é reforçado pelo fato de handlers como `BooksPrecedenceHandler` e `BooksPrecedence2Handler` terem corpos quase idênticos — a diferença de comportamento entre as rotas vem inteiramente de como o `ServeMux` interpreta cada padrão.

---

## ▶️ Como Executar

```bash
# Dentro da pasta aulas/20-extras/22.5-ServerMux-1.22
go run .
```

O servidor sobe em `http://localhost:9000`. Para explorar cada conceito na prática:

1. **Rotas ativas por padrão** — com o código como está, use o arquivo `test.http` (extensão REST Client do VSCode) ou `curl`:
   ```bash
   curl http://localhost:9000/books/2
   # → Books Precedence
   ```
   Repare que `/books/2` cai em `BooksPrecedenceHandler` porque `{s}` casa com qualquer segmento único — inclusive `"2"`.

2. **Teste o wildcard de segmento único (`{id}`)** — comente a linha 15 (`GET /books/{s}`) e descomente a linha 8 (`GET /books/{id}`), depois:
   ```bash
   go run .
   curl http://localhost:9000/books/42
   # → Book 42
   ```

3. **Teste o wildcard "resto do caminho" (`{d...}`)** — descomente a linha 9 (`GET /books/dir/{d...}`) e rode:
   ```bash
   curl http://localhost:9000/books/dir/a/b/c
   # → Accessing directory path: a/b/c
   ```
   Essa é justamente a requisição já presente em `test.http` (`GET /books/dir/djdjdjd/aa/add/dddre/gfgf`), que só funciona com este handler ativo.

4. **Teste a correspondência exata (`{$}`)** — descomente a linha 10 (`GET /books/{$}`) e compare:
   ```bash
   curl http://localhost:9000/books/
   # → Books (casa exatamente)

   curl http://localhost:9000/books/qualquer-coisa
   # → 404 not found (não casa mais, porque {$} restringe ao caminho exato)
   ```

5. **Teste a precedência entre literal e wildcard** — descomente as linhas 13 e 14 (`GET /books/precedence/latest` e `GET /books/precedence/{x}`), depois:
   ```bash
   curl http://localhost:9000/books/precedence/latest
   # → Books Precedence (o padrão literal, mais específico, vence)

   curl http://localhost:9000/books/precedence/qualquer-outra-coisa
   # → Books Precedence 2 (cai no wildcard, já que não bate com o literal)
   ```

6. **Envie um `POST` para uma rota registrada só com `GET`**, para ver o `405` automático:
   ```bash
   curl -i -X POST http://localhost:9000/books/1
   # → HTTP/1.1 405 Method Not Allowed
   ```
   Essa é exatamente a segunda requisição do `test.http`.

---

## ⚖️ Trade-offs

**✅ Vantagens**

- Zero dependências externas — todo o roteamento por método, wildcards e precedência vêm de `net/http`, já embutidos no Go 1.22+.
- API mais legível e direta: `"GET /books/{id}"` comunica método, caminho e variável capturada em uma única string, sem parsing manual dentro do handler.
- Regra de precedência determinística e verificada em tempo de registro — padrões genuinamente ambíguos causam `panic` cedo (ao subir o programa), em vez de um bug silencioso em produção.
- Migração de um `ServeMux` antigo costuma ser incremental: padrões sem método continuam funcionando como antes.

**❌ Desvantagens**

- Sem suporte nativo a middleware encadeado (`Use()`, por exemplo) — encadear middlewares ainda exige compor `http.Handler`s manualmente.
- Sem agrupamento/nesting de rotas por prefixo (como `r.Route("/books", func(r chi.Router) {...})` do `chi`) — cada padrão é registrado individualmente.
- Wildcards não suportam validação embutida por tipo ou regex (ex.: forçar que `{id}` seja só dígitos) — isso continua sendo responsabilidade do handler.
- Em APIs grandes, a ausência de agrupamento pode deixar o arquivo de rotas repetitivo e mais difícil de organizar do que com um router de terceiros.

---

## 🎯 Casos de Uso Ideais

**Use o `http.ServeMux` 1.22+ quando:**
- O projeto é uma API pequena a média, sem necessidade de middleware complexo encadeado;
- Você quer evitar uma dependência externa só para roteamento;
- As regras de rota são relativamente simples: método HTTP + alguns wildcards de path já resolvem o problema;
- É importante manter o projeto o mais próximo possível da biblioteca padrão (menos superfície de dependências para atualizar/auditar).

**Considere um router de terceiros (`chi`, `gorilla/mux`) quando:**
- A aplicação precisa de middleware reutilizável e componível (autenticação, logging, rate limiting) de forma nativa ao router;
- Rotas precisam ser agrupadas/organizadas em submódulos (ex.: `/api/v1/books`, `/api/v2/books` com middlewares diferentes por versão);
- É necessário matching de path por regex ou validação de tipo do wildcard direto no roteador;
- O projeto já usa um desses routers e migrar traria custo maior que o benefício.

---

## ⚠️ Principais Problemas ao Trabalhar com ServeMux 1.22+

### 1. Padrões Genuinamente Ambíguos Causam `panic` no Registro

```go
// ❌ Nenhum dos dois é mais específico que o outro para toda URL possível
mux.HandleFunc("GET /{a}/edit", handlerA)
mux.HandleFunc("GET /posts/{b}", handlerB)
```

Para uma URL como `/posts/edit`, ambos os padrões poderiam, em teoria, ser "o mais específico" dependendo de como se olha — o Go não tenta adivinhar: ele detecta a ambiguidade **em tempo de registro** (na chamada de `HandleFunc`/`Handle`) e dá `panic`, derrubando o programa imediatamente ao subir.

**Solução:** reestruturar os padrões para que um seja claramente mais específico que o outro, ou usar caminhos que não colidam:

```go
// ✅ Caminhos estruturalmente distintos, sem ambiguidade
mux.HandleFunc("GET /users/{a}/edit", handlerA)
mux.HandleFunc("GET /posts/{b}", handlerB)
```

### 2. Confundir Prefixo (`/books/`) com Exato (`/books/{$}`)

```go
// ❌ Pensando que "/books/" só casa com o caminho exato
mux.HandleFunc("GET /books/", BooksHandler)
// na prática, /books/1, /books/1/2, /books/qualquer-coisa TAMBÉM casam aqui
```

Um padrão terminado em `/` mantém o comportamento histórico de **prefixo**: ele casa com o caminho exato e com qualquer coisa depois dele. Isso pode surpreender quem espera que `"/books/"` case só com `/books/`.

**Solução:** usar `{$}` quando a intenção é realmente restringir ao caminho exato:

```go
// ✅ Casa apenas com /books/ — nada depois da barra
mux.HandleFunc("GET /books/{$}", BooksHandler)
```

### 3. `r.PathValue` Retorna String Vazia para Chave Inexistente

```go
// ❌ Erro de digitação no nome da variável não gera erro em tempo de compilação
mux.HandleFunc("GET /books/{id}", func(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("Id") // "Id" com I maiúsculo — não existe, o padrão usa "id"
	// bookID vem como "" silenciosamente
})
```

`r.PathValue(nome)` não falha nem gera erro se `nome` não corresponder a nenhuma variável do padrão registrado — ele simplesmente retorna string vazia. Isso pode mascarar erros de digitação por um bom tempo, já que o compilador não tem como validar essa string em tempo de build.

**Solução:** manter o nome da variável idêntico entre o padrão e a chamada de `PathValue`, e validar explicitamente quando o valor for obrigatório:

```go
// ✅
id := r.PathValue("id")
if id == "" {
	http.Error(w, "id ausente", http.StatusBadRequest)
	return
}
```

### 4. Esquecer que `{nome...}` Só Pode Ser o Último Segmento

```go
// ❌ Não compila / não é um padrão válido: wildcard "resto" no meio do caminho
mux.HandleFunc("GET /books/{d...}/edit", BooksPathHandler)
```

Como `{nome...}` captura tudo o que resta do caminho, não faz sentido — e não é permitido — ter segmentos fixos depois dele. O registro desse padrão falha (panic) ao subir o programa.

**Solução:** usar `{nome...}` apenas como último segmento do padrão, e se for necessário algo depois, repensar a estrutura da rota:

```go
// ✅
mux.HandleFunc("GET /books/dir/{d...}", BooksPathHandler)
```

---

## ❓ Perguntas de Entrevista

**O que mudou no `http.ServeMux` a partir do Go 1.22?**
O `ServeMux` ganhou três recursos que antes só existiam em routers de terceiros: (1) roteamento por método HTTP, escrevendo o verbo antes do caminho no padrão (`"GET /books/{id}"`); (2) wildcards de path, com `{nome}` capturando um segmento e `{nome...}` capturando o restante do caminho, ambos acessíveis via `r.PathValue`; e (3) uma regra de precedência determinística entre padrões concorrentes, resolvida em tempo de registro. Tudo isso continua sendo biblioteca padrão, sem nenhuma dependência externa.

**Como funciona a regra de precedência quando dois padrões poderiam casar com a mesma URL?**
O Go aplica a regra de "o padrão mais específico vence": um segmento literal é considerado mais específico que um wildcard na mesma posição, e um padrão mais restrito é preferido a um mais genérico. Quando essa comparação é possível de resolver univocamente, o roteamento é determinístico e independe da ordem em que as rotas foram registradas. Quando não é possível decidir (padrões genuinamente ambíguos), o Go não tenta adivinhar — ele gera um `panic` no momento do registro (`HandleFunc`), obrigando o desenvolvedor a resolver o conflito explicitamente antes mesmo do programa rodar.

**Qual a diferença entre `{nome}` e `{nome...}`?**
`{nome}` captura exatamente **um** segmento do caminho (o texto entre duas barras `/`). `{nome...}`, com reticências, captura **todo o restante** do caminho a partir daquele ponto, incluindo barras adicionais — por isso só pode aparecer como o último elemento de um padrão. Por exemplo, `/books/{id}` casa com `/books/42` mas não com `/books/42/comments`, enquanto `/books/{path...}` casaria com ambos, capturando `"42"` e `"42/comments"` respectivamente.

**O que é o `{$}` e por que ele existe?**
`{$}` é um marcador especial usado no final de um padrão para forçar correspondência **exata** de caminho, sem aceitar nada depois. Ele existe porque o comportamento padrão de um padrão terminado em `/` (como `"/books/"`) é casar como **prefixo** — ou seja, `/books/`, `/books/1` e `/books/1/2` casariam todos com esse mesmo padrão, comportamento herdado do `ServeMux` pré-1.22. Quando a intenção é restringir a rota apenas ao caminho exato, `{$}` resolve isso sem exigir lógica adicional dentro do handler.

**Quando compensa usar `chi` ou `gorilla/mux` em vez do `ServeMux` nativo do Go 1.22+?**
Compensa quando a aplicação precisa de recursos que o `ServeMux` nativo não oferece: middleware encadeado de forma nativa ao router (autenticação, logging, recuperação de panic aplicados a grupos de rotas), agrupamento/nesting de rotas por prefixo, ou matching de path via regex/validação de tipo. Para APIs pequenas a médias sem essas necessidades, o `ServeMux` nativo tende a ser suficiente e elimina uma dependência externa.

**Como capturar o valor de uma variável de path dentro do handler?**
Usando `r.PathValue(nome)`, onde `nome` é exatamente o identificador usado entre chaves no padrão registrado (por exemplo, para `"GET /books/{id}"`, o valor correto é `r.PathValue("id")`). O método retorna uma `string`; se o nome não corresponder a nenhuma variável do padrão que casou com a requisição, ele retorna string vazia silenciosamente, sem erro — por isso vale validar o resultado quando o valor for obrigatório para a lógica do handler.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **ServeMux** | O multiplexador de rotas HTTP da biblioteca padrão do Go (`*http.ServeMux`), responsável por decidir qual handler atende cada requisição. |
| **Pattern (padrão)** | A string passada para `Handle`/`HandleFunc` descrevendo método HTTP (opcional) + caminho, podendo incluir wildcards, ex.: `"GET /books/{id}"`. |
| **Wildcard** | Segmento do padrão entre chaves (`{nome}` ou `{nome...}`) que captura parte variável do caminho da URL. |
| **`{nome...}`** | Wildcard "resto do caminho" — captura tudo que resta da URL a partir daquele ponto, incluindo barras adicionais; só pode ser o último segmento do padrão. |
| **`{$}`** | Marcador de correspondência exata; usado no final de um padrão para impedir que ele case como prefixo. |
| **`r.PathValue(nome)`** | Método de `*http.Request` que retorna o valor capturado por um wildcard nomeado do padrão que casou com a requisição. |
| **Method matching** | Capacidade do `ServeMux` de filtrar requisições pelo verbo HTTP declarado no início do padrão, respondendo `405` automaticamente quando o método não bate. |
| **Precedência de padrões** | Regra que determina qual padrão registrado "vence" quando mais de um poderia casar com a mesma URL; padrões mais específicos (literais) vencem wildcards na mesma posição. |
| **Panic em tempo de registro** | Comportamento do Go ao detectar dois padrões genuinamente ambíguos durante a chamada de `HandleFunc`/`Handle` — o programa encerra imediatamente ao subir, em vez de rotear de forma imprevisível. |

---

## 🚀 Próximos Passos

- [ ] **Imediato:** rode o exemplo como está e use o `test.http` (ou `curl`) para confirmar o comportamento das duas rotas ativas (`GET /books/{s}` e `GET /{s}/latest`).
- [ ] **Imediato:** descomente, uma de cada vez, cada linha comentada em `main.go` (`{id}`, `{d...}`, `{$}`, precedência) e reproduza os testes descritos na seção [Como Executar](#️-como-executar).
- [ ] **Intermediário:** adicione uma nova rota `DELETE /books/{id}` com um handler próprio, e confirme que uma requisição `GET` para o mesmo caminho continua caindo em outro handler (ou em `405`, se não houver um `GET` registrado para ele).
- [ ] **Intermediário:** provoque intencionalmente um padrão ambíguo (ex.: registre `GET /{a}/x` e `GET /y/{b}` e teste com `/y/x`) e observe o `panic` gerado ao rodar `go run .`.
- [ ] **Avançado:** implemente um middleware simples manualmente (uma função que recebe e retorna `http.Handler`, aplicando log antes/depois de cada requisição) e aplique-o ao `mux` — para sentir na prática a ausência de suporte nativo a middleware encadeado.
- [ ] **Avançado:** reescreva este exemplo usando `chi` ou `gorilla/mux` e compare a ergonomia de agrupar as rotas de `/books` sob um único prefixo com middleware compartilhado.
