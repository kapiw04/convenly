package webapi

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateEventRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Date        string  `json:"date"` // ISO 8601 format
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Fee         float32 `json:"fee"`
}
