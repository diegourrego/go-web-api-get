package handler_test

import (
	"first_api/internal"
	"first_api/internal/handler"
	"first_api/internal/repository"
	"first_api/internal/service"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockProductStorage struct {
	LoadDataFunc func() (map[int]internal.Product, error)
	SaveDataFunc func(product internal.Product) error
}

func (m *MockProductStorage) LoadData() (map[int]internal.Product, error) {
	return m.LoadDataFunc()
}

func (m *MockProductStorage) SaveData(product internal.Product) error {
	return m.SaveDataFunc(product)
}

func TestDefaultProducts_GetProducts(t *testing.T) {
	// Arrange
	db := map[int]internal.Product{
		1: {ID: 1, Name: "Product 1", Quantity: 1, CodeValue: "CV01", IsPublished: false,
			Expiration: "12/12/2021", Price: 1.1},
		2: {ID: 2, Name: "Product 2", Quantity: 2, CodeValue: "CV02", IsPublished: false,
			Expiration: "22/03/2022", Price: 2.2},
		3: {ID: 3, Name: "Product 3", Quantity: 3, CodeValue: "CV03", IsPublished: true,
			Expiration: "03/08/2023", Price: 3.3},
	}

	mockStorage := &MockProductStorage{
		LoadDataFunc: func() (map[int]internal.Product, error) {
			return db, nil
		},
		SaveDataFunc: func(product internal.Product) error {
			return nil
		},
	}

	// repository
	rp := repository.NewProductMap(db, 0)
	//service
	sv := service.NewProductDefault(rp)
	//handler
	hd := handler.NewDefaultProducts(sv, mockStorage)

	t.Run("success case01 - should return a list of products", func(t *testing.T) {
		// Act
		rq := httptest.NewRequest("GET", "/products", nil)
		// Insert Authorization Token
		rq.Header.Add("Authorization", "123456")
		res := httptest.NewRecorder()
		hdFunc := hd.GetProducts()
		hdFunc(res, rq)

		// Assert

		expectedCode := http.StatusOK
		expectedBody := `{"message": "Products obtained successfully", "data": [
	{"id": 1, "name": "Product 1", "quantity": 1, "code_value": "CV01", "is_published": false,
	"expiration": "12/12/2021", "price":1.1 },
	{"id": 2, "name": "Product 2", "quantity": 2, "code_value": "CV02", "is_published": false,
	"expiration": "22/03/2022", "price":2.2 },
	{"id": 3, "name": "Product 3", "quantity": 3, "code_value": "CV03", "is_published": true,
	"expiration": "03/08/2023", "price":3.3 }
], "error": false}`
		actualBody := res.Body.String()
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header())
		require.JSONEq(t, expectedBody, actualBody)
	})

	t.Run("failure - case01: should returns unauthorized code", func(t *testing.T) {
		// Act
		rq := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()
		hdFunc := hd.GetProducts()
		hdFunc(res, rq)
		// Assert
		expectedCode := http.StatusUnauthorized
		expectedBody := `{"message": "Unauthorized", "data": null, "error": true}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}

		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedHeader, res.Header())
		require.JSONEq(t, expectedBody, res.Body.String())
	})

}
