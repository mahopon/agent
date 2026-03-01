package agent

import ()

// A session should contain history, as well as a summary of the initial message

type Session struct {
	Id      string
	History []map[string]any
}

func NewSession(existingHistory []map[string]any) *Session {
	var history []map[string]any
	if existingHistory != nil {
		history = existingHistory
	}
	return &Session{
		History: history,
	}
}
