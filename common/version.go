package common

import (
	"fmt"
	"os"
)

var (
	Version   string = "dev"
	BuildTime string = "<unset>"
	Commit    string = "<unset>"
)

func PrintVersion() {
	fmt.Printf("version: %s\n", Version)
	fmt.Printf("build time: %s\n", BuildTime)
	fmt.Printf("commit: %s\n", Commit)
	os.Exit(0)
}

func IsDev() bool {
	return Version == "dev" || os.Getenv("DEV") == "true"
}
