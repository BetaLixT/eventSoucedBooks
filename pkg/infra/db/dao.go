package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"eventSourcedBooks/pkg/domain/base"
	"eventSourcedBooks/pkg/domain/courtroom"

	"github.com/lib/pq"
)

// Models as stored in the database with required tags
type Courtroom struct {
	Id              string
	Title           string
	Description     string
	DateTimeStart   time.Time
	Location        string
	DateTimeEnd     time.Time
	Categories      pq.StringArray
	JoinInfo        JoinInfo
	ExternalCmsId   *string
	Metadata        JsonObj
	Version         int
	Deleted         bool
	DateTimeCreated time.Time
	DateTimeUpdated time.Time
}

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
			key: "initial",
			up: `
				CREATE TABLE courtrooms (
					id uuid PRIMARY KEY,
					title text,
					description text,
					datetimestart timestamp with time zone,
					datetimeend timestamp with time zone,
					location text,
					categories text[],
					joininfo jsonb,
					metadata jsonb,
					version integer DEFAULT 0,
					deleted boolean DEFAULT false,
					datetimecreated timestamp with time zone,
					datetimeupdated timestamp with time zone
				);

				CREATE TRIGGER set_courtrooms_datetimecreated
				BEFORE INSERT ON courtrooms
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_datetimecreated();
				
				CREATE TRIGGER set_courtrooms_datetimeupdated_in
				BEFORE INSERT ON courtrooms
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_datetimeupdated();

				CREATE TRIGGER set_courtrooms_datetimeupdated_up
				BEFORE UPDATE ON courtrooms
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_datetimeupdated();

				CREATE TRIGGER courtrooms_version_update
				BEFORE UPDATE ON courtrooms
				FOR EACH ROW
				EXECUTE PROCEDURE version_update();`,
			down: `
				DROP TRIGGER courtrooms_version_update on courtrooms;
				DROP TRIGGER set_courtrooms_datetimeupdated_up on courtrooms;
				DROP TRIGGER set_courtrooms_datetimeupdated_in on courtrooms;
				DROP TRIGGER set_courtrooms_datetimecreated on courtrooms;
				DROP TABLE courtrooms;`,
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
