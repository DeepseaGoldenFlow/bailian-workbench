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

		if req.Resolution == "" {
			req.Resolution = "1080P"
		}
		if req.Duration < 3 || req.Duration > 15 {
			req.Duration = 5
		}
		if req.Ratio == "" {
			req.Ratio = "16:9"
		}
		if req.AudioSetting == "" {
			req.AudioSetting = "auto"
		}
		if req.Watermark == nil {
			t := true
			req.Watermark = &t
		}

		// Validate
		switch req.Model {
		case "happyhorse-1.0-video-edit":
			if req.VideoURL == "" {
				writeError(w, 400, "video_url is required for video-edit")
				return
			}
			if len(req.RefImages) > 5 {
				writeError(w, 400, "video-edit supports at most 5 reference images")
				return
			}
		case "happyhorse-1.0-r2v":
			if len(req.RefImages) == 0 {
				writeError(w, 400, "at least 1 reference image required for r2v")
				return
			}
			if len(req.RefImages) > 9 {
				writeError(w, 400, "r2v supports at most 9 reference images")
				return
			}
		case "happyhorse-1.0-i2v":
			if req.FirstFrame == "" {
				writeError(w, 400, "first_frame is required for i2v")
				return
			}
		}

		// Build input
		input := map[string]any{"prompt": req.Prompt}
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

		// Build params
		params := map[string]any{"resolution": req.Resolution}
		if req.Model != "happyhorse-1.0-video-edit" {
			params["duration"] = req.Duration
			params["ratio"] = req.Ratio
		}
		if req.Seed != nil {
			params["seed"] = *req.Seed
		}
		if req.Watermark != nil {
			params["watermark"] = *req.Watermark
		}
		if req.AudioSetting != "" {
			params["audio_setting"] = req.AudioSetting
		}

		body := map[string]any{
			"model":      req.Model,
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
		repository.SaveGeneration("video", req.Model, req.Prompt, string(resultJSON), taskID, status)

		WriteJSON(w, 200, map[string]any{
			"task_id": taskID,
			"status":  status,
			"message": "Video generation started. Poll GET /api/tasks/" + taskID + " for results.",
		})
	}
}
