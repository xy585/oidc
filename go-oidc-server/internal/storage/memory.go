package storage

import (
	"sync"
)

// MemoryStorage is an in-memory storage for user sessions and tokens.
type MemoryStorage struct {
	sessions map[string]string
	tokens   map[string]string
	mu       sync.RWMutex
}

// NewMemoryStorage creates a new instance of MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		sessions: make(map[string]string),
		tokens:   make(map[string]string),
	}
}

// SaveSession saves a user session.
func (m *MemoryStorage) SaveSession(userID string, sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[userID] = sessionID
}

// GetSession retrieves a user session by user ID.
func (m *MemoryStorage) GetSession(userID string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sessionID, exists := m.sessions[userID]
	return sessionID, exists
}

// SaveToken saves a token for a user.
func (m *MemoryStorage) SaveToken(userID string, token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[userID] = token
}

// GetToken retrieves a token by user ID.
func (m *MemoryStorage) GetToken(userID string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	token, exists := m.tokens[userID]
	return token, exists
}

// DeleteSession removes a user session.
func (m *MemoryStorage) DeleteSession(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, userID)
}

// DeleteToken removes a token for a user.
func (m *MemoryStorage) DeleteToken(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tokens, userID)
}