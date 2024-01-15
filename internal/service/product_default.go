package service

import (
	"first_api/internal"
	"strconv"
	"strings"
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

func (pd *ProductDefault) GetProductWithPriceHigherThan(productPrice float64) (map[int]internal.Product, error) {
	// Validations
	if productPrice < 0.0 {
		return nil, internal.ErrInvalidPrice
	}

	productsFounded, err := pd.rp.GetProductWithPriceHigherThan(productPrice)
	if err != nil {
		return nil, err
	}

	return productsFounded, nil
}

func (pd *ProductDefault) Create(newProduct internal.Product) (internal.Product, error) {
	// Validaciones
	if err := validateBodyFields(newProduct); err != nil {
		return internal.Product{}, err
	}

	product, err := pd.rp.Create(newProduct)
	if err != nil {
		return internal.Product{}, err
	}

	return product, nil
}

func (pd *ProductDefault) Update(newProduct *internal.Product) (internal.Product, error) {
	// Validaciones
	//if err := validateBodyFields(newProduct); err != nil {
	//	return internal.Product{}, err
	//}

	productUpdated, err := pd.rp.Update(newProduct)
	if err != nil {
		return internal.Product{}, err
	}

	return productUpdated, nil
}

func (pd *ProductDefault) Delete(productID int) error {
	// Validations
	if productID <= 0 {
		return internal.ErrInvalidID
	}

	err := pd.rp.Delete(productID)
	if err != nil {
		return err
	}

	return nil
}

func validateBodyFields(product internal.Product) error {
	// Validations
	// Ningún campo puede estar vacío
	if product.Expiration == "" || product.Name == "" || product.CodeValue == "" ||
		product.Price == 0.0 || product.Quantity == 0 {
		return internal.ErrEmptyField
	}

	// La fecha de vencimiento debe tener el formato xx/xx/xxxx
	expirationArray := strings.Split(product.Expiration, "/")

	if len(expirationArray) != 3 || len(expirationArray[0]) != 2 || len(expirationArray[1]) != 2 || len(expirationArray[2]) != 4 {
		return internal.ErrInvalidDateFormat
	}

	// Día, mes y año deben ser valores válidos
	_, err := strconv.Atoi(expirationArray[0])
	_, err = strconv.Atoi(expirationArray[1])
	_, err = strconv.Atoi(expirationArray[2])
	if err != nil {
		return internal.ErrInvalidDateFormat
	}

	return nil
}
