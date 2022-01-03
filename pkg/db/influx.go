package db

import (
	"context"
	"fmt"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

var InfluxDateLayout = "2006-01-02T15:04:05Z"

type Influx struct {
	ServerURL string
	Org       string
	Bucket    string
	Client    influxdb2.Client
	Write     api.WriteAPIBlocking
	Query     api.QueryAPI
}

func NewInflux(serverURL string, org string, bucket string) *Influx {
	return &Influx{ServerURL: serverURL, Org: org, Bucket: bucket}
}

func (i *Influx) Create() error {
	var adminToken = os.Getenv("DOCKER_INFLUXDB_INIT_ADMIN_TOKEN")

	var client influxdb2.Client
	opts := influxdb2.DefaultOptions()
	opts.SetLogLevel(2)
	client = influxdb2.NewClientWithOptions(i.ServerURL, adminToken, opts)
	ready, err := client.Ready(context.TODO())
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("influxdb client is not ready")
	}

	i.Client = client
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(i.Org, i.Bucket)
	i.Write = writeAPI
	// Create normal query API
	i.Query = i.Client.QueryAPI(i.Org)
	return nil
}

func (i *Influx) Close() {
	i.Client.Close()
}
