package dto

type CreateProductInput struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
