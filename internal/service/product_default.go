package service

import (
	"first_api/internal"
)

// ¡Acá debemos hacer las validaciones!

type ProductDefault struct {
	rp internal.ProductRepository // Inyecto el repository
}

func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

func (pd *ProductDefault) GetProducts() map[int]internal.Product {
	products := pd.rp.GetProducts()
	// Acá se puede ejecutar lógica adicional como validaciones
	return products
}

func (pd *ProductDefault) GetProductByID(productID int) (internal.Product, error) {

	// Válido que el ID sea mayor a 0
	if productID <= 0 {
		return internal.Product{}, internal.ErrInvalidID
	}

	product, err := pd.rp.GetProductByID(productID)
	if err != nil {
		return internal.Product{}, err
	}

	return product, nil
}
