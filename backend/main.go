package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ============================================================
// Config
// ============================================================

type Config struct {
	ApiKey      string
	MySQLDSN    string
	StoragePath string
	BaseURL     string
}

func loadConfig() Config {
	return Config{
		ApiKey:      getEnv("BAILIAN_API_KEY", ""),
		MySQLDSN:    getEnv("MYSQL_DSN", "bailian:bailian_pass_2024@tcp(db:3306)/bailian?charset=utf8mb4&parseTime=true"),
		StoragePath: getEnv("STORAGE_PATH", "/data/storage"),
		BaseURL:     "https://dashscope.aliyuncs.com",
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// ============================================================
// Models
// ============================================================

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type ChatChoice struct {
	Index   int         `json:"index"`
	Message ChatMessage `json:"message"`
	Delta   ChatMessage `json:"delta,omitempty"`
	Finish  *string     `json:"finish_reason,omitempty"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
}

type DeltaResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ============================================================
// DB
// ============================================================

var db *sql.DB

func initDB(dsn string) error {
	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			db.SetMaxOpenConns(20)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(5 * time.Minute)
			log.Println("[DB] Connected")
			return nil
		}
		log.Printf("[DB] Waiting... (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return err
}

// ============================================================
// DashScope Client
// ============================================================

type DashScopeClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewDashScopeClient(apiKey string) *DashScopeClient {
	return &DashScopeClient{
		apiKey:  apiKey,
		baseURL: "https://dashscope.aliyuncs.com",
		client:  &http.Client{Timeout: 5 * time.Minute},
	}
}

func (c *DashScopeClient) do(path string, body any) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}

func (c *DashScopeClient) doAsync(path string, body any) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")
	return c.client.Do(req)
}

func (c *DashScopeClient) getTask(taskID string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v1/tasks/"+taskID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	return c.client.Do(req)
}

// ============================================================
// Middleware
// ============================================================

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s (%v)", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// ============================================================
// Helpers
// ============================================================

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, APIError{Code: status, Message: msg})
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// ============================================================
// Main
// ============================================================

func main() {
	cfg := loadConfig()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := initDB(cfg.MySQLDSN); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}
	defer db.Close()

	dashscope := NewDashScopeClient(cfg.ApiKey)

	mux := http.NewServeMux()

	// Chat
	mux.HandleFunc("POST /api/chat/completions", handleChat(dashscope))
	// Images
	mux.HandleFunc("POST /api/images/generate", handleImageGen(dashscope))
	// Videos
	mux.HandleFunc("POST /api/videos/generate", handleVideoGen(dashscope))
	// Audio
	mux.HandleFunc("POST /api/audio/speech", handleTTS(dashscope))
	mux.HandleFunc("POST /api/audio/transcribe", handleASR(dashscope))
	// Toolbox
	mux.HandleFunc("POST /api/toolbox/translate", handleTranslate(dashscope))
	mux.HandleFunc("POST /api/toolbox/ocr", handleOCR(dashscope))
	mux.HandleFunc("POST /api/toolbox/code", handleCode(dashscope))
	mux.HandleFunc("POST /api/toolbox/document", handleDocument(dashscope))
	// History
	mux.HandleFunc("GET /api/history/chat", handleChatHistory())
	mux.HandleFunc("GET /api/history/generations", handleGenHistory())
	mux.HandleFunc("DELETE /api/history/chat/{id}", handleDeleteChat())
	mux.HandleFunc("DELETE /api/history/generations/{id}", handleDeleteGen())
	// Proxy for generic API calls
	mux.HandleFunc("POST /api/proxy/", handleProxy(dashscope))
		// Task polling
		mux.HandleFunc("GET /api/tasks/{taskID}", handleTaskPoll(dashscope))
	// Health
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})

	handler := loggingMiddleware(corsMiddleware(mux))

	log.Println("[Server] Starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}