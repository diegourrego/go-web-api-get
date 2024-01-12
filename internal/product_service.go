package internal

import "errors"

var (
	ErrInvalidID = errors.New("id must be a positive number")
)

type ProductService interface {
	GetProducts() map[int]Product
	GetProductByID(productID string) (Product, error)
}
