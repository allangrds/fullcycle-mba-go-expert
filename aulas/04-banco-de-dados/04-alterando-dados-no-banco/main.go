package main

import (
	"database/sql"

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
}
