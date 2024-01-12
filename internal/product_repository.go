package internal

type ProductRepository interface {
	// GetProducts Get all the products
	GetProducts() map[int]Product
}
