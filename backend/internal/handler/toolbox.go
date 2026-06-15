package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"bailian-workbench/internal/client"
	"bailian-workbench/internal/model"
	"bailian-workbench/internal/repository"
)

func HandleTranslate(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.TranslateRequest
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
		resp, err := ds.Do("/api/v1/services/aigc/text-generation/generation", body)
		if err != nil {
			writeError(w, 500, "Translation error: "+err.Error())
			return
		}
		defer resp.Body.Close()
		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		translated := extractChoiceContent(result)
		WriteJSON(w, 200, map[string]any{"success": true, "translated_text": translated})
		translatedJSON, _ := json.Marshal(translated)
		repository.SaveGeneration("translate", "qwen-mt-plus", req.Text, string(translatedJSON), "", "completed")
	}
}

func HandleOCR(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.OCRRequest
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
			"messages": []map[string]any{{
				"role": "user",
				"content": []map[string]any{
					{"type": "image_url", "image_url": map[string]string{"url": "data:image/png;base64," + imageData}},
					{"type": "text", "text": "Read all text in this image. Return ONLY the text, nothing else."},
				},
			}},
		}
		resp, err := ds.Do("/compatible-mode/v1/chat/completions", body)
		if err != nil {
			writeError(w, 500, "OCR error: "+err.Error())
			return
		}
		defer resp.Body.Close()
		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		text := extractChoiceContent(result)
		WriteJSON(w, 200, map[string]any{"success": true, "text": text})
		ocrJSON, _ := json.Marshal(text)
		repository.SaveGeneration("ocr", "qwen-vl-ocr", "OCR request", string(ocrJSON), "", "completed")
	}
}

func HandleDocument(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.DocumentRequest
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
		resp, err := ds.Do("/compatible-mode/v1/chat/completions", body)
		if err != nil {
			writeError(w, 500, "Document analysis error: "+err.Error())
			return
		}
		defer resp.Body.Close()
		var result map[string]any
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		content := extractChoiceContent(result)
		WriteJSON(w, 200, map[string]any{"success": true, "result": content, "task": req.Task})
		docJSON, _ := json.Marshal(content)
		repository.SaveGeneration("document", "qwen-plus", req.Task+": "+req.Text, string(docJSON), "", "completed")
	}
}

func extractChoiceContent(result map[string]any) string {
	if choices, ok := result["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if msg, ok := choice["message"].(map[string]any); ok {
				if content, ok := msg["content"].(string); ok {
					return content
				}
			}
		}
	}
	return ""
}
