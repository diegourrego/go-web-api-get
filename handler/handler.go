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

//type ControllerProducts struct {
//	db map[int]*Product
//}
//
//func NewControllerProducts(products map[int]*Product) *ControllerProducts {
//	return &ControllerProducts{db: products}
//}

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

func (h *MyHandler) GetProductsGreaterThan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var productsHigherThan []Product
		priceGt := r.URL.Query().Get("priceGt")
		priceInt, err := strconv.Atoi(priceGt)

		// If id != int:
		if err != nil || priceInt < 0 {
			body := MyResponse{Message: "price must be a positive number", Data: nil}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(body); err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		for _, product := range h.data.Products {
			if product.Price > float64(priceInt) {
				productsHigherThan = append(productsHigherThan, product)
			}
		}

		body := MyResponse{Message: "All's good", Data: productsHigherThan}
		noProdFoundMessage := fmt.Sprintf("No products found with price greater than %d", priceInt)

		if len(productsHigherThan) == 0 {
			body = MyResponse{Message: noProdFoundMessage, Data: nil}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			fmt.Println(err)
			return
		}

	}
}

type RequestBodyProduct struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type ResponseBodyProduct struct {
	Message string
	Data    *struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Quantity    int     `json:"quantity"`
		CodeValue   string  `json:"code_value"`
		IsPublished bool    `json:"is_published"`
		Expiration  string  `json:"expiration"`
		Price       float64 `json:"price"`
	}
	Error bool
}

func (h *MyHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody RequestBodyProduct
		// Toca convertir el Json a estructura
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			code := http.StatusBadRequest
			body := &ResponseBodyProduct{
				Message: "Bad Request. check your body data",
				Data:    nil,
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		}

		// Serializar
		p := Product{
			ID:          len(h.data.Products) + 1,
			Name:        requestBody.Name,
			Quantity:    requestBody.Quantity,
			CodeValue:   requestBody.CodeValue,
			IsPublished: requestBody.IsPublished,
			Expiration:  requestBody.Expiration,
		}

		h.data.Products = append(h.data.Products, p)

		// Le respondemos al cliente
		code := http.StatusCreated
		body := ResponseBodyProduct{
			Message: "Product created successfully",
			Data: &struct {
				ID          int     `json:"id"`
				Name        string  `json:"name"`
				Quantity    int     `json:"quantity"`
				CodeValue   string  `json:"code_value"`
				IsPublished bool    `json:"is_published"`
				Expiration  string  `json:"expiration"`
				Price       float64 `json:"price"`
			}{ID: p.ID, Name: p.Name, Quantity: p.Quantity, CodeValue: p.CodeValue,
				IsPublished: p.IsPublished, Expiration: p.Expiration, Price: p.Price},
			Error: false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			fmt.Println(err)
			return
		}

	}
}
