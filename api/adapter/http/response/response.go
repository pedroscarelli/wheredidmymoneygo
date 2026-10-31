package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

// InternalError logs the real error and returns a generic message,
// so database or driver details never reach the client.
func InternalError(w http.ResponseWriter, err error) {
	log.Printf("internal error: %v", err)
	Error(w, http.StatusInternalServerError, "internal server error")
}
