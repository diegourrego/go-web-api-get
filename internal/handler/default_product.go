package handler

import (
	"encoding/json"
	"first_api/internal"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sort"
	"strconv"
)

type DefaultProducts struct {
	sv internal.ProductService
}

func NewDefaultProducts(sv internal.ProductService) *DefaultProducts {
	return &DefaultProducts{
		sv: sv,
	}
}

type BodyResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   bool   `json:"error"`
}

func (dp *DefaultProducts) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productsMap := dp.sv.GetProducts()

		type GetProductsBodyResponse struct {
			Message string             `json:"message"`
			Data    []internal.Product `json:"data"`
			Error   bool               `json:"error"`
		}

		var products []internal.Product

		for _, product := range productsMap {
			products = append(products, product)
		}

		// Ordeno lo productos por ID
		sort.Slice(products, func(i, j int) bool {
			return products[i].ID < products[j].ID
		})

		body := GetProductsBodyResponse{
			Message: "Products obtained successfully",
			Data:    products,
			Error:   false,
		}

		code := http.StatusOK
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}

func (dp *DefaultProducts) GetProductByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		// Válido que id sea un int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
				Message: internal.ErrInvalidID.Error(),
				Data:    nil,
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		product, err := dp.sv.GetProductByID(idInt)
		if err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
				//Message: "¡Ops, something went wrong! - Bad Request",
				Message: err.Error(),
				Data:    nil,
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		code := http.StatusOK
		body := BodyResponse{
			Message: "Product found - Everything's OK",
			Data:    product,
			Error:   false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}
