package application

import (
	"first_api/internal/handler"
	"first_api/internal/repository"
	"first_api/internal/service"
	"first_api/internal/storage"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type DefaultHTTP struct {
	addr string
}

func NewDefaultHTTP(addr string) *DefaultHTTP {
	return &DefaultHTTP{
		addr: addr,
	}
}

func (h *DefaultHTTP) Run() (err error) {
	// storage
	loader := storage.NewDataLoaded()
	loadData, err := loader.LoadData()
	if err != nil {
		fmt.Println(err)
		return
	}

	// repository
	rp := repository.NewProductMap(loadData, 0)
	// service
	sv := service.NewProductDefault(rp)
	// handler
	hd := handler.NewDefaultProducts(sv, loader)
	// router
	rt := chi.NewRouter()

	// endpoints
	rt.Route("/products", func(rt chi.Router) {
		rt.Get("/", hd.GetProducts())
		rt.Get("/{id}", hd.GetProductByID())
		rt.Get("/search", hd.GetProductsWithPriceHigherThan())
		rt.Post("/", hd.Create())
		rt.Put("/{id}", hd.Update())
		rt.Delete("/{id}", hd.Delete())
	})

	// run http service
	err = http.ListenAndServe(h.addr, rt)
	return
}
