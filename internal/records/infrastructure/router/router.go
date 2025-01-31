package router

import (
	"auth_records/pkg/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type APIHandler interface {
	Records(w http.ResponseWriter, r *http.Request)
}

func New(apiHandler APIHandler, authMiddleware *middleware.AuthMiddleware) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/records", authMiddleware.JwtMiddleware(apiHandler.Records))
	})

	return r
}
