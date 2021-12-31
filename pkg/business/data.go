package business

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"gitlab.void-ptr.org/go/reflection/pkg/sensors"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type DataSupport struct {
	Enabled bool
}

type Data struct {
	db.InfluxIdentifyable
	DeviceId  string    `json:"device_id"`
	Source    string    `json:"source"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"date_created"`
	UpdatedAt time.Time `json:"date_updated"`
}

func NewData(id *string, database *db.Influx) *Data {
	return &Data{InfluxIdentifyable: db.InfluxIdentifyable{Id: id, Database: database}}
}

func (d *Data) Create(deviceId string, source string, payload string) (*Data, int, error) {
	var parsed sensors.SensorValue
	err := json.Unmarshal([]byte(payload), &parsed)
	if err != nil {
		util.Log.Panic(err.Error())
	}
	point := influxdb2.NewPoint(deviceId+"/"+source,
		map[string]string{
			"deviceId": deviceId,
			"source":   source,
			"unit":     parsed.Unit,
			"unitName": parsed.UnitName,
		},
		map[string]interface{}{"value": parsed.Value},
		time.Now())
	err = d.Database.Write.WritePoint(context.TODO(), point)
	if err != nil {
		util.Log.Panic(err.Error())
	}
	return d, http.StatusCreated, nil
}

func (d *Data) Read(deviceId string, source string) (*Data, int, error) {
	return d, http.StatusCreated, nil
}
