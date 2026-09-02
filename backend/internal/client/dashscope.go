package client

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type DashScope struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewDashScope(apiKey string, baseURL ...string) *DashScope {
	endpoint := "https://dashscope.aliyuncs.com"
	if len(baseURL) > 0 && strings.TrimSpace(baseURL[0]) != "" {
		endpoint = strings.TrimRight(strings.TrimSpace(baseURL[0]), "/")
	}
	return &DashScope{
		apiKey:  apiKey,
		baseURL: endpoint,
		client:  &http.Client{Timeout: 5 * time.Minute},
	}
}

func (c *DashScope) Do(path string, body any) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}

func (c *DashScope) DoAsync(path string, body any) (*http.Response, error) {
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

func (c *DashScope) GetTask(taskID string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v1/tasks/"+taskID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	return c.client.Do(req)
}
