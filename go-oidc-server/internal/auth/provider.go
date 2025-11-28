package auth

import (
    "errors"
    "sync"
)

type User struct {
    ID       string
    Username string
    Password string
}

type Provider struct {
    mu    sync.RWMutex
    users map[string]User
}

func NewProvider() *Provider {
    return &Provider{
        users: make(map[string]User),
    }
}

func (p *Provider) Register(username, password string) (User, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if _, exists := p.users[username]; exists {
        return User{}, errors.New("user already exists")
    }

    user := User{
        ID:       generateID(),
        Username: username,
        Password: password, // In a real application, hash the password
    }
    p.users[username] = user
    return user, nil
}

func (p *Provider) Authenticate(username, password string) (User, error) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    user, exists := p.users[username]
    if !exists || user.Password != password {
        return User{}, errors.New("invalid credentials")
    }
    return user, nil
}

func generateID() string {
    // Implement a function to generate a unique ID for users
    return "some-unique-id" // Placeholder implementation
}