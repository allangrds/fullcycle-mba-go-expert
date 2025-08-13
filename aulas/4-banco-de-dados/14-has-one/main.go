package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID   int `gorm:"primaryKey"`
	Name string
	gorm.Model
}

type Product struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Price        float64
	CategoryId   int
	Category     Category
	SerialNumber SerialNumber //pra eu recuperar o serial number, mesmo a chave estando no serial number, 1 pra 1
	gorm.Model
}

type SerialNumber struct {
	ID        int `gorm:"primaryKey"`
	Number    string
	ProductID int
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//Cria a tabela com base na struct
	db.AutoMigrate(&Product{}, &Category{}, &SerialNumber{})

	//Create category
	// category := Category{Name: "Eletrônicos"}
	// db.Create(&category)

	// //Create products
	// productToCreate := Product{
	// 	Name:       "Notebook",
	// 	Price:      1000.00,
	// 	CategoryId: category.ID,
	// }
	// db.Create(&productToCreate)

	// //Create Serial Number
	// serialNumber := SerialNumber{
	// 	Number:    "123ABC",
	// 	ProductID: 3,
	// }
	// db.Create(&serialNumber)

	var products []Product
	db.Preload("Category").Preload("SerialNumber").Find(&products)
	for _, product := range products {
		fmt.Println(product.Name, product.Category.Name, product.SerialNumber.Number)
	}
}
