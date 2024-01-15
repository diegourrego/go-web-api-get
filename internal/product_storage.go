package internal

type ProductStorage interface {
	LoadData() (map[int]Product, error)
	SaveData(product Product) error
}
