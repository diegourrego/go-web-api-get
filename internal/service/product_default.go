package service

import "first_api/internal"

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
