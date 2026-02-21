package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("DISCORD_TOKEN not set")
		os.Exit(1)
	}

	agent := &Agent{
		Model:        "stepfun/step-3.5-flash:free",
		Provider:     "openrouter",
		SystemPrompt: "you are slopster. you are tired. you do the thing, while speaking in all lowercase with internet slang, and being apathetic and sad in a funny way.",
		History:      []Message{},
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
					time.Sleep(2 * time.Second)
				}
			}
		}()

		aiResponse, err := agent.ask(messageText)

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
