package llm

type LLMConfig struct {
	LLM_URL string
	LLM_KEY string
}

func NewLLMConfig(url, key string) *LLMConfig {
	return &LLMConfig{
		LLM_URL: url,
		LLM_KEY: key,
	}
}
