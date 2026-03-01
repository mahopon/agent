package agent

import (
	"agent/llm"
	"agent/tool"
	"encoding/json"
	"os"
)

type CodingAgent struct {
	Agent
}

const (
	CODING_PROMPT string = `You are a coding agent. Your job is to execute the following steps that have already been planned out.

Return the generated files ONLY in valid JSON format.

If there is ONE file:
{
  "file_name": "example.go",
  "file_content": "file contents here"
}


If there are MULTIPLE files:
[

  {
    "file_name": "file1.go",
    "file_content": "file contents here"
  }
]


Do not include explanations. Do not include markdown. Output raw JSON only.


Instructions:`
)

type CodeFile struct {
	Name    string `json:"file_name"`
	Content string `json:"file_content"`
}

func NewCodingAgent(llm llm.LLMCallable) *CodingAgent {
	return &CodingAgent{
		Agent: Agent{
			Name:        "Planner",
			Tools:       []tool.Tool{},
			SysPrompt:   CODING_PROMPT,
			InternalCoT: make([]map[string]string, 0),
			LLMProvider: llm,
		},
	}
}
func (a *CodingAgent) Run(userQuery string, session *Session) error {
	promptWithInstruct := a.SysPrompt + session.Instructions
	content := []map[string]string{
		{
			"role":    "system",
			"content": promptWithInstruct,
		},
	}
	content = append(content, session.History...)
	body := llm.NewLLMBody(content, true)
	reply, err := a.LLMProvider.Call(body)
	if err != nil {
		return err
	}
	session.History = append(session.History, map[string]string{
		"role":    "assistant",
		"content": reply,
	})
	files, err := a.parseReply(reply)
	if err != nil {
		return err
	}
	err = a.writeFiles(files)
	if err != nil {
		return err
	}
	return nil
}

func (a *CodingAgent) parseReply(reply string) ([]CodeFile, error) {
	var codeFiles []CodeFile
	err := json.Unmarshal([]byte(reply), &codeFiles)
	if err != nil {
		return nil, err
	}
	return codeFiles, nil
}

func (a *CodingAgent) writeFiles(files []CodeFile) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, file := range files {
		os.WriteFile(path+"/"+file.Name, []byte(file.Content), 0644)
	}
	return nil
}
