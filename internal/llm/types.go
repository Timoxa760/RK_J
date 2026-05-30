package llm

type generateRequest struct {
	SystemInstruction *contentBlock    `json:"systemInstruction,omitempty"`
	Contents          []contentMessage `json:"contents"`
	GenerationConfig  generationConfig `json:"generationConfig,omitempty"`
}

type contentBlock struct {
	Parts []part `json:"parts"`
}

type contentMessage struct {
	Role  string `json:"role,omitempty"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type generateResponse struct {
	Candidates []struct {
		Content contentBlock `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}
