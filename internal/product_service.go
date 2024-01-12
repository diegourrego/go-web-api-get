package internal

type ProductService interface {
	GetProducts() map[int]Product
}
