package common

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version string = "dev"
)

func IsDev() bool {
	return Version == "dev"
}

func PrintVersion() {
	fmt.Printf("version: %s\n", Version)
	os.Exit(0)
}

var VersionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		PrintVersion()
	},
}
