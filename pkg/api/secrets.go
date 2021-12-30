package api

import (
	"io/ioutil"

	"gitlab.void-ptr.org/go/schism/pkg/util"
)

var ApiSecret string

func ReadSecret(name string) string {
	bytes, err := ioutil.ReadFile("/run/secrets/" + name)
	if err != nil {
		util.Log.Panic(err)
	}
	return string(bytes)
}
