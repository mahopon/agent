package tool

import (
	"encoding/json"
	"fmt"
)

type ToolExecutor interface {
	Execute(name string, args map[string]any) (string, error)
	Schema() []Tool
}

func ToOpenAIScheme(executors []ToolExecutor) []Tool {
	tools := make([]Tool, 0)
	for _, executor := range executors {
		tools = append(tools, executor.Schema()...)
	}
	return tools
}

func FindExecutorForTool(executors []ToolExecutor, toolName string) (func(name string, args map[string]any) (string, error), error) {
	for _, executor := range executors {
		schema := executor.Schema()
		for _, tool := range schema {
			if tool.Function.Name == toolName {
				return executor.Execute, nil
			}
		}
	}
	return nil, fmt.Errorf("no executor found for tool: %s", toolName)
}

type ToolCallInfo struct {
	Name      string
	Arguments string
	ID        string
}

type ToolCallResult struct {
	ID      string
	Content string
	IsError bool
}

func ExecuteToolCalls(toolCalls []ToolCallInfo, executor func(name string, args map[string]any) (string, error)) ([]ToolCallResult, error) {
	results := make([]ToolCallResult, 0, len(toolCalls))
	for _, tc := range toolCalls {
		var args map[string]any
		if err := json.Unmarshal([]byte(tc.Arguments), &args); err != nil {
			results = append(results, ToolCallResult{
				ID:      tc.ID,
				Content: fmt.Sprintf("error: failed to parse arguments: %v", err),
				IsError: true,
			})
			continue
		}
		result, err := executor(tc.Name, args)
		if err != nil {
			results = append(results, ToolCallResult{
				ID:      tc.ID,
				Content: fmt.Sprintf("error: %v", err),
				IsError: true,
			})
			continue
		}
		results = append(results, ToolCallResult{
			ID:      tc.ID,
			Content: result,
		})
	}
	return results, nil
}
