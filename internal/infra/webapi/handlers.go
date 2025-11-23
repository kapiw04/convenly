package webapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
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
	user, err := rt.UserService.GetByEmail(loginRequest.Email)
	if err != nil {
		slog.Error("Failed to get user after login: %v", "err", err)
		ErrorResponse(w, http.StatusInternalServerError, "internal server error: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, user)
}

func (rt *Router) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := r.Context().Value(ctxSessionId).(string)
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
	uid := getUserId(r)

	e := &event.Event{
		EventID:     uuid.New().String(),
		Name:        addEventRequest.Name,
		Description: addEventRequest.Description,
		Date:        date,
		Latitude:    addEventRequest.Latitude,
		Longitude:   addEventRequest.Longitude,
		Fee:         addEventRequest.Fee,
		OrganizerID: uid,
	}

	err = rt.EventService.CreateEvent(e)
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
	user, err := rt.UserService.GetByUUID(getUserId(r))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	user.PasswordHash = ""

	JSONResponse(w, http.StatusOK, user)
}

func (rt *Router) BecomeHostHandler(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	err := rt.UserService.PromoteToHost(userId)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) EventDetailHandler(w http.ResponseWriter, r *http.Request) {
	eid := chi.URLParam(r, "id")
	uid := getUserId(r)
	e, err := rt.EventService.GetEventByID(eid)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	attendees, err := rt.EventService.GetAttendees(eid)

	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, struct {
		*event.Event
		AttendeesCount int  `json:"attendees_count"`
		UserRegistered bool `json:"user_registered"`
	}{
		Event:          e,
		AttendeesCount: len(attendees),
		UserRegistered: uid != "" && slices.Contains(attendees, uid),
	})
}

func (rt *Router) UnregisterFromEventHandler(w http.ResponseWriter, r *http.Request) {
	eventId := chi.URLParam(r, "id")
	userId := getUserId(r)

	err := rt.EventService.RemoveAttendance(userId, eventId)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) RegisterForEventHandler(w http.ResponseWriter, r *http.Request) {
	eventId := chi.URLParam(r, "id")
	userId := getUserId(r)

	err := rt.EventService.RegisterAttendance(userId, eventId)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) HealthHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusNotFound, "path not found")
}

func getUserId(r *http.Request) string {
	userId, ok := r.Context().Value(ctxUserId).(string)
	if !ok {
		return ""
	}
	return userId
}
