package agent

import "agent/tool"

const (
	TITLE_PROMPT string = "You are a text summary agent. Your goal is to summarise the user's query into a sentence of not more than 10 words to be displayed as the title to a chat session."
)

func NewTitleAgent(session *Session) *Agent {
	return &Agent{
		Name:      "Title",
		Tools:     []tool.Tool{},
		SysPrompt: TITLE_PROMPT,
		Session:   session,
	}
}
