# Packing em Go - Gerenciamento de Pacotes e Módulos

Este diretório contém exemplos práticos sobre gerenciamento de pacotes, módulos e workspaces em Go, abordando desde conceitos básicos de exportação até técnicas avançadas de desenvolvimento local.

## 📚 Conceitos Abordados

### 2. Acessando Pacotes Criados

**Localização:** `2-acessando-pacotes-criados/`

Fundamentos para criação e uso de pacotes personalizados:

**Estrutura do projeto:**

```text
2-acessando-pacotes-criados/
├── go.mod
├── cmd/
│   └── main.go
└── math/
    └── math.go
```

**Criando um pacote (`math/math.go`):**

```go
package math

type Math struct {
    A int
    B int
}

func (m Math) Add() int {
    return m.A + m.B
}
```

**Usando o pacote (`cmd/main.go`):**

```go
package main

import (
    "fmt"
    "github.com/allangrds/fullcycle-mba-go-expert/aulas/5-packing/2-acessando-pacotes-criados/math"
)

func main() {
    funcMath := math.Math{A: 1, B: 2}
    fmt.Println(funcMath.Add())
}
```

**Conceitos principais:**

- **Estrutura de módulo:** Use `go mod init` para criar um módulo
- **Organização:** Separe código em pacotes lógicos
- **Importação:** Use o caminho completo do módulo para importar pacotes
- **Convenção:** Use nomes descritivos para pacotes

### 3. Exportação de Objetos

**Localização:** `3-exportacao-de-objetos/`

Controle de visibilidade e exportação em Go:

**Regras de exportação:**

```go
package math

// Exportados (letra maiúscula no início)
var X string = "hello X"           // ✅ Exportado
type Math struct {                 // ✅ Exportado
    A int                         // ✅ Campo exportado
    B int                         // ✅ Campo exportado
}
func (m Math) Add() int { ... }   // ✅ Método exportado

// Não exportados (letra minúscula no início)
var x string = "hello x"          // ❌ Não exportado
type mathB struct {               // ❌ Não exportado
    a int                         // ❌ Campo não exportado
    b int                         // ❌ Campo não exportado
}
```

**Padrão Constructor:**

```go
// Struct não exportada com construtor público
type mathB struct {
    a int
    b int
}

// Função construtora exportada
func NewMathB(a, b int) mathB {
    return mathB{a: a, b: b}
}

func (m mathB) AddB() int {
    return m.a + m.b
}
```

**Uso do padrão Constructor:**

```go
// Criando instância através do construtor
funcMathB := math.NewMathB(1, 2)
fmt.Println(funcMathB.AddB())
```

**Benefícios do padrão Constructor:**

- **Encapsulamento:** Mantém estrutura interna privada
- **Validação:** Permite validar dados na criação
- **Flexibilidade:** Facilita mudanças futuras na estrutura
- **Controle:** Define exatamente como objetos são criados

### 5. Go Mod Replace

**Localização:** `5-go-mod-replace/`

Desenvolvimento de pacotes locais não publicados:

**Estrutura do projeto:**

```text
5-go-mod-replace/
├── math/
│   ├── go.mod
│   └── math.go
└── system/
    ├── go.mod
    └── main.go
```

**Pacote math (`math/go.mod`):**

```go.mod
module github.com/allangrds/fullcycle-mba-go-expert/aulas/5-packing/5-go-mod-replace/math

go 1.25.0
```

**Sistema consumidor (`system/go.mod`):**

```go.mod
module github.com/allangrds/fullcycle-mba-go-expert/aulas/5-packing/5-go-mod-replace/system

go 1.25.0

replace github.com/allangrds/fullcycle-mba-go-expert/aulas/5-packing/5-go-mod-replace/math => ../math

require github.com/allangrds/fullcycle-mba-go-expert/aulas/5-packing/5-go-mod-replace/math v0.0.0-00010101000000-000000000000
```

**Comandos para configurar replace:**

```bash
# Configurar o replace directive
go mod edit -replace github.com/allangrds/fullcycle-mba-go-expert/aulas/5-packing/5-go-mod-replace/math=../math

# Atualizar dependências
go mod tidy
```

**Casos de uso:**

- **Desenvolvimento local:** Testar pacotes antes da publicação
- **Debugging:** Modificar dependências temporariamente
- **Fork privado:** Usar versão customizada de uma biblioteca
- **Monorepo:** Gerenciar múltiplos módulos relacionados

### 6. Usando Workspaces

**Localização:** `6-usando-workspaces/`

Gerenciamento de múltiplos módulos em um workspace (Go 1.18+):

**Arquivo go.work:**

