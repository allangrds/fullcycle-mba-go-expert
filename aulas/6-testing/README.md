# Testing em Go

Este diretório contém exemplos práticos de conceitos fundamentais sobre testes em Go, abordando desde testes básicos até técnicas avançadas como mocking e fuzzing.

## 📚 Conceitos Abordados

### 1. Testes Básicos

**Localização:** `1-iniciando-testes/`

Fundamentos para criação de testes em Go:

- **Convenção de nomenclatura:** Arquivos de teste devem terminar com `_test.go`
- **Função de teste:** Deve começar com `Test` seguido do nome da função/funcionalidade sendo testada
- **Parâmetro:** Todas as funções de teste recebem `*testing.T` como parâmetro
- **Execução:** Use `go test` para executar os testes

```go
func TestCalculateTax(t *testing.T) {
    amount := 500.0
    expected := 5.0

    result := CalculateTax(amount)

    if result != expected {
        t.Errorf("Expected %f but got %f", expected, result)
    }
}
```

### 2. Testes em Batch (Table-Driven Tests)

**Localização:** `2-testando-em-batch/`

Técnica para testar múltiplos cenários de forma organizada:

- **Estrutura de dados:** Crie um slice de structs com casos de teste
- **Loop de teste:** Itere pelos casos testando cada um
- **Vantagens:** Facilita a adição de novos casos e melhora a legibilidade

```go
func TestCalculateTaxBatch(t *testing.T) {
    type calcTax struct {
        amount, expected float64
    }
    table := []calcTax{
        {250.0, 5.0},
        {500.0, 5.0},
        {1000.0, 10.0},
        {1500.0, 10.0},
    }

    for _, item := range table {
        result := CalculateTax(item.amount)
        if result != item.expected {
            t.Errorf("Expected %f but got %f", item.expected, result)
        }
    }
}
```

### 3. Cobertura de Testes

**Localização:** `3-verificando-cobertura-teste/`

Medição da cobertura de código pelos testes:

```bash
# Executar testes com cobertura
go test -cover

# Gerar arquivo de cobertura
go test -coverprofile=coverage.out

# Visualizar cobertura em HTML
go tool cover -html=coverage.out
```

**Benefícios:**

- Identifica código não testado
- Ajuda a garantir qualidade dos testes
- Métrica importante para CI/CD

### 4. Benchmarks

**Localização:** `4-trabalhando-com-benchmark/`

Testes de performance para medir desempenho:

- **Convenção:** Funções devem começar com `Benchmark`
- **Parâmetro:** Recebem `*testing.B`
- **Loop:** Use `b.N` para controlar iterações
- **Execução:** `go test -bench=.`

```go
func BenchmarkCalculateTax(b *testing.B) {
    for i := 0; i < b.N; i++ {
        CalculateTax(500.0)
    }
}
```

**Resultado típico:**

```text
BenchmarkCalculateTax-8    1000000000    0.25 ns/op
```

### 5. Fuzzing

**Localização:** `5-fuzzing/`

Teste automatizado com entradas aleatórias (Go 1.18+):

- **Convenção:** Funções devem começar com `Fuzz`
- **Parâmetro:** Recebem `*testing.F`
- **Seeds:** Defina valores iniciais com `f.Add()`
- **Execução:** `go test -fuzz=.`

```go
func FuzzCalculateTax(f *testing.F) {
    seed := []float64{-1, -2, -2.5, 500.0, 1000.0, 1501.0}
    for _, amount := range seed {
        f.Add(amount)
    }

    f.Fuzz(func(t *testing.T, amount float64) {
        result := CalculateTax(amount)

        if amount < 0 && result != 0 {
            t.Errorf("Received %f but expected 0", result)
        }
        if amount > 20000 && result != 20 {
            t.Errorf("Received %f but expected 20", result)
        }
    })
}
```

**Benefícios:**

- Encontra bugs edge cases
- Testa comportamento com entradas inesperadas
- Execução automática de milhares de casos

### 6. Testify - Biblioteca de Assertions

**Localização:** `6-iniciando-com-testify/`

Biblioteca popular que simplifica a escrita de testes:

```bash
go get github.com/stretchr/testify
```

**Principais funcionalidades:**

- **assert:** Assertions mais legíveis
- **require:** Para condições que devem parar o teste
- **mock:** Sistema de mocking robusto

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCalculateTax(t *testing.T) {
    tax, err := CalculateTax(1000)
    assert.Nil(t, err)
    assert.Equal(t, 10.0, tax)

    tax, err = CalculateTax(0)
    assert.Error(t, err, "amount must be greater than zero")
    assert.Equal(t, 0.0, tax)
}
```

### 7. Mocking

**Localização:** `7-trabalhando-com-mocks/`

Simulação de dependências para testes isolados:

**Criando um Mock:**

```go
type TaxRepositoryMock struct {
    mock.Mock
}

func (m *TaxRepositoryMock) SaveTax(tax float64) error {
    args := m.Called(tax)
    return args.Error(0)
}
```

**Usando o Mock em testes:**

```go
func TestCalculateTaxAndSave(t *testing.T) {
    repository := &TaxRepositoryMock{}

    // Configurar expectativas
    repository.On("SaveTax", 10.0).Return(nil).Once()
    repository.On("SaveTax", 0.0).Return(errors.New("error saving tax"))

    // Executar teste
    err := CalculateTaxAndSave(1000.0, repository)
    assert.Nil(t, err)

    // Verificar expectativas
    repository.AssertExpectations(t)
    repository.AssertNumberOfCalls(t, "SaveTax", 2)
}
```

## 🛠️ Comandos Úteis

```bash
# Executar todos os testes
go test

# Executar testes com verbose
go test -v

# Executar teste específico
go test -run TestCalculateTax

# Executar testes com cobertura
go test -cover

# Executar benchmarks
go test -bench=.

# Executar fuzzing
go test -fuzz=FuzzCalculateTax

# Executar testes em modo de observação
go test -watch
```

## 📊 Boas Práticas

### Estrutura de Testes

- **AAA Pattern:** Arrange, Act, Assert
- **Nomenclatura clara:** Nomes descritivos para testes
- **Um conceito por teste:** Testes focados e específicos

### Organização

- **Table-driven tests:** Para múltiplos cenários
- **Subtests:** Use `t.Run()` para agrupar casos relacionados
- **Setup/Teardown:** Prepare e limpe recursos adequadamente

### Qualidade

- **Cobertura:** Mantenha cobertura alta, mas foque na qualidade
- **Mocks:** Use para isolar unidades de teste
- **Benchmarks:** Monitore performance de código crítico

### Performance

- **Fuzzing:** Encontre edge cases automaticamente
- **Parallel tests:** Use `t.Parallel()` quando apropriado
- **Cleanup:** Use `t.Cleanup()` para limpeza automática

## 🎯 Objetivos de Aprendizado

Após estudar estes exemplos, você deve ser capaz de:

1. ✅ Criar testes básicos usando o pacote `testing`
2. ✅ Implementar table-driven tests para múltiplos cenários
3. ✅ Medir e interpretar cobertura de código
4. ✅ Escrever benchmarks para medir performance
5. ✅ Utilizar fuzzing para encontrar bugs automaticamente
6. ✅ Usar a biblioteca Testify para assertions mais legíveis
7. ✅ Implementar mocks para testar dependências isoladamente
8. ✅ Aplicar boas práticas de testing em projetos Go

## 📖 Recursos Adicionais

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Fuzzing Tutorial](https://go.dev/doc/tutorial/fuzz)
- [Effective Go - Testing](https://golang.org/doc/effective_go#testing)
