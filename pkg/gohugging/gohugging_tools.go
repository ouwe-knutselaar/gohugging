package gohugging

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (gh *GoHugging) EnableDebugging() {
	gh.debug = true
	gh.DebugLog("Debugging enabled")
}

func (gh *GoHugging) DebugLog(format string, v ...interface{}) {
	if gh.debug {
		fmt.Printf(format, v...)
		fmt.Println()
	}

}

func (gh *GoHugging) Clear(){
	gh.Context = []Message{}
	gh.DebugLog("Chat history cleared")
}

// AvailableModels fetches available models from the Hugging Face API
func (gh *GoHugging) AvailableModels() ([]HuggingFaceModel, error) {
	gh.DebugLog("Fetching available models from Hugging Face API")
	url := "https://huggingface.co/api/models"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if gh.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+gh.apiKey)
	}

	gh.DebugLog("Making request to %s", url)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var models []HuggingFaceModel

	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	gh.DebugLog("Fetched %d models", len(models))

	return models, nil
}
