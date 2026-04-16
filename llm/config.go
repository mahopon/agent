package llm

type LLMConfig struct {
	baseURL string
	apiKey string
	model string
}
func NewLLMConfig(url, key, model string) LLMConfig {
	return LLMConfig{
		baseURL: url,
		apiKey: key,
		model: model,
	}
}
