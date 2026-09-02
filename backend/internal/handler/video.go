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

func HandleVideoGen(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.VideoGenRequest
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request: "+err.Error())
			return
		}

		m := ModelCatalog.Find(req.Model)
		if m == nil {
			writeError(w, 400, "Unknown model: "+req.Model)
			return
		}
		if m.Category != catalog.CatVideo {
			writeError(w, 400, "Model "+req.Model+" is not a video model")
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
		setIfMissing(input, "prompt", prompt)
		if len(req.Media) > 0 {
			setIfMissing(input, "media", req.Media)
		}

		// Keep old clients working by translating their flat media fields. New clients
		// send input.media directly and can use every media type supported by the model.
		_, hasMedia := input["media"]
		if !hasMedia {
			switch req.Model {
			case "happyhorse-1.0-i2v":
				if req.FirstFrame != "" {
					input["media"] = []map[string]any{
						{"type": "first_frame", "url": req.FirstFrame},
					}
				}
			case "happyhorse-1.0-r2v":
				if len(req.RefImages) > 0 {
					var media []map[string]any
					for _, url := range req.RefImages {
						url = strings.TrimSpace(url)
						if url != "" {
							media = append(media, map[string]any{"type": "reference_image", "url": url})
						}
					}
					if len(media) > 0 {
						input["media"] = media
					}
				}
			case "happyhorse-1.0-video-edit":
				var media []map[string]any
				if req.VideoURL != "" {
					media = append(media, map[string]any{"type": "video", "url": req.VideoURL})
				}
				for _, url := range req.RefImages {
					url = strings.TrimSpace(url)
					if url != "" {
						media = append(media, map[string]any{"type": "reference_image", "url": url})
					}
				}
				if len(media) > 0 {
					input["media"] = media
				}
			}
		}

		setIfMissing(params, "resolution", req.Resolution)
		if req.Duration != 0 {
			setIfMissing(params, "duration", req.Duration)
		}
		setIfMissing(params, "ratio", req.Ratio)
		if req.Seed != nil {
			setIfMissing(params, "seed", *req.Seed)
		}
		if req.Watermark != nil {
			setIfMissing(params, "watermark", *req.Watermark)
		}
		if req.AudioSetting != "" {
			setIfMissing(params, "audio_setting", req.AudioSetting)
		}

		body := map[string]any{
			"model":      apiModel,
			"input":      input,
			"parameters": params,
		}

		resp, err := ds.DoAsync(m.Endpoint, body)
		if err != nil {
			writeError(w, 500, "Video API error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			w.Write(respBody)
			return
		}

		taskID := ""
		status := "pending"
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
		repository.SaveGeneration("video", apiModel, prompt, string(resultJSON), taskID, status)

		WriteJSON(w, 200, map[string]any{
			"task_id": taskID,
			"status":  status,
			"message": "Video generation started. Poll GET /api/tasks/" + taskID + " for results.",
		})
	}
}
