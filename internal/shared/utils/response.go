package utils

import (
	"encoding/json"
	"net/http"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/dto"
)

func base(w http.ResponseWriter, code int, status string, data, errs any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := dto.WebResponse{
		Status: status,
		Data: data,
		Errs: errs,
	}

	json.NewEncoder(w).Encode(res)
}

// --- Success Helpers ---
func Created(w http.ResponseWriter, data any) {
	base(w, http.StatusCreated, "CREATED", data, nil)
}

// --- Client Errors Helpers ---
func BadRequest(w http.ResponseWriter, message string) {
	base(w, http.StatusBadRequest, "BAD_REQUEST", nil, message)
}