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
	ctxUserID    ctxKey = "userID"
	ctxUserRole  ctxKey = "userRole"
	ctxSessionID ctxKey = "sessionID"
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
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	}))

	r.Get("/api/health", router.HealthHandler)
	r.Post("/api/register", router.RegisterUserHandler)
	r.Post("/api/login", router.LoginHandler)
	r.Get("/api/events", router.ListEventsHandler)
	r.NotFound(router.NotFoundHandler)

	r.Group(func(authR chi.Router) {
		authR.Use(AuthMiddleware(router.UserService))
		authR.Get("/api/me", router.GetUserInfoHandler)
		authR.Post("/api/logout", router.LogoutHandler)
		authR.Post("/api/become-host", router.BecomeHostHandler)
		authR.Get("/api/my-events", router.MyEventsHandler)
		authR.Get("/api/events/{id}", router.EventDetailHandler)
		authR.Post("/api/events/{id}/register", router.RegisterForEventHandler)
		authR.Delete("/api/events/{id}/unregister", router.UnregisterFromEventHandler)

		authR.Group(func(hostR chi.Router) {
			hostR.Use(AclMiddleware(user.HOST))
			hostR.Post("/api/events/add", router.CreateEventHandler)
			hostR.Delete("/api/events/{id}", router.DeleteEventHandler)
		})
	})

	return router
}
