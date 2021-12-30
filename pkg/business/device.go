package business

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.void-ptr.org/go/schism/pkg/db"
)

type DeviceSupport struct {
	Enabled bool
}

type Device struct {
	db.Identifyable
	Name    string `json:"name"`
	MacAddr string `json:"mac_address"`
}

func (d *Device) Exists() (bool, error) {
	if d.Id == nil {
		return false, fmt.Errorf("no device id given to read")
	}
	stmt, err := d.Database.Prepare("SELECT id from devices where id = ?")
	if err != nil {
		return false, err
	}
	row := stmt.QueryRow(*d.Id)
	if err := row.Err(); err != nil {
		return false, err
	}
	return true, nil
}

type DeviceCreate struct {
	Name    string `json:"name"`
	MacAddr string `json:"mac_address"`
}

// Create device
func (d *Device) Create(create *DeviceCreate) (*Device, int, error) {
	if d.Id != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("the device was already created with id '%s'", *d.Id)
	}
	u, err := uuid.NewUUID()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	id := u.String()
	d.Id = &id

	stmt, err := d.Database.Prepare("INSERT INTO devices (id, name, mac_address, date_created, date_updated) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	tNow := time.Now()
	now := tNow.UTC().Format(db.DateLayout)
	_, err = stmt.Exec(d.Id, create.Name, create.MacAddr, now, now)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	d.Name = create.Name
	d.MacAddr = create.MacAddr

	return d, http.StatusCreated, nil
}

// Read device
func (d *Device) Read() (*Device, int, error) {
	if d.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no device id given to read")
	}

	stmt, err := d.Database.Prepare("SELECT name, mac_address FROM devices WHERE id = ?")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	device := Device{Identifyable: db.Identifyable{
		Id: d.Id,
	}}
	err = stmt.QueryRow(*device.Id).Scan(&device.Name, &device.MacAddr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &device, http.StatusOK, nil
}

type DeviceUpdate struct {
	Name    *string `json:"name"`
	MacAddr *string `json:"mac_address"`
}

// Update device
func (d *Device) Update(update *DeviceUpdate) (*Device, int, error) {
	if d.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no device id given to update")
	}

	// Check if resource exists
	exists, err := d.Exists()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if !exists {
		return d, http.StatusNotFound, fmt.Errorf("device with id '%s' does not exist", *d.Id)
	}

	// Update resource properties
	if update.Name != nil {
		d.Name = *update.Name
	}
	if update.MacAddr != nil {
		d.MacAddr = *update.MacAddr
	}

	tNow := time.Now()
	now := tNow.UTC().Format(db.DateLayout)
	stmt, err := d.Database.Prepare("UPDATE devices SET name = ?, max_address = ?, date_updated = ? where id = ?")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	result, err := stmt.Exec(d.Name, d.MacAddr, now, *d.Id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rows != 1 {
		return nil, http.StatusInternalServerError, fmt.Errorf("update affected %d rows, only one expected", rows)
	}

	return d, http.StatusOK, nil
}

// Delete device
func (d *Device) Delete() (*Device, int, error) {
	if d.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no device id given to update")
	}

	exists, err := d.Exists()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if !exists {
		return d, http.StatusNotFound, fmt.Errorf("device with id '%s' does not exist", *d.Id)
	}

	stmt, err := d.Database.Prepare("DELETE FROM devices WHERE id = ?")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	result, err := stmt.Exec(*d.Id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if rows != 1 {
		return nil, http.StatusInternalServerError, fmt.Errorf("delete affected %d rows, only one expected", rows)
	}
	return d, http.StatusOK, nil
}
