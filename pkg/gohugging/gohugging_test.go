package gohugging

import (
	"fmt"
	"os"
	"testing"
)

func GetConfig() []byte {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configData, err := os.ReadFile(fmt.Sprintf("%s%c.clai.%c%s", home, os.PathSeparator, os.PathSeparator, "huggingface.yaml"))
	if err != nil {
		panic(err)
	}

	return configData
}

func TestNew(t *testing.T) {
	configData := []byte(`
huggingface:
  token: "your_api_key"
  model: "test/model"
  timeout: 30
  max_tokens: 1024
  base_url: "https://api.huggingface.co"
`)
	gohugging, err := New(configData)
	if err != nil {
		t.Fatalf("Failed to create GoHugging instance: %v", err)
	}

	if gohugging.BaseUrl != "https://api.huggingface.co" {
		t.Errorf("Expected BaseUrl to be 'https://api.huggingface.co', got %q", gohugging.BaseUrl)
	}
	if gohugging.apiKey != "your_api_key" {
		t.Errorf("Expected apiKey to be 'your_api_key', got %q", gohugging.apiKey)
	}

	if gohugging.Model != "test/model" {
		t.Errorf("Expected Model to be 'test/model', got %q", gohugging.Model)
	}
}

func TestCallHuggingFaceAPI(t *testing.T) {
	configData := GetConfig()

	gohugging, err := New(configData)
	if err != nil {
		t.Fatalf("Failed to create GoHugging instance: %v", err)
	}

	payload := &HuggingFacePayload{
		Model: gohugging.Model,
		Messages: []Message{
			{
				Content:          "Hello, world!",
				ReasoningContent: "This is a test message.",
				Role:             "user",
			},
		},
		MaxTokens:   gohugging.MaxTokens,
		Temperature: gohugging.Temperature,
	}

	// This test will fail unless you mock the HTTP request or use a valid token/model
	resp, err := gohugging.callHuggingFaceAPI(payload)
	if err == nil {
		t.Logf("Received response: %+v", resp)
	} else {
		fmt.Println("response is:", resp)
		t.Logf("Expected error (likely due to test token/model): %v", err)
	}
}

func TestSendChatMessage(t *testing.T) {
	configData := GetConfig()
	gohugging, err := New(configData)
	if err != nil {
		t.Fatalf("Failed to create GoHugging instance: %v", err)
	}

	reply, err := gohugging.SendChatMessage("Hello, Hugging Face!")
	if err == nil {
		t.Logf("Received reply: %s", reply)
	} else {
		t.Logf("Expected error (likely due to test token/model): %v", err)
	}
}

func TestAvailableModelsFromAPI(t *testing.T) {
	// Use a dummy or real API key for testing

	configData := GetConfig()
	gohugging, err := New(configData)
	if err != nil {
		t.Fatalf("Failed to create GoHugging instance: %v", err)
	}

	models, err := gohugging.AvailableModels()
	if err != nil {
		t.Logf("Expected error or empty result if no network or invalid key: %v", err)
	} else {
		t.Logf("Fetched %d models from API", len(models))
		if len(models) == 0 {
			t.Errorf("Expected at least one model, got zero")
		}
	}
}
