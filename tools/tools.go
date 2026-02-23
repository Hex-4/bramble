package tools

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]Parameter // for the JSON schema sent to the model
	Execute     func(args map[string]string) (string, error)
}

type Parameter struct {
	Type        string
	Description string
	Required    bool
}
