package api

type AgentListResponseAgent struct {
	Name string `json:"name"`
}

type SessionListResponseAgent struct {
	Name    string   `json:"name"`
	Agent   string   `json:"agent"`
	Masters []string `json:"masters"`
	State   string   `json:"state"`
}

type CreateSession struct {
	Agent   string
	Command string
}
