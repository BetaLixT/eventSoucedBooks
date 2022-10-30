package req

import "time"

type CreateUpdateCourtroom struct {
	Title         *string                `json:"title"`
	Description   *string                `json:"description"`
	DateTimeStart *time.Time             `json:"dateTimeStart"`
	DateTimeEnd   *time.Time             `json:"dateTimeEnd"`
	Categories    []string               `json:"categories"`
	Location      *string                `json:"location"`
	Metadata      map[string]interface{} `json:"metadata"`
}
