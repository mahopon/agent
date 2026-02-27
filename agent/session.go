package agent

// A session should contain history, as well as a summary of the initial message

type Session struct {
	History      []map[string]string
	Title        string
	Instructions string
}

func NewSession() *Session {
	return &Session{
		History:      make([]map[string]string, 0),
		Title:        "",
		Instructions: "",
	}
}
