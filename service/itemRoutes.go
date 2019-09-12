package service

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tjper/shoppingcart-server/service/item"
)

// ItemRoutes defines the item resource Rest endpoints.
func (svc *Service) ItemRoutes(r chi.Router) {
	r.Get("/items", svc.GetItemsHandler())
}

// GetItemsHandler retrieves all item resources from the service.
func (svc *Service) GetItemsHandler() http.HandlerFunc {
	type Response struct {
		Items []item.Item `json:"items"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()

		items, err := item.Items(ctx, svc.DB)
		if err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}

		var resp = Response{
			Items: items,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			svc.Error(w, err, http.StatusInternalServerError)
			return
		}
	}
}