```go.work
go 1.25.0

use (
    ./math
    ./system
)
```

**Comandos para workspaces:**

```bash
# Inicializar workspace
go work init ./math ./system

# Executar código no workspace
go run system/main.go

# Adicionar módulo ao workspace
go work use ./novo-modulo

# Sincronizar workspace
go work sync
```

**Vantagens dos workspaces:**

- **Simplicidade:** Não precisa de replace directives
- **Automático:** Go resolve dependências automaticamente
- **Flexibilidade:** Fácil adicionar/remover módulos
- **Desenvolvimento:** Ideal para projetos multi-módulo

**Considerações importantes:**

```bash
# Para projetos com dependências externas
go mod tidy -e  # Ignora pacotes não encontrados
```

**Problemas comuns:**

- **Dependências externas:** Podem não ser baixadas corretamente
- **CI/CD:** Workspaces são para desenvolvimento local
- **Publicação:** Módulos devem funcionar independentemente

## 🛠️ Comandos Úteis

### Gerenciamento de Módulos

```bash
# Inicializar módulo
go mod init nome-do-modulo

# Adicionar dependência
go get github.com/user/repo

# Atualizar dependências
go mod tidy

# Verificar dependências
go mod why pacote

# Baixar dependências
go mod download
```

### Replace Directives

```bash
# Adicionar replace
go mod edit -replace old=new

# Remover replace
go mod edit -dropreplace old

# Replace para caminho local
go mod edit -replace github.com/user/repo=../local-path
```

### Workspaces

```bash
# Criar workspace
go work init

# Adicionar módulos
go work use ./modulo1 ./modulo2

# Sincronizar workspace
go work sync

# Editar workspace
go work edit -use ./novo-modulo
```

### Publicação de Pacotes

```bash
# Criar tag de versão
git tag v1.0.0
git push origin v1.0.0

# Publicar no proxy do Go
GOPROXY=proxy.golang.org go list -m github.com/user/repo@v1.0.0
```

## 📊 Boas Práticas

### Estrutura de Projetos

- **Separação clara:** Um pacote por responsabilidade
- **Nomenclatura:** Use nomes descritivos e consistentes
- **Hierarquia:** Organize pacotes de forma lógica
- **Documentação:** Documente APIs públicas

### Exportação e Encapsulamento

- **Princípio da menor exposição:** Exporte apenas o necessário
- **Constructors:** Use para structs complexas
- **Interfaces:** Prefira interfaces públicas sobre tipos concretos
- **Validação:** Valide dados nos constructors

### Desenvolvimento Local

- **Workspaces:** Use para desenvolvimento multi-módulo
- **Replace temporário:** Para debugging e testes
- **Versionamento:** Use tags semânticas para releases
- **CI/CD:** Teste sem workspaces/replaces

### Organização de Código

```text
projeto/
├── go.mod
├── go.work              # Para desenvolvimento local
├── cmd/                 # Aplicações executáveis
│   └── main.go
├── internal/            # Código privado do projeto
│   └── service/
├── pkg/                 # Código reutilizável
│   └── utils/
└── api/                 # Definições de API
    └── proto/
```

## 🎯 Objetivos de Aprendizado

Após estudar estes exemplos, você deve ser capaz de:

1. ✅ Criar e organizar pacotes Go
2. ✅ Controlar exportação de tipos, funções e variáveis
3. ✅ Implementar padrões de constructor
4. ✅ Usar replace directives para desenvolvimento local
5. ✅ Configurar e usar Go workspaces
6. ✅ Gerenciar dependências com go mod
7. ✅ Aplicar boas práticas de estruturação de projetos
8. ✅ Publicar e versionar pacotes Go

## 🔄 Fluxo de Desenvolvimento

### Para Pacotes Simples

1. `go mod init nome-do-modulo`
2. Criar estrutura de pacotes
3. Implementar e testar
4. Publicar com tags git

### Para Projetos Multi-Módulo

1. `go work init`
2. `go work use ./modulo1 ./modulo2`
3. Desenvolver e testar localmente
4. Publicar módulos independentemente

### Para Desenvolvimento com Dependências Locais

1. Use workspaces (recomendado)
2. Ou use replace directives
3. Teste sem workspaces antes de publicar
4. Configure CI/CD sem dependências locais

## 📖 Recursos Adicionais

- [Go Modules Reference](https://golang.org/ref/mod)
- [Go Workspaces Tutorial](https://go.dev/doc/tutorial/workspaces)
- [Package Discovery](https://pkg.go.dev/)
- [Module Publishing](https://golang.org/doc/modules/publishing)
- [Semantic Versioning](https://semver.org/)
