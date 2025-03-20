package model

import "time"

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    time.Time `json:"last_login"`
}

// PlayerResources represents a player's in-game resources
type PlayerResources struct {
	UserID           string    `json:"user_id"`
	Gold             int       `json:"gold"`
	PremiumCurrency  int       `json:"premium_currency"`
	LastIdleClaim    time.Time `json:"last_idle_claim"`
}

// NewUser creates a new user instance
func NewUser(id, username, email, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		LastLogin:    now,
	}
}

// NewPlayerResources creates a new player resources instance
func NewPlayerResources(userID string, gold, premiumCurrency int) *PlayerResources {
	return &PlayerResources{
		UserID:          userID,
		Gold:            gold,
		PremiumCurrency: premiumCurrency,
		LastIdleClaim:   time.Now(),
	}
} 