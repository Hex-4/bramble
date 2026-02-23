package config

import _ "embed"

//go:embed default.toml
var DefaultConfig string

//go:embed agent.md
var DefaultAgentMD string
