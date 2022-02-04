package handlers

import "github.com/alexferl/air/storage"

type (
	// Handler represents the structure of our resource
	Handler struct {
		Storage storage.Storage
	}
)

// ErrorResponse holds an error message
type ErrorResponse struct {
	Message string `json:"error"`
}

type Response struct {
	Message string `json:"message"`
}
