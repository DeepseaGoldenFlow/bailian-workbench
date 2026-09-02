package handler

import (
	"encoding/json"
	"io"
	"net/http"

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
		apiModel := req.Model
		if req.ModelOverride != "" {
			apiModel = req.ModelOverride
		}

		input := copyMap(req.Input)
		params := copyMap(req.Parameters)
		prompt := req.Prompt
		if v, ok := input["prompt"].(string); ok && v != "" {
			prompt = v
		}

		// Compatibility with the original flat request format. Explicit input/parameters
		// maps always win, allowing every documented field to pass through unchanged.
		setIfMissing(input, "prompt", prompt)
		setIfMissing(input, "negative_prompt", req.NegativePrompt)
		setIfMissing(input, "ref_img", req.RefImg)
		if req.RefStrength > 0 {
			setIfMissing(input, "ref_strength", req.RefStrength)
		}
		setIfMissing(params, "size", req.Size)
		if req.N > 0 {
			setIfMissing(params, "n", req.N)
		}
		if req.Seed != nil {
			setIfMissing(params, "seed", *req.Seed)
		}
		if req.Steps > 0 {
			setIfMissing(params, "steps", req.Steps)
		}
		if req.PromptExtend != nil {
			setIfMissing(params, "prompt_extend", *req.PromptExtend)
		}

		if m.Payload == "messages" {
			if _, supplied := input["messages"]; !supplied {
				content := make([]map[string]any, 0, len(req.Images)+1)
				for _, imageURL := range req.Images {
					if imageURL != "" {
						content = append(content, map[string]any{"image": imageURL})
					}
				}
				if prompt != "" {
					content = append(content, map[string]any{"text": prompt})
				}
				input = map[string]any{"messages": []map[string]any{{"role": "user", "content": content}}}
			}
		}
		body := map[string]any{"model": apiModel, "input": input, "parameters": params}

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

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			w.Write(respBody)
			return
		}

		// For synchronous message-based images, download to local storage immediately.
		savedBody := respBody
		if !m.Async && m.Payload == "messages" {
			savedBody = HandleSyncImageResult(ds, respBody)
		}

		repository.SaveGeneration("image", apiModel, prompt, string(savedBody), taskID, status)

		if taskID != "" {
			WriteJSON(w, 200, map[string]any{"task_id": taskID, "status": status, "message": "Poll GET /api/tasks/" + taskID})
		} else {
			// Send the modified response with local URLs
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			w.Write(savedBody)
		}
	}
}
