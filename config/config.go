package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	API_KEY   string
	LLM_MODEL string
	LLM_URL   string
	DEBUG     bool
}

func NewConfig() *Config {
	var LLM_URL string
	var API_KEY string
	LLM_MODE := os.Getenv("LLM_MODE")
	if LLM_MODE == "LOCAL" {
		LLM_URL = os.Getenv("LLM_LOCAL_URL")
	} else if LLM_MODE == "API" {
		LLM_URL = os.Getenv("LLM_URL")
		API_KEY = os.Getenv("API_KEY")
	} else {
		panic(errors.New("LLM_MODE not specified"))
	}
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		panic(err)
	}
	return &Config{
		API_KEY:   API_KEY,
		LLM_MODEL: os.Getenv("LLM_MODEL"),
		LLM_URL:   LLM_URL,
		DEBUG:     debug,
	}
}
