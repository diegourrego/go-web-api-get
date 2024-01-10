package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type Products struct {
	Products []Product
}

func NewProducts() *Products {
	return &Products{}
}

func (p *Products) LoadProducts() {
	file, err := os.ReadFile("./data/data.json")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	err = json.Unmarshal(file, &p.Products)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
}

type MyHandler struct {
	data Products
}

type MyResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewHandler() *MyHandler {
	h := &MyHandler{
		data: *NewProducts(),
	}
	h.data.LoadProducts()
	return h
}

func (h *MyHandler) GetPong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (h *MyHandler) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		body := MyResponse{
			Message: "Todo OK",
			Data:    h.data.Products,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		err := json.NewEncoder(w).Encode(body)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (h *MyHandler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)

		// If id != int:
		if err != nil || idInt < 0 {
			body := MyResponse{Message: "id must be a positive number", Data: nil}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(body); err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		for _, product := range h.data.Products {
			if product.ID == idInt {
				body := MyResponse{Message: "product founded", Data: product}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(body); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		errMessage := fmt.Sprintf("Product with id %d not found", idInt)
		body := MyResponse{Message: errMessage, Data: nil}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			fmt.Println(err)
			return
		}

	}
}
