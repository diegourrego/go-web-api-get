package internal

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrCodeValueRepeated = errors.New("code value already exists in products")
)

type ProductRepository interface {
	GetProducts() map[int]Product
	GetProductByID(productID int) (Product, error)
	GetProductWithPriceHigherThan(productPrice float64) (map[int]Product, error)

	Create(newProduct Product) (Product, error)

	Update(newProduct Product) (Product, error)
	Delete(productID int) error
}
