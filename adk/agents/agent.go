package agents

// Agent represents an agent that can be used as a tool in LLM requests
type Agent struct {
	Name        string
	Description string
	Url         string
}
