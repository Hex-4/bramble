package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hex-4/bop/config"

	"github.com/charmbracelet/lipgloss"
)

const defaultEnv = `# bop secrets — loaded by ` + "`bop server run`" + ` at startup.
# keep this file private. never commit it, never paste it anywhere.

# your discord bot token (required — bop won't boot without this)
# discord developer portal → your application → bot → reset token
# https://discord.com/developers/applications
DISCORD_TOKEN=

# openrouter API key (required — this is how the agent thinks)
# https://openrouter.ai/keys
OPENROUTER_KEY=

# composio API key (required for now — external app tools like gmail/github)
# https://composio.dev
COMPOSIO_KEY=

# exa API key (optional — only the web_search / web_fetch tools use this)
# https://exa.ai
EXA_KEY=
`

const (
	Banner = `▄     ▄▄
█    █  ▀               ▄   ▄▄
█▀▀▄ █  ▄▀▀█  ▄ ▄  ▄    █  █▄▄█
▀▄▄▀    ▀▄▄█ █ █ █ █    █  ▀▄▄▄
             █ █ █ █▀▀▄ ▀▄
                   ▀▄▄▀`
)

func RunInit() {
	var bannerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#68A655")).
		Width(31)

	var taglineStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#B8A78B"))

	fmt.Println(bannerStyle.Render(Banner))
	fmt.Println(taglineStyle.Render("🥀 the hackable AI agent for people who'd rather make human slop"))
	fmt.Println()
	thingHappened("❧", "initializing your bop install...", "#68A655", 0)

	thingHappened("→", "creating ~/.bop directory... (this is where bop stores all of its data)", "#68A655", 0)

	home, _ := os.UserHomeDir()
	bopDir := filepath.Join(home, ".bop")

	if _, err := os.Stat(bopDir); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bop already exists, skipping", "#68A655", 3)
	} else if err := os.Mkdir(bopDir, 0755); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bop: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bop created successfully", "#68A655", 3)
	}

	thingHappened("→", "creating ~/.bop/config.toml... (set up with sensible defaults, but we'll configure it interactively later)", "#68A655", 0)

	if _, err := os.Stat(filepath.Join(bopDir, "config.toml")); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bop/config.toml already exists, skipping", "#68A655", 3)
	} else if err := os.WriteFile(filepath.Join(bopDir, "config.toml"), []byte(config.DefaultConfig), 0644); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bop/config.toml: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bop/config.toml created successfully", "#68A655", 3)
	}

	thingHappened("→", "creating ~/.bop/workspace/ (this is your agent's workspace, where it'll store its prompts + anything else it wants)", "#68A655", 0)

	if _, err := os.Stat(filepath.Join(bopDir, "workspace")); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bop/workspace already exists, skipping", "#68A655", 3)
	} else if err := os.Mkdir(filepath.Join(bopDir, "workspace"), 0755); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bop/workspace: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bop/workspace created successfully", "#68A655", 3)
	}

	thingHappened("→", "creating ~/.bop/workspace/agent.md (this is your agent's prompt)", "#68A655", 0)

	if _, err := os.Stat(filepath.Join(bopDir, "workspace", "agent.md")); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bop/workspace/agent.md already exists, skipping", "#68A655", 3)
	} else if err := os.WriteFile(filepath.Join(bopDir, "workspace", "agent.md"), []byte(config.DefaultAgentMD), 0644); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bop/workspace/agent.md: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bop/workspace/agent.md created successfully", "#68A655", 3)
	}

	thingHappened("→", "creating ~/.bop/.env (this is where your API keys live — fill it in!)", "#68A655", 0)

	if _, err := os.Stat(filepath.Join(bopDir, ".env")); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bop/.env already exists, skipping", "#68A655", 3)
	} else if err := os.WriteFile(filepath.Join(bopDir, ".env"), []byte(defaultEnv), 0600); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bop/.env: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bop/.env created (mode 0600, only you can read it)", "#68A655", 3)
	}

	thingHappened("❧", "done! next: add your keys to ~/.bop/.env, then run: bop server run", "#68A655", 0)
}

func thingHappened(icon string, message string, color string, indent int) {
	var iconStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Width(2)
	var textStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#b8a78b"))

	fmt.Println(strings.Repeat(" ", indent), iconStyle.Render(icon), textStyle.Render(message))
}
