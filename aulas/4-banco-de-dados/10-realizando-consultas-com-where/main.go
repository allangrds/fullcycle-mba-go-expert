package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID    int `gorm:"primaryKey"`
	Name  string
	Price float64
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/goexpert"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//Cria a tabela com base na struct
	db.AutoMigrate(&Product{})

	// db.Create(&Product{
	// 	Name:  "Notebook Dell",
	// 	Price: 1000.00,
	// })

	//create batch
	// products := []Product{
	// 	{
	// 		Name:  "Notebook Dell",
	// 		Price: 1000.00,
	// 	},
	// 	{
	// 		Name:  "Notebook Samsung",
	// 		Price: 1000.00,
	// 	},
	// }

	// db.Create(products)

	//select one
	// var product Product
	// db.First(&product, 1)
	// fmt.Println(product)
	// db.First(&product, "name = ?", "Notebook Samsung")
	// fmt.Println(product)

	//select all
	// var products []Product
	// db.Limit(2).Offset(1).Find(&products)
	// for _, product := range products {
	// 	fmt.Println(product)
	// }

	//Where explícito para preçø
	// var products []Product
	// db.Where("price > ?", 100).Find(&products)
	// for _, product := range products {
	// 	fmt.Println(product)
	// }

	//Where explícito para name
	var products []Product
	db.Where("name like ?", "%Notebook%").Find(&products)
	for _, product := range products {
		fmt.Println(product)
	}
}
