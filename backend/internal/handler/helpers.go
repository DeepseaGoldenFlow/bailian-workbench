package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"bailian-workbench/internal/model"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, model.APIError{Code: status, Message: msg})
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func copyMap(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func setIfMissing(dst map[string]any, key string, value any) {
	if _, exists := dst[key]; exists || value == nil {
		return
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return
		}
	}
	dst[key] = value
}
