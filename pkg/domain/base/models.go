package base

import "time"

type Event struct {
	Id        int
	SagaId    string
	Stream    string
	StreamId  string
	Event     string
	Version   int
	EventTime time.Time
	Data      map[string]interface{}
}
