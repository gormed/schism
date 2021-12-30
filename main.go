package main

import (
	"context"
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

	// Setup db connection
	db := db.NewSqlite()
	err := db.Create()
	if err != nil {
		util.Log.Panic(err)
		return
	}

	s := server.NewSaveServer(util.Log)
	if err := s.Serve(ctx, router.SchismRouter(db), func() {
		// Close db connection
		db.Close()
	}); err != nil {
		util.Log.Errorf("failed to serve:+%v\n", err)
	}
}
