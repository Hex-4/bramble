package tools

type Tool struct {
	Name        string
	Description string
	Emoji       string
	Parameters  map[string]Parameter // for the JSON schema sent to the model
	DetailParam string
	Execute     func(args map[string]string) (string, error)
}

type Parameter struct {
	Type        string
	Description string
	Required    bool
}

func NewRegistry(workspaceDir string) map[string]Tool {
	return map[string]Tool{
		"read_file":  newReadFile(workspaceDir),
		"write_file": newWriteFile(workspaceDir),
	}
}

func NewSchemaList(tools map[string]Tool) []map[string]any {
	schemas := make([]map[string]any, 0, len(tools))
	for _, t := range tools {
		schemas = append(schemas, t.ToSchema())
	}
	return schemas
}

func (t *Tool) ToSchema() map[string]any {

	properties := make(map[string]any)
	for name, param := range t.Parameters {
		properties[name] = map[string]any{
			"type":        param.Type,
			"description": param.Description,
		}
	}

	required := []string{}
	for name, param := range t.Parameters {
		if param.Required {
			required = append(required, name)
		}
	}

	return map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        t.Name,
			"description": t.Description,
			"parameters": map[string]any{
				"type":       "object",
				"properties": properties,
				"required":   required,
			},
		},
	}
}
