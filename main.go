package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"gitlab.void-ptr.org/go/reflection/pkg/server"
	"gitlab.void-ptr.org/go/schism/pkg/api/router"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

func main() {
	// logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	// logger.ChangePackageLogLevel("lateralus", logger.InfoLevel)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		util.Log.Infof("system call:%+v", oscall)
		cancel()
	}()

	// Setup sqlite connection
	sqlite := db.NewSqlite()
	err := sqlite.Create()
	if err != nil {
		util.Log.Panic(err)
		return
	}

	var influxHost = os.Getenv("INFLUXDB_HOST")
	var influxPort = os.Getenv("INFLUXDB_PORT")
	var influxOrg = os.Getenv("DOCKER_INFLUXDB_INIT_ORG")
	var influxBucket = os.Getenv("DOCKER_INFLUXDB_INIT_BUCKET")

	influxdb := db.NewInflux(fmt.Sprintf("http://%s:%s", influxHost, influxPort), influxOrg, influxBucket)
	err = influxdb.Create()
	if err != nil {
		util.Log.Panic(err)
		return
	}

	s := server.NewSaveServer(util.Log)
	if err := s.Serve(ctx, router.SchismRouter(sqlite, influxdb), func() {
		// Close db connection
		sqlite.Close()
		influxdb.Close()
	}); err != nil {
		util.Log.Errorf("failed to serve:+%v\n", err)
	}
}
