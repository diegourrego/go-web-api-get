package loader

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
