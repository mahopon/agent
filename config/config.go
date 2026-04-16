package config

import (
	"os"
)

type Config struct {
	API_KEY   string
	LLM_MODEL string
	LLM_URL   string
}

func NewConfig() *Config {
	var LLM_URL string
	var API_KEY string
	LLM_MODE := os.Getenv("LLM_MODE")
	if LLM_MODE == "API" {
		LLM_URL = os.Getenv("LLM_URL")
		API_KEY = os.Getenv("API_KEY")
	} else if LLM_MODE == "LOCAL" {
		LLM_URL = os.Getenv("LLM_LOCAL_URL")
	} else {
		LLM_URL = "http://192.168.1.251:8080/v1/chat/completions"
	}
	return &Config{
		API_KEY:   API_KEY,
		LLM_MODEL: os.Getenv("LLM_MODEL"),
		LLM_URL:   LLM_URL,
	}
}
