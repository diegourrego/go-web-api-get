package handler

import (
	"encoding/json"
	"errors"
	"first_api/internal"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type DefaultProducts struct {
	sv internal.ProductService
	ld internal.ProductStorage
}

func NewDefaultProducts(sv internal.ProductService, ld internal.ProductStorage) *DefaultProducts {
	return &DefaultProducts{
		sv: sv,
		ld: ld,
	}
}

type BodyResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   bool   `json:"error"`
}

type ProductJSON struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type BodyRequestProductJSON struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func (dp *DefaultProducts) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := validateToken(w, r); err != nil {
			log.Println(err)
			return
		}

		productsMap := dp.sv.GetProducts()

		var products []internal.Product

		for _, product := range productsMap {
			products = append(products, product)
		}

		// Ordeno lo productos por ID
		sort.Slice(products, func(i, j int) bool {
			return products[i].ID < products[j].ID
		})

		body := BodyResponse{
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

		if err := validateToken(w, r); err != nil {
			log.Println(err)
			return
		}

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

func (dp *DefaultProducts) GetProductsWithPriceHigherThan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := validateToken(w, r); err != nil {
			log.Println(err)
			return
		}

		priceGtStr := r.URL.Query().Get("priceGt")
		price, err := strconv.ParseFloat(priceGtStr, 64)
		if err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
				Message: internal.ErrInvalidPriceFormat.Error(),
				Data:    nil,
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		productsObtained, err := dp.sv.GetProductWithPriceHigherThan(price)
		if err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
				Message: err.Error(),
				Data:    nil,
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		var products []internal.Product

		for _, product := range productsObtained {
			products = append(products, product)
		}

		// Organizo los productos por precio
		sort.Slice(products, func(i, j int) bool {
			return products[i].Price < products[j].Price
		})

		code := http.StatusOK
		body := BodyResponse{
			Message: "Products obtained successfully",
			Data:    products,
			Error:   false,
		}

		if len(products) == 0 {
			body = BodyResponse{
				Message: "No products for the specified price found",
				Data:    products,
				Error:   false,
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)

	}
}

func (dp *DefaultProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := validateToken(w, r); err != nil {
			log.Println(err)
			return
		}

		var body BodyRequestProductJSON
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			code := http.StatusBadRequest
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			w.Write([]byte("Invalid body"))
			return
		}

		// Serializamos el producto
		product := internal.Product{
			CodeValue:   body.CodeValue,
			Expiration:  body.Expiration,
			Price:       body.Price,
			Name:        body.Name,
			Quantity:    body.Quantity,
			IsPublished: body.IsPublished,
		}

		// TODO: preguntar si el id se debería asignar acá o en el servicio
		productCreated, err := dp.sv.Create(product)
		// Save the product in file
		err = dp.ld.SaveData(productCreated)
		if err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
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
		respBody := BodyResponse{
			Message: "Product created successfully",
			Data:    productCreated,
			Error:   false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(respBody)

	}
}

func (dp *DefaultProducts) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := validateToken(w, r); err != nil {
			log.Println(err)
			return
		}

		// Obtengo el id
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
				Message: "Bad request - invalid ID",
				Data:    nil,
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		}

		var bodyReq BodyRequestProductJSON
		if err := json.NewDecoder(r.Body).Decode(&bodyReq); err != nil {
			code := http.StatusBadRequest
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			w.Write([]byte("Invalid body"))
			return
		}

		// Debo serializar
		product := internal.Product{
			ID:          id,
			CodeValue:   bodyReq.CodeValue,
			Expiration:  bodyReq.Expiration,
			Price:       bodyReq.Price,
			Name:        bodyReq.Name,
			Quantity:    bodyReq.Quantity,
			IsPublished: bodyReq.IsPublished,
		}

		productUpdated, err := dp.sv.Update(product)
		if err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
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
			Message: "Product updated successfully",
			Data:    productUpdated,
			Error:   false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)

	}
}

func (dp *DefaultProducts) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := validateToken(w, r); err != nil {
			log.Println(err)
			return
		}

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
				Message: "Bad request - invalid ID",
				Data:    nil,
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
			return
		}

		if err := dp.sv.Delete(id); err != nil {
			code := http.StatusBadRequest
			body := BodyResponse{
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
			Message: "Product deleted successfully",
			Data:    nil,
			Error:   false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}

//func ValidateKeyExist(mp map[string]any, keys ...string) (err error) {
//	for _, k := range keys {
//		if _, ok := mp[k]; !ok {
//			return fmt.Errorf("key %s not found", k)
//		}
//	}
//	return
//}

func validateToken(w http.ResponseWriter, r *http.Request) error {
	// Set token in header
	// Validate token
	token := r.Header.Get("Authorization")
	if token != os.Getenv("TOKEN") {
		body := BodyResponse{
			Message: "Unauthorized",
			Data:    nil,
			Error:   true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(body)
		return errors.New("unauthorized")
	}
	return nil
}
