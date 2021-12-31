package db

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

type InfluxIdentifyable struct {
	Database *Influx `json:"-"`
	Id       *string `json:"id"`
}

type Influx struct {
	ServerURL string
	Org       string
	Bucket    string
	Client    influxdb2.Client
	Write     api.WriteAPIBlocking
}

func NewInflux(serverURL string, org string, bucket string) *Influx {
	return &Influx{ServerURL: serverURL, Org: org, Bucket: bucket}
}

func (i *Influx) Create(token string) error {
	client := influxdb2.NewClient(i.ServerURL, token)
	i.Client = client
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(i.Org, i.Bucket)
	i.Write = writeAPI
	return nil
}
