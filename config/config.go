package config

import "github.com/BurntSushi/toml"

type Config struct {
    Agent AgentConfig `toml:"agent"`
}

type AgentConfig struct {
    Model        string   `toml:"model"`
    ContextFiles []string `toml:"context_files"`
}

func Load(path string) (Config, error) {
    var cfg Config
    _, err := toml.DecodeFile(path, &cfg)
    return cfg, err
}