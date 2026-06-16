package handler

import (
	"net/http"
	"strconv"

	"bailian-workbench/internal/repository"
)

func HandleChatHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		history, err := repository.GetChatHistory()
		if err != nil {
			writeError(w, 500, "Query error: "+err.Error())
			return
		}
		WriteJSON(w, 200, history)
	}
}

func HandleGenHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genType := r.URL.Query().Get("type")
		history, err := repository.GetGenerationHistory(genType)
		if err != nil {
			writeError(w, 500, "Query error: "+err.Error())
			return
		}
		WriteJSON(w, 200, history)
	}
}

func HandleUnifiedHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typ := r.URL.Query().Get("type")
		limitStr := r.URL.Query().Get("limit")
		limit := 100
		if limitStr != "" {
			if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 500 {
				limit = v
			}
		}
		entries, err := repository.GetUnifiedHistory(typ, limit)
		if err != nil {
			writeError(w, 500, "Query error: "+err.Error())
			return
		}
		if entries == nil {
			entries = make([]repository.UnifiedEntry, 0)
		}
		WriteJSON(w, 200, map[string]any{"entries": entries})
	}
}

func HandleDeleteChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		repository.DeleteChatMessage(id)
		WriteJSON(w, 200, map[string]string{"status": "deleted"})
	}
}

func HandleDeleteGen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		repository.DeleteGeneration(id)
		WriteJSON(w, 200, map[string]string{"status": "deleted"})
	}
}

func HandleDeleteUnified() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typ := r.URL.Query().Get("type")
		id := r.PathValue("id")
		if typ == "" || id == "" {
			writeError(w, 400, "type and id required")
			return
		}
		repository.DeleteUnifiedEntry(typ, id)
		WriteJSON(w, 200, map[string]string{"status": "deleted"})
	}
}
