package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Hex-4/bramble/ai"
	"github.com/Hex-4/bramble/config"
	"github.com/Hex-4/bramble/tools"

	"github.com/joho/godotenv"
)

func runServer() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting user home directory:", err)
		os.Exit(1)
	}

	defaultHome := filepath.Join(userHome, ".bramble")
	configFile, err := config.Load(filepath.Join(defaultHome, "config.toml"))

	configFile.BrambleDir = defaultHome

	if err != nil {
		fmt.Println("error loading config:", err)
		os.Exit(1)
	}

	godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("DISCORD_TOKEN not set")
		os.Exit(1)
	}

	agent := &ai.Agent{
		ActiveModel: configFile.Agent.Model,
		Config:      &configFile,
		Sessions:    make(map[string]*ai.Session),
		Tools:       tools.NewRegistry(filepath.Join(defaultHome, "workspace")),
	}

	agent.ToolSchemas = tools.NewSchemaList(agent.Tools)

	for sessionID, sessionDescription := range configFile.Agent.SessionDescriptions {
		session := &ai.Session{
			ID:          sessionID,
			Description: sessionDescription,
			History:     make([]ai.Message, 0),
		}
		agent.Sessions[session.ID] = session
	}

	discordBot, err := NewDiscordBot(token, agent)
	if err != nil {
		fmt.Println("error creating discord bot:", err)
		os.Exit(1)
	}

	err = discordBot.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		os.Exit(1)
	}

	fmt.Println("slopster is online 🫥")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	discordBot.Close()
}
