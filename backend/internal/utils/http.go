package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/prodanov17/znk/internal/types"
)

var Validate = validator.New()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) {
	fmt.Println(status)

	errorResponse := types.ErrorResponse{StatusCode: status, Error: err.Error(), Path: r.RequestURI, Timestamp: time.Now()}

	WriteJSON(w, status, errorResponse)
}

func WriteValidationError(w http.ResponseWriter, r *http.Request, status int, err error) {
	fmt.Println(status)

	errorResponse := types.ErrorResponse{StatusCode: status, Error: err.Error(), Path: r.RequestURI, Timestamp: time.Now()}

	WriteJSON(w, status, errorResponse)
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close()

	// Read the entire body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %v", err)
	}

	// Check if the body is empty
	if len(bodyBytes) == 0 {
		return fmt.Errorf("empty request body")
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	return nil
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return strings.TrimPrefix(tokenAuth, "Bearer ")
	}

	if tokenQuery != "" {
		return strings.TrimPrefix(tokenQuery, "Bearer ")
	}

	return ""
}

func ValidatePayload(v any) error {
	if err := Validate.Struct(v); err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}

	return nil
}
