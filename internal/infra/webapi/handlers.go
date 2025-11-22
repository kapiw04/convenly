package webapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kapiw04/convenly/internal/domain/event"
)

func (rt *Router) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var registerRequest RegisterRequest
	err := d.Decode(&registerRequest)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	if registerRequest.Name == "" || registerRequest.Email == "" || registerRequest.Password == "" {
		ErrorResponse(w, http.StatusBadRequest, "empty fields")
		return
	}

	err = rt.UserService.Register(registerRequest.Name, registerRequest.Email, registerRequest.Password)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (rt *Router) LoginHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var loginRequest LoginRequest
	err := d.Decode(&loginRequest)
	if err != nil {
		slog.Error("Failed to decode login request: %v", "err", err)
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	if loginRequest.Email == "" || loginRequest.Password == "" {
		slog.Error("Empty fields in login request")
		ErrorResponse(w, http.StatusBadRequest, "empty fields")
		return
	}
	sessionId, err := rt.UserService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		slog.Error("Login failed: %v", "err", err)
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	slog.Info("User logged in successfully: %s", "email", loginRequest.Email)
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Quoted:   false,
		Value:    sessionId,
		HttpOnly: true,
		Secure:   true,
	})
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := r.Header.Get("Authorization")
	if sessionId == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing session ID")
		return
	}
	err := rt.UserService.Logout(sessionId)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var addEventRequest CreateEventRequest
	err := d.Decode(&addEventRequest)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}

	date, err := time.Parse(time.RFC3339, addEventRequest.Date)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	userId := r.Context().Value(ctxUserId).(string)

	event := &event.Event{
		EventID:     uuid.New().String(),
		Name:        addEventRequest.Name,
		Description: addEventRequest.Description,
		Date:        date,
		Latitude:    addEventRequest.Latitude,
		Longitude:   addEventRequest.Longitude,
		Fee:         addEventRequest.Fee,
		OrganizerID: userId,
	}

	err = rt.EventService.CreateEvent(event)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (rt *Router) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := rt.EventService.GetAllEvents()
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponseSlice(w, http.StatusOK, events)
}

func (rt *Router) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user, err := rt.UserService.GetByUUID(r.Context().Value(ctxUserId).(string))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	user.PasswordHash = ""

	JSONResponse(w, http.StatusOK, user)
}

func (rt *Router) HealthHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusNotFound, "path not found")
}
