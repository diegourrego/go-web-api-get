package internal

type ProductRepository interface {
	GetProducts()
	GetProductByID()
	Save(product *Product) (err error)
}
