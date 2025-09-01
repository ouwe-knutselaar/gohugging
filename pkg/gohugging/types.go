package gohugging

// HuggingFacePayload represents the payload sent to Hugging Face API
type HuggingFacePayload struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float32   `json:"temperature"`
}

type Message struct {
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
	Role             string `json:"role"`
}

// HuggingFaceResponse represents the response from Hugging Face API
type HuggingFaceResponse struct {
	GeneratedText string `json:"generated_text,omitempty"`
	Error         string `json:"error,omitempty"`
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choice  `json:"choices"`
	Usage   UsageInfo `json:"usage"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}


// HuggingFaceConfig represents the configuration for the Hugging Face API
type HuggingFaceConfig struct {
	Params HuggingFaceConfigParms `yaml:"huggingface"`
}

type HuggingFaceConfigParms struct {
	Token     string `yaml:"token"`
	Model     string `yaml:"model"`
	Timeout   int    `yaml:"timeout"`
	MaxTokens int    `yaml:"max_tokens"`
	BaseUrl   string `yaml:"base_url"`
	Temperature float32 `yaml:"temperature"`
}


type HuggingFaceModel struct {
	Id		 string `json:"_id"`
	ModelId   string `json:"modelId"`
	PipelineTag	string `json:"pipeline_tag"`
	Tags []string `json:"tags"`
}
