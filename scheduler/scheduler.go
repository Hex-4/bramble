package scheduler

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Hex-4/bramble/ai"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Cron     *cron.Cron
	Jobs     map[string]Job
	sendFunc func(sessionID string, text string)
	Agent    *ai.Agent
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

func NewScheduler(agent *ai.Agent, sendFunc func(sessionID string, text string)) *Scheduler {
	return &Scheduler{
		Cron:     cron.New(),
		Jobs:     make(map[string]Job),
		sendFunc: sendFunc,
		Agent:    agent,
	}
}

func generateJobID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b) // e.g. "a3f1b20c"
}

func newCronStatusHandler() *ai.StatusHandler {
	var completedEmojis []string

	return &ai.StatusHandler{
		OnToolStart: func(emoji string, detail string) {
			completedEmojis = append(completedEmojis, emoji)
		},
		OnDone: func() {
		},
		Footer: func() string {
			if completedEmojis == nil {
				return ""
			}
			return strings.Join(completedEmojis, "") + " " + strconv.Itoa(len(completedEmojis)) + " tools"
		},
	}
}

func (s *Scheduler) AddCron(expression string, prompt string, sessionID string, silent bool) (string, error) {
	var cronFunc func()
	if silent {
		cronFunc = func() {
			_, err := s.Agent.Ask(sessionID, prompt, newCronStatusHandler())
			if err != nil {
				fmt.Printf("cron job failed: %v\n", err)
			}
		}
	} else {
		cronFunc = func() {
			statusHandler := newCronStatusHandler()
			response, err := s.Agent.Ask(sessionID, prompt, statusHandler)
			if err != nil {
				fmt.Printf("cron job failed: %v\n", err)
				s.sendFunc(sessionID, "⚠️ Cron job failed: "+err.Error())
				return
			}
			s.sendFunc(sessionID, response+"\n-# "+statusHandler.Footer())
		}
	}
	entryID, err := s.Cron.AddFunc(expression, cronFunc)
	jobID := generateJobID()
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
	var oneShotFunc func()
	if silent {
		oneShotFunc = func() {
			_, err := s.Agent.Ask(sessionID, prompt, newCronStatusHandler())
			if err != nil {
				fmt.Printf("one shot job failed: %v\n", err)
			}
			delete(s.Jobs, jobID)
		}
	} else {
		oneShotFunc = func() {
			statusHandler := newCronStatusHandler()
			response, err := s.Agent.Ask(sessionID, prompt, statusHandler)
			delete(s.Jobs, jobID)
			if err != nil {
				fmt.Printf("one shot job failed: %v\n", err)
				s.sendFunc(sessionID, "⚠️ One shot job failed: "+err.Error())
				return
			}
			s.sendFunc(sessionID, response+"\n-# "+statusHandler.Footer())

		}
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
		return fmt.Errorf("Could not find job IF")
	}
	if job.CronExpr != "" {
		s.Cron.Remove(job.CronEntryID)
	}
	delete(s.Jobs, jobID)
	return nil
}
