package tools

import (
	"os"
	"path/filepath"
)

func newReadFile(workspaceDir string) Tool {
	return Tool{
		Name:        "read_file",
		Description: "Read a file from the workspace",
		Emoji:       "📖",
		Parameters: map[string]Parameter{
			"path": {Type: "string", Description: "file path relative to workspace", Required: true},
		},
		DetailParam: "path",
		Execute: func(args map[string]string) (string, error) {
			data, err := os.ReadFile(filepath.Join(workspaceDir, args["path"]))
			return string(data), err
		},
	}
}

func newWriteFile(workspaceDir string) Tool {
	return Tool{
		Name:        "write_file",
		Description: "Write a file to the workspace",
		Emoji:       "📝",
		Parameters: map[string]Parameter{
			"path":    {Type: "string", Description: "file path relative to workspace", Required: true},
			"content": {Type: "string", Description: "file content", Required: true},
		},
		DetailParam: "path",
		Execute: func(args map[string]string) (string, error) {
			return "", os.WriteFile(filepath.Join(workspaceDir, args["path"]), []byte(args["content"]), 0644)
		},
	}
}
