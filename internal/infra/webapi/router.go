package webapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	}))

	r.Get("/api/health", router.HealthHandler)
	r.Post("/api/register", router.RegisterHandler)
	r.Post("/api/login", router.LoginHandler)
	r.NotFound(router.NotFoundHandler)

	r.Group(func(authR chi.Router) {
		authR.Use(AuthMiddleware(router.UserService))
		authR.Get("/api/events", router.ListEventsHandler)
		authR.Get("/api/me", router.GetUserInfoHandler)

		authR.Group(func(hostR chi.Router) {
			hostR.Use(AclMiddleware(user.HOST))
			hostR.Post("/api/events/add", router.CreateEventHandler)
		})
	})

	return router
}
