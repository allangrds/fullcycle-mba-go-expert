# 🔐 Trabalhando com Módulos Privados em Go

## 📑 Sumário

- [🎯 O que você vai aprender](#-o-que-você-vai-aprender)
- [💡 Por que usar módulos privados?](#-por-que-usar-módulos-privados)
- [🔤 Conceitos Fundamentais de Go](#-conceitos-fundamentais-de-go)
  - [📦 O que são Módulos](#-o-que-são-módulos)
  - [📂 O que são Pacotes](#-o-que-são-pacotes)
  - [📄 Arquivo go.mod](#-arquivo-gomod)
  - [🔒 Arquivo go.sum](#-arquivo-gosum)
  - [🔄 GOPATH vs Go Modules](#-gopath-vs-go-modules)
- [🏷️ Versionamento em Go](#️-versionamento-em-go)
  - [📊 Semantic Versioning](#-semantic-versioning)
  - [🔢 Pseudo-Versions](#-pseudo-versions)
  - [🎯 Como Go Resolve Versões](#-como-go-resolve-versões)
  - [🏷️ Tags vs Commits](#️-tags-vs-commits)
- [🔐 Repositórios Privados](#-repositórios-privados)
  - [🌍 Repositórios Públicos vs Privados](#-repositórios-públicos-vs-privados)
  - [💼 Casos de Uso](#-casos-de-uso)
  - [🏢 Por que empresas usam módulos privados](#-por-que-empresas-usam-módulos-privados)
- [🔑 Autenticação e Configuração](#-autenticação-e-configuração)
  - [🔧 Variáveis de Ambiente GOPRIVATE](#-variáveis-de-ambiente-goprivate)
  - [🔐 Autenticação com GitHub (Token)](#-autenticação-com-github-token)
  - [🔑 Autenticação com SSH](#-autenticação-com-ssh)
  - [📝 Arquivo .netrc](#-arquivo-netrc)
  - [⚙️ Configuração do Git](#️-configuração-do-git)
- [🔍 Análise do Código](#-análise-do-código)
  - [📄 Entendendo o go.mod](#-entendendo-o-gomod)
  - [📄 Entendendo o go.sum](#-entendendo-o-gosum)
  - [📄 Entendendo o main.go](#-entendendo-o-maingo)
  - [🎨 Fluxo de Execução](#-fluxo-de-execução)
- [🛠️ Comandos Úteis](#️-comandos-úteis)
  - [📥 go get - Baixar dependências](#-go-get---baixar-dependências)
  - [🧹 go mod tidy - Limpar dependências](#-go-mod-tidy---limpar-dependências)
  - [⬇️ go mod download - Download de módulos](#️-go-mod-download---download-de-módulos)
  - [✅ go mod verify - Verificar integridade](#-go-mod-verify---verificar-integridade)
  - [🗑️ go clean -modcache - Limpar cache](#️-go-clean--modcache---limpar-cache)
- [🐛 Troubleshooting Comum](#-troubleshooting-comum)
  - [❌ Erro: 410 Gone](#-erro-410-gone)
  - [❌ Erro: Authentication Failed](#-erro-authentication-failed)
  - [❌ Erro: Checksum Mismatch](#-erro-checksum-mismatch)
  - [❌ Erro: Cannot find module](#-erro-cannot-find-module)
  - [❌ Erro: Access Denied](#-erro-access-denied)
- [🎨 Design Patterns e Boas Práticas](#-design-patterns-e-boas-práticas)
  - [📌 Versionamento Semântico](#-versionamento-semântico)
  - [📁 Módulos Internos (internal/)](#-módulos-internos-internal)
  - [🏗️ Mono-repos vs Multi-repos](#️-mono-repos-vs-multi-repos)
  - [🔄 Estratégias de Atualização](#-estratégias-de-atualização)
- [💼 Aplicações Práticas](#-aplicações-práticas)
  - [📚 Bibliotecas Internas da Empresa](#-bibliotecas-internas-da-empresa)
  - [🔧 SDKs Proprietários](#-sdks-proprietários)
  - [🏗️ Microserviços Compartilhados](#️-microserviços-compartilhados)
- [📖 Glossário](#-glossário)
- [🎓 Conceitos Aprendidos](#-conceitos-aprendidos)
- [📚 Próximos Passos](#-próximos-passos)

---

## 🎯 O que você vai aprender

Nesta aula, você vai aprender a trabalhar com **módulos privados** em Go, um conceito essencial para qualquer desenvolvedor que trabalha em empresas ou projetos que precisam compartilhar código de forma segura e controlada.

Você irá entender:
- Como Go gerencia dependências através de módulos
- O que são repositórios privados e quando usá-los
- Como configurar autenticação para acessar módulos privados
- Como versionar e distribuir código interno da sua empresa
- Troubleshooting de problemas comuns com módulos privados

Este é um conhecimento **fundamental** para trabalhar com Go em ambientes corporativos e projetos reais.

---

## 💡 Por que usar módulos privados?

### 🍕 Analogia do Mundo Real

Imagine que você tem uma **receita secreta** de pizza que faz o sucesso do seu restaurante:

**🌍 Receita Pública** (Módulo Público):
```
Qualquer pessoa pode:
- Ver a receita completa
- Copiar os ingredientes
- Usar em seu próprio restaurante
- Modificar e redistribuir
```

**🔐 Receita Privada** (Módulo Privado):
```
Apenas funcionários autorizados podem:
- Acessar a receita
- Usar em filiais da empresa
- Manter o diferencial competitivo
- Controlar quem tem acesso
```

### ✅ Vantagens dos Módulos Privados

| Vantagem | Descrição |
|----------|-----------|
| 🔒 **Segurança** | Código sensível não fica exposto publicamente |
| 💼 **Propriedade Intelectual** | Protege algoritmos e lógica de negócio |
| 🎯 **Controle de Acesso** | Apenas membros autorizados podem usar |
| 🔄 **Reutilização Interna** | Compartilhar código entre projetos da empresa |
| 📦 **Gestão Centralizada** | Atualizar bibliotecas em um único lugar |

### ⚠️ Considerações

- ❌ Requer configuração de autenticação
- ❌ Acesso restrito à equipe/empresa
- ❌ Maior complexidade de setup inicial
- ✅ Essencial para ambientes corporativos
- ✅ Permite versionamento controlado de código interno

---

## 🔤 Conceitos Fundamentais de Go

Antes de trabalhar com módulos privados, é importante entender os conceitos básicos do sistema de módulos do Go.

### 📦 O que são Módulos

**Módulo** é a unidade de versionamento e distribuição de código em Go.

```
📦 Módulo = Coleção de pacotes relacionados versionados juntos
```

**Analogia**: Pense em um módulo como um **livro**:
- O livro tem um **título** (nome do módulo)
- O livro tem uma **edição/versão** (v1.2.3)
- O livro contém vários **capítulos** (pacotes)
- Você pode citar/usar o livro em sua pesquisa (importar)

**Exemplo de Módulo**:
```
github.com/devfullcycle/fcutils-secret
├── pkg/
│   └── events/          ← Pacote events
│       └── dispatcher.go
├── internal/
│   └── helpers/         ← Pacote helpers (interno)
└── go.mod               ← Define o módulo
```

### 📂 O que são Pacotes

**Pacote** é um conjunto de arquivos `.go` no mesmo diretório que trabalham juntos.

```
📂 Pacote = Diretório com arquivos .go
           Todos começam com: package nome_do_pacote
```

**Exemplo de Pacote**:
```go
// arquivo: events/dispatcher.go
package events

type EventDispatcher struct {
    handlers map[string][]EventHandler
}

func NewEventDispatcher() *EventDispatcher {
    return &EventDispatcher{
        handlers: make(map[string][]EventHandler),
    }
}
```

**Diferença Módulo vs Pacote**:
```
Módulo    = github.com/devfullcycle/fcutils-secret
Pacote    = github.com/devfullcycle/fcutils-secret/pkg/events
           └─────────────────┬────────────────┘ └──┬──┘
                      Nome do Módulo              Caminho do Pacote
```

### 📄 Arquivo go.mod

O arquivo `go.mod` é o **manifesto** do seu módulo. Ele define:

```go
module teste                    // Nome do seu módulo
go 1.19                        // Versão mínima do Go

require (
    // Dependências diretas
    github.com/devfullcycle/fcutils-secret v0.0.0-20221027133857-ba9b1434aac5
)
```

**Componentes do go.mod**:

| Componente | Descrição | Exemplo |
|------------|-----------|---------|
| `module` | Nome do módulo | `module teste` |
| `go` | Versão do Go | `go 1.19` |
| `require` | Dependências necessárias | `require github.com/...` |
| `replace` | Substituir dependências | `replace old => new` |
| `exclude` | Excluir versões | `exclude github.com/... v1.0.0` |

**📍 Importante**: O `go.mod` lista apenas dependências **diretas** (que você importa diretamente no código).

### 🔒 Arquivo go.sum

O arquivo `go.sum` contém **checksums criptográficos** de todas as dependências (diretas e transitivas).

```
github.com/devfullcycle/fcutils-secret v0.0.0-20221027133857-ba9b1434aac5 h1:aSx0qUKAue92XQMRh+Yxv3/fhHgHaXKoHRU6pcgI/Xs=
github.com/devfullcycle/fcutils-secret v0.0.0-20221027133857-ba9b1434aac5/go.mod h1:tiJ328pm49A5W4HTXM0ncW3bq61mgGJ8XugwihQ2eAU=
```

**Formato do go.sum**:
```
[módulo] [versão] [algoritmo]:[hash]
```

**Por que existe?**

🎯 **Segurança e Integridade**:
```
1. Primeira vez que baixa: Go calcula hash e salva no go.sum
2. Próximas vezes: Go verifica se o hash bate
3. Se diferente: ERRO! ⛔ Alguém alterou o código!
```

**Analogia**: É como o **lacre de um remédio**:
- ✅ Lacre intacto = Conteúdo não foi adulterado
- ❌ Lacre rompido = PERIGO! Não use!

**Dependências Transitivas**:
```
Seu projeto
  └── devfullcycle/fcutils-secret (dependência direta)
       └── stretchr/testify (dependência transitiva)
            └── davecgh/go-spew (dependência transitiva)
```

Todas aparecem no `go.sum`! 📦

### 🔄 GOPATH vs Go Modules

**Evolução do Go**:

```
Go < 1.11        →  GOPATH (antigo)
Go >= 1.11       →  Go Modules (moderno) ✅
```

**🗂️ GOPATH (Modo Antigo)**:
```
Problemas:
❌ Todo código tinha que estar em $GOPATH/src
❌ Sem versionamento de dependências
❌ Conflitos entre projetos
❌ Difícil reproduzir builds

Estrutura obrigatória:
$GOPATH/
  └── src/
      └── github.com/
          └── usuario/
              └── projeto/
```

**📦 Go Modules (Modo Moderno)**:
```
Vantagens:
✅ Projetos podem estar em qualquer lugar
✅ Versionamento explícito
✅ Builds reproduzíveis
✅ Cache de dependências
✅ Suporte a módulos privados

Estrutura livre:
/home/usuario/projetos/
  └── meu-app/
      ├── go.mod
      ├── go.sum
      └── main.go
```

**💡 Dica**: Sempre use Go Modules (go.mod) em projetos novos!

---

## 🏷️ Versionamento em Go

### 📊 Semantic Versioning

Go usa **Semantic Versioning** (SemVer) para versionar módulos:

```
v1.2.3
│ │ │
│ │ └─── PATCH   (correções de bugs)
│ └───── MINOR   (novas funcionalidades compatíveis)
└─────── MAJOR   (mudanças incompatíveis)
```

**Exemplos**:

| Mudança | Versão Antiga | Versão Nova | Tipo |
|---------|---------------|-------------|------|
| Corrigir bug | v1.2.3 | v1.2.4 | PATCH |
| Adicionar função nova | v1.2.4 | v1.3.0 | MINOR |
| Mudar assinatura de função | v1.3.0 | v2.0.0 | MAJOR |

**Regra de Compatibilidade**:
```
v1.x.x → v1.y.z  ✅ Compatível (minor/patch)
v1.x.x → v2.0.0  ❌ Incompatível (major)
```

**💡 Para versões v2+**, Go exige sufixo no module path:
```go
// go.mod
module github.com/user/repo/v2  // ← /v2 obrigatório
```

### 🔢 Pseudo-Versions

Quando um commit **não tem tag**, Go cria uma **pseudo-version**:

```
v0.0.0-20221027133857-ba9b1434aac5
│  │  │  │               │
│  │  │  │               └─── Hash do commit (primeiros 12 chars)
│  │  │  └──────────────────── Timestamp UTC (YYYYMMDDHHMMSS)
│  │  └─────────────────────── Build number
└──┴──────────────────────────── Versão base
```

**Componentes**:
```
v0.0.0         = Versão base (quando não há tag anterior)
20221027       = Data: 27 de Outubro de 2022
133857         = Hora: 13:38:57 UTC
ba9b1434aac5   = Hash abreviado do commit Git
```

**Quando aparecem?**:
```
✅ Repositório sem tags
✅ Desenvolvimento ativo (HEAD)
✅ Usar commits específicos não taggeados
```

**Exemplo Real**:
```go
require github.com/devfullcycle/fcutils-secret v0.0.0-20221027133857-ba9b1434aac5
                                                 └─────────────────┬────────────────┘
                                                         Pseudo-version
```

### 🎯 Como Go Resolve Versões

**Fluxo de Resolução**:

```
1. go get github.com/user/repo
   │
   ├─→ 2. Busca tags no repositório Git
   │      │
   │      ├─→ 3a. Achou tags? → Usa a maior versão compatível
   │      │
   │      └─→ 3b. Sem tags? → Gera pseudo-version do commit mais recente
   │
   └─→ 4. Baixa código
       │
       └─→ 5. Calcula hash e salva em go.sum
```

**Exemplo Prático**:
```bash
# Instalar versão específica
go get github.com/user/repo@v1.2.3

# Instalar commit específico
go get github.com/user/repo@abc1234

# Instalar branch
go get github.com/user/repo@main

# Instalar versão mais recente
go get -u github.com/user/repo
```

### 🏷️ Tags vs Commits

**Com Tags (Recomendado)**:
```bash
# No repositório do módulo:
git tag v1.0.0
git push origin v1.0.0

# Resultado no go.mod:
require github.com/user/repo v1.0.0  ✅ Limpo e claro
```

**Sem Tags (Pseudo-version)**:
```bash
# Nenhuma tag no repositório

# Resultado no go.mod:
require github.com/user/repo v0.0.0-20221027133857-ba9b1434aac5  ⚠️ Difícil de ler
```

**💡 Boa Prática**: Sempre use tags semânticas em módulos de produção!

---

## 🔐 Repositórios Privados

### 🌍 Repositórios Públicos vs Privados

**📖 Repositório Público**:
```
✅ Qualquer pessoa pode:
   - Ver o código
   - Clonar o repositório
   - Usar como dependência
   - Contribuir (via PR)

Exemplos:
- github.com/gin-gonic/gin
- github.com/stretchr/testify
- github.com/gorilla/mux
```

**🔐 Repositório Privado**:
```
🔒 Apenas usuários autorizados podem:
   - Ver o código
   - Clonar o repositório
   - Usar como dependência
   - Acessar via autenticação

Exemplos:
- github.com/sua-empresa/biblioteca-interna
- github.com/devfullcycle/fcutils-secret  ← Usado nesta aula
- gitlab.com/empresa-privada/sdk-pagamentos
```

**Comparação**:

| Aspecto | Público | Privado |
|---------|---------|---------|
| Acesso | Qualquer pessoa | Apenas autorizados |
| Autenticação | Não necessária | **Obrigatória** |
| Custo | Grátis | Geralmente pago |
| Uso | Open source, bibliotecas públicas | Código proprietário, interno |
| Exemplo | React, Vue, Gin | SDK da empresa, lógica de negócio |

### 💼 Casos de Uso

**Quando usar Módulos Privados?**

**1. 🏢 Bibliotecas Internas da Empresa**:
```go
// Utilitários compartilhados entre times
import "github.com/minhaempresa/utils/logger"
import "github.com/minhaempresa/utils/database"
import "github.com/minhaempresa/utils/auth"
```

**2. 🔧 SDKs Proprietários**:
```go
// SDK para integração com sistemas internos
import "github.com/minhaempresa/sdk-pagamentos"
import "github.com/minhaempresa/sdk-erp"
```

**3. 🎯 Lógica de Negócio**:
```go
// Algoritmos e regras específicas do negócio
import "github.com/minhaempresa/pricing-engine"
import "github.com/minhaempresa/fraud-detection"
```

**4. 🏗️ Microserviços Compartilhados**:
```go
// Código compartilhado entre microserviços
import "github.com/minhaempresa/proto-definitions"
import "github.com/minhaempresa/common-middleware"
```

### 🏢 Por que empresas usam módulos privados

**Razões Empresariais**:

```
1. 💰 Propriedade Intelectual
   └─→ Algoritmos proprietários não podem ser públicos

2. 🔒 Segurança
   └─→ Evitar exposição de credenciais ou lógica sensível

3. 🎯 Controle de Acesso
   └─→ Apenas funcionários/parceiros autorizados

4. 📦 Reutilização de Código
   └─→ DRY: Don't Repeat Yourself entre projetos

5. 🔄 Versionamento Interno
   └─→ Controlar quando e como bibliotecas são atualizadas

6. 🛡️ Conformidade Legal
   └─→ Regulações (LGPD, GDPR) podem exigir código privado
```

**Exemplo Real**:
```
🏦 Banco XYZ
└── Módulos Privados:
    ├── sdk-core-banking (regras bancárias)
    ├── fraud-detection (detecção de fraude)
    ├── kyc-validation (validação de clientes)
    └── crypto-utils (criptografia proprietária)

❌ NUNCA podem ser públicos!
✅ Compartilhados entre apps do banco
```

---

## 🔑 Autenticação e Configuração

Para usar módulos privados, você precisa configurar autenticação. Aqui estão todas as formas:

### 🔧 Variáveis de Ambiente GOPRIVATE

**GOPRIVATE** informa ao Go quais módulos são privados:

```bash
# Formato
export GOPRIVATE="github.com/sua-empresa/*,gitlab.com/seu-grupo/*"

# Exemplo
export GOPRIVATE="github.com/devfullcycle/*"
```

**O que faz?**:
```
✅ Desabilita proxy público (proxy.golang.org)
✅ Desabilita checksum database (sum.golang.org)
✅ Força download direto do repositório Git
```

**Outras variáveis relacionadas**:

| Variável | Função |
|----------|--------|
| `GOPRIVATE` | Atalho para `GONOPROXY` + `GONOSUMDB` |
| `GONOPROXY` | Módulos que não usam proxy |
| `GONOSUMDB` | Módulos que não usam checksum database |
| `GOPROXY` | Servidor proxy a usar (padrão: proxy.golang.org) |

**Exemplo completo**:
```bash
# .bashrc, .zshrc, ou .bash_profile
export GOPRIVATE="github.com/minhaempresa,gitlab.com/meugrupo"
export GOPROXY="https://proxy.golang.org,direct"
```

**💡 Dica**: Adicione ao seu shell profile para configuração permanente!

### 🔐 Autenticação com GitHub (Token)

**Método 1: Git Config (Recomendado)**

```bash
# 1. Criar Personal Access Token no GitHub
#    Settings → Developer settings → Personal access tokens → Tokens (classic)
#    Permissions: repo (Full control of private repositories)

# 2. Configurar Git para usar o token
git config --global url."https://<TOKEN>@github.com/".insteadOf "https://github.com/"

# Exemplo:
git config --global url."https://ghp_xxxxxxxxxxxxxxxxxxxx@github.com/".insteadOf "https://github.com/"
```

**Método 2: .netrc (Alternativo)**

Veja seção [📝 Arquivo .netrc](#-arquivo-netrc).

**Como obter o token**:
```
1. GitHub.com → Settings
2. Developer settings (final da barra lateral)
3. Personal access tokens → Tokens (classic)
4. Generate new token
5. Selecionar scopes:
   ✅ repo (acesso completo a repositórios)
6. Generate token
7. Copiar token (só mostra uma vez!)
```

**⚠️ Segurança**:
```
❌ NUNCA commite tokens no código
❌ NUNCA compartilhe tokens
✅ Use tokens com escopo mínimo necessário
✅ Rotacione tokens periodicamente
✅ Revogue tokens quando não precisar mais
```

### 🔑 Autenticação com SSH

**Configurar SSH para módulos privados**:

```bash
# 1. Gerar chave SSH (se não tiver)
ssh-keygen -t ed25519 -C "seu-email@exemplo.com"

# 2. Adicionar chave ao ssh-agent
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519

# 3. Copiar chave pública
cat ~/.ssh/id_ed25519.pub

# 4. Adicionar no GitHub:
#    Settings → SSH and GPG keys → New SSH key → Colar chave

# 5. Configurar Git para usar SSH
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

**Vantagens do SSH**:
```
✅ Mais seguro (chave criptográfica)
✅ Não expira (ao contrário de tokens)
✅ Não precisa digitar senha
✅ Padrão em ambientes corporativos
```

**Testar conexão SSH**:
```bash
ssh -T git@github.com
# Saída esperada: "Hi username! You've successfully authenticated..."
```

### 📝 Arquivo .netrc

O arquivo `.netrc` armazena credenciais para autenticação HTTP:

```bash
# Localização: ~/.netrc (Linux/Mac) ou %HOME%\_netrc (Windows)

# Criar arquivo
touch ~/.netrc
chmod 600 ~/.netrc  # Permissões restritas (importante!)

# Conteúdo
machine github.com
  login seu-usuario
  password ghp_seu_token_aqui

machine gitlab.com
  login seu-usuario
  password seu_token_gitlab
```

**Formato**:
```
machine <host>
  login <usuario>
  password <token>
```

**⚠️ Segurança do .netrc**:
```bash
# Arquivo DEVE ter permissão 600 (apenas você pode ler/escrever)
chmod 600 ~/.netrc

# Verificar permissões
ls -la ~/.netrc
# Saída esperada: -rw------- 1 usuario grupo 123 May 5 10:00 .netrc
```

**Vantagens**:
```
✅ Funciona automaticamente (Go lê o .netrc)
✅ Suporta múltiplos hosts
✅ Não precisa configurar git config
```

**Desvantagens**:
```
❌ Senha/token em texto plano
❌ Precisa proteger o arquivo
❌ Menos seguro que SSH
```

### ⚙️ Configuração do Git

**Configurações úteis para módulos privados**:

```bash
# 1. Usar HTTPS com token
git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

# 2. Usar SSH em vez de HTTPS
git config --global url."git@github.com:".insteadOf "https://github.com/"

# 3. Verificar configurações
git config --global --list | grep url

# 4. Remover configuração
git config --global --unset url."https://github.com/".insteadOf
```

**Configuração por repositório** (somente para um projeto):
```bash
# Dentro do diretório do projeto
git config url."https://${TOKEN}@github.com/".insteadOf "https://github.com/"
```

**Exemplo de configuração completa para empresa**:
```bash
# ~/.gitconfig
[url "https://TOKEN@github.com/minhaempresa/"]
    insteadOf = https://github.com/minhaempresa/

[url "git@github.com:minhaempresa/"]
    insteadOf = https://github.com/minhaempresa/
```

---

## 🔍 Análise do Código

Agora vamos analisar **linha por linha** o código desta aula.

### 📄 Entendendo o go.mod

```go
module teste
```
**O que significa?**:
- `module teste` → Nome do seu módulo
- Este é o identificador usado quando outros projetos importam seu código
- Geralmente usa-se URL do repositório: `module github.com/user/projeto`
- Aqui é apenas `teste` pois é um projeto de exemplo local

**💡 Em produção seria**:
```go
module github.com/minhaempresa/meu-projeto
```

---

```go
go 1.19
```
**O que significa?**:
- Versão **mínima** do Go necessária para compilar este módulo
- Código pode rodar em Go 1.19 ou superior
- Garante que features do Go 1.19 estão disponíveis

**Como definir**:
```bash
# Ao criar novo módulo
go mod init nome-do-modulo

# Go automaticamente detecta a versão instalada
# e coloca no go.mod
```

---

```go
require github.com/devfullcycle/fcutils-secret v0.0.0-20221027133857-ba9b1434aac5
```
**Decomposição completa**:

| Parte | Valor | Significado |
|-------|-------|-------------|
| Palavra-chave | `require` | Declara dependência |
| Caminho do módulo | `github.com/devfullcycle/fcutils-secret` | Localização do módulo |
| Versão | `v0.0.0-20221027133857-ba9b1434aac5` | Pseudo-version (sem tag) |

**Pseudo-version decomposta**:
```
v0.0.0-20221027133857-ba9b1434aac5
│      │               │
│      │               └─── ba9b1434aac5: Hash do commit
│      └───────────────────── 20221027133857: 27/Out/2022 13:38:57
└──────────────────────────── v0.0.0: Versão base
```

**Por que v0.0.0?**:
- O repositório `fcutils-secret` não tem tags Git
- Go gera automaticamente uma pseudo-version
- Se houvesse tag `v1.2.3`, apareceria: `v1.2.3` ou `v1.2.4-0.202210...` (se após a tag)

### 📄 Entendendo o go.sum

Vamos analisar cada linha do arquivo `go.sum`:

**Linha 1**:
```
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
```
**Decomposição**:
- `github.com/davecgh/go-spew` → Módulo (dependência transitiva)
- `v1.1.1` → Versão
- `h1:` → Algoritmo de hash (h1 = SHA-256 base64)
- `vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=` → Hash do conteúdo

**O que é go-spew?**:
- Biblioteca para debug e inspeção de estruturas
- Dependência transitiva (usada por `testify`)

---

**Linha 2**:
```
github.com/devfullcycle/fcutils-secret v0.0.0-20221027133857-ba9b1434aac5 h1:aSx0qUKAue92XQMRh+Yxv3/fhHgHaXKoHRU6pcgI/Xs=
```
**Decomposição**:
- Nossa dependência **direta** principal
- Hash do **código fonte** do módulo
- Garante que o código não foi alterado

---

**Linha 3**:
```
github.com/devfullcycle/fcutils-secret v0.0.0-20221027133857-ba9b1434aac5/go.mod h1:tiJ328pm49A5W4HTXM0ncW3bq61mgGJ8XugwihQ2eAU=
```
**Decomposição**:
- Mesma dependência, mas hash do **arquivo go.mod** dela
- `/go.mod` → Indica que é hash do manifesto, não do código
- Garante que as dependências declaradas não mudaram

**Por que dois hashes para o mesmo módulo?**
```
1️⃣ Hash do código (.zip)    → Garante integridade do código
2️⃣ Hash do go.mod          → Garante que dependências não mudaram
```

---

**Linhas seguintes** (dependências transitivas):
```
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/stretchr/objx v0.5.0 h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=
github.com/stretchr/testify v1.8.1 h1:w7B6lhMri9wdJUVmEZPGGhZzrYTPvgJArz7wNPgYKsk=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
```

**Árvore de Dependências**:
```
seu-projeto
└── devfullcycle/fcutils-secret (direta)
    └── stretchr/testify (transitiva - para testes)
        ├── stretchr/objx (transitiva)
        ├── pmezard/go-difflib (transitiva)
        └── gopkg.in/yaml.v3 (transitiva)
```

**💡 Importante**: Você NÃO importa essas bibliotecas diretamente, mas elas aparecem no `go.sum` para garantir integridade completa da árvore de dependências!

### 📄 Entendendo o main.go

```go
package main
```
**O que significa?**:
- `package main` → Declara que este é um pacote executável
- Todo programa Go executável **deve** ter um `package main`
- Se fosse `package utils`, seria uma biblioteca (não executável)

**Diferença**:
```go
package main  → Gera executável (pode rodar: go run main.go)
package utils → Gera biblioteca (apenas importável)
```

---

```go
import (
	"fmt"
	"github.com/devfullcycle/fcutils-secret/pkg/events"
)
```
**Decomposição**:

**Import 1**: `"fmt"`
- Pacote da **biblioteca padrão** do Go
- Fornece funções de formatação e impressão
- Não precisa de `go get` (já vem com Go)

**Import 2**: `"github.com/devfullcycle/fcutils-secret/pkg/events"`
```
github.com/devfullcycle/fcutils-secret/pkg/events
└───────────────┬──────────────────┘ └─┬─┘ └──┬──┘
          Módulo (go.mod)           Path Pacote
```

**Estrutura**:
- **Módulo**: `github.com/devfullcycle/fcutils-secret`
- **Caminho dentro do módulo**: `pkg/events`
- **Pacote usado**: `events` (usado como `events.NewEventDispatcher()`)

**Onde está esse código?**:
```
Repositório privado no GitHub:
github.com/devfullcycle/fcutils-secret
└── pkg/
    └── events/
        ├── event_dispatcher.go    ← Aqui está NewEventDispatcher()
        ├── event_handler.go
        └── event.go
```

---

```go
func main() {
```
**O que significa?**:
- Ponto de entrada do programa
- **Toda aplicação Go executável deve ter `func main()`**
- É a primeira função chamada quando o programa roda

---

```go
   ed := events.NewEventDispatcher()
```
**Decomposição completa**:

**Sintaxe**:
```go
ed := events.NewEventDispatcher()
│      │      │
│      │      └─────────────── Função/método construtora
│      └────────────────────── Pacote (do import)
└───────────────────────────── Variável (short declaration)
```

**O que faz?**:
1. Chama função `NewEventDispatcher()` do pacote `events`
2. Essa função retorna um ponteiro `*EventDispatcher`
3. Armazena o resultado na variável `ed`

**Código dentro do módulo privado** (aproximadamente):
```go
// Dentro de: github.com/devfullcycle/fcutils-secret/pkg/events/event_dispatcher.go

type EventDispatcher struct {
    handlers map[string][]EventHandler
}

func NewEventDispatcher() *EventDispatcher {
    return &EventDispatcher{
        handlers: make(map[string][]EventHandler),
    }
}
```

**💡 Pattern Construtor**: `NewXxx()` é o padrão Go para criar instâncias (equivalente a construtores em outras linguagens).

---

```go
   fmt.Println(ed)
```
**O que faz?**:
- `fmt.Println()` → Imprime no console com quebra de linha
- `ed` → O objeto EventDispatcher criado
- Vai imprimir algo como: `&{map[]}`

**Saída esperada**:
```
&{map[]}
  │  │
  │  └──────── Map vazio (sem handlers registrados)
  └─────────── Ponteiro (&)
```

---

```go
}
```
**Fim da função `main()`**.

### 🎨 Fluxo de Execução

**Diagrama completo do fluxo**:

```
┌─────────────────────────────────────────────┐
│ 1. go run main.go                           │
└────────────────┬────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────┐
│ 2. Go lê go.mod                             │
│    - Vê: require github.com/devfullcycle... │
└────────────────┬────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────┐
│ 3. Go verifica se módulo está no cache     │
│    - Cache: $GOPATH/pkg/mod/               │
└────────────────┬────────────────────────────┘
                 │
          ┌──────┴──────┐
          ▼             ▼
    ┌─────────┐   ┌──────────────┐
    │ Existe  │   │ Não existe   │
    └────┬────┘   └──────┬───────┘
         │                │
         │                ▼
         │     ┌────────────────────────────┐
         │     │ 4. Go tenta baixar módulo  │
         │     │    - Verifica GOPRIVATE    │
         │     │    - Usa autenticação      │
         │     │    - Clone do GitHub       │
         │     └──────┬──────────────────────┘
         │            │
         │            ▼
         │     ┌────────────────────────────┐
         │     │ 5. Verifica integridade    │
         │     │    - Calcula hash          │
         │     │    - Compara com go.sum    │
         │     └──────┬──────────────────────┘
         │            │
         └────────┬───┘
                  │
                  ▼
         ┌────────────────────────────┐
         │ 6. Compila o programa      │
         │    - Linka dependências    │
         │    - Gera executável       │
         └────────┬───────────────────┘
                  │
                  ▼
         ┌────────────────────────────┐
         │ 7. Executa func main()     │
         └────────┬───────────────────┘
                  │
                  ▼
         ┌────────────────────────────┐
         │ 8. events.NewEventDispatcher() │
         │    - Cria struct           │
         │    - Inicializa map        │
         │    - Retorna ponteiro      │
         └────────┬───────────────────┘
                  │
                  ▼
         ┌────────────────────────────┐
         │ 9. fmt.Println(ed)         │
         │    - Imprime: &{map[]}     │
         └────────┬───────────────────┘
                  │
                  ▼
         ┌────────────────────────────┐
         │ 10. Programa termina       │
         │     - Exit code: 0         │
         └────────────────────────────┘
```

**Possíveis pontos de falha**:

| Passo | Erro Possível | Solução |
|-------|---------------|---------|
| 3 | Módulo não no cache | Normal - Go vai baixar |
| 4 | 410 Gone | Ver [Troubleshooting](#-troubleshooting-comum) |
| 4 | Authentication failed | Configurar GOPRIVATE e credenciais |
| 5 | Checksum mismatch | Remover cache e baixar novamente |
| 6 | Import errors | Verificar import paths |

---

## 🛠️ Comandos Úteis

### 📥 go get - Baixar dependências

**Comando básico**:
```bash
# Adicionar/atualizar dependência
go get github.com/user/repo

# Versão específica
go get github.com/user/repo@v1.2.3

# Commit específico
go get github.com/user/repo@abc1234

# Branch específica
go get github.com/user/repo@main

# Versão mais recente
go get -u github.com/user/repo
```

**Flags úteis**:

| Flag | Descrição | Exemplo |
|------|-----------|---------|
| `-u` | Atualizar para versão mais recente | `go get -u` |
| `-u=patch` | Atualizar apenas patches | `go get -u=patch` |
| `-t` | Incluir dependências de teste | `go get -t` |
| `-d` | Apenas baixar (não instalar) | `go get -d` |
| `@none` | Remover dependência | `go get github.com/user/repo@none` |

**Exemplos práticos**:

```bash
# Instalar módulo privado (com GOPRIVATE configurado)
export GOPRIVATE="github.com/devfullcycle/*"
go get github.com/devfullcycle/fcutils-secret

# Atualizar todas dependências
go get -u ./...

# Atualizar apenas patches de segurança
go get -u=patch ./...

# Remover dependência
go get github.com/old/package@none
```

### 🧹 go mod tidy - Limpar dependências

**O que faz?**:
```
✅ Adiciona dependências que faltam (usadas mas não declaradas)
✅ Remove dependências não usadas
✅ Atualiza go.sum
✅ Organiza go.mod
```

**Uso básico**:
```bash
# Limpar e organizar dependências
go mod tidy

# Ver o que seria mudado (sem modificar)
go mod tidy -v
```

**Quando usar?**:
```
✅ Após adicionar novos imports
✅ Após remover imports
✅ Antes de commit (boa prática)
✅ Quando go.mod está desorganizado
✅ CI/CD pipelines (validação)
```

**Exemplo de uso**:
```bash
# 1. Você adiciona import no código
import "github.com/gin-gonic/gin"

# 2. Roda go mod tidy
go mod tidy

# 3. go.mod é atualizado automaticamente
# require github.com/gin-gonic/gin v1.9.1
```

**Flags**:

| Flag | Descrição |
|------|-----------|
| `-v` | Modo verbose (mostra o que está fazendo) |
| `-e` | Continua mesmo com erros |
| `-go=1.19` | Especifica versão do Go |

### ⬇️ go mod download - Download de módulos

**O que faz?**:
- Baixa módulos para o cache **sem compilar**
- Útil para CI/CD (cachear dependências)
- Verifica integridade com go.sum

**Uso**:
```bash
# Baixar todas dependências
go mod download

# Baixar módulo específico
go mod download github.com/gin-gonic/gin

# Download com verificação de JSON
go mod download -json
```

**Saída JSON** (`-json`):
```json
{
  "Path": "github.com/gin-gonic/gin",
  "Version": "v1.9.1",
  "Info": "/go/pkg/mod/cache/download/github.com/gin-gonic/gin/@v/v1.9.1.info",
  "GoMod": "/go/pkg/mod/cache/download/github.com/gin-gonic/gin/@v/v1.9.1.mod",
  "Zip": "/go/pkg/mod/cache/download/github.com/gin-gonic/gin/@v/v1.9.1.zip",
  "Sum": "h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg="
}
```

**Uso em CI/CD**:
```yaml
# .github/workflows/ci.yml
- name: Download dependencies
  run: go mod download

- name: Cache Go modules
  uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

### ✅ go mod verify - Verificar integridade

**O que faz?**:
```
✅ Verifica se módulos no cache não foram modificados
✅ Compara hashes com go.sum
✅ Detecta adulteração
```

**Uso**:
```bash
# Verificar todos os módulos
go mod verify
```

**Saídas possíveis**:

```bash
# ✅ Sucesso
all modules verified

# ❌ Falha
verifying github.com/gin-gonic/gin@v1.9.1: checksum mismatch
    downloaded: h1:OLD_HASH
    go.sum:     h1:EXPECTED_HASH
```

**Quando usar?**:
```
✅ Após clonar repositório
✅ Em CI/CD (garantir integridade)
✅ Após limpar cache
✅ Investigação de segurança
✅ Antes de deploy em produção
```

**Exemplo em script de deploy**:
```bash
#!/bin/bash
echo "Verificando integridade das dependências..."
if ! go mod verify; then
    echo "❌ ERRO: Dependências comprometidas!"
    exit 1
fi
echo "✅ Dependências verificadas com sucesso"
```

### 🗑️ go clean -modcache - Limpar cache

**O que faz?**:
- Remove **todo** cache de módulos
- Cache fica em `$GOPATH/pkg/mod/`
- Força re-download de todas dependências

**Uso**:
```bash
# Limpar todo cache de módulos
go clean -modcache

# Limpar cache e build cache
go clean -cache -modcache
```

**Tamanho do cache**:
```bash
# Ver tamanho do cache
du -sh $GOPATH/pkg/mod/
# Exemplo: 1.2G

# Listar módulos no cache
ls $GOPATH/pkg/mod/cache/download/
```

**Quando usar?**:
```
✅ Erro de checksum mismatch
✅ Módulos corrompidos
✅ Liberar espaço em disco
✅ Forçar re-download
✅ Troubleshooting de dependências
```

**⚠️ Cuidado**: Vai baixar **tudo** novamente na próxima build!

**Exemplo de troubleshooting**:
```bash
# 1. Erro de checksum
go build
# Error: verifying module: checksum mismatch

# 2. Limpar cache
go clean -modcache

# 3. Baixar novamente
go mod download

# 4. Tentar build novamente
go build
```

**Outros comandos `go clean`**:

| Comando | Descrição |
|---------|-----------|
| `go clean` | Remove binários compilados |
| `go clean -cache` | Limpa cache de build |
| `go clean -testcache` | Limpa cache de testes |
| `go clean -modcache` | Limpa cache de módulos |
| `go clean -i` | Remove pacotes instalados |

---

## 🐛 Troubleshooting Comum

### ❌ Erro: 410 Gone

**Mensagem de erro**:
```
go: downloading github.com/devfullcycle/fcutils-secret v0.0.0-20221027133857-ba9b1434aac5
go: github.com/devfullcycle/fcutils-secret@v0.0.0-20221027133857-ba9b1434aac5: reading https://sum.golang.org/lookup/github.com/devfullcycle/fcutils-secret@v0.0.0-20221027133857-ba9b1434aac5: 	410 Gone
```

**O que significa?**:
- Go está tentando buscar checksum no **sum.golang.org** (banco de dados público)
- Módulos **privados** não estão no sum.golang.org
- HTTP 410 Gone = "Este recurso não existe e nunca existirá"

**Solução**:
```bash
# Configurar GOPRIVATE para pular sum.golang.org
export GOPRIVATE="github.com/devfullcycle/*"

# Ou adicionar permanentemente ao shell
echo 'export GOPRIVATE="github.com/devfullcycle/*"' >> ~/.bashrc
source ~/.bashrc

# Tentar novamente
go get github.com/devfullcycle/fcutils-secret
```

**Explicação visual**:
```
SEM GOPRIVATE:
go get → sum.golang.org → 410 Gone ❌

COM GOPRIVATE:
go get → github.com (direto) → ✅ Sucesso
```

### ❌ Erro: Authentication Failed

**Mensagens possíveis**:
```
fatal: could not read Username for 'https://github.com': terminal prompts disabled
```
```
fatal: Authentication failed for 'https://github.com/devfullcycle/fcutils-secret/'
```
```
go: github.com/devfullcycle/fcutils-secret@v0.0.0-...: invalid version: git fetch failed
```

**Causa**:
- Repositório é privado
- Go não tem permissão para acessar
- Credenciais não configuradas

**Soluções**:

**Solução 1: Git Config com Token**
```bash
# Obter token no GitHub (Settings → Developer settings → Personal access tokens)
# Escopo necessário: repo

# Configurar Git
git config --global url."https://SEU_TOKEN@github.com/".insteadOf "https://github.com/"

# Testar
go get github.com/devfullcycle/fcutils-secret
```

**Solução 2: SSH**
```bash
# Configurar Git para usar SSH
git config --global url."git@github.com:".insteadOf "https://github.com/"

# Garantir que SSH está configurado
ssh -T git@github.com
# Hi username! You've successfully authenticated

# Testar
go get github.com/devfullcycle/fcutils-secret
```

**Solução 3: .netrc**
```bash
# Criar ~/.netrc
cat > ~/.netrc << EOF
machine github.com
  login seu-usuario
  password ghp_seu_token_aqui
EOF

# Proteger arquivo
chmod 600 ~/.netrc

# Testar
go get github.com/devfullcycle/fcutils-secret
```

**Verificar configuração**:
```bash
# Ver configurações Git
git config --global --list | grep url

# Testar clone manual
git clone https://github.com/devfullcycle/fcutils-secret.git
```

### ❌ Erro: Checksum Mismatch

**Mensagem de erro**:
```
verifying github.com/gin-gonic/gin@v1.9.1: checksum mismatch
    downloaded: h1:Abc123...
    go.sum:     h1:Xyz789...
```

**Causa**:
- Código no cache difere do esperado em go.sum
- Possível adulteração
- Cache corrompido
- go.sum desatualizado

**Soluções**:

**Solução 1: Limpar cache**
```bash
# Limpar cache de módulos
go clean -modcache

# Baixar novamente
go mod download

# Verificar
go mod verify
```

**Solução 2: Atualizar go.sum**
```bash
# Se você confia no código baixado
rm go.sum
go mod tidy
```

**Solução 3: Verificar integridade**
```bash
# Verificar se go.sum está corrompido
go mod verify

# Reconstruir go.sum do zero
rm go.sum
go mod download
```

**⚠️ IMPORTANTE**:
```
Se o erro persistir após limpar cache:
❌ PODE ser ataque man-in-the-middle
❌ PODE ser código adulterado
✅ Investigar antes de ignorar!
```

### ❌ Erro: Cannot find module

**Mensagem de erro**:
```
go: github.com/inexistente/pacote@v1.0.0: reading https://proxy.golang.org/github.com/inexistente/pacote/@v/v1.0.0.info: 404 Not Found
```

**Causas possíveis**:

**1. Módulo não existe**
```bash
# Verificar se URL está correta
curl https://github.com/usuario/repo
# Se 404 → Repositório não existe
```

**2. Módulo é privado (mas GOPRIVATE não configurado)**
```bash
# Configurar GOPRIVATE
export GOPRIVATE="github.com/sua-empresa/*"
go get github.com/sua-empresa/modulo-privado
```

**3. Versão não existe**
```bash
# Listar versões disponíveis
go list -m -versions github.com/user/repo

# Usar versão que existe
go get github.com/user/repo@v1.2.0
```

**4. Nome do módulo mudou**
```bash
# Ver o verdadeiro nome do módulo
curl https://raw.githubusercontent.com/user/repo/main/go.mod

# Usar o nome correto do go.mod
go get nome-correto-do-modulo
```

**Solução geral**:
```bash
# 1. Verificar se repositório existe
# 2. Verificar GOPRIVATE (se privado)
# 3. Verificar autenticação (se privado)
# 4. Verificar nome do módulo no go.mod do repo
# 5. Verificar versões disponíveis
```

### ❌ Erro: Access Denied

**Mensagem de erro**:
```
fatal: could not read Username for 'https://github.com': terminal prompts disabled
Confirm the import path was entered correctly.
If this is a private repository, see https://golang.org/doc/faq#git_https for additional information.
```

**Causa**:
- Tentando acessar repositório privado sem autenticação
- Terminal prompts desabilitados (CI/CD)

**Soluções**:

**Para desenvolvimento local**:
```bash
# Opção 1: Usar SSH
git config --global url."git@github.com:".insteadOf "https://github.com/"

# Opção 2: Token no Git config
git config --global url."https://TOKEN@github.com/".insteadOf "https://github.com/"

# Opção 3: .netrc
echo "machine github.com login usuario password TOKEN" > ~/.netrc
chmod 600 ~/.netrc
```

**Para CI/CD (GitHub Actions)**:
```yaml
# .github/workflows/ci.yml
- name: Configure Git for private modules
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

- name: Download dependencies
  env:
    GOPRIVATE: github.com/sua-empresa/*
  run: go mod download
```

**Para Docker**:
```dockerfile
# Dockerfile
FROM golang:1.19

# Argumento de build
ARG GITHUB_TOKEN

# Configurar Git
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

# Build
WORKDIR /app
COPY go.* ./
RUN go mod download

# Build da aplicação
CMD ["go", "run", "main.go"]
```

**Build do Docker**:
```bash
docker build --build-arg GITHUB_TOKEN=$GITHUB_TOKEN -t meu-app .
```

---

## 🎨 Design Patterns e Boas Práticas

### 📌 Versionamento Semântico

**Regras de Versionamento**:

```
v{MAJOR}.{MINOR}.{PATCH}

MAJOR: Mudanças incompatíveis (breaking changes)
MINOR: Novas funcionalidades compatíveis
PATCH: Correções de bugs
```

**Exemplos de mudanças**:

| Mudança | Tipo | Nova Versão | Razão |
|---------|------|-------------|-------|
| Corrigir bug | PATCH | v1.2.3 → v1.2.4 | Apenas correção |
| Adicionar função nova | MINOR | v1.2.4 → v1.3.0 | Nova feature compatível |
| Mudar assinatura de função | MAJOR | v1.3.0 → v2.0.0 | Breaking change |
| Remover função pública | MAJOR | v1.3.0 → v2.0.0 | Breaking change |
| Adicionar parâmetro opcional | MINOR | v1.3.0 → v1.4.0 | Compatível |
| Tornar parâmetro obrigatório | MAJOR | v1.4.0 → v2.0.0 | Breaking change |

**Como criar releases**:

```bash
# 1. Fazer mudanças no código
git add .
git commit -m "feat: adicionar nova funcionalidade"

# 2. Criar tag semântica
git tag v1.2.0

# 3. Push da tag
git push origin v1.2.0

# 4. Agora pode ser usada:
# go get github.com/user/repo@v1.2.0
```

**Convenção de commits** (Conventional Commits):
```bash
feat: nova funcionalidade (MINOR)
fix: correção de bug (PATCH)
docs: documentação (PATCH)
refactor: refatoração sem mudar API (PATCH)
perf: melhoria de performance (PATCH)
test: adicionar testes (PATCH)

# Breaking change (MAJOR):
feat!: mudar assinatura de função
```

**Compatibilidade em Go**:
```
v1.x.x → v1.y.z  ✅ Atualização segura
v1.x.x → v2.0.0  ⚠️ Requer mudanças no código

# Para v2+, modificar go.mod:
module github.com/user/repo/v2  // ← /v2 obrigatório
```

### 📁 Módulos Internos (internal/)

**Diretório `internal/`** é especial em Go:

```
✅ Código em internal/ só pode ser importado pelo módulo pai
❌ Outros módulos NÃO podem importar
```

**Estrutura recomendada**:
```
github.com/minhaempresa/meu-modulo/
├── go.mod
├── pkg/                    ← Público (qualquer um pode importar)
│   └── api/
│       └── client.go
├── internal/               ← Privado (só este módulo pode usar)
│   ├── database/
│   │   └── queries.go
│   └── helpers/
│       └── utils.go
└── cmd/
    └── server/
        └── main.go
```

**Exemplo de uso**:

```go
// ✅ PERMITIDO: Mesmo módulo
// Em: github.com/minhaempresa/meu-modulo/cmd/server/main.go
import "github.com/minhaempresa/meu-modulo/internal/database"

// ❌ PROIBIDO: Outro módulo
// Em: github.com/outraempresa/outro-modulo/main.go
import "github.com/minhaempresa/meu-modulo/internal/database"
// Erro: use of internal package not allowed
```

**Quando usar**:

| Usar internal/ | Usar pkg/ |
|---|---|
| Lógica de implementação | API pública |
| Helpers internos | Interfaces públicas |
| Código que pode mudar | Código estável |
| Não deve ser reutilizado | Feito para reutilização |

**💡 Boa Prática**: Comece com `internal/`, e mova para `pkg/` apenas quando tiver certeza que deve ser público.

### 🏗️ Mono-repos vs Multi-repos

**Multi-repos** (Repositório por módulo):

```
github.com/empresa/autenticacao/     (repo 1)
github.com/empresa/logger/           (repo 2)
github.com/empresa/database/         (repo 3)
github.com/empresa/api-cliente/      (repo 4)
```

**Vantagens**:
```
✅ Versionamento independente
✅ Permissões granulares
✅ CI/CD separado por módulo
✅ Equipes independentes
```

**Desvantagens**:
```
❌ Sincronizar mudanças entre repos é difícil
❌ Mais complexo para testar integração
❌ Muitos repositórios para gerenciar
```

---

**Mono-repo** (Vários módulos em um repositório):

```
github.com/empresa/backend/
├── go.work                    ← Go Workspace
├── services/
│   ├── auth/
│   │   ├── go.mod
│   │   └── main.go
│   ├── api/
│   │   ├── go.mod
│   │   └── main.go
│   └── worker/
│       ├── go.mod
│       └── main.go
└── libs/
    ├── logger/
    │   └── go.mod
    └── database/
        └── go.mod
```

**Vantagens**:
```
✅ Mudanças atômicas (um commit afeta tudo)
✅ Fácil refatoração cross-modules
✅ Testes de integração simplificados
✅ Um CI/CD para todos
```

**Desvantagens**:
```
❌ Permissões menos granulares
❌ CI pode ficar lento
❌ Versionamento mais complexo
```

**Go Workspace** (go.work):
```go
// go.work
go 1.19

use (
    ./services/auth
    ./services/api
    ./services/worker
    ./libs/logger
    ./libs/database
)
```

**Recomendação**:

| Cenário | Escolha |
|---------|---------|
| Startup pequena | Mono-repo |
| Empresa grande, times separados | Multi-repos |
| Microserviços acoplados | Mono-repo |
| Bibliotecas independentes | Multi-repos |

### 🔄 Estratégias de Atualização

**Estratégia 1: Versões Fixas** (Conservador)
```go
// go.mod
require (
    github.com/gin-gonic/gin v1.9.1      // Versão exata
    github.com/lib/pq v1.10.9             // Versão exata
)
```
**Quando usar**: Produção estável, zero downtime

---

**Estratégia 2: Atualizar Patches** (Recomendado)
```bash
# Atualizar apenas patches (v1.2.3 → v1.2.4)
go get -u=patch ./...
```
**Quando usar**: Segurança, correções de bugs

---

**Estratégia 3: Atualizar Minor** (Balanceado)
```bash
# Atualizar minor e patches (v1.2.3 → v1.3.0)
go get -u ./...
```
**Quando usar**: Novas features, sem breaking changes

---

**Estratégia 4: Atualizar Major** (Agressivo)
```bash
# Atualizar manualmente para v2
go get github.com/user/repo/v2@latest

# Modificar imports no código
# import "github.com/user/repo/v2"
```
**Quando usar**: Migração planejada, com testes

---

**Automação com Dependabot** (GitHub):
```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    target-branch: "develop"
    versioning-strategy: "increase-if-necessary"
    labels:
      - "dependencies"
      - "go"
```

**Workflow recomendado**:
```
1. Dependabot abre PR semanalmente
2. CI roda testes automaticamente
3. Se verde: Deploy em staging
4. QA valida em staging
5. Se OK: Merge e deploy em produção
```

---

## 💼 Aplicações Práticas

### 📚 Bibliotecas Internas da Empresa

**Cenário**: Sua empresa tem 20 microserviços que precisam de funcionalidades comuns.

**Problema sem módulos privados**:
```
❌ Copiar/colar código em cada serviço
❌ Código duplicado = bugs duplicados
❌ Atualizar em 20 lugares quando muda algo
❌ Inconsistência entre serviços
```

**Solução com módulos privados**:

```
📦 github.com/minhaempresa/stdlib (módulo privado)
├── pkg/
│   ├── logger/          ← Logging padronizado
│   ├── database/        ← Pool de conexões
│   ├── auth/            ← Middleware de autenticação
│   ├── metrics/         ← Prometheus, tracing
│   └── errors/          ← Tratamento de erros
└── go.mod
```

**Uso nos microserviços**:
```go
// Serviço 1: API de Produtos
import (
    "github.com/minhaempresa/stdlib/pkg/logger"
    "github.com/minhaempresa/stdlib/pkg/database"
)

// Serviço 2: API de Pedidos
import (
    "github.com/minhaempresa/stdlib/pkg/logger"
    "github.com/minhaempresa/stdlib/pkg/auth"
)
```

**Vantagens**:
```
✅ DRY: Don't Repeat Yourself
✅ Atualizar em um lugar reflete em todos
✅ Testes centralizados
✅ Padrões consistentes
✅ Onboarding mais rápido (desenvolvedores novos)
```

**Exemplo de logger padronizado**:
```go
// github.com/minhaempresa/stdlib/pkg/logger/logger.go
package logger

import "go.uber.org/zap"

func New(service string) (*zap.Logger, error) {
    cfg := zap.NewProductionConfig()
    cfg.InitialFields = map[string]interface{}{
        "service": service,
        "company": "MinhaEmpresa",
    }
    return cfg.Build()
}
```

### 🔧 SDKs Proprietários

**Cenário**: Sua empresa oferece uma API para parceiros, mas quer facilitar a integração.

**Solução**: SDK privado em Go

```
📦 github.com/minhaempresa/sdk-pagamentos (módulo privado)
├── pkg/
│   ├── client/
│   │   └── client.go       ← Cliente HTTP pré-configurado
│   ├── auth/
│   │   └── oauth.go        ← Autenticação OAuth2
│   ├── payments/
│   │   ├── create.go       ← Criar pagamento
│   │   ├── cancel.go       ← Cancelar pagamento
│   │   └── refund.go       ← Estornar pagamento
│   └── models/
│       └── payment.go      ← Structs compartilhadas
└── go.mod
```

**Uso pelo parceiro**:
```go
package main

import (
    "github.com/minhaempresa/sdk-pagamentos/pkg/client"
    "github.com/minhaempresa/sdk-pagamentos/pkg/payments"
)

func main() {
    // Cliente pré-configurado
    c := client.New("api-key-secreta")
    
    // Criar pagamento (SDK abstrai complexidade)
    payment, err := payments.Create(c, &payments.CreateRequest{
        Amount:      10000, // R$ 100,00
        Currency:    "BRL",
        Description: "Compra de produto",
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Pagamento criado: %s\n", payment.ID)
}
```

**Vantagens**:
```
✅ Parceiros integram mais rápido
✅ Menos erros de integração
✅ Você controla versões
✅ Pode adicionar analytics
✅ Documentação no código
```

### 🏗️ Microserviços Compartilhados

**Cenário**: Arquitetura de microserviços com código compartilhado.

**Estrutura**:
```
📦 github.com/minhaempresa/proto (módulo privado)
├── proto/
│   ├── user.proto              ← Definições gRPC
│   ├── product.proto
│   └── order.proto
├── gen/                         ← Código gerado
│   ├── user/
│   │   └── user.pb.go
│   ├── product/
│   │   └── product.pb.go
│   └── order/
│       └── order.pb.go
└── go.mod
```

**Uso nos microserviços**:

```go
// Microserviço: user-service
import userpb "github.com/minhaempresa/proto/gen/user"

func (s *Server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
    // Implementação
}
```

```go
// Microserviço: order-service (chama user-service)
import (
    userpb "github.com/minhaempresa/proto/gen/user"
    orderpb "github.com/minhaempresa/proto/gen/order"
)

func (s *Server) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.Order, error) {
    // Buscar usuário via gRPC
    user, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{
        UserId: req.UserId,
    })
    // ...
}
```

**Vantagens**:
```
✅ Contratos compartilhados (gRPC/Protobuf)
✅ Type-safe entre serviços
✅ Versionar contratos
✅ Gerar código automaticamente
✅ Evitar duplicação de modelos
```

**CI/CD para gerar código**:
```yaml
# .github/workflows/proto.yml
name: Generate Proto
on:
  push:
    paths:
      - 'proto/**'
jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Generate Go code
        run: |
          protoc --go_out=gen --go-grpc_out=gen proto/*.proto
      - name: Commit generated code
        run: |
          git config user.name "GitHub Actions"
          git add gen/
          git commit -m "chore: regenerate proto code" || true
          git push
```

---

## 📖 Glossário

| Termo | Tradução | Significado |
|-------|----------|-------------|
| **Module** | Módulo | Unidade de versionamento e distribuição de código em Go |
| **Package** | Pacote | Diretório com arquivos `.go` relacionados |
| **Dependency** | Dependência | Módulo que seu projeto precisa para funcionar |
| **go.mod** | - | Arquivo manifesto do módulo (metadados e dependências) |
| **go.sum** | - | Arquivo com checksums criptográficos das dependências |
| **GOPATH** | - | Diretório de trabalho do Go (modo antigo, pré-modules) |
| **Go Modules** | Módulos do Go | Sistema moderno de gerenciamento de dependências (Go ≥ 1.11) |
| **Module Path** | Caminho do módulo | Identificador único do módulo (ex: github.com/user/repo) |
| **Import Path** | Caminho de import | Caminho completo para importar pacote (module path + subpath) |
| **Semantic Versioning (SemVer)** | Versionamento Semântico | Esquema v{MAJOR}.{MINOR}.{PATCH} |
| **Pseudo-version** | Pseudo-versão | Versão gerada para commits sem tag (v0.0.0-timestamp-hash) |
| **Tag** | Tag/Etiqueta | Marcador Git para releases (ex: v1.0.0) |
| **Private Repository** | Repositório Privado | Repositório Git com acesso restrito |
| **Public Repository** | Repositório Público | Repositório Git acessível a todos |
| **GOPRIVATE** | - | Variável de ambiente que lista módulos privados |
| **GOPROXY** | - | Servidor proxy para cache de módulos |
| **Checksum Database** | Base de checksums | sum.golang.org - valida integridade de módulos públicos |
| **Checksum** | Soma de verificação | Hash criptográfico para verificar integridade |
| **Direct Dependency** | Dependência Direta | Módulo que você importa diretamente no código |
| **Transitive Dependency** | Dependência Transitiva | Dependência da sua dependência (indireta) |
| **Module Cache** | Cache de módulos | $GOPATH/pkg/mod/ - onde Go armazena módulos baixados |
| **Personal Access Token (PAT)** | Token de Acesso Pessoal | Credencial para acessar repos privados via HTTPS |
| **SSH Key** | Chave SSH | Par de chaves criptográficas para autenticação |
| **.netrc** | - | Arquivo com credenciais HTTP |
| **Breaking Change** | Mudança Incompatível | Alteração que quebra compatibilidade (requer MAJOR bump) |
| **internal/** | - | Diretório especial com código privado ao módulo |
| **pkg/** | - | Convenção para código público/importável |
| **Mono-repo** | Monorepo | Um repositório com múltiplos módulos/projetos |
| **Multi-repo** | Múltiplos repos | Cada módulo em seu próprio repositório |
| **Go Workspace** | Workspace do Go | go.work - agrupa múltiplos módulos locais |
| **Vendor** | Vendor/Vendoring | Copiar dependências para dentro do projeto |
| **Replace Directive** | Diretiva Replace | Substituir módulo por versão local/alternativa |

---

## 🎓 Conceitos Aprendidos

Após completar esta aula, você agora sabe:

### ✅ Fundamentos
- ✅ O que são módulos e pacotes em Go
- ✅ Diferença entre go.mod e go.sum
- ✅ Como funciona o sistema de versionamento do Go
- ✅ GOPATH vs Go Modules
- ✅ Import paths e module paths

### ✅ Versionamento
- ✅ Semantic Versioning (MAJOR.MINOR.PATCH)
- ✅ O que são pseudo-versions e quando aparecem
- ✅ Como Go resolve versões de dependências
- ✅ Diferença entre tags e commits

### ✅ Repositórios Privados
- ✅ Diferença entre repositórios públicos e privados
- ✅ Casos de uso para módulos privados
- ✅ Por que empresas usam módulos privados
- ✅ Vantagens e desvantagens

### ✅ Autenticação
- ✅ Configurar GOPRIVATE
- ✅ Autenticação via Token (GitHub PAT)
- ✅ Autenticação via SSH
- ✅ Configurar arquivo .netrc
- ✅ Git config para módulos privados

### ✅ Comandos
- ✅ `go get` - baixar e atualizar dependências
- ✅ `go mod tidy` - limpar e organizar dependências
- ✅ `go mod download` - baixar para cache
- ✅ `go mod verify` - verificar integridade
- ✅ `go clean -modcache` - limpar cache

### ✅ Troubleshooting
- ✅ Resolver erro "410 Gone"
- ✅ Resolver "Authentication Failed"
- ✅ Resolver "Checksum Mismatch"
- ✅ Resolver "Cannot find module"
- ✅ Resolver "Access Denied"

### ✅ Boas Práticas
- ✅ Usar versionamento semântico em releases
- ✅ Organizar código com internal/ e pkg/
- ✅ Escolher entre mono-repo e multi-repos
- ✅ Estratégias de atualização de dependências
- ✅ Automação com Dependabot

### ✅ Aplicações Práticas
- ✅ Criar bibliotecas internas da empresa
- ✅ Distribuir SDKs proprietários
- ✅ Compartilhar código entre microserviços
- ✅ Gerenciar contratos gRPC/Protobuf

---

## 📚 Próximos Passos

Agora que você domina módulos privados, aqui estão os próximos tópicos para aprofundar:

### 🔍 Go Workspaces (go.work)

Trabalhar com múltiplos módulos localmente:

```bash
# Criar workspace
go work init

# Adicionar módulos ao workspace
go work use ./service-a
go work use ./service-b
go work use ./shared-lib

# Agora pode desenvolver os 3 módulos simultaneamente
# sem precisar fazer replace no go.mod
```

**Quando usar**:
- Desenvolvimento de mono-repo
- Trabalhar em biblioteca e consumidor ao mesmo tempo
- Testar mudanças antes de publicar

**Recursos**:
- [Tutorial: Getting started with multi-module workspaces](https://go.dev/doc/tutorial/workspaces)

### 🔄 Replace Directive

Substituir dependências temporariamente:

```go
// go.mod
module meu-projeto

require github.com/empresa/biblioteca v1.2.3

// Substituir por versão local (desenvolvimento)
replace github.com/empresa/biblioteca => ../biblioteca-local

// Substituir por fork
replace github.com/empresa/biblioteca => github.com/meufork/biblioteca v1.2.4

// Substituir por commit específico
replace github.com/empresa/biblioteca => github.com/empresa/biblioteca abc1234
```

**Casos de uso**:
- Testar mudanças em biblioteca antes de publicar
- Usar fork enquanto aguarda merge de PR
- Desenvolvimento local

### 🏪 Module Proxy Privado

Hospedar seu próprio proxy de módulos:

**Opções populares**:
- [Athens](https://docs.gomods.io/) - Proxy e registry de módulos
- [Artifactory](https://jfrog.com/artifactory/) - Universal artifact repository
- [Sonatype Nexus](https://www.sonatype.com/products/repository-oss) - Repository manager

**Vantagens**:
```
✅ Cache local de módulos públicos (builds mais rápidos)
✅ Hospedar módulos privados internamente
✅ Controle de acesso granular
✅ Auditoria de dependências
✅ Proteção contra módulos desaparecendo (left-pad problem)
```

**Configuração**:
```bash
export GOPROXY="https://athens.empresa.com,https://proxy.golang.org,direct"
export GOPRIVATE="github.com/empresa/*"
```

### 🔐 Vendoring

Copiar dependências para dentro do projeto:

```bash
# Criar diretório vendor/
go mod vendor

# Estrutura gerada:
vendor/
├── github.com/
│   └── gin-gonic/
│       └── gin/
└── modules.txt

# Build usando vendor (ignora cache)
go build -mod=vendor
```

**Quando usar**:
- Ambientes sem acesso à internet
- Garantir build reproduzível 100%
- Compliance (manter cópia do código)

**Desvantagens**:
- Aumenta tamanho do repositório
- Precisa atualizar vendor/ quando muda dependências

### 📦 Publicar Seus Próprios Módulos

Criar e publicar módulo para sua equipe:

```bash
# 1. Criar repositório no GitHub
# 2. Inicializar módulo
go mod init github.com/empresa/minha-biblioteca

# 3. Criar código
mkdir pkg/utils
# ... escrever código ...

# 4. Commit e push
git add .
git commit -m "Initial commit"
git push

# 5. Criar release
git tag v1.0.0
git push origin v1.0.0

# 6. Agora pode ser usado:
# go get github.com/empresa/minha-biblioteca
```

**Melhores práticas**:
- README.md completo com exemplos
- Versionamento semântico estrito
- Changelog (CHANGELOG.md)
- Testes automatizados (CI)
- Documentação no pkg.go.dev

### 📚 Recursos de Aprendizado

**Documentação Oficial**:
- [Go Modules Reference](https://go.dev/ref/mod)
- [Module release and versioning workflow](https://go.dev/doc/modules/release-workflow)
- [Developing and publishing modules](https://go.dev/doc/modules/developing)

**Artigos Recomendados**:
- [Using Go Modules](https://go.dev/blog/using-go-modules)
- [Migrating to Go Modules](https://go.dev/blog/migrating-to-go-modules)
- [Module Mirror and Checksum Database](https://go.dev/blog/module-mirror-launch)

**Ferramentas**:
- [go mod graph](https://pkg.go.dev/cmd/go#hdr-Print_module_requirement_graph) - Visualizar árvore de dependências
- [go mod why](https://pkg.go.dev/cmd/go#hdr-Explain_why_packages_or_modules_are_needed) - Por que um módulo é necessário
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) - Verificar vulnerabilidades

---

## 🎉 Conclusão

Você completou o guia completo sobre **Módulos Privados em Go**!

Agora você está preparado para:
- ✅ Trabalhar em empresas que usam módulos privados
- ✅ Criar e distribuir bibliotecas internas
- ✅ Configurar autenticação para repos privados
- ✅ Troubleshoot problemas comuns
- ✅ Seguir boas práticas de versionamento
- ✅ Gerenciar dependências em projetos reais

**💡 Lembre-se**: Módulos privados são essenciais para desenvolvimento corporativo. Pratique os conceitos criando suas próprias bibliotecas internas!

**Próxima aula**: Continue sua jornada explorando tópicos avançados de Go! 🚀

---

<div align="center">

**Feito com ❤️ para estudantes de Go**

[⬆️ Voltar ao topo](#-trabalhando-com-módulos-privados-em-go)

</div>
