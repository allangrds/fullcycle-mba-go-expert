package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/dto"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/entity"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/infra/database"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (productHandler *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productDto dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&productDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fazer diretamente acesso da entidade no coração não é comum. Em vez disso, usaremos no futuro um use case(clean arch)
	product, err := entity.NewProduct(productDto.Name, productDto.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = productHandler.ProductDB.Create(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
