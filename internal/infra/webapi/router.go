package webapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/domain/user"
)

type ctxKey string

const (
	ctxUserId   ctxKey = "userId"
	ctxUserRole ctxKey = "userRole"
)

type Router struct {
	UserService  *app.UserService
	EventService *app.EventService
	Handler      http.Handler
}

func NewRouter(userService *app.UserService, eventService *app.EventService) *Router {
	r := chi.NewRouter()
	router := &Router{
		UserService:  userService,
		EventService: eventService,
		Handler:      r,
	}
	r.Get("/health", router.HealthHandler)
	r.Post("/register", router.RegisterHandler)
	r.Post("/login", router.LoginHandler)
	r.NotFound(router.NotFoundHandler)

	r.Group(func(authR chi.Router) {
		authR.Use(AuthMiddleware(router.UserService))
		authR.Get("/events", router.ListEventsHandler)

		authR.Group(func(hostR chi.Router) {
			hostR.Use(AclMiddleware(user.HOST))
			hostR.Post("/events/add", router.CreateEventHandler)
		})
	})

	return router
}
