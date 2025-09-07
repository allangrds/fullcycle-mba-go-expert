package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/dto"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/entity"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/infra/database"
	entityPkg "github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/pkg/entity"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// CreateProduct godoc
// @Summary Create product
// @Description Create product
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductInput true "product request"
// @Success 201
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /products [post]
// @Security ApiKeyAuth
func (productHandler *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productDto dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&productDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)

		return
	}

	// Fazer diretamente acesso da entidade no coração não é comum. Em vez disso, usaremos no futuro um use case(clean arch)
	product, err := entity.NewProduct(productDto.Name, productDto.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)

		return
	}

	err = productHandler.ProductDB.Create(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (productHandler *ProductHandler) GetProduct(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if id == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := productHandler.ProductDB.FindByID(id)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	//Aqui o json.NewEncoder já escreve diretamente no response writer
	//O .Encode faz o encode do objeto passado como parâmetro
	//Outra forma seria usar o json.Marshal(product) e depois w.Write(). Exemplo: data, _ := json.Marshal(product) w.Write(data)
	json.NewEncoder(writer).Encode(product)
}

func (productHandler *ProductHandler) UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if id == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product

	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = productHandler.ProductDB.Update(&product)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (productHandler *ProductHandler) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	if id == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err := productHandler.ProductDB.Delete(id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (productHandler *ProductHandler) GetProducts(writer http.ResponseWriter, request *http.Request) {
	pageString := request.URL.Query().Get("page")
	limitString := request.URL.Query().Get("limit")
	sort := request.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(pageString)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limitString)
	if err != nil {
		limitInt = 0
	}

	products, err := productHandler.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(products)
}
