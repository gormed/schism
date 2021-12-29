package main

import (
	"context"
	"os"
	"os/signal"

	_util "gitlab.void-ptr.org/go/lateralus/pkg/util"
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

	err := db.Create()
	if err != nil {
		util.Log.Panic(err)
	}

	if err := _util.Serve(ctx, router.SchismRouter(), func() {}); err != nil {
		util.Log.Errorf("failed to serve:+%v\n", err)
	}
}
