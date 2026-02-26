package llm

type LLMConfig struct {
	LLM_URL   string
	LLM_KEY   string
	LLM_MODEL string
}

func NewLLMConfig() *LLMConfig {
	return &LLMConfig{
		LLM_URL:   "http://localhost:8080/v1/chat/completions",
		LLM_KEY:   "",
		LLM_MODEL: "",
	}
}
