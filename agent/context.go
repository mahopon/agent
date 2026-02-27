package agent

type ContextBuilder interface {
	BuildContext(agent *Agent, session *Session) []map[string]string
}
