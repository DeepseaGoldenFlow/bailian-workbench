package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bailian-workbench/internal/client"
	"bailian-workbench/internal/repository"
)

var StoragePath = "/data/storage"

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

		var result map[string]any
		if json.Unmarshal(respBody, &result) == nil {
			if output, ok := result["output"].(map[string]any); ok {
				if status, ok := output["task_status"].(string); ok {
					if status == "SUCCEEDED" {
						// Download media to local storage
						localResult := downloadTaskMedia(output)
						if localResult != nil {
							result["output"] = localResult
							respBody, _ = json.Marshal(result)
						}
					}
					if status == "SUCCEEDED" || status == "FAILED" || status == "UNKNOWN" || status == "CANCELED" {
						repository.UpdateTaskResult(taskID, status, string(respBody))
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	}
}

func downloadTaskMedia(output map[string]any) map[string]any {
	local := make(map[string]any)
	for k, v := range output {
		local[k] = v
	}

	// Download video
	if videoURL, ok := output["video_url"].(string); ok && videoURL != "" {
		if localPath, err := downloadAndSave(videoURL, "video"); err == nil {
			local["video_url"] = "/api/storage/" + localPath
			local["_original_url"] = videoURL
		}
	}

	// Download result images
	if results, ok := output["results"].([]any); ok {
		var localResults []map[string]any
		for _, r := range results {
			if rm, ok := r.(map[string]any); ok {
				lm := make(map[string]any)
				for k, v := range rm {
					lm[k] = v
				}
				if url, ok := rm["url"].(string); ok && url != "" {
					if localPath, err := downloadAndSave(url, "image"); err == nil {
						lm["url"] = "/api/storage/" + localPath
						lm["_original_url"] = url
					}
				}
				localResults = append(localResults, lm)
			}
		}
		local["results"] = localResults
	}

	return local
}

func HandleSyncImageResult(ds *client.DashScope, respBody []byte) []byte {
	var result map[string]any
	if json.Unmarshal(respBody, &result) != nil {
		return respBody
	}

	output, ok := result["output"].(map[string]any)
	if !ok {
		return respBody
	}

	modified := false
	if choices, ok := output["choices"].([]any); ok {
		for _, c := range choices {
			if cm, ok := c.(map[string]any); ok {
				if msg, ok := cm["message"].(map[string]any); ok {
					if content, ok := msg["content"].([]any); ok {
						for _, item := range content {
							if im, ok := item.(map[string]any); ok {
								if imgURL, ok := im["image"].(string); ok && imgURL != "" {
									if localPath, err := downloadAndSave(imgURL, "image"); err == nil {
										im["image"] = "/api/storage/" + localPath
										modified = true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if modified {
		newBody, err := json.Marshal(result)
		if err == nil {
			return newBody
		}
	}
	return respBody
}

func downloadAndSave(url, mediaType string) (string, error) {
	if !strings.HasPrefix(url, "http") {
		return "", fmt.Errorf("not a remote URL")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed: %d", resp.StatusCode)
	}

	// Generate filename from URL hash + date
	hash := sha256.Sum256([]byte(url))
	hashStr := hex.EncodeToString(hash[:])[:12]
	dateDir := time.Now().Format("2006/01/02")
	dir := filepath.Join(StoragePath, mediaType, dateDir)
	os.MkdirAll(dir, 0755)

	ext := filepath.Ext(url)
	if idx := strings.Index(ext, "?"); idx > 0 {
		ext = ext[:idx]
	}
	if ext == "" || len(ext) > 5 {
		if mediaType == "video" {
			ext = ".mp4"
		} else {
			ext = ".png"
		}
	}

	filename := hashStr + ext
	fullPath := filepath.Join(dir, filename)

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}

	relPath := filepath.Join(mediaType, dateDir, filename)
	return relPath, nil
}


func HandleStorage(w http.ResponseWriter, r *http.Request) {
	relPath := strings.TrimPrefix(r.URL.Path, "/api/storage/")
	if relPath == "" || strings.Contains(relPath, "..") {
		writeError(w, 400, "invalid path")
		return
	}

	fullPath := filepath.Join(StoragePath, relPath)
	http.ServeFile(w, r, fullPath)
}
