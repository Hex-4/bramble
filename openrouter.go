package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type Agent struct {
	Model        string
	Provider     string
	SystemPrompt string
	History      []Message
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

func (a *Agent) ask(userMessage string) (string, error) {
	a.History = append(a.History, Message{Role: "user", Content: userMessage})

	messages := append([]Message{{Role: "system", Content: a.SystemPrompt}}, a.History...)
	body, _ := json.Marshal(ChatRequest{
		Model:    "stepfun/step-3.5-flash:free",
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
	a.History = append(a.History, Message{Role: "assistant", Content: result.Choices[0].Message.Content})
	return result.Choices[0].Message.Content, nil
}
