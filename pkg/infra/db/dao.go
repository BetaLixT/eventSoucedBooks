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
			key: "initial-books",
			up: `
				CREATE TABLE events (
					id bigserial PRIMARY KEY,
					saga_id text,
					stream text,
					stream_id text,
					version bigint,
					event text,
					event_time timestamp with time zone,
					data bytea,
					CONSTRAINT source_unique UNIQUE (stream, stream_id, version)
				);

				CREATE INDEX idx_events_stream_events ON events(stream, stream_id);

				CREATE TRIGGER set_events_event_time
				BEFORE INSERT ON events
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_event_time();

				CREATE TABLE uniques (
					stream text,
					stream_id text,
					property text,
					value text,
					CONSTRAINT source_unique UNIQUE (stream, property, value),
				)

				CREATE INDEX idx_uniques_stream_constraints
				ON uniques(stream, stream_id);
				`,
			down: `
			  DROP INDEX idx_uniques_stream_constraints;
			  DROP TABLE uniques;
				DROP TRIGGER set_events_event_time on events;
				DROP INDEX idx_events_stream_events;
				DROP TABLE events;
				`,
		},
	}
	return migrationScripts
}
