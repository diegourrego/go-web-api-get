package internal

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	// GetProducts Get all the products
	GetProducts() map[int]Product
	GetProductByID(productID int) (Product, error)
}
