package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// ============================================================
// Chat Handler (with streaming)
// ============================================================

func handleChat(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request")
			return
		}

		body := map[string]any{
			"model":    req.Model,
			"messages": req.Messages,
		}
		if req.MaxTokens > 0 {
			body["max_tokens"] = req.MaxTokens
		}
		if req.Temperature > 0 {
			body["temperature"] = req.Temperature
		}
		if req.Stream {
			body["stream"] = true
		}

		for _, msg := range req.Messages {
			if msg.Role == "user" {
				db.Exec("INSERT INTO chat_history (session_id, role, model, content) VALUES (?, 'user', ?, ?)",
					"default", req.Model, msg.Content)
			}
		}

		resp, err := c.do("/compatible-mode/v1/chat/completions", body)
		if err != nil {
			writeError(w, 500, "Chat API error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		if req.Stream {
			handleStreamResponse(w, resp)
		} else {
			handleNormalResponse(w, resp)
		}
	}
}

func handleStreamResponse(w http.ResponseWriter, resp *http.Response) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(200)

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	var fullContent strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(w, "%s\n", line)
		flusher.Flush()

		if strings.HasPrefix(line, "data: ") && line != "data: [DONE]" {
			data := strings.TrimPrefix(line, "data: ")
			var delta DeltaResponse
			if json.Unmarshal([]byte(data), &delta) == nil {
				for _, c := range delta.Choices {
					fullContent.WriteString(c.Delta.Content)
				}
			}
		}
	}
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()

	if fullContent.Len() > 0 {
		db.Exec("INSERT INTO chat_history (session_id, role, model, content) VALUES (?, 'assistant', ?, ?)",
			"default", "unknown", fullContent.String())
	}
}

func handleNormalResponse(w http.ResponseWriter, resp *http.Response) {
	var result ChatResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	for _, c := range result.Choices {
		db.Exec("INSERT INTO chat_history (session_id, role, model, content) VALUES (?, 'assistant', ?, ?)",
			"default", result.Model, c.Message.Content)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// ============================================================
// Image Generation Handler
// ============================================================

func handleImageGen(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Model          string  `json:"model"`
			Prompt         string  `json:"prompt"`
			NegativePrompt string  `json:"negative_prompt"`
			Size           string  `json:"size"`
			N              int     `json:"n"`
			Seed           *int    `json:"seed"`
			RefImg         string  `json:"ref_img"`
			RefStrength    float64 `json:"ref_strength"`
			Steps          int     `json:"steps"`
			PromptExtend   *bool   `json:"prompt_extend"`
		}
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

		// Route: qwen-image uses multimodal-generation (sync), wanx uses text2image (async)
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
			resp, apiErr = c.do(endpoint, body)
		} else {
			resp, apiErr = c.doAsync(endpoint, body)
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
		db.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, task_id, status) VALUES ('image', ?, ?, ?, ?, ?)",
			req.Model, req.Prompt, string(resultJSON), taskID, status)

		if taskID != "" {
			writeJSON(w, 200, map[string]any{
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

// ============================================================
// Video Generation Handler
// ============================================================

func handleVideoGen(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Model        string   `json:"model"`
			Prompt       string   `json:"prompt"`
			Duration     int      `json:"duration"`
			Resolution   string   `json:"resolution"`
			Ratio        string   `json:"ratio"`
			Seed         *int     `json:"seed"`
			Watermark    *bool    `json:"watermark"`
			FirstFrame   string   `json:"first_frame"`
			RefImages    []string `json:"ref_images"`
			VideoURL     string   `json:"video_url"`
			AudioSetting string   `json:"audio_setting"`
		}
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

		resp, err := c.doAsync("/api/v1/services/aigc/video-generation/video-synthesis", body)
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
		db.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, task_id, status) VALUES ('video', ?, ?, ?, ?, ?)",
			req.Model, req.Prompt, string(resultJSON), taskID, status)

		writeJSON(w, 200, map[string]any{
			"task_id": taskID,
			"status":  status,
			"message": "Video generation started. Poll GET /api/tasks/" + taskID + " for results.",
		})
	}
}

// ============================================================
// TTS Handler
// ============================================================

func handleTTS(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Input  string `json:"input"`
			Voice  string `json:"voice"`
			Format string `json:"format"`
		}
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

		resp, err := c.do("/api/v1/services/aigc/multimodal-generation/generation", body)
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

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Disposition", "attachment; filename=speech.mp3")
		io.Copy(w, audioResp.Body)

		db.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, status) VALUES ('audio', 'qwen-tts', ?, ?, 'completed')",
			req.Input, `["`+audioURL+`"]`)
	}
}

// ============================================================
// Toolbox Handlers
// ============================================================

