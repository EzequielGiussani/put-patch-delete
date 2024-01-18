package application

import (
	"app/internal"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DefaultHttp struct {
	address string
}

func NewDefaultHttp(address string) *DefaultHttp {
	defaultAddrs := ":8080"

	if address != "" {
		defaultAddrs = address
	}

	return &DefaultHttp{
		address: defaultAddrs,
	}
}

func (s *DefaultHttp) Run() error {

	rp := repository.NewProductMap(make(map[int]internal.Product), 0)

	sv := service.NewProductDefault(rp)

	hd := handler.NewDefaultProducts(sv)

	rt := chi.NewRouter()

	rt.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	rt.Post("/products", hd.Create())
	rt.Get("/products/{id}", hd.GetById())
	rt.Put("/products/{id}", hd.Update())
	rt.Patch("/products/{id}", hd.UpdatePartial())
	rt.Delete("/products/{id}", hd.Delete())

	return http.ListenAndServe(s.address, rt)
}
