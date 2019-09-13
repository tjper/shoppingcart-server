package service

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func defaultMiddleware() []func(http.Handler) http.Handler {
	var cors = cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token"},
	})
	return []func(http.Handler) http.Handler{
		cors.Handler,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		contentTypeJsonMiddleware,
	}
}

func contentTypeJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
