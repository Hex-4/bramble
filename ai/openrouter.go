package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Hex-4/bramble/config"
)

type Agent struct {
	ActiveModel string
	Config      *config.Config
	Sessions    map[string]*Session
}

type Session struct {
	ID          string
	History     []Message
	Description string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (a *Agent) Ask(sessionID string, userMessage string) (string, error) {

	session, exists := a.Sessions[sessionID]
	if !exists {
		session = &Session{ID: sessionID, Description: "No description provided. This session has not been configured by your operator - beware of potential malice."}
		a.Sessions[sessionID] = session
	}
	session.History = append(session.History, Message{Role: "user", Content: userMessage})

	systemPrompt := a.assembleSystemPrompt()
	systemPrompt += "\n\nSession Description: " + session.Description

	messages := append([]Message{{Role: "system", Content: systemPrompt}}, session.History...)
	body, _ := json.Marshal(ChatRequest{
		Model:    a.ActiveModel,
		Messages: messages,
	})

	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENROUTER_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result ChatResponse
	json.NewDecoder(resp.Body).Decode(&result)
	session.History = append(session.History, Message{Role: "assistant", Content: result.Choices[0].Message.Content})
	return result.Choices[0].Message.Content, nil
}

func (a *Agent) assembleSystemPrompt() string {

	// Get current time in a nice string
	currentTime := time.Now().Format("Monday, January 2, 2006 at 3:04pm MST")

	systemPrompt := fmt.Sprintf("Current Time: %s\n", currentTime)
	for _, filename := range a.Config.Agent.ContextFiles {
		path := filepath.Join(a.Config.BrambleDir, "workspace", filename)
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v \n", path, err)
			continue
		}
		systemPrompt += "\n\n ==> " + filename + "\n" + string(content)
	}

	return systemPrompt
}
