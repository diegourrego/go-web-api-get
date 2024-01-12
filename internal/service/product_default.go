package service

import (
	"first_api/internal"
	"strconv"
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

func (pd *ProductDefault) GetProductByID(productID string) (internal.Product, error) {
	// Válido que id sea un int
	id, err := strconv.Atoi(productID)
	if err != nil {
		return internal.Product{}, internal.ErrInvalidID
	}

	// Válido que el ID sea mayor a 0
	if id <= 0 {
		return internal.Product{}, internal.ErrInvalidID
	}

	product, err := pd.rp.GetProductByID(id)
	if err != nil {
		return internal.Product{}, err
	}

	return product, nil
}
