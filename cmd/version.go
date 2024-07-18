package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yoshino-s/go-framework/common"
)

var VersionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		common.PrintVersion()
	},
}
