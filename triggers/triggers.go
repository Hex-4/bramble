package triggers

import (
	"sync"

	"github.com/Hex-4/bop/ai"
)

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string][]ai.Message
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string][]ai.Message),
	}
}

func IsLastToolCall(messages []ai.Message, name string) bool {
	for i := len(messages) - 1; i >= 0; i-- {
		if len(messages[i].ToolCalls) > 0 {
			for _, call := range messages[i].ToolCalls {
				if call.Function.Name == name {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (store *SessionStore) Load(sessionID string) ([]ai.Message, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if messages, ok := store.sessions[sessionID]; ok {
		return messages, nil
	}
	return nil, nil
}

func (store *SessionStore) Save(sessionID string, messages []ai.Message) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.sessions[sessionID] = messages
}
