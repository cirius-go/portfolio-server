package model

import "time"

// Model.
type Model struct {
	ID        string    `gorm:"primaryKey;default:uuid_generate_v7();type:uuid" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ENUM(Debug)
//
//go:generate go-enum --marshal
type ContextKey string
