package main

import (
	"api_get/handler"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	// Creamos el router
	router := chi.NewRouter()
	// Creamos el handler
	h := handler.NewHandler()
	// Requirement 2:

	// a.
	router.Get("/ping", h.GetPong())

	//b.
	router.Get("/products", h.GetProducts())

	// c.
	router.Get("/products/{id}", h.GetProductById())

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
