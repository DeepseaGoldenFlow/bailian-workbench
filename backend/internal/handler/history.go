package handler

import (
	"net/http"

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