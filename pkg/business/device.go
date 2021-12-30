package business

import (
	"fmt"

	"gitlab.void-ptr.org/go/schism/pkg/db"
)

type DeviceSupport struct {
	Enabled bool
}

type Device struct {
	Id      *string `json:"id"`
	Name    string  `json:"name"`
	MacAddr string  `json:"mac_address"`
}

func (d *Device) Create() error {
	if d.Id != nil {
		return fmt.Errorf("the device was already created with id '%s'", *d.Id)
	}
	_, err := db.DB.Insert("devices", map[string]interface{}{
		"name":        d.Name,
		"mac_address": d.MacAddr,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) Read() (*Device, error) {
	if d.Id == nil {
		return nil, fmt.Errorf("no device id given to fetch")
	}
	stmt, err := db.DB.Prepare("SELECT name, mac_address from devices where id = ?")
	if err != nil {
		return nil, err
	}
	err = stmt.QueryRow(*d.Id).Scan(&d.Name, &d.MacAddr)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Device) Update() (*Device, error) {
	return d, nil
}

func (d *Device) Delete() (*Device, error) {
	return d, nil
}
