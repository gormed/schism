package business

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type DeviceSupport struct {
	Enabled bool
}

type Device struct {
	db.SqlIdentifyable
	Name      string    `json:"name"`
	MacAddr   string    `json:"mac_address"`
	CreatedAt time.Time `json:"date_created"`
	UpdatedAt time.Time `json:"date_updated"`
}

func NewDevice(id *string, database *db.Sqlite) *Device {
	return &Device{SqlIdentifyable: db.SqlIdentifyable{Id: id, Database: database}}
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

type DeviceCreate struct {
	Name    string `json:"name"`
	MacAddr string `json:"mac_address"`
}

// Create device
func (d *Device) Create(create *DeviceCreate) (*Device, int, error) {
	if d.Id != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("the device was already created with id '%s'", *d.Id)
	}

	// Create new unique id for device
	u, err := uuid.NewUUID()
	if err != nil {
		util.Log.Panic(err.Error())
	}
	id := u.String()
	d.Id = &id

	stmt, err := d.Database.Prepare("INSERT INTO devices (id, name, mac_address, date_created, date_updated) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		util.Log.Panic(err.Error())
	}

	tNow := time.Now()
	now := tNow.UTC().Format(db.DateLayout)
	_, err = stmt.Exec(d.Id, create.Name, create.MacAddr, now, now)
	if err != nil {
		util.Log.Panic(err.Error())
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
		util.Log.Panic(err.Error())
	}
	if !exists {
		return d, http.StatusNotFound, fmt.Errorf("device with id '%s' does not exist", id)
	}

	stmt, err := d.Database.Prepare("SELECT name, mac_address, date_created, date_updated FROM devices WHERE id = ?")
	if err != nil {
		util.Log.Panic(err.Error())
	}

	device := Device{SqlIdentifyable: db.SqlIdentifyable{Id: &id}}

	var name, mac_addr, date_created, date_updated string
	err = stmt.QueryRow(*device.Id).Scan(&name, &mac_addr, &date_created, &date_updated)
	if err != nil {
		util.Log.Panic(err.Error())
	}

	device.Name = name
	device.MacAddr = mac_addr
	device.CreatedAt, err = time.Parse(db.DateLayout, date_created)
	if err != nil {
		util.Log.Panic(err.Error())
	}
	device.UpdatedAt, err = time.Parse(db.DateLayout, date_updated)
	if err != nil {
		util.Log.Panic(err.Error())
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
	id := *d.Id

	// Check if resource exists
	exists, err := d.Exists()
	if err != nil {
		util.Log.Panic(err.Error())
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
	now := tNow.UTC().Format(db.DateLayout)

	stmt, err := d.Database.Prepare("UPDATE devices SET name = ?, mac_address = ?, date_updated = ? where id = ?")
	if err != nil {
		util.Log.Panic(err.Error())
	}
	result, err := stmt.Exec(d.Name, d.MacAddr, now, id)
	if err != nil {
		util.Log.Panic(err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		util.Log.Panic(err.Error())
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
		util.Log.Panic(err.Error())
	}
	if !exists {
		return d, http.StatusNotFound, fmt.Errorf("device with id '%s' does not exist", id)
	}

	stmt, err := d.Database.Prepare("DELETE FROM devices WHERE id = ?")
	if err != nil {
		util.Log.Panic(err.Error())
	}
	result, err := stmt.Exec(id)
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
	return d, http.StatusOK, nil
}
