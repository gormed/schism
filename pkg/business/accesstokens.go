package business

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	_business "gitlab.void-ptr.org/go/reflection/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

const Table = "accesstokens"

type Accesstoken struct {
	*_business.Accesstoken
	Database  *db.Sqlite `json:"-"`
	Token     *string    `json:"token"`
	DeviceId  string     `json:"device_id"`
	CreatedAt time.Time  `json:"date_created"`
	UpdatedAt time.Time  `json:"date_updated"`
}

func NewAccesstoken(id *string, database *db.Sqlite) *Accesstoken {
	return &Accesstoken{Accesstoken: _business.NewAccesstoken(id), Database: database}
}

// Exists an accesstoken
func (a *Accesstoken) Exists() (bool, error) {
	if a.Token == nil {
		return false, fmt.Errorf("no accesstoken token given to read")
	}

	stmt, err := a.Database.Prepare(fmt.Sprintf("SELECT id from %s where token = ?", Table))
	if err != nil {
		panic(err)
	}

	row := stmt.QueryRow(*a.Token)
	if err := row.Err(); err != nil {
		panic(err)
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
		panic(err)
	}

	id := u.String()
	a.Id = &id

	// Generate new random accesstoken
	token, err := util.RandomHex(64)
	if err != nil {
		panic(err)
	}

	tNow := time.Now()
	now := tNow.UTC().Format(db.SqliteDateLayout)

	stmt, err := a.Database.Prepare(fmt.Sprintf("INSERT INTO %s (id, device_id, token, date_created, date_updated) VALUES (?, ?, ?, ? ,?)", Table))
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(a.Id, create.DeviceId, token, now, now)
	if err != nil {
		panic(err)
	}

	a.DeviceId = create.DeviceId
	a.Token = &token
	a.CreatedAt = tNow
	a.UpdatedAt = tNow

	return a, http.StatusCreated, nil
}

// Authenticate accesstoken
func (a *Accesstoken) Authenticate(token string) (*Accesstoken, int, error) {
	a = NewAccesstoken(nil, a.Database)
	a.Token = &token

	exists, err := a.Exists()
	if err != nil {
		panic(err)
	}
	if !exists {
		return a, http.StatusNotFound, fmt.Errorf("accesstoken with id '%s' does not exist", *a.Id)
	}

	stmt, err := a.Database.Prepare(fmt.Sprintf("SELECT id, device_id FROM %s WHERE token = ?", Table))
	if err != nil {
		panic(err)
	}
	var id, deviceId string
	err = stmt.QueryRow(*a.Token).Scan(&id, &deviceId)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	a = NewAccesstoken(a.Id, a.Database)
	var token, deviceId string
	err = stmt.QueryRow(*a.Id).Scan(&token, &deviceId)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	if !exists {
		return a, http.StatusNotFound, fmt.Errorf("accesstoken with id '%s' does not exist", *a.Id)
	}

	stmt, err := a.Database.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = ?", Table))
	if err != nil {
		panic(err)
	}
	result, err := stmt.Exec(*a.Id)
	if err != nil {
		panic(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows != 1 {
		util.Log.Panicf("delete affected %d rows, only one expected", rows)
	}
	return a, http.StatusOK, nil
}
