package llm

type LLMConfig struct {
	LLM_URL   string
	LLM_KEY   string
	LLM_MODEL string
}

func NewLLMConfig(url, key, model string) *LLMConfig {
	return &LLMConfig{
		LLM_URL:   url,
		LLM_KEY:   key,
		LLM_MODEL: model,
	}
}

