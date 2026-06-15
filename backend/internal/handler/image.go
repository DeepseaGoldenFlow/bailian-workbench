package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"bailian-workbench/internal/catalog"
	"bailian-workbench/internal/client"
	"bailian-workbench/internal/model"
	"bailian-workbench/internal/repository"
)

func HandleImageGen(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.ImageGenRequest
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request: "+err.Error())
			return
		}

		m := ModelCatalog.Find(req.Model)
		if m == nil {
			writeError(w, 400, "Unknown model: "+req.Model)
			return
		}
		if m.Category != catalog.CatImage {
			writeError(w, 400, "Model "+req.Model+" is not an image model")
			return
		}

		if req.Size == "" {
			req.Size = "1024*1024"
		}
		if req.N <= 0 {
			req.N = 1
		}

		isQwen := strings.HasPrefix(req.Model, "qwen-image")
		params := map[string]any{
			"size": req.Size,
			"n":    req.N,
		}
		if req.Seed != nil {
			params["seed"] = *req.Seed
		}
		if req.Steps > 0 {
			params["steps"] = req.Steps
		}
		if req.PromptExtend != nil {
			params["prompt_extend"] = *req.PromptExtend
		}

		var body map[string]any
		if isQwen {
			content := []map[string]any{{"text": req.Prompt}}
			body = map[string]any{
				"model": req.Model,
				"input": map[string]any{
					"messages": []map[string]any{{"role": "user", "content": content}},
				},
				"parameters": params,
			}
		} else {
			input := map[string]any{"prompt": req.Prompt}
			if req.NegativePrompt != "" {
				input["negative_prompt"] = req.NegativePrompt
			}
			if req.RefImg != "" {
				input["ref_img"] = req.RefImg
				if req.RefStrength > 0 {
					input["ref_strength"] = req.RefStrength
				}
			}
			body = map[string]any{
				"model":      req.Model,
				"input":      input,
				"parameters": params,
			}
		}

		var resp *http.Response
		var apiErr error
		if m.Async {
			resp, apiErr = ds.DoAsync(m.Endpoint, body)
		} else {
			resp, apiErr = ds.Do(m.Endpoint, body)
		}
		if apiErr != nil {
			writeError(w, 500, "Image API error: "+apiErr.Error())
			return
		}
		defer resp.Body.Close()

		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)

		taskID := ""
		status := "completed"
		if output, ok := result["output"].(map[string]any); ok {
			if tid, ok := output["task_id"].(string); ok {
				taskID = tid
			}
			if s, ok := output["task_status"].(string); ok {
				status = s
			}
		}
		if code, _ := result["code"].(string); code != "" {
			status = "failed"
		}

		resultJSON, _ := json.Marshal(result)
		repository.SaveGeneration("image", req.Model, req.Prompt, string(resultJSON), taskID, status)

		if taskID != "" {
			WriteJSON(w, 200, map[string]any{
				"task_id": taskID,
				"status":  status,
				"message": "Image generation started. Poll GET /api/tasks/" + taskID + " for results.",
			})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			w.Write(respBody)
		}
	}
}
