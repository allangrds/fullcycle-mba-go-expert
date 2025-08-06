package main

import (
	"database/sql"
	"fmt"

	//O _ é um blank identifier, portanto vai deixar compilar
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

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

func insertProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("insert into products(id, name, price) values (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	// _ foi usado pois eu não quero um valor do resultado, então faço uso dele, e como o err já existe, faço uso do =
	_, err = stmt.Exec(product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func updateProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("update products set name = ?, price = ? where id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	// _ foi usado pois eu não quero um valor do resultado, então faço uso dele, e como o err já existe, faço uso do =
	_, err = stmt.Exec(product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func selectProduct(db *sql.DB, id string) (*Product, error) {
	stmt, err := db.Prepare("select id, name, price from products where id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var product Product

	//Para usar contexto para controlar um timeout, por exemplo
	// err = stmt.QueryRowContext(ctx, id).Scan(&product.ID, &product.Name, &product.Price)

	//QueryRow buscando apenas 1 linha
	//Scan para atribuir dados de coluna a campos de estrutura
	//& local na memória onde o product está armazenado para poder alterá-lo
	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

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

func deleteProduct(db *sql.DB, id string) error {
	stmt, err := db.Prepare("delete from products where id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	//buscar, usar query, executar, exec
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	productFirst := NewProduct("Notebook", 1999.80)
	err = insertProduct(db, productFirst)
	if err != nil {
		panic(err)
	}

	productFirst.Price = 100.00
	err = updateProduct(db, productFirst)
	if err != nil {
		panic(err)
	}

	// res, err := selectProduct(db, productFirst.ID)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Product: %v, possui o preço de %.2f", res.Name, res.Price)

	products, err := selectProducts(db)
	if err != nil {
		panic(err)
	}
	for _, product := range products {
		fmt.Printf("\nProduct: %v, possui o preço de %.2f e o ID %v", product.Name, product.Price, product.ID)
	}

	err = deleteProduct(db, productFirst.ID)
	if err != nil {
		panic(err)
	}

	products, err = selectProducts(db)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n")
	for _, product := range products {
		fmt.Printf("\nProduct: %v, possui o preço de %.2f e o ID %v", product.Name, product.Price, product.ID)
	}
}
