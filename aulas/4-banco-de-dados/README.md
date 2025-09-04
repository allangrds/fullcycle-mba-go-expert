# Banco de Dados em Go

Este diretório contém exemplos práticos sobre integração com banco de dados em Go, abordando desde SQL puro até ORM avançado com GORM, incluindo relacionamentos e otimizações.

## 🔤 Entendendo a Sintaxe Go para Banco de Dados

### Blank Identifier (`_`)

O `_` (underscore) é usado para ignorar valores que não queremos utilizar:

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql" // Importa apenas para side effects
)

// Ignorando valor de retorno que não precisamos
_, err = stmt.Exec(product.ID, product.Name, product.Price)
```

**Por que usar:**
- **Import side effects:** Driver MySQL registra-se automaticamente
- **Ignorar retornos:** Quando não precisamos do valor, apenas do erro
- **Evitar variáveis não utilizadas:** Go não compila com variáveis não usadas

### Ponteiros (`*` e `&`)

Ponteiros são fundamentais para eficiência e modificação de dados:

```go
// * indica que é um ponteiro para Product
func NewProduct(name string, price float64) *Product {
    return &Product{  // & retorna o endereço de memória
        ID:    uuid.New().String(),
        Name:  name,
        Price: price,
    }
}

// & passa o endereço para Scan modificar a variável original
err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
```

**Quando usar ponteiros:**
- **Eficiência:** Evita copiar estruturas grandes
- **Modificação:** Permite modificar a variável original
- **Nil checking:** Pode representar "ausência de valor"

### Defer Statement

`defer` executa código no final da função, mesmo se houver erro:

```go
func insertProduct(db *sql.DB, product *Product) error {
    stmt, err := db.Prepare("insert into products(id, name, price) values (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close() // ✅ Sempre executa, mesmo com erro

    // resto do código...
}
```

**Características do defer:**
- **LIFO:** Último defer é executado primeiro
- **Garantia:** Executa mesmo com panic ou return
- **Cleanup:** Ideal para liberar recursos

### Error Handling Pattern

Go usa retorno explícito de erros:

```go
// Padrão: última variável de retorno é sempre error
func selectProduct(db *sql.DB, id string) (*Product, error) {
    stmt, err := db.Prepare("select id, name, price from products where id = ?")
    if err != nil {
        return nil, err // ✅ Retorna erro imediatamente
    }
    defer stmt.Close()

    var product Product
    err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
    if err != nil {
        return nil, err // ✅ Retorna erro se scan falhar
    }

    return &product, nil // ✅ Sucesso: produto + erro nil
}
```

### Tags de Struct

Tags fornecem metadados para bibliotecas:

```go
type Product struct {
    ID    int    `gorm:"primaryKey"`                    // GORM: chave primária
    Name  string `gorm:"size:255;not null"`           // GORM: tamanho e not null
    Price float64 `gorm:"type:decimal(10,2)"`         // GORM: tipo específico do DB
    JSON  string `json:"json_name" db:"database_name"` // Tags múltiplas
}
```

**Tipos de tags comuns:**
- **gorm:** Configurações do ORM GORM
- **json:** Serialização JSON
- **db:** Mapeamento de colunas para database/sql
- **validate:** Validação de dados

### Slices vs Arrays

```go
// Array: tamanho fixo
var products [5]Product

// Slice: tamanho dinâmico (mais comum)
var products []Product
products = append(products, product) // Adiciona elemento

// Slice literal
products := []Product{
    {Name: "Notebook", Price: 1000.00},
    {Name: "Mouse", Price: 50.00},
}
```

### Interface{} e Type Assertion

```go
// interface{} aceita qualquer tipo
var value interface{} = "hello"

// Type assertion: converter para tipo específico
if str, ok := value.(string); ok {
    fmt.Println("É uma string:", str)
}

// Usado em GORM para updates dinâmicos
db.Model(&product).Updates(map[string]interface{}{
    "name":  "Novo Nome",
    "price": 999.99,
})
```

## 🏗️ Configuração do Ambiente

### Docker Compose Setup

**Arquivo:** `docker-compose.yml`

```yaml
services:
  mysql:
    image: mysql:8.0
    container_name: mysql
    restart: always
    platform: linux/amd64
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: goexpert
      MYSQL_PASSWORD: root
    ports:
      - 3306:3306
```

### Makefile para Automação

**Arquivo:** `Makefile`

```makefile
.PHONY: up mysql

up:
    docker compose up

mysql:
    docker compose exec mysql bash
```

**Comandos úteis:**

```bash
# Subir o ambiente
make up

# Acessar container MySQL
make mysql
```

## 📚 Conceitos Abordados

### 2. Preparando Base do Código

**Localização:** `2-preparando-base-codigo/`

Estrutura básica para trabalhar com bancos de dados:

**Modelo de dados:**

```go
package main

import "github.com/google/uuid"

type Product struct {
    ID    string
    Name  string
    Price float64
}

func NewProduct(name string, price float64) *Product {
    return &Product{
        ID:    uuid.New().String(),
        Name:  name,
        Price: price,
    }
}
```

**Dependências necessárias:**

```go.mod
module example

go 1.24.5

require github.com/google/uuid v1.6.0
```

### 3. Inserindo Dados no Banco

**Localização:** `3-inserindo-dados-no-banco/`

Operações básicas de inserção com SQL puro:

**Configuração da conexão:**

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql" // Driver MySQL
)

func main() {
    // sql.Open não conecta imediatamente, apenas valida argumentos
    db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
    if err != nil {
        panic(err) // Para em caso de erro
    }
    defer db.Close() // ✅ Fecha conexão no final da função
}
```

**Anatomia da Connection String:**
```
"root:root@tcp(localhost:3306)/goexpert"
 └─user:password@protocol(host:port)/database
```

**Função de inserção detalhada:**

```go
func insertProduct(db *sql.DB, product *Product) error {
    // Prepare: cria statement reutilizável e seguro
    // ? são placeholders para evitar SQL injection
    stmt, err := db.Prepare("insert into products(id, name, price) values (?, ?, ?)")
    if err != nil {
        return err // ❌ Falha na preparação
    }
    defer stmt.Close() // ✅ Libera recursos do statement

    // Exec: executa statement com valores reais
    // _ ignora sql.Result (lastInsertId, rowsAffected)
    _, err = stmt.Exec(product.ID, product.Name, product.Price)
    if err != nil {
        return err // ❌ Falha na execução
    }

    return nil // ✅ Sucesso
}
```

**Explicação das operações:**

1. **db.Prepare():**
   - Compila SQL uma vez, executa múltiplas vezes
   - Previne SQL injection automaticamente
   - Retorna `*sql.Stmt` e `error`

2. **stmt.Exec():**
   - Executa statement com parâmetros
   - Retorna `sql.Result` e `error`
   - Result contém LastInsertId() e RowsAffected()

3. **defer stmt.Close():**
   - Libera recursos do prepared statement
   - Executa automaticamente no final da função
   - Previne memory leaks

**Conceitos importantes:**

- **Prepared Statements:** Previnem SQL injection e melhoram performance
- **Defer Close:** Liberam recursos automaticamente, mesmo com erros
- **Error Handling:** Go força verificação explícita de erros
- **Blank Identifier:** `_` ignora valores não utilizados (sql.Result)

### 4. Alterando Dados no Banco

**Localização:** `4-alterando-dados-no-banco/`

Operações de atualização com SQL:

```go
func updateProduct(db *sql.DB, product *Product) error {
    stmt, err := db.Prepare("update products set name = ?, price = ? where id = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(product.Name, product.Price, product.ID)
    if err != nil {
        return err
    }

    return nil
}
```

### 5. Trabalhando com Query Row

**Localização:** `5-trabalhand-query-row/`

Consultas que retornam uma única linha:

```go
func selectProduct(db *sql.DB, id string) (*Product, error) {
    // Prepare: statement para buscar por ID
    stmt, err := db.Prepare("select id, name, price from products where id = ?")
    if err != nil {
        return nil, err // ❌ Erro na preparação
    }
    defer stmt.Close() // ✅ Cleanup automático

    var product Product // Variável para receber dados

    // QueryRow: busca APENAS 1 linha (primeira encontrada)
    // Scan: mapeia colunas SQL para campos Go
    // & passa endereço para Scan modificar as variáveis
    err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
    if err != nil {
        return nil, err // ❌ sql.ErrNoRows se não encontrar
    }

    return &product, nil // ✅ Retorna ponteiro para produto
}
```

**Explicação detalhada do Scan:**

```go
// A ordem DEVE corresponder à ordem das colunas no SELECT
// SELECT id, name, price FROM...
//        |    |     |
//        v    v     v
err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
```

**Por que usar ponteiros (&) no Scan:**
- Scan precisa **modificar** as variáveis originais
- Sem &, Scan receberia **cópias** e não conseguiria alterar
- & passa o **endereço de memória** onde estão as variáveis

**Diferenças entre métodos de query:**

| Método | Uso | Retorno |
|--------|-----|---------|
| `QueryRow()` | Uma linha apenas | `*sql.Row` |
| `Query()` | Múltiplas linhas | `*sql.Rows` |
| `Exec()` | INSERT/UPDATE/DELETE | `sql.Result` |

**Tratamento de erros específicos:**

```go
err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
if err != nil {
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("produto com ID %s não encontrado", id)
    }
    return nil, fmt.Errorf("erro ao buscar produto: %w", err)
}
```

**Métodos alternativos:**

```go
// Com contexto para timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err = stmt.QueryRowContext(ctx, id).Scan(&product.ID, &product.Name, &product.Price)
```

**Vantagens do QueryRow:**
- **Simples:** Para consultas que retornam um registro
- **Automático:** Não precisa iterar como Query()
- **Seguro:** sql.ErrNoRows quando não encontra
- **Eficiente:** Para para na primeira linha encontrada### 6. Selecionando Múltiplos Registros

**Localização:** `6-selecionando-multiplos-registros/`

Consultas que retornam múltiplas linhas:

```go
func selectProducts(db *sql.DB) ([]Product, error) {
    rows, err := db.Query("select id, name, price from products")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var product Product
        err = rows.Scan(&product.ID, &product.Name, &product.Price)
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }

    return products, nil
}
```

### 7. Removendo Registro

**Localização:** `7-removendo-registro/`

Operações de exclusão:

```go
func deleteProduct(db *sql.DB, id string) error {
    stmt, err := db.Prepare("delete from products where id = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    if err != nil {
        return err
    }

    return nil
}
```

### 8. Configurando e Criando Operações com GORM

**Localização:** `8-configurando-e-criando-operacoes/`

Introdução ao ORM GORM com explicações detalhadas:

**Setup inicial:**

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// Struct com tags GORM
type Product struct {
    ID    int `gorm:"primaryKey"`        // ✅ Chave primária
    Name  string                        // ✅ Campo VARCHAR padrão
    Price float64                       // ✅ Campo DECIMAL/FLOAT
}

func main() {
    // DSN: Data Source Name - string de conexão
    dsn := "root:root@tcp(localhost:3306)/goexpert"

    // Abre conexão usando driver MySQL
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // AutoMigrate: cria/atualiza tabelas automaticamente
    // Analisa struct e cria SQL correspondente
    db.AutoMigrate(&Product{})
}
```

**Entendendo as Tags GORM:**

```go
type Product struct {
    // primaryKey: define como chave primária
    ID int `gorm:"primaryKey"`

    // Configurações de campo
    Name string `gorm:"size:255;not null;uniqueIndex"`
    //              |     |        |      └── índice único
    //              |     |        └────────── não pode ser NULL
    //              |     └─────────────────── tamanho varchar
    //              └───────────────────────── tag GORM

    // Tipo específico do banco
    Price float64 `gorm:"type:decimal(10,2);not null"`
    //                   └── 10 dígitos, 2 decimais

    // Campo opcional (pode ser NULL)
    Description *string `gorm:"size:1000"`
    //          └── ponteiro permite NULL

    // Timestamps automáticos (se usar gorm.Model)
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

**AutoMigrate explicado:**

```go
// O que AutoMigrate faz:
db.AutoMigrate(&Product{})

// Equivale a algo como:
/*
CREATE TABLE products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
*/
```

**Operações básicas explicadas:**

```go
// CREATE - Criar um produto
product := Product{
    Name:  "Notebook Dell",
    Price: 1000.00,
}
result := db.Create(&product) // & passa ponteiro para GORM preencher ID

// Verificando resultado
if result.Error != nil {
    log.Fatal(result.Error)
}
fmt.Printf("Produto criado com ID: %d\n", product.ID) // GORM preenche o ID

// CREATE BATCH - Criar múltiplos produtos
products := []Product{
    {Name: "Notebook Dell", Price: 1000.00},
    {Name: "Notebook Samsung", Price: 1200.00},
}
db.Create(&products) // ✅ Insere todos de uma vez (mais eficiente)
```

**Vantagens do GORM sobre SQL puro:**

- **Auto Migration:** Cria/atualiza schema automaticamente
- **Type Safety:** Erros em tempo de compilação
- **Associations:** Relacionamentos automáticos
- **Hooks:** Before/After Create, Update, Delete
- **Soft Delete:** Exclusão lógica built-in
- **Query Builder:** Construção dinâmica de queries

**Configurações avançadas:**

```go
// Configuração personalizada
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // Log de queries
    DryRun: true,                                // Só mostra SQL, não executa
    DisableForeignKeyConstraintWhenMigrating: true, // Sem FK constraints
})
```

### 9. Realizando Primeiras Consultas

**Localização:** `9-realizando-primeiras-consultas/`

Consultas básicas com GORM:

```go
// Buscar por ID
var product Product
db.First(&product, 1)

// Buscar todos
var products []Product
db.Find(&products)

// Buscar um produto
db.First(&product, "name = ?", "Notebook")
```

### 10. Realizando Consultas com Where

**Localização:** `10-realizando-consultas-com-where/`

Consultas condicionais:

```go
// Where simples
var products []Product
db.Where("price > ?", 100).Find(&products)

// Múltiplas condições
db.Where("price > ? AND name LIKE ?", 100, "%Notebook%").Find(&products)

// Usando structs como condição
db.Where(&Product{Name: "Notebook"}).Find(&products)
```

### 11. Alterando e Removendo Registros

**Localização:** `11-alterando-e-removendo-registros/`

Operações de update e delete com GORM:

```go
// Atualizar um campo
db.Model(&product).Update("price", 1500.00)

// Atualizar múltiplos campos
db.Model(&product).Updates(Product{Name: "Notebook Pro", Price: 2000.00})

// Deletar
db.Delete(&product)

// Delete condicional
db.Where("price < ?", 100).Delete(&Product{})
```

### 12. Trabalhando com Soft Delete

**Localização:** `12-trabalhando-com-soft-delete/`

Exclusão lógica (soft delete):

```go
type Product struct {
    ID    int `gorm:"primaryKey"`
    Name  string
    Price float64
    gorm.DeletedAt `gorm:"index"` // Habilita soft delete
}

// Delete (soft) - apenas marca como deletado
db.Delete(&product)

// Para incluir registros deletados
db.Unscoped().Find(&products)

// Delete definitivo
db.Unscoped().Delete(&product)
```

### 13. Belongs To (Relacionamento N:1)

**Localização:** `13-belongs-to/`

Relacionamento onde muitos produtos pertencem a uma categoria:

```go
type Category struct {
    ID   int `gorm:"primaryKey"`
    Name string
    gorm.Model // Adiciona ID, CreatedAt, UpdatedAt, DeletedAt
}

type Product struct {
    ID         int `gorm:"primaryKey"`
    Name       string
    Price      float64

    // Foreign Key: referência à categoria
    CategoryId int      `gorm:"not null"`        // ✅ Campo FK obrigatório
    Category   Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    //         └── Struct aninhada para relacionamento

    gorm.Model
}
```

**Explicação do relacionamento:**

```sql
-- SQL gerado pelo GORM:
CREATE TABLE categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE TABLE products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    price DECIMAL,
    category_id BIGINT,  -- ✅ Foreign Key
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
```

**Criando dados com relacionamento:**

```go
// 1. Criar categoria primeiro
category := Category{Name: "Eletrônicos"}
db.Create(&category) // GORM preenche category.ID automaticamente

// 2. Criar produto associado
product := Product{
    Name:       "Notebook",
    Price:      1000.00,
    CategoryId: category.ID, // ✅ Referência à categoria criada
}
db.Create(&product)

// OU criar com associação automática:
product := Product{
    Name:     "Notebook",
    Price:    1000.00,
    Category: category, // ✅ GORM resolve automaticamente
}
db.Create(&product)
```

**Consultando com Preload:**

```go
var products []Product

// SEM Preload: CategoryId vem preenchido, mas Category fica vazio
db.Find(&products)
for _, product := range products {
    fmt.Println(product.Name)
    fmt.Println(product.CategoryId) // ✅ 1, 2, 3...
    fmt.Println(product.Category.Name) // ❌ "" (vazio)
}

// COM Preload: GORM carrega automaticamente os dados da categoria
db.Preload("Category").Find(&products)
for _, product := range products {
    fmt.Println(product.Name)         // ✅ "Notebook"
    fmt.Println(product.Category.Name) // ✅ "Eletrônicos"
}
```

**Preload explicado:**
- **Sem Preload:** Apenas CategoryId é carregado (eficiente mas limitado)
- **Com Preload:** GORM faz JOIN ou query separada para carregar Category
- **Lazy Loading:** Dados relacionados só carregam quando solicitados
- **N+1 Problem:** Preload previne múltiplas queries desnecessárias

**Consultas avançadas:**

```go
// Preload com condições
db.Preload("Category", "name = ?", "Eletrônicos").Find(&products)

// Preload aninhado (se Category tivesse subcategorias)
db.Preload("Category.Parent").Find(&products)

// Joins para filtros
db.Joins("Category").Where("categories.name = ?", "Eletrônicos").Find(&products)
```

### 14. Has One (Relacionamento 1:1)

**Localização:** `14-has-one/`

Relacionamento onde um produto tem um número de série:

```go
type Product struct {
    ID           int `gorm:"primaryKey"`
    Name         string
    SerialNumber SerialNumber
    gorm.Model
}

type SerialNumber struct {
    ID        int `gorm:"primaryKey"`
    Number    string
    ProductID int
}
```

### 15. Has Many (Relacionamento 1:N)

**Localização:** `15-has-many/`

Relacionamento onde uma categoria tem muitos produtos:

```go
type Category struct {
    ID       int `gorm:"primaryKey"`
    Name     string
    Products []Product // Has many
    gorm.Model
}

type Product struct {
    ID         int `gorm:"primaryKey"`
    Name       string
    Price      float64
    CategoryId int
    Category   Category
    gorm.Model
}

// Buscar categoria com produtos
var categories []Category
db.Preload("Products").Find(&categories)

for _, category := range categories {
    fmt.Println(category.Name, ":")
    for _, product := range category.Products {
        println("- ", product.Name)
    }
}
```

### 16. Pegadinhas Has Many

**Localização:** `16-pegadinhas-has-many/`

Problemas comuns e soluções no relacionamento has many:

- **N+1 Problem:** Use Preload adequadamente
- **Performance:** Considere lazy loading
- **Memory:** Cuidado com grandes volumes de dados

### 17. Many to Many (Relacionamento N:N)

**Localização:** `17-many-to-many/`

Relacionamento onde produtos podem ter múltiplas categorias e vice-versa:

```go
type Category struct {
    ID       int `gorm:"primaryKey"`
    Name     string
    // Slice indica relacionamento "muitos"
    Products []Product `gorm:"many2many:products_categories;"`
    //                        └── nome da tabela intermediária
    gorm.Model
}

type Product struct {
    ID         int `gorm:"primaryKey"`
    Name       string
    Price      float64
    // Relacionamento bidirecional
    Categories []Category `gorm:"many2many:products_categories;"`
    //         └── mesmo nome da tabela intermediária
    gorm.Model
}
```

**Tabela intermediária gerada:**

```sql
-- GORM cria automaticamente esta tabela:
CREATE TABLE products_categories (
    product_id BIGINT,     -- FK para products.id
    category_id BIGINT,    -- FK para categories.id
    PRIMARY KEY (product_id, category_id),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
```

**Criando relacionamentos Many-to-Many:**

```go
// 1. Criar categorias
category1 := Category{Name: "Eletrônicos"}
category2 := Category{Name: "Cozinha"}
db.Create(&category1)
db.Create(&category2)

// 2. Criar produto com múltiplas categorias
product := Product{
    Name:       "Notebook Gaming",
    Price:      2500.00,
    Categories: []Category{category1, category2}, // ✅ Associa a múltiplas categorias
}
db.Create(&product)

// GORM automaticamente:
// 1. Insere o produto na tabela products
// 2. Insere registros em products_categories
```

**Manipulando associações:**

```go
var product Product
db.First(&product, 1)

// Adicionar categoria ao produto existente
var newCategory Category
db.First(&newCategory, "name = ?", "Games")
db.Model(&product).Association("Categories").Append(&newCategory)

// Remover categoria
db.Model(&product).Association("Categories").Delete(&newCategory)

// Substituir todas as categorias
db.Model(&product).Association("Categories").Replace([]Category{category1})

// Contar associações
count := db.Model(&product).Association("Categories").Count()

// Limpar todas as associações
db.Model(&product).Association("Categories").Clear()
```

**Consultando Many-to-Many:**

```go
// Carregar categorias com seus produtos
var categories []Category
db.Preload("Products").Find(&categories)

for _, category := range categories {
    fmt.Printf("Categoria: %s\n", category.Name)
    for _, product := range category.Products {
        fmt.Printf("  - %s (R$ %.2f)\n", product.Name, product.Price)
    }
}

// Filtrar produtos por categoria
var products []Product
db.Joins("JOIN products_categories pc ON pc.product_id = products.id").
   Joins("JOIN categories c ON c.id = pc.category_id").
   Where("c.name = ?", "Eletrônicos").
   Find(&products)
```

**Preload com condições:**

```go
// Carregar apenas produtos acima de R$ 1000
db.Preload("Products", "price > ?", 1000).Find(&categories)

// Preload com ordenação
db.Preload("Products", func(db *gorm.DB) *gorm.DB {
    return db.Order("price DESC").Limit(5)
}).Find(&categories)
```

**Tabela intermediária customizada:**

```go
// Se precisar de campos extras na tabela intermediária
type ProductCategory struct {
    ProductID  uint `gorm:"primaryKey"`
    CategoryID uint `gorm:"primaryKey"`
    CreatedAt  time.Time
    Priority   int  // Campo extra
}

type Product struct {
    ID                int `gorm:"primaryKey"`
    Name              string
    ProductCategories []ProductCategory `gorm:"foreignKey:ProductID"`
    Categories        []Category        `gorm:"many2many:product_categories;joinForeignKey:ProductID;joinReferences:CategoryID"`
}
```

**Performance considerações:**

```go
// ❌ N+1 Problem - cada categoria faz query separada
for _, category := range categories {
    db.Model(&category).Association("Products").Find(&category.Products)
}

// ✅ Solução - usar Preload
db.Preload("Products").Find(&categories)

// ✅ Para grandes volumes - paginar
db.Preload("Products", func(db *gorm.DB) *gorm.DB {
    return db.Limit(10).Offset(page * 10)
}).Find(&categories)
```

### 18. Lock Otimista e Pessimista

**Localização:** `18-lock-otimista-e-pessimista/`

Controle de concorrência:

**Lock Pessimista:**

```go
// Lock durante a transação
tx := db.Begin()
var product Product
tx.Set("gorm:query_option", "FOR UPDATE").First(&product, 1)
// Operações...
tx.Commit()
```

**Lock Otimista:**

```go
type Product struct {
    ID      int `gorm:"primaryKey"`
    Name    string
    Price   float64
    Version int `gorm:"default:1"` // Campo de versão
}

// Atualização com verificação de versão
result := db.Model(&product).Where("version = ?", product.Version).
    Updates(map[string]interface{}{
        "price":   newPrice,
        "version": product.Version + 1,
    })

if result.RowsAffected == 0 {
    // Conflito detectado
    return errors.New("record was modified by another process")
}
```

## � Debugging e Dicas de Sintaxe

### Logs e Debug

```go
// Habilitar logs para ver SQL gerado
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// Log apenas queries lentas
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Warn).SlowThreshold(time.Second),
})

// Debug pontual
db.Debug().Create(&product) // Mostra SQL desta operação

// DryRun - mostra SQL sem executar
stmt := db.Session(&gorm.Session{DryRun: true}).Create(&product)
fmt.Println(stmt.Statement.SQL.String()) // SQL gerado
```

### Erros Comuns e Soluções

**1. Campo não atualiza:**

```go
// ❌ Não funciona - GORM ignora zero values
db.Model(&product).Updates(Product{Price: 0})

// ✅ Funciona - usa Select ou map
db.Model(&product).Select("price").Updates(Product{Price: 0})
// OU
db.Model(&product).Updates(map[string]interface{}{"price": 0})
```

**2. Preload não carrega:**

```go
// ❌ Nome errado do campo
db.Preload("Categorias").Find(&products) // Campo chama "Categories"

// ✅ Nome correto
db.Preload("Categories").Find(&products)

// ✅ Verificar se relacionamento está definido corretamente
type Product struct {
    CategoryID int
    Category   Category // ✅ Nome deve corresponder ao Preload
}
```

**3. Foreign Key não reconhecida:**

```go
// ❌ GORM não encontra FK automática
type Product struct {
    CatID    int
    Category Category
}

// ✅ Usar convenção ou especificar
type Product struct {
    CategoryID int      // ✅ Convenção: ModelNameID
    Category   Category `gorm:"foreignKey:CategoryID"` // ✅ Ou especificar
}
```

**4. Performance com relacionamentos:**

```go
// ❌ N+1 queries
var products []Product
db.Find(&products)
for _, product := range products {
    db.Model(&product).Association("Categories").Find(&product.Categories)
}

// ✅ Uma query com JOIN
db.Preload("Categories").Find(&products)

// ✅ Para casos complexos, usar Joins
db.Joins("LEFT JOIN categories ON categories.id = products.category_id").
   Select("products.*, categories.name as category_name").
   Find(&products)
```

### Sintaxe Avançada

**Where dinâmico:**

```go
// Construir query dinamicamente
query := db.Model(&Product{})

if name != "" {
    query = query.Where("name LIKE ?", "%"+name+"%")
}
if minPrice > 0 {
    query = query.Where("price >= ?", minPrice)
}
if categoryID > 0 {
    query = query.Where("category_id = ?", categoryID)
}

query.Find(&products)
```

**Scopes reutilizáveis:**

```go
// Definir scopes
func ExpensiveProducts(db *gorm.DB) *gorm.DB {
    return db.Where("price > ?", 1000)
}

func ByCategory(categoryName string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Joins("Category").Where("categories.name = ?", categoryName)
    }
}

// Usar scopes
db.Scopes(ExpensiveProducts, ByCategory("Eletrônicos")).Find(&products)
```

**Raw SQL quando necessário:**

```go
// Query raw
var result []Product
db.Raw("SELECT * FROM products WHERE price > ? AND created_at > ?", 1000, lastWeek).Scan(&result)

// Exec raw
db.Exec("UPDATE products SET price = price * 1.1 WHERE category_id = ?", categoryID)
```

## �🛠️ Comandos Úteis

### Configuração do Ambiente

```bash
# Subir container MySQL
docker compose up -d

# Conectar ao MySQL
mysql -h localhost -u root -p

# Criar database
CREATE DATABASE goexpert;
```

### Dependências Go

```bash
# Driver MySQL para database/sql
go get github.com/go-sql-driver/mysql

# GORM ORM
go get gorm.io/gorm
go get gorm.io/driver/mysql

# UUID
go get github.com/google/uuid
```

### Operações GORM

```bash
# Auto migration
db.AutoMigrate(&Product{}, &Category{})

# Drop table
db.Migrator().DropTable(&Product{})

# Check if table exists
db.Migrator().HasTable(&Product{})
```

## 📊 Boas Práticas

### Conexão com Banco

- **Pool de Conexões:** Configure adequadamente
- **Timeout:** Use contexto para controlar timeouts
- **SSL:** Configure SSL em produção
- **Prepared Statements:** Use para prevenir SQL injection

### Estrutura de Código

- **Repository Pattern:** Separe lógica de acesso a dados
- **Domain Models:** Mantenha modelos de domínio limpos
- **Migrations:** Use migrations para controle de schema
- **Environment Variables:** Configure conexão via variáveis

### Performance

- **Indexação:** Crie índices apropriados
- **Preload:** Use para evitar N+1 queries
- **Batch Operations:** Para operações em massa
- **Connection Pooling:** Configure pool adequadamente

### Segurança

- **SQL Injection:** Sempre use prepared statements
- **Validação:** Valide dados antes de persistir
- **Autorização:** Implemente controle de acesso
- **Audit Trail:** Mantenha log de alterações

### Relacionamentos

```go
// Estrutura recomendada
type User struct {
    ID       uint `gorm:"primaryKey"`
    Name     string
    Posts    []Post    `gorm:"foreignKey:UserID"`
    Profile  Profile   `gorm:"foreignKey:UserID"`
    Tags     []Tag     `gorm:"many2many:user_tags;"`
    gorm.Model
}

// Com índices
type Post struct {
    ID     uint   `gorm:"primaryKey"`
    Title  string `gorm:"index"`
    UserID uint   `gorm:"index"`
    User   User
}
```

## 🎯 Objetivos de Aprendizado

Após estudar estes exemplos, você deve ser capaz de:

1. ✅ Configurar ambiente de banco de dados com Docker
2. ✅ Conectar aplicação Go ao MySQL
3. ✅ Executar operações CRUD com SQL puro
4. ✅ Usar prepared statements para segurança
5. ✅ Implementar operações com GORM
6. ✅ Modelar relacionamentos de banco de dados
7. ✅ Implementar soft delete
8. ✅ Otimizar consultas e relacionamentos
9. ✅ Aplicar técnicas de controle de concorrência
10. ✅ Seguir boas práticas de segurança e performance

## 🔄 Fluxo de Desenvolvimento

### Setup Inicial

1. Configurar Docker Compose para MySQL
2. Criar modelos de dados
3. Configurar conexão com banco
4. Implementar migrations

### Desenvolvimento com SQL Puro

1. Criar prepared statements
2. Implementar CRUD operations
3. Tratar erros adequadamente
4. Fazer cleanup de recursos

### Desenvolvimento com GORM

1. Definir structs com tags GORM
2. Configurar auto migration
3. Implementar operações básicas
4. Modelar relacionamentos
5. Otimizar consultas

## 📖 Recursos Adicionais

- [GORM Documentation](https://gorm.io/docs/)
- [Go MySQL Driver](https://github.com/go-sql-driver/mysql)
- [Database/SQL Tutorial](https://go.dev/doc/tutorial/database-access)
- [SQL Best Practices](https://use-the-index-luke.com/)
- [MySQL Performance Tuning](https://dev.mysql.com/doc/refman/8.0/en/optimization.html)
