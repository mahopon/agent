package agent

import (
	"agent/memory"
)


type HistoryStore interface {
	Add(memory.ChatMessage)
	Remove(memory.ChatMessage)
	Size() int
	ToHistory() []map[string]string
	TokenCount(func(count int))
}

// A session should contain history, as well as a summary of the initial message

type Session struct {
	Id      string
	History HistoryStore
}

func NewSession(existingHistory HistoryStore) *Session {
	var history HistoryStore
	if existingHistory != nil {
		history = existingHistory
	}
	return &Session{
		History: history,
	}
}
