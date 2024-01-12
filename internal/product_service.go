package internal

import "errors"

var (
	ErrInvalidID          = errors.New("id must be a positive number")
	ErrInvalidPrice       = errors.New("price must be 0.0 or higher")
	ErrInvalidPriceFormat = errors.New("price must be a number")
)

type ProductService interface {
	GetProducts() map[int]Product
	GetProductByID(productID int) (Product, error)
	GetProductWithPriceHigherThan(productPrice float64) (map[int]Product, error)
}
