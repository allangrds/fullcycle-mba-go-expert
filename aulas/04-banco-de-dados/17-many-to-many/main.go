package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Products []Product `gorm:"many2many:products_categories;"`
	gorm.Model
}

type Product struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Price        float64
	Categories   []Category `gorm:"many2many:products_categories;"`
	SerialNumber SerialNumber
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

	// category2 := Category{Name: "Cozinha"}
	// db.Create(&category2)

	// // // //Create products
	// productToCreate := Product{
	// 	Name:       "Notebook",
	// 	Price:      1000.00,
	// 	Categories: []Category{category, category2},
	// }
	// db.Create(&productToCreate)

	// productToCreate = Product{
	// 	Name:       "Geladeira",
	// 	Price:      1200.00,
	// 	Categories: []Category{category, category2},
	// }
	// db.Create(&productToCreate)

	var categories []Category
	// err = db.Model(&Category{}).Preload("Products").Find(&categories).Error
	// err = db.Model(&Category{}).Preload("Products").Preload("Products.SerialNumber").Find(&categories).Error
	err = db.Model(&Category{}).Preload("Products").Find(&categories).Error
	if err != nil {
		panic(err)
	}

	for _, category := range categories {
		fmt.Println(category.Name, ":")
		for _, product := range category.Products {
			println("- ", product.Name)
		}
	}
}
