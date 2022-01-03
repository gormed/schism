package business

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	_business "gitlab.void-ptr.org/go/reflection/pkg/business"
	"gitlab.void-ptr.org/go/reflection/pkg/sensors"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type DataSupport struct {
	Enabled bool
}

type Data struct {
	*_business.Data
	Database *db.Influx `json:"-"`
}

func NewData(database *db.Influx) *Data {
	return &Data{Data: _business.NewData(), Database: database}
}

func (d *Data) Create(create *_business.DataCreate) (*Data, int, error) {
	now := time.Now()
	d.DeviceId = create.DeviceId
	d.Source = create.Source
	d.DataType = create.DataType
	d.Payload = create.Payload
	d.CreatedAt = now
	d.UpdatedAt = now

	switch create.DataType {
	case _business.SensorValue:
		var payload map[string]sensors.SensorValue
		err := json.Unmarshal([]byte(d.Payload), &payload)
		if err != nil {
			panic(err)
		}
		for name, val := range payload {
			tags := map[string]string{
				"deviceId": d.DeviceId,
				"source":   d.Source,
				"type":     _business.SensorValueType,
				"name":     name,
				"unit":     val.Unit,
				"unitName": val.UnitName,
			}
			fields := map[string]interface{}{"value": val.Value}
			point := influxdb2.NewPoint(d.DeviceId+"/"+d.Source, tags, fields, now)

			err = d.Database.Write.WritePoint(context.TODO(), point)
			if err != nil {
				panic(err)
			}
		}
		return d, http.StatusCreated, nil
	default:
		return d, http.StatusBadRequest, fmt.Errorf("unsupported data_type provided")
	}

}

type DataRead struct {
	Start string `json:"start"`
	Stop  string `json:"stop"`
}

type ReadResponse struct {
	Tables  []*TableMetadata
	Records []*Record
}

// FluxTableMetadata holds flux query result table information represented by collection of columns.
// Each new table is introduced by annotations
type TableMetadata struct {
	Position int       `json:"position"`
	Columns  []*Column `json:"columns"`
}

// FluxColumn holds flux query table column properties
type Column struct {
	Index        int    `json:"index"`
	Name         string `json:"name"`
	DataType     string `json:"data_type"`
	Group        bool   `json:"group"`
	DefaultValue string `json:"default_value"`
}

// FluxRecord represents row in the flux query result table
type Record struct {
	Table  int                    `json:"table"`
	Values map[string]interface{} `json:"values"`
}

type Values struct {
	Value       interface{} `json:"_value"`
	Field       string      `json:"_field"`
	Measurement string      `json:"_measurment"`
}

func (d *Data) Read(read *DataRead) (*ReadResponse, int, error) {
	queryRange := "range(start: -1h)"
	if len(read.Start) > 0 {
		queryRange = fmt.Sprintf("range(start: %s)", read.Start)
	}
	if len(read.Start) > 0 && len(read.Stop) > 0 {
		queryRange = fmt.Sprintf("range(start: %s, stop: %s)", read.Start, read.Stop)
	}
	measurement := d.DeviceId + "/" + d.Source
	filter := fmt.Sprintf(`filter(fn: (r) => r._measurement == "%s")`, measurement)
	// Get parser flux query result
	result, err := d.Database.Query.Query(
		context.Background(),
		fmt.Sprintf(`from(bucket:"%s")|> %s |> %s`,
			d.Database.Bucket,
			queryRange,
			filter),
	)
	if err != nil {
		panic(err)
	}
	if result.Err() != nil {
		panic(err)
	}

	var res = ReadResponse{}
	for result.Next() {
		// Observe when there is new grouping key producing new table
		if result.TableChanged() {
			fluxTableMetadata := result.TableMetadata()
			util.Log.Debugf("table: %s\n", fluxTableMetadata.String())

			tableMetadata := &TableMetadata{
				Position: fluxTableMetadata.Position(),
			}
			// Convert columns
			for _, c := range fluxTableMetadata.Columns() {
				tableMetadata.Columns = append(tableMetadata.Columns, &Column{
					Index:        c.Index(),
					Name:         c.Name(),
					DataType:     c.DataType(),
					Group:        c.IsGroup(),
					DefaultValue: c.DefaultValue(),
				})
			}
			res.Tables = append(res.Tables, tableMetadata)
		}
		// Convert records
		fluxRecord := result.Record()
		record := &Record{
			Table:  fluxRecord.Table(),
			Values: fluxRecord.Values(),
		}
		util.Log.Debugf("row: %s\n", fluxRecord.String())
		res.Records = append(res.Records, record)
	}
	return &res, http.StatusOK, nil
}
