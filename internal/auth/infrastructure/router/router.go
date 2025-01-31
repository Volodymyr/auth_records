package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type APIHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func New(apiHandler APIHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", apiHandler.Login)
	})

	return r
}
