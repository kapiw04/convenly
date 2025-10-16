package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kapiw04/convenly/internal/app"
)

type Router struct {
	UserService *app.UserService
	Handler     http.Handler
}

func NewRouter(userService *app.UserService) *Router {
	r := chi.NewRouter()
	r.Get("/health", HealthHandler)
	r.NotFound(NotFoundHandler)

	return &Router{
		UserService: userService,
		Handler:     r,
	}
}
