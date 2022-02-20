package business

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/influxdata/influxdb-client-go/domain"
	_business "gitlab.void-ptr.org/go/reflection/pkg/business"
	"gitlab.void-ptr.org/go/reflection/pkg/sensors"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type DataSupport struct {
	Enabled bool
}

type Data struct {
	Database *db.Influx `json:"-"`
}

func NewData(database *db.Influx) *Data {
	return &Data{Database: database}
}

func (d *Data) Create(createData _business.DataCreate) (*_business.DataCreateResponse, int, error) {
	now := time.Now()
	var points []*write.Point
	response := &_business.DataCreateResponse{}

	for _, create := range createData.Data {
		n := _business.NewData()
		n.DeviceId = create.DeviceId
		n.Source = create.Source
		n.DataType = create.DataType
		n.Payload = create.Payload
		n.CreatedAt = now
		n.UpdatedAt = now
		response.Data = append(response.Data, n)

		switch create.DataType {
		case _business.SensorValue:
			var payload map[string]sensors.SensorValue
			err := json.Unmarshal([]byte(n.Payload), &payload)
			if err != nil {
				util.Log.Error(err)
				return nil, http.StatusInternalServerError, fmt.Errorf("unmarshal error")
			}
			for name, val := range payload {
				tags := map[string]string{
					"deviceId": n.DeviceId,
					"source":   n.Source,
					"type":     _business.SensorValueType,
					"name":     name,
					"unit":     val.Unit,
					"unitName": val.UnitName,
				}
				fields := map[string]interface{}{"value": val.Value}
				point := influxdb2.NewPoint(n.DeviceId+"/"+n.Source, tags, fields, now)
				points = append(points, point)
			}
		default:
			return nil, http.StatusBadRequest, fmt.Errorf("unsupported data_type provided")
		}
	}

	health, err := d.Database.Client.Health(context.TODO())
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if health.Status != domain.HealthCheckStatusPass {
		return nil, http.StatusInternalServerError, fmt.Errorf("influxdb not healthy")
	}
	err = d.Database.Write.WritePoint(context.TODO(), points...)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return response, http.StatusCreated, nil
}

func (d *Data) Read(read *_business.DataRead) (*_business.DataReadResponse, int, error) {
	queryRange := "range(start: -1h)"
	if len(read.Start) > 0 {
		queryRange = fmt.Sprintf("range(start: %s)", read.Start)
	}
	if len(read.Start) > 0 && len(read.Stop) > 0 {
		queryRange = fmt.Sprintf("range(start: %s, stop: %s)", read.Start, read.Stop)
	}
	measurement := read.DeviceId + "/" + read.Source
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
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}
	if result.Err() != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("database error")
	}

	var res = _business.DataReadResponse{}
	for result.Next() {
		// Observe when there is new grouping key producing new table
		if result.TableChanged() {
			fluxTableMetadata := result.TableMetadata()
			util.Log.Debugf("table: %s\n", fluxTableMetadata.String())

			tableMetadata := &_business.TableMetadata{
				Position: fluxTableMetadata.Position(),
			}
			// Convert columns
			for _, c := range fluxTableMetadata.Columns() {
				tableMetadata.Columns = append(tableMetadata.Columns, &_business.Column{
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
		record := &_business.Record{
			Table:  fluxRecord.Table(),
			Values: fluxRecord.Values(),
		}
		util.Log.Debugf("row: %s\n", fluxRecord.String())
		res.Records = append(res.Records, record)
	}
	return &res, http.StatusOK, nil
}
