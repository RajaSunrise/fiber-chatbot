package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Prompt struct for request payload
type Prompt struct {
	Prompt string `json:"prompt"`
}

// BLACKBOXAI struct for handling requests
type BLACKBOXAI struct {
	ChatEndpoint string
	Timeout      int
	LastResponse map[string]interface{}
	Headers      map[string]string
}

// NewBLACKBOXAI creates a new BLACKBOXAI instance
func NewBLACKBOXAI() *BLACKBOXAI {
	return &BLACKBOXAI{
		ChatEndpoint: "https://www.blackbox.ai/api/chat",
		Timeout:      30,
		LastResponse: make(map[string]interface{}),
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "*/*",
			"Accept-Encoding": "Identity",
		},
	}
}

func (b *BLACKBOXAI) Ask(prompt string, stream bool, raw bool, optimizer string, conversationally bool) (map[string]interface{}, error) {
    payload := map[string]interface{}{
        "messages": []map[string]string{
            {"content": prompt, "role": "user"},
        },
        "id":                 "",
        "previewToken":       "",
        "userId":             "",
        "codeModelMode":      true,
        "agentMode":          make(map[string]interface{}),
        "trendingAgentMode":  make(map[string]interface{}),
        "isMicMode":          false,
    }

    client := &http.Client{Timeout: time.Duration(b.Timeout) * time.Second}
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("error marshalling payload: %w", err)
    }

    req, err := http.NewRequest("POST", b.ChatEndpoint, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    for key, value := range b.Headers {
        req.Header.Set(key, value)
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    // Log the response body for debugging
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }
    fmt.Printf("Raw response body: %s\n", bodyBytes)

    // Reset response body
    resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to generate response - (%d, %s)", resp.StatusCode, resp.Status)
    }

    var response map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&response)
    if err != nil {
        return nil, fmt.Errorf("error decoding response: %w", err)
    }

    b.LastResponse = response
    return response, nil
}
