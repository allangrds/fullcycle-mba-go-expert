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
	gorm.Model
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
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

	var product Product
	db.First(&product)
	product.Name = "Mouse"
	db.Save(&product)

	var product2 Product
	db.First(&product2)
	fmt.Println(product2.Name)
	db.Delete(&product)
}
