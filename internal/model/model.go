package model

import (
	"errors"
	"time"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrIdentifierExists = errors.New("identifier already exists")
)

type Shortening struct {
	Identifier string    `json:"identifier"`
	Original   string    `json:"original"`
	Visits     int64     `json:"visits"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
