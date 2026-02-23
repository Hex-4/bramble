package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hex-4/bramble/config"

	"github.com/charmbracelet/lipgloss"
)

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
	thingHappened("❧", "initializing your bramble install...", "#68A655", 0)

	thingHappened("→", "creating ~/.bramble directory... (this is where bramble stores all of its data)", "#68A655", 0)

	home, _ := os.UserHomeDir()
	brambleDir := filepath.Join(home, ".bramble")

	if _, err := os.Stat(brambleDir); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bramble already exists, skipping", "#68A655", 3)
	} else if err := os.Mkdir(brambleDir, 0755); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bramble: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bramble created successfully", "#68A655", 3)
	}

	thingHappened("→", "creating ~/.bramble/config.toml... (set up with sensible defaults, but we'll configure it interactively later)", "#68A655", 0)

	if _, err := os.Stat(filepath.Join(brambleDir, "config.toml")); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bramble/config.toml already exists, skipping", "#68A655", 3)
	} else if err := os.WriteFile(filepath.Join(brambleDir, "config.toml"), []byte(config.DefaultConfig), 0644); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bramble/config.toml: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bramble/config.toml created successfully", "#68A655", 3)
	}

	thingHappened("→", "creating ~/.bramble/workspace/ (this is your agent's workspace, where it'll store its prompts + anything else it wants)", "#68A655", 0)

	if _, err := os.Stat(filepath.Join(brambleDir, "workspace")); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bramble/workspace already exists, skipping", "#68A655", 3)
	} else if err := os.Mkdir(filepath.Join(brambleDir, "workspace"), 0755); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bramble/workspace: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bramble/workspace created successfully", "#68A655", 3)
	}

	thingHappened("→", "creating ~/.bramble/workspace/agent.md (this is your agent's prompt)", "#68A655", 0)

	if _, err := os.Stat(filepath.Join(brambleDir, "workspace", "agent.md")); !os.IsNotExist(err) {
		thingHappened("┄", "~/.bramble/workspace/agent.md already exists, skipping", "#68A655", 3)
	} else if err := os.WriteFile(filepath.Join(brambleDir, "workspace", "agent.md"), []byte(config.DefaultAgentMD), 0644); err != nil {
		thingHappened("⚠", fmt.Sprintf("failed to create ~/.bramble/workspace/agent.md: %v", err), "#CE7527", 3)
		return
	} else {
		thingHappened("✓", "~/.bramble/workspace/agent.md created successfully", "#68A655", 3)
	}

	thingHappened("❧", "done for now.", "#68A655", 0)
}

func thingHappened(icon string, message string, color string, indent int) {
	var iconStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Width(2)
	var textStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#b8a78b"))

	fmt.Println(strings.Repeat(" ", indent), iconStyle.Render(icon), textStyle.Render(message))
}
