# ⚡ fastjson em Go — Guia Didático

A biblioteca padrão do Go para JSON, `encoding/json`, tem um jeito de trabalhar muito comum: você define uma `struct` que descreve o formato esperado, chama `json.Unmarshal`, e o Go preenche a struct inteira — mesmo que você só precise de 2 campos dentro de um payload gigante com 200. Isso é ótimo pela previsibilidade e pelo type-safety, mas custa caro em CPU e alocação de memória quando o objetivo é só "ler um campinho aqui, outro ali" — por exemplo, decidir o roteamento de uma mensagem olhando só o campo `"type"` de um JSON grande.

A lib `github.com/valyala/fastjson` resolve esse problema específico: em vez de mapear o JSON inteiro para uma struct, ela faz um parsing "preguiçoso" (lazy) que devolve uma árvore de valores navegável — parecido com o DOM de um HTML — e você só "paga" o custo de extrair exatamente os campos que acessa. Este exemplo é minimalista de propósito — dois arquivos `main.go` pequenos — mas cobre o essencial: como parsear JSON sem struct, como navegar em campos aninhados, e como combinar `fastjson` com `encoding/json` quando você precisa das duas coisas.

---

## 📑 Sumário

- [🤔 O que é fastjson?](#-o-que-é-fastjson)
- [⚔️ fastjson vs Alternativas](#️-fastjson-vs-alternativas)
- [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
- [▶️ Como Executar](#️-como-executar)
- [⚖️ Trade-offs](#️-trade-offs)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [⚠️ Principais Problemas ao Trabalhar com fastjson](#️-principais-problemas-ao-trabalhar-com-fastjson)
- [❓ Perguntas de Entrevista](#-perguntas-de-entrevista)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é fastjson?

**Analogia:** pense num JSON grande como uma revista de 300 páginas. O jeito `encoding/json` (via `Unmarshal` numa struct) é como fotocopiar a revista inteira, página por página, encadernar tudo de novo num formato próprio — só depois disso você pode ler o que precisa. O `fastjson` é como folhear a revista direto pelo índice: você vai até a página que interessa, lê só aquele trecho, e nunca fotocopiou nada que não usou.

```
❌ encoding/json (parse completo em struct)
┌───────────────────────────────────────────┐
│  json.Unmarshal(bytes, &struct)            │
│  → percorre TODO o JSON                    │
│  → aloca memória para CADA campo da struct │
│  → converte tipos (string, int, bool...)   │
│  → você só usa 2 dos 20 campos preenchidos │
└───────────────────────────────────────────┘

✅ fastjson (parse lazy em árvore)
┌───────────────────────────────────────────┐
│  p.Parse(jsonData)                         │
│  → monta uma árvore de *fastjson.Value     │
│  → NADA é convertido/alocado ainda         │
│  → v.GetInt("campo") só converte           │
│    o campo que você pediu, na hora         │
└───────────────────────────────────────────┘
```

Tecnicamente:

- **`encoding/json`** exige que você declare o formato esperado via uma `struct` com [struct tags](#-glossário) (`json:"nome"`), e `Unmarshal` percorre o JSON inteiro convertendo cada campo para o tipo Go correspondente, alocando memória para a struct e para cada valor dentro dela.
- **`fastjson`** faz o parsing de um jeito mais "cru": ele monta uma árvore de nós (`*fastjson.Value`) que sabem *onde* cada valor está dentro do JSON original, mas só convertem o valor para um tipo Go (`string`, `int`, `bool`...) no momento em que você chama um método `Get*` — e reaproveita as mesmas estruturas internas entre chamadas de `Parse`, evitando muita alocação repetida.

Em código, a diferença fica clara comparando os dois estilos para o mesmo JSON:

```go
// ❌ encoding/json — precisa de uma struct que descreva TUDO
type Payload struct {
	Foo  string `json:"foo"`
	Num  int    `json:"num"`
	Bool bool   `json:"bool"`
	Arr  []int  `json:"arr"`
}

var payload Payload
json.Unmarshal([]byte(jsonData), &payload)
fmt.Println(payload.Foo) // "bar"
```

```go
// ✅ fastjson — sem struct nenhuma, acessa campo a campo (1/main.go)
var p fastjson.Parser
v, _ := p.Parse(jsonData)
fmt.Printf("foo=%s\n", v.GetStringBytes("foo")) // "bar"
```

---

## ⚔️ fastjson vs Alternativas

| Abordagem | O que faz | Quando usar |
|---|---|---|
| **`encoding/json`** | Biblioteca padrão do Go; usa reflection para mapear JSON ↔ struct via `Unmarshal`/`Marshal` | Praticamente sempre — é o padrão idiomático, tem melhor legibilidade e type-safety em tempo de compilação |
| **`fastjson`** (este exemplo) | Parsing lazy em árvore de `*fastjson.Value`, sem structs, focado em **leitura** de JSON | Quando você só precisa ler alguns campos de um JSON grande, sem precisar (des)serializar o objeto inteiro |
| **`jsoniter` (json-iterator/go)** | Drop-in replacement de `encoding/json` com API compatível, mas parsing mais rápido internamente | Quando você quer ganho de performance **sem reescrever** o código que já usa structs e `Unmarshal` |
| **`sonic` (bytedance/sonic)** | Parser JIT-compilado (usa JIT/SIMD), uma das opções mais rápidas do ecossistema Go, com API parecida com `encoding/json` | Cenários de altíssima performance (APIs com volume enorme de tráfego), geralmente em Linux/amd64 |

`fastjson` ocupa um nicho específico: ele não é um "substituto direto" de `encoding/json` (a API é bem diferente, não existe `Unmarshal` para structs arbitrárias) — é uma ferramenta para quando o objetivo principal é **navegar e extrair** campos de um JSON, não reconstruir o objeto inteiro como um tipo Go.

---

## 📚 Conceitos Fundamentais

### 1. `fastjson.Parser` e o Reuso do Parser

**Analogia:** é como usar sempre a mesma prancheta para anotar coisas, em vez de pegar uma prancheta nova a cada anotação — você economiza o "custo" de arranjar uma prancheta nova toda vez.

```go
// 1/main.go
var p fastjson.Parser

jsonData := `{"foo":"bar", "num":123, "bool":true, "arr": [1,2,3]}`

v, err := p.Parse(jsonData)
if err != nil {
	panic(err)
}
```

A lib oferece duas formas de parsear: a função solta `fastjson.Parse(jsonData)` (mais simples, mas aloca uma nova estrutura interna a cada chamada) e o jeito usado aqui, `var p fastjson.Parser` seguido de `p.Parse(...)`. Reutilizar a mesma instância de `Parser` em múltiplas chamadas permite que a biblioteca **reaproveise** os buffers e nós da árvore internos entre um parse e o próximo, evitando realocar memória do zero toda vez — esse é justamente o "pulo do gato" de performance que dá nome à lib.

> 💡 **Detalhe interessante:** essa é a razão pela qual o valor retornado por `Parse` (o `*fastjson.Value`) só é válido **até a próxima chamada de `Parse` no mesmo `Parser`** — veja a seção [Principais Problemas](#️-principais-problemas-ao-trabalhar-com-fastjson) para o porquê disso importar na prática.

### 2. Acesso via `Get*` — Parsing Sem Struct

```go
// 1/main.go
fmt.Printf("foo=%s\n", v.GetStringBytes("foo"))
fmt.Printf("num=%d\n", v.GetInt("num"))
fmt.Printf("bool=%v\n", v.GetBool("bool"))

a := v.GetArray("arr")
for i, value := range a {
	fmt.Printf("Index %d: %s\n", i, value)
}
```

Em vez de declarar uma struct Go que espelhe o formato do JSON, você navega direto na árvore retornada por `Parse` usando métodos como `GetStringBytes`, `GetInt`, `GetBool` e `GetArray`, passando o nome do campo como argumento. Cada `Get*` já devolve o valor no tipo Go esperado — e se o campo não existir ou tiver um tipo diferente, o método simplesmente devolve o "zero value" daquele tipo (`""`, `0`, `false`, `nil`), sem lançar erro. Isso é o que chamamos de parsing **schema-less**: o código não precisa conhecer o formato completo do JSON de antemão, só os caminhos que ele efetivamente usa.

> 💡 **Detalhe interessante:** repare que `GetStringBytes` devolve `[]byte`, não `string`. Isso é proposital — converter `[]byte` para `string` em Go normalmente exige uma cópia de memória; devolver `[]byte` permite que quem chama decida se precisa mesmo dessa cópia (por exemplo, para comparar bytes ou escrever direto num `io.Writer`, sem nunca precisar de uma `string`).

### 3. Padrão Híbrido: fastjson + encoding/json

**Analogia:** é como usar o índice de um livro para achar o capítulo certo (fastjson) e só então sentar e ler aquele capítulo com atenção, palavra por palavra (encoding/json) — você não lê o livro inteiro pra achar um capítulo, mas também não tenta "adivinhar" o conteúdo do capítulo sem realmente lê-lo.

```go
// 2/main.go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var p fastjson.Parser
jsonData := `{"user": {"name": "John Doe", "age": 30}}`

value, err := p.Parse(jsonData)
if err != nil {
	panic(err)
}
userJSON := value.Get("user").String()

var user User
if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
	panic(err)
}
fmt.Println(user.Name, user.Age)
```

Esse exemplo mistura as duas bibliotecas de propósito: `fastjson` é usado apenas para **navegar** até o sub-objeto `"user"` dentro de um JSON maior (`value.Get("user")`), sem se preocupar em tipar o resto do payload. Depois, esse trecho específico é convertido de volta para texto com `.String()` e passado para `json.Unmarshal`, que finalmente o converte para a struct `User` — com todo o type-safety e ergonomia do `encoding/json`. Esse padrão é útil quando um payload tem partes muito variáveis (você não sabe o formato inteiro, ou ele muda dependendo de outro campo) mas uma parte específica sempre tem um formato bem definido, que vale a pena tipar como struct.

---

## 🗂️ Estrutura do Projeto

```
22.4-fast-json/
├── 1/
│   ├── main.go   → parsing básico: lê campos tipados de um JSON sem struct
│   ├── go.mod
│   └── go.sum
└── 2/
    ├── main.go   → padrão híbrido: navega com fastjson, tipa o resultado com encoding/json
    ├── go.mod
    └── go.sum
```

Diferente de outros exemplos da pasta `20-extras`, aqui **cada subpasta é o seu próprio módulo Go** (tem `go.mod`/`go.sum` próprios), e ambos dependem de uma única biblioteca externa: `github.com/valyala/fastjson v1.6.4`.

---

## 🔍 Walkthrough do Código

### Exemplo 1 — `1/main.go`

```go
func main() {
	// 1. Uma única instância de Parser, reaproveitada em toda a função
	var p fastjson.Parser

	// 2. JSON "cru", sem struct nenhuma descrevendo o formato
	jsonData := `{"foo":"bar", "num":123, "bool":true, "arr": [1,2,3]}`

	// 3. Parse devolve uma árvore navegável (*fastjson.Value), não uma struct
	v, err := p.Parse(jsonData)
	if err != nil {
		panic(err)
	}

	// 4. Cada Get* converte só o campo pedido, no tipo Go esperado
	fmt.Printf("foo=%s\n", v.GetStringBytes("foo"))
	fmt.Printf("num=%d\n", v.GetInt("num"))
	fmt.Printf("bool=%v\n", v.GetBool("bool"))

	// 5. GetArray devolve um []*fastjson.Value, iterável normalmente
	a := v.GetArray("arr")
	for i, value := range a {
		fmt.Printf("Index %d: %s\n", i, value)
	}
}
```

Saída esperada:
```
foo=bar
num=123
bool=true
Index 0: 1
Index 1: 2
Index 2: 3
```

### Exemplo 2 — `2/main.go`

```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	var p fastjson.Parser
	jsonData := `{"user": {"name": "John Doe", "age": 30}}`

	// 1. Parseia o JSON inteiro só para navegação, sem tipar nada ainda
	value, err := p.Parse(jsonData)
	if err != nil {
		panic(err)
	}

	// 2. Navega até o sub-objeto "user" e o converte de volta para texto JSON
	userJSON := value.Get("user").String()

	// 3. Só ESSE trecho menor é convertido para uma struct Go tipada
	var user User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		panic(err)
	}
	fmt.Println(user.Name, user.Age)
}
```

Saída esperada: `John Doe 30`

O ponto-chave: `fastjson` entra apenas na etapa de **encontrar** o pedaço certo do JSON; a etapa de **tipar** esse pedaço continua sendo trabalho do `encoding/json`, que é mais ergonômico para isso.

---

## ▶️ Como Executar

Cada subpasta é um módulo independente (já tem `go.mod`/`go.sum`), então rode com `go run .` (ou `go run main.go`) dentro de cada uma:

```bash
# Exemplo 1 — dentro de aulas/20-extras/22.4-fast-json/1
go run main.go
# saída:
# foo=bar
# num=123
# bool=true
# Index 0: 1
# Index 1: 2
# Index 2: 3
```

```bash
# Exemplo 2 — dentro de aulas/20-extras/22.4-fast-json/2
go run main.go
# saída: John Doe 30
```

Para experimentar na prática:

1. No exemplo 1, troque `v.GetInt("num")` por `v.GetInt("num_errado")` (um campo que não existe) e observe que o programa não quebra — ele imprime `num=0`, o zero value de `int`.
2. No exemplo 2, edite o `jsonData` para incluir mais campos dentro de `"user"` que não estão na struct `User` (ex.: `"email"`) e confirme que `json.Unmarshal` simplesmente ignora o que não está mapeado na struct.

---

## ⚖️ Trade-offs

**✅ Vantagens**

- Evita alocar e converter campos que você nunca vai usar — relevante quando o JSON é grande e você só precisa de uma fração dele.
- Reaproveitar o `fastjson.Parser` entre chamadas reduz pressão sobre o garbage collector em código que parseia JSON com muita frequência (ex.: um consumidor de fila de alto volume).
- Não exige conhecer o formato completo do JSON de antemão — útil para payloads dinâmicos, polimórficos, ou quando você só precisa inspecionar um campo para decidir o que fazer a seguir (roteamento, validação rápida).

**❌ Desvantagens**

- Sem type-safety em tempo de compilação: `v.GetInt("campo")` compila normalmente mesmo que `"campo"` não exista ou não seja um número — erros de digitação só aparecem em runtime (e nem sempre como erro explícito, já que o zero value é retornado silenciosamente).
- API menos ergonômica que `encoding/json` para o caso comum de "tenho uma struct, quero preencher ela inteira" — nesse cenário, `fastjson` só adiciona complexidade.
- Os valores retornados (como `[]byte` de `GetStringBytes`) têm o tempo de vida ligado ao `*fastjson.Value` que os originou — reutilizar o mesmo `Parser` para um novo `Parse` invalida os dados anteriores (ver seção de problemas comuns).

---

## 🎯 Casos de Uso Ideais

**Use fastjson quando:**
- Você precisa ler só alguns campos de um JSON grande, e o custo de popular uma struct inteira via `encoding/json` é perceptível (hot path de performance, alto volume de requisições/mensagens);
- O formato do JSON é parcialmente desconhecido ou variável, e você só quer navegar até uma parte específica antes de decidir como tratá-la — como no padrão híbrido do exemplo 2;
- Você está construindo algo como um roteador de mensagens, que só precisa inspecionar um campo tipo `"type"` ou `"event"` para decidir o próximo passo, sem desserializar o payload inteiro.

**Evite fastjson quando:**
- O JSON é pequeno ou a aplicação não tem restrição real de performance — nesse caso, `encoding/json` com uma struct é mais legível, mais seguro e mais fácil de manter;
- Você precisa (de)serializar o objeto inteiro de forma estruturada e repetidamente (por exemplo, uma API REST tradicional trocando payloads tipados) — `encoding/json` (ou `jsoniter`/`sonic` como drop-in) é a escolha mais natural;
- Legibilidade e manutenção pesam mais que performance bruta para o time — a API baseada em `Get*` com strings de campo é mais propensa a erros de digitação silenciosos do que uma struct tipada.

---

## ⚠️ Principais Problemas ao Trabalhar com fastjson

### 1. Reutilizar um `*fastjson.Value` depois de um novo `Parse`

```go
// ❌ v1 é sobrescrito internamente quando p.Parse é chamado de novo
var p fastjson.Parser
v1, _ := p.Parse(`{"foo":"primeiro"}`)
v2, _ := p.Parse(`{"foo":"segundo"}`)
fmt.Println(v1.GetStringBytes("foo")) // pode não ser mais "primeiro"!
```

Como o `Parser` reaproveita seus buffers internos para economizar alocação, o valor devolvido por uma chamada de `Parse` só é garantidamente válido **até a próxima chamada de `Parse` no mesmo Parser**. Guardar um `*fastjson.Value` antigo e continuar usando depois de um novo parse é um bug silencioso — o dado pode ter sido sobrescrito.

**Solução:** ou processe (e descarte) cada `*fastjson.Value` completamente antes do próximo `Parse`, ou use uma instância de `Parser` separada (ou `fastjson.ParserPool`) para cada JSON que precisa coexistir na memória ao mesmo tempo.

### 2. Confundir zero value com "campo ausente"

```go
// ❌ Não dá para distinguir "num é 0" de "num não existe"
v.GetInt("num") // devolve 0 nos dois casos
```

Diferente de `encoding/json` com `json.RawMessage` ou ponteiros (`*int`), os métodos `Get*` do `fastjson` não têm como sinalizar "esse campo não existe" de forma distinta do zero value do tipo — ambos retornam `0`, `""`, `false` ou `nil`.

**Solução:** quando a distinção importa, use o método genérico `v.Get("campo")` primeiro (que devolve `nil` se o campo realmente não existir) antes de chamar o `Get*` tipado:

```go
// ✅
if field := v.Get("num"); field != nil {
	fmt.Println(field.GetInt())
} else {
	fmt.Println("campo 'num' não existe")
}
```

### 3. Ignorar o erro de `Parse`

```go
// ❌ Se o JSON for inválido, v fica nil e qualquer Get* seguinte causa panic
v, _ := p.Parse(jsonMalformado)
v.GetStringBytes("foo") // panic: nil pointer dereference
```

**Solução:** sempre checar o `err` retornado por `Parse` antes de navegar no resultado — como já fazem `1/main.go` e `2/main.go`, que chamam `panic(err)` explicitamente se o parse falhar, em vez de seguir em frente com um valor inválido.

---

## ❓ Perguntas de Entrevista

**Qual a principal diferença entre `fastjson` e `encoding/json`?**
`encoding/json` usa reflection para mapear um JSON inteiro para uma struct Go via `Unmarshal`, convertendo e alocando memória para todos os campos de uma vez. `fastjson` faz um parsing "lazy": ele monta uma árvore navegável de valores (`*fastjson.Value`) sem struct nenhuma, e só converte para um tipo Go concreto (via `GetInt`, `GetStringBytes` etc.) o campo que você efetivamente acessa — o que economiza trabalho quando você só precisa de uma fração do JSON.

**O que significa "parsing lazy" no contexto do fastjson?**
Significa que o `Parse` inicial só identifica a estrutura do JSON (onde cada valor começa e termina), sem converter nada para um tipo Go de fato. A conversão real (para `string`, `int`, `bool` etc.) só acontece no momento em que você chama um método `Get*` sobre um campo específico — por isso campos nunca acessados nunca chegam a ser convertidos.

**Por que reaproveitar a mesma instância de `fastjson.Parser` em vez de criar uma nova a cada `Parse`?**
Porque o `Parser` guarda internamente os buffers e nós da árvore usados no parsing anterior, e os reutiliza (sobrescrevendo-os) na próxima chamada de `Parse`, em vez de alocar tudo do zero. Isso reduz significativamente a pressão sobre o garbage collector em cenários que parseiam JSON repetidamente — o efeito colateral é que o `*fastjson.Value` de um parse anterior deixa de ser confiável assim que um novo `Parse` acontece no mesmo `Parser`.

**Quando faz sentido combinar `fastjson` com `encoding/json` no mesmo código, como no exemplo 2?**
Quando você tem um JSON cujo formato completo não vale a pena (ou não é possível) tipar inteiro como struct, mas uma parte específica dele tem um formato bem definido que você quer manipular com type-safety. `fastjson` entra só para "navegar" até esse trecho (`value.Get("user")`), e `encoding/json` entra para desserializar esse trecho menor numa struct tipada — combinando a flexibilidade de um com a ergonomia do outro.

**Por que `GetStringBytes` devolve `[]byte` em vez de `string`?**
Para evitar uma cópia de memória desnecessária. Converter `[]byte` para `string` em Go geralmente exige alocar uma nova string e copiar os bytes; devolvendo `[]byte` diretamente, `fastjson` deixa para quem chama decidir se essa cópia é realmente necessária (por exemplo, para comparar com `bytes.Equal` ou escrever direto num `io.Writer`, sem nunca precisar materializar uma `string`).

**`fastjson` é sempre a escolha certa quando performance importa?**
Não necessariamente. Existem alternativas como `jsoniter` (compatível com a API de `encoding/json`, então não exige reescrever o código) e `sonic` (parser JIT-compilado, geralmente ainda mais rápido, mas com restrições de plataforma). `fastjson` é uma boa escolha especificamente quando o padrão de acesso é "ler alguns campos de um JSON grande sem popular uma struct inteira" — para outros padrões de uso, as alternativas podem ser mais adequadas.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **Parsing lazy** | Estratégia de parsing em que a conversão de um valor para um tipo concreto só acontece no momento em que ele é efetivamente acessado, não durante o parse inicial. |
| **`*fastjson.Value`** | Nó de uma árvore que representa um valor JSON (objeto, array, string, número, etc.), navegável via métodos como `Get`, `GetInt`, `GetArray`. |
| **`fastjson.Parser`** | Struct que guarda o estado interno reutilizável do parsing; reaproveitá-la entre chamadas de `Parse` reduz alocação de memória. |
| **Zero value** | Valor padrão de um tipo Go quando nenhum valor explícito é atribuído (`0` para números, `""` para strings, `false` para bool, `nil` para ponteiros/slices). |
| **Struct tag** | Anotação como `` `json:"nome"` `` numa struct Go, usada por `encoding/json` para mapear o nome do campo Go ao nome do campo no JSON. |
| **Reflection** | Mecanismo que permite a um programa inspecionar/manipular seus próprios tipos em tempo de execução; é o que `encoding/json` usa internamente para preencher structs, com custo de performance associado. |
| **DOM (Document Object Model)** | Termo emprestado do mundo web: uma representação em árvore de um documento (aqui, um JSON) que permite navegar e consultar partes específicas sem processar o documento inteiro de uma vez. |
| **Schema-less** | Estilo de acesso a dados que não exige declarar de antemão a estrutura completa esperada — você só especifica os caminhos/campos que realmente vai usar. |

---

## 🚀 Próximos Passos

- [ ] **Imediato:** no exemplo 1, acesse um campo que não existe (ex.: `v.GetInt("inexistente")`) e confirme que o programa não quebra, apenas devolve o zero value.
- [ ] **Imediato:** no exemplo 2, adicione um novo campo (ex.: `"email": "john@example.com"`) dentro de `"user"` no JSON e na struct `User`, e confirme que ele aparece corretamente após o `Unmarshal`.
- [ ] **Intermediário:** reescreva o exemplo 1 usando `encoding/json` com uma struct equivalente, e compare a legibilidade e a quantidade de código entre as duas abordagens.
- [ ] **Intermediário:** provoque o problema descrito em [Principais Problemas #1](#️-principais-problemas-ao-trabalhar-com-fastjson): parseie dois JSONs diferentes com o mesmo `Parser`, guarde os dois `*fastjson.Value`, e observe o que acontece ao acessar o primeiro depois do segundo `Parse`.
- [ ] **Avançado:** escreva um benchmark (`*_test.go` com `func BenchmarkX(b *testing.B)`) comparando `fastjson` contra `encoding/json` para o mesmo JSON, variando o tamanho do payload e a quantidade de campos efetivamente lidos.
- [ ] **Avançado:** experimente substituir `fastjson` por `jsoniter` ou `sonic` num dos dois exemplos e compare a ergonomia da API e os resultados do benchmark do passo anterior.
