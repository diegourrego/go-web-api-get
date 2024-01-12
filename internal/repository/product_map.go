package repository

import (
	"first_api/internal"
)

type ProductMap struct {
	// Acá uso el map que usaré como bd en memoria
	db     map[int]internal.Product
	lastID int
}

func NewProductMap(db map[int]internal.Product, lastID int) *ProductMap {
	return &ProductMap{
		db:     db,
		lastID: lastID,
	}
}

func (pm *ProductMap) GetProducts() map[int]internal.Product {
	return pm.db
}

func (pm *ProductMap) GetProductByID(productID int) (internal.Product, error) {

	product, ok := pm.db[productID]
	if !ok {
		return internal.Product{}, internal.ErrProductNotFound
	}

	return product, nil
}

func (pm *ProductMap) GetProductWithPriceHigherThan(productPrice float64) (map[int]internal.Product, error) {
	productsFounded := make(map[int]internal.Product)

	for _, product := range pm.db {
		if product.Price > productPrice {
			productsFounded[product.ID] = product
		}
	}

	return productsFounded, nil
}
