package config

import "github.com/BurntSushi/toml"

type Config struct {
	BrambleDir string
	Agent      AgentConfig `toml:"agent"`
}

type AgentConfig struct {
	Name                string            `toml:"name"`
	Model               string            `toml:"model"`
	Provider            string            `toml:"provider"`
	SystemPrompt        string            `toml:"system_prompt"`
	SessionDescriptions map[string]string `toml:"session_descriptions"`
	ContextFiles        []string          `toml:"context_files"`
}

func Load(path string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	return cfg, err
}
