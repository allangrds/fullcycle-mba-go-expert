# 👀 fsnotify (File System Notifications) em Go — Guia Didático

Como saber que um arquivo mudou no disco? A resposta ingênua é: fica checando de tempos em tempos (`time.Sleep` + `os.Stat`, em loop). Funciona, mas desperdiça CPU o tempo todo e sempre existe um atraso entre a mudança real e o momento em que alguém percebe. A resposta mais elegante é deixar o próprio sistema operacional avisar, no exato instante em que algo acontece — é isso que a biblioteca `fsnotify` faz em Go. Este exemplo é minimalista de propósito — um único arquivo `main.go` com pouco mais de 60 linhas — mas cobre o essencial: como criar um watcher, como reagir a eventos de um arquivo específico, e o padrão mais comum de uso em produção: **hot-reload de configuração** sem reiniciar o processo.

---

## 📑 Sumário

- [🤔 O que é fsnotify?](#-o-que-é-fsnotify)
- [⚔️ fsnotify vs Alternativas](#️-fsnotify-vs-alternativas)
- [📚 Conceitos Fundamentais](#-conceitos-fundamentais)
- [🗂️ Estrutura do Projeto](#️-estrutura-do-projeto)
- [🔍 Walkthrough do Código](#-walkthrough-do-código)
- [▶️ Como Executar](#️-como-executar)
- [⚖️ Trade-offs](#️-trade-offs)
- [🎯 Casos de Uso Ideais](#-casos-de-uso-ideais)
- [⚠️ Principais Problemas ao Trabalhar com fsnotify](#️-principais-problemas-ao-trabalhar-com-fsnotify)
- [❓ Perguntas de Entrevista](#-perguntas-de-entrevista)
- [📖 Glossário](#-glossário)
- [🚀 Próximos Passos](#-próximos-passos)

---

## 🤔 O que é fsnotify?

**Analogia:** imagine que você precisa saber assim que alguém depositar uma carta na sua caixa de correio. Existem duas formas de fazer isso: você pode ir até a caixa a cada 5 minutos para conferir (gastando tempo e energia mesmo quando não chegou nada), ou pode pedir para o carteiro tocar a campainha exatamente no momento em que a carta for entregue. O `fsnotify` é a campainha: em vez de você ficar checando o arquivo, o sistema operacional avisa sua aplicação assim que algo muda.

```
❌ SEM fsnotify (polling manual)
┌─────────────────────────────────────────┐
│  loop infinito:                          │
│    - lê o arquivo                        │
│    - compara com a versão anterior       │
│    - dorme 1 segundo                     │
│    - repete...                           │
│  → gasta CPU mesmo sem nada mudar        │
│  → sempre existe um atraso (até 1s)      │
│  → não escala bem com muitos arquivos    │
└─────────────────────────────────────────┘

✅ COM fsnotify (eventos do SO)
┌─────────────────────────────────────────┐
│  registra o arquivo no watcher           │
│  fica bloqueado esperando um evento      │
│  → o SO avisa NO INSTANTE em que         │
│    o arquivo é escrito/criado/removido   │
│  → zero CPU gasta enquanto nada muda     │
│  → reação imediata, sem atraso artificial│
└─────────────────────────────────────────┘
```

Tecnicamente, `fsnotify` é uma biblioteca Go que expõe, de forma multiplataforma, as APIs nativas de notificação de sistema de arquivos de cada sistema operacional:

- **Linux** → `inotify`
- **macOS/BSD** → `kqueue`
- **Windows** → `ReadDirectoryChangesW`

Ela não implementa monitoramento "do zero" — apenas oferece uma API única em Go por cima de mecanismos que o kernel de cada SO já fornece.

Em código, a diferença fica bem clara:

```go
// ❌ Polling manual — funciona, mas gasta CPU e sempre tem atraso
for {
    data, _ := os.ReadFile("config.json")
    if string(data) != lastKnown {
        fmt.Println("mudou!")
        lastKnown = string(data)
    }
    time.Sleep(1 * time.Second)
}
```

```go
// ✅ fsnotify — este é o padrão do exemplo desta pasta
watcher, _ := fsnotify.NewWatcher()
watcher.Add("config.json")

for event := range watcher.Events {
    if event.Op&fsnotify.Write == fsnotify.Write {
        fmt.Println("mudou!", event.Name)
    }
}
```

---

## ⚔️ fsnotify vs Alternativas

| Abordagem | O que faz | Custo/Atraso | Quando usar |
|---|---|---|---|
| **Polling manual** (`time.Sleep` + `os.Stat`/`os.ReadFile`) | Fica checando o arquivo em intervalos fixos | Gasta CPU constantemente; atraso igual ao intervalo do sleep | Scripts simples, protótipos, ou quando o SO/filesystem não suporta bem notificações (ex.: alguns filesystems de rede) |
| **fsnotify** (este exemplo) | Recebe eventos do SO assim que algo muda | Praticamente zero CPU parado; reação quase instantânea | Aplicações Go que precisam reagir a mudanças em arquivos/diretórios locais |
| **Bibliotecas de config com watch embutido** (ex.: Viper `WatchConfig()`) | Envolve o fsnotify (ou equivalente) e já entrega o valor recarregado, com parsing e tipagem prontos | Mesmo custo do fsnotify, mas com mais conveniência e menos código boilerplate | Projetos que já usam uma lib de configuração e só querem hot-reload "de graça" |
| **Reload via sinal (`SIGHUP`)** | O processo recarrega a config apenas quando recebe um sinal explícito do operador (`kill -HUP <pid>`) | Sem monitoramento contínuo; depende de alguém (humano ou script) disparar o sinal | Ambientes onde reload automático é indesejado e o operador prefere decidir *quando* recarregar |

`fsnotify` fica no meio do caminho: mais eficiente e reativo que polling, porém mais "cru" (menos pronto para uso) que uma lib de configuração completa como Viper, que geralmente usa fsnotify por baixo dos panos.

---

## 📚 Conceitos Fundamentais

### 1. `fsnotify.NewWatcher()` — o Vigia Ligado ao Sistema Operacional

**Analogia:** é como contratar um segurança que fica de olho em portas específicas do prédio — ele não entra em ação sozinho, alguém precisa dizer quais portas ele deve vigiar.

```go
// main.go — linhas 21-25
watcher, err := fsnotify.NewWatcher()
if err != nil {
    panic(err)
}
defer watcher.Close()
```

`NewWatcher()` cria uma instância que se conecta diretamente ao mecanismo de notificação do sistema operacional (`inotify` no Linux, por exemplo). Esse watcher começa "vazio" — ele não observa nada até que algum caminho seja explicitamente adicionado com `Add()`. Assim como um arquivo aberto, ele consome um recurso do sistema operacional (um file descriptor, no caso do Linux) e por isso precisa ser fechado com `Close()` quando não for mais necessário — daí o `defer watcher.Close()`.

### 2. `watcher.Add(path)` — Arquivo vs Diretório

**Analogia:** é diferente pedir para o segurança vigiar "a porta 302" ou vigiar "o corredor inteiro do 3º andar". Vigiar uma porta específica é mais simples, mas se essa porta for trocada de lugar, o segurança continua olhando para o lugar errado.

```go
// main.go — linha 50
err = watcher.Add("config.json")
```

Neste exemplo, o watch é feito diretamente sobre o **arquivo** `config.json`, não sobre o diretório que o contém. Isso funciona bem enquanto o arquivo é editado *in place* (o conteúdo é sobrescrito, mas o arquivo continua sendo o "mesmo" no disco). O problema aparece quando um editor de texto salva usando **rename atômico** (grava um arquivo temporário e depois renomeia por cima do original) — nesse caso, o `config.json` original deixa de existir e o watch, que estava vigiando aquele arquivo específico, para de funcionar silenciosamente. Esse comportamento é detalhado na seção [⚠️ Principais Problemas](#️-principais-problemas-ao-trabalhar-com-fsnotify).

> 💡 **Detalhe interessante:** por isso é comum, em código de produção, fazer o watch no **diretório pai** e filtrar os eventos pelo nome do arquivo de interesse (`event.Name == "config.json"`) — assim, mesmo que o arquivo seja recriado por um rename atômico, o watch no diretório continua ativo e enxerga os novos eventos.

### 3. Goroutine + `select` sobre `Events` e `Errors`

**Analogia:** é o mesmo princípio de "colocar alguém para atender o telefone enquanto você segue com sua rotina" usado no exemplo de graceful shutdown — aqui, a goroutine fica de plantão ouvindo dois "telefones" ao mesmo tempo: um de eventos, outro de erros.

```go
// main.go — linhas 28-49
done := make(chan bool)
go func() {
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            fmt.Println("event :", event)
            if event.Op&fsnotify.Write == fsnotify.Write {
                MarshalConfig("config.json")
                fmt.Println("modified file:", event.Name)
                fmt.Println(config)
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            fmt.Println("error:", err)
        }
    }
}()
```

O watcher expõe dois canais: `Events` (mudanças detectadas) e `Errors` (falhas internas de monitoramento). Como ambos podem chegar a qualquer momento, o padrão idiomático é um `select` dentro de um loop infinito, rodando em uma **goroutine separada** — assim a função `main` fica livre para continuar (neste exemplo, ela só espera em `<-done`, mas em uma aplicação real ela poderia seguir subindo um servidor HTTP, por exemplo). O `ok` retornado por cada leitura do canal indica se o canal ainda está aberto; quando o watcher é fechado, os canais são fechados também, e `ok` vem `false` — é o sinal para a goroutine encerrar com `return`.

### 4. `event.Op` — um Bitmask de Operações

```go
// main.go — linha 37
if event.Op&fsnotify.Write == fsnotify.Write {
```

Cada evento carrega um campo `Op` do tipo `fsnotify.Op`, que é um **bitmask**: um único valor inteiro onde cada bit representa uma operação possível — `Create`, `Write`, `Remove`, `Rename`, `Chmod`. Um mesmo evento pode, em teoria, carregar mais de uma operação combinada. Por isso a forma correta de checar "esse evento inclui uma escrita?" não é comparar com `==` diretamente, e sim usar a operação bit a bit `&` (AND) para isolar o bit de interesse e comparar o resultado com a própria constante:

```go
event.Op & fsnotify.Write == fsnotify.Write
```

Se o bit `Write` estiver "ligado" em `event.Op`, o `&` preserva exatamente esse bit e o resultado é igual a `fsnotify.Write`; caso contrário, o resultado é zero (diferente de `fsnotify.Write`). Este exemplo verifica **apenas** o bit `Write` — eventos de `Create`, `Remove`, `Rename` e `Chmod` são impressos (`fmt.Println("event :", event)`), mas não disparam o reload da configuração.

> 💡 **Detalhe interessante:** ignorar os outros bits é uma simplificação didática. Em produção, um `Remove` ou `Rename` no arquivo observado costuma ser tão importante quanto um `Write` — é justamente o sinal de que o watch precisa ser refeito (veja o problema do rename atômico mais abaixo).

### 5. O Padrão de Hot-Reload

```go
// main.go — linha 26 (inicialização) e linha 38 (a cada evento)
MarshalConfig("config.json")
```

```go
// main.go — linhas 57-66
func MarshalConfig(file string) {
    data, err := os.ReadFile(file)
    if err != nil {
        panic(err)
    }
    err = json.Unmarshal(data, &config)
    if err != nil {
        panic(err)
    }
}
```

`MarshalConfig` é chamada duas vezes no fluxo do programa: uma vez na inicialização (para carregar a configuração inicial) e uma vez a cada evento de `Write` detectado no arquivo. O resultado é **hot-reload**: a variável global `config` (`var config DBConfig`, linha 18) é atualizada automaticamente sempre que `config.json` é salvo, sem precisar reiniciar o processo — o mesmo princípio usado por servidores web em desenvolvimento (`air`, usado neste próprio projeto para live-reload do binário) e por bibliotecas de configuração como Viper.

### 6. `defer watcher.Close()` e o Canal `done` que Nunca Recebe Nada

```go
// main.go — linha 25 e linhas 28, 54
defer watcher.Close()
// ...
done := make(chan bool)
// ...
<-done
```

O canal `done` é criado, mas em nenhum lugar do código algo é enviado para ele (`done <- true` nunca acontece). Isso significa que a linha `<-done`, no final da `main`, bloqueia **para sempre** — o programa só termina se for interrompido externamente (`Ctrl+C`, `kill`, etc.). É um padrão comum em exemplos didáticos para "manter o processo vivo" enquanto a goroutine de monitoramento roda em segundo plano, mas é também uma simplificação: em código de produção, esse `done` normalmente é fechado em resposta a um sinal do sistema operacional, seguindo o mesmo padrão do exemplo de [graceful shutdown](../22.1-graceful-shutdown/README.md) desta série.

---

## 🗂️ Estrutura do Projeto

```
22.3-fsnotify/
├── main.go       → todo o exemplo: watcher + goroutine de eventos + reload de config
├── config.json   → arquivo observado; editar e salvar dispara o evento de Write
├── go.mod        → módulo Go, com a dependência externa github.com/fsnotify/fsnotify
├── go.sum
└── .air.toml     → configuração do "air" (live-reload do binário durante desenvolvimento, opcional)
```

Diferente dos exemplos de graceful shutdown e panic/recover (que usam só a biblioteca padrão), este exemplo depende de um pacote de terceiros: `github.com/fsnotify/fsnotify` — porque monitorar o sistema de arquivos de forma multiplataforma exige lidar com APIs nativas diferentes por sistema operacional, algo que a biblioteca padrão do Go não abstrai.

---

## 🔍 Walkthrough do Código

Seguindo a ordem de execução real do programa:

```go
// 1. Cria o watcher (ainda não observa nada)
watcher, err := fsnotify.NewWatcher()
if err != nil {
    panic(err)
}
defer watcher.Close()

// 2. Carrega a configuração inicial, antes mesmo de começar a observar
MarshalConfig("config.json")

// 3. Cria o canal "done" (nunca recebe nada — mantém o processo vivo)
done := make(chan bool)

// 4. Sobe a goroutine que vai ficar escutando eventos e erros do watcher
go func() {
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            fmt.Println("event :", event)
            // 6. Se o evento for de escrita, recarrega a config e imprime o resultado
            if event.Op&fsnotify.Write == fsnotify.Write {
                MarshalConfig("config.json")
                fmt.Println("modified file:", event.Name)
                fmt.Println(config)
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            fmt.Println("error:", err)
        }
    }
}()

// 5. Só agora registra o arquivo a ser observado —
//    a partir daqui, qualquer mudança em config.json gera um evento
err = watcher.Add("config.json")
if err != nil {
    panic(err)
}

// 7. A main bloqueia aqui indefinidamente
<-done
```

O ponto-chave: os passos 4 e 5 acontecem em uma ordem específica — a goroutine que consome os eventos é iniciada **antes** de `watcher.Add()` ser chamado. Isso evita uma condição de corrida em que um evento pudesse, teoricamente, chegar antes de existir alguém lendo do canal (o que causaria bloqueio ou perda de eventos dependendo do buffer interno do watcher).

---

## ▶️ Como Executar

```bash
# Dentro da pasta aulas/20-extras/22.3-fsnotify
go run .

# Ou, com live-reload do binário durante o desenvolvimento:
air
```

Para observar o fsnotify em ação:

1. Rode `go run .` — o programa fica bloqueado, sem imprimir nada ainda (não há log de "iniciado").
2. Em outro terminal (ou editor de texto), abra `config.json` e altere um valor, por exemplo trocando `"host": "localhost"` por `"host": "127.0.0.1"`. Salve o arquivo.
3. No terminal do programa, você deve ver algo como:
   ```
   event : "config.json": WRITE
   modified file: config.json
   {mysql 127.0.0.1 root root}
   ```
4. Repita o passo 2 algumas vezes — cada salvamento gera um novo evento e um novo reload, sem reiniciar o processo.
5. **Experimente o gotcha:** se seu editor salvar usando rename atômico (comum em editores GUI e em alguns plugins de "salvar com segurança" do Vim/VS Code), você notará que, após o **primeiro** salvamento, os eventos seguintes deixam de aparecer — o watch morreu porque o arquivo original foi substituído por um novo inode. Esse comportamento está detalhado na seção de problemas abaixo.

---

## ⚖️ Trade-offs

**✅ Vantagens**

- Reação praticamente instantânea a mudanças no sistema de arquivos, sem o atraso artificial de um polling por intervalo.
- Zero consumo de CPU enquanto nada muda — o processo fica bloqueado esperando o SO notificar, não gastando ciclos em loop.
- API relativamente simples e multiplataforma (mesmo código roda em Linux, macOS e Windows, embora o comportamento de baixo nível varie).

**❌ Desvantagens**

- Comportamento não é 100% uniforme entre sistemas operacionais — o conjunto de eventos e sutilezas (como o problema de rename atômico) varia entre `inotify`, `kqueue` e `ReadDirectoryChangesW`.
- Watch em um arquivo específico é frágil: se o arquivo for recriado (rename, exclusão + criação), o watch pode parar de funcionar sem aviso explícito de erro.
- Exige uma goroutine dedicada e disciplina para drenar **dois** canais (`Events` e `Errors`) continuamente — esquecer de tratar `Errors` (como faz este próprio exemplo, que só imprime) pode esconder falhas reais de monitoramento.
- Em sistemas com muitos arquivos monitorados simultaneamente, há limites de recursos do SO a considerar (por exemplo, o limite de watches do `inotify` no Linux, configurável via `fs.inotify.max_user_watches`).

---

## 🎯 Casos de Uso Ideais

**Use fsnotify quando:**
- Você precisa de hot-reload de configuração sem reiniciar o processo (o cenário deste exemplo);
- Está construindo uma ferramenta de live-reload para desenvolvimento (como o próprio `air`, usado nesta pasta);
- Precisa sincronizar ou reagir a mudanças em arquivos locais (por exemplo, disparar um build, invalidar um cache, reindexar um arquivo);
- O volume de arquivos/diretórios observados é gerenciável e está em um filesystem local que suporta bem notificações nativas.

**Evite ou avalie alternativas quando:**
- O volume de arquivos a monitorar é muito grande (dezenas de milhares), o que pode esbarrar em limites de recursos do sistema operacional;
- O ambiente envolve filesystems de rede (NFS, alguns volumes montados em containers) onde notificações nativas podem não funcionar de forma confiável — nesses casos, polling manual pode ser mais previsível, mesmo sendo menos eficiente;
- Você só precisa de configuração com reload, sem lógica extra — nesse caso, uma biblioteca pronta como Viper (que já lida com os gotchas do fsnotify por baixo dos panos) costuma ser mais produtiva do que reimplementar o watch manualmente.

---

## ⚠️ Principais Problemas ao Trabalhar com fsnotify

### 1. Watch em Arquivo Único Quebra com Rename Atômico

```go
// ❌ Watch direto no arquivo — frágil se o editor salvar via rename
watcher.Add("config.json")
```

Muitos editores de texto (e algumas ferramentas de deploy) não sobrescrevem o arquivo original diretamente: eles escrevem um arquivo temporário e depois o renomeiam por cima do arquivo alvo. Do ponto de vista do sistema operacional, o `config.json` que estava sendo observado deixa de existir (um novo inode assume o nome), e o watch registrado especificamente para aquele arquivo para de funcionar — sem lançar nenhum erro explícito.

**Solução:** observar o **diretório** que contém o arquivo, e filtrar os eventos pelo nome:

```go
// ✅ Watch no diretório, filtrando pelo nome do arquivo de interesse
watcher.Add(".")

for event := range watcher.Events {
    if filepath.Base(event.Name) == "config.json" &&
        event.Op&fsnotify.Write == fsnotify.Write {
        MarshalConfig("config.json")
    }
}
```

### 2. Variável Global Sem Sincronização

```go
// ❌ config é lida e escrita sem nenhuma proteção
var config DBConfig
// escrita: dentro da goroutine do watcher
// leitura: potencialmente em qualquer outra parte do programa
```

Neste exemplo, `config` é uma variável global, escrita pela goroutine do watcher a cada `Write` detectado. Se qualquer outra parte da aplicação (por exemplo, um handler HTTP em uma versão mais completa deste programa) ler `config` ao mesmo tempo em que ela está sendo atualizada, isso é uma **condição de corrida (data race)** — não exercida no exemplo atual porque nada mais lê `config`, mas um risco real assim que o programa crescer.

**Solução:** proteger o acesso com `sync.RWMutex` (leitura concorrente livre, escrita exclusiva):

```go
// ✅
var (
    config   DBConfig
    configMu sync.RWMutex
)

func MarshalConfig(file string) {
    data, err := os.ReadFile(file)
    if err != nil {
        panic(err)
    }
    var newConfig DBConfig
    if err := json.Unmarshal(data, &newConfig); err != nil {
        panic(err)
    }
    configMu.Lock()
    config = newConfig
    configMu.Unlock()
}
```

### 3. Ignorar o Canal `Errors`

```go
// ❌ O erro só é impresso, nunca tratado de verdade
case err, ok := <-watcher.Errors:
    if !ok {
        return
    }
    fmt.Println("error:", err)
```

Um erro no canal `Errors` normalmente indica um problema real no monitoramento (por exemplo, o watcher perdeu a capacidade de observar um caminho). Apenas imprimir e seguir em frente significa que, se esses erros se acumularem, ninguém vai notar até que o comportamento de hot-reload simplesmente pare de funcionar.

**Solução:** tratar o erro de forma que ele seja observável (log estruturado, métrica, alerta) e, se fizer sentido, tentar reestabelecer o watch:

```go
// ✅
case err, ok := <-watcher.Errors:
    if !ok {
        return
    }
    log.Printf("erro no watcher: %v — tentando readicionar o watch\n", err)
    if addErr := watcher.Add("config.json"); addErr != nil {
        log.Printf("falha ao readicionar watch: %v\n", addErr)
    }
```

### 4. Usar `panic` para Erros Esperados

```go
// ❌ main.go — linhas 22-24, 51-53, 58-60, 63-65
if err != nil {
    panic(err)
}
```

O exemplo usa `panic` sempre que algo falha (criar o watcher, adicionar o watch, ler o arquivo, fazer o parse do JSON). Isso é aceitável para um exemplo didático curto, mas em código de produção esses são justamente os casos onde um `error` tratado (log + retry, ou encerramento controlado) é mais apropriado — um `config.json` temporariamente ilegível (por exemplo, sendo escrito no exato momento da leitura) não deveria derrubar o processo inteiro.

**Solução:** trocar `panic` por tratamento de erro que não interrompe o programa:

```go
// ✅
func MarshalConfig(file string) error {
    data, err := os.ReadFile(file)
    if err != nil {
        return fmt.Errorf("erro ao ler config: %w", err)
    }
    if err := json.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("erro ao parsear config: %w", err)
    }
    return nil
}
```

---

## ❓ Perguntas de Entrevista

**O que é fsnotify e qual problema ele resolve?**
`fsnotify` é uma biblioteca Go que permite que uma aplicação seja notificada por eventos do sistema operacional sempre que um arquivo ou diretório observado é criado, escrito, removido, renomeado ou tem suas permissões alteradas. Ele resolve o problema de saber que algo mudou no sistema de arquivos sem precisar ficar checando repetidamente (polling) — o que economiza CPU e reduz o atraso entre a mudança real e a reação da aplicação a ela.

**Como o fsnotify funciona por baixo dos panos?**
Ele não implementa monitoramento de arquivos do zero: apenas expõe, com uma API única em Go, os mecanismos nativos que cada sistema operacional já oferece para esse fim — `inotify` no Linux, `kqueue` no macOS/BSD, e `ReadDirectoryChangesW` no Windows. Isso significa que o comportamento exato de certos casos-limite pode variar entre sistemas operacionais, já que cada um implementa esses mecanismos de forma diferente.

**Por que o código verifica `event.Op&fsnotify.Write == fsnotify.Write` em vez de comparar diretamente com `==`?**
Porque `Op` é um **bitmask**: um único valor inteiro onde cada bit representa uma operação possível (`Create`, `Write`, `Remove`, `Rename`, `Chmod`), e teoricamente mais de um bit pode estar "ligado" no mesmo evento. O operador `&` (AND bit a bit) isola apenas o bit correspondente à operação de interesse; comparando o resultado desse AND com a própria constante, o código verifica "esse bit está ligado?" sem se importar com o estado dos outros bits. Comparar `event.Op == fsnotify.Write` diretamente falharia sempre que o evento tivesse qualquer outro bit adicional ligado.

**Por que observar um arquivo específico é considerado frágil, e qual a alternativa mais robusta?**
Porque muitos editores de texto e ferramentas de deploy não sobrescrevem o arquivo original: eles escrevem um arquivo temporário e o renomeiam por cima do arquivo alvo (um padrão chamado "rename atômico", usado para evitar que outros processos leiam um arquivo pela metade). Quando isso acontece, o inode original — para o qual o watch foi registrado — deixa de existir, e o watcher perde a referência silenciosamente, sem lançar um erro óbvio. A alternativa mais robusta é registrar o watch no **diretório** que contém o arquivo, e filtrar os eventos recebidos pelo nome do arquivo de interesse — assim, mesmo que o arquivo seja recriado, o watch no diretório continua ativo.

**Quando faz mais sentido usar fsnotify diretamente em vez de uma biblioteca de configuração como Viper?**
Usar fsnotify diretamente faz sentido quando você precisa de controle fino sobre o que observar e como reagir (por exemplo, observar arquivos que não são de configuração, como assets estáticos, ou implementar uma lógica de reload muito específica). Uma biblioteca como Viper, que usa fsnotify internamente através do método `WatchConfig()`, é preferível quando o objetivo é simplesmente "recarregar a configuração da aplicação quando o arquivo mudar" — ela já resolve boa parte dos gotchas de parsing, tipagem e observação de arquivo por você.

**O que aconteceria se a goroutine que lê `watcher.Events` parasse de rodar (por exemplo, por um bug ou um `return` prematuro)?**
O canal `Events` do fsnotify tem um buffer interno limitado. Se ninguém estiver lendo dele, novos eventos podem ser descartados silenciosamente assim que esse buffer se enche — a aplicação continuaria rodando normalmente, mas deixaria de perceber mudanças no arquivo observado, sem nenhum erro visível indicando o problema. É por isso que o padrão de sempre ter uma goroutine dedicada e ativa consumindo `Events` (e `Errors`) desde o início do programa é tão importante.

---

## 📖 Glossário

| Termo | Definição |
|---|---|
| **fsnotify** | Biblioteca Go que expõe, de forma multiplataforma, os mecanismos nativos de notificação de mudanças no sistema de arquivos de cada sistema operacional. |
| **inotify / kqueue** | Mecanismos nativos de notificação de sistema de arquivos do Linux (`inotify`) e do macOS/BSD (`kqueue`), usados pelo fsnotify por baixo dos panos. |
| **Watcher** | Instância (`fsnotify.Watcher`) responsável por observar um conjunto de caminhos (arquivos ou diretórios) e emitir eventos quando algo muda. |
| **Event** | Estrutura emitida pelo watcher a cada mudança detectada, contendo o caminho afetado (`Name`) e o tipo de operação (`Op`). |
| **Op (bitmask)** | Campo do `Event` que representa, como um conjunto de bits, quais operações ocorreram: `Create`, `Write`, `Remove`, `Rename`, `Chmod`. |
| **Hot-reload** | Padrão de recarregar dados (tipicamente configuração) em um processo já em execução, sem precisar reiniciá-lo. |
| **Rename atômico** | Técnica de salvar um arquivo escrevendo um temporário e depois renomeando-o por cima do original, evitando leituras parciais — mas que quebra watches registrados diretamente sobre o arquivo antigo. |
| **Canal (`chan`)** | Estrutura usada para comunicação e sincronização entre goroutines; aqui, `Events` e `Errors` são os canais expostos pelo watcher. |
| **Goroutine** | Unidade leve de execução concorrente do Go, usada neste exemplo para consumir eventos e erros do watcher em paralelo com o restante do programa. |

---

## 🚀 Próximos Passos

- [ ] **Imediato:** rode o exemplo e edite `config.json` algumas vezes, observando os campos impressos em `fmt.Println(config)` mudarem a cada salvamento.
- [ ] **Imediato:** modifique o `if` para também logar eventos de `Create` e `Remove` (`event.Op&fsnotify.Create == fsnotify.Create`), e observe a diferença de comportamento entre editores diferentes salvando o arquivo.
- [ ] **Intermediário:** implemente o watch no **diretório** em vez do arquivo (filtrando por `filepath.Base(event.Name)`), e teste com um editor que salva via rename atômico para confirmar que o watch continua funcionando onde antes falharia.
- [ ] **Intermediário:** adicione um `sync.RWMutex` para proteger a variável `config`, e escreva uma goroutine extra que leia `config` periodicamente para simular um consumidor concorrente.
- [ ] **Avançado:** substitua os `panic` por retornos de `error` tratados com log (ex.: `log/slog`), sem derrubar o processo quando a leitura ou o parse do JSON falhar temporariamente.
- [ ] **Avançado:** combine este exemplo com o padrão de [graceful shutdown](../22.1-graceful-shutdown/README.md): feche o canal `done` ao capturar `SIGINT`/`SIGTERM`, chamando `watcher.Close()` de forma explícita antes de encerrar o processo.
