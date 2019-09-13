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
	// r.Use(defaultMiddleware()...)
	r.Post("/cart/item", svc.AddCartItemHandler())
	r.Get("/cart/{userId}", svc.GetCartHandler())
	r.Put("/cart/item/{id}", svc.PutCartItemHandler())
	r.Delete("/cart/item/{id}", svc.DeleteCartItemHandler())
}

// PostCreatItemHandler creates a CartItem resource on the service.
func (svc *Service) AddCartItemHandler() http.HandlerFunc {
	type (
		Request struct {
			ItemId int `json:"itemId"`
			UserId int `json:"userId"`
			Count  int `json:"count"`
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

		v := new(validate)
		v.check("ItemId", intNotEmpty(req.ItemId))
		v.check("UserId", intNotEmpty(req.UserId))
		v.check("Count", intGreaterThan(req.Count, 0))
		if err := v.Err; err != nil {
			svc.Error(w, err, http.StatusBadRequest)
			return
		}

		id, err := cart.UserCartItemRelExists(ctx, svc.DB, req.UserId, req.ItemId)
		if err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}
		if id != 0 {
			cartItem, err := cart.FindCartItem(ctx, svc.DB, id)
			if err != nil {
				svc.Error(w, err, http.StatusInternalServerError)
				return
			}
			rel := cart.UserCartItemRel{
				ItemId: req.ItemId,
				UserId: req.UserId,
				Count:  cartItem.Count + req.Count,
			}
			if err := cart.UpdateUserCartItemRel(ctx, svc.DB, id, rel); err != nil {
				svc.Error(w, err, http.StatusInternalServerError)
				return
			}

		} else {
			rel := cart.UserCartItemRel{
				ItemId: req.ItemId,
				UserId: req.UserId,
				Count:  req.Count,
			}
			id, err = cart.CreateUserCartItemRel(ctx, svc.DB, rel)
			if err != nil {
				svc.Error(w, err, http.StatusInternalServerError)
				return
			}
		}

		cartItem, err := cart.FindCartItem(ctx, svc.DB, id)
		if err != nil {
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
		CartItems []cart.CartItem `json:"cartItems"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
		if err != nil || userId == 0 {
			svc.Error(w, err, http.StatusBadRequest)
		}

		cartItems, err := cart.CartItems(ctx, svc.DB, userId)
		if err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}

		var resp = Response{
			CartItems: cartItems,
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
			ItemId int `json:"itemId"`
			UserId int `json:"userId"`
			Count  int `json:"count"`
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

		v := new(validate)
		v.check("ItemId", intNotEmpty(req.ItemId))
		v.check("UserId", intNotEmpty(req.UserId))
		v.check("Count", intGreaterThan(req.Count, 0))
		if err := v.Err; err != nil {
			svc.Error(w, err, http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id == 0 {
			svc.Error(w, err, http.StatusBadRequest)
		}

		rel := cart.UserCartItemRel{
			ItemId: req.ItemId,
			UserId: req.UserId,
			Count:  req.Count,
		}
		if err := cart.UpdateUserCartItemRel(ctx, svc.DB, id, rel); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}
		cartItem, err := cart.FindCartItem(ctx, svc.DB, id)
		if err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}
		var resp = Response{
			CartItem: *cartItem,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
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

		if err := cart.DeleteCartItem(ctx, svc.DB, id); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}
	}
}
