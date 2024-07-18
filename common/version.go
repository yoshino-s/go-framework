package common

import (
	"fmt"
	"os"
)

var (
	Version string = "dev"
)

func PrintVersion() {
	fmt.Printf("version: %s\n", Version)
	os.Exit(0)
}

func IsDev() bool {
	return Version == "dev"
}
