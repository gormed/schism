package business

import (
	"context"
	"encoding/json"
	"fmt"
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
	DataType  DataType  `json:"data_type"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"date_created"`
	UpdatedAt time.Time `json:"date_updated"`
}

func NewData(database *db.Influx) *Data {
	return &Data{InfluxIdentifyable: db.InfluxIdentifyable{Database: database}}
}

type DataType int

const (
	SensorValue DataType = iota
)

type DataCreate struct {
	DeviceId string   `json:"device_id"`
	Source   string   `json:"source"`
	DataType DataType `json:"data_type"`
	Payload  string   `json:"payload"`
}

func (d *Data) Create(create *DataCreate) (*Data, int, error) {
	now := time.Now()
	d.DeviceId = create.DeviceId
	d.Source = create.Source
	d.DataType = create.DataType
	d.Payload = create.Payload
	d.CreatedAt = now
	d.UpdatedAt = now

	switch create.DataType {
	case SensorValue:

		var parsed sensors.SensorValue
		err := json.Unmarshal([]byte(d.Payload), &parsed)
		if err != nil {
			util.Log.Panic(err.Error())
		}
		tags := map[string]string{
			"deviceId": d.DeviceId,
			"source":   d.Source,
			"unit":     parsed.Unit,
			"unitName": parsed.UnitName,
		}
		fields := map[string]interface{}{"value": parsed.Value}
		point := influxdb2.NewPoint(d.DeviceId+"/"+d.Source, tags, fields, now)

		err = d.Database.Write.WritePoint(context.TODO(), point)
		if err != nil {
			util.Log.Panic(err.Error())
		}
		return d, http.StatusCreated, nil
	default:
		return d, http.StatusBadRequest, fmt.Errorf("unsupported data_type provided")
	}

}

type DataRead struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (d *Data) Read(read *DataRead) (*Data, int, error) {
	_measurement := d.DeviceId + "/" + d.Source
	// Get parser flux query result
	result, err := d.Database.Query.Query(
		context.Background(),
		fmt.Sprintf(`from(bucket:"%s")|> range(start: %s) |> filter(fn: (r) => r._measurement == "%s")`,
			d.Database.Bucket,
			read.Start,
			_measurement),
	)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// read result
			fmt.Printf("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}
	return d, http.StatusOK, nil
}
