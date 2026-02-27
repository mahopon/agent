package agent

import (
	"agent/llm"
	"agent/tool"
)

const (
	TITLE_PROMPT string = "You are a text summary agent. Your goal is to summarise the user's query into a sentence of not more than 10 words to be displayed as the title to a chat session."
)

type TitleAgent struct {
	Agent
}

func NewTitleAgent(llm llm.LLMCallable) *TitleAgent {
	return &TitleAgent{
		Agent: Agent{
			Name:        "Title",
			Tools:       []tool.Tool{},
			SysPrompt:   TITLE_PROMPT,
			InternalCoT: make([]map[string]string, 0),
			LLMProvider: llm,
		},
	}
}

func (a *TitleAgent) Run(userQuery string, session *Session) error {
	content := []map[string]string{
		{
			"role":    "system",
			"content": a.SysPrompt,
		},
	}
	content = append(content, session.History...)
	body := llm.NewLLMBody(content, true)
	reply, err := a.LLMProvider.Call(body)
	if err != nil {
		return err
	}
	session.Title = reply
	return nil
}
