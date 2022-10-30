package book

import "eventSourcedBooks/pkg/domain/base"

type BookData struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Author      *string  `json:"author"`
	Genres      []string `json:"genres"`
	Completed   bool     `json:"completed"`
}

type BookEvent struct {
	base.Event
	Data BookData `json:"data"`
}
