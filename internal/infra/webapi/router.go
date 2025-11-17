package webapi

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
	router := &Router{
		UserService: userService,
		Handler:     r,
	}
	r.Get("/health", router.HealthHandler)
	r.Post("/register", router.RegisterHandler)
	r.Post("/login", router.LoginHandler)
	r.NotFound(router.NotFoundHandler)

	return router
}
