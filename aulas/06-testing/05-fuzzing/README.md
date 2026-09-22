# Fuzz Testing em Go

## O que é Fuzz Testing?

Fuzz Testing (ou Fuzzing) é uma técnica de teste automatizada que envia dados aleatórios, inesperados ou malformados como entrada para um programa, com o objetivo de encontrar bugs, vulnerabilidades de segurança ou comportamentos inesperados.

No Go 1.18+, o Fuzzing foi incorporado nativamente à toolchain padrão, permitindo que desenvolvedores encontrem bugs que poderiam passar despercebidos em testes tradicionais.

## Como Funciona o Fuzz Testing

O Fuzz Testing no Go funciona da seguinte forma:

1. **Valores Seed (Sementes)**: Você fornece valores iniciais conhecidos que são bons casos de teste
2. **Geração Automática**: O Go gera automaticamente novos valores baseados nas sementes
3. **Execução em Massa**: A função é executada milhares/milhões de vezes com diferentes inputs
4. **Detecção de Problemas**: Quando um input causa um erro, ele é salvo para análise

## Exemplo Prático: CalculateTax

### Código Principal (`tax.go`)

```go
package tax

func CalculateTax(amount float64) float64 {
    if amount < 0 {
        return 0
    }

    if amount >= 1000 {
        return 10.0
    }

    if amount >= 20000 {  // BUG: Esta condição nunca será verdadeira!
        return 20.0
    }

    return 5.0
}
```

**Problema no código**: A lógica está incorreta. A condição `amount >= 20000` nunca será executada porque `amount >= 1000` é verificada primeiro e retorna 10.0 para qualquer valor >= 1000.

### Fuzz Test (`tax_test.go`)

```go
func FuzzCalculateTax(f *testing.F) {
    // 1. Definindo valores seed (sementes)
    seed := []float64{-1, -2, -2.5, 500.0, 1000.0, 1501.0}
    for _, amount := range seed {
        f.Add(amount) // Adiciona cada seed ao corpus de teste
    }

    // 2. Definindo a função de fuzz
    f.Fuzz(func(t *testing.T, amount float64) {
        result := CalculateTax(amount)

        // 3. Definindo invariantes (regras que sempre devem ser verdadeiras)
        if amount < 0 && result != 0 {
            t.Errorf("Received %f but expected 0", result)
        }
        if amount > 20000 && result != 20 {
            t.Errorf("Received %f but expected 20", result) // Este teste vai falhar!
        }
    })
}
```

## Pasta testdata e Corpus de Fuzz

### O que é a pasta testdata?

Quando você executa um fuzz test, o Go automaticamente cria uma pasta `testdata/fuzz/` no seu projeto. Esta pasta armazena:

1. **Corpus**: Conjunto de inputs que foram testados
2. **Casos que causaram falha**: Inputs específicos que fizeram o teste falhar
3. **Metadados**: Informações sobre a execução do fuzz

### Estrutura gerada:

```
testdata/
└── fuzz/
    └── FuzzCalculateTax/
        └── 19faafe021a066ef  # Hash único do input que causou falha
```

### Conteúdo do arquivo de falha:

```
go test fuzz v1
float64(165240)
```

Este arquivo mostra que o input `165240` causou uma falha no teste.

## Passo a Passo para Executar Fuzz Tests

### 1. Executar o Fuzz Test

```bash
# Executar fuzz por tempo determinado (ex: 30 segundos)
go test -fuzz=FuzzCalculateTax -fuzztime=30s

# Executar até encontrar uma falha
go test -fuzz=FuzzCalculateTax

# Executar sem fuzz (apenas com seeds)
go test -run=FuzzCalculateTax
```

### 2. Analisar Resultados

Quando uma falha é encontrada:

```bash
--- FAIL: FuzzCalculateTax (0.01s)
    --- FAIL: FuzzCalculateTax (0.00s)
        tax_test.go:58: Received 10.000000 but expected 20

    Failing input written to testdata/fuzz/FuzzCalculateTax/19faafe021a066ef
```

### 3. Reproduzir a Falha

```bash
# Executa apenas o caso que falhou
go test -run=FuzzCalculateTax/19faafe021a066ef
```

### 4. Corrigir o Bug

Após identificar o problema, corrija o código:

```go
func CalculateTax(amount float64) float64 {
    if amount < 0 {
        return 0
    }

    if amount >= 20000 {  // Mover para antes da condição >= 1000
        return 20.0
    }

    if amount >= 1000 {
        return 10.0
    }

    return 5.0
}
```

### 5. Validar a Correção

```bash
# Executar o teste que falhou para confirmar que agora passa
go test -run=FuzzCalculateTax/19faafe021a066ef

# Executar novamente o fuzz para encontrar outros problemas
go test -fuzz=FuzzCalculateTax -fuzztime=30s
```

## Vantagens do Fuzz Testing

1. **Encontra bugs edge-case**: Descobre problemas com inputs inesperados
2. **Automatizado**: Não precisa pensar em todos os casos manualmente
3. **Reproduzível**: Salva automaticamente os inputs que causaram falha
4. **Integrado**: Parte nativa do Go, sem dependências externas

## Boas Práticas

1. **Use seeds representativos**: Inclua valores conhecidos importantes
2. **Defina invariantes claros**: Regras que sempre devem ser verdadeiras
3. **Execute regularmente**: Inclua no CI/CD para execução contínua
4. **Limite o tempo**: Use `-fuzztime` para evitar execuções muito longas
5. **Mantenha o corpus**: Não delete a pasta `testdata`, ela é valiosa para regressions

## Comandos Úteis

```bash
# Executar fuzz por 1 minuto
go test -fuzz=. -fuzztime=1m

# Executar apenas testes normais (sem fuzz)
go test -run=^Test

# Limpar corpus de fuzz
rm -rf testdata/

# Ver apenas falhas de fuzz
go test -fuzz=. -run=^#
```
