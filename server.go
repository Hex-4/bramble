package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Hex-4/bramble/config"

	"github.com/bwmarrin/discordgo"
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

	agent := &Agent{
		Model:        configFile.Agent.Model,
		Provider:     configFile.Agent.Provider,
		SystemPrompt: configFile.Agent.SystemPrompt,
		Sessions:     make(map[string]*Session),
	}

	for sessionID, sessionDescription := range configFile.Agent.SessionDescriptions {
		session := &Session{
			ID:          sessionID,
			Description: sessionDescription,
			History:     make([]Message, 0),
		}
		agent.Sessions[session.ID] = session
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating session:", err)
		os.Exit(1)
	}

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		isDM := m.GuildID == ""

		if !isDM {
			mentioned := false
			for _, user := range m.Mentions {
				if user.ID == s.State.User.ID {
					mentioned = true
					break
				}
			}
			if !mentioned {
				return
			}
		}
		messageText := m.Message.Content
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				default:
					s.ChannelTyping(m.ChannelID)
					time.Sleep(10 * time.Second)
				}
			}
		}()

		// "log" the session id
		fmt.Println("Session ID:", "discord:"+m.ChannelID)
		aiResponse, err := agent.ask("discord:"+m.ChannelID, messageText)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "something broke. slopster is sorry. here's the error: "+err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, aiResponse)
		done <- true
	})
	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		os.Exit(1)
	}

	fmt.Println("slopster is online 🫥")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	dg.Close()
}
