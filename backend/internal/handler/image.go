package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"bailian-workbench/internal/client"
	"bailian-workbench/internal/model"
	"bailian-workbench/internal/repository"
)

func HandleImageGen(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.ImageGenRequest
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request")
			return
		}
		if req.Model == "" {
			req.Model = "qwen-image-2.0-pro"
		}
		if req.Size == "" {
			req.Size = "1024*1024"
		}
		if req.N <= 0 {
			req.N = 1
		}

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

		isQwen := strings.HasPrefix(req.Model, "qwen-image")
		var body map[string]any
		var endpoint string
		if isQwen {
			endpoint = "/api/v1/services/aigc/multimodal-generation/generation"
			content := []map[string]any{{"text": req.Prompt}}
			body = map[string]any{
				"model": req.Model,
				"input": map[string]any{
					"messages": []map[string]any{{"role": "user", "content": content}},
				},
				"parameters": params,
			}
		} else {
			endpoint = "/api/v1/services/aigc/text2image/image-synthesis"
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
		if isQwen {
			resp, apiErr = ds.Do(endpoint, body)
		} else {
			resp, apiErr = ds.DoAsync(endpoint, body)
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