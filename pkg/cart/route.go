package cart

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (svc *Service) Routes(r chi.Router) chi.Router {
	r.Post("/cart/item", svc.PostCartItemHandler)
	r.Get("/cart", svc.GetCartHandler)
	r.Put("/cart/item/:id", svc.PutCartItemHandler)
	r.Delete("/cart/item/:id", svc.DeleteCartItemHandler)
}

func (svc *Service) PostCartItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}

func (svc *Service) GetCartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}

func (svc *Service) PutCartItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}

func (svc *Service) DeleteCartItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}
