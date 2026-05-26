package handler

import (
	"io"
	"net/http"

	"bailian-workbench/internal/client"
)

func HandleTaskPoll(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := r.PathValue("taskID")
		if taskID == "" {
			writeError(w, 400, "taskID required")
			return
		}

		resp, err := ds.GetTask(taskID)
		if err != nil {
			writeError(w, 500, "Task query error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	}
}