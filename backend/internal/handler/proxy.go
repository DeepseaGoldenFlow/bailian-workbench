package handler

import (
	"io"
	"log"
	"net/http"
	"strings"

	"bailian-workbench/internal/client"
)

func HandleProxy(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/proxy/")
		log.Printf("[Proxy] Forwarding to: %s", path)

		var body map[string]any
		readJSON(r, &body)

		resp, err := ds.Do("/"+path, body)
		if err != nil {
			writeError(w, 500, "Proxy error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	}
}