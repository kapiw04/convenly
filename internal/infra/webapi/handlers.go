package webapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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
	sessionID, err := rt.UserService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		slog.Error("Login failed: %v", "err", err)
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	slog.Info("User logged in successfully: %s", "email", loginRequest.Email)
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Quoted:   false,
		Value:    sessionID,
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
	sessionID := r.Context().Value(ctxSessionID).(string)
	if sessionID == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing session ID")
		return
	}
	err := rt.UserService.Logout(sessionID)
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
	uid := getUserID(r)

	e := &event.Event{
		EventID:     uuid.New().String(),
		Name:        addEventRequest.Name,
		Description: addEventRequest.Description,
		Date:        date,
		Latitude:    addEventRequest.Latitude,
		Longitude:   addEventRequest.Longitude,
		Fee:         addEventRequest.Fee,
		OrganizerID: uid,
		Tags:        addEventRequest.Tags,
	}

	err = rt.EventService.CreateEvent(e)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (rt *Router) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	filter := &event.EventFilter{}

	if page := r.URL.Query().Get("page"); page != "" {
		p, err := strconv.Atoi(page)
		if err != nil || p < 1 {
			ErrorResponse(w, http.StatusBadRequest, "invalid page format")
			return
		}
		pageSize := 12
		if ps := r.URL.Query().Get("page_size"); ps != "" {
			pageSize, err = strconv.Atoi(ps)
			if err != nil || pageSize < 1 || pageSize > 100 {
				ErrorResponse(w, http.StatusBadRequest, "invalid page_size format (1-100)")
				return
			}
		}
		filter.Pagination = &event.Pagination{Page: p, PageSize: pageSize}
	}

	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		t, err := time.Parse(time.RFC3339, dateFrom)
		if err != nil {
			t, err = time.Parse("2006-01-02", dateFrom)
			if err != nil {
				ErrorResponse(w, http.StatusBadRequest, "invalid date_from format, use RFC3339 or YYYY-MM-DD")
				return
			}
		}
		filter.DateFrom = &t
	}

	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		t, err := time.Parse(time.RFC3339, dateTo)
		if err != nil {
			t, err = time.Parse("2006-01-02", dateTo)
			if err != nil {
				ErrorResponse(w, http.StatusBadRequest, "invalid date_to format, use RFC3339 or YYYY-MM-DD")
				return
			}
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		filter.DateTo = &t
	}

	if minFee := r.URL.Query().Get("min_fee"); minFee != "" {
		f, err := strconv.ParseFloat(minFee, 32)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "invalid min_fee format")
			return
		}
		fee := float32(f)
		filter.MinFee = &fee
	}

	if maxFee := r.URL.Query().Get("max_fee"); maxFee != "" {
		f, err := strconv.ParseFloat(maxFee, 32)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "invalid max_fee format")
			return
		}
		fee := float32(f)
		filter.MaxFee = &fee
	}

	if tags := r.URL.Query().Get("tags"); tags != "" {
		filter.Tags = strings.Split(tags, ",")
	}

	hasFilters := filter.DateFrom != nil || filter.DateTo != nil ||
		filter.MinFee != nil || filter.MaxFee != nil || len(filter.Tags) > 0 ||
		filter.Pagination != nil

	var events []*event.Event
	var err error

	if hasFilters {
		events, err = rt.EventService.GetEventsWithFilters(filter)
	} else {
		events, err = rt.EventService.GetAllEvents()
	}

	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponseSlice(w, http.StatusOK, events)
}

func (rt *Router) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user, err := rt.UserService.GetByUUID(getUserID(r))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	user.PasswordHash = ""

	JSONResponse(w, http.StatusOK, user)
}

func (rt *Router) BecomeHostHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	err := rt.UserService.PromoteToHost(userID)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) EventDetailHandler(w http.ResponseWriter, r *http.Request) {
	eid := chi.URLParam(r, "id")
	uid := getUserID(r)
	e, err := rt.EventService.GetEventByID(eid)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	attendeesCount, err := rt.EventService.GetAttendeesCount(eid)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "internal server error: "+err.Error())
		return
	}
	isUserAttending := rt.EventService.IsUserAttending(uid, eid)

	JSONResponse(w, http.StatusOK, struct {
		*event.Event
		AttendeesCount int  `json:"attendees_count"`
		UserRegistered bool `json:"user_registered"`
	}{
		Event:          e,
		AttendeesCount: attendeesCount,
		UserRegistered: uid != "" && isUserAttending,
	})
}

func (rt *Router) UnregisterFromEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	userID := getUserID(r)

	err := rt.EventService.RemoveAttendance(userID, eventID)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) RegisterForEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	userID := getUserID(r)

	err := rt.EventService.RegisterAttendance(userID, eventID)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad request: "+err.Error())
		return
	}
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) MyEventsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == "" {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	hosting, err := rt.EventService.GetHostingEvents(userID, nil)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to get hosting events: "+err.Error())
		return
	}

	attending, err := rt.EventService.GetAttendingEvents(userID, nil)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to get attending events: "+err.Error())
		return
	}

	if hosting == nil {
		hosting = []*event.Event{}
	}
	if attending == nil {
		attending = []*event.Event{}
	}

	JSONResponse(w, http.StatusOK, struct {
		Hosting   []*event.Event `json:"hosting"`
		Attending []*event.Event `json:"attending"`
	}{
		Hosting:   hosting,
		Attending: attending,
	})
}

func (rt *Router) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	userID := getUserID(r)

	eventData, err := rt.EventService.GetEventByID(eventID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "event not found")
		return
	}

	if eventData.OrganizerID != userID {
		ErrorResponse(w, http.StatusForbidden, "you can only delete your own events")
		return
	}

	err = rt.EventService.DeleteEvent(eventID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to delete event: "+err.Error())
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

func getUserID(r *http.Request) string {
	userID, ok := r.Context().Value(ctxUserID).(string)
	if !ok {
		return ""
	}
	return userID
}
