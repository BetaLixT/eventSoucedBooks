package book

import (
)

type BookData struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Author          string    `json:"author"`
	Genres          []string  `json:"genres"`
}
