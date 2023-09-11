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

func (d *Data) newSensorValuePoint(
	deviceId, source, name string,
	sensorName string,
	sensorType *sensors.SensorType,
	sensorValue sensors.SensorValue,
	t time.Time,
) *write.Point {
	// Influxdb tags
	tags := map[string]string{
		"deviceId":   deviceId,
		"source":     source,
		"type":       _business.SensorValueType,
		"name":       name,
		"unit":       sensorValue.Unit,
		"unitName":   sensorValue.UnitName,
		"sensorName": sensorName,
	}
	if sensorType != nil {
		tags["sensorType"] = fmt.Sprintf("%d", sensorType)
	}
	// Influxdb fields
	fields := map[string]interface{}{
		"value": sensorValue.Value,
	}
	util.Log.Debugf("NewPoint@%s %s/%s - %s - %s", t.Local().Format("23:05:00.000"), deviceId+"/"+source, tags, fields)
	return influxdb2.NewPoint(deviceId+"/"+source, tags, fields, t)
}

func (d *Data) parseSensorPayload(n *_business.Data, points []*write.Point) ([]*write.Point, int, error) {
	switch *n.SensorType {
	case sensors.SensorType_BMP:
		var payload sensors.BMPSensorData
		err := json.Unmarshal([]byte(n.Payload), &payload)
		if err != nil {
			util.Log.Error(err)
			return nil, http.StatusInternalServerError, fmt.Errorf("unmarshal error")
		}
		for name, t := range map[string]sensors.SensorValue{
			"temprature": *payload.Temprature,
			"humidity":   *payload.Humidity,
			"pressure":   *payload.Pressure,
		} {
			points = append(points, d.newSensorValuePoint(n.DeviceId, n.Source, name, payload.Sensor.Name, n.SensorType, t, n.CreatedAt))
		}
	case sensors.SensorType_SI1145:
		var payload sensors.SI1145SensorData
		err := json.Unmarshal([]byte(n.Payload), &payload)
		if err != nil {
			util.Log.Error(err)
			return nil, http.StatusInternalServerError, fmt.Errorf("unmarshal error")
		}
		for name, t := range map[string]sensors.SensorValue{
			"infraRed":     *payload.InfraRed,
			"ultraViolett": *payload.UltraViolett,
			"visible":      *payload.Visible,
		} {
			points = append(points, d.newSensorValuePoint(n.DeviceId, n.Source, name, payload.Sensor.Name, n.SensorType, t, n.CreatedAt))
		}
	case sensors.SensorType_NU40C16:
		var payload sensors.NU40C16SensorData
		err := json.Unmarshal([]byte(n.Payload), &payload)
		if err != nil {
			util.Log.Error(err)
			return nil, http.StatusInternalServerError, fmt.Errorf("unmarshal error")
		}
		for name, t := range map[string]sensors.SensorValue{
			"distance": *payload.Distance,
		} {
			points = append(points, d.newSensorValuePoint(n.DeviceId, n.Source, name, payload.Sensor.Name, n.SensorType, t, n.CreatedAt))
		}
	case sensors.SensorType_SoilMoisture:
		var payload sensors.SoilMoistureSensorData
		err := json.Unmarshal([]byte(n.Payload), &payload)
		if err != nil {
			util.Log.Error(err)
			return nil, http.StatusInternalServerError, fmt.Errorf("unmarshal error")
		}
		for name, t := range map[string]sensors.SensorValue{
			"soilMoisture": *payload.SoilMoisture,
		} {
			points = append(points, d.newSensorValuePoint(n.DeviceId, n.Source, name, payload.Sensor.Name, n.SensorType, t, n.CreatedAt))
		}
	case sensors.SensorType_AirQuality:
		var payload sensors.AirQualitySensorData
		err := json.Unmarshal([]byte(n.Payload), &payload)
		if err != nil {
			util.Log.Error(err)
			return nil, http.StatusInternalServerError, fmt.Errorf("unmarshal error")
		}
		for name, t := range map[string]sensors.SensorValue{
			"airQuality": *payload.AirQuality,
		} {
			points = append(points, d.newSensorValuePoint(n.DeviceId, n.Source, name, payload.Sensor.Name, n.SensorType, t, n.CreatedAt))
		}
	case sensors.SensorType_Loudness:
		var payload sensors.LoudnessSensorData
		err := json.Unmarshal([]byte(n.Payload), &payload)
		if err != nil {
			util.Log.Error(err)
			return nil, http.StatusInternalServerError, fmt.Errorf("unmarshal error")
		}
		for name, t := range map[string]sensors.SensorValue{
			"loudness": *payload.Loudness,
		} {
			points = append(points, d.newSensorValuePoint(n.DeviceId, n.Source, name, payload.Sensor.Name, n.SensorType, t, n.CreatedAt))
		}
	}
	return points, http.StatusCreated, nil
}

func (d *Data) parseGenericPayload(n *_business.Data, points []*write.Point) ([]*write.Point, int, error) {
	var payload map[string]sensors.SensorValue
	err := json.Unmarshal([]byte(n.Payload), &payload)
	if err != nil {
		util.Log.Error(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("unmarshal error")
	}
	for name, val := range payload {
		if name == "sensor" {
			continue
		}
		points = append(points, d.newSensorValuePoint(n.DeviceId, n.Source, name, "generic", nil, val, n.CreatedAt))
	}
	return points, http.StatusCreated, nil
}

func (d *Data) Create(createData _business.DataCreate) (*_business.DataCreateResponse, int, error) {
	var status = http.StatusInternalServerError
	var points []*write.Point
	var err error
	response := &_business.DataCreateResponse{}

	for _, create := range createData.Data {
		n := _business.NewData()
		n.DeviceId = create.DeviceId
		n.Source = create.Source
		n.DataType = create.DataType
		n.SensorType = create.SensorType
		n.Payload = create.Payload
		n.CreatedAt = create.Meta.MeasuredAt
		n.UpdatedAt = time.Now()
		response.Data = append(response.Data, n)

		switch create.DataType {
		case _business.SensorValue:
			if n.SensorType != nil {
				points, status, err = d.parseSensorPayload(n, points)
				if err != nil {
					util.Log.Error(err)
					return nil, status, err
				}
			} else {
				points, status, err = d.parseGenericPayload(n, points)
				if err != nil {
					util.Log.Error(err)
					return nil, status, err
				}
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

	return response, status, nil
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
