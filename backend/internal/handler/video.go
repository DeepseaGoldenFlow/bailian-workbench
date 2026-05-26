package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"bailian-workbench/internal/client"
	"bailian-workbench/internal/model"
	"bailian-workbench/internal/repository"
)

func HandleVideoGen(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.VideoGenRequest
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request")
			return
		}
		if req.Model == "" {
			req.Model = "happyhorse-1.0-t2v"
		}
		if req.Duration < 3 || req.Duration > 15 {
			req.Duration = 5
		}
		if req.Resolution == "" {
			req.Resolution = "1080P"
		}
		if req.Ratio == "" {
			req.Ratio = "16:9"
		}
		if req.AudioSetting == "" {
			req.AudioSetting = "auto"
		}
		if req.Model == "happyhorse-1.0-video-edit" && len(req.RefImages) > 5 {
			writeError(w, 400, "video-edit supports at most 5 reference images")
			return
		}
		if req.Model == "happyhorse-1.0-r2v" && len(req.RefImages) > 9 {
			writeError(w, 400, "r2v supports at most 9 reference images")
			return
		}

		input := map[string]any{"prompt": req.Prompt}

		switch req.Model {
		case "happyhorse-1.0-i2v":
			if req.FirstFrame != "" {
				input["media"] = []map[string]any{{"type": "first_frame", "url": req.FirstFrame}}
			}
		case "happyhorse-1.0-r2v":
			if len(req.RefImages) > 0 {
				var media []map[string]any
				for _, url := range req.RefImages {
					media = append(media, map[string]any{"type": "reference_image", "url": url})
				}
				input["media"] = media
			}
		case "happyhorse-1.0-video-edit":
			var media []map[string]any
			if req.VideoURL != "" {
				media = append(media, map[string]any{"type": "video", "url": req.VideoURL})
			}
			for _, url := range req.RefImages {
				media = append(media, map[string]any{"type": "reference_image", "url": url})
			}
			if len(media) > 0 {
				input["media"] = media
			}
		}

		params := map[string]any{"resolution": req.Resolution}
		if req.Model != "happyhorse-1.0-video-edit" {
			params["duration"] = req.Duration
			params["ratio"] = req.Ratio
		} else {
			params["audio_setting"] = req.AudioSetting
		}
		if req.Seed != nil {
			params["seed"] = *req.Seed
		}
		if req.Watermark != nil {
			params["watermark"] = *req.Watermark
		} else {
			params["watermark"] = false
		}

		body := map[string]any{
			"model":      req.Model,
			"input":      input,
			"parameters": params,
		}

		resp, err := ds.DoAsync("/api/v1/services/aigc/video-generation/video-synthesis", body)
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