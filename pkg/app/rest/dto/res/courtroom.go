package res

import (
	"time"

	"eventSourcedBooks/pkg/domain/base"
	"eventSourcedBooks/pkg/domain/courtroom"
)

type Courtroom struct {
	Id              string                 `json:"id"`
	Title           *string                `json:"title"`
	Description     *string                `json:"description"`
	DateTimeStart   *time.Time             `json:"dateTimeStart"`
	DateTimeEnd     *time.Time             `json:"dateTimeEnd"`
	Categories      []string               `json:"categories"`
	Location        *string                `json:"location"`
	ExternalCmsId   *string                `json:"externalCmsId"`
	Metadata        map[string]interface{} `json:"metadata"`
	JoinInfo        base.JoinInfo          `json:"joinInfo"`
	DateTimeCreated *time.Time             `json:"dateTimeCreated"`
	DateTimeUpdated *time.Time             `json:"dateTimeUpdated"`
}

func MapCourtroomToDto(cr courtroom.Courtroom) Courtroom {
	return Courtroom{
		Id:              cr.Id,
		Title:           &cr.Title,
		Description:     &cr.Description,
		DateTimeStart:   &cr.DateTimeStart,
		DateTimeEnd:     &cr.DateTimeEnd,
		Categories:      cr.Categories,
		Location:        &cr.Location,
		Metadata:        cr.Metadata,
		JoinInfo:        cr.JoinInfo,
		ExternalCmsId:   &cr.ExternalCmsId,
		DateTimeCreated: &cr.DateTimeCreated,
		DateTimeUpdated: &cr.DateTimeUpdated,
	}
}

func MapCourtroomToDtoList(cr []courtroom.Courtroom) (dto []Courtroom) {
	dto = make([]Courtroom, len(cr))
	for idx, c := range cr {
		dto[idx] = MapCourtroomToDto(c)
	}
	return
}
