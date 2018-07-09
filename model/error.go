package model

// ErrorResponse is a generic response for sending a error.
type ErrorResponse struct {
	Err string `json:"err,omitempty"`
}