func handleTranslate(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Text       string `json:"text"`
			SourceLang string `json:"source_lang"`
			TargetLang string `json:"target_lang"`
		}
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request")
			return
		}
		if req.SourceLang == "" {
			req.SourceLang = "auto"
		}
		if req.TargetLang == "" {
			req.TargetLang = "en"
		}

		body := map[string]any{
			"model": "qwen-mt-plus",
			"input": map[string]any{
				"messages": []map[string]string{{"role": "user", "content": req.Text}},
			},
			"parameters": map[string]any{
				"translation_options": map[string]string{
					"source_lang": req.SourceLang,
					"target_lang": req.TargetLang,
				},
			},
		}

		resp, err := c.do("/api/v1/services/aigc/text-generation/generation", body)
		if err != nil {
			writeError(w, 500, "Translation error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)

		translated := ""
		if output, ok := result["output"].(map[string]any); ok {
			if choices, ok := output["choices"].([]any); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]any); ok {
					if msg, ok := choice["message"].(map[string]any); ok {
						if content, ok := msg["content"].(string); ok {
							translated = content
						}
					}
				}
			}
		}

		writeJSON(w, 200, map[string]any{
			"success":         true,
			"translated_text": translated,
		})
		translatedJSON, _ := json.Marshal(translated)
		db.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, status) VALUES ('translate', 'qwen-mt-plus', ?, ?, 'completed')",
			req.Text, string(translatedJSON))
	}
}

func handleOCR(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ImageBase64 string `json:"image_base64"`
			ImageURL    string `json:"image_url"`
		}
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request")
			return
		}

		imageData := req.ImageBase64
		if imageData == "" && req.ImageURL != "" {
			imgResp, err := http.Get(req.ImageURL)
			if err != nil {
				writeError(w, 500, "Failed to download image: "+err.Error())
				return
			}
			defer imgResp.Body.Close()
			imgBytes, _ := io.ReadAll(imgResp.Body)
			imageData = base64Encode(imgBytes)
		}
		if imageData == "" {
			writeError(w, 400, "No image provided (image_base64 or image_url required)")
			return
		}

		body := map[string]any{
			"model": "qwen-vl-ocr-2025-11-20",
			"messages": []map[string]any{
				{
					"role": "user",
					"content": []map[string]any{
						{
							"type":      "image_url",
							"image_url": map[string]string{"url": "data:image/png;base64," + imageData},
						},
						{
							"type": "text",
							"text": "Read all text in this image. Return ONLY the text, nothing else.",
						},
					},
				},
			},
		}

		resp, err := c.do("/compatible-mode/v1/chat/completions", body)
		if err != nil {
			writeError(w, 500, "OCR error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)

		text := ""
		if choices, ok := result["choices"].([]any); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]any); ok {
				if msg, ok := choice["message"].(map[string]any); ok {
					if content, ok := msg["content"].(string); ok {
						text = content
					}
				}
			}
		}

		writeJSON(w, 200, map[string]any{
			"success": true,
			"text":    text,
		})
		ocrJSON, _ := json.Marshal(text)
		db.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, status) VALUES ('ocr', 'qwen-vl-ocr', ?, ?, 'completed')",
			"OCR request", string(ocrJSON))
	}
}

func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func handleCode(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Language string `json:"language"`
			Code     string `json:"code"`
		}
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request")
			return
		}
		if req.Language == "" {
			req.Language = "python"
		}

		var output string
		var err error
		success := true

		switch strings.ToLower(req.Language) {
		case "python":
			output, err = runCmd("python3", "-c", req.Code)
		case "bash", "sh":
			output, err = runCmd("sh", "-c", req.Code)
		default:
			output = "Unsupported language. Use python or bash."
		}

		if err != nil {
			output = err.Error() + "\n" + output
			success = false
		}

		writeJSON(w, 200, map[string]any{
			"success":  success,
			"output":   output,
			"language": req.Language,
		})
		codeJSON, _ := json.Marshal(output)
		db.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, status) VALUES ('code', ?, ?, ?, 'completed')",
			req.Language, req.Code, string(codeJSON))
	}
}

func runCmd(name string, args ...string) (string, error) {
	ctx := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	ctx.Stdout = &stdout
	ctx.Stderr = &stderr
	err := ctx.Run()
	result := stdout.String()
	if stderr.Len() > 0 {
		result += stderr.String()
	}
	return result, err
}

