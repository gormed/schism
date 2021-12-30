package meta

import (
	"os"
)

// Information definition of the Lateralus API
type Information struct {
	Commit  string `json:"commit"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

var gitCommitHash = os.Getenv("GIT_COMMIT")
var gitTag = os.Getenv("GIT_TAG")

// MetaInfo of the Lateralus API
var MetaInfo = Information{
	Commit:  gitCommitHash,
	Version: gitTag,
	Name:    "Schism API",
}
