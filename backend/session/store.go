package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"nimbus/client"

	"sync"
	"time"
)

type Session struct {
	Client    *client.LibrusClient
	ExpiresAt time.Time
}

type Store struct {
	mu       sync.Mutex
	sessions map[string]Session
}

func NewStore() *Store {
	return &Store{sessions: make(map[string]Session)}
}

func (s *Store) Create(c *client.LibrusClient) (string, error) {
	if c == nil {
		return "", errors.New("nil client")
	}

	read := make([]byte, 64)
	_, err := rand.Read(read)
	if err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(read)
	s.mu.Lock()
	s.sessions[token] = Session{
		Client:    c,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.mu.Unlock()
	return token, nil
}

func (s *Store) Get(token string) (*client.LibrusClient, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[token]
	if !ok {
		return nil, false
	}
	if time.Now().Before(session.ExpiresAt) {
		return session.Client, true
	}
	delete(s.sessions, token)
	return nil, false
}