func handleDocument(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Task     string `json:"task"`
			Text     string `json:"text"`
			Question string `json:"question"`
		}
		if err := readJSON(r, &req); err != nil {
			writeError(w, 400, "Invalid request")
			return
		}
		if req.Task == "" {
			req.Task = "summarize"
		}

		var prompt string
		switch req.Task {
		case "summarize":
			prompt = "Summarize in Chinese:\n" + req.Text
		case "qa":
			prompt = "Answer based on doc. Q: " + req.Question + "\nDoc:\n" + req.Text
		case "extract":
			prompt = "Extract key info as JSON:\n" + req.Text
		case "translate":
			prompt = "Translate to English:\n" + req.Text
		default:
			prompt = "Analyze:\n" + req.Text
		}

		if len(prompt) > 12000 {
			prompt = prompt[:12000]
		}

		body := map[string]any{
			"model": "qwen-plus",
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
			"max_tokens": 4000,
		}

		resp, err := c.do("/compatible-mode/v1/chat/completions", body)
		if err != nil {
			writeError(w, 500, "Document analysis error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)

		content := ""
		if choices, ok := result["choices"].([]any); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]any); ok {
				if msg, ok := choice["message"].(map[string]any); ok {
					if c, ok := msg["content"].(string); ok {
						content = c
					}
				}
			}
		}

		writeJSON(w, 200, map[string]any{
			"success": true,
			"result":  content,
			"task":    req.Task,
		})
		docJSON, _ := json.Marshal(content)
		db.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, status) VALUES ('document', 'qwen-plus', ?, ?, 'completed')",
			req.Task+": "+req.Text, string(docJSON))
	}
}

// ============================================================
// ASR Handler
// ============================================================

func handleASR(c *DashScopeClient) http.HandlerFunc {
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

		resp, err := c.do("/api/v1/services/aigc/multimodal-generation/generation", body)
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
		db.Exec("INSERT INTO generation_history (type, model, prompt, result_urls, status) VALUES ('asr', 'paraformer-v2', ?, ?, 'completed')",
			"ASR request", string(asrJSON))

		writeJSON(w, resp.StatusCode, result)
	}
}

// ============================================================
// History Handlers
// ============================================================

func handleChatHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, session_id, role, model, content, created_at FROM chat_history ORDER BY created_at DESC LIMIT 200")
		if err != nil {
			writeError(w, 500, "Query error: "+err.Error())
			return
		}
		defer rows.Close()

		history := make([]map[string]any, 0)
		for rows.Next() {
			var id int64
			var sessionID, role, model, content string
			var createdAt time.Time
			rows.Scan(&id, &sessionID, &role, &model, &content, &createdAt)
			history = append(history, map[string]any{
				"id":         id,
				"session_id": sessionID,
				"role":       role,
				"model":      model,
				"content":    content,
				"created_at": createdAt,
			})
		}
		writeJSON(w, 200, history)
	}
}

func handleGenHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genType := r.URL.Query().Get("type")
		query := "SELECT id, type, model, prompt, result_urls, task_id, status, created_at FROM generation_history"
		args := []any{}
		if genType != "" {
			query += " WHERE type = ?"
			args = append(args, genType)
		}
		query += " ORDER BY created_at DESC LIMIT 100"

		rows, err := db.Query(query, args...)
		if err != nil {
			writeError(w, 500, "Query error: "+err.Error())
			return
		}
		defer rows.Close()

		history := make([]map[string]any, 0)
		for rows.Next() {
			var id int64
			var typ, model, prompt, resultJSON, taskID, status string
			var createdAt time.Time
			var promptNull, taskNull, resultNull sql.NullString
			rows.Scan(&id, &typ, &model, &promptNull, &resultNull, &taskNull, &status, &createdAt)
			prompt = promptNull.String
			resultJSON = resultNull.String
			taskID = taskNull.String
			history = append(history, map[string]any{
				"id":         id,
				"type":       typ,
				"model":      model,
				"prompt":     prompt,
				"result":     resultJSON,
				"task_id":    taskID,
				"status":     status,
				"created_at": createdAt,
			})
		}
		writeJSON(w, 200, history)
	}
}

func handleDeleteChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		db.Exec("DELETE FROM chat_history WHERE id = ?", id)
		writeJSON(w, 200, map[string]string{"status": "deleted"})
	}
}

func handleDeleteGen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		db.Exec("DELETE FROM generation_history WHERE id = ?", id)
		writeJSON(w, 200, map[string]string{"status": "deleted"})
	}
}

// ============================================================
// Generic Proxy Handler
// ============================================================

func handleProxy(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/proxy/")
		log.Printf("[Proxy] Forwarding to: %s", path)

		var body map[string]any
		readJSON(r, &body)

		resp, err := c.do("/"+path, body)
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

// ============================================================
// Task Polling Handler
// ============================================================

func handleTaskPoll(c *DashScopeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := r.PathValue("taskID")
		if taskID == "" {
			writeError(w, 400, "taskID required")
			return
		}

		resp, err := c.getTask(taskID)
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