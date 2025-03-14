package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// respondWithJSON writes a JSON response to the client
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error marshaling JSON: %v", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError writes an error response to the client
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ApiResponse{
		Success: false,
		Message: message,
	})
}

// decodeJSONBody decodes a JSON request body
func decodeJSONBody(r *http.Request, dst interface{}) error {
	// Verify content type
	contentType := r.Header.Get("Content-Type")
	if contentType != "" {
		if !strings.Contains(contentType, "application/json") {
			return errors.New("content-type header is not application/json")
		}
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// If the body is empty, return early
	if len(body) == 0 {
		return errors.New("request body is empty")
	}

	// Decode the JSON
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("request body contains invalid JSON: %v", err)
	}

	return nil
}
