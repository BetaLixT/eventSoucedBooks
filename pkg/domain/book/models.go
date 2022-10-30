package book

import (
	"eventSourcedBooks/pkg/domain/base"
	"time"
)

type BookState struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Author          string    `json:"author"`
	Genres          []string  `json:"genres"`
	Completed       bool      `json:"completed"`
	CurrentPosition float32   `json:"currentPosition"`
	Version         uint64    `json:"version"`
	DateTimeCreated time.Time `json:"dateTimeCreated"`
	DateTimeUpdated time.Time `json:"dateTimeUpdated"`
}

type BookData struct {
	Title           *string  `json:"title"`
	Description     *string  `json:"description"`
	Author          *string  `json:"author"`
	Genres          []string `json:"genres"`
	Completed       *bool    `json:"completed"`
	CurrentPosition *float32 `json:"currentPosition"`
}

type BookEvent struct {
	base.Event
	Data BookData `json:"data"`
}
