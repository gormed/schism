package business

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

const Table = "accesstokens"

type Accesstoken struct {
	db.SqlIdentifyable
	Token     *string   `json:"token"`
	DeviceId  string    `json:"device_id"`
	CreatedAt time.Time `json:"date_created"`
	UpdatedAt time.Time `json:"date_updated"`
}

func NewAccesstoken(id *string, database *db.Sqlite) *Accesstoken {
	return &Accesstoken{SqlIdentifyable: db.SqlIdentifyable{Id: id, Database: database}}
}

// Exists an accesstoken
func (a *Accesstoken) Exists() (bool, error) {
	if a.Token == nil {
		return false, fmt.Errorf("no accesstoken token given to read")
	}

	stmt, err := a.Database.Prepare(fmt.Sprintf("SELECT id from %s where token = ?", Table))
	if err != nil {
		util.Log.Panic(err.Error())
	}

	row := stmt.QueryRow(*a.Token)
	if err := row.Err(); err != nil {
		util.Log.Panic(err.Error())
	}

	return true, nil
}

type AccesstokenCreate struct {
	DeviceId string `json:"device_id"`
}

// Create accesstoken
func (a *Accesstoken) Create(create *AccesstokenCreate) (*Accesstoken, int, error) {
	if a.Id != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("the accesstoken was already created with id '%s'", *a.Id)
	}

	u, err := uuid.NewUUID()
	if err != nil {
		util.Log.Panic(err.Error())
	}

	id := u.String()
	a.Id = &id

	// Generate new random accesstoken
	token, err := util.RandomHex(64)
	if err != nil {
		util.Log.Panic(err.Error())
	}

	tNow := time.Now()
	now := tNow.UTC().Format(db.DateLayout)

	stmt, err := a.Database.Prepare(fmt.Sprintf("INSERT INTO %s (id, device_id, token, date_created, date_updated) VALUES (?, ?, ?, ? ,?)", Table))
	if err != nil {
		util.Log.Panic(err.Error())
	}

	_, err = stmt.Exec(a.Id, create.DeviceId, token, now, now)
	if err != nil {
		util.Log.Panic(err.Error())
	}

	a.DeviceId = create.DeviceId
	a.Token = &token
	a.CreatedAt = tNow
	a.UpdatedAt = tNow

	return a, http.StatusCreated, nil
}

// Authenticate accesstoken
func (a *Accesstoken) Authenticate(token string) (*Accesstoken, int, error) {
	stmt, err := a.Database.Prepare(fmt.Sprintf("SELECT id, device_id FROM %s WHERE token = ?", Table))
	if err != nil {
		util.Log.Panic(err.Error())
	}

	a = &Accesstoken{Token: &token}
	var id, deviceId string
	err = stmt.QueryRow(*a.Token).Scan(&id, &deviceId)
	if err != nil {
		util.Log.Panic(err.Error())
	}

	a.Id = &id
	a.DeviceId = deviceId

	return a, http.StatusOK, nil
}

// Read accesstoken
func (a *Accesstoken) Read() (*Accesstoken, int, error) {
	if a.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no accesstoken id given to read")
	}

	stmt, err := a.Database.Prepare(fmt.Sprintf("SELECT token, device_id FROM %s WHERE id = ?", Table))
	if err != nil {
		util.Log.Panic(err.Error())
	}

	a = &Accesstoken{SqlIdentifyable: db.SqlIdentifyable{
		Id: a.Id,
	}}
	var token, deviceId string
	err = stmt.QueryRow(*a.Id).Scan(&token, &deviceId)
	if err != nil {
		util.Log.Panic(err.Error())
	}

	a.Token = &token
	a.DeviceId = deviceId

	return a, http.StatusOK, nil
}

// Delete accesstoken
func (a *Accesstoken) Delete() (*Accesstoken, int, error) {
	if a.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no accesstoken id given to update")
	}

	exists, err := a.Exists()
	if err != nil {
		util.Log.Panic(err.Error())
	}
	if !exists {
		return a, http.StatusNotFound, fmt.Errorf("accesstoken with id '%s' does not exist", *a.Id)
	}

	stmt, err := a.Database.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = ?", Table))
	if err != nil {
		util.Log.Panic(err.Error())
	}
	result, err := stmt.Exec(*a.Id)
	if err != nil {
		util.Log.Panic(err.Error())
	}
	rows, err := result.RowsAffected()
	if err != nil {
		util.Log.Panic(err.Error())
	}
	if rows != 1 {
		util.Log.Panicf("delete affected %d rows, only one expected", rows)
	}
	return a, http.StatusOK, nil
}
