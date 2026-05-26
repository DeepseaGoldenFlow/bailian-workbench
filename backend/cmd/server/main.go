package main

import (
	"log"
	"net/http"

	"bailian-workbench/internal/client"
	"bailian-workbench/internal/config"
	"bailian-workbench/internal/handler"
	"bailian-workbench/internal/middleware"
	"bailian-workbench/internal/repository"
)

func main() {
	cfg := config.Load()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := repository.InitDB(cfg.MySQLDSN); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}
	defer repository.DB.Close()

	ds := client.NewDashScope(cfg.ApiKey)

	mux := http.NewServeMux()

	// Chat
	mux.HandleFunc("POST /api/chat/completions", handler.HandleChat(ds))
	// Images
	mux.HandleFunc("POST /api/images/generate", handler.HandleImageGen(ds))
	// Videos
	mux.HandleFunc("POST /api/videos/generate", handler.HandleVideoGen(ds))
	// Audio
	mux.HandleFunc("POST /api/audio/speech", handler.HandleTTS(ds))
	mux.HandleFunc("POST /api/audio/transcribe", handler.HandleASR(ds))
	// Toolbox
	mux.HandleFunc("POST /api/toolbox/translate", handler.HandleTranslate(ds))
	mux.HandleFunc("POST /api/toolbox/ocr", handler.HandleOCR(ds))
	mux.HandleFunc("POST /api/toolbox/code", handler.HandleCode(ds))
	mux.HandleFunc("POST /api/toolbox/document", handler.HandleDocument(ds))
	// History
	mux.HandleFunc("GET /api/history/chat", handler.HandleChatHistory())
	mux.HandleFunc("GET /api/history/generations", handler.HandleGenHistory())
	mux.HandleFunc("DELETE /api/history/chat/{id}", handler.HandleDeleteChat())
	mux.HandleFunc("DELETE /api/history/generations/{id}", handler.HandleDeleteGen())
	// Proxy
	mux.HandleFunc("POST /api/proxy/", handler.HandleProxy(ds))
	// Task polling
	mux.HandleFunc("GET /api/tasks/{taskID}", handler.HandleTaskPoll(ds))
	// Health
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteJSON(w, 200, map[string]string{"status": "ok"})
	})

	handler := middleware.Logging(middleware.CORS(mux))

	log.Println("[Server] Starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}