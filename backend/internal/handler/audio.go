package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"bailian-workbench/internal/client"
	"bailian-workbench/internal/model"
	"bailian-workbench/internal/repository"
)

func HandleTTS(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.TTSRequest
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request")
			return
		}
		if req.Voice == "" {
			req.Voice = "Cherry"
		}
		if req.Format == "" {
			req.Format = "mp3"
		}

		body := map[string]any{
			"model": "qwen-tts",
			"input": map[string]any{"text": req.Input},
			"parameters": map[string]any{
				"voice":  req.Voice,
				"format": req.Format,
			},
		}

		resp, err := ds.Do("/api/v1/services/aigc/multimodal-generation/generation", body)
		if err != nil {
			writeError(w, 500, "TTS API error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			writeError(w, resp.StatusCode, string(respBody))
			return
		}

		var result map[string]any
		json.Unmarshal(respBody, &result)

		audioURL := ""
		if output, ok := result["output"].(map[string]any); ok {
			if audio, ok := output["audio"].(map[string]any); ok {
				if url, ok := audio["url"].(string); ok {
					audioURL = url
				}
			}
		}

		if audioURL == "" {
			writeError(w, 500, "No audio URL in TTS response")
			return
		}

		audioResp, err := http.Get(audioURL)
		if err != nil {
			writeError(w, 500, "Audio download error: "+err.Error())
			return
		}
		defer audioResp.Body.Close()

		contentType := audioResp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		extension := req.Format
		switch strings.ToLower(strings.Split(contentType, ";")[0]) {
		case "audio/wav", "audio/x-wav", "audio/wave":
			extension = "wav"
		case "audio/mpeg", "audio/mp3":
			extension = "mp3"
		case "audio/ogg", "audio/opus":
			extension = "opus"
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=speech.%s", extension))
		io.Copy(w, audioResp.Body)

		repository.SaveGeneration("audio", "qwen-tts", req.Input, `["`+audioURL+`"]`, "", "completed")
	}
}

func HandleASR(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(50 << 20)
		file, _, err := r.FormFile("file")
		if err != nil {
			writeError(w, 400, "No audio file provided (form field: file)")
			return
		}
		defer file.Close()

		fileBytes, _ := io.ReadAll(file)
		fileData := base64Encode(fileBytes)

		body := map[string]any{
			"model": "paraformer-v2",
			"input": map[string]any{
				"messages": []map[string]any{
					{
						"role": "user",
						"content": []map[string]any{
							{"audio": "data:audio/wav;base64," + fileData},
						},
					},
				},
			},
		}

		resp, err := ds.Do("/api/v1/services/aigc/multimodal-generation/generation", body)
		if err != nil {
			writeError(w, 500, "ASR error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)

		text := ""
		if output, ok := result["output"].(map[string]any); ok {
			if choices, ok := output["choices"].([]any); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]any); ok {
					if msg, ok := choice["message"].(map[string]any); ok {
						if c, ok := msg["content"].(string); ok {
							text = c
						}
					}
				}
			}
		}
		asrJSON, _ := json.Marshal(text)
		repository.SaveGeneration("asr", "paraformer-v2", "ASR request", string(asrJSON), "", "completed")

		WriteJSON(w, resp.StatusCode, result)
	}
}
