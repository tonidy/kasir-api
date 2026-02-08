package httputil

import (
	"encoding/json"
	stdErrors "errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/internal/model"
	errorsPkg "kasir-api/pkg/errors"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, map[string]string{"error": message})
}

// WriteAppError writes AppError response
func WriteAppError(w http.ResponseWriter, appErr *errorsPkg.AppError) error {
	return WriteJSON(w, appErr.Code, map[string]interface{}{
		"type":    appErr.Type,
		"message": appErr.Message,
	})
}

func ParseID(r *http.Request) (int, error) {
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return 0, stdErrors.New("invalid path")
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID: %w", err)
	}

	return id, nil
}

func ErrorStatus(err error) int {
	// Check if it's an AppError
	var appErr errorsPkg.AppError
	if errorsPkg.As(err, &appErr) {
		return appErr.Code
	}

	// Check if it's a model error
	if stdErrors.Is(err, model.ErrNotFound) {
		return http.StatusNotFound
	}
	if stdErrors.Is(err, model.ErrValidation) {
		return http.StatusBadRequest
	}
	if stdErrors.Is(err, model.ErrConflict) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

// HandleError handles errors consistently across the application
func HandleError(w http.ResponseWriter, err error) {
	var appErr errorsPkg.AppError
	if errorsPkg.As(err, &appErr) {
		WriteAppError(w, &appErr)
		return
	}

	// Convert standard errors to AppError
	errStr := err.Error()
	switch {
	case stdErrors.Is(err, model.ErrNotFound):
		WriteAppError(w, errorsPkg.NotFoundError(errStr))
	case stdErrors.Is(err, model.ErrValidation):
		WriteAppError(w, errorsPkg.ValidationError(errStr))
	case stdErrors.Is(err, model.ErrConflict):
		WriteAppError(w, errorsPkg.ConflictError(errStr))
	default:
		WriteAppError(w, errorsPkg.InternalError("internal server error"))
	}
}
