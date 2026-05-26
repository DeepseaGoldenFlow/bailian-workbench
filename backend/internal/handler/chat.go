package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"bailian-workbench/internal/client"
	"bailian-workbench/internal/model"
	"bailian-workbench/internal/repository"
)

func HandleChat(ds *client.DashScope) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.ChatRequest
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
				repository.SaveChatMessage("default", "user", req.Model, msg.Content)
			}
		}

		resp, err := ds.Do("/compatible-mode/v1/chat/completions", body)
		if err != nil {
			writeError(w, 500, "Chat API error: "+err.Error())
			return
		}
		defer resp.Body.Close()

		if req.Stream {
			handleStreamResponse(w, resp, req.Model)
		} else {
			handleNormalResponse(w, resp)
		}
	}
}

func handleStreamResponse(w http.ResponseWriter, resp *http.Response, model string) {
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
			var delta model.DeltaResponse
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
		repository.SaveChatMessage("default", "assistant", model, fullContent.String())
	}
}

func handleNormalResponse(w http.ResponseWriter, resp *http.Response) {
	var result model.ChatResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	for _, c := range result.Choices {
		repository.SaveChatMessage("default", "assistant", result.Model, c.Message.Content)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}