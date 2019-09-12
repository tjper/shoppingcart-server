package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tjper/shoppingcart-server/service/cart"

	"github.com/go-chi/chi"
)

// CartRoutes defines the cart resources REST endpoints.
func (svc *Service) CartRoutes(r chi.Router) {
	r.Post("/cart/item", svc.PostCartItemHandler())
	r.Get("/cart/{userId}", svc.GetCartHandler())
	r.Put("/cart/item/{id}", svc.PutCartItemHandler())
	r.Delete("/cart/item/{id}", svc.DeleteCartItemHandler())
}

// PostCreatItemHandler creates a CartItem resource on the service.
func (svc *Service) PostCartItemHandler() http.HandlerFunc {
	type (
		Request struct {
			ItemId int64 `json:"itemId"`
			UserId int64 `json:"userId"`
			Count  int64 `json:"count"`
		}
		Response struct {
			CartItem cart.CartItem `json:"cartItem"`
		}
	)
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
			req Request
		)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}

		var err error
		validateInt64(err, "ItemId", req.ItemId)
		validateInt64(err, "UserId", req.UserId)
		validateInt64(err, "Count", req.Count)
		if err != nil {
			svc.Error(w, err, http.StatusBadRequest)
			return
		}

		var cartItem = &cart.CartItem{
			ItemId: req.ItemId,
			UserId: req.UserId,
			Count:  req.Count,
		}
		if err := cartItem.Insert(ctx, svc.DB); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		var resp = Response{
			CartItem: *cartItem,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
		}

	}
}

// GetCartHandler retrieves a user's cart from the service.
func (svc *Service) GetCartHandler() http.HandlerFunc {
	type Response struct {
		Cart cart.Cart `json:"cart"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
		if err != nil || userId == 0 {
			svc.Error(w, err, http.StatusBadRequest)
		}

		var cart = new(cart.Cart)
		if err := cart.Get(ctx, svc.DB, int64(userId)); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}

		var resp = Response{
			Cart: *cart,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
		}
	}
}

// PutCartItemHandler updates a cart item in the service.
func (svc *Service) PutCartItemHandler() http.HandlerFunc {
	type (
		Request struct {
			ItemId int64 `json:"itemId"`
			UserId int64 `json:"userId"`
			Count  int64 `json:"count"`
		}
		Response struct {
			CartItem cart.CartItem `json:"cartItem"`
		}
	)
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
			req Request
		)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}

		var err error
		validateInt64(err, "ItemId", req.ItemId)
		validateInt64(err, "UserId", req.UserId)
		validateInt64(err, "Count", req.Count)
		if err != nil {
			svc.Error(w, err, http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id == 0 {
			svc.Error(w, err, http.StatusBadRequest)
		}

		var cartItem = &cart.CartItem{
			Id:     int64(id),
			ItemId: req.ItemId,
			UserId: req.UserId,
			Count:  req.Count,
		}
		if err := cartItem.Update(ctx, svc.DB); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}
	}
}

// DeleteCartItemHandler deletes a cart item from the service.
func (svc *Service) DeleteCartItemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id == 0 {
			svc.Error(w, err, http.StatusBadRequest)
			return
		}

		if err := cart.DeleteCartItem(ctx, svc.DB, int64(id)); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}
	}
}
