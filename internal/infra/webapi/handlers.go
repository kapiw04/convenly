package webapi

import (
	"encoding/json"
	"net/http"
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

func (rt *Router) HealthHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (rt *Router) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusNotFound, "path not found")
}
