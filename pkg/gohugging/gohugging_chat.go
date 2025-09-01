package gohugging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

// Initial implementation for the gohugging package

// ...existing code...

// BasicStruct is a simple example struct
type GoHugging struct {
	ID          int
	Name        string
	BaseUrl     string
	apiKey      string
	MaxTokens   int
	Temperature float32
	Model       string
	Context     []Message
	debug       bool
}

// New returns a pointer to a new GoHugging instance
func New(configData []byte) (*GoHugging, error) {
	var config HuggingFaceConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return &GoHugging{
		ID:          0,
		Name:        "",
		BaseUrl:     config.Params.BaseUrl,
		apiKey:      config.Params.Token,
		Model:       config.Params.Model,
		MaxTokens:   config.Params.MaxTokens,
		Temperature: config.Params.Temperature,
		debug:       false,
	}, nil
}

// CallHuggingFaceAPI sends a payload to the Hugging Face API and returns the response or an error
func (gh *GoHugging) callHuggingFaceAPI(payload *HuggingFacePayload) (*HuggingFaceResponse, error) {
	url := gh.BaseUrl
	if url == "" {
		return nil, fmt.Errorf("BaseUrl is not set")
	}

	gh.DebugLog("Calling Hugging Face API at %s with model %s", url, payload.Model)
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	gh.DebugLog("Make request")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+gh.apiKey)
	req.Header.Set("Content-Type", "application/json")

	gh.DebugLog("Sending request")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	gh.DebugLog("Reading response body")
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var hfResp HuggingFaceResponse
	if err := json.Unmarshal(response, &hfResp); err != nil {
		return &HuggingFaceResponse{Error: string(response)}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &HuggingFaceResponse{Error: string(response)}, fmt.Errorf("HTTP error: status code %d, error: %s", resp.StatusCode, hfResp.Error)
	}

	stopReason := hfResp.Choices[0].FinishReason
	if stopReason != "stop" {
		payload.MaxTokens = payload.MaxTokens * 2
		if payload.MaxTokens > 8192 {
			return &hfResp, fmt.Errorf("max tokens exceeded the limit of 8192")
		}
		gh.DebugLog("Max tokens exceeded, retrying with increased limit")
		return gh.callHuggingFaceAPI(payload)
	}

	return &hfResp, nil
}

func (gh *GoHugging) SendChatMessage(message string) (string, error) {
	gh.DebugLog("Creating payload for chat message")
	payload := &HuggingFacePayload{
		Model: gh.Model,
		Messages: []Message{
			{
				Content: message,
				Role:    "user",
			},
		},
		MaxTokens:   gh.MaxTokens,
		Temperature: gh.Temperature,
	}

	// Append context messages to the payload
	gh.DebugLog("Appending context messages, count: %d", len(gh.Context))
	payload.Messages = append(payload.Messages, gh.Context...)

	// Call the Hugging Face API
	gh.DebugLog("Calling Hugging Face API")
	resp, err := gh.callHuggingFaceAPI(payload)
	if err != nil {
		return "", err
	}

	gh.DebugLog("Appending user message to context")
	gh.Context = append(gh.Context, Message{
		Content: message,
		Role:    "user",
	})

	gh.DebugLog("Appending assistant message to context")
	gh.Context = append(gh.Context, Message{
		Content: resp.Choices[0].Message.Content,
		Role:    "assistant",
	})
	gh.DebugLog("History size=%d", len(gh.Context))

	return resp.Choices[0].Message.Content, nil
}
