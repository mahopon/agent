package agent

// A session should contain history, as well as a summary of the initial message

type Session struct {
	History []map[string]any
}

func NewSession() *Session {
	return &Session{
		History: make([]map[string]any, 0),
	}
}
