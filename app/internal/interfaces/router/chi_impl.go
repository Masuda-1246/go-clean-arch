package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewChiImpl() http.Handler {
	r := chi.NewRouter()
	setMiddlewareImplChi(r)
	return r
}

func setMiddlewareImplChi(r *chi.Mux) {
	r.Use(middleware.DefaultLogger)
}
