package http

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {

// }
