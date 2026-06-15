package handler

import (
	"net/http"

	"bailian-workbench/internal/catalog"
)

var ModelCatalog *catalog.Catalog

func InitCatalog() {
	ModelCatalog = catalog.BuildCatalog()
}

func HandleModels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cat := r.URL.Query().Get("category")
		if cat != "" {
			WriteJSON(w, 200, map[string]any{"models": ModelCatalog.ByCategory(cat)})
			return
		}
		WriteJSON(w, 200, map[string]any{"models": ModelCatalog.Models})
	}
}

func HandleModelByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		m := ModelCatalog.Find(id)
		if m == nil {
			writeError(w, 404, "model not found: "+id)
			return
		}
		WriteJSON(w, 200, m)
	}
}
