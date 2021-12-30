package business

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gitlab.void-ptr.org/go/schism/pkg/db"
)

type Accesstoken struct {
	db.Identifyable
	Token    *string `json:"token"`
	DeviceId string  `json:"device_id"`
}

func (a *Accesstoken) Exists() (bool, error) {
	if a.Id == nil {
		return false, fmt.Errorf("no accesstoken id given to read")
	}

	stmt, err := a.Database.Prepare("SELECT id from accesstokens where id = ?")
	if err != nil {
		return false, err
	}

	row := stmt.QueryRow(*a.Id)
	if err := row.Err(); err != nil {
		return false, err
	}

	return true, nil
}

type AccesstokenCreate struct {
	DeviceId string `json:"device_id"`
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *Accesstoken) Create(create *AccesstokenCreate) (*Accesstoken, int, error) {
	if a.Id != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("the accesstoken was already created with id '%s'", *a.Id)
	}

	u, err := uuid.NewUUID()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	*a.Id = u.String()

	token, err := randomHex(64)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	stmt, err := a.Database.Prepare("INSERT INTO devices (id, device_id, token) VALUES (?, ?, ?)")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	_, err = stmt.Exec(a.Id, create.DeviceId, token)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	a.DeviceId = create.DeviceId
	a.Token = &token

	return a, http.StatusCreated, nil
}

func (a *Accesstoken) Authenticate(token string) (*Accesstoken, int, error) {
	if a.Token == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no accesstoken token given to read")
	}

	stmt, err := a.Database.Prepare("SELECT id, device_id FROM accesstokens WHERE token = ?")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	accesstoken := Accesstoken{Token: &token}
	err = stmt.QueryRow(*accesstoken.Token).Scan(accesstoken.Id, &accesstoken.DeviceId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &accesstoken, http.StatusOK, nil
}

func (a *Accesstoken) Read() (*Accesstoken, int, error) {
	if a.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no accesstoken id given to read")
	}

	stmt, err := a.Database.Prepare("SELECT token, device_id FROM accesstokens WHERE id = ?")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	accesstoken := Accesstoken{Identifyable: db.Identifyable{
		Id: a.Id,
	}}
	err = stmt.QueryRow(*accesstoken.Id).Scan(accesstoken.Token, &accesstoken.DeviceId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &accesstoken, http.StatusOK, nil
}
