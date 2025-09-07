package database

import (
	"fmt"
	"testing"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, _ := entity.NewProduct("Product Test", 10)
	productDb := NewProduct(db)
	err = productDb.Create(product)

	assert.Nil(t, err)
	assert.NotEmpty(t, product.ID)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestFindProductById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, _ := entity.NewProduct("Product Test", 10)
	productDb := NewProduct(db)
	err = productDb.Create(product)

	productFound, err := productDb.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestFindProductAll(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	productDb := NewProduct(db)
	for i := 1; i < 25; i++ {
		product, _ := entity.NewProduct(fmt.Sprintf("Product %d", i), 10+i)
		productDb.Create(product)
	}

	products, err := productDb.FindAll(1, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDb.FindAll(2, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, _ := entity.NewProduct("Product Test", 10)
	productDb := NewProduct(db)
	productDb.Create(product)

	product.Name = "Product Test Updated"
	err = productDb.Update(product)
	assert.Nil(t, err)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestDeleteProduct(t *testing.T) {
	//Criação de registro em memória
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product Test", 10)
	productDb := NewProduct(db)
	productDb.Create(product)

	//Exclusão do registro
	err = productDb.Delete(product.ID.String())
	assert.Nil(t, err)

	//Verificação se o registro foi excluído
	productFound, err := productDb.FindByID(product.ID.String())
	assert.Nil(t, productFound)
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, "record not found", err.Error())
}
