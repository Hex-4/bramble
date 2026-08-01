package scheduler

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/Hex-4/bop/ai"
	"github.com/Hex-4/bop/tools"
	"github.com/Hex-4/bop/triggers"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Cron               *cron.Cron
	Jobs               map[string]Job
	sessions           *triggers.SessionStore
	ExtraToolsProvider func(channelID string) map[string]tools.Tool
	Agent              *ai.Agent
}

type Job struct {
	ID          string       `json:"id"`
	CronExpr    string       `json:"cronExpr,omitempty"` // Only for recurring jobs
	CronEntryID cron.EntryID `json:"cronEntryId,omitempty"`
	FireAt      time.Time    `json:"fireAt,omitempty"` // Only for one-shots
	Prompt      string       `json:"prompt"`
	SessionID   string       `json:"sessionId"`
	Silent      bool         `json:"silent"`
}

func NewScheduler(agent *ai.Agent) *Scheduler {
	sessions := triggers.NewSessionStore()

	return &Scheduler{
		Cron:     cron.New(),
		Jobs:     make(map[string]Job),
		sessions: sessions,
		Agent:    agent,
	}
}

func generateJobID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b) // e.g. "a3f1b20c"
}

func (s *Scheduler) fire(jobID string, wrappedPrompt string, sessionID string, silent bool) {
	prompt := s.Agent.SystemPrompt()
	niceTimeString := time.Now().Format("2006-01-02T15:04:05")
	prompt = prompt + " (current time: " + niceTimeString + ", local time, YYYY-MM-DDTHH:MM:SS)"

	history, _ := s.sessions.Load(jobID)
	history = append(history, ai.Message{Role: "system", Content: wrappedPrompt})

	messages := []ai.Message{
		{Role: "system", Content: prompt},
	}
	messages = append(messages, history...)

	var extraTools map[string]tools.Tool

	if !silent {
		channelID := strings.TrimPrefix(sessionID, "discord:")
		extraTools = s.ExtraToolsProvider(channelID)
	}

	newMessages, err := s.Agent.Ask(messages, extraTools)

	if err != nil {
		fmt.Printf("cron job failed: %v\n", err)
		return
	}
	history = append(history, newMessages...)

	if !silent {
		if !triggers.IsLastToolCall(newMessages, "send_message") {
			reminder := ai.Message{Role: "system", Content: "You ended your turn without calling send_message. Your last response was never delivered — the user only sees messages sent via send_message. If it contained anything the user needs, call send_message now with that content (you can reuse what you already wrote). If this was intentional, end your turn without calling any tools."}
			history = append(history, reminder)
			retryMessages := append([]ai.Message{{Role: "system", Content: prompt}}, history...)
			retriedNewMessages, err := s.Agent.Ask(retryMessages, extraTools)
			if err == nil {
				history = append(history, retriedNewMessages...)
			}
		}
	}

	s.sessions.Save(jobID, history)
}

func (s *Scheduler) AddCron(expression string, prompt string, sessionID string, silent bool) (string, error) {
	jobID := generateJobID()
	wrappedPrompt := "The following is a background cron job, not a live user message. Your history may include past iterations of this job. Job ID: " + jobID + ". Use tools as normal. Execute the following: " + prompt

	cronFunc := func() {
		s.fire(jobID, wrappedPrompt, sessionID, silent)
	}

	entryID, err := s.Cron.AddFunc(expression, cronFunc)
	if err != nil {
		return "", fmt.Errorf("invalid cron expression: %w", err)
	}
	s.Jobs[jobID] = Job{
		CronExpr:    expression,
		ID:          jobID,
		CronEntryID: entryID,
		Prompt:      prompt,
		SessionID:   sessionID,
		Silent:      silent,
	}
	return jobID, nil
}

func (s *Scheduler) AddOneShot(fireAt time.Time, prompt string, sessionID string, silent bool) string {
	jobID := generateJobID()
	wrappedPrompt := "The following is a scheduled one-shot job, not a live user message. Job ID: " + jobID + ". Use tools as normal. Execute the following: " + prompt

	oneShotFunc := func() {
		s.fire(jobID, wrappedPrompt, sessionID, silent)
		delete(s.Jobs, jobID)
	}
	time.AfterFunc(time.Until(fireAt), oneShotFunc)

	s.Jobs[jobID] = Job{
		ID:        jobID,
		FireAt:    fireAt,
		Prompt:    prompt,
		SessionID: sessionID,
		Silent:    silent,
	}
	return jobID
}

func (s *Scheduler) RemoveJob(jobID string) error {
	job, ok := s.Jobs[jobID]
	if !ok {
		return fmt.Errorf("Could not find job ID")
	}
	if job.CronExpr != "" {
		s.Cron.Remove(job.CronEntryID)
	}
	delete(s.Jobs, jobID)
	return nil
}
