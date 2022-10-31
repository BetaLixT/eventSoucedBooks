package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

type JsonObj map[string]interface{}

var _ driver.Value = (*JsonObj)(nil)

func (a JsonObj) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonObj) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type JoinInfo base.JoinInfo

var _ driver.Value = (*JoinInfo)(nil)

func (a JoinInfo) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JoinInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

// Mapping Functions
func MapDaoToCourtroom(dao *Courtroom) *courtroom.Courtroom {
	externalCmsId := ""
	if dao.ExternalCmsId != nil {
		externalCmsId = *dao.ExternalCmsId
	}

	return &courtroom.Courtroom{
		Id:              dao.Id,
		Title:           dao.Title,
		Description:     dao.Description,
		DateTimeStart:   dao.DateTimeStart,
		DateTimeEnd:     dao.DateTimeEnd,
		Categories:      dao.Categories,
		Location:        dao.Location,
		Metadata:        dao.Metadata,
		ExternalCmsId:   externalCmsId,
		JoinInfo:        base.JoinInfo(dao.JoinInfo),
		DateTimeCreated: dao.DateTimeCreated,
		DateTimeUpdated: dao.DateTimeUpdated,
	}
}

func MapDaoToCourtroomList(daos []Courtroom) (crs []courtroom.Courtroom) {
	crs = make([]courtroom.Courtroom, len(daos))
	for idx, dao := range daos {
		crs[idx] = *MapDaoToCourtroom(&dao)
	}
	return
}

// Generic stuff
type ExistsEntity struct {
	Exists bool `db:"exists"`
}

// Migration script
func GetMigrationScripts() []MigrationScript {
	migrationScripts := []MigrationScript{
		{
			key: "initial books",
			up: `
				CREATE TABLE events (
					id bigserial PRIMARY KEY,
					saga_id text,
					stream text,
					stream_id text,
					event text,
					version bigint,
					event_time timestamp with time zone,
					data bytea,
					CONSTRAINT source_unique UNIQUE (id, version)
				);

				CREATE TRIGGER set_events_event_time
				BEFORE INSERT ON events
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_event_time();`,
			down: `
				DROP TRIGGER set_events_event_time on events;
				DROP TABLE events;`,
		},
		{
			key: "externalCmsId",
			up: `
				ALTER TABLE courtrooms
				ADD COLUMN externalCmsId TEXT UNIQUE`,
			down: `
				ALTER TABLE courtrooms
				DROP COLUMN externalCmsId`,
		},
	}
	return migrationScripts
}
