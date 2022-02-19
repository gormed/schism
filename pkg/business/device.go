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

type DeviceSupport struct {
	Enabled bool
}

type Device struct {
	*_business.Device
	Database *db.Sqlite `json:"-"`
}

func NewDevice(id *string, database *db.Sqlite) *Device {
	return &Device{Device: _business.NewDevice(id), Database: database}
}

func (d *Device) Exists() (bool, error) {
	if d.Id == nil {
		return false, fmt.Errorf("no device id given to read")
	}
	stmt, err := d.Database.Prepare("SELECT id from devices where id = ?")
	if err != nil {
		return false, err
	}
	id := *d.Id
	row := stmt.QueryRow(id)
	if err := row.Err(); err != nil {
		return false, err
	}
	return true, nil
}

// Create device
func (d *Device) Create(create *_business.DeviceCreate) (*Device, int, error) {
	if d.Id != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("the device was already created with id '%s'", *d.Id)
	}

	// Create new unique id for device
	u, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	id := u.String()
	d.Id = &id

	stmt, err := d.Database.Prepare("INSERT INTO devices (id, name, mac_address, date_created, date_updated) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}

	tNow := time.Now()
	now := tNow.UTC().Format(db.SqliteDateLayout)
	_, err = stmt.Exec(d.Id, create.Name, create.MacAddr, now, now)
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}

	d.Name = create.Name
	d.MacAddr = create.MacAddr
	d.CreatedAt = tNow
	d.UpdatedAt = tNow

	return d, http.StatusCreated, nil
}

// Read device
func (d *Device) Read() (*Device, int, error) {
	if d.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no device id given to read")
	}
	id := *d.Id

	// Check if resource exists
	exists, err := d.Exists()
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}
	if !exists {
		return d, http.StatusNotFound, fmt.Errorf("device with id '%s' does not exist", id)
	}

	stmt, err := d.Database.Prepare("SELECT name, mac_address, date_created, date_updated FROM devices WHERE id = ?")
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}

	device := NewDevice(&id, d.Database)

	var name, mac_addr, date_created, date_updated string
	err = stmt.QueryRow(*device.Id).Scan(&name, &mac_addr, &date_created, &date_updated)
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}

	device.Name = name
	device.MacAddr = mac_addr
	device.CreatedAt, err = time.Parse(db.SqliteDateLayout, date_created)
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("parse error")
	}
	device.UpdatedAt, err = time.Parse(db.SqliteDateLayout, date_updated)
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("parse error")
	}

	return device, http.StatusOK, nil
}

// Update device
func (d *Device) Update(update *_business.DeviceUpdate) (*Device, int, error) {
	if d.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no device id given to update")
	}
	id := *d.Id

	// Check if resource exists
	exists, err := d.Exists()
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}
	if !exists {
		return d, http.StatusNotFound, fmt.Errorf("device with id '%s' does not exist", id)
	}

	// Update resource properties
	if update.Name != nil {
		d.Name = *update.Name
	}
	if update.MacAddr != nil {
		d.MacAddr = *update.MacAddr
	}
	tNow := time.Now()
	d.UpdatedAt = tNow
	now := tNow.UTC().Format(db.SqliteDateLayout)

	stmt, err := d.Database.Prepare("UPDATE devices SET name = ?, mac_address = ?, date_updated = ? where id = ?")
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}
	result, err := stmt.Exec(d.Name, d.MacAddr, now, id)
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows != 1 {
		util.Log.Panicf("update affected %d rows, only one expected", rows)
	}

	return d, http.StatusOK, nil
}

// Delete device
func (d *Device) Delete() (*Device, int, error) {
	if d.Id == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("no device id given to update")
	}
	id := *d.Id

	exists, err := d.Exists()
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}
	if !exists {
		return d, http.StatusNotFound, fmt.Errorf("device with id '%s' does not exist", id)
	}

	stmt, err := d.Database.Prepare("DELETE FROM devices WHERE id = ?")
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}
	result, err := stmt.Exec(id)
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}
	if rows != 1 {
		util.Log.Panicf("delete affected %d rows, only one expected", rows)
	}
	return d, http.StatusOK, nil
}
