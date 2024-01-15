package storage

import (
	"encoding/json"
	"first_api/internal"
	"os"
)

type DataLoaded struct {
	data map[int]internal.Product `json:"data"`
}

func NewDataLoaded() *DataLoaded {
	return &DataLoaded{}
}

func (dl *DataLoaded) LoadData() (map[int]internal.Product, error) {
	file, err := os.ReadFile("./data/data.json")
	if err != nil {
		return nil, err
	}

	var products []internal.Product
	if err := json.Unmarshal(file, &products); err != nil {
		return nil, err
	}

	dl.data = make(map[int]internal.Product)
	for _, product := range products {
		dl.data[product.ID] = product
	}

	return dl.data, nil
}

func (dl *DataLoaded) SaveData(product internal.Product) error {
	dl.data[product.ID] = product

	var products []internal.Product
	for _, product := range dl.data {
		products = append(products, product)
	}

	// Guardo el map en un archivo
	file, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		return err
	}

	if err := os.WriteFile("./data/data.json", file, 0644); err != nil {
		return err
	}
	return nil
}
